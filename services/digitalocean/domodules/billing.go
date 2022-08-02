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

// Billing is billing
type Billing struct {
	modules.Base
	BillingHistory []godo.BillingHistoryEntry
	doClient       *godo.Client
}

// NewBilling creates and returns an instance of Billing
func NewBilling(title string, refreshChan chan bool, client *godo.Client, logger *modules.Logger) *Billing {
	mod := &Billing{
		Base:           modules.NewBase(title, modules.WithTextView, refreshChan, 5*time.Second, logger),
		BillingHistory: []godo.BillingHistoryEntry{},
		doClient:       client,
	}

	mod.Enabled = true

	mod.PositionData = pieces.PositionData{
		Row:       2,
		Col:       7,
		RowSpan:   2,
		ColSpan:   4,
		MinHeight: 0,
		MinWidth:  0,
	}

	mod.RefreshFunc = mod.Refresh

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

	b.Logger.Log(fmt.Sprintf("refreshing %s", b.GetTitle()))

	b.SetAvailable(false)

	billingHistory, err := b.fetch()
	if err != nil {
		b.LastError = err
	} else {
		b.LastError = nil
		b.BillingHistory = billingHistory
	}

	b.SetAvailable(true)

	b.Render()

	// Tell the Refresher that there's new data to display
	b.RefreshChan <- true
}

// Render draws the current string representation into the view
func (b *Billing) Render() {
	str := b.ToStr()
	b.GetView().(*tview.TextView).SetText(str)
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
		str += fmt.Sprintf(
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
