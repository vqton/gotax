-- Align journal_entries.company_id with VARCHAR companies.id.
-- Migration 000001 declared it UUID (legacy design); companies.id is VARCHAR(20) ("CMP001"),
-- so every journal write on PG failed with a uuid parse error.
ALTER TABLE journal_entries ALTER COLUMN company_id DROP DEFAULT;
ALTER TABLE journal_entries ALTER COLUMN company_id TYPE VARCHAR(20);
