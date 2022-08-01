package modules

import (
	"github.com/rivo/tview"
)

const (
	// EmptyContentLabel is the content to display if there is no content
	EmptyContentLabel = "none"
)

// Base is base
type Base struct {
	Available bool // If a module is Available, it can be refreshed
	Enabled   bool
	Focus     bool
	LastError error
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
	view.SetBorderPadding(0, 0, 1, 1)

	return Base{
		Available: true,  // Modules are available unless they're fetching data
		Enabled:   false, // Modules are disabled by default, enabled explicitly
		Focus:     false, // Modules are unfoused by default, receiving focus explicitly
		Title:     title,
		View:      view,
	}
}

/* -------------------- Exported Functions -------------------- */

// GetAvailable returns whether or not this module is available for refreshing
func (b *Base) GetAvailable() bool {
	return b.Available
}

// GetEnabled returns whether or not a module is enabled
func (b *Base) GetEnabled() bool {
	return b.Enabled
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
