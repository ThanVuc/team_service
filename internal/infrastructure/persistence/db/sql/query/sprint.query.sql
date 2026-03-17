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
