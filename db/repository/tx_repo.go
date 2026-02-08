package repository

import "context"

type TxRepo interface {
	BeginTx(ctx context.Context) error
	Commit() error
	Rollback() error
}
