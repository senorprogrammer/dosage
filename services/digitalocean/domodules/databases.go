package domodules

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/rivo/tview"
	"github.com/senorprogrammer/dosage/formatting"
	"github.com/senorprogrammer/dosage/modules"
	"github.com/senorprogrammer/dosage/pieces"
)

// Databases is database
type Databases struct {
	modules.Base
	Databases []godo.Database
	doClient  *godo.Client
}

// NewDatabases creates and returns an instance of Databases
func NewDatabases(title string, refreshChan chan bool, client *godo.Client, logger *modules.Logger) *Databases {
	mod := &Databases{
		Base:      modules.NewBase(title, modules.WithTableView, refreshChan, modules.DefaultRefreshSeconds*time.Second, logger),
		Databases: []godo.Database{},
		doClient:  client,
	}

	mod.Enabled = true

	mod.PositionData = pieces.PositionData{
		Row:       4,
		Col:       3,
		RowSpan:   2,
		ColSpan:   7,
		MinHeight: 0,
		MinWidth:  0,
	}

	mod.RefreshFunc = mod.Refresh

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

	d.Logger.Log(fmt.Sprintf("refreshing %s", d.GetTitle()))

	d.SetAvailable(false)

	databases, err := d.fetch()
	if err != nil {
		d.LastError = err
	} else {
		d.LastError = nil
		d.Databases = databases
	}

	d.SetAvailable(true)

	d.Render()

	// Tell the Refresher that there's new data to display
	d.RefreshChan <- true
}

// Render draws the current string representation into the view
func (d *Databases) Render() {
	table := d.GetView().(*tview.Table)
	table.Clear()

	if d.Databases == nil {
		table.SetCell(0, 0, tview.NewTableCell("Databases are nil"))
		return
	}

	if d.LastError != nil {
		table.SetCell(0, 0, tview.NewTableCell(d.LastError.Error()))
		return
	}

	for idx, header := range []string{"Name", "Size", "Status", "Engine", "Tags"} {
		table.SetCell(0, idx, tview.NewTableCell(formatting.Bold(formatting.Underline(header))).SetAlign(tview.AlignCenter))
	}

	for idx, database := range d.Databases {
		row := idx + 1
		table.SetCell(row, 0, tview.NewTableCell(database.Name))
		table.SetCell(row, 1, tview.NewTableCell(database.SizeSlug))
		table.SetCell(row, 2, tview.NewTableCell(formatting.ColorForState(database.Status, database.Status)))
		table.SetCell(row, 3, tview.NewTableCell(database.EngineSlug))
		table.SetCell(row, 4, tview.NewTableCell(strings.Join(database.Tags, ", ")))
	}
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

		databaseList = append(databaseList, doDbs...)

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
