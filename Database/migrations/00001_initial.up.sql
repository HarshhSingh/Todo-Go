BEGIN;

-- CREATE TABLE IF NOT EXISTS users (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     name TEXT NOT NULL,
--     email TEXT NOT NULL,
--     password TEXT NOT NULL,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
--     archived_at TIMESTAMP WITH TIME ZONE
-- );
-- CREATE UNIQUE INDEX IF NOT EXISTS active_user ON users(TRIM(LOWER(email))) WHERE archived_at IS NULL;
--
-- CREATE TYPE role_type AS ENUM (
--     'admin',
--     'user'
-- );
--
-- CREATE TABLE IF NOT EXISTS user_roles (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     user_id UUID REFERENCES users(id) NOT NULL,
--     role role_type NOT NULL,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
--     archived_at TIMESTAMP WITH TIME ZONE
-- );
-- CREATE UNIQUE INDEX IF NOT EXISTS unique_role ON user_roles(user_id, role) WHERE archived_at IS NULL;
--
-- CREATE TABLE IF NOT EXISTS user_session (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     user_id UUID REFERENCES users(id) NOT NULL,
--     session_token TEXT NOT NULL,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
-- );
--
-- CREATE TABLE IF NOT EXISTS todo (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     user_id UUID REFERENCES users(id) NOT NULL,
--     name TEXT NOT NULL,
--     description TEXT NOT NULL,
--     pending_at TIMESTAMP WITH TIME ZONE,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
--     archived_at TIMESTAMP WITH TIME ZONE
-- );
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task
        TEXT
        NOT NULL,
    task_status
        TEXT
        NOT NULL,
    date
        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    due_date
        TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);
COMMIT;