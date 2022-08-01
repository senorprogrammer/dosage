package core

import (
	"time"

	"github.com/rivo/tview"
	"github.com/senorprogrammer/dosage/services"
)

const (
	refreshInterval = 5
)

// Refresher handles refreshing everything
type Refresher struct {
	QuitChan      chan struct{}
	RefreshTicker *time.Ticker
	Services      map[services.ServiceName]services.Service
	TViewApp      *tview.Application
}

// NewRefresher creates and returns an instance of Refresher
func NewRefresher(svcs map[services.ServiceName]services.Service, tviewApp *tview.Application) *Refresher {
	return &Refresher{
		QuitChan:      make(chan struct{}),
		RefreshTicker: time.NewTicker(refreshInterval * time.Second),
		Services:      svcs,
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
