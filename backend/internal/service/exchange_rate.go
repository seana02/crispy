package service

import (
    "database/sql"
    "log/slog"

    "github.com/crispy-2/backend/internal/dal/generated"
    "github.com/crispy-2/backend/internal/model"
    "github.com/shopspring/decimal"
)

// ExchangeRateService handles exchange rate business logic.
type ExchangeRateService struct {
    db      *sql.DB
    queries *generated.Queries
}

// NewExchangeRateService creates a new ExchangeRateService.
func NewExchangeRateService(db *sql.DB) *ExchangeRateService {
    return &ExchangeRateService{
        db:      db,
        queries: generated.New(db),
    }
}

// CreateRate creates a new exchange rate.
func (s *ExchangeRateService) CreateRate(baseCurrency, quoteCurrency, rate, effectiveDate, source string) (*model.ExchangeRate, error) {
    if baseCurrency == "" || quoteCurrency == "" {
        return nil, model.NewValidationError("base_currency and quote_currency are required")
    }
    if baseCurrency == quoteCurrency {
        return nil, model.NewValidationError("base_currency and quote_currency must be different")
    }
    rateDecimal, err := decimal.NewFromString(rate)
    if err != nil {
        return nil, model.NewValidationError("rate must be a positive decimal value")
    }
    if rateDecimal.IsZero() || rateDecimal.IsNegative() {
        return nil, model.NewValidationError("rate must be a positive decimal value")
    }
    if effectiveDate == "" {
        return nil, model.NewValidationError("effective_date is required")
    }
    if source == "" {
        source = "manual"
    }
    if source != "manual" && source != "api" {
        return nil, model.NewValidationError("source must be 'manual' or 'api'")
    }

    row, err := s.queries.CreateExchangeRate(ctx(), generated.CreateExchangeRateParams{
        BaseCurrency:  baseCurrency,
        QuoteCurrency: quoteCurrency,
        Rate:          rate,
        EffectiveDate: effectiveDate,
        Source:        source,
    })
    if err != nil {
        return nil, model.NewInternalError("failed to create exchange rate: " + err.Error())
    }

    slog.Info("exchange rate created", "id", row.ID, "pair", baseCurrency+"/"+quoteCurrency)
    return toModelExchangeRate(row), nil
}

// GetRate gets the effective exchange rate for a currency pair on a given date.
// Falls back to inverse pair if the direct pair is not found.
func (s *ExchangeRateService) GetRate(baseCurrency, quoteCurrency, date string) (*model.ExchangeRate, error) {
    if baseCurrency == quoteCurrency {
        // Same currency: rate is 1
        return &model.ExchangeRate{
            BaseCurrency:  baseCurrency,
            QuoteCurrency: quoteCurrency,
            Rate:          model.StringDecimal{Decimal: decimal.NewFromInt(1)},
            EffectiveDate: date,
            Source:        "identity",
        }, nil
    }

    // Try direct pair
    row, err := s.queries.GetEffectiveRate(ctx(), generated.GetEffectiveRateParams{
        BaseCurrency:  baseCurrency,
        QuoteCurrency: quoteCurrency,
        EffectiveDate: date,
    })
    if err == nil {
        return toModelExchangeRate(row), nil
    }
    if err != sql.ErrNoRows {
        return nil, model.NewInternalError("failed to get exchange rate: " + err.Error())
    }

    // Try inverse pair
    inverseRow, err := s.queries.GetEffectiveRate(ctx(), generated.GetEffectiveRateParams{
        BaseCurrency:  quoteCurrency,
        QuoteCurrency: baseCurrency,
        EffectiveDate: date,
    })
    if err == nil {
        inverseRate, _ := decimal.NewFromString(inverseRow.Rate)
        slog.Warn("using inverse exchange rate", "pair", baseCurrency+"/"+quoteCurrency)
        return &model.ExchangeRate{
            ID:            inverseRow.ID,
            BaseCurrency:  baseCurrency,
            QuoteCurrency: quoteCurrency,
            Rate:          model.StringDecimal{Decimal: decimal.NewFromInt(1).Div(inverseRate)},
            EffectiveDate: inverseRow.EffectiveDate,
            Source:        inverseRow.Source,
        }, nil
    }
    if err != sql.ErrNoRows {
        return nil, model.NewInternalError("failed to get inverse exchange rate: " + err.Error())
    }

    return nil, model.NewNotFoundError("exchange rate", baseCurrency+"/"+quoteCurrency+" on "+date)
}

// ListRates returns a paginated list of exchange rates.
func (s *ExchangeRateService) ListRates(baseCurrency, quoteCurrency string, pg model.Pagination) (*model.PaginatedResult[model.ExchangeRate], error) {
    pg = normalizePagination(pg)
    offset := int64((pg.Page - 1) * pg.PageSize)
    limit := int64(pg.PageSize)

    var rows []generated.ExchangeRate
    var total int64
    var err error

    if baseCurrency != "" && quoteCurrency != "" {
        rows, err = s.queries.ListExchangeRatesByPair(ctx(), generated.ListExchangeRatesByPairParams{
            BaseCurrency:  baseCurrency,
            QuoteCurrency: quoteCurrency,
            Limit:         limit,
            Offset:        offset,
        })
        if err == nil {
            total, err = s.queries.CountExchangeRatesByPair(ctx(), generated.CountExchangeRatesByPairParams{
                BaseCurrency:  baseCurrency,
                QuoteCurrency: quoteCurrency,
            })
        }
    } else {
        rows, err = s.queries.ListExchangeRates(ctx(), generated.ListExchangeRatesParams{
            Limit:  limit,
            Offset: offset,
        })
        if err == nil {
            total, err = s.queries.CountExchangeRates(ctx())
        }
    }

    if err != nil {
        return nil, model.NewInternalError("failed to list exchange rates: " + err.Error())
    }

    items := make([]model.ExchangeRate, len(rows))
    for i, r := range rows {
        items[i] = *toModelExchangeRate(r)
    }

    totalPages := int(total) / pg.PageSize
    if int(total)%pg.PageSize != 0 {
        totalPages++
    }

    return &model.PaginatedResult[model.ExchangeRate]{
        Items:      items,
        Total:      int(total),
        Page:       pg.Page,
        PageSize:   pg.PageSize,
        TotalPages: totalPages,
    }, nil
}

// DeleteRate deletes an exchange rate.
func (s *ExchangeRateService) DeleteRate(id int64) error {
    _, err := s.queries.GetExchangeRate(ctx(), id)
    if err != nil {
        if err == sql.ErrNoRows {
            return model.NewNotFoundError("exchange rate", id)
        }
        return model.NewInternalError("failed to get exchange rate: " + err.Error())
    }

    if err := s.queries.DeleteExchangeRate(ctx(), id); err != nil {
        return model.NewInternalError("failed to delete exchange rate: " + err.Error())
    }

    slog.Info("exchange rate deleted", "id", id)
    return nil
}

func toModelExchangeRate(row generated.ExchangeRate) *model.ExchangeRate {
	rate, _ := decimal.NewFromString(row.Rate)
	return &model.ExchangeRate{
		ID:            row.ID,
		BaseCurrency:  row.BaseCurrency,
		QuoteCurrency: row.QuoteCurrency,
		Rate:          model.StringDecimal{Decimal: rate},
		EffectiveDate: row.EffectiveDate,
		Source:        row.Source,
		CreatedAt:     row.CreatedAt,
	}
}

// FetchRatesFromAPI fetches exchange rates from an external API.
// TODO: Implement actual API provider (e.g., exchangerate.host, Open Exchange Rates).
func (s *ExchangeRateService) FetchRatesFromAPI(baseCurrency string, quoteCurrencies []string) error {
    return model.NewValidationError("exchange rate API provider not configured")
}
