package doc

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/attr"
	"github.com/maragudk/gomponents/el"
)

// Props for documents.
// Title is set no matter what, Description and Language only if non-empty.
type Props struct {
	Title       string
	Description string
	Language    string
	Head        []g.Node
	Body        []g.Node
}

// HTML5 document template.
func HTML5(p Props) g.NodeFunc {
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
