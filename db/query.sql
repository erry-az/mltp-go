-- name: CreateUser :one
INSERT INTO users (username, fullname)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1
LIMIT 1;

-- name: GetTopUserTransaction :many
SELECT u.username, t.amount, t.type, t.created_at
FROM users u
INNER JOIN transactions t
ON u.id = t.user_id
WHERE t.type=$2
ORDER BY t.amount DESC
LIMIT $1;

-- name: CreateTransaction :one
INSERT INTO transactions (user_id, amount, name, type)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateBalance :one
UPDATE users
SET balance = balance + $2
WHERE id = $1
RETURNING *;

-- name: TopSummaryTransactions :many
SELECT u.username, sum(t.amount) transacted_value
FROM transactions t
INNER JOIN users u on u.id = t.user_id
WHERE type='debit'
GROUP BY u.username
ORDER BY sum(t.amount) DESC
LIMIT $1;

-- name: TopTransactionByUsername :many
SELECT u.username,
    cast(CASE
       WHEN t.type='debit' THEN -t.amount
       ELSE t.amount
    END AS BIGINT) amount_value
FROM transactions t
INNER JOIN users u on u.id = t.user_id
WHERE u.username = $1
ORDER BY t.amount DESC
LIMIT $2;