package main

import (
	"fmt"
	"log"
	"os"

	"github.com/senorprogrammer/dosage/core"
	"github.com/senorprogrammer/dosage/flags"
	"github.com/senorprogrammer/dosage/modules"
	"github.com/senorprogrammer/dosage/splashscreen"

	"github.com/rivo/tview"
)

const appName = "dosage"

var (
	logger    *modules.Logger
	refresher *core.Refresher
	servicer  *core.Servicer
	tviewApp  = tview.NewApplication()
)

/* -------------------- Main -------------------- */

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load and parse the flags
	flags := flags.NewFlags()
	err := flags.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	splashscreen.Show()

	// Create the TView app that handles onscreen drawing
	tviewPages := tview.NewPages()
	tviewApp.SetRoot(tviewPages, true)

	// Create the refresher, which handles the refresh loop
	refresher = core.NewRefresher(tviewApp)
	defer close(refresher.RefreshChan)
	refresher.Run()

	logger = modules.NewLogger(" logger ", refresher.RefreshChan)

	// Create the Servicer, which manages loading of services
	servicer = core.NewServicer()
	servicer.LoadServices(flags, tviewPages, refresher.RefreshChan, logger)

	logger.Log("starting app...")

	// Run the underlying app loop
	if err := tviewApp.Run(); err != nil {
		panic(err)
	}
}
