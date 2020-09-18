package attr_test

import (
	"testing"

	"github.com/maragudk/gomponents/assert"
	"github.com/maragudk/gomponents/attr"
)

func TestID(t *testing.T) {
	t.Run("given a value, returns id=value", func(t *testing.T) {
		assert.Equal(t, ` id="hat"`, attr.ID("hat"))
	})
}

func TestClass(t *testing.T) {
	t.Run("given a value, returns class=value", func(t *testing.T) {
		assert.Equal(t, ` class="hat"`, attr.Class("hat"))
	})
}

func TestClasses(t *testing.T) {
	t.Run("given a map, returns sorted keys from the map with value true", func(t *testing.T) {
		assert.Equal(t, ` class="boheme-hat hat partyhat"`, attr.Classes(map[string]bool{
			"boheme-hat": true,
			"hat":        true,
			"partyhat":   true,
			"turtlehat":  false,
		}))
	})
}
