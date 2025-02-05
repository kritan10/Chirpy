-- name: CreateRefreshToken :one
INSERT INTO refresh_token(token, created_at, updated_at, user_id, expires_at)
VALUES ($1, NOW(), NOW(), $2, NOW() + INTERVAL '60 days') 
RETURNING token;

-- name: GetUserFromRefreshToken :one
SELECT token, created_at, updated_at, user_id, expires_at, revoked_at
FROM refresh_token 
WHERE token=$1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_token
SET revoked_at=NOW(), updated_at=NOW()
WHERE token=$1;