package service

import "crispy/db/repository"

type serviceDeps interface {
	repository.TransactionRepository
	repository.PostingRepository
	repository.AccountRepository
	repository.TxRepo
}

type Service struct {
	repo serviceDeps
}

func NewService(repo serviceDeps) *Service {
	return &Service{repo: repo}
}
