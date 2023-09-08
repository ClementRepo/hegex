CPU Model: Apple M1, 8 Core

Go Version: 1.20.2 darwin/arm64

hegex Version: v0.0.1

```bash
$ go test ./... -bench=. -cpu=1,8
goos: darwin
goarch: arm64
pkg: github.com/proxoar/hegex
BenchmarkSimple/regex                    3704464               331.4 ns/op
BenchmarkSimple/regex-8                  3859987               316.3 ns/op
BenchmarkSimple/hegex                    3892382               319.5 ns/op
BenchmarkSimple/hegex-8                  3778665               309.3 ns/op
BenchmarkComplex/regex                    616395              1952 ns/op
BenchmarkComplex/regex-8                  603943              1947 ns/op
BenchmarkComplex/hegex                   1371399               877.7 ns/op
BenchmarkComplex/hegex-8                 1384042               866.5 ns/op
PASS
ok      github.com/proxoar/hegex        13.675s
```

`*` `BenchmarkComplex` uses a longer and more complex string than `BenchmarkSimple` does.
See [hegex_bench_test.go](hegex_bench_test.go)