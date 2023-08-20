package is

import "testing"

func TestIs(t *testing.T) {
	t.Parallel()
	is := New(t)
	is.True(1 != 2)
}
