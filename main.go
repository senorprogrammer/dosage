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

// Create the tview app containers and load the modules into it
func newTViewApp() (*tview.Application, *tview.Grid) {
	root := tview.NewGrid()
	root.SetBorder(true)
	root.SetTitle(" dosage ")
	root.SetRows(8, 8, 8, 8, 8, 8, 8, 0)
	root.SetColumns(48, 12, 12, 12, 12, 12, 12, 12, 12, 0)

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
	client := godo.NewFromToken(flags.APIKey)

	// Create the modules
	logger = modules.NewLogger(" logger ")
	account := modules.NewAccount(" account ", client)
	droplets := modules.NewDroplets(" droplets ", client)
	reservedIPs := modules.NewReservedIPs(" reserved ips ", client)
	databases := modules.NewDatabases(" databases ", client)
	storage := modules.NewStorage(" volumes ", client)

	mods = append(mods, logger)
	mods = append(mods, account)
	mods = append(mods, droplets)
	mods = append(mods, reservedIPs)
	mods = append(mods, databases)
	mods = append(mods, storage)

	ll("starting...")

	ll("adding logger")
	root.AddItem(
		logger.GetView(),
		logger.PositionData.Row,
		logger.PositionData.Col,
		logger.PositionData.RowSpan,
		logger.PositionData.ColSpan,
		logger.PositionData.MinWidth,
		logger.PositionData.MinHeight,
		false,
	)

	ll("adding account")
	root.AddItem(
		account.GetView(),
		account.PositionData.Row,
		account.PositionData.Col,
		account.PositionData.RowSpan,
		account.PositionData.ColSpan,
		account.PositionData.MinWidth,
		account.PositionData.MinHeight,
		false,
	)

	ll("adding droplets")
	root.AddItem(
		droplets.GetView(),
		droplets.PositionData.Row,
		droplets.PositionData.Col,
		droplets.PositionData.RowSpan,
		droplets.PositionData.ColSpan,
		droplets.PositionData.MinWidth,
		droplets.PositionData.MinHeight,
		false,
	)

	ll("adding reservedips")
	root.AddItem(
		reservedIPs.GetView(),
		reservedIPs.PositionData.Row,
		reservedIPs.PositionData.Col,
		reservedIPs.PositionData.RowSpan,
		reservedIPs.PositionData.ColSpan,
		reservedIPs.PositionData.MinWidth,
		reservedIPs.PositionData.MinHeight,
		false,
	)

	ll("adding databases")
	root.AddItem(
		databases.GetView(),
		databases.PositionData.Row,
		databases.PositionData.Col,
		databases.PositionData.RowSpan,
		databases.PositionData.ColSpan,
		databases.PositionData.MinWidth,
		databases.PositionData.MinHeight,
		false,
	)

	ll("adding storage")
	root.AddItem(
		storage.GetView(),
		storage.PositionData.Row,
		storage.PositionData.Col,
		storage.PositionData.RowSpan,
		storage.PositionData.ColSpan,
		storage.PositionData.MinWidth,
		storage.PositionData.MinHeight,
		false,
	)

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
