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
	Deprecated string // If set, this is a deprecated alias pointing to the given function
}

// elementsWithGroup are elements that should use g.Group(children) instead of children...
// These are typically inline/phrasing content elements.
var elementsWithGroup = map[string]bool{
	"abbr": true, "b": true, "caption": true, "dd": true, "del": true,
	"dfn": true, "dt": true, "em": true, "figcaption": true,
	"h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
	"i": true, "ins": true, "kbd": true, "mark": true, "q": true,
	"s": true, "samp": true, "small": true, "strong": true, "sub": true,
	"sup": true, "time": true, "title": true, "u": true, "var": true, "video": true,
}

// booleanAttributes are attributes that don't take a value.
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
var nameConflicts = map[string][2]string{
	"cite":  {"Cite", "CiteAttr"},
	"data":  {"DataEl", ""},     // Data is special (prefix helper)
	"form":  {"Form", "FormAttr"},
	"label": {"Label", "LabelAttr"},
	"map":   {"MapEl", ""},      // Map conflicts with gomponents.Map helper
	"slot":  {"SlotEl", "SlotAttr"},
	"style": {"StyleEl", "Style"},
	"title": {"TitleEl", "Title"},
}

// deprecatedElementAliases are deprecated element function names.
var deprecatedElementAliases = map[string]string{
	"CiteEl":  "Cite",
	"FormEl":  "Form",
	"LabelEl": "Label",
}

// deprecatedAttrAliases are deprecated attribute function names.
var deprecatedAttrAliases = map[string]string{
	"DataAttr":  "Data",
	"StyleAttr": "Style",
	"TitleAttr": "Title",
}

// skipElements are elements that shouldn't be generated (deprecated/obsolete).
var skipElements = map[string]bool{
	"acronym": true, "applet": true, "basefont": true, "bgsound": true,
	"big": true, "blink": true, "center": true, "content": true,
	"dir": true, "font": true, "frame": true, "frameset": true,
	"image": true, "keygen": true, "marquee": true, "menuitem": true,
	"nobr": true, "noembed": true, "noframes": true,
	"plaintext": true, "rb": true, "rtc": true, "shadow": true,
	"spacer": true, "strike": true, "tt": true, "xmp": true,
}

// extraElements are elements that should be generated but aren't in MDN HTML data.
// SVG is commonly used inline in HTML documents.
var extraElements = []Element{
	{Tag: "svg", FuncName: "SVG"},
}

// extraAttributes are attributes that should be generated but aren't in MDN data.
// Role is an ARIA attribute that's very commonly used.
var extraAttributes = []Attribute{
	{Key: "role", FuncName: "Role"},
}

// knownGlobalAttributes are the global attributes we want to generate.
// This is curated to match the current html package plus useful additions.
var knownGlobalAttributes = map[string]bool{
	"accesskey": true, "autocapitalize": true, "autofocus": true,
	"class": true, "contenteditable": true, "dir": true, "draggable": true,
	"enterkeyhint": true, "hidden": true, "id": true, "inert": true,
	"inputmode": true, "is": true, "itemid": true, "itemprop": true,
	"itemref": true, "itemscope": true, "itemtype": true, "lang": true,
	"nonce": true, "part": true, "popover": true, "role": true, "slot": true,
	"spellcheck": true, "style": true, "tabindex": true, "title": true,
	"translate": true, "virtualkeyboardpolicy": true, "writingsuggestions": true,
}

// knownElementAttributes are element-specific attributes we want to generate.
var knownElementAttributes = map[string]bool{
	// Form-related
	"accept": true, "accept-charset": true, "action": true, "autocomplete": true,
	"checked": true, "cols": true, "disabled": true, "enctype": true,
	"for": true, "form": true, "formaction": true, "formenctype": true,
	"formmethod": true, "formnovalidate": true, "formtarget": true,
	"list": true, "max": true, "maxlength": true, "method": true,
	"min": true, "minlength": true, "multiple": true, "name": true,
	"novalidate": true, "pattern": true, "placeholder": true, "readonly": true,
	"required": true, "rows": true, "selected": true, "size": true,
	"step": true, "type": true, "value": true, "wrap": true,

	// Link/navigation
	"download": true, "href": true, "hreflang": true, "ping": true,
	"referrerpolicy": true, "rel": true, "target": true,

	// Media
	"alt": true, "autoplay": true, "controls": true, "crossorigin": true,
	"height": true, "loop": true, "muted": true, "playsinline": true,
	"poster": true, "preload": true, "src": true, "srcset": true,
	"width": true,

	// iframe/embed
	"allow": true, "allowfullscreen": true, "loading": true,
	"sandbox": true, "srcdoc": true,

	// Table
	"colspan": true, "headers": true, "rowspan": true, "scope": true,

	// Meta/link
	"as": true, "charset": true, "content": true, "http-equiv": true,
	"integrity": true, "media": true, "sizes": true,

	// Script
	"async": true, "defer": true, "nomodule": true,

	// Other
	"cite": true, "data": true, "datetime": true, "default": true,
	"high": true, "label": true, "low": true, "open": true, "optimum": true,
	"reversed": true, "start": true, "usemap": true,

	// Popover
	"popovertarget": true, "popovertargetaction": true,
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
	// Find module root by looking for go.mod
	moduleRoot, err := findModuleRoot()
	if err != nil {
		return fmt.Errorf("finding module root: %w", err)
	}

	mdnDataPath = filepath.Join(moduleRoot, "tmp", "browser-compat-data")
	htmlOutputDir = filepath.Join(moduleRoot, "html")

	// Check if MDN data exists
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

		var doc map[string]any
		if err := json.Unmarshal(data, &doc); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", entry.Name(), err)
		}

		// Navigate to html.elements
		htmlData, ok := doc["html"].(map[string]any)
		if !ok {
			continue
		}
		elementsData, ok := htmlData["elements"].(map[string]any)
		if !ok {
			continue
		}

		for tag := range elementsData {
			if seen[tag] || skipElements[tag] {
				continue
			}
			seen[tag] = true

			elem := Element{
				Tag:      tag,
				FuncName: tagToFuncName(tag),
				UseGroup: elementsWithGroup[tag],
			}

			// Handle name conflicts
			if conflict, ok := nameConflicts[tag]; ok {
				elem.FuncName = conflict[0]
			}

			elements = append(elements, elem)
		}
	}

	// Add extra elements (like SVG which is commonly used but not in HTML spec)
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

	// Sort by function name
	sort.Slice(elements, func(i, j int) bool {
		// Put deprecated at the end
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

	var doc map[string]any
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, err
	}

	htmlData, ok := doc["html"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("missing html key in global_attributes.json")
	}
	globalAttrs, ok := htmlData["global_attributes"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("missing global_attributes key")
	}

	// Process global attributes
	for key := range globalAttrs {
		if key == "__compat" {
			continue
		}
		if knownGlobalAttributes[key] {
			addAttribute(attrs, key)
		}
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

		var doc map[string]any
		if err := json.Unmarshal(data, &doc); err != nil {
			continue
		}

		htmlData, ok := doc["html"].(map[string]any)
		if !ok {
			continue
		}
		elementsData, ok := htmlData["elements"].(map[string]any)
		if !ok {
			continue
		}

		// For each element, look at its direct children (which are attributes)
		for _, elemData := range elementsData {
			elemMap, ok := elemData.(map[string]any)
			if !ok {
				continue
			}
			for attrKey := range elemMap {
				if attrKey == "__compat" {
					continue
				}
				if knownElementAttributes[attrKey] {
					addAttribute(attrs, attrKey)
				}
			}
		}
	}

	// Add special prefix helper attributes
	attrs["aria"] = Attribute{Key: "aria", FuncName: "Aria"}
	attrs["data-custom"] = Attribute{Key: "data", FuncName: "Data"}

	// Add extra attributes (like role which is ARIA but commonly used)
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

	// Add deprecated aliases
	for alias, target := range deprecatedAttrAliases {
		result = append(result, Attribute{
			FuncName:   alias,
			Deprecated: target,
		})
	}

	// Sort: boolean first, then string, then variadic, then special, then deprecated
	sort.Slice(result, func(i, j int) bool {
		// Deprecated at end
		if result[i].Deprecated != "" && result[j].Deprecated == "" {
			return false
		}
		if result[i].Deprecated == "" && result[j].Deprecated != "" {
			return true
		}
		// Special (Aria, Data) near end
		iSpecial := result[i].FuncName == "Aria" || result[i].FuncName == "Data"
		jSpecial := result[j].FuncName == "Aria" || result[j].FuncName == "Data"
		if iSpecial && !jSpecial {
			return false
		}
		if !iSpecial && jSpecial {
			return true
		}
		// Variadic before special
		if result[i].IsVariadic && !result[j].IsVariadic {
			return false
		}
		if !result[i].IsVariadic && result[j].IsVariadic {
			return true
		}
		// Boolean before string
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
	// Skip if contains hyphen (need special handling)
	if strings.Contains(key, "-") && key != "http-equiv" && key != "accept-charset" {
		return
	}

	funcName := attrToFuncName(key)

	// Handle name conflicts - skip if it's meant to be an element
	if conflict, ok := nameConflicts[key]; ok {
		if conflict[1] == "" {
			return // No attr function for this conflict
		}
		funcName = conflict[1]
	}

	attr := Attribute{
		Key:        key,
		FuncName:   funcName,
		IsBoolean:  booleanAttributes[key],
		IsVariadic: variadicAttributes[key],
	}

	attrs[key] = attr
}

func tagToFuncName(tag string) string {
	switch tag {
	case "a":
		return "A"
	case "b":
		return "B"
	case "i":
		return "I"
	case "p":
		return "P"
	case "q":
		return "Q"
	case "s":
		return "S"
	case "u":
		return "U"
	case "br":
		return "Br"
	case "dd":
		return "Dd"
	case "dl":
		return "Dl"
	case "dt":
		return "Dt"
	case "em":
		return "Em"
	case "h1", "h2", "h3", "h4", "h5", "h6":
		return strings.ToUpper(tag[:1]) + tag[1:]
	case "hr":
		return "Hr"
	case "li":
		return "Li"
	case "ol":
		return "Ol"
	case "rp":
		return "Rp"
	case "rt":
		return "Rt"
	case "td":
		return "Td"
	case "th":
		return "Th"
	case "tr":
		return "Tr"
	case "ul":
		return "Ul"
	case "html":
		return "HTML"
	case "svg":
		return "SVG"
	case "bdi":
		return "Bdi"
	case "bdo":
		return "Bdo"
	case "del":
		return "Del"
	case "dfn":
		return "Dfn"
	case "img":
		return "Img"
	case "ins":
		return "Ins"
	case "kbd":
		return "Kbd"
	case "map":
		return "Map"
	case "nav":
		return "Nav"
	case "pre":
		return "Pre"
	case "sub":
		return "Sub"
	case "sup":
		return "Sup"
	case "var":
		return "Var"
	case "wbr":
		return "Wbr"
	case "abbr":
		return "Abbr"
	case "area":
		return "Area"
	case "base":
		return "Base"
	case "body":
		return "Body"
	case "code":
		return "Code"
	case "head":
		return "Head"
	case "link":
		return "Link"
	case "main":
		return "Main"
	case "mark":
		return "Mark"
	case "menu":
		return "Menu"
	case "meta":
		return "Meta"
	case "ruby":
		return "Ruby"
	case "samp":
		return "Samp"
	case "span":
		return "Span"
	case "time":
		return "Time"
	case "audio":
		return "Audio"
	case "aside":
		return "Aside"
	case "embed":
		return "Embed"
	case "input":
		return "Input"
	case "meter":
		return "Meter"
	case "small":
		return "Small"
	case "table":
		return "Table"
	case "tbody":
		return "TBody"
	case "tfoot":
		return "TFoot"
	case "thead":
		return "THead"
	case "track":
		return "Track"
	case "video":
		return "Video"
	case "button":
		return "Button"
	case "canvas":
		return "Canvas"
	case "dialog":
		return "Dialog"
	case "figure":
		return "Figure"
	case "footer":
		return "Footer"
	case "header":
		return "Header"
	case "hgroup":
		return "HGroup"
	case "iframe":
		return "IFrame"
	case "legend":
		return "Legend"
	case "object":
		return "Object"
	case "option":
		return "Option"
	case "output":
		return "Output"
	case "script":
		return "Script"
	case "search":
		return "Search"
	case "select":
		return "Select"
	case "source":
		return "Source"
	case "strong":
		return "Strong"
	case "address":
		return "Address"
	case "article":
		return "Article"
	case "caption":
		return "Caption"
	case "col":
		return "Col"
	case "colgroup":
		return "ColGroup"
	case "datalist":
		return "DataList"
	case "details":
		return "Details"
	case "div":
		return "Div"
	case "fieldset":
		return "FieldSet"
	case "figcaption":
		return "FigCaption"
	case "noscript":
		return "NoScript"
	case "optgroup":
		return "OptGroup"
	case "picture":
		return "Picture"
	case "progress":
		return "Progress"
	case "section":
		return "Section"
	case "summary":
		return "Summary"
	case "template":
		return "Template"
	case "textarea":
		return "Textarea"
	case "blockquote":
		return "BlockQuote"
	case "fencedframe":
		return "FencedFrame"
	}

	// Default: capitalize first letter
	return capitalize(tag)
}

func attrToFuncName(key string) string {
	switch key {
	case "accept-charset":
		return "AcceptCharset"
	case "http-equiv":
		return "HTTPEquiv"
	case "id":
		return "ID"
	case "for":
		return "For"
	case "class":
		return "Class"
	case "href":
		return "Href"
	case "src":
		return "Src"
	case "alt":
		return "Alt"
	case "type":
		return "Type"
	case "name":
		return "Name"
	case "value":
		return "Value"
	case "rel":
		return "Rel"
	case "as":
		return "As"
	case "dir":
		return "Dir"
	case "max":
		return "Max"
	case "min":
		return "Min"
	case "low":
		return "Low"
	case "high":
		return "High"
	case "is":
		return "Is"
	case "srcset":
		return "SrcSet"
	case "colspan":
		return "ColSpan"
	case "rowspan":
		return "RowSpan"
	case "maxlength":
		return "MaxLength"
	case "minlength":
		return "MinLength"
	case "datetime":
		return "DateTime"
	case "tabindex":
		return "TabIndex"
	case "accesskey":
		return "AccessKey"
	case "autocomplete":
		return "AutoComplete"
	case "autocapitalize":
		return "AutoCapitalize"
	case "autofocus":
		return "AutoFocus"
	case "autoplay":
		return "AutoPlay"
	case "crossorigin":
		return "CrossOrigin"
	case "enctype":
		return "EncType"
	case "formaction":
		return "FormAction"
	case "formenctype":
		return "FormEncType"
	case "formmethod":
		return "FormMethod"
	case "formnovalidate":
		return "FormNoValidate"
	case "formtarget":
		return "FormTarget"
	case "playsinline":
		return "PlaysInline"
	case "readonly":
		return "ReadOnly"
	case "referrerpolicy":
		return "ReferrerPolicy"
	case "spellcheck":
		return "SpellCheck"
	case "popovertarget":
		return "PopoverTarget"
	case "popovertargetaction":
		return "PopoverTargetAction"
	case "contenteditable":
		return "ContentEditable"
	case "enterkeyhint":
		return "EnterKeyHint"
	case "inputmode":
		return "InputMode"
	case "itemid":
		return "ItemID"
	case "itemprop":
		return "ItemProp"
	case "itemref":
		return "ItemRef"
	case "itemscope":
		return "ItemScope"
	case "itemtype":
		return "ItemType"
	case "novalidate":
		return "NoValidate"
	case "nomodule":
		return "NoModule"
	case "allowfullscreen":
		return "AllowFullscreen"
	case "hreflang":
		return "HrefLang"
	case "srcdoc":
		return "SrcDoc"
	case "usemap":
		return "UseMap"
	case "ismap":
		return "IsMap"
	case "virtualkeyboardpolicy":
		return "VirtualKeyboardPolicy"
	case "writingsuggestions":
		return "WritingSuggestions"
	}

	// Default: capitalize first letter
	return capitalize(key)
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

		if attr.Deprecated != "" {
			fmt.Fprintf(&buf, "// Deprecated: Use [%s] instead.\n", attr.Deprecated)
			if attr.FuncName == "DataAttr" {
				fmt.Fprintf(&buf, "func %s(name, v string) g.Node {\n", attr.FuncName)
				fmt.Fprintf(&buf, "\treturn %s(name, v)\n", attr.Deprecated)
			} else {
				fmt.Fprintf(&buf, "func %s(v string) g.Node {\n", attr.FuncName)
				fmt.Fprintf(&buf, "\treturn %s(v)\n", attr.Deprecated)
			}
			buf.WriteString("}\n")
			continue
		}

		if attr.FuncName == "Aria" {
			buf.WriteString("// Aria attributes automatically have their name prefixed with \"aria-\".\n")
			fmt.Fprintf(&buf, "func %s(name, v string) g.Node {\n", attr.FuncName)
			buf.WriteString("\treturn g.Attr(\"aria-\"+name, v)\n")
			buf.WriteString("}\n")
			continue
		}

		if attr.FuncName == "Data" {
			buf.WriteString("// Data attributes automatically have their name prefixed with \"data-\".\n")
			fmt.Fprintf(&buf, "func %s(name, v string) g.Node {\n", attr.FuncName)
			buf.WriteString("\treturn g.Attr(\"data-\"+name, v)\n")
			buf.WriteString("}\n")
			continue
		}

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
