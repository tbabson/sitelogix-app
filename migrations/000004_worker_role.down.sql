DROP INDEX IF EXISTS idx_workers_user_id;
ALTER TABLE workers DROP COLUMN IF EXISTS user_id;
-- NOTE: PostgreSQL does not support DROP VALUE from an enum; recreate if needed.
