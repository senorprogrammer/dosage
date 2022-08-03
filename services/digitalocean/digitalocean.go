package digitalocean

import (
	"github.com/digitalocean/godo"
	"github.com/senorprogrammer/dosage/services"
	"github.com/senorprogrammer/dosage/services/digitalocean/domodules"
)

const (
	// GridTitle defines the title to display at the top of the tview.Page
	GridTitle = " DigitalOcean "

	// ServiceName defines the unique name for this service
	ServiceName = "digitalocean"
)

// DigitalOcean is the container application that handles all things DigitalOcean
type DigitalOcean struct {
	services.Base
	DOClient *godo.Client
}

// NewDigitalOcean creates and returns an instance of a DigitalOcean page container
func NewDigitalOcean(apiKey string, serviceOpts *services.ServiceOptions) *DigitalOcean {
	svc := &DigitalOcean{
		Base:     services.NewBase(ServiceName, GridTitle, serviceOpts.Logger),
		DOClient: godo.NewFromToken(apiKey),
	}

	svc.LoadModules(serviceOpts.RefreshChan)

	return svc
}

/* -------------------- Exported Functions -------------------- */

// LoadModules instantiates each module and attaches it to the TView app
// Pass the logger in because it's common across everything and needs to
// be instantiated before the rest of the modules
func (d *DigitalOcean) LoadModules(refreshChan chan bool) {
	d.Modules = append(d.Modules, d.Logger)
	d.Modules = append(d.Modules, domodules.NewAccount(" account ", refreshChan, d.DOClient, d.Logger))
	d.Modules = append(d.Modules, domodules.NewBilling(" billing ", refreshChan, d.DOClient, d.Logger))
	d.Modules = append(d.Modules, domodules.NewCertificates(" certs ", refreshChan, d.DOClient, d.Logger))
	d.Modules = append(d.Modules, domodules.NewDatabases(" databases ", refreshChan, d.DOClient, d.Logger))
	d.Modules = append(d.Modules, domodules.NewDroplets(" droplets ", refreshChan, d.DOClient, d.Logger))
	d.Modules = append(d.Modules, domodules.NewReservedIPs(" reserved ips ", refreshChan, d.DOClient, d.Logger))
	d.Modules = append(d.Modules, domodules.NewSSHKeys(" ssh keys ", refreshChan, d.DOClient, d.Logger))
	d.Modules = append(d.Modules, domodules.NewVolumes(" volumes ", refreshChan, d.DOClient, d.Logger))

	for _, mod := range d.Modules {
		if !mod.GetEnabled() {
			continue
		}

		// Enabled modules get add to the gridview
		d.Grid.AddItem(
			mod.GetView(),
			mod.GetPositionData().GetRow(),
			mod.GetPositionData().GetCol(),
			mod.GetPositionData().GetRowSpan(),
			mod.GetPositionData().GetColSpan(),
			mod.GetPositionData().GetMinWidth(),
			mod.GetPositionData().GetMinHeight(),
			false,
		)
	}

	// Start running the modules
	for _, mod := range d.Modules {
		if !mod.GetEnabled() {
			continue
		}

		mod.Run()
	}
}
