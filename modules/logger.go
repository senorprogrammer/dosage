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
		Base: NewBase(title),
		PositionData: pieces.PositionData{
			Row:       2,
			Col:       0,
			RowSpan:   8,
			ColSpan:   2,
			MinHeight: 0,
			MinWidth:  0,
		},
		Messages: []string{},
	}

	mod.Enabled = true

	return mod
}

/* -------------------- Exported Functions -------------------- */

// GetPositionData returns PositionData
func (l *Logger) GetPositionData() *pieces.PositionData {
	return &l.PositionData
}

// Log prepends a log message to the Messages slice
func (l *Logger) Log(msg string) {
	l.Messages = append(l.Messages, msg)
	l.Refresh()
}

// Refresh updates the view content with the latest data
func (l *Logger) Refresh() {
	l.GetView().SetText(l.data())
}

/* -------------------- Unexported Functions -------------------- */

// data returns a string representation of the module
// suitable for display onscreen
func (l *Logger) data() string {
	data := ""

	// Messages are displayed in LIFO order
	for _, msg := range l.Messages {
		data = msg + "\n" + data
	}

	return data
}
