BEGIN;

CREATE TABLE IF NOT EXISTS regisuser (
    userid uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX IF NOT EXISTS active_user ON regisuser(TRIM(LOWER(email))) WHERE archived_at IS NULL;

CREATE TABLE IF NOT EXISTS todos (
    todo_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    id uuid REFERENCES regisuser(userid) NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    is_completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS User_sessions (
    session_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id UUID REFERENCES regisuser(userid) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE
);

COMMIT;