package domodules

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/senorprogrammer/dosage/modules"
	"github.com/senorprogrammer/dosage/pieces"
)

// SSHKeys is SSH keys
type SSHKeys struct {
	modules.Base
	SSHKeys      []godo.Key
	PositionData pieces.PositionData
	doClient     *godo.Client
}

// NewSSHKeys creates and returns an instance of SSHKeys
func NewSSHKeys(title string, client *godo.Client) *SSHKeys {
	mod := &SSHKeys{
		Base: modules.NewBase(title),
		PositionData: pieces.PositionData{
			Row:       2,
			Col:       11,
			RowSpan:   2,
			ColSpan:   4,
			MinHeight: 0,
			MinWidth:  0,
		},
		SSHKeys:  []godo.Key{},
		doClient: client,
	}

	mod.Enabled = true

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

	s.SetAvailable(false)

	sshKeys, err := s.fetch()
	if err != nil {
		s.LastError = err
	} else {
		s.LastError = nil
		s.SSHKeys = sshKeys
	}

	s.SetAvailable(true)
}

// Render draws the current string representation into the view
func (s *SSHKeys) Render() {
	str := s.ToStr()
	s.GetView().SetText(str)
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

		for _, doKey := range doKeys {
			keysList = append(keysList, doKey)
		}

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
