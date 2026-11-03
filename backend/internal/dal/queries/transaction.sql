-- name: CreateTransaction :one
INSERT INTO transactions (date, description, status, recurring_template_id)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transactions WHERE id = ?;

-- name: ListTransactions :many
SELECT * FROM transactions
ORDER BY date DESC, id DESC
LIMIT ? OFFSET ?;

-- name: CountTransactions :one
SELECT COUNT(*) FROM transactions;

-- name: ListTransactionsByDateRange :many
SELECT * FROM transactions
WHERE date >= ? AND date <= ?
ORDER BY date DESC, id DESC
LIMIT ? OFFSET ?;

-- name: CountTransactionsByDateRange :one
SELECT COUNT(*) FROM transactions WHERE date >= ? AND date <= ?;

-- name: ListTransactionsByStatus :many
SELECT * FROM transactions
WHERE status = ?
ORDER BY date DESC, id DESC
LIMIT ? OFFSET ?;

-- name: ListTransactionsByCurrency :many
SELECT DISTINCT t.* FROM transactions t
JOIN postings p ON p.transaction_id = t.id
WHERE p.currency = ?
ORDER BY t.date DESC, t.id DESC
LIMIT ? OFFSET ?;

-- name: CountTransactionsByCurrency :one
SELECT COUNT(DISTINCT t.id) FROM transactions t
JOIN postings p ON p.transaction_id = t.id
WHERE p.currency = ?;

-- name: CountTransactionsByStatus :one
SELECT COUNT(*) FROM transactions WHERE status = ?;

-- name: ListTransactionsByDescription :many
SELECT * FROM transactions
WHERE description LIKE ?
ORDER BY date DESC, id DESC
LIMIT ? OFFSET ?;

-- name: CountTransactionsByDescription :one
SELECT COUNT(*) FROM transactions WHERE description LIKE ?;

-- name: ListTransactionsByAccount :many
SELECT DISTINCT t.* FROM transactions t
JOIN postings p ON p.transaction_id = t.id
WHERE p.account_id = ?
ORDER BY t.date DESC, t.id DESC
LIMIT ? OFFSET ?;

-- name: CountTransactionsByAccount :one
SELECT COUNT(DISTINCT t.id) FROM transactions t
JOIN postings p ON p.transaction_id = t.id
WHERE p.account_id = ?;

-- name: ListTransactionsByTagKey :many
SELECT DISTINCT t.* FROM transactions t
JOIN transaction_tags tt ON tt.transaction_id = t.id
JOIN tags tg ON tg.id = tt.tag_id
WHERE tg.key = ?
ORDER BY t.date DESC, t.id DESC
LIMIT ? OFFSET ?;

-- name: CountTransactionsByTagKey :one
SELECT COUNT(DISTINCT t.id) FROM transactions t
JOIN transaction_tags tt ON tt.transaction_id = t.id
JOIN tags tg ON tg.id = tt.tag_id
WHERE tg.key = ?;

-- name: ListTransactionsByTagKeyValue :many
SELECT DISTINCT t.* FROM transactions t
JOIN transaction_tags tt ON tt.transaction_id = t.id
JOIN tags tg ON tg.id = tt.tag_id
WHERE tg.key = ? AND tg.value = ?
ORDER BY t.date DESC, t.id DESC
LIMIT ? OFFSET ?;

-- name: CountTransactionsByTagKeyValue :one
SELECT COUNT(DISTINCT t.id) FROM transactions t
JOIN transaction_tags tt ON tt.transaction_id = t.id
JOIN tags tg ON tg.id = tt.tag_id
WHERE tg.key = ? AND tg.value = ?;

-- name: UpdateTransaction :one
UPDATE transactions
SET date = ?, description = ?, status = ?,
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteTransaction :exec
DELETE FROM transactions WHERE id = ?;
