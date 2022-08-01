package core

import (
	"time"

	"github.com/rivo/tview"
)

const (
	refreshInterval = 5
)

// Refresher handles refreshing everything
type Refresher struct {
	QuitChan      chan struct{}
	RefreshTicker *time.Ticker
	TViewApp      *tview.Application
}

// NewRefresher creates and returns an instance of Refresher
func NewRefresher(tviewApp *tview.Application) *Refresher {
	return &Refresher{
		QuitChan:      make(chan struct{}),
		RefreshTicker: time.NewTicker(refreshInterval * time.Second),
		TViewApp:      tviewApp,
	}
}

/* -------------------- Exported Functions -------------------- */

// Run starts the refresh loop
func (r *Refresher) Run() {
	go func(refreshFunc func()) {
		refreshFunc()

		for {
			select {
			case <-r.RefreshTicker.C:
				refreshFunc()
			case <-r.QuitChan:
				r.RefreshTicker.Stop()
				return
			}
		}
	}(r.refresh)
}

/* -------------------- Unexported Functions -------------------- */

// refresh loops through all the modules and updates their contents
func (r *Refresher) refresh() {
	r.TViewApp.Draw()
}
