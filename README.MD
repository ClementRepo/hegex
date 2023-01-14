<h1 align="center">
  <br>
  <img src="https://socialify.git.ci/proxfly/hegex/image?font=Inter&language=1&name=1&owner=1&pattern=Circuit%20Board&theme=Light" alt="Markdownify" width="500">
  <br>
  hegex
  <br>
</h1>

<h4 align="center">A straightforward regular expression for HTTP, built on top
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

**hegex** helps you write simple expressions to match hostnames or paths.

Check out what **hegex** is capable of:

| Expression                       | Match                                                           |
|----------------------------------|-----------------------------------------------------------------|
| `*.example.com`                  | `string-01.example.com`, `string-01.string-02.example.com`, etc |
| `{site}.example.com`             | `string-01.example.com`                                         |
| `{site[a\|bb\|ccc]}.example.com` | only `a.example.com`, `bb.example.com` and `ccc.example.com`    |

The same rules apply to paths:

| Expression                           | Match                                                                 |
|--------------------------------------|-----------------------------------------------------------------------|
| `/path/*/api.json`                   | `/path/string-01/api.json`, `/path/string-01/string-02/api.json`, etc |
| `/path/{middle}/api.json`            | `/path/string-01/api.json`                                            |
| `/path/{middle[a\|bb\|ccc]/api.json` | only `/path/a/api.json`, `/path/bb/api.json` and `/path/ccc/api.json` |

**Q: Why hegex?**

**A:** hegex streamlines the process of formulating rules for an API gateway, making it more intuitive and
user-friendly. Regular expressions, such as `[a-zA-Z0-9-]+\.example\.com` or `[^\s\./]+\.example\.com`, can be mentally
taxing to interpret. However, with hegex, you can simply write the expression as `{site}.example.com`, which is
significantly more straightforward.

Despite making things simpler, hegex doesn't match the versatility and efficiency provided by regular expressions.

## NOTE: hegex is NOT yet prepared for production use

## Key features

* Write expression to match hostname or path
* Use `{}` in an expression to match any characters except for `.` or `/`
* Use `{string}` in an expression
* Use `{string[opt-1|opt-2]}` in an expression to limit the characters to certain options
* Use `*` in an expression to match absolutely anything, including `.` and `/`
* Use `**`, `***`, `****`, etc in an expression
* Use `{}`, `{string}`, `*`, `**`, `***`, etc in the template, and substitute them with the corresponding characters

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

Match and substitute a hostname

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
	substitute, ok := he.MatchAndSubstitute("api.example.com", "example.com/*")
	fmt.Println(ok)         // true
	fmt.Println(substitute) // example.com/api
}
```

Match and substitute a path

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
	// In the process of substitution, hegex knows that * corresponds to "path" and ** corresponds to "data"
	// You can use longer consecutive * in your code. For example: /***/image/*/**.jpg
	substitute, ok := he.MatchAndSubstitute("/path/data", "/*/to/**")
	fmt.Println(ok)         // true
	fmt.Println(substitute) // /path/to/data
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

The benchmark results clearly show that as the expressions get longer and more complex, hegex's performance drops off
noticeably compared to regex.

## Roadmap

- [x] Multi-asterisk placeholder
- [x] Curly brace placeholder with options
- [x] Replace placeholder with string
- [x] Benchmark and compare hegex with golang regexp
- [ ] Support non-ASCII characters
- [ ] Production ready: hegex needs more feedback to be fully prepared for production use

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