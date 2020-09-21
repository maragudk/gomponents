# gomponents

[![GoDoc](https://godoc.org/github.com/maragudk/gomponents?status.svg)](https://godoc.org/github.com/maragudk/gomponents)
[![codecov](https://codecov.io/gh/maragudk/gomponents/branch/master/graph/badge.svg)](https://codecov.io/gh/maragudk/gomponents)

gomponents are components of DOM nodes for Go, that can render to an HTML Document.
gomponents aims to make it easy to build HTML pages of reusable components,
without the use of a template language. Think server-side-rendered React,
but without the virtual DOM and diffing.

The implementation is still incomplete, but usable. The API may change until version 1 is reached.

## Usage

Get the library using `go get`:

```shell script
go get -u github.com/maragudk/gomponents
```

Then do something like this:

```go
package foo

import (
    "net/http"

    g "github.com/maragudk/gomponents"
    "github.com/maragudk/gomponents/el"
)

func main() {
    _ = http.ListenAndServe("localhost:8080", handler())
}

func handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = g.Write(w, el.Document(el.HTML(el.Head(el.Title("Hi!")), el.Body(el.H1("Hi!")))))
	}
}
```

For more complete examples, see [the examples directory](examples/).
