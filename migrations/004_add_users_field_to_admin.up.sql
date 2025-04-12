-- Add users count field to admin table
ALTER TABLE admin ADD COLUMN IF NOT EXISTS users INTEGER DEFAULT 0;