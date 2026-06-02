CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    address VARCHAR(42),
    username VARCHAR(64),
    password_hash VARCHAR(256),
    name VARCHAR(128),
    nickname VARCHAR(128),
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    status VARCHAR(16) NOT NULL DEFAULT 'active',
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_address ON users (address) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users (username) WHERE deleted_at IS NULL;
