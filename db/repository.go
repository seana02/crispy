package db

import "context"

type Repository interface {
	BeginTx(ctx context.Context) error
	Rollback() error
	Commit() error

	SQLBuilder(table Table) *SQLite_Builder
	TransactionQueryBuilder() *SQLite_TransactionBuilder
	AccountQueryBuilder() *SQLite_AccountBuilder
	PostingQueryBuilder() *SQLite_PostingBuilder
}
