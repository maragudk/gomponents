// Package components provides high-level components that are composed of low-level elements and attributes.
package components

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/attr"
	"github.com/maragudk/gomponents/el"
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
		lang = attr.Lang(p.Language)
	}
	if p.Description != "" {
		description = el.Meta(attr.Name("description"), attr.Content(p.Description))
	}
	return el.Document(
		el.HTML(lang,
			el.Head(
				el.Meta(attr.Charset("utf-8")),
				el.Meta(attr.Name("viewport"), attr.Content("width=device-width, initial-scale=1")),
				el.Title(p.Title),
				description,
				g.Group(p.Head),
			),
			el.Body(g.Group(p.Body)),
		),
	)
}
