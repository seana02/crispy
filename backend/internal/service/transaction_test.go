package service

import (
	"fmt"
	"testing"

	"github.com/crispy-2/backend/internal/model"
	"github.com/crispy-2/backend/internal/testutil"
	"github.com/shopspring/decimal"
)

func setupTransactionTestAccounts(t *testing.T, svc *AccountService) (checking, expenses *model.Account) {
	t.Helper()
	assets, err := svc.CreateAccount("Assets", nil, "asset")
	if err != nil {
		t.Fatalf("failed to create Assets: %v", err)
	}
	c, err := svc.CreateAccount("Checking", &assets.ID, "")
	if err != nil {
		t.Fatalf("failed to create Checking: %v", err)
	}
	e, err := svc.CreateAccount("Expenses", nil, "expense")
	if err != nil {
		t.Fatalf("failed to create Expenses: %v", err)
	}
	return c, e
}

func TestCreateTransaction_Success(t *testing.T) {
	db := testutil.NewTestDB(t)
	acctSvc := NewAccountService(db)
	txSvc := NewTransactionService(db, nil)

	checking, expenses := setupTransactionTestAccounts(t, acctSvc)

	txn, err := txSvc.CreateTransaction("2026-01-15", "Groceries", "posted", []PostingInput{
		{AccountID: checking.ID, Currency: "USD", Amount: "-50.00"},
		{AccountID: expenses.ID, Currency: "USD", Amount: "50.00"},
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if txn.Date != "2026-01-15" {
		t.Errorf("expected date 2026-01-15, got %q", txn.Date)
	}
	if txn.Description != "Groceries" {
		t.Errorf("expected description 'Groceries', got %q", txn.Description)
	}
	if len(txn.Postings) != 2 {
		t.Fatalf("expected 2 postings, got %d", len(txn.Postings))
	}
	testutil.AssertDecimalEqual(t, decimal.RequireFromString("-50.00"), txn.Postings[0].Amount.Decimal)
	testutil.AssertDecimalEqual(t, decimal.RequireFromString("50.00"), txn.Postings[1].Amount.Decimal)
}

func TestCreateTransaction_AutoInferLastPosting(t *testing.T) {
	db := testutil.NewTestDB(t)
	acctSvc := NewAccountService(db)
	txSvc := NewTransactionService(db, nil)

	checking, expenses := setupTransactionTestAccounts(t, acctSvc)

	txn, err := txSvc.CreateTransaction("2026-01-15", "Groceries", "posted", []PostingInput{
		{AccountID: checking.ID, Currency: "USD", Amount: "-50.00"},
		{AccountID: expenses.ID, Currency: "USD", Amount: ""},
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// The blank posting should be inferred as 50.00
	testutil.AssertDecimalEqual(t, decimal.RequireFromString("50.00"), txn.Postings[1].Amount.Decimal)
}

func TestCreateTransaction_UnbalancedRejected(t *testing.T) {
	db := testutil.NewTestDB(t)
	acctSvc := NewAccountService(db)
	txSvc := NewTransactionService(db, nil)

	checking, expenses := setupTransactionTestAccounts(t, acctSvc)

	_, err := txSvc.CreateTransaction("2026-01-15", "Bad", "posted", []PostingInput{
		{AccountID: checking.ID, Currency: "USD", Amount: "-50.00"},
		{AccountID: expenses.ID, Currency: "USD", Amount: "30.00"},
	}, nil)
	if err == nil {
		t.Fatal("expected error for unbalanced transaction")
	}
	appErr, ok := err.(*model.AppError)
	if !ok {
		t.Fatalf("expected *model.AppError, got %T", err)
	}
	if appErr.Code != 400 {
		t.Errorf("expected code 400, got %d", appErr.Code)
	}
}

func TestCreateTransaction_LessThanTwoPostingsRejected(t *testing.T) {
	db := testutil.NewTestDB(t)
	txSvc := NewTransactionService(db, nil)

	_, err := txSvc.CreateTransaction("2026-01-15", "Bad", "posted", []PostingInput{
		{AccountID: 1, Currency: "USD", Amount: "50.00"},
	}, nil)
	if err == nil {
		t.Fatal("expected error for less than two postings")
	}
}

func TestCreateTransaction_MultiCurrencyAutoInferRejected(t *testing.T) {
	db := testutil.NewTestDB(t)
	acctSvc := NewAccountService(db)
	txSvc := NewTransactionService(db, nil)

	checking, expenses := setupTransactionTestAccounts(t, acctSvc)

	_, err := txSvc.CreateTransaction("2026-01-15", "Multi", "posted", []PostingInput{
		{AccountID: checking.ID, Currency: "USD", Amount: "-50.00"},
		{AccountID: expenses.ID, Currency: "JPY", Amount: ""},
	}, nil)
	if err == nil {
		t.Fatal("expected error for auto-infer on multi-currency")
	}
}

func TestCreateTransaction_WithTags(t *testing.T) {
	db := testutil.NewTestDB(t)
	acctSvc := NewAccountService(db)
	txSvc := NewTransactionService(db, nil)

	checking, expenses := setupTransactionTestAccounts(t, acctSvc)

	txn, err := txSvc.CreateTransaction("2026-01-15", "Japan trip meal", "posted", []PostingInput{
		{AccountID: checking.ID, Currency: "USD", Amount: "-25.00"},
		{AccountID: expenses.ID, Currency: "USD", Amount: "25.00"},
	}, []TagInput{
		{Key: "trip", Value: "japan-2026"},
		{Key: "category", Value: "food"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(txn.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(txn.Tags))
	}
}

func TestCreateTransaction_InvalidStatusRejected(t *testing.T) {
	db := testutil.NewTestDB(t)
	txSvc := NewTransactionService(db, nil)

	_, err := txSvc.CreateTransaction("2026-01-15", "Bad", "invalid", []PostingInput{
		{AccountID: 1, Currency: "USD", Amount: "-50.00"},
		{AccountID: 2, Currency: "USD", Amount: "50.00"},
	}, nil)
	if err == nil {
		t.Fatal("expected error for invalid status")
	}
}

func TestCreateTransaction_EmptyDateRejected(t *testing.T) {
	db := testutil.NewTestDB(t)
	txSvc := NewTransactionService(db, nil)

	_, err := txSvc.CreateTransaction("", "Bad", "posted", []PostingInput{
		{AccountID: 1, Currency: "USD", Amount: "-50.00"},
		{AccountID: 2, Currency: "USD", Amount: "50.00"},
	}, nil)
	if err == nil {
		t.Fatal("expected error for empty date")
	}
}

func TestGetTransaction(t *testing.T) {
	db := testutil.NewTestDB(t)
	acctSvc := NewAccountService(db)
	txSvc := NewTransactionService(db, nil)

	checking, expenses := setupTransactionTestAccounts(t, acctSvc)

	created, _ := txSvc.CreateTransaction("2026-01-15", "Groceries", "posted", []PostingInput{
		{AccountID: checking.ID, Currency: "USD", Amount: "-50.00"},
		{AccountID: expenses.ID, Currency: "USD", Amount: "50.00"},
	}, []TagInput{{Key: "trip", Value: "japan-2026"}})

	got, err := txSvc.GetTransaction(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Description != "Groceries" {
		t.Errorf("expected 'Groceries', got %q", got.Description)
	}
	if len(got.Postings) != 2 {
		t.Errorf("expected 2 postings, got %d", len(got.Postings))
	}
	if len(got.Tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(got.Tags))
	}
}

func TestGetTransaction_NotFound(t *testing.T) {
	db := testutil.NewTestDB(t)
	txSvc := NewTransactionService(db, nil)

	_, err := txSvc.GetTransaction(9999)
	if err == nil {
		t.Fatal("expected not found error")
	}
}

func TestListTransactions(t *testing.T) {
	db := testutil.NewTestDB(t)
	acctSvc := NewAccountService(db)
	txSvc := NewTransactionService(db, nil)

	checking, expenses := setupTransactionTestAccounts(t, acctSvc)

	txSvc.CreateTransaction("2026-01-15", "Groceries", "posted", []PostingInput{
		{AccountID: checking.ID, Currency: "USD", Amount: "-50.00"},
		{AccountID: expenses.ID, Currency: "USD", Amount: "50.00"},
	}, nil)
	txSvc.CreateTransaction("2026-01-16", "Coffee", "posted", []PostingInput{
		{AccountID: checking.ID, Currency: "USD", Amount: "-5.00"},
		{AccountID: expenses.ID, Currency: "USD", Amount: "5.00"},
	}, nil)

	result, err := txSvc.ListTransactions(TransactionFilters{}, model.Pagination{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 2 {
		t.Errorf("expected 2 transactions, got %d", result.Total)
	}
}

func TestListTransactions_FilterByAccount(t *testing.T) {
	db := testutil.NewTestDB(t)
	acctSvc := NewAccountService(db)
	txSvc := NewTransactionService(db, nil)

	checking, expenses := setupTransactionTestAccounts(t, acctSvc)
	income, _ := acctSvc.CreateAccount("Income", nil, "revenue")

	// Transaction 1: checking -> expenses
	txSvc.CreateTransaction("2026-01-15", "Groceries", "posted", []PostingInput{
		{AccountID: checking.ID, Currency: "USD", Amount: "-50.00"},
		{AccountID: expenses.ID, Currency: "USD", Amount: "50.00"},
	}, nil)
	// Transaction 2: income -> checking
	txSvc.CreateTransaction("2026-01-16", "Salary", "posted", []PostingInput{
		{AccountID: income.ID, Currency: "USD", Amount: "-3000.00"},
		{AccountID: checking.ID, Currency: "USD", Amount: "3000.00"},
	}, nil)

	// Filter by expenses account should return 1 transaction
	result, err := txSvc.ListTransactions(TransactionFilters{AccountID: &expenses.ID}, model.Pagination{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 1 {
		t.Errorf("expected 1 transaction for expenses, got %d", result.Total)
	}
}

func TestUpdateTransaction(t *testing.T) {
	db := testutil.NewTestDB(t)
	acctSvc := NewAccountService(db)
	txSvc := NewTransactionService(db, nil)

	checking, expenses := setupTransactionTestAccounts(t, acctSvc)

	created, _ := txSvc.CreateTransaction("2026-01-15", "Groceries", "posted", []PostingInput{
		{AccountID: checking.ID, Currency: "USD", Amount: "-50.00"},
		{AccountID: expenses.ID, Currency: "USD", Amount: "50.00"},
	}, nil)

	updated, err := txSvc.UpdateTransaction(created.ID, "2026-01-16", "Updated Groceries", "cleared", []PostingInput{
		{AccountID: checking.ID, Currency: "USD", Amount: "-75.00"},
		{AccountID: expenses.ID, Currency: "USD", Amount: "75.00"},
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Description != "Updated Groceries" {
		t.Errorf("expected 'Updated Groceries', got %q", updated.Description)
	}
	if updated.Status != "cleared" {
		t.Errorf("expected status 'cleared', got %q", updated.Status)
	}
	testutil.AssertDecimalEqual(t, decimal.RequireFromString("-75.00"), updated.Postings[0].Amount.Decimal)
}

func TestDeleteTransaction(t *testing.T) {
	db := testutil.NewTestDB(t)
	acctSvc := NewAccountService(db)
	txSvc := NewTransactionService(db, nil)

	checking, expenses := setupTransactionTestAccounts(t, acctSvc)

	created, _ := txSvc.CreateTransaction("2026-01-15", "Groceries", "posted", []PostingInput{
		{AccountID: checking.ID, Currency: "USD", Amount: "-50.00"},
		{AccountID: expenses.ID, Currency: "USD", Amount: "50.00"},
	}, nil)

	err := txSvc.DeleteTransaction(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = txSvc.GetTransaction(created.ID)
	if err == nil {
		t.Fatal("expected not found after deletion")
	}
}

func TestCreateTransaction_MultipleBlankAmountsRejected(t *testing.T) {
	db := testutil.NewTestDB(t)
	acctSvc := NewAccountService(db)
	txSvc := NewTransactionService(db, nil)

	checking, expenses := setupTransactionTestAccounts(t, acctSvc)
	income, _ := acctSvc.CreateAccount("Income", nil, "revenue")

	_, err := txSvc.CreateTransaction("2026-01-15", "Bad", "posted", []PostingInput{
		{AccountID: checking.ID, Currency: "USD", Amount: "50.00"},
		{AccountID: expenses.ID, Currency: "USD", Amount: ""},
		{AccountID: income.ID, Currency: "USD", Amount: ""},
	}, nil)
	if err == nil {
		t.Fatal("expected error for multiple blank amounts")
	}
}

func TestCreateTransaction_MultiCurrencyBalancedWithRate(t *testing.T) {
	db := testutil.NewTestDB(t)
	acctSvc := NewAccountService(db)
	rateSvc := NewExchangeRateService(db)
	txSvc := NewTransactionService(db, rateSvc)

	checking, expenses := setupTransactionTestAccounts(t, acctSvc)

	_, err := rateSvc.CreateRate("USD", "JPY", "150.00", "2026-01-15", "manual")
	if err != nil {
		t.Fatalf("failed to create exchange rate: %v", err)
	}

	_, err = txSvc.CreateTransaction("2026-01-15", "Multi-currency balanced", "posted", []PostingInput{
		{AccountID: checking.ID, Currency: "USD", Amount: "-100.00"},
		{AccountID: expenses.ID, Currency: "JPY", Amount: "15000.00"},
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error for balanced multi-currency transaction: %v", err)
	}
}

func TestCreateTransaction_MultiCurrencyUnbalancedRejected(t *testing.T) {
	db := testutil.NewTestDB(t)
	acctSvc := NewAccountService(db)
	rateSvc := NewExchangeRateService(db)
	txSvc := NewTransactionService(db, rateSvc)

	checking, expenses := setupTransactionTestAccounts(t, acctSvc)

	_, err := rateSvc.CreateRate("USD", "JPY", "150.00", "2026-01-15", "manual")
	if err != nil {
		t.Fatalf("failed to create exchange rate: %v", err)
	}

	_, err = txSvc.CreateTransaction("2026-01-15", "Multi-currency unbalanced", "posted", []PostingInput{
		{AccountID: checking.ID, Currency: "USD", Amount: "-100.00"},
		{AccountID: expenses.ID, Currency: "JPY", Amount: "14000.00"},
	}, nil)
	if err == nil {
		t.Fatal("expected error for unbalanced multi-currency transaction")
	}
	appErr, ok := err.(*model.AppError)
	if !ok {
		t.Fatalf("expected *model.AppError, got %T", err)
	}
	if appErr.Code != 400 {
		t.Errorf("expected code 400, got %d", appErr.Code)
	}
}

func TestCreateTransaction_MultiCurrencyMissingRateRejected(t *testing.T) {
	db := testutil.NewTestDB(t)
	acctSvc := NewAccountService(db)
	rateSvc := NewExchangeRateService(db)
	txSvc := NewTransactionService(db, rateSvc)

	checking, expenses := setupTransactionTestAccounts(t, acctSvc)

	_, err := txSvc.CreateTransaction("2026-01-15", "Multi-currency no rate", "posted", []PostingInput{
		{AccountID: checking.ID, Currency: "USD", Amount: "-100.00"},
		{AccountID: expenses.ID, Currency: "EUR", Amount: "100.00"},
	}, nil)
	if err == nil {
		t.Fatal("expected error for missing exchange rate")
	}
	appErr, ok := err.(*model.AppError)
	if !ok {
		t.Fatalf("expected *model.AppError, got %T", err)
	}
	if appErr.Code != 404 {
		t.Errorf("expected code 404 for missing rate, got %d", appErr.Code)
	}
}

// Regression: GetTransaction must return account_full_name on postings, not a raw ID.
func TestGetTransaction_PostingsHaveAccountFullName(t *testing.T) {
	db := testutil.NewTestDB(t)
	acctSvc := NewAccountService(db)
	rateSvc := NewExchangeRateService(db)
	txSvc := NewTransactionService(db, rateSvc)

	checking, expenses := setupTransactionTestAccounts(t, acctSvc)

	created, err := txSvc.CreateTransaction("2026-01-15", "name test", "cleared", []PostingInput{
		{AccountID: checking.ID, Currency: "USD", Amount: "-50.00"},
		{AccountID: expenses.ID, Currency: "USD", Amount: "50.00"},
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tx, err := txSvc.GetTransaction(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tx.Postings) != 2 {
		t.Fatalf("expected 2 postings, got %d", len(tx.Postings))
	}
	for _, p := range tx.Postings {
		if p.AccountFullName == "" {
			t.Errorf("posting account_full_name is empty for account_id %d — expected name, got empty string", p.AccountID)
		}
		// Must not be the numeric ID stringified
		if p.AccountFullName == string(rune(p.AccountID)) || p.AccountFullName == fmt.Sprintf("%d", p.AccountID) {
			t.Errorf("posting account_full_name looks like a raw ID: %q", p.AccountFullName)
		}
	}
	if tx.Postings[0].AccountFullName != "Assets:Checking" {
		t.Errorf("expected 'Assets:Checking', got %q", tx.Postings[0].AccountFullName)
	}
	if tx.Postings[1].AccountFullName != "Expenses" {
		t.Errorf("expected 'Expenses', got %q", tx.Postings[1].AccountFullName)
	}
}

// Regression: edit form must show account full names, not numeric IDs.
// This tests the backend side: GetTransaction returns account_full_name.
func TestGetTransaction_EditFormReceivesAccountNames(t *testing.T) {
	db := testutil.NewTestDB(t)
	acctSvc := NewAccountService(db)
	rateSvc := NewExchangeRateService(db)
	txSvc := NewTransactionService(db, rateSvc)

	// Create a deeper hierarchy to verify full path
	assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
	checking, _ := acctSvc.CreateAccount("Checking", &assets.ID, "")
	food, _ := acctSvc.CreateAccount("Food", nil, "expense")

	created, _ := txSvc.CreateTransaction("2026-06-01", "lunch", "cleared", []PostingInput{
		{AccountID: checking.ID, Currency: "USD", Amount: "-20.00"},
		{AccountID: food.ID, Currency: "USD", Amount: "20.00"},
	}, nil)

	tx, err := txSvc.GetTransaction(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	names := make([]string, len(tx.Postings))
	for i, p := range tx.Postings {
		names[i] = p.AccountFullName
	}
	if names[0] != "Assets:Checking" {
		t.Errorf("expected 'Assets:Checking', got %q", names[0])
	}
	if names[1] != "Food" {
		t.Errorf("expected 'Food', got %q", names[1])
	}
}

// --- Tests for resolvePostingAccountID and auto-creation of accounts from transaction form ---

func TestResolvePostingAccountID_ExistingAccount(t *testing.T) {
	db := testutil.NewTestDB(t)
	acctSvc := NewAccountService(db)
	txSvc := NewTransactionService(db, nil)
	txSvc.SetAccountService(acctSvc)

	acct, err := acctSvc.CreateAccount("Assets", nil, "asset")
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}

	// Resolve by ID
	id, err := txSvc.resolvePostingAccountID(PostingInput{AccountID: acct.ID, Currency: "USD", Amount: "100"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != acct.ID {
		t.Errorf("expected ID %d, got %d", acct.ID, id)
	}
}

func TestResolvePostingAccountID_ExistingAccountByName(t *testing.T) {
	db := testutil.NewTestDB(t)
	acctSvc := NewAccountService(db)
	txSvc := NewTransactionService(db, nil)
	txSvc.SetAccountService(acctSvc)

	acct, err := acctSvc.CreateAccount("Assets", nil, "asset")
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}

	// Resolve by exact name
	id, err := txSvc.resolvePostingAccountID(PostingInput{AccountName: "Assets", Currency: "USD", Amount: "100"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != acct.ID {
		t.Errorf("expected ID %d, got %d", acct.ID, id)
	}
}

func TestResolvePostingAccountID_CreatesNewRootAccount(t *testing.T) {
	db := testutil.NewTestDB(t)
	acctSvc := NewAccountService(db)
	txSvc := NewTransactionService(db, nil)
	txSvc.SetAccountService(acctSvc)

	// "Groceries" doesn't exist — should be auto-created with a suggested type
	id, err := txSvc.resolvePostingAccountID(PostingInput{AccountName: "Groceries", Currency: "USD", Amount: "50"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == 0 {
		t.Fatal("expected non-zero account ID")
	}

	// Verify the account was created
	acct, err := acctSvc.GetAccount(id)
	if err != nil {
		t.Fatalf("failed to get auto-created account: %v", err)
	}
	if acct.FullName != "Groceries" {
		t.Errorf("expected full_name 'Groceries', got %q", acct.FullName)
	}
	// "Groceries" doesn't match any suggestion, so it gets the default "expense"
	if acct.AccountType != "expense" {
		t.Errorf("expected account_type 'expense' (default), got %q", acct.AccountType)
	}
}

func TestResolvePostingAccountID_CreatesNewRootWithSuggestedType(t *testing.T) {
	db := testutil.NewTestDB(t)
	acctSvc := NewAccountService(db)
	txSvc := NewTransactionService(db, nil)
	txSvc.SetAccountService(acctSvc)

	// "Assets" should be suggested as type "asset"
	id, err := txSvc.resolvePostingAccountID(PostingInput{AccountName: "Assets", Currency: "USD", Amount: "100"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	acct, err := acctSvc.GetAccount(id)
	if err != nil {
		t.Fatalf("failed to get account: %v", err)
	}
	if acct.AccountType != "asset" {
		t.Errorf("expected account_type 'asset', got %q", acct.AccountType)
	}
}

func TestResolvePostingAccountID_CreatesHierarchy(t *testing.T) {
	db := testutil.NewTestDB(t)
	acctSvc := NewAccountService(db)
	txSvc := NewTransactionService(db, nil)
	txSvc.SetAccountService(acctSvc)

	// "Expenses:Food:Restaurants" — nothing exists, should create all three
	id, err := txSvc.resolvePostingAccountID(PostingInput{AccountName: "Expenses:Food:Restaurants", Currency: "USD", Amount: "25"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == 0 {
		t.Fatal("expected non-zero account ID")
	}

	leaf, err := acctSvc.GetAccount(id)
	if err != nil {
		t.Fatalf("failed to get leaf account: %v", err)
	}
	if leaf.FullName != "Expenses:Food:Restaurants" {
		t.Errorf("expected full_name 'Expenses:Food:Restaurants', got %q", leaf.FullName)
	}
	if leaf.Name != "Restaurants" {
		t.Errorf("expected name 'Restaurants', got %q", leaf.Name)
	}

	// Verify the intermediate "Food" was created with correct parent
	results, err := acctSvc.SearchAccounts("Expenses:Food")
	if err != nil {
		t.Fatalf("failed to search: %v", err)
	}
	found := false
	for _, a := range results {
		if a.FullName == "Expenses:Food" {
			found = true
			if a.ParentID == nil {
				t.Error("expected Food to have a parent")
			}
		}
	}
	if !found {
		t.Error("intermediate account 'Expenses:Food' was not created")
	}

	// Verify the root "Expenses" was created with type "expense" (from SuggestAccountType)
	results, err = acctSvc.SearchAccounts("Expenses")
	if err != nil {
		t.Fatalf("failed to search: %v", err)
	}
	for _, a := range results {
		if a.FullName == "Expenses" {
			if a.AccountType != "expense" {
				t.Errorf("expected root 'Expenses' to have type 'expense', got %q", a.AccountType)
			}
			if a.ParentID != nil {
				t.Error("expected root 'Expenses' to have nil parent_id")
			}
		}
	}
}

func TestResolvePostingAccountID_PartialHierarchyExists(t *testing.T) {
	db := testutil.NewTestDB(t)
	acctSvc := NewAccountService(db)
	txSvc := NewTransactionService(db, nil)
	txSvc.SetAccountService(acctSvc)

	// Pre-create "Expenses" as root
	expenses, err := acctSvc.CreateAccount("Expenses", nil, "expense")
	if err != nil {
		t.Fatalf("failed to create Expenses: %v", err)
	}

	// Now resolve "Expenses:Food:Restaurants" — Expenses exists, Food and Restaurants don't
	id, err := txSvc.resolvePostingAccountID(PostingInput{AccountName: "Expenses:Food:Restaurants", Currency: "USD", Amount: "25"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	leaf, err := acctSvc.GetAccount(id)
	if err != nil {
		t.Fatalf("failed to get leaf: %v", err)
	}
	if leaf.FullName != "Expenses:Food:Restaurants" {
		t.Errorf("expected 'Expenses:Food:Restaurants', got %q", leaf.FullName)
	}

	// Verify Food was created under the pre-existing Expenses
	results, _ := acctSvc.SearchAccounts("Expenses:Food")
	for _, a := range results {
		if a.FullName == "Expenses:Food" {
			if a.ParentID == nil || *a.ParentID != expenses.ID {
				t.Errorf("expected Food's parent to be Expenses (id=%d), got %v", expenses.ID, a.ParentID)
			}
		}
	}
}

func TestResolvePostingAccountID_MissingNameAndID(t *testing.T) {
	db := testutil.NewTestDB(t)
	txSvc := NewTransactionService(db, nil)

	_, err := txSvc.resolvePostingAccountID(PostingInput{Currency: "USD", Amount: "10"})
	if err == nil {
		t.Fatal("expected error for posting with no account_id or account_name")
	}
}

func TestCreateTransaction_WithAutoCreatedAccounts(t *testing.T) {
	db := testutil.NewTestDB(t)
	acctSvc := NewAccountService(db)
	txSvc := NewTransactionService(db, nil)
	txSvc.SetAccountService(acctSvc)

	// Create transaction where both accounts are specified by name, not ID
	txn, err := txSvc.CreateTransaction("2026-03-01", "Auto-create test", "posted", []PostingInput{
		{AccountName: "Assets:Checking", Currency: "USD", Amount: "-100.00"},
		{AccountName: "Expenses:Food:Restaurants", Currency: "USD", Amount: "100.00"},
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(txn.Postings) != 2 {
		t.Fatalf("expected 2 postings, got %d", len(txn.Postings))
	}

	// Verify both accounts exist and have correct hierarchy
	got, err := txSvc.GetTransaction(txn.ID)
	if err != nil {
		t.Fatalf("failed to get transaction: %v", err)
	}
	if got.Postings[0].AccountFullName != "Assets:Checking" {
		t.Errorf("expected 'Assets:Checking', got %q", got.Postings[0].AccountFullName)
	}
	if got.Postings[1].AccountFullName != "Expenses:Food:Restaurants" {
		t.Errorf("expected 'Expenses:Food:Restaurants', got %q", got.Postings[1].AccountFullName)
	}
}
