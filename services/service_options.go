package services

import (
	"github.com/senorprogrammer/dosage/modules"
)

// ServiceOptions are used to pass common system-level options to the services
type ServiceOptions struct {
	Logger      *modules.Logger
	RefreshChan chan bool
}

// NewServiceOptions creates and returns an instance of ServiceOptions
func NewServiceOptions(logger *modules.Logger, refreshChan chan bool) *ServiceOptions {
	return &ServiceOptions{
		Logger:      logger,
		RefreshChan: refreshChan,
	}
}
