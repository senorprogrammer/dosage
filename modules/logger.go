package modules

import "github.com/rivo/tview"

// Logger stores and displays log messages emitted by other parts of the system.
type Logger struct {
	Focus    bool
	Title    string
	View     *tview.TextView
	Messages []string
}

// NewLogger creates and returns an instance of Logger
func NewLogger(title string) *Logger {
	view := tview.NewTextView()
	view.SetTitle(title)
	view.SetWrap(false)
	view.SetBorder(true)
	view.SetScrollable(true)

	return &Logger{
		Focus: false,
		View:  view,

		Messages: []string{},
	}
}

/* -------------------- Exported Functions -------------------- */

// GetFocus returns the focus val for display
func (l *Logger) GetFocus() bool {
	return l.Focus
}

// GetView returns the tview.TextView used to display this module's data
func (l *Logger) GetView() *tview.TextView {
	return l.View
}

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
