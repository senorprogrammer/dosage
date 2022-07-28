package modules

import (
	"github.com/digitalocean/godo"
	"github.com/rivo/tview"
	"github.com/senorprogrammer/dosage/do"
)

// Droplets displays a list of all your available DigitalOcean droplets.
type Droplets struct {
	FixedSize  int
	Focus      bool
	Proportion int
	Title      string
	View       *tview.TextView

	doClient *godo.Client
}

// NewDroplets creates and returns an instance of Droplets
func NewDroplets(apiKey string) *Droplets {
	view := tview.NewTextView()
	view.SetTitle("droplets")
	view.SetWrap(false)
	view.SetBorder(true)

	return &Droplets{
		FixedSize:  0,
		Proportion: 1,
		Focus:      false,
		View:       view,

		doClient: do.NewClient(apiKey),
	}
}

/* -------------------- Exported Functions -------------------- */

// GetFixedSize returns the fixedSize val for display
func (d *Droplets) GetFixedSize() int {
	return d.FixedSize
}

// GetFocus returns the focus val for display
func (d *Droplets) GetFocus() bool {
	return d.Focus
}

// GetProportion returns the proportion for display
func (d *Droplets) GetProportion() int {
	return d.Proportion
}

// GetView returns the tview.TextView used to display this module's data
func (d *Droplets) GetView() *tview.TextView {
	return d.View
}

// Refresh updates the view content with the latest data
func (d *Droplets) Refresh() {
	d.GetView().SetText(d.data())
}

/* -------------------- Unexported Functions -------------------- */

func (d *Droplets) data() string {
	return "this is droplets"
}
