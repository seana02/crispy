-- name: CreateAccount :one
INSERT INTO accounts (name, full_name, parent_id, account_type)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts WHERE id = ?;

-- name: GetAccountByFullName :one
SELECT * FROM accounts WHERE full_name = ?;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY full_name ASC
LIMIT ? OFFSET ?;

-- name: CountAccounts :one
SELECT COUNT(*) FROM accounts;

-- name: ListAccountsByType :many
SELECT * FROM accounts
WHERE account_type = ?
ORDER BY full_name ASC
LIMIT ? OFFSET ?;

-- name: CountAccountsByType :one
SELECT COUNT(*) FROM accounts WHERE account_type = ?;

-- name: ListAccountsByName :many
SELECT * FROM accounts
WHERE full_name LIKE ?
ORDER BY full_name ASC
LIMIT ? OFFSET ?;

-- name: CountAccountsByName :one
SELECT COUNT(*) FROM accounts WHERE full_name LIKE ?;

-- name: ListAccountsByTypeAndName :many
SELECT * FROM accounts
WHERE account_type = ? AND full_name LIKE ?
ORDER BY full_name ASC
LIMIT ? OFFSET ?;

-- name: CountAccountsByTypeAndName :one
SELECT COUNT(*) FROM accounts WHERE account_type = ? AND full_name LIKE ?;

-- name: UpdateAccount :one
UPDATE accounts
SET name = ?, full_name = ?, parent_id = ?, account_type = ?,
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = ?;

-- name: GetChildAccounts :many
SELECT * FROM accounts WHERE parent_id = ?;

-- name: SearchAccounts :many
SELECT * FROM accounts
WHERE full_name LIKE ?
ORDER BY full_name ASC
LIMIT 20;

-- name: GetDescendantAccounts :many
SELECT * FROM accounts WHERE full_name LIKE ? ORDER BY full_name ASC;

-- name: HasPostingsForAccount :one
SELECT EXISTS(SELECT 1 FROM postings WHERE account_id = ?) AS has_postings;

-- name: GetRootAncestor :one
WITH RECURSIVE ancestors AS (
    SELECT accounts.id, accounts.name, accounts.full_name, accounts.parent_id, accounts.account_type
    FROM accounts WHERE accounts.id = ?
    UNION ALL
    SELECT a.id, a.name, a.full_name, a.parent_id, a.account_type
    FROM accounts a JOIN ancestors anc ON a.id = anc.parent_id
)
SELECT ancestors.id, ancestors.name, ancestors.full_name, ancestors.parent_id, ancestors.account_type FROM ancestors WHERE ancestors.parent_id IS NULL;
