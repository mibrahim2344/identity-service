-- Add missing columns to users table
ALTER TABLE users
ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'pending',
ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;
