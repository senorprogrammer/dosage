package modules

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/senorprogrammer/dosage/pieces"
)

// Droplets displays a list of all your available DigitalOcean droplets.
type Droplets struct {
	Base
	PositionData pieces.PositionData
	doClient     *godo.Client
}

// NewDroplets creates and returns an instance of Droplets
func NewDroplets(title string, client *godo.Client) *Droplets {
	return &Droplets{
		Base: NewBase(title),
		PositionData: pieces.PositionData{
			Row:       0,
			Col:       2,
			RowSpan:   2,
			ColSpan:   5,
			MinHeight: 0,
			MinWidth:  0,
		},
		doClient: client,
	}
}

/* -------------------- Exported Functions -------------------- */

// Refresh updates the view content with the latest data
func (d *Droplets) Refresh() {
	d.GetView().SetText(d.data())
}

/* -------------------- Unexported Functions -------------------- */

func (d *Droplets) data() string {
	droplets, err := d.fetch()
	if err != nil {
		return err.Error()
	}

	if len(droplets) == 0 {
		return "none"
	}

	data := ""

	for idx, droplet := range droplets {
		data = data + fmt.Sprintf("%3d\t%12d\t%s\n", (idx+1), droplet.ID, droplet.Name)
	}

	return data
}

// fetch uses the DigitalOcean API to fetch information about all the available droplets
func (d *Droplets) fetch() ([]godo.Droplet, error) {
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
