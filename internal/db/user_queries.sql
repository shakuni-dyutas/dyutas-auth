-- name: CreateUser :one
INSERT INTO auth.users (
    code,
    google_id,
    email,
    profile_image_url,
    username,
    signed_up_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetUserByGoogleID :one
SELECT * FROM auth.users 
WHERE google_id = $1 AND deleted_at IS NULL;

-- name: GetUserByCode :one
SELECT * FROM auth.users 
WHERE code = $1 AND deleted_at IS NULL;

-- name: UpdateUserProfileByCode :one
UPDATE auth.users 
SET 
    email = $2,
    profile_image_url = $3,
    username = $4
WHERE code = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteUserByCode :exec
UPDATE auth.users 
SET deleted_at = $2
WHERE code = $1 AND deleted_at IS NULL;

-- ========== User Snapshots Queries ==========

-- name: CreateUserSnapshot :one
INSERT INTO auth.user_snapshots (
    user_id,
    email,
    profile_image_url,
    username,
    created_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;