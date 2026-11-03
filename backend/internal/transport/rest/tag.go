package rest

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/crispy-2/backend/internal/dal/generated"
	"github.com/crispy-2/backend/internal/model"
)

// TagHandler provides HTTP handlers for the tags resource.
// Tags are managed inline with transactions; this handler only provides listing/search.
type TagHandler struct {
	queries *generated.Queries
}

func NewTagHandler(db *sql.DB) *TagHandler {
	return &TagHandler{queries: generated.New(db)}
}

// tagWithCount is the wire type for a tag with its usage count.
type tagWithCount struct {
	ID               int64  `json:"id"`
	Key              string `json:"key"`
	Value            string `json:"value"`
	TransactionCount int64  `json:"transaction_count"`
}

func (h *TagHandler) listTagsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := context.Background()
		query := r.URL.Query()
		search := query.Get("search")
		page := parseInt(query.Get("page"), 1)
		perPage := parseInt(query.Get("per_page"), parseInt(query.Get("page_size"), 50))
		offset := int64((page - 1) * perPage)
		limit := int64(perPage)

		var rows []tagWithCount
		var total int64

		if search != "" {
			pattern := "%" + search + "%"
			rawRows, err := h.queries.ListTagsByKey(rctx, generated.ListTagsByKeyParams{
				Key:    pattern,
				Limit:  limit,
				Offset: offset,
			})
			if err != nil {
				writeError(w, model.NewInternalError("failed to list tags: "+err.Error()))
				return
			}
			cnt, err := h.queries.CountTagsByKey(rctx, pattern)
			if err != nil {
				writeError(w, model.NewInternalError("failed to count tags: "+err.Error()))
				return
			}
			total = cnt
			for _, row := range rawRows {
				rows = append(rows, tagWithCount{ID: row.ID, Key: row.Key, Value: row.Value, TransactionCount: row.TransactionCount})
			}
		} else {
			rawRows, err := h.queries.ListTags(rctx, generated.ListTagsParams{Limit: limit, Offset: offset})
			if err != nil {
				writeError(w, model.NewInternalError("failed to list tags: "+err.Error()))
				return
			}
			cnt, err := h.queries.CountTags(rctx)
			if err != nil {
				writeError(w, model.NewInternalError("failed to count tags: "+err.Error()))
				return
			}
			total = cnt
			for _, row := range rawRows {
				rows = append(rows, tagWithCount{ID: row.ID, Key: row.Key, Value: row.Value, TransactionCount: row.TransactionCount})
			}
		}

		if rows == nil {
			rows = []tagWithCount{}
		}
		totalPages := int(total) / perPage
		if int(total)%perPage != 0 {
			totalPages++
		}
		if totalPages == 0 {
			totalPages = 1
		}
		resp := paginatedResponse[tagWithCount]{
			Data: rows,
			Meta: paginatedMeta{
				Page:       page,
				PerPage:    perPage,
				Total:      int(total),
				TotalPages: totalPages,
			},
		}
		writeOK(w, resp)
	}
}
