package services

import (
	"fmt"
	"time"

	"github.com/digitalocean/godo"
	"github.com/rivo/tview"
	"github.com/senorprogrammer/dosage/modules"
)

const (
	refreshInterval = 5
)

// DigitalOcean is the container application that handles all things DigitalOcean
type DigitalOcean struct {
	DOClient      *godo.Client
	FocusedModule *modules.Module
	Modules       []modules.Module
	RefreshTicker *time.Ticker

	// TViewPages *tview.Pages
	Grid *tview.Grid
}

// NewDosage creates and returns an instance of a Dosage app
func NewDosage(apiKey string, appName string, tviewPages *tview.Pages) *DigitalOcean {
	grid := newGrid(appName)
	tviewPages.AddPage(appName, grid, true, true)

	return &DigitalOcean{
		DOClient:      godo.NewFromToken(apiKey),
		FocusedModule: nil,
		Modules:       []modules.Module{},
		RefreshTicker: time.NewTicker(refreshInterval * time.Second),
		Grid:          grid,
	}
}

/* -------------------- Exported Functions -------------------- */

// LoadModules instantiates each module and attaches it to the TView app
// Pass the logger in because it's common across everything and needs to
// be instantiated before the rest of the modules
func (d *DigitalOcean) LoadModules(logger *modules.Logger) []modules.Module {
	mods := []modules.Module{}

	account := modules.NewAccount(" account ", d.DOClient)
	droplets := modules.NewDroplets(" droplets ", d.DOClient)
	reservedIPs := modules.NewReservedIPs(" reserved ips ", d.DOClient)
	databases := modules.NewDatabases(" databases ", d.DOClient)
	storage := modules.NewVolumes(" volumes ", d.DOClient)
	billing := modules.NewBilling(" billing ", d.DOClient)

	mods = append(mods, logger)
	mods = append(mods, account)
	mods = append(mods, droplets)
	mods = append(mods, reservedIPs)
	mods = append(mods, databases)
	mods = append(mods, storage)
	mods = append(mods, billing)

	for _, mod := range mods {
		d.Grid.AddItem(
			mod.GetView(),
			mod.GetPositionData().GetRow(),
			mod.GetPositionData().GetCol(),
			mod.GetPositionData().GetRowSpan(),
			mod.GetPositionData().GetColSpan(),
			mod.GetPositionData().GetMinWidth(),
			mod.GetPositionData().GetMinHeight(),
			false,
		)
	}

	return mods
}

/* -------------------- Unexported Functions -------------------- */

// Create the grid that holds the module views
func newGrid(appName string) *tview.Grid {
	grid := tview.NewGrid()
	grid.SetBorder(true)
	grid.SetTitle(fmt.Sprintf(" %s ", appName))
	grid.SetRows(8, 8, 8, 8, 8, 8, 8, 8, 8, 0)
	grid.SetColumns(12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 0)

	return grid
}
