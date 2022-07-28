package modules

import (
	"context"
	"fmt"

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
	droplets, err := d.dropletsFetch()
	if err != nil {
		return err.Error()
	}

	data := ""

	for idx, droplet := range droplets {
		data = data + fmt.Sprintf("%3d\t%12d\t%s\n", (idx+1), droplet.ID, droplet.Name)
	}

	return data
}

// dropletsFetch uses the DigitalOcean API to fetch information about all the available droplets
func (d *Droplets) dropletsFetch() ([]godo.Droplet, error) {
	dropletList := []godo.Droplet{}
	opts := &godo.ListOptions{}

	for {
		doDroplets, resp, err := d.doClient.Droplets.List(context.Background(), opts)
		if err != nil {
			return dropletList, err
		}

		for _, doDroplet := range doDroplets {
			dropletList = append(dropletList, doDroplet)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return dropletList, err
		}

		// Set the page we want for the next request
		opts.Page = page + 1
	}

	return dropletList, nil
}
