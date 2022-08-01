package domodules

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/senorprogrammer/dosage/modules"
	"github.com/senorprogrammer/dosage/pieces"
)

// Account is account
type Account struct {
	modules.Base

	AccountInfo  *godo.Account
	PositionData pieces.PositionData
	doClient     *godo.Client
}

// NewAccount creates and returns an instance of Account
func NewAccount(title string, client *godo.Client) *Account {
	mod := &Account{
		Base:        modules.NewBase(title),
		AccountInfo: &godo.Account{},
		PositionData: pieces.PositionData{
			Row:       0,
			Col:       0,
			RowSpan:   2,
			ColSpan:   2,
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
func (a *Account) GetPositionData() *pieces.PositionData {
	return &a.PositionData
}

// Refresh updates the view content with the latest data
func (a *Account) Refresh() {
	if !a.GetAvailable() || !a.GetEnabled() {
		return
	}

	a.SetAvailable(false)

	accountInfo, err := a.fetch()
	if err != nil {
		a.LastError = err
	} else {
		a.LastError = nil
		a.AccountInfo = accountInfo
	}

	a.SetAvailable(true)
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
