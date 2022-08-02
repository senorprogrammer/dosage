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

// Certificates is certificates
type Certificates struct {
	modules.Base
	Certificates []godo.Certificate
	doClient     *godo.Client
}

// NewCertificates creates and returns an instance of Certificates
func NewCertificates(title string, refreshChan chan bool, client *godo.Client, logger *modules.Logger) *Certificates {
	mod := &Certificates{
		Base:         modules.NewBase(title, modules.WithTextView, refreshChan, 5*time.Second, logger),
		Certificates: []godo.Certificate{},
		doClient:     client,
	}

	mod.Enabled = true

	mod.PositionData = pieces.PositionData{
		Row:       0,
		Col:       11,
		RowSpan:   2,
		ColSpan:   4,
		MinHeight: 0,
		MinWidth:  0,
	}

	mod.RefreshFunc = mod.Refresh

	return mod
}

/* -------------------- Exported Functions -------------------- */

// GetPositionData returns PositionData
func (c *Certificates) GetPositionData() *pieces.PositionData {
	return &c.PositionData
}

// Refresh updates the view content with the latest data
func (c *Certificates) Refresh() {
	if !c.GetAvailable() || !c.GetEnabled() {
		return
	}

	c.Logger.Log(fmt.Sprintf("refreshing %s", c.GetTitle()))

	c.SetAvailable(false)

	certs, err := c.fetch()
	if err != nil {
		c.LastError = err
	} else {
		c.LastError = nil
		c.Certificates = certs
	}

	c.SetAvailable(true)

	c.Render()

	// Tell the Refresher that there's new data to display
	c.RefreshChan <- true
}

// Render draws the current string representation into the view
func (c *Certificates) Render() {
	str := c.ToStr()
	c.GetView().(*tview.TextView).SetText(str)
}

// ToStr returns a string representation of the module suitable for display onscreen
func (c *Certificates) ToStr() string {
	if c.LastError != nil {
		return c.LastError.Error()
	}

	if len(c.Certificates) == 0 {
		return modules.EmptyContentLabel
	}

	str := ""

	for _, cert := range c.Certificates {
		str = str + fmt.Sprintf(
			"%s\t%s\t%s\t%s\n",
			cert.ID,
			cert.Name,
			cert.Type,
			cert.State,
		)
	}

	return str
}

/* -------------------- Unexported Functions -------------------- */

func (c *Certificates) fetch() ([]godo.Certificate, error) {
	certsList := []godo.Certificate{}
	opts := &godo.ListOptions{}

	for {
		doCerts, resp, err := c.doClient.Certificates.List(context.Background(), opts)
		if err != nil {
			return certsList, err
		}

		for _, doCert := range doCerts {
			certsList = append(certsList, doCert)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return certsList, err
		}

		// Set the page we want for the next request
		opts.Page = page + 1
	}

	return certsList, nil
}
