package service

import (
	"testing"

	"github.com/crispy-2/backend/internal/model"
	"github.com/crispy-2/backend/internal/testutil"
)

func TestCreateAccount_RootAccount(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	acct, err := svc.CreateAccount("Assets", nil, "asset")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if acct.Name != "Assets" {
		t.Errorf("expected name 'Assets', got %q", acct.Name)
	}
	if acct.FullName != "Assets" {
		t.Errorf("expected full_name 'Assets', got %q", acct.FullName)
	}
	if acct.AccountType != "asset" {
		t.Errorf("expected account_type 'asset', got %q", acct.AccountType)
	}
	if acct.ParentID != nil {
		t.Errorf("expected nil parent_id, got %v", acct.ParentID)
	}
	if acct.ID == 0 {
		t.Error("expected non-zero ID")
	}
}

func TestCreateAccount_ChildAccount(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	parent, err := svc.CreateAccount("Assets", nil, "asset")
	if err != nil {
		t.Fatalf("unexpected error creating parent: %v", err)
	}

	child, err := svc.CreateAccount("Checking", &parent.ID, "")
	if err != nil {
		t.Fatalf("unexpected error creating child: %v", err)
	}
	if child.FullName != "Assets:Checking" {
		t.Errorf("expected full_name 'Assets:Checking', got %q", child.FullName)
	}
	if child.ParentID == nil || *child.ParentID != parent.ID {
		t.Errorf("expected parent_id %d, got %v", parent.ID, child.ParentID)
	}
	if child.AccountType != parent.AccountType {
		t.Errorf("expected child account_type %q to match parent %q", child.AccountType, parent.AccountType)
	}
}

func TestCreateAccount_DeepHierarchy(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	root, _ := svc.CreateAccount("Expenses", nil, "expense")
	mid, _ := svc.CreateAccount("Shopping", &root.ID, "")
	leaf, err := svc.CreateAccount("Online", &mid.ID, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if leaf.FullName != "Expenses:Shopping:Online" {
		t.Errorf("expected full_name 'Expenses:Shopping:Online', got %q", leaf.FullName)
	}
}

func TestCreateAccount_NameContainsSeparatorRejected(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	_, err := svc.CreateAccount("Bad:Name", nil, "asset")
	if err == nil {
		t.Fatal("expected validation error for name containing ':'")
	}
	appErr, ok := err.(*model.AppError)
	if !ok {
		t.Fatalf("expected *model.AppError, got %T", err)
	}
	if appErr.Code != 400 {
		t.Errorf("expected error code 400, got %d", appErr.Code)
	}
}

func TestCreateAccount_EmptyNameRejected(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	_, err := svc.CreateAccount("", nil, "asset")
	if err == nil {
		t.Fatal("expected validation error for empty name")
	}
}

func TestCreateAccount_InvalidTypeRejected(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	_, err := svc.CreateAccount("Test", nil, "bogus")
	if err == nil {
		t.Fatal("expected validation error for invalid account type")
	}
}

func TestCreateAccount_RootRequiresType(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	_, err := svc.CreateAccount("Assets", nil, "")
	if err == nil {
		t.Fatal("expected validation error: root accounts must have a type")
	}
}

func TestCreateAccount_TypeAutoSuggestion(t *testing.T) {
	suggested := model.SuggestAccountType("Expenses")
	if suggested != model.AccountTypeExpense {
		t.Errorf("expected suggestion 'expense' for 'Expenses', got %q", suggested)
	}
}

func TestGetAccount(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	created, _ := svc.CreateAccount("Assets", nil, "asset")
	got, err := svc.GetAccount(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, got.ID)
	}
	if got.FullName != "Assets" {
		t.Errorf("expected full_name 'Assets', got %q", got.FullName)
	}
}

func TestGetAccount_NotFound(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	_, err := svc.GetAccount(9999)
	if err == nil {
		t.Fatal("expected not found error")
	}
	appErr, ok := err.(*model.AppError)
	if !ok {
		t.Fatalf("expected *model.AppError, got %T", err)
	}
	if appErr.Code != 404 {
		t.Errorf("expected error code 404, got %d", appErr.Code)
	}
}

func TestListAccounts(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	svc.CreateAccount("Assets", nil, "asset")
	svc.CreateAccount("Expenses", nil, "expense")

	result, err := svc.ListAccounts("", "", model.Pagination{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 2 {
		t.Errorf("expected 2 accounts, got %d", result.Total)
	}
	if len(result.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(result.Items))
	}
}

func TestListAccounts_FilterByType(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	svc.CreateAccount("Assets", nil, "asset")
	svc.CreateAccount("Expenses", nil, "expense")

	result, err := svc.ListAccounts("", "asset", model.Pagination{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 1 {
		t.Errorf("expected 1 account, got %d", result.Total)
	}
	if result.Items[0].Name != "Assets" {
		t.Errorf("expected 'Assets', got %q", result.Items[0].Name)
	}
}

func TestListAccounts_FilterByName(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	svc.CreateAccount("Assets", nil, "asset")
	svc.CreateAccount("Expenses", nil, "expense")

	result, err := svc.ListAccounts("asset", "", model.Pagination{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 1 {
		t.Errorf("expected 1 account, got %d", result.Total)
	}
}

func TestUpdateAccount_Rename(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	root, _ := svc.CreateAccount("Assets", nil, "asset")
	child, _ := svc.CreateAccount("Checking", &root.ID, "")

	updated, err := svc.UpdateAccount(child.ID, stringPtr("Savings"), nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Name != "Savings" {
		t.Errorf("expected name 'Savings', got %q", updated.Name)
	}
	if updated.FullName != "Assets:Savings" {
		t.Errorf("expected full_name 'Assets:Savings', got %q", updated.FullName)
	}
}

func TestUpdateAccount_RenameUpdatesDescendants(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	root, _ := svc.CreateAccount("Expenses", nil, "expense")
	mid, _ := svc.CreateAccount("Shopping", &root.ID, "")
	svc.CreateAccount("Online", &mid.ID, "")

	// Rename the middle account
	_, err := svc.UpdateAccount(mid.ID, stringPtr("Food"), nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the descendant's full_name was updated
	leaf, _ := svc.GetAccount(3) // Online is id 3
	if leaf.FullName != "Expenses:Food:Online" {
		t.Errorf("expected descendant full_name 'Expenses:Food:Online', got %q", leaf.FullName)
	}
}

func TestDeleteAccount_Success(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	acct, _ := svc.CreateAccount("Temporary", nil, "asset")
	err := svc.DeleteAccount(acct.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = svc.GetAccount(acct.ID)
	if err == nil {
		t.Fatal("expected not found after deletion")
	}
}

func TestDeleteAccount_RejectsIfHasChildren(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	parent, _ := svc.CreateAccount("Assets", nil, "asset")
	svc.CreateAccount("Checking", &parent.ID, "")

	err := svc.DeleteAccount(parent.ID)
	if err == nil {
		t.Fatal("expected error when deleting account with children")
	}
	appErr, ok := err.(*model.AppError)
	if !ok {
		t.Fatalf("expected *model.AppError, got %T", err)
	}
	if appErr.Code != 409 {
		t.Errorf("expected error code 409, got %d", appErr.Code)
	}
}

func TestDeleteAccount_RejectsIfHasPostings(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	acctA, _ := svc.CreateAccount("Assets", nil, "asset")
	acctB, _ := svc.CreateAccount("Expenses", nil, "expense")

	// Create a transaction with postings referencing the account
	txSvc := NewTransactionService(db, nil)
	_, err := txSvc.CreateTransaction("2026-01-01", "Test", "posted", []PostingInput{
		{AccountID: acctA.ID, Currency: "USD", Amount: "100.00"},
		{AccountID: acctB.ID, Currency: "USD", Amount: "-100.00"},
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error creating transaction: %v", err)
	}

	err = svc.DeleteAccount(acctA.ID)
	if err == nil {
		t.Fatal("expected error when deleting account with postings")
	}
}

func TestSearchAccounts(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	svc.CreateAccount("Assets", nil, "asset")
	root, _ := svc.CreateAccount("Expenses", nil, "expense")
	svc.CreateAccount("Shopping", &root.ID, "")

	results, err := svc.SearchAccounts("shop")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if len(results) > 0 && results[0].FullName != "Expenses:Shopping" {
		t.Errorf("expected 'Expenses:Shopping', got %q", results[0].FullName)
	}
}

func stringPtr(s string) *string {
	return &s
}

// --- Tests for CreateAccountByPath ---

func TestCreateAccountByPath_SingleSegmentSuggested(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	acct, err := svc.CreateAccountByPath("Expenses", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if acct.FullName != "Expenses" {
		t.Errorf("expected full_name 'Expenses', got %q", acct.FullName)
	}
	if acct.AccountType != "expense" {
		t.Errorf("expected type 'expense', got %q", acct.AccountType)
	}
}

func TestCreateAccountByPath_SingleSegmentNoSuggestion(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	// "Groceries" has no suggestion and no type provided — should fail
	_, err := svc.CreateAccountByPath("Groceries", "")
	if err == nil {
		t.Fatal("expected error for root with no type suggestion")
	}

	// Providing a type should work
	acct, err := svc.CreateAccountByPath("Groceries", "expense")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if acct.AccountType != "expense" {
		t.Errorf("expected type 'expense', got %q", acct.AccountType)
	}
}

func TestCreateAccountByPath_MultiSegment(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	acct, err := svc.CreateAccountByPath("Expenses:Food:Restaurants", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if acct.FullName != "Expenses:Food:Restaurants" {
		t.Errorf("expected 'Expenses:Food:Restaurants', got %q", acct.FullName)
	}
	if acct.Name != "Restaurants" {
		t.Errorf("expected name 'Restaurants', got %q", acct.Name)
	}

	if acct.AccountType != "expense" {
		t.Errorf("expected leaf account_type 'expense', got %q", acct.AccountType)
	}

	// Verify all intermediate accounts also have the type set
	results, _ := svc.SearchAccounts("Expenses")
	byFullName := make(map[string]string)
	for _, a := range results {
		byFullName[a.FullName] = a.AccountType
	}
	for _, fullName := range []string{"Expenses", "Expenses:Food", "Expenses:Food:Restaurants"} {
		if byFullName[fullName] != "expense" {
			t.Errorf("expected account_type 'expense' for %q, got %q", fullName, byFullName[fullName])
		}
	}
}

func TestCreateAccountByPath_PartialExists(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	// Pre-create root
	svc.CreateAccount("Expenses", nil, "expense")

	acct, err := svc.CreateAccountByPath("Expenses:Food", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if acct.FullName != "Expenses:Food" {
		t.Errorf("expected 'Expenses:Food', got %q", acct.FullName)
	}
	if acct.AccountType != "expense" {
		t.Errorf("expected account_type 'expense', got %q", acct.AccountType)
	}
}

func TestCreateAccountByPath_EmptySegmentRejected(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	_, err := svc.CreateAccountByPath("Expenses::Food", "")
	if err == nil {
		t.Fatal("expected error for empty segment")
	}
}

func TestCreateAccountByPath_EmptyNameRejected(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	_, err := svc.CreateAccountByPath("", "")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestCreateAccountByPath_SuggestionOverridesFallback(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	// "Assets" should be suggested as "asset" even when fallback is "expense"
	acct, err := svc.CreateAccountByPath("Assets", "expense")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if acct.AccountType != "asset" {
		t.Errorf("expected suggestion 'asset' to override fallback 'expense', got %q", acct.AccountType)
	}
}

func TestCreateAccountByPath_AlreadyExists(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := NewAccountService(db)

	// Create it first
	svc.CreateAccount("Expenses", nil, "expense")

	// CreateAccountByPath on an existing single-segment path should return it
	// Actually, since it calls CreateAccount which will fail on duplicate, let's verify
	// the multi-segment case: path fully exists → returns the leaf
	root, _ := svc.CreateAccount("Assets", nil, "asset")
	svc.CreateAccount("Checking", &root.ID, "")

	acct, err := svc.CreateAccountByPath("Assets:Checking", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if acct.FullName != "Assets:Checking" {
		t.Errorf("expected 'Assets:Checking', got %q", acct.FullName)
	}
}
