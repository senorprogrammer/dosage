package modules

import (
	"github.com/rivo/tview"
	"github.com/senorprogrammer/dosage/pieces"
)

// Module defines the interface that all modules must adhere to
type Module interface {
	GetFocus() bool
	GetPositionData() *pieces.PositionData
	GetView() *tview.TextView

	Refresh()
}
