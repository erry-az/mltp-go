// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package query

import (
	"context"
	"time"
)

const createTransaction = `-- name: CreateTransaction :one
INSERT INTO transactions (user_id, amount, name, type)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, amount, name, type, created_at
`

type CreateTransactionParams struct {
	UserID int32
	Amount int64
	Name   TransactionName
	Type   TransactionType
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, createTransaction,
		arg.UserID,
		arg.Amount,
		arg.Name,
		arg.Type,
	)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Amount,
		&i.Name,
		&i.Type,
		&i.CreatedAt,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (username, fullname)
VALUES ($1, $2)
RETURNING id, username, fullname, balance, created_at, updated_at
`

type CreateUserParams struct {
	Username string
	Fullname string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.Fullname)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Fullname,
		&i.Balance,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getTopUserTransaction = `-- name: GetTopUserTransaction :many
SELECT u.username, t.amount, t.type, t.created_at
FROM users u
INNER JOIN transactions t
ON u.id = t.user_id
WHERE t.type=$2
ORDER BY t.amount DESC
LIMIT $1
`

type GetTopUserTransactionParams struct {
	Limit int32
	Type  TransactionType
}

type GetTopUserTransactionRow struct {
	Username  string
	Amount    int64
	Type      TransactionType
	CreatedAt time.Time
}

func (q *Queries) GetTopUserTransaction(ctx context.Context, arg GetTopUserTransactionParams) ([]GetTopUserTransactionRow, error) {
	rows, err := q.db.QueryContext(ctx, getTopUserTransaction, arg.Limit, arg.Type)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTopUserTransactionRow
	for rows.Next() {
		var i GetTopUserTransactionRow
		if err := rows.Scan(
			&i.Username,
			&i.Amount,
			&i.Type,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, username, fullname, balance, created_at, updated_at
FROM users
WHERE username = $1
LIMIT 1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Fullname,
		&i.Balance,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const topSummaryTransactions = `-- name: TopSummaryTransactions :many
SELECT u.username, sum(t.amount) transacted_value
FROM transactions t
INNER JOIN users u on u.id = t.user_id
WHERE type='debit'
GROUP BY u.username
ORDER BY sum(t.amount) DESC
LIMIT $1
`

type TopSummaryTransactionsRow struct {
	Username        string
	TransactedValue int64
}

func (q *Queries) TopSummaryTransactions(ctx context.Context, limit int32) ([]TopSummaryTransactionsRow, error) {
	rows, err := q.db.QueryContext(ctx, topSummaryTransactions, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TopSummaryTransactionsRow
	for rows.Next() {
		var i TopSummaryTransactionsRow
		if err := rows.Scan(&i.Username, &i.TransactedValue); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const topTransactionByUsername = `-- name: TopTransactionByUsername :many
SELECT u.username,
    cast(CASE
       WHEN t.type='debit' THEN -t.amount
       ELSE t.amount
    END AS BIGINT) amount_value
FROM transactions t
INNER JOIN users u on u.id = t.user_id
WHERE u.username = $1
ORDER BY t.amount DESC
LIMIT $2
`

type TopTransactionByUsernameParams struct {
	Username string
	Limit    int32
}

type TopTransactionByUsernameRow struct {
	Username    string
	AmountValue int64
}

func (q *Queries) TopTransactionByUsername(ctx context.Context, arg TopTransactionByUsernameParams) ([]TopTransactionByUsernameRow, error) {
	rows, err := q.db.QueryContext(ctx, topTransactionByUsername, arg.Username, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TopTransactionByUsernameRow
	for rows.Next() {
		var i TopTransactionByUsernameRow
		if err := rows.Scan(&i.Username, &i.AmountValue); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateBalance = `-- name: UpdateBalance :one
UPDATE users
SET balance = balance + $2
WHERE id = $1
RETURNING id, username, fullname, balance, created_at, updated_at
`

type UpdateBalanceParams struct {
	ID      int32
	Balance int64
}

func (q *Queries) UpdateBalance(ctx context.Context, arg UpdateBalanceParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateBalance, arg.ID, arg.Balance)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Fullname,
		&i.Balance,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
