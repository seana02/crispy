package service

import (
    "database/sql"
    "log/slog"
    "time"

    "github.com/crispy-2/backend/internal/dal/generated"
    "github.com/crispy-2/backend/internal/model"
)

// WidgetService handles home page widget configuration.
type WidgetService struct {
    db      *sql.DB
    queries *generated.Queries
}

// NewWidgetService creates a new WidgetService.
func NewWidgetService(db *sql.DB) *WidgetService {
    return &WidgetService{
        db:      db,
        queries: generated.New(db),
    }
}

// CreateWidget creates a new widget configuration.
func (s *WidgetService) CreateWidget(widgetType string, posX, posY, width, height int, config string) (*model.WidgetConfig, error) {
    if widgetType == "" {
        return nil, model.NewValidationError("widget_type is required")
    }
    if config == "" {
        config = "{}"
    }

    row, err := s.queries.CreateWidget(ctx(), generated.CreateWidgetParams{
        WidgetType: widgetType,
        PositionX:  int64(posX),
        PositionY:  int64(posY),
        Width:      int64(width),
        Height:     int64(height),
        Config:     config,
    })
    if err != nil {
        return nil, model.NewInternalError("failed to create widget: " + err.Error())
    }

    slog.Info("widget created", "id", row.ID, "type", widgetType)
    return toModelWidget(row), nil
}

// ListWidgets returns all widget configurations.
func (s *WidgetService) ListWidgets() ([]model.WidgetConfig, error) {
    rows, err := s.queries.ListWidgets(ctx())
    if err != nil {
        return nil, model.NewInternalError("failed to list widgets: " + err.Error())
    }

    widgets := make([]model.WidgetConfig, len(rows))
    for i, r := range rows {
        widgets[i] = *toModelWidget(r)
    }
    return widgets, nil
}

// UpdateWidget updates a widget configuration.
func (s *WidgetService) UpdateWidget(id int64, widgetType string, posX, posY, width, height int, config string) (*model.WidgetConfig, error) {
    _, err := s.queries.GetWidget(ctx(), id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, model.NewNotFoundError("widget", id)
        }
        return nil, model.NewInternalError("failed to get widget: " + err.Error())
    }

    row, err := s.queries.UpdateWidget(ctx(), generated.UpdateWidgetParams{
        WidgetType: widgetType,
        PositionX:  int64(posX),
        PositionY:  int64(posY),
        Width:      int64(width),
        Height:     int64(height),
        Config:     config,
        UpdatedAt:  time.Now().Format(time.RFC3339),
        ID:         id,
    })
    if err != nil {
        return nil, model.NewInternalError("failed to update widget: " + err.Error())
    }

    slog.Info("widget updated", "id", id)
    return toModelWidget(row), nil
}

// DeleteWidget deletes a widget configuration.
func (s *WidgetService) DeleteWidget(id int64) error {
    _, err := s.queries.GetWidget(ctx(), id)
    if err != nil {
        if err == sql.ErrNoRows {
            return model.NewNotFoundError("widget", id)
        }
        return model.NewInternalError("failed to get widget: " + err.Error())
    }

    if err := s.queries.DeleteWidget(ctx(), id); err != nil {
        return model.NewInternalError("failed to delete widget: " + err.Error())
    }

    slog.Info("widget deleted", "id", id)
    return nil
}

func toModelWidget(row generated.WidgetConfiguration) *model.WidgetConfig {
	return &model.WidgetConfig{
		ID:         row.ID,
		WidgetType: row.WidgetType,
		PositionX:  int(row.PositionX),
		PositionY:  int(row.PositionY),
		Width:      int(row.Width),
		Height:     int(row.Height),
		Config:     row.Config,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}
}
