ALTER TABLE goods_receipt_notes ADD COLUMN IF NOT EXISTS return_of_grn_id VARCHAR(36) NOT NULL DEFAULT '';
CREATE INDEX IF NOT EXISTS idx_grn_return_of ON goods_receipt_notes (return_of_grn_id);

ALTER TABLE supplier_invoices ADD COLUMN IF NOT EXISTS original_invoice_id VARCHAR(36) NOT NULL DEFAULT '';
CREATE INDEX IF NOT EXISTS idx_inv_original ON supplier_invoices (original_invoice_id);
