// Command htmlgen generates the html package elements and attributes from MDN browser-compat-data.
//
// Usage:
//
//	go run ./cmd/htmlgen
//
// This requires the MDN browser-compat-data repository to be cloned locally.
// Run scripts/fetch-mdn-data.sh first to clone the data.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

var mdnDataPath = "tmp/browser-compat-data"

// Element represents an HTML element for code generation.
type Element struct {
	Tag        string
	FuncName   string
	UseGroup   bool
	Deprecated string // If set, this is a deprecated alias pointing to the given function
}

// Attribute represents an HTML attribute for code generation.
type Attribute struct {
	Key        string
	FuncName   string
	IsBoolean  bool
	IsVariadic bool
}

// elementsWithGroup are elements that should use g.Group(children) instead of children...
// These are typically inline/phrasing content elements. This is a gomponents-specific design
// choice that can't be derived from MDN data.
var elementsWithGroup = map[string]bool{
	"abbr": true, "b": true, "caption": true, "dd": true, "del": true,
	"dfn": true, "dt": true, "em": true, "figcaption": true,
	"h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
	"i": true, "ins": true, "kbd": true, "mark": true, "q": true,
	"s": true, "samp": true, "small": true, "strong": true, "sub": true,
	"sup": true, "time": true, "title": true, "u": true, "var": true, "video": true,
}

// booleanAttributes are attributes that don't take a value.
// MDN doesn't provide this information, so we need to maintain this list.
var booleanAttributes = map[string]bool{
	"allowfullscreen": true, "async": true, "autofocus": true, "autoplay": true,
	"checked": true, "controls": true, "default": true, "defer": true,
	"disabled": true, "formnovalidate": true, "inert": true, "ismap": true,
	"itemscope": true, "loop": true, "multiple": true, "muted": true,
	"nomodule": true, "novalidate": true, "open": true, "playsinline": true,
	"readonly": true, "required": true, "reversed": true, "selected": true,
}

// variadicAttributes are attributes that can optionally take a value.
var variadicAttributes = map[string]bool{
	"popover": true,
}

// nameConflicts defines how to handle element/attribute name conflicts.
// Key is the HTML name, value is [elementFuncName, attrFuncName].
// Empty string means skip generation for that type.
var nameConflicts = map[string][2]string{
	"abbr":    {"Abbr", "AbbrAttr"},
	"cite":    {"Cite", "CiteAttr"},
	"data":    {"DataEl", ""},        // Data is special (prefix helper)
	"form":    {"Form", "FormAttr"},
	"label":   {"Label", "LabelAttr"},
	"link":    {"Link", ""},          // link attr is deprecated (body link color)
	"map":     {"MapEl", ""},         // Map conflicts with gomponents.Map helper
	"text":    {"", ""},              // text attr deprecated; conflicts with gomponents.Text
	"slot":    {"SlotEl", "SlotAttr"},
	"span":    {"Span", "SpanAttr"},
	"style":   {"StyleEl", "Style"},
	"summary": {"Summary", "SummaryAttr"},
	"title":   {"TitleEl", "Title"},
}

// deprecatedElementAliases are deprecated element function names for backwards compatibility.
var deprecatedElementAliases = map[string]string{
	"CiteEl":  "Cite",
	"FormEl":  "Form",
	"LabelEl": "Label",
}

// extraElements are elements not in MDN HTML data but commonly used.
var extraElements = []Element{
	{Tag: "svg", FuncName: "SVG"},   // SVG is in a different namespace but commonly inlined
	{Tag: "param", FuncName: "Param"}, // Deprecated in HTML5, but kept for backwards compatibility
}

// extraAttributes are attributes not in MDN HTML data but commonly used.
var extraAttributes = []Attribute{
	{Key: "role", FuncName: "Role"}, // ARIA role attribute
	// Microdata attributes (https://html.spec.whatwg.org/multipage/microdata.html)
	{Key: "itemid", FuncName: "ItemID"},
	{Key: "itemprop", FuncName: "ItemProp"},
	{Key: "itemref", FuncName: "ItemRef"},
	{Key: "itemscope", FuncName: "ItemScope", IsBoolean: true},
	{Key: "itemtype", FuncName: "ItemType"},
}

// skipAttributes are attributes that shouldn't be auto-generated.
var skipAttributes = map[string]bool{
	// Deprecated body attributes (universally deprecated)
	"alink": true, "background": true, "bgcolor": true, "border": true,
	"bottommargin": true, "color": true, "leftmargin": true, "marginheight": true,
	"marginwidth": true, "rightmargin": true, "topmargin": true, "vlink": true,
	// Experimental (can add later if stabilized)
	"anchor": true, "exportparts": true,
}

// funcNameOverrides provides Go function names for cases where the generic algorithm doesn't work.
// This includes acronyms, compound words, and other special cases.
var funcNameOverrides = map[string]string{
	// Acronyms
	"html": "HTML", "svg": "SVG", "id": "ID",
	// Hyphenated (keep together)
	"http-equiv": "HTTPEquiv", "accept-charset": "AcceptCharset",
	// Compound words (camelCase)
	"tbody": "TBody", "thead": "THead", "tfoot": "TFoot",
	"colgroup": "ColGroup", "optgroup": "OptGroup", "fieldset": "FieldSet",
	"datalist": "DataList", "textarea": "Textarea", "blockquote": "BlockQuote",
	"noscript": "NoScript", "iframe": "IFrame", "hgroup": "HGroup",
	"figcaption": "FigCaption", "fencedframe": "FencedFrame",
	// Attributes with compound words
	"srcset": "SrcSet", "colspan": "ColSpan", "rowspan": "RowSpan",
	"maxlength": "MaxLength", "minlength": "MinLength", "datetime": "DateTime",
	"tabindex": "TabIndex", "accesskey": "AccessKey", "autocomplete": "AutoComplete",
	"autocapitalize": "AutoCapitalize", "autofocus": "AutoFocus", "autoplay": "AutoPlay",
	"crossorigin": "CrossOrigin", "enctype": "EncType", "formaction": "FormAction",
	"formenctype": "FormEncType", "formmethod": "FormMethod", "formnovalidate": "FormNoValidate",
	"formtarget": "FormTarget", "playsinline": "PlaysInline", "readonly": "ReadOnly",
	"referrerpolicy": "ReferrerPolicy", "spellcheck": "SpellCheck",
	"popovertarget": "PopoverTarget", "popovertargetaction": "PopoverTargetAction",
	"contenteditable": "ContentEditable", "enterkeyhint": "EnterKeyHint",
	"inputmode": "InputMode", "itemid": "ItemID", "itemprop": "ItemProp",
	"itemref": "ItemRef", "itemscope": "ItemScope", "itemtype": "ItemType",
	"novalidate": "NoValidate", "nomodule": "NoModule", "allowfullscreen": "AllowFullscreen",
	"hreflang": "HrefLang", "srcdoc": "SrcDoc", "usemap": "UseMap", "ismap": "IsMap",
	"virtualkeyboardpolicy": "VirtualKeyboardPolicy", "writingsuggestions": "WritingSuggestions",
	// Elements
	"selectedcontent": "SelectedContent",
	// More compound word attributes
	"allowpaymentrequest": "AllowPaymentRequest", "attributionsourceid": "AttributionSourceID",
	"attributionsrc": "AttributionSrc", "autocorrect": "AutoCorrect",
	"browsingtopics": "BrowsingTopics", "cellpadding": "CellPadding", "cellspacing": "CellSpacing",
	"charoff": "CharOff", "controlslist": "ControlsList",
	"disablepictureinpicture": "DisablePictureInPicture", "disableremoteplayback": "DisableRemotePlayback",
	"fetchpriority": "FetchPriority", "hreftranslate": "HrefTranslate",
	"imagesizes": "ImageSizes", "imagesrcset": "ImageSrcSet",
	"shadowrootclonable": "ShadowRootClonable", "shadowrootdelegatesfocus": "ShadowRootDelegatesFocus",
	"shadowrootmode": "ShadowRootMode", "shadowrootserializable": "ShadowRootSerializable",
	"srclang": "SrcLang", "webkitdirectory": "WebkitDirectory",
}

var htmlOutputDir string

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func findModuleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find go.mod in any parent directory")
		}
		dir = parent
	}
}

func run() error {
	moduleRoot, err := findModuleRoot()
	if err != nil {
		return fmt.Errorf("finding module root: %w", err)
	}

	mdnDataPath = filepath.Join(moduleRoot, "tmp", "browser-compat-data")
	htmlOutputDir = filepath.Join(moduleRoot, "html")

	if _, err := os.Stat(mdnDataPath); os.IsNotExist(err) {
		return fmt.Errorf("MDN data not found at %s. Run scripts/fetch-mdn-data.sh first", mdnDataPath)
	}

	elements, err := loadElements()
	if err != nil {
		return fmt.Errorf("loading elements: %w", err)
	}

	attrs, err := loadAttributes()
	if err != nil {
		return fmt.Errorf("loading attributes: %w", err)
	}

	if err := generateElementsFile(elements); err != nil {
		return fmt.Errorf("generating elements.go: %w", err)
	}

	if err := generateAttributesFile(attrs); err != nil {
		return fmt.Errorf("generating attributes.go: %w", err)
	}

	log.Printf("Generated %d elements and %d attributes", len(elements), len(attrs))
	return nil
}

// MDNCompat represents the __compat structure in MDN data.
type MDNCompat struct {
	Status struct {
		Deprecated   bool `json:"deprecated"`
		Experimental bool `json:"experimental"`
	} `json:"status"`
}

func loadElements() ([]Element, error) {
	elementsDir := filepath.Join(mdnDataPath, "html", "elements")
	entries, err := os.ReadDir(elementsDir)
	if err != nil {
		return nil, err
	}

	var elements []Element
	seen := make(map[string]bool)

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(elementsDir, entry.Name()))
		if err != nil {
			return nil, err
		}

		var doc map[string]json.RawMessage
		if err := json.Unmarshal(data, &doc); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", entry.Name(), err)
		}

		var htmlData map[string]json.RawMessage
		if err := json.Unmarshal(doc["html"], &htmlData); err != nil {
			continue
		}

		var elementsData map[string]json.RawMessage
		if err := json.Unmarshal(htmlData["elements"], &elementsData); err != nil {
			continue
		}

		for tag, elemJSON := range elementsData {
			if seen[tag] {
				continue
			}

			// Check deprecation status from MDN
			var elemData map[string]json.RawMessage
			if err := json.Unmarshal(elemJSON, &elemData); err != nil {
				continue
			}

			if compatJSON, ok := elemData["__compat"]; ok {
				var compat MDNCompat
				if err := json.Unmarshal(compatJSON, &compat); err == nil {
					if compat.Status.Deprecated {
						continue // Skip deprecated elements
					}
				}
			}

			seen[tag] = true

			funcName := toFuncName(tag)
			if conflict, ok := nameConflicts[tag]; ok {
				funcName = conflict[0]
			}

			elements = append(elements, Element{
				Tag:      tag,
				FuncName: funcName,
				UseGroup: elementsWithGroup[tag],
			})
		}
	}

	// Add extra elements
	for _, extra := range extraElements {
		if !seen[extra.Tag] {
			elements = append(elements, extra)
		}
	}

	// Add deprecated aliases
	for alias, target := range deprecatedElementAliases {
		elements = append(elements, Element{
			FuncName:   alias,
			Deprecated: target,
		})
	}

	// Sort by function name (deprecated at end)
	sort.Slice(elements, func(i, j int) bool {
		if elements[i].Deprecated != "" && elements[j].Deprecated == "" {
			return false
		}
		if elements[i].Deprecated == "" && elements[j].Deprecated != "" {
			return true
		}
		return elements[i].FuncName < elements[j].FuncName
	})

	return elements, nil
}

func loadAttributes() ([]Attribute, error) {
	attrs := make(map[string]Attribute)

	// Load global attributes
	globalPath := filepath.Join(mdnDataPath, "html", "global_attributes.json")
	data, err := os.ReadFile(globalPath)
	if err != nil {
		return nil, err
	}

	var doc map[string]json.RawMessage
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, err
	}

	var htmlData map[string]json.RawMessage
	if err := json.Unmarshal(doc["html"], &htmlData); err != nil {
		return nil, fmt.Errorf("missing html key in global_attributes.json")
	}

	var globalAttrs map[string]json.RawMessage
	if err := json.Unmarshal(htmlData["global_attributes"], &globalAttrs); err != nil {
		return nil, fmt.Errorf("missing global_attributes key")
	}

	for key := range globalAttrs {
		if key == "__compat" {
			continue
		}
		addAttribute(attrs, key)
	}

	// Load element-specific attributes
	elementsDir := filepath.Join(mdnDataPath, "html", "elements")
	entries, err := os.ReadDir(elementsDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(elementsDir, entry.Name()))
		if err != nil {
			return nil, err
		}

		var doc map[string]json.RawMessage
		if err := json.Unmarshal(data, &doc); err != nil {
			continue
		}

		var htmlData map[string]json.RawMessage
		if err := json.Unmarshal(doc["html"], &htmlData); err != nil {
			continue
		}

		var elementsData map[string]json.RawMessage
		if err := json.Unmarshal(htmlData["elements"], &elementsData); err != nil {
			continue
		}

		for _, elemJSON := range elementsData {
			var elemData map[string]json.RawMessage
			if err := json.Unmarshal(elemJSON, &elemData); err != nil {
				continue
			}
			for attrKey := range elemData {
				if attrKey == "__compat" {
					continue
				}
				addAttribute(attrs, attrKey)
			}
		}
	}

	// Add extra attributes
	for _, extra := range extraAttributes {
		if _, exists := attrs[extra.Key]; !exists {
			attrs[extra.Key] = extra
		}
	}

	// Convert to slice
	var result []Attribute
	for _, attr := range attrs {
		result = append(result, attr)
	}

	// Sort: boolean first, then string, then variadic
	sort.Slice(result, func(i, j int) bool {
		if result[i].IsVariadic && !result[j].IsVariadic {
			return false
		}
		if !result[i].IsVariadic && result[j].IsVariadic {
			return true
		}
		if result[i].IsBoolean && !result[j].IsBoolean {
			return true
		}
		if !result[i].IsBoolean && result[j].IsBoolean {
			return false
		}
		return result[i].FuncName < result[j].FuncName
	})

	return result, nil
}

func addAttribute(attrs map[string]Attribute, key string) {
	if skipAttributes[key] {
		return
	}

	// Skip hyphenated attributes except known ones
	if strings.Contains(key, "-") && key != "http-equiv" && key != "accept-charset" {
		return
	}

	// Skip MDN meta-features (contain underscores)
	if strings.Contains(key, "_") {
		return
	}

	funcName := toFuncName(key)

	// Handle name conflicts
	if conflict, ok := nameConflicts[key]; ok {
		if conflict[1] == "" {
			return // No attr function for this conflict
		}
		funcName = conflict[1]
	}

	attrs[key] = Attribute{
		Key:        key,
		FuncName:   funcName,
		IsBoolean:  booleanAttributes[key],
		IsVariadic: variadicAttributes[key],
	}
}

// toFuncName converts an HTML tag or attribute name to a Go function name.
// It first checks overrides, then applies a generic capitalize algorithm.
func toFuncName(name string) string {
	if override, ok := funcNameOverrides[name]; ok {
		return override
	}
	return capitalize(name)
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func generateElementsFile(elements []Element) error {
	var buf bytes.Buffer

	buf.WriteString(`// Code generated by cmd/htmlgen; DO NOT EDIT.

// Package html provides common HTML elements and attributes.
//
// See https://developer.mozilla.org/en-US/docs/Web/HTML/Element for a list of elements.
//
// See https://developer.mozilla.org/en-US/docs/Web/HTML/Attributes for a list of attributes.
package html

import (
	"io"

	g "maragu.dev/gomponents"
)

// Doctype returns a special kind of [g.Node] that prefixes its sibling with the string "<!doctype html>".
func Doctype(sibling g.Node) g.Node {
	return g.NodeFunc(func(w io.Writer) error {
		if _, err := w.Write([]byte("<!doctype html>")); err != nil {
			return err
		}
		return sibling.Render(w)
	})
}
`)

	for _, elem := range elements {
		buf.WriteString("\n")

		if elem.Deprecated != "" {
			fmt.Fprintf(&buf, "// Deprecated: Use [%s] instead.\n", elem.Deprecated)
			fmt.Fprintf(&buf, "func %s(children ...g.Node) g.Node {\n", elem.FuncName)
			fmt.Fprintf(&buf, "\treturn %s(children...)\n", elem.Deprecated)
			buf.WriteString("}\n")
		} else if elem.UseGroup {
			fmt.Fprintf(&buf, "func %s(children ...g.Node) g.Node {\n", elem.FuncName)
			fmt.Fprintf(&buf, "\treturn g.El(%q, g.Group(children))\n", elem.Tag)
			buf.WriteString("}\n")
		} else {
			fmt.Fprintf(&buf, "func %s(children ...g.Node) g.Node {\n", elem.FuncName)
			fmt.Fprintf(&buf, "\treturn g.El(%q, children...)\n", elem.Tag)
			buf.WriteString("}\n")
		}
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("formatting: %w\n%s", err, buf.String())
	}

	return os.WriteFile(filepath.Join(htmlOutputDir, "elements.go"), formatted, 0644)
}

func generateAttributesFile(attrs []Attribute) error {
	var buf bytes.Buffer

	buf.WriteString(`// Code generated by cmd/htmlgen; DO NOT EDIT.

package html

import (
	g "maragu.dev/gomponents"
)
`)

	for _, attr := range attrs {
		buf.WriteString("\n")

		if attr.IsVariadic {
			fmt.Fprintf(&buf, "func %s(v ...string) g.Node {\n", attr.FuncName)
			fmt.Fprintf(&buf, "\treturn g.Attr(%q, v...)\n", attr.Key)
			buf.WriteString("}\n")
		} else if attr.IsBoolean {
			fmt.Fprintf(&buf, "func %s() g.Node {\n", attr.FuncName)
			fmt.Fprintf(&buf, "\treturn g.Attr(%q)\n", attr.Key)
			buf.WriteString("}\n")
		} else {
			fmt.Fprintf(&buf, "func %s(v string) g.Node {\n", attr.FuncName)
			fmt.Fprintf(&buf, "\treturn g.Attr(%q, v)\n", attr.Key)
			buf.WriteString("}\n")
		}
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("formatting: %w\n%s", err, buf.String())
	}

	return os.WriteFile(filepath.Join(htmlOutputDir, "attributes.go"), formatted, 0644)
}
