-- Create super_admin table
CREATE TABLE IF NOT EXISTS super_admin (
    id SERIAL PRIMARY KEY,
    login VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create admin table
CREATE TABLE IF NOT EXISTS admin (
    id SERIAL PRIMARY KEY,
    user_name VARCHAR(100) NOT NULL,
    company_name VARCHAR(100) NOT NULL,
    system_id VARCHAR(100) NOT NULL,
    system_token VARCHAR(255) NOT NULL,
    system_token_updated_time TIMESTAMP WITH TIME ZONE NOT NULL,
    sms_token VARCHAR(255) NOT NULL,
    sms_email VARCHAR(100) NOT NULL,
    sms_password VARCHAR(255) NOT NULL,
    sms_message TEXT NOT NULL,
    payment_username VARCHAR(100) NOT NULL,
    payment_password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create banner table
CREATE TABLE IF NOT EXISTS banner (
    id SERIAL PRIMARY KEY,
    admin_id INTEGER NOT NULL REFERENCES admin(id) ON DELETE CASCADE,
    image VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create notification table
CREATE TABLE IF NOT EXISTS notification (
    id SERIAL PRIMARY KEY,
    admin_id INTEGER NOT NULL REFERENCES admin(id) ON DELETE CASCADE,
    payload TEXT NOT NULL,
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create fcm_token table
CREATE TABLE IF NOT EXISTS fcm_token (
    id SERIAL PRIMARY KEY,
    admin_id INTEGER NOT NULL REFERENCES admin(id) ON DELETE CASCADE,
    fcm_token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = CURRENT_TIMESTAMP;
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for updating updated_at timestamp
CREATE TRIGGER update_super_admin_timestamp BEFORE UPDATE ON super_admin FOR EACH ROW EXECUTE PROCEDURE update_timestamp();
CREATE TRIGGER update_admin_timestamp BEFORE UPDATE ON admin FOR EACH ROW EXECUTE PROCEDURE update_timestamp();
CREATE TRIGGER update_banner_timestamp BEFORE UPDATE ON banner FOR EACH ROW EXECUTE PROCEDURE update_timestamp();
CREATE TRIGGER update_notification_timestamp BEFORE UPDATE ON notification FOR EACH ROW EXECUTE PROCEDURE update_timestamp();
CREATE TRIGGER update_fcm_token_timestamp BEFORE UPDATE ON fcm_token FOR EACH ROW EXECUTE PROCEDURE update_timestamp();

-- Insert default super admin
-- Password will be hashed in the application code
INSERT INTO super_admin (login, password) 
VALUES ('superadmin', '$2a$10$YourHashedPasswordWillBeSetInCode') 
ON CONFLICT (login) DO NOTHING;