-- name: GetExpenseBreakdownByAccount :many
SELECT a.id, a.full_name, p.currency, SUM(p.amount) as amount
FROM postings p
JOIN transactions t ON t.id = p.transaction_id
JOIN accounts a ON a.id = p.account_id
WHERE t.date >= ? AND t.date <= ?
  AND p.amount > '0'
  AND a.account_type = 'expense'
GROUP BY a.id, a.full_name, p.currency
ORDER BY SUM(p.amount) DESC;

-- name: GetExpenseBreakdownByTag :many
SELECT tg.value, p.currency, SUM(p.amount) as amount
FROM postings p
JOIN transactions tx ON tx.id = p.transaction_id
JOIN transaction_tags tt ON tt.transaction_id = tx.id
JOIN tags tg ON tg.id = tt.tag_id
JOIN accounts a ON a.id = p.account_id
WHERE tx.date >= ? AND tx.date <= ?
  AND p.amount > '0'
  AND a.account_type = 'expense'
  AND tg.key = ?
GROUP BY tg.value, p.currency
ORDER BY SUM(p.amount) DESC;

-- name: GetTagSpendingSummary :many
SELECT tg.value, p.currency, SUM(p.amount) as amount
FROM postings p
JOIN transactions tx ON tx.id = p.transaction_id
JOIN transaction_tags tt ON tt.transaction_id = tx.id
JOIN tags tg ON tg.id = tt.tag_id
JOIN accounts a ON a.id = p.account_id
WHERE tx.date >= ? AND tx.date <= ?
  AND p.amount > '0'
  AND a.account_type = 'expense'
  AND tg.key = ?
GROUP BY tg.value, p.currency
ORDER BY SUM(p.amount) DESC;
