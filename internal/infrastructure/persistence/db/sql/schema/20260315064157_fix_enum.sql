-- +goose Up
-- +goose StatementBegin

BEGIN;

-- =========================
-- USERS STATUS FIX
-- =========================

ALTER TABLE users
DROP CONSTRAINT IF EXISTS users_status_check;

ALTER TABLE users
ADD CONSTRAINT users_status_check
CHECK (status IN ('active', 'inactive'));

-- =========================
-- GROUP MEMBERS ROLE FIX
-- =========================

ALTER TABLE group_members
DROP CONSTRAINT IF EXISTS group_members_role_check;

ALTER TABLE group_members
ADD CONSTRAINT group_members_role_check
CHECK (role IN ('owner','manager','member','viewer'));

-- =========================
-- INVITES ROLE FIX
-- =========================

ALTER TABLE invites
DROP CONSTRAINT IF EXISTS invites_role_check;

ALTER TABLE invites
ADD CONSTRAINT invites_role_check
CHECK (role IN ('owner','manager','member','viewer'));

-- =========================
-- SPRINT STATUS FIX
-- =========================

ALTER TABLE sprints
DROP CONSTRAINT IF EXISTS sprints_status_check;

ALTER TABLE sprints
ADD CONSTRAINT sprints_status_check
CHECK (status IN ('draft','active','completed','cancelled'));

-- =========================
-- WORK STATUS FIX
-- =========================

ALTER TABLE works
DROP CONSTRAINT IF EXISTS works_status_check;

ALTER TABLE works
ADD CONSTRAINT works_status_check
CHECK (status IN ('todo','inprogress','inreview','done'));

-- =========================
-- WORK PRIORITY FIX
-- =========================

ALTER TABLE works
DROP CONSTRAINT IF EXISTS works_priority_check;

ALTER TABLE works
ADD CONSTRAINT works_priority_check
CHECK (priority IN ('low','medium','high') OR priority IS NULL);

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

BEGIN;

-- USERS
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_status_check;
ALTER TABLE users ADD CONSTRAINT users_status_check
CHECK (status IN ('active', 'inactive', 'suspended'));

-- GROUP MEMBERS
ALTER TABLE group_members DROP CONSTRAINT IF EXISTS group_members_role_check;
ALTER TABLE group_members ADD CONSTRAINT group_members_role_check
CHECK (role IN ('owner','admin','member'));

-- INVITES
ALTER TABLE invites DROP CONSTRAINT IF EXISTS invites_role_check;
ALTER TABLE invites ADD CONSTRAINT invites_role_check
CHECK (role IN ('owner','admin','member'));

-- SPRINTS
ALTER TABLE sprints DROP CONSTRAINT IF EXISTS sprints_status_check;
ALTER TABLE sprints ADD CONSTRAINT sprints_status_check
CHECK (status IN ('planned','active','completed'));

-- WORKS STATUS
ALTER TABLE works DROP CONSTRAINT IF EXISTS works_status_check;
ALTER TABLE works ADD CONSTRAINT works_status_check
CHECK (status IN ('todo','in_progress','done'));

-- WORKS PRIORITY
ALTER TABLE works DROP CONSTRAINT IF EXISTS works_priority_check;
ALTER TABLE works ADD CONSTRAINT works_priority_check
CHECK (priority IN ('low','medium','high','urgent') OR priority IS NULL);

COMMIT;

-- +goose StatementEnd
