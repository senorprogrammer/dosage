package modules

import (
	"github.com/rivo/tview"
)

// Base is base
type Base struct {
	Focus bool
	Title string
	View  *tview.TextView
}

// NewBase creates and returns an instance of Base
func NewBase(title string) Base {
	view := tview.NewTextView()
	view.SetBorder(true)
	view.SetScrollable(true)
	view.SetTitle(title)
	view.SetWrap(false)

	return Base{
		Focus: false,
		Title: title,
		View:  view,
	}
}

/* -------------------- Exported Functions -------------------- */

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
