package htmx_test

import (
	"fmt"
	"testing"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/htmx"
	"github.com/maragudk/gomponents/internal/assert"
)

func TestSimpleAttribute(t *testing.T) {
	cases := map[string]func() g.Node{
		"hx-history-elt": HistoryELT,
		"hx-preserve":    Preserve,
		"hx-validate":    Validate,
	}

	for name, fn := range cases {
		t.Run(fmt.Sprintf("%v", name), func(t *testing.T) {
			n := g.El("div", fn())
			assert.Equal(t, fmt.Sprintf(`<div %v></div>`, name), n)
		})
	}
}

func TestStringAttribute(t *testing.T) {
	cases := map[string]func(string) g.Node{
		"_":             Hyper,
		"hx-confirm":    Confirm,
		"hx-delete":     Delete,
		"hx-encoding":   Encoding,
		"hx-ext":        Ext,
		"hx-get":        Get,
		"hx-include":    Include,
		"hx-indicator":  Indicator,
		"hx-on":         On,
		"hx-patch":      Patch,
		"hx-post":       Post,
		"hx-prompt":     Prompt,
		"hx-push-url":   PushUrl,
		"hx-put":        Put,
		"hx-request":    Request,
		"hx-select-oob": SelectOOB,
		"hx-select":     Select,
		"hx-sync":       Sync,
		"hx-target":     Target,
		"hx-trigger":    Trigger,
	}

	for name, fn := range cases {
		t.Run(fmt.Sprintf("%s", name), func(t *testing.T) {
			n := g.El("div", fn("hat"))
			assert.Equal(t, fmt.Sprintf(`<div %v='hat'></div>`, name), n)
		})
	}
}

//
// Helper functions that abstract some of the boilerplate
//

func TestDisable(t *testing.T) {
	t.Run("hx-disable", func(t *testing.T) {
		n := g.El("div", Disable())
		assert.Equal(t, `<div hx-disable='true'></div>`, n)
	})
}

func TestExtIgnoreAttribute(t *testing.T) {
	t.Run("hx-ext/ignore", func(t *testing.T) {
		n := g.El("div", ExtIgnore("hat"))
		assert.Equal(t, `<div hx-ext='ignore:hat'></div>`, n)
	})
}

func TestNoHistory(t *testing.T) {
	t.Run("hx-history/false", func(t *testing.T) {
		n := g.El("div", NoHistory())
		assert.Equal(t, `<div hx-history='false'></div>`, n)
	})
}

func TestHeadersAttribute(t *testing.T) {
	t.Run("hx-headers", func(t *testing.T) {
		n := g.El("div", Headers(map[string]string{"foo": "bar"}))
		assert.Equal(t, `<div hx-headers='{"foo":"bar"}'></div>`, n)
	})
}

//
// SSE attributes
//

func TestSSEConnectAttribute(t *testing.T) {
	t.Run("sse-connect", func(t *testing.T) {
		n := g.El("div", SSEConnect("blah"))
		assert.Equal(t, `<div sse-connect='blah'></div>`, n)
	})
}

func TestSSESwapAttribute(t *testing.T) {
	t.Run("sse-swap", func(t *testing.T) {
		n := g.El("div", SSESwap("blah"))
		assert.Equal(t, `<div sse-swap='blah'></div>`, n)
	})
}

//
// dynamic param attributes
//

func TestBoost(t *testing.T) {
	t.Run("hx-boost", func(t *testing.T) {
		n := g.El("div", Boost())
		assert.Equal(t, `<div hx-boost='true'></div>`, n)
	})
}

func TestDisinheritAttribute(t *testing.T) {
	t.Run("hx-disinherit", func(t *testing.T) {
		n := g.El("div", Disinherit())
		assert.Equal(t, `<div hx-disinherit='*'></div>`, n)
	})
}

func TestNoParamsAttribute(t *testing.T) {
	t.Run("hx-params/none", func(t *testing.T) {
		n := g.El("div", NoParams())
		assert.Equal(t, `<div hx-params='none'></div>`, n)
	})
}

func TestParamsAttribute(t *testing.T) {
	t.Run("hx-params/multiple", func(t *testing.T) {
		n := g.El("div", Params("param1", "param2"))
		assert.Equal(t, `<div hx-params='param1,param2'></div>`, n)
	})
}

func TestReplaceURLAttribute(t *testing.T) {
	t.Run("hx-replace-url", func(t *testing.T) {
		n := g.El("div", ReplaceURL())
		assert.Equal(t, `<div hx-replace-url='true'></div>`, n)
	})
}

func TestSwapAttribute(t *testing.T) {
	t.Run("hx-swap", func(t *testing.T) {
		n := g.El("div", Swap())
		assert.Equal(t, `<div hx-swap='innerHTML'></div>`, n)
	})
}

func TestSwapOOBAttribute(t *testing.T) {
	t.Run("hx-swap-oob", func(t *testing.T) {
		n := g.El("div", SwapOOB())
		assert.Equal(t, `<div hx-swap-oob='true'></div>`, n)
	})
}

func TestValsAttribute(t *testing.T) {
	t.Run("hx-vals/js", func(t *testing.T) {
		n := g.El("div", Vals(true, map[string]string{"foo": "bar"}))
		assert.Equal(t, `<div hx-vals='js:{"foo":"bar"}'></div>`, n)
	})
}
