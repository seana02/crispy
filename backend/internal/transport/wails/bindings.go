//go:build wails

// Package wails provides the Wails desktop transport layer.
// This file exposes the service layer as Wails-bound methods, callable
// directly from the SolidJS frontend via the generated Wails JS bindings
// (instead of HTTP fetch calls used in web mode).
//
// To build the desktop app:
//   1. Install Wails v2: go install github.com/wailsapp/wails/v2/cmd/wails@latest
//   2. Add Wails to go.mod: go get github.com/wailsapp/wails/v2
//   3. Run: wails build  (from backend/cmd/desktop/)
//
// The service layer is identical to web mode — only the transport differs.
package wails

import (
	"context"

	"github.com/crispy-2/backend/internal/model"
	"github.com/crispy-2/backend/internal/service"
)

// App is the Wails application struct. All public methods on this struct
// are automatically exposed to the frontend as callable JS functions.
type App struct {
	ctx context.Context

	accountSvc      *service.AccountService
	transactionSvc  *service.TransactionService
	budgetSvc       *service.BudgetService
	exchangeRateSvc *service.ExchangeRateService
	recurringSvc    *service.RecurringService
	reportSvc       *service.ReportService
	widgetSvc       *service.WidgetService
	settingsSvc     *service.SettingsService
	searchSvc       *service.SearchService
}

// NewApp creates a new Wails App with all services injected.
func NewApp(
	accountSvc *service.AccountService,
	transactionSvc *service.TransactionService,
	budgetSvc *service.BudgetService,
	exchangeRateSvc *service.ExchangeRateService,
	recurringSvc *service.RecurringService,
	reportSvc *service.ReportService,
	widgetSvc *service.WidgetService,
	settingsSvc *service.SettingsService,
	searchSvc *service.SearchService,
) *App {
	return &App{
		accountSvc:      accountSvc,
		transactionSvc:  transactionSvc,
		budgetSvc:       budgetSvc,
		exchangeRateSvc: exchangeRateSvc,
		recurringSvc:    recurringSvc,
		reportSvc:       reportSvc,
		widgetSvc:       widgetSvc,
		settingsSvc:     settingsSvc,
		searchSvc:       searchSvc,
	}
}

// Startup is called by Wails when the application starts.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// ── Account Service ──────────────────────────────────────────────────────────

func (a *App) ListAccounts(nameFilter, typeFilter string, page, pageSize int) (*model.PaginatedResult[model.Account], error) {
	return a.accountSvc.ListAccounts(nameFilter, typeFilter, model.Pagination{Page: page, PageSize: pageSize})
}

func (a *App) GetAccount(id int64) (*model.Account, error) {
	return a.accountSvc.GetAccount(id)
}

func (a *App) CreateAccount(name string, parentID *int64, accountType string) (*model.Account, error) {
	return a.accountSvc.CreateAccount(name, parentID, accountType)
}

func (a *App) UpdateAccount(id int64, name *string, parentID **int64, accountType string) (*model.Account, error) {
	return a.accountSvc.UpdateAccount(id, name, parentID, accountType)
}

func (a *App) DeleteAccount(id int64) error {
	return a.accountSvc.DeleteAccount(id)
}

func (a *App) SearchAccounts(query string) ([]model.Account, error) {
	return a.accountSvc.SearchAccounts(query)
}

// ── Transaction Service ──────────────────────────────────────────────────────

func (a *App) ListTransactions(filters service.TransactionFilters, page, pageSize int) (*model.PaginatedResult[model.Transaction], error) {
	return a.transactionSvc.ListTransactions(filters, model.Pagination{Page: page, PageSize: pageSize})
}

func (a *App) GetTransaction(id int64) (*model.Transaction, error) {
	return a.transactionSvc.GetTransaction(id)
}

func (a *App) CreateTransaction(date, description, status string, postings []service.PostingInput, tags []service.TagInput) (*model.Transaction, error) {
	return a.transactionSvc.CreateTransaction(date, description, status, postings, tags)
}

func (a *App) UpdateTransaction(id int64, date, description, status string, postings []service.PostingInput, tags []service.TagInput) (*model.Transaction, error) {
	return a.transactionSvc.UpdateTransaction(id, date, description, status, postings, tags)
}

func (a *App) DeleteTransaction(id int64) error {
	return a.transactionSvc.DeleteTransaction(id)
}

// ── Budget Service ───────────────────────────────────────────────────────────

func (a *App) ListBudgets() ([]model.Budget, error) {
	return a.budgetSvc.ListBudgets()
}

func (a *App) GetBudget(id int64) (*model.Budget, error) {
	return a.budgetSvc.GetBudget(id)
}

func (a *App) CreateBudget(name, direction, period, currency, targetAmount string, accountIDs []int64, tagIDs []int64) (*model.Budget, error) {
	return a.budgetSvc.CreateBudget(name, direction, period, currency, targetAmount, accountIDs, tagIDs)
}

func (a *App) UpdateBudget(id int64, name, direction, period, currency, targetAmount string, accountIDs []int64, tagIDs []int64) (*model.Budget, error) {
	return a.budgetSvc.UpdateBudget(id, name, direction, period, currency, targetAmount, accountIDs, tagIDs)
}

func (a *App) DeleteBudget(id int64) error {
	return a.budgetSvc.DeleteBudget(id)
}

// ── Exchange Rate Service ────────────────────────────────────────────────────

func (a *App) ListRates(baseCurrency, quoteCurrency string, page, pageSize int) (*model.PaginatedResult[model.ExchangeRate], error) {
	return a.exchangeRateSvc.ListRates(baseCurrency, quoteCurrency, model.Pagination{Page: page, PageSize: pageSize})
}

func (a *App) GetRate(baseCurrency, quoteCurrency, date string) (*model.ExchangeRate, error) {
	return a.exchangeRateSvc.GetRate(baseCurrency, quoteCurrency, date)
}

func (a *App) CreateRate(baseCurrency, quoteCurrency, rate, effectiveDate, source string) (*model.ExchangeRate, error) {
	return a.exchangeRateSvc.CreateRate(baseCurrency, quoteCurrency, rate, effectiveDate, source)
}

func (a *App) DeleteRate(id int64) error {
	return a.exchangeRateSvc.DeleteRate(id)
}

func (a *App) FetchRatesFromAPI(baseCurrency string, quoteCurrencies []string) error {
	return a.exchangeRateSvc.FetchRatesFromAPI(baseCurrency, quoteCurrencies)
}

// ── Recurring Service ────────────────────────────────────────────────────────

func (a *App) ListTemplates() ([]model.RecurringTemplate, error) {
	return a.recurringSvc.ListTemplates()
}

func (a *App) GetTemplate(id int64) (*model.RecurringTemplate, error) {
	return a.recurringSvc.GetTemplate(id)
}

func (a *App) DeleteTemplate(id int64) error {
	return a.recurringSvc.DeleteTemplate(id)
}

func (a *App) PreviewSchedule(scheduleType string, intervalN int, dayOfMonth *int, startDate string, endDate *string, count int) []string {
	return a.recurringSvc.PreviewSchedule(scheduleType, intervalN, dayOfMonth, startDate, endDate, count)
}

func (a *App) GenerateDueTransactions(asOfDate string) (int, error) {
	return a.recurringSvc.GenerateDueTransactions(asOfDate)
}

// ── Report Service ───────────────────────────────────────────────────────────

func (a *App) GetBalanceSheet(dateFrom, dateTo string, depth int, tagFilters []int64) (*service.BalanceSheetReport, error) {
	return a.reportSvc.GetBalanceSheet(dateFrom, dateTo, depth, tagFilters)
}

func (a *App) GetIncomeStatement(dateFrom, dateTo string, depth int, tagFilters []int64) (*service.IncomeStatementReport, error) {
	return a.reportSvc.GetIncomeStatement(dateFrom, dateTo, depth, tagFilters)
}

func (a *App) GetNetWorth(asOfDate string) (*service.NetWorthWithTrend, error) {
	return a.reportSvc.GetNetWorth(asOfDate)
}

func (a *App) GetCashFlow(dateFrom, dateTo string) (model.StringDecimal, error) {
	return a.reportSvc.GetCashFlow(dateFrom, dateTo)
}

func (a *App) GetTopSpendingCategories(dateFrom, dateTo string, limit int) ([]service.SpendingCategory, error) {
	return a.reportSvc.GetTopSpendingCategories(dateFrom, dateTo, limit)
}

// ── Widget Service ───────────────────────────────────────────────────────────

func (a *App) ListWidgets() ([]model.WidgetConfig, error) {
	return a.widgetSvc.ListWidgets()
}

func (a *App) CreateWidget(widgetType string, posX, posY, width, height int, config string) (*model.WidgetConfig, error) {
	return a.widgetSvc.CreateWidget(widgetType, posX, posY, width, height, config)
}

func (a *App) UpdateWidget(id int64, widgetType string, posX, posY, width, height int, config string) (*model.WidgetConfig, error) {
	return a.widgetSvc.UpdateWidget(id, widgetType, posX, posY, width, height, config)
}

func (a *App) DeleteWidget(id int64) error {
	return a.widgetSvc.DeleteWidget(id)
}

// ── Settings Service ─────────────────────────────────────────────────────────

func (a *App) GetSettings() ([]model.Setting, error) {
	return a.settingsSvc.GetSettings()
}

func (a *App) UpdateSetting(key, value string) error {
	return a.settingsSvc.UpdateSetting(key, value)
}

// ── Search Service ───────────────────────────────────────────────────────────

func (a *App) GlobalSearch(query string) (*service.SearchResult, error) {
	return a.searchSvc.GlobalSearch(query)
}
