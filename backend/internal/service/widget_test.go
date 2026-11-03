package service

import (
    "testing"

    "github.com/crispy-2/backend/internal/testutil"
)

func TestCreateWidget(t *testing.T) {
    db := testutil.NewTestDB(t)
    svc := NewWidgetService(db)

    widget, err := svc.CreateWidget("account_balance", 0, 0, 2, 1, `{"account_id": 1}`)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if widget.WidgetType != "account_balance" {
        t.Errorf("expected type 'account_balance', got %q", widget.WidgetType)
    }
    if widget.Width != 2 {
        t.Errorf("expected width 2, got %d", widget.Width)
    }
}

func TestCreateWidget_EmptyTypeRejected(t *testing.T) {
    db := testutil.NewTestDB(t)
    svc := NewWidgetService(db)

    _, err := svc.CreateWidget("", 0, 0, 1, 1, "{}")
    if err == nil {
        t.Fatal("expected error for empty widget type")
    }
}

func TestListWidgets(t *testing.T) {
    db := testutil.NewTestDB(t)
    svc := NewWidgetService(db)

    svc.CreateWidget("account_balance", 0, 0, 2, 1, `{"account_id": 1}`)
    svc.CreateWidget("budget_progress", 2, 0, 2, 1, `{"budget_id": 1}`)

    widgets, err := svc.ListWidgets()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(widgets) != 2 {
        t.Errorf("expected 2 widgets, got %d", len(widgets))
    }
}

func TestUpdateWidget(t *testing.T) {
    db := testutil.NewTestDB(t)
    svc := NewWidgetService(db)

    created, _ := svc.CreateWidget("account_balance", 0, 0, 2, 1, `{"account_id": 1}`)
    updated, err := svc.UpdateWidget(created.ID, "account_balance", 1, 1, 3, 2, `{"account_id": 2}`)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if updated.PositionX != 1 {
        t.Errorf("expected position_x 1, got %d", updated.PositionX)
    }
    if updated.Width != 3 {
        t.Errorf("expected width 3, got %d", updated.Width)
    }
}

func TestDeleteWidget(t *testing.T) {
    db := testutil.NewTestDB(t)
    svc := NewWidgetService(db)

    created, _ := svc.CreateWidget("account_balance", 0, 0, 2, 1, `{}`)
    err := svc.DeleteWidget(created.ID)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    widgets, _ := svc.ListWidgets()
    if len(widgets) != 0 {
        t.Errorf("expected 0 widgets after delete, got %d", len(widgets))
    }
}

func TestDeleteWidget_NotFound(t *testing.T) {
    db := testutil.NewTestDB(t)
    svc := NewWidgetService(db)

    err := svc.DeleteWidget(9999)
    if err == nil {
        t.Fatal("expected not found error")
    }
}
