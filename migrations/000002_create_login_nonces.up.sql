CREATE TABLE IF NOT EXISTS login_nonces (
    id BIGSERIAL PRIMARY KEY,
    address VARCHAR(42) NOT NULL,
    nonce VARCHAR(64) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_login_nonces_address ON login_nonces (address);
CREATE INDEX IF NOT EXISTS idx_login_nonces_expires_at ON login_nonces (expires_at);
