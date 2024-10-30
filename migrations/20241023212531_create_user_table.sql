-- +goose Up
CREATE TYPE enum_role AS ENUM ('USER', 'ADMIN');

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    role enum_role NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

-- +goose Down
DROP TYPE enum_role;
DROP TABLE IF EXISTS users;
