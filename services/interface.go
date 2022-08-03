package services

import "github.com/rivo/tview"

// Serviceable defines the interface for external services
type Serviceable interface {
	GetGrid() *tview.Grid
	GetName() ServiceName
	LoadModules(refreshChan chan bool)
}
