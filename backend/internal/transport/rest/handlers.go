package rest

import (
	"encoding/json"
	"net/http"

	"github.com/crispy-2/backend/internal/model"
)

// paginatedResponse is the wire format for all paginated list responses.
// Frontend expects: {"data": [...], "meta": {"page": N, "per_page": N, "total": N, "total_pages": N}}
type paginatedMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type paginatedResponse[T any] struct {
	Data []T           `json:"data"`
	Meta paginatedMeta `json:"meta"`
}

// writePaginated serialises a model.PaginatedResult in the frontend-expected wire format.
func writePaginated[T any](w http.ResponseWriter, result *model.PaginatedResult[T]) {
	resp := paginatedResponse[T]{
		Data: result.Items,
		Meta: paginatedMeta{
			Page:       result.Page,
			PerPage:    result.PageSize,
			Total:      result.Total,
			TotalPages: result.TotalPages,
		},
	}
	if resp.Data == nil {
		resp.Data = []T{}
	}
	writeOK(w, resp)
}

type ResponseWriter struct {
	w http.ResponseWriter
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w: w}
}

func (rw *ResponseWriter) writeJSON(status int, data any) {
	rw.w.Header().Set("Content-Type", "application/json")
	rw.w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(rw.w).Encode(data)
	}
}

func (rw *ResponseWriter) writeOK(data any) {
	rw.writeJSON(http.StatusOK, data)
}

func (rw *ResponseWriter) writeCreated(data any) {
	rw.writeJSON(http.StatusCreated, data)
}

func (rw *ResponseWriter) writeError(err error) {
	appErr, ok := err.(*model.AppError)
	if ok {
		rw.writeJSON(appErr.Code, appErr)
		return
	}
	rw.writeJSON(http.StatusInternalServerError, model.NewInternalError(err.Error()))
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	NewResponseWriter(w).writeJSON(status, data)
}

func writeOK(w http.ResponseWriter, data any) {
	NewResponseWriter(w).writeJSON(http.StatusOK, data)
}

func writeCreated(w http.ResponseWriter, data any) {
	NewResponseWriter(w).writeJSON(http.StatusCreated, data)
}

func writeError(w http.ResponseWriter, err error) {
	NewResponseWriter(w).writeError(err)
}

func parseInt64PathParam(r *http.Request, param string) (int64, error) {
	var id int64
	val := r.PathValue(param)
	if val == "" {
		return 0, model.NewValidationError("missing path parameter: " + param)
	}
	var err error
	id, err = deserializeInt64(val)
	if err != nil {
		return 0, model.NewValidationError("invalid " + param + ": " + err.Error())
	}
	return id, nil
}

func deserializeInt64(s string) (int64, error) {
	var id int64
	err := json.Unmarshal([]byte(s), &id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func deserializeString(s string) (string, error) {
	var val string
	err := json.Unmarshal([]byte(s), &val)
	if err != nil {
		return "", err
	}
	return val, nil
}

func deserializeBool(s string) (bool, error) {
	var val bool
	err := json.Unmarshal([]byte(s), &val)
	if err != nil {
		return false, err
	}
	return val, nil
}

func deserializeSlice[T any](r *http.Request, param string) ([]T, error) {
	var val []T
	queryVal := r.URL.Query().Get(param)
	if queryVal == "" {
		return nil, nil
	}
	err := json.Unmarshal([]byte(queryVal), &val)
	if err != nil {
		return nil, err
	}
	return val, nil
}

