-- Migration to add bot_token and bot_chat_id columns to the admin table

-- Check if columns don't exist before adding them
DO $$
BEGIN
    -- Add bot_token column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                 WHERE table_name='admin' AND column_name='bot_token') THEN
        ALTER TABLE admin ADD COLUMN bot_token VARCHAR(255) DEFAULT '';
    END IF;

    -- Add bot_chat_id column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                 WHERE table_name='admin' AND column_name='bot_chat_id') THEN
        ALTER TABLE admin ADD COLUMN bot_chat_id VARCHAR(255) DEFAULT '';
    END IF;
END
$$;