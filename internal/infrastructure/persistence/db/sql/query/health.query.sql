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


-- name: GetUserByID :one 
SELECT id, email, avatar_url
FROM users
WHERE id = $1;


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
SELECT id, name, description, avatar_url, owner_id, created_at, updated_at
FROM groups
WHERE id = $1;


-- name: CountGroupMembersByGroupID :one
SELECT COUNT(*)
FROM group_members
WHERE group_id = $1;


-- name: GetSprintByGroupID :one
SELECT id, name,group_id,goal, start_date, end_date,velocity_work,velocity_estimate,work_deleted
FROM sprints
WHERE group_id = $1
ORDER BY created_at DESC
LIMIT 1; 