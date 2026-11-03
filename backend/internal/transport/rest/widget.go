package rest

import (
	"encoding/json"
	"net/http"

	"github.com/crispy-2/backend/internal/model"
	"github.com/crispy-2/backend/internal/service"
)

type createWidgetRequest struct {
	WidgetType string `json:"widget_type"`
	PositionX  int    `json:"position_x"`
	PositionY  int    `json:"position_y"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Config     string `json:"config"`
}

func createWidgetHandler(svc *service.WidgetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createWidgetRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, model.NewValidationError("invalid request body"))
			return
		}

		widget, err := svc.CreateWidget(req.WidgetType, req.PositionX, req.PositionY, req.Width, req.Height, req.Config)
		if err != nil {
			writeError(w, err)
			return
		}

		writeCreated(w, widget)
	}
}

type updateWidgetRequest struct {
	WidgetType string `json:"widget_type"`
	PositionX  int    `json:"position_x"`
	PositionY  int    `json:"position_y"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Config     string `json:"config"`
}

func updateWidgetHandler(svc *service.WidgetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseInt64PathParam(r, "id")
		if err != nil {
			writeError(w, err)
			return
		}

		var req updateWidgetRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, model.NewValidationError("invalid request body"))
			return
		}

		widget, err := svc.UpdateWidget(id, req.WidgetType, req.PositionX, req.PositionY, req.Width, req.Height, req.Config)
		if err != nil {
			writeError(w, err)
			return
		}

		writeOK(w, widget)
	}
}

func listWidgetsHandler(svc *service.WidgetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		widgets, err := svc.ListWidgets()
		if err != nil {
			writeError(w, err)
			return
		}

		writeOK(w, widgets)
	}
}

func deleteWidgetHandler(svc *service.WidgetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseInt64PathParam(r, "id")
		if err != nil {
			writeError(w, err)
			return
		}

		err = svc.DeleteWidget(id)
		if err != nil {
			writeError(w, err)
			return
		}

		writeStatusNoContent(w)
	}
}
