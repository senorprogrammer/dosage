package services

import "github.com/senorprogrammer/dosage/modules"

// Service defines the interface for external services
type Service interface {
	GetName() string

	LoadModules(logger *modules.Logger)
	Refresh()
}
