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


-- name: UpdateGroup :one
UPDATE groups
SET
  name = COALESCE(sqlc.narg('name'), name),
  description = COALESCE(sqlc.narg('description'), description),
  updated_at = NOW()
WHERE id = sqlc.arg('id') AND deleted_at IS NULL
RETURNING *;

-- name: CheckGroupExists :one
SELECT EXISTS (
    SELECT 1
    FROM groups
    WHERE id = $1 AND deleted_at IS NULL
);


-- name: DeleteGroup :exec 
UPDATE groups 
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: CountManagerAndMemberByGroupID :one
SELECT COUNT(*)
FROM group_members
WHERE group_id = $1 AND role IN ('manager', 'member');


-- name: UpdateRoleMember :one
UPDATE group_members gm
SET role = $1
FROM users u
WHERE gm.group_id = $2
AND gm.user_id = $3
AND gm.user_id = u.id
RETURNING 
    u.id,
    u.email,
    u.avatar_url,
    gm.role,
    gm.joined_at;

-- name: RemoveMember :exec
DELETE FROM group_members
WHERE group_id = $1
AND user_id = $2;

-- name: CheckMemberExistsByEmail :one
SELECT EXISTS (
    SELECT 1
    FROM group_members gm
    JOIN users u ON gm.user_id = u.id
    WHERE gm.group_id = $1
    AND u.email = $2
);

-- name: GetSimpleUserByGroupID :many
SELECT u.id, u.email, u.avatar_url
FROM group_members gm
JOIN users u ON gm.user_id = u.id
WHERE gm.group_id = $1
AND gm.role IN ('manager', 'member', 'owner');