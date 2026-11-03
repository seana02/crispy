package service

import (
    "database/sql"
    "fmt"
    "log/slog"
    "time"

    "github.com/crispy-2/backend/internal/dal/generated"
    "github.com/crispy-2/backend/internal/model"
    "github.com/shopspring/decimal"
)

// BudgetService handles budget business logic.
type BudgetService struct {
    db              *sql.DB
    queries         *generated.Queries
    exchangeRateSvc *ExchangeRateService
}

// NewBudgetService creates a new BudgetService.
func NewBudgetService(db *sql.DB, exchangeRateSvc *ExchangeRateService) *BudgetService {
    return &BudgetService{
        db:              db,
        queries:         generated.New(db),
        exchangeRateSvc: exchangeRateSvc,
    }
}

// CreateBudget creates a new budget.
func (s *BudgetService) CreateBudget(name, direction, period, currency, targetAmount string, accountIDs []int64, tagIDs []int64) (*model.Budget, error) {
    if err := validateBudget(name, direction, period, targetAmount, accountIDs); err != nil {
        return nil, err
    }

    tx, err := s.db.Begin()
    if err != nil {
        return nil, model.NewInternalError("failed to begin transaction: " + err.Error())
    }
    defer tx.Rollback()

    q := generated.New(tx)

    row, err := q.CreateBudget(ctx(), generated.CreateBudgetParams{
        Name:         name,
        Direction:    direction,
        Period:       period,
        Currency:     currency,
        TargetAmount: targetAmount,
    })
    if err != nil {
        return nil, model.NewInternalError("failed to create budget: " + err.Error())
    }

    for _, aid := range accountIDs {
        if err := q.AddBudgetAccount(ctx(), generated.AddBudgetAccountParams{
            BudgetID:  row.ID,
            AccountID: aid,
        }); err != nil {
            return nil, model.NewInternalError("failed to add budget account: " + err.Error())
        }
    }

    for _, tid := range tagIDs {
        if err := q.AddBudgetTag(ctx(), generated.AddBudgetTagParams{
            BudgetID: row.ID,
            TagID:    tid,
        }); err != nil {
            return nil, model.NewInternalError("failed to add budget tag: " + err.Error())
        }
    }

    if err := tx.Commit(); err != nil {
        return nil, model.NewInternalError("failed to commit: " + err.Error())
    }

	slog.Info("budget created", "id", row.ID, "name", name)
	return s.GetBudget(row.ID)
}

// GetBudget retrieves a budget by ID.
func (s *BudgetService) GetBudget(id int64) (*model.Budget, error) {
    row, err := s.queries.GetBudget(ctx(), id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, model.NewNotFoundError("budget", id)
        }
        return nil, model.NewInternalError("failed to get budget: " + err.Error())
    }

    accounts, err := s.queries.GetBudgetAccounts(ctx(), id)
    if err != nil {
        return nil, model.NewInternalError("failed to get budget accounts: " + err.Error())
    }

    tags, err := s.queries.GetBudgetTags(ctx(), id)
    if err != nil {
        return nil, model.NewInternalError("failed to get budget tags: " + err.Error())
    }

    accountIDs := make([]int64, len(accounts))
    accountNames := make([]string, len(accounts))
    for i, a := range accounts {
        accountIDs[i] = a.ID
        accountNames[i] = a.FullName
    }

    tagIDs := make([]int64, len(tags))
    for i, t := range tags {
        tagIDs[i] = t.ID
    }

	amt, _ := decimal.NewFromString(row.TargetAmount)
	return &model.Budget{
		ID:           row.ID,
		Name:         row.Name,
		Direction:    row.Direction,
		Period:       row.Period,
		Currency:     row.Currency,
		TargetAmount: model.StringDecimal{Decimal: amt},
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
		AccountIDs:   accountIDs,
		AccountNames: accountNames,
		TagIDs:       tagIDs,
	}, nil
}

// ListBudgets returns all budgets.
func (s *BudgetService) ListBudgets() ([]model.Budget, error) {
    rows, err := s.queries.ListBudgets(ctx())
    if err != nil {
        return nil, model.NewInternalError("failed to list budgets: " + err.Error())
    }

    budgets := make([]model.Budget, len(rows))
    for i, r := range rows {
        b, err := s.GetBudget(r.ID)
        if err != nil {
            return nil, err
        }
        budgets[i] = *b
    }
    return budgets, nil
}

// UpdateBudget updates a budget.
func (s *BudgetService) UpdateBudget(id int64, name, direction, period, currency, targetAmount string, accountIDs []int64, tagIDs []int64) (*model.Budget, error) {
    if err := validateBudget(name, direction, period, targetAmount, accountIDs); err != nil {
        return nil, err
    }

    _, err := s.queries.GetBudget(ctx(), id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, model.NewNotFoundError("budget", id)
        }
        return nil, model.NewInternalError("failed to get budget: " + err.Error())
    }

    tx, err := s.db.Begin()
    if err != nil {
        return nil, model.NewInternalError("failed to begin transaction: " + err.Error())
    }
    defer tx.Rollback()

    q := generated.New(tx)

    _, err = q.UpdateBudget(ctx(), generated.UpdateBudgetParams{
        Name:         name,
        Direction:    direction,
        Period:       period,
        Currency:     currency,
        TargetAmount: targetAmount,
        UpdatedAt:    time.Now().Format(time.RFC3339),
        ID:           id,
    })
    if err != nil {
        return nil, model.NewInternalError("failed to update budget: " + err.Error())
    }

    // Replace accounts
    if err := q.RemoveBudgetAccounts(ctx(), id); err != nil {
        return nil, model.NewInternalError("failed to remove budget accounts: " + err.Error())
    }
    for _, aid := range accountIDs {
        if err := q.AddBudgetAccount(ctx(), generated.AddBudgetAccountParams{
            BudgetID:  id,
            AccountID: aid,
        }); err != nil {
            return nil, model.NewInternalError("failed to add budget account: " + err.Error())
        }
    }

    // Replace tags
    if err := q.RemoveBudgetTags(ctx(), id); err != nil {
        return nil, model.NewInternalError("failed to remove budget tags: " + err.Error())
    }
    for _, tid := range tagIDs {
        if err := q.AddBudgetTag(ctx(), generated.AddBudgetTagParams{
            BudgetID: id,
            TagID:    tid,
        }); err != nil {
            return nil, model.NewInternalError("failed to add budget tag: " + err.Error())
        }
    }

    if err := tx.Commit(); err != nil {
        return nil, model.NewInternalError("failed to commit: " + err.Error())
    }

	slog.Info("budget updated", "id", id)
	return s.GetBudget(id)
}

// DeleteBudget deletes a budget.
func (s *BudgetService) DeleteBudget(id int64) error {
    _, err := s.queries.GetBudget(ctx(), id)
    if err != nil {
        if err == sql.ErrNoRows {
            return model.NewNotFoundError("budget", id)
        }
        return model.NewInternalError("failed to get budget: " + err.Error())
    }

    if err := s.queries.DeleteBudget(ctx(), id); err != nil {
        return model.NewInternalError("failed to delete budget: " + err.Error())
    }

    slog.Info("budget deleted", "id", id)
    return nil
}

// ComputeBudgetCurrent computes the current spent/saved amount for a budget within a date range.
func (s *BudgetService) ComputeBudgetCurrent(budgetID int64, dateFrom, dateTo string) (model.StringDecimal, error) {
    budget, err := s.GetBudget(budgetID)
    if err != nil {
        return model.StringDecimal{}, err
    }

    // Get all accounts and their descendants
    accounts, err := s.queries.GetBudgetAccounts(ctx(), budgetID)
    if err != nil {
        return model.StringDecimal{}, model.NewInternalError("failed to get budget accounts: " + err.Error())
    }

    // Collect all account IDs (direct + descendants)
    var allAccountIDs []int64
    for _, a := range accounts {
        allAccountIDs = append(allAccountIDs, a.ID)
        // Get descendants
        descendants, err := s.queries.GetDescendantAccounts(ctx(), a.FullName+":%" )
        if err != nil {
            return model.StringDecimal{}, model.NewInternalError("failed to get descendants: " + err.Error())
        }
        for _, d := range descendants {
            allAccountIDs = append(allAccountIDs, d.ID)
        }
    }

    // Get budget tag filters
    budgetTags, err := s.queries.GetBudgetTags(ctx(), budgetID)
    if err != nil {
        return model.StringDecimal{}, model.NewInternalError("failed to get budget tags: " + err.Error())
    }

    // Sum postings matching account criteria within the date range
    // We need to query postings directly since we have complex criteria
    total := decimal.Zero
    for _, accountID := range allAccountIDs {
        postings, err := s.getPostingsForAccountInRange(accountID, dateFrom, dateTo)
        if err != nil {
            return model.StringDecimal{}, err
        }
        for _, p := range postings {
            // If there are tag filters, check that the transaction has at least one matching tag
            if len(budgetTags) > 0 {
                hasMatchingTag, err := s.transactionHasAnyTag(p.transactionID, budgetTags)
                if err != nil {
                    return model.StringDecimal{}, err
                }
                if !hasMatchingTag {
                    continue
                }
            }

            amt, _ := decimal.NewFromString(p.amount)
            if p.currency == budget.Currency {
                total = total.Add(amt)
            } else if s.exchangeRateSvc != nil {
                rate, err := s.exchangeRateSvc.GetRate(p.currency, budget.Currency, p.date)
                if err != nil {
                    slog.Warn("skipping posting due to missing exchange rate", "currency", p.currency, "target", budget.Currency)
                    continue
                }
                total = total.Add(amt.Mul(rate.Rate.Decimal))
            }
        }
    }

    return model.StringDecimal{Decimal: total}, nil
}

type postingWithContext struct {
    transactionID int64
    amount        string
    currency      string
    date          string
}

func (s *BudgetService) getPostingsForAccountInRange(accountID int64, dateFrom, dateTo string) ([]postingWithContext, error) {
    query := `SELECT p.transaction_id, p.amount, p.currency, t.date
              FROM postings p
              JOIN transactions t ON t.id = p.transaction_id
              WHERE p.account_id = ? AND t.date >= ? AND t.date <= ?`

    rows, err := s.db.Query(query, accountID, dateFrom, dateTo)
    if err != nil {
        return nil, model.NewInternalError(fmt.Sprintf("failed to query postings for account %d: %v", accountID, err))
    }
    defer rows.Close()

    var results []postingWithContext
    for rows.Next() {
        var p postingWithContext
        if err := rows.Scan(&p.transactionID, &p.amount, &p.currency, &p.date); err != nil {
            return nil, model.NewInternalError("failed to scan posting: " + err.Error())
        }
        results = append(results, p)
    }
    return results, rows.Err()
}

func (s *BudgetService) transactionHasAnyTag(transactionID int64, tags []generated.Tag) (bool, error) {
    txTags, err := s.queries.GetTagsByTransaction(ctx(), transactionID)
    if err != nil {
        return false, model.NewInternalError("failed to get transaction tags: " + err.Error())
    }
    for _, tt := range txTags {
        for _, bt := range tags {
            if tt.ID == bt.ID {
                return true, nil
            }
        }
    }
    return false, nil
}

func validateBudget(name, direction, period, targetAmount string, accountIDs []int64) error {
    if name == "" {
        return model.NewValidationError("budget name is required")
    }
    if direction != model.DirectionSpending && direction != model.DirectionSaving {
        return model.NewValidationError("direction must be 'spending' or 'saving'")
    }
    if period != model.PeriodWeekly && period != model.PeriodMonthly && period != model.PeriodYearly {
        return model.NewValidationError("period must be 'weekly', 'monthly', or 'yearly'")
    }
    if _, err := decimal.NewFromString(targetAmount); err != nil {
        return model.NewValidationError("invalid target_amount: " + targetAmount)
    }
    if len(accountIDs) == 0 {
        return model.NewValidationError("at least one account must be specified")
    }
    return nil
}
