package modules

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
)

// Account is account
type Account struct {
	Base
	doClient *godo.Client
}

// NewAccount creates and returns an instance of Account
func NewAccount(title string, client *godo.Client) *Account {
	return &Account{
		Base:     NewBase(title),
		doClient: client,
	}
}

/* -------------------- Exported Functions -------------------- */

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
	data += fmt.Sprintf("Droplet limit: %d\n", accountInfo.DropletLimit)
	data += fmt.Sprintf("Reserved IP limit: %d\n", accountInfo.ReservedIPLimit)
	data += fmt.Sprintf("Volume limit: %d\n", accountInfo.VolumeLimit)

	return data
}

func (a *Account) fetch() (*godo.Account, error) {
	acct, _, err := a.doClient.Account.Get(context.Background())
	if err != nil {
		return nil, err
	}

	return acct, nil
}
