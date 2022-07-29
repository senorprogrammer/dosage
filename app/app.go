package app

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

// Dosage is the container application that handles everything. Everything.
type Dosage struct {
	DOClient      *godo.Client
	FocusedModule *modules.Module
	Modules       []modules.Module
	RefreshTicker *time.Ticker

	TViewApp *tview.Application
	Root     *tview.Grid
}

// NewDosage creates and returns an instance of a Dosage app
func NewDosage(apiKey string, appName string) *Dosage {
	tviewApp, root := newTViewApp(appName)

	return &Dosage{
		DOClient:      godo.NewFromToken(apiKey),
		FocusedModule: nil,
		Modules:       []modules.Module{},
		RefreshTicker: time.NewTicker(refreshInterval * time.Second),

		TViewApp: tviewApp,
		Root:     root,
	}
}

/* -------------------- Exported Functions -------------------- */

func (d *Dosage) LoadModules(root *tview.Grid, logger *modules.Logger) []modules.Module {
	mods := []modules.Module{}

	// logger := modules.NewLogger(" logger ")
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
		root.AddItem(
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

// Refresh refreshes the content in all the modules
func (d *Dosage) Refresh() {
	for _, mod := range d.Modules {
		mod.Refresh()
	}

	d.TViewApp.Draw()
}

/* -------------------- Unexported Functions -------------------- */

// Create the tview app containers and load the modules into it
func newTViewApp(appName string) (*tview.Application, *tview.Grid) {
	root := tview.NewGrid()
	root.SetBorder(true)
	root.SetTitle(fmt.Sprintf(" %s ", appName))
	root.SetRows(8, 8, 8, 8, 8, 8, 8, 8, 8, 0)
	root.SetColumns(12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 0)

	tviewApp := tview.NewApplication()
	tviewApp.SetRoot(root, true).SetFocus(root)

	return tviewApp, root
}
