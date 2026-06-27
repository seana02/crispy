package db

import (
	"context"
	"crispy/domain"
	"database/sql"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type SQLite_PostingBuilder struct {
	SQLite_Builder
}

func (s *SQLiteDB) PostingQueryBuilder() *SQLite_PostingBuilder {
	return &SQLite_PostingBuilder{
		SQLite_Builder: *s.SQLBuilder(Table_Posting),
	}
}

func (s *SQLite_PostingBuilder) AddColumn(c Column) *SQLite_PostingBuilder {
	s.SQLite_Builder.AddColumn(c)
	return s
}

func (s *SQLite_PostingBuilder) AddColumns(columns []Column) *SQLite_PostingBuilder {
	s.SQLite_Builder.AddColumns(columns)
	return s
}

func (s *SQLite_PostingBuilder) SetCondition(cond Condition) *SQLite_PostingBuilder {
	s.SQLite_Builder.SetCondition(cond)
	return s
}

func (s *SQLite_PostingBuilder) AddJoin(tableName Table, condition string) *SQLite_PostingBuilder {
	s.SQLite_Builder.AddJoin(tableName, condition)
	return s
}

func (s *SQLite_PostingBuilder) AddOrderBy(c []Column, r bool) *SQLite_PostingBuilder {
	s.SQLite_Builder.AddOrderBy(c, r)
	return s
}

func (s *SQLite_PostingBuilder) Select(ctx context.Context, page, perPage int) ([]*domain.Posting, error) {
	var id int64
	var transactionID int64
	var accountID int64
	var amount decimal.Decimal
	var currency string
	var dateCreated, dateUpdated time.Time

	var ps []*domain.Posting
	callback := func(nextRow *sql.Rows) error {
		err := nextRow.Scan(
			&id, &transactionID, &accountID, &amount, &currency, &dateCreated, &dateUpdated,
		)
		if err != nil {
			return fmt.Errorf("Scan failed: %w", err)
		}
		newPosting := domain.NewPosting(
			id, transactionID, accountID, amount, currency, dateCreated, dateUpdated,
		)
		ps = append(ps, newPosting)
		return nil
	}
	if err := s.SQLite_Builder.Select(ctx, page, perPage, callback); err != nil {
		return nil, fmt.Errorf("SELECT in \"Posting\" failed: %w", err)
	}
	return ps, nil
}
