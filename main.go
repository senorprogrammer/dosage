package main

import (
	"log"
	"os"

	"github.com/senorprogrammer/dosage/flags"
	"github.com/senorprogrammer/dosage/modules"
	"github.com/senorprogrammer/dosage/pieces"
	"github.com/senorprogrammer/dosage/services"
	digitalocean "github.com/senorprogrammer/dosage/services/digitalocean"
	"github.com/senorprogrammer/dosage/splash"

	"github.com/rivo/tview"
)

const appName = "dosage"

var (
	logger    = modules.NewLogger(" logger ")
	refresher *pieces.Refresher
	svcs      = []services.Service{}
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

	// Load the services
	digitalOcean := digitalocean.NewDigitalOcean(flags.APIKey, tviewPages, logger)
	digitalOcean.LoadModules()
	svcs = append(svcs, digitalOcean)

	// Create the refresher, which handles the refresh loop
	refresher = pieces.NewRefresher(svcs, tviewApp)
	defer close(refresher.QuitChan)
	refresher.Run()

	ll("starting app...")

	// Run the underlying app loop
	if err := tviewApp.Run(); err != nil {
		panic(err)
	}
}
