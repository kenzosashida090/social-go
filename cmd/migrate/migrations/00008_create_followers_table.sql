-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS followers (
  user_id bigint NOT NULL,
  follower_id bigint NOT NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),

  PRIMARY KEY(user_id, follower_id),
  FOREIGN KEY(user_id) REFERENCES users(id),
  FOREIGN KEY(follower_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE followers;
-- +goose StatementEnd
