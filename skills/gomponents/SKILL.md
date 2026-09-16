---
name: gomponents
description: Building, rendering, or editing any HTML or UI in a Go application means using gomponents — ALWAYS invoke this skill as your FIRST action, before any Read, Write, Edit, or Bash. This fires for any view, page, layout, component, form, table, navbar, footer, or list that renders to HTML; any function returning a `Node`; anything in the `html` package; and converting HTML or templates into Go. It is a hard requirement — the user writes all HTML through gomponents, a pure-Go component library whose conventions (dot imports, `Node` composition, `Map`/`If`/`Iff`, `Group`, HTML5 documents) ordinary Go and HTML habits get wrong. Load it even when a general Go skill also applies. Skip only when no HTML or view code is touched (pure database, handler, or business logic; plain `.css`/`.js`/`.html` files; conceptual questions with no code).
---

# gomponents

gomponents is HTML components in pure Go. A view is a Go function that returns a `Node`, and a `Node` renders itself to an `io.Writer` as HTML5. There is no template language, no code generation, and no dependencies.

```sh
go get maragu.dev/gomponents
```

The import path is `maragu.dev/gomponents`, not the GitHub URL. It needs Go 1.18 or later. The API is stable and feature-complete: new elements and attributes get added, control-flow helpers such as `IfElse` do not. Keep small helpers like that in your own project.

## Mental model

Everything is a `Node`:

```go
type Node interface {
	Render(w io.Writer) error
}
```

Elements, attributes, text, and groups all implement it, and all are passed as children to the same variadic element functions. Four rules explain nearly all behaviour:

1. **Attributes go in the tag, everything else between the tags.** An element renders its attribute children inside the opening tag and the rest between the tags, in whatever order you pass them. By convention, write attributes first.
2. **`nil` children are skipped.** `If` depends on this: it returns `nil` when the condition is false.
3. **Nodes are lazy.** Building a tree captures values in closures; nothing is written until `Render` runs. Build a fresh tree per request and render it straight to the response.
4. **Nothing is added.** No whitespace appears between nodes: `Span(Text("a")), Span(Text("b"))` renders `<span>a</span><span>b</span>`. Add `Text(" ")` where inline spacing matters.

```go
package html

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Navbar(authenticated bool, currentPath string) Node {
	return Nav(Class("navbar"),
		navbarLink("/", "Home", currentPath),
		navbarLink("/about", "About", currentPath),
		If(authenticated, navbarLink("/profile", "Profile", currentPath)),
	)
}

func navbarLink(href, name, currentPath string) Node {
	return A(Href(href), Classes{"is-active": currentPath == href}, Text(name))
}
```

## Imports and package layout

Dot-import the three main packages so views read like HTML. This is the strongly preferred style unless the project already imports them another way. Whichever style the project uses, use it in every file that mentions a `Node`, handler files included: write `Node`, not `g.Node`.

Put views in a package named `html`, as the gomponents example app and starter kit do, unless the project already has a views package under another name. A package named `html` can dot-import `maragu.dev/gomponents/html` without conflict.

Inside that package, export only what another package uses, which is usually just the pages:

- **Exported pages** are named `XxxPage` and take a props struct: `func LoginPage(props LoginPageProps) Node`.
- **Unexported building blocks** are named in lower case after the element they wrap: `button`, `input`, `label`, `a`, `card`, `container`. Lower-case names cannot collide with the dot-imported `Button`, `Input`, `Label`, and `A`. When a block does need exporting, give it a specific name (`SubmitButton`, not `Button`) rather than dropping the dot import.

`staticcheck` flags dot imports by default, so whitelist them in `.golangci.yml`:

```yaml
version: "2"
linters:
  settings:
    staticcheck:
      dot-import-whitelist:
        - "maragu.dev/gomponents"
        - "maragu.dev/gomponents/components"
        - "maragu.dev/gomponents/html"
```

The `http` package clashes with `net/http`, so alias it: `ghttp "maragu.dev/gomponents/http"`.

## Core package: `maragu.dev/gomponents`

| Function | What it does |
|---|---|
| `El(name, children...)` | Any element. Use it for elements the `html` package lacks: custom elements, or SVG children such as `El("path", Attr("d", "..."))`. |
| `Attr(name)` / `Attr(name, value)` | Boolean or valued attribute. Use it for attributes the `html` package lacks: `Attr("hx-get", "/items")`, `Attr("onclick", "...")`. More than one value panics. |
| `Text(s)` / `Textf(format, args...)` | HTML-escaped text. The default for all content, and the only choice for anything user-controlled. |
| `Raw(s)` / `Rawf(format, args...)` | Unescaped text. For markup you wrote yourself: inline SVG, `<script>` and `<style>` bodies, entities like `&copy;`. Never for user content. |
| `Map(slice, func(T) Node) Group` | Turns a slice of data into nodes. Pass the result directly as a child; it is already a `Node`. |
| `Group{...}` / `Group(nodes)` | A `[]Node` that renders as one node. Use it to return several siblings, or to pass a `children ...Node` slice on. |
| `If(cond, node)` | `node` if `cond`, else `nil`. The node argument is evaluated either way. |
| `Iff(cond, func() Node)` | Like `If`, but the function only runs when `cond` is true. |

Element and attribute *names* are written verbatim, so they must be trusted compile-time values. Attribute *values*, `Text`, and `Textf` are escaped. Never build a name from user input, as in `Data(userKey, v)` or `El(tagFromRequest)`.

Every built-in node implements `fmt.Stringer`, so `fmt.Println(node)` prints its HTML when debugging.

### `If` evaluates eagerly, `Iff` doesn't

`If` is a plain function call, so Go builds its second argument before `If` sees the condition. When that argument reads through a nil pointer, it panics:

```go
var user *User // nil when nobody is logged in

// Panics when user is nil: user.Name runs before If checks the condition.
If(user != nil, Text(user.Name))

// Safe: the function only runs when the condition is true.
Iff(user != nil, func() Node { return Text(user.Name) })
```

Use `If` when building the node is cheap and cannot fail. Use `Iff` when the node reads through a pointer, indexes a slice, or does work worth skipping.

### Conditionals and loops

There is no `IfElse`, on purpose. Use two `If`s with opposite conditions, or a function with a `switch` when the branches are more than one node:

```go
P(
	If(ok, Text("Saved.")),
	If(!ok, Text("Something went wrong.")),
)
```

`Map` has no index. When you need one, keep a counter in the closure, or use the experimental `maragu.dev/gomponents/x/slices` package, whose `Map` and `Filter` pass the index and return plain slices to spread with `...`. Packages under `x/` may change without notice.

```go
Ul(Map(items, func(item Item) Node {
	return Li(Text(item.Name))
}))
```

## The `html` package: elements and attributes

Element functions take `...Node`. Attribute functions take one `string` value, or nothing for boolean attributes (`Required()`, `Disabled()`, `Checked()`). All values are strings, so convert numbers yourself: `Width(strconv.Itoa(w))`.

Names follow the HTML names with Go casing: word boundaries are capitalized (`ColSpan`, `TabIndex`, `MaxLength`, `FieldSet`, `FigCaption`, `SrcSet`, `AutoComplete`) and initialisms are upper-case (`ID`, `HTML`, `SVG`, `IFrame`, `THead`, `TBody`, `TFoot`, `HGroup`). A few HTML names are both an element and an attribute, so one side gets a suffix:

| HTML name | Element | Attribute |
|---|---|---|
| `cite` | `Cite` | `CiteAttr` |
| `data` | `DataEl` | `Data(name, value)`, renders `data-<name>` |
| `form` | `Form` | `FormAttr` |
| `label` | `Label` | `LabelAttr` |
| `slot` | `SlotEl` | `SlotAttr` |
| `style` | `StyleEl` | `Style` |
| `title` | `TitleEl` | `Title` |

`CiteEl`, `DataAttr`, `FormEl`, `LabelEl`, `StyleAttr`, and `TitleAttr` still compile but are deprecated; do not use them in new code.

`Data("id", v)` renders `data-id="..."` and `Aria("label", v)` renders `aria-label="..."`. `Doctype(node)` prefixes `<!doctype html>`; `HTML5` already calls it.

When unsure whether a helper exists or how it is cased, check rather than guess:

```sh
go doc maragu.dev/gomponents/html | grep -i colspan
```

If nothing turns up, `El` and `Attr` produce the same output a dedicated helper would. Left out on purpose: SVG child elements, `on*` event handler attributes, and framework attributes. For htmx and Datastar there are typed helpers in `maragu.dev/gomponents-htmx` and `maragu.dev/gomponents-datastar`; otherwise `Attr("hx-get", ...)` is fine.

## The `components` package

**`HTML5`** renders a complete document: doctype, `<html>`, a `<head>` with charset, viewport, title, and optional description, and the `<body>`. Attribute nodes in `Body` or `Head` are placed on that element. Every app has one layout function like this, and every page calls it:

```go
func page(title string, body ...Node) Node {
	return HTML5(HTML5Props{
		Title:    title + " - MyApp",
		Language: "en",
		Head: Group{
			Link(Rel("stylesheet"), Href("/static/app.css")),
			Script(Src("/static/app.js"), Defer()),
		},
		Body: Group{Class("min-h-full bg-gray-50"),
			navbar(),
			Main(Group(body)),
			footer(),
		},
	})
}
```

**`Classes`** is a `map[string]bool` that renders as one `class` attribute holding the keys whose value is true, sorted. Use it for conditional classes instead of building the string yourself. It works anywhere `Class` does.

```go
Button(Classes{"btn": true, "btn-primary": primary, "opacity-50": disabled}, Text(label))
```

**`JoinAttrs`** merges every attribute with the given name among the direct children into one attribute. It looks through groups at any depth, but not into child elements. This lets blocks build on each other, each adding its own classes, while the caller adds more:

```go
func button(children ...Node) Node {
	return Button(JoinAttrs("class", Group(children), Class("btn")))
}

func primaryButton(children ...Node) Node {
	return button(Class("btn-primary"), Group(children))
}

primaryButton(Class("mt-4"), Text("Save"))
// <button class="btn-primary mt-4 btn">Save</button>
```

## The `http` package

`Adapt` turns a handler that returns `(Node, error)` into an `http.HandlerFunc`:

- The node is rendered even when there is an error, so return an error page along with the error.
- If the error has a `StatusCode() int` method, that status is sent; any other error sends 500.
- A `nil` node writes nothing; return it after you have issued a redirect yourself.

```go
type notFoundError struct{}

func (notFoundError) Error() string   { return "not found" }
func (notFoundError) StatusCode() int { return http.StatusNotFound }

mux.Handle("GET /users/{id}", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
	u, err := users.Get(r.Context(), r.PathValue("id"))
	if errors.Is(err, ErrNotFound) {
		return html.NotFoundPage(), notFoundError{}
	}
	if err != nil {
		return html.ErrorPage(), err
	}
	return html.UserPage(u), nil
}))
```

## Common patterns

- **Building blocks take `children ...Node`** and pass them on with `Group(children)`. Groups are transparent, so a caller's attributes (`ID`, `Attr("hx-target", ...)`) are placed on the root element: `input(Type("email"), Name("email"), Required())` needs no parameters of its own.
- **Placement classes come from the caller.** A block used in several places takes a `class string` first and joins it with its own classes. Use `JoinAttrs` instead when callers pass attribute nodes rather than a string.
- **One look, several elements: pass the element function.** `heading(H1, "", Text(title))` and `heading(Dt, "", Text(name))` share one style.
- **Many options: use a props struct**, as `HTML5Props` does, rather than a long parameter list.
- **Views take plain values.** Load data and handle errors before calling the view; `Render` fails only on a write error, so a component has no error path.
- **Dynamic attributes** work like dynamic elements, because `nil` is skipped inside the tag too: `If(disabled, Disabled())`, `Classes{...}`, `Value(u.Email)`.
- **Fragments** for htmx and similar are blocks rendered without the layout. Return a `Group` when a fragment has several root elements.

```go
func container(class string, children ...Node) Node {
	return Div(classes(class, "mx-auto max-w-7xl px-6 lg:px-8"), Group(children))
}

func heading(el func(...Node) Node, class string, children ...Node) Node {
	return el(classes(class, "text-4xl font-bold tracking-tight"), Group(children))
}

// classes puts the caller's placement classes before the block's own.
func classes(own, base string) Node {
	if own == "" {
		return Class(base)
	}
	return Class(own + " " + base)
}
```

## Gotchas

- **`<script>` and `<style>` bodies need `Raw`.** `Text` would turn `&&` into `&amp;&amp;` and break the code. Keep user data out of those bodies; pass it through `Data` attributes instead, which are escaped and unescaped by the browser.
- **Two `Class` calls make two attributes.** Combine them into one string, use `Classes`, or use `JoinAttrs`.
- **Void elements silently ignore non-attribute children.** `Img(Text("x"))` renders `<img>` with no error.
- **`Attr` with two or more values panics** at construction time, not at render time.
- **`Classes` sorts.** `Classes{"tier": true, "popular": true}` renders `class="popular tier"`. Fine for CSS, but a test that pins the exact string must expect the sorted order.

## Testing components

Test exported pages and components: one `TestComponent` function per exported component, with subtests for one happy path, the error cases, and the edge cases. Test what matters in each: the branch that depends on input, the value that must be escaped, the link that appears for one kind of user and not another. An expected string for a whole page restates the component and breaks on every unrelated change, so reserve exact-string comparison for small blocks and fragments, even when converting existing HTML. gomponents has no public test helpers, so render to a `strings.Builder` through a small helper and check with `strings.Contains`:

```go
func TestNavbar(t *testing.T) {
	t.Run("links to the profile only when authenticated", func(t *testing.T) {
		if got := render(t, Navbar(false, "/")); strings.Contains(got, `href="/profile"`) {
			t.Fatalf("unauthenticated navbar has a profile link: %s", got)
		}
		if got := render(t, Navbar(true, "/")); !strings.Contains(got, `<a href="/profile"`) {
			t.Fatalf("authenticated navbar has no profile link: %s", got)
		}
	})
}

func render(t *testing.T, n Node) string {
	t.Helper()
	var b strings.Builder
	if err := n.Render(&b); err != nil {
		t.Fatal(err)
	}
	return b.String()
}
```

Test handlers through `Adapt` with `httptest.NewRecorder`, and assert on the status code and a few body markers.

Unexported blocks are tested from inside the package. To reach them from an external `_test` package, re-export them in an `export_internal_test.go` file with `var Card = card`. That re-export collides with dot-imported element names just as an exported function would, so it works for `card` but not for `button`.

## Further reading

- API reference: https://pkg.go.dev/maragu.dev/gomponents
- Guides: https://www.gomponents.com
- Full application template: https://github.com/maragudk/gomponents-starter-kit
