CREATE TABLE sessions (
    id TEXT NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users (id),
    expires_at TIMESTAMPTZ NOT NULL,

    PRIMARY KEY (id)
);