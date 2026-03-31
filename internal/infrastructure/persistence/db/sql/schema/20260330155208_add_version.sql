-- +goose Up
-- +goose StatementBegin
ALTER TABLE works ADD COLUMN version INTEGER NOT NULL DEFAULT 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE works DROP COLUMN version;
-- +goose StatementEnd
