-- name: UnassignWorksByMember :exec
UPDATE works
SET assignee_id = NULL
WHERE group_id = $1
AND assignee_id = $2;