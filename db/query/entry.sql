-- name: CreateEntry :one
INSERT INTO entries (
    account_id, amount
) VALUES ($1, $2)
returning *;

-- name: ListEntry :many
SELECT * FROM entries
WHERE account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1
LIMIT 1;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;

-- name: DeleteAccountEntry :exec
DELETE FROM entries
WHERE account_id = $1;