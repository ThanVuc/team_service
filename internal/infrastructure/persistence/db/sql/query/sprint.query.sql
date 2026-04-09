-- name: DeleteDraftSprintsByGroupID :exec
DELETE FROM sprints
WHERE id = $1 AND status = 'draft';

-- name: CompleteActiveSprintsByGroupID :exec
UPDATE sprints
SET status = 'completed', updated_at = NOW()
WHERE id = $1 AND status = 'active';

-- name: CancelActiveSprintsByGroupID :exec
UPDATE sprints
SET status = 'canceled', updated_at = NOW()
WHERE id = $1 AND status = 'active';

-- name: GetSprintsByGroupID :many
SELECT *
FROM sprints
WHERE group_id = $1
ORDER BY created_at DESC;

-- name: IsSprintOverlap :one
SELECT EXISTS (
    SELECT 1
    FROM sprints
    WHERE group_id = $1
      AND status != 'cancelled'
      AND (
            daterange(start_date, end_date, '[]')
            && daterange($2::date, $3::date, '[]')
          )
) AS is_overlap;

-- name: CreateSprint :one
INSERT INTO sprints (
    id,
    group_id,
    name,
    goal,
    start_date,
    end_date,
    status,
    velocity_work,
    velocity_estimate,
    work_deleted,
    created_at,
    updated_at
) VALUES (
    $1, -- id
    $2, -- group_id
    $3, -- name
    $4, -- goal
    $5, -- start_date
    $6, -- end_date
    'draft',
    0,
    0,
    0,
    NOW(),
    NOW()
)
RETURNING
    id,
    group_id,
    name,
    goal,
    start_date,
    end_date,
    status,
    velocity_work,
    velocity_estimate,
    work_deleted,
    created_at,
    updated_at;

-- name: GetSprintByID :one
SELECT *
FROM sprints
WHERE id = $1;

-- name: UpdateSprint :one
UPDATE sprints
SET
        name = COALESCE(sqlc.narg('name'), name),
        goal = COALESCE(sqlc.narg('goal'), goal),
        start_date = COALESCE(sqlc.narg('start_date'), start_date),
        end_date = COALESCE(sqlc.narg('end_date'), end_date),
        updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateSprintStatus :one
UPDATE sprints
SET
        status = sqlc.arg('status'),
        updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteSprint :execrows
WITH target_sprint AS (
    SELECT sprints.id
    FROM sprints
    WHERE sprints.id = $1
      AND sprints.status = 'draft'
), moved_works AS (
    UPDATE works
    SET sprint_id = NULL,
        updated_at = NOW()
    WHERE works.sprint_id IN (SELECT target_sprint.id FROM target_sprint)
)
DELETE FROM sprints
WHERE sprints.id IN (SELECT target_sprint.id FROM target_sprint);

-- name: GetSimpleSprintsByGroupID :many
SELECT id, name, status
FROM sprints
WHERE group_id = $1;

-- name: DeleteDraftSprintByID :execrows
WITH target_sprint AS (
    SELECT sprints.id
    FROM sprints
    WHERE sprints.id = $1
      AND sprints.status = 'draft'
), deleted_works AS (
    DELETE FROM works
    WHERE works.sprint_id IN (SELECT target_sprint.id FROM target_sprint)
)
DELETE FROM sprints
WHERE sprints.id IN (SELECT target_sprint.id FROM target_sprint);