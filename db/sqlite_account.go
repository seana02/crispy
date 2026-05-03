package db

import (
	"context"
	"crispy/domain"
	"database/sql"
	"fmt"
	"time"
)

type SQLite_AccountBuilder struct {
	SQLite_Builder
}

func (s *SQLiteDB) AccountQueryBuilder() *SQLite_AccountBuilder {
	return &SQLite_AccountBuilder{
		SQLite_Builder: *s.SQLBuilder(Table_Account),
	}
}

func (s *SQLite_AccountBuilder) AddColumn(c Column) *SQLite_AccountBuilder {
	s.SQLite_Builder.AddColumn(c)
	return s
}

func (s *SQLite_AccountBuilder) AddColumns(columns []Column) *SQLite_AccountBuilder {
	s.SQLite_Builder.AddColumns(columns)
	return s
}

func (s *SQLite_AccountBuilder) SetCondition(cond Condition) *SQLite_AccountBuilder {
	s.SQLite_Builder.SetCondition(cond)
	return s
}

func (s *SQLite_AccountBuilder) AddJoin(tableName, condition string) *SQLite_AccountBuilder {
	s.SQLite_Builder.AddJoin(tableName, condition)
	return s
}

func (s *SQLite_AccountBuilder) AddOrderBy(c []Column, r bool) *SQLite_AccountBuilder {
	s.SQLite_Builder.AddOrderBy(c, r)
	return s
}

func (s *SQLite_AccountBuilder) Select(ctx context.Context, page, perPage int) ([]*domain.Account, error) {
	var id int64
	var parentID int64
	var name string
	var type_ domain.Type
	var currency string
	var description string
	var active bool
	var dateCreated, dateUpdated time.Time

	var accts []*domain.Account
	callback := func(rows *sql.Rows) error {
		err := rows.Scan(
			&id, &parentID, &name, &type_, &currency, &description, &active, &dateCreated, &dateUpdated,
		)
		if err != nil {
			return fmt.Errorf("Scan failed: %w", err)
		}
		newAcct := domain.NewAccount(
			id, parentID, name, type_, currency, description, active, dateCreated, dateUpdated,
		)
		accts = append(accts, newAcct)
		return nil
	}
	if err := s.SQLite_Builder.Select(ctx, page, perPage, callback); err != nil {
		return nil, fmt.Errorf("SELECT in \"Account\" failed: %w", err)
	}
	return accts, nil
}
