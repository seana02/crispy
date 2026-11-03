package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/crispy-2/backend/internal/model"
	"github.com/crispy-2/backend/internal/service"
)

type createAccountRequest struct {
	Name        string  `json:"name"`
	ParentID    *int64  `json:"parent_id"`
	AccountType string  `json:"account_type"`
}

func createAccountHandler(svc *service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createAccountRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, model.NewValidationError("invalid request body"))
			return
		}

		var account *model.Account
		var err error

		if strings.Contains(req.Name, ":") {
			// Hierarchical path — create the full chain, auto-creating parents
			account, err = svc.CreateAccountByPath(req.Name, req.AccountType)
		} else {
			account, err = svc.CreateAccount(req.Name, req.ParentID, req.AccountType)
		}

		if err != nil {
			writeError(w, err)
			return
		}

		writeCreated(w, account)
	}
}

type updateAccountRequest struct {
	Name        *string  `json:"name,omitempty"`
	ParentID    **int64  `json:"parent_id,omitempty"`
	AccountType string   `json:"account_type,omitempty"`
}

func updateAccountHandler(svc *service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseInt64PathParam(r, "id")
		if err != nil {
			writeError(w, err)
			return
		}

		var req updateAccountRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, model.NewValidationError("invalid request body"))
			return
		}

		account, err := svc.UpdateAccount(id, req.Name, req.ParentID, req.AccountType)
		if err != nil {
			writeError(w, err)
			return
		}

		writeOK(w, account)
	}
}

func getAccountHandler(svc *service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseInt64PathParam(r, "id")
		if err != nil {
			writeError(w, err)
			return
		}

		account, err := svc.GetAccount(id)
		if err != nil {
			writeError(w, err)
			return
		}

		writeOK(w, account)
	}
}

func listAccountsHandler(svc *service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		// Accept both "search" (frontend) and "name_filter" (legacy)
		nameFilter := query.Get("search")
		if nameFilter == "" {
			nameFilter = query.Get("name_filter")
		}
		// Accept both "account_type" (frontend) and "type_filter" (legacy)
		typeFilter := query.Get("account_type")
		if typeFilter == "" {
			typeFilter = query.Get("type_filter")
		}
		pg := model.Pagination{
			Page:     parseInt(query.Get("page"), 1),
			PageSize: parseInt(query.Get("per_page"), parseInt(query.Get("page_size"), 20)),
		}

		result, err := svc.ListAccounts(nameFilter, typeFilter, pg)
		if err != nil {
			writeError(w, err)
			return
		}

		writePaginated(w, result)
	}
}

func deleteAccountHandler(svc *service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseInt64PathParam(r, "id")
		if err != nil {
			writeError(w, err)
			return
		}

		err = svc.DeleteAccount(id)
		if err != nil {
			writeError(w, err)
			return
		}

		writeStatusNoContent(w)
	}
}

func searchAccountsHandler(svc *service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			writeError(w, model.NewValidationError("missing query parameter 'q'"))
			return
		}

		accounts, err := svc.SearchAccounts(query)
		if err != nil {
			writeError(w, err)
			return
		}

		writeOK(w, accounts)
	}
}

func parseInt(s string, def int) int {
	if s == "" {
		return def
	}
	var val int
	err := json.Unmarshal([]byte(s), &val)
	if err != nil {
		return def
	}
	return val
}

func writeStatusNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
