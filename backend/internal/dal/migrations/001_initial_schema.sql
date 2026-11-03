-- +goose Up

CREATE TABLE accounts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    full_name TEXT NOT NULL,
    parent_id INTEGER REFERENCES accounts(id),
    account_type TEXT CHECK(account_type IN ('asset', 'liability', 'revenue', 'expense', 'equity')),
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE INDEX idx_accounts_parent_id ON accounts(parent_id);
CREATE UNIQUE INDEX idx_accounts_full_name ON accounts(full_name);

CREATE TABLE transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending', 'posted', 'cleared')),
    recurring_template_id INTEGER REFERENCES recurring_templates(id),
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE INDEX idx_transactions_date ON transactions(date);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_recurring_template_id ON transactions(recurring_template_id);

CREATE TABLE postings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    transaction_id INTEGER NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    account_id INTEGER NOT NULL REFERENCES accounts(id),
    currency TEXT NOT NULL DEFAULT 'USD',
    amount TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE INDEX idx_postings_transaction_id ON postings(transaction_id);
CREATE INDEX idx_postings_account_id ON postings(account_id);

CREATE TABLE tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    UNIQUE(key, value)
);

CREATE INDEX idx_tags_key ON tags(key);
CREATE INDEX idx_tags_key_value ON tags(key, value);

CREATE TABLE transaction_tags (
    transaction_id INTEGER NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags(id),
    PRIMARY KEY (transaction_id, tag_id)
);

CREATE INDEX idx_transaction_tags_transaction_id ON transaction_tags(transaction_id);
CREATE INDEX idx_transaction_tags_tag_id ON transaction_tags(tag_id);

CREATE TABLE budgets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    direction TEXT NOT NULL CHECK(direction IN ('spending', 'saving')),
    period TEXT NOT NULL CHECK(period IN ('weekly', 'monthly', 'yearly')),
    currency TEXT NOT NULL DEFAULT 'USD',
    target_amount TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE TABLE budget_accounts (
    budget_id INTEGER NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    account_id INTEGER NOT NULL REFERENCES accounts(id),
    PRIMARY KEY (budget_id, account_id)
);

CREATE TABLE budget_tags (
    budget_id INTEGER NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags(id),
    PRIMARY KEY (budget_id, tag_id)
);

CREATE TABLE exchange_rates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    base_currency TEXT NOT NULL,
    quote_currency TEXT NOT NULL,
    rate TEXT NOT NULL,
    effective_date TEXT NOT NULL,
    source TEXT NOT NULL DEFAULT 'manual' CHECK(source IN ('manual', 'api')),
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE INDEX idx_exchange_rates_lookup ON exchange_rates(base_currency, quote_currency, effective_date);

CREATE TABLE recurring_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending', 'posted', 'cleared')),
    schedule_type TEXT NOT NULL CHECK(schedule_type IN ('every_n_days', 'every_n_weeks', 'every_n_months', 'monthly_on_day')),
    interval_n INTEGER NOT NULL DEFAULT 1,
    day_of_month INTEGER,
    start_date TEXT NOT NULL,
    end_date TEXT,
    last_generated_date TEXT,
    is_active INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE INDEX idx_recurring_templates_is_active ON recurring_templates(is_active);

CREATE TABLE recurring_template_postings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    template_id INTEGER NOT NULL REFERENCES recurring_templates(id) ON DELETE CASCADE,
    account_id INTEGER NOT NULL REFERENCES accounts(id),
    currency TEXT NOT NULL DEFAULT 'USD',
    amount TEXT NOT NULL
);

CREATE TABLE recurring_template_tags (
    template_id INTEGER NOT NULL REFERENCES recurring_templates(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags(id),
    PRIMARY KEY (template_id, tag_id)
);

CREATE TABLE widget_configurations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    widget_type TEXT NOT NULL,
    position_x INTEGER NOT NULL DEFAULT 0,
    position_y INTEGER NOT NULL DEFAULT 0,
    width INTEGER NOT NULL DEFAULT 1,
    height INTEGER NOT NULL DEFAULT 1,
    config TEXT NOT NULL DEFAULT '{}',
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE TABLE settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);

-- Seed default settings
INSERT INTO settings (key, value) VALUES
    ('default_currency', 'USD'),
    ('display_currency', 'USD'),
    ('available_currencies', '["USD","EUR","GBP","JPY","CAD","AUD","CHF"]'),
    ('date_format', 'YYYY-MM-DD'),
    ('theme', 'neon-night'),
    ('fiscal_year_start_month', '1'),
    ('fiscal_year_start_day', '1');

-- +goose Down

DROP TABLE IF EXISTS settings;
DROP TABLE IF EXISTS widget_configurations;
DROP TABLE IF EXISTS recurring_template_tags;
DROP TABLE IF EXISTS recurring_template_postings;
DROP TABLE IF EXISTS recurring_templates;
DROP TABLE IF EXISTS exchange_rates;
DROP TABLE IF EXISTS budget_tags;
DROP TABLE IF EXISTS budget_accounts;
DROP TABLE IF EXISTS budgets;
DROP TABLE IF EXISTS transaction_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS postings;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS accounts;
