// Package components provides high-level components that are composed of low-level elements and attributes.
package components

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

// DocumentProps for HTML5.
// Title is set no matter what, Description and Language elements only if the strings are non-empty.
type DocumentProps struct {
	Title       string
	Description string
	Language    string
	Head        []g.Node
	Body        []g.Node
}

// HTML5 document template.
func HTML5(p DocumentProps) g.NodeFunc {
	var lang, description g.Node
	if p.Language != "" {
		lang = Lang(p.Language)
	}
	if p.Description != "" {
		description = Meta(Name("description"), Content(p.Description))
	}
	return Document(
		HTML(lang,
			Head(
				Meta(Charset("utf-8")),
				Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
				TitleEl(p.Title),
				description,
				g.Group(p.Head),
			),
			Body(g.Group(p.Body)),
		),
	)
}
