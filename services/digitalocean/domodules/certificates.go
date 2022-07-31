package domodules

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/senorprogrammer/dosage/modules"
	"github.com/senorprogrammer/dosage/pieces"
)

type Certificates struct {
	modules.Base
	PositionData pieces.PositionData
	doClient     *godo.Client
}

// NewCertificates creates and returns an instance of Certificates
func NewCertificates(title string, client *godo.Client) *Certificates {
	mod := &Certificates{
		Base: modules.NewBase(title),
		PositionData: pieces.PositionData{
			Row:       0,
			Col:       11,
			RowSpan:   2,
			ColSpan:   4,
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
func (c *Certificates) GetPositionData() *pieces.PositionData {
	return &c.PositionData
}

// Refresh updates the view content with the latest data
func (c *Certificates) Refresh() {
	if !c.GetAvailable() || !c.GetEnabled() {
		return
	}

	c.SetAvailable(false)
	c.GetView().SetText(c.data())
	c.SetAvailable(true)
}

/* -------------------- Unexported Functions -------------------- */

func (c *Certificates) data() string {
	certs, err := c.fetch()
	if err != nil {
		return err.Error()
	}

	if len(certs) == 0 {
		return modules.EmptyDataLabel
	}

	data := ""

	for idx, cert := range certs {
		data = data + fmt.Sprintf(
			"%3d\t%s\t%s\t%s\t%s\n",
			(idx+1),
			cert.ID,
			cert.Name,
			cert.Type,
			cert.State,
		)
	}

	return data
}

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
