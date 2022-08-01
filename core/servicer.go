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

// GetService returns the service for the given name
func (s *Servicer) GetService(name services.ServiceName) services.Serviceable {
	return s.services[name]
}

// GetServices returns a slice of the loaded services
func (s *Servicer) GetServices() []services.Serviceable {
	svcs := []services.Serviceable{}

	for _, svc := range s.services {
		svcs = append(svcs, svc)
	}

	return svcs
}

// LoadServices creates and stores instances of all the supported services
func (s *Servicer) LoadServices(flags *flags.Flags, tviewPages *tview.Pages, logger *modules.Logger) []services.Serviceable {
	digitalOcean := digitalocean.NewDigitalOcean(flags.APIKey, tviewPages, logger)
	digitalOcean.LoadModules()

	s.services[digitalOcean.GetName()] = digitalOcean

	return s.GetServices()
}
