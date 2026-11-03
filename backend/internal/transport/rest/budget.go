package rest

import (
	"encoding/json"
	"net/http"

	"github.com/crispy-2/backend/internal/model"
	"github.com/crispy-2/backend/internal/service"
)

type createBudgetRequest struct {
	Name         string   `json:"name"`
	Direction    string   `json:"direction"`
	Period       string   `json:"period"`
	Currency     string   `json:"currency"`
	TargetAmount string   `json:"target_amount"`
	AccountIDs   []int64  `json:"account_ids"`
	TagIDs       []int64  `json:"tag_ids"`
}

func createBudgetHandler(svc *service.BudgetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createBudgetRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, model.NewValidationError("invalid request body"))
			return
		}

		budget, err := svc.CreateBudget(req.Name, req.Direction, req.Period, req.Currency, req.TargetAmount, req.AccountIDs, req.TagIDs)
		if err != nil {
			writeError(w, err)
			return
		}

		writeCreated(w, budget)
	}
}

type updateBudgetRequest struct {
	Name         string  `json:"name"`
	Direction    string  `json:"direction"`
	Period       string  `json:"period"`
	Currency     string  `json:"currency"`
	TargetAmount string  `json:"target_amount"`
	AccountIDs   []int64 `json:"account_ids"`
	TagIDs       []int64 `json:"tag_ids"`
}

func updateBudgetHandler(svc *service.BudgetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseInt64PathParam(r, "id")
		if err != nil {
			writeError(w, err)
			return
		}

		var req updateBudgetRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, model.NewValidationError("invalid request body"))
			return
		}

		budget, err := svc.UpdateBudget(id, req.Name, req.Direction, req.Period, req.Currency, req.TargetAmount, req.AccountIDs, req.TagIDs)
		if err != nil {
			writeError(w, err)
			return
		}

		writeOK(w, budget)
	}
}

func getBudgetHandler(svc *service.BudgetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseInt64PathParam(r, "id")
		if err != nil {
			writeError(w, err)
			return
		}

		budget, err := svc.GetBudget(id)
		if err != nil {
			writeError(w, err)
			return
		}

		writeOK(w, budget)
	}
}

func listBudgetsHandler(svc *service.BudgetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		budgets, err := svc.ListBudgets()
		if err != nil {
			writeError(w, err)
			return
		}

		// Wrap as paginated response to match frontend expectations
		page := parseInt(r.URL.Query().Get("page"), 1)
		perPage := parseInt(r.URL.Query().Get("per_page"), 50)
		total := len(budgets)
		totalPages := total / perPage
		if total%perPage != 0 {
			totalPages++
		}
		if totalPages == 0 {
			totalPages = 1
		}
		result := &model.PaginatedResult[model.Budget]{
			Items:      budgets,
			Total:      total,
			Page:       page,
			PageSize:   perPage,
			TotalPages: totalPages,
		}
		writePaginated(w, result)
	}
}

func deleteBudgetHandler(svc *service.BudgetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseInt64PathParam(r, "id")
		if err != nil {
			writeError(w, err)
			return
		}

		err = svc.DeleteBudget(id)
		if err != nil {
			writeError(w, err)
			return
		}

		writeStatusNoContent(w)
	}
}

func computeBudgetCurrentHandler(svc *service.BudgetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseInt64PathParam(r, "id")
		if err != nil {
			writeError(w, err)
			return
		}

		dateFrom := r.URL.Query().Get("date_from")
		dateTo := r.URL.Query().Get("date_to")

		if dateFrom == "" || dateTo == "" {
			writeError(w, model.NewValidationError("date_from and date_to are required"))
			return
		}

		current, err := svc.ComputeBudgetCurrent(id, dateFrom, dateTo)
		if err != nil {
			writeError(w, err)
			return
		}

		writeOK(w, map[string]string{
			"current_amount": current.String(),
		})
	}
}
