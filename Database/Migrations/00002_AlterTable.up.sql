BEGIN;

-- no need to make transaction for a single db statement
ALTER TABLE todos ADD COLUMN archieved_at TIMESTAMP WITH TIME ZONE;

COMMIT;