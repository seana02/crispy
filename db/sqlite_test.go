package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"crispy/domain"
)

// newTestDB creates an in-memory SQLite database for testing.
func newTestDB(t *testing.T) *SQLiteDB {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	db, err := InitSQLite(":memory:", logger)
	if err != nil {
		t.Fatalf("InitSQLite failed: %v", err)
	}
	// SQLite in-memory databases are per-connection. Without this, the pool
	// may open a second connection that has no schema, causing "no such table"
	// errors on any operation that doesn't happen to reuse the first connection.
	db.db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.db.Close() })
	return db
}

// seedAccount inserts a minimal Account row inside the current tx and returns its id.
func seedAccount(t *testing.T, ctx context.Context, db *SQLiteDB, name string, acctType domain.Type) int64 {
	t.Helper()
	b := db.SQLBuilder(Table_Account)
	b.AddColumns([]Column{
		Column_ParentId, Column_Name, Column_Type, Column_Currency,
		Column_Description, Column_Active, Column_DateCreated, Column_DateUpdated,
	})
	now := time.Now().Truncate(time.Second)
	id, err := b.Insert(ctx, int64(0), name, string(acctType), "USD", "test account", true, now, now)
	if err != nil {
		t.Fatalf("seedAccount(%q) insert failed: %v", name, err)
	}
	return id
}

// seedTransaction inserts a Transaction row inside the current tx and returns its id.
func seedTransaction(t *testing.T, ctx context.Context, db *SQLiteDB, desc string) int64 {
	t.Helper()
	b := db.SQLBuilder(Table_Transaction)
	b.AddColumns([]Column{
		Column_Description, Column_Date, Column_Status,
		Column_DateCreated, Column_DateUpdated,
	})
	now := time.Now().Truncate(time.Second)
	id, err := b.Insert(ctx, desc, now, string(domain.Cleared), now, now)
	if err != nil {
		t.Fatalf("seedTransaction(%q) insert failed: %v", desc, err)
	}
	return id
}

// seedPosting inserts a Posting row inside the current tx and returns its id.
func seedPosting(t *testing.T, ctx context.Context, db *SQLiteDB, txID, acctID int64, amount string) int64 {
	t.Helper()
	b := db.SQLBuilder(Table_Posting)
	b.AddColumns([]Column{
		Column_TransactionId, Column_AccountId, Column_Amount,
		Column_Currency, Column_DateCreated, Column_DateUpdated,
	})
	now := time.Now().Truncate(time.Second)
	id, err := b.Insert(ctx, txID, acctID, amount, "USD", now, now)
	if err != nil {
		t.Fatalf("seedPosting insert failed: %v", err)
	}
	return id
}

// idStr converts an int64 to a decimal string suitable for use as a Where value.
func idStr(id int64) string { return fmt.Sprintf("%d", id) }

func TestInitSQLite(t *testing.T) {
	db := newTestDB(t)
	if db == nil {
		t.Fatal("expected non-nil SQLiteDB")
	}
	if db.db == nil {
		t.Fatal("expected non-nil internal sql.DB handle")
	}

	rows, err := db.db.Query(`SELECT id FROM Account LIMIT 1`)
	if err != nil {
		t.Fatalf("Account table not present after migration: %v", err)
	}
	rows.Close()

	var fkEnabled int
	err = db.db.QueryRow(`PRAGMA foreign_keys`).Scan(&fkEnabled)
	if err != nil {
		t.Fatalf("PRAGMA foreign_keys query failed: %v", err)
	}
	if fkEnabled != 1 {
		t.Errorf("expected foreign_keys = 1, got %d", fkEnabled)
	}
}

func TestTransaction(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	if err := db.Commit(); err == nil {
		t.Fatal("expected error when committing without active tx")
	}
	if err := db.Rollback(); err == nil {
		t.Fatal("expected error when rolling back without active tx")
	}

	if err := db.BeginTx(ctx); err != nil {
		t.Fatalf("BeginTx failed: %v", err)
	}

	if db.tx == nil {
		t.Fatal("expected tx to be set after BeginTx")
	}

	if err := db.BeginTx(ctx); err == nil {
		t.Fatal("expected error when beginning tx inside active tx")
	}

	seedAccount(t, ctx, db, "ShouldPersist", domain.Asset)
	if err := db.Commit(); err != nil {
		t.Fatalf("Commit failed: %v", err)
	}
	if db.tx != nil {
		t.Fatal("expected tx to be nil after Commit")
	}

	var count int
	db.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM Account WHERE name = ?`, "ShouldPersist").Scan(&count)
	if count != 1 {
		t.Errorf("expected 1 row after commit, got %d", count)
	}

	db.BeginTx(ctx)
	seedAccount(t, ctx, db, "ShouldDisappear", domain.Asset)
	if err := db.Rollback(); err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}
	if db.tx != nil {
		t.Fatal("expected tx to be nil after Rollback")
	}

	db.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM Account WHERE name = ?`, "ShouldDisappear").Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 rows after rollback, got %d", count)
	}
}

func TestComparisonOperator_String(t *testing.T) {
	cases := []struct {
		op   ComparisonOperator
		want string
	}{
		{Equal, "="},
		{NotEqual, "<>"},
		{LessThan, "<"},
		{GreaterThan, ">"},
		{LessThanOrEqual, "<="},
		{GreaterThanOrEqual, ">="},
		{Like, "LIKE"},
		{Is, "IS"},
		{IsNot, "IS NOT"},
	}
	for _, tc := range cases {
		if got := tc.op.String(); got != tc.want {
			t.Errorf("ComparisonOperator(%d).String() = %q, want %q", tc.op, got, tc.want)
		}
	}
}

func TestTable_String(t *testing.T) {
	cases := []struct {
		tbl  Table
		want string
	}{
		{Table_Transaction, `"Transaction"`},
		{Table_Account, `"Account"`},
		{Table_Posting, `"Posting"`},
	}
	for _, tc := range cases {
		if got := tc.tbl.String(); got != tc.want {
			t.Errorf("Table(%d).String() = %q, want %q", tc.tbl, got, tc.want)
		}
	}
}

func TestColumn_String(t *testing.T) {
	cases := []struct {
		col  Column
		want string
	}{
		{Column_All, "*"},
		{Column_Id, "id"},
		{Column_Description, "description"},
		{Column_Date, "date"},
		{Column_Status, "status"},
		{Column_ReferenceId, "reference_id"},
		{Column_DateCreated, "date_created"},
		{Column_DateUpdated, "date_updated"},
		{Column_ParentId, "parent_id"},
		{Column_Name, "name"},
		{Column_Type, "type"},
		{Column_Currency, "currency"},
		{Column_Active, "active"},
		{Column_TransactionId, "transaction_id"},
		{Column_AccountId, "account_id"},
		{Column_Amount, "amount"},
	}
	for _, tc := range cases {
		if got := tc.col.String(); got != tc.want {
			t.Errorf("Column(%d).String() = %q, want %q", tc.col, got, tc.want)
		}
	}
}

// ── Where / Condition ToSQL ───────────────────────────────────────────────────

func TestWhere_ToSQL(t *testing.T) {
	cases := []struct {
		w         Condition
		want      string
		argLength int
	}{
		{NewWhere(Column_Id, Equal, "42"), "(id = ?)", 1},
		{NewWhere(Column_Description, Like, "%food%"), "(description LIKE ?)", 1},
		{
			AndCondition{Conditions: []Condition{
				NewWhere(Column_Id, Equal, "1"),
				NewWhere(Column_Status, Equal, "Cleared"),
			}},
			"((id = ?) AND (status = ?))",
			2,
		},
		{
			OrCondition{Conditions: []Condition{
				NewWhere(Column_Status, Equal, "Pending"),
				NewWhere(Column_Status, Equal, "Posted"),
			}},
			"((status = ?) OR (status = ?))",
			2,
		},
		{
			AndCondition{Conditions: []Condition{
				AndCondition{Conditions: []Condition{
					NewWhere(Column_Id, Equal, "1"),
					NewWhere(Column_Description, Equal, "x"),
				}},
				NewWhere(Column_Status, Equal, "Cleared"),
			}},
			"(((id = ?) AND (description = ?)) AND (status = ?))",
			3,
		},
		{
			OrCondition{Conditions: []Condition{
				AndCondition{Conditions: []Condition{
					NewWhere(Column_Id, Equal, "1"),
					OrCondition{Conditions: []Condition{
						NewWhere(Column_Description, Like, "%desc%"),
						NewWhere(Column_Status, Equal, "Cleared"),
					}},
				}},
				AndCondition{Conditions: []Condition{
					NewWhere(Column_Status, Equal, "Cleared"),
					NewWhere(Column_Id, Equal, "14"),
				}},
			}},
			"(((id = ?) AND ((description LIKE ?) OR (status = ?))) OR ((status = ?) AND (id = ?)))",
			5,
		},
	}

	for _, c := range cases {
		got, args := c.w.ToSQL()
		if got != c.want {
			t.Errorf("ToSQL() sql = %q, want %q", got, c.want)
		}
		if len(args) != c.argLength {
			t.Errorf("ToSQL() args = %v, wanted length %d", args, c.argLength)
		}
	}
}

// ── SQLite_Builder column management ──────────────────────────────────────────

func TestSQLBuilder_AddColumn(t *testing.T) {
	db := newTestDB(t)
	b := db.SQLBuilder(Table_Transaction)
	b.AddColumn(Column_Id)
	b.AddColumn(Column_Description)
	b.AddColumn(Column_Date)
	b.AddColumn(Column_Id)

	if len(b.Columns) != 3 {
		t.Errorf("expected 3 unique columns, got %v", b.Columns)
	}
	if b.Columns[0] != Column_Id || b.Columns[1] != Column_Description || b.Columns[2] != Column_Date {
		t.Errorf("column order not preserved: %v", b.Columns)
	}
}

func TestSQLBuilder_AddColumns(t *testing.T) {
	db := newTestDB(t)
	b := db.SQLBuilder(Table_Account)
	b.AddColumns([]Column{Column_Name, Column_Type, Column_Name}) // Name duplicated

	if len(b.Columns) != 2 {
		t.Errorf("expected 2 unique columns, got %d", len(b.Columns))
	}
}

func TestSQLBuilder_SetCondition(t *testing.T) {
	db := newTestDB(t)
	b := db.SQLBuilder(Table_Transaction)
	b.SetCondition(NewWhere(Column_Id, Equal, "1"))
	b.SetCondition(NewWhere(Column_Id, Equal, "2"))

	if b.Where == nil {
		t.Fatal("expected Where to be set after SetCondition")
	}
	_, args := b.Where.ToSQL()
	if len(args) != 1 || args[0] != "2" {
		t.Errorf("expected second condition to overwrite first, args = %v", args)
	}
}

func TestSQLBuilder_AddJoin(t *testing.T) {
	db := newTestDB(t)
	b := db.SQLBuilder(Table_Transaction)
	b.AddJoin("Posting", `"Transaction".id = Posting.transaction_id`)

	if b.Join == nil {
		t.Fatal("expected Join to be set")
	}
	if b.Join.TableName != "Posting" {
		t.Errorf("Join.TableName = %q, want %q", b.Join.TableName, "Posting")
	}
	if b.Join.Condition != `"Transaction".id = Posting.transaction_id` {
		t.Errorf("Join.Condition = %q", b.Join.Condition)
	}
}

func TestSQLBuilder_AddOrderBy(t *testing.T) {
	db := newTestDB(t)
	b := db.SQLBuilder(Table_Transaction)
	b.AddOrderBy([]Column{Column_Date, Column_Id}, true)

	if len(b.OrderBy) != 2 {
		t.Errorf("expected 2 order-by columns, got %d", len(b.OrderBy))
	}
	if !b.Reverse {
		t.Error("expected Reverse = true")
	}
}

func TestSQLBuilder_AddOrderBy_NotReversed(t *testing.T) {
	db := newTestDB(t)
	b := db.SQLBuilder(Table_Transaction)
	b.AddOrderBy([]Column{Column_Date}, false)

	if b.Reverse {
		t.Error("expected Reverse = false")
	}
}

// ── Insert ────────────────────────────────────────────────────────────────────

func TestInsert_ErrorWithoutTx(t *testing.T) {
	db := newTestDB(t)
	b := db.SQLBuilder(Table_Account)
	b.AddColumn(Column_All)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	_, err := b.Insert(ctx, "test")
	if err == nil {
		t.Fatal("expected error when inserting without active tx")
	}
	_, err = b.Insert(ctx, "anything")
	if err == nil {
		t.Fatal("expected error when inserting with Column_All")
	}
}

func TestInsert_MultipleRows_UniqueIDs(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	db.BeginTx(ctx)    
	defer db.Rollback()

	id1 := seedAccount(t, ctx, db, "Checking", domain.Asset)
	id2 := seedAccount(t, ctx, db, "Savings", domain.Asset)
	if id1 <= 0 || id2 <= 0 {
		t.Errorf("expected positive ID, got %d and %d", id1, id2)
	}
	if id1 == id2 {
		t.Errorf("expected distinct IDs, got id1=%d id2=%d", id1, id2)
	}
}

func TestInsert_ForeignKeyViolation_ReturnsError(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	db.BeginTx(ctx)    
	defer db.Rollback()

	// Insert a Posting referencing a nonexistent transaction.
	b := db.SQLBuilder(Table_Posting)
	b.AddColumns([]Column{
		Column_TransactionId, Column_AccountId, Column_Amount,
		Column_Currency, Column_DateCreated, Column_DateUpdated,
	})
	now := time.Now()
	_, err := b.Insert(ctx, int64(99999), int64(0), "100.00", "USD", now, now)
	if err == nil {
		t.Error("expected foreign key error when referencing nonexistent transaction")
	}
}

// ── Select (raw builder) ──────────────────────────────────────────────────────

func TestSelect_WithWhereClause(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	db.BeginTx(ctx)
	seedAccount(t, ctx, db, "Alpha", domain.Asset)
	seedAccount(t, ctx, db, "Beta", domain.Liability)
	db.Commit()

	b := db.SQLBuilder(Table_Account)
	b.SetCondition(NewWhere(Column_Name, Equal, "Alpha"))

	var count int
	err := b.Select(ctx, 0, 10, func(_ *sql.Rows) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatalf("Select failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 row, got %d", count)
	}
}

func TestSelect_NoWhereClause_ReturnsAllRows(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	db.BeginTx(ctx)
	seedAccount(t, ctx, db, "A", domain.Asset)
	seedAccount(t, ctx, db, "B", domain.Asset)
	db.Commit()

	b := db.SQLBuilder(Table_Account)
	var count int
	err := b.Select(ctx, 0, 100, func(_ *sql.Rows) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatalf("Select failed: %v", err)
	}
	// root account + A + B = 3
	if count < 3 {
		t.Errorf("expected at least 3 rows, got %d", count)
	}
}

func TestSelect_Pagination_LimitsResults(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	db.BeginTx(ctx)
	for _, name := range []string{"P1", "P2", "P3", "P4", "P5"} {
		seedTransaction(t, ctx, db, name)
	}
	db.Commit()

	b := db.SQLBuilder(Table_Transaction)
	var count int
	err := b.Select(ctx, 0, 2, func(_ *sql.Rows) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatalf("Select failed: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 rows with perPage=2, got %d", count)
	}

	count = 0
	err = b.Select(ctx, 1, 2, func(_ *sql.Rows) error { // page 1, 2 per page
		count++
		return nil
	})
	if err != nil {
		t.Fatalf("Select page 2 failed: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 rows on page 2, got %d", count)
	}
}

func TestSelect_NegativePagePerPage_NoLimit(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	db.BeginTx(ctx)
	for _, name := range []string{"X", "Y", "Z"} {
		seedAccount(t, ctx, db, name, domain.Asset)
	}
	db.Commit()

	b := db.SQLBuilder(Table_Account)
	var count int
	err := b.Select(ctx, -1, -1, func(_ *sql.Rows) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatalf("Select with negative pagination failed: %v", err)
	}
	// root account + X + Y + Z = 4
	if count < 4 {
		t.Errorf("expected at least 4 rows with no limit, got %d", count)
	}
}

func TestSelect_WithJoin(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	db.BeginTx(ctx)
	acctID := seedAccount(t, ctx, db, "Checking", domain.Asset)
	txID := seedTransaction(t, ctx, db, "Buy coffee")
	seedPosting(t, ctx, db, txID, acctID, "5.00")
	db.Commit()

	b := db.SQLBuilder(Table_Transaction)
	b.AddJoin(`"Posting"`, `"Transaction".id = "Posting".transaction_id`)

	var count int
	err := b.Select(ctx, 0, 100, func(_ *sql.Rows) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatalf("Select with JOIN failed: %v", err)
	}
	if count < 1 {
		t.Error("expected at least one row from joined query")
	}
}

func TestSelect_WithOrderBy_AscendingIDs(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	db.BeginTx(ctx)
	seedTransaction(t, ctx, db, "First")
	seedTransaction(t, ctx, db, "Second")
	seedTransaction(t, ctx, db, "Third")
	db.Commit()

	b := db.SQLBuilder(Table_Transaction)
	b.AddOrderBy([]Column{Column_Id}, false)

	var ids []int64
	err := b.Select(ctx, 0, 100, func(rows *sql.Rows) error {
		var id int64
		var desc string
		var date time.Time
		var status string
		var refID *int64
		var created, updated time.Time
		if err := rows.Scan(&id, &desc, &date, &status, &refID, &created, &updated); err != nil {
			return err
		}
		ids = append(ids, id)
		return nil
	})
	if err != nil {
		t.Fatalf("Select with ORDER BY failed: %v", err)
	}
	for i := 1; i < len(ids); i++ {
		if ids[i] < ids[i-1] {
			t.Errorf("expected ascending order, got ids[%d]=%d before ids[%d]=%d", i-1, ids[i-1], i, ids[i])
		}
	}
}

func TestSelect_AndCondition(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	db.BeginTx(ctx)
	seedTransaction(t, ctx, db, "Match")
	seedTransaction(t, ctx, db, "NoMatch")
	db.Commit()

	b := db.SQLBuilder(Table_Transaction)
	b.SetCondition(AndCondition{Conditions: []Condition{
		NewWhere(Column_Description, Equal, "Match"),
		NewWhere(Column_Status, Equal, string(domain.Cleared)),
	}})

	var count int
	err := b.Select(ctx, 0, 10, func(_ *sql.Rows) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatalf("Select with AND condition failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 row, got %d", count)
	}
}

func TestSelect_CallbackError_Propagates(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	db.BeginTx(ctx)
	seedAccount(t, ctx, db, "Fail", domain.Asset)
	db.Commit()

	b := db.SQLBuilder(Table_Account)
	err := b.Select(ctx, 0, 10, func(_ *sql.Rows) error {
		return fmt.Errorf("intentional callback error")
	})
	if err == nil {
		t.Fatal("expected error from callback to propagate")
	}
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestUpdate_ErrorWithoutTx(t *testing.T) {
	db := newTestDB(t)
	b := db.SQLBuilder(Table_Account)
	b.AddColumn(Column_Name)
	b.SetCondition(NewWhere(Column_Id, Equal, "1"))
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	if err := b.Update(ctx, "New Name"); err == nil {
		t.Fatal("expected error when updating without active tx")
	}
}

func TestUpdate_ErrorWithoutWhereClause(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	db.BeginTx(ctx)    
	defer db.Rollback()

	b := db.SQLBuilder(Table_Account)
	b.AddColumn(Column_Name)
	// No SetCondition.

	if err := b.Update(ctx, "New Name"); err == nil {
		t.Fatal("expected error when updating without WHERE clause")
	}
}

func TestUpdate_ModifiesRow(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	db.BeginTx(ctx)
	id := seedAccount(t, ctx, db, "OldName", domain.Asset)
	db.Commit()

	db.BeginTx(ctx)

	b := db.SQLBuilder(Table_Account)
	b.AddColumn(Column_Name)
	b.SetCondition(NewWhere(Column_Id, Equal, idStr(id)))

	if err := b.Update(ctx, "NewName"); err != nil {
		db.Rollback()
		t.Fatalf("Update failed: %v", err)
	}
	db.Commit()

	var name string
	db.db.QueryRowContext(ctx, `SELECT name FROM Account WHERE id = ?`, id).Scan(&name)
	if name != "NewName" {
		t.Errorf("expected name %q after committed update, got %q", "NewName", name)
	}

	db.BeginTx(ctx)
	if err := b.Update(ctx, "ShouldNotStick"); err != nil {
		db.Rollback()
		t.Fatalf("Update failed: %v", err)
	}
	db.Rollback()

	db.db.QueryRowContext(ctx, `SELECT name FROM Account WHERE id = ?`, id).Scan(&name)
	if name != "NewName" {
		t.Errorf("expected rollback to restore name %q, got %q", "Original", name)
	}
}

// ── Delete ────────────────────────────────────────────────────────────────────

func TestDelete_RemovesMatchingRow(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	db.BeginTx(ctx)
	id := seedAccount(t, ctx, db, "ToDelete", domain.Asset)
	db.Commit()

	db.BeginTx(ctx)
	b := db.SQLBuilder(Table_Account)
	if err := b.Delete(ctx); err == nil {
		t.Fatal("expected error when deleting without WHERE clause")
	}

	b.SetCondition(NewWhere(Column_Id, Equal, idStr(id)))

	if err := b.Delete(ctx); err != nil {
		db.Rollback()
		t.Fatalf("Delete failed: %v", err)
	}
	db.Commit()

	var count int
	db.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM Account WHERE id = ?`, id).Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 rows after delete, got %d", count)
	}
}

func TestDelete_LeavesNonMatchingRows(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	db.BeginTx(ctx)
	id1 := seedAccount(t, ctx, db, "Keep", domain.Asset)
	id2 := seedAccount(t, ctx, db, "Remove", domain.Asset)
	db.Commit()

	b := db.SQLBuilder(Table_Account)
	b.SetCondition(NewWhere(Column_Id, Equal, idStr(id2)))

	if err := b.Delete(ctx); err == nil {
		t.Fatal("expected error when deleting without active tx")
	}

	db.BeginTx(ctx)

	if err := b.Delete(ctx); err != nil {
		db.Rollback()
		t.Fatalf("Delete failed: %v", err)
	}
	db.Commit()

	var count int
	db.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM Account WHERE id = ?`, id1).Scan(&count)
	if count != 1 {
		t.Errorf("expected row id=%d to survive delete, got count=%d", id1, count)
	}

	db.BeginTx(ctx)
	b.SetCondition(NewWhere(Column_Id, Equal, idStr(id1)))
	if err := b.Delete(ctx); err != nil {
		db.Rollback()
		t.Fatalf("Delete failed: %v", err)
	}
	db.Rollback()

	db.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM Account WHERE id = ?`, id1).Scan(&count)
	if count != 1 {
		t.Errorf("expected rollback to restore the deleted row, got count=%d", count)
	}
}

// ── AccountQueryBuilder ───────────────────────────────────────────────────────

func TestAccountQueryBuilder_SelectAll(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	db.BeginTx(ctx)
	seedAccount(t, ctx, db, "Cash", domain.Asset)
	db.Commit()

	accts, err := db.AccountQueryBuilder().Select(ctx, 0, 100)
	if err != nil {
		t.Fatalf("AccountQueryBuilder.Select failed: %v", err)
	}
	// Migration seeds a root account; at minimum 2 rows.
	if len(accts) < 2 {
		t.Errorf("expected at least 2 accounts, got %d", len(accts))
	}
}

func TestAccountQueryBuilder_FilterByName(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	db.BeginTx(ctx)
	seedAccount(t, ctx, db, "Wallet", domain.Asset)
	seedAccount(t, ctx, db, "Other", domain.Liability)
	db.Commit()

	accts, err := db.AccountQueryBuilder().
		SetCondition(NewWhere(Column_Name, Equal, "Wallet")).
		Select(ctx, 0, 10)
	if err != nil {
		t.Fatalf("AccountQueryBuilder.Select failed: %v", err)
	}
	if len(accts) != 1 {
		t.Fatalf("expected 1 account, got %d", len(accts))
	}
	a := accts[0]
	if a.Name() != "Wallet" {
		t.Errorf("Name() = %q, want Wallet", a.Name())
	}
	if a.Currency() != "USD" {
		t.Errorf("Currency() = %q, want USD", a.Currency())
	}
	if !a.Active() {
		t.Error("expected Active() = true")
	}
	if a.Type() != domain.Asset {
		t.Errorf("Type() = %q, want Asset", a.Type())
	}
}

func TestAccountQueryBuilder_FilterByType(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	db.BeginTx(ctx)
	seedAccount(t, ctx, db, "CreditCard", domain.Liability)
	db.Commit()

	accts, err := db.AccountQueryBuilder().
		SetCondition(NewWhere(Column_Type, Equal, string(domain.Liability))).
		Select(ctx, 0, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, a := range accts {
		if a.Type() != domain.Liability {
			t.Errorf("expected Liability, got %q", a.Type())
		}
	}
}

func TestAccountQueryBuilder_ChainedMethodsReturnCorrectType(t *testing.T) {
	db := newTestDB(t)
	b := db.AccountQueryBuilder().
		AddColumn(Column_Name).
		SetCondition(NewWhere(Column_Active, Equal, "1")).
		AddOrderBy([]Column{Column_Name}, false)

	if b == nil {
		t.Fatal("chained builder must not return nil")
	}
}

// ── TransactionQueryBuilder ───────────────────────────────────────────────────

func TestTransactionQueryBuilder_SelectAll(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	db.BeginTx(ctx)
	seedTransaction(t, ctx, db, "Groceries")
	seedTransaction(t, ctx, db, "Rent")
	db.Commit()

	txs, err := db.TransactionQueryBuilder().Select(ctx, 0, 100)
	if err != nil {
		t.Fatalf("TransactionQueryBuilder.Select failed: %v", err)
	}
	if len(txs) < 2 {
		t.Errorf("expected at least 2 transactions, got %d", len(txs))
	}
}

func TestTransactionQueryBuilder_FilterByDescription(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	db.BeginTx(ctx)
	seedTransaction(t, ctx, db, "Groceries")
	seedTransaction(t, ctx, db, "Rent")
	db.Commit()

	txs, err := db.TransactionQueryBuilder().
		SetCondition(NewWhere(Column_Description, Equal, "Rent")).
		Select(ctx, 0, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(txs) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(txs))
	}
	if txs[0].Description() != "Rent" {
		t.Errorf("Description() = %q, want Rent", txs[0].Description())
	}
}

func TestTransactionQueryBuilder_FilterByStatus(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	db.BeginTx(ctx)
	seedTransaction(t, ctx, db, "A cleared tx")
	db.Commit()

	txs, err := db.TransactionQueryBuilder().
		SetCondition(NewWhere(Column_Status, Equal, string(domain.Cleared))).
		Select(ctx, 0, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, tx := range txs {
		if tx.Status() != domain.Cleared {
			t.Errorf("expected Cleared, got %q", tx.Status())
		}
	}
}

func TestTransactionQueryBuilder_LikeFilter(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	db.BeginTx(ctx)
	seedTransaction(t, ctx, db, "Coffee shop")
	seedTransaction(t, ctx, db, "Coffee beans")
	seedTransaction(t, ctx, db, "Rent payment")
	db.Commit()

	txs, err := db.TransactionQueryBuilder().
		SetCondition(NewWhere(Column_Description, Like, "Coffee%")).
		Select(ctx, 0, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(txs) != 2 {
		t.Errorf("expected 2 coffee transactions, got %d", len(txs))
	}
}

// ── PostingQueryBuilder ───────────────────────────────────────────────────────

func TestPostingQueryBuilder_SelectAll(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	db.BeginTx(ctx)
	acctID := seedAccount(t, ctx, db, "Savings", domain.Asset)
	txID := seedTransaction(t, ctx, db, "Deposit")
	seedPosting(t, ctx, db, txID, acctID, "500.00")
	db.Commit()

	postings, err := db.PostingQueryBuilder().Select(ctx, 0, 100)
	if err != nil {
		t.Fatalf("PostingQueryBuilder.Select failed: %v", err)
	}
	if len(postings) < 1 {
		t.Errorf("expected at least 1 posting, got %d", len(postings))
	}
}

func TestPostingQueryBuilder_FilterByTransactionID(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	db.BeginTx(ctx)
	acctID := seedAccount(t, ctx, db, "Checking", domain.Asset)
	txID1 := seedTransaction(t, ctx, db, "Paycheck")
	txID2 := seedTransaction(t, ctx, db, "Other")
	seedPosting(t, ctx, db, txID1, acctID, "1000.00")
	seedPosting(t, ctx, db, txID1, acctID, "-1000.00")
	seedPosting(t, ctx, db, txID2, acctID, "50.00")
	db.Commit()

	postings, err := db.PostingQueryBuilder().
		SetCondition(NewWhere(Column_TransactionId, Equal, idStr(txID1))).
		Select(ctx, 0, 10)
	if err != nil {
		t.Fatalf("PostingQueryBuilder.Select failed: %v", err)
	}
	if len(postings) != 2 {
		t.Errorf("expected 2 postings for txID1, got %d", len(postings))
	}
	for _, p := range postings {
		if p.TransactionID() != txID1 {
			t.Errorf("expected TransactionID=%d, got %d", txID1, p.TransactionID())
		}
	}
}

func TestPostingQueryBuilder_AmountPreserved(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	db.BeginTx(ctx)
	acctID := seedAccount(t, ctx, db, "Checking", domain.Asset)
	txID := seedTransaction(t, ctx, db, "Paycheck")
	seedPosting(t, ctx, db, txID, acctID, "1234.56")
	db.Commit()

	postings, err := db.PostingQueryBuilder().
		SetCondition(NewWhere(Column_TransactionId, Equal, idStr(txID))).
		Select(ctx, 0, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(postings) != 1 {
		t.Fatalf("expected 1 posting, got %d", len(postings))
	}
	want := decimal.RequireFromString("1234.56")
	if !postings[0].Amount().Equal(want) {
		t.Errorf("Amount() = %s, want %s", postings[0].Amount(), want)
	}
}

func TestPostingQueryBuilder_FilterByAccountID(t *testing.T) {
	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	db.BeginTx(ctx)
	acct1 := seedAccount(t, ctx, db, "Checking", domain.Asset)
	acct2 := seedAccount(t, ctx, db, "Savings", domain.Asset)
	txID := seedTransaction(t, ctx, db, "Transfer")
	seedPosting(t, ctx, db, txID, acct1, "-200.00")
	seedPosting(t, ctx, db, txID, acct2, "200.00")
	db.Commit()

	postings, err := db.PostingQueryBuilder().
		SetCondition(NewWhere(Column_AccountId, Equal, idStr(acct2))).
		Select(ctx, 0, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(postings) != 1 {
		t.Fatalf("expected 1 posting for acct2, got %d", len(postings))
	}
	if postings[0].AccountID() != acct2 {
		t.Errorf("AccountID() = %d, want %d", postings[0].AccountID(), acct2)
	}
}

// ── Repository interface compliance ──────────────────────────────────────────

// This is a compile-time check. If SQLiteDB stops satisfying Repository, this fails.
func TestSQLiteDB_ImplementsRepository(t *testing.T) {
	db := newTestDB(t)
	var _ Repository = db
	_ = db
}
