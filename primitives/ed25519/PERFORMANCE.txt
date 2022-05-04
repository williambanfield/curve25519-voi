go test -bench=. -test.v -run BenchmarkVerifyBatchOnly
goos: linux
goarch: amd64
pkg: github.com/oasisprotocol/curve25519-voi/primitives/ed25519
cpu: Intel(R) Core(TM) i7-4790K CPU @ 4.00GHz
BenchmarkVerifyBatchOnly
BenchmarkVerifyBatchOnly/1
BenchmarkVerifyBatchOnly/1-8       15415             68863 ns/op            5816 B/op         14 allocs/op
BenchmarkVerifyBatchOnly/2
BenchmarkVerifyBatchOnly/2-8       12628             92898 ns/op           10248 B/op         17 allocs/op
BenchmarkVerifyBatchOnly/4
BenchmarkVerifyBatchOnly/4-8        9768            141113 ns/op           18720 B/op         22 allocs/op
BenchmarkVerifyBatchOnly/8
BenchmarkVerifyBatchOnly/8-8        4975            226853 ns/op           36448 B/op         31 allocs/op
BenchmarkVerifyBatchOnly/16
BenchmarkVerifyBatchOnly/16-8               2773            443891 ns/op           70256 B/op         48 allocs/op
BenchmarkVerifyBatchOnly/32
BenchmarkVerifyBatchOnly/32-8               1386            760545 ns/op          145824 B/op         81 allocs/op
BenchmarkVerifyBatchOnly/64
BenchmarkVerifyBatchOnly/64-8                818           1513406 ns/op          283264 B/op        146 allocs/op
BenchmarkVerifyBatchOnly/128
BenchmarkVerifyBatchOnly/128-8               436           2816722 ns/op          334081 B/op        205 allocs/op
BenchmarkVerifyBatchOnly/256
BenchmarkVerifyBatchOnly/256-8               240           4961097 ns/op          520960 B/op        206 allocs/op
BenchmarkVerifyBatchOnly/384
BenchmarkVerifyBatchOnly/384-8               169           7067180 ns/op          819841 B/op        207 allocs/op
BenchmarkVerifyBatchOnly/512
BenchmarkVerifyBatchOnly/512-8               132           9323969 ns/op          898817 B/op        207 allocs/op
BenchmarkVerifyBatchOnly/768
BenchmarkVerifyBatchOnly/768-8                97          13713648 ns/op         1398530 B/op        208 allocs/op
BenchmarkVerifyBatchOnly/1024
BenchmarkVerifyBatchOnly/1024-8               73          16969990 ns/op         2071810 B/op        209 allocs/op
BenchmarkGenerateKey
BenchmarkGenerateKey/voi
BenchmarkGenerateKey/voi-8                 92394             12761 ns/op             128 B/op          3 allocs/op
BenchmarkGenerateKey/stdlib
BenchmarkGenerateKey/stdlib-8              55874             20989 ns/op             128 B/op          3 allocs/op
BenchmarkNewKeyFromSeed
BenchmarkNewKeyFromSeed/voi
BenchmarkNewKeyFromSeed/voi-8              95295             12668 ns/op               0 B/op          0 allocs/op
BenchmarkNewKeyFromSeed/stdlib
BenchmarkNewKeyFromSeed/stdlib-8           56394             20930 ns/op               0 B/op          0 allocs/op
BenchmarkSigning
BenchmarkSigning/voi
BenchmarkSigning/voi-8                     86106             13914 ns/op              64 B/op          1 allocs/op
BenchmarkSigning/stdlib
BenchmarkSigning/stdlib-8                  47311             25364 ns/op               0 B/op          0 allocs/op
BenchmarkVerification
BenchmarkVerification/voi
BenchmarkVerification/voi-8                31423             38106 ns/op               0 B/op          0 allocs/op
BenchmarkVerification/voi_stdlib
BenchmarkVerification/voi_stdlib-8         30446             39506 ns/op               0 B/op          0 allocs/op
BenchmarkVerification/stdlib
BenchmarkVerification/stdlib-8             19455             60375 ns/op               0 B/op          0 allocs/op
BenchmarkExpanded
BenchmarkExpanded/NewExpandedPublicKey
BenchmarkExpanded/NewExpandedPublicKey-8                  194943              6727 ns/op            1504 B/op          2 allocs/op
BenchmarkExpanded/Verification
BenchmarkExpanded/Verification/voi
BenchmarkExpanded/Verification/voi-8                       36622             32231 ns/op               0 B/op          0 allocs/op
BenchmarkExpanded/Verification/voi_stdlib
BenchmarkExpanded/Verification/voi_stdlib-8                34561             34662 ns/op               0 B/op          0 allocs/op
BenchmarkExpanded/VerifyBatchOnly
BenchmarkExpanded/VerifyBatchOnly/1
BenchmarkExpanded/VerifyBatchOnly/1-8                      25112             48267 ns/op            4312 B/op         12 allocs/op
BenchmarkExpanded/VerifyBatchOnly/2
BenchmarkExpanded/VerifyBatchOnly/2-8                      18632             65312 ns/op            7240 B/op         13 allocs/op
BenchmarkExpanded/VerifyBatchOnly/4
BenchmarkExpanded/VerifyBatchOnly/4-8                      12596            115578 ns/op           12704 B/op         14 allocs/op
BenchmarkExpanded/VerifyBatchOnly/8
BenchmarkExpanded/VerifyBatchOnly/8-8                       6494            168266 ns/op           24416 B/op         15 allocs/op
BenchmarkExpanded/VerifyBatchOnly/16
BenchmarkExpanded/VerifyBatchOnly/16-8                      4186            293298 ns/op           46192 B/op         16 allocs/op
BenchmarkExpanded/VerifyBatchOnly/32
BenchmarkExpanded/VerifyBatchOnly/32-8                      2241            571748 ns/op           97696 B/op         17 allocs/op
BenchmarkExpanded/VerifyBatchOnly/64
BenchmarkExpanded/VerifyBatchOnly/64-8                       945           1109797 ns/op          187008 B/op         18 allocs/op
BenchmarkExpanded/VerifyBatchOnly/128
BenchmarkExpanded/VerifyBatchOnly/128-8                      597           2036238 ns/op          192704 B/op         17 allocs/op
BenchmarkExpanded/VerifyBatchOnly/256
BenchmarkExpanded/VerifyBatchOnly/256-8                      327           3773224 ns/op          379584 B/op         18 allocs/op
BenchmarkExpanded/VerifyBatchOnly/384
BenchmarkExpanded/VerifyBatchOnly/384-8                      224           5282515 ns/op          678464 B/op         19 allocs/op
BenchmarkExpanded/VerifyBatchOnly/512
BenchmarkExpanded/VerifyBatchOnly/512-8                      171           6836644 ns/op          757440 B/op         19 allocs/op
BenchmarkExpanded/VerifyBatchOnly/768
BenchmarkExpanded/VerifyBatchOnly/768-8                      121           9768945 ns/op         1257153 B/op         20 allocs/op
BenchmarkExpanded/VerifyBatchOnly/1024
BenchmarkExpanded/VerifyBatchOnly/1024-8                      78          13066307 ns/op         1930433 B/op         21 allocs/op
PASS
ok      github.com/oasisprotocol/curve25519-voi/primitives/ed25519      58.467s
