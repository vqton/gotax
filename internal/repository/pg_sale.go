package repository

import (
	"context"
	"fmt"
	"gotax/internal/domain"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGSaleRepo struct {
	pool *pgxpool.Pool
}

func NewPGSaleRepo(pool *pgxpool.Pool) *PGSaleRepo {
	return &PGSaleRepo{pool: pool}
}

// ─── Customer ────────────────────────────────────────────────────────────

func (r *PGSaleRepo) CreateCustomer(ctx context.Context, c *domain.Customer) error {
	now := time.Now().UTC().Format(time.RFC3339)
	if c.Status == "" {
		c.Status = domain.CustomerActive
	}
	return r.pool.QueryRow(ctx, `INSERT INTO customers (id,company_id,code,name,tax_code,address,phone,email,bank_account_name,bank_account_number,bank_name,payment_terms,credit_limit,currency,customer_type,status,notes,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19) RETURNING id`,
		c.ID, c.CompanyID, c.Code, c.Name, c.TaxCode, c.Address, c.Phone, c.Email, c.BankAccountName, c.BankAccountNumber, c.BankName, c.PaymentTerms, c.CreditLimit, c.Currency, c.CustomerType, c.Status, c.Notes, now, now,
	).Scan(&c.ID)
}

func scanCustomer(row pgx.Row) (domain.Customer, error) {
	var c domain.Customer
	err := row.Scan(&c.ID, &c.CompanyID, &c.Code, &c.Name, &c.TaxCode, &c.Address, &c.Phone, &c.Email, &c.BankAccountName, &c.BankAccountNumber, &c.BankName, &c.PaymentTerms, &c.CreditLimit, &c.Currency, &c.CustomerType, &c.Status, &c.Notes, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *PGSaleRepo) GetCustomer(ctx context.Context, id string) (*domain.Customer, error) {
	c, err := scanCustomer(r.pool.QueryRow(ctx, `SELECT id,company_id,code,name,tax_code,address,phone,email,bank_account_name,bank_account_number,bank_name,payment_terms,credit_limit,currency,customer_type,status,notes,created_at,updated_at FROM customers WHERE id=$1`, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrCustomerNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *PGSaleRepo) GetCustomerByCode(ctx context.Context, companyID, code string) (*domain.Customer, error) {
	c, err := scanCustomer(r.pool.QueryRow(ctx, `SELECT id,company_id,code,name,tax_code,address,phone,email,bank_account_name,bank_account_number,bank_name,payment_terms,credit_limit,currency,customer_type,status,notes,created_at,updated_at FROM customers WHERE company_id=$1 AND code=$2`, companyID, code))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrCustomerNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *PGSaleRepo) ListCustomers(ctx context.Context, companyID string) ([]domain.Customer, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,company_id,code,name,tax_code,address,phone,email,bank_account_name,bank_account_number,bank_name,payment_terms,credit_limit,currency,customer_type,status,notes,created_at,updated_at FROM customers WHERE company_id=$1 ORDER BY code`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Customer
	for rows.Next() {
		c, err := scanCustomer(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}

func (r *PGSaleRepo) UpdateCustomer(ctx context.Context, c *domain.Customer) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := r.pool.Exec(ctx, `UPDATE customers SET code=$1,name=$2,tax_code=$3,address=$4,phone=$5,email=$6,bank_account_name=$7,bank_account_number=$8,bank_name=$9,payment_terms=$10,credit_limit=$11,currency=$12,customer_type=$13,status=$14,notes=$15,updated_at=$16 WHERE id=$17`,
		c.Code, c.Name, c.TaxCode, c.Address, c.Phone, c.Email, c.BankAccountName, c.BankAccountNumber, c.BankName, c.PaymentTerms, c.CreditLimit, c.Currency, c.CustomerType, c.Status, c.Notes, now, c.ID)
	return err
}

func (r *PGSaleRepo) DeleteCustomer(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE customers SET status=$1 WHERE id=$2`, domain.CustomerSuspended, id)
	return err
}

// ─── Sales Order ─────────────────────────────────────────────────────────

func (r *PGSaleRepo) CreateSO(ctx context.Context, so *domain.SalesOrder) error {
	now := time.Now().UTC().Format(time.RFC3339)
	if so.Status == "" {
		so.Status = domain.SODraft
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := tx.QueryRow(ctx, `INSERT INTO sales_orders (id,company_id,so_number,quotation_id,customer_id,order_date,expected_date,currency,exchange_rate,payment_terms,delivery_terms,shipping_address,subtotal,discount_amount,tax_amount,total_amount,status,approved_by,approved_at,cancelled_reason,notes,created_by,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24) RETURNING id`,
		so.ID, so.CompanyID, so.SONumber, so.QuotationID, so.CustomerID, so.OrderDate.Format(time.RFC3339), pt(so.ExpectedDate), so.Currency, so.ExchangeRate, so.PaymentTerms, so.DeliveryTerms, so.ShippingAddress, so.Subtotal, so.DiscountAmount, so.TaxAmount, so.TotalAmount, so.Status, so.ApprovedBy, pt(so.ApprovedAt), so.CancelledReason, so.Notes, so.CreatedBy, now, now,
	).Scan(&so.ID); err != nil {
		return err
	}
	for _, l := range so.Lines {
		_, err = tx.Exec(ctx, `INSERT INTO so_lines (id,so_id,line_number,item_code,item_name,unit,quantity,unit_price,discount_pct,vat_rate,vat_type,revenue_account_id,vat_account_id,line_total,line_vat_amount,delivered_qty,invoiced_qty) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`,
			l.ID, so.ID, l.LineNumber, l.ItemCode, l.ItemName, l.Unit, l.Quantity, l.UnitPrice, l.DiscountPct, l.VATRate, l.VATType, l.RevenueAccount, l.VATAccountID, l.LineTotal, l.LineVATAmount, l.DeliveredQty, l.InvoicedQty)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func scanSO(row pgx.Row) (domain.SalesOrder, error) {
	var so domain.SalesOrder
	var orderDate, expDate, appAt, crAt, upAt string
	err := row.Scan(&so.ID, &so.CompanyID, &so.SONumber, &so.QuotationID, &so.CustomerID, &orderDate, &expDate, &so.Currency, &so.ExchangeRate, &so.PaymentTerms, &so.DeliveryTerms, &so.ShippingAddress, &so.Subtotal, &so.DiscountAmount, &so.TaxAmount, &so.TotalAmount, &so.Status, &so.ApprovedBy, &appAt, &so.CancelledReason, &so.Notes, &so.CreatedBy, &crAt, &upAt)
	if err != nil {
		return so, err
	}
	so.OrderDate, _ = time.Parse(time.RFC3339, orderDate)
	if expDate != "" {
		t, _ := time.Parse(time.RFC3339, expDate)
		so.ExpectedDate = &t
	}
	if appAt != "" {
		t, _ := time.Parse(time.RFC3339, appAt)
		so.ApprovedAt = &t
	}
	so.CreatedAt, _ = time.Parse(time.RFC3339, crAt)
	so.UpdatedAt, _ = time.Parse(time.RFC3339, upAt)
	return so, nil
}

func (r *PGSaleRepo) GetSO(ctx context.Context, id string) (*domain.SalesOrder, error) {
	so, err := scanSO(r.pool.QueryRow(ctx, `SELECT id,company_id,so_number,quotation_id,customer_id,order_date,expected_date,currency,exchange_rate,payment_terms,delivery_terms,shipping_address,subtotal,discount_amount,tax_amount,total_amount,status,approved_by,approved_at,cancelled_reason,notes,created_by,created_at,updated_at FROM sales_orders WHERE id=$1`, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrSONotFound
		}
		return nil, err
	}
	so.Lines, _ = r.GetSOLines(ctx, id)
	return &so, nil
}

func (r *PGSaleRepo) GetSOByNumber(ctx context.Context, companyID, soNumber string) (*domain.SalesOrder, error) {
	so, err := scanSO(r.pool.QueryRow(ctx, `SELECT id,company_id,so_number,quotation_id,customer_id,order_date,expected_date,currency,exchange_rate,payment_terms,delivery_terms,shipping_address,subtotal,discount_amount,tax_amount,total_amount,status,approved_by,approved_at,cancelled_reason,notes,created_by,created_at,updated_at FROM sales_orders WHERE company_id=$1 AND so_number=$2`, companyID, soNumber))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrSONotFound
		}
		return nil, err
	}
	so.Lines, _ = r.GetSOLines(ctx, so.ID)
	return &so, nil
}

type soCountArgs struct {
	companyID  string
	customerID string
	status     string
	fromDate   string
	toDate     string
}

func soCountClauses(f domain.SalesOrderFilter) (string, []interface{}) {
	args := []interface{}{f.CompanyID}
	n := 2
	q := ""
	if f.CustomerID != "" {
		q += fmt.Sprintf(" AND customer_id=$%d", n); args = append(args, f.CustomerID); n++
	}
	if f.Status != "" {
		q += fmt.Sprintf(" AND status=$%d", n); args = append(args, string(f.Status)); n++
	}
	if f.FromDate != "" {
		q += fmt.Sprintf(" AND order_date>=$%d", n); args = append(args, f.FromDate); n++
	}
	if f.ToDate != "" {
		q += fmt.Sprintf(" AND order_date<=$%d", n); args = append(args, f.ToDate); n++
	}
	return q, args
}

func (r *PGSaleRepo) ListSOs(ctx context.Context, filter domain.SalesOrderFilter) ([]domain.SalesOrder, int, error) {
	where, args := soCountClauses(filter)
	var total int
	cq := `SELECT COUNT(*) FROM sales_orders WHERE company_id=$1` + where
	if err := r.pool.QueryRow(ctx, cq, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	q := `SELECT id,company_id,so_number,quotation_id,customer_id,order_date,expected_date,currency,exchange_rate,payment_terms,delivery_terms,shipping_address,subtotal,discount_amount,tax_amount,total_amount,status,approved_by,approved_at,cancelled_reason,notes,created_by,created_at,updated_at FROM sales_orders WHERE company_id=$1` + where
	dArgs := make([]interface{}, len(args))
	copy(dArgs, args)
	n := len(args) + 1
	q += " ORDER BY order_date DESC"
	if filter.Limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", n); dArgs = append(dArgs, filter.Limit); n++
	}
	if filter.Offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", n); dArgs = append(dArgs, filter.Offset)
	}
	rows, err := r.pool.Query(ctx, q, dArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []domain.SalesOrder
	for rows.Next() {
		so, err := scanSO(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, so)
	}
	return out, total, nil
}

func (r *PGSaleRepo) UpdateSO(ctx context.Context, so *domain.SalesOrder) error {
	now := time.Now().UTC().Format(time.RFC3339)
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `UPDATE sales_orders SET so_number=$1,customer_id=$2,order_date=$3,expected_date=$4,currency=$5,exchange_rate=$6,payment_terms=$7,delivery_terms=$8,shipping_address=$9,subtotal=$10,discount_amount=$11,tax_amount=$12,total_amount=$13,status=$14,notes=$15,updated_at=$16 WHERE id=$17`,
		so.SONumber, so.CustomerID, so.OrderDate.Format(time.RFC3339), so.ExpectedDate.Format(time.RFC3339), so.Currency, so.ExchangeRate, so.PaymentTerms, so.DeliveryTerms, so.ShippingAddress, so.Subtotal, so.DiscountAmount, so.TaxAmount, so.TotalAmount, so.Status, so.Notes, now, so.ID)
	if err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `DELETE FROM so_lines WHERE so_id=$1`, so.ID); err != nil {
		return err
	}
	for _, l := range so.Lines {
		_, err = tx.Exec(ctx, `INSERT INTO so_lines (id,so_id,line_number,item_code,item_name,unit,quantity,unit_price,discount_pct,vat_rate,vat_type,revenue_account_id,vat_account_id,line_total,line_vat_amount,delivered_qty,invoiced_qty) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`,
			l.ID, so.ID, l.LineNumber, l.ItemCode, l.ItemName, l.Unit, l.Quantity, l.UnitPrice, l.DiscountPct, l.VATRate, l.VATType, l.RevenueAccount, l.VATAccountID, l.LineTotal, l.LineVATAmount, l.DeliveredQty, l.InvoicedQty)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *PGSaleRepo) UpdateSOStatus(ctx context.Context, id string, status domain.SOStatus) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := r.pool.Exec(ctx, `UPDATE sales_orders SET status=$1,updated_at=$2 WHERE id=$3`, string(status), now, id)
	return err
}

func (r *PGSaleRepo) ApproveSO(ctx context.Context, id, approvedBy string, approvedAt time.Time) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := r.pool.Exec(ctx, `UPDATE sales_orders SET status=$1,approved_by=$2,approved_at=$3,updated_at=$4 WHERE id=$5`,
		string(domain.SOApproved), approvedBy, approvedAt.Format(time.RFC3339), now, id)
	return err
}

func (r *PGSaleRepo) CancelSO(ctx context.Context, id, cancelReason string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := r.pool.Exec(ctx, `UPDATE sales_orders SET status=$1,cancelled_reason=$2,updated_at=$3 WHERE id=$4`,
		string(domain.SOCancelled), cancelReason, now, id)
	return err
}

func (r *PGSaleRepo) GetSOLines(ctx context.Context, soID string) ([]domain.SOLine, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,so_id,line_number,item_code,item_name,unit,quantity,unit_price,discount_pct,vat_rate,vat_type,revenue_account_id,vat_account_id,line_total,line_vat_amount,delivered_qty,invoiced_qty FROM so_lines WHERE so_id=$1 ORDER BY line_number`, soID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.SOLine
	for rows.Next() {
		var l domain.SOLine
		if err := rows.Scan(&l.ID, &l.SOID, &l.LineNumber, &l.ItemCode, &l.ItemName, &l.Unit, &l.Quantity, &l.UnitPrice, &l.DiscountPct, &l.VATRate, &l.VATType, &l.RevenueAccount, &l.VATAccountID, &l.LineTotal, &l.LineVATAmount, &l.DeliveredQty, &l.InvoicedQty); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, nil
}

func (r *PGSaleRepo) CreateSOLines(ctx context.Context, items []domain.SOLine) error {
	for _, l := range items {
		_, err := r.pool.Exec(ctx, `INSERT INTO so_lines (id,so_id,line_number,item_code,item_name,unit,quantity,unit_price,discount_pct,vat_rate,vat_type,revenue_account_id,vat_account_id,line_total,line_vat_amount,delivered_qty,invoiced_qty) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`,
			l.ID, l.SOID, l.LineNumber, l.ItemCode, l.ItemName, l.Unit, l.Quantity, l.UnitPrice, l.DiscountPct, l.VATRate, l.VATType, l.RevenueAccount, l.VATAccountID, l.LineTotal, l.LineVATAmount, l.DeliveredQty, l.InvoicedQty)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PGSaleRepo) UpdateSOLines(ctx context.Context, items []domain.SOLine) error {
	for _, l := range items {
		_, err := r.pool.Exec(ctx, `UPDATE so_lines SET item_code=$1,item_name=$2,unit=$3,quantity=$4,unit_price=$5,discount_pct=$6,vat_rate=$7,vat_type=$8,revenue_account_id=$9,vat_account_id=$10,line_total=$11,line_vat_amount=$12,delivered_qty=$13,invoiced_qty=$14 WHERE id=$15`,
			l.ItemCode, l.ItemName, l.Unit, l.Quantity, l.UnitPrice, l.DiscountPct, l.VATRate, l.VATType, l.RevenueAccount, l.VATAccountID, l.LineTotal, l.LineVATAmount, l.DeliveredQty, l.InvoicedQty, l.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PGSaleRepo) NextSONumber(ctx context.Context, companyID, yyyymm string) (string, error) {
	var n int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*)+1 FROM sales_orders WHERE company_id=$1 AND so_number LIKE $2`, companyID, "SO-"+yyyymm+"-%").Scan(&n)
	if err != nil {
		return fmt.Sprintf("SO-%s-00001", yyyymm), nil
	}
	return fmt.Sprintf("SO-%s-%05d", yyyymm, n), nil
}

// ─── Delivery Note ────────────────────────────────────────────────────────

func (r *PGSaleRepo) CreateDN(ctx context.Context, dn *domain.DeliveryNote) error {
	now := time.Now().UTC().Format(time.RFC3339)
	if dn.Status == "" {
		dn.Status = domain.DNDraft
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := tx.QueryRow(ctx, `INSERT INTO delivery_notes (id,company_id,dn_number,so_id,delivery_date,warehouse,shipping_method,carrier_name,tracking_number,delivery_address,status,notes,created_by,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING id`,
		dn.ID, dn.CompanyID, dn.DNNumber, dn.SOID, dn.DeliveryDate.Format(time.RFC3339), dn.Warehouse, dn.ShippingMethod, dn.CarrierName, dn.TrackingNumber, dn.DeliveryAddress, dn.Status, dn.Notes, dn.CreatedBy, now,
	).Scan(&dn.ID); err != nil {
		return err
	}
	for _, l := range dn.Lines {
		_, err = tx.Exec(ctx, `INSERT INTO dn_lines (id,dn_id,so_line_id,item_code,item_name,unit,qty_delivered,qty_returned,unit_price,line_total,cost_price) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
			l.ID, dn.ID, l.SOLineID, l.ItemCode, l.ItemName, l.Unit, l.QtyDelivered, l.QtyReturned, l.UnitPrice, l.LineTotal, l.CostPrice)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func scanDN(row pgx.Row) (domain.DeliveryNote, error) {
	var dn domain.DeliveryNote
	var delDate, crAt string
	err := row.Scan(&dn.ID, &dn.CompanyID, &dn.DNNumber, &dn.SOID, &delDate, &dn.Warehouse, &dn.ShippingMethod, &dn.CarrierName, &dn.TrackingNumber, &dn.DeliveryAddress, &dn.Status, &dn.Notes, &dn.CreatedBy, &crAt)
	if err != nil {
		return dn, err
	}
	dn.DeliveryDate, _ = time.Parse(time.RFC3339, delDate)
	dn.CreatedAt, _ = time.Parse(time.RFC3339, crAt)
	return dn, nil
}

func (r *PGSaleRepo) GetDN(ctx context.Context, id string) (*domain.DeliveryNote, error) {
	dn, err := scanDN(r.pool.QueryRow(ctx, `SELECT id,company_id,dn_number,so_id,delivery_date,warehouse,shipping_method,carrier_name,tracking_number,delivery_address,status,notes,created_by,created_at FROM delivery_notes WHERE id=$1`, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDNNotFound
		}
		return nil, err
	}
	dn.Lines, _ = r.GetDNLines(ctx, id)
	return &dn, nil
}

func (r *PGSaleRepo) GetDNByNumber(ctx context.Context, companyID, dnNumber string) (*domain.DeliveryNote, error) {
	dn, err := scanDN(r.pool.QueryRow(ctx, `SELECT id,company_id,dn_number,so_id,delivery_date,warehouse,shipping_method,carrier_name,tracking_number,delivery_address,status,notes,created_by,created_at FROM delivery_notes WHERE company_id=$1 AND dn_number=$2`, companyID, dnNumber))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDNNotFound
		}
		return nil, err
	}
	dn.Lines, _ = r.GetDNLines(ctx, dn.ID)
	return &dn, nil
}

func (r *PGSaleRepo) ListDNs(ctx context.Context, filter domain.DeliveryNoteFilter) ([]domain.DeliveryNote, int, error) {
	where := ""
	args := []interface{}{filter.CompanyID}
	n := 2
	if filter.SOID != "" {
		where += fmt.Sprintf(" AND so_id=$%d", n); args = append(args, filter.SOID); n++
	}
	if filter.Status != "" {
		where += fmt.Sprintf(" AND status=$%d", n); args = append(args, string(filter.Status)); n++
	}
	if filter.FromDate != "" {
		where += fmt.Sprintf(" AND delivery_date>=$%d", n); args = append(args, filter.FromDate); n++
	}
	if filter.ToDate != "" {
		where += fmt.Sprintf(" AND delivery_date<=$%d", n); args = append(args, filter.ToDate); n++
	}
	var total int
	cq := `SELECT COUNT(*) FROM delivery_notes WHERE company_id=$1` + where
	if err := r.pool.QueryRow(ctx, cq, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	q := `SELECT id,company_id,dn_number,so_id,delivery_date,warehouse,shipping_method,carrier_name,tracking_number,delivery_address,status,notes,created_by,created_at FROM delivery_notes WHERE company_id=$1` + where
	dArgs := make([]interface{}, len(args))
	copy(dArgs, args)
	q += " ORDER BY delivery_date DESC"
	if filter.Limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", n); dArgs = append(dArgs, filter.Limit); n++
	}
	if filter.Offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", n); dArgs = append(dArgs, filter.Offset)
	}
	rows, err := r.pool.Query(ctx, q, dArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []domain.DeliveryNote
	for rows.Next() {
		dn, err := scanDN(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, dn)
	}
	return out, total, nil
}

func (r *PGSaleRepo) UpdateDN(ctx context.Context, dn *domain.DeliveryNote) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `UPDATE delivery_notes SET dn_number=$1,so_id=$2,delivery_date=$3,warehouse=$4,shipping_method=$5,carrier_name=$6,tracking_number=$7,delivery_address=$8,status=$9,notes=$10 WHERE id=$11`,
		dn.DNNumber, dn.SOID, dn.DeliveryDate.Format(time.RFC3339), dn.Warehouse, dn.ShippingMethod, dn.CarrierName, dn.TrackingNumber, dn.DeliveryAddress, dn.Status, dn.Notes, dn.ID)
	if err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `DELETE FROM dn_lines WHERE dn_id=$1`, dn.ID); err != nil {
		return err
	}
	for _, l := range dn.Lines {
		_, err = tx.Exec(ctx, `INSERT INTO dn_lines (id,dn_id,so_line_id,item_code,item_name,unit,qty_delivered,qty_returned,unit_price,line_total,cost_price) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
			l.ID, dn.ID, l.SOLineID, l.ItemCode, l.ItemName, l.Unit, l.QtyDelivered, l.QtyReturned, l.UnitPrice, l.LineTotal, l.CostPrice)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *PGSaleRepo) UpdateDNStatus(ctx context.Context, id string, status domain.DNStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE delivery_notes SET status=$1 WHERE id=$2`, string(status), id)
	return err
}

func (r *PGSaleRepo) GetDNLines(ctx context.Context, dnID string) ([]domain.DNLine, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,dn_id,so_line_id,item_code,item_name,unit,qty_delivered,qty_returned,unit_price,line_total,cost_price FROM dn_lines WHERE dn_id=$1`, dnID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.DNLine
	for rows.Next() {
		var l domain.DNLine
		if err := rows.Scan(&l.ID, &l.DNID, &l.SOLineID, &l.ItemCode, &l.ItemName, &l.Unit, &l.QtyDelivered, &l.QtyReturned, &l.UnitPrice, &l.LineTotal, &l.CostPrice); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, nil
}

func (r *PGSaleRepo) CreateDNLines(ctx context.Context, items []domain.DNLine) error {
	for _, l := range items {
		_, err := r.pool.Exec(ctx, `INSERT INTO dn_lines (id,dn_id,so_line_id,item_code,item_name,unit,qty_delivered,qty_returned,unit_price,line_total,cost_price) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
			l.ID, l.DNID, l.SOLineID, l.ItemCode, l.ItemName, l.Unit, l.QtyDelivered, l.QtyReturned, l.UnitPrice, l.LineTotal, l.CostPrice)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PGSaleRepo) UpdateDNLines(ctx context.Context, items []domain.DNLine) error {
	for _, l := range items {
		_, err := r.pool.Exec(ctx, `UPDATE dn_lines SET item_code=$1,item_name=$2,unit=$3,qty_delivered=$4,qty_returned=$5,unit_price=$6,line_total=$7,cost_price=$8 WHERE id=$9`,
			l.ItemCode, l.ItemName, l.Unit, l.QtyDelivered, l.QtyReturned, l.UnitPrice, l.LineTotal, l.CostPrice, l.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PGSaleRepo) NextDNNumber(ctx context.Context, companyID, yyyymm string) (string, error) {
	var n int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*)+1 FROM delivery_notes WHERE company_id=$1 AND dn_number LIKE $2`, companyID, "DN-"+yyyymm+"-%").Scan(&n)
	if err != nil {
		return fmt.Sprintf("DN-%s-00001", yyyymm), nil
	}
	return fmt.Sprintf("DN-%s-%05d", yyyymm, n), nil
}

func (r *PGSaleRepo) NextInvNumber(ctx context.Context, companyID, yyyymm string) (string, error) {
	var n int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*)+1 FROM customer_invoices WHERE company_id=$1 AND invoice_number LIKE $2`, companyID, "INV-"+yyyymm+"-%").Scan(&n)
	if err != nil {
		return fmt.Sprintf("INV-%s-00001", yyyymm), nil
	}
	return fmt.Sprintf("INV-%s-%05d", yyyymm, n), nil
}

// ─── Customer Invoice ─────────────────────────────────────────────────────

func (r *PGSaleRepo) CreateInvoice(ctx context.Context, inv *domain.CustomerInvoice) error {
	now := time.Now().UTC().Format(time.RFC3339)
	if inv.Status == "" {
		inv.Status = domain.SInvDraft
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := tx.QueryRow(ctx, `INSERT INTO customer_invoices (id,company_id,invoice_number,invoice_date,so_id,dn_id,customer_id,customer_name,customer_tax_code,customer_address,invoice_type,currency,exchange_rate,subtotal,discount_amount,tax_amount,total_amount,amount_received,balance_due,due_date,invoice_note,e_invoice_data,e_invoice_code,e_invoice_status,digital_signature_id,signed_data,gdt_response,original_invoice_id,adjustment_type,status,gl_posted,gl_posted_at,notes,created_by,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35) RETURNING id`,
		inv.ID, inv.CompanyID, inv.InvoiceNumber, inv.InvoiceDate.Format(time.RFC3339), inv.SOID, inv.DNID, inv.CustomerID, inv.CustomerName, inv.CustomerTaxCode, inv.CustomerAddress, inv.InvoiceType, inv.Currency, inv.ExchangeRate, inv.Subtotal, inv.DiscountAmount, inv.TaxAmount, inv.TotalAmount, inv.AmountReceived, inv.BalanceDue, pt(inv.DueDate), inv.InvoiceNote, inv.EInvoiceData, inv.EInvoiceCode, inv.EInvStatus, inv.DigitalSignatureID, inv.SignedData, inv.GDTResponse, inv.OriginalInvoiceID, inv.AdjustmentType, inv.Status, boolToInt(inv.GLPosted), pt(inv.GLPostedAt), inv.Notes, inv.CreatedBy, now,
	).Scan(&inv.ID); err != nil {
		return err
	}
	for _, l := range inv.Lines {
		_, err = tx.Exec(ctx, `INSERT INTO inv_lines (id,invoice_id,so_line_id,dn_line_id,item_code,item_name,unit,quantity,unit_price,discount_pct,vat_rate,vat_type,line_total,line_vat_amount,revenue_account_id,vat_account_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`,
			l.ID, inv.ID, l.SOLineID, l.DNLineID, l.ItemCode, l.ItemName, l.Unit, l.Quantity, l.UnitPrice, l.DiscountPct, l.VATRate, l.VATType, l.LineTotal, l.LineVATAmount, l.RevenueAccount, l.VATAccountID)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func pt(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func scanInv(row pgx.Row) (domain.CustomerInvoice, error) {
	var inv domain.CustomerInvoice
	var glPosted int
	var dueDate, glPostedAt, invDate, soDate, dnDate string
	err := row.Scan(&inv.ID, &inv.CompanyID, &inv.InvoiceNumber, &invDate, &soDate, &dnDate, &inv.CustomerID, &inv.CustomerName, &inv.CustomerTaxCode, &inv.CustomerAddress, &inv.InvoiceType, &inv.Currency, &inv.ExchangeRate, &inv.Subtotal, &inv.DiscountAmount, &inv.TaxAmount, &inv.TotalAmount, &inv.AmountReceived, &inv.BalanceDue, &dueDate, &inv.InvoiceNote, &inv.EInvoiceData, &inv.EInvoiceCode, &inv.EInvStatus, &inv.DigitalSignatureID, &inv.SignedData, &inv.GDTResponse, &inv.OriginalInvoiceID, &inv.AdjustmentType, &inv.Status, &glPosted, &glPostedAt, &inv.Notes, &inv.CreatedBy, &inv.CreatedAt)
	if err != nil {
		return inv, err
	}
	inv.GLPosted = glPosted != 0
	inv.InvoiceDate, _ = time.Parse(time.RFC3339, invDate)
	if dueDate != "" {
		d, _ := time.Parse(time.RFC3339, dueDate)
		inv.DueDate = &d
	}
	if glPostedAt != "" {
		t, _ := time.Parse(time.RFC3339, glPostedAt)
		inv.GLPostedAt = &t
	}
	return inv, nil
}

func (r *PGSaleRepo) GetInvoice(ctx context.Context, id string) (*domain.CustomerInvoice, error) {
	inv, err := scanInv(r.pool.QueryRow(ctx, `SELECT id,company_id,invoice_number,invoice_date,so_id,dn_id,customer_id,customer_name,customer_tax_code,customer_address,invoice_type,currency,exchange_rate,subtotal,discount_amount,tax_amount,total_amount,amount_received,balance_due,due_date,invoice_note,e_invoice_data,e_invoice_code,e_invoice_status,digital_signature_id,signed_data,gdt_response,original_invoice_id,adjustment_type,status,gl_posted,gl_posted_at,notes,created_by,created_at FROM customer_invoices WHERE id=$1`, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrInvNotFound
		}
		return nil, err
	}
	inv.Lines, _ = r.GetInvoiceLines(ctx, id)
	return &inv, nil
}

func (r *PGSaleRepo) GetInvoiceByNumber(ctx context.Context, companyID, invoiceNumber string) (*domain.CustomerInvoice, error) {
	inv, err := scanInv(r.pool.QueryRow(ctx, `SELECT id,company_id,invoice_number,invoice_date,so_id,dn_id,customer_id,customer_name,customer_tax_code,customer_address,invoice_type,currency,exchange_rate,subtotal,discount_amount,tax_amount,total_amount,amount_received,balance_due,due_date,invoice_note,e_invoice_data,e_invoice_code,e_invoice_status,digital_signature_id,signed_data,gdt_response,original_invoice_id,adjustment_type,status,gl_posted,gl_posted_at,notes,created_by,created_at FROM customer_invoices WHERE company_id=$1 AND invoice_number=$2`, companyID, invoiceNumber))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrInvNotFound
		}
		return nil, err
	}
	inv.Lines, _ = r.GetInvoiceLines(ctx, inv.ID)
	return &inv, nil
}

func (r *PGSaleRepo) ListInvoices(ctx context.Context, filter domain.CustomerInvoiceFilter) ([]domain.CustomerInvoice, int, error) {
	where := ""
	args := []interface{}{filter.CompanyID}
	n := 2
	if filter.CustomerID != "" {
		where += fmt.Sprintf(" AND customer_id=$%d", n); args = append(args, filter.CustomerID); n++
	}
	if filter.Status != "" {
		where += fmt.Sprintf(" AND status=$%d", n); args = append(args, string(filter.Status)); n++
	}
	if filter.FromDate != "" {
		where += fmt.Sprintf(" AND invoice_date>=$%d", n); args = append(args, filter.FromDate); n++
	}
	if filter.ToDate != "" {
		where += fmt.Sprintf(" AND invoice_date<=$%d", n); args = append(args, filter.ToDate); n++
	}
	var total int
	cq := `SELECT COUNT(*) FROM customer_invoices WHERE company_id=$1` + where
	if err := r.pool.QueryRow(ctx, cq, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	q := `SELECT id,company_id,invoice_number,invoice_date,so_id,dn_id,customer_id,customer_name,customer_tax_code,customer_address,invoice_type,currency,exchange_rate,subtotal,discount_amount,tax_amount,total_amount,amount_received,balance_due,due_date,invoice_note,e_invoice_data,e_invoice_code,e_invoice_status,digital_signature_id,signed_data,gdt_response,original_invoice_id,adjustment_type,status,gl_posted,gl_posted_at,notes,created_by,created_at FROM customer_invoices WHERE company_id=$1` + where
	dArgs := make([]interface{}, len(args))
	copy(dArgs, args)
	q += " ORDER BY invoice_date DESC"
	if filter.Limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", n); dArgs = append(dArgs, filter.Limit); n++
	}
	if filter.Offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", n); dArgs = append(dArgs, filter.Offset)
	}
	rows, err := r.pool.Query(ctx, q, dArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []domain.CustomerInvoice
	for rows.Next() {
		inv, err := scanInv(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, inv)
	}
	return out, total, nil
}

func (r *PGSaleRepo) UpdateInvoice(ctx context.Context, inv *domain.CustomerInvoice) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `UPDATE customer_invoices SET invoice_number=$1,invoice_date=$2,so_id=$3,dn_id=$4,customer_id=$5,customer_name=$6,customer_tax_code=$7,customer_address=$8,invoice_type=$9,currency=$10,exchange_rate=$11,subtotal=$12,discount_amount=$13,tax_amount=$14,total_amount=$15,amount_received=$16,balance_due=$17,due_date=$18,invoice_note=$19,status=$20,notes=$21 WHERE id=$22`,
		inv.InvoiceNumber, inv.InvoiceDate.Format(time.RFC3339), inv.SOID, inv.DNID, inv.CustomerID, inv.CustomerName, inv.CustomerTaxCode, inv.CustomerAddress, inv.InvoiceType, inv.Currency, inv.ExchangeRate, inv.Subtotal, inv.DiscountAmount, inv.TaxAmount, inv.TotalAmount, inv.AmountReceived, inv.BalanceDue, inv.DueDate.Format(time.RFC3339), inv.InvoiceNote, inv.Status, inv.Notes, inv.ID)
	if err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `DELETE FROM inv_lines WHERE invoice_id=$1`, inv.ID); err != nil {
		return err
	}
	for _, l := range inv.Lines {
		_, err = tx.Exec(ctx, `INSERT INTO inv_lines (id,invoice_id,so_line_id,dn_line_id,item_code,item_name,unit,quantity,unit_price,discount_pct,vat_rate,vat_type,line_total,line_vat_amount,revenue_account_id,vat_account_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`,
			l.ID, inv.ID, l.SOLineID, l.DNLineID, l.ItemCode, l.ItemName, l.Unit, l.Quantity, l.UnitPrice, l.DiscountPct, l.VATRate, l.VATType, l.LineTotal, l.LineVATAmount, l.RevenueAccount, l.VATAccountID)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *PGSaleRepo) UpdateInvoiceStatus(ctx context.Context, id string, status domain.SaleInvoiceStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE customer_invoices SET status=$1 WHERE id=$2`, string(status), id)
	return err
}

func (r *PGSaleRepo) PostInvoice(ctx context.Context, id string, postedAt time.Time) error {
	_, err := r.pool.Exec(ctx, `UPDATE customer_invoices SET status=$1 WHERE id=$2`,
		string(domain.SInvPosted), id)
	return err
}

func (r *PGSaleRepo) SetInvoiceGLPosted(ctx context.Context, id string, postedAt time.Time) error {
	_, err := r.pool.Exec(ctx, `UPDATE customer_invoices SET gl_posted=1,gl_posted_at=$1 WHERE id=$2`,
		postedAt.Format(time.RFC3339), id)
	return err
}

func (r *PGSaleRepo) AllocateToInvoice(ctx context.Context, invoiceID string, amount float64) error {
	_, err := r.pool.Exec(ctx, `UPDATE customer_invoices SET amount_received=amount_received+$1,balance_due=balance_due-$1 WHERE id=$2`, amount, invoiceID)
	return err
}

func (r *PGSaleRepo) GetInvoiceLines(ctx context.Context, invoiceID string) ([]domain.InvLine, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,invoice_id,so_line_id,dn_line_id,item_code,item_name,unit,quantity,unit_price,discount_pct,vat_rate,vat_type,line_total,line_vat_amount,revenue_account_id,vat_account_id FROM inv_lines WHERE invoice_id=$1`, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.InvLine
	for rows.Next() {
		var l domain.InvLine
		if err := rows.Scan(&l.ID, &l.InvoiceID, &l.SOLineID, &l.DNLineID, &l.ItemCode, &l.ItemName, &l.Unit, &l.Quantity, &l.UnitPrice, &l.DiscountPct, &l.VATRate, &l.VATType, &l.LineTotal, &l.LineVATAmount, &l.RevenueAccount, &l.VATAccountID); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, nil
}

func (r *PGSaleRepo) CreateInvoiceLines(ctx context.Context, items []domain.InvLine) error {
	for _, l := range items {
		_, err := r.pool.Exec(ctx, `INSERT INTO inv_lines (id,invoice_id,so_line_id,dn_line_id,item_code,item_name,unit,quantity,unit_price,discount_pct,vat_rate,vat_type,line_total,line_vat_amount,revenue_account_id,vat_account_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`,
			l.ID, l.InvoiceID, l.SOLineID, l.DNLineID, l.ItemCode, l.ItemName, l.Unit, l.Quantity, l.UnitPrice, l.DiscountPct, l.VATRate, l.VATType, l.LineTotal, l.LineVATAmount, l.RevenueAccount, l.VATAccountID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PGSaleRepo) UpdateInvoiceLines(ctx context.Context, items []domain.InvLine) error {
	for _, l := range items {
		_, err := r.pool.Exec(ctx, `UPDATE inv_lines SET item_code=$1,item_name=$2,unit=$3,quantity=$4,unit_price=$5,discount_pct=$6,vat_rate=$7,vat_type=$8,line_total=$9,line_vat_amount=$10,revenue_account_id=$11,vat_account_id=$12 WHERE id=$13`,
			l.ItemCode, l.ItemName, l.Unit, l.Quantity, l.UnitPrice, l.DiscountPct, l.VATRate, l.VATType, l.LineTotal, l.LineVATAmount, l.RevenueAccount, l.VATAccountID, l.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// ─── Receipt ──────────────────────────────────────────────────────────────

func (r *PGSaleRepo) CreateReceipt(ctx context.Context, rcpt *domain.CustomerReceipt) error {
	now := time.Now().UTC().Format(time.RFC3339)
	if rcpt.Status == "" {
		rcpt.Status = domain.RcpDraft
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := tx.QueryRow(ctx, `INSERT INTO customer_receipts (id,company_id,receipt_number,customer_id,receipt_date,payment_method,bank_account_id,currency,exchange_rate,amount,unallocated_amount,reference,notes,status,created_by,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) RETURNING id`,
		rcpt.ID, rcpt.CompanyID, rcpt.ReceiptNumber, rcpt.CustomerID, rcpt.ReceiptDate.Format(time.RFC3339), rcpt.PaymentMethod, rcpt.BankAccountID, rcpt.Currency, rcpt.ExchangeRate, rcpt.Amount, rcpt.UnallocatedAmount, rcpt.Reference, rcpt.Notes, rcpt.Status, rcpt.CreatedBy, now,
	).Scan(&rcpt.ID); err != nil {
		return err
	}
	for _, a := range rcpt.Allocations {
		_, err = tx.Exec(ctx, `INSERT INTO receipt_allocations (id,receipt_id,invoice_id,allocated_amount,discount_amount) VALUES ($1,$2,$3,$4,$5)`,
			a.ID, rcpt.ID, a.InvoiceID, a.AllocatedAmount, a.DiscountAmount)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func scanReceipt(row pgx.Row) (domain.CustomerReceipt, error) {
	var r domain.CustomerReceipt
	var rcptDate, createdAt string
	err := row.Scan(&r.ID, &r.CompanyID, &r.ReceiptNumber, &r.CustomerID, &rcptDate, &r.PaymentMethod, &r.BankAccountID, &r.Currency, &r.ExchangeRate, &r.Amount, &r.UnallocatedAmount, &r.Reference, &r.Notes, &r.Status, &r.CreatedBy, &createdAt)
	if err != nil {
		return r, err
	}
	r.ReceiptDate, _ = time.Parse(time.RFC3339, rcptDate)
	r.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return r, nil
}

func (r *PGSaleRepo) GetReceipt(ctx context.Context, id string) (*domain.CustomerReceipt, error) {
	rcpt, err := scanReceipt(r.pool.QueryRow(ctx, `SELECT id,company_id,receipt_number,customer_id,receipt_date,payment_method,bank_account_id,currency,exchange_rate,amount,unallocated_amount,reference,notes,status,created_by,created_at FROM customer_receipts WHERE id=$1`, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrRcpNotFound
		}
		return nil, err
	}
	rcpt.Allocations, _ = r.GetReceiptAllocations(ctx, id)
	return &rcpt, nil
}

func (r *PGSaleRepo) GetReceiptByNumber(ctx context.Context, companyID, receiptNumber string) (*domain.CustomerReceipt, error) {
	rcpt, err := scanReceipt(r.pool.QueryRow(ctx, `SELECT id,company_id,receipt_number,customer_id,receipt_date,payment_method,bank_account_id,currency,exchange_rate,amount,unallocated_amount,reference,notes,status,created_by,created_at FROM customer_receipts WHERE company_id=$1 AND receipt_number=$2`, companyID, receiptNumber))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrRcpNotFound
		}
		return nil, err
	}
	rcpt.Allocations, _ = r.GetReceiptAllocations(ctx, rcpt.ID)
	return &rcpt, nil
}

func (r *PGSaleRepo) ListReceipts(ctx context.Context, filter domain.ReceiptFilter) ([]domain.CustomerReceipt, int, error) {
	where := ""
	args := []interface{}{filter.CompanyID}
	n := 2
	if filter.CustomerID != "" {
		where += fmt.Sprintf(" AND customer_id=$%d", n); args = append(args, filter.CustomerID); n++
	}
	if filter.Status != "" {
		where += fmt.Sprintf(" AND status=$%d", n); args = append(args, string(filter.Status)); n++
	}
	if filter.FromDate != "" {
		where += fmt.Sprintf(" AND receipt_date>=$%d", n); args = append(args, filter.FromDate); n++
	}
	if filter.ToDate != "" {
		where += fmt.Sprintf(" AND receipt_date<=$%d", n); args = append(args, filter.ToDate); n++
	}
	var total int
	cq := `SELECT COUNT(*) FROM customer_receipts WHERE company_id=$1` + where
	if err := r.pool.QueryRow(ctx, cq, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	q := `SELECT id,company_id,receipt_number,customer_id,receipt_date,payment_method,bank_account_id,currency,exchange_rate,amount,unallocated_amount,reference,notes,status,created_by,created_at FROM customer_receipts WHERE company_id=$1` + where
	dArgs := make([]interface{}, len(args))
	copy(dArgs, args)
	q += " ORDER BY receipt_date DESC"
	if filter.Limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", n); dArgs = append(dArgs, filter.Limit); n++
	}
	if filter.Offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", n); dArgs = append(dArgs, filter.Offset)
	}
	rows, err := r.pool.Query(ctx, q, dArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []domain.CustomerReceipt
	for rows.Next() {
		rcpt, err := scanReceipt(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, rcpt)
	}
	return out, total, nil
}

func (r *PGSaleRepo) UpdateReceipt(ctx context.Context, rcpt *domain.CustomerReceipt) error {
	_, err := r.pool.Exec(ctx, `UPDATE customer_receipts SET receipt_number=$1,customer_id=$2,receipt_date=$3,payment_method=$4,bank_account_id=$5,currency=$6,exchange_rate=$7,amount=$8,unallocated_amount=$9,reference=$10,notes=$11,status=$12 WHERE id=$13`,
		rcpt.ReceiptNumber, rcpt.CustomerID, rcpt.ReceiptDate.Format(time.RFC3339), rcpt.PaymentMethod, rcpt.BankAccountID, rcpt.Currency, rcpt.ExchangeRate, rcpt.Amount, rcpt.UnallocatedAmount, rcpt.Reference, rcpt.Notes, rcpt.Status, rcpt.ID)
	return err
}

func (r *PGSaleRepo) UpdateReceiptStatus(ctx context.Context, id string, status domain.ReceiptStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE customer_receipts SET status=$1 WHERE id=$2`, string(status), id)
	return err
}

func (r *PGSaleRepo) SetReceiptGLPosted(ctx context.Context, id string, postedAt time.Time) error {
	_, err := r.pool.Exec(ctx, `UPDATE customer_receipts SET gl_posted=1,gl_posted_at=$1 WHERE id=$2`,
		postedAt.Format(time.RFC3339), id)
	return err
}

func (r *PGSaleRepo) CreateReceiptAllocations(ctx context.Context, allocs []domain.RcpAllocation) error {
	for _, a := range allocs {
		_, err := r.pool.Exec(ctx, `INSERT INTO receipt_allocations (id,receipt_id,invoice_id,allocated_amount,discount_amount) VALUES ($1,$2,$3,$4,$5)`,
			a.ID, a.ReceiptID, a.InvoiceID, a.AllocatedAmount, a.DiscountAmount)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PGSaleRepo) GetReceiptAllocations(ctx context.Context, receiptID string) ([]domain.RcpAllocation, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,receipt_id,invoice_id,allocated_amount,discount_amount FROM receipt_allocations WHERE receipt_id=$1`, receiptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.RcpAllocation
	for rows.Next() {
		var a domain.RcpAllocation
		if err := rows.Scan(&a.ID, &a.ReceiptID, &a.InvoiceID, &a.AllocatedAmount, &a.DiscountAmount); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

// ─── Credit Note ──────────────────────────────────────────────────────────

func (r *PGSaleRepo) CreateCN(ctx context.Context, cn *domain.CreditNote) error {
	now := time.Now().UTC().Format(time.RFC3339)
	if cn.Status == "" {
		cn.Status = domain.CNDraft
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := tx.QueryRow(ctx, `INSERT INTO credit_notes (id,company_id,cn_number,original_invoice_id,customer_id,return_date,return_reason,return_type,dn_id,subtotal,tax_amount,total_amount,e_invoice_data,e_invoice_code,status,gl_posted,gl_posted_at,notes,created_by,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20) RETURNING id`,
		cn.ID, cn.CompanyID, cn.CNNumber, cn.OriginalInvoiceID, cn.CustomerID, cn.ReturnDate.Format(time.RFC3339), cn.ReturnReason, cn.ReturnType, cn.DNID, cn.Subtotal, cn.TaxAmount, cn.TotalAmount, cn.EInvoiceData, cn.EInvoiceCode, cn.Status, boolToInt(cn.GLPosted), pt(cn.GLPostedAt), cn.Notes, cn.CreatedBy, now,
	).Scan(&cn.ID); err != nil {
		return err
	}
	for _, l := range cn.Lines {
		_, err = tx.Exec(ctx, `INSERT INTO cn_lines (id,cn_id,invoice_line_id,item_name,unit,quantity,unit_price,vat_rate,line_total,line_vat_amount) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
			l.ID, cn.ID, l.InvLineID, l.ItemName, l.Unit, l.Quantity, l.UnitPrice, l.VATRate, l.LineTotal, l.LineVATAmount)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func scanCN(row pgx.Row) (domain.CreditNote, error) {
	var cn domain.CreditNote
	var glPosted int
	var glPostedAt, retDate string
	err := row.Scan(&cn.ID, &cn.CompanyID, &cn.CNNumber, &cn.OriginalInvoiceID, &cn.CustomerID, &retDate, &cn.ReturnReason, &cn.ReturnType, &cn.DNID, &cn.Subtotal, &cn.TaxAmount, &cn.TotalAmount, &cn.EInvoiceData, &cn.EInvoiceCode, &cn.Status, &glPosted, &glPostedAt, &cn.Notes, &cn.CreatedBy, &cn.CreatedAt)
	if err != nil {
		return cn, err
	}
	cn.GLPosted = glPosted != 0
	cn.ReturnDate, _ = time.Parse(time.RFC3339, retDate)
	if glPostedAt != "" {
		t, _ := time.Parse(time.RFC3339, glPostedAt)
		cn.GLPostedAt = &t
	}
	return cn, nil
}

func (r *PGSaleRepo) GetCN(ctx context.Context, id string) (*domain.CreditNote, error) {
	cn, err := scanCN(r.pool.QueryRow(ctx, `SELECT id,company_id,cn_number,original_invoice_id,customer_id,return_date,return_reason,return_type,dn_id,subtotal,tax_amount,total_amount,e_invoice_data,e_invoice_code,status,gl_posted,gl_posted_at,notes,created_by,created_at FROM credit_notes WHERE id=$1`, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrCNNotFound
		}
		return nil, err
	}
	cn.Lines, _ = r.GetCNLines(ctx, id)
	return &cn, nil
}

func (r *PGSaleRepo) GetCNByNumber(ctx context.Context, companyID, cnNumber string) (*domain.CreditNote, error) {
	cn, err := scanCN(r.pool.QueryRow(ctx, `SELECT id,company_id,cn_number,original_invoice_id,customer_id,return_date,return_reason,return_type,dn_id,subtotal,tax_amount,total_amount,e_invoice_data,e_invoice_code,status,gl_posted,gl_posted_at,notes,created_by,created_at FROM credit_notes WHERE company_id=$1 AND cn_number=$2`, companyID, cnNumber))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrCNNotFound
		}
		return nil, err
	}
	cn.Lines, _ = r.GetCNLines(ctx, cn.ID)
	return &cn, nil
}

func (r *PGSaleRepo) ListCNs(ctx context.Context, filter domain.CreditNoteFilter) ([]domain.CreditNote, int, error) {
	where := ""
	args := []interface{}{filter.CompanyID}
	n := 2
	if filter.CustomerID != "" {
		where += fmt.Sprintf(" AND customer_id=$%d", n); args = append(args, filter.CustomerID); n++
	}
	if filter.Status != "" {
		where += fmt.Sprintf(" AND status=$%d", n); args = append(args, string(filter.Status)); n++
	}
	if filter.FromDate != "" {
		where += fmt.Sprintf(" AND return_date>=$%d", n); args = append(args, filter.FromDate); n++
	}
	if filter.ToDate != "" {
		where += fmt.Sprintf(" AND return_date<=$%d", n); args = append(args, filter.ToDate); n++
	}
	var total int
	cq := `SELECT COUNT(*) FROM credit_notes WHERE company_id=$1` + where
	if err := r.pool.QueryRow(ctx, cq, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	q := `SELECT id,company_id,cn_number,original_invoice_id,customer_id,return_date,return_reason,return_type,dn_id,subtotal,tax_amount,total_amount,e_invoice_data,e_invoice_code,status,gl_posted,gl_posted_at,notes,created_by,created_at FROM credit_notes WHERE company_id=$1` + where
	dArgs := make([]interface{}, len(args))
	copy(dArgs, args)
	q += " ORDER BY return_date DESC"
	if filter.Limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", n); dArgs = append(dArgs, filter.Limit); n++
	}
	if filter.Offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", n); dArgs = append(dArgs, filter.Offset)
	}
	rows, err := r.pool.Query(ctx, q, dArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []domain.CreditNote
	for rows.Next() {
		cn, err := scanCN(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, cn)
	}
	return out, total, nil
}

func (r *PGSaleRepo) UpdateCN(ctx context.Context, cn *domain.CreditNote) error {
	_, err := r.pool.Exec(ctx, `UPDATE credit_notes SET cn_number=$1,original_invoice_id=$2,customer_id=$3,return_date=$4,return_reason=$5,return_type=$6,dn_id=$7,subtotal=$8,tax_amount=$9,total_amount=$10,status=$11,notes=$12 WHERE id=$13`,
		cn.CNNumber, cn.OriginalInvoiceID, cn.CustomerID, cn.ReturnDate.Format(time.RFC3339), cn.ReturnReason, cn.ReturnType, cn.DNID, cn.Subtotal, cn.TaxAmount, cn.TotalAmount, cn.Status, cn.Notes, cn.ID)
	return err
}

func (r *PGSaleRepo) UpdateCNStatus(ctx context.Context, id string, status domain.CNStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE credit_notes SET status=$1 WHERE id=$2`, string(status), id)
	return err
}

func (r *PGSaleRepo) PostCN(ctx context.Context, id string, postedAt time.Time) error {
	_, err := r.pool.Exec(ctx, `UPDATE credit_notes SET status=$1 WHERE id=$2`,
		string(domain.CNPosted), id)
	return err
}

func (r *PGSaleRepo) SetCNGLPosted(ctx context.Context, id string, postedAt time.Time) error {
	_, err := r.pool.Exec(ctx, `UPDATE credit_notes SET gl_posted=1,gl_posted_at=$1 WHERE id=$2`,
		postedAt.Format(time.RFC3339), id)
	return err
}

func (r *PGSaleRepo) GetCNLines(ctx context.Context, cnID string) ([]domain.CNLine, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,cn_id,invoice_line_id,item_name,unit,quantity,unit_price,vat_rate,line_total,line_vat_amount FROM cn_lines WHERE cn_id=$1`, cnID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.CNLine
	for rows.Next() {
		var l domain.CNLine
		if err := rows.Scan(&l.ID, &l.CNID, &l.InvLineID, &l.ItemName, &l.Unit, &l.Quantity, &l.UnitPrice, &l.VATRate, &l.LineTotal, &l.LineVATAmount); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, nil
}

func (r *PGSaleRepo) CreateCNLines(ctx context.Context, items []domain.CNLine) error {
	for _, l := range items {
		_, err := r.pool.Exec(ctx, `INSERT INTO cn_lines (id,cn_id,invoice_line_id,item_name,unit,quantity,unit_price,vat_rate,line_total,line_vat_amount) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
			l.ID, l.CNID, l.InvLineID, l.ItemName, l.Unit, l.Quantity, l.UnitPrice, l.VATRate, l.LineTotal, l.LineVATAmount)
		if err != nil {
			return err
		}
	}
	return nil
}

// ─── AR Transaction ───────────────────────────────────────────────────────

func (r *PGSaleRepo) CreateARTransaction(ctx context.Context, t *domain.ARTransaction) error {
	now := time.Now().UTC().Format(time.RFC3339)
	if t.Currency == "" {
		t.Currency = "VND"
	}
	return r.pool.QueryRow(ctx, `INSERT INTO ar_transactions (id,company_id,customer_id,invoice_id,transaction_type,transaction_date,amount,currency,reference_type,reference_id,notes,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id`,
		t.ID, t.CompanyID, t.CustomerID, t.InvoiceID, t.TransactionType, t.TransactionDate.Format(time.RFC3339), t.Amount, t.Currency, t.ReferenceType, t.ReferenceID, t.Notes, now,
	).Scan(&t.ID)
}

func scanARTransaction(row pgx.Row) (domain.ARTransaction, error) {
	var t domain.ARTransaction
	err := row.Scan(&t.ID, &t.CompanyID, &t.CustomerID, &t.InvoiceID, &t.TransactionType, &t.TransactionDate, &t.Amount, &t.Currency, &t.ReferenceType, &t.ReferenceID, &t.Notes, &t.CreatedAt)
	return t, err
}

func (r *PGSaleRepo) GetARTransaction(ctx context.Context, id string) (*domain.ARTransaction, error) {
	t, err := scanARTransaction(r.pool.QueryRow(ctx, `SELECT id,company_id,customer_id,invoice_id,transaction_type,transaction_date,amount,currency,reference_type,reference_id,notes,created_at FROM ar_transactions WHERE id=$1`, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrARTransNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *PGSaleRepo) ListARTransactions(ctx context.Context, companyID, customerID string) ([]domain.ARTransaction, error) {
	q := `SELECT id,company_id,customer_id,invoice_id,transaction_type,transaction_date,amount,currency,reference_type,reference_id,notes,created_at FROM ar_transactions WHERE customer_id=$1`
	args := []interface{}{customerID}
	if companyID != "" {
		q = `SELECT art.id,art.company_id,art.customer_id,art.invoice_id,art.transaction_type,art.transaction_date,art.amount,art.currency,art.reference_type,art.reference_id,art.notes,art.created_at FROM ar_transactions art JOIN customers c ON c.id=art.customer_id WHERE c.company_id=$1 AND art.customer_id=$2`
		args = []interface{}{companyID, customerID}
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.ARTransaction
	for rows.Next() {
		t, err := scanARTransaction(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}

func (r *PGSaleRepo) ListARTransactionsAll(ctx context.Context, companyID string, offset, limit int) ([]domain.ARTransaction, int, error) {
	var total int
	q := `SELECT COUNT(*) FROM ar_transactions art JOIN customers c ON c.id=art.customer_id WHERE c.company_id=$1`
	if err := r.pool.QueryRow(ctx, q, companyID).Scan(&total); err != nil {
		return nil, 0, err
	}
	q = `SELECT art.id,art.company_id,art.customer_id,art.invoice_id,art.transaction_type,art.transaction_date,art.amount,art.currency,art.reference_type,art.reference_id,art.notes,art.created_at FROM ar_transactions art JOIN customers c ON c.id=art.customer_id WHERE c.company_id=$1 ORDER BY art.transaction_date DESC`
	args := []interface{}{companyID}
	n := 2
	if limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", n)
		args = append(args, limit)
		n++
	}
	if offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", n)
		args = append(args, offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []domain.ARTransaction
	for rows.Next() {
		t, err := scanARTransaction(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, t)
	}
	return out, total, nil
}

// ─── Sales Quotation ──────────────────────────────────────────────────────

func (r *PGSaleRepo) CreateSQ(ctx context.Context, sq *domain.SalesQuotation) error {
	now := time.Now().UTC().Format(time.RFC3339)
	return r.pool.QueryRow(ctx, `INSERT INTO sales_quotations (id,company_id,qn_number,customer_id,valid_until,status,total_amount,created_by,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`,
		sq.ID, sq.CompanyID, sq.QNNumber, sq.CustomerID, sq.ValidUntil.Format(time.RFC3339), sq.Status, sq.TotalAmount, sq.CreatedBy, now,
	).Scan(&sq.ID)
}

func scanSQ(row pgx.Row) (domain.SalesQuotation, error) {
	var sq domain.SalesQuotation
	var validUntil, createdAt string
	err := row.Scan(&sq.ID, &sq.CompanyID, &sq.QNNumber, &sq.CustomerID, &validUntil, &sq.Status, &sq.TotalAmount, &sq.CreatedBy, &createdAt)
	if err != nil {
		return sq, err
	}
	if validUntil != "" {
		t, _ := time.Parse(time.RFC3339, validUntil)
		sq.ValidUntil = &t
	}
	sq.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return sq, nil
}

func (r *PGSaleRepo) GetSQ(ctx context.Context, id string) (*domain.SalesQuotation, error) {
	sq, err := scanSQ(r.pool.QueryRow(ctx, `SELECT id,company_id,qn_number,customer_id,valid_until,status,total_amount,created_by,created_at FROM sales_quotations WHERE id=$1`, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &sq, nil
}

func (r *PGSaleRepo) ListSQs(ctx context.Context, companyID string) ([]domain.SalesQuotation, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,company_id,qn_number,customer_id,valid_until,status,total_amount,created_by,created_at FROM sales_quotations WHERE company_id=$1 ORDER BY created_at DESC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.SalesQuotation
	for rows.Next() {
		sq, err := scanSQ(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, sq)
	}
	return out, nil
}

func (r *PGSaleRepo) UpdateSQ(ctx context.Context, sq *domain.SalesQuotation) error {
	_, err := r.pool.Exec(ctx, `UPDATE sales_quotations SET qn_number=$1,customer_id=$2,valid_until=$3,status=$4,total_amount=$5 WHERE id=$6`,
		sq.QNNumber, sq.CustomerID, sq.ValidUntil.Format(time.RFC3339), sq.Status, sq.TotalAmount, sq.ID)
	return err
}
