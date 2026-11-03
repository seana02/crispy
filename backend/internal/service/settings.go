package service

import (
    "database/sql"
    "log/slog"

    "github.com/crispy-2/backend/internal/dal/generated"
    "github.com/crispy-2/backend/internal/model"
)

// SettingsService handles application settings.
type SettingsService struct {
    db      *sql.DB
    queries *generated.Queries
}

// NewSettingsService creates a new SettingsService.
func NewSettingsService(db *sql.DB) *SettingsService {
    return &SettingsService{
        db:      db,
        queries: generated.New(db),
    }
}

// GetSettings returns all settings.
func (s *SettingsService) GetSettings() ([]model.Setting, error) {
    rows, err := s.queries.GetAllSettings(ctx())
    if err != nil {
        return nil, model.NewInternalError("failed to get settings: " + err.Error())
    }

    settings := make([]model.Setting, len(rows))
    for i, r := range rows {
        settings[i] = model.Setting{Key: r.Key, Value: r.Value}
    }
    return settings, nil
}

// UpdateSetting updates a single setting.
func (s *SettingsService) UpdateSetting(key, value string) error {
    if key == "" {
        return model.NewValidationError("setting key is required")
    }

    if err := s.queries.UpsertSetting(ctx(), generated.UpsertSettingParams{
        Key:   key,
        Value: value,
    }); err != nil {
        return model.NewInternalError("failed to update setting: " + err.Error())
    }

    slog.Info("setting updated", "key", key)
    return nil
}
