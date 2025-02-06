package html

import (
	g "maragu.dev/gomponents"
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

func Checked() g.Node {
	return g.Attr("checked")
}

func Controls() g.Node {
	return g.Attr("controls")
}

func CrossOrigin(v string) g.Node {
	return g.Attr("crossorigin", v)
}

func DateTime(v string) g.Node {
	return g.Attr("datetime", v)
}

func Defer() g.Node {
	return g.Attr("defer")
}

func Disabled() g.Node {
	return g.Attr("disabled")
}

func Download(v string) g.Node {
	return g.Attr("download", v)
}

func Draggable(v string) g.Node {
	return g.Attr("draggable", v)
}

func Loop() g.Node {
	return g.Attr("loop")
}

func Multiple() g.Node {
	return g.Attr("multiple")
}

func Muted() g.Node {
	return g.Attr("muted")
}

func PlaysInline() g.Node {
	return g.Attr("playsinline")
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

func As(v string) g.Node {
	return g.Attr("as", v)
}

func AutoComplete(v string) g.Node {
	return g.Attr("autocomplete", v)
}

func Charset(v string) g.Node {
	return g.Attr("charset", v)
}

func CiteAttr(v string) g.Node {
	return g.Attr("cite", v)
}

func Class(v string) g.Node {
	return g.Attr("class", v)
}

func Cols(v string) g.Node {
	return g.Attr("cols", v)
}

func ColSpan(v string) g.Node {
	return g.Attr("colspan", v)
}

func Content(v string) g.Node {
	return g.Attr("content", v)
}

// Data attributes automatically have their name prefixed with "data-".
func Data(name, v string) g.Node {
	return g.Attr("data-"+name, v)
}

// DataAttr attributes automatically have their name prefixed with "data-".
//
// Deprecated: Use [Data] instead.
func DataAttr(name, v string) g.Node {
	return Data(name, v)
}

func SlotAttr(v string) g.Node {
  return g.Attr("slot", v)
}

func For(v string) g.Node {
	return g.Attr("for", v)
}

func FormAction(v string) g.Node {
	return g.Attr("formaction", v)
}

func FormAttr(v string) g.Node {
	return g.Attr("form", v)
}

func FormEncType(v string) g.Node {
	return g.Attr("formenctype", v)
}

func FormMethod(v string) g.Node {
	return g.Attr("formmethod", v)
}

func FormNoValidate() g.Node {
	return g.Attr("formnovalidate")
}

func FormTarget(v string) g.Node {
	return g.Attr("formtarget", v)
}

func Height(v string) g.Node {
	return g.Attr("height", v)
}

func Hidden(v string) g.Node {
	return g.Attr("hidden", v)
}

func Href(v string) g.Node {
	return g.Attr("href", v)
}

func ID(v string) g.Node {
	return g.Attr("id", v)
}

func Integrity(v string) g.Node {
	return g.Attr("integrity", v)
}

func LabelAttr(v string) g.Node {
	return g.Attr("label", v)
}

func Lang(v string) g.Node {
	return g.Attr("lang", v)
}

func List(v string) g.Node {
	return g.Attr("list", v)
}

func Loading(v string) g.Node {
	return g.Attr("loading", v)
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

func Placeholder(v string) g.Node {
	return g.Attr("placeholder", v)
}

func Popover(value ...string) g.Node {
	return g.Attr("popover", value...)
}

func PopoverTarget(v string) g.Node {
	return g.Attr("popovertarget", v)
}

func PopoverTargetAction(v string) g.Node {
	return g.Attr("popovertargetaction", v)
}

func Poster(v string) g.Node {
	return g.Attr("poster", v)
}

func Preload(v string) g.Node {
	return g.Attr("preload", v)
}

func ReferrerPolicy(v string) g.Node {
	return g.Attr("referrerpolicy", v)
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

func RowSpan(v string) g.Node {
	return g.Attr("rowspan", v)
}

func Src(v string) g.Node {
	return g.Attr("src", v)
}

func SrcSet(v string) g.Node {
	return g.Attr("srcset", v)
}

func Step(v string) g.Node {
	return g.Attr("step", v)
}

func Style(v string) g.Node {
	return g.Attr("style", v)
}

// Deprecated: Use [Style] instead.
func StyleAttr(v string) g.Node {
	return Style(v)
}

func TabIndex(v string) g.Node {
	return g.Attr("tabindex", v)
}

func Target(v string) g.Node {
	return g.Attr("target", v)
}

func Title(v string) g.Node {
	return g.Attr("title", v)
}

// Deprecated: Use [Title] instead.
func TitleAttr(v string) g.Node {
	return Title(v)
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

func EncType(v string) g.Node {
	return g.Attr("enctype", v)
}

func Dir(v string) g.Node {
	return g.Attr("dir", v)
}
