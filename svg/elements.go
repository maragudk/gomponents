// Package svg provides common SVG elements and attributes.
// See https://developer.mozilla.org/en-US/docs/Web/SVG/Element for an overview.
package svg

import (
	g "github.com/maragudk/gomponents"
)

func Path(children ...g.Node) g.Node {
	return g.El("path", children...)
}

func SVG(children ...g.Node) g.Node {
	return g.El("svg", g.Attr("xmlns", "http://www.w3.org/2000/svg"), g.Group(children))
}

func A(children ...g.Node) g.Node {
	return g.El("a", children...)
}

func Animate(children ...g.Node) g.Node {
	return g.El("animate", children...)
}

func AnimateMotion(children ...g.Node) g.Node {
	return g.El("animateMotion", children...)
}

func AnimateTransform(children ...g.Node) g.Node {
	return g.El("animateTransform", children...)
}

func Circle(children ...g.Node) g.Node {
	return g.El("circle", children...)
}

func ClipPathEl(children ...g.Node) g.Node {
	return g.El("clipPath", children...)
}

func Defs(children ...g.Node) g.Node {
	return g.El("defs", children...)
}

func Desc(children ...g.Node) g.Node {
	return g.El("desc", children...)
}

func Ellipse(children ...g.Node) g.Node {
	return g.El("ellipse", children...)
}

func FeBlend(children ...g.Node) g.Node {
	return g.El("feBlend", children...)
}

func FeColorMatrix(children ...g.Node) g.Node {
	return g.El("feColorMatrix", children...)
}

func FeComponentTransfer(children ...g.Node) g.Node {
	return g.El("feComponentTransfer", children...)
}

func FeComposite(children ...g.Node) g.Node {
	return g.El("feComposite", children...)
}

func FeConvolveMatrix(children ...g.Node) g.Node {
	return g.El("feConvolveMatrix", children...)
}

func FeDiffuseLighting(children ...g.Node) g.Node {
	return g.El("feDiffuseLighting", children...)
}

func FeDisplacementMap(children ...g.Node) g.Node {
	return g.El("feDisplacementMap", children...)
}

func FeDistantLight(children ...g.Node) g.Node {
	return g.El("feDistantLight", children...)
}

func FeDropShadow(children ...g.Node) g.Node {
	return g.El("feDropShadow", children...)
}

func FeFlood(children ...g.Node) g.Node {
	return g.El("feFlood", children...)
}

func FeFuncA(children ...g.Node) g.Node {
	return g.El("feFuncA", children...)
}

func FeFuncB(children ...g.Node) g.Node {
	return g.El("feFuncB", children...)
}

func FeFuncG(children ...g.Node) g.Node {
	return g.El("feFuncG", children...)
}

func FeFuncR(children ...g.Node) g.Node {
	return g.El("feFuncR", children...)
}

func FeGaussianBlur(children ...g.Node) g.Node {
	return g.El("feGaussianBlur", children...)
}

func FeImage(children ...g.Node) g.Node {
	return g.El("feImage", children...)
}

func FeMerge(children ...g.Node) g.Node {
	return g.El("feMerge", children...)
}

func FeMergeNode(children ...g.Node) g.Node {
	return g.El("feMergeNode", children...)
}

func FeMorphology(children ...g.Node) g.Node {
	return g.El("feMorphology", children...)
}

func FeOffset(children ...g.Node) g.Node {
	return g.El("feOffset", children...)
}

func FePointLight(children ...g.Node) g.Node {
	return g.El("fePointLight", children...)
}

func FeSpecularLighting(children ...g.Node) g.Node {
	return g.El("feSpecularLighting", children...)
}

func FeSpotLight(children ...g.Node) g.Node {
	return g.El("feSpotLight", children...)
}

func FeTile(children ...g.Node) g.Node {
	return g.El("feTile", children...)
}

func FeTurbulence(children ...g.Node) g.Node {
	return g.El("feTurbulence", children...)
}

func FilterEl(children ...g.Node) g.Node {
	return g.El("filter", children...)
}

func Filter(children ...g.Node) g.Node {
	return FilterEl(children...)
}

func ForeignObject(children ...g.Node) g.Node {
	return g.El("foreignObject", children...)
}

func G(children ...g.Node) g.Node {
	return g.El("g", children...)
}

func Image(children ...g.Node) g.Node {
	return g.El("image", children...)
}

func Line(children ...g.Node) g.Node {
	return g.El("line", children...)
}

func LinearGradient(children ...g.Node) g.Node {
	return g.El("linearGradient", children...)
}

func Marker(children ...g.Node) g.Node {
	return g.El("marker", children...)
}

func MaskEl(children ...g.Node) g.Node {
	return g.El("mask", children...)
}

func Mask(children ...g.Node) g.Node {
	return MaskEl(children...)
}

func Metadata(children ...g.Node) g.Node {
	return g.El("metadata", children...)
}

func Mpath(children ...g.Node) g.Node {
	return g.El("mpath", children...)
}

func Pattern(children ...g.Node) g.Node {
	return g.El("pattern", children...)
}

func Polygon(children ...g.Node) g.Node {
	return g.El("polygon", children...)
}

func Polyline(children ...g.Node) g.Node {
	return g.El("polyline", children...)
}

func RadialGradient(children ...g.Node) g.Node {
	return g.El("radialGradient", children...)
}

func Rect(children ...g.Node) g.Node {
	return g.El("rect", children...)
}

func Script(children ...g.Node) g.Node {
	return g.El("script", children...)
}

func Set(children ...g.Node) g.Node {
	return g.El("set", children...)
}

func Stop(children ...g.Node) g.Node {
	return g.El("stop", children...)
}

func StyleEl(children ...g.Node) g.Node {
	return g.El("style", children...)
}

func Style(children ...g.Node) g.Node {
	return StyleEl(children...)
}

func Switch(children ...g.Node) g.Node {
	return g.El("switch", children...)
}

func Symbol(children ...g.Node) g.Node {
	return g.El("symbol", children...)
}

func Text(children ...g.Node) g.Node {
	return g.El("text", children...)
}

func TextPath(children ...g.Node) g.Node {
	return g.El("textPath", children...)
}

func Title(children ...g.Node) g.Node {
	return g.El("title", children...)
}

func Tspan(children ...g.Node) g.Node {
	return g.El("tspan", children...)
}

func Use(children ...g.Node) g.Node {
	return g.El("use", children...)
}

func View(children ...g.Node) g.Node {
	return g.El("view", children...)
}
