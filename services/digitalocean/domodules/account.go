package domodules

import (
	"context"
	"fmt"
	"time"

	"github.com/digitalocean/godo"
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
func NewAccount(title string, client *godo.Client, logger *modules.Logger) *Account {
	mod := &Account{
		Base:        modules.NewBase(title, 15*time.Second, logger),
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
}

// Render draws the current string representation into the view
func (a *Account) Render() {
	str := a.ToStr()
	a.GetView().SetText(str)
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
