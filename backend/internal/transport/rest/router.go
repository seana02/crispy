package rest

import (
	"database/sql"
	"net/http"

	"github.com/crispy-2/backend/internal/service"
)

func SetupRouter(
	db *sql.DB,
	accountSvc *service.AccountService,
	transactionSvc *service.TransactionService,
	budgetSvc *service.BudgetService,
	exchangeRateSvc *service.ExchangeRateService,
	recurringSvc *service.RecurringService,
	reportSvc *service.ReportService,
	widgetSvc *service.WidgetService,
	settingsSvc *service.SettingsService,
	searchSvc *service.SearchService,
) *http.ServeMux {
	mux := http.NewServeMux()

	tagHandler := NewTagHandler(db)

	// Account routes
	mux.HandleFunc("GET /api/accounts", listAccountsHandler(accountSvc))
	mux.HandleFunc("POST /api/accounts", createAccountHandler(accountSvc))
	mux.HandleFunc("GET /api/accounts/search", searchAccountsHandler(accountSvc))
	mux.HandleFunc("GET /api/accounts/{id}", getAccountHandler(accountSvc))
	mux.HandleFunc("PUT /api/accounts/{id}", updateAccountHandler(accountSvc))
	mux.HandleFunc("DELETE /api/accounts/{id}", deleteAccountHandler(accountSvc))

	// Transaction routes
	mux.HandleFunc("POST /api/transactions", createTransactionHandler(transactionSvc))
	mux.HandleFunc("GET /api/transactions/{id}", getTransactionHandler(transactionSvc))
	mux.HandleFunc("GET /api/transactions", listTransactionsHandler(transactionSvc))
	mux.HandleFunc("PUT /api/transactions/{id}", updateTransactionHandler(transactionSvc))
	mux.HandleFunc("DELETE /api/transactions/{id}", deleteTransactionHandler(transactionSvc))

	// Budget routes
	mux.HandleFunc("POST /api/budgets", createBudgetHandler(budgetSvc))
	mux.HandleFunc("GET /api/budgets/{id}", getBudgetHandler(budgetSvc))
	mux.HandleFunc("GET /api/budgets", listBudgetsHandler(budgetSvc))
	mux.HandleFunc("PUT /api/budgets/{id}", updateBudgetHandler(budgetSvc))
	mux.HandleFunc("DELETE /api/budgets/{id}", deleteBudgetHandler(budgetSvc))
	mux.HandleFunc("GET /api/budgets/{id}/compute-current", computeBudgetCurrentHandler(budgetSvc))

	// Exchange rate routes
	mux.HandleFunc("POST /api/exchange-rates", createExchangeRateHandler(exchangeRateSvc))
	mux.HandleFunc("GET /api/exchange-rates/{base}/{quote}/{date}", getExchangeRateHandler(exchangeRateSvc))
	mux.HandleFunc("GET /api/exchange-rates", listExchangeRatesHandler(exchangeRateSvc))
	mux.HandleFunc("DELETE /api/exchange-rates/{id}", deleteExchangeRateHandler(exchangeRateSvc))
	mux.HandleFunc("POST /api/exchange-rates/fetch", fetchExchangeRatesHandler(exchangeRateSvc))

	// Recurring template routes
	mux.HandleFunc("POST /api/recurring-templates", createRecurringTemplateHandler(recurringSvc))
	mux.HandleFunc("GET /api/recurring-templates/{id}", getRecurringTemplateHandler(recurringSvc))
	mux.HandleFunc("GET /api/recurring-templates", listRecurringTemplatesHandler(recurringSvc))
	mux.HandleFunc("PUT /api/recurring-templates/{id}", updateRecurringTemplateHandler(recurringSvc))
	mux.HandleFunc("DELETE /api/recurring-templates/{id}", deleteRecurringTemplateHandler(recurringSvc))
	mux.HandleFunc("POST /api/recurring-templates/preview", previewRecurringScheduleHandler(recurringSvc))
	mux.HandleFunc("POST /api/recurring-templates/generate", generateRecurringTransactionsHandler(recurringSvc))

	// Tag routes
	mux.HandleFunc("GET /api/tags", tagHandler.listTagsHandler())

	// Report routes
	mux.HandleFunc("GET /api/reports/balance-sheet", getBalanceSheetHandler(reportSvc))
	mux.HandleFunc("GET /api/reports/income-statement", getIncomeStatementHandler(reportSvc))
	mux.HandleFunc("GET /api/reports/net-worth", getNetWorthHandler(reportSvc))
	mux.HandleFunc("GET /api/reports/net-worth-with-trend", getNetWorthHandlerWithTrend(reportSvc))
	mux.HandleFunc("GET /api/reports/cash-flow", getCashFlowHandler(reportSvc))
	mux.HandleFunc("GET /api/reports/top-spending", getTopSpendingHandler(reportSvc))
	mux.HandleFunc("GET /api/reports/expense-breakdown", getExpenseBreakdownHandler(reportSvc))
	mux.HandleFunc("GET /api/reports/tag-spending-summary", getTagSpendingSummaryHandler(reportSvc))

	// Widget routes
	mux.HandleFunc("POST /api/widgets", createWidgetHandler(widgetSvc))
	mux.HandleFunc("GET /api/widgets", listWidgetsHandler(widgetSvc))
	mux.HandleFunc("PUT /api/widgets/{id}", updateWidgetHandler(widgetSvc))
	mux.HandleFunc("DELETE /api/widgets/{id}", deleteWidgetHandler(widgetSvc))

	// Settings routes
	mux.HandleFunc("GET /api/settings", getSettingsHandler(settingsSvc))
	mux.HandleFunc("PUT /api/settings", updateSettingHandler(settingsSvc))
	mux.HandleFunc("PATCH /api/settings/bulk", updateSettingsBulkHandler(settingsSvc))

	// Search routes
	mux.HandleFunc("GET /api/search", globalSearchHandler(searchSvc))

	return mux
}
