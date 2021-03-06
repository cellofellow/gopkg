package colorprint

import (
	"fmt"
	"strings"
)

func Printf(colorCode, format string, a ...interface{}) {
	fmt.Printf(colourText(colorCode, format), a...)
}

func Print(colorCode string, a ...interface{}) {
	fmt.Print(colourText(colorCode, fmt.Sprint(a...)))
}

func Println(colorCode string, a ...interface{}) {
	Println(fmt.Sprint(a...))
}

// ColorText returns the text colored according to the color code
// (exported version of colourText).
func ColorText(colorCode, text string) string {
	return colourText(colorCode, text)
}

// colourText return a colours text based on the given colour code.
func colourText(colorCode, text string) string {
	var coloredText string

	if len(colorCode) < 1 {
		return text // nothing to do
	}

	text = strings.Replace(text, "%", "%%", -1)

	switch colorCode[0] {
	// foreground
	case 'r':
		coloredText = Red(text).String()
	case 'g':
		coloredText = Green(text).String()
	case 'y':
		coloredText = Yellow(text).String()
	case 'b':
		coloredText = Blue(text).String()
	case 'x':
		coloredText = Black(text).String()
	case 'm':
		coloredText = Magenta(text).String()
	case 'c':
		coloredText = Cyan(text).String()
	case 'w':
		coloredText = White(text).String()
	case 'd':
		coloredText = Default(text).String()

	// background
	case 'R':
		coloredText = BgRed(text).String()
	case 'G':
		coloredText = BgGreen(text).String()
	case 'Y':
		coloredText = BgYellow(text).String()
	case 'B':
		coloredText = BgBlue(text).String()

	// case 'X': -> not implemented
	case 'M':
		coloredText = BgMagenta(text).String()
	case 'C':
		coloredText = BgCyan(text).String()
	case 'W':
		coloredText = BgWhite(text).String()
	case 'D':
		coloredText = BgDefault(text).String()
	// specials
	case '+':
		coloredText = Bold(text).String()
	case '*':
		coloredText = Italic(text).String()
	case '~':
		coloredText = Reverse(text).String()
	case '_':
		coloredText = Underline(text).String()
	case '#':
		coloredText = Blink(text).String()
	case '?':
		coloredText = Concealed(text).String()

	// default -> panic
	default:
		panic("colorprint: unknown color code")
	}

	coloredText = strings.Replace(coloredText, "%%", "%", -1)
	return colourText(colorCode[1:], coloredText)
}
