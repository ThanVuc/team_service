-- +goose Up
-- +goose StatementBegin

BEGIN;

-- =========================
-- USERS
-- =========================
CREATE TABLE users (
id UUID PRIMARY KEY,
email VARCHAR(255) NOT NULL UNIQUE,
status VARCHAR(20) NOT NULL,
time_zone VARCHAR(64) NOT NULL,
created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
CHECK (status IN ('active', 'inactive', 'suspended'))
);

CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);

-- =========================
-- GROUPS
-- =========================
CREATE TABLE groups (
id UUID PRIMARY KEY,
name VARCHAR(255) NOT NULL,
description TEXT,
owner_id UUID NOT NULL REFERENCES users(id),
created_at TIMESTAMPTZ NOT NULL,
updated_at TIMESTAMPTZ NOT NULL,
deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_groups_owner_id ON groups(owner_id);
CREATE INDEX idx_groups_deleted_at ON groups(deleted_at) WHERE deleted_at IS NULL;

-- =========================
-- GROUP MEMBERS
-- =========================
CREATE TABLE group_members (
id UUID PRIMARY KEY,
group_id UUID NOT NULL REFERENCES groups(id),
user_id UUID NOT NULL REFERENCES users(id),
role VARCHAR(20) NOT NULL,
joined_at TIMESTAMPTZ NOT NULL,
CHECK (role IN ('owner','admin','member'))
);

CREATE INDEX idx_group_members_group_id ON group_members(group_id);
CREATE INDEX idx_group_members_user_id ON group_members(user_id);
CREATE UNIQUE INDEX idx_group_members_group_user ON group_members(group_id, user_id);
CREATE INDEX idx_group_members_role ON group_members(group_id, role);

-- =========================
-- INVITES
-- =========================
CREATE TABLE invites (
id UUID PRIMARY KEY,
group_id UUID NOT NULL REFERENCES groups(id),
token VARCHAR(255) NOT NULL UNIQUE,
role VARCHAR(20) NOT NULL,
email VARCHAR(255),
expires_at TIMESTAMPTZ NOT NULL,
created_by UUID NOT NULL REFERENCES users(id),
created_at TIMESTAMPTZ NOT NULL,
CHECK (role IN ('owner','admin','member'))
);

CREATE UNIQUE INDEX idx_invites_token ON invites(token);
CREATE INDEX idx_invites_group_id ON invites(group_id);
CREATE INDEX idx_invites_email ON invites(email);
CREATE INDEX idx_invites_expires_at ON invites(expires_at);

-- =========================
-- SPRINTS
-- =========================
CREATE TABLE sprints (
id UUID PRIMARY KEY,
group_id UUID NOT NULL REFERENCES groups(id),
name VARCHAR(255) NOT NULL,
goal TEXT,
start_date DATE NOT NULL,
end_date DATE NOT NULL,
status VARCHAR(20) NOT NULL,
velocity_work INT,
velocity_estimate FLOAT,
work_deleted INT,
created_at TIMESTAMPTZ NOT NULL,
updated_at TIMESTAMPTZ NOT NULL,
CHECK (status IN ('planned','active','completed'))
);

CREATE INDEX idx_sprints_group_id ON sprints(group_id);
CREATE INDEX idx_sprints_group_status ON sprints(group_id, status);
CREATE INDEX idx_sprints_active ON sprints(group_id) WHERE status='active';

-- =========================
-- WORKS
-- =========================
CREATE TABLE works (
id UUID PRIMARY KEY,
group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
sprint_id UUID REFERENCES sprints(id) ON DELETE SET NULL,
name VARCHAR(500) NOT NULL,
description TEXT,
status VARCHAR(20) NOT NULL,
assignee_id UUID REFERENCES users(id),
creator_id UUID NOT NULL REFERENCES users(id),
estimate_hours FLOAT CHECK (estimate_hours > 0),
story_point INT CHECK (story_point > 0),
priority VARCHAR(20),
due_date DATE,
created_at TIMESTAMPTZ NOT NULL,
updated_at TIMESTAMPTZ NOT NULL,
CHECK (status IN ('todo','in_progress','done')),
CHECK (priority IN ('low','medium','high','urgent') OR priority IS NULL)
);

CREATE INDEX idx_works_group_id ON works(group_id);
CREATE INDEX idx_works_sprint_id ON works(sprint_id);
CREATE INDEX idx_works_assignee_id ON works(assignee_id);
CREATE INDEX idx_works_group_status ON works(group_id, status);
CREATE INDEX idx_works_backlog ON works(group_id) WHERE sprint_id IS NULL;
CREATE INDEX idx_works_due_date ON works(due_date);

-- =========================
-- CHECKLIST ITEMS
-- =========================
CREATE TABLE checklist_items (
id UUID PRIMARY KEY,
work_id UUID NOT NULL REFERENCES works(id) ON DELETE CASCADE,
name VARCHAR(500) NOT NULL,
is_completed BOOLEAN NOT NULL DEFAULT false,
created_at TIMESTAMPTZ NOT NULL,
updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_checklist_work_id ON checklist_items(work_id);

-- =========================
-- COMMENTS
-- =========================
CREATE TABLE comments (
id UUID PRIMARY KEY,
work_id UUID NOT NULL REFERENCES works(id) ON DELETE CASCADE,
creator_id UUID NOT NULL REFERENCES users(id),
content TEXT NOT NULL,
created_at TIMESTAMPTZ NOT NULL,
updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_comments_work_id ON comments(work_id);
CREATE INDEX idx_comments_work_created ON comments(work_id, created_at);
CREATE INDEX idx_comments_creator_id ON comments(creator_id);

COMMIT;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS checklist_items;
DROP TABLE IF EXISTS works;
DROP TABLE IF EXISTS sprints;
DROP TABLE IF EXISTS invites;
DROP TABLE IF EXISTS group_members;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd
