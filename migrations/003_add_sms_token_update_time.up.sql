-- Add sms_token_updated_time column if it doesn't exist
ALTER TABLE admin ADD COLUMN IF NOT EXISTS sms_token_updated_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;

-- Set NOT NULL constraint
ALTER TABLE admin ALTER COLUMN sms_token_updated_time SET NOT NULL;