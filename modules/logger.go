package modules

import (
	"time"

	"github.com/rivo/tview"
	"github.com/senorprogrammer/dosage/pieces"
)

// Logger stores and displays log messages emitted by other parts of the system.
type Logger struct {
	Base
	PositionData pieces.PositionData
	Messages     []string
}

// NewLogger creates and returns an instance of Logger
func NewLogger(title string, refreshChan chan bool) *Logger {
	mod := &Logger{
		Base:     NewBase(title, WithTextView, refreshChan, 1*time.Second, nil),
		Messages: []string{},
	}

	mod.Enabled = true

	mod.PositionData = pieces.PositionData{
		Row:       2,
		Col:       0,
		RowSpan:   8,
		ColSpan:   3,
		MinHeight: 0,
		MinWidth:  0,
	}

	mod.RefreshFunc = mod.Refresh

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
func (l *Logger) Refresh() {
	l.Render()

	// Tell the Refresher that there's new data to display
	l.RefreshChan <- true
}

// Render draws the current string representation into the view
func (l *Logger) Render() {
	str := l.ToStr()
	l.GetView().(*tview.TextView).SetText(str)
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
