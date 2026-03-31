-- name: CreateWork :one
INSERT INTO works (
	id,
	group_id,
	sprint_id,
	name,
	description,
	status,
	assignee_id,
	creator_id,
	estimate_hours,
	story_point,
	priority,
	due_date,
	created_at,
	updated_at
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	'todo',
	NULL,
	$6,
	NULL,
	NULL,
	NULL,
	NULL,
	NOW(),
	NOW()
)
RETURNING *;

-- name: UpdateWork :one
UPDATE works
SET
	name = COALESCE(sqlc.narg('name'), name),
	description = COALESCE(sqlc.narg('description'), description),
	sprint_id = COALESCE(sqlc.narg('sprint_id'), sprint_id),
	assignee_id = COALESCE(sqlc.narg('assignee_id'), assignee_id),
	status = COALESCE(sqlc.narg('status'), status),
	story_point = COALESCE(sqlc.narg('story_point'), story_point),
	due_date = COALESCE(sqlc.narg('due_date'), due_date),
	priority = COALESCE(sqlc.narg('priority'), priority),
	version = version + 1,
	updated_at = NOW()
WHERE id = sqlc.arg('id') AND version = sqlc.arg('version')
RETURNING
	id,
	CASE WHEN sqlc.narg('name')::text IS NOT NULL THEN name ELSE NULL END AS name,
	CASE WHEN sqlc.narg('description')::text IS NOT NULL THEN description ELSE NULL END AS description,
	CASE WHEN sqlc.narg('sprint_id')::uuid IS NOT NULL THEN sprint_id ELSE NULL END AS sprint_id,
	CASE WHEN sqlc.narg('assignee_id')::uuid IS NOT NULL THEN assignee_id ELSE NULL END AS assignee_id,
	CASE WHEN sqlc.narg('status')::text IS NOT NULL THEN status ELSE NULL END AS status,
	CASE WHEN sqlc.narg('story_point')::int IS NOT NULL THEN story_point ELSE NULL END AS story_point,
	CASE WHEN sqlc.narg('due_date')::date IS NOT NULL THEN due_date ELSE NULL END AS due_date,
	CASE WHEN sqlc.narg('priority')::text IS NOT NULL THEN priority ELSE NULL END AS priority,
	version,
	updated_at;

-- name: DeleteWork :one
WITH deleted AS (
	DELETE FROM works
	WHERE id = $1
	RETURNING id
)
SELECT
	EXISTS(SELECT 1 FROM deleted) AS success,
	NOW()::timestamptz AS deleted_at;

-- name: GetWorksBySprint :many
SELECT
    w.id,
    w.group_id,
    w.sprint_id,
    w.name,
    w.description,
    w.status,
    w.assignee_id,
    w.creator_id,
    w.estimate_hours,
    w.story_point,
    w.priority,
    w.due_date,
    w.created_at,
    w.updated_at,
	w.version,
    u.email AS assignee_email,
    u.avatar_url AS assignee_avatar_url
FROM works w
LEFT JOIN users u ON u.id = w.assignee_id
WHERE w.group_id = sqlc.arg('group_id')
AND (sqlc.narg('sprint_id')::uuid IS NULL OR w.sprint_id = sqlc.narg('sprint_id'))
AND (sqlc.narg('assignee_id')::uuid IS NULL OR w.assignee_id = sqlc.narg('assignee_id'))
ORDER BY w.updated_at DESC;

-- name: GetWorksBySprintWithoutAggregation :many
SELECT *
FROM works
WHERE group_id = sqlc.arg('group_id')
AND sprint_id IS NOT DISTINCT FROM sqlc.narg('sprint_id')
ORDER BY
completed_at ASC NULLS LAST;

-- name: GetWork :one
SELECT
	w.id,
	w.group_id,
	w.sprint_id,
	w.name,
	w.description,
	w.status,
	w.assignee_id,
	w.creator_id,
	w.estimate_hours,
	w.story_point,
	w.priority,
	w.due_date,
	w.created_at,
	w.updated_at,
	w.version,
	s.name AS sprint_name,
	u.email AS assignee_email,
	u.avatar_url AS assignee_avatar_url
FROM works w
LEFT JOIN sprints s ON s.id = w.sprint_id
LEFT JOIN users u ON u.id = w.assignee_id
WHERE w.id = $1;

-- name: GetCheckListByWorkId :many
SELECT
	id,
	work_id,
	name,
	is_completed,
	created_at,
	updated_at
FROM checklist_items
WHERE work_id = $1
ORDER BY created_at ASC;

-- name: GetChecklistItemMeta :one
SELECT
	id,
	work_id
FROM checklist_items
WHERE id = $1;

-- name: GetCommentsByWorkId :many
SELECT
	c.id,
	c.work_id,
	c.creator_id,
	c.content,
	c.created_at,
	c.updated_at,
	u.email AS creator_email,
	u.avatar_url AS creator_avatar_url
FROM comments c
JOIN users u ON u.id = c.creator_id
WHERE c.work_id = $1
ORDER BY c.created_at ASC;

-- name: GetCommentMeta :one
SELECT
	id,
	work_id,
	creator_id
FROM comments
WHERE id = $1;

-- name: CreateChecklistItem :one
INSERT INTO checklist_items (
	id,
	work_id,
	name,
	is_completed,
	created_at,
	updated_at
) VALUES (
	$1,
	$2,
	$3,
	false,
	NOW(),
	NOW()
)
RETURNING *;

-- name: UpdateChecklistItem :one
UPDATE checklist_items
SET
	name = COALESCE(sqlc.narg('name'), name),
	is_completed = COALESCE(sqlc.narg('is_completed'), is_completed),
	updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING
	id,
	CASE WHEN sqlc.narg('name')::text IS NOT NULL THEN name ELSE NULL END AS name,
	CASE WHEN sqlc.narg('is_completed')::boolean IS NOT NULL THEN is_completed ELSE NULL END AS is_completed,
	updated_at;

-- name: DeleteChecklistItem :one
WITH deleted AS (
	DELETE FROM checklist_items
	WHERE id = $1
	RETURNING id
)
SELECT
	EXISTS(SELECT 1 FROM deleted) AS success,
	(SELECT id FROM deleted) AS id,
	NOW()::timestamptz AS deleted_at;

-- name: CreateComment :one
WITH inserted AS (
	INSERT INTO comments (
		id,
		work_id,
		creator_id,
		content,
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
	RETURNING
		id,
		work_id,
		creator_id,
		content,
		created_at,
		updated_at
)
SELECT
	i.id,
	i.work_id,
	i.creator_id,
	i.content,
	i.created_at,
	i.updated_at,
	u.email AS creator_email,
	u.avatar_url AS creator_avatar_url
FROM inserted i
JOIN users u ON u.id = i.creator_id;

-- name: UpdateComment :one
UPDATE comments
SET
	content = $1,
	updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: DeleteComment :one
WITH deleted AS (
	DELETE FROM comments
	WHERE id = $1
	RETURNING id
)
SELECT
	EXISTS(SELECT 1 FROM deleted) AS success,
	NOW()::timestamptz AS deleted_at;

-- name: UnassignWorksByMember :exec
UPDATE works
SET assignee_id = NULL
WHERE group_id = $1
AND assignee_id = $2;