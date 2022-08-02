package services

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/senorprogrammer/dosage/modules"
)

// ServiceName is a specialized string that names a service
type ServiceName string

// Base is base
type Base struct {
	Grid    *tview.Grid
	Logger  *modules.Logger
	Modules []modules.Module
	Name    ServiceName
}

// NewBase creates and returns an instance of Base
func NewBase(serviceName ServiceName, gridTitle string, logger *modules.Logger) Base {
	b := Base{
		Grid:    newGrid(gridTitle),
		Modules: []modules.Module{},
		Name:    serviceName,
		Logger:  logger,
	}

	return b
}

/* -------------------- Exported Functions -------------------- */

// GetName returns the name of the service
func (b *Base) GetName() ServiceName {
	return b.Name
}

/* -------------------- Unexported Functions -------------------- */

// Create the grid that holds the module views
func newGrid(gridTitle string) *tview.Grid {
	grid := tview.NewGrid()
	grid.SetBorder(true)
	grid.SetTitle(fmt.Sprintf(" %s ", gridTitle))
	grid.SetRows(8, 8, 8, 8, 8, 8, 8, 8, 8, 0)
	grid.SetColumns(12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 0)

	return grid
}
