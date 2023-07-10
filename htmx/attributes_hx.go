package htmx

import (
	"encoding/json"
	"fmt"
	"strings"

	g "github.com/maragudk/gomponents"
	_ "github.com/maragudk/gomponents/html"
)

const (
	DefaultHyperScriptVersion = `0.9.8`
)

func Boost(b ...string) g.Node {
	if len(b) == 1 {
		return HxAttr("hx-boost", b[0])
	}
	return HxAttr("hx-boost", "true")
}

//
// requests
//

func Get(path string) g.Node {
	return HxAttr("hx-get", path)
}

func Post(path string) g.Node {
	return HxAttr("hx-post", path)
}

func Delete(path string) g.Node {
	return HxAttr("hx-delete", path)
}

func Patch(path string) g.Node {
	return HxAttr("hx-patch", path)
}

func Put(path string) g.Node {
	return HxAttr("hx-put", path)
}

//
// core attributes
//

func On(v string) g.Node {
	return HxAttr("hx-on", v)
}

func PushUrl(url string) g.Node {
	return HxAttr("hx-push-url", url)
}

func Select(v string) g.Node {
	return HxAttr("hx-select", v)
}

func SelectOOB(v string) g.Node {
	return HxAttr("hx-select-oob", v)
}

func Swap(v ...string) g.Node {
	if len(v) > 0 {
		if isValidSwapValue(v[0]) {
			return HxAttr("hx-swap", v[0])
		}
	}
	return HxAttr("hx-swap", "innerHTML")
}

func SwapOOB(v ...string) g.Node {
	switch len(v) {
	case 1:
		if isValidSwapValue(v[0]) {
			return HxAttr("hx-swap-oob", v[0])
		}
		return HxAttr("hx-swap-oob", "true")
	case 2:
		if isValidSwapValue(v[0]) {
			return HxAttr("hx-swap-oob", fmt.Sprintf("%s,%s", v[0], v[1]))
		}
		return HxAttr("hx-swap-oob", "true")
	}
	return HxAttr("hx-swap-oob", "true")
}

func Target(v string) g.Node {
	return HxAttr("hx-target", v)
}

func Trigger(v string) g.Node {
	return HxAttr("hx-trigger", v)
}

//
// additional attributes
//

func Confirm(v string) g.Node {
	return HxAttr("hx-confirm", v)
}

func Disable() g.Node {
	return HxAttr("hx-disable", "true")
}

func Disinherit(v ...string) g.Node {
	if len(v) > 0 {
		if isValidSwapValue(v[0]) {
			return HxAttr("hx-disinherit", v[0])
		}
	}
	return HxAttr("hx-disinherit", "*")
}

func Encoding(v string) g.Node {
	return HxAttr("hx-encoding", v)
}

func Ext(name string) g.Node {
	return HxAttr("hx-ext", name)
}

func ExtIgnore(name string) g.Node {
	return HxAttr("hx-ext", "ignore:"+name)
}

func Headers(m map[string]string) g.Node {
	headers, _ := json.Marshal(m)
	return HxAttr("hx-headers", string(headers))
}

func NoHistory() g.Node {
	return HxAttr("hx-history", "false")
}

func HistoryELT() g.Node {
	return HxAttr("hx-history-elt")
}

func Include(v string) g.Node {
	return HxAttr("hx-include", v)
}

func Indicator(v string) g.Node {
	return HxAttr("hx-indicator", v)
}

func Params(v ...string) g.Node {
	if len(v) > 0 {
		return HxAttr("hx-params", strings.Join(v, ","))
	}
	return HxAttr("hx-params", "*")
}

func NoParams(v ...string) g.Node {
	if len(v) > 0 {
		return HxAttr("hx-params", strings.Join(v, ","))
	}
	return HxAttr("hx-params", "none")
}

func Preserve() g.Node {
	return HxAttr("hx-preserve")
}

func Prompt(v string) g.Node {
	return HxAttr("hx-prompt", v)
}

func ReplaceURL(v ...string) g.Node {
	if len(v) > 0 {
		return HxAttr("hx-replace-url", v[0])
	}
	return HxAttr("hx-replace-url", "true")
}

func Request(v string) g.Node {
	return HxAttr("hx-request", v)
}

func Sync(v string) g.Node {
	return HxAttr("hx-sync", v)
}

func Validate() g.Node {
	return HxAttr("hx-validate")
}

func Vals(js bool, m map[string]string) g.Node {
	var jsPfx string
	if js {
		jsPfx = "js:"
	}
	v, _ := json.Marshal(m)
	return HxAttr("hx-vals", fmt.Sprintf("%s%s", jsPfx, string(v)))
}

//
// neco
//

func isValidSwapValue(in string) bool {
	for _, v := range []string{"innerHTML", "outerHTML", "beforebegin", "afterbegin", "beforeend", "afterend", "delete", "none"} {
		if in == v {
			return true
		}
	}
	return false
}
