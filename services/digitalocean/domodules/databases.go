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
	Databases    []godo.Database
	PositionData pieces.PositionData
	doClient     *godo.Client
}

// NewDatabases creates and returns an instance of Databases
func NewDatabases(title string, client *godo.Client) *Databases {
	mod := &Databases{
		Base:      modules.NewBase(title),
		Databases: []godo.Database{},
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

	mod.Enabled = true

	return mod
}

/* -------------------- Exported Functions -------------------- */

// GetPositionData returns PositionData
func (d *Databases) GetPositionData() *pieces.PositionData {
	return &d.PositionData
}

// Refresh updates the view content with the latest data
func (d *Databases) Refresh() {
	if !d.GetAvailable() || !d.GetEnabled() {
		return
	}

	d.SetAvailable(false)

	databases, err := d.fetch()
	if err != nil {
		d.LastError = err
	} else {
		d.LastError = nil
		d.Databases = databases
	}

	d.SetAvailable(true)
}

// Render draws the current string representation into the view
func (d *Databases) Render() {
	str := d.ToStr()
	d.GetView().SetText(str)
}

// ToStr returns a string representation of the module suitable for display onscreen
func (d *Databases) ToStr() string {
	if d.LastError != nil {
		return d.LastError.Error()
	}

	if len(d.Databases) == 0 {
		return modules.EmptyContentLabel
	}

	str := ""

	for idx, database := range d.Databases {
		str = str + fmt.Sprintf(
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

	return str
}

/* -------------------- Unexported Functions -------------------- */

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
