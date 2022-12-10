<h1 align="center">
  <br>
  <img src="https://socialify.git.ci/proxfly/hegex/image?font=Inter&language=1&name=1&owner=1&pattern=Circuit%20Board&theme=Light" alt="Markdownify" width="500">
  <br>
  hegex
  <br>
</h1>

<h4 align="center">A regex for HTTP, built on top
of <a href="https://github.com/golang/go/blob/master/src/regexp/regexp.go" target="_blank">golang regexp</a>.</h4>

<p align="center">
  <a href="#key-features">Key features</a> •
  <a href="#how-to-use">How to use</a> •
  <a href="#benchmark">Benchmark</a> •
  <a href="#roadmap">Roadmap</a> •
  <a href="#contact">Contact</a> •
  <a href="#contributing">Contributing</a> •
  <a href="#credits">Credits</a> •
  <a href="#license">License</a>
</p>

**hegex** helps you write simple patterns to match hostnames or paths.

Take a look at what **hegex** can do when matching a hostname:

| Pattern for hostname             | Match                                                                 |
|----------------------------------|-----------------------------------------------------------------------|
| `*.example.com`                  | `string-01.example.com`, `string-01.string-02.example.com`, and so on |
| `{site}.example.com`             | `string-01.example.com`                                               |
| `{site[a\|bb\|ccc]}.example.com` | only `a.example.com`, `bb.example.com`, `ccc.example.com`             |

The same rules apply to paths:

| Pattern for path                     | Match                                                                       |
|--------------------------------------|-----------------------------------------------------------------------------|
| `/path/*/api.json`                   | `/path/string-01/api.json`, `/path/string-01/string-02/api.json`, and so on |
| `/path/{middle}/api.json`            | `/path/string-01/api.json`                                                  |
| `/path/{middle[a\|bb\|ccc]/api.json` | only `/path/a/api.json`, `/path/bb/api.json`, `/path/ccc/api.json`          |

**Q: Who should use hegex?**

**A:** hegex can be used to write rules for an API gateway. Regular expressions (regex) may require a lot of mental
effort to
read when it comes to `[a-zA-Z0-9-]+\.example\.com`. By using hegex, the pattern be rewritten as `{site}.example.com`,
which is more intuitive.

Though hegex brings simplicity, it is less flexible and has lower performance than regex.

## NOTE: hegex is not ready for production

## Key features

* Match hostname and path
* `{}` supports candidates
* Treat `*`, `**`, and `***` as different placeholders in the process of rewriting

## How to use

With [Go module](https://github.com/golang/go/wiki/Modules) support, simply add the following import to your code:

```go
import "github.com/proxfly/hegex"
```

Then, `go [build|run|test]` will automatically fetch the necessary dependencies.

Otherwise, run the following Go command to install the `hegex` package:

```sh
$ go get -u github.com/proxfly/hegex
```

We currently support the most recent major Go versions from 1.19 onward.

### Examples

Match a hostname

```go
package main

import (
	"fmt"
	"github.com/proxfly/hegex"
)

func main() {
	he, err := hegex.Compile("*.example.com")
	if err != nil {
		panic(err)
	}
	ok := he.MatchString("string-01.example.com")
	fmt.Println(ok) // true
}
```

Match a path

```go
package main

import (
	"fmt"
	"github.com/proxfly/hegex"
)

func main() {
	he, err := hegex.Compile("/path/*/api.{postfix[json|yml]}")
	if err != nil {
		panic(err)
	}
	ok := he.MatchString("/path/doc/api.json")
	fmt.Println(ok) // true
}
```

Match and rewrite a hostname

```go
package main

import (
	"fmt"
	"github.com/proxfly/hegex"
)

func main() {
	he, err := hegex.Compile("*.example.com")
	if err != nil {
		panic(err)
	}
	rewritten, ok := he.MatchAndRewrite("api.example.com", "example.com/*")
	fmt.Println(ok)        // true
	fmt.Println(rewritten) // example.com/api
}
```

Match and rewrite a path

```go
package main

import (
	"fmt"
	"github.com/proxfly/hegex"
)

func main() {
	he, err := hegex.Compile("/*/**")
	if err != nil {
		panic(err)
	}
	// Both * and ** match any string
	// In the process of rewriting, hegex knows * points to "path" and ** points to "data"
	// You can use longer repeated * in your code. For example: /***/image/*/**.jpg
	rewritten, ok := he.MatchAndRewrite("/path/data", "/*/to/**")
	fmt.Println(ok)        // true
	fmt.Println(rewritten) // /path/to/data
}
```

For more examples, see [hegexp_test.go](hegexp_test.go)

## Benchmark

[See benchmark details](BENCHMARKS.MD)

| Benchmark     |         1 CPU |        8 CPUs |
|---------------|--------------:|--------------:|
| Simple/regex  | 3704464 ns/op | 3859987 ns/op |
| Simple/hegex  | 3892382 ns/op | 3778665 ns/op |
| Complex/regex |  616395 ns/op |  603943 ns/op |
| Complex/hegex | 1371399 ns/op | 1384042 ns/op |

The benchmark results clearly show that as the patterns get longer and more complex, hegex's performance drops off
noticeably compared to regex.

## Roadmap

- [x] Multi-asterisk placeholder
- [x] Curly brace placeholder with candidates
- [x] Replace placeholder with string
- [x] Benchmark and compare hegex with golang regexp
- [ ] Support non-ASCII characters

## Contact

* [Email](mailto:github-contributing.ra2t5@simplelogin.com?subject=hegex%20Feedback)
* Slack: [#hegex](https://proxfly.slack.com/archives/C05G65BLYUA) channel

## Contributing

hegex is MIT licensed and accepts contributions via GitHub pull requests.

## Credits

This software uses the following open-source packages:

- [testify](https://github.com/stretchr/testify)

## License

MIT