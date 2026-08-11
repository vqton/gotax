-- Down migration 000017: sale_schema
-- Drop all 13 tables in reverse dependency order.

DROP TABLE IF EXISTS ar_transactions;
DROP TABLE IF EXISTS cn_lines;
DROP TABLE IF EXISTS credit_notes;
DROP TABLE IF EXISTS receipt_allocations;
DROP TABLE IF EXISTS customer_receipts;
DROP TABLE IF EXISTS inv_lines;
DROP TABLE IF EXISTS customer_invoices;
DROP TABLE IF EXISTS dn_lines;
DROP TABLE IF EXISTS delivery_notes;
DROP TABLE IF EXISTS so_lines;
DROP TABLE IF EXISTS sales_orders;
DROP TABLE IF EXISTS sales_quotations;
DROP TABLE IF EXISTS customers;
