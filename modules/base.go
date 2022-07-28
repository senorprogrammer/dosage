package modules

import "github.com/rivo/tview"

type BaseModule struct {
	FixedSize  int
	Focus      bool
	Proportion int

	View *tview.TextView
}

/* -------------------- Exported Functions -------------------- */

// GetFixedSize returns the fixedSize val for display
func (b *BaseModule) GetFixedSize() int {
	return b.FixedSize
}

// GetFocus returns the focus val for display
func (b *BaseModule) GetFocus() bool {
	return b.Focus
}

// GetProportion returns the proportion for display
func (b *BaseModule) GetProportion() int {
	return b.Proportion
}

// GetView returns the tview.TextView used to display this module's data
func (b *BaseModule) GetView() *tview.TextView {
	return b.View
}
