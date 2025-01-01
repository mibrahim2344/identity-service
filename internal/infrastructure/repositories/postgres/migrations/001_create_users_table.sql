-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TYPE user_status AS ENUM ('active', 'inactive', 'pending');
CREATE TYPE user_role AS ENUM ('admin', 'user');

CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    status user_status NOT NULL DEFAULT 'pending',
    role user_role NOT NULL DEFAULT 'user',
    email_verified BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    last_login_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_created_at ON users(created_at);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS user_status;
DROP TYPE IF EXISTS user_role;
