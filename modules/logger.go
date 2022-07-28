package modules

import "github.com/rivo/tview"

// Logger stores and displays log messages emitted by other parts of the system.
type Logger struct {
	FixedSize  int
	Focus      bool
	Proportion int
	Title      string
	View       *tview.TextView
	Messages   []string
}

// NewLogger creates and returns an instance of Logger
func NewLogger() *Logger {
	view := tview.NewTextView()
	view.SetTitle("logger")
	view.SetWrap(false)
	view.SetBorder(true)

	return &Logger{
		FixedSize:  30,
		Proportion: 1,
		Focus:      false,
		View:       view,

		Messages: []string{},
	}
}

/* -------------------- Exported Functions -------------------- */

// GetFixedSize returns the fixedSize val for display
func (l *Logger) GetFixedSize() int {
	return l.FixedSize
}

// GetFocus returns the focus val for display
func (l *Logger) GetFocus() bool {
	return l.Focus
}

// GetProportion returns the proportion for display
func (l *Logger) GetProportion() int {
	return l.Proportion
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
