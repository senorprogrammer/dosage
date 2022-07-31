package domodules

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/senorprogrammer/dosage/modules"
	"github.com/senorprogrammer/dosage/pieces"
)

// ReservedIPs displays a list of all your reserved IPs and which droplet they're attached to
type ReservedIPs struct {
	modules.Base
	PositionData pieces.PositionData
	doClient     *godo.Client
}

// NewReservedIPs creates and returns an instance of Droplets
func NewReservedIPs(title string, client *godo.Client) *ReservedIPs {
	mod := &ReservedIPs{
		Base: modules.NewBase(title),
		PositionData: pieces.PositionData{
			Row:       2,
			Col:       2,
			RowSpan:   2,
			ColSpan:   5,
			MinHeight: 0,
			MinWidth:  0,
		},
		doClient: client,
	}

	mod.Enabled = true

	return mod
}

/* -------------------- Exported Functions -------------------- */

// GetPositionData returns PositionData
func (r *ReservedIPs) GetPositionData() *pieces.PositionData {
	return &r.PositionData
}

// Refresh updates the view content with the latest data
func (r *ReservedIPs) Refresh() {
	if !r.GetAvailable() || !r.GetEnabled() {
		return
	}

	r.SetAvailable(false)
	r.GetView().SetText(r.data())
	r.SetAvailable(true)
}

/* -------------------- Unexported Functions -------------------- */

func (r *ReservedIPs) data() string {
	reservedIPs, err := r.fetch()
	if err != nil {
		return err.Error()
	}

	if len(reservedIPs) == 0 {
		return modules.EmptyDataLabel
	}

	data := ""

	for idx, reservedIP := range reservedIPs {
		dropletID := 0
		if reservedIP.Droplet != nil {
			dropletID = reservedIP.Droplet.ID
		}

		data = data + fmt.Sprintf("%3d\t%10d\t%16s\t%s\n", (idx+1), dropletID, reservedIP.IP, reservedIP.Region.Slug)
	}

	return data
}

// fetch uses the DigitalOcean API to fetch information about all the available droplets
func (r *ReservedIPs) fetch() ([]godo.ReservedIP, error) {
	reservedIPsList := []godo.ReservedIP{}
	opts := &godo.ListOptions{}

	for {
		doReservedIPs, resp, err := r.doClient.ReservedIPs.List(context.Background(), opts)
		if err != nil {
			return reservedIPsList, err
		}

		for _, reservedIP := range doReservedIPs {
			reservedIPsList = append(reservedIPsList, reservedIP)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return reservedIPsList, err
		}

		// Set the page we want for the next request
		opts.Page = page + 1
	}

	return reservedIPsList, nil
}
