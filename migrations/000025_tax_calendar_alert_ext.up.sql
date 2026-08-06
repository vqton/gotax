-- Extend tax_calendar with full fields used by service layer
ALTER TABLE tax_calendar ADD COLUMN IF NOT EXISTS period_type VARCHAR(10) NOT NULL DEFAULT 'MONTHLY';
ALTER TABLE tax_calendar ADD COLUMN IF NOT EXISTS start_date DATE;
ALTER TABLE tax_calendar ADD COLUMN IF NOT EXISTS end_date DATE;
ALTER TABLE tax_calendar ADD COLUMN IF NOT EXISTS declaration_due DATE;
ALTER TABLE tax_calendar ADD COLUMN IF NOT EXISTS payment_due DATE;
ALTER TABLE tax_calendar ADD COLUMN IF NOT EXISTS declaration_id VARCHAR(36);

-- Extend tax_alerts with calendar link + acknowledgement
ALTER TABLE tax_alerts ADD COLUMN IF NOT EXISTS calendar_id VARCHAR(36);
ALTER TABLE tax_alerts ADD COLUMN IF NOT EXISTS channel VARCHAR(20) NOT NULL DEFAULT 'ALL';
ALTER TABLE tax_alerts ADD COLUMN IF NOT EXISTS acknowledged_at TIMESTAMPTZ;
ALTER TABLE tax_alerts ADD COLUMN IF NOT EXISTS acknowledged_by VARCHAR(36);

-- Extend tax_audit_cases with full audit fields
ALTER TABLE tax_audit_cases ADD COLUMN IF NOT EXISTS audit_period_start DATE;
ALTER TABLE tax_audit_cases ADD COLUMN IF NOT EXISTS audit_period_end DATE;
ALTER TABLE tax_audit_cases ADD COLUMN IF NOT EXISTS audit_decision_number VARCHAR(50);
ALTER TABLE tax_audit_cases ADD COLUMN IF NOT EXISTS auditor_contact VARCHAR(255);
ALTER TABLE tax_audit_cases ADD COLUMN IF NOT EXISTS findings TEXT;
ALTER TABLE tax_audit_cases ADD COLUMN IF NOT EXISTS penalty_amount DOUBLE PRECISION DEFAULT 0;
