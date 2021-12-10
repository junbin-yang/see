package see

// github.com/julienschmidt/httprouter
/*
goos: linux
goarch: amd64
pkg: github.com/julienschmidt/go-http-routing-benchmark
cpu: Intel(R) Xeon(R) CPU E5-26xx v4
BenchmarkSee_Param          12098535         103.0 ns/op         0 B/op        0 allocs/op
BenchmarkSee_Param5          5566252         225.2 ns/op         0 B/op        0 allocs/op
BenchmarkSee_Param20         1622751         748.9 ns/op         0 B/op        0 allocs/op
BenchmarkSee_ParamWrite      9367747         133.0 ns/op         0 B/op        0 allocs/op
BenchmarkSee_GithubStatic    4433935         279.5 ns/op         0 B/op        0 allocs/op
BenchmarkSee_GithubParam     3829358         305.4 ns/op         0 B/op        0 allocs/op
BenchmarkSee_GithubAll         17100       70916 ns/op         2 B/op        0 allocs/op
BenchmarkSee_GPlusStatic    13158685          92.75 ns/op        0 B/op        0 allocs/op
BenchmarkSee_GPlusParam      9106490         126.7 ns/op         0 B/op        0 allocs/op
BenchmarkSee_GPlus2Params    5727114         210.7 ns/op         0 B/op        0 allocs/op
BenchmarkSee_GPlusAll         550324        2165 ns/op         0 B/op        0 allocs/op
BenchmarkSee_ParseStatic    10578208         115.6 ns/op         0 B/op        0 allocs/op
BenchmarkSee_ParseParam     11151218         108.0 ns/op         0 B/op        0 allocs/op
BenchmarkSee_Parse2Params    7978048         151.3 ns/op         0 B/op        0 allocs/op
BenchmarkSee_ParseAll         293788        4112 ns/op         0 B/op        0 allocs/op
BenchmarkSee_StaticAll         19059       63878 ns/op         0 B/op        0 allocs/op
PASS
ok    github.com/julienschmidt/go-http-routing-benchmark  24.768s
*/
