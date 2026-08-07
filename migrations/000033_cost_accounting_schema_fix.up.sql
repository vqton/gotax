-- Add missing columns for cost accounting hierarchy and pool classification

ALTER TABLE cost_objects ADD COLUMN IF NOT EXISTS parent_id TEXT;
ALTER TABLE cost_objects ADD COLUMN IF NOT EXISTS gl_account_code TEXT;

ALTER TABLE cost_pools ADD COLUMN IF NOT EXISTS code TEXT;
ALTER TABLE cost_pools ADD COLUMN IF NOT EXISTS pool_type TEXT DEFAULT 'OVERHEAD';
ALTER TABLE cost_pools ADD COLUMN IF NOT EXISTS allocation_base TEXT;
