package service

import (
	"context"
	"crispy/domain"
	"fmt"
	"time"
)

func (s *Service) CreateTransaction(ctx context.Context, newTransaction *domain.Transaction) (*domain.Transaction, error) {
	if err := newTransaction.Validate(); err != nil {
		return nil, err
	}
	err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error starting tx in CreateTransaction:\n%s", err)
	}
	lastID, err := s.repo.CreateTransaction(ctx, newTransaction)
	if err != nil {
		s.repo.Rollback()
		return nil, fmt.Errorf("Error creating transaction:\n%s", err)
	}
	for _, p := range newTransaction.Postings() {
		p.SetTransactionID(lastID)
		_, err := s.repo.CreatePosting(ctx, p)
		if err != nil {
			s.repo.Rollback()
			return nil, fmt.Errorf("Error creating posting %v:\n%s", p, err)
		}
	}
	if err := s.repo.Commit(); err != nil {
		s.repo.Rollback()
		return nil, fmt.Errorf("Error committing CreateTransaction:\n%s", err)
	}
	return s.repo.GetTransactionById(ctx, lastID)
}

func (s *Service) GetAllTransactions(ctx context.Context) ([]*domain.Transaction, error) {
	t_list, err := s.repo.GetTransactionByDate(ctx, time.Date(1970, time.January, 1, 0, 0, 0, 0, time.Local), time.Now())
	if err != nil {
		return nil, err
	}
	return t_list, nil
}

func (s *Service) GetTransactionByID(ctx context.Context, id int64) (*domain.Transaction, error) {
	tx, err := s.repo.GetTransactionById(ctx, id)
	if err != nil {
		return nil, err
	}
	return tx, err
}

func (s *Service) UpdateTransaction(ctx context.Context, newTransaction *domain.Transaction) (*domain.Transaction, error) {
	if err := newTransaction.Validate(); err != nil {
		return nil, err
	}
	err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error starting tx in UpdateTransaction:\n%s", err)
	}
	lastID, err := s.repo.UpdateTransaction(ctx, newTransaction)
	if err != nil {
		s.repo.Rollback()
		return nil, fmt.Errorf("Error updating transaction:\n%s", err)
	}
	for _, p := range newTransaction.Postings() {
		p.SetTransactionID(lastID)
		err := s.repo.UpdatePosting(ctx, p)
		if err != nil {
			s.repo.Rollback()
			return nil, fmt.Errorf("Error updating posting %v:\n%s", p, err)
		}
	}
	if err := s.repo.Commit(); err != nil {
		s.repo.Rollback()
		return nil, fmt.Errorf("Error committing CreateTransaction:\n%s", err)
	}
	return s.repo.GetTransactionById(ctx, lastID)
}

func (s *Service) DeletePosting(ctx context.Context, postingID int64) error {
	return fmt.Errorf("TODO")
}
