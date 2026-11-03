package service

import (
    "database/sql"
    "log/slog"
    "strings"

    "github.com/crispy-2/backend/internal/dal/generated"
    "github.com/crispy-2/backend/internal/model"
    "github.com/shopspring/decimal"
)

// ReportService handles financial reports.
type ReportService struct {
    db              *sql.DB
    queries         *generated.Queries
    exchangeRateSvc *ExchangeRateService
}

// NewReportService creates a new ReportService.
func NewReportService(db *sql.DB, exchangeRateSvc *ExchangeRateService) *ReportService {
    return &ReportService{
        db:              db,
        queries:         generated.New(db),
        exchangeRateSvc: exchangeRateSvc,
    }
}

// BalanceSheetRow represents a row in a balance sheet report.
type BalanceSheetRow struct {
    AccountID   int64       `json:"account_id"`
    AccountName string      `json:"account_name"`
    Balance     model.StringDecimal `json:"balance"`
    Currency    string      `json:"currency"`
    Depth       int         `json:"depth"`
}

// BalanceSheetReport is the complete balance sheet.
type BalanceSheetReport struct {
    AssetRows     []BalanceSheetRow `json:"asset_rows"`
    LiabilityRows []BalanceSheetRow `json:"liability_rows"`
    TotalAssets   model.StringDecimal `json:"total_assets"`
    TotalLiab     model.StringDecimal `json:"total_liabilities"`
}

// IncomeStatementRow represents a row in an income statement.
type IncomeStatementRow struct {
    AccountID   int64       `json:"account_id"`
    AccountName string      `json:"account_name"`
    Amount      model.StringDecimal `json:"amount"`
    Currency    string      `json:"currency"`
    Depth       int         `json:"depth"`
}

// IncomeStatementReport is the complete income statement.
type IncomeStatementReport struct {
    RevenueRows  []IncomeStatementRow `json:"revenue_rows"`
    ExpenseRows  []IncomeStatementRow `json:"expense_rows"`
    TotalRevenue model.StringDecimal `json:"total_revenue"`
    TotalExpense model.StringDecimal `json:"total_expenses"`
    NetIncome    model.StringDecimal `json:"net_income"`
}

// SpendingCategory represents a top spending category.
type SpendingCategory struct {
    AccountID   int64       `json:"account_id"`
    AccountName string      `json:"account_name"`
    Amount      model.StringDecimal `json:"amount"`
    Currency    string      `json:"currency"`
}

// NetWorthWithTrend represents net worth with comparison to another date.
type NetWorthWithTrend struct {
    AsOfDate           string              `json:"as_of_date"`
    NetWorth           model.StringDecimal `json:"net_worth"`
    ComparisonDate     string              `json:"comparison_date,omitempty"`
    ComparisonNetWorth model.StringDecimal `json:"comparison_net_worth,omitempty"`
    Trend              model.StringDecimal `json:"trend,omitempty"`
    TrendPercentage    string              `json:"trend_percentage,omitempty"`
}

// GetBalanceSheet generates a balance sheet for the given date range.
func (s *ReportService) GetBalanceSheet(dateFrom, dateTo string, depth int, tagFilters []int64) (*BalanceSheetReport, error) {
    if len(tagFilters) > 0 {
        assetRows, totalAssets, err := s.getAccountBalancesWithTags("asset", dateFrom, dateTo, depth, tagFilters)
        if err != nil {
            return nil, err
        }
        liabilityRows, totalLiab, err := s.getAccountBalancesWithTags("liability", dateFrom, dateTo, depth, tagFilters)
        if err != nil {
            return nil, err
        }
        slog.Debug("balance sheet generated", "assets", totalAssets, "liabilities", totalLiab)
        return &BalanceSheetReport{
            AssetRows:     assetRows,
            LiabilityRows: liabilityRows,
            TotalAssets:   totalAssets,
            TotalLiab:     totalLiab,
        }, nil
    }

    assetRows, totalAssets, err := s.getAccountBalances("asset", dateFrom, dateTo, depth)
    if err != nil {
        return nil, err
    }

    liabilityRows, totalLiab, err := s.getAccountBalances("liability", dateFrom, dateTo, depth)
    if err != nil {
        return nil, err
    }

    slog.Debug("balance sheet generated", "assets", totalAssets, "liabilities", totalLiab)
    return &BalanceSheetReport{
        AssetRows:     assetRows,
        LiabilityRows: liabilityRows,
        TotalAssets:   totalAssets,
        TotalLiab:     totalLiab,
    }, nil
}

// GetIncomeStatement generates an income statement for the given date range.
func (s *ReportService) GetIncomeStatement(dateFrom, dateTo string, depth int, tagFilters []int64) (*IncomeStatementReport, error) {
    if len(tagFilters) > 0 {
        revenueRows, totalRevenue, err := s.getAccountAmountsWithTags("revenue", dateFrom, dateTo, depth, tagFilters)
        if err != nil {
            return nil, err
        }
        expenseRows, totalExpense, err := s.getAccountAmountsWithTags("expense", dateFrom, dateTo, depth, tagFilters)
        if err != nil {
            return nil, err
        }
        netIncome := model.StringDecimal{Decimal: totalRevenue.Decimal.Neg().Sub(totalExpense.Decimal)}
        slog.Debug("income statement generated", "revenue", totalRevenue, "expenses", totalExpense, "net", netIncome)
        return &IncomeStatementReport{
            RevenueRows:  revenueRows,
            ExpenseRows:  expenseRows,
            TotalRevenue: totalRevenue,
            TotalExpense: totalExpense,
            NetIncome:    netIncome,
        }, nil
    }

    revenueRows, totalRevenue, err := s.getAccountAmounts("revenue", dateFrom, dateTo, depth)
    if err != nil {
        return nil, err
    }

    expenseRows, totalExpense, err := s.getAccountAmounts("expense", dateFrom, dateTo, depth)
    if err != nil {
        return nil, err
    }

    netIncome := model.StringDecimal{Decimal: totalRevenue.Decimal.Neg().Sub(totalExpense.Decimal)}

    slog.Debug("income statement generated", "revenue", totalRevenue, "expenses", totalExpense, "net", netIncome)
    return &IncomeStatementReport{
        RevenueRows:  revenueRows,
        ExpenseRows:  expenseRows,
        TotalRevenue: totalRevenue,
        TotalExpense: totalExpense,
        NetIncome:    netIncome,
    }, nil
}

// GetNetWorth computes total assets minus total liabilities as of a given date.
// Optional comparisonDate parameter allows computing trend vs. a previous date.
func (s *ReportService) GetNetWorth(asOfDate string, comparisonDate ...string) (*NetWorthWithTrend, error) {
    totalAssets, err := s.sumAccountBalances("asset", asOfDate)
    if err != nil {
        return nil, err
    }

    totalLiabilities, err := s.sumAccountBalances("liability", asOfDate)
    if err != nil {
        return nil, err
    }

    netWorth := model.StringDecimal{Decimal: totalAssets.Decimal.Sub(totalLiabilities.Decimal)}

    result := &NetWorthWithTrend{
        AsOfDate: asOfDate,
        NetWorth: netWorth,
    }

    if len(comparisonDate) > 0 && comparisonDate[0] != "" {
        compAssets, err := s.sumAccountBalances("asset", comparisonDate[0])
        if err != nil {
            return nil, err
        }
        compLiabilities, err := s.sumAccountBalances("liability", comparisonDate[0])
        if err != nil {
            return nil, err
        }
        compNetWorth := model.StringDecimal{Decimal: compAssets.Decimal.Sub(compLiabilities.Decimal)}

        result.ComparisonDate = comparisonDate[0]
        result.ComparisonNetWorth = compNetWorth

        trend := model.StringDecimal{Decimal: netWorth.Decimal.Sub(compNetWorth.Decimal)}
        result.Trend = trend

        if !compNetWorth.Decimal.IsZero() {
            trendPct := trend.Decimal.Div(compNetWorth.Decimal).Mul(decimal.NewFromInt(100))
            result.TrendPercentage = trendPct.String()
        } else if !netWorth.Decimal.IsZero() {
            result.TrendPercentage = "100"
        } else {
            result.TrendPercentage = "0"
        }

        slog.Debug("net worth with trend computed", "as_of_date", asOfDate, "comparison_date", comparisonDate[0], "net_worth", netWorth, "trend", trend, "trend_pct", result.TrendPercentage)
    } else {
        slog.Debug("net worth computed", "assets", totalAssets, "liabilities", totalLiabilities, "net_worth", netWorth)
    }

    return result, nil
}

// GetCashFlow computes net income minus expenses for a date range.
func (s *ReportService) GetCashFlow(dateFrom, dateTo string) (model.StringDecimal, error) {
    revSum, err := s.sumPostingsByType("revenue", dateFrom, dateTo)
    if err != nil {
        return model.StringDecimal{}, err
    }

    expSum, err := s.sumPostingsByType("expense", dateFrom, dateTo)
    if err != nil {
        return model.StringDecimal{}, err
    }

    cashFlow := model.StringDecimal{Decimal: revSum.Decimal.Neg().Sub(expSum.Decimal)}
    return cashFlow, nil
}

// GetTopSpendingCategories returns the top N expense accounts by amount.
func (s *ReportService) GetTopSpendingCategories(dateFrom, dateTo string, limit int) ([]SpendingCategory, error) {
    query := `
        SELECT a.id, a.full_name, p.currency, p.amount
        FROM postings p
        JOIN transactions t ON t.id = p.transaction_id
        JOIN accounts a ON a.id = p.account_id
        WHERE t.date >= ? AND t.date <= ?
          AND p.amount > '0'
          AND a.id IN (
              SELECT id FROM accounts WHERE account_type = 'expense'
              UNION
              SELECT id FROM accounts WHERE full_name LIKE (
                  SELECT full_name || ':%' FROM accounts WHERE account_type = 'expense'
              )
          )
    `

    rows, err := s.db.Query(query, dateFrom, dateTo)
    if err != nil {
        return nil, model.NewInternalError("failed to query top spending: " + err.Error())
    }
    defer rows.Close()

    categoryTotals := make(map[string]*SpendingCategory)
    for rows.Next() {
        var accountID int64
        var accountName, currency, amount string
        if err := rows.Scan(&accountID, &accountName, &currency, &amount); err != nil {
            return nil, model.NewInternalError("failed to scan: " + err.Error())
        }

        key := accountName
        cat, exists := categoryTotals[key]
        if !exists {
            cat = &SpendingCategory{
                AccountID:   accountID,
                AccountName: accountName,
                Currency:    currency,
                Amount:      model.StringDecimal{Decimal: decimal.Zero},
            }
            categoryTotals[key] = cat
        }
        amt, err := decimal.NewFromString(amount)
        if err != nil {
            return nil, model.NewInternalError("failed to parse amount: " + err.Error())
        }
        cat.Amount = model.StringDecimal{Decimal: cat.Amount.Decimal.Add(amt)}
    }
    if err := rows.Err(); err != nil {
        return nil, model.NewInternalError("failed to iterate rows: " + err.Error())
    }

    var categories []SpendingCategory
    for _, cat := range categoryTotals {
        categories = append(categories, *cat)
    }

    for i := 0; i < len(categories)-1; i++ {
        for j := i + 1; j < len(categories); j++ {
            if categories[j].Amount.Decimal.GreaterThan(categories[i].Amount.Decimal) {
                categories[i], categories[j] = categories[j], categories[i]
            }
        }
    }

    if len(categories) > limit {
        categories = categories[:limit]
    }

    return categories, nil
}

// GetExpenseBreakdown returns spending breakdown grouped by account or tag.
func (s *ReportService) GetExpenseBreakdown(startDate, endDate string, groupBy string) (map[string]model.StringDecimal, error) {
    if groupBy != "account" && groupBy != "tag" {
        return nil, model.NewValidationError("group_by must be 'account' or 'tag'")
    }

    result := make(map[string]model.StringDecimal)

    if groupBy == "account" {
        query := `
            SELECT a.full_name, p.currency, SUM(p.amount) as amount
            FROM postings p
            JOIN transactions t ON t.id = p.transaction_id
            JOIN accounts a ON a.id = p.account_id
            WHERE t.date >= ? AND t.date <= ?
              AND p.amount > '0'
              AND (a.account_type = 'expense' OR a.full_name LIKE (
                  SELECT full_name || ':%' FROM accounts WHERE account_type = 'expense' LIMIT 1
              ))
            GROUP BY a.full_name, p.currency
            ORDER BY SUM(p.amount) DESC
        `

        rows, err := s.db.Query(query, startDate, endDate)
        if err != nil {
            return nil, model.NewInternalError("failed to query expense breakdown: " + err.Error())
        }
        defer rows.Close()

        for rows.Next() {
            var accountName, currency, amount string
            if err := rows.Scan(&accountName, &currency, &amount); err != nil {
                return nil, model.NewInternalError("failed to scan: " + err.Error())
            }
            key := accountName + " (" + currency + ")"
            amt, err := decimal.NewFromString(amount)
            if err != nil {
                return nil, model.NewInternalError("failed to parse amount: " + err.Error())
            }
            result[key] = model.StringDecimal{Decimal: amt}
        }
        if err := rows.Err(); err != nil {
            return nil, model.NewInternalError("failed to iterate rows: " + err.Error())
        }
    } else {
        query := `
            SELECT tg.value, p.currency, SUM(p.amount) as amount
            FROM postings p
            JOIN transactions t ON t.id = p.transaction_id
            JOIN transaction_tags tt ON tt.transaction_id = t.id
            JOIN tags tg ON tg.id = tt.tag_id
            JOIN accounts a ON a.id = p.account_id
            WHERE t.date >= ? AND t.date <= ?
              AND p.amount > '0'
              AND (a.account_type = 'expense' OR a.full_name LIKE (
                  SELECT full_name || ':%' FROM accounts WHERE account_type = 'expense' LIMIT 1
              ))
            GROUP BY tg.value, p.currency
            ORDER BY SUM(p.amount) DESC
        `

        rows, err := s.db.Query(query, startDate, endDate)
        if err != nil {
            return nil, model.NewInternalError("failed to query expense breakdown by tag: " + err.Error())
        }
        defer rows.Close()

        for rows.Next() {
            var tagValue, currency, amount string
            if err := rows.Scan(&tagValue, &currency, &amount); err != nil {
                return nil, model.NewInternalError("failed to scan: " + err.Error())
            }
            key := tagValue + " (" + currency + ")"
            amt, err := decimal.NewFromString(amount)
            if err != nil {
                return nil, model.NewInternalError("failed to parse amount: " + err.Error())
            }
            result[key] = model.StringDecimal{Decimal: amt}
        }
        if err := rows.Err(); err != nil {
            return nil, model.NewInternalError("failed to iterate rows: " + err.Error())
        }
    }

    slog.Debug("expense breakdown computed", "start_date", startDate, "end_date", endDate, "group_by", groupBy)
    return result, nil
}

// GetTagSpendingSummary returns spending grouped by values of a tag key.
func (s *ReportService) GetTagSpendingSummary(tagKey string, startDate, endDate string) (map[string]model.StringDecimal, error) {
    result := make(map[string]model.StringDecimal)

    query := `
        SELECT tg.value, p.currency, SUM(p.amount) as amount
        FROM postings p
        JOIN transactions t ON t.id = p.transaction_id
        JOIN transaction_tags tt ON tt.transaction_id = t.id
        JOIN tags tg ON tg.id = tt.tag_id
        JOIN accounts a ON a.id = p.account_id
        WHERE t.date >= ? AND t.date <= ?
          AND p.amount > '0'
          AND (a.account_type = 'expense' OR a.full_name LIKE (
              SELECT full_name || ':%' FROM accounts WHERE account_type = 'expense' LIMIT 1
          ))
          AND tg.key = ?
        GROUP BY tg.value, p.currency
        ORDER BY SUM(p.amount) DESC
    `

    rows, err := s.db.Query(query, startDate, endDate, tagKey)
    if err != nil {
        return nil, model.NewInternalError("failed to query tag spending summary: " + err.Error())
    }
    defer rows.Close()

    for rows.Next() {
        var tagValue, currency, amount string
        if err := rows.Scan(&tagValue, &currency, &amount); err != nil {
            return nil, model.NewInternalError("failed to scan: " + err.Error())
        }
        key := tagValue + " (" + currency + ")"
        amt, err := decimal.NewFromString(amount)
        if err != nil {
            return nil, model.NewInternalError("failed to parse amount: " + err.Error())
        }
        result[key] = model.StringDecimal{Decimal: amt}
    }
    if err := rows.Err(); err != nil {
        return nil, model.NewInternalError("failed to iterate rows: " + err.Error())
    }

    slog.Debug("tag spending summary computed", "tag_key", tagKey, "start_date", startDate, "end_date", endDate)
    return result, nil
}

func (s *ReportService) getAccountBalances(accountType, dateFrom, dateTo string, depth int) ([]BalanceSheetRow, model.StringDecimal, error) {
    query := `
        SELECT a.id, a.full_name, p.currency, p.amount
        FROM postings p
        JOIN transactions t ON t.id = p.transaction_id
        JOIN accounts a ON a.id = p.account_id
        WHERE t.date >= ? AND t.date <= ?
          AND (a.account_type = ? OR a.full_name LIKE (
              SELECT full_name || ':%' FROM accounts WHERE account_type = ? LIMIT 1
          ))
    `

    rows, err := s.db.Query(query, dateFrom, dateTo, accountType, accountType)
    if err != nil {
        return nil, model.StringDecimal{}, model.NewInternalError("failed to query balances: " + err.Error())
    }
    defer rows.Close()

    accountTotals := make(map[int64]struct {
        row     BalanceSheetRow
        total   decimal.Decimal
    })

    for rows.Next() {
        var accountID int64
        var accountName, currency, amount string
        if err := rows.Scan(&accountID, &accountName, &currency, &amount); err != nil {
            return nil, model.StringDecimal{}, model.NewInternalError("failed to scan: " + err.Error())
        }

        data, exists := accountTotals[accountID]
        if !exists {
            data.row = BalanceSheetRow{
                AccountID:   accountID,
                AccountName: accountName,
                Currency:    currency,
                Depth:       depth,
            }
            data.total = decimal.Zero
        }

        amt, err := decimal.NewFromString(amount)
        if err != nil {
            return nil, model.StringDecimal{}, model.NewInternalError("failed to parse amount: " + err.Error())
        }
        data.total = data.total.Add(amt)
        accountTotals[accountID] = data
    }
    if err := rows.Err(); err != nil {
        return nil, model.StringDecimal{}, model.NewInternalError("failed to iterate rows: " + err.Error())
    }

    var result []BalanceSheetRow
    var total decimal.Decimal

    for _, data := range accountTotals {
        data.row.Balance = model.StringDecimal{Decimal: data.total}
        result = append(result, data.row)
        total = total.Add(data.total)
    }

    return result, model.StringDecimal{Decimal: total}, nil
}

func (s *ReportService) getAccountAmounts(accountType, dateFrom, dateTo string, depth int) ([]IncomeStatementRow, model.StringDecimal, error) {
    query := `
        SELECT a.id, a.full_name, p.currency, p.amount
        FROM postings p
        JOIN transactions t ON t.id = p.transaction_id
        JOIN accounts a ON a.id = p.account_id
        WHERE t.date >= ? AND t.date <= ?
          AND (a.account_type = ? OR a.full_name LIKE (
              SELECT full_name || ':%' FROM accounts WHERE account_type = ? LIMIT 1
          ))
    `

    rows, err := s.db.Query(query, dateFrom, dateTo, accountType, accountType)
    if err != nil {
        return nil, model.StringDecimal{}, model.NewInternalError("failed to query amounts: " + err.Error())
    }
    defer rows.Close()

    accountTotals := make(map[int64]struct {
        row     IncomeStatementRow
        total   decimal.Decimal
    })

    for rows.Next() {
        var accountID int64
        var accountName, currency, amount string
        if err := rows.Scan(&accountID, &accountName, &currency, &amount); err != nil {
            return nil, model.StringDecimal{}, model.NewInternalError("failed to scan: " + err.Error())
        }

        data, exists := accountTotals[accountID]
        if !exists {
            data.row = IncomeStatementRow{
                AccountID:   accountID,
                AccountName: accountName,
                Currency:    currency,
                Depth:       depth,
            }
            data.total = decimal.Zero
        }

        amt, err := decimal.NewFromString(amount)
        if err != nil {
            return nil, model.StringDecimal{}, model.NewInternalError("failed to parse amount: " + err.Error())
        }
        data.total = data.total.Add(amt)
        accountTotals[accountID] = data
    }
    if err := rows.Err(); err != nil {
        return nil, model.StringDecimal{}, model.NewInternalError("failed to iterate rows: " + err.Error())
    }

    var result []IncomeStatementRow
    var total decimal.Decimal

    for _, data := range accountTotals {
        data.row.Amount = model.StringDecimal{Decimal: data.total}
        result = append(result, data.row)
        total = total.Add(data.total)
    }

    return result, model.StringDecimal{Decimal: total}, nil
}

func (s *ReportService) getAccountBalancesWithTags(accountType, dateFrom, dateTo string, depth int, tagFilters []int64) (rows []BalanceSheetRow, total model.StringDecimal, err error) {
    if len(tagFilters) == 0 {
        return s.getAccountBalances(accountType, dateFrom, dateTo, depth)
    }

    placeholders := make([]string, len(tagFilters))
    for i := range tagFilters {
        placeholders[i] = "?"
    }

    query := `
        SELECT DISTINCT a.id, a.full_name, p.currency, p.amount
        FROM postings p
        JOIN transactions t ON t.id = p.transaction_id
        JOIN accounts a ON a.id = p.account_id
        JOIN transaction_tags tt ON tt.transaction_id = t.id
        WHERE t.date >= ? AND t.date <= ?
          AND (a.account_type = ? OR a.full_name LIKE (
              SELECT full_name || ':%' FROM accounts WHERE account_type = ? LIMIT 1
          ))
          AND tt.tag_id IN (` + strings.Join(placeholders, ",") + `)
    `

    args := []interface{}{dateFrom, dateTo, accountType, accountType}
    for _, tf := range tagFilters {
        args = append(args, tf)
    }

    rowsStmt, err := s.db.Query(query, args...)
    if err != nil {
        return nil, model.StringDecimal{}, model.NewInternalError("failed to query balances with tags: " + err.Error())
    }
    defer rowsStmt.Close()

    accountTotals := make(map[int64]struct {
        row     BalanceSheetRow
        total   decimal.Decimal
    })

    for rowsStmt.Next() {
        var accountID int64
        var accountName, currency, amount string
        if err := rowsStmt.Scan(&accountID, &accountName, &currency, &amount); err != nil {
            return nil, model.StringDecimal{}, model.NewInternalError("failed to scan: " + err.Error())
        }

        data, exists := accountTotals[accountID]
        if !exists {
            data.row = BalanceSheetRow{
                AccountID:   accountID,
                AccountName: accountName,
                Currency:    currency,
                Depth:       depth,
            }
            data.total = decimal.Zero
        }

        amt, err := decimal.NewFromString(amount)
        if err != nil {
            return nil, model.StringDecimal{}, model.NewInternalError("failed to parse amount: " + err.Error())
        }
        data.total = data.total.Add(amt)
        accountTotals[accountID] = data
    }
    if err := rowsStmt.Err(); err != nil {
        return nil, model.StringDecimal{}, model.NewInternalError("failed to iterate rows: " + err.Error())
    }

    var result []BalanceSheetRow
    var totalAmt decimal.Decimal

    for _, data := range accountTotals {
        data.row.Balance = model.StringDecimal{Decimal: data.total}
        result = append(result, data.row)
        totalAmt = totalAmt.Add(data.total)
    }

    return result, model.StringDecimal{Decimal: totalAmt}, nil
}

func (s *ReportService) getAccountAmountsWithTags(accountType, dateFrom, dateTo string, depth int, tagFilters []int64) (rows []IncomeStatementRow, total model.StringDecimal, err error) {
    if len(tagFilters) == 0 {
        return s.getAccountAmounts(accountType, dateFrom, dateTo, depth)
    }

    placeholders := make([]string, len(tagFilters))
    for i := range tagFilters {
        placeholders[i] = "?"
    }

    query := `
        SELECT DISTINCT a.id, a.full_name, p.currency, p.amount
        FROM postings p
        JOIN transactions t ON t.id = p.transaction_id
        JOIN accounts a ON a.id = p.account_id
        JOIN transaction_tags tt ON tt.transaction_id = t.id
        WHERE t.date >= ? AND t.date <= ?
          AND (a.account_type = ? OR a.full_name LIKE (
              SELECT full_name || ':%' FROM accounts WHERE account_type = ? LIMIT 1
          ))
          AND tt.tag_id IN (` + strings.Join(placeholders, ",") + `)
    `

    args := []interface{}{dateFrom, dateTo, accountType, accountType}
    for _, tf := range tagFilters {
        args = append(args, tf)
    }

    rowsStmt, err := s.db.Query(query, args...)
    if err != nil {
        return nil, model.StringDecimal{}, model.NewInternalError("failed to query amounts with tags: " + err.Error())
    }
    defer rowsStmt.Close()

    accountTotals := make(map[int64]struct {
        row     IncomeStatementRow
        total   decimal.Decimal
    })

    for rowsStmt.Next() {
        var accountID int64
        var accountName, currency, amount string
        if err := rowsStmt.Scan(&accountID, &accountName, &currency, &amount); err != nil {
            return nil, model.StringDecimal{}, model.NewInternalError("failed to scan: " + err.Error())
        }

        data, exists := accountTotals[accountID]
        if !exists {
            data.row = IncomeStatementRow{
                AccountID:   accountID,
                AccountName: accountName,
                Currency:    currency,
                Depth:       depth,
            }
            data.total = decimal.Zero
        }

        amt, err := decimal.NewFromString(amount)
        if err != nil {
            return nil, model.StringDecimal{}, model.NewInternalError("failed to parse amount: " + err.Error())
        }
        data.total = data.total.Add(amt)
        accountTotals[accountID] = data
    }
    if err := rowsStmt.Err(); err != nil {
        return nil, model.StringDecimal{}, model.NewInternalError("failed to iterate rows: " + err.Error())
    }

    var result []IncomeStatementRow
    var totalAmt decimal.Decimal

    for _, data := range accountTotals {
        data.row.Amount = model.StringDecimal{Decimal: data.total}
        result = append(result, data.row)
        totalAmt = totalAmt.Add(data.total)
    }

    return result, model.StringDecimal{Decimal: totalAmt}, nil
}

func (s *ReportService) sumAccountBalances(accountType, asOfDate string) (model.StringDecimal, error) {
    query := `
        SELECT p.amount
        FROM postings p
        JOIN transactions t ON t.id = p.transaction_id
        JOIN accounts a ON a.id = p.account_id
        WHERE t.date <= ?
          AND (a.account_type = ? OR a.full_name LIKE (
              SELECT full_name || ':%' FROM accounts WHERE account_type = ? LIMIT 1
          ))
    `

    rows, err := s.db.Query(query, asOfDate, accountType, accountType)
    if err != nil {
        return model.StringDecimal{}, model.NewInternalError("failed to query balances: " + err.Error())
    }
    defer rows.Close()

    var total decimal.Decimal
    for rows.Next() {
        var amount string
        if err := rows.Scan(&amount); err != nil {
            return model.StringDecimal{}, model.NewInternalError("failed to scan: " + err.Error())
        }
        amt, err := decimal.NewFromString(amount)
        if err != nil {
            return model.StringDecimal{}, model.NewInternalError("failed to parse amount: " + err.Error())
        }
        total = total.Add(amt)
    }
    if err := rows.Err(); err != nil {
        return model.StringDecimal{}, model.NewInternalError("failed to iterate rows: " + err.Error())
    }

    return model.StringDecimal{Decimal: total}, nil
}

func (s *ReportService) sumPostingsByType(accountType, dateFrom, dateTo string) (model.StringDecimal, error) {
    query := `
        SELECT p.amount
        FROM postings p
        JOIN transactions t ON t.id = p.transaction_id
        JOIN accounts a ON a.id = p.account_id
        WHERE t.date >= ? AND t.date <= ?
          AND (a.account_type = ? OR a.full_name LIKE (
              SELECT full_name || ':%' FROM accounts WHERE account_type = ? LIMIT 1
          ))
    `

    rows, err := s.db.Query(query, dateFrom, dateTo, accountType, accountType)
    if err != nil {
        return model.StringDecimal{}, model.NewInternalError("failed to query postings: " + err.Error())
    }
    defer rows.Close()

    var total decimal.Decimal
    for rows.Next() {
        var amount string
        if err := rows.Scan(&amount); err != nil {
            return model.StringDecimal{}, model.NewInternalError("failed to scan: " + err.Error())
        }
        amt, err := decimal.NewFromString(amount)
        if err != nil {
            return model.StringDecimal{}, model.NewInternalError("failed to parse amount: " + err.Error())
        }
        total = total.Add(amt)
    }
    if err := rows.Err(); err != nil {
        return model.StringDecimal{}, model.NewInternalError("failed to iterate rows: " + err.Error())
    }

    return model.StringDecimal{Decimal: total}, nil
}
