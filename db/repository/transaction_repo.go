package repository

import (
	"context"
	"crispy/domain"
	"time"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, t *domain.Transaction) (int64, error)
	GetTransactionById(ctx context.Context, id int64) (*domain.Transaction, error)
	GetTransactionByTag(ctx context.Context, tagID int64) ([]*domain.Transaction, error)
	GetTransactionByDate(ctx context.Context, from time.Time, to time.Time) ([]*domain.Transaction, error)
	UpdateTransaction(ctx context.Context, a *domain.Transaction) (int64, error)
	DeleteTransaction(ctx context.Context, id int64) error
}
