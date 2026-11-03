package rest

import (
	"net/http"

	"github.com/crispy-2/backend/internal/model"
	"github.com/crispy-2/backend/internal/service"
)

func globalSearchHandler(svc *service.SearchService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			writeError(w, model.NewValidationError("missing query parameter 'q'"))
			return
		}

		result, err := svc.GlobalSearch(query)
		if err != nil {
			writeError(w, err)
			return
		}

		writeOK(w, result)
	}
}
