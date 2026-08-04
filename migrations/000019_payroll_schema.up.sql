-- Payroll module schema (Phase 1)
-- Tables: employee_payroll_info, dependants, salary_components, payroll_periods,
--         payroll_runs, timekeeping_records, leave_requests, leave_balances,
--         payslips, payroll_config

-- Employee payroll info (extends employees table)
CREATE TABLE IF NOT EXISTS employee_payroll_info (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id),
    contract_type VARCHAR(20) NOT NULL DEFAULT 'INDEFINITE',
    contract_start_date DATE,
    contract_end_date DATE,
    salary_type VARCHAR(20) NOT NULL DEFAULT 'TIME_BASED',
    base_salary NUMERIC(15,2) NOT NULL DEFAULT 0,
    salary_coefficient NUMERIC(5,2) DEFAULT 0,
    position_allowance NUMERIC(15,2) DEFAULT 0,
    responsibility_allowance NUMERIC(15,2) DEFAULT 0,
    seniority_allowance NUMERIC(15,2) DEFAULT 0,
    other_allowances NUMERIC(15,2) DEFAULT 0,
    insurance_base_salary NUMERIC(15,2) DEFAULT 0,
    region VARCHAR(5) NOT NULL DEFAULT 'I',
    is_foreign_employee BOOLEAN DEFAULT FALSE,
    is_trade_union_member BOOLEAN DEFAULT FALSE,
    is_high_tech_talent BOOLEAN DEFAULT FALSE,
    bank_account_no VARCHAR(50),
    bank_code VARCHAR(20),
    effective_date DATE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(employee_id, effective_date)
);

CREATE INDEX IF NOT EXISTS idx_payroll_info_employee ON employee_payroll_info(employee_id);
CREATE INDEX IF NOT EXISTS idx_payroll_info_region ON employee_payroll_info(region);

-- Dependants
CREATE TABLE IF NOT EXISTS dependants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id),
    full_name VARCHAR(200) NOT NULL,
    relationship VARCHAR(20) NOT NULL,
    date_of_birth DATE,
    tax_code VARCHAR(20),
    is_disabled BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_dependants_employee ON dependants(employee_id);

-- Salary components
CREATE TABLE IF NOT EXISTS salary_components (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    code VARCHAR(20) NOT NULL,
    name VARCHAR(200) NOT NULL,
    type VARCHAR(20) NOT NULL,
    calculation VARCHAR(20) NOT NULL,
    formula TEXT,
    default_value NUMERIC(15,2) DEFAULT 0,
    is_taxable BOOLEAN DEFAULT TRUE,
    is_insurable BOOLEAN DEFAULT TRUE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(company_id, code)
);

CREATE INDEX IF NOT EXISTS idx_salary_components_company ON salary_components(company_id);

-- Payroll periods
CREATE TABLE IF NOT EXISTS payroll_periods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    year INT NOT NULL,
    month INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    prepared_by VARCHAR(100),
    prepared_at TIMESTAMPTZ,
    reviewed_by VARCHAR(100),
    reviewed_at TIMESTAMPTZ,
    approved_by VARCHAR(100),
    approved_at TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,
    closed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(company_id, year, month)
);

CREATE INDEX IF NOT EXISTS idx_payroll_periods_company ON payroll_periods(company_id);

-- Payroll runs (per employee per period)
CREATE TABLE IF NOT EXISTS payroll_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    period_id UUID NOT NULL REFERENCES payroll_periods(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    company_id UUID NOT NULL,

    -- Input
    working_days INT DEFAULT 0,
    ot_hours NUMERIC(5,2) DEFAULT 0,
    night_shift_hours NUMERIC(5,2) DEFAULT 0,
    leave_days NUMERIC(5,2) DEFAULT 0,
    unpaid_leave_days NUMERIC(5,2) DEFAULT 0,
    absent_days INT DEFAULT 0,

    -- Income
    base_salary NUMERIC(15,2) DEFAULT 0,
    ot_pay NUMERIC(15,2) DEFAULT 0,
    night_shift_pay NUMERIC(15,2) DEFAULT 0,
    leave_pay NUMERIC(15,2) DEFAULT 0,
    holiday_pay NUMERIC(15,2) DEFAULT 0,
    allowances NUMERIC(15,2) DEFAULT 0,
    bonuses NUMERIC(15,2) DEFAULT 0,
    other_income NUMERIC(15,2) DEFAULT 0,
    gross_salary NUMERIC(15,2) DEFAULT 0,

    -- Employee deductions
    si_deduction NUMERIC(15,2) DEFAULT 0,
    hi_deduction NUMERIC(15,2) DEFAULT 0,
    ui_deduction NUMERIC(15,2) DEFAULT 0,
    trade_union_dues NUMERIC(15,2) DEFAULT 0,
    pit_amount NUMERIC(15,2) DEFAULT 0,
    other_deductions NUMERIC(15,2) DEFAULT 0,
    total_deductions NUMERIC(15,2) DEFAULT 0,

    -- Net
    net_pay NUMERIC(15,2) DEFAULT 0,

    -- Employer costs
    employer_si NUMERIC(15,2) DEFAULT 0,
    employer_hi NUMERIC(15,2) DEFAULT 0,
    employer_ui NUMERIC(15,2) DEFAULT 0,
    employer_trade_union NUMERIC(15,2) DEFAULT 0,
    total_employer_cost NUMERIC(15,2) DEFAULT 0,

    -- Insurance base
    insurance_base NUMERIC(15,2) DEFAULT 0,
    ui_base NUMERIC(15,2) DEFAULT 0,

    -- Status
    status VARCHAR(20) DEFAULT 'DRAFT',
    adjustment_reason TEXT,
    adjusted_by VARCHAR(100),
    adjusted_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(period_id, employee_id)
);

CREATE INDEX IF NOT EXISTS idx_payroll_runs_period ON payroll_runs(period_id);
CREATE INDEX IF NOT EXISTS idx_payroll_runs_employee ON payroll_runs(employee_id);
CREATE INDEX IF NOT EXISTS idx_payroll_runs_company ON payroll_runs(company_id);

-- Timekeeping records
CREATE TABLE IF NOT EXISTS timekeeping_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id),
    company_id UUID NOT NULL,
    date DATE NOT NULL,
    clock_in TIME,
    clock_out TIME,
    hours_worked NUMERIC(5,2) DEFAULT 0,
    ot_hours NUMERIC(5,2) DEFAULT 0,
    night_hours NUMERIC(5,2) DEFAULT 0,
    is_holiday BOOLEAN DEFAULT FALSE,
    is_rest_day BOOLEAN DEFAULT FALSE,
    leave_type VARCHAR(20),
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(employee_id, date)
);

CREATE INDEX IF NOT EXISTS idx_timekeeping_employee ON timekeeping_records(employee_id);
CREATE INDEX IF NOT EXISTS idx_timekeeping_date ON timekeeping_records(date);

-- Leave requests
CREATE TABLE IF NOT EXISTS leave_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id),
    company_id UUID NOT NULL,
    leave_type VARCHAR(20) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    days NUMERIC(5,2) NOT NULL,
    reason TEXT,
    status VARCHAR(20) DEFAULT 'PENDING',
    approved_by VARCHAR(100),
    approved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_leave_requests_employee ON leave_requests(employee_id);
CREATE INDEX IF NOT EXISTS idx_leave_requests_company ON leave_requests(company_id);
CREATE INDEX IF NOT EXISTS idx_leave_requests_status ON leave_requests(status);

-- Leave balances
CREATE TABLE IF NOT EXISTS leave_balances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id),
    year INT NOT NULL,
    leave_type VARCHAR(20) NOT NULL,
    entitled NUMERIC(5,2) DEFAULT 0,
    used NUMERIC(5,2) DEFAULT 0,
    remaining NUMERIC(5,2) DEFAULT 0,
    carried_over NUMERIC(5,2) DEFAULT 0,
    UNIQUE(employee_id, year, leave_type)
);

CREATE INDEX IF NOT EXISTS idx_leave_balances_employee ON leave_balances(employee_id);

-- Payslips
CREATE TABLE IF NOT EXISTS payslips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    run_id UUID NOT NULL REFERENCES payroll_runs(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    period_id UUID NOT NULL REFERENCES payroll_periods(id),
    payslip_no VARCHAR(50) NOT NULL,
    pdf_path TEXT,
    sent_at TIMESTAMPTZ,
    viewed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payslips_run ON payslips(run_id);
CREATE INDEX IF NOT EXISTS idx_payslips_employee ON payslips(employee_id);
CREATE INDEX IF NOT EXISTS idx_payslips_period ON payslips(period_id);

-- Payroll config
CREATE TABLE IF NOT EXISTS payroll_config (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    config_key VARCHAR(100) NOT NULL,
    config_value TEXT NOT NULL,
    effective_from DATE NOT NULL,
    effective_to DATE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(company_id, config_key)
);

CREATE INDEX IF NOT EXISTS idx_payroll_config_company ON payroll_config(company_id);
