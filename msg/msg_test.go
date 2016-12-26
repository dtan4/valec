package msg

import (
	"testing"

	"github.com/fatih/color"
)

func TestDisableColor(t *testing.T) {
	color.NoColor = false

	DisableColor()

	if !color.NoColor {
		t.Error("color.Nocolor should be true.")
	}
}
