package htmx

import (
	"io"
	"strings"

	g "github.com/maragudk/gomponents"
)

// HxAttr creates an attribute DOM Node with a name and optional value.
// If only a name is passed, it's a name-only (boolean) attribute (like "required").
// If a name and value are passed, it's a name-value attribute (like `class='header'`).
// More than one value make Attr panic.
// Use this if you want to set attribute in DOM with double quotes in it. The content will always be raw
func HxAttr(name string, value ...string) g.Node {
	switch len(value) {
	case 0:
		return &attr{name: name}
	case 1:
		return &attr{name: name, value: &value[0]}
	default:
		panic("attribute must be just name or name and value pair")
	}
}

type attr struct {
	name  string
	value *string
}

// Render satisfies Node.
func (a *attr) Render(w io.Writer) error {
	if a.value == nil {
		_, err := w.Write([]byte(" " + a.name))
		return err
	}
	_, err := w.Write([]byte(" " + a.name + `='` + *a.value + `'`))
	return err
}

// Type satisfies nodeTypeDescriber.
func (a *attr) Type() g.NodeType {
	return g.AttributeType
}

// String satisfies fmt.Stringer.
func (a *attr) String() string {
	var b strings.Builder
	_ = a.Render(&b)
	return b.String()
}
