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


