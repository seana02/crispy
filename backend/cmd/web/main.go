package main

import (
    "database/sql"
    "log/slog"
    "net/http"
    "os"
    "time"

    "github.com/crispy-2/backend/internal/dal"
    "github.com/crispy-2/backend/internal/service"
    "github.com/crispy-2/backend/internal/transport/rest"
    _ "github.com/mattn/go-sqlite3"
    "github.com/pressly/goose/v3"
)

// responseRecorder wraps http.ResponseWriter to capture the status code for logging.
type responseRecorder struct {
    http.ResponseWriter
    status int
}

func (r *responseRecorder) WriteHeader(status int) {
    r.status = status
    r.ResponseWriter.WriteHeader(status)
}

// requestLogger logs every HTTP request with method, path, status, and duration.
func requestLogger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        rec := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
        next.ServeHTTP(rec, r)
        duration := time.Since(start)
        level := slog.LevelDebug
        if rec.status >= 500 {
            level = slog.LevelError
        } else if rec.status >= 400 {
            level = slog.LevelWarn
        }
        slog.Log(r.Context(), level, "request",
            "method", r.Method,
            "path", r.URL.Path,
            "status", rec.status,
            "duration", duration,
        )
    })
}

func main() {
    // Configure logging
    logLevel := slog.LevelInfo
    if os.Getenv("LOG_LEVEL") == "debug" {
        logLevel = slog.LevelDebug
    }

    var handler slog.Handler
    if os.Getenv("ENV") == "production" {
        handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
    } else {
        handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
    }
    slog.SetDefault(slog.New(handler))

    // Database setup
    dbPath := os.Getenv("DB_PATH")
    if dbPath == "" {
        dbPath = "crispy.db"
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

    // Initialize services
    exchangeRateSvc := service.NewExchangeRateService(db)
    accountSvc := service.NewAccountService(db)
    txSvc := service.NewTransactionService(db, exchangeRateSvc)
    txSvc.SetAccountService(accountSvc) // enables auto-creation of accounts by name in postings
    budgetSvc := service.NewBudgetService(db, exchangeRateSvc)
    recurringSvc := service.NewRecurringService(db, txSvc)
    reportSvc := service.NewReportService(db, exchangeRateSvc)
    widgetSvc := service.NewWidgetService(db)
    settingsSvc := service.NewSettingsService(db)
    searchSvc := service.NewSearchService(db)

    // Generate any overdue recurring transactions on startup
    today := time.Now().Format("2006-01-02")
    if count, err := recurringSvc.GenerateDueTransactions(today); err != nil {
        slog.Warn("failed to generate recurring transactions on startup", "error", err)
    } else if count > 0 {
        slog.Info("generated recurring transactions on startup", "count", count)
    }

    // Set up router
    mux := rest.SetupRouter(
        db,
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

    // Serve the frontend static files
    frontendDir := os.Getenv("FRONTEND_DIR")
    if frontendDir == "" {
        frontendDir = "../frontend/dist"
    }
    if _, err := os.Stat(frontendDir); err == nil {
        fs := http.FileServer(http.Dir(frontendDir))
        mux.Handle("/", fs)
    }

    // Wrap with request logging middleware
    httpHandler := requestLogger(mux)

    // Start server
    addr := os.Getenv("ADDR")
    if addr == "" {
        addr = ":8080"
    }

    slog.Info("starting server", "addr", addr)
    if err := http.ListenAndServe(addr, httpHandler); err != nil {
        slog.Error("server error", "error", err)
        os.Exit(1)
    }
}
