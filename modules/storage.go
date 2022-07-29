package modules

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/senorprogrammer/dosage/pieces"
)

// Storage is storage
type Storage struct {
	Base
	PositionData pieces.PositionData
	doClient     *godo.Client
}

// NewStorage creates and returns an instance of Storage
func NewStorage(title string, client *godo.Client) *Storage {
	return &Storage{
		Base: NewBase(title),
		PositionData: pieces.PositionData{
			Row:       0,
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
func (s *Storage) GetPositionData() *pieces.PositionData {
	return &s.PositionData
}

// Refresh updates the view content with the latest data
func (s *Storage) Refresh() {
	s.GetView().SetText(s.data())
}

/* -------------------- Unexported Functions -------------------- */

func (s *Storage) data() string {
	volumes, err := s.fetch()
	if err != nil {
		return err.Error()
	}

	if len(volumes) == 0 {
		return "none"
	}

	data := ""

	for idx, vol := range volumes {
		data = data + fmt.Sprintf(
			"%3d\t%s\t%d\t%s\t%s\n",
			(idx+1),
			vol.Name,
			vol.SizeGigaBytes,
			vol.Description,
			vol.Region.Slug,
		)
	}

	return data
}

// fetch uses the DigitalOcean API to fetch information about all the available droplets
func (s *Storage) fetch() ([]godo.Volume, error) {
	volumesList := []godo.Volume{}
	opts := &godo.ListVolumeParams{}

	doVols, _, err := s.doClient.Storage.ListVolumes(context.Background(), opts)
	if err != nil {
		return volumesList, err
	}

	for _, doVol := range doVols {
		volumesList = append(volumesList, doVol)
	}

	return volumesList, nil
}
