package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var elements = map[string]string{
	"a":          "A",
	"abbr":       "Abbr",
	"address":    "Address",
	"article":    "Article",
	"aside":      "Aside",
	"audio":      "Audio",
	"b":          "B",
	"blockquote": "BlockQuote",
	"body":       "Body",
	"button":     "Button",
	"canvas":     "Canvas",
	"caption":    "Caption",
	"cite":       "Cite",
	"code":       "Code",
	"colgroup":   "ColGroup",
	"data":       "DataEl",
	"datalist":   "DataList",
	"dd":         "Dd",
	"del":        "Del",
	"details":    "Details",
	"dfn":        "Dfn",
	"dialog":     "Dialog",
	"div":        "Div",
	"dl":         "Dl",
	"dt":         "Dt",
	"em":         "Em",
	"fieldset":   "FieldSet",
	"figcaption": "FigCaption",
	"figure":     "Figure",
	"footer":     "Footer",
	"form":       "Form",
	"h1":         "H1",
	"h2":         "H2",
	"h3":         "H3",
	"h4":         "H4",
	"h5":         "H5",
	"h6":         "H6",
	"head":       "Head",
	"header":     "Header",
	"hgroup":     "HGroup",
	"html":       "HTML",
	"i":          "I",
	"iframe":     "IFrame",
	"ins":        "Ins",
	"kbd":        "Kbd",
	"label":      "Label",
	"legend":     "Legend",
	"li":         "Li",
	"main":       "Main",
	"mark":       "Mark",
	"menu":       "Menu",
	"meter":      "Meter",
	"nav":        "Nav",
	"noscript":   "NoScript",
	"object":     "Object",
	"ol":         "Ol",
	"optgroup":   "OptGroup",
	"option":     "Option",
	"p":          "P",
	"picture":    "Picture",
	"pre":        "Pre",
	"progress":   "Progress",
	"q":          "Q",
	"s":          "S",
	"samp":       "Samp",
	"script":     "Script",
	"section":    "Section",
	"select":     "Select",
	"small":      "Small",
	"span":       "Span",
	"strong":     "Strong",
	"style":      "StyleEl",
	"sub":        "Sub",
	"summary":    "Summary",
	"sup":        "Sup",
	"svg":        "SVG",
	"table":      "Table",
	"tbody":      "TBody",
	"td":         "Td",
	"textarea":   "Textarea",
	"tfoot":      "TFoot",
	"th":         "Th",
	"thead":      "THead",
	"time":       "Time",
	"title":      "TitleEl",
	"tr":         "Tr",
	"u":          "U",
	"ul":         "Ul",
	"var":        "Var",
	"video":      "Video",
	"area":       "Area",
	"base":       "Base",
	"br":         "Br",
	"col":        "Col",
	"embed":      "Embed",
	"hr":         "Hr",
	"img":        "Img",
	"input":      "Input",
	"link":       "Link",
	"meta":       "Meta",
	"param":      "Param",
	"source":     "Source",
	"wbr":        "Wbr",
}

var boolAttrs = map[string]string{
	"async":       "Async",
	"autofocus":   "AutoFocus",
	"autoplay":    "AutoPlay",
	"checked":     "Checked",
	"controls":    "Controls",
	"defer":       "Defer",
	"disabled":    "Disabled",
	"loop":        "Loop",
	"multiple":    "Multiple",
	"muted":       "Muted",
	"playsinline": "PlaysInline",
	"readonly":    "ReadOnly",
	"required":    "Required",
	"selected":    "Selected",
}

var simpleAttr = map[string]string{
	"accept":       "Accept",
	"action":       "Action",
	"alt":          "Alt",
	"as":           "As",
	"autocomplete": "AutoComplete",
	"charset":      "Charset",
	"class":        "Class",
	"cols":         "Cols",
	"colspan":      "ColSpan",
	"content":      "Content",
	"crossorigin":  "CrossOrigin",
	"datetime":     "DateTime",
	"enctype":      "EncType",
	"dir":          "Dir",
	"for":          "For",
	"form":         "FormAttr",
	"height":       "Height",
	"href":         "Href",
	"id":           "ID",
	"integrity":    "Integrity",
	"label":        "LabelAttr",
	"lang":         "Lang",
	"list":         "List",
	"loading":      "Loading",
	"max":          "Max",
	"maxlength":    "MaxLength",
	"method":       "Method",
	"min":          "Min",
	"minlength":    "MinLength",
	"name":         "Name",
	"pattern":      "Pattern",
	"placeholder":  "Placeholder",
	"poster":       "Poster",
	"preload":      "Preload",
	"rel":          "Rel",
	"role":         "Role",
	"rows":         "Rows",
	"rowspan":      "RowSpan",
	"src":          "Src",
	"srcset":       "SrcSet",
	"step":         "Step",
	"style":        "Style",
	"tabindex":     "TabIndex",
	"target":       "Target",
	"title":        "Title",
	"type":         "Type",
	"value":        "Value",
	"width":        "Width",
}

// TODO(Amr Ojjeh): Integrate SVG and HTMX
// TODO(Amr Ojjeh): Write tests

func main() {
	p := flag.String("p", "ui", "sets the package for the output")
	flag.Parse()
	if fileName := flag.Arg(0); fileName != "" {
		file, err := os.Open(fileName)
		if err != nil {
			panic("could not find file: " + os.Args[1])
		}
		convert(file, os.Stdout, *p)
		return
	}
	convert(os.Stdin, os.Stdout, *p)
}

func convert(r io.Reader, w io.Writer, packageName string) {

	w.Write([]byte(fmt.Sprintf(`package %s

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func Page() g.Node {
	return `, packageName)))

	dec := xml.NewDecoder(r)
	dec.Strict = false
	dec.AutoClose = xml.HTMLAutoClose
	dec.Entity = xml.HTMLEntity
	tok, err := dec.Token()
	for !errors.Is(err, io.EOF) {
		if err != nil {
			panic(err)
		}
		switch t := tok.(type) {
		case xml.StartElement:
			element(w, dec, t, 2)
		case xml.CharData:
		case xml.Directive:
			if strings.ToLower(string(t)) == "doctype html" {
				doctype(w, dec)
			}
		case xml.Comment:
		default:
			panic(fmt.Sprintf("unexpected token: %T", tok))
		}
		tok, err = dec.Token()
	}

	w.Write([]byte("\n}\n"))
}

func element(w io.Writer, dec *xml.Decoder, startEl xml.StartElement, indent int) {
	newline(w, indent)
	if el, ok := elements[startEl.Name.Local]; ok {
		w.Write([]byte(el + "("))
	} else {
		w.Write([]byte(fmt.Sprintf(`g.El("%s", `, startEl.Name.Local)))
	}
	for _, a := range startEl.Attr {
		attribute(w, a)
		w.Write([]byte(", "))
	}
	hasNested := false
	tok, err := dec.Token()
	for !errors.Is(err, io.EOF) {
		if err != nil {
			panic(err)
		}
		switch t := tok.(type) {
		case xml.StartElement:
			element(w, dec, t, indent+1)
			hasNested = true
		case xml.CharData:
			text(w, string(t), indent+1)
			hasNested = true
		case xml.EndElement:
			if hasNested {
				newline(w, indent)
			}
			w.Write([]byte("),"))
			return
		default:
			panic(fmt.Sprintf("unexpected token: %T", tok))
		}
		tok, err = dec.Token()
	}
	panic("expected end tag")
}

func attribute(w io.Writer, a xml.Attr) {
	if attr, ok := boolAttrs[a.Name.Local]; ok {
		w.Write([]byte(attr + "()"))
		return
	}

	if attr, ok := simpleAttr[a.Name.Local]; ok {
		w.Write([]byte(fmt.Sprintf(`%s("%s")`, attr, a.Value)))
		return
	}

	if after, found := cutPrefix(a.Name.Local, "aria-"); found {
		w.Write([]byte(fmt.Sprintf(`Aria("%s", "%s")`, after, a.Value)))
		return
	}

	if after, found := cutPrefix(a.Name.Local, "data-"); found {
		w.Write([]byte(fmt.Sprintf(`Data("%s", "%s")`, after, a.Value)))
		return
	}

	w.Write([]byte(fmt.Sprintf(`g.Attr("%s", "%s")`, a.Name.Local, a.Value)))
}

func text(w io.Writer, text string, indent int) {
	if strings.TrimSpace(text) == "" {
		return
	}
	newline(w, indent)
	w.Write([]byte(fmt.Sprintf(`g.Text("%s"),`, text)))
}

func doctype(w io.Writer, dec *xml.Decoder) {
	w.Write([]byte("Doctype("))
	hasNested := false
	tok, err := dec.Token()
	for !errors.Is(err, io.EOF) {
		switch t := tok.(type) {
		case xml.StartElement:
			element(w, dec, t, 2)
			hasNested = true
		case xml.CharData:
			text(w, string(t), 2)
			hasNested = true
		case xml.Comment:
		default:
			panic(fmt.Sprintf("unexpected token: %T", tok))
		}
		tok, err = dec.Token()
	}
	if hasNested {
		newline(w, 1)
	}
	w.Write([]byte(")"))
}

func newline(w io.Writer, indent int) {
	w.Write([]byte("\n"))
	for x := 0; x < indent; x++ {
		w.Write([]byte("\t"))
	}
}

func cutPrefix(s, prefix string) (after string, found bool) {
	if strings.HasPrefix(s, prefix) {
		return s[len(prefix):], true
	}
	return s, false
}
