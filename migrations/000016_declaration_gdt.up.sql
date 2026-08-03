-- Round C: GDT declaration submission tracking
ALTER TABLE tax_declarations
    ADD COLUMN IF NOT EXISTS gdt_submission_id VARCHAR(100);
