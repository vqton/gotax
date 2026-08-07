ALTER TABLE cost_objects DROP COLUMN IF EXISTS parent_id;
ALTER TABLE cost_objects DROP COLUMN IF EXISTS gl_account_code;

ALTER TABLE cost_pools DROP COLUMN IF EXISTS code;
ALTER TABLE cost_pools DROP COLUMN IF EXISTS pool_type;
ALTER TABLE cost_pools DROP COLUMN IF EXISTS allocation_base;
