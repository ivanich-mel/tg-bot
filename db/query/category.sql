-- name: CreateCategory :one
INSERT INTO categories (
  name,
  balance,
  permanent_balance

) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetCategory :one
SELECT * FROM categories
WHERE id = $1 LIMIT 1;

-- name: ListCategories :many
SELECT * FROM categories
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateCategory :exec
UPDATE categories
SET
    name = $2,
    balance = $3,
    permanent_balance = $4
WHERE id = $1;
-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;