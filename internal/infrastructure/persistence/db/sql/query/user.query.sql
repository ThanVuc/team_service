-- name: GetUserByID :one 
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserWithPermissionByID :one
SELECT
    u.id,
    u.email,
    u.status,
    u.created_at,
    gm.group_id,
    gm.role,
    gm.joined_at
FROM users u
JOIN group_members gm
    ON gm.user_id = u.id
WHERE gm.group_id = $1
AND u.id = $2
LIMIT 1;

-- name: UpsertUser :exec
INSERT INTO users (
    id,
    email,
    status,
    created_at,
    avatar_url,
    has_email_notification,
    has_push_notification
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
ON CONFLICT (id) DO UPDATE
SET
    email = EXCLUDED.email,
    status = EXCLUDED.status,
    avatar_url = EXCLUDED.avatar_url,
    has_email_notification = EXCLUDED.has_email_notification,
    has_push_notification = EXCLUDED.has_push_notification;
