-- name: CreateRecurringTemplate :one
INSERT INTO recurring_templates (description, status, schedule_type, interval_n, day_of_month, start_date, end_date, is_active)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetRecurringTemplate :one
SELECT * FROM recurring_templates WHERE id = ?;

-- name: ListRecurringTemplates :many
SELECT * FROM recurring_templates ORDER BY created_at DESC;

-- name: UpdateRecurringTemplate :one
UPDATE recurring_templates
SET description = ?, status = ?, schedule_type = ?, interval_n = ?,
    day_of_month = ?, start_date = ?, end_date = ?, is_active = ?,
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: UpdateRecurringTemplateLastGenerated :exec
UPDATE recurring_templates
SET last_generated_date = ?,
    updated_at = ?
WHERE id = ?;

-- name: DeleteRecurringTemplate :exec
DELETE FROM recurring_templates WHERE id = ?;

-- name: GetActiveTemplates :many
SELECT * FROM recurring_templates WHERE is_active = 1;

-- name: CreateTemplatePosting :one
INSERT INTO recurring_template_postings (template_id, account_id, currency, amount)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetTemplatePostings :many
SELECT * FROM recurring_template_postings WHERE template_id = ?;

-- name: GetTemplatePostingsWithAccounts :many
SELECT p.id, p.template_id, p.account_id, p.currency, p.amount,
       a.full_name AS account_full_name
FROM recurring_template_postings p
JOIN accounts a ON a.id = p.account_id
WHERE p.template_id = ?;

-- name: DeleteTemplatePostings :exec
DELETE FROM recurring_template_postings WHERE template_id = ?;

-- name: AddTemplateTag :exec
INSERT OR IGNORE INTO recurring_template_tags (template_id, tag_id) VALUES (?, ?);

-- name: RemoveTemplateTags :exec
DELETE FROM recurring_template_tags WHERE template_id = ?;

-- name: GetTemplateTags :many
SELECT t.* FROM tags t
JOIN recurring_template_tags rtt ON rtt.tag_id = t.id
WHERE rtt.template_id = ?;
