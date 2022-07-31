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
	PositionData pieces.PositionData
	doClient     *godo.Client
}

// NewBilling creates and returns an instance of Billing
func NewBilling(title string, client *godo.Client) *Billing {
	return &Billing{
		Base: modules.NewBase(title),
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
}

/* -------------------- Exported Functions -------------------- */

// GetPositionData returns PositionData
func (b *Billing) GetPositionData() *pieces.PositionData {
	return &b.PositionData
}

// Refresh updates the view content with the latest data
func (b *Billing) Refresh() {
	if !b.Available {
		return
	}

	b.SetAvailable(false)
	b.GetView().SetText(b.data())
	b.SetAvailable(true)
}

/* -------------------- Unexported Functions -------------------- */

func (b *Billing) data() string {
	billingHistory, err := b.fetch()
	if err != nil {
		return err.Error()
	}

	if len(billingHistory) == 0 {
		return "none"
	}

	data := ""

	for idx, bhe := range billingHistory {
		data = fmt.Sprintf(
			"%3d\t%s\t%v\t%v\t%8s\n",
			(idx+1),
			*bhe.InvoiceID,
			bhe.Date,
			bhe.Amount,
			bhe.Type,
		) + data
	}

	return data
}

func (b *Billing) fetch() ([]godo.BillingHistoryEntry, error) {
	billingHistory := []godo.BillingHistoryEntry{}
	opts := &godo.ListOptions{}

	bh, _, err := b.doClient.BillingHistory.List(context.Background(), opts)
	if err != nil {
		return billingHistory, err
	}

	return bh.BillingHistory, nil
}
