package service

import (
	"context"
	"crispy/db"
	"crispy/domain"
	"fmt"
	"strconv"
	"time"
)

func (s *Service) CreateTransaction(ctx context.Context, newTransaction *domain.Transaction) error {
	const errorMsg = "CreateTransaction failed: %w"
	if err := newTransaction.Validate(); err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	if err := s.repo.BeginTx(ctx); err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	// Transaction
	lastId, err := s.repo.TransactionQueryBuilder().
		AddColumns([]db.Column{
			db.Column_Description,
			db.Column_Date,
			db.Column_Status,
			db.Column_DateCreated,
			db.Column_DateUpdated,
		}).
		Insert(ctx, newTransaction.Description(), newTransaction.Date(), newTransaction.Status(), now(), now())
	if err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, err)
	}

	// Tags
	for _, tag := range newTransaction.Tags() {
		if err := s.insertTag(ctx, tag); err != nil {
			s.repo.Rollback()
			return fmt.Errorf(errorMsg, err)
		}
	}

	// Postings
	if err := s.addPostings(ctx, newTransaction.Postings(), lastId); err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, err)
	}

	if err := s.repo.Commit(); err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, err)
	}

	return nil
}

func (s *Service) GetAllTransactions(ctx context.Context, page, perPage int) ([]*domain.Transaction, error) {
	const errorMsg = "GetAllTransactions failed: %w"

	tx_arr, err := s.repo.TransactionQueryBuilder().
		Select(ctx, page, perPage)
	if err != nil {
		return nil, fmt.Errorf(errorMsg, err)
	}

	return tx_arr, nil
}

func (s *Service) GetTransactionByID(ctx context.Context, id int64, includeFullData bool) (*domain.Transaction, error) {
	const errorMsg = "GetTransactionsByID for ID %d failed: %w"

	tx_arr, err := s.repo.TransactionQueryBuilder().
		SetCondition(
			db.NewWhere(db.Column_Id, db.Equal, strconv.FormatInt(id, 10)),
		).
		Select(ctx, 0, 1)
	if err != nil {
		return nil, fmt.Errorf(errorMsg, id, err)
	}

	if !includeFullData {
		return tx_arr[0], nil
	}

	t, err := s.getTags(ctx, id)
	if err != nil {
		return nil, fmt.Errorf(errorMsg, id, err)
	}
	tx_arr[0].SetTags(t)

	p, err := s.getPostings(ctx, id)
	if err != nil {
		return nil, fmt.Errorf(errorMsg, id, err)
	}
	tx_arr[0].SetPostings(p)
	return tx_arr[0], nil
}

func (s *Service) GetTransactionPostings(ctx context.Context, id int64) ([]*domain.Posting, error) {
	return s.getPostings(ctx, id)
}

func (s *Service) UpdateTransaction(ctx context.Context, newTransaction *domain.Transaction) error {
	const errorMsg = "UpdateTranasction failed: %w"
	if err := newTransaction.Validate(); err != nil {
		return fmt.Errorf(errorMsg, err)
	}
	if err := s.repo.BeginTx(ctx); err != nil {
		return fmt.Errorf(errorMsg, err)
	}
	err := s.repo.TransactionQueryBuilder().
		AddColumns([]db.Column{
			db.Column_Description,
			db.Column_Date,
			db.Column_Status,
			db.Column_DateUpdated,
		}).
		SetCondition(
			db.NewWhere(db.Column_Id, db.Equal, strconv.FormatInt(newTransaction.ID(), 10)),
		).
		Update(ctx, newTransaction.Description(), newTransaction.Date(), newTransaction.Status(), now())
	if err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, err)
	}

	if err := s.updateTransactionTags(ctx, newTransaction.ID(), newTransaction.Tags()); err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, err)
	}

	err = s.repo.PostingQueryBuilder().
		SetCondition(db.NewWhere(db.Column_TransactionId, db.Equal, strconv.FormatInt(newTransaction.ID(), 10))).
		Delete(ctx)
	if err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, err)
	}

	if err := s.addPostings(ctx, newTransaction.Postings(), newTransaction.ID()); err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, err)
	}

	if err := s.repo.Commit(); err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, err)
	}
	return nil
}

func (s *Service) DeleteTransaction(ctx context.Context, transactionId int64) error {
	const errorMsg = "DeleteTransaction on ID %d failed: %w"
	if err := s.repo.BeginTx(ctx); err != nil {
		return fmt.Errorf(errorMsg, transactionId, err)
	}
	err := s.repo.TransactionQueryBuilder().
		SetCondition(
			db.NewWhere(db.Column_Id, db.Equal, strconv.FormatInt(transactionId, 10)),
		).
		Delete(ctx)
	if err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, transactionId, err)
	}
	if err := s.repo.Commit(); err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, transactionId, err)
	}
	return nil
}

func (s *Service) getPostings(ctx context.Context, transactionID int64) ([]*domain.Posting, error) {
	p, err := s.repo.PostingQueryBuilder().
		SetCondition(
			db.NewWhere(db.Column_TransactionId, db.Equal, strconv.FormatInt(transactionID, 10)),
		).
		Select(ctx, -1, -1)
	if err != nil {
		return nil, fmt.Errorf("Failed to get posting %d: %w", transactionID, err)
	}
	return p, nil
}

func (s *Service) addPostings(ctx context.Context, postingList []*domain.Posting, transactionId int64) error {
	for _, p := range postingList {
		_, err := s.repo.PostingQueryBuilder().
			AddColumns([]db.Column{
				db.Column_TransactionId,
				db.Column_AccountId,
				db.Column_Amount,
				db.Column_Currency,
				db.Column_DateCreated,
				db.Column_DateUpdated,
			}).
			Insert(ctx, transactionId, p.AccountID(), p.Amount(), p.Currency(), time.Now(), time.Now())
		if err != nil {
			return err
		}
	}
	return nil
}