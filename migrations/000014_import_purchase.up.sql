ALTER TABLE supplier_invoices ADD COLUMN import_duty DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE supplier_invoices ADD COLUMN import_vat DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE supplier_invoices ADD COLUMN customs_declaration_number VARCHAR(50) NOT NULL DEFAULT '';
CREATE INDEX idx_inv_customs ON supplier_invoices (customs_declaration_number);
