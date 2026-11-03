package service

import (
    "testing"

    "github.com/crispy-2/backend/internal/testutil"
    "github.com/shopspring/decimal"
)

func TestCreateBudget_Success(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    budgetSvc := NewBudgetService(db, NewExchangeRateService(db))

    expenses, _ := acctSvc.CreateAccount("Expenses", nil, "expense")

    budget, err := budgetSvc.CreateBudget("Monthly Food", "spending", "monthly", "USD", "500.00", []int64{expenses.ID}, nil)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if budget.Name != "Monthly Food" {
        t.Errorf("expected name 'Monthly Food', got %q", budget.Name)
    }
    if budget.Direction != "spending" {
        t.Errorf("expected direction 'spending', got %q", budget.Direction)
    }
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("500.00"), budget.TargetAmount.Decimal)
    if len(budget.AccountIDs) != 1 {
        t.Errorf("expected 1 account, got %d", len(budget.AccountIDs))
    }
}

func TestCreateBudget_NoAccountsRejected(t *testing.T) {
    db := testutil.NewTestDB(t)
    budgetSvc := NewBudgetService(db, NewExchangeRateService(db))

    _, err := budgetSvc.CreateBudget("Bad", "spending", "monthly", "USD", "500.00", nil, nil)
    if err == nil {
        t.Fatal("expected error for no accounts")
    }
}

func TestCreateBudget_InvalidDirectionRejected(t *testing.T) {
    db := testutil.NewTestDB(t)
    budgetSvc := NewBudgetService(db, NewExchangeRateService(db))

    _, err := budgetSvc.CreateBudget("Bad", "bogus", "monthly", "USD", "500.00", []int64{1}, nil)
    if err == nil {
        t.Fatal("expected error for invalid direction")
    }
}

func TestCreateBudget_InvalidPeriodRejected(t *testing.T) {
    db := testutil.NewTestDB(t)
    budgetSvc := NewBudgetService(db, NewExchangeRateService(db))

    _, err := budgetSvc.CreateBudget("Bad", "spending", "biweekly", "USD", "500.00", []int64{1}, nil)
    if err == nil {
        t.Fatal("expected error for invalid period")
    }
}

func TestGetBudget(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    budgetSvc := NewBudgetService(db, NewExchangeRateService(db))

    expenses, _ := acctSvc.CreateAccount("Expenses", nil, "expense")

    created, _ := budgetSvc.CreateBudget("Monthly Food", "spending", "monthly", "USD", "500.00", []int64{expenses.ID}, nil)
    got, err := budgetSvc.GetBudget(created.ID)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if got.Name != "Monthly Food" {
        t.Errorf("expected 'Monthly Food', got %q", got.Name)
    }
}

func TestGetBudget_NotFound(t *testing.T) {
    db := testutil.NewTestDB(t)
    budgetSvc := NewBudgetService(db, NewExchangeRateService(db))

    _, err := budgetSvc.GetBudget(9999)
    if err == nil {
        t.Fatal("expected not found error")
    }
}

func TestListBudgets(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    budgetSvc := NewBudgetService(db, NewExchangeRateService(db))

    expenses, _ := acctSvc.CreateAccount("Expenses", nil, "expense")
    budgetSvc.CreateBudget("Food", "spending", "monthly", "USD", "500.00", []int64{expenses.ID}, nil)
    budgetSvc.CreateBudget("Savings", "saving", "monthly", "USD", "1000.00", []int64{expenses.ID}, nil)

    budgets, err := budgetSvc.ListBudgets()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(budgets) != 2 {
        t.Errorf("expected 2 budgets, got %d", len(budgets))
    }
}

func TestUpdateBudget(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    budgetSvc := NewBudgetService(db, NewExchangeRateService(db))

    expenses, _ := acctSvc.CreateAccount("Expenses", nil, "expense")
    created, _ := budgetSvc.CreateBudget("Food", "spending", "monthly", "USD", "500.00", []int64{expenses.ID}, nil)

    updated, err := budgetSvc.UpdateBudget(created.ID, "Food & Drinks", "spending", "monthly", "USD", "600.00", []int64{expenses.ID}, nil)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if updated.Name != "Food & Drinks" {
        t.Errorf("expected 'Food & Drinks', got %q", updated.Name)
    }
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("600.00"), updated.TargetAmount.Decimal)
}

func TestDeleteBudget(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    budgetSvc := NewBudgetService(db, NewExchangeRateService(db))

    expenses, _ := acctSvc.CreateAccount("Expenses", nil, "expense")
    created, _ := budgetSvc.CreateBudget("Food", "spending", "monthly", "USD", "500.00", []int64{expenses.ID}, nil)

    err := budgetSvc.DeleteBudget(created.ID)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    _, err = budgetSvc.GetBudget(created.ID)
    if err == nil {
        t.Fatal("expected not found after deletion")
    }
}

func TestBudgetComputation(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    txSvc := NewTransactionService(db, nil)
    budgetSvc := NewBudgetService(db, NewExchangeRateService(db))

    assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
    checking, _ := acctSvc.CreateAccount("Checking", &assets.ID, "")
    expenses, _ := acctSvc.CreateAccount("Expenses", nil, "expense")
    food, _ := acctSvc.CreateAccount("Food", &expenses.ID, "")

    // Create a budget for Expenses (includes children)
    budget, _ := budgetSvc.CreateBudget("Monthly Expenses", "spending", "monthly", "USD", "1000.00", []int64{expenses.ID}, nil)

    // Create transactions within the period
    txSvc.CreateTransaction("2026-10-05", "Groceries", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-50.00"},
        {AccountID: food.ID, Currency: "USD", Amount: "50.00"},
    }, nil)
    txSvc.CreateTransaction("2026-10-10", "Restaurant", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-30.00"},
        {AccountID: food.ID, Currency: "USD", Amount: "30.00"},
    }, nil)

    // Compute the budget for October 2026
    current, err := budgetSvc.ComputeBudgetCurrent(budget.ID, "2026-10-01", "2026-10-31")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    // Should be 80.00 (50 + 30 in the Food sub-account of Expenses)
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("80"), current.Decimal)
}

func TestBudgetComputation_WithTagFilter(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    txSvc := NewTransactionService(db, nil)
    budgetSvc := NewBudgetService(db, NewExchangeRateService(db))

    assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
    checking, _ := acctSvc.CreateAccount("Checking", &assets.ID, "")
    expenses, _ := acctSvc.CreateAccount("Expenses", nil, "expense")
    food, _ := acctSvc.CreateAccount("Food", &expenses.ID, "")
    transport, _ := acctSvc.CreateAccount("Transport", &expenses.ID, "")

    // Create transactions with tags
    tx1, _ := txSvc.CreateTransaction("2026-10-05", "Groceries", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-50.00"},
        {AccountID: food.ID, Currency: "USD", Amount: "50.00"},
    }, []TagInput{{Key: "trip", Value: "japan-2026"}})

    txSvc.CreateTransaction("2026-10-10", "Restaurant", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-30.00"},
        {AccountID: food.ID, Currency: "USD", Amount: "30.00"},
    }, nil)

    txSvc.CreateTransaction("2026-10-15", "Bus Ticket", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-10.00"},
        {AccountID: transport.ID, Currency: "USD", Amount: "10.00"},
    }, []TagInput{{Key: "trip", Value: "japan-2026"}})

    // Create budget with tag filter
    budget, _ := budgetSvc.CreateBudget("Food Budget", "spending", "monthly", "USD", "100.00", []int64{food.ID}, []int64{tx1.Tags[0].ID})

    // Compute the budget for October 2026
    current, err := budgetSvc.ComputeBudgetCurrent(budget.ID, "2026-10-01", "2026-10-31")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    // Should be 50.00 (only the transaction with the matching tag)
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("50"), current.Decimal)
}

func TestBudgetComputation_PeriodBoundary(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    txSvc := NewTransactionService(db, nil)
    budgetSvc := NewBudgetService(db, NewExchangeRateService(db))

    assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
    checking, _ := acctSvc.CreateAccount("Checking", &assets.ID, "")
    expenses, _ := acctSvc.CreateAccount("Expenses", nil, "expense")
    food, _ := acctSvc.CreateAccount("Food", &expenses.ID, "")

    // Create a budget for October 2026
    budget, _ := budgetSvc.CreateBudget("October Budget", "spending", "monthly", "USD", "100.00", []int64{food.ID}, nil)

    // Create transaction on the last day of October (should be included)
    txSvc.CreateTransaction("2026-10-31", "Groceries Oct", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-20.00"},
        {AccountID: food.ID, Currency: "USD", Amount: "20.00"},
    }, nil)

    // Create transaction on the first day of November (should be excluded)
    txSvc.CreateTransaction("2026-11-01", "Groceries Nov", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-30.00"},
        {AccountID: food.ID, Currency: "USD", Amount: "30.00"},
    }, nil)

    // Create transaction on the last day of September (should be excluded)
    txSvc.CreateTransaction("2026-09-30", "Groceries Sep", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-10.00"},
        {AccountID: food.ID, Currency: "USD", Amount: "10.00"},
    }, nil)

    // Compute the budget for October 2026
    current, err := budgetSvc.ComputeBudgetCurrent(budget.ID, "2026-10-01", "2026-10-31")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    // Should be 20.00 (only the Oct 31 transaction)
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("20"), current.Decimal)
}

func TestBudgetComputation_AccountChildren(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    txSvc := NewTransactionService(db, nil)
    budgetSvc := NewBudgetService(db, NewExchangeRateService(db))

    assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
    checking, _ := acctSvc.CreateAccount("Checking", &assets.ID, "")
    expenses, _ := acctSvc.CreateAccount("Expenses", nil, "expense")
    food, _ := acctSvc.CreateAccount("Food", &expenses.ID, "")
    shopping, _ := acctSvc.CreateAccount("Shopping", &expenses.ID, "")
    online, _ := acctSvc.CreateAccount("Online", &shopping.ID, "")

    // Create a budget targeting the parent Expenses account
    budget, _ := budgetSvc.CreateBudget("All Expenses", "spending", "monthly", "USD", "500.00", []int64{expenses.ID}, nil)

    // Create transaction against Food (child of Expenses)
    txSvc.CreateTransaction("2026-10-05", "Groceries", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-50.00"},
        {AccountID: food.ID, Currency: "USD", Amount: "50.00"},
    }, nil)

    // Create transaction against Shopping (child of Expenses)
    txSvc.CreateTransaction("2026-10-10", "Book", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-30.00"},
        {AccountID: shopping.ID, Currency: "USD", Amount: "30.00"},
    }, nil)

    // Create transaction against Online (grandchild of Expenses)
    txSvc.CreateTransaction("2026-10-15", "Amazon", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-20.00"},
        {AccountID: online.ID, Currency: "USD", Amount: "20.00"},
    }, nil)

    // Compute the budget for October 2026
    current, err := budgetSvc.ComputeBudgetCurrent(budget.ID, "2026-10-01", "2026-10-31")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    // Should be 100.00 (50 + 30 + 20 from all child accounts)
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("100"), current.Decimal)
}
