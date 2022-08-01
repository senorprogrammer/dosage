package domodules

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/senorprogrammer/dosage/modules"
	"github.com/senorprogrammer/dosage/pieces"
)

// Billing is billing
type Billing struct {
	modules.Base
	BillingHistory []godo.BillingHistoryEntry
	PositionData   pieces.PositionData
	doClient       *godo.Client
}

// NewBilling creates and returns an instance of Billing
func NewBilling(title string, client *godo.Client) *Billing {
	mod := &Billing{
		Base:           modules.NewBase(title),
		BillingHistory: []godo.BillingHistoryEntry{},
		PositionData: pieces.PositionData{
			Row:       2,
			Col:       7,
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
func (b *Billing) GetPositionData() *pieces.PositionData {
	return &b.PositionData
}

// Refresh updates the view content with the latest data
func (b *Billing) Refresh() {
	if !b.GetAvailable() || !b.GetEnabled() {
		return
	}

	b.SetAvailable(false)

	billingHistory, err := b.fetch()
	if err != nil {
		b.LastError = err
	} else {
		b.LastError = nil
		b.BillingHistory = billingHistory
	}

	b.SetAvailable(true)
}

// Render draws the current string representation into the view
func (b *Billing) Render() {
	str := b.ToStr()
	b.GetView().SetText(str)
}

// ToStr returns a string representation of the module suitable for display onscreen
func (b *Billing) ToStr() string {
	if b.LastError != nil {
		return b.LastError.Error()
	}

	if len(b.BillingHistory) == 0 {
		return modules.EmptyContentLabel
	}

	str := ""

	for _, bhe := range b.BillingHistory {
		str = str + fmt.Sprintf(
			"%s\t%v\t%v\t%8s\n",
			*bhe.InvoiceID,
			bhe.Date,
			bhe.Amount,
			bhe.Type,
		) + str
	}

	return str
}

/* -------------------- Unexported Functions -------------------- */

func (b *Billing) fetch() ([]godo.BillingHistoryEntry, error) {
	billingHistory := []godo.BillingHistoryEntry{}
	opts := &godo.ListOptions{}

	bh, _, err := b.doClient.BillingHistory.List(context.Background(), opts)
	if err != nil {
		return billingHistory, err
	}

	return bh.BillingHistory, nil
}
