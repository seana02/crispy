package db

import (
	"context"
	"crispy/currency"
	"crispy/domain"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shopspring/decimal"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

type SQLiteDB struct {
	db *sql.DB
	tx *sql.Tx
}

func InitSQLite(name string) (*SQLiteDB, error) {
	handle, err := getHandle(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get database handle: %v", err)
	}

	return &SQLiteDB{db: handle}, nil
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
	return nil
}

func (s *SQLiteDB) CreateTransaction(ctx context.Context, t *domain.Transaction) (int64, error) {
	const funcError = "Error in CreateTransaction"
	if s.tx == nil {
		return -1, fmt.Errorf("%s: no Tx started. Call SQLiteDB.BeginTx", funcError)
	}
	now := time.Now().Local().Truncate(time.Second)
	result, err := s.tx.ExecContext(ctx, "INSERT INTO \"Transaction\" (description, date, status, date_created, date_updated) VALUES (?, ?, ?, ?, ?)",
		t.Description(),
		t.Date(),
		t.Status(),
		now,
		now,
	)
	// TODO: TAGS!!!!
	if err != nil {
		return -1, fmt.Errorf("%s inserting transaction: %s", funcError, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s getting LastInsertID: %s", funcError, err)
	}
	return id, nil
}

func (s *SQLiteDB) GetTransactionById(ctx context.Context, id int64) (*domain.Transaction, error) {
	const funcError = "Error in GetTransactionById"
	var description string
	var date time.Time
	var status domain.Status
	var referenceID *int64
	var dateCreated, dateUpdated time.Time
	var created, updated sql.NullTime
	err := s.db.QueryRowContext(ctx, "SELECT * FROM \"Transaction\" WHERE id = ?", id).Scan(
		&id, &description, &date, &status, &referenceID, &created, &updated,
	)
	if err != nil {
		return nil, fmt.Errorf("%s creating query: %s", funcError, err)
	}
	postings, err := s.GetPostingsByTransactionId(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s getting postings: %s", funcError, err)
	}
	tags, err := s.getTagsById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s getting tags: %s", funcError, err)
	}
	if created.Valid {
		dateCreated = created.Time
	}
	if updated.Valid {
		dateUpdated = created.Time
	}
	return domain.NewTransaction(
		id, description, date, status, referenceID, dateCreated, dateUpdated, postings, tags,
	), nil
}

func (s *SQLiteDB) GetTransactionByTag(ctx context.Context, tagID int64) ([]*domain.Transaction, error) {
	const funcError = "Error in GetTransactionByTag"
	var id int64
	var description string
	var date time.Time
	var status domain.Status
	var referenceID *int64
	var dateCreated, dateUpdated time.Time
	rows, err := s.db.QueryContext(ctx, "SELECT \"Transaction\".* FROM \"Transaction\" JOIN Transaction_Tag ON Transaction_Tag.transaction_id = \"Transaction\".id WHERE Transaction_Tag.tag_id = ?", tagID)
	if err != nil {
		return nil, fmt.Errorf("%s creating query: %s", funcError, err)
	}

	var ts []*domain.Transaction
	for rows.Next() {
		rows.Scan(
			&id, &description, &date, &status, &referenceID, &dateCreated, &dateUpdated,
		)
		postings, err := s.GetPostingsByTransactionId(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("%s getting postings: %s", funcError, err)
		}
		tags, err := s.getTagsById(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("%s getting tags: %s", funcError, err)
		}
		newTx := domain.NewTransaction(
			id, description, date, status, referenceID, dateCreated, dateUpdated, postings, tags,
		)
		ts = append(ts, newTx)
	}

	return ts, nil
}

func (s *SQLiteDB) GetTransactionByDate(ctx context.Context, from time.Time, to time.Time) ([]*domain.Transaction, error) {
	const funcError = "Error in GetTransactionByDate"
	var id int64
	var description string
	var date time.Time
	var status domain.Status
	var referenceID *int64
	var dateCreated, dateUpdated time.Time
	// rows, err := s.db.QueryContext(ctx, "SELECT * FROM \"Transaction\" WHERE date(\"date\") >= ? and date(\"date\") < ?", from.Format("2006-01-02"), to.Format("2006-01-02"))
	rows, err := s.db.QueryContext(ctx, "SELECT * FROM \"Transaction\"")
	if err != nil {
		return nil, fmt.Errorf("%s creating query: %s", funcError, err)
	}

	var ts []*domain.Transaction
	for rows.Next() {
		rows.Scan(
			&id, &description, &date, &status, &referenceID, &dateCreated, &dateUpdated,
		)
		postings, err := s.GetPostingsByTransactionId(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("%s getting postings: %s", funcError, err)
		}
		tags, err := s.getTagsById(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("%s getting tags: %s", funcError, err)
		}
		newTx := domain.NewTransaction(
			id, description, date, status, referenceID, dateCreated, dateUpdated, postings, tags,
		)
		ts = append(ts, newTx)
	}

	return ts, nil
}

func (s *SQLiteDB) UpdateTransaction(ctx context.Context, t *domain.Transaction) (int64, error) {
	const funcError = "Error in UpdateTransaction"
	if s.tx == nil {
		return -1, fmt.Errorf("%s: no Tx started. Call SQLiteDB.BeginTx", funcError)
	}

	now := time.Now().Local().Truncate(time.Second)
	_, err := s.tx.ExecContext(ctx, "UPDATE \"Transaction\" SET description = ?, date = ?, status = ?, date_updated = ? WHERE id = ?",
		t.Description(),
		t.Date(),
		t.Status(),
		now,
		t.ID(),
	)
	// TODO: TAGS!!!!!
	if err != nil {
		return -1, fmt.Errorf("%s: %s", funcError, err)
	}
	return t.ID(), nil
}

func (s *SQLiteDB) DeleteTransaction(ctx context.Context, id int64) error {
	const funcError = "Error in DeleteTransction"
	if s.tx == nil {
		return fmt.Errorf("%s: no Tx started. Call SQLiteDB.BeginTx", funcError)
	}

	_, err := s.tx.ExecContext(ctx, "DELETE FROM \"Transaction\" WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("%s: %s", funcError, err)
	}
	return nil
}

func (s *SQLiteDB) getTagsById(ctx context.Context, id int64) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT * FROM Transaction_Tags WHERE transaction_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var t_id int64
	var name string
	var tags []string
	for rows.Next() {
		err := rows.Scan(&t_id, &name)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan tag name: %s", err)
		}
		tags = append(tags, name)
	}
	return tags, nil
}

func (s *SQLiteDB) addTagsToTransaction(ctx context.Context, transactionID int64, tags []string) error {
	const funcError = "Error in addTagsToTransaction"
	if s.tx == nil {
		return fmt.Errorf("%s: no Tx started. Call SQLiteDB.BeginTx", funcError)
	}
	for _, tag := range tags {
		query, err := s.tx.QueryContext(ctx, "SELECT id FROM Tag WHERE name = ?", tag)
		if err != nil {
			return fmt.Errorf("%s querying for tag \"%s\"", funcError, tag)
		}
		var tagId int64
		if query.Next() {
			err := query.Scan(&tagId)
			if err != nil {
				return fmt.Errorf("%s failed to get tag ID for tag \"%s\"", funcError, tag)
			}
		} else {
			result, err := s.tx.ExecContext(ctx, "INSERT INTO Tag (name) VALUES (?)", tag)
			if err != nil {
				return fmt.Errorf("%s failed to insert new tag \"%s\"", funcError, tag)
			}
			tagId, err = result.LastInsertId()
			if err != nil {
				return fmt.Errorf("%s failed to get last insert id for new tag \"%s\"", funcError, tag)
			}
		}
		s.tx.ExecContext(ctx, "INSERT INTO Transaction_Tag (transaction_id, tag_id) VALUES (?, ?)", transactionID, tagId)
	}
	return nil
}

func (s *SQLiteDB) CreateAccount(ctx context.Context, a *domain.Account) (int64, error) {
	const funcError = "Error in CreateAccount"
	if s.tx == nil {
		return -1, fmt.Errorf("%s: no Tx started. Call SQLiteDB.BeginTx", funcError)
	}
	now := time.Now().Local().Truncate(time.Second)
	result, err := s.tx.ExecContext(ctx, "INSERT INTO Account (parent_id, name, type, currency, description, active, date_created, date_updated) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		a.ParentID(),
		a.Name(),
		a.Type(),
		a.Currency(),
		a.Description(),
		a.Active(),
		now,
		now,
	)
	if err != nil {
		return -1, fmt.Errorf("%s creating statement: %s", funcError, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s getting new account id: %s", funcError, err)
	}
	return id, nil
}

func (s *SQLiteDB) GetAccountById(ctx context.Context, id int64) (*domain.Account, error) {
	const funcError = "Error in GetAccountById"
	var parentID int64
	var name string
	var type_ domain.Type
	var currency string
	var description string
	var active bool
	var dateCreated, dateUpdated time.Time
	err := s.db.QueryRowContext(ctx, "SELECT * FROM \"Account\" WHERE id = ?", id).Scan(
		&id, &parentID, &name, &type_, &currency, &description, &active, &dateCreated, &dateUpdated,
	)
	if err != nil {
		return nil, fmt.Errorf("%s creating query: %s", funcError, err)
	}
	return domain.NewAccount(
		id, parentID, name, type_, currency, description, active, dateCreated, dateUpdated,
	)
}

func (s *SQLiteDB) GetAccountByFullName(ctx context.Context, fullName string) (*domain.Account, error) {
	const funcError = "Error in GetAccountByFullName"
	var id int64
	var parentID int64
	var name string
	var type_ domain.Type
	var currency string
	var description string
	var active bool
	var dateCreated, dateUpdated time.Time

	query := `
	WITH RECURSIVE
		path_parts(step, name) AS (
			SELECT key+1, value FROM json_each(?)
		),
		path_len AS (
			SELECT COUNT(*) as total FROM path_parts
		),
		traversal(id, name, step) AS (
			SELECT a.id, a.name, 1
			FROM "Account" a
			WHERE a.name = (SELECT name FROM path_parts WHERE step = 1)
			  AND a.parent_id = 0
		
			UNION ALL

			SELECT a.id, a.name, p.step+1
			FROM "Account" a
			JOIN traversal p ON a.parent_id = p.id
			JOIN path_parts pp ON a.name = pp.name AND pp.step = p.step+1
		)
	SELECT a.* FROM "Account" a
	JOIN traversal t ON t.id = a.id
	WHERE step = (SELECT total FROM path_len)
	ORDER BY step DESC
	LIMIT 1;
	`

	parts := strings.Split(fullName, ":")
	pathJson, _ := json.Marshal(parts)

	err := s.db.QueryRowContext(ctx, query, pathJson).Scan(
		&id, &parentID, &name, &type_, &currency, &description, &active, &dateCreated, &dateUpdated,
	)
	// err := s.db.QueryRowContext(ctx, query, pathJson).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("Error finding account %s: %s", fullName, err)
	}

	return domain.NewAccount(
		id, parentID, name, type_, currency, description, active, dateCreated, dateUpdated,
	)
	// return s.GetAccountById(ctx, id)
}

func (s *SQLiteDB) UpdateAccount(ctx context.Context, a *domain.Account) error {
	const funcError = "Error in UpdateAccount"
	if s.tx == nil {
		return fmt.Errorf("%s: no Tx started. Call SQLiteDB.BeginTx", funcError)
	}
	now := time.Now().Local().Truncate(time.Second)
	_, err := s.tx.ExecContext(ctx, "UPDATE \"Posting\" SET parent_id = ?, name = ?, type = ?, currency = ?, description = ?, active = ?, date_updated = ? WHERE id = ?",
		a.ParentID(),
		a.Name(),
		a.Type(),
		a.Currency(),
		a.Description(),
		a.Active(),
		now,
		a.ID(),
	)
	if err != nil {
		return fmt.Errorf("%s updating row: %s", funcError, err)
	}
	return nil
}

func (s *SQLiteDB) DeleteAccount(ctx context.Context, id int64) error {
	const funcError = "Error in UpdateAccount"
	if s.tx == nil {
		return fmt.Errorf("%s: no Tx started. Call SQLiteDB.BeginTx", funcError)
	}
	_, err := s.tx.ExecContext(ctx, "DELETE FROM \"Posting\" WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("%s deleting row: %s", funcError, err)
	}
	return nil
}

func (s *SQLiteDB) CreatePosting(ctx context.Context, p *domain.Posting) (int64, error) {
	const funcError = "Error in CreatePosting"
	if s.tx == nil {
		return -1, fmt.Errorf("%s: no Tx started. Call SQLiteDB.BeginTx", funcError)
	}
	now := time.Now().Local().Truncate(time.Second)
	result, err := s.tx.ExecContext(ctx, "INSERT INTO Posting (transaction_id, account_id, amount, currency, date_created, date_updated) VALUES (?, ?, ?, ?, ?, ?)",
		p.TransactionID(),
		p.AccountID(),
		p.Amount().StringFixed(currency.PrecisionFor(p.Currency())),
		p.Currency(),
		now,
		now,
	)
	if err != nil {
		return -1, fmt.Errorf("%s inserting posting: %s", funcError, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s getting LastInsertID: %s", funcError, err)
	}
	return id, nil
}

func (s *SQLiteDB) GetPostingById(ctx context.Context, id int64) (*domain.Posting, error) {
	const funcError = "Error in GetPostingsByID"
	var postingID, transactionID, accountID int64
	var amount decimal.Decimal
	var currency string
	var dateCreated, dateUpdated time.Time
	err := s.db.QueryRowContext(ctx, "SELECT * FROM Posting WHERE id = ?", id).Scan(
		&postingID, &transactionID, &accountID, &amount, &currency, &dateCreated, &dateUpdated,
	)
	if err != nil {
		return nil, fmt.Errorf("%s getting row: %s", funcError, err)
	}
	return domain.NewPosting(postingID, transactionID, accountID, amount, currency, dateCreated, dateUpdated), nil
}

func (s *SQLiteDB) GetPostingsByTransactionId(ctx context.Context, id int64) ([]*domain.Posting, error) {
	const funcError = "Error in GetPostingsByTransactionId"
	rows, err := s.db.QueryContext(ctx, "SELECT * FROM Posting WHERE transaction_id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("%s creating query: %s", funcError, err)
	}
	defer rows.Close()
	var p_id, transactionID, accountID int64
	var amount decimal.Decimal
	var currency string
	var dateCreated, dateUpdated time.Time
	var postings []*domain.Posting
	for rows.Next() {
		err := rows.Scan(&p_id, &transactionID, &accountID, &amount, &currency, &dateCreated, &dateUpdated)
		if err != nil {
			return nil, fmt.Errorf("%s scanning posting row: %s", funcError, err)
		}
		postings = append(postings, domain.NewPosting(
			p_id, transactionID, accountID, amount, currency, dateCreated, dateUpdated,
		))
	}
	return postings, nil
}

func (s *SQLiteDB) UpdatePosting(ctx context.Context, p *domain.Posting) error {
	const funcError = "Error in UpdatePosting"
	if s.tx == nil {
		return fmt.Errorf("%s: no Tx started. Call SQLiteDB.BeginTx", funcError)
	}

	now := time.Now().Local().Truncate(time.Second)
	_, err := s.tx.ExecContext(ctx, "UPDATE Posting SET transaction_id = ?, account_id = ?, amount = ?, currency = ?, date_updated = ? WHERE id = ?",
		p.TransactionID(),
		p.AccountID(),
		p.Amount(),
		p.Currency(),
		now,
		p.ID(),
	)
	if err != nil {
		return fmt.Errorf("%s: %s", funcError, err)
	}
	return nil
}

func (s *SQLiteDB) DeletePosting(ctx context.Context, id int64) error {
	const funcError = "Error in DeletePosting"
	if s.tx == nil {
		return fmt.Errorf("%s: no Tx started. Call SQLiteDB.BeginTx", funcError)
	}

	_, err := s.tx.ExecContext(ctx, "DELETE FROM Posting WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("%s: %s", funcError, err)
	}
	return nil
}
