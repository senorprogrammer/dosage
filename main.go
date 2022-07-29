package main

import (
	"log"
	"os"

	"github.com/digitalocean/godo"
	"github.com/senorprogrammer/dosage/flags"
	"github.com/senorprogrammer/dosage/modules"
	"github.com/senorprogrammer/dosage/pieces"
	digitalocean "github.com/senorprogrammer/dosage/services/digitalocean"

	"github.com/rivo/tview"
)

const appName = "dosage"

var (
	doClient *godo.Client
	logger   *modules.Logger
	mods     = []modules.Module{}
	tviewApp = tview.NewApplication()
)

func ll(msg string) {
	logger.Log(msg)
}

// refresh loops through all the modules and updates their contents
func refresh(tviewApp *tview.Application) {
	ll("refreshing...")

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

	pieces.DisplaySplashScreen()

	tviewPages := tview.NewPages()
	dosage := digitalocean.NewDosage(flags.APIKey, appName, tviewPages)
	tviewApp.SetRoot(tviewPages, true)

	// Load the individual modules
	logger = modules.NewLogger(" logger ")
	mods = dosage.LoadModules(logger)

	ll("starting app...")

	// Start the go routine that updates the module content on a timer
	quit := make(chan struct{})
	defer close(quit)

	go func(refreshFunc func(tviewApp *tview.Application), tviewApp *tview.Application) {
		refreshFunc(tviewApp)

		for {
			select {
			case <-dosage.RefreshTicker.C:
				refreshFunc(tviewApp)
			case <-quit:
				dosage.RefreshTicker.Stop()
				return
			}
		}
	}(refresh, tviewApp)

	// Run the underlying app loop
	if err := tviewApp.Run(); err != nil {
		panic(err)
	}
}
