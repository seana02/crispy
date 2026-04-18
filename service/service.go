package service

import (
	"crispy/db"
	"time"
)

type serviceDeps interface {
	db.Repository
}

type Service struct {
	repo serviceDeps
}

func NewService(repo serviceDeps) *Service {
	return &Service{repo: repo}
}

// Returns the current time with sub-seconds truncated
func now() time.Time {
	return time.Now().Local().Truncate(time.Second)
}
