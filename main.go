package main

import (
	"log"
	"os"

	"github.com/senorprogrammer/dosage/flags"
	"github.com/senorprogrammer/dosage/modules"
	"github.com/senorprogrammer/dosage/services"
	digitalocean "github.com/senorprogrammer/dosage/services/digitalocean"
	"github.com/senorprogrammer/dosage/splash"

	"github.com/rivo/tview"
)

const appName = "dosage"

var (
	logger   *modules.Logger
	svcs     = []services.Service{}
	tviewApp = tview.NewApplication()
)

func ll(msg string) {
	logger.Log(msg)
}

// refresh loops through all the modules and updates their contents
func refresh(tviewApp *tview.Application) {
	ll("refreshing...")

	logger.Refresh()

	for _, svc := range svcs {
		svc.Refresh()
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

	splash.DisplaySplashScreen()

	tviewPages := tview.NewPages()
	tviewApp.SetRoot(tviewPages, true)

	logger = modules.NewLogger(" logger ")

	// Load the services
	digitalOcean := digitalocean.NewDigitalOcean(flags.APIKey, tviewPages, logger)
	digitalOcean.LoadModules()
	svcs = append(svcs, digitalOcean)

	ll("starting app...")

	// Start the go routine that updates the module content on a timer
	quit := make(chan struct{})
	defer close(quit)

	go func(refreshFunc func(tviewApp *tview.Application), tviewApp *tview.Application) {
		refreshFunc(tviewApp)

		for {
			select {
			case <-digitalOcean.RefreshTicker.C:
				refreshFunc(tviewApp)
			case <-quit:
				digitalOcean.RefreshTicker.Stop()
				return
			}
		}
	}(refresh, tviewApp)

	// Run the underlying app loop
	if err := tviewApp.Run(); err != nil {
		panic(err)
	}
}
