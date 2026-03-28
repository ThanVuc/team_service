-- name: CreateInvite :one
INSERT INTO invites (
    id,
    group_id,
    token,
    role,
    email,
    expires_at,
    created_by,
    created_at
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    NOW()
)
RETURNING token, expires_at, created_at;


-- name: CheckPendingInvite :one
SELECT EXISTS(
    SELECT 1
    FROM invites
    WHERE group_id = $1
    AND email = $2
    AND expires_at > NOW()
);


-- name: GetInviteByToken :one
SELECT id, group_id, token, role, email, expires_at, created_by, created_at
FROM invites
WHERE token = $1
AND expires_at > NOW();