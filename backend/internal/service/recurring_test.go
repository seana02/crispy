package service

import (
    "testing"

    "github.com/crispy-2/backend/internal/model"
    "github.com/crispy-2/backend/internal/testutil"
)

func TestCreateTemplate(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    recSvc := NewRecurringService(db, NewTransactionService(db, nil))

    assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
    checking, _ := acctSvc.CreateAccount("Checking", &assets.ID, "")
    expenses, _ := acctSvc.CreateAccount("Expenses", nil, "expense")

    tmpl, err := recSvc.CreateTemplate("Monthly Rent", "posted", "every_n_months", 1, nil, "2026-01-01", nil, []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-1500.00"},
        {AccountID: expenses.ID, Currency: "USD", Amount: "1500.00"},
    }, nil)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if tmpl.Description != "Monthly Rent" {
        t.Errorf("expected 'Monthly Rent', got %q", tmpl.Description)
    }
    if tmpl.ScheduleType != "every_n_months" {
        t.Errorf("expected 'every_n_months', got %q", tmpl.ScheduleType)
    }
    if len(tmpl.Postings) != 2 {
        t.Errorf("expected 2 postings, got %d", len(tmpl.Postings))
    }
}

func TestCreateTemplate_InvalidScheduleTypeRejected(t *testing.T) {
    db := testutil.NewTestDB(t)
    recSvc := NewRecurringService(db, NewTransactionService(db, nil))

    _, err := recSvc.CreateTemplate("Bad", "posted", "invalid_type", 1, nil, "2026-01-01", nil, []PostingInput{
        {AccountID: 1, Currency: "USD", Amount: "-100.00"},
        {AccountID: 2, Currency: "USD", Amount: "100.00"},
    }, nil)
    if err == nil {
        t.Fatal("expected error for invalid schedule type")
    }
}

func TestPreviewSchedule_EveryNDays(t *testing.T) {
    db := testutil.NewTestDB(t)
    recSvc := NewRecurringService(db, NewTransactionService(db, nil))

    dates := recSvc.PreviewSchedule("every_n_days", 7, nil, "2026-01-01", nil, 5)
    if len(dates) != 5 {
        t.Fatalf("expected 5 dates, got %d", len(dates))
    }
    expected := []string{"2026-01-01", "2026-01-08", "2026-01-15", "2026-01-22", "2026-01-29"}
    for i, d := range dates {
        if d != expected[i] {
            t.Errorf("date[%d]: expected %s, got %s", i, expected[i], d)
        }
    }
}

func TestPreviewSchedule_EveryNWeeks(t *testing.T) {
    db := testutil.NewTestDB(t)
    recSvc := NewRecurringService(db, NewTransactionService(db, nil))

    dates := recSvc.PreviewSchedule("every_n_weeks", 2, nil, "2026-01-01", nil, 3)
    expected := []string{"2026-01-01", "2026-01-15", "2026-01-29"}
    if len(dates) != 3 {
        t.Fatalf("expected 3 dates, got %d", len(dates))
    }
    for i, d := range dates {
        if d != expected[i] {
            t.Errorf("date[%d]: expected %s, got %s", i, expected[i], d)
        }
    }
}

func TestPreviewSchedule_MonthlyOnDay(t *testing.T) {
    db := testutil.NewTestDB(t)
    recSvc := NewRecurringService(db, NewTransactionService(db, nil))

    day := 15
    dates := recSvc.PreviewSchedule("monthly_on_day", 1, &day, "2026-01-15", nil, 4)
    expected := []string{"2026-01-15", "2026-02-15", "2026-03-15", "2026-04-15"}
    if len(dates) != 4 {
        t.Fatalf("expected 4 dates, got %d", len(dates))
    }
    for i, d := range dates {
        if d != expected[i] {
            t.Errorf("date[%d]: expected %s, got %s", i, expected[i], d)
        }
    }
}

func TestPreviewSchedule_EveryNMonths(t *testing.T) {
    db := testutil.NewTestDB(t)
    recSvc := NewRecurringService(db, NewTransactionService(db, nil))

    dates := recSvc.PreviewSchedule("every_n_months", 3, nil, "2026-01-15", nil, 4)
    expected := []string{"2026-01-15", "2026-04-15", "2026-07-15", "2026-10-15"}
    if len(dates) != 4 {
        t.Fatalf("expected 4 dates, got %d", len(dates))
    }
    for i, d := range dates {
        if d != expected[i] {
            t.Errorf("date[%d]: expected %s, got %s", i, expected[i], d)
        }
    }
}

func TestPreviewSchedule_RespectsEndDate(t *testing.T) {
    db := testutil.NewTestDB(t)
    recSvc := NewRecurringService(db, NewTransactionService(db, nil))

    endDate := "2026-01-20"
    dates := recSvc.PreviewSchedule("every_n_days", 7, nil, "2026-01-01", &endDate, 10)
    // Only 3 dates should fit: Jan 1, Jan 8, Jan 15 (Jan 22 > end date)
    if len(dates) != 3 {
        t.Errorf("expected 3 dates before end date, got %d: %v", len(dates), dates)
    }
}

func TestGenerateDueTransactions(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    txSvc := NewTransactionService(db, nil)
    recSvc := NewRecurringService(db, txSvc)

    assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
    checking, _ := acctSvc.CreateAccount("Checking", &assets.ID, "")
    expenses, _ := acctSvc.CreateAccount("Expenses", nil, "expense")

    _, err := recSvc.CreateTemplate("Weekly Groceries", "posted", "every_n_weeks", 1, nil, "2026-10-01", nil, []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-100.00"},
        {AccountID: expenses.ID, Currency: "USD", Amount: "100.00"},
    }, nil)
    if err != nil {
        t.Fatalf("failed to create template: %v", err)
    }

    // Generate for "today" = 2026-10-29 (should generate Oct 1, 8, 15, 22, 29 = 5 transactions)
    count, err := recSvc.GenerateDueTransactions("2026-10-29")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if count != 5 {
        t.Errorf("expected 5 generated transactions, got %d", count)
    }

    // Verify transactions were actually created
    result, _ := txSvc.ListTransactions(TransactionFilters{}, model.Pagination{Page: 1, PageSize: 20})
    if result.Total != 5 {
        t.Errorf("expected 5 transactions in DB, got %d", result.Total)
    }
}

func TestGenerateDueTransactions_DoesNotDuplicate(t *testing.T) {
    db := testutil.NewTestDB(t)
    acctSvc := NewAccountService(db)
    txSvc := NewTransactionService(db, nil)
    recSvc := NewRecurringService(db, txSvc)

    assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
    checking, _ := acctSvc.CreateAccount("Checking", &assets.ID, "")
    expenses, _ := acctSvc.CreateAccount("Expenses", nil, "expense")

    _, err := recSvc.CreateTemplate("Weekly", "posted", "every_n_weeks", 1, nil, "2026-10-01", nil, []PostingInput{
        {AccountID: checking.ID, Currency: "USD", Amount: "-50.00"},
        {AccountID: expenses.ID, Currency: "USD", Amount: "50.00"},
    }, nil)
    if err != nil {
        t.Fatalf("failed to create template: %v", err)
    }

    // Generate once
    recSvc.GenerateDueTransactions("2026-10-15")

    // Generate again for the same date - should not create duplicates
    count, err := recSvc.GenerateDueTransactions("2026-10-15")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if count != 0 {
        t.Errorf("expected 0 new transactions on second run, got %d", count)
    }
}
