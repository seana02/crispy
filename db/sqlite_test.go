package db

import (
	"context"
	"crispy/domain"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestGetHandle(t *testing.T) {
	filename := ":memory:"

	h, err := getHandle(filename)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if h == nil {
		t.Fatalf("expected non-nil handle, got nil")
	}

	if err := h.Ping(); err != nil {
		t.Fatalf("expected Ping to success, got %v", err)
	}

	h.Close()
}

func TestCreateTransaction(t *testing.T) {
	ctx := context.Background()
	filename := ":memory:"

	h, err := InitSQLite(filename)
	if err != nil {
		t.Fatalf("expected no error initializing sqlite, got %v", err)
	}

	currentTime := time.Now()

	// new transaction 1
	// new transaction 2 with duplicate values
	h.BeginTx(ctx)
	newTx, err := domain.NewTransaction(
		-1,
		"Tx1",
		currentTime,
		domain.Pending,
		-1,
		currentTime,
		currentTime,
		[]*domain.Posting{domain.NewPosting(-1, -1, 1, decimal.New(2, 0), "USD", time.Now(), time.Now())},
		[]string{"tag"},
	)
	if err != nil {
		t.Fatalf("expected no error creating transaction, got %v", err)
	}
	id1, err := h.CreateTransaction(ctx, newTx)
	if err != nil {
		t.Fatalf("expected no error inserting transaction, got %v", err)
	}
	if id1 != 1 {
		t.Errorf("expected transaction id of 1, got %d", id1)
	}

	t.Run("inserted 1 transaction", func (t *testing.T) {
		rows, err := h.tx.QueryContext(ctx, "SELECT * FROM \"Transaction\"")
		if err != nil {
			t.Errorf("expected successful query, got %v", err)
		}
		count := 0
		for rows.Next() {
			count += 1
			var id int64
			var description string
			var date time.Time
			var status domain.Status
			var referenceID int64
			var dateCreated, dateUpdated time.Time
			rows.Scan(&id, &description, &date, &status, &referenceID, &dateCreated, &dateUpdated)
			if id != 1 {
				t.Errorf("expected row id of 1, got %d", id)
			}
			if description != "Tx1" {
				t.Errorf("expected desc of \"Tx1\", got %s", description)
			}
			if date != currentTime {
				t.Errorf("expected time of %s, got %s", currentTime, date)
			}
			if status != domain.Pending {
				t.Errorf("expected status of %s, got %s", domain.Pending, status)
			}
			if referenceID != -1 {
				t.Errorf("expected ref id of -1, got %d", referenceID)
			}
			if dateCreated != currentTime {
				t.Errorf("expected date created of %s, got %s", currentTime, dateCreated)
			}
			if dateUpdated != currentTime {
				t.Errorf("expected date created of %s, got %s", currentTime, dateUpdated)
			}
		}
		if count != 1 {
			t.Errorf("expected 1 row, got %v", count)
		}
	})

	newTx, err = domain.NewTransaction(
		-1,
		"Tx2",
		currentTime,
		domain.Pending,
		-1,
		currentTime,
		currentTime,
		[]*domain.Posting{domain.NewPosting(-1, -1, 1, decimal.New(3, 1), "USD", time.Now(), time.Now())},
		[]string{"tag"},
	)
	if err != nil {
		t.Fatalf("expected no error creating transaction, got %v", err)
	}
	id2, err := h.CreateTransaction(ctx, newTx)
	if err != nil {
		t.Fatalf("expected no error inserting transaction, got %v", err)
	}
	if id2 != 2 {
		t.Errorf("expected transaction id of 2, got %d", id2)
	}

}
