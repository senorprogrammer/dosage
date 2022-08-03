package core

import (
	"github.com/rivo/tview"
	"github.com/senorprogrammer/dosage/flags"
	"github.com/senorprogrammer/dosage/services"
	"github.com/senorprogrammer/dosage/services/digitalocean"
)

// Servicer manages the services
type Servicer struct {
	services map[services.ServiceName]services.Serviceable
}

// NewServicer creates and returns an instance of Servicer
func NewServicer() *Servicer {
	s := &Servicer{
		services: make(map[services.ServiceName]services.Serviceable),
	}

	return s
}

/* -------------------- Exported Functions -------------------- */

// LoadServices creates and stores instances of the supported services
func (s *Servicer) LoadServices(flags *flags.Flags, tviewPages *tview.Pages, serviceOpts *services.ServiceOptions) {
	// Create the services
	s.services[digitalocean.ServiceName] = digitalocean.NewDigitalOcean(flags.APIKey, serviceOpts)

	// Add the services's Grids to the Pages
	for _, svc := range s.services {
		tviewPages.AddPage(string(svc.GetName()), svc.GetGrid(), true, true)
	}
}
