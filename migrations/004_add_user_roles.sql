-- +migrate Up
CREATE TYPE user_role AS ENUM ('user', 'admin', 'moderator');

ALTER TABLE users 
ADD COLUMN role user_role DEFAULT 'user',
ADD COLUMN last_login TIMESTAMP WITH TIME ZONE;

CREATE INDEX idx_users_role ON users(role);

-- +migrate Down
ALTER TABLE users 
DROP COLUMN role,
DROP COLUMN last_login;

DROP TYPE user_role;