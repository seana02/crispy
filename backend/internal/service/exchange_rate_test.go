package service

import (
    "testing"

    "github.com/crispy-2/backend/internal/model"
    "github.com/crispy-2/backend/internal/testutil"
    "github.com/shopspring/decimal"
)

func TestCreateRate_Success(t *testing.T) {
    db := testutil.NewTestDB(t)
    svc := NewExchangeRateService(db)

    rate, err := svc.CreateRate("USD", "JPY", "150.50", "2026-01-01", "manual")
    if err != nil {
        t.Fatalf("failed to create rate: %v", err)
    }
    if rate.BaseCurrency != "USD" {
        t.Errorf("expected base USD, got %q", rate.BaseCurrency)
    }
    if rate.QuoteCurrency != "JPY" {
        t.Errorf("expected quote JPY, got %q", rate.QuoteCurrency)
    }
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("150.50"), rate.Rate.Decimal)
    if rate.Source != "manual" {
        t.Errorf("expected source manual, got %q", rate.Source)
    }
}

func TestCreateRate_SameCurrencyRejected(t *testing.T) {
    db := testutil.NewTestDB(t)
    svc := NewExchangeRateService(db)

    _, err := svc.CreateRate("USD", "USD", "1.0", "2026-01-01", "manual")
    if err == nil {
        t.Fatal("expected error for same currency pair")
    }
}

func TestCreateRate_InvalidRateRejected(t *testing.T) {
    db := testutil.NewTestDB(t)
    svc := NewExchangeRateService(db)

    _, err := svc.CreateRate("USD", "JPY", "0", "2026-01-01", "manual")
    if err == nil {
        t.Fatal("expected error for zero rate")
    }

    _, err = svc.CreateRate("USD", "JPY", "-10", "2026-01-01", "manual")
    if err == nil {
        t.Fatal("expected error for negative rate")
    }
}

func TestGetRate_DirectPair(t *testing.T) {
    db := testutil.NewTestDB(t)
    svc := NewExchangeRateService(db)

    svc.CreateRate("USD", "JPY", "150.00", "2026-01-01", "manual")
    svc.CreateRate("USD", "JPY", "152.00", "2026-01-15", "manual")

    // Should get the most recent rate on or before the queried date
    rate, err := svc.GetRate("USD", "JPY", "2026-01-10")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("150.00"), rate.Rate.Decimal)

    // On Jan 20, should get the Jan 15 rate
    rate, err = svc.GetRate("USD", "JPY", "2026-01-20")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    testutil.AssertDecimalEqual(t, decimal.RequireFromString("152.00"), rate.Rate.Decimal)
}

func TestGetRate_InverseFallback(t *testing.T) {
    db := testutil.NewTestDB(t)
    svc := NewExchangeRateService(db)

    svc.CreateRate("USD", "JPY", "150.00", "2026-01-01", "manual")

    // Query the inverse pair - should compute 1/150
    rate, err := svc.GetRate("JPY", "USD", "2026-01-01")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    expected := decimal.NewFromInt(1).Div(decimal.RequireFromString("150.00"))
    testutil.AssertDecimalEqual(t, expected, rate.Rate.Decimal)
}

func TestGetRate_SameCurrencyReturnsOne(t *testing.T) {
    db := testutil.NewTestDB(t)
    svc := NewExchangeRateService(db)

    rate, err := svc.GetRate("USD", "USD", "2026-01-01")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    testutil.AssertDecimalEqual(t, decimal.NewFromInt(1), rate.Rate.Decimal)
}

func TestGetRate_NotFound(t *testing.T) {
    db := testutil.NewTestDB(t)
    svc := NewExchangeRateService(db)

    _, err := svc.GetRate("USD", "JPY", "2026-01-01")
    if err == nil {
        t.Fatal("expected not found error")
    }
    appErr, ok := err.(*model.AppError)
    if !ok {
        t.Fatalf("expected *model.AppError, got %T", err)
    }
    if appErr.Code != 404 {
        t.Errorf("expected code 404, got %d", appErr.Code)
    }
}

func TestListRates(t *testing.T) {
    db := testutil.NewTestDB(t)
    svc := NewExchangeRateService(db)

    svc.CreateRate("USD", "JPY", "150.00", "2026-01-01", "manual")
    svc.CreateRate("USD", "EUR", "0.92", "2026-01-01", "manual")

    result, err := svc.ListRates("", "", model.Pagination{Page: 1, PageSize: 10})
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if result.Total != 2 {
        t.Errorf("expected 2 rates, got %d", result.Total)
    }
}

func TestListRates_FilterByPair(t *testing.T) {
    db := testutil.NewTestDB(t)
    svc := NewExchangeRateService(db)

    svc.CreateRate("USD", "JPY", "150.00", "2026-01-01", "manual")
    svc.CreateRate("USD", "EUR", "0.92", "2026-01-01", "manual")

    result, err := svc.ListRates("USD", "JPY", model.Pagination{Page: 1, PageSize: 10})
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if result.Total != 1 {
        t.Errorf("expected 1 rate, got %d", result.Total)
    }
}

func TestDeleteRate(t *testing.T) {
    db := testutil.NewTestDB(t)
    svc := NewExchangeRateService(db)

    rate, _ := svc.CreateRate("USD", "JPY", "150.00", "2026-01-01", "manual")
    err := svc.DeleteRate(rate.ID)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    result, _ := svc.ListRates("", "", model.Pagination{Page: 1, PageSize: 10})
    if result.Total != 0 {
        t.Errorf("expected 0 rates after delete, got %d", result.Total)
    }
}

func TestDeleteRate_NotFound(t *testing.T) {
    db := testutil.NewTestDB(t)
    svc := NewExchangeRateService(db)

    err := svc.DeleteRate(9999)
    if err == nil {
        t.Fatal("expected not found error")
    }
}
