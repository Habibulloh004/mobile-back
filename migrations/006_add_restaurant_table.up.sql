-- Create restaurant table
CREATE TABLE IF NOT EXISTS restaurant (
    id SERIAL PRIMARY KEY,
    admin_id INTEGER NOT NULL REFERENCES admin(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    contacts JSONB NOT NULL DEFAULT '{"phone":"", "gmail":"", "location":""}',
    social_media JSONB NOT NULL DEFAULT '{"instagram":"", "telegram":"", "facebook":""}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_admin_restaurant UNIQUE (admin_id)
);

-- Create trigger for updating timestamp
CREATE TRIGGER update_restaurant_timestamp BEFORE UPDATE ON restaurant 
FOR EACH ROW EXECUTE PROCEDURE update_timestamp();

-- Add index for faster lookup by admin_id
CREATE INDEX idx_restaurant_admin_id ON restaurant(admin_id);