package msg

import (
	"github.com/fatih/color"
)

var (
	// Bold represents bold font
	Bold = color.New(color.Bold)
	// Green represents green font
	Green = color.New(color.FgGreen)
	// GreenBold represents green + bold font
	GreenBold = color.New(color.FgGreen, color.Bold)
	// Red represents red font
	Red = color.New(color.FgRed)
	// RedBold represents red + bold font
	RedBold = color.New(color.FgRed, color.Bold)
	// Yellow represents yellow font
	Yellow = color.New(color.FgYellow)
)

// DisableColor disables colorized output
func DisableColor() {
	color.NoColor = true
}
