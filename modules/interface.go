package modules

import "github.com/rivo/tview"

// Module defines the interface that all modules must adhere to
type Module interface {
	GetFixedSize() int
	GetFocus() bool
	GetProportion() int
	GetView() *tview.TextView

	Refresh()
}
