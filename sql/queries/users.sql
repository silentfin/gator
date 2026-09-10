-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;


-- name: GetUser :one
SELECT * FROM users where name=$1;

-- name: Reset :exec
TRUNCATE TABLE users CASCADE;

-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUserFromID :one
SELECT * FROM users where id=$1;
