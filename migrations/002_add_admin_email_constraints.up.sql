-- Add email column to admin table if it doesn't exist
ALTER TABLE admin ADD COLUMN IF NOT EXISTS email VARCHAR(100);

-- Create a temporary table to identify duplicates
CREATE TEMP TABLE duplicate_admins AS
SELECT 
    id,
    user_name,
    system_id,
    ROW_NUMBER() OVER (PARTITION BY user_name, system_id ORDER BY id) as row_num
FROM 
    admin;

-- Display duplicates (for logging)
SELECT id, user_name, system_id 
FROM duplicate_admins 
WHERE row_num > 1;

-- Update duplicate rows to make the combination unique by appending a suffix to system_id
UPDATE admin
SET system_id = admin.system_id || '_' || t.row_num 
FROM duplicate_admins t
WHERE admin.id = t.id AND t.row_num > 1;

-- Now make sure all emails are unique by setting default unique values
UPDATE admin
SET email = CONCAT('admin_', admin.id, '@example.com')
WHERE email IS NULL OR email = '';

-- Deduplicate emails if any existing duplicates
WITH duplicate_emails AS (
    SELECT 
        id,
        email,
        ROW_NUMBER() OVER (PARTITION BY email ORDER BY id) as row_num
    FROM 
        admin
    WHERE 
        email IS NOT NULL
)
UPDATE admin
SET email = CONCAT(admin.email, '.', d.row_num) 
FROM duplicate_emails d
WHERE admin.id = d.id AND d.row_num > 1;

-- Add NOT NULL constraint
ALTER TABLE admin ALTER COLUMN email SET NOT NULL;

-- Add unique constraints
ALTER TABLE admin ADD CONSTRAINT admin_email_unique UNIQUE (email);
ALTER TABLE admin ADD CONSTRAINT admin_username_systemid_unique UNIQUE (user_name, system_id);

-- Drop the temporary table
DROP TABLE duplicate_admins;