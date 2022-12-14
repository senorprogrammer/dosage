package modules

import (
	"time"

	"github.com/rivo/tview"
	"github.com/senorprogrammer/dosage/pieces"
)

const (
	// DefaultRefreshSeconds defines how often the module data should refresh
	DefaultRefreshSeconds = 30

	// EmptyContentLabel is the content to display if there is no content
	EmptyContentLabel = "none"
)

// ViewType defines the enum that defines which TViews can be instantiated by modules
type ViewType int64

const (
	// WithTableView is used to indicate a table should be instantiated
	WithTableView ViewType = iota

	// WithTextView is used to indicate a textView should be instantiated
	WithTextView
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
	View         tview.Primitive // The view to display the module data in

	// Properties relevant to refreshing the module data
	QuitChan        chan struct{} // The channel that's used to stop the RefreshTicker
	RefreshChan     chan bool     // The channel to call into when a refresh completes
	RefreshFunc     func()        // The function to execute when the RefreshTicker ticks
	RefreshInterval time.Duration // Defines the number of seconds between data refreshes
	RefreshTicker   *time.Ticker  // Controls how often, in seconds, the module will refresh its data
}

// NewBase creates and returns an instance of Base
func NewBase(title string, viewType ViewType, refreshChan chan bool, refreshInterval time.Duration, logger *Logger) Base {
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
	}

	// This ticker controls how often the module's data is refreshed and redrawn
	base.RefreshTicker = time.NewTicker(base.RefreshInterval)

	// Create and store the data view that'll be used by the underlying module to display data
	switch viewType {
	case WithTableView:
		base.View = base.newTableView(title)
	case WithTextView:
		base.View = base.newTextView(title)
	default:
		base.View = base.newTextView(title)
	}

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

// GetView returns the tview.Primitive used to display this module's data
func (b *Base) GetView() tview.Primitive {
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

func (b *Base) newTableView(title string) *tview.Table {
	view := tview.NewTable()
	view.SetBorder(true)
	view.SetBorderPadding(0, 0, 1, 1)
	view.SetTitle(title)

	return view
}

func (b *Base) newTextView(title string) *tview.TextView {
	view := tview.NewTextView()
	view.SetBorder(true)
	view.SetScrollable(true)
	view.SetTitle(title)
	view.SetWrap(false)
	view.SetDynamicColors(true)
	view.SetBorderPadding(0, 0, 1, 1)

	return view
}
