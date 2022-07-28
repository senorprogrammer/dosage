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
func newTViewApp() (*tview.Application, *tview.Flex) {
	root := tview.NewFlex()
	root.SetBorder(true)
	root.SetTitle(" dosage ")

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

	// Create the modules
	logger = modules.NewLogger()
	mods = append(mods, logger)

	droplets := modules.NewDroplets(flags.APIKey)
	mods = append(mods, droplets)

	ll("starting...")

	ll("adding logger")
	root.AddItem(logger.GetView(), 0, 1, true)

	ll("adding droplets")
	root.AddItem(droplets.GetView(), 0, 1, false)

	ll(fmt.Sprintf("using api key %s", flags.APIKey))

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
