-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE LOWER(email) = LOWER($1);

-- name: ListUsers :many
SELECT *
FROM users
ORDER BY created_at DESC;

-- name: CreateUser :one
INSERT INTO users (
    first_name,
    last_name,
    email,
    password_hash,
    mfa_secret,
    mfa_enabled,
    created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;
