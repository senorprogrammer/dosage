package modules

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/rivo/tview"
	"github.com/senorprogrammer/dosage/do"
)

type ReservedIPs struct {
	FixedSize  int
	Focus      bool
	Proportion int
	Title      string
	View       *tview.TextView

	doClient *godo.Client
}

// NewReservedIPs creates and returns an instance of Droplets
func NewReservedIPs(apiKey string) *ReservedIPs {
	view := tview.NewTextView()
	view.SetTitle("reserved ips")
	view.SetWrap(false)
	view.SetBorder(true)

	return &ReservedIPs{
		FixedSize:  0,
		Proportion: 1,
		Focus:      false,
		View:       view,

		doClient: do.NewClient(apiKey),
	}
}

/* -------------------- Exported Functions -------------------- */

// GetFixedSize returns the fixedSize val for display
func (r *ReservedIPs) GetFixedSize() int {
	return r.FixedSize
}

// GetFocus returns the focus val for display
func (r *ReservedIPs) GetFocus() bool {
	return r.Focus
}

// GetProportion returns the proportion for display
func (r *ReservedIPs) GetProportion() int {
	return r.Proportion
}

// GetView returns the tview.TextView used to display this module's data
func (r *ReservedIPs) GetView() *tview.TextView {
	return r.View
}

// Refresh updates the view content with the latest data
func (r *ReservedIPs) Refresh() {
	r.GetView().SetText(r.data())
}

/* -------------------- Unexported Functions -------------------- */

func (r *ReservedIPs) data() string {
	reservedIPs, err := r.fetch()
	if err != nil {
		return err.Error()
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
