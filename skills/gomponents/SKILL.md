---
name: gomponents
description: Building, rendering, or editing any HTML or UI in a Go application means using gomponents — ALWAYS invoke this skill as your FIRST action, before any Read, Write, Edit, or Bash. This fires for any view, page, layout, component, form, table, navbar, footer, or list that renders to HTML; any function returning a `Node`; anything in the `html` package; and converting HTML or templates into Go. It is a hard requirement — the user writes all HTML through gomponents, a pure-Go component library whose conventions (dot imports, `Node` composition, `Map`/`If`/`Iff`, `Group`, HTML5 documents) ordinary Go and HTML habits get wrong. Skip only when no HTML or view code is touched (pure database, handler, or business logic; plain `.css`/`.js`/`.html` files; conceptual questions with no code).
---

# gomponents

gomponents is HTML components in pure Go. A component is a Go function that returns a `Node`, and a `Node` renders itself to an `io.Writer` as HTML5. There is no template language, no code generation, and no dependencies.

```sh
go get maragu.dev/gomponents@latest
```

## Mental model

Everything is a `Node`:

```go
type Node interface {
	Render(w io.Writer) error
}
```

Elements, attributes, and text all implement it, and all are passed as children to the same variadic element functions. The result is like a DSL for HTML, in valid Go: an element function takes children, an attribute function takes a value, and nesting calls mirrors nesting tags. Attribute children go in the tag and everything else goes between the tags, whatever order you pass them in; by convention, write attributes first.

```html
<a href="/about" class="nav-link">About</a>
```

```go
A(Href("/about"), Class("nav-link"), Text("About"))
```

The idiomatic way is to write it declaratively, but a component is an ordinary function and can be as imperative as it needs to be. gomponents does not check that the result is valid HTML, by design, just as writing HTML by hand doesn't.

The core package provides these functions. Everything else is in the `html` and `components` packages.

| Function | What it does |
|---|---|
| `Text(s)` / `Textf(format, args...)` | HTML-escaped text. The default for all content, and the strongly recommended choice for unsanitized user content. |
| `Raw(s)` / `Rawf(format, args...)` | Unescaped text. For markup you wrote yourself: inline SVG, `<script>` and `<style>` bodies. Never for unsanitized user content. |
| `Map(slice, func(T) Node) Group` | Turns a slice of data into nodes. |
| `Group{...}` / `Group(nodes)` | A `[]Node` that renders as one node. Use it to return several siblings, or to pass a `children ...Node` slice on. |
| `If(cond, node)` | `node` if `cond`, else `nil`, and `nil` children are skipped. The node argument is evaluated either way. |
| `Iff(cond, func() Node)` | Like `If`, but the function only runs when `cond` is true. |
| `El(name, children...)` | Any element. Use it for elements the `html` package lacks: custom elements, or SVG children such as `El("path", Attr("d", "..."))`. |
| `Attr(name)` / `Attr(name, value)` | Boolean or valued attribute. Use it for attributes the `html` package lacks: `Attr("hx-get", "/items")`, `Attr("onclick", "...")`. More than one value panics. |

Element and attribute *names* are written verbatim, so they must be trusted compile-time values. Attribute *values*, `Text`, and `Textf` are escaped. Never build a name from user input, as in `Data(userKey, v)` or `El(tagFromRequest)`.

An example component using all of the above except `El` and `Attr`:

```go
package html

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type NavLink struct {
	Href, Text string
}

type User struct {
	Name   string
	Unread int
}

// Navbar shows the site links and, for a logged-in user, their name and a log out link.
// user is nil for a visitor.
func Navbar(links []NavLink, currentPath string, user *User) Node {
	return Nav(Class("navbar"),
		Map(links, func(l NavLink) Node {
			return A(Href(l.Href), Classes{"link": true, "is-active": l.Href == currentPath}, Text(l.Text))
		}),

		If(user == nil, A(Href("/login"), Text("Log in"))),

		// user.Name must only run when user is not nil.
		Iff(user != nil, func() Node {
			return Group{
				Span(Textf("%v (%v unread)", user.Name, user.Unread)),
				A(Href("/logout"), Text("Log out")),
			}
		}),

		Raw(`<svg viewBox="0 0 16 16" width="16" height="16"><circle cx="8" cy="8" r="7"/></svg>`),
	)
}
```

## Common patterns

- **`If` builds its node before checking the condition.** `If(user != nil, Text(user.Name))` panics when `user` is nil, because `user.Name` runs first. Use `Iff` when the node reads through a pointer, indexes a slice, or does work worth skipping.
- **There is no `IfElse`.** Use two `If`s with opposite conditions, or a function with a `switch` when the branches are more than one node.
- **`Map` has no index.** When you need one, keep a counter in the closure, or use `maragu.dev/gomponents/x/slices`, whose `Map` and `Filter` pass the index and return plain slices to spread with `...`. Packages under `x/` are experimental and may change, although they probably won't.
- **Building blocks take `children ...Node`** and pass them on with `Group(children)`. Groups are transparent, so a caller's attributes (`ID`, `Name`) are placed on the root element. A block like `input(children ...Node)` that only adds classes to `Input` needs no parameters, because callers pass `Type("email")`, `Name("email")`, or `Required()` as children and they land on the input.
- **Callers add classes through `JoinAttrs`.** A block joins its own `Class` with whatever the caller passes, so `card(Class("mt-4"), ...)` renders one `class` attribute. See the `components` package below.
- **Dynamic attributes** work like dynamic elements, because `nil` is skipped inside the tag too: `If(disabled, Disabled())`, `Classes{...}`, `Value(u.Email)`.
- **A fragment is just a component** returned without the layout; gomponents has no special notion of it. Return a `Group` when it has several root elements.
- **Print a node to see its HTML.** Every built-in node implements `fmt.Stringer`, so `fmt.Println(node)` works when debugging.
- **`<script>` and `<style>` bodies need `Raw`.** `Text` would turn `&&` into `&amp;&amp;` and break the code. Keep unsanitized user data out of those bodies; pass it through `Data` attributes instead, which are escaped and unescaped by the browser.
- **Two `Class` calls make two attributes**, and the browser keeps only the first. Combine them into one string, use `Classes`, or use `JoinAttrs`.
- **Void elements silently ignore non-attribute children.** `Img(Text("x"))` renders `<img>` with no error.

## Imports and package layout

Dot-import the three main packages so components read like HTML. This is the strongly preferred style unless the project already imports them another way. Whichever style the project uses, use it in every file that mentions a `Node`, handler files included: write `Node`, not `g.Node`.

Put components in a package named `html`, unless the project already has one under another name. A package named `html` can dot-import `maragu.dev/gomponents/html` without conflict.

Inside that package, export only what another package uses, which is usually just the pages and fragments:

- **Exported pages and fragments** are named for what they are and take a props struct: `func LoginPage(props LoginPageProps) Node`, `func UserRow(u User) Node`.
- **Unexported building blocks** are camelCase with a lower-case first letter, named after the element they wrap: `button`, `input`, `label`, `a`, `card`, `container`. Names with a lower-case first letter cannot collide with the dot-imported `Button`, `Input`, `Label`, and `A`. When a block does need exporting, give it a specific name (`SubmitButton`, not `Button`) rather than dropping the dot import.

The `http` package often clashes with `net/http`, so consider aliasing it: `ghttp "maragu.dev/gomponents/http"`.

## The `html` package: elements and attributes

Element functions take `...Node`. Attribute functions take one `string` value, or nothing for boolean attributes (`Required()`, `Disabled()`, `Checked()`). A few take `...string`, for attributes that are valid both with and without a value. All values are strings, so convert numbers yourself: `Width(strconv.Itoa(w))`.

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

`CiteEl`, `DataAttr`, `FormEl`, `LabelEl`, `StyleAttr`, and `TitleAttr` still compile but are deprecated.

`Data("id", v)` renders `data-id="..."` and `Aria("label", v)` renders `aria-label="..."`.

When unsure whether a helper exists or how it is cased, check rather than guess:

```sh
go doc maragu.dev/gomponents/html | grep -i colspan
```

If nothing turns up, `El` and `Attr` produce the same output a dedicated helper would.

## The `components` package

**`HTML5`** renders a complete document: doctype, `<html>`, a `<head>` with charset, viewport, title, and optional description, and the `<body>`. Attribute nodes in `Body` or `Head` are placed on that element. Most apps have one layout function like this, and their pages call it:

```go
func page(title, description string, body ...Node) Node {
	return HTML5(HTML5Props{
		Title:       title + " - MyApp",
		Description: description,
		Language:    "en",
		HTMLAttrs:   Group{Class("h-full")},
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

- The node is rendered even when there is an error. For pages, return an error page along with the error, as for a 403, 404, or 500. For fragments, `nil, err` is common and sends only the status.
- If the error has a `StatusCode() int` method, that status is sent; any other error sends 500.
- A `nil` node writes nothing, so also return it after writing the response yourself, such as a redirect.

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

## Testing components

Test exported pages and components: one `TestComponent` function per exported component, with subtests for one happy path, the error cases, and the edge cases. Test what matters in each, for example the branch that depends on input, the value that must be escaped, or the link that appears for one kind of user and not another. An expected string for a whole page restates the component and breaks on every unrelated change, so reserve exact-string comparison for small blocks and fragments, even when converting existing HTML. Programming errors in a component, such as `Attr` given two values, panic when the component is built, so rendering each component in a test catches them. gomponents has no public test helpers, so render to a `strings.Builder` through a small helper and check with `strings.Contains`:

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

To test unexported blocks, use an internal test package, or re-export them for an external `_test` package in an `export_internal_test.go` file with `var Card = card`, according to project preference. That re-export collides with dot-imported element names just as an exported function would, so it works for `card` but not for `button`.

For linter setup, see `references/linting.md`.

## Further reading

- API reference: `go doc maragu.dev/gomponents`, and likewise for `/html`, `/components`, and `/http`.
- Full application template: https://github.com/maragudk/gomponents-starter-kit
