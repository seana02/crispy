package main

import (
	"crispy/automation"
	"crispy/db"
	"crispy/domain"
	"crispy/scheduler"
	"crispy/service"
	"fmt"
)

type Dependencies struct {
	Service          *service.Service
	RecurringManager *scheduler.Manager
	RuleEngine       *automation.Engine
	// CurrencyService    *currency.Service
}

func BuildDependencies(cfg *Config) (*Dependencies, error) {
	fmt.Println("Building dependencies")
	db, err := db.InitSQLite(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize sqlite database: %s", err)
	}
	fmt.Println("created db")

	service := service.NewService(db)

	recurringManager := scheduler.NewManager(service)
	ruleEngine := automation.NewEngine(service)
	// currencyService := currency.NewService(db)

	return &Dependencies{
		Service:          service,
		RecurringManager: recurringManager,
		RuleEngine:       ruleEngine,
		// CurrencyService:    currencyService,
	}, nil
}

// Returns the list of Go functions bound to Wails frontend
// func (a *Dependencies) WailsBindings() []interface{} {
// 	return []interface{}{
// 		a.TransactionService,
// 		a.AccountService,
// 		a.RecurringManager,
// 		// a.CurrencyService,
// 		a.RuleEngine,
// 	}
// }

// Returns list of Enum bindings
func (a *Dependencies) WailsEnumBindings() []interface{} {
	var StatusBindings = []struct {
		Value  domain.Status
		TSName string
	}{
		{domain.Pending, "Pending"},
		{domain.Posted, "Posted"},
		{domain.Cleared, "Cleared"},
	}

	var AccountTypeBindings = []struct {
		Value  domain.Type
		TSName string
	}{
		{domain.Asset, "Asset"},
		{domain.Liability, "Liability"},
		{domain.Revenue, "Revenue"},
		{domain.Expense, "Expense"},
		{domain.Equity, "Equity"},
	}

	return []interface{}{
		StatusBindings,
		AccountTypeBindings,
	}
}
