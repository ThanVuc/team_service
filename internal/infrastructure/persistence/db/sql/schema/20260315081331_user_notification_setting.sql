-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN has_email_notification BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN has_push_notification BOOLEAN NOT NULL DEFAULT true;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN has_email_notification,
DROP COLUMN has_push_notification;
-- +goose StatementEnd
