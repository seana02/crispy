-- name: CreateExchangeRate :one
INSERT INTO exchange_rates (base_currency, quote_currency, rate, effective_date, source)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetExchangeRate :one
SELECT * FROM exchange_rates WHERE id = ?;

-- name: GetEffectiveRate :one
SELECT * FROM exchange_rates
WHERE base_currency = ? AND quote_currency = ? AND effective_date <= ?
ORDER BY effective_date DESC
LIMIT 1;

-- name: GetEffectiveRateInverse :one
SELECT * FROM exchange_rates
WHERE base_currency = ? AND quote_currency = ? AND effective_date <= ?
ORDER BY effective_date DESC
LIMIT 1;

-- name: ListExchangeRates :many
SELECT * FROM exchange_rates
ORDER BY effective_date DESC, base_currency ASC, quote_currency ASC
LIMIT ? OFFSET ?;

-- name: CountExchangeRates :one
SELECT COUNT(*) FROM exchange_rates;

-- name: ListExchangeRatesByPair :many
SELECT * FROM exchange_rates
WHERE base_currency = ? AND quote_currency = ?
ORDER BY effective_date DESC
LIMIT ? OFFSET ?;

-- name: ListExchangeRatesByDateRange :many
SELECT * FROM exchange_rates
WHERE effective_date >= ? AND effective_date <= ?
ORDER BY effective_date DESC, base_currency ASC, quote_currency ASC
LIMIT ? OFFSET ?;

-- name: CountExchangeRatesByPair :one
SELECT COUNT(*) FROM exchange_rates WHERE base_currency = ? AND quote_currency = ?;

-- name: DeleteExchangeRate :exec
DELETE FROM exchange_rates WHERE id = ?;
