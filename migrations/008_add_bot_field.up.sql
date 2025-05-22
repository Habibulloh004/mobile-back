-- Migration to add bot_token and bot_chat_id fields to admin table

-- Add bot_token column
ALTER TABLE admin 
ADD COLUMN bot_token VARCHAR(255) DEFAULT '';

-- Add bot_chat_id column  
ALTER TABLE admin 
ADD COLUMN bot_chat_id VARCHAR(255) DEFAULT '';

-- Add comments for documentation
COMMENT ON COLUMN admin.bot_token IS 'Telegram bot token for notifications';
COMMENT ON COLUMN admin.bot_chat_id IS 'Telegram chat ID for bot notifications';

-- Optional: Create index on bot_token if needed for performance
-- CREATE INDEX idx_admin_bot_token ON admin(bot_token) WHERE bot_token != '';

-- Optional: Create index on bot_chat_id if needed for performance  
-- CREATE INDEX idx_admin_bot_chat_id ON admin(bot_chat_id) WHERE bot_chat_id != '';