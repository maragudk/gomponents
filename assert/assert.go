package assert

import (
	"testing"

	g "github.com/maragudk/gomponents"
)

// Equal checks for equality between the given expected string and the rendered Node string.
func Equal(t *testing.T, expected string, actual g.Node) {
	if expected != actual.Render() {
		t.Errorf("expected `%v` but got `%v`", expected, actual)
		t.FailNow()
	}
}
