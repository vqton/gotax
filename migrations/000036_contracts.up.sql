CREATE TABLE IF NOT EXISTS contracts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    contract_type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    value NUMERIC(19,4) DEFAULT 0,
    start_date DATE,
    end_date DATE,
    counterparty_name VARCHAR(255),
    counterparty_tax_code VARCHAR(20),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, code)
);

CREATE INDEX IF NOT EXISTS idx_contracts_company ON contracts(company_id);

CREATE TABLE IF NOT EXISTS contract_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_id UUID NOT NULL REFERENCES contracts(id) ON DELETE CASCADE,
    payment_date DATE NOT NULL,
    amount NUMERIC(19,4) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_contract_payments_contract ON contract_payments(contract_id);
