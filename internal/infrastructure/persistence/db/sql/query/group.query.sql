-- name: Ping :one
SELECT 1;

-- name: CreateGroup :one
INSERT INTO groups (
    id,
    name,
    description,
    owner_id,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    NOW(),
    NOW()
)
RETURNING *;


-- name: CountGroupsByOwner :one
SELECT COUNT(*) 
FROM groups
WHERE owner_id = $1;


-- name: CreateGroupMember :exec
INSERT INTO group_members (
    id,
    group_id,
    user_id,
    role,
    joined_at
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    now()
);

-- name: GetRoleByGroupIDAndUserID :one
SELECT role
FROM group_members
WHERE group_id = $1
AND user_id = $2;


-- name: GetGroupByID :one
SELECT *
FROM groups
WHERE id = $1;


-- name: CountGroupMembersByGroupID :one
SELECT COUNT(*)
FROM group_members
WHERE group_id = $1;


-- name: GetSprintByGroupID :one
SELECT *
FROM sprints
WHERE group_id = $1
ORDER BY created_at DESC
LIMIT 1; 