package core

import (
	"github.com/rivo/tview"
	"github.com/senorprogrammer/dosage/flags"
	"github.com/senorprogrammer/dosage/modules"
	"github.com/senorprogrammer/dosage/services"
	"github.com/senorprogrammer/dosage/services/digitalocean"
)

// Servicer manages the services
type Servicer struct {
	Services map[services.ServiceName]services.Service
}

// NewServicer creates and returns an instance of Servicer
func NewServicer() *Servicer {
	s := &Servicer{
		Services: make(map[services.ServiceName]services.Service),
	}

	return s
}

/* -------------------- Exported Functions -------------------- */

// GetServices returns a slice of the loaded services
func (s *Servicer) GetServices() map[services.ServiceName]services.Service {
	return s.Services
}

// LoadServices creates and stores instances of all the supported services
func (s *Servicer) LoadServices(flags *flags.Flags, tviewPages *tview.Pages, logger *modules.Logger) map[services.ServiceName]services.Service {
	digitalOcean := digitalocean.NewDigitalOcean(flags.APIKey, tviewPages, logger)
	digitalOcean.LoadModules()

	s.Services[digitalOcean.GetName()] = digitalOcean

	return s.Services
}
