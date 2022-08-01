package modules

import (
	"time"

	"github.com/rivo/tview"
	"github.com/senorprogrammer/dosage/pieces"
)

const (
	// EmptyContentLabel is the content to display if there is no content
	EmptyContentLabel = "none"
)

// Base is base
type Base struct {
	Available    bool  // If a module is Available, it can be refreshed
	Enabled      bool  // If a module is Enabled, it can be refreshed
	Focus        bool  // Whether or not this module should have the keyboard focus
	LastError    error // If a refresh generates an error, the error will be stored here
	Logger       *Logger
	PositionData pieces.PositionData
	Title        string          // The text string to be displayed at the top of the module view
	View         *tview.TextView // The view to display the module data in

	// Properties relevant to refreshing the module data
	QuitChan        chan struct{} // The channel that's used to stop the RefreshTicker
	RefreshChan     chan bool     // The channel to call into when a refresh completes
	RefreshFunc     func()        // The function to execute when the RefreshTicker ticks
	RefreshInterval time.Duration // Defines the number of seconds between data refreshes
	RefreshTicker   *time.Ticker  // Controls how often, in seconds, the module will refresh its data
}

// NewBase creates and returns an instance of Base
func NewBase(title string, refreshChan chan bool, refreshInterval time.Duration, logger *Logger) Base {
	view := tview.NewTextView()
	view.SetBorder(true)
	view.SetScrollable(true)
	view.SetTitle(title)
	view.SetWrap(false)
	view.SetDynamicColors(true)
	view.SetBorderPadding(0, 0, 1, 1)

	base := Base{
		Available:       true,  // Modules are available unless they're fetching data
		Enabled:         false, // Modules are disabled by default, enabled explicitly
		Focus:           false, // Modules are unfoused by default, receiving focus explicitly
		Logger:          logger,
		QuitChan:        make(chan struct{}),
		RefreshChan:     refreshChan,
		RefreshFunc:     nil,
		RefreshInterval: refreshInterval,
		Title:           title,
		View:            view,
	}

	// This ticker controls how often the module's data is refreshed and redrawn
	base.RefreshTicker = time.NewTicker(base.RefreshInterval)

	return base
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

// Run starts the refresh loop for this module
func (b *Base) Run() {
	go func(refreshFunc func()) {
		refreshFunc()

		for {
			select {
			case <-b.RefreshTicker.C:
				refreshFunc()
			case <-b.QuitChan:
				b.RefreshTicker.Stop()
				return
			}
		}
	}(b.RefreshFunc)
}

/* -------------------- Unexported Functions -------------------- */
