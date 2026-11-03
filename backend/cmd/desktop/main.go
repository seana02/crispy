//go:build wails

// Desktop entry point for the Wails standalone application.
//
// Prerequisites:
//   - Install Wails CLI: go install github.com/wailsapp/wails/v2/cmd/wails@latest
//   - Add Wails to go.mod: go get github.com/wailsapp/wails/v2
//   - Place frontend build output at frontend/dist/ (run: cd frontend && npm run build)
//   - From this directory, run: wails build
//
// The frontend is the same SolidJS app as web mode — Wails embeds the
// frontend/dist/ directory as a static webview. The frontend's API layer
// is swapped from HTTP fetch calls to the Wails-generated JS bindings,
// which call directly into the Go service layer.
package main

import (
	"database/sql"
	"log/slog"
	"os"
	"time"

	"github.com/crispy-2/backend/internal/dal"
	"github.com/crispy-2/backend/internal/service"
	wailsbindings "github.com/crispy-2/backend/internal/transport/wails"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

func main() {
	// Configure logging
	logLevel := slog.LevelInfo
	if os.Getenv("LOG_LEVEL") == "debug" {
		logLevel = slog.LevelDebug
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(handler))

	// Database setup
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			userConfigDir = "."
		}
		dbPath = userConfigDir + "/crispy/crispy.db"
		if err := os.MkdirAll(userConfigDir+"/crispy", 0755); err != nil {
			slog.Warn("failed to create config directory, using current dir", "error", err)
			dbPath = "crispy.db"
		}
	}

	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		slog.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Run migrations
	goose.SetBaseFS(dal.Migrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		slog.Error("failed to set goose dialect", "error", err)
		os.Exit(1)
	}
	if err := goose.Up(db, "migrations"); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Initialise services
	exchangeRateSvc := service.NewExchangeRateService(db)
	accountSvc := service.NewAccountService(db)
	txSvc := service.NewTransactionService(db, exchangeRateSvc)
	budgetSvc := service.NewBudgetService(db, exchangeRateSvc)
	recurringSvc := service.NewRecurringService(db, txSvc)
	reportSvc := service.NewReportService(db, exchangeRateSvc)
	widgetSvc := service.NewWidgetService(db)
	settingsSvc := service.NewSettingsService(db)
	searchSvc := service.NewSearchService(db)

	// Catch up on any missed recurring transactions
	today := time.Now().Format("2006-01-02")
	if count, err := recurringSvc.GenerateDueTransactions(today); err != nil {
		slog.Warn("failed to generate recurring transactions on startup", "error", err)
	} else if count > 0 {
		slog.Info("generated recurring transactions on startup", "count", count)
	}

	// Create the Wails app binding
	app := wailsbindings.NewApp(
		accountSvc,
		txSvc,
		budgetSvc,
		exchangeRateSvc,
		recurringSvc,
		reportSvc,
		widgetSvc,
		settingsSvc,
		searchSvc,
	)

	// Run the Wails application
	if err := wails.Run(&options.App{
		Title:  "Crispy",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			// Wails embeds the frontend build into the binary.
			// Run `cd frontend && npm run build` before `wails build`.
			Assets: os.DirFS("../../frontend/dist"),
		},
		OnStartup: app.Startup,
		Bind:      []interface{}{app},
	}); err != nil {
		slog.Error("wails run error", "error", err)
		os.Exit(1)
	}
}
