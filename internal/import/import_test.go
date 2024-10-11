package import_test

import (
	"testing"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
	. "maragu.dev/gomponents/http"
)

func TestImports(t *testing.T) {
	t.Run("this is just a test that does nothing, but I need the dot imports above", func(t *testing.T) {
		_ = El("div")
		_ = A()
		_ = HTML5(HTML5Props{})
		_ = Adapt(nil)
	})
}
