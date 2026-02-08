package scheduler

import "crispy/service"

// checks and generates pending recurring transactions
type Manager struct {
	service *service.Service
}

func NewManager(service *service.Service) *Manager {
	return &Manager{
		service,
	}
}
