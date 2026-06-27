package db

import (
	"context"
	"crispy/domain"
	"database/sql"
	"fmt"
	"time"
)

type SQLite_TransactionBuilder struct {
	SQLite_Builder
}

func (s *SQLiteDB) TransactionQueryBuilder() *SQLite_TransactionBuilder {
	return &SQLite_TransactionBuilder{
		SQLite_Builder: *s.SQLBuilder(Table_Transaction),
	}
}

func (s *SQLite_TransactionBuilder) AddColumn(c Column) *SQLite_TransactionBuilder {
	s.SQLite_Builder.AddColumn(c)
	return s
}

func (s *SQLite_TransactionBuilder) AddColumns(columns []Column) *SQLite_TransactionBuilder {
	s.SQLite_Builder.AddColumns(columns)
	return s
}

func (s *SQLite_TransactionBuilder) SetCondition(cond Condition) *SQLite_TransactionBuilder {
	s.SQLite_Builder.SetCondition(cond)
	return s
}

func (s *SQLite_TransactionBuilder) AddJoin(tableName Table, condition string) *SQLite_TransactionBuilder {
	s.SQLite_Builder.AddJoin(tableName, condition)
	return s
}

func (s *SQLite_TransactionBuilder) AddOrderBy(c []Column, r bool) *SQLite_TransactionBuilder {
	s.SQLite_Builder.AddOrderBy(c, r)
	return s
}

func (s *SQLite_TransactionBuilder) Select(ctx context.Context, page, perPage int) ([]*domain.Transaction, error) {
	var id int64
	var description string
	var date time.Time
	var status domain.Status
	var referenceID *int64
	var dateCreated, dateUpdated time.Time

	var ts []*domain.Transaction
	callback := func(nextRow *sql.Rows) error {
		err := nextRow.Scan(
			&id, &description, &date, &status, &referenceID, &dateCreated, &dateUpdated,
		)
		if err != nil {
			return fmt.Errorf("Scan failed: %w", err)
		}
		newTx := domain.NewTransaction(
			id, description, date, status, referenceID, dateCreated, dateUpdated, []*domain.Posting{}, []string{},
		)
		ts = append(ts, newTx)
		return nil
	}
	if err := s.SQLite_Builder.Select(ctx, page, perPage, callback); err != nil {
		return nil, fmt.Errorf("SELECT in \"Transaction\" failed: %w", err)
	}

	return ts, nil
}
