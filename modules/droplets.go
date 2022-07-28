package modules

import (
	"github.com/digitalocean/godo"
	"github.com/rivo/tview"
)

// Droplets displays a list of all your available DigitalOcean droplets.
type Droplets struct {
	BaseModule

	doClient *godo.Client
	view     *tview.TextView
}

// NewDroplets creates and returns an instance of Droplets
func NewDroplets(doClient *godo.Client) *Droplets {
	view := tview.NewTextView()
	view.SetTitle(" droplets ")
	view.SetBorder(true)
	view.SetWrap(false)

	return &Droplets{
		BaseModule: BaseModule{
			FixedSize:  0,
			Focus:      false,
			Proportion: 1,
			View:       view,
		},

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
	return d.view
}

/* -------------------- Unexported Functions -------------------- */

func (d *Droplets) data() string {
	return "this is droplets"
}
