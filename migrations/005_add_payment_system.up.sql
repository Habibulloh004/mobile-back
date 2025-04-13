-- Drop the existing subscription_tier table
DROP TABLE IF EXISTS subscription_tier CASCADE;

-- Create the subscription_tier table with the correct column names
CREATE TABLE subscription_tier (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,     -- This should be 'name', not 'tier_name'
    min_users INTEGER NOT NULL,
    max_users INTEGER,
    price DECIMAL(10, 2) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger for updating timestamp
CREATE TRIGGER update_subscription_tier_timestamp BEFORE UPDATE ON subscription_tier 
FOR EACH ROW EXECUTE PROCEDURE update_timestamp();

-- Insert default subscription tiers
INSERT INTO subscription_tier (name, min_users, max_users, price, description) VALUES
('Free', 0, 100, 0.00, 'Free tier for up to 100 users'),
('Basic', 101, 1000, 5.00, 'Basic tier for 101-1000 users'),
('Professional', 1001, 5000, 10.00, 'Professional tier for 1001-5000 users'),
('Enterprise', 5001, NULL, 20.00, 'Enterprise tier for 5001+ users');