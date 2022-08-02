package pieces

import "fmt"

// ColorForState returns the 'body' text wrapped in colour tags
// for the defined state
func ColorForState(state string, body string) string {
	left := "[white:]"
	right := "[white:]"
	str := body

	switch state {
	case "active":
		left = "[green:]"
		right = "[white:]"
	case "off":
		left = "[red:]"
		right = "[white:]"
	case "offline":
		left = "[red:]"
		right = "[white:]"
	case "online":
		left = "[green:]"
		right = "[white:]"
	default:
		left = "[white:]"
		right = "[white:]"
	}

	return left + str + right
}

/* -------------------- Color and Format Modifiers -------------------- */

func Bold(body string) string {
	return fmt.Sprintf("%s%s%s", "[::b]", body, "[::-]")
}

func Green(body string) string {
	return fmt.Sprintf("%s%s%s", "[green:]", body, "[white:]")
}
