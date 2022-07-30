package digitalocean

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
	Logger        *modules.Logger
	Modules       []modules.Module
	Name          string
	RefreshTicker *time.Ticker

	Grid *tview.Grid
}

// NewDigitalOcean creates and returns an instance of a DigitalOcean page container
func NewDigitalOcean(apiKey string, tviewPages *tview.Pages, logger *modules.Logger) *DigitalOcean {
	grid := newGrid(" DigitalOcean ")
	tviewPages.AddPage("digitalocean", grid, true, true)

	return &DigitalOcean{
		DOClient:      godo.NewFromToken(apiKey),
		FocusedModule: nil,
		Logger:        logger,
		Modules:       []modules.Module{},
		Name:          "DigitalOcean",
		RefreshTicker: time.NewTicker(refreshInterval * time.Second),
		Grid:          grid,
	}
}

/* -------------------- Exported Functions -------------------- */

// GetName returns the name of the service
func (d *DigitalOcean) GetName() string {
	return d.Name
}

// LoadModules instantiates each module and attaches it to the TView app
// Pass the logger in because it's common across everything and needs to
// be instantiated before the rest of the modules
func (d *DigitalOcean) LoadModules() {
	account := modules.NewAccount(" account ", d.DOClient)
	billing := modules.NewBilling(" billing ", d.DOClient)
	databases := modules.NewDatabases(" databases ", d.DOClient)
	droplets := modules.NewDroplets(" droplets ", d.DOClient)
	reservedIPs := modules.NewReservedIPs(" reserved ips ", d.DOClient)
	storage := modules.NewVolumes(" volumes ", d.DOClient)

	d.Modules = append(d.Modules, d.Logger)

	d.Modules = append(d.Modules, account)
	d.Modules = append(d.Modules, billing)
	d.Modules = append(d.Modules, databases)
	d.Modules = append(d.Modules, droplets)
	d.Modules = append(d.Modules, reservedIPs)
	d.Modules = append(d.Modules, storage)

	for _, mod := range d.Modules {
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
}

// Refresh tells each module to update its contents
func (d *DigitalOcean) Refresh() {
	for _, mod := range d.Modules {
		mod.Refresh()
	}
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
