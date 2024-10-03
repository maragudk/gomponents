// Package assert provides testing helpers.
package assert

import (
	"strings"
	"testing"

	g "maragu.dev/gomponents"
)

// Equal checks for equality between the given expected string and the rendered Node string.
func Equal(t *testing.T, expected string, actual g.Node) {
	t.Helper()

	var b strings.Builder
	_ = actual.Render(&b)
	if expected != b.String() {
		t.Fatalf(`expected "%v" but got "%v"`, expected, b.String())
	}
}

// Error checks for a non-nil error.
func Error(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatal("error is nil")
	}
}
