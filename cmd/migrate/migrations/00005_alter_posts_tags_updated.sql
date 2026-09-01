-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE posts
ADD tags VARCHAR(100)[];

ALTER TABLE posts
ADD updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(); 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE posts DROP COLUMN tags;
ALTER TABLE posts DROP COLUMN updated_at;
-- +goose StatementEnd
