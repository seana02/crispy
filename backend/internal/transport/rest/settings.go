package rest

import (
	"encoding/json"
	"net/http"

	"github.com/crispy-2/backend/internal/model"
	"github.com/crispy-2/backend/internal/service"
)

type updateSettingRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func getSettingsHandler(svc *service.SettingsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		settings, err := svc.GetSettings()
		if err != nil {
			writeError(w, err)
			return
		}

		writeOK(w, settings)
	}
}

func updateSettingHandler(svc *service.SettingsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req updateSettingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, model.NewValidationError("invalid request body"))
			return
		}

		err := svc.UpdateSetting(req.Key, req.Value)
		if err != nil {
			writeError(w, err)
			return
		}

		writeStatusNoContent(w)
	}
}

// updateSettingsBulkHandler handles bulk updates of multiple settings at once.
// Frontend sends: {"key1": "value1", "key2": "value2", ...}
func updateSettingsBulkHandler(svc *service.SettingsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req map[string]string
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, model.NewValidationError("invalid request body"))
			return
		}

		for key, value := range req {
			if err := svc.UpdateSetting(key, value); err != nil {
				writeError(w, err)
				return
			}
		}

		writeStatusNoContent(w)
	}
}
