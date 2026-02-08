package automation

import "crispy/service"

type Engine struct {
	service *service.Service
}

func NewEngine(service *service.Service) *Engine {
	return &Engine{
		service,
	}
}
