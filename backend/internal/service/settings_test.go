package service

import (
    "testing"

    "github.com/crispy-2/backend/internal/testutil"
)

func TestGetSettings(t *testing.T) {
    db := testutil.NewTestDB(t)
    svc := NewSettingsService(db)

    settings, err := svc.GetSettings()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    // Default settings are seeded by the migration
    if len(settings) == 0 {
        t.Fatal("expected default settings to be seeded")
    }

    // Check for known defaults
    found := false
    for _, s := range settings {
        if s.Key == "default_currency" && s.Value == "USD" {
            found = true
            break
        }
    }
    if !found {
        t.Error("expected default_currency=USD in settings")
    }
}

func TestUpdateSetting(t *testing.T) {
    db := testutil.NewTestDB(t)
    svc := NewSettingsService(db)

    err := svc.UpdateSetting("default_currency", "EUR")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    settings, _ := svc.GetSettings()
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
}

func TestUpdateSetting_EmptyKeyRejected(t *testing.T) {
    db := testutil.NewTestDB(t)
    svc := NewSettingsService(db)

    err := svc.UpdateSetting("", "value")
    if err == nil {
        t.Fatal("expected error for empty key")
    }
}
