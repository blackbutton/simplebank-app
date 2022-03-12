-- name: CreateTransfer :one
INSERT INTO transfers (
    from_account_id, to_account_id, amount
) VALUES ($1, $2, $3)
returning *;

-- name: ListFRomAccountTransfers :many
SELECT * FROM transfers
WHERE from_account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListToAccountTransfers :many
SELECT * FROM transfers
WHERE to_account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = $1;

-- name: DeleteFromAccountTransfer :exec
DELETE FROM transfers
WHERE from_account_id = $1;

-- name: DeleteToAccountTransfer :exec
DELETE FROM transfers
WHERE to_account_id = $1;