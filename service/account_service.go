package service

import (
	"context"
	"crispy/domain"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

const (
	ROOT_PARENT_ID = 0
)

func (s *Service) CreateAccount(ctx context.Context, newAccount *domain.Account) (*domain.Account, error) {
	if err := newAccount.Validate(); err != nil {
		return nil, err
	}
	err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error starting tx in CreateAccount:\n%s", err)
	}

	i := strings.LastIndex(newAccount.Name(), ":")
	if i == -1 {
		// no colon, parent is root
		newAccount.SetParentID(ROOT_PARENT_ID)
	} else {
		parentName := newAccount.Name()[:i]
		p_id, err := s.GetID(ctx, parentName)
		if errors.Is(err, sql.ErrNoRows) {
			// Handle parent does not exist
			return nil, err
		} else if err != nil {
			return nil, err
		}
		newAccount.SetParentID(p_id)
		newAccount.SetName(newAccount.Name()[i+1:])
	}

	lastID, err := s.repo.CreateAccount(ctx, newAccount)
	if err != nil {
		s.repo.Rollback()
		return nil, err
	}

	if err := s.repo.Commit(); err != nil {
		s.repo.Rollback()
		return nil, fmt.Errorf("Error committing CreateAccount:\n%s", err)
	}
	return s.repo.GetAccountById(ctx, lastID)
}

func (s *Service) GetID(ctx context.Context, accountName string) (int64, error) {
	acct, err := s.repo.GetAccountByFullName(ctx, accountName)
	if err != nil {
		return -9, err
	}
	return acct.ID(), err
}

func (s *Service) GetAccountName(ctx context.Context, id int64) (string, error) {
	acct, err := s.repo.GetAccountById(ctx, id)
	if err != nil {
		return "", err
	}
	return acct.Name(), nil
}
