package service

import (
    "database/sql"
    "log/slog"

    "github.com/crispy-2/backend/internal/dal/generated"
    "github.com/crispy-2/backend/internal/model"
)

// SearchResult holds the results of a global search.
type SearchResult struct {
    Accounts     []model.Account     `json:"accounts"`
    Transactions []model.Transaction `json:"transactions"`
    Tags         []model.Tag         `json:"tags"`
}

// SearchService handles global search across entities.
type SearchService struct {
    db      *sql.DB
    queries *generated.Queries
}

// NewSearchService creates a new SearchService.
func NewSearchService(db *sql.DB) *SearchService {
    return &SearchService{
        db:      db,
        queries: generated.New(db),
    }
}

// GlobalSearch searches across accounts, transactions, and tags.
func (s *SearchService) GlobalSearch(query string) (*SearchResult, error) {
    pattern := "%" + query + "%"

    // Search accounts
    accountRows, err := s.queries.SearchAccounts(ctx(), pattern)
    if err != nil {
        return nil, model.NewInternalError("failed to search accounts: " + err.Error())
    }
    accounts := make([]model.Account, len(accountRows))
    for i, r := range accountRows {
        accounts[i] = *toModelAccount(r)
    }

    // Search transactions by description
    txRows, err := s.queries.ListTransactionsByDescription(ctx(), generated.ListTransactionsByDescriptionParams{
        Description: pattern,
        Limit:       20,
        Offset:      0,
    })
    if err != nil {
        return nil, model.NewInternalError("failed to search transactions: " + err.Error())
    }
    transactions := make([]model.Transaction, len(txRows))
    for i, r := range txRows {
        transactions[i] = model.Transaction{
            ID:          r.ID,
            Date:        r.Date,
            Description: r.Description,
            Status:      r.Status,
        }
    }

    // Search tags
    tagRows, err := s.queries.SearchTags(ctx(), generated.SearchTagsParams{
        Key:   pattern,
        Value: pattern,
    })
    if err != nil {
        return nil, model.NewInternalError("failed to search tags: " + err.Error())
    }
    tags := make([]model.Tag, len(tagRows))
    for i, r := range tagRows {
        tags[i] = model.Tag{ID: r.ID, Key: r.Key, Value: r.Value}
    }

    slog.Debug("search completed", "query", query, "accounts_found", len(accounts), "transactions_found", len(transactions), "tags_found", len(tags))
    
    return &SearchResult{
        Accounts:     accounts,
        Transactions: transactions,
        Tags:         tags,
    }, nil
}
