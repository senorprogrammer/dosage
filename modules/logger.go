package modules

import (
	"github.com/rivo/tview"
)

// Logger stores and displays log messages emitted by other parts of the system.
type Logger struct {
	BaseModule

	Messages []string

	view *tview.TextView
}

// NewLogger creates and returns an instance of Logger
func NewLogger() *Logger {
	view := tview.NewTextView()
	view.SetTitle(" logger ")
	view.SetBorder(true)
	view.SetWrap(false)

	return &Logger{
		BaseModule: BaseModule{
			FixedSize:  0,
			Focus:      true,
			Proportion: 1,
			View:       view,
		},

		Messages: []string{},
	}
}

/* -------------------- Exported Functions -------------------- */

// Log prepends a log message to the Messages slice
func (l *Logger) Log(msg string) {
	l.Messages = append(l.Messages, msg)
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
