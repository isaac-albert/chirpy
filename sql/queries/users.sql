-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1
)
RETURNING *;

-- name: DeleteTable :execrows
DELETE FROM users;

-- name: CreateUserMessage :one
INSERT INTO userchirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetUsers :many
SELECT * FROM userchirps 
ORDER BY created_at ASC;

-- name: GetMessage :one
SELECT * FROM userchirps
WHERE id = $1;