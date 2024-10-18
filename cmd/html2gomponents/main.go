package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io"
	"log/slog"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var attrs = map[string]string{
	"autocomplete": "AutoComplete",
	"autofocus":    "AutoFocus",
	"autoplay":     "AutoPlay",
	"cite":         "CiteAttr",
	"colspan":      "ColSpan",
	"crossorigin":  "CrossOrigin",
	"datetime":     "DateTime",
	"enctype":      "EncType",
	"form":         "FormAttr",
	"id":           "ID",
	"label":        "LabelAttr",
	"maxlength":    "MaxLength",
	"minlength":    "MinLength",
	"playsinline":  "PlaysInline",
	"readonly":     "ReadOnly",
	"rowspan":      "RowSpan",
	"srcset":       "SrcSet",
	"tabindex":     "TabIndex",
}

var els = map[string]string{
	"blockquote": "BlockQuote",
	"colgroup":   "ColGroup",
	"data":       "DataEl",
	"datalist":   "DataList",
	"fieldset":   "FieldSet",
	"figcaption": "FigCaption",
	"hgroup":     "HGroup",
	"html":       "HTML",
	"iframe":     "IFrame",
	"noscript":   "NoScript",
	"optgroup":   "OptGroup",
	"style":      "StyleEl",
	"svg":        "SVG",
	"tbody":      "TBody",
	"tfoot":      "TFoot",
	"thead":      "THead",
	"title":      "TitleEl",
}

func main() {
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	if err := start(os.Stdin, os.Stdout); err != nil {
		log.Info("Error", "error", err)
		os.Exit(1)
	}
}

func start(r io.Reader, w2 io.Writer) error {
	var b bytes.Buffer
	w := &statefulWriter{w: &b}

	w.Write("package html\n")
	w.Write("\n")
	w.Write("import (\n")
	w.Write(". \"maragu.dev/gomponents\"\n")
	w.Write(". \"maragu.dev/gomponents/html\"\n")
	w.Write(")\n")
	w.Write("\n")
	w.Write("func Component() Node {\n")

	z := html.NewTokenizer(r)

	var hasContent bool
	var depth int
loop:
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			if err := z.Err(); err != nil {
				if errors.Is(err, io.EOF) {
					if !hasContent {
						w.Write("return nil")
					}
					break loop
				}
				return err
			}

		case html.TextToken:
			text := strings.TrimSpace(string(z.Text()))
			if text == "" {
				continue
			}

			if !hasContent {
				w.Write("return ")
			}

			hasContent = true
			w.Write(fmt.Sprintf("Text(%q)", text))
			if depth > 0 {
				w.Write(",")
			}

		case html.StartTagToken, html.SelfClosingTagToken:
			if !hasContent {
				w.Write("return ")
			}

			if hasContent {
				w.Write("\n")
			}
			hasContent = true

			name, hasAttr := z.TagName()
			if el, ok := els[string(name)]; ok {
				w.Write(el)
			} else {
				w.Write(strings.ToTitle(string(name[0])))
				w.Write(string(name[1:]))
			}
			w.Write("(")

			if hasAttr {
				for {
					key, val, moreAttr := z.TagAttr()

					if attr, ok := attrs[string(key)]; ok {
						w.Write(attr)
					} else {
						w.Write(strings.ToTitle(string(key[0])))
						w.Write(string(key[1:]))
					}
					w.Write("(")

					if len(val) > 0 {
						w.Write(`"` + string(val) + `"`)
					}

					w.Write("),")

					if !moreAttr {
						break
					}
				}
				w.Write("\n")
			}
			depth++

			if tt == html.SelfClosingTagToken {
				depth--
				w.Write("\n)")
				if depth > 0 {
					w.Write(",")
				}
			}

		case html.EndTagToken:
			depth--
			w.Write("\n)")
			if depth > 0 {
				w.Write(",")
			}

		case html.CommentToken:
			w.Write("// " + string(z.Text()) + "\n")

		case html.DoctypeToken:
			// TODO Ignore for now
		}
	}

	w.Write("\n}\n")

	if w.err != nil {
		return w.err
	}

	formatted, err := format.Source(b.Bytes())
	if err != nil {
		return fmt.Errorf("error formatting output: %w", err)
	}

	if _, err = w2.Write(formatted); err != nil {
		return err
	}

	return nil
}

// statefulWriter only writes if no errors have occurred earlier in its lifetime.
type statefulWriter struct {
	w   io.Writer
	err error
}

func (w *statefulWriter) Write(s string) {
	if w.err != nil {
		return
	}
	_, w.err = w.w.Write([]byte(s))
}
