package domodules

import (
	"github.com/rivo/tview"
)

// Base is base
type Base struct {
	Available bool // if a module is Available, it can be refreshed
	Focus     bool
	Title     string
	View      *tview.TextView
}

// NewBase creates and returns an instance of Base
func NewBase(title string) Base {
	view := tview.NewTextView()
	view.SetBorder(true)
	view.SetScrollable(true)
	view.SetTitle(title)
	view.SetWrap(false)
	view.SetDynamicColors(true)

	return Base{
		Available: true,
		Focus:     false,
		Title:     title,
		View:      view,
	}
}

/* -------------------- Exported Functions -------------------- */

// GetAvailable returns whether or not this module is available for refreshing
func (b *Base) GetAvailable() bool {
	return b.Available
}

// GetFocus returns the focus val for display
func (b *Base) GetFocus() bool {
	return b.Focus
}

// GetTitle returns the view title
func (b *Base) GetTitle() string {
	return b.Title
}

// GetView returns the tview.TextView used to display this module's data
func (b *Base) GetView() *tview.TextView {
	return b.View
}

// SetAvailable sets whether or not this module is available for refreshing
func (b *Base) SetAvailable(isAvailable bool) {
	b.Available = isAvailable
}
