package digitalocean

import (
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/rivo/tview"
	"github.com/senorprogrammer/dosage/modules"
	"github.com/senorprogrammer/dosage/services"
	"github.com/senorprogrammer/dosage/services/digitalocean/domodules"
)

// DigitalOcean is the container application that handles all things DigitalOcean
type DigitalOcean struct {
	services.Base

	DOClient      *godo.Client
	FocusedModule *modules.Module
	Logger        *modules.Logger
	Modules       []modules.Module

	Grid *tview.Grid
}

// NewDigitalOcean creates and returns an instance of a DigitalOcean page container
func NewDigitalOcean(apiKey string, tviewPages *tview.Pages, logger *modules.Logger) *DigitalOcean {
	grid := newGrid(" DigitalOcean ")
	tviewPages.AddPage("digitalocean", grid, true, true)

	return &DigitalOcean{
		Base:          services.Base{Name: "digitalocean"},
		DOClient:      godo.NewFromToken(apiKey),
		FocusedModule: nil,
		Logger:        logger,
		Modules:       []modules.Module{},
		Grid:          grid,
	}
}

/* -------------------- Exported Functions -------------------- */

// LoadModules instantiates each module and attaches it to the TView app
// Pass the logger in because it's common across everything and needs to
// be instantiated before the rest of the modules
func (d *DigitalOcean) LoadModules(refreshChan chan bool) {
	account := domodules.NewAccount(" account ", refreshChan, d.DOClient, d.Logger)
	billing := domodules.NewBilling(" billing ", refreshChan, d.DOClient, d.Logger)
	certs := domodules.NewCertificates(" certs ", refreshChan, d.DOClient, d.Logger)
	databases := domodules.NewDatabases(" databases ", refreshChan, d.DOClient, d.Logger)
	droplets := domodules.NewDroplets(" droplets ", refreshChan, d.DOClient, d.Logger)
	reservedIPs := domodules.NewReservedIPs(" reserved ips ", refreshChan, d.DOClient, d.Logger)
	sshKeys := domodules.NewSSHKeys(" ssh keys ", refreshChan, d.DOClient, d.Logger)
	storage := domodules.NewVolumes(" volumes ", refreshChan, d.DOClient, d.Logger)

	d.Modules = append(d.Modules, d.Logger)

	d.Modules = append(d.Modules, account)
	d.Modules = append(d.Modules, billing)
	d.Modules = append(d.Modules, certs)
	d.Modules = append(d.Modules, databases)
	d.Modules = append(d.Modules, droplets)
	d.Modules = append(d.Modules, reservedIPs)
	d.Modules = append(d.Modules, sshKeys)
	d.Modules = append(d.Modules, storage)

	for _, mod := range d.Modules {
		// Disabled modules do not get loaded
		if !mod.GetEnabled() {
			continue
		}

		// Enabled modules get add to the gridview
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

	// Start running the modules
	for _, mod := range d.Modules {
		if !mod.GetEnabled() {
			continue
		}

		mod.Run()
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
