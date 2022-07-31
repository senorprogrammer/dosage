package domodules

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/senorprogrammer/dosage/modules"
	"github.com/senorprogrammer/dosage/pieces"
)

// Databases is database
type Databases struct {
	modules.Base
	PositionData pieces.PositionData
	doClient     *godo.Client
}

// NewDatabases creates and returns an instance of Databases
func NewDatabases(title string, client *godo.Client) *Databases {
	return &Databases{
		Base: modules.NewBase(title),
		PositionData: pieces.PositionData{
			Row:       4,
			Col:       2,
			RowSpan:   2,
			ColSpan:   9,
			MinHeight: 0,
			MinWidth:  0,
		},
		doClient: client,
	}
}

/* -------------------- Exported Functions -------------------- */

// GetPositionData returns PositionData
func (d *Databases) GetPositionData() *pieces.PositionData {
	return &d.PositionData
}

// Refresh updates the view content with the latest data
func (d *Databases) Refresh() {
	if !d.Available {
		return
	}

	d.SetAvailable(false)
	d.GetView().SetText(d.data())
	d.SetAvailable(true)
}

/* -------------------- Unexported Functions -------------------- */

func (d *Databases) data() string {
	databases, err := d.fetch()
	if err != nil {
		return err.Error()
	}

	if len(databases) == 0 {
		return "none"
	}

	data := ""

	for idx, database := range databases {
		data = data + fmt.Sprintf(
			"%3d\t%s\t%s\t%s\t%s\t%s\t%v\n",
			(idx+1),
			database.Name,
			database.EngineSlug,
			pieces.ColorForState(database.Status, database.Status),
			database.SizeSlug,
			database.RegionSlug,
			database.Tags,
		)
	}

	return data
}

// fetch uses the DigitalOcean API to fetch information about all the available droplets
func (d *Databases) fetch() ([]godo.Database, error) {
	databaseList := []godo.Database{}
	opts := &godo.ListOptions{}

	for {
		doDbs, resp, err := d.doClient.Databases.List(context.Background(), opts)
		if err != nil {
			return databaseList, err
		}

		for _, doDb := range doDbs {
			databaseList = append(databaseList, doDb)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return databaseList, err
		}

		// Set the page we want for the next request
		opts.Page = page + 1
	}

	return databaseList, nil
}