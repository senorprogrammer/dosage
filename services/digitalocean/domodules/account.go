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

// Account is account
type Account struct {
	modules.Base
	AccountInfo *godo.Account
	doClient    *godo.Client
}

// NewAccount creates and returns an instance of Account
func NewAccount(title string, refreshChan chan bool, client *godo.Client, logger *modules.Logger) *Account {
	mod := &Account{
		Base:        modules.NewBase(title, modules.WithTableView, refreshChan, 15*time.Second, logger),
		AccountInfo: nil,
		doClient:    client,
	}

	mod.Enabled = true

	mod.PositionData = pieces.PositionData{
		Row:       0,
		Col:       0,
		RowSpan:   2,
		ColSpan:   2,
		MinHeight: 0,
		MinWidth:  0,
	}

	mod.RefreshFunc = mod.Refresh

	return mod
}

/* -------------------- Exported Functions -------------------- */

// GetPositionData returns PositionData
func (a *Account) GetPositionData() *pieces.PositionData {
	return &a.PositionData
}

// Refresh updates the view content with the latest data
func (a *Account) Refresh() {
	if !a.GetAvailable() || !a.GetEnabled() {
		return
	}

	a.Logger.Log(fmt.Sprintf("refreshing %s", a.GetTitle()))

	a.SetAvailable(false)

	accountInfo, err := a.fetch()
	if err != nil {
		a.LastError = err
	} else {
		a.LastError = nil
		a.AccountInfo = accountInfo
	}

	a.SetAvailable(true)

	a.Render()

	// Tell the Refresher that there's new data to display
	a.RefreshChan <- true
}

// Render draws the current string representation into the view
func (a *Account) Render() {
	statusText := pieces.ColorForState(a.AccountInfo.Status, a.AccountInfo.Status)
	limitsLabel := pieces.Bold(pieces.Green("Limits"))

	// Create the table
	table := a.GetView().(*tview.Table)
	table.Clear()

	table.SetCell(0, 0, tview.NewTableCell("Status:").SetAlign(tview.AlignRight))
	table.SetCell(0, 1, tview.NewTableCell(fmt.Sprint(statusText)))

	table.SetCell(1, 0, tview.NewTableCell("Team:").SetAlign(tview.AlignRight))
	table.SetCell(1, 1, tview.NewTableCell(a.AccountInfo.Team.Name))

	table.SetCell(2, 0, tview.NewTableCell(""))
	table.SetCell(3, 0, tview.NewTableCell(limitsLabel))

	table.SetCell(4, 0, tview.NewTableCell("Droplets:").SetAlign(tview.AlignRight))
	table.SetCell(4, 1, tview.NewTableCell(fmt.Sprint(a.AccountInfo.DropletLimit)))

	table.SetCell(5, 0, tview.NewTableCell("Reserved IPs:").SetAlign(tview.AlignRight))
	table.SetCell(5, 1, tview.NewTableCell(fmt.Sprint(a.AccountInfo.ReservedIPLimit)))

	table.SetCell(6, 0, tview.NewTableCell("Volumes:").SetAlign(tview.AlignRight))
	table.SetCell(6, 1, tview.NewTableCell(fmt.Sprint(a.AccountInfo.VolumeLimit)))
}

// ToStr returns a string representation of the module suitable for display onscreen
func (a *Account) ToStr() string {
	if a.LastError != nil {
		return a.LastError.Error()
	}

	if a.AccountInfo == nil {
		return modules.EmptyContentLabel
	}

	str := ""
	str += fmt.Sprintf("Status: %s %s\n", pieces.ColorForState(a.AccountInfo.Status, a.AccountInfo.Status), a.AccountInfo.StatusMessage)
	str += fmt.Sprintf("Team: %s\n", a.AccountInfo.Team.Name)
	str += "\n"
	str += fmt.Sprintf("[green:]%s[white:]\n", "Limits")
	str += fmt.Sprintf("%12s: %d\n", "Droplets", a.AccountInfo.DropletLimit)
	str += fmt.Sprintf("%12s: %d\n", "Reserved IPs", a.AccountInfo.ReservedIPLimit)
	str += fmt.Sprintf("%12s: %d\n", "Volumes", a.AccountInfo.VolumeLimit)

	return str
}

/* -------------------- Unexported Functions -------------------- */

func (a *Account) fetch() (*godo.Account, error) {
	acct, _, err := a.doClient.Account.Get(context.Background())
	if err != nil {
		return nil, err
	}

	return acct, nil
}
