package services

import (
	"github.com/rivo/tview"
	"github.com/senorprogrammer/dosage/modules"
)

type ServiceOptionable interface {
	GetLogger() *modules.Logger
	GetPages() *tview.Pages
	GetRefreshChan() chan bool
}

// ServiceOptions are used to pass common system-level options to the services
type ServiceOptions struct {
	Logger      *modules.Logger
	Pages       *tview.Pages
	RefreshChan chan bool
}

func (so *ServiceOptions) GetLogger() *modules.Logger {
	return so.Logger
}

func (so *ServiceOptions) GetPages() *tview.Pages {
	return so.Pages
}

func (so *ServiceOptions) GetRefreshChan() chan bool {
	return so.RefreshChan
}
