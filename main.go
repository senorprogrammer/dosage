package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/digitalocean/godo"
	"github.com/senorprogrammer/dosage/flags"
	"github.com/senorprogrammer/dosage/modules"

	"github.com/logrusorgru/aurora"
	"github.com/rivo/tview"
)

const appName = "dosage"
const splashMsg string = `
 ______   _____  _______ _______  ______ _______
 |     \ |     | |______ |_____| |  ____ |______
 |_____/ |_____| ______| |     | |_____| |______
                                                
`
const (
	refreshInterval = 5
	splashInterval  = 0
)

var (
	doClient *godo.Client
	logger   *modules.Logger
	mods     = []modules.Module{}
)

// Instantiate and store all the modules
func loadModules(doClient *godo.Client, root *tview.Grid) []modules.Module {
	mods := []modules.Module{}

	logger = modules.NewLogger(" logger ")
	account := modules.NewAccount(" account ", doClient)
	droplets := modules.NewDroplets(" droplets ", doClient)
	reservedIPs := modules.NewReservedIPs(" reserved ips ", doClient)
	databases := modules.NewDatabases(" databases ", doClient)
	storage := modules.NewStorage(" volumes ", doClient)

	mods = append(mods, logger)
	mods = append(mods, account)
	mods = append(mods, droplets)
	mods = append(mods, reservedIPs)
	mods = append(mods, databases)
	mods = append(mods, storage)

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

// Create the tview app containers and load the modules into it
func newTViewApp() (*tview.Application, *tview.Grid) {
	root := tview.NewGrid()
	root.SetBorder(true)
	root.SetTitle(" dosage ")
	root.SetRows(8, 8, 8, 8, 8, 8, 8, 8, 8, 0)
	root.SetColumns(12, 12, 12, 12, 12, 12, 12, 12, 12, 0)

	tviewApp := tview.NewApplication()
	tviewApp.SetRoot(root, true).SetFocus(root)

	return tviewApp, root
}

func ll(msg string) {
	logger.Log(msg)
}

// refresh loops through all the modules and updates their contents
func refresh(tviewApp *tview.Application) {
	for _, mod := range mods {
		mod.Refresh()
	}

	tviewApp.Draw()
}

/* -------------------- Main -------------------- */

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load and parse the flags
	flags := flags.NewFlags()
	err := flags.Parse()
	if err != nil {
		os.Exit(1)
	}

	// Display the splash message
	fmt.Println(aurora.BrightGreen(splashMsg))
	time.Sleep(splashInterval * time.Second)

	// Create the tview application
	tviewApp, root := newTViewApp()

	// Create the DO client
	doClient := godo.NewFromToken(flags.APIKey)

	// Load the individual modules
	mods = loadModules(doClient, root)

	ll("starting app...")

	// Start the go routine that updates the module content on a timer
	ticker := time.NewTicker(refreshInterval * time.Second)
	quit := make(chan struct{})
	defer close(quit)

	go func(refreshFunc func(tviewApp *tview.Application), tviewApp *tview.Application) {
		refreshFunc(tviewApp)

		for {
			select {
			case <-ticker.C:
				refreshFunc(tviewApp)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}(refresh, tviewApp)

	// Run the underlying app loop
	if err := tviewApp.Run(); err != nil {
		panic(err)
	}
}
