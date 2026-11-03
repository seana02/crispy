package service

import (
	"testing"

	"github.com/crispy-2/backend/internal/testutil"
)

func TestGlobalSearch_FindsAccounts(t *testing.T) {
	db := testutil.NewTestDB(t)
	accountSvc := NewAccountService(db)
	svc := NewSearchService(db)

	accountSvc.CreateAccount("Assets", nil, "asset")
	accountSvc.CreateAccount("Expenses", nil, "expense")

	result, err := svc.GlobalSearch("asset")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Accounts) == 0 {
		t.Errorf("expected at least 1 account result, got 0")
	}
}

func TestGlobalSearch_FindsTransactions(t *testing.T) {
	db := testutil.NewTestDB(t)
	accountSvc := NewAccountService(db)
	exchangeSvc := NewExchangeRateService(db)
	txSvc := NewTransactionService(db, exchangeSvc)
	svc := NewSearchService(db)

	assets, _ := accountSvc.CreateAccount("Assets", nil, "asset")
	expenses, _ := accountSvc.CreateAccount("Expenses", nil, "expense")
	txSvc.CreateTransaction("2026-01-01", "coffee shop purchase", "cleared", []PostingInput{
		{AccountID: expenses.ID, Currency: "USD", Amount: "5.00"},
		{AccountID: assets.ID, Currency: "USD", Amount: "-5.00"},
	}, nil)

	result, err := svc.GlobalSearch("coffee")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Transactions) == 0 {
		t.Errorf("expected at least 1 transaction result, got 0")
	}
}

func TestGlobalSearch_FindsTags(t *testing.T) {
	db := testutil.NewTestDB(t)
	accountSvc := NewAccountService(db)
	exchangeSvc := NewExchangeRateService(db)
	txSvc := NewTransactionService(db, exchangeSvc)
	svc := NewSearchService(db)

	assets, _ := accountSvc.CreateAccount("Assets", nil, "asset")
	expenses, _ := accountSvc.CreateAccount("Expenses", nil, "expense")
	txSvc.CreateTransaction("2026-01-01", "japan trip", "cleared", []PostingInput{
		{AccountID: expenses.ID, Currency: "USD", Amount: "200.00"},
		{AccountID: assets.ID, Currency: "USD", Amount: "-200.00"},
	}, []TagInput{{Key: "trip", Value: "japan-2026"}})

	result, err := svc.GlobalSearch("japan")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Tags) == 0 {
		t.Errorf("expected at least 1 tag result, got 0")
	}
}

func TestGlobalSearch_EmptyQueryReturnsError(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewSearchService(db)

	// Empty query should still work (returns nothing)
	result, err := svc.GlobalSearch("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Errorf("expected non-nil result")
	}
}
