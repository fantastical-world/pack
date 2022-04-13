package pack

import "testing"

func TestPackError_Error(t *testing.T) {
	t.Run("validate that error message is correct...", func(t *testing.T) {
		got := Error("this is what i want")
		if got.Error() != "this is what i want" {
			t.Errorf("want this is what i want, got %s", got)
		}
	})
}
