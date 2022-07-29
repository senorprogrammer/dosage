package modules

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/rivo/tview"
)

// Databases is databases
type Databases struct {
	Focus bool
	Title string
	View  *tview.TextView

	doClient *godo.Client
}

// NewDatabases creates and returns an instance of Databases
func NewDatabases(title string, client *godo.Client) *Databases {
	view := tview.NewTextView()
	view.SetTitle(title)
	view.SetWrap(false)
	view.SetBorder(true)
	view.SetScrollable(true)

	return &Databases{
		Focus: false,
		View:  view,

		doClient: client,
	}
}

/* -------------------- Exported Functions -------------------- */

// GetFocus returns the focus val for display
func (d *Databases) GetFocus() bool {
	return d.Focus
}

// GetView returns the tview.TextView used to display this module's data
func (d *Databases) GetView() *tview.TextView {
	return d.View
}

// Refresh updates the view content with the latest data
func (d *Databases) Refresh() {
	d.GetView().SetText(d.data())
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
			"%3d\t%s\t%s\t%s\t%s\t%s\n",
			(idx+1),
			database.Name,
			database.EngineSlug,
			database.Status,
			database.SizeSlug,
			database.RegionSlug,
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
