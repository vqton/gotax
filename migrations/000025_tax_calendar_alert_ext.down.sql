ALTER TABLE tax_calendar DROP COLUMN IF EXISTS period_type;
ALTER TABLE tax_calendar DROP COLUMN IF EXISTS start_date;
ALTER TABLE tax_calendar DROP COLUMN IF EXISTS end_date;
ALTER TABLE tax_calendar DROP COLUMN IF EXISTS declaration_due;
ALTER TABLE tax_calendar DROP COLUMN IF EXISTS payment_due;
ALTER TABLE tax_calendar DROP COLUMN IF EXISTS declaration_id;

ALTER TABLE tax_alerts DROP COLUMN IF EXISTS calendar_id;
ALTER TABLE tax_alerts DROP COLUMN IF EXISTS channel;
ALTER TABLE tax_alerts DROP COLUMN IF EXISTS acknowledged_at;
ALTER TABLE tax_alerts DROP COLUMN IF EXISTS acknowledged_by;

ALTER TABLE tax_audit_cases DROP COLUMN IF EXISTS audit_period_start;
ALTER TABLE tax_audit_cases DROP COLUMN IF EXISTS audit_period_end;
ALTER TABLE tax_audit_cases DROP COLUMN IF EXISTS audit_decision_number;
ALTER TABLE tax_audit_cases DROP COLUMN IF EXISTS auditor_contact;
ALTER TABLE tax_audit_cases DROP COLUMN IF EXISTS findings;
ALTER TABLE tax_audit_cases DROP COLUMN IF EXISTS penalty_amount;
