CREATE TABLE password_resets (
    token TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    expired_at TIMESTAMPTZ NOT NULL,
    update_time TIMESTAMPTZ DEFAULT NOW()
)