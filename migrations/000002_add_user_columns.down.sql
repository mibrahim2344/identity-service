-- Remove added columns from users table
ALTER TABLE users
DROP COLUMN IF EXISTS status,
DROP COLUMN IF EXISTS deleted_at;
