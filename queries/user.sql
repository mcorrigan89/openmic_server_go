-- name: GetUserByID :one
SELECT sqlc.embed(users) FROM users
WHERE users.id = $1;

-- name: GetUserByEmail :one
SELECT sqlc.embed(users) FROM users
WHERE users.email = $1;

-- name: GetUserByHandle :one
SELECT sqlc.embed(users) FROM users
WHERE users.user_handle = $1;

-- name: GetUserBySessionToken :one
SELECT sqlc.embed(users), sqlc.embed(user_session) FROM users
JOIN user_session ON users.id = user_session.user_id
WHERE user_session.token = $1;

-- name: CreateUser :one
INSERT INTO users (id, given_name, family_name, email, email_verified, claimed, user_handle) 
VALUES (
    sqlc.arg(id), 
    sqlc.narg(given_name), 
    sqlc.narg(family_name), 
    sqlc.arg(email), 
    sqlc.arg(email_verified)::boolean,
    sqlc.arg(claimed),
    sqlc.arg(user_handle)
) RETURNING *;
 
-- name: UpdateUser :one
UPDATE users SET 
    id = sqlc.arg(id), 
    given_name = sqlc.narg(given_name), 
    family_name = sqlc.narg(family_name),
    email = sqlc.arg(email),
    email_verified = sqlc.arg(email_verified)::boolean,
    claimed = sqlc.arg(claimed),
    user_handle = sqlc.arg(user_handle)
WHERE id = sqlc.arg(id) RETURNING *;

-- name: SetAvatarImage :one
UPDATE users SET avatar_id = sqlc.arg(image_id) WHERE id = sqlc.arg(user_id) RETURNING *;

-- name: CreateUserSession :one
INSERT INTO user_session (user_id, token, expires_at) VALUES (sqlc.arg(user_id), sqlc.arg(token), sqlc.arg(expires_at)) RETURNING *;

-- name: ExpireUserSession :exec
UPDATE user_session SET user_expired = TRUE WHERE user_session.id = $1;
