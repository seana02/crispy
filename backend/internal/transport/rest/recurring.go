package rest

import (
	"encoding/json"
	"net/http"

	"github.com/crispy-2/backend/internal/model"
	"github.com/crispy-2/backend/internal/service"
)

type createRecurringTemplateRequest struct {
	Description  string                `json:"description"`
	Status       string                `json:"status"`
	ScheduleType string                `json:"schedule_type"`
	IntervalN    int                   `json:"interval_n"`
	DayOfMonth   *int                  `json:"day_of_month"`
	StartDate    string                `json:"start_date"`
	EndDate      *string               `json:"end_date"`
	Postings     []service.PostingInput `json:"postings"`
	Tags         []service.TagInput    `json:"tags"`
}

func createRecurringTemplateHandler(svc *service.RecurringService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createRecurringTemplateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, model.NewValidationError("invalid request body"))
			return
		}

		template, err := svc.CreateTemplate(
			req.Description,
			req.Status,
			req.ScheduleType,
			req.IntervalN,
			req.DayOfMonth,
			req.StartDate,
			req.EndDate,
			req.Postings,
			req.Tags,
		)
		if err != nil {
			writeError(w, err)
			return
		}

		writeCreated(w, template)
	}
}

type updateRecurringTemplateRequest struct {
	Description  string                `json:"description"`
	Status       string                `json:"status"`
	ScheduleType string                `json:"schedule_type"`
	IntervalN    int                   `json:"interval_n"`
	DayOfMonth   *int                  `json:"day_of_month"`
	StartDate    string                `json:"start_date"`
	EndDate      *string               `json:"end_date"`
	IsActive     bool                  `json:"is_active"`
	Postings     []service.PostingInput `json:"postings"`
	Tags         []service.TagInput    `json:"tags"`
}

func updateRecurringTemplateHandler(svc *service.RecurringService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseInt64PathParam(r, "id")
		if err != nil {
			writeError(w, err)
			return
		}

		var req updateRecurringTemplateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, model.NewValidationError("invalid request body"))
			return
		}

		template, err := svc.UpdateTemplate(
			id,
			req.Description,
			req.Status,
			req.ScheduleType,
			req.IntervalN,
			req.DayOfMonth,
			req.StartDate,
			req.EndDate,
			req.IsActive,
			req.Postings,
			req.Tags,
		)
		if err != nil {
			writeError(w, err)
			return
		}

		writeOK(w, template)
	}
}

func getRecurringTemplateHandler(svc *service.RecurringService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseInt64PathParam(r, "id")
		if err != nil {
			writeError(w, err)
			return
		}

		template, err := svc.GetTemplate(id)
		if err != nil {
			writeError(w, err)
			return
		}

		writeOK(w, template)
	}
}

func listRecurringTemplatesHandler(svc *service.RecurringService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templates, err := svc.ListTemplates()
		if err != nil {
			writeError(w, err)
			return
		}

		page := parseInt(r.URL.Query().Get("page"), 1)
		perPage := parseInt(r.URL.Query().Get("per_page"), 50)
		total := len(templates)
		totalPages := total / perPage
		if total%perPage != 0 {
			totalPages++
		}
		if totalPages == 0 {
			totalPages = 1
		}
		result := &model.PaginatedResult[model.RecurringTemplate]{
			Items:      templates,
			Total:      total,
			Page:       page,
			PageSize:   perPage,
			TotalPages: totalPages,
		}
		writePaginated(w, result)
	}
}

func deleteRecurringTemplateHandler(svc *service.RecurringService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseInt64PathParam(r, "id")
		if err != nil {
			writeError(w, err)
			return
		}

		err = svc.DeleteTemplate(id)
		if err != nil {
			writeError(w, err)
			return
		}

		writeStatusNoContent(w)
	}
}

type previewScheduleRequest struct {
	ScheduleType string `json:"schedule_type"`
	IntervalN    int    `json:"interval_n"`
	DayOfMonth   *int   `json:"day_of_month"`
	StartDate    string `json:"start_date"`
	EndDate      *string `json:"end_date"`
	Count        int    `json:"count"`
}

func previewRecurringScheduleHandler(svc *service.RecurringService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req previewScheduleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, model.NewValidationError("invalid request body"))
			return
		}

		dates := svc.PreviewSchedule(
			req.ScheduleType,
			req.IntervalN,
			req.DayOfMonth,
			req.StartDate,
			req.EndDate,
			req.Count,
		)

		writeOK(w, map[string][]string{
			"scheduled_dates": dates,
		})
	}
}

type generateRecurringRequest struct {
	AsOfDate string `json:"as_of_date"`
}

func generateRecurringTransactionsHandler(svc *service.RecurringService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req generateRecurringRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, model.NewValidationError("invalid request body"))
			return
		}

		count, err := svc.GenerateDueTransactions(req.AsOfDate)
		if err != nil {
			writeError(w, err)
			return
		}

		writeOK(w, map[string]int{
			"generated": count,
		})
	}
}
