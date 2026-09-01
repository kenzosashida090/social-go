-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE users
  ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE users
   DROP COLUMN is_active;
-- +goose StatementEnd
