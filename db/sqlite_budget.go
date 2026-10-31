package db

import (
	"context"
	"crispy/domain"
	"database/sql"
	"fmt"

	"github.com/shopspring/decimal"
)

type SQLite_BudgetBuilder struct {
	SQLite_Builder
}

func (s *SQLiteDB) BudgetQueryBuilder() *SQLite_BudgetBuilder {
	return &SQLite_BudgetBuilder{
		SQLite_Builder: *s.SQLBuilder(Table_Budget),
	}
}

func (s *SQLite_BudgetBuilder) AddColumn(c Column) *SQLite_BudgetBuilder {
	s.SQLite_Builder.AddColumn(c)
	return s
}

func (s *SQLite_BudgetBuilder) AddColumns(columns []Column) *SQLite_BudgetBuilder {
	s.SQLite_Builder.AddColumns(columns)
	return s
}

func (s *SQLite_BudgetBuilder) SetCondition(cond Condition) *SQLite_BudgetBuilder {
	s.SQLite_Builder.SetCondition(cond)
	return s
}

func (s *SQLite_BudgetBuilder) AddJoin(tableName Table, condition string) *SQLite_BudgetBuilder {
	s.SQLite_Builder.AddJoin(tableName, condition)
	return s
}

func (s *SQLite_BudgetBuilder) AddOrderBy(c []Column, r bool) *SQLite_BudgetBuilder {
	s.SQLite_Builder.AddOrderBy(c, r)
	return s
}

func (s *SQLite_BudgetBuilder) Select(ctx context.Context, page, perPage int) ([]*domain.Budget, error) {
	var id int64
	var name string
	var description string
	var target decimal.Decimal
	var tagID int64
	var recurrenceID int64

	var ts []*domain.Budget
	callback := func(nextRow *sql.Rows) error {
		err := nextRow.Scan(
			&id, &name, &description, &target, &tagID, &recurrenceID,
		)
		if err != nil {
			return fmt.Errorf("Scan failed: %w", err)
		}
		newRecurrence := domain.Recurrence{}
		newRecurrence.SetID(recurrenceID)
		newBudget := domain.NewBudget(
			id, name, description, target, tagID, &newRecurrence,
		)
		ts = append(ts, newBudget)
		return nil
	}
	if err := s.SQLite_Builder.Select(ctx, page, perPage, callback); err != nil {
		return nil, fmt.Errorf("SELECT in \"Budget\" failed: %w", err)
	}

	return ts, nil
}
