-- name: CreateTag :one
INSERT INTO tags (key, value) VALUES (?, ?)
ON CONFLICT(key, value) DO UPDATE SET key = key
RETURNING *;

-- name: GetTag :one
SELECT * FROM tags WHERE id = ?;

-- name: GetTagByKeyValue :one
SELECT * FROM tags WHERE key = ? AND value = ?;

-- name: ListTags :many
SELECT t.*, COUNT(tt.transaction_id) as transaction_count
FROM tags t
LEFT JOIN transaction_tags tt ON tt.tag_id = t.id
GROUP BY t.id
ORDER BY t.key ASC, t.value ASC
LIMIT ? OFFSET ?;

-- name: CountTags :one
SELECT COUNT(*) FROM tags;

-- name: ListTagsByKey :many
SELECT t.*, COUNT(tt.transaction_id) as transaction_count
FROM tags t
LEFT JOIN transaction_tags tt ON tt.tag_id = t.id
WHERE t.key LIKE ?
GROUP BY t.id
ORDER BY t.key ASC, t.value ASC
LIMIT ? OFFSET ?;

-- name: CountTagsByKey :one
SELECT COUNT(*) FROM tags WHERE key LIKE ?;

-- name: SearchTags :many
SELECT * FROM tags
WHERE key LIKE ? OR value LIKE ?
ORDER BY key ASC, value ASC
LIMIT 20;

-- name: AddTransactionTag :exec
INSERT OR IGNORE INTO transaction_tags (transaction_id, tag_id) VALUES (?, ?);

-- name: RemoveTransactionTags :exec
DELETE FROM transaction_tags WHERE transaction_id = ?;

-- name: GetTagsByTransaction :many
SELECT t.* FROM tags t
JOIN transaction_tags tt ON tt.tag_id = t.id
WHERE tt.transaction_id = ?
ORDER BY t.key ASC, t.value ASC;
