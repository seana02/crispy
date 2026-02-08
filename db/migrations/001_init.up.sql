CREATE TABLE IF NOT EXISTS Account (
    id INTEGER PRIMARY KEY,
    parent_id INTEGER REFERENCES Account (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    type TEXT NOT NULL CHECK( type IN ('Asset', 'Liability', 'Revenue', 'Expense', 'Equity' ) ),
    currency TEXT NOT NULL,
    description TEXT,
    active BOOLEAN,
    date_created DATETIME,
    date_updated DATETIME,
    UNIQUE(parent_id, name)
);
CREATE INDEX IF NOT EXISTS idx_account_parent ON Account(parent_id, name);
INSERT INTO Account (id, name, type, currency) VALUES (0, 'root', 'Asset', 'USD');

CREATE TABLE IF NOT EXISTS "Transaction" (
    id INTEGER PRIMARY KEY,
    description TEXT,
    date DATETIME NOT NULL,
    status TEXT NOT NULL CHECK( status IN ('Pending', 'Posted', 'Cleared') ),
    reference_id INTEGER DEFAULT NULL REFERENCES RecurringTransaction(id) ON DELETE SET DEFAULT,
    date_created DATETIME,
    date_updated DATETIME
);
CREATE INDEX IF NOT EXISTS idx_transaction_date ON "Transaction"(date);

CREATE TABLE IF NOT EXISTS Posting (
    id INTEGER PRIMARY KEY,
    transaction_id INTEGER NOT NULL REFERENCES "Transaction"(id) ON DELETE CASCADE,
    account_id INTEGER NOT NULL REFERENCES Account(id),
    amount TEXT NOT NULL,
    currency TEXT NOT NULL,
    date_created DATETIME,
    date_updated DATETIME
);
CREATE INDEX IF NOT EXISTS idx_posting_transactionid ON Posting(transaction_id);
CREATE INDEX IF NOT EXISTS idx_posting_accountid ON Posting(account_id);

CREATE TABLE IF NOT EXISTS RecurringTransaction (
    id INTEGER PRIMARY KEY,
    description TEXT,
    recurrence_id INTEGER NOT NULL REFERENCES Recurrence(id),
    last_run TEXT NOT NULL DEFAULT '1970-1-1',
    active BOOLEAN NOT NULL,
    status TEXT NOT NULL CHECK( status IN ('Pending', 'Posted', 'Cleared') ),
    date_created DATETIME,
    date_updated DATETIME
);

CREATE TABLE IF NOT EXISTS RecurringPosting (
    id INTEGER PRIMARY KEY,
    template_id INTEGER NOT NULL REFERENCES RecurringTransaction(id),
    account_id INTEGER NOT NULL REFERENCES Account(id),
    amount_expression TEXT NOT NULL,
    currency TEXT NOT NULL,
    date_created DATETIME
);

CREATE TABLE IF NOT EXISTS Recurrence (
    id INTEGER PRIMARY KEY,
    frequency TEXT NOT NULL CHECK ( frequency IN ('Daily', 'Weekly', 'Monthly') ),
    interval INTEGER NOT NULL,
    day_of_week INTEGER,
    day_of_month INTEGER,
    start_date DATETIME,
    end_date DATETIME,
    CHECK (
        (frequency = 'Weekly' AND day_of_week BETWEEN 0 AND 6)
        OR (frequency <> 'Weekly' AND day_of_week IS NULL)
    ),
    CHECK (
        (frequency = 'Monthly' AND day_of_month BETWEEN 1 AND 31)
        OR (frequency <> 'Monthly' AND day_of_month IS NULL)
    )
);

CREATE TABLE IF NOT EXISTS Budget (
    id INTEGER PRIMARY KEY,
    name TEXT,
    description TEXT,
    target TEXT,
    tag INTEGER NOT NULL REFERENCES Tag(id),
    reset_frequency INTEGER NOT NULL REFERENCES Recurrence(id)
);

CREATE TABLE IF NOT EXISTS Rule (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    condition TEXT NOT NULL,
    "action" TEXT NOT NULL,
    enabled BOOLEAN NOT NULL,
    date_created DATETIME,
    date_updated DATETIME
);

CREATE TABLE IF NOT EXISTS Tag (
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS Transaction_Tags (
    transaction_id INTEGER NOT NULL REFERENCES "Transaction"(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES Tag(id) ON DELETE CASCADE,
    UNIQUE(transaction_id, tag_id)
);

CREATE TABLE IF NOT EXISTS Attachment (
    id INTEGER PRIMARY KEY,
    transaction_id INTEGER NOT NULL REFERENCES "Transaction"(id),
    filename TEXT NOT NULL,
    path TEXT NOT NULL,
    mime_type TEXT,
    date_created DATETIME
);

CREATE TABLE IF NOT EXISTS ExchangeRate (
    id INTEGER PRIMARY KEY,
    from_currency TEXT NOT NULL,
    to_currency TEXT NOT NULL,
    rate TEXT NOT NULL,
    date DATETIME NOT NULL,
    source TEXT
);
