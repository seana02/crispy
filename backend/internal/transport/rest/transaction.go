package rest

import (
	"encoding/json"
	"net/http"

	"github.com/crispy-2/backend/internal/model"
	"github.com/crispy-2/backend/internal/service"
)

type createTransactionRequest struct {
	Date        string                `json:"date"`
	Description string                `json:"description"`
	Status      string                `json:"status"`
	Postings    []service.PostingInput `json:"postings"`
	Tags        []service.TagInput    `json:"tags"`
}

func createTransactionHandler(svc *service.TransactionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createTransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, model.NewValidationError("invalid request body"))
			return
		}

		transaction, err := svc.CreateTransaction(req.Date, req.Description, req.Status, req.Postings, req.Tags)
		if err != nil {
			writeError(w, err)
			return
		}

		writeCreated(w, transaction)
	}
}

type updateTransactionRequest struct {
	Date        string                `json:"date"`
	Description string                `json:"description"`
	Status      string                `json:"status"`
	Postings    []service.PostingInput `json:"postings"`
	Tags        []service.TagInput    `json:"tags"`
}

func updateTransactionHandler(svc *service.TransactionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseInt64PathParam(r, "id")
		if err != nil {
			writeError(w, err)
			return
		}

		var req updateTransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, model.NewValidationError("invalid request body"))
			return
		}

		transaction, err := svc.UpdateTransaction(id, req.Date, req.Description, req.Status, req.Postings, req.Tags)
		if err != nil {
			writeError(w, err)
			return
		}

		writeOK(w, transaction)
	}
}

func getTransactionHandler(svc *service.TransactionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseInt64PathParam(r, "id")
		if err != nil {
			writeError(w, err)
			return
		}

		transaction, err := svc.GetTransaction(id)
		if err != nil {
			writeError(w, err)
			return
		}

		writeOK(w, transaction)
	}
}

type listTransactionsHandlerWrapper struct {
	svc *service.TransactionService
}

func listTransactionsHandler(svc *service.TransactionService) http.HandlerFunc {
	return listTransactionsHandlerWrapper{svc: svc}.handle
}

func (h listTransactionsHandlerWrapper) handle(w http.ResponseWriter, r *http.Request) {
	filters := service.TransactionFilters{}
	query := r.URL.Query()

	filters.Description = query.Get("description")
	filters.DateFrom = query.Get("date_from")
	filters.DateTo = query.Get("date_to")
	filters.Status = query.Get("status")
	filters.TagKey = query.Get("tag_key")
	filters.TagValue = query.Get("tag_value")
	filters.Currency = query.Get("currency")

	if accountIDStr := query.Get("account_id"); accountIDStr != "" {
		accountID, err := deserializeInt64(accountIDStr)
		if err != nil {
			writeError(w, model.NewValidationError("invalid account_id: "+err.Error()))
			return
		}
		filters.AccountID = &accountID
	}

	pg := model.Pagination{
		Page:     parseInt(query.Get("page"), 1),
		PageSize: parseInt(query.Get("per_page"), parseInt(query.Get("page_size"), 20)),
	}

	result, err := h.svc.ListTransactions(filters, pg)
	if err != nil {
		writeError(w, err)
		return
	}

	writePaginated(w, result)
}

func deleteTransactionHandler(svc *service.TransactionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseInt64PathParam(r, "id")
		if err != nil {
			writeError(w, err)
			return
		}

		err = svc.DeleteTransaction(id)
		if err != nil {
			writeError(w, err)
			return
		}

		writeStatusNoContent(w)
	}
}
