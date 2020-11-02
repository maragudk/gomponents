package assert

import (
	"strings"
	"testing"

	g "github.com/maragudk/gomponents"
)

// Equal checks for equality between the given expected string and the rendered Node string.
func Equal(t *testing.T, expected string, actual g.Node) {
	var b strings.Builder
	_ = actual.Render(&b)
	if expected != b.String() {
		t.Errorf("expected `%v` but got `%v`", expected, actual)
		t.FailNow()
	}
}

// Error checks for a non-nil error.
func Error(t *testing.T, err error) {
	if err == nil {
		t.FailNow()
	}
}
