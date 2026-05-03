package service

import (
	"context"
	"crispy/db"
	"crispy/domain"
	"fmt"
	"strconv"
	"strings"
)

const (
	ROOT_PARENT_ID = int64(0)
)

func (s *Service) CreateAccount(ctx context.Context, newAccount *domain.Account) error {
	const errorMsg = "CreateAccount failed: %w"
	if err := newAccount.Validate(); err != nil {
		return fmt.Errorf(errorMsg, err)
	}
	err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	last := strings.LastIndex(newAccount.Name(), ":")
	if last == -1 {
		newAccount.SetParentID(ROOT_PARENT_ID)
	} else {
		parentID, err := s.GetAccountID(ctx, newAccount.Name()[:last])
		if err != nil {
			return fmt.Errorf(errorMsg, err)
		}
		newAccount.SetParentID(parentID)
		newAccount.SetName(newAccount.Name()[last+1:])
	}

	_, err = s.repo.AccountQueryBuilder().
		AddColumns([]db.Column{
			db.Column_ParentId,
			db.Column_Name,
			db.Column_Type,
			db.Column_Currency,
			db.Column_Description,
			db.Column_Active,
			db.Column_DateCreated,
			db.Column_DateUpdated,
		}).
		Insert(ctx, newAccount.ParentID(), newAccount.Name(), newAccount.Type(), newAccount.Currency(), newAccount.Description(), newAccount.Active(), now(), now())
	if err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, err)
	}
	if err := s.repo.Commit(); err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, err)
	}
	return nil
}

func (s *Service) GetAllAccounts(ctx context.Context) ([]*domain.Account, error) {
	const errorMsg = "GetAllAccounts failed: %w"
	acct_arr, err := s.repo.AccountQueryBuilder().
		SetCondition(
			db.NewWhere(db.Column_Id, db.NotEqual, "0"),
		).
		Select(ctx, -1, -1)
	if err != nil {
		return nil, fmt.Errorf(errorMsg, err)
	}
	return acct_arr, nil
}

func (s *Service) GetAccountByID(ctx context.Context, id int64) (*domain.Account, error) {
	const errorMsg = "GetAccountByID for ID %d failed: %w"
	acct_arr, err := s.repo.AccountQueryBuilder().
		SetCondition(db.AndCondition{
			Conditions: []db.Condition{
				db.NewWhere(db.Column_Id, db.Equal, strconv.FormatInt(id, 10)),
				db.NewWhere(db.Column_Id, db.NotEqual, "0"),
			},
		}).
		Select(ctx, 0, 1)
	if err != nil {
		return nil, fmt.Errorf(errorMsg, id, err)
	}
	return acct_arr[0], nil
}

func (s *Service) GetAccountID(ctx context.Context, accountName string) (int64, error) {
	const errorMsg = "GetAccountID for account %s failed: %w"
	accountTree := strings.Split(accountName, ":")
	curr_id := ROOT_PARENT_ID
	for _, next := range accountTree {
		acct_arr, err := s.repo.AccountQueryBuilder().
			SetCondition(db.AndCondition{
				Conditions: []db.Condition{
					db.NewWhere(db.Column_ParentId, db.Equal, strconv.FormatInt(curr_id, 10)),
					db.NewWhere(db.Column_Name, db.Equal, next),
				},
			}).
			Select(ctx, 0, 1)
		if err != nil || len(acct_arr) == 0 {
			return -9, fmt.Errorf(errorMsg, accountName, err)
		}
		curr_id = acct_arr[0].ID()
	}
	return curr_id, nil
}

func (s *Service) GetAccountName(ctx context.Context, id int64) (string, error) {
	const errorMsg = "GetAccountName for ID %d failed: %w"
	if id <= 0 {
		return "", fmt.Errorf(errorMsg, id, "invalid ID")
	}
	acct, err := s.GetAccountByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf(errorMsg, id, err)
	}
	return acct.Name(), nil
}

func (s *Service) SearchAccount(ctx context.Context, query string) ([]*domain.Account, error) {
	const errorMsg = "SearchAccount for query %s failed: %w"
	tx_arr, err := s.repo.AccountQueryBuilder().
		SetCondition(
			db.AndCondition{Conditions: []db.Condition{
				db.NewWhere(db.Column_Id, db.NotEqual, "0"),
				db.OrCondition{Conditions: []db.Condition{
					db.NewWhere(db.Column_Name, db.Like, query+"%"),
					db.NewWhere(db.Column_Name, db.Like, "%"+query+"%"),
					db.NewWhere(db.Column_Name, db.Like, "%"+strings.Join(strings.Split(query, ""), "%")+"%"),
				}},
			}},
		).
		Select(ctx, 0, 10)
	if err != nil {
		return nil, fmt.Errorf(errorMsg, query, err)
	}
	return append(tx_arr), nil
}

func (s *Service) UpdateAccount(ctx context.Context, newAccount *domain.Account) error {
	const errorMsg = "UpdateAccount failed: %w"
	if err := newAccount.Validate(); err != nil {
		return fmt.Errorf(errorMsg, err)
	}
	if err := s.repo.BeginTx(ctx); err != nil {
		return fmt.Errorf(errorMsg, err)
	}
	err := s.repo.AccountQueryBuilder().
		AddColumns([]db.Column{
			db.Column_ParentId,
			db.Column_Name,
			db.Column_Type,
			db.Column_Currency,
			db.Column_Description,
			db.Column_Active,
			db.Column_DateUpdated,
		}).
		SetCondition(
			db.NewWhere(db.Column_Id, db.Equal, strconv.FormatInt(newAccount.ID(), 10)),
		).
		Update(ctx, newAccount.ParentID(), newAccount.Name(), newAccount.Type(), newAccount.Currency(), newAccount.Description(), newAccount.Active(), now())
	if err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, err)
	}
	if err := s.repo.Commit(); err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, err)
	}
	return nil
}

func (s *Service) DeleteAccount(ctx context.Context, accountID int64) error {
	const errorMsg = "DeleteAccount for ID %d failed: %w"
	if accountID <= 0 {
		return fmt.Errorf(errorMsg, accountID, "Invalid ID")
	}
	if err := s.repo.BeginTx(ctx); err != nil {
		return fmt.Errorf(errorMsg, accountID, err)
	}
	err := s.repo.AccountQueryBuilder().
		SetCondition(
			db.NewWhere(db.Column_Id, db.Equal, strconv.FormatInt(accountID, 10)),
		).
		Delete(ctx)
	if err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, accountID, err)
	}
	if err := s.repo.Commit(); err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, accountID, err)
	}
	return nil
}
