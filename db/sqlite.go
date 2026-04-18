package db

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

type ComparisonOperator int
type Table int

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

const (
	Table_Transaction Table = iota
	Table_Account
	Table_Posting
)

func (t Table) String() string {
	switch t {
	case Table_Transaction:
		return "\"Transaction\""
	case Table_Account:
		return "\"Account\""
	case Table_Posting:
		return "\"Posting\""
	}
	return ""
}

type Where struct {
	Column Column
	Op     ComparisonOperator
	Value  string
}

type Join struct {
	TableName string
	Condition string
}

type SQLite_Builder struct {
	Handle    *SQLiteDB
	table     Table
	ColumnSet map[Column]bool
	Columns   []Column
	Join      *Join
	Where     []Where
	OrderBy   []Column
	Reverse   bool
}

type SQLiteDB struct {
	db     *sql.DB
	tx     *sql.Tx
	Logger *slog.Logger
}

func InitSQLite(name string, logger *slog.Logger) (*SQLiteDB, error) {
	handle, err := getHandle(name)
	if err != nil {
		return nil, err
	}

	logger.Debug("SQLite initialized successfully")
	return &SQLiteDB{db: handle, Logger: logger}, nil
}

func getHandle(name string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?_journal=WAL&_synchronous=NORMAL&_loc=auto&parseTime=true", name))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		db.Close()
		return nil, err
	}

	sourceDriver, err := iofs.New(migrationFiles, "migrations")

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "sqlite3", driver)
	if err != nil {
		db.Close()
		return nil, err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		db.Close()
		return nil, err
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func (s *SQLiteDB) BeginTx(ctx context.Context) error {
	if s.tx != nil {
		return fmt.Errorf("Cannot begin- Tx already started")
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	s.tx = tx
	return nil
}

func (s *SQLiteDB) Rollback() error {
	if s.tx == nil {
		return fmt.Errorf("Cannot rollback- no Tx started")
	}
	err := s.tx.Rollback()
	if err != nil {
		return fmt.Errorf("Rollback failed: %s", err)
	}
	s.tx = nil
	return nil
}

func (s *SQLiteDB) Commit() error {
	if s.tx == nil {
		return fmt.Errorf("Cannot commit- no Tx started")
	}
	err := s.tx.Commit()
	if err != nil {
		s.tx.Rollback()
		return fmt.Errorf("Commit failed- Tx rolled back: %s", err)
	}
	s.tx = nil
	return nil
}

// Build a SQL query.
//
// Insert creates a new row using the columns specified. It may return an error.
// Select returns a list of Transactions with all columns, regardless of columns specified.
// Update sets the specified columns to all rows found. It requires at least one WHERE clause
// Delete deletes all rows satisfying the given WHERE clause(s).
func (s *SQLiteDB) SQLBuilder(table Table) *SQLite_Builder {
	return &SQLite_Builder{
		Handle:    s,
		table:     table,
		ColumnSet: make(map[Column]bool),
	}
}

func (s *SQLite_Builder) AddColumn(c Column) {
	if !s.ColumnSet[c] {
		s.ColumnSet[c] = true
		s.Columns = append(s.Columns, c)
	}
}

func (s *SQLite_Builder) AddColumns(columns []Column) {
	for _, c := range columns {
		s.AddColumn(c)
	}
}

func (s *SQLite_Builder) AddCondition(c Column, op ComparisonOperator, val string) {
	s.Where = append(s.Where, Where{c, op, val})
}

func (s *SQLite_Builder) AddJoin(tableName, condition string) {
	s.Join = &Join{
		TableName: tableName,
		Condition: condition,
	}
}

func (s *SQLite_Builder) AddOrderBy(c []Column, r bool) {
	s.OrderBy = c
	s.Reverse = r
}

func (s *SQLite_Builder) Insert(ctx context.Context, args ...any) (int64, error) {
	if s.Handle.tx == nil {
		return -1, errors.New("Error inserting transaction: no Tx started. Call SQLiteDB.BeginTx")
	}
	var columnList strings.Builder
	var valueList strings.Builder
	skip := true
	for _, k := range s.Columns {
		if k == Column_All {
			return -1, errors.New("Cannot insert using All")
		}
		if !skip {
			columnList.Write([]byte(", "))
			valueList.Write([]byte(", "))
		}
		columnList.Write([]byte(k.String()))
		valueList.Write([]byte("?"))
		skip = false
	}
	var stmt strings.Builder
	stmt.Write([]byte("INSERT INTO "))
	stmt.Write([]byte(s.table.String()))
	stmt.Write([]byte(" ("))
	stmt.Write([]byte(columnList.String()))
	stmt.Write([]byte(") VALUES ("))
	stmt.Write([]byte(valueList.String()))
	stmt.Write([]byte(")"))

	s.Handle.Logger.Debug("Inserting", "Query", stmt.String())
	result, err := s.Handle.tx.ExecContext(ctx, stmt.String(), args...)
	if err != nil {
		return -1, fmt.Errorf("Insert failed: %w", err)
	}
	return result.LastInsertId()
}

func (s *SQLite_Builder) Select(ctx context.Context, page, perPage int, callback func(*sql.Rows) error) error {
	var stmt strings.Builder
	stmt.Reset()
	stmt.Write([]byte("SELECT "))
	stmt.Write([]byte(s.table.String()))
	stmt.Write([]byte(".* FROM "))
	stmt.Write([]byte(s.table.String()))

	writeJoin(&stmt, s.Join)
	writeWhere(&stmt, s.Where)
	writeOrderBy(&stmt, s.OrderBy, s.Reverse)
	writePagination(&stmt, page, perPage)

	s.Handle.Logger.Debug("Selecting", "Query", stmt.String())
	rows, err := s.Handle.db.QueryContext(ctx, stmt.String())
	if err != nil {
		return fmt.Errorf("Select failed: %w", err)
	}
	for rows.Next() {
		if err := callback(rows); err != nil {
			return fmt.Errorf("Callback failed: %w", err)
		}
	}
	return nil
}

func (s *SQLite_Builder) Update(ctx context.Context, args ...any) error {
	if s.Handle.tx == nil {
		return errors.New("Error updating transaction: no Tx started. Call SQLiteDB.BeginTx")
	}
	if len(s.Where) <= 0 {
		return errors.New("Update requires a WHERE clause")
	}
	var columnList strings.Builder
	skip := true
	for _, k := range s.Columns {
		if !skip {
			columnList.Write([]byte(", "))
		}
		columnList.Write([]byte(k.String()))
		columnList.Write([]byte(" = ?"))
		skip = false
	}
	var stmt strings.Builder
	stmt.Write([]byte("UPDATE "))
	stmt.Write([]byte(s.table.String()))
	stmt.Write([]byte(" SET "))
	stmt.Write([]byte(columnList.String()))
	writeWhere(&stmt, s.Where)
	fmt.Printf("Query: %+v\nArgs: %v\n", stmt.String(), args)
	s.Handle.Logger.Debug("Updating", "Query", stmt.String(), "Args", args)
	_, err := s.Handle.db.ExecContext(ctx, stmt.String(), args...)
	if err != nil {
		return fmt.Errorf("Update failed: %w", err)
	}
	return nil
}

func (s *SQLite_Builder) Delete(ctx context.Context) error {
	if len(s.Where) <= 0 {
		return errors.New("Cannot unconditionally delete unless explicitly specified")
	}
	var stmt strings.Builder
	stmt.Write([]byte("DELETE FROM "))
	stmt.Write([]byte(s.table.String()))
	writeWhere(&stmt, s.Where)
	s.Handle.Logger.Debug("Deleting", "Query", stmt.String())
	_, err := s.Handle.db.ExecContext(ctx, stmt.String())
	if err != nil {
		return fmt.Errorf("Delete failed: %w", err)
	}
	return nil
}

func writeJoin(stmt *strings.Builder, join *Join) {
	if join != nil {
		stmt.Write([]byte(" JOIN "))
		stmt.Write([]byte(join.TableName))
		stmt.Write([]byte(" on "))
		stmt.Write([]byte(join.Condition))
	}
}

func writeWhere(stmt *strings.Builder, where []Where) {
	if len(where) > 0 {
		stmt.Write([]byte(" WHERE "))
		for i, w := range where {
			if i > 0 {
				stmt.Write([]byte(" AND "))
			}
			stmt.Write([]byte(w.Column.String()))
			stmt.Write([]byte(" "))
			stmt.Write([]byte(w.Op.String()))
			stmt.Write([]byte(" \""))
			stmt.Write([]byte(w.Value))
			stmt.Write([]byte("\""))
		}
	}
}

func writeOrderBy(stmt *strings.Builder, orderby []Column, reverse bool) {
	if orderby == nil {
		return
	}
	stmt.Write([]byte(" ORDER BY "))
	for i, c := range orderby {
		if i > 0 {
			stmt.Write([]byte(", "))
		}
		stmt.Write([]byte(c.String()))
	}
	if reverse {
		stmt.Write([]byte("DESC"))
	}
}

func writePagination(stmt *strings.Builder, page, perPage int) {
	if page < 0 {
		return
	}
	if perPage < 0 {
		return
	}
	stmt.Write([]byte(" LIMIT "))
	stmt.Write([]byte(strconv.Itoa(perPage)))
	stmt.Write([]byte(" OFFSET "))
	stmt.Write([]byte(strconv.Itoa(page * perPage)))
}
