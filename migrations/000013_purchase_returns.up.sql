ALTER TABLE goods_receipt_notes ADD COLUMN return_of_grn_id VARCHAR(36) NOT NULL DEFAULT '';
CREATE INDEX idx_grn_return_of ON goods_receipt_notes (return_of_grn_id);

ALTER TABLE supplier_invoices ADD COLUMN original_invoice_id VARCHAR(36) NOT NULL DEFAULT '';
CREATE INDEX idx_inv_original ON supplier_invoices (original_invoice_id);
