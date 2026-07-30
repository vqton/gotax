package repository

import (
	"context"
	"fmt"
	"gotax/internal/domain"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// ─── Customer ──────────────────────────────────────────────────────────

type PGCustomerRepo struct {
	db *gorm.DB
}

func NewPGCustomerRepo(db *gorm.DB) *PGCustomerRepo {
	return &PGCustomerRepo{db: db}
}

func customerToGORM(c *domain.Customer) *domain.CustomerGORM {
	tc := c.TaxCode
	if tc == "" {
		tc = " "
	}
	addr := c.Address
	ph := c.Phone
	em := c.Email
	return &domain.CustomerGORM{
		ID:         c.ID,
		CompanyID:  c.CompanyID,
		Code:       c.Code,
		Name:       c.Name,
		TaxCode:    &tc,
		Address:    &addr,
		Phone:      &ph,
		Email:      &em,
		IsActive:   c.Status == domain.CustomerActive,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
}

func customerFromGORM(g *domain.CustomerGORM) *domain.Customer {
	c := &domain.Customer{
		ID:        g.ID,
		CompanyID: g.CompanyID,
		Code:      g.Code,
		Name:      g.Name,
		Status:    domain.CustomerActive,
		Currency:  "VND",
		CreatedBy: "",
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}
	if g.TaxCode != nil && *g.TaxCode != " " {
		c.TaxCode = *g.TaxCode
	}
	if g.Address != nil {
		c.Address = *g.Address
	}
	if g.Phone != nil {
		c.Phone = *g.Phone
	}
	if g.Email != nil {
		c.Email = *g.Email
	}
	if !g.IsActive {
		c.Status = domain.CustomerSuspended
	}
	return c
}

func (r *PGCustomerRepo) CreateCustomer(ctx context.Context, c *domain.Customer) error {
	g := customerToGORM(c)
	if err := r.db.WithContext(ctx).Create(g).Error; err != nil {
		return err
	}
	c.ID = g.ID
	c.CreatedAt = g.CreatedAt
	c.UpdatedAt = g.UpdatedAt
	return nil
}

func (r *PGCustomerRepo) GetCustomer(ctx context.Context, id string) (*domain.Customer, error) {
	var g domain.CustomerGORM
	if err := r.db.WithContext(ctx).First(&g, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrCustomerNotFound
		}
		return nil, err
	}
	return customerFromGORM(&g), nil
}

func (r *PGCustomerRepo) GetCustomerByCode(ctx context.Context, companyID, code string) (*domain.Customer, error) {
	var g domain.CustomerGORM
	if err := r.db.WithContext(ctx).Where("company_id = ? AND code = ?", companyID, code).First(&g).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrCustomerNotFound
		}
		return nil, err
	}
	return customerFromGORM(&g), nil
}

func (r *PGCustomerRepo) ListCustomers(ctx context.Context, companyID string) ([]domain.Customer, error) {
	var gs []domain.CustomerGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("code").Find(&gs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Customer, len(gs))
	for i := range gs {
		out[i] = *customerFromGORM(&gs[i])
	}
	return out, nil
}

func (r *PGCustomerRepo) UpdateCustomer(ctx context.Context, c *domain.Customer) error {
	g := customerToGORM(c)
	return r.db.WithContext(ctx).Model(&domain.CustomerGORM{}).Where("id = ?", c.ID).Updates(map[string]interface{}{
		"code":      g.Code,
		"name":      g.Name,
		"tax_code":  g.TaxCode,
		"address":   g.Address,
		"phone":     g.Phone,
		"email":     g.Email,
		"is_active": g.IsActive,
	}).Error
}

func (r *PGCustomerRepo) DeleteCustomer(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&domain.CustomerGORM{}).Where("id = ?", id).Update("is_active", false).Error
}

// ─── Sales Order ────────────────────────────────────────────────────────

type PGSaleOrderRepo struct {
	db *gorm.DB
}

func NewPGSaleOrderRepo(db *gorm.DB) *PGSaleOrderRepo {
	return &PGSaleOrderRepo{db: db}
}

func soToGORM(so *domain.SalesOrder) *domain.SalesOrderGORM {
	g := &domain.SalesOrderGORM{
		ID:          so.ID,
		CompanyID:   so.CompanyID,
		SONumber:    so.SONumber,
		OrderDate:   so.OrderDate,
		CustomerID:  so.CustomerID,
		Subtotal:    so.Subtotal,
		TaxAmount:   so.TaxAmount,
		TotalAmount: so.TotalAmount,
		Currency:    so.Currency,
		Status:      string(so.Status),
		ApprovedBy:  &so.ApprovedBy,
		CreatedBy:   so.CreatedBy,
		CreatedAt:   so.CreatedAt,
		UpdatedAt:   so.UpdatedAt,
	}
	if so.ExpectedDate != nil {
		g.DeliveryDate = so.ExpectedDate
	}
	if so.ApprovedAt != nil {
		g.ApprovedAt = so.ApprovedAt
	}
	if so.ShippingAddress != "" {
		g.DeliveryAddress = &so.ShippingAddress
	}
	if so.PaymentTerms != "" {
		g.PaymentTerms = &so.PaymentTerms
	}
	return g
}

func soFromGORM(g *domain.SalesOrderGORM) *domain.SalesOrder {
	so := &domain.SalesOrder{
		ID:          g.ID,
		CompanyID:   g.CompanyID,
		SONumber:    g.SONumber,
		OrderDate:   g.OrderDate,
		CustomerID:  g.CustomerID,
		Subtotal:    g.Subtotal,
		TaxAmount:   g.TaxAmount,
		TotalAmount: g.TotalAmount,
		Currency:    g.Currency,
		Status:      domain.SOStatus(g.Status),
		CreatedBy:   g.CreatedBy,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}
	if g.DeliveryDate != nil {
		so.ExpectedDate = g.DeliveryDate
	}
	if g.ApprovedAt != nil {
		so.ApprovedAt = g.ApprovedAt
	}
	if g.ApprovedBy != nil {
		so.ApprovedBy = *g.ApprovedBy
	}
	if g.DeliveryAddress != nil {
		so.ShippingAddress = *g.DeliveryAddress
	}
	if g.PaymentTerms != nil {
		so.PaymentTerms = *g.PaymentTerms
	}
	return so
}

func (r *PGSaleOrderRepo) CreateSO(ctx context.Context, so *domain.SalesOrder) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		g := soToGORM(so)
		if err := tx.Create(g).Error; err != nil {
			return err
		}
		so.ID = g.ID
		so.CreatedAt = g.CreatedAt
		so.UpdatedAt = g.UpdatedAt
		for _, l := range so.Lines {
			lg := soLineToGORM(&l)
			lg.SOID = so.ID
			if err := tx.Create(lg).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *PGSaleOrderRepo) GetSO(ctx context.Context, id string) (*domain.SalesOrder, error) {
	var g domain.SalesOrderGORM
	if err := r.db.WithContext(ctx).Preload("Lines").First(&g, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrSONotFound
		}
		return nil, err
	}
	so := soFromGORM(&g)
	so.Lines = make([]domain.SOLine, len(g.Lines))
	for i, l := range g.Lines {
		so.Lines[i] = *soLineFromGORM(&l)
	}
	return so, nil
}

func (r *PGSaleOrderRepo) GetSOByNumber(ctx context.Context, companyID, soNumber string) (*domain.SalesOrder, error) {
	var g domain.SalesOrderGORM
	if err := r.db.WithContext(ctx).Preload("Lines").Where("company_id = ? AND so_number = ?", companyID, soNumber).First(&g).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrSONotFound
		}
		return nil, err
	}
	so := soFromGORM(&g)
	so.Lines = make([]domain.SOLine, len(g.Lines))
	for i, l := range g.Lines {
		so.Lines[i] = *soLineFromGORM(&l)
	}
	return so, nil
}

func (r *PGSaleOrderRepo) ListSOs(ctx context.Context, filter domain.SalesOrderFilter) ([]domain.SalesOrder, int, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&domain.SalesOrderGORM{}).Where("company_id = ?", filter.CompanyID)
	if filter.CustomerID != "" {
		q = q.Where("customer_id = ?", filter.CustomerID)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", string(filter.Status))
	}
	if filter.FromDate != "" {
		q = q.Where("order_date >= ?", filter.FromDate)
	}
	if filter.ToDate != "" {
		q = q.Where("order_date <= ?", filter.ToDate)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var gs []domain.SalesOrderGORM
	dq := r.db.WithContext(ctx).Where("company_id = ?", filter.CompanyID)
	if filter.CustomerID != "" {
		dq = dq.Where("customer_id = ?", filter.CustomerID)
	}
	if filter.Status != "" {
		dq = dq.Where("status = ?", string(filter.Status))
	}
	if filter.FromDate != "" {
		dq = dq.Where("order_date >= ?", filter.FromDate)
	}
	if filter.ToDate != "" {
		dq = dq.Where("order_date <= ?", filter.ToDate)
	}
	dq = dq.Order("order_date DESC")
	if filter.Limit > 0 {
		dq = dq.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		dq = dq.Offset(filter.Offset)
	}
	if err := dq.Find(&gs).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.SalesOrder, len(gs))
	for i := range gs {
		so := soFromGORM(&gs[i])
		r.db.WithContext(ctx).Model(&domain.SOLineGORM{}).Where("so_id = ?", gs[i].ID).Find(&so.Lines)
		out[i] = *so
	}
	return out, int(total), nil
}

func (r *PGSaleOrderRepo) UpdateSO(ctx context.Context, so *domain.SalesOrder) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.SalesOrderGORM{}).Where("id = ?", so.ID).Updates(map[string]interface{}{
			"customer_id":     so.CustomerID,
			"order_date":      so.OrderDate,
			"currency":        so.Currency,
			"payment_terms":   so.PaymentTerms,
			"delivery_address": so.ShippingAddress,
			"status":          string(so.Status),
			"notes":           so.Notes,
		}).Error; err != nil {
			return err
		}
		if err := tx.Where("so_id = ?", so.ID).Delete(&domain.SOLineGORM{}).Error; err != nil {
			return err
		}
		for _, l := range so.Lines {
			lg := soLineToGORM(&l)
			lg.SOID = so.ID
			if err := tx.Create(lg).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *PGSaleOrderRepo) UpdateSOStatus(ctx context.Context, id string, status domain.SOStatus) error {
	return r.db.WithContext(ctx).Model(&domain.SalesOrderGORM{}).Where("id = ?", id).Update("status", string(status)).Error
}

func (r *PGSaleOrderRepo) ApproveSO(ctx context.Context, id, approvedBy string, approvedAt time.Time) error {
	return r.db.WithContext(ctx).Model(&domain.SalesOrderGORM{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      string(domain.SOApproved),
		"approved_by": approvedBy,
		"approved_at": approvedAt,
	}).Error
}

func (r *PGSaleOrderRepo) CancelSO(ctx context.Context, id, cancelReason string) error {
	return r.db.WithContext(ctx).Model(&domain.SalesOrderGORM{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":   string(domain.SOCancelled),
		"cancelled_reason": cancelReason,
	}).Error
}

func soLineToGORM(l *domain.SOLine) *domain.SOLineGORM {
	disc := l.DiscountPct
	tax := l.VATRate
	return &domain.SOLineGORM{
		LineNumber: l.LineNumber,
		ItemCode:   l.ItemCode,
		ItemName:   l.ItemName,
		Quantity:   l.Quantity,
		UnitPrice:  l.UnitPrice,
		LineTotal:  l.LineTotal,
		Discount:   &disc,
		TaxRate:    &tax,
	}
}

func soLineFromGORM(g *domain.SOLineGORM) *domain.SOLine {
	l := &domain.SOLine{
		ID:        strconv.Itoa(int(g.ID)),
		SOID:      g.SOID,
		LineNumber: g.LineNumber,
		ItemCode:  g.ItemCode,
		ItemName:  g.ItemName,
		Unit:      "",
		Quantity:  g.Quantity,
		UnitPrice: g.UnitPrice,
		LineTotal: g.LineTotal,
		VATRate:   10,
	}
	if g.Discount != nil {
		l.DiscountPct = *g.Discount
	}
	if g.TaxRate != nil {
		l.VATRate = *g.TaxRate
	}
	return l
}

func (r *PGSaleOrderRepo) GetSOLines(ctx context.Context, soID string) ([]domain.SOLine, error) {
	var gs []domain.SOLineGORM
	if err := r.db.WithContext(ctx).Where("so_id = ?", soID).Order("line_number").Find(&gs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.SOLine, len(gs))
	for i := range gs {
		out[i] = *soLineFromGORM(&gs[i])
	}
	return out, nil
}

func (r *PGSaleOrderRepo) CreateSOLines(ctx context.Context, items []domain.SOLine) error {
	if len(items) == 0 {
		return nil
	}
	gs := make([]domain.SOLineGORM, len(items))
	for i := range items {
		gs[i] = *soLineToGORM(&items[i])
	}
	return r.db.WithContext(ctx).Create(&gs).Error
}

func (r *PGSaleOrderRepo) UpdateSOLines(ctx context.Context, items []domain.SOLine) error {
	for _, l := range items {
		g := soLineToGORM(&l)
		if err := r.db.WithContext(ctx).Model(&domain.SOLineGORM{}).Where("id = ?", l.ID).Updates(map[string]interface{}{
			"item_code":  g.ItemCode,
			"item_name":  g.ItemName,
			"quantity":   g.Quantity,
			"unit_price": g.UnitPrice,
			"line_total": g.LineTotal,
			"discount":   g.Discount,
			"tax_rate":   g.TaxRate,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *PGSaleOrderRepo) NextSONumber(ctx context.Context, companyID, yyyymm string) (string, error) {
	prefix := fmt.Sprintf("SO-%s-", yyyymm)
	var maxNum int
	if err := r.db.WithContext(ctx).Raw(`SELECT COALESCE(MAX(CAST(SUBSTRING(so_number FROM '([0-9]+)$') AS INTEGER)), 0) FROM sales_orders WHERE company_id = ? AND so_number LIKE ?`, companyID, prefix+"%").Scan(&maxNum).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%05d", prefix, maxNum+1), nil
}

// ─── Delivery Note ──────────────────────────────────────────────────────

type PGDeliveryNoteRepo struct {
	db *gorm.DB
}

func NewPGDeliveryNoteRepo(db *gorm.DB) *PGDeliveryNoteRepo {
	return &PGDeliveryNoteRepo{db: db}
}

func dnToGORM(dn *domain.DeliveryNote) *domain.DeliveryNoteGORM {
	g := &domain.DeliveryNoteGORM{
		ID:        dn.ID,
		CompanyID: dn.CompanyID,
		DNNumber:  dn.DNNumber,
		DNDate:    dn.DeliveryDate,
		SOID:      &dn.SOID,
		Status:    string(dn.Status),
		CreatedBy: dn.CreatedBy,
		CreatedAt: dn.CreatedAt,
		UpdatedAt: dn.UpdatedAt,
	}
	return g
}

func dnFromGORM(g *domain.DeliveryNoteGORM) *domain.DeliveryNote {
	dn := &domain.DeliveryNote{
		ID:           g.ID,
		CompanyID:    g.CompanyID,
		DNNumber:     g.DNNumber,
		DeliveryDate: g.DNDate,
		Status:       domain.DNStatus(g.Status),
		CreatedBy:    g.CreatedBy,
		CreatedAt:    g.CreatedAt,
	}
	if g.SOID != nil {
		dn.SOID = *g.SOID
	}
	return dn
}

func (r *PGDeliveryNoteRepo) CreateDN(ctx context.Context, dn *domain.DeliveryNote) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		g := dnToGORM(dn)
		if err := tx.Create(g).Error; err != nil {
			return err
		}
		dn.ID = g.ID
		dn.CreatedAt = g.CreatedAt
		for _, l := range dn.Lines {
			lg := dnLineToGORM(&l)
			lg.DNID = dn.ID
			if err := tx.Create(lg).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *PGDeliveryNoteRepo) GetDN(ctx context.Context, id string) (*domain.DeliveryNote, error) {
	var g domain.DeliveryNoteGORM
	if err := r.db.WithContext(ctx).Preload("Lines").First(&g, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrDNNotFound
		}
		return nil, err
	}
	dn := dnFromGORM(&g)
	dn.Lines = make([]domain.DNLine, len(g.Lines))
	for i, l := range g.Lines {
		dn.Lines[i] = *dnLineFromGORM(&l)
	}
	return dn, nil
}

func (r *PGDeliveryNoteRepo) GetDNByNumber(ctx context.Context, companyID, dnNumber string) (*domain.DeliveryNote, error) {
	var g domain.DeliveryNoteGORM
	if err := r.db.WithContext(ctx).Preload("Lines").Where("company_id = ? AND dn_number = ?", companyID, dnNumber).First(&g).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrDNNotFound
		}
		return nil, err
	}
	dn := dnFromGORM(&g)
	dn.Lines = make([]domain.DNLine, len(g.Lines))
	for i, l := range g.Lines {
		dn.Lines[i] = *dnLineFromGORM(&l)
	}
	return dn, nil
}

func (r *PGDeliveryNoteRepo) ListDNs(ctx context.Context, filter domain.DeliveryNoteFilter) ([]domain.DeliveryNote, int, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&domain.DeliveryNoteGORM{}).Where("company_id = ?", filter.CompanyID)
	if filter.SOID != "" {
		q = q.Where("so_id = ?", filter.SOID)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", string(filter.Status))
	}
	if filter.FromDate != "" {
		q = q.Where("dn_date >= ?", filter.FromDate)
	}
	if filter.ToDate != "" {
		q = q.Where("dn_date <= ?", filter.ToDate)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var gs []domain.DeliveryNoteGORM
	dq := r.db.WithContext(ctx).Where("company_id = ?", filter.CompanyID)
	if filter.SOID != "" {
		dq = dq.Where("so_id = ?", filter.SOID)
	}
	if filter.Status != "" {
		dq = dq.Where("status = ?", string(filter.Status))
	}
	if filter.FromDate != "" {
		dq = dq.Where("dn_date >= ?", filter.FromDate)
	}
	if filter.ToDate != "" {
		dq = dq.Where("dn_date <= ?", filter.ToDate)
	}
	dq = dq.Order("dn_date DESC")
	if filter.Limit > 0 {
		dq = dq.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		dq = dq.Offset(filter.Offset)
	}
	if err := dq.Find(&gs).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.DeliveryNote, len(gs))
	for i := range gs {
		dn := dnFromGORM(&gs[i])
		r.db.WithContext(ctx).Model(&domain.DNLineGORM{}).Where("dn_id = ?", gs[i].ID).Find(&dn.Lines)
		out[i] = *dn
	}
	return out, int(total), nil
}

func (r *PGDeliveryNoteRepo) UpdateDN(ctx context.Context, dn *domain.DeliveryNote) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.DeliveryNoteGORM{}).Where("id = ?", dn.ID).Updates(map[string]interface{}{
			"dn_number": dn.DNNumber,
			"dn_date":   dn.DeliveryDate,
			"so_id":     dn.SOID,
			"status":    string(dn.Status),
		}).Error; err != nil {
			return err
		}
		if err := tx.Where("dn_id = ?", dn.ID).Delete(&domain.DNLineGORM{}).Error; err != nil {
			return err
		}
		for _, l := range dn.Lines {
			lg := dnLineToGORM(&l)
			lg.DNID = dn.ID
			if err := tx.Create(lg).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *PGDeliveryNoteRepo) UpdateDNStatus(ctx context.Context, id string, status domain.DNStatus) error {
	return r.db.WithContext(ctx).Model(&domain.DeliveryNoteGORM{}).Where("id = ?", id).Update("status", string(status)).Error
}

func dnLineToGORM(l *domain.DNLine) *domain.DNLineGORM {
	solID := l.SOLineID
	return &domain.DNLineGORM{
		LineNumber: 0,
		SOLineID:   &solID,
		ItemCode:   l.ItemCode,
		QtyShipped: l.QtyDelivered,
		Unit:       &l.Unit,
	}
}

func dnLineFromGORM(g *domain.DNLineGORM) *domain.DNLine {
	l := &domain.DNLine{
		ID:          strconv.Itoa(int(g.ID)),
		DNID:        g.DNID,
		ItemCode:    g.ItemCode,
		QtyDelivered: g.QtyShipped,
		UnitPrice:   0,
		LineTotal:   0,
	}
	if g.SOLineID != nil {
		l.SOLineID = *g.SOLineID
	}
	if g.Unit != nil {
		l.Unit = *g.Unit
	}
	return l
}

func (r *PGDeliveryNoteRepo) GetDNLines(ctx context.Context, dnID string) ([]domain.DNLine, error) {
	var gs []domain.DNLineGORM
	if err := r.db.WithContext(ctx).Where("dn_id = ?", dnID).Find(&gs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.DNLine, len(gs))
	for i := range gs {
		out[i] = *dnLineFromGORM(&gs[i])
	}
	return out, nil
}

func (r *PGDeliveryNoteRepo) CreateDNLines(ctx context.Context, items []domain.DNLine) error {
	if len(items) == 0 {
		return nil
	}
	gs := make([]domain.DNLineGORM, len(items))
	for i := range items {
		gs[i] = *dnLineToGORM(&items[i])
	}
	return r.db.WithContext(ctx).Create(&gs).Error
}

func (r *PGDeliveryNoteRepo) UpdateDNLines(ctx context.Context, items []domain.DNLine) error {
	for _, l := range items {
		g := dnLineToGORM(&l)
		if err := r.db.WithContext(ctx).Model(&domain.DNLineGORM{}).Where("id = ?", l.ID).Updates(map[string]interface{}{
			"item_code":  g.ItemCode,
			"qty_shipped": g.QtyShipped,
			"unit":       g.Unit,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *PGDeliveryNoteRepo) NextDNNumber(ctx context.Context, companyID, yyyymm string) (string, error) {
	prefix := fmt.Sprintf("DN-%s-", yyyymm)
	var maxNum int
	if err := r.db.WithContext(ctx).Raw(`SELECT COALESCE(MAX(CAST(SUBSTRING(dn_number FROM '([0-9]+)$') AS INTEGER)), 0) FROM delivery_notes WHERE company_id = ? AND dn_number LIKE ?`, companyID, prefix+"%").Scan(&maxNum).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%05d", prefix, maxNum+1), nil
}

// ─── Customer Invoice ───────────────────────────────────────────────────

type PGCustomerInvoiceRepo struct {
	db *gorm.DB
}

func NewPGCustomerInvoiceRepo(db *gorm.DB) *PGCustomerInvoiceRepo {
	return &PGCustomerInvoiceRepo{db: db}
}

func cinvToGORM(inv *domain.CustomerInvoice) *domain.CustomerInvoiceGORM {
	g := &domain.CustomerInvoiceGORM{
		ID:            inv.ID,
		CompanyID:     inv.CompanyID,
		InvoiceNumber: inv.InvoiceNumber,
		InvoiceDate:   inv.InvoiceDate,
		CustomerID:    inv.CustomerID,
		Subtotal:      inv.Subtotal,
		TaxAmount:     inv.TaxAmount,
		GrandTotal:    inv.TotalAmount,
		Currency:      inv.Currency,
		Status:        string(inv.Status),
		GLPosted:      inv.GLPosted,
		CreatedBy:     inv.CreatedBy,
		CreatedAt:     inv.CreatedAt,
		UpdatedAt:     inv.UpdatedAt,
	}
	if inv.DueDate != nil {
		g.DueDate = inv.DueDate
	}
	if inv.GLPostedAt != nil {
		g.PostedAt = inv.GLPostedAt
	}
	return g
}

func cinvFromGORM(g *domain.CustomerInvoiceGORM) *domain.CustomerInvoice {
	inv := &domain.CustomerInvoice{
		ID:            g.ID,
		CompanyID:     g.CompanyID,
		InvoiceNumber: g.InvoiceNumber,
		InvoiceDate:   g.InvoiceDate,
		CustomerID:    g.CustomerID,
		Subtotal:      g.Subtotal,
		TaxAmount:     g.TaxAmount,
		TotalAmount:   g.GrandTotal,
		Currency:      g.Currency,
		Status:        domain.SaleInvoiceStatus(g.Status),
		GLPosted:      g.GLPosted,
		CreatedBy:     g.CreatedBy,
		CreatedAt:     g.CreatedAt,
		AmountReceived: 0,
		BalanceDue:     g.GrandTotal,
		CustomerName:   "",
		CustomerTaxCode: "",
		InvoiceType:    "",
		ExchangeRate:   1,
	}
	if g.DueDate != nil {
		inv.DueDate = g.DueDate
	}
	if g.PostedAt != nil {
		inv.GLPostedAt = g.PostedAt
	}
	return inv
}

func (r *PGCustomerInvoiceRepo) CreateInvoice(ctx context.Context, inv *domain.CustomerInvoice) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		g := cinvToGORM(inv)
		if err := tx.Create(g).Error; err != nil {
			return err
		}
		inv.ID = g.ID
		inv.CreatedAt = g.CreatedAt
		inv.UpdatedAt = g.UpdatedAt
		for _, l := range inv.Lines {
			lg := invLineToGORM(&l)
			lg.InvoiceID = inv.ID
			if err := tx.Create(lg).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *PGCustomerInvoiceRepo) GetInvoice(ctx context.Context, id string) (*domain.CustomerInvoice, error) {
	var g domain.CustomerInvoiceGORM
	if err := r.db.WithContext(ctx).Preload("Lines").First(&g, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrInvNotFound
		}
		return nil, err
	}
	inv := cinvFromGORM(&g)
	inv.Lines = make([]domain.InvLine, len(g.Lines))
	for i, l := range g.Lines {
		inv.Lines[i] = *invLineFromGORM(&l)
	}
	return inv, nil
}

func (r *PGCustomerInvoiceRepo) GetInvoiceByNumber(ctx context.Context, companyID, invoiceNumber string) (*domain.CustomerInvoice, error) {
	var g domain.CustomerInvoiceGORM
	if err := r.db.WithContext(ctx).Preload("Lines").Where("company_id = ? AND invoice_number = ?", companyID, invoiceNumber).First(&g).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrInvNotFound
		}
		return nil, err
	}
	inv := cinvFromGORM(&g)
	inv.Lines = make([]domain.InvLine, len(g.Lines))
	for i, l := range g.Lines {
		inv.Lines[i] = *invLineFromGORM(&l)
	}
	return inv, nil
}

func (r *PGCustomerInvoiceRepo) ListInvoices(ctx context.Context, filter domain.CustomerInvoiceFilter) ([]domain.CustomerInvoice, int, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&domain.CustomerInvoiceGORM{}).Where("company_id = ?", filter.CompanyID)
	if filter.CustomerID != "" {
		q = q.Where("customer_id = ?", filter.CustomerID)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", string(filter.Status))
	}
	if filter.FromDate != "" {
		q = q.Where("invoice_date >= ?", filter.FromDate)
	}
	if filter.ToDate != "" {
		q = q.Where("invoice_date <= ?", filter.ToDate)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var gs []domain.CustomerInvoiceGORM
	dq := r.db.WithContext(ctx).Where("company_id = ?", filter.CompanyID)
	if filter.CustomerID != "" {
		dq = dq.Where("customer_id = ?", filter.CustomerID)
	}
	if filter.Status != "" {
		dq = dq.Where("status = ?", string(filter.Status))
	}
	if filter.FromDate != "" {
		dq = dq.Where("invoice_date >= ?", filter.FromDate)
	}
	if filter.ToDate != "" {
		dq = dq.Where("invoice_date <= ?", filter.ToDate)
	}
	dq = dq.Order("invoice_date DESC")
	if filter.Limit > 0 {
		dq = dq.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		dq = dq.Offset(filter.Offset)
	}
	if err := dq.Find(&gs).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.CustomerInvoice, len(gs))
	for i := range gs {
		inv := cinvFromGORM(&gs[i])
		r.db.WithContext(ctx).Model(&domain.InvLineGORM{}).Where("invoice_id = ?", gs[i].ID).Find(&inv.Lines)
		out[i] = *inv
	}
	return out, int(total), nil
}

func (r *PGCustomerInvoiceRepo) UpdateInvoice(ctx context.Context, inv *domain.CustomerInvoice) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.CustomerInvoiceGORM{}).Where("id = ?", inv.ID).Updates(map[string]interface{}{
			"invoice_date": inv.InvoiceDate,
			"customer_id":  inv.CustomerID,
			"subtotal":     inv.Subtotal,
			"tax_amount":   inv.TaxAmount,
			"grand_total":  inv.TotalAmount,
			"due_date":     inv.DueDate,
			"status":       string(inv.Status),
			"notes":        inv.Notes,
		}).Error; err != nil {
			return err
		}
		if err := tx.Where("invoice_id = ?", inv.ID).Delete(&domain.InvLineGORM{}).Error; err != nil {
			return err
		}
		for _, l := range inv.Lines {
			lg := invLineToGORM(&l)
			lg.InvoiceID = inv.ID
			if err := tx.Create(lg).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *PGCustomerInvoiceRepo) UpdateInvoiceStatus(ctx context.Context, id string, status domain.SaleInvoiceStatus) error {
	return r.db.WithContext(ctx).Model(&domain.CustomerInvoiceGORM{}).Where("id = ?", id).Update("status", string(status)).Error
}

func (r *PGCustomerInvoiceRepo) PostInvoice(ctx context.Context, id string, postedAt time.Time) error {
	return r.db.WithContext(ctx).Model(&domain.CustomerInvoiceGORM{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status": string(domain.SInvPosted),
	}).Error
}

func (r *PGCustomerInvoiceRepo) SetInvoiceGLPosted(ctx context.Context, id string, postedAt time.Time) error {
	return r.db.WithContext(ctx).Model(&domain.CustomerInvoiceGORM{}).Where("id = ?", id).Updates(map[string]interface{}{
		"gl_posted": true,
		"posted_at": postedAt,
	}).Error
}

func (r *PGCustomerInvoiceRepo) AllocateToInvoice(ctx context.Context, invoiceID string, amount float64) error {
	return r.db.WithContext(ctx).Model(&domain.CustomerInvoiceGORM{}).Where("id = ?", invoiceID).Update("grand_total", gorm.Expr("grand_total - ?", amount)).Error
}

func invLineToGORM(l *domain.InvLine) *domain.InvLineGORM {
	disc := l.DiscountPct
	tax := l.VATRate
	taxAmt := l.LineVATAmount
	return &domain.InvLineGORM{
		LineNumber: 0,
		ItemCode:   l.ItemCode,
		ItemName:   l.ItemName,
		Quantity:   l.Quantity,
		UnitPrice:  l.UnitPrice,
		Discount:   &disc,
		LineTotal:  l.LineTotal,
		TaxRate:    &tax,
		TaxAmount:  &taxAmt,
	}
}

func invLineFromGORM(g *domain.InvLineGORM) *domain.InvLine {
	l := &domain.InvLine{
		ID:        strconv.Itoa(int(g.ID)),
		InvoiceID: g.InvoiceID,
		ItemCode:  g.ItemCode,
		ItemName:  g.ItemName,
		Unit:      "",
		Quantity:  g.Quantity,
		UnitPrice: g.UnitPrice,
		LineTotal: g.LineTotal,
		VATRate:   10,
		RevenueAccount: "",
		VATAccountID:  "",
	}
	if g.Discount != nil {
		l.DiscountPct = *g.Discount
	}
	if g.TaxRate != nil {
		l.VATRate = *g.TaxRate
	}
	if g.TaxAmount != nil {
		l.LineVATAmount = *g.TaxAmount
	}
	return l
}

func (r *PGCustomerInvoiceRepo) GetInvoiceLines(ctx context.Context, invoiceID string) ([]domain.InvLine, error) {
	var gs []domain.InvLineGORM
	if err := r.db.WithContext(ctx).Where("invoice_id = ?", invoiceID).Find(&gs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.InvLine, len(gs))
	for i := range gs {
		out[i] = *invLineFromGORM(&gs[i])
	}
	return out, nil
}

func (r *PGCustomerInvoiceRepo) CreateInvoiceLines(ctx context.Context, items []domain.InvLine) error {
	if len(items) == 0 {
		return nil
	}
	gs := make([]domain.InvLineGORM, len(items))
	for i := range items {
		gs[i] = *invLineToGORM(&items[i])
	}
	return r.db.WithContext(ctx).Create(&gs).Error
}

func (r *PGCustomerInvoiceRepo) UpdateInvoiceLines(ctx context.Context, items []domain.InvLine) error {
	for _, l := range items {
		g := invLineToGORM(&l)
		if err := r.db.WithContext(ctx).Model(&domain.InvLineGORM{}).Where("id = ?", l.ID).Updates(map[string]interface{}{
			"item_code":  g.ItemCode,
			"item_name":  g.ItemName,
			"quantity":   g.Quantity,
			"unit_price": g.UnitPrice,
			"line_total": g.LineTotal,
			"discount":   g.Discount,
			"tax_rate":   g.TaxRate,
			"tax_amount": g.TaxAmount,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *PGCustomerInvoiceRepo) NextInvNumber(ctx context.Context, companyID, yyyymm string) (string, error) {
	prefix := fmt.Sprintf("INV-%s-", yyyymm)
	var maxNum int
	if err := r.db.WithContext(ctx).Raw(`SELECT COALESCE(MAX(CAST(SUBSTRING(invoice_number FROM '([0-9]+)$') AS INTEGER)), 0) FROM customer_invoices WHERE company_id = ? AND invoice_number LIKE ?`, companyID, prefix+"%").Scan(&maxNum).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%05d", prefix, maxNum+1), nil
}

// ─── Customer Receipt ──────────────────────────────────────────────────

type PGCustomerReceiptRepo struct {
	db *gorm.DB
}

func NewPGCustomerReceiptRepo(db *gorm.DB) *PGCustomerReceiptRepo {
	return &PGCustomerReceiptRepo{db: db}
}

func rcptToGORM(r *domain.CustomerReceipt) *domain.CustomerReceiptGORM {
	return &domain.CustomerReceiptGORM{
		ID:            r.ID,
		CompanyID:     r.CompanyID,
		ReceiptNumber: r.ReceiptNumber,
		ReceiptDate:   r.ReceiptDate,
		CustomerID:    r.CustomerID,
		PaymentMethod: r.PaymentMethod,
		TotalReceived: r.Amount,
		Status:        string(r.Status),
		GLPosted:      r.GLPosted,
		CreatedBy:     r.CreatedBy,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
}

func rcptFromGORM(g *domain.CustomerReceiptGORM) *domain.CustomerReceipt {
	r := &domain.CustomerReceipt{
		ID:            g.ID,
		CompanyID:     g.CompanyID,
		ReceiptNumber: g.ReceiptNumber,
		ReceiptDate:   g.ReceiptDate,
		CustomerID:    g.CustomerID,
		PaymentMethod: g.PaymentMethod,
		Amount:        g.TotalReceived,
		UnallocatedAmount: g.TotalReceived,
		Status:        domain.ReceiptStatus(g.Status),
		GLPosted:      g.GLPosted,
		CreatedBy:     g.CreatedBy,
		CreatedAt:     g.CreatedAt,
		Currency:      "VND",
		ExchangeRate:  1,
	}
	return r
}

func (r *PGCustomerReceiptRepo) CreateReceipt(ctx context.Context, rcpt *domain.CustomerReceipt) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		g := rcptToGORM(rcpt)
		if err := tx.Create(g).Error; err != nil {
			return err
		}
		rcpt.ID = g.ID
		rcpt.CreatedAt = g.CreatedAt
		for _, a := range rcpt.Allocations {
			ag := rcpAllocToGORM(&a)
			ag.ReceiptID = rcpt.ID
			if err := tx.Create(ag).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *PGCustomerReceiptRepo) GetReceipt(ctx context.Context, id string) (*domain.CustomerReceipt, error) {
	var g domain.CustomerReceiptGORM
	if err := r.db.WithContext(ctx).Preload("Allocations").First(&g, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrRcpNotFound
		}
		return nil, err
	}
	rcpt := rcptFromGORM(&g)
	rcpt.Allocations = make([]domain.RcpAllocation, len(g.Allocations))
	for i, a := range g.Allocations {
		rcpt.Allocations[i] = *rcpAllocFromGORM(&a)
	}
	return rcpt, nil
}

func (r *PGCustomerReceiptRepo) GetReceiptByNumber(ctx context.Context, companyID, receiptNumber string) (*domain.CustomerReceipt, error) {
	var g domain.CustomerReceiptGORM
	if err := r.db.WithContext(ctx).Preload("Allocations").Where("company_id = ? AND receipt_number = ?", companyID, receiptNumber).First(&g).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrRcpNotFound
		}
		return nil, err
	}
	rcpt := rcptFromGORM(&g)
	rcpt.Allocations = make([]domain.RcpAllocation, len(g.Allocations))
	for i, a := range g.Allocations {
		rcpt.Allocations[i] = *rcpAllocFromGORM(&a)
	}
	return rcpt, nil
}

func (r *PGCustomerReceiptRepo) ListReceipts(ctx context.Context, filter domain.ReceiptFilter) ([]domain.CustomerReceipt, int, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&domain.CustomerReceiptGORM{}).Where("company_id = ?", filter.CompanyID)
	if filter.CustomerID != "" {
		q = q.Where("customer_id = ?", filter.CustomerID)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", string(filter.Status))
	}
	if filter.FromDate != "" {
		q = q.Where("receipt_date >= ?", filter.FromDate)
	}
	if filter.ToDate != "" {
		q = q.Where("receipt_date <= ?", filter.ToDate)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var gs []domain.CustomerReceiptGORM
	dq := r.db.WithContext(ctx).Where("company_id = ?", filter.CompanyID)
	if filter.CustomerID != "" {
		dq = dq.Where("customer_id = ?", filter.CustomerID)
	}
	if filter.Status != "" {
		dq = dq.Where("status = ?", string(filter.Status))
	}
	if filter.FromDate != "" {
		dq = dq.Where("receipt_date >= ?", filter.FromDate)
	}
	if filter.ToDate != "" {
		dq = dq.Where("receipt_date <= ?", filter.ToDate)
	}
	dq = dq.Order("receipt_date DESC")
	if filter.Limit > 0 {
		dq = dq.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		dq = dq.Offset(filter.Offset)
	}
	if err := dq.Find(&gs).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.CustomerReceipt, len(gs))
	for i := range gs {
		rcpt := rcptFromGORM(&gs[i])
		r.db.WithContext(ctx).Model(&domain.RcpAllocationGORM{}).Where("receipt_id = ?", gs[i].ID).Find(&rcpt.Allocations)
		out[i] = *rcpt
	}
	return out, int(total), nil
}

func (r *PGCustomerReceiptRepo) UpdateReceipt(ctx context.Context, rcpt *domain.CustomerReceipt) error {
	g := rcptToGORM(rcpt)
	return r.db.WithContext(ctx).Model(&domain.CustomerReceiptGORM{}).Where("id = ?", rcpt.ID).Updates(map[string]interface{}{
		"receipt_number": g.ReceiptNumber,
		"receipt_date":   g.ReceiptDate,
		"customer_id":    g.CustomerID,
		"payment_method": g.PaymentMethod,
		"total_received": g.TotalReceived,
		"status":         g.Status,
	}).Error
}

func (r *PGCustomerReceiptRepo) UpdateReceiptStatus(ctx context.Context, id string, status domain.ReceiptStatus) error {
	return r.db.WithContext(ctx).Model(&domain.CustomerReceiptGORM{}).Where("id = ?", id).Update("status", string(status)).Error
}

func (r *PGCustomerReceiptRepo) SetReceiptGLPosted(ctx context.Context, id string, postedAt time.Time) error {
	return r.db.WithContext(ctx).Model(&domain.CustomerReceiptGORM{}).Where("id = ?", id).Updates(map[string]interface{}{
		"gl_posted": true,
	}).Error
}

func rcpAllocToGORM(a *domain.RcpAllocation) *domain.RcpAllocationGORM {
	return &domain.RcpAllocationGORM{
		InvoiceID: a.InvoiceID,
		Allocated: a.AllocatedAmount,
	}
}

func rcpAllocFromGORM(g *domain.RcpAllocationGORM) *domain.RcpAllocation {
	return &domain.RcpAllocation{
		ID:              strconv.Itoa(int(g.ID)),
		ReceiptID:       g.ReceiptID,
		InvoiceID:       g.InvoiceID,
		AllocatedAmount: g.Allocated,
	}
}

func (r *PGCustomerReceiptRepo) CreateReceiptAllocations(ctx context.Context, allocs []domain.RcpAllocation) error {
	if len(allocs) == 0 {
		return nil
	}
	gs := make([]domain.RcpAllocationGORM, len(allocs))
	for i := range allocs {
		gs[i] = *rcpAllocToGORM(&allocs[i])
	}
	return r.db.WithContext(ctx).Create(&gs).Error
}

func (r *PGCustomerReceiptRepo) GetReceiptAllocations(ctx context.Context, receiptID string) ([]domain.RcpAllocation, error) {
	var gs []domain.RcpAllocationGORM
	if err := r.db.WithContext(ctx).Where("receipt_id = ?", receiptID).Find(&gs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.RcpAllocation, len(gs))
	for i := range gs {
		out[i] = *rcpAllocFromGORM(&gs[i])
	}
	return out, nil
}

// ─── Credit Note ────────────────────────────────────────────────────────

type PGCreditNoteRepo struct {
	db *gorm.DB
}

func NewPGCreditNoteRepo(db *gorm.DB) *PGCreditNoteRepo {
	return &PGCreditNoteRepo{db: db}
}

func cnToGORM(cn *domain.CreditNote) *domain.CreditNoteGORM {
	g := &domain.CreditNoteGORM{
		ID:         cn.ID,
		CompanyID:  cn.CompanyID,
		CNNumber:   cn.CNNumber,
		CNDate:     cn.ReturnDate,
		InvoiceID:  &cn.OriginalInvoiceID,
		CustomerID: cn.CustomerID,
		Reason:     &cn.ReturnReason,
		Subtotal:   cn.Subtotal,
		TaxAmount:  cn.TaxAmount,
		GrandTotal: cn.TotalAmount,
		Status:     string(cn.Status),
		GLPosted:   cn.GLPosted,
		CreatedBy:  cn.CreatedBy,
		CreatedAt:  cn.CreatedAt,
		UpdatedAt:  cn.UpdatedAt,
	}
	if cn.GLPostedAt != nil {
		g.PostedAt = cn.GLPostedAt
	}
	return g
}

func cnFromGORM(g *domain.CreditNoteGORM) *domain.CreditNote {
	cn := &domain.CreditNote{
		ID:         g.ID,
		CompanyID:  g.CompanyID,
		CNNumber:   g.CNNumber,
		ReturnDate: g.CNDate,
		CustomerID: g.CustomerID,
		Subtotal:   g.Subtotal,
		TaxAmount:  g.TaxAmount,
		TotalAmount: g.GrandTotal,
		Status:     domain.CNStatus(g.Status),
		GLPosted:   g.GLPosted,
		CreatedBy:  g.CreatedBy,
		CreatedAt:  g.CreatedAt,
	}
	if g.InvoiceID != nil {
		cn.OriginalInvoiceID = *g.InvoiceID
	}
	if g.Reason != nil {
		cn.ReturnReason = *g.Reason
	}
	if g.PostedAt != nil {
		cn.GLPostedAt = g.PostedAt
	}
	return cn
}

func (r *PGCreditNoteRepo) CreateCN(ctx context.Context, cn *domain.CreditNote) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		g := cnToGORM(cn)
		if err := tx.Create(g).Error; err != nil {
			return err
		}
		cn.ID = g.ID
		cn.CreatedAt = g.CreatedAt
		for _, l := range cn.Lines {
			lg := cnLineToGORM(&l)
			lg.CNID = cn.ID
			if err := tx.Create(lg).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *PGCreditNoteRepo) GetCN(ctx context.Context, id string) (*domain.CreditNote, error) {
	var g domain.CreditNoteGORM
	if err := r.db.WithContext(ctx).Preload("Lines").First(&g, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrCNNotFound
		}
		return nil, err
	}
	cn := cnFromGORM(&g)
	cn.Lines = make([]domain.CNLine, len(g.Lines))
	for i, l := range g.Lines {
		cn.Lines[i] = *cnLineFromGORM(&l)
	}
	return cn, nil
}

func (r *PGCreditNoteRepo) GetCNByNumber(ctx context.Context, companyID, cnNumber string) (*domain.CreditNote, error) {
	var g domain.CreditNoteGORM
	if err := r.db.WithContext(ctx).Preload("Lines").Where("company_id = ? AND cn_number = ?", companyID, cnNumber).First(&g).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrCNNotFound
		}
		return nil, err
	}
	cn := cnFromGORM(&g)
	cn.Lines = make([]domain.CNLine, len(g.Lines))
	for i, l := range g.Lines {
		cn.Lines[i] = *cnLineFromGORM(&l)
	}
	return cn, nil
}

func (r *PGCreditNoteRepo) ListCNs(ctx context.Context, filter domain.CreditNoteFilter) ([]domain.CreditNote, int, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&domain.CreditNoteGORM{}).Where("company_id = ?", filter.CompanyID)
	if filter.CustomerID != "" {
		q = q.Where("customer_id = ?", filter.CustomerID)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", string(filter.Status))
	}
	if filter.FromDate != "" {
		q = q.Where("cn_date >= ?", filter.FromDate)
	}
	if filter.ToDate != "" {
		q = q.Where("cn_date <= ?", filter.ToDate)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var gs []domain.CreditNoteGORM
	dq := r.db.WithContext(ctx).Where("company_id = ?", filter.CompanyID)
	if filter.CustomerID != "" {
		dq = dq.Where("customer_id = ?", filter.CustomerID)
	}
	if filter.Status != "" {
		dq = dq.Where("status = ?", string(filter.Status))
	}
	if filter.FromDate != "" {
		dq = dq.Where("cn_date >= ?", filter.FromDate)
	}
	if filter.ToDate != "" {
		dq = dq.Where("cn_date <= ?", filter.ToDate)
	}
	dq = dq.Order("cn_date DESC")
	if filter.Limit > 0 {
		dq = dq.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		dq = dq.Offset(filter.Offset)
	}
	if err := dq.Find(&gs).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.CreditNote, len(gs))
	for i := range gs {
		cn := cnFromGORM(&gs[i])
		r.db.WithContext(ctx).Model(&domain.CNLineGORM{}).Where("cn_id = ?", gs[i].ID).Find(&cn.Lines)
		out[i] = *cn
	}
	return out, int(total), nil
}

func (r *PGCreditNoteRepo) UpdateCN(ctx context.Context, cn *domain.CreditNote) error {
	return r.db.WithContext(ctx).Model(&domain.CreditNoteGORM{}).Where("id = ?", cn.ID).Updates(map[string]interface{}{
		"cn_number": cn.CNNumber,
		"cn_date":   cn.ReturnDate,
		"customer_id": cn.CustomerID,
		"reason":    cn.ReturnReason,
		"subtotal":  cn.Subtotal,
		"tax_amount": cn.TaxAmount,
		"grand_total": cn.TotalAmount,
		"status":    string(cn.Status),
	}).Error
}

func (r *PGCreditNoteRepo) UpdateCNStatus(ctx context.Context, id string, status domain.CNStatus) error {
	return r.db.WithContext(ctx).Model(&domain.CreditNoteGORM{}).Where("id = ?", id).Update("status", string(status)).Error
}

func (r *PGCreditNoteRepo) PostCN(ctx context.Context, id string, postedAt time.Time) error {
	return r.db.WithContext(ctx).Model(&domain.CreditNoteGORM{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status": string(domain.CNPosted),
	}).Error
}

func (r *PGCreditNoteRepo) SetCNGLPosted(ctx context.Context, id string, postedAt time.Time) error {
	return r.db.WithContext(ctx).Model(&domain.CreditNoteGORM{}).Where("id = ?", id).Updates(map[string]interface{}{
		"gl_posted": true,
		"posted_at": postedAt,
	}).Error
}

func cnLineToGORM(l *domain.CNLine) *domain.CNLineGORM {
	return &domain.CNLineGORM{
		LineNumber: 0,
		ItemCode:   l.ItemName,
		ItemName:   l.ItemName,
		Qty:        l.Quantity,
		UnitPrice:  l.UnitPrice,
		LineTotal:  l.LineTotal,
	}
}

func cnLineFromGORM(g *domain.CNLineGORM) *domain.CNLine {
	return &domain.CNLine{
		ID:       strconv.Itoa(int(g.ID)),
		CNID:     g.CNID,
		ItemName: g.ItemName,
		Unit:     "",
		Quantity: g.Qty,
		UnitPrice: g.UnitPrice,
		LineTotal: g.LineTotal,
		VATRate:  10,
	}
}

func (r *PGCreditNoteRepo) GetCNLines(ctx context.Context, cnID string) ([]domain.CNLine, error) {
	var gs []domain.CNLineGORM
	if err := r.db.WithContext(ctx).Where("cn_id = ?", cnID).Find(&gs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.CNLine, len(gs))
	for i := range gs {
		out[i] = *cnLineFromGORM(&gs[i])
	}
	return out, nil
}

func (r *PGCreditNoteRepo) CreateCNLines(ctx context.Context, items []domain.CNLine) error {
	if len(items) == 0 {
		return nil
	}
	gs := make([]domain.CNLineGORM, len(items))
	for i := range items {
		gs[i] = *cnLineToGORM(&items[i])
	}
	return r.db.WithContext(ctx).Create(&gs).Error
}

// ─── AR Transaction ─────────────────────────────────────────────────────

type PGARTransactionRepo struct {
	db *gorm.DB
}

func NewPGARTransactionRepo(db *gorm.DB) *PGARTransactionRepo {
	return &PGARTransactionRepo{db: db}
}

func artToGORM(t *domain.ARTransaction) *domain.ARTransactionGORM {
	return &domain.ARTransactionGORM{
		ID:         t.ID,
		CompanyID:  t.CompanyID,
		CustomerID: t.CustomerID,
		Amount:     t.Amount,
		Currency:   t.Currency,
		Status:     "OPEN",
		CreatedBy:  "",
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.CreatedAt,
	}
}

func artFromGORM(g *domain.ARTransactionGORM) *domain.ARTransaction {
	return &domain.ARTransaction{
		ID:              g.ID,
		CompanyID:       g.CompanyID,
		CustomerID:      g.CustomerID,
		TransactionType: domain.ARTransInvoice,
		Amount:          g.Amount,
		Currency:        g.Currency,
		CreatedAt:       g.CreatedAt,
		TransactionDate: g.CreatedAt,
	}
}

func (r *PGARTransactionRepo) CreateARTransaction(ctx context.Context, t *domain.ARTransaction) error {
	g := artToGORM(t)
	if err := r.db.WithContext(ctx).Create(g).Error; err != nil {
		return err
	}
	t.ID = g.ID
	t.CreatedAt = g.CreatedAt
	return nil
}

func (r *PGARTransactionRepo) GetARTransaction(ctx context.Context, id string) (*domain.ARTransaction, error) {
	var g domain.ARTransactionGORM
	if err := r.db.WithContext(ctx).First(&g, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrARTransNotFound
		}
		return nil, err
	}
	return artFromGORM(&g), nil
}

func (r *PGARTransactionRepo) ListARTransactions(ctx context.Context, companyID, customerID string) ([]domain.ARTransaction, error) {
	var gs []domain.ARTransactionGORM
	q := r.db.WithContext(ctx).Joins("JOIN customers ON customers.id = ar_transactions.customer_id").Where("ar_transactions.customer_id = ? AND customers.company_id = ?", customerID, companyID)
	if err := q.Order("ar_transactions.created_at DESC").Find(&gs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.ARTransaction, len(gs))
	for i := range gs {
		out[i] = *artFromGORM(&gs[i])
	}
	return out, nil
}

func (r *PGARTransactionRepo) ListARTransactionsAll(ctx context.Context, companyID string, offset, limit int) ([]domain.ARTransaction, int, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&domain.ARTransactionGORM{}).Joins("JOIN customers ON customers.id = ar_transactions.customer_id").Where("customers.company_id = ?", companyID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var gs []domain.ARTransactionGORM
	dq := r.db.WithContext(ctx).Joins("JOIN customers ON customers.id = ar_transactions.customer_id").Where("customers.company_id = ?", companyID).Order("ar_transactions.created_at DESC")
	if limit > 0 {
		dq = dq.Limit(limit)
	}
	if offset > 0 {
		dq = dq.Offset(offset)
	}
	if err := dq.Find(&gs).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.ARTransaction, len(gs))
	for i := range gs {
		out[i] = *artFromGORM(&gs[i])
	}
	return out, int(total), nil
}

// ─── Sales Quotation ────────────────────────────────────────────────────

type PGSalesQuotationRepo struct {
	db *gorm.DB
}

func NewPGSalesQuotationRepo(db *gorm.DB) *PGSalesQuotationRepo {
	return &PGSalesQuotationRepo{db: db}
}

func sqToGORM(sq *domain.SalesQuotation) *domain.SalesQuotationGORM {
	return &domain.SalesQuotationGORM{
		ID:          sq.ID,
		CompanyID:   sq.CompanyID,
		QuoteNumber: sq.QNNumber,
		QuoteDate:   time.Now(),
		CustomerID:  sq.CustomerID,
		ValidUntil:  *sq.ValidUntil,
		Subtotal:    sq.TotalAmount,
		TotalAmount: sq.TotalAmount,
		Status:      sq.Status,
		CreatedBy:   sq.CreatedBy,
		CreatedAt:   sq.CreatedAt,
		UpdatedAt:   sq.CreatedAt,
	}
}

func sqFromGORM(g *domain.SalesQuotationGORM) *domain.SalesQuotation {
	return &domain.SalesQuotation{
		ID:          g.ID,
		CompanyID:   g.CompanyID,
		QNNumber:    g.QuoteNumber,
		CustomerID:  g.CustomerID,
		ValidUntil:  &g.ValidUntil,
		Status:      g.Status,
		TotalAmount: g.TotalAmount,
		CreatedBy:   g.CreatedBy,
		CreatedAt:   g.CreatedAt,
	}
}

func (r *PGSalesQuotationRepo) CreateSQ(ctx context.Context, sq *domain.SalesQuotation) error {
	g := sqToGORM(sq)
	if err := r.db.WithContext(ctx).Create(g).Error; err != nil {
		return err
	}
	sq.ID = g.ID
	sq.CreatedAt = g.CreatedAt
	return nil
}

func (r *PGSalesQuotationRepo) GetSQ(ctx context.Context, id string) (*domain.SalesQuotation, error) {
	var g domain.SalesQuotationGORM
	if err := r.db.WithContext(ctx).First(&g, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return sqFromGORM(&g), nil
}

func (r *PGSalesQuotationRepo) ListSQs(ctx context.Context, companyID string) ([]domain.SalesQuotation, error) {
	var gs []domain.SalesQuotationGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("created_at DESC").Find(&gs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.SalesQuotation, len(gs))
	for i := range gs {
		out[i] = *sqFromGORM(&gs[i])
	}
	return out, nil
}

func (r *PGSalesQuotationRepo) UpdateSQ(ctx context.Context, sq *domain.SalesQuotation) error {
	g := sqToGORM(sq)
	return r.db.WithContext(ctx).Model(&domain.SalesQuotationGORM{}).Where("id = ?", sq.ID).Updates(map[string]interface{}{
		"quote_number": g.QuoteNumber,
		"customer_id":  g.CustomerID,
		"valid_until":  g.ValidUntil,
		"status":       g.Status,
		"total_amount": g.TotalAmount,
	}).Error
}
