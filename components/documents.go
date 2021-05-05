// Package components provides high-level components and helpers that are composed of low-level elements and attributes.
package components

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

// HTML5Props for HTML5.
// Title is set no matter what, Description and Language elements only if the strings are non-empty.
type HTML5Props struct {
	Title       string
	Description string
	Language    string
	Head        []g.Node
	Body        []g.Node
}

// HTML5 document template.
func HTML5(p HTML5Props) g.Node {
	return Doctype(
		HTML(g.If(p.Language != "", Lang(p.Language)),
			Head(
				Meta(Charset("utf-8")),
				Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
				TitleEl(g.Text(p.Title)),
				g.If(p.Description != "", Meta(Name("description"), Content(p.Description))),
				g.Group(p.Head),
			),
			Body(g.Group(p.Body)),
		),
	)
}
