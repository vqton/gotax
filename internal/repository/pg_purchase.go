package repository

import (
	"context"
	"fmt"
	"gotax/internal/domain"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGPurchaseRepo struct {
	pool *pgxpool.Pool
}

func NewPGPurchaseRepo(pool *pgxpool.Pool) *PGPurchaseRepo {
	return &PGPurchaseRepo{pool: pool}
}

// ─── Supplier ────────────────────────────────────────────────────────────

func (r *PGPurchaseRepo) CreateSupplier(ctx context.Context, s *domain.Supplier) error {
	now := time.Now().UTC().Format(time.RFC3339)
	if s.Currency == "" {
		s.Currency = "VND"
	}
	if s.Status == "" {
		s.Status = domain.SupplierActive
	}
	return r.pool.QueryRow(ctx, `INSERT INTO suppliers (id, company_id, code, name, tax_code, address, phone, email, bank_account_name, bank_account_number, bank_name, payment_terms, credit_limit, currency, supplier_type, status, notes, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19) RETURNING id`,
		s.ID, s.CompanyID, s.Code, s.Name, s.TaxCode, s.Address, s.Phone, s.Email, s.BankAccountName, s.BankAccountNumber, s.BankName, s.PaymentTerms, s.CreditLimit, s.Currency, s.SupplierType, s.Status, s.Notes, now, now,
	).Scan(&s.ID)
}

func (r *PGPurchaseRepo) scanSupplierFromRow(row pgx.Row) (domain.Supplier, error) {
	var s domain.Supplier
	err := row.Scan(&s.ID, &s.CompanyID, &s.Code, &s.Name, &s.TaxCode, &s.Address, &s.Phone, &s.Email, &s.BankAccountName, &s.BankAccountNumber, &s.BankName, &s.PaymentTerms, &s.CreditLimit, &s.Currency, &s.SupplierType, &s.Status, &s.Notes, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func (r *PGPurchaseRepo) GetSupplier(ctx context.Context, id string) (*domain.Supplier, error) {
	s, err := r.scanSupplierFromRow(r.pool.QueryRow(ctx, `SELECT id, company_id, code, name, tax_code, address, phone, email, bank_account_name, bank_account_number, bank_name, payment_terms, credit_limit, currency, supplier_type, status, notes, created_at, updated_at FROM suppliers WHERE id=$1`, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrSupplierNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *PGPurchaseRepo) GetSupplierByCode(ctx context.Context, companyID, code string) (*domain.Supplier, error) {
	s, err := r.scanSupplierFromRow(r.pool.QueryRow(ctx, `SELECT id, company_id, code, name, tax_code, address, phone, email, bank_account_name, bank_account_number, bank_name, payment_terms, credit_limit, currency, supplier_type, status, notes, created_at, updated_at FROM suppliers WHERE company_id=$1 AND code=$2`, companyID, code))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrSupplierNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *PGPurchaseRepo) ListSuppliers(ctx context.Context, filter domain.PurchaseOrderFilter) ([]domain.Supplier, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM suppliers WHERE company_id=$1`, filter.CompanyID).Scan(&total); err != nil {
		return nil, 0, err
	}
	q := `SELECT id, company_id, code, name, tax_code, address, phone, email, bank_account_name, bank_account_number, bank_name, payment_terms, credit_limit, currency, supplier_type, status, notes, created_at, updated_at FROM suppliers WHERE company_id=$1`
	args := []interface{}{filter.CompanyID}
	n := 2
	if filter.SupplierID != "" {
		q += fmt.Sprintf(" AND supplier_id=$%d", n)
		args = append(args, filter.SupplierID)
		n++
	}
	q += " ORDER BY code"
	if filter.Limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", n)
		args = append(args, filter.Limit)
		n++
	}
	if filter.Offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", n)
		args = append(args, filter.Offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	out := make([]domain.Supplier, 0)
	for rows.Next() {
		s, err := r.scanSupplierFromRow(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, s)
	}
	return out, total, nil
}

func (r *PGPurchaseRepo) UpdateSupplier(ctx context.Context, s *domain.Supplier) error {
	_, err := r.pool.Exec(ctx, `UPDATE suppliers SET code=$1, name=$2, tax_code=$3, address=$4, phone=$5, email=$6, bank_account_name=$7, bank_account_number=$8, bank_name=$9, payment_terms=$10, credit_limit=$11, currency=$12, supplier_type=$13, status=$14, notes=$15, updated_at=NOW() WHERE id=$16`,
		s.Code, s.Name, s.TaxCode, s.Address, s.Phone, s.Email, s.BankAccountName, s.BankAccountNumber, s.BankName, s.PaymentTerms, s.CreditLimit, s.Currency, s.SupplierType, s.Status, s.Notes, s.ID)
	return err
}

func (r *PGPurchaseRepo) DeleteSupplier(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM suppliers WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrSupplierNotFound
	}
	return nil
}

// ─── Purchase Order ──────────────────────────────────────────────────────

func (r *PGPurchaseRepo) CreatePO(ctx context.Context, po *domain.PurchaseOrder) error {
	now := time.Now().UTC().Format(time.RFC3339)
	if po.Status == "" {
		po.Status = domain.POStatusDraft
	}
	po.CalculateTotals()
	return r.pool.QueryRow(ctx, `INSERT INTO purchase_orders (id, company_id, po_number, supplier_id, requisition_id, order_date, expected_date, currency, exchange_rate, payment_terms, delivery_terms, subtotal, discount_amount, tax_amount, total_amount, status, notes, created_by, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19) RETURNING id`,
		po.ID, po.CompanyID, po.PONumber, po.SupplierID, po.RequisitionID, po.OrderDate, po.ExpectedDate, po.Currency, po.ExchangeRate, po.PaymentTerms, po.DeliveryTerms, po.Subtotal, po.DiscountAmount, po.TaxAmount, po.TotalAmount, po.Status, po.Notes, po.CreatedBy, now, now,
	).Scan(&po.ID)
}

func (r *PGPurchaseRepo) scanPOFromRow(row pgx.Row) (domain.PurchaseOrder, error) {
	var po domain.PurchaseOrder
	err := row.Scan(&po.ID, &po.CompanyID, &po.PONumber, &po.SupplierID, &po.RequisitionID, &po.OrderDate, &po.ExpectedDate, &po.Currency, &po.ExchangeRate, &po.PaymentTerms, &po.DeliveryTerms, &po.Subtotal, &po.DiscountAmount, &po.TaxAmount, &po.TotalAmount, &po.Status, &po.ApprovedBy, &po.ApprovedAt, &po.CancelledReason, &po.Notes, &po.CreatedBy, &po.CreatedAt, &po.UpdatedAt)
	return po, err
}

func (r *PGPurchaseRepo) GetPO(ctx context.Context, id string) (*domain.PurchaseOrder, error) {
	po, err := r.scanPOFromRow(r.pool.QueryRow(ctx, `SELECT id, company_id, po_number, supplier_id, requisition_id, order_date, expected_date, currency, exchange_rate, payment_terms, delivery_terms, subtotal, discount_amount, tax_amount, total_amount, status, approved_by, approved_at, cancelled_reason, notes, created_by, created_at, updated_at FROM purchase_orders WHERE id=$1`, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrPONotFound
		}
		return nil, err
	}
	lines, _ := r.GetPOLines(ctx, id)
	po.Lines = lines
	return &po, nil
}

func (r *PGPurchaseRepo) GetPOByNumber(ctx context.Context, companyID, poNumber string) (*domain.PurchaseOrder, error) {
	po, err := r.scanPOFromRow(r.pool.QueryRow(ctx, `SELECT id, company_id, po_number, supplier_id, requisition_id, order_date, expected_date, currency, exchange_rate, payment_terms, delivery_terms, subtotal, discount_amount, tax_amount, total_amount, status, approved_by, approved_at, cancelled_reason, notes, created_by, created_at, updated_at FROM purchase_orders WHERE company_id=$1 AND po_number=$2`, companyID, poNumber))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrPONotFound
		}
		return nil, err
	}
	lines, _ := r.GetPOLines(ctx, po.ID)
	po.Lines = lines
	return &po, nil
}

func (r *PGPurchaseRepo) ListPOs(ctx context.Context, filter domain.PurchaseOrderFilter) ([]domain.PurchaseOrder, int, error) {
	var total int
	countQ := `SELECT COUNT(*) FROM purchase_orders WHERE company_id=$1`
	countArgs := []interface{}{filter.CompanyID}
	n := 2
	if filter.SupplierID != "" {
		countQ += fmt.Sprintf(" AND supplier_id=$%d", n)
		countArgs = append(countArgs, filter.SupplierID)
		n++
	}
	if filter.Status != "" {
		countQ += fmt.Sprintf(" AND status=$%d", n)
		countArgs = append(countArgs, filter.Status)
		n++
	}
	if filter.FromDate != "" {
		countQ += fmt.Sprintf(" AND order_date>=$%d", n)
		countArgs = append(countArgs, filter.FromDate)
		n++
	}
	if filter.ToDate != "" {
		countQ += fmt.Sprintf(" AND order_date<=$%d", n)
		countArgs = append(countArgs, filter.ToDate)
		n++
	}
	if err := r.pool.QueryRow(ctx, countQ, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}
	q := `SELECT id, company_id, po_number, supplier_id, requisition_id, order_date, expected_date, currency, exchange_rate, payment_terms, delivery_terms, subtotal, discount_amount, tax_amount, total_amount, status, approved_by, approved_at, cancelled_reason, notes, created_by, created_at, updated_at FROM purchase_orders WHERE company_id=$1`
	args := []interface{}{filter.CompanyID}
	n = 2
	if filter.SupplierID != "" {
		q += fmt.Sprintf(" AND supplier_id=$%d", n)
		args = append(args, filter.SupplierID)
		n++
	}
	if filter.Status != "" {
		q += fmt.Sprintf(" AND status=$%d", n)
		args = append(args, filter.Status)
		n++
	}
	if filter.FromDate != "" {
		q += fmt.Sprintf(" AND order_date>=$%d", n)
		args = append(args, filter.FromDate)
		n++
	}
	if filter.ToDate != "" {
		q += fmt.Sprintf(" AND order_date<=$%d", n)
		args = append(args, filter.ToDate)
		n++
	}
	q += " ORDER BY order_date DESC"
	if filter.Limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", n)
		args = append(args, filter.Limit)
		n++
	}
	if filter.Offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", n)
		args = append(args, filter.Offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	out := make([]domain.PurchaseOrder, 0)
	for rows.Next() {
		po, err := r.scanPOFromRow(rows)
		if err != nil {
			return nil, 0, err
		}
		lines, _ := r.GetPOLines(ctx, po.ID)
		po.Lines = lines
		out = append(out, po)
	}
	return out, total, nil
}

func (r *PGPurchaseRepo) UpdatePO(ctx context.Context, po *domain.PurchaseOrder) error {
	if _, err := r.pool.Exec(ctx, `UPDATE purchase_orders SET supplier_id=$1, requisition_id=$2, order_date=$3, expected_date=$4, currency=$5, exchange_rate=$6, payment_terms=$7, delivery_terms=$8, discount_amount=$9, notes=$10, updated_at=NOW() WHERE id=$11`,
		po.SupplierID, po.RequisitionID, po.OrderDate, po.ExpectedDate, po.Currency, po.ExchangeRate, po.PaymentTerms, po.DeliveryTerms, po.DiscountAmount, po.Notes, po.ID); err != nil {
		return err
	}
	return r.recalcPOLines(ctx, po.ID)
}

func (r *PGPurchaseRepo) UpdatePOStatus(ctx context.Context, id string, status domain.POStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE purchase_orders SET status=$1, updated_at=NOW() WHERE id=$2`, status, id)
	return err
}

func (r *PGPurchaseRepo) recalcPOLines(ctx context.Context, poID string) error {
	rows, err := r.pool.Query(ctx, `SELECT quantity, unit_price, discount_pct, vat_rate FROM po_lines WHERE po_id=$1`, poID)
	if err != nil {
		return err
	}
	defer rows.Close()
	var subtotal, discount, tax float64
	for rows.Next() {
		var qty, price, disc, vat float64
		if err := rows.Scan(&qty, &price, &disc, &vat); err != nil {
			return err
		}
		lt := qty * price * (1 - disc/100)
		subtotal += lt
		discount += qty * price * disc / 100
		tax += lt * vat / 100
	}
	_, err = r.pool.Exec(ctx, `UPDATE purchase_orders SET subtotal=$1, discount_amount=$2, tax_amount=$3, total_amount=$4, updated_at=NOW() WHERE id=$5`, subtotal, discount, tax, subtotal+tax, poID)
	return err
}

func (r *PGPurchaseRepo) GetPOLines(ctx context.Context, poID string) ([]domain.POItem, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, po_id, line_number, item_code, item_name, unit, quantity, unit_price, discount_pct, vat_rate, vat_type, account_id, vat_account_id, line_total, line_vat_amount, received_qty, invoiced_qty FROM po_lines WHERE po_id=$1 ORDER BY line_number`, poID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]domain.POItem, 0)
	for rows.Next() {
		var l domain.POItem
		if err := rows.Scan(&l.ID, &l.POID, &l.LineNumber, &l.ItemCode, &l.ItemName, &l.Unit, &l.Quantity, &l.UnitPrice, &l.DiscountPct, &l.VATRate, &l.VATType, &l.AccountID, &l.VATAccountID, &l.LineTotal, &l.LineVATAmount, &l.ReceivedQty, &l.InvoicedQty); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, nil
}

func (r *PGPurchaseRepo) CreatePOLines(ctx context.Context, items []domain.POItem) error {
	if len(items) == 0 {
		return nil
	}
	src := pgx.CopyFromSlice(len(items), func(i int) ([]interface{}, error) {
		it := items[i]
		lt := it.Quantity * it.UnitPrice * (1 - it.DiscountPct/100)
		lv := lt * it.VATRate / 100
		return []interface{}{it.ID, it.POID, it.LineNumber, it.ItemCode, it.ItemName, it.Unit, it.Quantity, it.UnitPrice, it.DiscountPct, it.VATRate, it.VATType, it.AccountID, it.VATAccountID, lt, lv, it.ReceivedQty, it.InvoicedQty}, nil
	})
	_, err := r.pool.CopyFrom(ctx, pgx.Identifier{"po_lines"}, []string{"id", "po_id", "line_number", "item_code", "item_name", "unit", "quantity", "unit_price", "discount_pct", "vat_rate", "vat_type", "account_id", "vat_account_id", "line_total", "line_vat_amount", "received_qty", "invoiced_qty"}, src)
	return err
}

func (r *PGPurchaseRepo) CreatePOLinesTx(ctx context.Context, tx pgx.Tx, items []domain.POItem) error {
	if len(items) == 0 {
		return nil
	}
	src := pgx.CopyFromSlice(len(items), func(i int) ([]interface{}, error) {
		it := items[i]
		lt := it.Quantity * it.UnitPrice * (1 - it.DiscountPct/100)
		lv := lt * it.VATRate / 100
		return []interface{}{it.ID, it.POID, it.LineNumber, it.ItemCode, it.ItemName, it.Unit, it.Quantity, it.UnitPrice, it.DiscountPct, it.VATRate, it.VATType, it.AccountID, it.VATAccountID, lt, lv, it.ReceivedQty, it.InvoicedQty}, nil
	})
	_, err := tx.CopyFrom(ctx, pgx.Identifier{"po_lines"}, []string{"id", "po_id", "line_number", "item_code", "item_name", "unit", "quantity", "unit_price", "discount_pct", "vat_rate", "vat_type", "account_id", "vat_account_id", "line_total", "line_vat_amount", "received_qty", "invoiced_qty"}, src)
	return err
}

func (r *PGPurchaseRepo) UpdatePOLines(ctx context.Context, items []domain.POItem) error {
	if len(items) == 0 {
		return nil
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `DELETE FROM po_lines WHERE po_id=$1`, items[0].POID); err != nil {
		return err
	}
	if err := r.CreatePOLinesTx(ctx, tx, items); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *PGPurchaseRepo) NextPONumber(ctx context.Context, companyID, yyyymm string) (string, error) {
	prefix := fmt.Sprintf("PO-%s-", yyyymm)
	var maxNum *string
	if err := r.pool.QueryRow(ctx, `SELECT MAX(po_number) FROM purchase_orders WHERE company_id=$1 AND po_number LIKE $2`, companyID, prefix+"%").Scan(&maxNum); err != nil {
		return "", err
	}
	if maxNum == nil || *maxNum == "" {
		return prefix + "0001", nil
	}
	if !strings.HasPrefix(*maxNum, prefix) {
		return prefix + "0001", nil
	}
	suffix, err := strconv.Atoi(strings.TrimPrefix(*maxNum, prefix))
	if err != nil {
		return prefix + "0001", nil
	}
	return fmt.Sprintf("%s%04d", prefix, suffix+1), nil
}

// ─── GRN ────────────────────────────────────────────────────────────────

func (r *PGPurchaseRepo) CreateGRN(ctx context.Context, g *domain.GRN) error {
	now := time.Now().UTC().Format(time.RFC3339)
	if g.Status == "" {
		g.Status = domain.GRNDraft
	}
	return r.pool.QueryRow(ctx, `INSERT INTO goods_receipt_notes (id, company_id, grn_number, po_id, receipt_date, warehouse, status, notes, created_by, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`,
		g.ID, g.CompanyID, g.GRNNumber, g.POID, g.ReceiptDate, g.Warehouse, g.Status, g.Notes, g.CreatedBy, now,
	).Scan(&g.ID)
}

func (r *PGPurchaseRepo) scanGRNFromRow(row pgx.Row) (domain.GRN, error) {
	var g domain.GRN
	err := row.Scan(&g.ID, &g.CompanyID, &g.GRNNumber, &g.POID, &g.ReceiptDate, &g.Warehouse, &g.Status, &g.Notes, &g.CreatedBy, &g.CreatedAt)
	return g, err
}

func (r *PGPurchaseRepo) GetGRN(ctx context.Context, id string) (*domain.GRN, error) {
	g, err := r.scanGRNFromRow(r.pool.QueryRow(ctx, `SELECT id, company_id, grn_number, po_id, receipt_date, warehouse, status, notes, created_by, created_at FROM goods_receipt_notes WHERE id=$1`, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrGRNNotFound
		}
		return nil, err
	}
	lines, _ := r.GetGRNLines(ctx, id)
	g.Lines = lines
	return &g, nil
}

func (r *PGPurchaseRepo) GetGRNByNumber(ctx context.Context, companyID, grnNumber string) (*domain.GRN, error) {
	g, err := r.scanGRNFromRow(r.pool.QueryRow(ctx, `SELECT id, company_id, grn_number, po_id, receipt_date, warehouse, status, notes, created_by, created_at FROM goods_receipt_notes WHERE company_id=$1 AND grn_number=$2`, companyID, grnNumber))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrGRNNotFound
		}
		return nil, err
	}
	lines, _ := r.GetGRNLines(ctx, g.ID)
	g.Lines = lines
	return &g, nil
}

func (r *PGPurchaseRepo) ListGRNs(ctx context.Context, filter domain.GRNFilter) ([]domain.GRN, int, error) {
	var total int
	countQ := `SELECT COUNT(*) FROM goods_receipt_notes WHERE company_id=$1`
	countArgs := []interface{}{filter.CompanyID}
	n := 2
	if filter.POID != "" {
		countQ += fmt.Sprintf(" AND po_id=$%d", n)
		countArgs = append(countArgs, filter.POID)
		n++
	}
	if filter.Status != "" {
		countQ += fmt.Sprintf(" AND status=$%d", n)
		countArgs = append(countArgs, filter.Status)
		n++
	}
	if filter.FromDate != "" {
		countQ += fmt.Sprintf(" AND receipt_date>=$%d", n)
		countArgs = append(countArgs, filter.FromDate)
		n++
	}
	if filter.ToDate != "" {
		countQ += fmt.Sprintf(" AND receipt_date<=$%d", n)
		countArgs = append(countArgs, filter.ToDate)
		n++
	}
	if err := r.pool.QueryRow(ctx, countQ, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}
	q := `SELECT id, company_id, grn_number, po_id, receipt_date, warehouse, status, notes, created_by, created_at FROM goods_receipt_notes WHERE company_id=$1`
	args := []interface{}{filter.CompanyID}
	n = 2
	if filter.POID != "" {
		q += fmt.Sprintf(" AND po_id=$%d", n)
		args = append(args, filter.POID)
		n++
	}
	if filter.Status != "" {
		q += fmt.Sprintf(" AND status=$%d", n)
		args = append(args, filter.Status)
		n++
	}
	if filter.FromDate != "" {
		q += fmt.Sprintf(" AND receipt_date>=$%d", n)
		args = append(args, filter.FromDate)
		n++
	}
	if filter.ToDate != "" {
		q += fmt.Sprintf(" AND receipt_date<=$%d", n)
		args = append(args, filter.ToDate)
		n++
	}
	q += " ORDER BY receipt_date DESC"
	if filter.Limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", n)
		args = append(args, filter.Limit)
		n++
	}
	if filter.Offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", n)
		args = append(args, filter.Offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	out := make([]domain.GRN, 0)
	for rows.Next() {
		g, err := r.scanGRNFromRow(rows)
		if err != nil {
			return nil, 0, err
		}
		lines, _ := r.GetGRNLines(ctx, g.ID)
		g.Lines = lines
		out = append(out, g)
	}
	return out, total, nil
}

func (r *PGPurchaseRepo) UpdateGRN(ctx context.Context, g *domain.GRN) error {
	_, err := r.pool.Exec(ctx, `UPDATE goods_receipt_notes SET receipt_date=$1, warehouse=$2, notes=$3 WHERE id=$4`, g.ReceiptDate, g.Warehouse, g.Notes, g.ID)
	return err
}

func (r *PGPurchaseRepo) UpdateGRNStatus(ctx context.Context, id string, status domain.GRNStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE goods_receipt_notes SET status=$1 WHERE id=$2`, status, id)
	return err
}

func (r *PGPurchaseRepo) GetGRNLines(ctx context.Context, grnID string) ([]domain.GRNItem, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, grn_id, po_line_id, item_code, item_name, unit, quantity_received, quantity_rejected, unit_price, line_total FROM grn_lines WHERE grn_id=$1 ORDER BY id`, grnID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]domain.GRNItem, 0)
	for rows.Next() {
		var l domain.GRNItem
		if err := rows.Scan(&l.ID, &l.GRNID, &l.POLineID, &l.ItemCode, &l.ItemName, &l.Unit, &l.QuantityReceived, &l.QuantityRejected, &l.UnitPrice, &l.LineTotal); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, nil
}

func (r *PGPurchaseRepo) CreateGRNLines(ctx context.Context, items []domain.GRNItem) error {
	if len(items) == 0 {
		return nil
	}
	src := pgx.CopyFromSlice(len(items), func(i int) ([]interface{}, error) {
		it := items[i]
		lt := it.QuantityReceived * it.UnitPrice
		return []interface{}{it.ID, it.GRNID, it.POLineID, it.ItemCode, it.ItemName, it.Unit, it.QuantityReceived, it.QuantityRejected, it.UnitPrice, lt}, nil
	})
	_, err := r.pool.CopyFrom(ctx, pgx.Identifier{"grn_lines"}, []string{"id", "grn_id", "po_line_id", "item_code", "item_name", "unit", "quantity_received", "quantity_rejected", "unit_price", "line_total"}, src)
	return err
}

func (r *PGPurchaseRepo) CreateGRNLinesTx(ctx context.Context, tx pgx.Tx, items []domain.GRNItem) error {
	if len(items) == 0 {
		return nil
	}
	src := pgx.CopyFromSlice(len(items), func(i int) ([]interface{}, error) {
		it := items[i]
		lt := it.QuantityReceived * it.UnitPrice
		return []interface{}{it.ID, it.GRNID, it.POLineID, it.ItemCode, it.ItemName, it.Unit, it.QuantityReceived, it.QuantityRejected, it.UnitPrice, lt}, nil
	})
	_, err := tx.CopyFrom(ctx, pgx.Identifier{"grn_lines"}, []string{"id", "grn_id", "po_line_id", "item_code", "item_name", "unit", "quantity_received", "quantity_rejected", "unit_price", "line_total"}, src)
	return err
}

func (r *PGPurchaseRepo) UpdateGRNLines(ctx context.Context, items []domain.GRNItem) error {
	if len(items) == 0 {
		return nil
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `DELETE FROM grn_lines WHERE grn_id=$1`, items[0].GRNID); err != nil {
		return err
	}
	if err := r.CreateGRNLinesTx(ctx, tx, items); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *PGPurchaseRepo) NextGRNNumber(ctx context.Context, companyID, yyyymm string) (string, error) {
	prefix := fmt.Sprintf("GRN-%s-", yyyymm)
	var maxNum *string
	if err := r.pool.QueryRow(ctx, `SELECT MAX(grn_number) FROM goods_receipt_notes WHERE company_id=$1 AND grn_number LIKE $2`, companyID, prefix+"%").Scan(&maxNum); err != nil {
		return "", err
	}
	if maxNum == nil || *maxNum == "" {
		return prefix + "0001", nil
	}
	if !strings.HasPrefix(*maxNum, prefix) {
		return prefix + "0001", nil
	}
	suffix, err := strconv.Atoi(strings.TrimPrefix(*maxNum, prefix))
	if err != nil {
		return prefix + "0001", nil
	}
	return fmt.Sprintf("%s%04d", prefix, suffix+1), nil
}

// ─── Supplier Invoice ──────────────────────────────────────────────────

func (r *PGPurchaseRepo) CreateInvoice(ctx context.Context, inv *domain.SupplierInvoice) error {
	now := time.Now().UTC().Format(time.RFC3339)
	if inv.Status == "" {
		inv.Status = domain.InvoiceDraft
	}
	inv.BalanceDue = inv.TotalAmount - inv.AmountPaid
	return r.pool.QueryRow(ctx, `INSERT INTO supplier_invoices (id, company_id, invoice_number, invoice_date, po_id, grn_id, supplier_id, supplier_name, supplier_tax_code, invoice_type, currency, exchange_rate, subtotal, discount_amount, tax_amount, total_amount, amount_paid, balance_due, due_date, vat_deduction_status, e_invoice_data, e_invoice_code, status, gl_posted, gl_posted_at, notes, created_by, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28) RETURNING id`,
		inv.ID, inv.CompanyID, inv.InvoiceNumber, inv.InvoiceDate, inv.POID, inv.GRNID, inv.SupplierID, inv.SupplierName, inv.SupplierTaxCode, inv.InvoiceType, inv.Currency, inv.ExchangeRate, inv.Subtotal, inv.DiscountAmount, inv.TaxAmount, inv.TotalAmount, inv.AmountPaid, inv.BalanceDue, inv.DueDate, inv.VATDeductionStatus, inv.EInvoiceData, inv.EInvoiceCode, inv.Status, inv.GLPosted, inv.GLPostedAt, inv.Notes, inv.CreatedBy, now,
	).Scan(&inv.ID)
}

func (r *PGPurchaseRepo) scanInvoiceFromRow(row pgx.Row) (domain.SupplierInvoice, error) {
	var inv domain.SupplierInvoice
	err := row.Scan(&inv.ID, &inv.CompanyID, &inv.InvoiceNumber, &inv.InvoiceDate, &inv.POID, &inv.GRNID, &inv.SupplierID, &inv.SupplierName, &inv.SupplierTaxCode, &inv.InvoiceType, &inv.Currency, &inv.ExchangeRate, &inv.Subtotal, &inv.DiscountAmount, &inv.TaxAmount, &inv.TotalAmount, &inv.AmountPaid, &inv.BalanceDue, &inv.DueDate, &inv.VATDeductionStatus, &inv.EInvoiceData, &inv.EInvoiceCode, &inv.Status, &inv.GLPosted, &inv.GLPostedAt, &inv.Notes, &inv.CreatedBy, &inv.CreatedAt)
	return inv, err
}

func (r *PGPurchaseRepo) GetInvoice(ctx context.Context, id string) (*domain.SupplierInvoice, error) {
	inv, err := r.scanInvoiceFromRow(r.pool.QueryRow(ctx, `SELECT id, company_id, invoice_number, invoice_date, po_id, grn_id, supplier_id, supplier_name, supplier_tax_code, invoice_type, currency, exchange_rate, subtotal, discount_amount, tax_amount, total_amount, amount_paid, balance_due, due_date, vat_deduction_status, e_invoice_data, e_invoice_code, status, gl_posted, gl_posted_at, notes, created_by, created_at FROM supplier_invoices WHERE id=$1`, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrSupplierInvoiceNotFound
		}
		return nil, err
	}
	lines, _ := r.GetInvoiceLines(ctx, id)
	inv.Lines = lines
	return &inv, nil
}

func (r *PGPurchaseRepo) GetInvoiceByNumber(ctx context.Context, companyID, invoiceNumber string) (*domain.SupplierInvoice, error) {
	inv, err := r.scanInvoiceFromRow(r.pool.QueryRow(ctx, `SELECT id, company_id, invoice_number, invoice_date, po_id, grn_id, supplier_id, supplier_name, supplier_tax_code, invoice_type, currency, exchange_rate, subtotal, discount_amount, tax_amount, total_amount, amount_paid, balance_due, due_date, vat_deduction_status, e_invoice_data, e_invoice_code, status, gl_posted, gl_posted_at, notes, created_by, created_at FROM supplier_invoices WHERE company_id=$1 AND invoice_number=$2`, companyID, invoiceNumber))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrSupplierInvoiceNotFound
		}
		return nil, err
	}
	lines, _ := r.GetInvoiceLines(ctx, inv.ID)
	inv.Lines = lines
	return &inv, nil
}

func (r *PGPurchaseRepo) ListInvoices(ctx context.Context, filter domain.SupplierInvoiceFilter) ([]domain.SupplierInvoice, int, error) {
	var total int
	countQ := `SELECT COUNT(*) FROM supplier_invoices WHERE company_id=$1`
	countArgs := []interface{}{filter.CompanyID}
	n := 2
	if filter.SupplierID != "" {
		countQ += fmt.Sprintf(" AND supplier_id=$%d", n)
		countArgs = append(countArgs, filter.SupplierID)
		n++
	}
	if filter.Status != "" {
		countQ += fmt.Sprintf(" AND status=$%d", n)
		countArgs = append(countArgs, filter.Status)
		n++
	}
	if filter.FromDate != "" {
		countQ += fmt.Sprintf(" AND invoice_date>=$%d", n)
		countArgs = append(countArgs, filter.FromDate)
		n++
	}
	if filter.ToDate != "" {
		countQ += fmt.Sprintf(" AND invoice_date<=$%d", n)
		countArgs = append(countArgs, filter.ToDate)
		n++
	}
	if err := r.pool.QueryRow(ctx, countQ, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}
	q := `SELECT id, company_id, invoice_number, invoice_date, po_id, grn_id, supplier_id, supplier_name, supplier_tax_code, invoice_type, currency, exchange_rate, subtotal, discount_amount, tax_amount, total_amount, amount_paid, balance_due, due_date, vat_deduction_status, e_invoice_data, e_invoice_code, status, gl_posted, gl_posted_at, notes, created_by, created_at FROM supplier_invoices WHERE company_id=$1`
	args := []interface{}{filter.CompanyID}
	n = 2
	if filter.SupplierID != "" {
		q += fmt.Sprintf(" AND supplier_id=$%d", n)
		args = append(args, filter.SupplierID)
		n++
	}
	if filter.Status != "" {
		q += fmt.Sprintf(" AND status=$%d", n)
		args = append(args, filter.Status)
		n++
	}
	if filter.FromDate != "" {
		q += fmt.Sprintf(" AND invoice_date>=$%d", n)
		args = append(args, filter.FromDate)
		n++
	}
	if filter.ToDate != "" {
		q += fmt.Sprintf(" AND invoice_date<=$%d", n)
		args = append(args, filter.ToDate)
		n++
	}
	q += " ORDER BY invoice_date DESC"
	if filter.Limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", n)
		args = append(args, filter.Limit)
		n++
	}
	if filter.Offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", n)
		args = append(args, filter.Offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	out := make([]domain.SupplierInvoice, 0)
	for rows.Next() {
		inv, err := r.scanInvoiceFromRow(rows)
		if err != nil {
			return nil, 0, err
		}
		lines, _ := r.GetInvoiceLines(ctx, inv.ID)
		inv.Lines = lines
		out = append(out, inv)
	}
	return out, total, nil
}

func (r *PGPurchaseRepo) UpdateInvoice(ctx context.Context, inv *domain.SupplierInvoice) error {
	_, err := r.pool.Exec(ctx, `UPDATE supplier_invoices SET po_id=$1, grn_id=$2, invoice_date=$3, subtotal=$4, discount_amount=$5, tax_amount=$6, total_amount=$7, amount_paid=$8, balance_due=$9, due_date=$10, vat_deduction_status=$11, e_invoice_data=$12, e_invoice_code=$13, notes=$14 WHERE id=$15`,
		inv.POID, inv.GRNID, inv.InvoiceDate, inv.Subtotal, inv.DiscountAmount, inv.TaxAmount, inv.TotalAmount, inv.AmountPaid, inv.BalanceDue, inv.DueDate, inv.VATDeductionStatus, inv.EInvoiceData, inv.EInvoiceCode, inv.Notes, inv.ID)
	return err
}

func (r *PGPurchaseRepo) UpdateInvoiceStatus(ctx context.Context, id string, status domain.InvoiceStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE supplier_invoices SET status=$1 WHERE id=$2`, status, id)
	return err
}

func (r *PGPurchaseRepo) PostInvoice(ctx context.Context, id string, postedAt time.Time) error {
	_, err := r.pool.Exec(ctx, `UPDATE supplier_invoices SET status=$1, gl_posted=true, gl_posted_at=$2 WHERE id=$3`,
		domain.InvoicePosted, postedAt, id)
	return err
}

func (r *PGPurchaseRepo) GetInvoiceLines(ctx context.Context, invoiceID string) ([]domain.SupplierInvoiceLine, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, invoice_id, po_line_id, grn_line_id, item_code, item_name, unit, quantity, unit_price, vat_rate, vat_type, line_total, line_vat_amount, account_id, vat_account_id FROM invoice_lines WHERE invoice_id=$1 ORDER BY id`, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]domain.SupplierInvoiceLine, 0)
	for rows.Next() {
		var l domain.SupplierInvoiceLine
		if err := rows.Scan(&l.ID, &l.InvoiceID, &l.POLineID, &l.GRNLineID, &l.ItemCode, &l.ItemName, &l.Unit, &l.Quantity, &l.UnitPrice, &l.VATRate, &l.VATType, &l.LineTotal, &l.LineVATAmount, &l.AccountID, &l.VATAccountID); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, nil
}

func (r *PGPurchaseRepo) CreateInvoiceLines(ctx context.Context, items []domain.SupplierInvoiceLine) error {
	if len(items) == 0 {
		return nil
	}
	src := pgx.CopyFromSlice(len(items), func(i int) ([]interface{}, error) {
		it := items[i]
		lt := it.Quantity * it.UnitPrice
		lv := lt * it.VATRate / 100
		return []interface{}{it.ID, it.InvoiceID, it.POLineID, it.GRNLineID, it.ItemCode, it.ItemName, it.Unit, it.Quantity, it.UnitPrice, it.VATRate, it.VATType, lt, lv, it.AccountID, it.VATAccountID}, nil
	})
	_, err := r.pool.CopyFrom(ctx, pgx.Identifier{"invoice_lines"}, []string{"id", "invoice_id", "po_line_id", "grn_line_id", "item_code", "item_name", "unit", "quantity", "unit_price", "vat_rate", "vat_type", "line_total", "line_vat_amount", "account_id", "vat_account_id"}, src)
	return err
}

func (r *PGPurchaseRepo) CreateInvoiceLinesTx(ctx context.Context, tx pgx.Tx, items []domain.SupplierInvoiceLine) error {
	if len(items) == 0 {
		return nil
	}
	src := pgx.CopyFromSlice(len(items), func(i int) ([]interface{}, error) {
		it := items[i]
		lt := it.Quantity * it.UnitPrice
		lv := lt * it.VATRate / 100
		return []interface{}{it.ID, it.InvoiceID, it.POLineID, it.GRNLineID, it.ItemCode, it.ItemName, it.Unit, it.Quantity, it.UnitPrice, it.VATRate, it.VATType, lt, lv, it.AccountID, it.VATAccountID}, nil
	})
	_, err := tx.CopyFrom(ctx, pgx.Identifier{"invoice_lines"}, []string{"id", "invoice_id", "po_line_id", "grn_line_id", "item_code", "item_name", "unit", "quantity", "unit_price", "vat_rate", "vat_type", "line_total", "line_vat_amount", "account_id", "vat_account_id"}, src)
	return err
}

func (r *PGPurchaseRepo) UpdateInvoiceLines(ctx context.Context, items []domain.SupplierInvoiceLine) error {
	if len(items) == 0 {
		return nil
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `DELETE FROM invoice_lines WHERE invoice_id=$1`, items[0].InvoiceID); err != nil {
		return err
	}
	if err := r.CreateInvoiceLinesTx(ctx, tx, items); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// ─── AP Transaction ────────────────────────────────────────────────────

func (r *PGPurchaseRepo) CreateAPTransaction(ctx context.Context, t *domain.APTransaction) error {
	now := time.Now().UTC().Format(time.RFC3339)
	if t.Currency == "" {
		t.Currency = "VND"
	}
	return r.pool.QueryRow(ctx, `INSERT INTO ap_transactions (id, company_id, supplier_id, invoice_id, transaction_type, transaction_date, amount, currency, reference_type, reference_id, notes, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id`,
		t.ID, t.CompanyID, t.SupplierID, t.InvoiceID, t.TransactionType, t.TransactionDate, t.Amount, t.Currency, t.ReferenceType, t.ReferenceID, t.Notes, now,
	).Scan(&t.ID)
}

func (r *PGPurchaseRepo) GetAPTransaction(ctx context.Context, id string) (*domain.APTransaction, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, company_id, supplier_id, invoice_id, transaction_type, transaction_date, amount, currency, reference_type, reference_id, notes, created_at FROM ap_transactions WHERE id=$1`, id)
	var t domain.APTransaction
	err := row.Scan(&t.ID, &t.CompanyID, &t.SupplierID, &t.InvoiceID, &t.TransactionType, &t.TransactionDate, &t.Amount, &t.Currency, &t.ReferenceType, &t.ReferenceID, &t.Notes, &t.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrAPTransNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *PGPurchaseRepo) ListAPTransactionsBySupplier(ctx context.Context, companyID, supplierID string) ([]domain.APTransaction, error) {
	q := `SELECT id, company_id, supplier_id, invoice_id, transaction_type, transaction_date, amount, currency, reference_type, reference_id, notes, created_at FROM ap_transactions WHERE supplier_id=$1`
	args := []interface{}{supplierID}
	if companyID != "" {
		q += " AND company_id=$2"
		args = append(args, companyID)
	}
	q += " ORDER BY transaction_date DESC"
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]domain.APTransaction, 0)
	for rows.Next() {
		var t domain.APTransaction
		if err := rows.Scan(&t.ID, &t.CompanyID, &t.SupplierID, &t.InvoiceID, &t.TransactionType, &t.TransactionDate, &t.Amount, &t.Currency, &t.ReferenceType, &t.ReferenceID, &t.Notes, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}

func (r *PGPurchaseRepo) ListAPTransactions(ctx context.Context, companyID string, offset, limit int) ([]domain.APTransaction, int, error) {
	var total int
	countQ := `SELECT COUNT(*) FROM ap_transactions`
	countArgs := []interface{}{}
	selQ := `SELECT id, company_id, supplier_id, invoice_id, transaction_type, transaction_date, amount, currency, reference_type, reference_id, notes, created_at FROM ap_transactions`
	sArgs := []interface{}{}
	if companyID != "" {
		countQ = `SELECT COUNT(*) FROM ap_transactions t JOIN suppliers s ON s.id=t.supplier_id WHERE s.company_id=$1`
		countArgs = append(countArgs, companyID)
		selQ = `SELECT t.id, t.company_id, t.supplier_id, t.invoice_id, t.transaction_type, t.transaction_date, t.amount, t.currency, t.reference_type, t.reference_id, t.notes, t.created_at FROM ap_transactions t JOIN suppliers s ON s.id=t.supplier_id WHERE s.company_id=$1`
		sArgs = append(sArgs, companyID)
	}
	if err := r.pool.QueryRow(ctx, countQ, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}
	selQ += " ORDER BY transaction_date DESC"
	if limit > 0 {
		selQ += fmt.Sprintf(" LIMIT $%d", len(sArgs)+1)
		sArgs = append(sArgs, limit)
	}
	if offset > 0 {
		selQ += fmt.Sprintf(" OFFSET $%d", len(sArgs)+1)
		sArgs = append(sArgs, offset)
	}
	rows, err := r.pool.Query(ctx, selQ, sArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	out := make([]domain.APTransaction, 0)
	for rows.Next() {
		var t domain.APTransaction
		if err := rows.Scan(&t.ID, &t.CompanyID, &t.SupplierID, &t.InvoiceID, &t.TransactionType, &t.TransactionDate, &t.Amount, &t.Currency, &t.ReferenceType, &t.ReferenceID, &t.Notes, &t.CreatedAt); err != nil {
			return nil, 0, err
		}
		out = append(out, t)
	}
	return out, total, nil
}

// ─── Cost Allocation ────────────────────────────────────────────────────

func (r *PGPurchaseRepo) CreateCostAllocation(ctx context.Context, c *domain.CostAllocation) error {
	return r.pool.QueryRow(ctx, `INSERT INTO cost_allocations (id, company_id, invoice_id, cost_type, cost_amount, allocation_method, allocated_lines, notes) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		c.ID, c.CompanyID, c.InvoiceID, c.CostType, c.CostAmount, c.AllocationMethod, c.AllocatedLines, c.Notes,
	).Scan(&c.ID)
}

func (r *PGPurchaseRepo) GetCostAllocation(ctx context.Context, id string) (*domain.CostAllocation, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, company_id, invoice_id, cost_type, cost_amount, allocation_method, allocated_lines, notes FROM cost_allocations WHERE id=$1`, id)
	var c domain.CostAllocation
	err := row.Scan(&c.ID, &c.CompanyID, &c.InvoiceID, &c.CostType, &c.CostAmount, &c.AllocationMethod, &c.AllocatedLines, &c.Notes)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("cost allocation %q not found", id)
		}
		return nil, err
	}
	return &c, nil
}

func (r *PGPurchaseRepo) ListCostAllocationsByInvoice(ctx context.Context, invoiceID string) ([]domain.CostAllocation, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, company_id, invoice_id, cost_type, cost_amount, allocation_method, allocated_lines, notes FROM cost_allocations WHERE invoice_id=$1`, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]domain.CostAllocation, 0)
	for rows.Next() {
		var c domain.CostAllocation
		if err := rows.Scan(&c.ID, &c.CompanyID, &c.InvoiceID, &c.CostType, &c.CostAmount, &c.AllocationMethod, &c.AllocatedLines, &c.Notes); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}
