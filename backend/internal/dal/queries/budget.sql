-- name: CreateBudget :one
INSERT INTO budgets (name, direction, period, currency, target_amount)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetBudget :one
SELECT * FROM budgets WHERE id = ?;

-- name: ListBudgets :many
SELECT * FROM budgets ORDER BY name ASC;

-- name: UpdateBudget :one
UPDATE budgets
SET name = ?, direction = ?, period = ?, currency = ?, target_amount = ?,
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteBudget :exec
DELETE FROM budgets WHERE id = ?;

-- name: AddBudgetAccount :exec
INSERT OR IGNORE INTO budget_accounts (budget_id, account_id) VALUES (?, ?);

-- name: RemoveBudgetAccounts :exec
DELETE FROM budget_accounts WHERE budget_id = ?;

-- name: GetBudgetAccounts :many
SELECT a.* FROM accounts a
JOIN budget_accounts ba ON ba.account_id = a.id
WHERE ba.budget_id = ?;

-- name: AddBudgetTag :exec
INSERT OR IGNORE INTO budget_tags (budget_id, tag_id) VALUES (?, ?);

-- name: RemoveBudgetTags :exec
DELETE FROM budget_tags WHERE budget_id = ?;

-- name: GetBudgetTags :many
SELECT t.* FROM tags t
JOIN budget_tags bt ON bt.tag_id = t.id
WHERE bt.budget_id = ?;
