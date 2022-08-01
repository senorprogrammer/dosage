package services

// Service defines the interface for external services
type Service interface {
	GetName() ServiceName

	LoadModules()
	Refresh()
}
