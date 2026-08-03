CREATE TABLE IF NOT EXISTS purchase_requisitions (
    id VARCHAR(36) PRIMARY KEY,
    company_id VARCHAR(36) NOT NULL,
    requisition_number VARCHAR(50) NOT NULL,
    requester_id VARCHAR(36) NOT NULL,
    requester_name VARCHAR(255) NOT NULL DEFAULT '',
    department_id VARCHAR(36) NOT NULL DEFAULT '',
    need_by_date DATE,
    priority VARCHAR(20) NOT NULL DEFAULT '',
    reason TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    total_estimated DOUBLE PRECISION NOT NULL DEFAULT 0,
    approved_by VARCHAR(36) NOT NULL DEFAULT '',
    approved_at TIMESTAMPTZ,
    rejected_reason TEXT,
    created_by VARCHAR(36) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_req_company ON purchase_requisitions(company_id);
CREATE INDEX IF NOT EXISTS idx_req_status ON purchase_requisitions(status);
CREATE INDEX IF NOT EXISTS idx_req_number ON purchase_requisitions(requisition_number);

CREATE TABLE IF NOT EXISTS requisition_lines (
    id VARCHAR(36) PRIMARY KEY,
    requisition_id VARCHAR(36) NOT NULL REFERENCES purchase_requisitions(id) ON DELETE CASCADE,
    line_number INTEGER NOT NULL DEFAULT 0,
    item_code VARCHAR(50) NOT NULL DEFAULT '',
    item_name VARCHAR(255) NOT NULL,
    unit VARCHAR(20) NOT NULL DEFAULT '',
    quantity DOUBLE PRECISION NOT NULL DEFAULT 0,
    estimated_price DOUBLE PRECISION NOT NULL DEFAULT 0,
    estimated_total DOUBLE PRECISION NOT NULL DEFAULT 0,
    account_id VARCHAR(36) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_reql_req ON requisition_lines(requisition_id);
