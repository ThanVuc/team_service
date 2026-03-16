-- +goose Up
-- +goose StatementBegin
ALTER TABLE groups
ADD COLUMN IF NOT EXISTS avatar_url VARCHAR(128);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE groups
DROP COLUMN IF EXISTS avatar_url;
-- +goose StatementEnd