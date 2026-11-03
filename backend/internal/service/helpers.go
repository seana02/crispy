package service

import (
    "context"

    "github.com/crispy-2/backend/internal/model"
)

// ctx returns a background context. In a real request context, this would
// come from the HTTP request or Wails call.
func ctx() context.Context {
    return context.Background()
}

// normalizePagination sets defaults for pagination parameters.
func normalizePagination(pg model.Pagination) model.Pagination {
    if pg.Page < 1 {
        pg.Page = 1
    }
    if pg.PageSize < 1 {
        pg.PageSize = 20
    }
    if pg.PageSize > 100 {
        pg.PageSize = 100
    }
    return pg
}
