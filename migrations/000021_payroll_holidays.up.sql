CREATE TABLE IF NOT EXISTS payroll_holidays (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    company_id TEXT NOT NULL,
    name TEXT NOT NULL,
    date DATE NOT NULL,
    year INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, name, date)
);

CREATE INDEX IF NOT EXISTS idx_payroll_holidays_company_year ON payroll_holidays(company_id, year);
