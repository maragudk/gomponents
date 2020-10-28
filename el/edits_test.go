package el_test

import (
	"testing"

	"github.com/maragudk/gomponents/assert"
	"github.com/maragudk/gomponents/el"
)

func TestDel(t *testing.T) {
	t.Run("returns a del element", func(t *testing.T) {
		assert.Equal(t, `<del>hat</del>`, el.Del("hat"))
	})
}

func TestIns(t *testing.T) {
	t.Run("returns an ins element", func(t *testing.T) {
		assert.Equal(t, `<ins>hat</ins>`, el.Ins("hat"))
	})
}
