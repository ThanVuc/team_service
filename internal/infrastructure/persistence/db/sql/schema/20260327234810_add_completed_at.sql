-- +goose Up
-- +goose StatementBegin
ALTER TABLE works
ADD COLUMN completed_at TIMESTAMPTZ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE works
DROP COLUMN completed_at;
-- +goose StatementEnd
