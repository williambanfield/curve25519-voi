// Copyright (c) 2016-2019 Isis Agora Lovecruft, Henry de Valence. All rights reserved.
// Copyright (c) 2021 Oasis Labs Inc.  All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
// 1. Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright
// notice, this list of conditions and the following disclaimer in the
// documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS
// IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED
// TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A
// PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
// TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
// PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
// LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
// NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// +build amd64,!purego,!forcenoasm,!force32bit

package curve

import (
	"golang.org/x/sys/cpu"

	"github.com/oasisprotocol/curve25519-voi/internal/disalloweq"
	"github.com/oasisprotocol/curve25519-voi/internal/field"
)

var supportsVectorizedEdwards bool

// This is the dalek AVX2 backend, ported to a mix of Go and Go's assembly
// dialect.
//
// What originally started out as nicely abstracted Rust using intrinsics
// is now essentially a fucked gigantic assembly file of doom, along with
// some stubs because there is no way to maintain the original abstraction
// (split between point and field arithmetic), without spending an excessive
// amount of time shuffling data between memory and YMM registers.
//
// At this point, this is more of an elaborate cry for help than anything
// else.

//go:noescape
func vecConditionalSelect_AVX2(out, a, b *fieldElement2625x4, mask uint32)

//go:noescape
func vecNegate_AVX2(out *fieldElement2625x4)

//go:noescape
func vecReduce_AVX2(out *fieldElement2625x4)

//go:noescape
func vecMul_AVX2(out, a, b *fieldElement2625x4)

//go:noescape
func vecSquareAndNegateD_AVX2(out *fieldElement2625x4)

//go:noescape
func vecDoubleExtended_Step1_AVX2(out *fieldElement2625x4, vec *extendedPoint)

//go:noescape
func vecDoubleExtended_Step2_AVX2(tmp0, tmp1 *fieldElement2625x4)

//go:noescape
func vecAddSubExtendedCached_Step1_AVX2(out *fieldElement2625x4, vec *extendedPoint)

//go:noescape
func vecAddSubExtendedCached_Step2_AVX2(tmp0, tmp1 *fieldElement2625x4)

//go:noescape
func vecNegateLazyCached_AVX2(out *fieldElement2625x4, vec *cachedPoint)

//go:noescape
func vecCachedFromExtended_Step1_AVX2(out *cachedPoint, vec *extendedPoint)

type extendedPoint struct {
	disalloweq.DisallowEqual //nolint:unused
	inner                    fieldElement2625x4
}

func (p *EdwardsPoint) setExtended(ep *extendedPoint) *EdwardsPoint {
	ep.inner.Split(&p.inner.X, &p.inner.Y, &p.inner.Z, &p.inner.T)
	return p
}

func (p *extendedPoint) SetEdwards(ep *EdwardsPoint) *extendedPoint {
	p.inner = newFieldElement2625x4(&ep.inner.X, &ep.inner.Y, &ep.inner.Z, &ep.inner.T)
	return p
}

func (p *extendedPoint) ConditionalSelect(a, b *extendedPoint, choice int) {
	p.inner.ConditionalSelect(&a.inner, &b.inner, choice)
}

func (p *extendedPoint) ConditionalAssign(other *extendedPoint, choice int) {
	p.inner.ConditionalAssign(&other.inner, choice)
}

// Note: dalek has the identity point as the default ctor for
// ExtendedPoint/CachedPoint.

func (p *extendedPoint) Identity() *extendedPoint {
	*p = constEXTENDEDPOINT_IDENTITY
	return p
}

func (p *extendedPoint) Double(t *extendedPoint) *extendedPoint {
	var tmp1 fieldElement2625x4
	vecDoubleExtended_Step1_AVX2(&tmp1, t)
	vecSquareAndNegateD_AVX2(&tmp1)

	var tmp0 fieldElement2625x4
	vecDoubleExtended_Step2_AVX2(&tmp0, &tmp1)
	p.inner.Mul(&tmp0, &tmp1)

	return p
}

func (p *extendedPoint) MulByPow2(t *extendedPoint, k uint) *extendedPoint {
	// Note: Assumes `k > 0`, but the panic is elided.
	p.Double(t)
	for i := uint(0); i < k-1; i++ {
		p.Double(p)
	}
	return p
}

func (p *extendedPoint) AddExtendedCached(a *extendedPoint, b *cachedPoint) *extendedPoint {
	var tmp0 fieldElement2625x4
	vecAddSubExtendedCached_Step1_AVX2(&tmp0, a)
	tmp0.Mul(&tmp0, &b.inner)

	var tmp1 fieldElement2625x4
	vecAddSubExtendedCached_Step2_AVX2(&tmp0, &tmp1)
	p.inner.Mul(&tmp0, &tmp1)

	return p
}

func (p *extendedPoint) SubExtendedCached(a *extendedPoint, b *cachedPoint) *extendedPoint {
	var tmp0, other fieldElement2625x4
	vecAddSubExtendedCached_Step1_AVX2(&tmp0, a)
	vecNegateLazyCached_AVX2(&other, b)
	tmp0.Mul(&tmp0, &other)

	var tmp1 fieldElement2625x4
	vecAddSubExtendedCached_Step2_AVX2(&tmp0, &tmp1)
	p.inner.Mul(&tmp0, &tmp1)

	return p
}

type cachedPoint struct {
	disalloweq.DisallowEqual //nolint:unused
	inner                    fieldElement2625x4
}

func (p *cachedPoint) SetExtended(ep *extendedPoint) *cachedPoint {
	vecCachedFromExtended_Step1_AVX2(p, ep)

	neg_x := *p
	neg_x.inner.Neg()

	// x = x.blend(-x, Lanes::D)
	//
	// Having to take the overhead of a function call, move x and neg_x
	// into the YMM registers, execute 5 VPBLENDDs, and then write back
	// x, is probably more expensive than just doing this serially.
	for i := 0; i < 5; i++ {
		p.inner.inner[i][5] = neg_x.inner.inner[i][5] // d_2i
		p.inner.inner[i][7] = neg_x.inner.inner[i][7] // d_2i_1
	}

	return p
}

func (p *cachedPoint) ConditionalSelect(a, b *cachedPoint, choice int) {
	p.inner.ConditionalSelect(&a.inner, &b.inner, choice)
}

func (p *cachedPoint) ConditionalAssign(other *cachedPoint, choice int) {
	p.inner.ConditionalAssign(&other.inner, choice)
}

func (p *cachedPoint) ConditionalNegate(choice int) {
	var pNeg cachedPoint
	vecNegateLazyCached_AVX2(&pNeg.inner, p)
	p.ConditionalAssign(&pNeg, choice)
}

type fieldElement2625x4 struct {
	disalloweq.DisallowEqual //nolint:unused
	inner                    [5][8]uint32
}

// ConditionalSelect sets the field elements to a iff choice == 0 and
// b iff choice == 1.
func (vec *fieldElement2625x4) ConditionalSelect(a, b *fieldElement2625x4, choice int) {
	mask := uint32(-choice)
	vecConditionalSelect_AVX2(vec, a, b, mask)
}

// ConditionalAssign conditionally assigns the field elements according to choice.
func (vec *fieldElement2625x4) ConditionalAssign(other *fieldElement2625x4, choice int) {
	vec.ConditionalSelect(vec, other, choice)
}

// Split splits the vector into four (serial) field elements.
func (vec *fieldElement2625x4) Split(fe0, fe1, fe2, fe3 *field.FieldElement) {
	fe0i, fe1i, fe2i, fe3i := fe0.UnsafeInner(), fe1.UnsafeInner(), fe2.UnsafeInner(), fe3.UnsafeInner()
	for i := 0; i < 5; i++ {
		fe0i[i] = uint64(vec.inner[i][0]) + (uint64(vec.inner[i][2]) << 26) // a_2i + (a_2i_1 << 26)
		fe1i[i] = uint64(vec.inner[i][1]) + (uint64(vec.inner[i][3]) << 26) // b_2i + (b_2i_1 << 26)
		fe2i[i] = uint64(vec.inner[i][4]) + (uint64(vec.inner[i][6]) << 26) // c_2i + (c_2i_1 << 26)
		fe3i[i] = uint64(vec.inner[i][5]) + (uint64(vec.inner[i][7]) << 26) // d_2i + (d_2i_1 << 26)
	}
}

// Neg computes `(-A, -B, -C, -D)`.
func (vec *fieldElement2625x4) Neg() {
	vecNegate_AVX2(vec)
}

// Reduce reduces the vector of field elements.
func (vec *fieldElement2625x4) Reduce() {
	vecReduce_AVX2(vec)
}

// Mul computes `a * b`.
func (vec *fieldElement2625x4) Mul(a, b *fieldElement2625x4) {
	vecMul_AVX2(vec, a, b)
}

// SquareAndNegateD squares the field elements and negates the result's D value.
func (vec *fieldElement2625x4) SquareAndNegateD() {
	vecSquareAndNegateD_AVX2(vec)
}

// newFieldElement2625 constructs a field element vector from its raw components.
func newFieldElement2625x4(fe0, fe1, fe2, fe3 *field.FieldElement) fieldElement2625x4 {
	const low_26_bit_mask uint64 = (1 << 26) - 1

	fe0i, fe1i, fe2i, fe3i := fe0.UnsafeInner(), fe1.UnsafeInner(), fe2.UnsafeInner(), fe3.UnsafeInner()

	var fe fieldElement2625x4
	for i := 0; i < 5; i++ {
		fe.inner[i][0] = uint32(fe0i[i] & low_26_bit_mask) // a_2i
		fe.inner[i][1] = uint32(fe1i[i] & low_26_bit_mask) // b_2i
		fe.inner[i][2] = uint32(fe0i[i] >> 26)             // a_2i_1
		fe.inner[i][3] = uint32(fe1i[i] >> 26)             // b_2i_1
		fe.inner[i][4] = uint32(fe2i[i] & low_26_bit_mask) // c_2i
		fe.inner[i][5] = uint32(fe3i[i] & low_26_bit_mask) // d_2i
		fe.inner[i][6] = uint32(fe2i[i] >> 26)             // c_2i_1
		fe.inner[i][7] = uint32(fe3i[i] >> 26)             // d_2i_1
	}

	fe.Reduce()

	return fe
}

func init() {
	supportsVectorizedEdwards = cpu.Initialized && cpu.X86.HasAVX2

	// Enable the vector backend for the hardcoded basepoint table,
	// if the vector backend is enabled for everything else.
	if supportsVectorizedEdwards {
		ED25519_BASEPOINT_TABLE.inner = constVECTOR_ED25519_BASEPOINT_TABLE
	}
}
