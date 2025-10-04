CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    permission_name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT
);