package repository

import (
	"context"

	"crispy/domain"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, a *domain.Account) (int64, error)
	GetAccountById(ctx context.Context, id int64) (*domain.Account, error)
	GetAccountByFullName(ctx context.Context, fullName string) (*domain.Account, error)
	UpdateAccount(ctx context.Context, a *domain.Account) error
	DeleteAccount(ctx context.Context, id int64) error
}
