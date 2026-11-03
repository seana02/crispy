# Overall Vision
This app is a GUI app that can work as a standalone desktop app or as a self-hosted web app.
  - **Standalone mode:** a desktop application using a framework like Wails (Go backend bundled with a webview, runs as a native window).
  - **Web mode:** the Go backend serves the SolidJS frontend as a self-hosted web server (user accesses via browser at e.g. `localhost:8080`).
The frontend is the same SolidJS app in both modes; the difference is only in how it's delivered (embedded webview vs. browser). The architecture should be modular enough that switching between modes requires minimal changes.
The financial model is based on hledger, including the double-entry bookkeeping system, and basic transaction tracking functionality, but with extra features as well.

# User Stories

1. The user can create Accounts.
Users can create Accounts, with any name they choose (except for any special characters with technical limitations), and can be a "root" level account or have any other account as parent.
The user can create a hierarchy of accounts for tracking finances, such as Expenses:Shopping:Online.
Just like hledger, the full name of the account will have parents separated by ":".
Root-level accounts have an explicit **account type** field: one of Asset, Liability, Revenue, Expense, or Equity. This type determines how the account (and all its children) behaves in reports (e.g. Asset/Liability accounts appear on the balance sheet, Revenue/Expense accounts appear on the income statement).
Subaccounts inherit their type from their root ancestor — the type field is only set on root accounts.
The app should auto-suggest a type based on the account name (e.g. an account named "Assets" defaults to Asset), but the user can override it.

2. The user can create Transactions.
The transactions can have a date, a short description, a status, zero or more key-value tags (e.g. `trip:japan-2026`, `project:kitchen-reno`), and an adjustable number of postings (minimum two).
Tags have a key and a value, allowing filtering by key alone (e.g. all `trip` tags), by key+value (e.g. `trip:japan-2026`), or grouping/reporting by tag values.
The status can be pending, posted, or cleared (free-form, no enforced progression).
Each posting can specify an account, the currency denomination (e.g. USD or JPY), and the amount.
Transactions must obey double-entry bookkeeping: all postings within a transaction must balance to zero. For multi-currency postings, balancing is evaluated using exchange rates (see User Story 9).
As a convenience, the last posting's amount may be left blank and will be auto-inferred as the balancing amount. This auto-infer only works when all postings in the transaction are the same currency; multi-currency transactions require all amounts to be explicitly entered.
The user should be able to select a date using a date selector.
The account should be able to be auto-completed, with fuzzy matching, but should also allow for new accounts that don't yet exist.
The account should be specified by the full name.

3. The user can create Budgets.
A budget has a time period (e.g. weekly, monthly, yearly), a currency, a target amount, a **direction** (spending or saving), and a target defined by:
  - **Accounts (required):** one or more accounts (each including all children). At least one account must be specified.
  - **Tags (optional):** zero or more tags to further narrow which transactions are included.
The budget measures the sum of postings matching (account1 OR account2 OR ...) AND optionally (tag1 OR tag2 OR ...).
The user can view the budget to see how much has been spent in the current period vs. the target amount.
A **spending** budget treats the target as a ceiling (warn when approaching/exceeding). A **saving** budget treats the target as a floor/goal (track progress toward reaching it). The underlying data model and computation are the same; only the display logic differs.

4. The user can view a balance sheet report.
The user can see their assets and liabilities as a table or graph.
They can see this data over multiple weeks/months/years simultaneously to compare or see the changes over time.
The account hierarchy depth is configurable (e.g. show only top 2 levels, collapsing everything below).
The report can be filtered by tags (e.g. show only transactions tagged `trip:japan-2026`).
Clicking on an account row navigates to that account's page.
It should replace hledger's balancesheet.

5. The user can view an income statement report.
The user can see their income and expenses as a table or graph.
They can see this data over multiple weeks/months/years simultaneously to compare or see the changes over time.
The account hierarchy depth is configurable (e.g. show only top 2 levels, collapsing everything below).
The report can be filtered by tags (e.g. show only transactions tagged `project:kitchen-reno`).
Clicking on an account row navigates to that account's page.
It should replace hledger's income statement.

6. The user can add and modify the home page using widgets.
Widgets can be placed on a home page grid to customize the home page.
Each widget covers a certain size of the grid, and can be dragged around.
The initial set of widget types:
  - **Account balance:** shows current balance of a specific account.
  - **Budget progress:** shows a specific budget's spent vs. target.
  - **Balance sheet summary:** a condensed version of the balance sheet report.
  - **Income statement summary:** a condensed version of the income statement report.
  - **Recent transactions:** list of N most recent transactions.
  - **Chart:** a configurable graph (e.g. spending over time for selected accounts).
  - **Net worth:** total assets minus total liabilities, with a trend indicator vs. last period.
  - **Cash flow:** net income minus expenses for the current period (week/month).
  - **Top spending categories:** bar or pie chart of top N expense accounts for the current period.
  - **Expense breakdown:** proportional breakdown of spending by account or tag (e.g. donut chart).
  - **Upcoming recurring transactions:** shows expected recurring transactions and whether they've appeared this period (depends on recurring transactions module, User Story 10).
  - **Tag spending summary:** spending totals grouped by a specific tag key.
  - **Exchange rate ticker:** shows current rates for configured currency pairs, with change indicators.
  - **Quick entry:** an inline form to quickly add a transaction without navigating away from the home page.

7. The user should be able to adjust global settings.
  - **Default currency:** the currency used by default for new postings.
  - **Display currency:** the currency used for converted totals in reports.
  - **Available currencies:** limit the currency options shown in the UI (for users who only work with certain currencies).
  - **Date format:** configurable display format for dates throughout the app.
  - **Theme selection:** choose from available themes (see Requirements #8).
  - **Fiscal year start:** configurable start month/day for annual reports (defaults to January 1).

8. The user should be able to inspect and modify the underlying data.
The user should be able to open the database in a database viewer and understand and modify the data.
In other words, it should not be obfuscated, encrypted, or otherwise rely on the backend to parse.

9. The user can manage exchange rates and commodity prices.
The user can manually enter exchange rates or commodity prices (e.g. USD/JPY, stock ticker prices) with an effective date, similar to hledger's `P` price directives.
The user can also optionally fetch current exchange rates from an external API for common fiat currencies.
These rates are used to:
  - Evaluate whether multi-currency transactions balance.
  - Convert amounts for display in reports (e.g. showing all values in a single "display currency").
  - Track the value of commodities/stocks over time.
The exchange rate module should be its own independent service/module on the backend.

10. The user can create recurring transactions.
A recurring transaction is a template that is structurally identical to a regular transaction (date, description, status, tags, postings) but is not itself recorded as a transaction.
The user specifies a **schedule** for when it triggers, and the app automatically creates real transactions according to that schedule.
The schedule should be flexible, supporting patterns like:
  - Every N days/weeks/months starting on a given date.
  - Every N-th month on a specific day of the month (e.g. "every 1 month on the 15th").
The user can optionally specify an **end date**; if omitted, the recurring template runs indefinitely until manually stopped or deleted.
The UI should show a **preview** of the next several transactions that would be generated (especially their dates) so the user can verify the schedule is correct before saving.

11. The app has interconnected detail pages.
  - **Account page:** shows the account's details (name, type, parent) and lists all transactions with postings to that account. Each transaction is clickable, navigating to its transaction page.
  - **Transaction page:** shows the transaction's full details (date, description, status, tags, all postings). Each posting's account is clickable, navigating to its account page.
  - **Tag page:** shows all transactions associated with a given tag key or key-value pair. Each transaction is clickable.
  - **Accounts list page:** shows all accounts as styled rows with searching, filtering, and pagination. Each row supports view, edit, and delete actions. Each account is clickable to its detail page.
  - **Transactions list page:** shows all transactions as styled rows with searching, filtering, and pagination. Each row supports view, edit, and delete actions. Each transaction is clickable to its detail page.
  - **Tags list page:** shows all tag keys (and their values) as styled rows with searching and filtering. Each tag is clickable to its tag page.
Throughout the app, accounts, transactions, and tags should be clickable links navigating to their respective detail pages wherever they appear.
This creates a navigable web: reports → accounts → transactions → accounts, with tags as an alternative entry point.

# Requirements
In addition to the requirements that come naturally from the user stories, there are other requirements

1. The app must be reasonably performant.
2. The app must be reasonably secure.
3. The code must be easily readable and maintainable by a human.
4. The code must be modular in design to allow for expansion and module-independent updates.
5. The frontend should be in TypeScript using Solid.
6. The backend should be in Golang.
7. The database for the standalone app should be SQLite.
8. The visual design of the UI should be easily customizable with themes (and new themes should be able to be set from within the app). Use plain CSS with nesting. It should be primarily be a dark-mode theme.
9. The visual design of charts should use more colors, like a terminal palette.

# Functional Design

## Cross-Cutting Concerns

### Decimal Handling
All monetary values (posting amounts, budget targets, exchange rates) are stored as text-encoded exact decimal strings in the database (e.g. "12.34", "0.00000834") and handled consistently across the stack:
- **Database:** stored as `TEXT` columns.
- **Backend:** handled using `shopspring/decimal`. All arithmetic (balancing checks, budget computation, currency conversion) uses this library. JSON serialization sends decimal values as strings.
- **Frontend:** received as strings from the API. **Never converted to JavaScript `number`/`float`** for monetary arithmetic (e.g. balance checking, summing). A JavaScript exact decimal library (`decimal.js` or `big.js`) is used for any client-side arithmetic (e.g. computing the auto-inferred balancing amount in the transaction form). For display-only formatting (commas, currency symbols, trailing zeros), conversion to `Number` is acceptable when the value is only used for `Intl.NumberFormat` display — this is safe for typical personal finance amounts. However, arithmetic or comparison operations must always use the exact decimal library.

### Date and Timestamp Handling
There are two distinct date types used throughout the app:

**Date-only fields** (transaction date, exchange rate effective date, schedule start/end dates, etc.):
- Represent a calendar date with no time or timezone component.
- **Database:** stored as `TEXT` in `YYYY-MM-DD` format.
- **Backend:** Go's `time.Time` with the time portion zeroed. JSON serialization outputs `YYYY-MM-DD` strings.
- **Frontend:** kept as plain `YYYY-MM-DD` strings throughout. **Never passed through JavaScript's `new Date()` constructor**, which interprets date-only strings as UTC midnight and can shift to the previous day in negative-offset timezones. Parsing, formatting, and comparison use `date-fns` or plain string operations. The date picker component must output a plain `YYYY-MM-DD` string, not a `Date` object. To get today's date, use a utility function (e.g. `todayLocalDate()`) that uses `date-fns` format rather than `new Date().toISOString()`.

**Timestamps** (`created_at`, `updated_at`):
- Represent actual instants in time for audit purposes.
- **Database:** stored as `TEXT` in ISO 8601 format with timezone offset (e.g. `2026-10-30T15:30:00-07:00`). The local timezone is used, with the offset preserved so both local time and UTC can be reconstructed. Timestamps are generated in Go using `time.Now().Format(time.RFC3339)` and passed as parameters to SQL queries.
- **Backend:** Go's `time.Time` type, serialized with timezone offset.
- **Frontend:** parsed and displayed using `date-fns` with the user's local timezone for display.

### Logging
The backend uses structured logging for debugging and operational visibility.
- **Library:** Go's built-in `log/slog` (available since Go 1.21). No external logging library needed.
- **Format:** JSON in production (machine-parseable), text in development (human-readable).
- **Levels:** DEBUG, INFO, WARN, ERROR. The log level is configurable (e.g. via a flag or environment variable).
- **What to log:**
  - **INFO:** service-level operations (transaction created, budget computed, recurring transactions generated, exchange rates fetched from API). Include relevant identifiers (e.g. transaction ID, template ID).
  - **WARN:** recoverable issues (e.g. exchange rate not found for a pair, falling back to inverse; recurring template skipped due to missing account).
  - **ERROR:** failures (database errors, validation failures that indicate a bug rather than user input, external API failures).
  - **DEBUG:** detailed internal state useful during development (query parameters, computed values, schedule date calculations).
- **Where to log:** in the service layer (where business logic decisions happen) and in the transport layer (incoming request summary, response status). The DAL layer does not log — errors bubble up to the service layer.
- **Context:** log entries should include contextual fields (e.g. `slog.With("account_id", id)`) so that related log lines can be correlated.

## Data Model

**General conventions:**
- All many-to-many relationships are represented as junction tables, not embedded arrays or JSON columns.
- All entities have an `id` primary key.
- Timestamps (`created_at`, `updated_at`) are stored as ISO 8601 with local timezone offset.
- All monetary values (posting amounts, budget targets, exchange rates) are stored as text-encoded exact decimal strings in the database (e.g. "12.34", "0.00000834") and handled using `shopspring/decimal` in the backend. This avoids all floating-point/base-10 rounding errors. A single approach is used consistently across all monetary fields.

### Account
- `id`: unique identifier
- `name`: the short name of this account (e.g. "Checking", not the full path)
- `full_name`: the full hierarchical name, auto-generated from the ancestry chain (e.g. "Assets:Checking"). This is a denormalized field stored in the database for query convenience and human readability (someone inspecting the DB can see the full path directly). It is maintained by application code — any rename or re-parenting triggers an update of `full_name` for the affected account and all its descendants.
- `parent_id`: nullable reference to the parent account. Null means this is a root account.
- `account_type`: one of Asset, Liability, Revenue, Expense, Equity. Only stored on root accounts; child accounts inherit from their root ancestor.
- `created_at`, `updated_at`: timestamps

**Rules:**
- Account names may not contain ":" (since that is the hierarchy separator).
- Deleting an account is only allowed if it has no postings referencing it and no child accounts. Otherwise the user must reassign or delete those first.

### Transaction
- `id`: unique identifier
- `date`: the date of the transaction (YYYY-MM-DD string)
- `description`: short text description
- `status`: one of pending, posted, cleared (free-form, no enforced progression)
- `recurring_template_id`: nullable reference to the recurring transaction template that generated this transaction. Null for manually created transactions.
- `created_at`, `updated_at`: timestamps

**Rules:**
- A transaction must have at least two postings.
- All postings must balance to zero. For single-currency transactions, this means the sum of amounts is zero. For multi-currency transactions, balancing is evaluated by converting all amounts to a common currency using the exchange rate effective on the transaction date.
- Validation is performed on create and update; the backend rejects any transaction that does not balance.

### Posting
- `id`: unique identifier
- `transaction_id`: reference to the parent transaction
- `account_id`: reference to the target account
- `currency`: the currency code (e.g. "USD", "JPY")
- `amount`: signed exact decimal value, stored as text in the database and handled using `shopspring/decimal` in the backend. Examples: "12.34", "-1500.00", "0.00000834".

**Rules:**
- If a transaction has all postings in the same currency, the last posting's amount may be omitted on input and will be auto-inferred as the negated sum of the other postings.
- Multi-currency transactions require all posting amounts to be explicitly provided.

### Tag
Tags are stored using a normalized approach with two tables:
- **Tags table:** stores unique tag key-value pairs.
  - `id`: unique identifier
  - `key`: the tag key (e.g. "trip", "project")
  - `value`: the tag value (e.g. "japan-2026", "kitchen-reno")
- **Transaction-Tags junction table:** maps transactions to tags (many-to-many).
  - `transaction_id`: reference to the transaction
  - `tag_id`: reference to the tag

This avoids duplicating tag strings across rows, enables efficient querying of all unique tags (for the Tags list page), and allows future extension of tag metadata (e.g. description, color).

**Rules:**
- A transaction may have zero or more tags.
- The same key may appear multiple times on a transaction with different values (e.g. a transaction could have both `trip:japan-2026` and `trip:side-trip-osaka`).
- Tags are created and managed inline with the transaction form. When a user enters a tag on a transaction, the backend reuses an existing tag record if one matches, or creates a new one.
- Orphaned tags (not referenced by any transaction) may be cleaned up periodically or left for reuse.

### Budget
- `id`: unique identifier
- `name`: user-assigned name for the budget
- `direction`: spending (ceiling) or saving (floor)
- `period`: weekly, monthly, or yearly
- `currency`: the currency the target amount is denominated in
- `target_amount`: the target amount for the period, stored as text (exact decimal)
- `created_at`, `updated_at`: timestamps

**Related tables:**
- **Budget-Accounts junction table:** maps budgets to their target accounts (many-to-many). Each entry means "this budget tracks this account and all its children."
- **Budget-Tags junction table:** maps budgets to tag filters (many-to-many). Each entry references a tag record. Optional — an empty set means no tag filtering.

**Rules:**
- At least one account must be specified.
- The budget's current value is computed on demand: sum of postings matching the account/tag criteria within the current period, converted to the budget's currency using exchange rates if needed.

### Exchange Rate
- `id`: unique identifier
- `base_currency`: the currency being priced (e.g. "USD")
- `quote_currency`: the currency it's priced in (e.g. "JPY")
- `rate`: the exchange rate (1 unit of base = rate units of quote). Stored as text (exact decimal), same as all monetary values.
- `effective_date`: the date this rate is effective (YYYY-MM-DD)
- `source`: manual or API (to distinguish user-entered from fetched rates)

**Rules:**
- When converting amounts, the system uses the rate with the most recent `effective_date` that is on or before the date in question.
- If no rate exists for a currency pair, the system should check for the inverse pair (e.g. if JPY/USD exists but USD/JPY does not, compute the inverse).
- If no rate exists at all for a required conversion, the system should report an error rather than silently using a default.

### Recurring Transaction Template
- `id`: unique identifier
- `description`: short text description (same as a transaction)
- `status`: default status for generated transactions
- `schedule_type`: the type of recurrence (e.g. "every_n_days", "every_n_weeks", "every_n_months", "monthly_on_day")
- `interval_n`: the N value (e.g. 2 for "every 2 weeks"). Used by all schedule types.
- `day_of_month`: nullable; the specific day of the month (e.g. 15). Only used for monthly-on-day schedules. For "every N months on day D", both `interval_n` and `day_of_month` are used.
- `start_date`: when the schedule begins (YYYY-MM-DD)
- `end_date`: nullable; when the schedule ends (null = indefinite)
- `last_generated_date`: the date of the most recently generated transaction (used to determine what needs to be generated next)
- `is_active`: boolean to allow pausing without deleting
- `created_at`, `updated_at`: timestamps

**Related tables:**
- **Recurring Template Postings table:** stores the template's postings, with the same structure as regular postings (account_id, currency, amount) but referencing the template instead of a transaction.
- **Recurring Template Tags:** templates reference tags via a junction table (similar to transaction-tags), so generated transactions inherit the template's tags.

**Rules:**
- The system generates transactions automatically based on the schedule. On each app start (or on a periodic check), it looks at all active templates, compares `last_generated_date` to the current date, and generates any transactions that are due.
- Generated transactions are normal transactions in every way — they appear in reports, affect balances, etc. They carry a `recurring_template_id` reference back to the template that created them (for traceability), but are otherwise independent and can be edited or deleted without affecting the template.
- If the app is offline for a period, it catches up by generating all missed transactions on next startup.

### Widget Configuration
- `id`: unique identifier
- `widget_type`: one of the supported widget types (account_balance, budget_progress, chart, etc.)
- `position_x`, `position_y`: grid position
- `width`, `height`: grid size
- `config`: JSON blob containing widget-specific configuration (e.g. which account for an account_balance widget, which tag key for a tag_spending_summary widget, time period for charts, etc.)

### Settings
- Stored as key-value pairs in a settings table.
- Keys: `default_currency`, `display_currency`, `available_currencies` (JSON array), `date_format`, `theme`, `fiscal_year_start_month`, `fiscal_year_start_day`.

## API Surface

The backend exposes a service layer with the following operations. In web mode, these are exposed as REST endpoints. In standalone (Wails) mode, they are exposed as direct Go function bindings. The service layer is identical in both cases.

### Account Service
- `ListAccounts(filters, pagination)` → paginated list of accounts
  Supported filters: name/full_name (text match), account type.
  Filtering is the primary mechanism for searching accounts.
- `GetAccount(id)` → account details
- `CreateAccount(name, parent_id, account_type?)` → new account
- `UpdateAccount(id, name?, parent_id?, account_type?)` → updated account
- `DeleteAccount(id)` → success/error
- `SearchAccounts(query)` → fuzzy-matched list (for autocomplete)

### Transaction Service
- `ListTransactions(filters, pagination)` → paginated list with postings and tags
  Supported filters: description (text match), date range, account(s), tag key, tag key+value, status, amount range, currency.
  Filtering is the primary mechanism for searching transactions — there is no separate search endpoint.
- `GetTransaction(id)` → full transaction details with postings and tags
- `CreateTransaction(date, description, status, postings, tags)` → new transaction (validates balance)
- `UpdateTransaction(id, ...)` → updated transaction (re-validates balance)
- `DeleteTransaction(id)` → success/error

### Budget Service
- `ListBudgets()` → all budgets with current computed values
- `GetBudget(id)` → budget details with current period computation
- `CreateBudget(...)` → new budget
- `UpdateBudget(id, ...)` → updated budget
- `DeleteBudget(id)` → success/error

### Exchange Rate Service
- `ListRates(filters)` → exchange rates, filterable by currency pair and date range
- `GetRate(base, quote, date)` → the effective rate for a currency pair on a given date
- `CreateRate(base, quote, rate, effective_date, source)` → new rate
- `DeleteRate(id)` → success/error
- `FetchRatesFromAPI(base, quotes[])` → fetches current rates from external API and stores them

### Recurring Transaction Service
- `ListTemplates()` → all recurring transaction templates
- `GetTemplate(id)` → template details
- `CreateTemplate(...)` → new template
- `UpdateTemplate(id, ...)` → updated template
- `DeleteTemplate(id)` → success/error
- `PreviewSchedule(template, count)` → returns the next N dates that would be generated
- `GenerateDueTransactions()` → checks all active templates and generates any overdue transactions

### Report Service
- `GetBalanceSheet(date_range, period_grouping, depth, tag_filters?)` → balance sheet data (assets/liabilities by account, optionally grouped by period)
- `GetIncomeStatement(date_range, period_grouping, depth, tag_filters?)` → income statement data (revenues/expenses by account, optionally grouped by period)
- `GetNetWorth(date?, comparison_date?)` → net worth with optional trend
- `GetCashFlow(date_range)` → net income minus expenses
- `GetTopSpendingCategories(date_range, limit)` → top N expense accounts
- `GetExpenseBreakdown(date_range, group_by)` → spending breakdown by account or tag
- `GetTagSpendingSummary(tag_key, date_range)` → spending grouped by values of a tag key

### Widget Service
- `ListWidgets()` → all widget configurations for the home page
- `CreateWidget(widget_type, position, size, config)` → new widget
- `UpdateWidget(id, position?, size?, config?)` → updated widget (for dragging, resizing, reconfiguring)
- `DeleteWidget(id)` → success/error

### Settings Service
- `GetSettings()` → all settings
- `UpdateSettings(key, value)` → update a single setting

### Search Service
- `GlobalSearch(query)` → searches across multiple entity types (accounts by name/full_name, transactions by description, tags by key/value). Returns grouped results (e.g. `{ accounts: [...], transactions: [...], tags: [...] }`). Used by the top bar unified search.

# Architectural Design

The architecture is separated into three tiers: **Frontend**, **Backend**, and **Database**. The backend exposes the same service layer through two transport mechanisms: a REST API (for web mode) and Wails Go-to-JS bindings (for standalone desktop mode). The frontend is identical in both modes.

## Backend

### Layers

The backend follows a layered architecture:

1. **Transport layer** — translates external requests into service calls and service responses into external responses.
   - **REST handlers** (web mode): HTTP route handlers using Go 1.22+ `net/http` built-in routing. Handles JSON serialization/deserialization, HTTP status codes, request validation.
   - **Wails bindings** (standalone mode): thin wrapper functions that Wails exposes to the frontend. Calls directly into the service layer.
   - Both transports call the same service methods with the same parameter types.

2. **Service layer** — contains all business logic. One service per domain module:
   - `AccountService`
   - `TransactionService`
   - `BudgetService`
   - `ExchangeRateService`
   - `RecurringTransactionService`
   - `ReportService`
   - `WidgetService`
   - `SettingsService`
   - `SearchService`
   Services may depend on each other (e.g. `TransactionService` uses `ExchangeRateService` to validate multi-currency balancing; `ReportService` uses `ExchangeRateService` for currency conversion). Dependencies are injected via constructor parameters.

3. **Data access layer (DAL)** — generated by **sqlc** from handwritten SQL queries. Each service has a corresponding set of SQL queries in `.sql` files. sqlc generates type-safe Go functions and structs from these queries. This layer is the only layer that touches the database directly.

4. **Model layer** — shared domain types used across the service and transport layers. Provides translation between sqlc-generated types and domain types. All monetary fields use a `StringDecimal` wrapper type (wrapping `shopspring/decimal`) that marshals to/from JSON as strings.

### Libraries and Tools
- **Router (web mode):** Go 1.22+ `net/http` built-in routing (method-based routing and path parameters, no external dependency needed)
- **Desktop framework (standalone mode):** Wails v2
- **Database driver:** `mattn/go-sqlite3` (CGO-based, mature and battle-tested)
- **SQL code generation:** sqlc
- **Schema migrations:** goose (versioned `.sql` migration files)
- **Logging:** Go's built-in `log/slog` (structured logging, no external dependency)
- **Decimal arithmetic:** `shopspring/decimal` for exact decimal math on all monetary values
- **External exchange rate API:** a configurable HTTP client; specific API provider to be determined (e.g. exchangerate.host, Open Exchange Rates, or similar free-tier API)

### Project Structure (Backend)
```
backend/
  cmd/
    web/              # main.go for web server mode
    desktop/          # main.go for Wails desktop mode
  internal/
    service/          # business logic services
      account.go
      transaction.go
      budget.go
      exchange_rate.go
      recurring.go
      report.go
      widget.go
      settings.go
      search.go
    transport/
      rest/            # HTTP handlers and route registration
        account.go
        transaction.go
        ...
        router.go
      wails/           # Wails binding wrappers
        bindings.go
    dal/               # data access layer
      queries/         # handwritten .sql files for sqlc
        account.sql
        transaction.sql
        ...
      generated/       # sqlc-generated Go code (do not edit)
      migrations/      # goose migration .sql files
    model/             # shared domain types and enums
      model.go
      errors.go
  sqlc.yaml            # sqlc configuration
```

## Frontend

### Layers

The frontend follows a three-layer architecture:

1. **API layer** — handles communication with the backend. In web mode, this makes HTTP `fetch` calls to the REST API. In standalone mode, this calls Wails-generated JS bindings. The rest of the frontend is agnostic to which transport is used.

2. **State/data layer** — manages client-side state using SolidJS's reactive primitives (signals, stores, resources). Each domain has a corresponding state module:
   - `useAccounts()`, `useTransactions()`, `useBudgets()`, etc.
   - These use SolidJS `createResource` for data fetching and caching, and `createSignal`/`createStore` for local UI state.
   - Provides the reactive data that components consume.

3. **UI layer** — SolidJS components organized by feature area. Uses reusable shared components for common patterns (forms, lists, modals, etc.).

### Libraries and Tools
- **Framework:** SolidJS
- **Routing:** `@solidjs/router`
- **Build tool:** Vite (with `vite-plugin-solid`)
- **Styling:** plain CSS with nesting (no CSS-in-JS, no preprocessors). CSS custom properties for theming.
- **Charts:** a lightweight charting library compatible with SolidJS (e.g. Chart.js with a thin Solid wrapper, or uPlot for performance)
- **Drag-and-drop (widgets):** a lightweight library or custom implementation for the home page grid
- **Date handling:** `date-fns` for date formatting, parsing, and comparison. Used for display formatting according to the user's date format preference. Never uses JavaScript `Date` for date-only values. A utility function `todayLocalDate()` is used to get the current date as YYYY-MM-DD in local time.
- **Decimal arithmetic:** `big.js` for exact decimal math on the client side (e.g. transaction form balancing preview). Monetary values are never used in arithmetic as JavaScript `number`.

### Project Structure (Frontend)
```
frontend/
  src/
    api/              # API layer (fetch calls or Wails bindings)
      account.ts
      transaction.ts
      ...
      client.ts       # shared HTTP client config (web mode)
      types.ts        # TypeScript interfaces matching backend model
    state/            # reactive state/data layer
      accounts.ts
      transactions.ts
      budgets.ts
      decimal.ts      # exact decimal utilities using big.js
      dates.ts        # date utilities (todayLocalDate, formatDateDisplay, etc.)
      ...
    components/       # reusable UI components
      forms/          # form inputs, date picker, autocomplete, tag input
      lists/          # styled list component with filtering, pagination
      charts/         # chart wrappers
      widgets/        # home page widget components
      layout/         # shell, sidebar, header, breadcrumbs
    pages/            # page-level components (one per route)
      Home.tsx
      Accounts.tsx
      AccountDetail.tsx
      Transactions.tsx
      TransactionDetail.tsx
      Tags.tsx
      TagDetail.tsx
      Budgets.tsx
      BudgetDetail.tsx
      BalanceSheet.tsx
      IncomeStatement.tsx
      ExchangeRates.tsx
      RecurringTransactions.tsx
      RecurringTransactionDetail.tsx
      Settings.tsx
    styles/           # global styles and theme definitions
      themes/         # theme CSS files
        dark.css      # default dark theme ("Neon Night")
      global.css      # reset, base styles, CSS custom properties
    App.tsx           # root component with router
    index.tsx         # entry point
  index.html
  vite.config.ts
  tsconfig.json
  package.json
```

## Database

### Engine
SQLite, stored as a single `.db` file. The user can open this file in any SQLite viewer (e.g. DB Browser for SQLite) to inspect or modify data directly.

### Schema Overview
Tables map directly to the data model entities. Junction tables handle all many-to-many relationships. The schema is managed through goose migration files (versioned, incremental SQL).

**Tables:**
- `accounts` — id, name, full_name (unique), parent_id (FK → accounts.id, nullable), account_type (nullable, only for root accounts), created_at, updated_at
- `transactions` — id, date, description, status, recurring_template_id (FK → recurring_templates.id, nullable), created_at, updated_at
- `postings` — id, transaction_id (FK → transactions.id), account_id (FK → accounts.id), currency, amount (text, exact decimal), created_at
- `tags` — id, key, value (unique constraint on key+value pair)
- `transaction_tags` — transaction_id (FK → transactions.id), tag_id (FK → tags.id), primary key (transaction_id, tag_id)
- `budgets` — id, name, direction, period, currency, target_amount (text, exact decimal), created_at, updated_at
- `budget_accounts` — budget_id (FK → budgets.id ON DELETE CASCADE), account_id (FK → accounts.id)
- `budget_tags` — budget_id (FK → budgets.id ON DELETE CASCADE), tag_id (FK → tags.id)
- `exchange_rates` — id, base_currency, quote_currency, rate (text, exact decimal), effective_date, source, created_at
- `recurring_templates` — id, description, status, schedule_type, interval_n, day_of_month (nullable), start_date, end_date (nullable), last_generated_date (nullable), is_active, created_at, updated_at
- `recurring_template_postings` — id, template_id (FK → recurring_templates.id), account_id (FK → accounts.id), currency, amount (text, exact decimal)
- `recurring_template_tags` — template_id (FK → recurring_templates.id), tag_id (FK → tags.id)
- `widget_configurations` — id, widget_type, position_x, position_y, width, height, config (JSON text), created_at, updated_at
- `settings` — key (primary key), value (text)

### Indexing Strategy
- `accounts.parent_id` — for tree traversal queries
- `accounts.full_name` (unique) — for text search and sorting on the accounts list page
- `transactions.date` — for date range queries in reports and list filtering
- `transactions.status` — for status filtering
- `transactions.recurring_template_id` — for reverse lookup of generated transactions
- `postings.transaction_id` — for loading postings with their transaction
- `postings.account_id` — for account-based queries (account page, reports, budgets)
- `tags.key`, `tags.key + tags.value` — for tag filtering
- `transaction_tags.transaction_id`, `transaction_tags.tag_id` — for junction lookups
- `exchange_rates.base_currency + exchange_rates.quote_currency + exchange_rates.effective_date` — for rate lookups
- `recurring_templates.is_active` — for the generation check on startup

# UI/UX

## App Shell Layout

The app uses a persistent two-part layout:

**Sidebar (left):**
- Fixed-width left sidebar, always visible.
- Contains the primary navigation links, ordered top to bottom:
  - Home
  - Accounts
  - Transactions
  - Budgets
  - Reports (expandable submenu: Balance Sheet, Income Statement)
  - Tags
  - Recurring Transactions
  - Exchange Rates
  - (spacer)
  - Settings
- Each link has an icon and a text label.
- The sidebar can be collapsed to icon-only mode to save space.
- The active page is visually highlighted.

**Top bar (horizontal, spans the content area to the right of the sidebar):**
- **Left:** Breadcrumbs showing the current navigation path (e.g. `Accounts > Assets > Checking`). Each breadcrumb segment is clickable.
- **Center/Right:** Global search bar. As the user types, a dropdown appears showing grouped results (Accounts, Transactions, Tags). Clicking a result navigates to its detail page. Pressing Enter navigates to a full search results page.
- **Right:** Quick action button(s) — at minimum a "New Transaction" button for fast entry without navigating away.

**Content area:** fills the remaining space to the right of the sidebar and below the top bar. All page content renders here.

## Theming

- Themes are implemented using **CSS custom properties** (variables) defined at the `:root` level.
- Each theme is a separate CSS file that overrides the custom properties.
- The default theme is a **dark theme**.
- Switching themes applies a different CSS file, which changes all custom properties and therefore all component colors.
- New themes can be added by creating a new CSS file following the same variable structure.

### Default Dark Theme — "Neon Night"

The design aesthetic is a modern dark UI with bright "neon sign" accent colors on a dark background.

**Backgrounds:**
| Token | Hex | Usage |
|-------|-----|-------|
| `--color-bg-primary` | `#1a1a2e` | Page background |
| `--color-bg-surface` | `#22223a` | Cards, rows, sidebar, modals |
| `--color-bg-elevated` | `#2a2a4a` | Hover states, active items, input fields |
| `--color-border` | `#333366` | Subtle borders |

**Text:**
| Token | Hex | Usage |
|-------|-----|-------|
| `--color-text-primary` | `#e0e0f0` | Primary body text |
| `--color-text-secondary` | `#8888aa` | Parent paths, timestamps, faint indicators |
| `--color-text-disabled` | `#555577` | Disabled/inactive elements |

**Primary accent (blue):**
| Token | Hex | Usage |
|-------|-----|-------|
| `--color-primary` | `#4d8eff` | Links, active nav, primary buttons |
| `--color-primary-hover` | `#6aa3ff` | Hover state for primary elements |
| `--color-primary-muted` | `#2a4a8a` | Selected/active backgrounds |

**Neon accent palette:**
Bright, saturated colors designed to pop against the dark background — used for charts, badges, and visual indicators.
| Token | Name | Hex | Badge / Usage |
|-------|------|-----|---------------|
| `--color-neon-green` | Neon green | `#39ff85` | Asset badge, positive amounts, saving goals on track |
| `--color-neon-red` | Neon red | `#ff3b5c` | Liability badge, negative amounts, over-budget |
| `--color-neon-cyan` | Neon cyan | `#22d3ee` | Revenue badge |
| `--color-neon-yellow` | Neon yellow | `#ffe03d` | Expense badge, pending indicator, warnings |
| `--color-neon-magenta` | Neon magenta | `#f472b6` | Equity badge |
| `--color-neon-violet` | Neon violet | `#a78bfa` | Tag badges |
| `--color-neon-blue` | Neon blue | `#3b82f6` | Chart series |
| `--color-neon-orange` | Neon orange | `#ff8c42` | Chart series |

**Semantic colors:**
| Token | Hex | Usage |
|-------|-----|-------|
| `--color-success` | `#39ff85` | Success states |
| `--color-warning` | `#ffe03d` | Warnings, approaching limits |
| `--color-error` | `#ff3b5c` | Errors, destructive actions |
| `--color-info` | `#4d8eff` | Informational |

**Chart palette order:** Charts cycle through colors in this order for distinct visual separation:
`--color-primary`, `--color-neon-green`, `--color-neon-cyan`, `--color-neon-yellow`, `--color-neon-magenta`, `--color-neon-violet`, `--color-neon-blue`, `--color-neon-orange` (8 colors before repeating).

**Typography:**
| Token | Value | Usage |
|-------|-------|-------|
| `--font-family` | `Inter, -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif` | All UI text |
| `--font-family-mono` | `'JetBrains Mono', 'Fira Code', monospace` | Monetary amounts, numbers |

**Visual properties:**
| Token | Value | Usage |
|-------|-------|-------|
| `--radius-sm` | `4px` | Inputs, small elements |
| `--radius-md` | `6px` | Cards, buttons |
| `--radius-lg` | `12px` | Badges, chips, pills |
| `--shadow-card` | `0 2px 8px rgba(0, 0, 0, 0.3)` | Card elevation |

## Page Descriptions

### Home Page
- A grid of draggable, resizable widgets.
- An "Edit mode" toggle to enter layout editing (add/remove/resize/reposition widgets). In view mode, widgets are static.
- When adding a widget, a picker dialog shows available widget types. After selecting a type, a configuration dialog appears for widget-specific settings (e.g. which account for an account balance widget, time range for a chart).

### Accounts List Page
- A styled list of account rows (not a spreadsheet-style table).
- Each row shows: parent path as small muted text above, short account name as primary text below, colored account type badge, current balance, currency.
- Account type badge colors: Asset → neon green, Liability → neon red, Revenue → neon cyan, Expense → neon yellow, Equity → neon magenta.
- Supports: text search (by name), filtering by account type, sorting options, pagination.
- Each row has action buttons: view (navigates to detail page), edit (opens edit form/modal), delete (with confirmation).
- A "New Account" button above the list opens the account creation form.

### Account Detail Page
- Header: full account name with colored type badge, parent account (clickable link).
- Summary section: current balance, currency breakdown if multi-currency.
- Transactions list: styled rows for all transactions that have at least one posting to this account, sorted by date descending. Each row shows date, description, faint pending/posted indicator (hidden if cleared), and the posting amount relevant to this account. Each row is clickable to navigate to the transaction detail page.
- The transactions list supports the same filtering/search/pagination as the main Transactions list page.

### Transactions List Page
- A styled list of transaction rows (not a spreadsheet-style table).
- Each row shows a compact summary: date, description, and a faint pending/posted indicator (only shown if the transaction is not cleared — cleared transactions show no indicator). Action buttons for view, edit, delete.
- Supports: text search (by description), filtering by date range, account, tag (key or key+value), status, amount range, and currency. Sorting options. Pagination.
- Filtering changes which rows are displayed; the visual style remains the same.
- A "New Transaction" button above the list opens the transaction creation form.

### Transaction Detail Page
- Header: date, description, status (with an inline status selector to change it).
- Tags section: displays all tags as chips/badges (`key:value`). Each is clickable to navigate to the tag page.
- Postings list: lists all postings with: account (shown as parent path in small muted text above, short name below, with colored type badge — clickable link), currency, amount. Debit/credit amounts are visually distinguished (e.g. color-coded).
- Edit button to modify the transaction. Delete button with confirmation.
- If generated from a recurring template: a link to the source template.

### Transaction Creation/Edit Form
- Date: date picker input.
- Description: text input.
- Status: dropdown (pending, posted, cleared).
- Tags: a tag input component — user types `key:value` and it appears as a chip. Existing tags are suggested via autocomplete as the user types.
- Postings: a dynamic list of posting rows. Each row has:
  - Account: autocomplete input with fuzzy matching against existing accounts. Allows typing a new account name that doesn't exist yet (will be created on save).
  - Currency: dropdown populated from the available currencies setting.
  - Amount: numeric input. The last row's amount can be left blank for auto-inference (single-currency only).
  - A remove button per row, and an "Add Posting" button to add more rows.
- Save and Cancel buttons. Save validates balance before submitting.
- Balance validation uses exact decimal arithmetic (big.js), not JavaScript floating-point.

### Budgets List Page
- Shows all budgets as styled cards in a responsive grid.
- Each budget card shows: name, direction (spending/saving), period, target amount, current amount, progress bar/indicator.
- Visual distinction for budgets approaching or exceeding their target (spending) or nearing their goal (saving).
- Action buttons: view, edit, delete.
- "New Budget" button.

### Budget Detail Page
- Header: budget name, direction, period, currency, target amount.
- Progress visualization: progress bar or gauge showing current vs. target, color-coded by status (under/near/over for spending; progress toward goal for saving).
- Target definition: lists the tracked accounts and tag filters.
- Transaction breakdown: styled list of matching transactions in the current period, so the user can see what contributed to the budget total.

### Balance Sheet Report Page
- Controls at top: date range selector, period grouping (weekly/monthly/yearly), depth slider/selector, tag filter input.
- Main content: a table showing the account hierarchy for Asset and Liability accounts, with period columns showing balances. Totals row at the bottom.
- Toggle between table view and chart view (line chart or stacked bar chart showing balances over time).
- Account names in the table are clickable links to account detail pages.

### Income Statement Report Page
- Same layout and controls as the Balance Sheet, but for Revenue and Expense accounts.
- Shows income, expenses, and net income (income minus expenses) per period.
- Toggle between table and chart views.

### Tags List Page
- A styled list of tag rows.
- Each row shows: key, value, transaction count (how many transactions use this tag).
- Filterable by key or value text search. Sortable. Pagination.
- Each row is clickable to navigate to the tag detail page.

### Tag Detail Page
- Header: tag key and value.
- Transactions list: styled rows for all transactions that carry this tag, with the same layout and features as the main Transactions list page. Each row is clickable.

### Recurring Transactions Page
- A styled list of all recurring templates.
- Each row shows: description, schedule summary (human-readable, e.g. "Every month on the 15th"), next occurrence date, status (active/paused), last generated date.
- Action buttons: view/edit, pause/resume, delete.
- "New Recurring Transaction" button.

### Recurring Transaction Detail/Edit Form
- Same fields as the transaction form (description, status, tags, postings), plus:
- Schedule configuration: schedule type dropdown, interval input, day-of-month input (conditional), start date picker, optional end date picker.
- Preview section: shows the next 5-10 dates that would be generated, updating live as the user changes schedule settings.

### Exchange Rates Page
- A styled list of all exchange rates.
- Each row shows: base currency, quote currency, rate, effective date, source badge (manual/API).
- Filterable by currency pair and date range.
- "Add Rate" button for manual entry.
- "Fetch Rates" button that opens a dialog to select a base currency and target currencies, then fetches from the external API.

### Settings Page
- A form with sections for each setting group:
  - **Currency:** default currency dropdown, display currency dropdown, available currencies multi-select.
  - **Display:** date format dropdown (e.g. YYYY-MM-DD, MM/DD/YYYY, DD/MM/YYYY), fiscal year start (month and day selectors).
  - **Appearance:** theme selector with preview.
- Save button. Changes take effect immediately or on save.

## Shared Components

The following reusable components appear across multiple pages:

- **Styled list:** reusable list component with styled rows, search bar, filter controls, sorting options, pagination, and per-row action buttons. Each page provides a `renderRow` function that defines the visual layout of each row; the list infrastructure (filtering, pagination, etc.) is shared. Used by Accounts, Transactions, Tags, Exchange Rates, Recurring Transactions, and sub-lists on detail pages.
- **Report table:** tabular component (actual HTML table) used only for report pages (Balance Sheet, Income Statement) where showing account hierarchies with period columns requires a grid layout.
- **Autocomplete input:** text input with fuzzy-matched dropdown suggestions. Used for account selection in posting forms and tag input.
- **Date picker:** calendar-based date selection. Used in transaction forms, report filters, recurring template forms. Outputs a plain `YYYY-MM-DD` string, never a `Date` object.
- **Tag input:** specialized input for entering key:value tags as chips, with autocomplete against existing tags. Used in transaction forms and filter controls.
- **Progress indicator:** bar or gauge for budget progress display. Used in budget widgets and budget detail pages.
- **Chart components:** wrappers around the charting library for line charts, bar charts, pie/donut charts. Used in reports and chart widgets.
- **Confirmation dialog:** modal dialog for confirming destructive actions (delete). Used throughout.
- **Form controls:** standard styled inputs (text, number, dropdown/select, multi-select) consistent with the theming system.

# Testing

Development follows **red-green TDD**: write a failing test first, implement the minimum code to pass it, then refactor. Tests are written before the code they test.

## Philosophy
- Tests should provide good coverage with a **manageable number of tests**. Avoid generating excessive tests that become a maintenance burden.
- **One test per meaningful behavior.** A simple function with one code path needs one test. A function with several distinct behaviors (happy path, validation error, edge case) gets one test per behavior. The goal is that each test has a clear, specific name describing the single behavior it verifies.
- Each test should test only the behavior it is named for — no multi-purpose tests.
- Tests are organized by module/domain to match the codebase structure, making it easy to find and maintain tests for a given feature.

## Backend Testing

### Tools
- **Test framework:** Go's built-in `testing` package. No external assertion libraries. Helper functions should be used to produce clear, descriptive failure messages (e.g. `assertDecimalEqual(t, expected, actual)` that formats the diff readably).
- **Database:** In-memory SQLite for all tests that touch the database. Each test (or test suite) gets a fresh in-memory database with migrations applied. This tests the full path from service → DAL → database without mocks.

### Test Organization
```
backend/
  internal/
    service/
      account_test.go        # tests for AccountService
      transaction_test.go    # tests for TransactionService
      budget_test.go         # tests for BudgetService
      exchange_rate_test.go  # tests for ExchangeRateService
      recurring_test.go      # tests for RecurringTransactionService
      report_test.go         # tests for ReportService
      widget_test.go         # tests for WidgetService
      settings_test.go       # tests for SettingsService
      search_test.go         # tests for SearchService
      integration_test.go    # end-to-end workflow tests
    transport/
      rest/
        account_test.go      # tests for HTTP handler request/response
        transaction_test.go
        ...
    testutil/
      db.go                  # helper to create fresh in-memory DB with migrations
```

### What to Test

**Service layer (primary focus):**
Each service method gets tests covering:
- Happy path (valid input produces correct result).
- Validation failures (invalid input is rejected with appropriate error).
- Edge cases specific to the domain logic.

Key areas with more thorough coverage:
- **Transaction balancing:** single-currency balance check, multi-currency balance check using exchange rates, auto-infer of last posting, rejection of unbalanced transactions, rejection of auto-infer on multi-currency.
- **Budget computation:** correct summing of matching postings, account children inclusion, tag filter narrowing, period boundary handling, currency conversion.
- **Exchange rate lookup:** most recent rate on-or-before date, inverse pair fallback, error when no rate exists.
- **Recurring transaction generation:** correct date calculation for each schedule type, catch-up generation after offline period, end date enforcement, pause/resume behavior.
- **Account hierarchy:** full_name computation and sync on rename/re-parent, type inheritance, deletion constraints.

**Transport layer (lighter coverage):**
- Test that HTTP handlers correctly deserialize requests, call the right service method, and serialize responses.
- Test error responses (400 for validation errors, 404 for not found, etc.).
- These tests can use a mock or stub service layer since the service logic is already tested.

### What Not to Test
- sqlc-generated code (it's generated; test the queries indirectly through service tests).
- Trivial getters/setters or struct construction.

## Frontend Testing

### Tools
- **Test framework:** Vitest (compatible with Vite, fast).
- **Component testing:** `@solidjs/testing-library` for rendering SolidJS components and interacting with them.
- **Assertions:** Vitest's built-in `expect` API.

### Test Organization
```
frontend/
  src/
    state/
      decimal.test.ts        # tests for decimal utility functions
      dates.test.ts          # tests for date utility functions
    components/
      forms/
        TransactionForm.test.tsx
        TagInput.test.tsx
        DatePicker.test.tsx
      widgets/
        BudgetProgress.test.tsx
```

### What to Test

**State/data layer:**
- Correct transformation of API responses into reactive state.
- Client-side decimal arithmetic (e.g. auto-infer balancing amount uses exact decimal math, not floating-point).
- Date utilities: `todayLocalDate()` returns YYYY-MM-DD, `formatDateDisplay()` handles dates without timezone shifting.

**Components with non-trivial logic:**
- **Transaction form:** adding/removing posting rows, auto-infer of last posting amount using big.js, balance validation feedback.
- **Budget progress display:** correct progress calculation and visual state (under/near/over for spending, progress for saving).
- **Date picker:** outputs `YYYY-MM-DD` string, not a `Date` object.
- **Tag input:** adding/removing key:value chips, autocomplete behavior.

### What Not to Test
- Pure layout/display components with no logic (e.g. a styled card wrapper).
- Page components that are just composition of tested sub-components.
- Routing (trust the router library).

## Integration Tests

Integration tests cover end-to-end workflows through the backend (service + DAL + database), verifying that multi-step operations produce correct results.

Key integration workflows:
- **Create accounts → create transaction → verify balances:** create an account hierarchy, post a transaction with multiple postings, verify account balances reflect the postings correctly.
- **Multi-currency transaction flow:** create exchange rates, create a multi-currency transaction, verify it balances correctly using the rates.
- **Budget lifecycle:** create accounts, create a budget, post transactions, verify budget computation reflects the transactions.
- **Recurring transaction generation:** create a template, advance time, trigger generation, verify correct transactions are created.
- **Report accuracy:** create a set of transactions spanning multiple periods, verify balance sheet and income statement produce correct aggregated values.
- **Tag filtering in reports:** create transactions with tags, run reports with tag filters, verify filtered vs. unfiltered results differ correctly.
- **Widget CRUD:** create, list, update, delete widgets.
- **Settings persistence:** update settings, verify they persist across reads.

Integration tests use the same in-memory SQLite setup as service tests. They are kept in a separate `integration_test.go` file to distinguish them from unit tests.
