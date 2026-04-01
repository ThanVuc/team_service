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
WHERE owner_id = $1
AND deleted_at IS NULL;


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

-- name: GetGroupsByUserID :many
SELECT 
    g.id,
    g.name,
    gm.role AS my_role,
    gm_count.member_total,
    COALESCE(g.avatar_url, '') AS avatar_url,
    g.created_at,
    g.updated_at
FROM groups g
JOIN group_members gm ON g.id = gm.group_id
JOIN (
    SELECT group_id, COUNT(*)::bigint AS member_total
    FROM group_members
    GROUP BY group_id
) gm_count ON gm_count.group_id = g.id
WHERE gm.user_id = $1
AND g.deleted_at IS NULL
ORDER BY g.created_at DESC;

-- name: GetOwnerByGroupID :one
SELECT
    u.id AS owner_id,
    u.email AS owner_email,
    COALESCE(u.avatar_url, '') AS owner_image
FROM groups g
JOIN group_members gm ON gm.group_id = g.id
    AND gm.role = 'owner'
    AND gm.user_id = g.owner_id
JOIN users u ON u.id = g.owner_id
WHERE g.id = $1
AND g.deleted_at IS NULL;
-- name: CountViewerByGroupID :one
SELECT COUNT(*)
FROM group_members
WHERE group_id = $1 AND role = 'viewer';


-- name: GetListUserIDByGroupID :many
SELECT user_id
FROM group_members
WHERE group_id = $1;