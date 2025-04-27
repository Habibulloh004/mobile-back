-- Add delivery field to admin table
ALTER TABLE admin ADD COLUMN IF NOT EXISTS delivery INTEGER DEFAULT 0;

-- Set NOT NULL constraint
ALTER TABLE admin ALTER COLUMN delivery SET DEFAULT 0;

-- Update any existing records to have the default value
UPDATE admin SET delivery = 0 WHERE delivery IS NULL;