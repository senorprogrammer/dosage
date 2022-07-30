package services

// Service defines the interface for external services
type Service interface {
	GetName() string

	LoadModules()
	Refresh()
}
