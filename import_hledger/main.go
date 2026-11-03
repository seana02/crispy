// import_hledger is a standalone script that parses an hledger journal file and
// imports accounts and transactions into the crispy SQLite database.
//
// Usage:
//
//	go run main.go -db <path-to-crispy.db> -journal <path-to-file.journal>
//
// Rules:
//   - Accounts are created lazily from posting lines, including all ancestors.
//   - account_type is inferred from the first segment of the full_name:
//     assets→asset, liabilities/liability→liability, revenues/income/revenue→revenue,
//     expenses/expense→expense, equity→equity.
//   - Transactions dated before 2026-10-01 get status "cleared"; others get "posted".
//   - Only transactions and their postings are imported. Directives, comments,
//     price directives (P), account declarations, etc. are silently ignored.
//   - If a transaction has a posting with no amount (the auto-balancing posting),
//     its amount is computed as the negation of the sum of all other postings
//     in the same currency. Mixed-currency auto-balance is not supported.
//   - The script is idempotent on accounts (upsert by full_name) but inserts
//     transactions fresh each run; run against an empty DB or handle duplicates
//     externally.
package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode"

	_ "github.com/mattn/go-sqlite3"
)

// ---- domain types ----------------------------------------------------------

type Posting struct {
	Account  string
	Currency string
	Amount   *big.Rat // nil means "auto-balance"
	Negative bool     // true when the original amount was negative
}

type Transaction struct {
	Date        string
	Description string
	Status      string // "cleared" or "posted"
	Postings    []Posting
}

// ---- hledger parsing -------------------------------------------------------

var (
	// Transaction header: date [status-flag] description
	// e.g.  2024-01-15 * Grocery store
	//       2024-01-15 ! Pending payment
	//       2024-01-15 Cleared transaction
	txnHeaderRe = regexp.MustCompile(`^(\d{4}[-/]\d{2}[-/]\d{2})(?:\s+([*!]))?\s+(.*)$`)

	// Posting line: two or more spaces then account name, optionally followed by
	// two or more spaces and an amount.
	// amount examples:  $123.45  $-12  EUR 100  100 EUR  -100.00 USD
	postingRe = regexp.MustCompile(`^(\s{2,})([^\s].+?)(?:\s{2,}(.+))?$`)

	// Amount: optional sign, optional currency symbol before, number, optional
	// currency code after.
	// Handles: $1,234.56  $-1,234.56  1234.56 USD  -1234.56 USD  EUR 100
	amountRe = regexp.MustCompile(
		`^(?:([A-Z]{3}|[$€£¥₹])\s*)?` + // optional currency before number
			`(-?)` + // minus
			`([\d,]+(?:\.\d+)?)` + // number (with optional commas)
			`(?:\s+([A-Z]{3}))?` + // optional currency after number
			`$`,
	)
)

func normalizeDate(d string) string {
	return strings.ReplaceAll(d, "/", "-")
}

func parseAmount(raw string) (currency string, amount *big.Rat, err error) {
	raw = strings.TrimSpace(raw)
	// Remove @ cost annotations (e.g. "100 EUR @ $1.10")
	if idx := strings.Index(raw, " @ "); idx >= 0 {
		raw = strings.TrimSpace(raw[:idx])
	}
	if idx := strings.Index(raw, " @@ "); idx >= 0 {
		raw = strings.TrimSpace(raw[:idx])
	}

	m := amountRe.FindStringSubmatch(raw)
	if m == nil {
		return "", nil, fmt.Errorf("cannot parse amount %q", raw)
	}
	currBefore := m[1]
	sign := m[2]
	numStr := strings.ReplaceAll(m[3], ",", "") // strip thousand-separators
	currAfter := m[4]

	currency = currBefore
	if currency == "" {
		currency = currAfter
	}
	// Normalise currency symbols to codes
	switch currency {
	case "$":
		currency = "USD"
	case "€":
		currency = "EUR"
	case "£":
		currency = "GBP"
	case "¥":
		currency = "JPY"
	case "₹":
		currency = "INR"
	case "":
		currency = "USD" // default
	}

	r := new(big.Rat)
	if _, ok := r.SetString(sign + numStr); !ok {
		return "", nil, fmt.Errorf("cannot parse number %q", sign+numStr)
	}
	return currency, r, nil
}

// parseJournal reads an hledger journal file and returns all transactions.
func parseJournal(path string) ([]Transaction, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var txns []Transaction
	var current *Transaction

	flush := func() {
		if current == nil {
			return
		}
		if err := balanceTransaction(current); err != nil {
			log.Printf("warning: skipping unbalanced transaction %q on %s: %v",
				current.Description, current.Date, err)
			current = nil
			return
		}
		txns = append(txns, *current)
		current = nil
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		// Strip inline comments (; or #)
		if idx := strings.Index(line, ";"); idx >= 0 {
			line = line[:idx]
		}
		trimmed := strings.TrimRightFunc(line, unicode.IsSpace)

		// Blank line = end of transaction block
		if trimmed == "" {
			flush()
			continue
		}

		// Skip top-level directives and comment lines
		if len(trimmed) > 0 && !unicode.IsSpace(rune(trimmed[0])) {
			// Could be a transaction header or a directive
			if m := txnHeaderRe.FindStringSubmatch(trimmed); m != nil {
				flush()
				date := normalizeDate(m[1])
				statusFlag := m[2]
				desc := strings.TrimSpace(m[3])

				// Determine status based on date rule:
				// before 2026-10-01 → cleared, otherwise → posted.
				// hledger's own * flag also means cleared; ! means pending.
				// We apply the date rule as the primary rule and only use
				// the flag as a secondary hint for the "posted" bucket.
				txStatus := statusFromDate(date)
				_ = statusFlag // date rule takes precedence per requirements

				current = &Transaction{
					Date:        date,
					Description: desc,
					Status:      txStatus,
				}
			}
			// else: directive line — skip
			continue
		}

		// Posting line (starts with whitespace)
		if current == nil {
			continue
		}
		m := postingRe.FindStringSubmatch(trimmed)
		if m == nil {
			continue
		}
		accountRaw := strings.TrimSpace(m[2])
		amountRaw := strings.TrimSpace(m[3])

		account := accountRaw

		p := Posting{Account: account}
		if amountRaw == "" {
			// auto-balance posting
			p.Amount = nil
			p.Currency = ""
		} else {
			cur, amt, err := parseAmount(amountRaw)
			if err != nil {
				log.Printf("warning: cannot parse amount %q in transaction %q: %v — skipping posting",
					amountRaw, current.Description, err)
				continue
			}
			p.Currency = cur
			p.Amount = amt
		}
		current.Postings = append(current.Postings, p)
	}
	flush()

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return txns, nil
}

func statusFromDate(date string) string {
	cutoff := "2026-10-01"
	if date < cutoff {
		return "cleared"
	}
	return "posted"
}

// balanceTransaction fills in a nil-amount posting by computing the negation
// of the sum of the other postings. If there is no nil posting, it validates
// that the transaction balances (sum == 0) within each currency group.
func balanceTransaction(t *Transaction) error {
	// Find auto-balance posting index
	autoIdx := -1
	for i, p := range t.Postings {
		if p.Amount == nil {
			if autoIdx >= 0 {
				return fmt.Errorf("more than one auto-balance posting")
			}
			autoIdx = i
		}
	}

	// Sum per currency
	sums := make(map[string]*big.Rat)
	for i, p := range t.Postings {
		if i == autoIdx {
			continue
		}
		if _, ok := sums[p.Currency]; !ok {
			sums[p.Currency] = new(big.Rat)
		}
		sums[p.Currency].Add(sums[p.Currency], p.Amount)
	}

	if autoIdx >= 0 {
		if len(sums) != 1 {
			return fmt.Errorf("auto-balance posting with mixed currencies is not supported")
		}
		for cur, sum := range sums {
			neg := new(big.Rat).Neg(sum)
			t.Postings[autoIdx].Currency = cur
			t.Postings[autoIdx].Amount = neg
		}
		return nil
	}

	// No auto-balance: verify each currency sums to zero
	for cur, sum := range sums {
		if sum.Sign() != 0 {
			return fmt.Errorf("currency %s does not balance (sum=%s)", cur, sum.RatString())
		}
	}
	return nil
}

// ---- account helpers -------------------------------------------------------

// accountType infers the type from the first segment of the full_name.
func accountType(fullName string) string {
	first := strings.ToLower(strings.SplitN(fullName, ":", 2)[0])
	switch first {
	case "assets", "asset":
		return "asset"
	case "liabilities", "liability":
		return "liability"
	case "revenues", "revenue", "income":
		return "revenue"
	case "expenses", "expense":
		return "expense"
	case "equity":
		return "equity"
	}
	return ""
}

// ancestors returns all ancestor full_names for a given full account name,
// from the root down (not including the account itself).
// e.g. "Assets:Bank:Checking" → ["Assets", "Assets:Bank"]
func ancestors(fullName string) []string {
	parts := strings.Split(fullName, ":")
	var result []string
	for i := 1; i < len(parts); i++ {
		result = append(result, strings.Join(parts[:i], ":"))
	}
	return result
}

// shortName returns the last segment of a full account name.
func shortName(fullName string) string {
	parts := strings.Split(fullName, ":")
	return parts[len(parts)-1]
}

// ---- database operations ---------------------------------------------------

// ensureAccount upserts an account and returns its id.
// It recursively ensures all ancestors exist first.
func ensureAccount(db *sql.DB, cache map[string]int64, fullName string) (int64, error) {
	if id, ok := cache[fullName]; ok {
		return id, nil
	}

	// Ensure parent first
	var parentID sql.NullInt64
	parts := strings.Split(fullName, ":")
	if len(parts) > 1 {
		parentFull := strings.Join(parts[:len(parts)-1], ":")
		pid, err := ensureAccount(db, cache, parentFull)
		if err != nil {
			return 0, err
		}
		parentID = sql.NullInt64{Int64: pid, Valid: true}
	}

	name := shortName(fullName)
	atype := ""
	if len(parts) == 1 {
		atype = accountType(fullName)
	}

	// Try to get existing
	var id int64
	err := db.QueryRow(`SELECT id FROM accounts WHERE full_name = ?`, fullName).Scan(&id)
	if err == sql.ErrNoRows {
		// Insert
		var res sql.Result
		if atype != "" {
			res, err = db.Exec(
				`INSERT INTO accounts (name, full_name, parent_id, account_type) VALUES (?, ?, ?, ?)`,
				name, fullName, parentID, atype,
			)
		} else {
			res, err = db.Exec(
				`INSERT INTO accounts (name, full_name, parent_id) VALUES (?, ?, ?)`,
				name, fullName, parentID,
			)
		}
		if err != nil {
			return 0, fmt.Errorf("insert account %q: %w", fullName, err)
		}
		id, err = res.LastInsertId()
		if err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, fmt.Errorf("query account %q: %w", fullName, err)
	}

	cache[fullName] = id
	return id, nil
}

func importTransactions(db *sql.DB, txns []Transaction) error {
	accountCache := make(map[string]int64)

	dbTx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			dbTx.Rollback()
		}
	}()

	for _, t := range txns {
		// Insert transaction
		res, err := dbTx.Exec(
			`INSERT INTO transactions (date, description, status) VALUES (?, ?, ?)`,
			t.Date, t.Description, t.Status,
		)
		if err != nil {
			return fmt.Errorf("insert transaction %q on %s: %w", t.Description, t.Date, err)
		}
		txID, err := res.LastInsertId()
		if err != nil {
			return err
		}

		for _, p := range t.Postings {
			// Ensure account (and ancestors) exist — run outside the deferred
			// db transaction so we can re-use existing accounts across txns.
			// We need to use the same connection though; pass dbTx-capable db.
			accID, err := ensureAccountTx(dbTx, accountCache, p.Account)
			if err != nil {
				return err
			}

			amtStr := ratToDecimal(p.Amount)
			_, err = dbTx.Exec(
				`INSERT INTO postings (transaction_id, account_id, currency, amount) VALUES (?, ?, ?, ?)`,
				txID, accID, p.Currency, amtStr,
			)
			if err != nil {
				return fmt.Errorf("insert posting for account %q: %w", p.Account, err)
			}
		}
	}

	return dbTx.Commit()
}

// ensureAccountTx is like ensureAccount but operates within a *sql.Tx.
func ensureAccountTx(tx *sql.Tx, cache map[string]int64, fullName string) (int64, error) {
	if id, ok := cache[fullName]; ok {
		return id, nil
	}

	var parentID sql.NullInt64
	parts := strings.Split(fullName, ":")
	if len(parts) > 1 {
		parentFull := strings.Join(parts[:len(parts)-1], ":")
		pid, err := ensureAccountTx(tx, cache, parentFull)
		if err != nil {
			return 0, err
		}
		parentID = sql.NullInt64{Int64: pid, Valid: true}
	}

	name := shortName(fullName)
	atype := ""
	if len(parts) == 1 {
		atype = accountType(fullName)
	}

	var id int64
	err := tx.QueryRow(`SELECT id FROM accounts WHERE full_name = ?`, fullName).Scan(&id)
	if err == sql.ErrNoRows {
		var res sql.Result
		if atype != "" {
			res, err = tx.Exec(
				`INSERT INTO accounts (name, full_name, parent_id, account_type) VALUES (?, ?, ?, ?)`,
				name, fullName, parentID, atype,
			)
		} else {
			res, err = tx.Exec(
				`INSERT INTO accounts (name, full_name, parent_id) VALUES (?, ?, ?)`,
				name, fullName, parentID,
			)
		}
		if err != nil {
			return 0, fmt.Errorf("insert account %q: %w", fullName, err)
		}
		id, err = res.LastInsertId()
		if err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, fmt.Errorf("query account %q: %w", fullName, err)
	}

	cache[fullName] = id
	return id, nil
}

// ratToDecimal converts a *big.Rat to a decimal string with enough precision.
func ratToDecimal(r *big.Rat) string {
	if r.IsInt() {
		return r.Num().String()
	}
	// Convert to float string with 10 decimal places then trim trailing zeros.
	f, _ := r.Float64()
	s := fmt.Sprintf("%.10f", f)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

// ---- main ------------------------------------------------------------------

func main() {
	dbPath := flag.String("db", "", "path to crispy.db (required)")
	journalPath := flag.String("journal", "", "path to hledger journal file (required)")
	flag.Parse()

	if *dbPath == "" || *journalPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Parse the journal
	log.Printf("Parsing journal: %s", *journalPath)
	txns, err := parseJournal(*journalPath)
	if err != nil {
		log.Fatalf("parse journal: %v", err)
	}
	log.Printf("Parsed %d transactions", len(txns))

	// Open the database
	dsn := *dbPath + "?_foreign_keys=on&_journal_mode=WAL"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	// Import
	start := time.Now()
	if err := importTransactions(db, txns); err != nil {
		log.Fatalf("import: %v", err)
	}
	log.Printf("Import complete in %s", time.Since(start).Round(time.Millisecond))
}
