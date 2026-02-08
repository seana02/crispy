# Functional Requirements

## Accounts & Transactions
### Account Management
- [ ? ] Create, edit, and delete accounts
- [ ? ] Support for nested/subaccounts (e.g., Assets:Bank:Checking)
- [ ? ] Ability to mark accounts as active/inactive
- [ ? ] Account types: Asset, Liability, Income, Expense, Equity
- [ ? ] Each account has a currency
### Transactions
- [ ? ] Double-entry bookkeeping: Each transaction must balance
- [ ? ] Support for multiple postings per transaction
- [ ? ] Support for multiple currencies in a transaction
- [ ? ] Stores exchange rates
- [ ? ] Transaction metadata: date, description, tags, notes, attachments
- [ ? ] Transaction import/export (CSV, OFX, QIF, Ledger)
- [ ? ] Search and filter transactions by account, tag, amount, date range, etc
### Recurring Transactions
- [ ? ] Create, edit, delete recurring transactions
- [ ? ] Flexible scheduling (cron-like syntax, daily/weekly/monthly/yearly options)
- [ ? ] Auto-generation of instances with optional user confirmation (auto-insert transaction, or mark as "pending" and require user to confirm)
- [ ? ] Smart handling of skipped periods
### Automation Rules
- [ ? ] Define rules to automatically modify or tag transactions
- Examples:
    - [ ? ] Auto-assign category based on payee or memo
    - [ ? ] Auto-calculate cashback, fees, or interest and create postings
    - [ ? ] Auto-split transactions into multiple subaccounts

## Reporting & Analysis
### Standard Reports
- [ ? ] Balance sheet
- [ ? ] Income statement / Profit & Loss
- [ ? ] Cash flow report
- [ ? ] Account register view
- [ ? ] Transaction history by tag/category
- [ ? ] Reports can be generated in a base currency, using historical rates if necessary
### Custom Reports
- [ ? ] Filterable by date ranges, accounts, tags, payees
- [ ? ] Aggregation of subaccounts
- [ ? ] Exportable to CSV, PDF
### Dashboard
- [ ? ] Overview of balances, recent transactions, budget progress, upcoming recurring transactions
- [ ? ] Visualizations: charts for spending trends, income, net worth over time

## Budgeting & Alerts (Optional Early Features)
- [ ? ] Track planned spending vs. actuals
- [ ? ] Set budgets for accounts or tags
- [ ? ] Notify users when approaching/exceeding budgets

## User Experience / Quality of Life
- [ ? ] Keyboard shortcuts for transaction entry
- [ ? ] Auto-complete for accounts, tags, and payees
- [ ? ] Templates/presets for common transaction types
- [ ? ] Undo/redo functionality
- [ ? ] Support for attachments (e.g., receipts)
- [ ? ] Auto-fetch external rates from an external service

## Import/Export & Interoperability
- [ ? ] Ledger/hledger-style text import/export
- [ ? ] CSV/Excel import/export
- [ ? ] Potential API for programmatic access (for future integrations)


# Non-Functional Requirements
## Performance
- [ ? ] Handle tens of thousands of transactions efficiently on SQLite
## Security
- [ ? ] Optional encryption of the database
- [ ? ] Local-first design; optional cloud sync in the future
## Reliability
- [ ? ] Safe handling of database writes (transactions should be atomic)
## Usability
- [ ? ] Clean, responsive UI with low learning curve
## Portability
- [ ? ] Cross-platform desktop/web support

# Technical Requirements
## Backend
- [ ? ] Golang for API/server logic
- [ ? ] SQLite database for local storage
- [ ? ] RESTful or GraphQL API
- [ ? ] Transaction handling logic with double-entry validation
- [ ? ] Scheduler/cron engine for recurring transactions
## Frontend
- [ ? ] SolidJS (or similar reactive framework like React/Vue/Svelte)
- [ ? ] Responsive, modern UI design
- [ ? ] Client-side caching for fast access to accounts/transactions
- [ ? ] Charts and visualizations
## Data Model (rough draft)
- [ ? ] Account: id, parent_id, name, type, description, active
- [ ? ] Account: supports nested structure and multi-currency
- [ ? ] Transaction: id, date, description, payee, metadata, recurrence_id
- [ ? ] Posting: id, transaction_id, account_id, amount, currency
- [ ? ] Rule: id, condition, action
- [ ? ] Recurring Transaction: id, template_transaction_id, cron_expression, next_run_date

## Future/Optional Features
- [ ? ] Cloud sync / multi-device support
- [ ? ] Mobile app
- [ ? ] Advanced AI-powered transaction tagging and recommendations
- [ ? ] Collaborative/family budgeting


# Database Schema
## Account
|Column|Type|Notes|
|:--|:--|:--|
|id|INTEGER PK|Unique account ID|
|parent_id|INTEGER FK|Nullable, self-reference for subaccounts|
|name|TEXT|Account name|
|type|TEXT|ENUM: Asset, Liability, Revenue, Expense, Equity|
|currency|TEXT|ISO currency code|
|description|TEXT|optional desc|
|active|BOOLEAN|account active or not|
|date_created|TEXT|timestamp|
|date_updated|TEXT|timestamp|
CREATE INDEX idx_account_parent ON Account(parent_id); -> speed up subaccount queries
(parent_id, name) must be unique

## Transaction
|Column|Type|Notes|
|:--|:--|:--|
|id|INTEGER PK|unique transaction ID|
|description|TEXT|memo/desc|
|date|TEXT|transaction date|
|status|TEXT|ENUM: pending, posted, cleared; see below|
|reference_id|INTEGER FK|nullable, links to a recurring transaction template|
|date_created|TEXT|timestamp|
|date_updated|TEXT|timestamp|
pending: not finalized, maybe user needs to approve, or awaiting bank posting
posted: recorded in ledger, affects balances
cleared: matches the external statement (bank/credit card)
-> reconciliation
Index Transaction(date) for chronological queries
can use SQLite Recursive CTEs (using WITH clause probably) to do recursive queries

## Posting
|Column|Type|Notes|
|:--|:--|:--|
|id|INTEGER PK|unique posting ID|
|transaction_id|INTEGER FK|ref to transactions|
|account_id|INTEGER FK|ref to accounts|
|amount|TEXT|(decimal) amount (positive or negative)|
|currency|TEXT|ISO currency code|
|date_created|TEXT|timestamp|
|date_updated|TEXT|timestamp|
note: sum of all postings in a transaction must sum to 0 (in a currency)
Indexes Postings(transaction_id), Postings(account_id) for fast reports

## RecurringTransaction
|Column|Type|Notes|
|:--|:--|:--|
|id|INTEGER PK|unique posting ID|
|description|TEXT|memo|
|recurrence_id|INTEGER FK|id of recurrence|
|last_run|TEXT|last time this template was executed, default: '1970-1-1'|
|active|BOOLEAN|enable/disable template|
|status|TEXT|ENUM: pending, posted, cleared, status to create as|
|date_created|TEXT|timestamp|
|date_updated|TEXT|timestamp|

## RecurringPosting
|Column|Type|Notes|
|:--|:--|:--|
|id|INTEGER PK|unique recurring posting ID|
|template_id|INTEGER FK|ref to recurring transaction template|
|account_id|INTEGER FK|ref to accounts|
|amount_expression|TEXT|value like "50.00" or formula like "amount*0.01"|
|currency|TEXT|ISO currency code|
|date_created|TEXT|ISO8601 timestamp|
when generating a transaction, evaluate amount_expression to compute actual posting amounts

## Recurrence
|Column|Type|Notes|
|:--|:--|:--|
|id|INTEGER PK|unique recurrence ID|
|frequency|TEXT NOT NULL|must be in 'daily', 'weekly', or 'monthly'|
|interval|INTEGER NOT NULL|x in "every x days/weeks/months"|
|day_of_week|INTEGER|must be 0-6 if 'weekly', or NULL|
|day_of_month|INTEGER|must be 1-31 if 'monthly', or NULL|
|start_date|TEXT|first date that is triggered|
|end_date|TEXT|last possible date to trigger|

## Budget
|Column|Type|Notes|
|:--|:--|:--|
|id|INTEGER PK|unique budget ID|
|name|TEXT|label for the budget|
|description|TEXT|description or additional notes|
|target|TEXT|the target amount for budget|
|tag|INTEGER FK|associated tag for expense tracking|
|reset_frequency|INTEGER FK|id of recurrence for when the budget is reset|

## Rule
rules for auto-tagging, posting adjustments, or cashback calculations
|Column|Type|Notes|
|:--|:--|:--|
|id|INTEGER PK|unique rule ID|
|name|TEXT|human-readable name, unique|
|condition|TEXT|e.g. payee = "starbucks"|
|action|TEXT|e.g. add_tag("coffee")|
|enabled|BOOLEAN|true if rule is active|
|date_created|TEXT|timestamp|
|date_updated|TEXT|timestamp|
condition/action as JSON or small DSL (Domain-Specific Lang)
parse and check each, apply action if condition matches
evaluated when transaction is inserted or updated
use a (tiny) expression parser library, e.g. Knetic/govaluate

## Tag
*:* relation between transactions and tags
|Column|Type|Notes|
|:--|:--|:--|
|id|INTEGER PK|tag ID|
|name|TEXT UNIQUE|tag name|

## Transaction_Tag
associative table linking transactions and tags
|Column|Type|Notes|
|:--|:--|:--|
|transaction_id|INTEGER FK|transaction ref|
|tag_id|INTEGER FK|tag ref|
Index Transaction_Tags(tag_id) for filtered tag reports

## Attachment
|Column|Type|Notes|
|:--|:--|:--|
|id|INTEGER PK|attachment ID|
|transaction_id|INTEGER FK|transaction ref|
|filename|TEXT|name of file|
|path|TEXT|nullable, copy file to location controlled by app (e.g. application directory), then store path|
|mime_type|TEXT|e.g. image/png|
|date_created|TEXT|timestamp|

## Exchange Rate
|Column|Type|Notes|
|:--|:--|:--|
|id|INTEGER PK|unique ID|
|from_currency|TEXT|ISO code|
|to_currency|TEXT|ISO code|
|rate|TEXT|(decimal) conversion rate|
|date|TEXT|date of rate|
|source|TEXT|optional (manual or API)|

## Optional Utility Tables
- audit logs (for transaction changes)
    - CRUDing anything in the database

# ER Diagram by ChatGPT
Key Legend:
→ = foreign key reference
< = one-to-many direction (parent → child)
(1:N) = cardinality from parent to child
PK = primary key

Account (id PK)
├─ parent_id → Account.id (1:N self-reference, subaccounts)
├─ name
├─ type
├─ currency
├─ description
├─ active
├─ created_at
└─ updated_at
    ├─< Posting.account_id (1:N)
    └─< RecurringPosting.account_id (1:N)

Transaction (id PK)
├─ description
├─ date
├─ status (Pending, Posted, Cleared)
├─ created_at
└─ updated_at
    ├─< Posting.transaction_id (1:N)
    ├─< TransactionTag.transaction_id (1:N)
    └─< Attachment.transaction_id (1:N)

Posting (id PK)
├─ transaction_id → Transaction.id
├─ account_id → Account.id
├─ amount
├─ currency
└─ created_at

RecurringTransaction (id PK)
├─ description
├─ recurrence_id → Recurrence.id
├─ last_run
├─ base_amount
├─ active
├─ status (Pending, Posted, Cleared)
├─ created_at
└─ updated_at
    │
    └─< RecurringPosting.recurring_transaction_id (1:N)

RecurringPosting (id PK)
├─ recurring_transaction_id → RecurringTransaction.id
├─ account_id → Account.id
├─ amount_expression
├─ currency
└─ created_at

Recurrence (id PK)
├─ frequency (Daily, Weekly, Monthly)
├─ interval
├─ day_of_week (0-6)
├─ day_of_month (1-31)
├─ start_date
├─ end_date
├─ account_id → Account.id
├─ amount_expression
├─ currency
└─ created_at

Recurrence (id PK)
├─ name
├─ description
├─ target
├─ tag → Tag.id
├─ reset_frequency → Recurrence.id

Rule (id PK)
├─ name
├─ condition
├─ action
├─ enabled
├─ created_at
└─ updated_at
    └─applied to transactions via logic (no direct FK)

Tag (id PK)
├─ name
├─ created_at
└─ updated_at
    │
    └─< TransactionTag.tag_id (1:N)

TransactionTag (transaction_id PK, tag_id PK)
├─ transaction_id → Transaction.id
└─ tag_id → Tag.id

Attachment (id PK)
├─ transaction_id → Transaction.id
├─ path
├─ filename
├─ mime_type
└─ created_at

ExchangeRate (id PK)
├─ from_currency
├─ to_currency
├─ rate
├─ date
└─ source


# Basic File Structure
internal/
├── app/           # Application bootstrap & dependency wiring
    ├─ app.go          # Init app (entry point for main.go), literally just starts
    ├─ config.go       # Loads and validates configuration into a Config struct
    ├─ dependencies.go # Dependency injection & wiring, creates the App struct with all the dependencies injected
    └─ lifecycle.go    # Graceful shutdown, signals, etc.
├── domain/        # Core business entities and rules (pure logic)
    ├─ account.go           # No SQL, no HTTP, no frameworks- pure Go structs and methods to serve as API contract
    ├─ transaction.go
    ├─ posting.go
    ├─ recurring.go
    ├─ budget.go
    ├─ recurrence.go
    ├─ rule.go
    ├─ tag.go
    ├─ exchange_rate.go
    └─ errors.go
├── db/            # Data access, migrations, persistence logic
    ├─ sqlite.go               # DB initializations, migrations
    ├─ migrations/             # SQL migration files, see golang-migrate/migrate
    ├─ repository/             # Defines the interfaces for the various queries/commands that can be performed
    │  ├─ account_repo.go
    │  ├─ transaction_repo.go
    │  ├─ posting_repo.go
    │  ├─ recurring_repo.go
    │  ├─ rule_repo.go
    │  ├─ budget_repo.go
    │  ├─ tag_repo.go
    │  ├─ exchange_repo.go
    └─ seed.go                 # optional demo data loader
├── scheduler/     # Handles recurring transactions; each RecurringTransaction spawns a job, jobs evaluate amount_expression using library
    ├─ manager.go     # Checks and generates pending recurring transactions
    ├─ job_runner.go  # Executes recurring templates
    └─ evaluator.go   # Evaluates amount expressions
├── automation/         # Rule engine for automation
    ├─ engine.go    # core rule evaluator
    ├─ condition.go # parses & checks rule conditions
    ├─ action.go    # defines actions (e.g. add_tag, add_posting); actions are Go functions implementing an interface "type Action interface { Apply(*domain.Transaction) error }"
    ├─ context.go   # transcation/posting context passed into rules - data passed into condition/action so action knows what transaction they're modifying
    ├─ cashback.go
    ├─ tagging.go
    └─ templates.go
├── currency/      # Multi-currency logic, conversions, formatting, basically single function files following SRP
    ├─ exchange_service.go # interacts with ExchangeRate table, provides APIs to look up/update exchange rates; handles fallback (e.g. use latest rate if date not found)
    ├─ precision.go        # defines how many decimal places each currency uses "func PrecisionFor(currency string) int"
    ├─ converter.go        # performs conversions between currencies safely using "shopspring/decimal"; may use exchange_service under the hood
    └─ formatter.go        # handles display formatting, possible expansion for locale support
└── utils/         # Common helpers (config, logging, time, etc.)
    ├─ logger.go
    ├─ timeutil.go
    ├─ errors.go   # wrapping, formatting, logging errors, not related to business errors, basically error helpers
    └─ fileutil.go


# Basic Build Plan
1. build frontend (svelte, shadcn-svelte, wails)
2. embed frontend static files and SQL migration files into go
3. build go app





# Frontend Requirements
## Example Color Palettes
"Professional Dark"
| Role              | Color                   |
| ----------------- | ----------------------- |
| Primary           | `#2A9D8F` (teal green)  |
| Secondary         | `#E9C46A` (warm yellow) |
| Background (dark) | `#1E1E1E`               |
| Surface           | `#242424`               |
| Text (primary)    | `#EAEAEA`               |
| Text (secondary)  | `#A9A9A9`               |
| Accent (error)    | `#E76F51`               |
| Accent (info)     | `#4ECDC4`               |

"Sleek Neutral"
| Role              | Color                   |
| ----------------- | ----------------------- |
| Primary           | `#3B82F6` (modern blue) |
| Secondary         | `#F59E0B` (amber)       |
| Background (dark) | `#121212`               | light: #F9FAFB
| Surface           | `#1F2937`               | light: #FFFFFF
| Text (primary)    | `#F9FAFB`               | light: #1F2937
| Text (secondary)  | `#9CA3AF`               | light: #4B5563
| Accent (success)  | `#10B981`               |
| Accent (error)    | `#EF4444`               |
| Accent (info)     | `#4ECDC4`               |

### Usage
- `--color-bg`: outermost background
    - eg: page body bg
- `--color-surface`: "floating" bg elements (cards, panels, modals) that sit on top of bg
    - eg: dashboard widgets, card panels, tables
- `--color-text-primary`: default text color for main content
    - eg: titles, main labels, table content
- `--color-text-secondary`: less emphasized or supporting text
    - eg: sub-labels, descriptions, timestamps, placeholder text
- `--color-primary`: brand or action color
    - eg: buttons active sidebar item, links, progress bars
- `--color-secondary`: supporting accent color used sparingly
    - eg: highlights, chart accents, secondary buttons
- `--color-accent-*` (info/success/error): contextual semantic colors to represent meaning
    - eg: info: tips, hover hints, charts

## Example Fonts
Primary:
- [Inter](https://rsms.me/inter/)
- [IBM Plex Sans](https://fonts.google.com/specimen/IBM+Plex+Sans)

- Headings: 600-700 weight
- Body: 400-500
- Number/Charts: Use monospaced variant (like IBM Plex Mono)

## Example Icon Sets
- [Lucide](https://lucide.dev/)
- [Tabler Icons](https://tabler.io/icons)
- [Phosphor Icons](https://phosphoricons.com/)
- implement via Solid's' JSX wrapper to dynamically recolor icons using currentColor

## Navigation
- Collapsible sidebar on the left
- Icons + labels
- When collapsed, icons only
- Active item highlight
- Dashboard, Transactions, Budgets, Exchange Rates, Recurring, Rules, Settings
- Header Bar:
    - optional top bar for search
    - quick add ("+" button)
    - theme toggle
    - user info
- Main Area:
    - flexible grid or card layout
    - modular Solid components for each card/widget
    - ex: balance overview, recent transactions, spending by category, upcoming recurring payments
- Responsive Handling:
    - desktop-first but sidebar collapses on narrow width (<900px)
    - cards stack vertically or switch to scrollable layout for smaller widths

## Data Visualization
- ApexCharts (with solid-apexcharts wrapper)
    - reactive, lightweight, interactive (tooltips, zoom, hover)
- Alt: Chart.js with solid-chartjs (simplicity)

## Other
- solid-form-handler, or custom Solid signals with small validation helpers (zod or valibot)
- solid-toast for notifications
- solid-transition-group for page transitions or element entrance animations
- lazy-load charts and large data tables
- memoization (creaeMemo) for computed values
- defer non-critical animations until idle
- all icons and interactive elements have `aria-label`s
- maintain good contrast
- keyboard navigation in sidebar/modals

## Routing
- `@solidjs/router`
- `<Link>` components

## Folder Structure
frontend/
├─ assets/                     # static assets (icons, fonts, images)
├─ components/                 # reusable building blocks
│  ├─ layout/                  # shared structural UI (sidebar, topbar)
│  │  ├─ Sidebar/
│  │  │  ├─ Sidebar.tsx
│  │  │  ├─ SidebarItem.tsx
│  │  │  └─ sidebar.css
│  │  ├─ Topbar/
│  │  │  ├─ Topbar.tsx
│  │  │  └─ topbar.css
│  │  └─ PageLayout.tsx        # wraps sidebar + main content
│  ├─ cards/                   # modular dashboard cards/widgets
│  │  ├─ BalanceCard.tsx
│  │  ├─ IncomeExpenseChart.tsx
│  │  ├─ BudgetsPieChart.tsx
│  │  ├─ RecentTransactionsCard.tsx
│  │  └─ UpcomingBillsCard.tsx
│  ├─ ui/                      # small generic components
│  │  ├─ Button.tsx
│  │  ├─ Card.tsx
│  │  ├─ Modal.tsx
│  │  ├─ Table.tsx
│  │  ├─ Toast.tsx
│  │  └─ Tooltip.tsx
│  └─ charts/                  # wrappers for chart libs (ApexCharts)
│     ├─ BarChart.tsx
│     ├─ PieChart.tsx
│     ├─ LineChart.tsx
│     └─ ChartTheme.ts
├─ pages/                      # top-level route components
│  ├─ Dashboard.tsx
│  ├─ Transactions.tsx
│  ├─ Budgets.tsx
│  ├─ ExchangeRates.tsx
│  ├─ Recurring.tsx
│  ├─ Rules.tsx
│  └─ Settings.tsx
├─ stores/                     # Solid signals/state; reactive global state
│  ├─ themeStore.ts            # handle dark/light mode, color theme switching, persist user pref
│  ├─ userStore.ts             # User preferences/profile state
│  ├─ currencyStore.ts         # hold ex rates, base currency, utility conversion functions
│  ├─ transactionStore.ts      # hold ex rates, base currency, utility conversion functions
│  ├─ budgetStore.ts           # hold ex rates, base currency, utility conversion functions
│  └─ notificationStore.ts     # centralized reactive state for in-app notifs
├─ styles/
│  ├─ theme.css                # CSS variables for color system
│  ├─ globals.css
│  ├─ layout.css
│  └─ neumorphic.css           # optional style variant
├─ utils/
│  ├─ formatters.ts            # currency/date helpers
│  ├─ api.ts                   # bridge to Wails backend functions
│  └─ chartUtils.ts
├─ App.tsx                     # main entry
└─ index.tsx                   # Solid render entry


### Example ThemeStore
```typescript
import { createSignal, onMount } from "solid-js";

const [theme, setTheme] = createSignal(localStorage.getItem("theme") || "dark");

const toggleTheme = () => {
  const newTheme = theme() === "dark" ? "light" : "dark";
  setTheme(newTheme);
  document.documentElement.setAttribute("data-theme", newTheme);
  localStorage.setItem("theme", newTheme);
};

onMount(() => {
  document.documentElement.setAttribute("data-theme", theme());
});

export { theme, toggleTheme };
```

### Example UserStore
```typescript
import { createStore } from "solid-js/store";

const [user, setUser] = createStore({
  name: "Local User",
  preferredCurrency: "USD",
  locale: "en-US",
  defaultPage: "dashboard",
});

const updateUser = (updates) => setUser(updates);

export { user, updateUser };
```

### Example CurrencyStore
```typescript
import { createSignal, createStore } from "solid-js/store";
import { getExchangeRates } from "../utils/api";

const [baseCurrency, setBaseCurrency] = createSignal("USD");
const [rates, setRates] = createStore({});

async function fetchRates() {
  const result = await getExchangeRates(baseCurrency());
  setRates(result);
}

function convert(amount: number, from: string, to: string) {
  if (from === to) return amount;
  const fromRate = rates[from] || 1;
  const toRate = rates[to] || 1;
  return (amount / fromRate) * toRate;
}

export { baseCurrency, setBaseCurrency, rates, fetchRates, convert };
```

### Example NotificationStore
```typescript
import { createStore } from "solid-js/store";

let idCounter = 0;
const [notifications, setNotifications] = createStore([]);

function addNotification(type, message, duration = 4000) {
  const id = ++idCounter;
  setNotifications([...notifications, { id, type, message }]);
  setTimeout(() => removeNotification(id), duration);
}

function removeNotification(id) {
  setNotifications(notifications.filter(n => n.id !== id));
}

export { notifications, addNotification, removeNotification };
```
