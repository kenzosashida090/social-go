-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE user_invitation
  ADD COLUMN expiry TIMESTAMP(0) WITH TIME ZONE NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE user_invitation
  DROP COLUMN expiry;
-- +goose StatementEnd
