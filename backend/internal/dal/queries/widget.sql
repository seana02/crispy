-- name: CreateWidget :one
INSERT INTO widget_configurations (widget_type, position_x, position_y, width, height, config)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetWidget :one
SELECT * FROM widget_configurations WHERE id = ?;

-- name: ListWidgets :many
SELECT * FROM widget_configurations ORDER BY position_y ASC, position_x ASC;

-- name: UpdateWidget :one
UPDATE widget_configurations
SET widget_type = ?, position_x = ?, position_y = ?, width = ?, height = ?, config = ?,
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteWidget :exec
DELETE FROM widget_configurations WHERE id = ?;
