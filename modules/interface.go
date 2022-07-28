package modules

import "github.com/rivo/tview"

// Module defines the interface that all modules must adhere to
type Module interface {
	Data() string
	Refresh()
	View() *tview.TextView
}
