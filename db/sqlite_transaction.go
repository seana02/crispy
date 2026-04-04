package db

import (
	"context"
	"crispy/domain"
	"errors"
	"strings"
	"time"
)

type StatementType int
type TransactionColumn int
type ComparisonOperator int

const (
	StatementType_Insert StatementType = iota
	StatementType_Select
	StatementType_Update
	StatementType_Delete
)

func (st StatementType) String() string {
	switch st {
	case StatementType_Insert:
		return "INSERT INTO"
	case StatementType_Select:
		return "SELECT"
	case StatementType_Update:
		return "UPDATE"
	case StatementType_Delete:
		return "DELETE"
	}
	return ""
}

const (
	TransactionColumn_Id TransactionColumn = iota
	TransactionColumn_Description
	TransactionColumn_Date
	TransactionColumn_Status
	TransactionColumn_ReferenceId
	TransactionColumn_DateCreated
	TransactionColumn_DateUpdated
	TransactionColumn_All
)

func (tc TransactionColumn) String() string {
	switch tc {
	case TransactionColumn_Id:
		return "id"
	case TransactionColumn_Description:
		return "description"
	case TransactionColumn_Date:
		return "date"
	case TransactionColumn_Status:
		return "status"
	case TransactionColumn_ReferenceId:
		return "reference_id"
	case TransactionColumn_DateCreated:
		return "date_created"
	case TransactionColumn_DateUpdated:
		return "date_updated"
	case TransactionColumn_All:
		return "*"
	}
	return ""
}

const (
	Equal ComparisonOperator = iota
	NotEqual
	LessThan
	GreaterThan
	LessThanOrEqual
	GreaterThanOrEqual
	Like
	Is
	IsNot
)

func (co ComparisonOperator) String() string {
	switch co {
	case Equal:
		return "="
	case NotEqual:
		return "<>"
	case LessThan:
		return "<"
	case GreaterThan:
		return ">"
	case LessThanOrEqual:
		return "<="
	case GreaterThanOrEqual:
		return ">="
	case Like:
		return "LIKE"
	case Is:
		return "IS"
	case IsNot:
		return "IS NOT"
	}
	return ""
}

type Where struct {
	Column TransactionColumn
	Op     ComparisonOperator
	Value  string
}

type Join struct {
	TableName string
	Condition string
}

type SQLite_TransactionBuilder struct {
	handle        *SQLiteDB
	statementType StatementType
	columns       map[TransactionColumn]struct{}
	join          *Join
	where         []Where
	inputTx       *domain.Transaction
}

// Build a query for the Transaction table.
//
// Insert creates a new row using the columns specified. It may return an error.
// Select returns a list of Transactions with all columns, regardless of columns specified.
// Update sets the specified columns to all rows found. It requires at least one WHERE clause
// Delete deletes all rows satisfying the given WHERE clause(s).
func (s *SQLiteDB) TransactionBuilder(op StatementType) *SQLite_TransactionBuilder {
	return &SQLite_TransactionBuilder{
		handle:        s,
		statementType: op,
		columns:       make(map[TransactionColumn]struct{}),
	}
}

func (s *SQLite_TransactionBuilder) AddColumn(c TransactionColumn) *SQLite_TransactionBuilder {
	s.columns[c] = struct{}{}
	return s
}

func (s *SQLite_TransactionBuilder) AddCondition(c TransactionColumn, op ComparisonOperator, val string) *SQLite_TransactionBuilder {
	s.where = append(s.where, Where{c, op, val})
	return s
}

func (s *SQLite_TransactionBuilder) AddJoin(tableName, condition string) *SQLite_TransactionBuilder {
	s.join = &Join{
		TableName: tableName,
		Condition: condition,
	}
	return s
}

func (s *SQLite_TransactionBuilder) Execute(ctx context.Context, args ...any) ([]*domain.Transaction, error) {
	switch s.statementType {
	case StatementType_Insert:
		return nil, insertTransaction(s, ctx, args...)
	case StatementType_Select:
		return selectTransaction(s, ctx, args...)
	case StatementType_Update:
		return nil, updateTransaction(s, ctx, args...)
	case StatementType_Delete:
		return nil, deleteTransaction(s, ctx, args...)
	}
	return nil, nil
}

func insertTransaction(s *SQLite_TransactionBuilder, ctx context.Context, args ...any) error {
	var columnList strings.Builder
	var valueList strings.Builder
	skip := true
	for k, _ := range s.columns {
		if k == TransactionColumn_All {
			return errors.New("Cannot insert using All")
		}
		if !skip {
			columnList.Write([]byte(", "))
			valueList.Write([]byte(", "))
			skip = false
		}
		columnList.Write([]byte(k.String()))
		valueList.Write([]byte("?"))
	}
	var stmt strings.Builder
	stmt.Write([]byte("INSERT INTO \"Transaction\" ("))
	stmt.Write([]byte(columnList.String()))
	stmt.Write([]byte(") VALUES ("))
	stmt.Write([]byte(valueList.String()))
	stmt.Write([]byte(")"))
	return dbExec(s, ctx, stmt.String(), args...)
}

func selectTransaction(s *SQLite_TransactionBuilder, ctx context.Context, args ...any) ([]*domain.Transaction, error) {
	var stmt strings.Builder
	stmt.Write([]byte("SELECT \"Transaction\".* FROM \"Transaction\""))
	if s.join != nil {
		stmt.Write([]byte(" JOIN "))
		stmt.Write([]byte(s.join.TableName))
		stmt.Write([]byte(" on "))
		stmt.Write([]byte(s.join.Condition))
	}
	if len(s.where) > 0 {
		stmt.Write([]byte(" WHERE "))
		for i, w := range s.where {
			if i > 0 {
				stmt.Write([]byte(" AND "))
			}
			stmt.Write([]byte(w.Column.String()))
			stmt.Write([]byte(" "))
			stmt.Write([]byte(w.Op.String()))
			stmt.Write([]byte(" "))
			stmt.Write([]byte(w.Value))
		}
	}
	rows, err := s.handle.db.QueryContext(ctx, stmt.String(), args...)
	if err != nil {
		return nil, err
	}
	var id int64
	var description string
	var date time.Time
	var status domain.Status
	var referenceID *int64
	var dateCreated, dateUpdated time.Time

	var ts []*domain.Transaction
	for rows.Next() {
		err := rows.Scan(
			&id, &description, &date, &status, &referenceID, &dateCreated, &dateUpdated,
		)
		if err != nil {
			return nil, err
		}
		newTx := domain.NewTransaction(
			id, description, date, status, referenceID, dateCreated, dateUpdated, []*domain.Posting{}, []string{},
		)
		ts = append(ts, newTx)
	}
	return ts, nil
}

func updateTransaction(s *SQLite_TransactionBuilder, ctx context.Context, args ...any) error {
	if len(s.where) <= 0 {
		return errors.New("Update requires a WHERE clause")
	}
	var columnList strings.Builder
	skip := true
	for k, _ := range s.columns {
		if !skip {
			columnList.Write([]byte(", "))
			skip = false
		}
		columnList.Write([]byte(k.String()))
		columnList.Write([]byte(" = ?"))
	}
	var stmt strings.Builder
	stmt.Write([]byte("UPDATE \"Transaction\" SET "))
	stmt.Write([]byte(columnList.String()))
	stmt.Write([]byte(" WHERE "))
	for i, w := range s.where {
		if i > 0 {
			stmt.Write([]byte(" AND "))
		}
		stmt.Write([]byte(w.Column.String()))
		stmt.Write([]byte(" "))
		stmt.Write([]byte(w.Op.String()))
		stmt.Write([]byte(" "))
		stmt.Write([]byte(w.Value))
	}
	return dbExec(s, ctx, stmt.String(), args...)
}

func deleteTransaction(s *SQLite_TransactionBuilder, ctx context.Context, args ...any) error {
	if len(s.where) <= 0 {
		return errors.New("Cannot unconditionally delete")
	}
	var stmt strings.Builder
	stmt.Write([]byte("DELETE FROM \"Transaction\" WHERE "))
	for i, w := range s.where {
		if i > 0 {
			stmt.Write([]byte(" AND "))
		}
		stmt.Write([]byte(w.Column.String()))
		stmt.Write([]byte(" "))
		stmt.Write([]byte(w.Op.String()))
		stmt.Write([]byte(" "))
		stmt.Write([]byte(w.Value))
	}
	return dbExec(s, ctx, stmt.String(), args...)
}

func dbExec(s *SQLite_TransactionBuilder, ctx context.Context, stmt string, args ...any) error {
	tx, err := s.handle.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
