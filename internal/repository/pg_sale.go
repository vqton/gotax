package repository

import (
	"context"
	"fmt"
	"gotax/internal/domain"
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
	return &domain.CustomerGORM{
		ID:                c.ID,
		CompanyID:         c.CompanyID,
		Code:              c.Code,
		Name:              c.Name,
		TaxCode:           c.TaxCode,
		Address:           c.Address,
		Phone:             c.Phone,
		Email:             c.Email,
		BankAccountName:   c.BankAccountName,
		BankAccountNumber: c.BankAccountNumber,
		BankName:          c.BankName,
		PaymentTerms:      string(c.PaymentTerms),
		CreditLimit:       c.CreditLimit,
		Currency:          c.Currency,
		CustomerType:      string(c.CustomerType),
		CustomerGroup:     string(c.CustomerGroup),
		PriceListID:       c.PriceListID,
		Status:            string(c.Status),
		Notes:             c.Notes,
		CreatedBy:         c.CreatedBy,
		CreatedAt:         c.CreatedAt,
		UpdatedAt:         c.UpdatedAt,
	}
}

func customerFromGORM(g *domain.CustomerGORM) *domain.Customer {
	return &domain.Customer{
		ID:                g.ID,
		CompanyID:         g.CompanyID,
		Code:              g.Code,
		Name:              g.Name,
		TaxCode:           g.TaxCode,
		Address:           g.Address,
		Phone:             g.Phone,
		Email:             g.Email,
		BankAccountName:   g.BankAccountName,
		BankAccountNumber: g.BankAccountNumber,
		BankName:          g.BankName,
		PaymentTerms:      domain.PaymentTerms(g.PaymentTerms),
		CreditLimit:       g.CreditLimit,
		Currency:          g.Currency,
		CustomerType:      domain.CustomerType(g.CustomerType),
		CustomerGroup:     domain.CustomerGroup(g.CustomerGroup),
		PriceListID:       g.PriceListID,
		Status:            domain.CustomerStatus(g.Status),
		Notes:             g.Notes,
		CreatedBy:         g.CreatedBy,
		CreatedAt:         g.CreatedAt,
		UpdatedAt:         g.UpdatedAt,
	}
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
	return r.db.WithContext(ctx).Model(&domain.CustomerGORM{}).Where("id = ?", c.ID).Select("*").Updates(g).Error
}

func (r *PGCustomerRepo) DeleteCustomer(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.CustomerGORM{}, "id = ?", id).Error
}

// ─── Sales Order ───────────────────────────────────────────────────────

type PGSaleOrderRepo struct {
	db *gorm.DB
}

func NewPGSaleOrderRepo(db *gorm.DB) *PGSaleOrderRepo {
	return &PGSaleOrderRepo{db: db}
}

func soToGORM(so *domain.SalesOrder) *domain.SalesOrderGORM {
	g := &domain.SalesOrderGORM{
		ID:              so.ID,
		CompanyID:       so.CompanyID,
		SONumber:        so.SONumber,
		QuotationID:     so.QuotationID,
		CustomerID:      so.CustomerID,
		OrderDate:       so.OrderDate,
		ExpectedDate:    so.ExpectedDate,
		Currency:        so.Currency,
		ExchangeRate:    so.ExchangeRate,
		PaymentTerms:    so.PaymentTerms,
		DeliveryTerms:   so.DeliveryTerms,
		ShippingAddress: so.ShippingAddress,
		Subtotal:        so.Subtotal,
		DiscountAmount:  so.DiscountAmount,
		TaxAmount:       so.TaxAmount,
		TotalAmount:     so.TotalAmount,
		Status:          string(so.Status),
		ApprovedBy:      so.ApprovedBy,
		ApprovedAt:      so.ApprovedAt,
		CancelledReason: so.CancelledReason,
		Notes:           so.Notes,
		CreatedBy:       so.CreatedBy,
		CreatedAt:       so.CreatedAt,
		UpdatedAt:       so.UpdatedAt,
	}
	if len(so.Lines) > 0 {
		g.Lines = make([]domain.SOLineGORM, len(so.Lines))
		for i := range so.Lines {
			g.Lines[i] = *soLineToGORM(&so.Lines[i])
		}
	}
	return g
}

func soFromGORM(g *domain.SalesOrderGORM) *domain.SalesOrder {
	so := &domain.SalesOrder{
		ID:              g.ID,
		CompanyID:       g.CompanyID,
		SONumber:        g.SONumber,
		QuotationID:     g.QuotationID,
		CustomerID:      g.CustomerID,
		OrderDate:       g.OrderDate,
		ExpectedDate:    g.ExpectedDate,
		Currency:        g.Currency,
		ExchangeRate:    g.ExchangeRate,
		PaymentTerms:    g.PaymentTerms,
		DeliveryTerms:   g.DeliveryTerms,
		ShippingAddress: g.ShippingAddress,
		Subtotal:        g.Subtotal,
		DiscountAmount:  g.DiscountAmount,
		TaxAmount:       g.TaxAmount,
		TotalAmount:     g.TotalAmount,
		Status:          domain.SOStatus(g.Status),
		ApprovedBy:      g.ApprovedBy,
		ApprovedAt:      g.ApprovedAt,
		CancelledReason: g.CancelledReason,
		Notes:           g.Notes,
		CreatedBy:       g.CreatedBy,
		CreatedAt:       g.CreatedAt,
		UpdatedAt:       g.UpdatedAt,
	}
	if len(g.Lines) > 0 {
		so.Lines = make([]domain.SOLine, len(g.Lines))
		for i := range g.Lines {
			so.Lines[i] = *soLineFromGORM(&g.Lines[i])
		}
	}
	return so
}

func soLineToGORM(l *domain.SOLine) *domain.SOLineGORM {
	return &domain.SOLineGORM{
		ID:             l.ID,
		SOID:           l.SOID,
		LineNumber:     l.LineNumber,
		ItemCode:       l.ItemCode,
		ItemName:       l.ItemName,
		Unit:           l.Unit,
		Quantity:       l.Quantity,
		UnitPrice:      l.UnitPrice,
		DiscountPct:    l.DiscountPct,
		VATRate:        l.VATRate,
		VATType:        string(l.VATType),
		RevenueAccount: l.RevenueAccount,
		VATAccountID:   l.VATAccountID,
		LineTotal:      l.LineTotal,
		LineVATAmount:  l.LineVATAmount,
		DeliveredQty:   l.DeliveredQty,
		InvoicedQty:    l.InvoicedQty,
	}
}

func soLineFromGORM(g *domain.SOLineGORM) *domain.SOLine {
	return &domain.SOLine{
		ID:             g.ID,
		SOID:           g.SOID,
		LineNumber:     g.LineNumber,
		ItemCode:       g.ItemCode,
		ItemName:       g.ItemName,
		Unit:           g.Unit,
		Quantity:       g.Quantity,
		UnitPrice:      g.UnitPrice,
		DiscountPct:    g.DiscountPct,
		VATRate:        g.VATRate,
		VATType:        domain.VATType(g.VATType),
		RevenueAccount: g.RevenueAccount,
		VATAccountID:   g.VATAccountID,
		LineTotal:      g.LineTotal,
		LineVATAmount:  g.LineVATAmount,
		DeliveredQty:   g.DeliveredQty,
		InvoicedQty:    g.InvoicedQty,
	}
}

func (r *PGSaleOrderRepo) CreateSO(ctx context.Context, so *domain.SalesOrder) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		g := soToGORM(so)
		g.Lines = nil
		if err := tx.Create(g).Error; err != nil {
			return err
		}
		so.ID = g.ID
		so.CreatedAt = g.CreatedAt
		so.UpdatedAt = g.UpdatedAt
		for i, l := range so.Lines {
			lg := soLineToGORM(&l)
			lg.ID = ""
			lg.SOID = so.ID
			lg.LineNumber = i + 1
			if err := tx.Create(lg).Error; err != nil {
				return err
			}
			so.Lines[i].ID = lg.ID
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
	for i := range g.Lines {
		so.Lines[i] = *soLineFromGORM(&g.Lines[i])
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
	for i := range g.Lines {
		so.Lines[i] = *soLineFromGORM(&g.Lines[i])
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
	var gs []domain.SalesOrderGORM
	if err := dq.Find(&gs).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.SalesOrder, len(gs))
	for i := range gs {
		so := soFromGORM(&gs[i])
		var glines []domain.SOLineGORM
		r.db.WithContext(ctx).Model(&domain.SOLineGORM{}).Where("so_id = ?", gs[i].ID).Order("line_number").Find(&glines)
		so.Lines = make([]domain.SOLine, len(glines))
		for j := range glines {
			so.Lines[j] = *soLineFromGORM(&glines[j])
		}
		out[i] = *so
	}
	return out, int(total), nil
}

func (r *PGSaleOrderRepo) UpdateSO(ctx context.Context, so *domain.SalesOrder) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		g := soToGORM(so)
		if err := tx.Model(&domain.SalesOrderGORM{}).Where("id = ?", so.ID).Select("*").Updates(g).Error; err != nil {
			return err
		}
		if err := tx.Where("so_id = ?", so.ID).Delete(&domain.SOLineGORM{}).Error; err != nil {
			return err
		}
		for i, l := range so.Lines {
			lg := soLineToGORM(&l)
			lg.ID = ""
			lg.SOID = so.ID
			lg.LineNumber = i + 1
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
		"status":           string(domain.SOCancelled),
		"cancelled_reason": cancelReason,
	}).Error
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
		gs[i].LineNumber = i + 1
	}
	return r.db.WithContext(ctx).Create(&gs).Error
}

func (r *PGSaleOrderRepo) UpdateSOLines(ctx context.Context, items []domain.SOLine) error {
	for _, l := range items {
		g := soLineToGORM(&l)
		if err := r.db.WithContext(ctx).Model(&domain.SOLineGORM{}).Where("id = ?", l.ID).Select("*").Updates(g).Error; err != nil {
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

// ─── Delivery Note ─────────────────────────────────────────────────────

type PGDeliveryNoteRepo struct {
	db *gorm.DB
}

func NewPGDeliveryNoteRepo(db *gorm.DB) *PGDeliveryNoteRepo {
	return &PGDeliveryNoteRepo{db: db}
}

func dnToGORM(dn *domain.DeliveryNote) *domain.DeliveryNoteGORM {
	g := &domain.DeliveryNoteGORM{
		ID:              dn.ID,
		CompanyID:       dn.CompanyID,
		DNNumber:        dn.DNNumber,
		SOID:            dn.SOID,
		DeliveryDate:    dn.DeliveryDate,
		Warehouse:       dn.Warehouse,
		ShippingMethod:  dn.ShippingMethod,
		CarrierName:     dn.CarrierName,
		TrackingNumber:  dn.TrackingNumber,
		DeliveryAddress: dn.DeliveryAddress,
		Status:          string(dn.Status),
		Notes:           dn.Notes,
		CreatedBy:       dn.CreatedBy,
		CreatedAt:       dn.CreatedAt,
		UpdatedAt:       dn.UpdatedAt,
	}
	if len(dn.Lines) > 0 {
		g.Lines = make([]domain.DNLineGORM, len(dn.Lines))
		for i := range dn.Lines {
			g.Lines[i] = *dnLineToGORM(&dn.Lines[i])
		}
	}
	return g
}

func dnFromGORM(g *domain.DeliveryNoteGORM) *domain.DeliveryNote {
	dn := &domain.DeliveryNote{
		ID:              g.ID,
		CompanyID:       g.CompanyID,
		DNNumber:        g.DNNumber,
		SOID:            g.SOID,
		DeliveryDate:    g.DeliveryDate,
		Warehouse:       g.Warehouse,
		ShippingMethod:  g.ShippingMethod,
		CarrierName:     g.CarrierName,
		TrackingNumber:  g.TrackingNumber,
		DeliveryAddress: g.DeliveryAddress,
		Status:          domain.DNStatus(g.Status),
		Notes:           g.Notes,
		CreatedBy:       g.CreatedBy,
		CreatedAt:       g.CreatedAt,
		UpdatedAt:       g.UpdatedAt,
	}
	if len(g.Lines) > 0 {
		dn.Lines = make([]domain.DNLine, len(g.Lines))
		for i := range g.Lines {
			dn.Lines[i] = *dnLineFromGORM(&g.Lines[i])
		}
	}
	return dn
}

func dnLineToGORM(l *domain.DNLine) *domain.DNLineGORM {
	return &domain.DNLineGORM{
		ID:           l.ID,
		DNID:         l.DNID,
		SOLineID:     l.SOLineID,
		ItemCode:     l.ItemCode,
		ItemName:     l.ItemName,
		Unit:         l.Unit,
		QtyDelivered: l.QtyDelivered,
		QtyReturned:  l.QtyReturned,
		UnitPrice:    l.UnitPrice,
		LineTotal:    l.LineTotal,
		CostPrice:    l.CostPrice,
	}
}

func dnLineFromGORM(g *domain.DNLineGORM) *domain.DNLine {
	return &domain.DNLine{
		ID:           g.ID,
		DNID:         g.DNID,
		SOLineID:     g.SOLineID,
		ItemCode:     g.ItemCode,
		ItemName:     g.ItemName,
		Unit:         g.Unit,
		QtyDelivered: g.QtyDelivered,
		QtyReturned:  g.QtyReturned,
		UnitPrice:    g.UnitPrice,
		LineTotal:    g.LineTotal,
		CostPrice:    g.CostPrice,
	}
}

func (r *PGDeliveryNoteRepo) CreateDN(ctx context.Context, dn *domain.DeliveryNote) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		g := dnToGORM(dn)
		g.Lines = nil
		if err := tx.Create(g).Error; err != nil {
			return err
		}
		dn.ID = g.ID
		dn.CreatedAt = g.CreatedAt
		dn.UpdatedAt = g.UpdatedAt
		for i, l := range dn.Lines {
			lg := dnLineToGORM(&l)
			lg.ID = ""
			lg.DNID = dn.ID
			if err := tx.Create(lg).Error; err != nil {
				return err
			}
			dn.Lines[i].ID = lg.ID
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
	for i := range g.Lines {
		dn.Lines[i] = *dnLineFromGORM(&g.Lines[i])
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
	for i := range g.Lines {
		dn.Lines[i] = *dnLineFromGORM(&g.Lines[i])
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
		q = q.Where("delivery_date >= ?", filter.FromDate)
	}
	if filter.ToDate != "" {
		q = q.Where("delivery_date <= ?", filter.ToDate)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	dq := r.db.WithContext(ctx).Where("company_id = ?", filter.CompanyID)
	if filter.SOID != "" {
		dq = dq.Where("so_id = ?", filter.SOID)
	}
	if filter.Status != "" {
		dq = dq.Where("status = ?", string(filter.Status))
	}
	if filter.FromDate != "" {
		dq = dq.Where("delivery_date >= ?", filter.FromDate)
	}
	if filter.ToDate != "" {
		dq = dq.Where("delivery_date <= ?", filter.ToDate)
	}
	dq = dq.Order("delivery_date DESC")
	if filter.Limit > 0 {
		dq = dq.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		dq = dq.Offset(filter.Offset)
	}
	var gs []domain.DeliveryNoteGORM
	if err := dq.Find(&gs).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.DeliveryNote, len(gs))
	for i := range gs {
		dn := dnFromGORM(&gs[i])
		var glines []domain.DNLineGORM
		r.db.WithContext(ctx).Model(&domain.DNLineGORM{}).Where("dn_id = ?", gs[i].ID).Find(&glines)
		dn.Lines = make([]domain.DNLine, len(glines))
		for j := range glines {
			dn.Lines[j] = *dnLineFromGORM(&glines[j])
		}
		out[i] = *dn
	}
	return out, int(total), nil
}

func (r *PGDeliveryNoteRepo) UpdateDN(ctx context.Context, dn *domain.DeliveryNote) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		g := dnToGORM(dn)
		if err := tx.Model(&domain.DeliveryNoteGORM{}).Where("id = ?", dn.ID).Select("*").Updates(g).Error; err != nil {
			return err
		}
		if err := tx.Where("dn_id = ?", dn.ID).Delete(&domain.DNLineGORM{}).Error; err != nil {
			return err
		}
		for _, l := range dn.Lines {
			lg := dnLineToGORM(&l)
			lg.ID = ""
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
		if err := r.db.WithContext(ctx).Model(&domain.DNLineGORM{}).Where("id = ?", l.ID).Select("*").Updates(g).Error; err != nil {
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
		ID:                 inv.ID,
		CompanyID:          inv.CompanyID,
		InvoiceNumber:      inv.InvoiceNumber,
		InvoiceDate:        inv.InvoiceDate,
		SOID:               inv.SOID,
		DNID:               inv.DNID,
		CustomerID:         inv.CustomerID,
		CustomerName:       inv.CustomerName,
		CustomerTaxCode:    inv.CustomerTaxCode,
		CustomerAddress:    inv.CustomerAddress,
		InvoiceType:        inv.InvoiceType,
		Currency:           inv.Currency,
		ExchangeRate:       inv.ExchangeRate,
		Subtotal:           inv.Subtotal,
		DiscountAmount:     inv.DiscountAmount,
		TaxAmount:          inv.TaxAmount,
		TotalAmount:        inv.TotalAmount,
		AmountReceived:     inv.AmountReceived,
		BalanceDue:         inv.BalanceDue,
		DueDate:            inv.DueDate,
		InvoiceNote:        inv.InvoiceNote,
		EInvoiceData:       inv.EInvoiceData,
		EInvoiceCode:       inv.EInvoiceCode,
		EInvStatus:         string(inv.EInvStatus),
		DigitalSignatureID: inv.DigitalSignatureID,
		SignedData:         inv.SignedData,
		GDTResponse:        inv.GDTResponse,
		OriginalInvoiceID:  inv.OriginalInvoiceID,
		AdjustmentType:     string(inv.AdjustmentType),
		Status:             string(inv.Status),
		GLPosted:           inv.GLPosted,
		GLPostedAt:         inv.GLPostedAt,
		Notes:              inv.Notes,
		CreatedBy:          inv.CreatedBy,
		CreatedAt:          inv.CreatedAt,
		UpdatedAt:          inv.UpdatedAt,
	}
	if len(inv.Lines) > 0 {
		g.Lines = make([]domain.InvLineGORM, len(inv.Lines))
		for i := range inv.Lines {
			g.Lines[i] = *invLineToGORM(&inv.Lines[i])
		}
	}
	return g
}

func cinvFromGORM(g *domain.CustomerInvoiceGORM) *domain.CustomerInvoice {
	inv := &domain.CustomerInvoice{
		ID:                 g.ID,
		CompanyID:          g.CompanyID,
		InvoiceNumber:      g.InvoiceNumber,
		InvoiceDate:        g.InvoiceDate,
		SOID:               g.SOID,
		DNID:               g.DNID,
		CustomerID:         g.CustomerID,
		CustomerName:       g.CustomerName,
		CustomerTaxCode:    g.CustomerTaxCode,
		CustomerAddress:    g.CustomerAddress,
		InvoiceType:        g.InvoiceType,
		Currency:           g.Currency,
		ExchangeRate:       g.ExchangeRate,
		Subtotal:           g.Subtotal,
		DiscountAmount:     g.DiscountAmount,
		TaxAmount:          g.TaxAmount,
		TotalAmount:        g.TotalAmount,
		AmountReceived:     g.AmountReceived,
		BalanceDue:         g.BalanceDue,
		DueDate:            g.DueDate,
		InvoiceNote:        g.InvoiceNote,
		EInvoiceData:       g.EInvoiceData,
		EInvoiceCode:       g.EInvoiceCode,
		EInvStatus:         domain.InvEStatus(g.EInvStatus),
		DigitalSignatureID: g.DigitalSignatureID,
		SignedData:         g.SignedData,
		GDTResponse:        g.GDTResponse,
		OriginalInvoiceID:  g.OriginalInvoiceID,
		AdjustmentType:     domain.SaleAdjustmentType(g.AdjustmentType),
		Status:             domain.SaleInvoiceStatus(g.Status),
		GLPosted:           g.GLPosted,
		GLPostedAt:         g.GLPostedAt,
		Notes:              g.Notes,
		CreatedBy:          g.CreatedBy,
		CreatedAt:          g.CreatedAt,
		UpdatedAt:          g.UpdatedAt,
	}
	if len(g.Lines) > 0 {
		inv.Lines = make([]domain.InvLine, len(g.Lines))
		for i := range g.Lines {
			inv.Lines[i] = *invLineFromGORM(&g.Lines[i])
		}
	}
	return inv
}

func invLineToGORM(l *domain.InvLine) *domain.InvLineGORM {
	return &domain.InvLineGORM{
		ID:             l.ID,
		InvoiceID:      l.InvoiceID,
		SOLineID:       l.SOLineID,
		DNLineID:       l.DNLineID,
		ItemCode:       l.ItemCode,
		ItemName:       l.ItemName,
		Unit:           l.Unit,
		Quantity:       l.Quantity,
		UnitPrice:      l.UnitPrice,
		DiscountPct:    l.DiscountPct,
		VATRate:        l.VATRate,
		VATType:        string(l.VATType),
		LineTotal:      l.LineTotal,
		LineVATAmount:  l.LineVATAmount,
		RevenueAccount: l.RevenueAccount,
		VATAccountID:   l.VATAccountID,
	}
}

func invLineFromGORM(g *domain.InvLineGORM) *domain.InvLine {
	return &domain.InvLine{
		ID:             g.ID,
		InvoiceID:      g.InvoiceID,
		SOLineID:       g.SOLineID,
		DNLineID:       g.DNLineID,
		ItemCode:       g.ItemCode,
		ItemName:       g.ItemName,
		Unit:           g.Unit,
		Quantity:       g.Quantity,
		UnitPrice:      g.UnitPrice,
		DiscountPct:    g.DiscountPct,
		VATRate:        g.VATRate,
		VATType:        domain.VATType(g.VATType),
		LineTotal:      g.LineTotal,
		LineVATAmount:  g.LineVATAmount,
		RevenueAccount: g.RevenueAccount,
		VATAccountID:   g.VATAccountID,
	}
}

func (r *PGCustomerInvoiceRepo) CreateInvoice(ctx context.Context, inv *domain.CustomerInvoice) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		g := cinvToGORM(inv)
		g.Lines = nil
		if err := tx.Create(g).Error; err != nil {
			return err
		}
		inv.ID = g.ID
		inv.CreatedAt = g.CreatedAt
		inv.UpdatedAt = g.UpdatedAt
		for i, l := range inv.Lines {
			lg := invLineToGORM(&l)
			lg.ID = ""
			lg.InvoiceID = inv.ID
			if err := tx.Create(lg).Error; err != nil {
				return err
			}
			inv.Lines[i].ID = lg.ID
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
	for i := range g.Lines {
		inv.Lines[i] = *invLineFromGORM(&g.Lines[i])
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
	for i := range g.Lines {
		inv.Lines[i] = *invLineFromGORM(&g.Lines[i])
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
	var gs []domain.CustomerInvoiceGORM
	if err := dq.Find(&gs).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.CustomerInvoice, len(gs))
	for i := range gs {
		inv := cinvFromGORM(&gs[i])
		var glines []domain.InvLineGORM
		r.db.WithContext(ctx).Model(&domain.InvLineGORM{}).Where("invoice_id = ?", gs[i].ID).Find(&glines)
		inv.Lines = make([]domain.InvLine, len(glines))
		for j := range glines {
			inv.Lines[j] = *invLineFromGORM(&glines[j])
		}
		out[i] = *inv
	}
	return out, int(total), nil
}

func (r *PGCustomerInvoiceRepo) UpdateInvoice(ctx context.Context, inv *domain.CustomerInvoice) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		g := cinvToGORM(inv)
		if err := tx.Model(&domain.CustomerInvoiceGORM{}).Where("id = ?", inv.ID).Select("*").Updates(g).Error; err != nil {
			return err
		}
		if err := tx.Where("invoice_id = ?", inv.ID).Delete(&domain.InvLineGORM{}).Error; err != nil {
			return err
		}
		for _, l := range inv.Lines {
			lg := invLineToGORM(&l)
			lg.ID = ""
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
		"status":    string(domain.SInvPosted),
		"posted_at": postedAt,
	}).Error
}

func (r *PGCustomerInvoiceRepo) SetInvoiceGLPosted(ctx context.Context, id string, postedAt time.Time) error {
	return r.db.WithContext(ctx).Model(&domain.CustomerInvoiceGORM{}).Where("id = ?", id).Updates(map[string]interface{}{
		"gl_posted":    true,
		"gl_posted_at": postedAt,
	}).Error
}

func (r *PGCustomerInvoiceRepo) AllocateToInvoice(ctx context.Context, invoiceID string, amount float64) error {
	return r.db.WithContext(ctx).Model(&domain.CustomerInvoiceGORM{}).Where("id = ?", invoiceID).Updates(map[string]interface{}{
		"amount_received": gorm.Expr("amount_received + ?", amount),
		"balance_due":     gorm.Expr("GREATEST(balance_due - ?, 0)", amount),
	}).Error
}

func (r *PGCustomerInvoiceRepo) ReduceInvoiceBalance(ctx context.Context, invoiceID string, amount float64) error {
	return r.db.WithContext(ctx).Model(&domain.CustomerInvoiceGORM{}).Where("id = ?", invoiceID).Updates(map[string]interface{}{
		"balance_due": gorm.Expr("GREATEST(balance_due - ?, 0)", amount),
	}).Error
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
		if err := r.db.WithContext(ctx).Model(&domain.InvLineGORM{}).Where("id = ?", l.ID).Select("*").Updates(g).Error; err != nil {
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
	g := &domain.CustomerReceiptGORM{
		ID:                r.ID,
		CompanyID:         r.CompanyID,
		ReceiptNumber:     r.ReceiptNumber,
		CustomerID:        r.CustomerID,
		ReceiptDate:       r.ReceiptDate,
		PaymentMethod:     r.PaymentMethod,
		BankAccountID:     r.BankAccountID,
		Currency:          r.Currency,
		ExchangeRate:      r.ExchangeRate,
		Amount:            r.Amount,
		UnallocatedAmount: r.UnallocatedAmount,
		Reference:         r.Reference,
		Notes:             r.Notes,
		Status:            string(r.Status),
		GLPosted:          r.GLPosted,
		GLPostedAt:        r.GLPostedAt,
		CreatedBy:         r.CreatedBy,
		CreatedAt:         r.CreatedAt,
		UpdatedAt:         r.UpdatedAt,
	}
	if len(r.Allocations) > 0 {
		g.Allocations = make([]domain.RcpAllocationGORM, len(r.Allocations))
		for i := range r.Allocations {
			g.Allocations[i] = *rcpAllocToGORM(&r.Allocations[i])
		}
	}
	return g
}

func rcptFromGORM(g *domain.CustomerReceiptGORM) *domain.CustomerReceipt {
	rcpt := &domain.CustomerReceipt{
		ID:                g.ID,
		CompanyID:         g.CompanyID,
		ReceiptNumber:     g.ReceiptNumber,
		CustomerID:        g.CustomerID,
		ReceiptDate:       g.ReceiptDate,
		PaymentMethod:     g.PaymentMethod,
		BankAccountID:     g.BankAccountID,
		Currency:          g.Currency,
		ExchangeRate:      g.ExchangeRate,
		Amount:            g.Amount,
		UnallocatedAmount: g.UnallocatedAmount,
		Reference:         g.Reference,
		Notes:             g.Notes,
		Status:            domain.ReceiptStatus(g.Status),
		GLPosted:          g.GLPosted,
		GLPostedAt:        g.GLPostedAt,
		CreatedBy:         g.CreatedBy,
		CreatedAt:         g.CreatedAt,
		UpdatedAt:         g.UpdatedAt,
	}
	if len(g.Allocations) > 0 {
		rcpt.Allocations = make([]domain.RcpAllocation, len(g.Allocations))
		for i := range g.Allocations {
			rcpt.Allocations[i] = *rcpAllocFromGORM(&g.Allocations[i])
		}
	}
	return rcpt
}

func rcpAllocToGORM(a *domain.RcpAllocation) *domain.RcpAllocationGORM {
	return &domain.RcpAllocationGORM{
		ID:              a.ID,
		ReceiptID:       a.ReceiptID,
		InvoiceID:       a.InvoiceID,
		AllocatedAmount: a.AllocatedAmount,
		DiscountAmount:  a.DiscountAmount,
	}
}

func rcpAllocFromGORM(g *domain.RcpAllocationGORM) *domain.RcpAllocation {
	return &domain.RcpAllocation{
		ID:              g.ID,
		ReceiptID:       g.ReceiptID,
		InvoiceID:       g.InvoiceID,
		AllocatedAmount: g.AllocatedAmount,
		DiscountAmount:  g.DiscountAmount,
	}
}

func (r *PGCustomerReceiptRepo) CreateReceipt(ctx context.Context, rcpt *domain.CustomerReceipt) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		g := rcptToGORM(rcpt)
		g.Allocations = nil // avoid GORM cascade
		if err := tx.Create(g).Error; err != nil {
			return err
		}
		rcpt.ID = g.ID
		rcpt.CreatedAt = g.CreatedAt
		rcpt.UpdatedAt = g.UpdatedAt
		for i, a := range rcpt.Allocations {
			ag := rcpAllocToGORM(&a)
			ag.ID = ""
			ag.ReceiptID = rcpt.ID
			if err := tx.Create(ag).Error; err != nil {
				return err
			}
			rcpt.Allocations[i].ID = ag.ID
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
	for i := range g.Allocations {
		rcpt.Allocations[i] = *rcpAllocFromGORM(&g.Allocations[i])
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
	for i := range g.Allocations {
		rcpt.Allocations[i] = *rcpAllocFromGORM(&g.Allocations[i])
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
	var gs []domain.CustomerReceiptGORM
	if err := dq.Find(&gs).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.CustomerReceipt, len(gs))
	for i := range gs {
		rcpt := rcptFromGORM(&gs[i])
		var galls []domain.RcpAllocationGORM
		r.db.WithContext(ctx).Model(&domain.RcpAllocationGORM{}).Where("receipt_id = ?", gs[i].ID).Find(&galls)
		rcpt.Allocations = make([]domain.RcpAllocation, len(galls))
		for j := range galls {
			rcpt.Allocations[j] = *rcpAllocFromGORM(&galls[j])
		}
		out[i] = *rcpt
	}
	return out, int(total), nil
}

func (r *PGCustomerReceiptRepo) UpdateReceipt(ctx context.Context, rcpt *domain.CustomerReceipt) error {
	g := rcptToGORM(rcpt)
	return r.db.WithContext(ctx).Model(&domain.CustomerReceiptGORM{}).Where("id = ?", rcpt.ID).Select("*").Updates(g).Error
}

func (r *PGCustomerReceiptRepo) UpdateReceiptStatus(ctx context.Context, id string, status domain.ReceiptStatus) error {
	return r.db.WithContext(ctx).Model(&domain.CustomerReceiptGORM{}).Where("id = ?", id).Update("status", string(status)).Error
}

func (r *PGCustomerReceiptRepo) SetReceiptGLPosted(ctx context.Context, id string, postedAt time.Time) error {
	return r.db.WithContext(ctx).Model(&domain.CustomerReceiptGORM{}).Where("id = ?", id).Updates(map[string]interface{}{
		"gl_posted":    true,
		"gl_posted_at": postedAt,
	}).Error
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
		ID:                cn.ID,
		CompanyID:         cn.CompanyID,
		CNNumber:          cn.CNNumber,
		OriginalInvoiceID: cn.OriginalInvoiceID,
		CustomerID:        cn.CustomerID,
		ReturnDate:        cn.ReturnDate,
		ReturnReason:      cn.ReturnReason,
		ReturnType:        string(cn.ReturnType),
		DNID:              cn.DNID,
		Subtotal:          cn.Subtotal,
		TaxAmount:         cn.TaxAmount,
		TotalAmount:       cn.TotalAmount,
		EInvoiceData:      cn.EInvoiceData,
		EInvoiceCode:      cn.EInvoiceCode,
		Status:            string(cn.Status),
		GLPosted:          cn.GLPosted,
		GLPostedAt:        cn.GLPostedAt,
		Notes:             cn.Notes,
		CreatedBy:         cn.CreatedBy,
		CreatedAt:         cn.CreatedAt,
		UpdatedAt:         cn.UpdatedAt,
	}
	if len(cn.Lines) > 0 {
		g.Lines = make([]domain.CNLineGORM, len(cn.Lines))
		for i := range cn.Lines {
			g.Lines[i] = *cnLineToGORM(&cn.Lines[i])
		}
	}
	return g
}

func cnFromGORM(g *domain.CreditNoteGORM) *domain.CreditNote {
	cn := &domain.CreditNote{
		ID:                g.ID,
		CompanyID:         g.CompanyID,
		CNNumber:          g.CNNumber,
		OriginalInvoiceID: g.OriginalInvoiceID,
		CustomerID:        g.CustomerID,
		ReturnDate:        g.ReturnDate,
		ReturnReason:      g.ReturnReason,
		ReturnType:        domain.ReturnType(g.ReturnType),
		DNID:              g.DNID,
		Subtotal:          g.Subtotal,
		TaxAmount:         g.TaxAmount,
		TotalAmount:       g.TotalAmount,
		EInvoiceData:      g.EInvoiceData,
		EInvoiceCode:      g.EInvoiceCode,
		Status:            domain.CNStatus(g.Status),
		GLPosted:          g.GLPosted,
		GLPostedAt:        g.GLPostedAt,
		Notes:             g.Notes,
		CreatedBy:         g.CreatedBy,
		CreatedAt:         g.CreatedAt,
		UpdatedAt:         g.UpdatedAt,
	}
	if len(g.Lines) > 0 {
		cn.Lines = make([]domain.CNLine, len(g.Lines))
		for i := range g.Lines {
			cn.Lines[i] = *cnLineFromGORM(&g.Lines[i])
		}
	}
	return cn
}

func cnLineToGORM(l *domain.CNLine) *domain.CNLineGORM {
	return &domain.CNLineGORM{
		ID:            l.ID,
		CNID:          l.CNID,
		InvLineID:     l.InvLineID,
		ItemName:      l.ItemName,
		Unit:          l.Unit,
		Quantity:      l.Quantity,
		UnitPrice:     l.UnitPrice,
		VATRate:       l.VATRate,
		LineTotal:     l.LineTotal,
		LineVATAmount: l.LineVATAmount,
	}
}

func cnLineFromGORM(g *domain.CNLineGORM) *domain.CNLine {
	return &domain.CNLine{
		ID:            g.ID,
		CNID:          g.CNID,
		InvLineID:     g.InvLineID,
		ItemName:      g.ItemName,
		Unit:          g.Unit,
		Quantity:      g.Quantity,
		UnitPrice:     g.UnitPrice,
		VATRate:       g.VATRate,
		LineTotal:     g.LineTotal,
		LineVATAmount: g.LineVATAmount,
	}
}

func (r *PGCreditNoteRepo) CreateCN(ctx context.Context, cn *domain.CreditNote) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		g := cnToGORM(cn)
		g.Lines = nil
		if err := tx.Create(g).Error; err != nil {
			return err
		}
		cn.ID = g.ID
		cn.CreatedAt = g.CreatedAt
		cn.UpdatedAt = g.UpdatedAt
		for i, l := range cn.Lines {
			lg := cnLineToGORM(&l)
			lg.ID = ""
			lg.CNID = cn.ID
			if err := tx.Create(lg).Error; err != nil {
				return err
			}
			cn.Lines[i].ID = lg.ID
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
	for i := range g.Lines {
		cn.Lines[i] = *cnLineFromGORM(&g.Lines[i])
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
	for i := range g.Lines {
		cn.Lines[i] = *cnLineFromGORM(&g.Lines[i])
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
		q = q.Where("return_date >= ?", filter.FromDate)
	}
	if filter.ToDate != "" {
		q = q.Where("return_date <= ?", filter.ToDate)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	dq := r.db.WithContext(ctx).Where("company_id = ?", filter.CompanyID)
	if filter.CustomerID != "" {
		dq = dq.Where("customer_id = ?", filter.CustomerID)
	}
	if filter.Status != "" {
		dq = dq.Where("status = ?", string(filter.Status))
	}
	if filter.FromDate != "" {
		dq = dq.Where("return_date >= ?", filter.FromDate)
	}
	if filter.ToDate != "" {
		dq = dq.Where("return_date <= ?", filter.ToDate)
	}
	dq = dq.Order("return_date DESC")
	if filter.Limit > 0 {
		dq = dq.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		dq = dq.Offset(filter.Offset)
	}
	var gs []domain.CreditNoteGORM
	if err := dq.Find(&gs).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.CreditNote, len(gs))
	for i := range gs {
		cn := cnFromGORM(&gs[i])
		var glines []domain.CNLineGORM
		r.db.WithContext(ctx).Model(&domain.CNLineGORM{}).Where("cn_id = ?", gs[i].ID).Find(&glines)
		cn.Lines = make([]domain.CNLine, len(glines))
		for j := range glines {
			cn.Lines[j] = *cnLineFromGORM(&glines[j])
		}
		out[i] = *cn
	}
	return out, int(total), nil
}

func (r *PGCreditNoteRepo) UpdateCN(ctx context.Context, cn *domain.CreditNote) error {
	g := cnToGORM(cn)
	return r.db.WithContext(ctx).Model(&domain.CreditNoteGORM{}).Where("id = ?", cn.ID).Select("*").Updates(g).Error
}

func (r *PGCreditNoteRepo) UpdateCNStatus(ctx context.Context, id string, status domain.CNStatus) error {
	return r.db.WithContext(ctx).Model(&domain.CreditNoteGORM{}).Where("id = ?", id).Update("status", string(status)).Error
}

func (r *PGCreditNoteRepo) PostCN(ctx context.Context, id string, postedAt time.Time) error {
	return r.db.WithContext(ctx).Model(&domain.CreditNoteGORM{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":    string(domain.CNPosted),
		"posted_at": postedAt,
	}).Error
}

func (r *PGCreditNoteRepo) SetCNGLPosted(ctx context.Context, id string, postedAt time.Time) error {
	return r.db.WithContext(ctx).Model(&domain.CreditNoteGORM{}).Where("id = ?", id).Updates(map[string]interface{}{
		"gl_posted":    true,
		"gl_posted_at": postedAt,
	}).Error
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
		ID:              t.ID,
		CompanyID:       t.CompanyID,
		CustomerID:      t.CustomerID,
		InvoiceID:       t.InvoiceID,
		TransactionType: string(t.TransactionType),
		TransactionDate: t.TransactionDate,
		Amount:          t.Amount,
		Currency:        t.Currency,
		ReferenceType:   t.ReferenceType,
		ReferenceID:     t.ReferenceID,
		Notes:           t.Notes,
		CreatedAt:       t.CreatedAt,
	}
}

func artFromGORM(g *domain.ARTransactionGORM) *domain.ARTransaction {
	return &domain.ARTransaction{
		ID:              g.ID,
		CompanyID:       g.CompanyID,
		CustomerID:      g.CustomerID,
		InvoiceID:       g.InvoiceID,
		TransactionType: domain.ARTransType(g.TransactionType),
		TransactionDate: g.TransactionDate,
		Amount:          g.Amount,
		Currency:        g.Currency,
		ReferenceType:   g.ReferenceType,
		ReferenceID:     g.ReferenceID,
		Notes:           g.Notes,
		CreatedAt:       g.CreatedAt,
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
	if err := q.Order("ar_transactions.transaction_date DESC").Find(&gs).Error; err != nil {
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
	dq := r.db.WithContext(ctx).Joins("JOIN customers ON customers.id = ar_transactions.customer_id").Where("customers.company_id = ?", companyID).Order("ar_transactions.transaction_date DESC")
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
		QNNumber:    sq.QNNumber,
		CustomerID:  sq.CustomerID,
		ValidUntil:  sq.ValidUntil,
		Status:      sq.Status,
		TotalAmount: sq.TotalAmount,
		CreatedBy:   sq.CreatedBy,
		CreatedAt:   sq.CreatedAt,
		UpdatedAt:   sq.CreatedAt,
	}
}

func sqFromGORM(g *domain.SalesQuotationGORM) *domain.SalesQuotation {
	return &domain.SalesQuotation{
		ID:          g.ID,
		CompanyID:   g.CompanyID,
		QNNumber:    g.QNNumber,
		CustomerID:  g.CustomerID,
		ValidUntil:  g.ValidUntil,
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
	return r.db.WithContext(ctx).Model(&domain.SalesQuotationGORM{}).Where("id = ?", sq.ID).Select("*").Updates(g).Error
}
