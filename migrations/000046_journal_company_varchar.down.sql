-- Destructive: VARCHAR company ids ("CMP001") cannot cast back to UUID.
DELETE FROM journal_lines;
DELETE FROM journal_entries;
ALTER TABLE journal_entries ALTER COLUMN company_id TYPE UUID USING '00000000-0000-0000-0000-000000000000'::uuid;
ALTER TABLE journal_entries ALTER COLUMN company_id SET DEFAULT '00000000-0000-0000-0000-000000000000';
