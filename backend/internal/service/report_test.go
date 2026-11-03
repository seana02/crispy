package service

import (
    "testing"

    "github.com/crispy-2/backend/internal/testutil"
    "github.com/shopspring/decimal"
)

func TestGetBalanceSheet(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    txSvc := NewTransactionService(db, nil)
    reportSvc := NewReportService(db, NewExchangeRateService(db))

    // Create accounts
    assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
    checking, _ := acctSvc.CreateAccount("Checking", &assets.ID, "")
    liabilities, _ := acctSvc.CreateAccount("Liabilities", nil, "liability")
    creditCard, _ := acctSvc.CreateAccount("CreditCard", &liabilities.ID, "")

    // Create transactions
    txSvc.CreateTransaction("2026-01-15", "Deposit", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "5000.00"},
        {AccountID: creditCard.ID, Currency: "USD", Amount: "-5000.00"},
    }, nil)

    report, err := reportSvc.GetBalanceSheet("2026-01-01", "2026-01-31", 3, nil)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(report.AssetRows) == 0 {
        t.Error("expected asset rows in balance sheet")
    }
    if len(report.LiabilityRows) == 0 {
        t.Error("expected liability rows in balance sheet")
    }
}

func TestGetIncomeStatement(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    txSvc := NewTransactionService(db, nil)
    reportSvc := NewReportService(db, NewExchangeRateService(db))

    // Create accounts
    assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
    checking, _ := acctSvc.CreateAccount("Checking", &assets.ID, "")
    revenue, _ := acctSvc.CreateAccount("Revenue", nil, "revenue")
    salary, _ := acctSvc.CreateAccount("Salary", &revenue.ID, "")
    expenses, _ := acctSvc.CreateAccount("Expenses", nil, "expense")
    food, _ := acctSvc.CreateAccount("Food", &expenses.ID, "")

    // Income
    txSvc.CreateTransaction("2026-01-15", "Paycheck", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "3000.00"},
        {AccountID: salary.ID, Currency: "USD", Amount: "-3000.00"},
    }, nil)
    // Expense
    txSvc.CreateTransaction("2026-01-20", "Groceries", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-200.00"},
        {AccountID: food.ID, Currency: "USD", Amount: "200.00"},
    }, nil)

    report, err := reportSvc.GetIncomeStatement("2026-01-01", "2026-01-31", 3, nil)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(report.RevenueRows) == 0 {
        t.Error("expected revenue rows")
    }
    if len(report.ExpenseRows) == 0 {
        t.Error("expected expense rows")
    }
}

func TestGetNetWorth(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    txSvc := NewTransactionService(db, nil)
    reportSvc := NewReportService(db, NewExchangeRateService(db))

    assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
    checking, _ := acctSvc.CreateAccount("Checking", &assets.ID, "")
    liabilities, _ := acctSvc.CreateAccount("Liabilities", nil, "liability")
    creditCard, _ := acctSvc.CreateAccount("CreditCard", &liabilities.ID, "")
    equity, _ := acctSvc.CreateAccount("Equity", nil, "equity")

    // Deposit 5000
    txSvc.CreateTransaction("2026-01-15", "Deposit", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "5000.00"},
        {AccountID: equity.ID, Currency: "USD", Amount: "-5000.00"},
    }, nil)
    // Charge 1000 to credit card
    txSvc.CreateTransaction("2026-01-20", "Charge", "posted", []PostingInput{
        {AccountID: creditCard.ID, Currency: "USD", Amount: "1000.00"},
        {AccountID: checking.ID, Currency: "USD", Amount: "-1000.00"},
    }, nil)

    netWorth, err := reportSvc.GetNetWorth("2026-01-31")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    // Net worth = Assets(4000) - Liabilities(1000) = 3000
    // Assets: Checking = 5000 - 1000 = 4000
    // Liabilities: CreditCard = 1000
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("3000"), netWorth.NetWorth.Decimal)
}

func TestGetTopSpendingCategories(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    txSvc := NewTransactionService(db, nil)
    reportSvc := NewReportService(db, NewExchangeRateService(db))

    assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
    checking, _ := acctSvc.CreateAccount("Checking", &assets.ID, "")
    expenses, _ := acctSvc.CreateAccount("Expenses", nil, "expense")
    food, _ := acctSvc.CreateAccount("Food", &expenses.ID, "")
    rent, _ := acctSvc.CreateAccount("Rent", &expenses.ID, "")

    txSvc.CreateTransaction("2026-01-10", "Groceries", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-200.00"},
        {AccountID: food.ID, Currency: "USD", Amount: "200.00"},
    }, nil)
    txSvc.CreateTransaction("2026-01-15", "Rent", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-1500.00"},
        {AccountID: rent.ID, Currency: "USD", Amount: "1500.00"},
    }, nil)

    categories, err := reportSvc.GetTopSpendingCategories("2026-01-01", "2026-01-31", 5)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(categories) < 2 {
        t.Fatalf("expected at least 2 categories, got %d", len(categories))
    }
    // Rent should be first (highest)
    if categories[0].AccountName != "Expenses:Rent" {
        t.Errorf("expected top category 'Expenses:Rent', got %q", categories[0].AccountName)
    }
}

func TestSearchGlobal(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    txSvc := NewTransactionService(db, nil)
    searchSvc := NewSearchService(db)

    assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
    checking, _ := acctSvc.CreateAccount("Checking", &assets.ID, "")
    expenses, _ := acctSvc.CreateAccount("Expenses", nil, "expense")

    txSvc.CreateTransaction("2026-01-15", "Grocery shopping", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-50.00"},
        {AccountID: expenses.ID, Currency: "USD", Amount: "50.00"},
    }, []TagInput{{Key: "category", Value: "groceries"}})

    results, err := searchSvc.GlobalSearch("grocer")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    // Should find the transaction (description) and the tag (value)
    if len(results.Transactions) == 0 {
        t.Error("expected to find transaction matching 'grocer'")
    }
    if len(results.Tags) == 0 {
        t.Error("expected to find tag matching 'grocer'")
    }
}

func TestGetExpenseBreakdownByAccount(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    txSvc := NewTransactionService(db, nil)
    reportSvc := NewReportService(db, NewExchangeRateService(db))

    assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
    checking, _ := acctSvc.CreateAccount("Checking", &assets.ID, "")
    expenses, _ := acctSvc.CreateAccount("Expenses", nil, "expense")
    food, _ := acctSvc.CreateAccount("Food", &expenses.ID, "")
    rent, _ := acctSvc.CreateAccount("Rent", &expenses.ID, "")

    txSvc.CreateTransaction("2026-01-10", "Groceries", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-200.00"},
        {AccountID: food.ID, Currency: "USD", Amount: "200.00"},
    }, nil)
    txSvc.CreateTransaction("2026-01-15", "Rent Payment", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-1500.00"},
        {AccountID: rent.ID, Currency: "USD", Amount: "1500.00"},
    }, nil)

    breakdown, err := reportSvc.GetExpenseBreakdown("2026-01-01", "2026-01-31", "account")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if len(breakdown) != 2 {
        t.Fatalf("expected 2 expense accounts, got %d", len(breakdown))
    }

    total := decimal.Zero
    for _, amt := range breakdown {
        total = total.Add(amt.Decimal)
    }
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("1700"), total)
}

func TestGetExpenseBreakdownByTag(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    txSvc := NewTransactionService(db, nil)
    reportSvc := NewReportService(db, NewExchangeRateService(db))

    assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
    checking, _ := acctSvc.CreateAccount("Checking", &assets.ID, "")
    expenses, _ := acctSvc.CreateAccount("Expenses", nil, "expense")
    food, _ := acctSvc.CreateAccount("Food", &expenses.ID, "")

    txSvc.CreateTransaction("2026-01-10", "Groceries", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-200.00"},
        {AccountID: food.ID, Currency: "USD", Amount: "200.00"},
    }, []TagInput{{Key: "trip", Value: "japan-2026"}})
    txSvc.CreateTransaction("2026-01-15", "Local Food", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-50.00"},
        {AccountID: food.ID, Currency: "USD", Amount: "50.00"},
    }, []TagInput{{Key: "trip", Value: "home"}})

    breakdown, err := reportSvc.GetExpenseBreakdown("2026-01-01", "2026-01-31", "tag")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if len(breakdown) != 2 {
        t.Fatalf("expected 2 tag values, got %d", len(breakdown))
    }
}

func TestGetTagSpendingSummary(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    txSvc := NewTransactionService(db, nil)
    reportSvc := NewReportService(db, NewExchangeRateService(db))

    assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
    checking, _ := acctSvc.CreateAccount("Checking", &assets.ID, "")
    expenses, _ := acctSvc.CreateAccount("Expenses", nil, "expense")
    food, _ := acctSvc.CreateAccount("Food", &expenses.ID, "")
    rent, _ := acctSvc.CreateAccount("Rent", &expenses.ID, "")

    txSvc.CreateTransaction("2026-01-10", "Groceries Japan", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-200.00"},
        {AccountID: food.ID, Currency: "USD", Amount: "200.00"},
    }, []TagInput{{Key: "trip", Value: "japan-2026"}})
    txSvc.CreateTransaction("2026-01-15", "Groceries Home", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-50.00"},
        {AccountID: food.ID, Currency: "USD", Amount: "50.00"},
    }, []TagInput{{Key: "trip", Value: "home"}})
    txSvc.CreateTransaction("2026-01-20", "Rent", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-1500.00"},
        {AccountID: rent.ID, Currency: "USD", Amount: "1500.00"},
    }, nil)

    summary, err := reportSvc.GetTagSpendingSummary("trip", "2026-01-01", "2026-01-31")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if len(summary) != 2 {
        t.Fatalf("expected 2 tag values, got %d", len(summary))
    }

    japanAmt := summary["japan-2026 (USD)"]
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("200"), japanAmt.Decimal)

    homeAmt := summary["home (USD)"]
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("50"), homeAmt.Decimal)
}

func TestGetNetWorthWithComparison(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    txSvc := NewTransactionService(db, nil)
    reportSvc := NewReportService(db, NewExchangeRateService(db))

    assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
    checking, _ := acctSvc.CreateAccount("Checking", &assets.ID, "")
    liabilities, _ := acctSvc.CreateAccount("Liabilities", nil, "liability")
    creditCard, _ := acctSvc.CreateAccount("CreditCard", &liabilities.ID, "")
    equity, _ := acctSvc.CreateAccount("Equity", nil, "equity")

    txSvc.CreateTransaction("2026-01-15", "Initial Deposit", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "5000.00"},
        {AccountID: equity.ID, Currency: "USD", Amount: "-5000.00"},
    }, nil)

    netWorth1, err := reportSvc.GetNetWorth("2026-01-31")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("5000"), netWorth1.NetWorth.Decimal)

    txSvc.CreateTransaction("2026-02-15", "Charge Credit Card", "posted", []PostingInput{
        {AccountID: creditCard.ID, Currency: "USD", Amount: "1000.00"},
        {AccountID: checking.ID, Currency: "USD", Amount: "-1000.00"},
    }, nil)

    netWorth2, err := reportSvc.GetNetWorth("2026-02-28", "2026-01-31")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    // Feb 28: Assets=4000, Liabilities=1000, NetWorth=3000
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("3000"), netWorth2.NetWorth.Decimal)
    // Trend: 3000 - 5000 = -2000
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("-2000"), netWorth2.Trend.Decimal)
}

func TestGetBalanceSheetWithTags(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    txSvc := NewTransactionService(db, nil)
    reportSvc := NewReportService(db, NewExchangeRateService(db))

    assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
    checking, _ := acctSvc.CreateAccount("Checking", &assets.ID, "")
    liabilities, _ := acctSvc.CreateAccount("Liabilities", nil, "liability")
    creditCard, _ := acctSvc.CreateAccount("CreditCard", &liabilities.ID, "")

    txSvc.CreateTransaction("2026-01-15", "Deposit for Japan Trip", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "5000.00"},
        {AccountID: creditCard.ID, Currency: "USD", Amount: "-5000.00"},
    }, []TagInput{{Key: "trip", Value: "japan-2026"}})

    txSvc.CreateTransaction("2026-01-20", "General Deposit", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "1000.00"},
        {AccountID: creditCard.ID, Currency: "USD", Amount: "-1000.00"},
    }, nil)

    report, err := reportSvc.GetBalanceSheet("2026-01-01", "2026-01-31", 3, []int64{})
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(report.AssetRows) == 0 {
        t.Error("expected asset rows in balance sheet")
    }

    var tagID int64
    tagName := "trip:japan-2026"
    row := db.QueryRow(`SELECT id FROM tags WHERE key || ':' || value = ?`, tagName)
    err = row.Scan(&tagID)
    if err != nil {
        t.Fatalf("failed to get tag ID: %v", err)
    }

    reportWithFilter, err := reportSvc.GetBalanceSheet("2026-01-01", "2026-01-31", 3, []int64{tagID})
    if err != nil {
        t.Fatalf("unexpected error with tag filter: %v", err)
    }
    if len(reportWithFilter.AssetRows) == 0 {
        t.Error("expected asset rows with tag filter")
    }
}

func TestGetIncomeStatementWithTags(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    txSvc := NewTransactionService(db, nil)
    reportSvc := NewReportService(db, NewExchangeRateService(db))

    assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
    checking, _ := acctSvc.CreateAccount("Checking", &assets.ID, "")
    revenue, _ := acctSvc.CreateAccount("Revenue", nil, "revenue")
    salary, _ := acctSvc.CreateAccount("Salary", &revenue.ID, "")
    expenses, _ := acctSvc.CreateAccount("Expenses", nil, "expense")
    food, _ := acctSvc.CreateAccount("Food", &expenses.ID, "")

    txSvc.CreateTransaction("2026-01-15", "Paycheck Japan Trip", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "3000.00"},
        {AccountID: salary.ID, Currency: "USD", Amount: "-3000.00"},
    }, []TagInput{{Key: "trip", Value: "japan-2026"}})
    txSvc.CreateTransaction("2026-01-20", "Groceries", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-200.00"},
        {AccountID: food.ID, Currency: "USD", Amount: "200.00"},
    }, nil)

    report, err := reportSvc.GetIncomeStatement("2026-01-01", "2026-01-31", 3, []int64{})
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(report.RevenueRows) == 0 {
        t.Error("expected revenue rows")
    }

    var tagID int64
    row := db.QueryRow(`SELECT id FROM tags WHERE key || ':' || value = ?`, "trip:japan-2026")
    err = row.Scan(&tagID)
    if err != nil {
        t.Fatalf("failed to get tag ID: %v", err)
    }

    reportWithFilter, err := reportSvc.GetIncomeStatement("2026-01-01", "2026-01-31", 3, []int64{tagID})
    if err != nil {
        t.Fatalf("unexpected error with tag filter: %v", err)
    }
    if len(reportWithFilter.RevenueRows) == 0 {
        t.Error("expected revenue rows with tag filter")
    }
}

func TestExchangeRateFetchNotConfigured(t *testing.T) {
    db := testutil.NewTestDB(t)
    svc := NewExchangeRateService(db)

    err := svc.FetchRatesFromAPI("USD", []string{"JPY", "EUR"})
    if err == nil {
        t.Error("expected error for unconfigured API")
    }
    if err.Error() != "exchange rate API provider not configured" {
        t.Errorf("unexpected error message: %v", err)
    }
}

