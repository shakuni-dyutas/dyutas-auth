-- name: CreateRefreshToken :one
INSERT INTO auth.user_refresh_tokens (
    user_id,
    token_hash,
    created_at,
    expires_at
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetUserRefreshTokenByTokenHash :one
SELECT * FROM auth.user_refresh_tokens 
WHERE user_id = $1
    AND token_hash = $2;

-- name: GetUserRefreshTokensExpiresAfter :many
SELECT * FROM auth.user_refresh_tokens 
WHERE user_id = $1 
    AND expires_at > $2
ORDER BY expires_at DESC;

-- name: RevokeUserRefreshToken :exec
UPDATE auth.user_refresh_tokens 
SET revoked_at = COALESCE($2, NOW())
WHERE token_hash = $1 AND revoked_at IS NULL;

-- name: RevokeUserAllRefreshTokensAfter :exec
UPDATE auth.user_refresh_tokens 
SET revoked_at = COALESCE($2, NOW())
WHERE user_id = $1 AND expires_at > COALESCE($2, NOW());

-- name: CleanupRefreshTokensExpiredBefore :exec
DELETE FROM auth.user_refresh_tokens 
WHERE expires_at < COALESCE($1, NOW());