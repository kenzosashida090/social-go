-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE posts
ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE posts
DROP CONSTRAINT fk_user
-- +goose StatementEnd
