DROP INDEX IF EXISTS idx_media_issue_id;
ALTER TABLE media DROP COLUMN IF EXISTS issue_id;
