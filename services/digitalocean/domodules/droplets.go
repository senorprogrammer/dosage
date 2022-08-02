package domodules

import (
	"context"
	"fmt"
	"time"

	"github.com/digitalocean/godo"
	"github.com/rivo/tview"
	"github.com/senorprogrammer/dosage/modules"
	"github.com/senorprogrammer/dosage/pieces"
)

// Droplets displays a list of all your available DigitalOcean droplets.
type Droplets struct {
	modules.Base
	Droplets []godo.Droplet
	doClient *godo.Client
}

// NewDroplets creates and returns an instance of Droplets
func NewDroplets(title string, refreshChan chan bool, client *godo.Client, logger *modules.Logger) *Droplets {
	mod := &Droplets{
		Base:     modules.NewBase(title, modules.WithTextView, refreshChan, 5*time.Second, logger),
		Droplets: []godo.Droplet{},
		doClient: client,
	}

	mod.Enabled = true

	mod.PositionData = pieces.PositionData{
		Row:       0,
		Col:       2,
		RowSpan:   2,
		ColSpan:   5,
		MinHeight: 0,
		MinWidth:  0,
	}

	mod.RefreshFunc = mod.Refresh

	return mod
}

/* -------------------- Exported Functions -------------------- */

// GetPositionData returns PositionData
func (d *Droplets) GetPositionData() *pieces.PositionData {
	return &d.PositionData
}

// Refresh updates the view content with the latest data
func (d *Droplets) Refresh() {
	if !d.GetAvailable() || !d.GetEnabled() {
		return
	}

	d.Logger.Log(fmt.Sprintf("refreshing %s", d.GetTitle()))

	d.SetAvailable(false)

	droplets, err := d.fetch()
	if err != nil {
		d.LastError = err
	} else {
		d.LastError = nil
		d.Droplets = droplets
	}

	d.SetAvailable(true)

	d.Render()

	// Tell the Refresher that there's new data to display
	d.RefreshChan <- true
}

// Render draws the current string representation into the view
func (d *Droplets) Render() {
	str := d.ToStr()
	d.GetView().(*tview.TextView).SetText(str)
}

// ToStr returns a string representation of the module suitable for display onscreen
func (d *Droplets) ToStr() string {
	if d.LastError != nil {
		return d.LastError.Error()
	}

	if len(d.Droplets) == 0 {
		return modules.EmptyContentLabel
	}

	str := ""

	for _, droplet := range d.Droplets {
		str += fmt.Sprintf(
			"%10d\t%s\t%s\n",
			droplet.ID,
			droplet.Name,
			pieces.ColorForState(droplet.Status, droplet.Status),
		)
	}

	return str
}

/* -------------------- Unexported Functions -------------------- */

// fetch uses the DigitalOcean API to fetch information about all the available droplets
func (d *Droplets) fetch() ([]godo.Droplet, error) {
	dropletList := []godo.Droplet{}
	opts := &godo.ListOptions{}

	for {
		doDroplets, resp, err := d.doClient.Droplets.List(context.Background(), opts)
		if err != nil {
			return dropletList, err
		}

		dropletList = append(dropletList, doDroplets...)

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
