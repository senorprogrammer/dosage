package modules

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/senorprogrammer/dosage/pieces"
)

// Account is account
type Account struct {
	Base
	PositionData pieces.PositionData
	doClient     *godo.Client
}

// NewAccount creates and returns an instance of Account
func NewAccount(title string, client *godo.Client) *Account {
	return &Account{
		Base: NewBase(title),
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
}

/* -------------------- Exported Functions -------------------- */

// GetPositionData returns PositionData
func (a *Account) GetPositionData() *pieces.PositionData {
	return &a.PositionData
}

// Refresh updates the view content with the latest data
func (a *Account) Refresh() {
	a.GetView().SetText(a.data())
}

/* -------------------- Unexported Functions -------------------- */

// data returns a string representation of the module
// suitable for display onscreen
func (a *Account) data() string {
	accountInfo, err := a.fetch()
	if err != nil {
		return err.Error()
	}

	data := ""

	data += fmt.Sprintf("Status: %s %s\n", accountInfo.Status, accountInfo.StatusMessage)
	data += fmt.Sprintf("Team: %s\n", accountInfo.Team.Name)
	data += "\n"
	data += fmt.Sprintf("[green:]%s[white:]\n", "Limits")
	data += fmt.Sprintf("%12s: %d\n", "Droplets", accountInfo.DropletLimit)
	data += fmt.Sprintf("%12s: %d\n", "Reserved IPs", accountInfo.ReservedIPLimit)
	data += fmt.Sprintf("%12s: %d\n", "Volumes", accountInfo.VolumeLimit)

	return data
}

func (a *Account) fetch() (*godo.Account, error) {
	acct, _, err := a.doClient.Account.Get(context.Background())
	if err != nil {
		return nil, err
	}

	return acct, nil
}
