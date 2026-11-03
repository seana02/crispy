package service

import (
    "testing"

    "github.com/crispy-2/backend/internal/model"
    "github.com/crispy-2/backend/internal/testutil"
    "github.com/shopspring/decimal"
)

// Integration tests cover end-to-end workflows through the backend (service + DAL + database).

func TestIntegration_CreateAccountsAndTransaction(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    txSvc := NewTransactionService(db, nil)

    // Create account hierarchy
    assets, err := acctSvc.CreateAccount("Assets", nil, "asset")
    if err != nil {
        t.Fatalf("failed to create Assets: %v", err)
    }
    checking, err := acctSvc.CreateAccount("Checking", &assets.ID, "")
    if err != nil {
        t.Fatalf("failed to create Checking: %v", err)
    }
    expenses, err := acctSvc.CreateAccount("Expenses", nil, "expense")
    if err != nil {
        t.Fatalf("failed to create Expenses: %v", err)
    }
    food, err := acctSvc.CreateAccount("Food", &expenses.ID, "")
    if err != nil {
        t.Fatalf("failed to create Food: %v", err)
    }

    // Verify full name hierarchy
    if checking.FullName != "Assets:Checking" {
        t.Errorf("expected 'Assets:Checking', got %q", checking.FullName)
    }
    if food.FullName != "Expenses:Food" {
        t.Errorf("expected 'Expenses:Food', got %q", food.FullName)
    }

    // Create a transaction
    txn, err := txSvc.CreateTransaction("2026-01-15", "Groceries", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-50.00"},
        {AccountID: food.ID, Currency: "USD", Amount: "50.00"},
    }, []TagInput{
        {Key: "trip", Value: "japan-2026"},
    })
    if err != nil {
        t.Fatalf("failed to create transaction: %v", err)
    }

    // Verify full transaction retrieval
    got, err := txSvc.GetTransaction(txn.ID)
    if err != nil {
        t.Fatalf("failed to get transaction: %v", err)
    }
    if len(got.Postings) != 2 {
        t.Fatalf("expected 2 postings, got %d", len(got.Postings))
    }
    if len(got.Tags) != 1 {
        t.Fatalf("expected 1 tag, got %d", len(got.Tags))
    }
    if got.Tags[0].Key != "trip" || got.Tags[0].Value != "japan-2026" {
        t.Errorf("expected tag trip:japan-2026, got %s:%s", got.Tags[0].Key, got.Tags[0].Value)
    }

    // Verify transaction appears in list filtered by account
    result, err := txSvc.ListTransactions(TransactionFilters{AccountID: &food.ID}, model.Pagination{Page: 1, PageSize: 10})
    if err != nil {
        t.Fatalf("failed to list by account: %v", err)
    }
    if result.Total != 1 {
        t.Errorf("expected 1 transaction for food account, got %d", result.Total)
    }

    // Verify transaction appears when filtering by tag
    result, err = txSvc.ListTransactions(TransactionFilters{TagKey: "trip", TagValue: "japan-2026"}, model.Pagination{Page: 1, PageSize: 10})
    if err != nil {
        t.Fatalf("failed to list by tag: %v", err)
    }
    if result.Total != 1 {
        t.Errorf("expected 1 transaction for tag, got %d", result.Total)
    }
}

func TestIntegration_MultiCurrencyTransactionWithExchangeRates(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    erSvc := NewExchangeRateService(db)
    txSvc := NewTransactionService(db, erSvc)

    assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
    usdAccount, _ := acctSvc.CreateAccount("USD Bank", &assets.ID, "")
    jpyAccount, _ := acctSvc.CreateAccount("JPY Bank", &assets.ID, "")

    // Create exchange rate
    _, err := erSvc.CreateRate("USD", "JPY", "150.00", "2026-01-01", "manual")
    if err != nil {
        t.Fatalf("failed to create rate: %v", err)
    }

    // Create a multi-currency transaction (all amounts explicit)
    _, err = txSvc.CreateTransaction("2026-01-15", "Currency exchange", "posted", []PostingInput{
        {AccountID: usdAccount.ID, Currency: "USD", Amount: "-100.00"},
        {AccountID: jpyAccount.ID, Currency: "JPY", Amount: "15000"},
    }, nil)
    if err != nil {
        t.Fatalf("failed to create multi-currency transaction: %v", err)
    }

    // Verify the rate lookup works
    rate, err := erSvc.GetRate("USD", "JPY", "2026-01-15")
    if err != nil {
        t.Fatalf("failed to get rate: %v", err)
    }
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("150.00"), rate.Rate.Decimal)

    // Verify inverse rate works
    inverseRate, err := erSvc.GetRate("JPY", "USD", "2026-01-15")
    if err != nil {
        t.Fatalf("failed to get inverse rate: %v", err)
    }
    expected := decimal.NewFromInt(1).Div(decimal.RequireFromString("150.00"))
    testutil.AssertDecimalEqual(t, expected, inverseRate.Rate.Decimal)
}

func TestIntegration_BudgetLifecycle(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    txSvc := NewTransactionService(db, nil)
    budgetSvc := NewBudgetService(db, NewExchangeRateService(db))

    // Set up accounts
    assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
    checking, _ := acctSvc.CreateAccount("Checking", &assets.ID, "")
    expenses, _ := acctSvc.CreateAccount("Expenses", nil, "expense")
    food, _ := acctSvc.CreateAccount("Food", &expenses.ID, "")
    transport, _ := acctSvc.CreateAccount("Transport", &expenses.ID, "")

    // Create budget tracking Expenses (and all children)
    budget, err := budgetSvc.CreateBudget("Monthly Budget", "spending", "monthly", "USD", "500.00", []int64{expenses.ID}, nil)
    if err != nil {
        t.Fatalf("failed to create budget: %v", err)
    }

    // Post some transactions
    txSvc.CreateTransaction("2026-10-05", "Groceries", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-200.00"},
        {AccountID: food.ID, Currency: "USD", Amount: "200.00"},
    }, nil)
    txSvc.CreateTransaction("2026-10-10", "Bus", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-50.00"},
        {AccountID: transport.ID, Currency: "USD", Amount: "50.00"},
    }, nil)

    // Compute budget
    current, err := budgetSvc.ComputeBudgetCurrent(budget.ID, "2026-10-01", "2026-10-31")
    if err != nil {
        t.Fatalf("failed to compute budget: %v", err)
    }
    // 200 + 50 = 250 spent
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("250"), current.Decimal)

    // Verify budget is under target (500)
    if current.Decimal.GreaterThan(budget.TargetAmount.Decimal) {
        t.Error("budget should be under target")
    }
}

func TestIntegration_RecurringTransactionGeneration(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    txSvc := NewTransactionService(db, nil)
    recSvc := NewRecurringService(db, txSvc)

    assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
    checking, _ := acctSvc.CreateAccount("Checking", &assets.ID, "")
    expenses, _ := acctSvc.CreateAccount("Expenses", nil, "expense")

    // Create weekly recurring template
    _, err := recSvc.CreateTemplate("Weekly Groceries", "posted", "every_n_weeks", 1, nil, "2026-10-01", nil, []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-100.00"},
        {AccountID: expenses.ID, Currency: "USD", Amount: "100.00"},
    }, []TagInput{{Key: "category", Value: "groceries"}})
    if err != nil {
        t.Fatalf("failed to create template: %v", err)
    }

    // Simulate the app being offline until Oct 22 - should catch up
    count, err := recSvc.GenerateDueTransactions("2026-10-22")
    if err != nil {
        t.Fatalf("failed to generate: %v", err)
    }
    // Oct 1, 8, 15, 22 = 4 transactions
    if count != 4 {
        t.Errorf("expected 4 catch-up transactions, got %d", count)
    }

    // Verify transactions exist with correct tags
    result, _ := txSvc.ListTransactions(TransactionFilters{TagKey: "category", TagValue: "groceries"}, model.Pagination{Page: 1, PageSize: 20})
    if result.Total != 4 {
        t.Errorf("expected 4 tagged transactions, got %d", result.Total)
    }

    // Run again - should generate nothing new
    count2, _ := recSvc.GenerateDueTransactions("2026-10-22")
    if count2 != 0 {
        t.Errorf("expected 0 on second run, got %d", count2)
    }

    // Advance to Oct 29 - should generate 1 more
    count3, _ := recSvc.GenerateDueTransactions("2026-10-29")
    if count3 != 1 {
        t.Errorf("expected 1 new transaction for Oct 29, got %d", count3)
    }
}

func TestIntegration_ReportAccuracy(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    txSvc := NewTransactionService(db, nil)
    reportSvc := NewReportService(db, NewExchangeRateService(db))

    // Set up accounts
    assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
    checking, _ := acctSvc.CreateAccount("Checking", &assets.ID, "")
    liabilities, _ := acctSvc.CreateAccount("Liabilities", nil, "liability")
    creditCard, _ := acctSvc.CreateAccount("CreditCard", &liabilities.ID, "")
    revenue, _ := acctSvc.CreateAccount("Revenue", nil, "revenue")
    salary, _ := acctSvc.CreateAccount("Salary", &revenue.ID, "")
    expenses, _ := acctSvc.CreateAccount("Expenses", nil, "expense")
    food, _ := acctSvc.CreateAccount("Food", &expenses.ID, "")
    equity, _ := acctSvc.CreateAccount("Equity", nil, "equity")

    // January: Salary deposit
    txSvc.CreateTransaction("2026-01-15", "Paycheck", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "5000.00"},
        {AccountID: salary.ID, Currency: "USD", Amount: "-5000.00"},
    }, nil)

    // January: Groceries
    txSvc.CreateTransaction("2026-01-20", "Groceries", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-300.00"},
        {AccountID: food.ID, Currency: "USD", Amount: "300.00"},
    }, nil)

    // January: Credit card charge
    txSvc.CreateTransaction("2026-01-25", "Online purchase", "posted", []PostingInput{
        {AccountID: creditCard.ID, Currency: "USD", Amount: "200.00"},
        {AccountID: food.ID, Currency: "USD", Amount: "200.00"},
        {AccountID: checking.ID, Currency: "USD", Amount: "-400.00"},
    }, nil)

    _ = equity // unused but available for opening balance scenarios

    // Balance sheet: Assets = 5000 - 300 - 400 = 4300, Liabilities = 200
    bs, err := reportSvc.GetBalanceSheet("2026-01-01", "2026-01-31", 3, nil)
    if err != nil {
        t.Fatalf("failed to get balance sheet: %v", err)
    }
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("4300"), bs.TotalAssets.Decimal)
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("200"), bs.TotalLiab.Decimal)

    // Net worth = 4300 - 200 = 4100
    netWorth, err := reportSvc.GetNetWorth("2026-01-31")
    if err != nil {
        t.Fatalf("failed to get net worth: %v", err)
    }
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("4100"), netWorth.NetWorth.Decimal)

    // Income statement: Revenue = -5000, Expenses = 500
    is, err := reportSvc.GetIncomeStatement("2026-01-01", "2026-01-31", 3, nil)
    if err != nil {
        t.Fatalf("failed to get income statement: %v", err)
    }
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("-5000"), is.TotalRevenue.Decimal)
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("500"), is.TotalExpense.Decimal)
}

func TestIntegration_TagFilteringInReports(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    txSvc := NewTransactionService(db, nil)
    reportSvc := NewReportService(db, NewExchangeRateService(db))

    // Set up accounts
    assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
    checking, _ := acctSvc.CreateAccount("Checking", &assets.ID, "")
    expenses, _ := acctSvc.CreateAccount("Expenses", nil, "expense")
    food, _ := acctSvc.CreateAccount("Food", &expenses.ID, "")
    transport, _ := acctSvc.CreateAccount("Transport", &expenses.ID, "")

    // Create transactions with different tags
    // Japan trip transactions
    tx1, _ := txSvc.CreateTransaction("2026-10-05", "Groceries Japan", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-100.00"},
        {AccountID: food.ID, Currency: "USD", Amount: "100.00"},
    }, []TagInput{{Key: "trip", Value: "japan-2026"}})

    txSvc.CreateTransaction("2026-10-06", "Bus Ticket Japan", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-20.00"},
        {AccountID: transport.ID, Currency: "USD", Amount: "20.00"},
    }, []TagInput{{Key: "trip", Value: "japan-2026"}})

    // Home transactions (no tag filter)
    txSvc.CreateTransaction("2026-10-10", "Groceries Home", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-50.00"},
        {AccountID: food.ID, Currency: "USD", Amount: "50.00"},
    }, []TagInput{{Key: "trip", Value: "home"}})

    txSvc.CreateTransaction("2026-10-11", "Restaurant Home", "posted", []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-30.00"},
        {AccountID: food.ID, Currency: "USD", Amount: "30.00"},
    }, nil)

    // Get tag IDs
    tagJapan := tx1.Tags[0].ID

    // Balance sheet with tag filter (should only include Japan transactions)
    bs, err := reportSvc.GetBalanceSheet("2026-10-01", "2026-10-31", 3, []int64{tagJapan})
    if err != nil {
        t.Fatalf("failed to get balance sheet: %v", err)
    }
    // Checking only has Japan transactions: -100 + -20 = -120
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("-120"), bs.TotalAssets.Decimal)

    // Income statement with tag filter (should only include Japan transactions)
    is, err := reportSvc.GetIncomeStatement("2026-10-01", "2026-10-31", 3, []int64{tagJapan})
    if err != nil {
        t.Fatalf("failed to get income statement: %v", err)
    }
    // Revenue: 0, Expenses: 100 + 20 = 120 (from Japan transactions)
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("0"), is.TotalRevenue.Decimal)
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("120"), is.TotalExpense.Decimal)

    // Verify that without tag filter, we get all transactions
    bsAll, _ := reportSvc.GetBalanceSheet("2026-10-01", "2026-10-31", 3, nil)
    // Checking: -100 + -20 + -50 + -30 = -200
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("-200"), bsAll.TotalAssets.Decimal)

    isAll, _ := reportSvc.GetIncomeStatement("2026-10-01", "2026-10-31", 3, nil)
    // Expenses: 100 + 20 + 50 + 30 = 200
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("200"), isAll.TotalExpense.Decimal)
}

func TestIntegration_WidgetCRUD(t *testing.T) {
    db := testutil.NewTestDB(t)
    widgetSvc := NewWidgetService(db)

    // Create widgets
    widget1, err := widgetSvc.CreateWidget("account_balance", 0, 0, 2, 1, `{"account_id": 1}`)
    if err != nil {
        t.Fatalf("failed to create widget1: %v", err)
    }
    if widget1.WidgetType != "account_balance" {
        t.Errorf("expected type 'account_balance', got %q", widget1.WidgetType)
    }

    widget2, err := widgetSvc.CreateWidget("budget_progress", 2, 0, 2, 1, `{"budget_id": 1}`)
    if err != nil {
        t.Fatalf("failed to create widget2: %v", err)
    }

    // List widgets
    widgets, err := widgetSvc.ListWidgets()
    if err != nil {
        t.Fatalf("failed to list widgets: %v", err)
    }
    if len(widgets) != 2 {
        t.Errorf("expected 2 widgets, got %d", len(widgets))
    }

    // Update widget
    updated, err := widgetSvc.UpdateWidget(widget1.ID, "account_balance", 1, 1, 3, 2, `{"account_id": 2}`)
    if err != nil {
        t.Fatalf("failed to update widget: %v", err)
    }
    if updated.PositionX != 1 {
        t.Errorf("expected position_x 1, got %d", updated.PositionX)
    }
    if updated.Width != 3 {
        t.Errorf("expected width 3, got %d", updated.Width)
    }
    if updated.Config != `{"account_id": 2}` {
        t.Errorf("expected config `{\"account_id\": 2}`, got %q", updated.Config)
    }

    // Delete widget
    err = widgetSvc.DeleteWidget(widget2.ID)
    if err != nil {
        t.Fatalf("failed to delete widget: %v", err)
    }

    // Verify deletion
    widgets, _ = widgetSvc.ListWidgets()
    if len(widgets) != 1 {
        t.Errorf("expected 1 widget after delete, got %d", len(widgets))
    }
}

func TestIntegration_SettingsPersistence(t *testing.T) {
    db := testutil.NewTestDB(t)
    settingsSvc := NewSettingsService(db)

    // Update a setting
    err := settingsSvc.UpdateSetting("default_currency", "EUR")
    if err != nil {
        t.Fatalf("failed to update setting: %v", err)
    }

    // Verify it persists across reads
    settings, err := settingsSvc.GetSettings()
    if err != nil {
        t.Fatalf("failed to get settings: %v", err)
    }

    found := false
    for _, s := range settings {
        if s.Key == "default_currency" && s.Value == "EUR" {
            found = true
            break
        }
    }
    if !found {
        t.Error("expected default_currency=EUR after update")
    }

    // Update another setting
    err = settingsSvc.UpdateSetting("available_currencies", "USD,EUR,JPY")
    if err != nil {
        t.Fatalf("failed to update second setting: %v", err)
    }

    // Verify both settings are present
    settings, _ = settingsSvc.GetSettings()
    foundCurrency := false
    foundAvailable := false
    for _, s := range settings {
        if s.Key == "default_currency" && s.Value == "EUR" {
            foundCurrency = true
        }
        if s.Key == "available_currencies" && s.Value == "USD,EUR,JPY" {
            foundAvailable = true
        }
    }
    if !foundCurrency || !foundAvailable {
        t.Error("expected both settings to persist")
    }
}
