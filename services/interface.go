package services

// Serviceable defines the interface for external services
type Serviceable interface {
	GetName() ServiceName
	LoadModules(refreshChan chan bool)
}
