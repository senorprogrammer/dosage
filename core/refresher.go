package core

import (
	"github.com/rivo/tview"
)

const (
	refreshInterval = 5
)

// Refresher handles refreshing everything
type Refresher struct {
	RefreshChan chan bool
	TViewApp    *tview.Application
}

// NewRefresher creates and returns an instance of Refresher
func NewRefresher(tviewApp *tview.Application) *Refresher {
	return &Refresher{
		RefreshChan: make(chan bool),
		TViewApp:    tviewApp,
	}
}

/* -------------------- Exported Functions -------------------- */

// Run starts the refresh loop
func (r *Refresher) Run() {
	go func(refreshFunc func()) {
		refreshFunc()

		for {
			select {
			case <-r.RefreshChan:
				refreshFunc()
			}
		}
	}(r.refresh)
}

/* -------------------- Unexported Functions -------------------- */

// refresh loops through all the modules and updates their contents
func (r *Refresher) refresh() {
	r.TViewApp.Draw()
}
