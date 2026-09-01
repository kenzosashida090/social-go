-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
  CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    level int NOT NULL DEFAULT 0,
    description TEXT 
  );
  INSERT INTO roles (
    name,
    description,
    level
  ) VALUES ( 
    'user',
    'A user can create posts and comments',
    1
  );   
   INSERT INTO roles (
    name,
    description,
    level
  ) VALUES ( 
    'moderator',
    'A moderator can update any comments',
    1
  );   
 INSERT INTO roles (
    name,
    description,
    level
  ) VALUES ( 
    'admin',
    'An admin can delete update posts',
    1
  );    

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
