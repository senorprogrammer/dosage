package pieces

import "fmt"

// ColorForState returns the 'body' text wrapped in colour tags
// for the defined state
func ColorForState(state string, body string) string {
	str := body

	switch state {
	case "active":
		str = Green(str)
	case "off":
		str = Red(str)
	case "offline":
		str = Red(str)
	case "online":
		str = Green(str)
	default:
		str = White(str)
	}

	return str
}

/* -------------------- Color and Format Modifiers -------------------- */

func Bold(body string) string {
	return fmt.Sprintf("%s%s%s", "[::b]", body, "[::-]")
}

func Green(body string) string {
	return fmt.Sprintf("%s%s%s", "[green:]", body, "[-:]")
}

func Red(body string) string {
	return fmt.Sprintf("%s%s%s", "[red:]", body, "[-:]")
}

func White(body string) string {
	return fmt.Sprintf("%s%s%s", "[white:]", body, "[-:]")
}
