# Tired of complex template languages?

<img src="logo.png" alt="Logo" width="300" align="right">

[![GoDoc](https://pkg.go.dev/badge/maragu.dev/gomponents)](https://pkg.go.dev/maragu.dev/gomponents)
[![CI](https://github.com/maragudk/gomponents/actions/workflows/ci.yml/badge.svg)](https://github.com/maragudk/gomponents/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/maragudk/gomponents/branch/main/graph/badge.svg)](https://codecov.io/gh/maragudk/gomponents)

Try HTML components in pure Go.

_gomponents_ are HTML components written in pure Go.
They render to HTML 5, and make it easy for you to build reusable components.
So you can focus on building your app instead of learning yet another templating language.

```shell
go get maragu.dev/gomponents
```

Made with ✨sparkles✨ by [maragu](https://www.maragu.dev/): independent software consulting for cloud-native Go apps & AI engineering.

[Contact me at markus@maragu.dk](mailto:markus@maragu.dk) for consulting work, or perhaps an invoice to support this project?

## Features

Check out [www.gomponents.com](https://www.gomponents.com) for an introduction or [pkg.go.dev/maragu.dev/gomponents](https://pkg.go.dev/maragu.dev/gomponents) for the official docs.

- Build reusable HTML components
- Write declarative HTML 5 in Go without all the strings, so you get
  - Type safety from the compiler
  - Auto-completion from the IDE
  - Easy debugging with the standard Go debugger
  - Automatic formatting with `gofmt`/`goimports`
- Simple API that's easy to learn and use (you know most already if you know HTML)
- Useful helpers like
  - `Text` and `Textf` that insert HTML-escaped text,
  - `Raw` and `Rawf` for inserting raw strings,
  - `Map` for mapping data to components and `Group` for grouping components,
  - and `If`/`Iff` for conditional rendering.
- No external dependencies
- Mature and stable, no breaking changes

## Usage

```shell
go get maragu.dev/gomponents
```

```go
package main

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Navbar(authenticated bool, currentPath string) Node {
	return Nav(
		NavbarLink("/", "Home", currentPath),
		NavbarLink("/about", "About", currentPath),
		If(authenticated, NavbarLink("/profile", "Profile", currentPath)),
	)
}

func NavbarLink(href, name, currentPath string) Node {
	return A(Href(href), Classes{"is-active": currentPath == href}, Text(name))
}
```

(Some people don't like dot-imports, and luckily it's completely optional.)

For a more complete example, see [the examples directory](internal/examples/).
There's also the [gomponents-starter-kit](https://github.com/maragudk/gomponents-starter-kit) for a full application template.

### Coding agents

There's a [skill](skills/gomponents/SKILL.md) that teaches coding agents how to use gomponents. It's also a good, concise introduction to gomponents for humans. Install it as a plugin:

<details>
<summary>Claude Code</summary>

```shell
/plugin marketplace add maragudk/gomponents
/plugin install gomponents@gomponents
```

</details>

<details>
<summary>Codex</summary>

Run these commands from a local clone:

```shell
codex plugin marketplace add maragudk/gomponents
codex plugin add gomponents@gomponents
```

</details>

## Architecture

gomponents is organized into several packages:

- `gomponents`: Core interfaces and functions like `Node`, `El`, `Attr`, and helpers like `Map`, `Group`, `If`, `Text`, `Raw`.
- `gomponents/html`: HTML elements and attributes.
- `gomponents/components`: Higher-level components and utilities.
- `gomponents/http`: HTTP-related utilities for web servers.
- `gomponents/x/...`: Experimental packages. These do not have the same compatibility guarantees as the core library, and in particular, may get breaking changes.

### Void Elements

Void elements in HTML (like `<br>`, `<img>`, `<input>`) don't have closing tags.
gomponents handles these correctly by checking against an internal list of void elements during rendering.
When you create a void element, any child nodes that are not attributes will be ignored automatically to ensure valid HTML output.

## FAQ

### Is gomponents production-ready?

Yes! gomponents is mature, stable, fully tested with 100% coverage, and is used in production by myself and many others, and has been for years.

### Should I choose `html/template`, Templ, or gomponents?

These are all good choices, and it largely comes down to preference.
I wrote gomponents because I didn't like how I think it's hard to pass data around between templates in `html/template`.
gomponents is pure Go, with no extra build step like Templ, so it works with all tools that already support Go.

That said, both `html/template` and Templ will do the same thing as gomponents in the end. Try them all and choose what you like!

### Is gomponents fast?

Yes. gomponents renders directly to an `io.Writer`, making it efficient for server-side rendering.
The library avoids unnecessary allocations where possible.
There's also an extensive benchmark suite to keep it that way, which you can run with `make benchmark`.

<details>
<summary>Benchmark results on an Apple M4 (2026-09-17)</summary>

```
go test -bench . -benchmem ./...
goos: darwin
goarch: arm64
pkg: maragu.dev/gomponents
cpu: Apple M4
BenchmarkAttr/construct_and_render/discarded/boolean-10         	100000000	        10.74 ns/op	      24 B/op	       1 allocs/op
BenchmarkAttr/render_pre-built/discarded/boolean-10             	320929339	         3.750 ns/op	       0 B/op	       0 allocs/op
BenchmarkAttr/construct_and_render/discarded/short_value,_no_escaping-10         	60663632	        19.50 ns/op	      48 B/op	       1 allocs/op
BenchmarkAttr/render_pre-built/discarded/short_value,_no_escaping-10             	100000000	        11.81 ns/op	       0 B/op	       0 allocs/op
BenchmarkAttr/construct_and_render/discarded/long_value,_no_escaping-10          	13451110	        88.10 ns/op	      48 B/op	       1 allocs/op
BenchmarkAttr/render_pre-built/discarded/long_value,_no_escaping-10              	15595132	        77.74 ns/op	       0 B/op	       0 allocs/op
BenchmarkAttr/construct_and_render/discarded/short_value,_little_escaping-10     	46189747	        25.86 ns/op	      48 B/op	       1 allocs/op
BenchmarkAttr/render_pre-built/discarded/short_value,_little_escaping-10         	64064064	        18.33 ns/op	       0 B/op	       0 allocs/op
BenchmarkAttr/construct_and_render/discarded/long_value,_little_escaping-10      	 8555257	       140.3 ns/op	      48 B/op	       1 allocs/op
BenchmarkAttr/render_pre-built/discarded/long_value,_little_escaping-10          	 8940240	       134.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkAttr/construct_and_render/discarded/short_value,_much_escaping-10       	36048033	        32.60 ns/op	      48 B/op	       1 allocs/op
BenchmarkAttr/render_pre-built/discarded/short_value,_much_escaping-10           	50799472	        23.73 ns/op	       0 B/op	       0 allocs/op
BenchmarkAttr/construct_and_render/discarded/long_value,_much_escaping-10        	 4690027	       254.6 ns/op	      48 B/op	       1 allocs/op
BenchmarkAttr/render_pre-built/discarded/long_value,_much_escaping-10            	 4764429	       252.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkAttr/construct_and_render/buffered/boolean-10                           	94576621	        13.01 ns/op	      24 B/op	       1 allocs/op
BenchmarkAttr/render_pre-built/buffered/boolean-10                               	166734266	         7.216 ns/op	       0 B/op	       0 allocs/op
BenchmarkAttr/construct_and_render/buffered/short_value,_no_escaping-10          	39031530	        28.45 ns/op	      48 B/op	       1 allocs/op
BenchmarkAttr/render_pre-built/buffered/short_value,_no_escaping-10              	56037778	        21.36 ns/op	       0 B/op	       0 allocs/op
BenchmarkAttr/construct_and_render/buffered/long_value,_no_escaping-10           	12038955	        99.30 ns/op	      48 B/op	       1 allocs/op
BenchmarkAttr/render_pre-built/buffered/long_value,_no_escaping-10               	13049937	        91.44 ns/op	       0 B/op	       0 allocs/op
BenchmarkAttr/construct_and_render/buffered/short_value,_little_escaping-10      	31760734	        37.96 ns/op	      48 B/op	       1 allocs/op
BenchmarkAttr/render_pre-built/buffered/short_value,_little_escaping-10          	39149636	        30.90 ns/op	       0 B/op	       0 allocs/op
BenchmarkAttr/construct_and_render/buffered/long_value,_little_escaping-10       	 6229292	       193.5 ns/op	      48 B/op	       1 allocs/op
BenchmarkAttr/render_pre-built/buffered/long_value,_little_escaping-10           	 6390865	       187.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkAttr/construct_and_render/buffered/short_value,_much_escaping-10        	21732404	        52.81 ns/op	      48 B/op	       1 allocs/op
BenchmarkAttr/render_pre-built/buffered/short_value,_much_escaping-10            	25714606	        44.95 ns/op	       0 B/op	       0 allocs/op
BenchmarkAttr/construct_and_render/buffered/long_value,_much_escaping-10         	 2656693	       451.3 ns/op	      48 B/op	       1 allocs/op
BenchmarkAttr/render_pre-built/buffered/long_value,_much_escaping-10             	 2729870	       437.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkAttr/construct_and_render/write-only/boolean-10                         	62418456	        18.83 ns/op	      32 B/op	       2 allocs/op
BenchmarkAttr/render_pre-built/write-only/boolean-10                             	91025748	        12.75 ns/op	       8 B/op	       1 allocs/op
BenchmarkAttr/construct_and_render/write-only/short_value,_no_escaping-10        	31461930	        37.61 ns/op	      72 B/op	       3 allocs/op
BenchmarkAttr/render_pre-built/write-only/short_value,_no_escaping-10            	39774226	        30.04 ns/op	      24 B/op	       2 allocs/op
BenchmarkAttr/construct_and_render/write-only/long_value,_no_escaping-10         	 9534724	       126.2 ns/op	     312 B/op	       3 allocs/op
BenchmarkAttr/render_pre-built/write-only/long_value,_no_escaping-10             	10368916	       116.0 ns/op	     264 B/op	       2 allocs/op
BenchmarkAttr/construct_and_render/write-only/short_value,_little_escaping-10    	17865394	        66.69 ns/op	     104 B/op	       4 allocs/op
BenchmarkAttr/render_pre-built/write-only/short_value,_little_escaping-10        	20061423	        56.85 ns/op	      56 B/op	       3 allocs/op
BenchmarkAttr/construct_and_render/write-only/long_value,_little_escaping-10     	 3384784	       356.7 ns/op	    1016 B/op	       5 allocs/op
BenchmarkAttr/render_pre-built/write-only/long_value,_little_escaping-10         	 3452178	       347.8 ns/op	     968 B/op	       4 allocs/op
BenchmarkAttr/construct_and_render/write-only/short_value,_much_escaping-10      	17231866	        69.34 ns/op	     120 B/op	       4 allocs/op
BenchmarkAttr/render_pre-built/write-only/short_value,_much_escaping-10          	19882389	        60.26 ns/op	      72 B/op	       3 allocs/op
BenchmarkAttr/construct_and_render/write-only/long_value,_much_escaping-10       	 2594019	       461.7 ns/op	    1592 B/op	       5 allocs/op
BenchmarkAttr/render_pre-built/write-only/long_value,_much_escaping-10           	 2712394	       443.6 ns/op	    1544 B/op	       4 allocs/op
BenchmarkAttr/construct_and_render/write-only_buffered/boolean-10                	53303227	        22.46 ns/op	      32 B/op	       2 allocs/op
BenchmarkAttr/render_pre-built/write-only_buffered/boolean-10                    	71694579	        16.33 ns/op	       8 B/op	       1 allocs/op
BenchmarkAttr/construct_and_render/write-only_buffered/short_value,_no_escaping-10         	23188666	        51.15 ns/op	      72 B/op	       3 allocs/op
BenchmarkAttr/render_pre-built/write-only_buffered/short_value,_no_escaping-10             	28709564	        42.28 ns/op	      24 B/op	       2 allocs/op
BenchmarkAttr/construct_and_render/write-only_buffered/long_value,_no_escaping-10          	 8706303	       137.1 ns/op	     312 B/op	       3 allocs/op
BenchmarkAttr/render_pre-built/write-only_buffered/long_value,_no_escaping-10              	 9322614	       127.5 ns/op	     264 B/op	       2 allocs/op
BenchmarkAttr/construct_and_render/write-only_buffered/short_value,_little_escaping-10     	15925530	        75.94 ns/op	     104 B/op	       4 allocs/op
BenchmarkAttr/render_pre-built/write-only_buffered/short_value,_little_escaping-10         	17756732	        66.76 ns/op	      56 B/op	       3 allocs/op
BenchmarkAttr/construct_and_render/write-only_buffered/long_value,_little_escaping-10      	 3242295	       371.4 ns/op	    1016 B/op	       5 allocs/op
BenchmarkAttr/render_pre-built/write-only_buffered/long_value,_little_escaping-10          	 3279267	       363.7 ns/op	     968 B/op	       4 allocs/op
BenchmarkAttr/construct_and_render/write-only_buffered/short_value,_much_escaping-10       	15746222	        77.78 ns/op	     120 B/op	       4 allocs/op
BenchmarkAttr/render_pre-built/write-only_buffered/short_value,_much_escaping-10           	16958872	        69.24 ns/op	      72 B/op	       3 allocs/op
BenchmarkAttr/construct_and_render/write-only_buffered/long_value,_much_escaping-10        	 2475848	       480.3 ns/op	    1592 B/op	       5 allocs/op
BenchmarkAttr/render_pre-built/write-only_buffered/long_value,_much_escaping-10            	 2583634	       464.7 ns/op	    1544 B/op	       4 allocs/op
BenchmarkEl/normal_elements-10                                                             	72882795	        16.74 ns/op	      48 B/op	       1 allocs/op
BenchmarkRaw/raw_element-10                                                                	415332338	         2.878 ns/op	       0 B/op	       0 allocs/op
BenchmarkRawf/formatted_raw_element-10                                                     	31366227	        38.21 ns/op	      40 B/op	       2 allocs/op
BenchmarkText/construct_and_render/discarded/short_text,_no_escaping-10                    	87209301	        13.86 ns/op	      16 B/op	       1 allocs/op
BenchmarkText/render_pre-built/discarded/short_text,_no_escaping-10                        	514631895	         2.316 ns/op	       0 B/op	       0 allocs/op
BenchmarkText/construct_and_render/discarded/long_text,_no_escaping-10                     	14761682	        83.14 ns/op	      16 B/op	       1 allocs/op
BenchmarkText/render_pre-built/discarded/long_text,_no_escaping-10                         	519436578	         2.310 ns/op	       0 B/op	       0 allocs/op
BenchmarkText/construct_and_render/discarded/short_text,_little_escaping-10                	30175236	        39.09 ns/op	      40 B/op	       2 allocs/op
BenchmarkText/render_pre-built/discarded/short_text,_little_escaping-10                    	518872826	         2.307 ns/op	       0 B/op	       0 allocs/op
BenchmarkText/construct_and_render/discarded/long_text,_little_escaping-10                 	 3867835	       308.8 ns/op	     656 B/op	       3 allocs/op
BenchmarkText/render_pre-built/discarded/long_text,_little_escaping-10                     	523249729	         2.302 ns/op	       0 B/op	       0 allocs/op
BenchmarkText/construct_and_render/discarded/short_text,_much_escaping-10                  	28601136	        41.65 ns/op	      48 B/op	       2 allocs/op
BenchmarkText/render_pre-built/discarded/short_text,_much_escaping-10                      	517962228	         2.294 ns/op	       0 B/op	       0 allocs/op
BenchmarkText/construct_and_render/discarded/long_text,_much_escaping-10                   	 3021069	       395.8 ns/op	    1040 B/op	       3 allocs/op
BenchmarkText/render_pre-built/discarded/long_text,_much_escaping-10                       	516403984	         2.311 ns/op	       0 B/op	       0 allocs/op
BenchmarkText/construct_and_render/buffered/short_text,_no_escaping-10                     	75498212	        15.28 ns/op	      16 B/op	       1 allocs/op
BenchmarkText/render_pre-built/buffered/short_text,_no_escaping-10                         	326405378	         3.672 ns/op	       0 B/op	       0 allocs/op
BenchmarkText/construct_and_render/buffered/long_text,_no_escaping-10                      	13552480	        88.19 ns/op	      16 B/op	       1 allocs/op
BenchmarkText/render_pre-built/buffered/long_text,_no_escaping-10                          	200795988	         5.950 ns/op	       0 B/op	       0 allocs/op
BenchmarkText/construct_and_render/buffered/short_text,_little_escaping-10                 	29852200	        40.31 ns/op	      40 B/op	       2 allocs/op
BenchmarkText/render_pre-built/buffered/short_text,_little_escaping-10                     	342237326	         3.514 ns/op	       0 B/op	       0 allocs/op
BenchmarkText/construct_and_render/buffered/long_text,_little_escaping-10                  	 3820880	       313.0 ns/op	     656 B/op	       3 allocs/op
BenchmarkText/render_pre-built/buffered/long_text,_little_escaping-10                      	173351754	         6.911 ns/op	       0 B/op	       0 allocs/op
BenchmarkText/construct_and_render/buffered/short_text,_much_escaping-10                   	27716479	        42.79 ns/op	      48 B/op	       2 allocs/op
BenchmarkText/render_pre-built/buffered/short_text,_much_escaping-10                       	343751961	         3.506 ns/op	       0 B/op	       0 allocs/op
BenchmarkText/construct_and_render/buffered/long_text,_much_escaping-10                    	 2857876	       414.3 ns/op	    1040 B/op	       3 allocs/op
BenchmarkText/render_pre-built/buffered/long_text,_much_escaping-10                        	129388539	         9.279 ns/op	       0 B/op	       0 allocs/op
BenchmarkText/construct_and_render/write-only/short_text,_no_escaping-10                   	54852649	        22.02 ns/op	      32 B/op	       2 allocs/op
BenchmarkText/render_pre-built/write-only/short_text,_no_escaping-10                       	100000000	        10.26 ns/op	      16 B/op	       1 allocs/op
BenchmarkText/construct_and_render/write-only/long_text,_no_escaping-10                    	11116381	       106.6 ns/op	     272 B/op	       2 allocs/op
BenchmarkText/render_pre-built/write-only/long_text,_no_escaping-10                        	45116952	        26.79 ns/op	     256 B/op	       1 allocs/op
BenchmarkText/construct_and_render/write-only/short_text,_little_escaping-10               	24978750	        47.70 ns/op	      64 B/op	       3 allocs/op
BenchmarkText/render_pre-built/write-only/short_text,_little_escaping-10                   	100000000	        10.90 ns/op	      24 B/op	       1 allocs/op
BenchmarkText/construct_and_render/write-only/long_text,_little_escaping-10                	 3546586	       339.4 ns/op	     976 B/op	       4 allocs/op
BenchmarkText/render_pre-built/write-only/long_text,_little_escaping-10                    	38719200	        31.80 ns/op	     320 B/op	       1 allocs/op
BenchmarkText/construct_and_render/write-only/short_text,_much_escaping-10                 	23761924	        50.16 ns/op	      80 B/op	       3 allocs/op
BenchmarkText/render_pre-built/write-only/short_text,_much_escaping-10                     	100000000	        10.80 ns/op	      32 B/op	       1 allocs/op
BenchmarkText/construct_and_render/write-only/long_text,_much_escaping-10                  	 2709938	       440.7 ns/op	    1552 B/op	       4 allocs/op
BenchmarkText/render_pre-built/write-only/long_text,_much_escaping-10                      	25247544	        47.16 ns/op	     512 B/op	       1 allocs/op
BenchmarkText/construct_and_render/write-only_buffered/short_text,_no_escaping-10          	50406577	        23.52 ns/op	      32 B/op	       2 allocs/op
BenchmarkText/render_pre-built/write-only_buffered/short_text,_no_escaping-10              	100000000	        11.66 ns/op	      16 B/op	       1 allocs/op
BenchmarkText/construct_and_render/write-only_buffered/long_text,_no_escaping-10           	10739191	       111.4 ns/op	     272 B/op	       2 allocs/op
BenchmarkText/render_pre-built/write-only_buffered/long_text,_no_escaping-10               	40478546	        30.95 ns/op	     256 B/op	       1 allocs/op
BenchmarkText/construct_and_render/write-only_buffered/short_text,_little_escaping-10      	23306042	        49.53 ns/op	      64 B/op	       3 allocs/op
BenchmarkText/render_pre-built/write-only_buffered/short_text,_little_escaping-10          	98144455	        12.33 ns/op	      24 B/op	       1 allocs/op
BenchmarkText/construct_and_render/write-only_buffered/long_text,_little_escaping-10       	 3528684	       341.9 ns/op	     976 B/op	       4 allocs/op
BenchmarkText/render_pre-built/write-only_buffered/long_text,_little_escaping-10           	33327006	        37.02 ns/op	     320 B/op	       1 allocs/op
BenchmarkText/construct_and_render/write-only_buffered/short_text,_much_escaping-10        	23347684	        51.74 ns/op	      80 B/op	       3 allocs/op
BenchmarkText/render_pre-built/write-only_buffered/short_text,_much_escaping-10            	98707548	        12.13 ns/op	      32 B/op	       1 allocs/op
BenchmarkText/construct_and_render/write-only_buffered/long_text,_much_escaping-10         	 2674753	       447.7 ns/op	    1552 B/op	       4 allocs/op
BenchmarkText/render_pre-built/write-only_buffered/long_text,_much_escaping-10             	23118806	        52.78 ns/op	     512 B/op	       1 allocs/op
BenchmarkTextf/formatted_text_element-10                                                   	26295630	        44.37 ns/op	      40 B/op	       2 allocs/op
PASS
ok  	maragu.dev/gomponents	128.939s
goos: darwin
goarch: arm64
pkg: maragu.dev/gomponents/components
cpu: Apple M4
BenchmarkJoinAttrs/flat/construct-10         	 3877149	       295.3 ns/op	     512 B/op	      10 allocs/op
BenchmarkJoinAttrs/flat/construct_and_render-10         	 3965832	       302.1 ns/op	     512 B/op	      10 allocs/op
BenchmarkJoinAttrs/nested_groups/construct-10           	 4192413	       286.1 ns/op	     448 B/op	      10 allocs/op
BenchmarkJoinAttrs/nested_groups/construct_and_render-10         	 4006548	       298.7 ns/op	     448 B/op	      10 allocs/op
BenchmarkJoinAttrs/no_match/construct-10                         	 7406504	       161.5 ns/op	     200 B/op	       4 allocs/op
BenchmarkJoinAttrs/no_match/construct_and_render-10              	 6806730	       175.3 ns/op	     200 B/op	       4 allocs/op
PASS
ok  	maragu.dev/gomponents/components	7.574s
goos: darwin
goarch: arm64
pkg: maragu.dev/gomponents/html
cpu: Apple M4
BenchmarkRealisticPage/construct_and_render/discarded-10         	   17484	     68624 ns/op	  129686 B/op	    3044 allocs/op
BenchmarkRealisticPage/construct_and_render/buffered-10          	   13102	     91613 ns/op	  129686 B/op	    3044 allocs/op
BenchmarkRealisticPage/render_pre-built_tree/discarded-10        	   27015	     44298 ns/op	    1408 B/op	      18 allocs/op
BenchmarkRealisticPage/render_pre-built_tree/buffered-10         	   17697	     68030 ns/op	    1408 B/op	      18 allocs/op
PASS
ok  	maragu.dev/gomponents/html	5.099s
```

</details>

### I don't like how HTML looks in Go.

First of all, that's not a question. 😉

More seriously, think of gomponents like a DSL for HTML. You're building UI components. Give it a day, and it'll feel natural.

### I'd like to add feature X, can I do that?

First of all, thank you for wanting to contribute! 😊 See [CONTRIBUTING.md](CONTRIBUTING.md) for a guide.

I accept code contributions, especially with new HTML elements and attributes.
I always welcome issues discussing interesting aspects of gomponents, and obviously bug reports and the like.
But otherwise, I consider gomponents pretty much feature complete.

New features to the core library are unlikely to be merged, since I like keeping it simple and the API small.
In particular, new flow control functions (`IfElse`/`Else`) will not be added to the core library.
For generic collection utilities like `Map`, `Filter`, and `Reduce`, see the experimental `x/slices` package.

If there's something missing that you need, I would recommend to keep small helper functions around in your own projects.
And if all else fails, you can always use an [IIFE](https://developer.mozilla.org/en-US/docs/Glossary/IIFE):

```go
func list(ordered bool) Node {
	return func() Node {
		// Do whatever you need to do, imperatively
		if ordered {
			return Ol()
		} else {
			return Ul()
		}
	}()
}
```

### What's up with the specially named elements and attributes?

Unfortunately, there are some name clashes in HTML elements and attributes, so they need an `El` or `Attr` suffix,
to be able to co-exist in the same package in Go.

I've chosen one or the other based on what I think is the common usage:

- `cite`: `Cite` (element) / `CiteAttr` (attribute)
- `data`: `DataEl` (element) / `Data` (attribute)
- `form`: `Form` (element) / `FormAttr` (attribute)
- `label`: `Label` (element) / `LabelAttr` (attribute)
- `style`: `StyleEl` (element) / `Style` (attribute)
- `title`: `TitleEl` (element) / `Title` (attribute)

Deprecated aliases (`CiteEl`, `DataAttr`, `FormEl`, `LabelEl`, `StyleAttr`, `TitleAttr`) also exist for backwards compatibility but should not be used in new code.

<details>
	<summary>Example with `Style` and `StyleEl`</summary>

```go
package html

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func MyPage() Node {
	return HTML5(HTML5Props{
		Title: "My Page",
		Head: []Node{
			StyleEl(Raw("body {background-color: #fff; }")),
		},
		Body: []Node{
			H1(Style("color: #000"), Text("My Page")),
		},
	})
}
```

</details>
