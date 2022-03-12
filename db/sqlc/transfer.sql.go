// Code generated by sqlc. DO NOT EDIT.
// source: transfer.sql

package db

import (
	"context"
)

const createTransfer = `-- name: CreateTransfer :one
INSERT INTO transfers (
    from_account_id, to_account_id, amount
) VALUES ($1, $2, $3)
returning id, from_account_id, to_account_id, amount, created_at
`

type CreateTransferParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

func (q *Queries) CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, createTransfer, arg.FromAccountID, arg.ToAccountID, arg.Amount)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const deleteFromAccountTransfer = `-- name: DeleteFromAccountTransfer :exec
DELETE FROM transfers
WHERE from_account_id = $1
`

func (q *Queries) DeleteFromAccountTransfer(ctx context.Context, fromAccountID int64) error {
	_, err := q.db.ExecContext(ctx, deleteFromAccountTransfer, fromAccountID)
	return err
}

const deleteToAccountTransfer = `-- name: DeleteToAccountTransfer :exec
DELETE FROM transfers
WHERE to_account_id = $1
`

func (q *Queries) DeleteToAccountTransfer(ctx context.Context, toAccountID int64) error {
	_, err := q.db.ExecContext(ctx, deleteToAccountTransfer, toAccountID)
	return err
}

const deleteTransfer = `-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = $1
`

func (q *Queries) DeleteTransfer(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTransfer, id)
	return err
}

const getTransfer = `-- name: GetTransfer :one
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTransfer(ctx context.Context, id int64) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, getTransfer, id)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const listFRomAccountTransfers = `-- name: ListFRomAccountTransfers :many
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
WHERE from_account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListFRomAccountTransfersParams struct {
	FromAccountID int64 `json:"from_account_id"`
	Limit         int32 `json:"limit"`
	Offset        int32 `json:"offset"`
}

func (q *Queries) ListFRomAccountTransfers(ctx context.Context, arg ListFRomAccountTransfersParams) ([]Transfer, error) {
	rows, err := q.db.QueryContext(ctx, listFRomAccountTransfers, arg.FromAccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transfer
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Amount,
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

const listToAccountTransfers = `-- name: ListToAccountTransfers :many
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
WHERE to_account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListToAccountTransfersParams struct {
	ToAccountID int64 `json:"to_account_id"`
	Limit       int32 `json:"limit"`
	Offset      int32 `json:"offset"`
}

func (q *Queries) ListToAccountTransfers(ctx context.Context, arg ListToAccountTransfersParams) ([]Transfer, error) {
	rows, err := q.db.QueryContext(ctx, listToAccountTransfers, arg.ToAccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transfer
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Amount,
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