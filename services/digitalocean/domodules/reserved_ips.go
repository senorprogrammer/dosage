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

// ReservedIPs displays a list of all your reserved IPs and which droplet they're attached to
type ReservedIPs struct {
	modules.Base
	ReservedIPs []godo.ReservedIP
	doClient    *godo.Client
}

// NewReservedIPs creates and returns an instance of Droplets
func NewReservedIPs(title string, refreshChan chan bool, client *godo.Client, logger *modules.Logger) *ReservedIPs {
	mod := &ReservedIPs{
		Base:        modules.NewBase(title, modules.WithTextView, refreshChan, 5*time.Second, logger),
		ReservedIPs: []godo.ReservedIP{},
		doClient:    client,
	}

	mod.Enabled = true

	mod.PositionData = pieces.PositionData{
		Row:       2,
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
func (r *ReservedIPs) GetPositionData() *pieces.PositionData {
	return &r.PositionData
}

// Refresh updates the view content with the latest data
func (r *ReservedIPs) Refresh() {
	if !r.GetAvailable() || !r.GetEnabled() {
		return
	}

	r.Logger.Log(fmt.Sprintf("refreshing %s", r.GetTitle()))

	r.SetAvailable(false)

	rip, err := r.fetch()
	if err != nil {
		r.LastError = err
	} else {
		r.LastError = nil
		r.ReservedIPs = rip
	}

	r.SetAvailable(true)

	r.Render()

	// Tell the Refresher that there's new data to display
	r.RefreshChan <- true
}

// Render draws the current string representation into the view
func (r *ReservedIPs) Render() {
	str := r.ToStr()
	r.GetView().(*tview.TextView).SetText(str)
}

// ToStr returns a string representation of the module suitable for display onscreen
func (r *ReservedIPs) ToStr() string {
	if r.LastError != nil {
		return r.LastError.Error()
	}

	if len(r.ReservedIPs) == 0 {
		return modules.EmptyContentLabel
	}

	str := ""

	for _, reservedIP := range r.ReservedIPs {
		dropletID := 0
		if reservedIP.Droplet != nil {
			dropletID = reservedIP.Droplet.ID
		}

		str = str + fmt.Sprintf(
			"%10d\t%16s\t%s\n",
			dropletID,
			reservedIP.IP,
			reservedIP.Region.Slug,
		)
	}

	return str
}

/* -------------------- Unexported Functions -------------------- */

// fetch uses the DigitalOcean API to fetch information about all the available droplets
func (r *ReservedIPs) fetch() ([]godo.ReservedIP, error) {
	reservedIPsList := []godo.ReservedIP{}
	opts := &godo.ListOptions{}

	for {
		doReservedIPs, resp, err := r.doClient.ReservedIPs.List(context.Background(), opts)
		if err != nil {
			return reservedIPsList, err
		}

		reservedIPsList = append(reservedIPsList, doReservedIPs...)

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
