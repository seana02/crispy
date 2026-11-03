package model

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

// AccountType represents the type of a root-level account.
type AccountType string

const (
    AccountTypeAsset     AccountType = "asset"
    AccountTypeLiability AccountType = "liability"
    AccountTypeRevenue   AccountType = "revenue"
    AccountTypeExpense   AccountType = "expense"
    AccountTypeEquity    AccountType = "equity"
)

func ValidAccountType(s string) bool {
    switch AccountType(s) {
    case AccountTypeAsset, AccountTypeLiability, AccountTypeRevenue, AccountTypeExpense, AccountTypeEquity:
        return true
    }
    return false
}

// SuggestAccountType returns a suggested AccountType based on the account name.
func SuggestAccountType(name string) AccountType {
    switch {
    case matchesAny(name, "asset", "assets", "checking", "savings", "bank", "cash", "wallet"):
        return AccountTypeAsset
    case matchesAny(name, "liability", "liabilities", "credit", "loan", "loans", "mortgage", "debt"):
        return AccountTypeLiability
    case matchesAny(name, "revenue", "revenues", "income", "salary", "wages"):
        return AccountTypeRevenue
    case matchesAny(name, "expense", "expenses", "spending", "cost", "costs"):
        return AccountTypeExpense
    case matchesAny(name, "equity", "opening", "retained"):
        return AccountTypeEquity
    }
    return ""
}

func matchesAny(name string, options ...string) bool {
    lower := toLower(name)
    for _, opt := range options {
        if lower == opt {
            return true
        }
    }
    return false
}

func toLower(s string) string {
    b := make([]byte, len(s))
    for i := range s {
        c := s[i]
        if c >= 'A' && c <= 'Z' {
            c += 'a' - 'A'
        }
        b[i] = c
    }
    return string(b)
}

// StringDecimal wraps decimal.Decimal with custom JSON marshaling to output strings.
type StringDecimal struct {
    decimal.Decimal
}

func (d StringDecimal) MarshalJSON() ([]byte, error) {
    return json.Marshal(d.String())
}

func (d *StringDecimal) UnmarshalJSON(data []byte) error {
    var s string
    if err := json.Unmarshal(data, &s); err != nil {
        var f float64
        if err2 := json.Unmarshal(data, &f); err2 != nil {
            return err
        }
        d.Decimal = decimal.NewFromFloat(f)
        return nil
    }
    dec, err := decimal.NewFromString(s)
    if err != nil {
        return err
    }
    d.Decimal = dec
    return nil
}

func (d StringDecimal) String() string {
    return d.Decimal.String()
}

// Account represents a financial account.
type Account struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	FullName    string     `json:"full_name"`
	ParentID    *int64     `json:"parent_id"`
	AccountType string     `json:"account_type"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
}

// Transaction represents a financial transaction.
type Transaction struct {
	ID                   int64      `json:"id"`
	Date                 string     `json:"date"`
	Description          string     `json:"description"`
	Status               string     `json:"status"`
	RecurringTemplateID  *int64     `json:"recurring_template_id"`
	CreatedAt            string     `json:"created_at"`
	UpdatedAt            string     `json:"updated_at"`
	Postings             []Posting  `json:"postings"`
	Tags                 []Tag      `json:"tags"`
}

// TransactionStatus values.
const (
    StatusPending = "pending"
    StatusPosted  = "posted"
    StatusCleared = "cleared"
)

func ValidStatus(s string) bool {
    return s == StatusPending || s == StatusPosted || s == StatusCleared
}

// Posting represents a single posting within a transaction.
type Posting struct {
	ID              int64         `json:"id"`
	TransactionID   int64         `json:"transaction_id"`
	AccountID       int64         `json:"account_id"`
	AccountFullName string        `json:"account_full_name"` // populated on reads, empty on creates
	Currency        string        `json:"currency"`
	Amount          StringDecimal `json:"amount"`
}

// Tag represents a key-value tag.
type Tag struct {
    ID    int64  `json:"id"`
    Key   string `json:"key"`
    Value string `json:"value"`
}

// Budget represents a budget for tracking spending or saving.
type Budget struct {
	ID           int64         `json:"id"`
	Name         string        `json:"name"`
	Direction    string        `json:"direction"`
	Period       string        `json:"period"`
	Currency     string        `json:"currency"`
	TargetAmount StringDecimal `json:"target_amount"`
	CreatedAt    string        `json:"created_at"`
	UpdatedAt    string        `json:"updated_at"`
	AccountIDs   []int64       `json:"account_ids"`
	AccountNames []string      `json:"account_names"` // full_names of tracked accounts
	TagIDs       []int64       `json:"tag_ids"`
}

// BudgetDirection values.
const (
    DirectionSpending = "spending"
    DirectionSaving   = "saving"
)

// BudgetPeriod values.
const (
    PeriodWeekly  = "weekly"
    PeriodMonthly = "monthly"
    PeriodYearly  = "yearly"
)

// ExchangeRate represents an exchange rate between two currencies.
type ExchangeRate struct {
	ID            int64       `json:"id"`
	BaseCurrency  string      `json:"base_currency"`
	QuoteCurrency string      `json:"quote_currency"`
	Rate          StringDecimal `json:"rate"`
	EffectiveDate string      `json:"effective_date"`
	Source        string      `json:"source"`
	CreatedAt     string      `json:"created_at"`
}

// RecurringTemplate represents a recurring transaction template.
type RecurringTemplate struct {
	ID                int64      `json:"id"`
	Description       string     `json:"description"`
	Status            string     `json:"status"`
	ScheduleType      string     `json:"schedule_type"`
	IntervalN         int        `json:"interval_n"`
	DayOfMonth        *int       `json:"day_of_month"`
	StartDate         string     `json:"start_date"`
	EndDate           *string    `json:"end_date"`
	LastGeneratedDate *string    `json:"last_generated_date"`
	IsActive          bool       `json:"is_active"`
	CreatedAt         string     `json:"created_at"`
	UpdatedAt         string     `json:"updated_at"`
	Postings          []TemplatePosting `json:"postings"`
	Tags              []Tag      `json:"tags"`
}

// ScheduleType values.
const (
    ScheduleEveryNDays   = "every_n_days"
    ScheduleEveryNWeeks  = "every_n_weeks"
    ScheduleEveryNMonths = "every_n_months"
    ScheduleMonthlyOnDay = "monthly_on_day"
)

// TemplatePosting represents a posting in a recurring template.
type TemplatePosting struct {
	ID              int64         `json:"id"`
	TemplateID      int64         `json:"template_id"`
	AccountID       int64         `json:"account_id"`
	AccountFullName string        `json:"account_full_name"`
	Currency        string        `json:"currency"`
	Amount          StringDecimal `json:"amount"`
}

// WidgetConfig represents a widget on the home page.
type WidgetConfig struct {
	ID         int64     `json:"id"`
	WidgetType string    `json:"widget_type"`
	PositionX  int       `json:"position_x"`
	PositionY  int       `json:"position_y"`
	Width      int       `json:"width"`
	Height     int       `json:"height"`
	Config     string    `json:"config"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
}

// Setting represents a key-value setting.
type Setting struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

// Pagination parameters for list queries.
type Pagination struct {
    Page     int `json:"page"`
    PageSize int `json:"page_size"`
}

// PaginatedResult wraps a result set with pagination info.
type PaginatedResult[T any] struct {
    Items      []T `json:"items"`
    Total      int `json:"total"`
    Page       int `json:"page"`
    PageSize   int `json:"page_size"`
    TotalPages int `json:"total_pages"`
}
