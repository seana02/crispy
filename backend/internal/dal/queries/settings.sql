-- name: GetSetting :one
SELECT * FROM settings WHERE key = ?;

-- name: GetAllSettings :many
SELECT * FROM settings ORDER BY key ASC;

-- name: UpsertSetting :exec
INSERT INTO settings (key, value) VALUES (?, ?)
ON CONFLICT(key) DO UPDATE SET value = excluded.value;
