package service

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/crispy-2/backend/internal/dal/generated"
	"github.com/crispy-2/backend/internal/model"
	"github.com/shopspring/decimal"
)

// PostingInput represents input data for a posting.
// If AccountID is 0 and AccountName is set, the account will be created automatically.
type PostingInput struct {
	AccountID   int64  `json:"account_id"`
	AccountName string `json:"account_name,omitempty"` // used when AccountID is 0 to create a new account
	Currency    string `json:"currency"`
	Amount      string `json:"amount"` // may be empty for auto-infer
}

// TagInput represents input data for a tag.
type TagInput struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

// TransactionService handles transaction business logic.
type TransactionService struct {
	db              *sql.DB
	queries         *generated.Queries
	exchangeRateSvc *ExchangeRateService
	accountSvc      *AccountService // injected after construction to resolve account names
}

// NewTransactionService creates a new TransactionService.
func NewTransactionService(db *sql.DB, exchangeRateSvc *ExchangeRateService) *TransactionService {
	return &TransactionService{
		db:              db,
		queries:         generated.New(db),
		exchangeRateSvc: exchangeRateSvc,
	}
}

// SetAccountService injects the AccountService for on-the-fly account name resolution.
func (s *TransactionService) SetAccountService(accountSvc *AccountService) {
	s.accountSvc = accountSvc
}

// resolvePostingAccountID resolves a posting's account ID.
// If AccountID > 0 it is used directly.
// If AccountID is 0 and AccountName is set, it looks up the account by full_name.
// If not found, it creates the account — including the full hierarchy for names like
// "Assets:Checking:Savings" (creates Assets, then Assets:Checking, then the leaf).
func (s *TransactionService) resolvePostingAccountID(p PostingInput) (int64, error) {
	if p.AccountID > 0 {
		return p.AccountID, nil
	}
	if p.AccountName == "" {
		return 0, model.NewValidationError("posting must have an account_id or account_name")
	}
	if s.accountSvc == nil {
		return 0, model.NewValidationError("account not found: " + p.AccountName)
	}

	// Search by exact full_name first — may already exist
	accounts, err := s.accountSvc.SearchAccounts(p.AccountName)
	if err == nil {
		for _, a := range accounts {
			if a.FullName == p.AccountName {
				return a.ID, nil
			}
		}
	}

	// Not found. Create the full hierarchy via AccountService.
	// Pass "expense" as default type for ad-hoc accounts created from transactions.
	acct, err := s.accountSvc.CreateAccountByPath(p.AccountName, "expense")
	if err != nil {
		return 0, err
	}
	return acct.ID, nil
}

// CreateTransaction creates a new transaction with postings and tags.
func (s *TransactionService) CreateTransaction(date, description, status string, postings []PostingInput, tags []TagInput) (*model.Transaction, error) {
	if err := validateTransaction(date, description, status, postings); err != nil {
		return nil, err
	}

	// Resolve account IDs before opening the DB transaction (account creation may open its own tx)
	resolvedAccountIDs := make([]int64, len(postings))
	for i, p := range postings {
		id, err := s.resolvePostingAccountID(p)
		if err != nil {
			return nil, err
		}
		resolvedAccountIDs[i] = id
	}

	// Parse and validate postings, handle auto-infer
	parsedAmounts, err := parseAndBalancePostings(postings)
	if err != nil {
		return nil, err
	}

    // Validate multi-currency transactions with exchange rates
    if len(postings) > 0 && s.exchangeRateSvc != nil {
        currencies := make(map[string]bool)
        for _, p := range postings {
            currencies[p.Currency] = true
        }
        if len(currencies) > 1 {
            baseCurrency := postings[0].Currency
            sum := parsedAmounts[0]
            for i := 1; i < len(postings); i++ {
                if postings[i].Currency != baseCurrency {
                    rate, err := s.exchangeRateSvc.GetRate(baseCurrency, postings[i].Currency, date)
                    if err != nil {
                        return nil, err
                    }
                    convertedAmt := parsedAmounts[i].Div(rate.Rate.Decimal)
                    sum = sum.Add(convertedAmt)
                } else {
                    sum = sum.Add(parsedAmounts[i])
                }
            }
            if !sum.IsZero() {
                return nil, model.NewValidationError("postings do not balance after currency conversion: sum is " + sum.String())
            }
        }
    }

    tx, err := s.db.Begin()
    if err != nil {
        return nil, model.NewInternalError("failed to begin transaction: " + err.Error())
    }
    defer tx.Rollback()

    q := generated.New(tx)

    // Create the transaction row
    row, err := q.CreateTransaction(ctx(), generated.CreateTransactionParams{
        Date:        date,
        Description: description,
        Status:      status,
    })
    if err != nil {
        return nil, model.NewInternalError("failed to create transaction: " + err.Error())
    }

	// Create postings
	modelPostings := make([]model.Posting, len(postings))
	for i, p := range postings {
		accountID := resolvedAccountIDs[i]
		pRow, err := q.CreatePosting(ctx(), generated.CreatePostingParams{
			TransactionID: row.ID,
			AccountID:     accountID,
			Currency:      p.Currency,
			Amount:        parsedAmounts[i].String(),
		})
		if err != nil {
			return nil, model.NewInternalError("failed to create posting: " + err.Error())
		}
		modelPostings[i] = model.Posting{
			ID:            pRow.ID,
			TransactionID: row.ID,
			AccountID:     accountID,
			Currency:      p.Currency,
			Amount:        model.StringDecimal{Decimal: parsedAmounts[i]},
		}
	}

    // Create/link tags
    var modelTags []model.Tag
    for _, t := range tags {
        tagRow, err := q.CreateTag(ctx(), generated.CreateTagParams{
            Key:   t.Key,
            Value: t.Value,
        })
        if err != nil {
            return nil, model.NewInternalError("failed to create tag: " + err.Error())
        }
        if err := q.AddTransactionTag(ctx(), generated.AddTransactionTagParams{
            TransactionID: row.ID,
            TagID:         tagRow.ID,
        }); err != nil {
            return nil, model.NewInternalError("failed to link tag: " + err.Error())
        }
        modelTags = append(modelTags, model.Tag{
            ID:    tagRow.ID,
            Key:   t.Key,
            Value: t.Value,
        })
    }

    if err := tx.Commit(); err != nil {
        return nil, model.NewInternalError("failed to commit transaction: " + err.Error())
    }

	slog.Info("transaction created", "id", row.ID, "date", date)
	return &model.Transaction{
		ID:          row.ID,
		Date:        date,
		Description: description,
		Status:      status,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		Postings:    modelPostings,
		Tags:        modelTags,
	}, nil
}

// GetTransaction retrieves a transaction with its postings and tags.
func (s *TransactionService) GetTransaction(id int64) (*model.Transaction, error) {
    row, err := s.queries.GetTransaction(ctx(), id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, model.NewNotFoundError("transaction", id)
        }
        return nil, model.NewInternalError("failed to get transaction: " + err.Error())
    }

    postingRows, err := s.queries.GetPostingsByTransactionWithAccounts(ctx(), id)
    if err != nil {
        return nil, model.NewInternalError("failed to get postings: " + err.Error())
    }

    tagRows, err := s.queries.GetTagsByTransaction(ctx(), id)
    if err != nil {
        return nil, model.NewInternalError("failed to get tags: " + err.Error())
    }

    postings := make([]model.Posting, len(postingRows))
    for i, p := range postingRows {
        amt, _ := decimal.NewFromString(p.Amount)
        postings[i] = model.Posting{
            ID:              p.ID,
            TransactionID:   p.TransactionID,
            AccountID:       p.AccountID,
            AccountFullName: p.AccountFullName,
            Currency:        p.Currency,
            Amount:          model.StringDecimal{Decimal: amt},
        }
    }

    tags := make([]model.Tag, len(tagRows))
    for i, t := range tagRows {
        tags[i] = model.Tag{ID: t.ID, Key: t.Key, Value: t.Value}
    }

    var recurringID *int64
    if row.RecurringTemplateID.Valid {
        recurringID = &row.RecurringTemplateID.Int64
    }

	return &model.Transaction{
		ID:                  row.ID,
		Date:                row.Date,
		Description:         row.Description,
		Status:              row.Status,
		RecurringTemplateID: recurringID,
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           row.UpdatedAt,
		Postings:            postings,
		Tags:                tags,
	}, nil
}

// ListTransactions returns a paginated list of transactions.
func (s *TransactionService) ListTransactions(filters TransactionFilters, pg model.Pagination) (*model.PaginatedResult[model.Transaction], error) {
    pg = normalizePagination(pg)
    offset := int64((pg.Page - 1) * pg.PageSize)
    limit := int64(pg.PageSize)

    var rows []generated.Transaction
    var total int64
    var err error

    switch {
    case filters.AccountID != nil:
        rows, err = s.queries.ListTransactionsByAccount(ctx(), generated.ListTransactionsByAccountParams{
            AccountID: *filters.AccountID,
            Limit:     limit,
            Offset:    offset,
        })
        if err == nil {
            total, err = s.queries.CountTransactionsByAccount(ctx(), *filters.AccountID)
        }
    case filters.Description != "":
        pattern := "%" + filters.Description + "%"
        rows, err = s.queries.ListTransactionsByDescription(ctx(), generated.ListTransactionsByDescriptionParams{
            Description: pattern,
            Limit:       limit,
            Offset:      offset,
        })
        if err == nil {
            total, err = s.queries.CountTransactionsByDescription(ctx(), pattern)
        }
    case filters.Status != "":
        rows, err = s.queries.ListTransactionsByStatus(ctx(), generated.ListTransactionsByStatusParams{
            Status: filters.Status,
            Limit:  limit,
            Offset: offset,
        })
        if err == nil {
            total, err = s.queries.CountTransactionsByStatus(ctx(), filters.Status)
        }
    case filters.DateFrom != "" && filters.DateTo != "":
        rows, err = s.queries.ListTransactionsByDateRange(ctx(), generated.ListTransactionsByDateRangeParams{
            Date:   filters.DateFrom,
            Date_2: filters.DateTo,
            Limit:  limit,
            Offset: offset,
        })
        if err == nil {
            total, err = s.queries.CountTransactionsByDateRange(ctx(), generated.CountTransactionsByDateRangeParams{
                Date:   filters.DateFrom,
                Date_2: filters.DateTo,
            })
        }
    case filters.TagKey != "" && filters.TagValue != "":
        rows, err = s.queries.ListTransactionsByTagKeyValue(ctx(), generated.ListTransactionsByTagKeyValueParams{
            Key:    filters.TagKey,
            Value:  filters.TagValue,
            Limit:  limit,
            Offset: offset,
        })
        if err == nil {
            total, err = s.queries.CountTransactionsByTagKeyValue(ctx(), generated.CountTransactionsByTagKeyValueParams{
                Key:   filters.TagKey,
                Value: filters.TagValue,
            })
        }
    case filters.TagKey != "":
        rows, err = s.queries.ListTransactionsByTagKey(ctx(), generated.ListTransactionsByTagKeyParams{
            Key:    filters.TagKey,
            Limit:  limit,
            Offset: offset,
        })
        if err == nil {
            total, err = s.queries.CountTransactionsByTagKey(ctx(), filters.TagKey)
        }
    default:
        rows, err = s.queries.ListTransactions(ctx(), generated.ListTransactionsParams{
            Limit:  limit,
            Offset: offset,
        })
        if err == nil {
            total, err = s.queries.CountTransactions(ctx())
        }
    }

    if err != nil {
        return nil, model.NewInternalError("failed to list transactions: " + err.Error())
    }

    items := make([]model.Transaction, len(rows))
    for i, r := range rows {
        // Load postings and tags for each transaction
        txn, err := s.GetTransaction(r.ID)
        if err != nil {
            return nil, err
        }
        items[i] = *txn
    }

    totalPages := int(total) / pg.PageSize
    if int(total)%pg.PageSize != 0 {
        totalPages++
    }

    return &model.PaginatedResult[model.Transaction]{
        Items:      items,
        Total:      int(total),
        Page:       pg.Page,
        PageSize:   pg.PageSize,
        TotalPages: totalPages,
    }, nil
}

// UpdateTransaction updates a transaction and its postings/tags.
func (s *TransactionService) UpdateTransaction(id int64, date, description, status string, postings []PostingInput, tags []TagInput) (*model.Transaction, error) {
	// Verify it exists
	_, err := s.queries.GetTransaction(ctx(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.NewNotFoundError("transaction", id)
		}
		return nil, model.NewInternalError("failed to get transaction: " + err.Error())
	}

	if err := validateTransaction(date, description, status, postings); err != nil {
		return nil, err
	}

	// Resolve account IDs before opening the DB transaction
	resolvedAccountIDsUpdate := make([]int64, len(postings))
	for i, p := range postings {
		accountID, err := s.resolvePostingAccountID(p)
		if err != nil {
			return nil, err
		}
		resolvedAccountIDsUpdate[i] = accountID
	}

	parsedAmounts, err := parseAndBalancePostings(postings)
    if err != nil {
        return nil, err
    }

    // Validate multi-currency transactions with exchange rates
    if len(postings) > 0 && s.exchangeRateSvc != nil {
        currencies := make(map[string]bool)
        for _, p := range postings {
            currencies[p.Currency] = true
        }
        if len(currencies) > 1 {
            baseCurrency := postings[0].Currency
            sum := parsedAmounts[0]
            for i := 1; i < len(postings); i++ {
                if postings[i].Currency != baseCurrency {
                    rate, err := s.exchangeRateSvc.GetRate(baseCurrency, postings[i].Currency, date)
                    if err != nil {
                        return nil, err
                    }
                    convertedAmt := parsedAmounts[i].Div(rate.Rate.Decimal)
                    sum = sum.Add(convertedAmt)
                } else {
                    sum = sum.Add(parsedAmounts[i])
                }
            }
            if !sum.IsZero() {
                return nil, model.NewValidationError("postings do not balance after currency conversion: sum is " + sum.String())
            }
        }
    }

    dbTx, err := s.db.Begin()
    if err != nil {
        return nil, model.NewInternalError("failed to begin db transaction: " + err.Error())
    }
    defer dbTx.Rollback()

    q := generated.New(dbTx)

    // Update transaction row
    row, err := q.UpdateTransaction(ctx(), generated.UpdateTransactionParams{
        Date:        date,
        Description: description,
        Status:      status,
        UpdatedAt:   time.Now().Format(time.RFC3339),
        ID:          id,
    })
    if err != nil {
        return nil, model.NewInternalError("failed to update transaction: " + err.Error())
    }

	// Replace postings
	if err := q.DeletePostingsByTransaction(ctx(), id); err != nil {
		return nil, model.NewInternalError("failed to delete old postings: " + err.Error())
	}

	modelPostings := make([]model.Posting, len(postings))
	for i, p := range postings {
		accountID := resolvedAccountIDsUpdate[i]
		pRow, err := q.CreatePosting(ctx(), generated.CreatePostingParams{
			TransactionID: id,
			AccountID:     accountID,
			Currency:      p.Currency,
			Amount:        parsedAmounts[i].String(),
		})
		if err != nil {
			return nil, model.NewInternalError("failed to create posting: " + err.Error())
		}
		modelPostings[i] = model.Posting{
			ID:            pRow.ID,
			TransactionID: id,
			AccountID:     accountID,
			Currency:      p.Currency,
			Amount:        model.StringDecimal{Decimal: parsedAmounts[i]},
		}
	}

    // Replace tags
    if err := q.RemoveTransactionTags(ctx(), id); err != nil {
        return nil, model.NewInternalError("failed to remove old tags: " + err.Error())
    }

    var modelTags []model.Tag
    for _, t := range tags {
        tagRow, err := q.CreateTag(ctx(), generated.CreateTagParams{
            Key:   t.Key,
            Value: t.Value,
        })
        if err != nil {
            return nil, model.NewInternalError("failed to create tag: " + err.Error())
        }
        if err := q.AddTransactionTag(ctx(), generated.AddTransactionTagParams{
            TransactionID: id,
            TagID:         tagRow.ID,
        }); err != nil {
            return nil, model.NewInternalError("failed to link tag: " + err.Error())
        }
        modelTags = append(modelTags, model.Tag{
            ID:    tagRow.ID,
            Key:   t.Key,
            Value: t.Value,
        })
    }

    if err := dbTx.Commit(); err != nil {
        return nil, model.NewInternalError("failed to commit: " + err.Error())
    }

	slog.Info("transaction updated", "id", id)
	return &model.Transaction{
		ID:          row.ID,
		Date:        date,
		Description: description,
		Status:      status,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		Postings:    modelPostings,
		Tags:        modelTags,
	}, nil
}

// DeleteTransaction deletes a transaction and its postings/tags.
func (s *TransactionService) DeleteTransaction(id int64) error {
    _, err := s.queries.GetTransaction(ctx(), id)
    if err != nil {
        if err == sql.ErrNoRows {
            return model.NewNotFoundError("transaction", id)
        }
        return model.NewInternalError("failed to get transaction: " + err.Error())
    }

    if err := s.queries.DeleteTransaction(ctx(), id); err != nil {
        return model.NewInternalError("failed to delete transaction: " + err.Error())
    }

    slog.Info("transaction deleted", "id", id)
    return nil
}

// TransactionFilters holds the optional filter parameters for listing transactions.
type TransactionFilters struct {
	Description string
	DateFrom    string
	DateTo      string
	AccountID   *int64
	TagKey      string
	TagValue    string
	Status      string
	Currency    string // filter by currency denomination
}

func validateTransaction(date, description, status string, postings []PostingInput) error {
    if date == "" {
        return model.NewValidationError("date is required")
    }
    if !model.ValidStatus(status) {
        return model.NewValidationError("invalid status: " + status)
    }
    if len(postings) < 2 {
        return model.NewValidationError("a transaction must have at least two postings")
    }
    return nil
}

func parseAndBalancePostings(postings []PostingInput) ([]decimal.Decimal, error) {
    amounts := make([]decimal.Decimal, len(postings))
    currencies := make(map[string]bool)
    emptyIndex := -1

    for i, p := range postings {
        currencies[p.Currency] = true
        if p.Amount == "" {
            if emptyIndex >= 0 {
                return nil, model.NewValidationError("only one posting amount may be left blank")
            }
            emptyIndex = i
        } else {
            amt, err := decimal.NewFromString(p.Amount)
            if err != nil {
                return nil, model.NewValidationError("invalid amount: " + p.Amount)
            }
            amounts[i] = amt
        }
    }

    isMultiCurrency := len(currencies) > 1

    if emptyIndex >= 0 {
        if isMultiCurrency {
            return nil, model.NewValidationError("multi-currency transactions require all posting amounts to be specified")
        }
        // Auto-infer the blank posting
        sum := decimal.Zero
        for i, amt := range amounts {
            if i != emptyIndex {
                sum = sum.Add(amt)
            }
        }
        amounts[emptyIndex] = sum.Neg()
    }

    // Check balance (single-currency case: sum must be zero)
    if !isMultiCurrency {
        sum := decimal.Zero
        for _, amt := range amounts {
            sum = sum.Add(amt)
        }
        if !sum.IsZero() {
            return nil, model.NewValidationError("postings do not balance: sum is " + sum.String())
        }
    }

    return amounts, nil
}
