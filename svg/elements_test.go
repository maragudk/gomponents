package svg_test

import (
	"fmt"
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/internal/assert"
	. "github.com/maragudk/gomponents/svg"
)

func TestSimpleElements(t *testing.T) {
	cases := map[string]func(...g.Node) g.Node{
		"path":                Path,
		"a":                   A,
		"animate":             Animate,
		"animateMotion":       AnimateMotion,
		"animateTransform":    AnimateTransform,
		"circle":              Circle,
		"clipPath":            ClipPathEl,
		"defs":                Defs,
		"desc":                Desc,
		"ellipse":             Ellipse,
		"feBlend":             FeBlend,
		"feColorMatrix":       FeColorMatrix,
		"feComponentTransfer": FeComponentTransfer,
		"feComposite":         FeComposite,
		"feConvolveMatrix":    FeConvolveMatrix,
		"feDiffuseLighting":   FeDiffuseLighting,
		"feDisplacementMap":   FeDisplacementMap,
		"feDistantLight":      FeDistantLight,
		"feDropShadow":        FeDropShadow,
		"feFlood":             FeFlood,
		"feFuncA":             FeFuncA,
		"feFuncB":             FeFuncB,
		"feFuncG":             FeFuncG,
		"feFuncR":             FeFuncR,
		"feGaussianBlur":      FeGaussianBlur,
		"feImage":             FeImage,
		"feMerge":             FeMerge,
		"feMergeNode":         FeMergeNode,
		"feMorphology":        FeMorphology,
		"feOffset":            FeOffset,
		"fePointLight":        FePointLight,
		"feSpecularLighting":  FeSpecularLighting,
		"feSpotLight":         FeSpotLight,
		"feTile":              FeTile,
		"feTurbulence":        FeTurbulence,
		"filter":              FilterEl,
		"foreignObject":       ForeignObject,
		"g":                   G,
		"image":               Image,
		"line":                Line,
		"linearGradient":      LinearGradient,
		"marker":              Marker,
		"mask":                MaskEl,
		"metadata":            Metadata,
		"mpath":               Mpath,
		"pattern":             Pattern,
		"polygon":             Polygon,
		"polyline":            Polyline,
		"radialGradient":      RadialGradient,
		"rect":                Rect,
		"script":              Script,
		"set":                 Set,
		"stop":                Stop,
		"style":               StyleEl,
		"switch":              Switch,
		"symbol":              Symbol,
		"text":                Text,
		"textPath":            TextPath,
		"title":               Title,
		"tspan":               Tspan,
		"use":                 Use,
		"view":                View,
	}

	for name, fn := range cases {
		t.Run(fmt.Sprintf("should output %v", name), func(t *testing.T) {
			n := fn(g.Attr("id", "hat"))
			assert.Equal(t, fmt.Sprintf(`<%v id="hat"></%v>`, name, name), n)
		})
	}
}

func TestSVG(t *testing.T) {
	t.Run("outputs svg element with xml namespace attribute", func(t *testing.T) {
		assert.Equal(t, `<svg xmlns="http://www.w3.org/2000/svg"><path></path></svg>`, SVG(g.El("path")))
	})
}
