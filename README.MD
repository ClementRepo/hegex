<h1 align="center">
  <br>
  <img src="https://socialify.git.ci/ClementRepo/hegex/image?font=Inter&language=1&name=1&owner=1&pattern=Circuit%20Board&theme=Light" alt="Markdownify" width="500">
  <br>
  hegex
  <br>
</h1>

<h4 align="center">A regex for HTTP, built on top
of <a href="https://github.com/golang/go/blob/master/src/regexp/regexp.go" target="_blank">golang regexp</a>.</h4>

<p align="center">
  <a href="#key-features">Key features</a> •
  <a href="#how-to-use">How to use</a> •
  <a href="#roadmap">Roadmap</a> •
  <a href="#contact">Contact</a> •
  <a href="#contributing">Contributing</a> •
  <a href="#credits">Credits</a> •
  <a href="#license">License</a>
</p>

**Hegex** helps you write simple patterns to match HTTP host or HTTP path.

Take a look at what **hegex** can do when matching a host:

| pattern for host                                            | match                                                                                |
|-------------------------------------------------------------|--------------------------------------------------------------------------------------|
| `*.example.com`                                             | `string-01.example.com`,`string-01.string-02.example.com` and so on                  |
| `{site}.example.com`                                        | `string-01.example.com`                                                              |
| `{site[candidate01\|candidate02\|candidate03]}.example.com` | only `candidate01.example.com`, `candidate02.example.com`, `candidate03.example.com` |

The same rules apply to path:

| pattern for path                                                | match                                                                                         |
|-----------------------------------------------------------------|-----------------------------------------------------------------------------------------------|
| `/path/*/api.json`                                              | `/path/string-01/api.json`,`/path/string-01/string-02/api.json` and so on                     |
| `/path/{middle}/api.json`                                       | `/path/string-01/api.json`                                                                    |
| `/path/{middle[candidate01\|candidate02\|candidate03]/api.json` | only `/path/candidate01/api.json`, `/path/candidate02/api.json`, `/path/candidate03/api.json` |

Q: What is the purpose of **hegex**?

A: **Hegex** can be used to write API gateway rules. **Regex** may require a lot of mental effort to read when
it comes to `[a-zA-Z0-9-]+\.example\.com`, which can be rewritten more intuitively as `{site}.example.com` by using *
*hegex**.

Though **hegex** brings simplicity, it is less flexible and has lower performance than **regex**.

## NOTE: Hegex is not ready for production

## Key features

* Match hostname and HTTP path
* `*`, `**` and `***`... are treated as different placeholder
* `{}` placeholder support candidates

## How to use

With [Go module](https://github.com/golang/go/wiki/Modules) support, simply add the following import

```
import "github.com/ClementRepo/hegex"
```

to your code, and then `go [build|run|test]` will automatically fetch the necessary dependencies.

Otherwise, run the following Go command to install the `hegex` package:

```sh
$ go get -u github.com/ClementRepo/hegex
```

### Code Example
* match an HTTP host
```go
package main

import (
	"fmt"
	"github.com/ClementRepo/hegex"
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
* match an HTTP path
```go
package main

import (
	"fmt"
	"github.com/ClementRepo/hegex"
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
* match and rewrite an HTTP host
```go
package main

import (
	"fmt"
	"github.com/ClementRepo/hegex"
)

func main() {
	he, err := hegex.Compile("*.example.com")
	if err != nil {
		panic(err)
	}
	rewritten,ok := he.MatchAndRewrite("api.example.com", "example.com/*")
	fmt.Println(ok) // true
	fmt.Println(rewritten) // example.com/api
}
```
* match and rewrite an HTTP host
```go
package main

import (
	"fmt"
	"github.com/ClementRepo/hegex"
)

func main() {
	he, err := hegex.Compile("/*/**")
	if err != nil {
		panic(err)
	}
	// Both * and ** match any string
	// In the process of rewriting, hegex knows * points to "path" and ** points to "data"
	// You can use longer repeated * in your code. For example: /***/image/*/**.jpg
	rewritten,ok := he.MatchAndRewrite("/path/data", "/*/to/**")
	fmt.Println(ok) // true
	fmt.Println(rewritten) // /path/to/data
}
```

For more examples, see [hegexp_test.go](hegexp_test.go)

## Roadmap

- [x] Multi-sterix placeholder
- [x] Curly brace placeholder with candidates
- [x] Replace placeholder with string
- [] Benchmark and compare hegex with golang regexp
- [] Support non-ASCII characters

## Contact

* <a href="mailto:github-contributing.ra2t5@simplelogin.com?subject=hegex Feedback"> Email </a>
* Slack: [#hegex](https://clementrepo.slack.com/archives/C05G65BLYUA) channel

## Contributing

hegex is MIT licensed and accepts contributions via GitHub pull requests.

## Credits

This software uses the following open source packages:

- [testify](github.com/stretchr/testify)

## License

MIT