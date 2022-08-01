package domodules

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/senorprogrammer/dosage/modules"
	"github.com/senorprogrammer/dosage/pieces"
)

// Volumes is storage
type Volumes struct {
	modules.Base
	PositionData pieces.PositionData
	Volumes      []godo.Volume
	doClient     *godo.Client
}

// NewVolumes creates and returns an instance of Storage
func NewVolumes(title string, client *godo.Client) *Volumes {
	mod := &Volumes{
		Base: modules.NewBase(title),
		PositionData: pieces.PositionData{
			Row:       0,
			Col:       7,
			RowSpan:   2,
			ColSpan:   4,
			MinHeight: 0,
			MinWidth:  0,
		},
		Volumes:  []godo.Volume{},
		doClient: client,
	}

	mod.Enabled = true

	return mod
}

/* -------------------- Exported Functions -------------------- */

// GetPositionData returns PositionData
func (v *Volumes) GetPositionData() *pieces.PositionData {
	return &v.PositionData
}

// Refresh updates the view content with the latest data
func (v *Volumes) Refresh() {
	if !v.GetAvailable() || !v.GetEnabled() {
		return
	}

	v.SetAvailable(false)

	vols, err := v.fetch()
	if err != nil {
		v.LastError = err
	} else {
		v.LastError = nil
		v.Volumes = vols
	}

	v.SetAvailable(true)
}

// Render draws the current string representation into the view
func (v *Volumes) Render() {
	str := v.ToStr()
	v.GetView().SetText(str)
}

// ToStr returns a string representation of the module suitable for display onscreen
func (v *Volumes) ToStr() string {
	if v.LastError != nil {
		return v.LastError.Error()
	}

	if len(v.Volumes) == 0 {
		return modules.EmptyContentLabel
	}

	str := ""

	for idx, vol := range v.Volumes {
		str = str + fmt.Sprintf(
			"%3d\t%s\t%d\t%s\t%s\n",
			(idx+1),
			vol.Name,
			vol.SizeGigaBytes,
			vol.Description,
			vol.Region.Slug,
		)
	}

	return str
}

/* -------------------- Unexported Functions -------------------- */

// fetch uses the DigitalOcean API to fetch information about all the available droplets
func (v *Volumes) fetch() ([]godo.Volume, error) {
	volumesList := []godo.Volume{}
	opts := &godo.ListVolumeParams{}

	doVols, _, err := v.doClient.Storage.ListVolumes(context.Background(), opts)
	if err != nil {
		return volumesList, err
	}

	for _, doVol := range doVols {
		volumesList = append(volumesList, doVol)
	}

	return volumesList, nil
}
