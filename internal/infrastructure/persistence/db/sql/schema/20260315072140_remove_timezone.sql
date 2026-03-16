-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN IF EXISTS time_zone;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN time_zone VARCHAR(64);
-- +goose StatementEnd
