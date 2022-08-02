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

// SSHKeys is SSH keys
type SSHKeys struct {
	modules.Base
	SSHKeys  []godo.Key
	doClient *godo.Client
}

// NewSSHKeys creates and returns an instance of SSHKeys
func NewSSHKeys(title string, refreshChan chan bool, client *godo.Client, logger *modules.Logger) *SSHKeys {
	mod := &SSHKeys{
		Base:     modules.NewBase(title, modules.WithTextView, refreshChan, 5*time.Second, logger),
		SSHKeys:  []godo.Key{},
		doClient: client,
	}

	mod.Enabled = true

	mod.PositionData = pieces.PositionData{
		Row:       2,
		Col:       11,
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
func (s *SSHKeys) GetPositionData() *pieces.PositionData {
	return &s.PositionData
}

// Refresh updates the view content with the latest data
func (s *SSHKeys) Refresh() {
	if !s.GetAvailable() || !s.GetEnabled() {
		return
	}

	s.Logger.Log(fmt.Sprintf("refreshing %s", s.GetTitle()))

	s.SetAvailable(false)

	sshKeys, err := s.fetch()
	if err != nil {
		s.LastError = err
	} else {
		s.LastError = nil
		s.SSHKeys = sshKeys
	}

	s.SetAvailable(true)

	s.Render()

	// Tell the Refresher that there's new data to display
	s.RefreshChan <- true
}

// Render draws the current string representation into the view
func (s *SSHKeys) Render() {
	str := s.ToStr()
	s.GetView().(*tview.TextView).SetText(str)
}

// ToStr returns a string representation of the module suitable for display onscreen
func (s *SSHKeys) ToStr() string {
	if s.LastError != nil {
		return s.LastError.Error()
	}

	if len(s.SSHKeys) == 0 {
		return modules.EmptyContentLabel
	}

	str := ""

	for _, key := range s.SSHKeys {
		str = str + fmt.Sprintf(
			"%d\t%s\t%s\n",
			key.ID,
			key.Name,
			key.Fingerprint,
		)
	}

	return str
}

/* -------------------- Unexported Functions -------------------- */

func (s *SSHKeys) fetch() ([]godo.Key, error) {
	keysList := []godo.Key{}
	opts := &godo.ListOptions{}

	for {
		doKeys, resp, err := s.doClient.Keys.List(context.Background(), opts)
		if err != nil {
			return keysList, err
		}

		keysList = append(keysList, doKeys...)

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return keysList, err
		}

		// Set the page we want for the next request
		opts.Page = page + 1
	}

	return keysList, nil
}
