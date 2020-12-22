package html

import (
	g "github.com/maragudk/gomponents"
)

func Async() g.Node {
	return g.Attr("async")
}

func AutoFocus() g.Node {
	return g.Attr("autofocus")
}

func AutoPlay() g.Node {
	return g.Attr("autoplay")
}

func Controls() g.Node {
	return g.Attr("controls")
}

func Defer() g.Node {
	return g.Attr("defer")
}

func Disabled() g.Node {
	return g.Attr("disabled")
}

func Multiple() g.Node {
	return g.Attr("multiple")
}

func ReadOnly() g.Node {
	return g.Attr("readonly")
}

func Required() g.Node {
	return g.Attr("required")
}

func Selected() g.Node {
	return g.Attr("selected")
}

func Accept(v string) g.Node {
	return g.Attr("accept", v)
}

func Action(v string) g.Node {
	return g.Attr("action", v)
}

func Alt(v string) g.Node {
	return g.Attr("alt", v)
}

// Aria attributes automatically have their name prefixed with "aria-".
func Aria(name, v string) g.Node {
	return g.Attr("aria-"+name, v)
}

func AutoComplete(v string) g.Node {
	return g.Attr("autocomplete", v)
}

func Charset(v string) g.Node {
	return g.Attr("charset", v)
}

func Class(v string) g.Node {
	return g.Attr("class", v)
}

func Cols(v string) g.Node {
	return g.Attr("cols", v)
}

func Content(v string) g.Node {
	return g.Attr("content", v)
}

func For(v string) g.Node {
	return g.Attr("for", v)
}

func FormAttr(v string) g.Node {
	return g.Attr("form", v)
}

func Height(v string) g.Node {
	return g.Attr("height", v)
}

func Href(v string) g.Node {
	return g.Attr("href", v)
}

func ID(v string) g.Node {
	return g.Attr("id", v)
}

func Lang(v string) g.Node {
	return g.Attr("lang", v)
}

func Max(v string) g.Node {
	return g.Attr("max", v)
}

func MaxLength(v string) g.Node {
	return g.Attr("maxlength", v)
}

func Method(v string) g.Node {
	return g.Attr("method", v)
}

func Min(v string) g.Node {
	return g.Attr("min", v)
}

func MinLength(v string) g.Node {
	return g.Attr("minlength", v)
}

func Name(v string) g.Node {
	return g.Attr("name", v)
}

func Pattern(v string) g.Node {
	return g.Attr("pattern", v)
}

func Preload(v string) g.Node {
	return g.Attr("preload", v)
}

func Placeholder(v string) g.Node {
	return g.Attr("placeholder", v)
}

func Rel(v string) g.Node {
	return g.Attr("rel", v)
}

func Role(v string) g.Node {
	return g.Attr("role", v)
}

func Rows(v string) g.Node {
	return g.Attr("rows", v)
}

func Src(v string) g.Node {
	return g.Attr("src", v)
}

func StyleAttr(v string) g.Node {
	return g.Attr("style", v)
}

func TabIndex(v string) g.Node {
	return g.Attr("tabindex", v)
}

func Target(v string) g.Node {
	return g.Attr("target", v)
}

func TitleAttr(v string) g.Node {
	return g.Attr("title", v)
}

func Type(v string) g.Node {
	return g.Attr("type", v)
}

func Value(v string) g.Node {
	return g.Attr("value", v)
}

func Width(v string) g.Node {
	return g.Attr("width", v)
}