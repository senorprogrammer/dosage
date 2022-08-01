package services

// ServiceName is a specialized string that names a service
type ServiceName string

// Base is base
type Base struct {
	Name ServiceName
}

// GetName returns the name of the service
func (b *Base) GetName() ServiceName {
	return b.Name
}
