package modules

import (
	"github.com/rivo/tview"
)

// Logger stores the log messages sent to it and displays them in the Logger view
type Logger struct {
	WidthPercent  int
	HeightPercent int

	Messages []string

	view *tview.TextView
}

// NewLogger creates and returns an instance of Logger
func NewLogger(widthPerc, heightPerc int) *Logger {
	view := tview.NewTextView()
	view.SetTitle(" logger ")
	view.SetBorder(true)
	view.SetWrap(false)

	return &Logger{
		WidthPercent:  widthPerc,
		HeightPercent: heightPerc,

		Messages: []string{},

		view: view,
	}
}

// Log prepends a log message to the Messages slice
func (l *Logger) Log(msg string) {
	l.Messages = append(l.Messages, msg)
}

// Data returns a string representation of the module
// suitable for display onscreen
func (l *Logger) Data() string {
	data := ""

	// Messages are displayed in LIFO order
	for _, msg := range l.Messages {
		data = msg + "\n" + data
	}

	return data
}

// Refresh updates the view content with the latest data
func (l *Logger) Refresh() {
	l.view.SetText(l.Data())
}

// View returns the tview.TextView used to display this module's data
func (l *Logger) View() *tview.TextView {
	return l.view
}
