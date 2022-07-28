package modules

import (
	"github.com/digitalocean/godo"
	"github.com/rivo/tview"
)

// Droplets displays a list of all your available DigitalOcean droplets.
type Droplets struct {
	BaseModule

	doClient *godo.Client
}

// NewDroplets creates and returns an instance of Droplets
func NewDroplets(doClient *godo.Client) *Droplets {
	return &Droplets{
		BaseModule: NewBaseModule("Droplets", 0, 1, false),

		doClient: doClient,
	}
}

/* -------------------- Exported Functions -------------------- */

// Refresh updates the view content with the latest data
func (d *Droplets) Refresh() {
	d.GetView().SetText(d.data())
}

// View returns the tview.TextView used to display this module's data
func (d *Droplets) View() *tview.TextView {
	return d.GetView()
}

/* -------------------- Unexported Functions -------------------- */

func (d *Droplets) data() string {
	return "this is droplets"
}
