package main

import (
	"log"
	"os"

	"github.com/senorprogrammer/dosage/core"
	"github.com/senorprogrammer/dosage/flags"
	"github.com/senorprogrammer/dosage/modules"
	"github.com/senorprogrammer/dosage/splash"

	"github.com/rivo/tview"
)

const appName = "dosage"

var (
	logger    = modules.NewLogger(" logger ")
	refresher *core.Refresher
	servicer  *core.Servicer
	tviewApp  = tview.NewApplication()
)

func ll(msg string) {
	logger.Log(msg)
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

	splash.DisplaySplashScreen()

	// Create the TView app that handles onscreen drawing
	tviewPages := tview.NewPages()
	tviewApp.SetRoot(tviewPages, true)

	// Create the Servicer, which manages loading of services
	servicer = core.NewServicer()
	servicer.LoadServices(flags, tviewPages, logger)

	// Create the refresher, which handles the refresh loop
	refresher = core.NewRefresher(tviewApp)
	defer close(refresher.QuitChan)
	refresher.Run()

	ll("starting app...")

	// Run the underlying app loop
	if err := tviewApp.Run(); err != nil {
		panic(err)
	}
}
