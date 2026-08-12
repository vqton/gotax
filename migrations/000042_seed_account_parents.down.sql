-- Revert 000042: restore pre-042 state (parent_code unset for all accounts).
UPDATE accounts SET parent_code = NULL WHERE parent_code IS NOT NULL;
