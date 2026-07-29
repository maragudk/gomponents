package gomponents

import "testing"

func TestRaw(t *testing.T) {
	t.Run("r.String() == string(r)", func(t *testing.T) {
		r := raw("<p>raw</p>")
		if r.String() != string(r) {
			t.Fail()
		}
	})
}
