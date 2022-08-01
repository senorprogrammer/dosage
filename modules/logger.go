package modules

import "github.com/senorprogrammer/dosage/pieces"

// Logger stores and displays log messages emitted by other parts of the system.
type Logger struct {
	Base
	PositionData pieces.PositionData
	Messages     []string
}

// NewLogger creates and returns an instance of Logger
func NewLogger(title string) *Logger {
	mod := &Logger{
		Base:     NewBase(title),
		Messages: []string{},
	}

	mod.Enabled = true

	mod.PositionData = pieces.PositionData{
		Row:       2,
		Col:       0,
		RowSpan:   8,
		ColSpan:   2,
		MinHeight: 0,
		MinWidth:  0,
	}

	return mod
}

/* -------------------- Exported Functions -------------------- */

// Clear deletes all messages from the Logger
func (l *Logger) Clear() {
	l.Messages = []string{}
}

// GetPositionData returns PositionData
func (l *Logger) GetPositionData() *pieces.PositionData {
	return &l.PositionData
}

// Log prepends a log message to the Messages slice
func (l *Logger) Log(msg string) {
	l.Messages = append(l.Messages, msg)
}

// Refresh updates the view content with the latest data
// In the Logger, Refresh() is a null op, it doesn't call out to anything
func (l *Logger) Refresh() {}

// Render draws the current string representation into the view
func (l *Logger) Render() {
	str := l.ToStr()
	l.GetView().SetText(str)
}

// ToStr returns a string representation of the module suitable for display onscreen
func (l *Logger) ToStr() string {
	str := ""

	// Messages are displayed in LIFO order
	for _, msg := range l.Messages {
		str = msg + "\n" + str
	}

	return str
}
