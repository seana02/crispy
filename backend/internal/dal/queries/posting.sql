-- name: CreatePosting :one
INSERT INTO postings (transaction_id, account_id, currency, amount)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetPostingsByTransaction :many
SELECT * FROM postings WHERE transaction_id = ? ORDER BY id ASC;

-- name: GetPostingsByTransactionWithAccounts :many
SELECT p.id, p.transaction_id, p.account_id, p.currency, p.amount, p.created_at,
       a.full_name AS account_full_name
FROM postings p
JOIN accounts a ON a.id = p.account_id
WHERE p.transaction_id = ?
ORDER BY p.id ASC;

-- name: DeletePostingsByTransaction :exec
DELETE FROM postings WHERE transaction_id = ?;

-- name: GetPostingsByAccount :many
SELECT * FROM postings WHERE account_id = ? ORDER BY id ASC;
