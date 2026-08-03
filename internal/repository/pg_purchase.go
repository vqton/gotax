package repository

import (
	"context"
	"fmt"
	"gotax/internal/domain"
	"time"

	"gorm.io/gorm"
)

// ─── Supplier ────────────────────────────────────────────────────────────

type PGSupplierRepo struct {
	db *gorm.DB
}

func NewPGSupplierRepo(db *gorm.DB) *PGSupplierRepo {
	return &PGSupplierRepo{db: db}
}

func supplierToGORM(s *domain.Supplier) *domain.SupplierGORM {
	return &domain.SupplierGORM{
		ID:                s.ID,
		CompanyID:         s.CompanyID,
		Code:              s.Code,
		Name:              s.Name,
		TaxCode:           s.TaxCode,
		Address:           s.Address,
		Phone:             s.Phone,
		Email:             s.Email,
		BankAccountName:   s.BankAccountName,
		BankAccountNumber: s.BankAccountNumber,
		BankName:          s.BankName,
		PaymentTerms:      string(s.PaymentTerms),
		CreditLimit:       s.CreditLimit,
		Currency:          s.Currency,
		SupplierType:      string(s.SupplierType),
		Status:            string(s.Status),
		Notes:             s.Notes,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}
}

func supplierFromGORM(g *domain.SupplierGORM) *domain.Supplier {
	return &domain.Supplier{
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
		SupplierType:      domain.SupplierType(g.SupplierType),
		Status:            domain.SupplierStatus(g.Status),
		Notes:             g.Notes,
		CreatedAt:         g.CreatedAt,
		UpdatedAt:         g.UpdatedAt,
	}
}

func (r *PGSupplierRepo) CreateSupplier(ctx context.Context, s *domain.Supplier) error {
	g := supplierToGORM(s)
	if err := r.db.WithContext(ctx).Create(g).Error; err != nil {
		return err
	}
	s.ID = g.ID
	s.CreatedAt = g.CreatedAt
	s.UpdatedAt = g.UpdatedAt
	return nil
}

func (r *PGSupplierRepo) GetSupplier(ctx context.Context, id string) (*domain.Supplier, error) {
	var g domain.SupplierGORM
	if err := r.db.WithContext(ctx).First(&g, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrSupplierNotFound
		}
		return nil, err
	}
	return supplierFromGORM(&g), nil
}

func (r *PGSupplierRepo) GetSupplierByCode(ctx context.Context, companyID, code string) (*domain.Supplier, error) {
	var g domain.SupplierGORM
	if err := r.db.WithContext(ctx).Where("company_id = ? AND code = ?", companyID, code).First(&g).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrSupplierNotFound
		}
		return nil, err
	}
	return supplierFromGORM(&g), nil
}

func (r *PGSupplierRepo) ListSuppliers(ctx context.Context, filter domain.PurchaseOrderFilter) ([]domain.Supplier, int, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&domain.SupplierGORM{}).Where("company_id = ?", filter.CompanyID)
	if filter.SupplierID != "" {
		q = q.Where("id = ?", filter.SupplierID)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var gs []domain.SupplierGORM
	dq := r.db.WithContext(ctx).Where("company_id = ?", filter.CompanyID)
	if filter.SupplierID != "" {
		dq = dq.Where("id = ?", filter.SupplierID)
	}
	dq = dq.Order("code")
	if filter.Limit > 0 {
		dq = dq.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		dq = dq.Offset(filter.Offset)
	}
	if err := dq.Find(&gs).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.Supplier, len(gs))
	for i := range gs {
		out[i] = *supplierFromGORM(&gs[i])
	}
	return out, int(total), nil
}

func (r *PGSupplierRepo) ListSuppliersByIDs(ctx context.Context, ids []string) ([]domain.Supplier, error) {
	var gs []domain.SupplierGORM
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&gs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Supplier, len(gs))
	for i := range gs {
		out[i] = *supplierFromGORM(&gs[i])
	}
	return out, nil
}

func (r *PGSupplierRepo) UpdateSupplier(ctx context.Context, s *domain.Supplier) error {
	g := supplierToGORM(s)
	return r.db.WithContext(ctx).Model(&domain.SupplierGORM{}).Where("id = ?", s.ID).Updates(map[string]interface{}{
		"code":                g.Code,
		"name":                g.Name,
		"tax_code":            g.TaxCode,
		"address":             g.Address,
		"phone":               g.Phone,
		"email":               g.Email,
		"bank_account_name":   g.BankAccountName,
		"bank_account_number": g.BankAccountNumber,
		"bank_name":           g.BankName,
		"payment_terms":       g.PaymentTerms,
		"credit_limit":        g.CreditLimit,
		"currency":            g.Currency,
		"supplier_type":       g.SupplierType,
		"status":              g.Status,
		"notes":               g.Notes,
	}).Error
}

func (r *PGSupplierRepo) DeleteSupplier(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&domain.SupplierGORM{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrSupplierNotFound
	}
	return nil
}

// ─── Purchase Order ──────────────────────────────────────────────────────

type PGPurchaseOrderRepo struct {
	db *gorm.DB
}

func NewPGPurchaseRepo(db *gorm.DB) *PGPurchaseOrderRepo {
	return &PGPurchaseOrderRepo{db: db}
}

func poToGORM(po *domain.PurchaseOrder) *domain.PurchaseOrderGORM {
	return &domain.PurchaseOrderGORM{
		ID:              po.ID,
		CompanyID:       po.CompanyID,
		PONumber:        po.PONumber,
		SupplierID:      po.SupplierID,
		RequisitionID:   po.RequisitionID,
		OrderDate:       safeTimeStr(po.OrderDate),
		ExpectedDate:    safeTimePtrStr(po.ExpectedDate),
		Currency:        po.Currency,
		ExchangeRate:    po.ExchangeRate,
		PaymentTerms:    po.PaymentTerms,
		DeliveryTerms:   po.DeliveryTerms,
		Subtotal:        po.Subtotal,
		DiscountAmount:  po.DiscountAmount,
		TaxAmount:       po.TaxAmount,
		TotalAmount:     po.TotalAmount,
		Status:          string(po.Status),
		ApprovedBy:      po.ApprovedBy,
		ApprovedAt:      safeTimePtrRFC3339(po.ApprovedAt),
		CancelledReason: po.CancelledReason,
		Notes:           po.Notes,
		CreatedBy:       po.CreatedBy,
		CreatedAt:       po.CreatedAt,
		UpdatedAt:       po.UpdatedAt,
	}
}

func poFromGORM(g *domain.PurchaseOrderGORM) *domain.PurchaseOrder {
	return &domain.PurchaseOrder{
		ID:              g.ID,
		CompanyID:       g.CompanyID,
		PONumber:        g.PONumber,
		SupplierID:      g.SupplierID,
		RequisitionID:   g.RequisitionID,
		OrderDate:       parseDate(g.OrderDate),
		ExpectedDate:    timePtr(parseDate(g.ExpectedDate)),
		Currency:        g.Currency,
		ExchangeRate:    g.ExchangeRate,
		PaymentTerms:    g.PaymentTerms,
		DeliveryTerms:   g.DeliveryTerms,
		Subtotal:        g.Subtotal,
		DiscountAmount:  g.DiscountAmount,
		TaxAmount:       g.TaxAmount,
		TotalAmount:     g.TotalAmount,
		Status:          domain.POStatus(g.Status),
		ApprovedBy:      g.ApprovedBy,
		ApprovedAt:      timePtr(parseDateTime(g.ApprovedAt)),
		CancelledReason: g.CancelledReason,
		Notes:           g.Notes,
		CreatedBy:       g.CreatedBy,
		CreatedAt:       g.CreatedAt,
		UpdatedAt:       g.UpdatedAt,
	}
}

func (r *PGPurchaseOrderRepo) CreatePO(ctx context.Context, po *domain.PurchaseOrder) error {
	g := poToGORM(po)
	if err := r.db.WithContext(ctx).Create(g).Error; err != nil {
		return err
	}
	po.ID = g.ID
	po.CreatedAt = g.CreatedAt
	po.UpdatedAt = g.UpdatedAt
	return nil
}

func (r *PGPurchaseOrderRepo) GetPO(ctx context.Context, id string) (*domain.PurchaseOrder, error) {
	var g domain.PurchaseOrderGORM
	if err := r.db.WithContext(ctx).Preload("Lines").First(&g, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPONotFound
		}
		return nil, err
	}
	po := poFromGORM(&g)
	po.Lines = make([]domain.POItem, len(g.Lines))
	for i, l := range g.Lines {
		po.Lines[i] = *poItemFromGORM(&l)
	}
	return po, nil
}

func (r *PGPurchaseOrderRepo) GetPOByNumber(ctx context.Context, companyID, poNumber string) (*domain.PurchaseOrder, error) {
	var g domain.PurchaseOrderGORM
	if err := r.db.WithContext(ctx).Preload("Lines").Where("company_id = ? AND po_number = ?", companyID, poNumber).First(&g).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPONotFound
		}
		return nil, err
	}
	po := poFromGORM(&g)
	po.Lines = make([]domain.POItem, len(g.Lines))
	for i, l := range g.Lines {
		po.Lines[i] = *poItemFromGORM(&l)
	}
	return po, nil
}

func (r *PGPurchaseOrderRepo) ListPOs(ctx context.Context, filter domain.PurchaseOrderFilter) ([]domain.PurchaseOrder, int, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&domain.PurchaseOrderGORM{}).Where("company_id = ?", filter.CompanyID)
	if filter.SupplierID != "" {
		q = q.Where("supplier_id = ?", filter.SupplierID)
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
	var gs []domain.PurchaseOrderGORM
	dq := r.db.WithContext(ctx).Where("company_id = ?", filter.CompanyID)
	if filter.SupplierID != "" {
		dq = dq.Where("supplier_id = ?", filter.SupplierID)
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
	out := make([]domain.PurchaseOrder, len(gs))
	for i := range gs {
		po := poFromGORM(&gs[i])
		r.db.WithContext(ctx).Model(&domain.POItemGORM{}).Where("po_id = ?", gs[i].ID).Find(&po.Lines)
		out[i] = *po
	}
	return out, int(total), nil
}

func (r *PGPurchaseOrderRepo) UpdatePO(ctx context.Context, po *domain.PurchaseOrder) error {
	g := poToGORM(po)
	return r.db.WithContext(ctx).Model(&domain.PurchaseOrderGORM{}).Where("id = ?", po.ID).Updates(map[string]interface{}{
		"supplier_id":    g.SupplierID,
		"order_date":     g.OrderDate,
		"expected_date":  g.ExpectedDate,
		"currency":       g.Currency,
		"delivery_terms": g.DeliveryTerms,
		"payment_terms":  g.PaymentTerms,
		"notes":          g.Notes,
	}).Error
}

func (r *PGPurchaseOrderRepo) UpdatePOStatus(ctx context.Context, id string, status domain.POStatus) error {
	return r.db.WithContext(ctx).Model(&domain.PurchaseOrderGORM{}).Where("id = ?", id).Update("status", string(status)).Error
}

func (r *PGPurchaseOrderRepo) ApprovePO(ctx context.Context, id string, approvedBy string, approvedAt time.Time) error {
	return r.db.WithContext(ctx).Model(&domain.PurchaseOrderGORM{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      string(domain.POStatusApproved),
		"approved_by": approvedBy,
		"approved_at": approvedAt.Format(time.RFC3339),
	}).Error
}

func (r *PGPurchaseOrderRepo) CancelPO(ctx context.Context, id string, cancelReason string) error {
	return r.db.WithContext(ctx).Model(&domain.PurchaseOrderGORM{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":           string(domain.POStatusCancelled),
		"cancelled_reason": cancelReason,
	}).Error
}

func poItemToGORM(l *domain.POItem) *domain.POItemGORM {
	return &domain.POItemGORM{
		POID:          l.POID,
		LineNumber:    l.LineNumber,
		ItemCode:      l.ItemCode,
		ItemName:      l.ItemName,
		Unit:          l.Unit,
		Quantity:      l.Quantity,
		UnitPrice:     l.UnitPrice,
		DiscountPct:   l.DiscountPct,
		VATRate:       l.VATRate,
		VATType:       string(l.VATType),
		AccountID:     l.AccountID,
		VATAccountID:  l.VATAccountID,
		LineTotal:     l.LineTotal,
		LineVATAmount: l.LineVATAmount,
		ReceivedQty:   l.ReceivedQty,
		InvoicedQty:   l.InvoicedQty,
	}
}

func poItemFromGORM(g *domain.POItemGORM) *domain.POItem {
	return &domain.POItem{
		ID:            g.ID,
		POID:          g.POID,
		LineNumber:    g.LineNumber,
		ItemCode:      g.ItemCode,
		ItemName:      g.ItemName,
		Unit:          g.Unit,
		Quantity:      g.Quantity,
		UnitPrice:     g.UnitPrice,
		DiscountPct:   g.DiscountPct,
		VATRate:       g.VATRate,
		VATType:       domain.VATType(g.VATType),
		AccountID:     g.AccountID,
		VATAccountID:  g.VATAccountID,
		LineTotal:     g.LineTotal,
		LineVATAmount: g.LineVATAmount,
		ReceivedQty:   g.ReceivedQty,
		InvoicedQty:   g.InvoicedQty,
	}
}

func (r *PGPurchaseOrderRepo) GetPOLines(ctx context.Context, poID string) ([]domain.POItem, error) {
	var gs []domain.POItemGORM
	if err := r.db.WithContext(ctx).Where("po_id = ?", poID).Order("line_number").Find(&gs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.POItem, len(gs))
	for i := range gs {
		out[i] = *poItemFromGORM(&gs[i])
	}
	return out, nil
}

func (r *PGPurchaseOrderRepo) CreatePOLines(ctx context.Context, items []domain.POItem) error {
	if len(items) == 0 {
		return nil
	}
	gs := make([]domain.POItemGORM, len(items))
	for i := range items {
		gs[i] = *poItemToGORM(&items[i])
	}
	return r.db.WithContext(ctx).Create(&gs).Error
}

func (r *PGPurchaseOrderRepo) UpdatePOLines(ctx context.Context, items []domain.POItem) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("po_id = ?", items[0].POID).Delete(&domain.POItemGORM{}).Error; err != nil {
			return err
		}
		gs := make([]domain.POItemGORM, len(items))
		for i := range items {
			gs[i] = *poItemToGORM(&items[i])
		}
		return tx.Create(&gs).Error
	})
}

func (r *PGPurchaseOrderRepo) NextPONumber(ctx context.Context, companyID, yyyymm string) (string, error) {
	prefix := fmt.Sprintf("PO-%s-", yyyymm)
	var maxNum int
	if err := r.db.WithContext(ctx).Raw(`SELECT COALESCE(MAX(CAST(SUBSTRING(po_number FROM '([0-9]+)$') AS INTEGER)), 0) FROM purchase_orders WHERE company_id = ? AND po_number LIKE ?`, companyID, prefix+"%").Scan(&maxNum).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%04d", prefix, maxNum+1), nil
}

// ─── GRN ────────────────────────────────────────────────────────────────

type PGGRNRepo struct {
	db *gorm.DB
}

func NewPGGRNRepo(db *gorm.DB) *PGGRNRepo {
	return &PGGRNRepo{db: db}
}

func grnToGORM(g *domain.GRN) *domain.GRNGORM {
	return &domain.GRNGORM{
		ID:          g.ID,
		CompanyID:   g.CompanyID,
		GRNNumber:   g.GRNNumber,
		POID:        g.POID,
		ReceiptDate: safeTimeStr(g.ReceiptDate),
		Warehouse:   g.Warehouse,
		Status:      string(g.Status),
		Notes:       g.Notes,
		CreatedBy:   g.CreatedBy,
		CreatedAt:   g.CreatedAt,
	}
}

func grnFromGORM(gg *domain.GRNGORM) *domain.GRN {
	return &domain.GRN{
		ID:          gg.ID,
		CompanyID:   gg.CompanyID,
		GRNNumber:   gg.GRNNumber,
		POID:        gg.POID,
		ReceiptDate: parseDate(gg.ReceiptDate),
		Warehouse:   gg.Warehouse,
		Status:      domain.GRNStatus(gg.Status),
		Notes:       gg.Notes,
		CreatedBy:   gg.CreatedBy,
		CreatedAt:   gg.CreatedAt,
	}
}

func (r *PGGRNRepo) CreateGRN(ctx context.Context, g *domain.GRN) error {
	gg := grnToGORM(g)
	if err := r.db.WithContext(ctx).Create(gg).Error; err != nil {
		return err
	}
	g.ID = gg.ID
	g.CreatedAt = gg.CreatedAt
	return nil
}

func (r *PGGRNRepo) GetGRN(ctx context.Context, id string) (*domain.GRN, error) {
	var gg domain.GRNGORM
	if err := r.db.WithContext(ctx).Preload("Lines").First(&gg, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrGRNNotFound
		}
		return nil, err
	}
	g := grnFromGORM(&gg)
	g.Lines = make([]domain.GRNItem, len(gg.Lines))
	for i, l := range gg.Lines {
		g.Lines[i] = *grnItemFromGORM(&l)
	}
	return g, nil
}

func (r *PGGRNRepo) GetGRNByNumber(ctx context.Context, companyID, grnNumber string) (*domain.GRN, error) {
	var gg domain.GRNGORM
	if err := r.db.WithContext(ctx).Preload("Lines").Where("company_id = ? AND grn_number = ?", companyID, grnNumber).First(&gg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrGRNNotFound
		}
		return nil, err
	}
	g := grnFromGORM(&gg)
	g.Lines = make([]domain.GRNItem, len(gg.Lines))
	for i, l := range gg.Lines {
		g.Lines[i] = *grnItemFromGORM(&l)
	}
	return g, nil
}

func (r *PGGRNRepo) ListGRNs(ctx context.Context, filter domain.GRNFilter) ([]domain.GRN, int, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&domain.GRNGORM{}).Where("company_id = ?", filter.CompanyID)
	if filter.POID != "" {
		q = q.Where("po_id = ?", filter.POID)
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
	var gs []domain.GRNGORM
	dq := r.db.WithContext(ctx).Where("company_id = ?", filter.CompanyID)
	if filter.POID != "" {
		dq = dq.Where("po_id = ?", filter.POID)
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
	out := make([]domain.GRN, len(gs))
	for i := range gs {
		g := grnFromGORM(&gs[i])
		r.db.WithContext(ctx).Model(&domain.GRNItemGORM{}).Where("grn_id = ?", gs[i].ID).Find(&g.Lines)
		out[i] = *g
	}
	return out, int(total), nil
}

func (r *PGGRNRepo) UpdateGRN(ctx context.Context, g *domain.GRN) error {
	gg := grnToGORM(g)
	return r.db.WithContext(ctx).Model(&domain.GRNGORM{}).Where("id = ?", g.ID).Updates(map[string]interface{}{
		"receipt_date": gg.ReceiptDate,
		"po_id":        gg.POID,
		"status":       gg.Status,
		"notes":        gg.Notes,
	}).Error
}

func (r *PGGRNRepo) UpdateGRNStatus(ctx context.Context, id string, status domain.GRNStatus) error {
	return r.db.WithContext(ctx).Model(&domain.GRNGORM{}).Where("id = ?", id).Update("status", string(status)).Error
}

func grnItemToGORM(l *domain.GRNItem) *domain.GRNItemGORM {
	return &domain.GRNItemGORM{
		GRNID:            l.GRNID,
		POLineID:         l.POLineID,
		ItemCode:         l.ItemCode,
		ItemName:         l.ItemName,
		Unit:             l.Unit,
		QuantityReceived: l.QuantityReceived,
		QuantityRejected: l.QuantityRejected,
		UnitPrice:        l.UnitPrice,
		LineTotal:        l.LineTotal,
	}
}

func grnItemFromGORM(g *domain.GRNItemGORM) *domain.GRNItem {
	return &domain.GRNItem{
		ID:               g.ID,
		GRNID:            g.GRNID,
		POLineID:         g.POLineID,
		ItemCode:         g.ItemCode,
		ItemName:         g.ItemName,
		Unit:             g.Unit,
		QuantityReceived: g.QuantityReceived,
		QuantityRejected: g.QuantityRejected,
		UnitPrice:        g.UnitPrice,
		LineTotal:        g.LineTotal,
	}
}

func (r *PGGRNRepo) GetGRNLines(ctx context.Context, grnID string) ([]domain.GRNItem, error) {
	var gs []domain.GRNItemGORM
	if err := r.db.WithContext(ctx).Where("grn_id = ?", grnID).Find(&gs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.GRNItem, len(gs))
	for i := range gs {
		out[i] = *grnItemFromGORM(&gs[i])
	}
	return out, nil
}

func (r *PGGRNRepo) CreateGRNLines(ctx context.Context, items []domain.GRNItem) error {
	if len(items) == 0 {
		return nil
	}
	gs := make([]domain.GRNItemGORM, len(items))
	for i := range items {
		gs[i] = *grnItemToGORM(&items[i])
	}
	return r.db.WithContext(ctx).Create(&gs).Error
}

func (r *PGGRNRepo) UpdateGRNLines(ctx context.Context, items []domain.GRNItem) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("grn_id = ?", items[0].GRNID).Delete(&domain.GRNItemGORM{}).Error; err != nil {
			return err
		}
		gs := make([]domain.GRNItemGORM, len(items))
		for i := range items {
			gs[i] = *grnItemToGORM(&items[i])
		}
		return tx.Create(&gs).Error
	})
}

func (r *PGGRNRepo) NextGRNNumber(ctx context.Context, companyID, yyyymm string) (string, error) {
	prefix := fmt.Sprintf("GRN-%s-", yyyymm)
	var maxNum int
	if err := r.db.WithContext(ctx).Raw(`SELECT COALESCE(MAX(CAST(SUBSTRING(grn_number FROM '([0-9]+)$') AS INTEGER)), 0) FROM goods_receipt_notes WHERE company_id = ? AND grn_number LIKE ?`, companyID, prefix+"%").Scan(&maxNum).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%04d", prefix, maxNum+1), nil
}

// ─── Supplier Invoice ──────────────────────────────────────────────────

type PGSupplierInvoiceRepo struct {
	db *gorm.DB
}

func NewPGSupplierInvoiceRepo(db *gorm.DB) *PGSupplierInvoiceRepo {
	return &PGSupplierInvoiceRepo{db: db}
}

func sinvToGORM(inv *domain.SupplierInvoice) *domain.SupplierInvoiceGORM {
	return &domain.SupplierInvoiceGORM{
		ID:                inv.ID,
		CompanyID:         inv.CompanyID,
		InvoiceNumber:     inv.InvoiceNumber,
		InvoiceDate:       safeTimeStr(inv.InvoiceDate),
		POID:              inv.POID,
		GRNID:             inv.GRNID,
		SupplierID:        inv.SupplierID,
		SupplierName:      inv.SupplierName,
		SupplierTaxCode:   inv.SupplierTaxCode,
		InvoiceType:       inv.InvoiceType,
		Currency:          inv.Currency,
		ExchangeRate:      inv.ExchangeRate,
		Subtotal:          inv.Subtotal,
		DiscountAmount:    inv.DiscountAmount,
		TaxAmount:         inv.TaxAmount,
		TotalAmount:       inv.TotalAmount,
		AmountPaid:        inv.AmountPaid,
		BalanceDue:        inv.BalanceDue,
		DueDate:           safeTimePtrStr(inv.DueDate),
		VATDeductionStatus: string(inv.VATDeductionStatus),
		EInvoiceData:      inv.EInvoiceData,
		EInvoiceCode:      inv.EInvoiceCode,
		Status:            string(inv.Status),
		GLPosted:          inv.GLPosted,
		GLPostedAt:        safeTimePtrRFC3339(inv.GLPostedAt),
		Notes:             inv.Notes,
		CreatedBy:         inv.CreatedBy,
		CreatedAt:         inv.CreatedAt,
	}
}

func sinvFromGORM(g *domain.SupplierInvoiceGORM) *domain.SupplierInvoice {
	return &domain.SupplierInvoice{
		ID:                g.ID,
		CompanyID:         g.CompanyID,
		InvoiceNumber:     g.InvoiceNumber,
		InvoiceDate:       parseDate(g.InvoiceDate),
		POID:              g.POID,
		GRNID:             g.GRNID,
		SupplierID:        g.SupplierID,
		SupplierName:      g.SupplierName,
		SupplierTaxCode:   g.SupplierTaxCode,
		InvoiceType:       g.InvoiceType,
		Currency:          g.Currency,
		ExchangeRate:      g.ExchangeRate,
		Subtotal:          g.Subtotal,
		DiscountAmount:    g.DiscountAmount,
		TaxAmount:         g.TaxAmount,
		TotalAmount:       g.TotalAmount,
		AmountPaid:        g.AmountPaid,
		BalanceDue:        g.BalanceDue,
		DueDate:           timePtr(parseDate(g.DueDate)),
		VATDeductionStatus: domain.VATDeductionStatus(g.VATDeductionStatus),
		EInvoiceData:      g.EInvoiceData,
		EInvoiceCode:      g.EInvoiceCode,
		Status:            domain.InvoiceStatus(g.Status),
		GLPosted:          g.GLPosted,
		GLPostedAt:        timePtr(parseDateTime(g.GLPostedAt)),
		Notes:             g.Notes,
		CreatedBy:         g.CreatedBy,
		CreatedAt:         g.CreatedAt,
	}
}

func (r *PGSupplierInvoiceRepo) CreateInvoice(ctx context.Context, inv *domain.SupplierInvoice) error {
	g := sinvToGORM(inv)
	if err := r.db.WithContext(ctx).Create(g).Error; err != nil {
		return err
	}
	inv.ID = g.ID
	inv.CreatedAt = g.CreatedAt
	return nil
}

func (r *PGSupplierInvoiceRepo) GetInvoice(ctx context.Context, id string) (*domain.SupplierInvoice, error) {
	var g domain.SupplierInvoiceGORM
	if err := r.db.WithContext(ctx).Preload("Lines").First(&g, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrSupplierInvoiceNotFound
		}
		return nil, err
	}
	inv := sinvFromGORM(&g)
	inv.Lines = make([]domain.SupplierInvoiceLine, len(g.Lines))
	for i, l := range g.Lines {
		inv.Lines[i] = *sinvLineFromGORM(&l)
	}
	return inv, nil
}

func (r *PGSupplierInvoiceRepo) GetInvoiceByNumber(ctx context.Context, companyID, invoiceNumber string) (*domain.SupplierInvoice, error) {
	var g domain.SupplierInvoiceGORM
	if err := r.db.WithContext(ctx).Preload("Lines").Where("company_id = ? AND invoice_number = ?", companyID, invoiceNumber).First(&g).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrSupplierInvoiceNotFound
		}
		return nil, err
	}
	inv := sinvFromGORM(&g)
	inv.Lines = make([]domain.SupplierInvoiceLine, len(g.Lines))
	for i, l := range g.Lines {
		inv.Lines[i] = *sinvLineFromGORM(&l)
	}
	return inv, nil
}

func (r *PGSupplierInvoiceRepo) ListInvoices(ctx context.Context, filter domain.SupplierInvoiceFilter) ([]domain.SupplierInvoice, int, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&domain.SupplierInvoiceGORM{}).Where("company_id = ?", filter.CompanyID)
	if filter.SupplierID != "" {
		q = q.Where("supplier_id = ?", filter.SupplierID)
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
	var gs []domain.SupplierInvoiceGORM
	dq := r.db.WithContext(ctx).Where("company_id = ?", filter.CompanyID)
	if filter.SupplierID != "" {
		dq = dq.Where("supplier_id = ?", filter.SupplierID)
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
	out := make([]domain.SupplierInvoice, len(gs))
	for i := range gs {
		inv := sinvFromGORM(&gs[i])
		r.db.WithContext(ctx).Model(&domain.SupplierInvoiceLineGORM{}).Where("invoice_id = ?", gs[i].ID).Find(&inv.Lines)
		out[i] = *inv
	}
	return out, int(total), nil
}

func (r *PGSupplierInvoiceRepo) UpdateInvoice(ctx context.Context, inv *domain.SupplierInvoice) error {
	g := sinvToGORM(inv)
	return r.db.WithContext(ctx).Model(&domain.SupplierInvoiceGORM{}).Where("id = ?", inv.ID).Updates(map[string]interface{}{
		"invoice_date": g.InvoiceDate,
		"subtotal":     g.Subtotal,
		"tax_amount":   g.TaxAmount,
		"total_amount": g.TotalAmount,
		"due_date":     g.DueDate,
		"notes":        g.Notes,
	}).Error
}

func (r *PGSupplierInvoiceRepo) UpdateInvoiceStatus(ctx context.Context, id string, status domain.InvoiceStatus) error {
	return r.db.WithContext(ctx).Model(&domain.SupplierInvoiceGORM{}).Where("id = ?", id).Update("status", string(status)).Error
}

func (r *PGSupplierInvoiceRepo) PostInvoice(ctx context.Context, id string, postedAt time.Time) error {
	return r.db.WithContext(ctx).Model(&domain.SupplierInvoiceGORM{}).Where("id = ?", id).Update("status", string(domain.InvoicePosted)).Error
}

func (r *PGSupplierInvoiceRepo) SetInvoiceGLPosted(ctx context.Context, id string, postedAt time.Time) error {
	return r.db.WithContext(ctx).Model(&domain.SupplierInvoiceGORM{}).Where("id = ?", id).Updates(map[string]interface{}{
		"gl_posted":    true,
		"gl_posted_at": postedAt.Format(time.RFC3339),
	}).Error
}

func sinvLineToGORM(l *domain.SupplierInvoiceLine) *domain.SupplierInvoiceLineGORM {
	return &domain.SupplierInvoiceLineGORM{
		InvoiceID:     l.InvoiceID,
		POLineID:      l.POLineID,
		GRNLineID:     l.GRNLineID,
		ItemCode:      l.ItemCode,
		ItemName:      l.ItemName,
		Unit:          l.Unit,
		Quantity:      l.Quantity,
		UnitPrice:     l.UnitPrice,
		VATRate:       l.VATRate,
		VATType:       string(l.VATType),
		LineTotal:     l.LineTotal,
		LineVATAmount: l.LineVATAmount,
		AccountID:     l.AccountID,
		VATAccountID:  l.VATAccountID,
	}
}

func sinvLineFromGORM(g *domain.SupplierInvoiceLineGORM) *domain.SupplierInvoiceLine {
	return &domain.SupplierInvoiceLine{
		ID:            g.ID,
		InvoiceID:     g.InvoiceID,
		POLineID:      g.POLineID,
		GRNLineID:     g.GRNLineID,
		ItemCode:      g.ItemCode,
		ItemName:      g.ItemName,
		Unit:          g.Unit,
		Quantity:      g.Quantity,
		UnitPrice:     g.UnitPrice,
		VATRate:       g.VATRate,
		VATType:       domain.VATType(g.VATType),
		LineTotal:     g.LineTotal,
		LineVATAmount: g.LineVATAmount,
		AccountID:     g.AccountID,
		VATAccountID:  g.VATAccountID,
	}
}

func (r *PGSupplierInvoiceRepo) GetInvoiceLines(ctx context.Context, invoiceID string) ([]domain.SupplierInvoiceLine, error) {
	var gs []domain.SupplierInvoiceLineGORM
	if err := r.db.WithContext(ctx).Where("invoice_id = ?", invoiceID).Find(&gs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.SupplierInvoiceLine, len(gs))
	for i := range gs {
		out[i] = *sinvLineFromGORM(&gs[i])
	}
	return out, nil
}

func (r *PGSupplierInvoiceRepo) CreateInvoiceLines(ctx context.Context, items []domain.SupplierInvoiceLine) error {
	if len(items) == 0 {
		return nil
	}
	gs := make([]domain.SupplierInvoiceLineGORM, len(items))
	for i := range items {
		gs[i] = *sinvLineToGORM(&items[i])
	}
	return r.db.WithContext(ctx).Create(&gs).Error
}

func (r *PGSupplierInvoiceRepo) UpdateInvoiceLines(ctx context.Context, items []domain.SupplierInvoiceLine) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("invoice_id = ?", items[0].InvoiceID).Delete(&domain.SupplierInvoiceLineGORM{}).Error; err != nil {
			return err
		}
		gs := make([]domain.SupplierInvoiceLineGORM, len(items))
		for i := range items {
			gs[i] = *sinvLineToGORM(&items[i])
		}
		return tx.Create(&gs).Error
	})
}

// ─── AP Transaction ────────────────────────────────────────────────────

type PGAPTransactionRepo struct {
	db *gorm.DB
}

func NewPGAPTransactionRepo(db *gorm.DB) *PGAPTransactionRepo {
	return &PGAPTransactionRepo{db: db}
}

func aptToGORM(t *domain.APTransaction) *domain.APTransactionGORM {
	return &domain.APTransactionGORM{
		ID:              t.ID,
		CompanyID:       t.CompanyID,
		SupplierID:      t.SupplierID,
		InvoiceID:       t.InvoiceID,
		TransactionType: string(t.TransactionType),
		TransactionDate: t.TransactionDate.Format(time.RFC3339),
		Amount:          t.Amount,
		Currency:        t.Currency,
		ReferenceType:   t.ReferenceType,
		ReferenceID:     t.ReferenceID,
		Notes:           t.Notes,
		CreatedAt:       t.CreatedAt,
	}
}

func aptFromGORM(g *domain.APTransactionGORM) *domain.APTransaction {
	return &domain.APTransaction{
		ID:              g.ID,
		CompanyID:       g.CompanyID,
		SupplierID:      g.SupplierID,
		InvoiceID:       g.InvoiceID,
		TransactionType: domain.APTransactionType(g.TransactionType),
		TransactionDate: parseDateTime(g.TransactionDate),
		Amount:          g.Amount,
		Currency:        g.Currency,
		ReferenceType:   g.ReferenceType,
		ReferenceID:     g.ReferenceID,
		Notes:           g.Notes,
		CreatedAt:       g.CreatedAt,
	}
}

func (r *PGAPTransactionRepo) CreateAPTransaction(ctx context.Context, t *domain.APTransaction) error {
	g := aptToGORM(t)
	if err := r.db.WithContext(ctx).Create(g).Error; err != nil {
		return err
	}
	t.ID = g.ID
	t.CreatedAt = g.CreatedAt
	return nil
}

func (r *PGAPTransactionRepo) GetAPTransaction(ctx context.Context, id string) (*domain.APTransaction, error) {
	var g domain.APTransactionGORM
	if err := r.db.WithContext(ctx).First(&g, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrAPTransNotFound
		}
		return nil, err
	}
	return aptFromGORM(&g), nil
}

func (r *PGAPTransactionRepo) ListAPTransactionsBySupplier(ctx context.Context, companyID, supplierID string) ([]domain.APTransaction, error) {
	var gs []domain.APTransactionGORM
	q := r.db.WithContext(ctx).Where("supplier_id = ?", supplierID)
	if companyID != "" {
		q = q.Where("company_id = ?", companyID)
	}
	if err := q.Order("created_at DESC").Find(&gs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.APTransaction, len(gs))
	for i := range gs {
		out[i] = *aptFromGORM(&gs[i])
	}
	return out, nil
}

func (r *PGAPTransactionRepo) ListAPTransactions(ctx context.Context, companyID string, offset, limit int) ([]domain.APTransaction, int, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&domain.APTransactionGORM{}).Joins("JOIN suppliers ON suppliers.id = ap_transactions.supplier_id").Where("suppliers.company_id = ?", companyID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var gs []domain.APTransactionGORM
	dq := r.db.WithContext(ctx).Joins("JOIN suppliers ON suppliers.id = ap_transactions.supplier_id").Where("suppliers.company_id = ?", companyID).Order("ap_transactions.created_at DESC")
	if limit > 0 {
		dq = dq.Limit(limit)
	}
	if offset > 0 {
		dq = dq.Offset(offset)
	}
	if err := dq.Find(&gs).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.APTransaction, len(gs))
	for i := range gs {
		out[i] = *aptFromGORM(&gs[i])
	}
	return out, int(total), nil
}

// ─── Cost Allocation ────────────────────────────────────────────────────

type PGCostAllocationRepo struct {
	db *gorm.DB
}

func NewPGCostAllocationRepo(db *gorm.DB) *PGCostAllocationRepo {
	return &PGCostAllocationRepo{db: db}
}

func caToGORM(c *domain.CostAllocation) *domain.CostAllocationGORM {
	return &domain.CostAllocationGORM{
		CompanyID:        c.CompanyID,
		InvoiceID:        c.InvoiceID,
		CostType:         string(c.CostType),
		CostAmount:       c.CostAmount,
		AllocationMethod: string(c.AllocationMethod),
		AllocatedLines:   c.AllocatedLines,
		Notes:            c.Notes,
	}
}

func caFromGORM(g *domain.CostAllocationGORM) *domain.CostAllocation {
	return &domain.CostAllocation{
		ID:               g.ID,
		CompanyID:        g.CompanyID,
		InvoiceID:        g.InvoiceID,
		CostType:         domain.CostAllocationType(g.CostType),
		CostAmount:       g.CostAmount,
		AllocationMethod: domain.CostAllocationMethod(g.AllocationMethod),
		AllocatedLines:   g.AllocatedLines,
		Notes:            g.Notes,
	}
}

func (r *PGCostAllocationRepo) CreateCostAllocation(ctx context.Context, c *domain.CostAllocation) error {
	g := caToGORM(c)
	if err := r.db.WithContext(ctx).Create(g).Error; err != nil {
		return err
	}
	return nil
}

func (r *PGCostAllocationRepo) GetCostAllocation(ctx context.Context, id string) (*domain.CostAllocation, error) {
	var g domain.CostAllocationGORM
	if err := r.db.WithContext(ctx).First(&g, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("cost allocation %q not found", id)
		}
		return nil, err
	}
	return caFromGORM(&g), nil
}

func (r *PGCostAllocationRepo) ListCostAllocationsByInvoice(ctx context.Context, invoiceID string) ([]domain.CostAllocation, error) {
	var gs []domain.CostAllocationGORM
	if err := r.db.WithContext(ctx).Where("invoice_id = ?", invoiceID).Find(&gs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.CostAllocation, len(gs))
	for i := range gs {
		out[i] = *caFromGORM(&gs[i])
	}
	return out, nil
}

// ─── Doubtful Debt Provisions ───────────────────────────────────────────

type PGDoubtfulDebtProvisionRepo struct{ db *gorm.DB }

func NewPGDoubtfulDebtProvisionRepo(db *gorm.DB) *PGDoubtfulDebtProvisionRepo {
	return &PGDoubtfulDebtProvisionRepo{db: db}
}

func ddpToGORM(p *domain.DoubtfulDebtProvision) *domain.DoubtfulDebtProvisionGORM {
	return &domain.DoubtfulDebtProvisionGORM{
		ID:               p.ID,
		CompanyID:        p.CompanyID,
		AsOfDate:         parseDate(p.AsOfDate),
		TotalOutstanding: p.TotalOutstanding,
		TotalProvision:   p.TotalProvision,
		Status:           string(p.Status),
		CreatedBy:        p.CreatedBy,
	}
}

func ddpFromGORM(g *domain.DoubtfulDebtProvisionGORM) *domain.DoubtfulDebtProvision {
	return &domain.DoubtfulDebtProvision{
		ID:               g.ID,
		CompanyID:        g.CompanyID,
		AsOfDate:         g.AsOfDate.Format(time.DateOnly),
		TotalOutstanding: g.TotalOutstanding,
		TotalProvision:   g.TotalProvision,
		Status:           domain.ProvisionStatus(g.Status),
		CreatedBy:        g.CreatedBy,
		CreatedAt:        g.CreatedAt,
	}
}

func ddplToGORM(l *domain.DoubtfulDebtProvisionLine) *domain.DoubtfulDebtProvisionLineGORM {
	return &domain.DoubtfulDebtProvisionLineGORM{
		ID:                l.ID,
		ProvisionID:       l.ProvisionID,
		SupplierID:        l.SupplierID,
		SupplierName:      l.SupplierName,
		TaxCode:           l.TaxCode,
		OutstandingAmount: l.OutstandingAmount,
		AgeMonths:         l.AgeMonths,
		RatePct:           l.RatePct,
		ProvisionAmount:   l.ProvisionAmount,
	}
}

func ddplFromGORM(g *domain.DoubtfulDebtProvisionLineGORM) *domain.DoubtfulDebtProvisionLine {
	return &domain.DoubtfulDebtProvisionLine{
		ID:                g.ID,
		ProvisionID:       g.ProvisionID,
		SupplierID:        g.SupplierID,
		SupplierName:      g.SupplierName,
		TaxCode:           g.TaxCode,
		OutstandingAmount: g.OutstandingAmount,
		AgeMonths:         g.AgeMonths,
		RatePct:           g.RatePct,
		ProvisionAmount:   g.ProvisionAmount,
	}
}

func (r *PGDoubtfulDebtProvisionRepo) CreateProvision(ctx context.Context, p *domain.DoubtfulDebtProvision) error {
	g := ddpToGORM(p)
	if err := r.db.WithContext(ctx).Create(g).Error; err != nil {
		return err
	}
	return nil
}

func (r *PGDoubtfulDebtProvisionRepo) CreateProvisionLines(ctx context.Context, lines []domain.DoubtfulDebtProvisionLine) error {
	gs := make([]domain.DoubtfulDebtProvisionLineGORM, len(lines))
	for i := range lines {
		gs[i] = *ddplToGORM(&lines[i])
	}
	return r.db.WithContext(ctx).Create(&gs).Error
}

func (r *PGDoubtfulDebtProvisionRepo) GetProvision(ctx context.Context, id string) (*domain.DoubtfulDebtProvision, error) {
	var g domain.DoubtfulDebtProvisionGORM
	if err := r.db.WithContext(ctx).First(&g, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrProvisionNotFound
		}
		return nil, err
	}
	p := ddpFromGORM(&g)
	lines, err := r.GetProvisionLines(ctx, id)
	if err != nil {
		return nil, err
	}
	p.Lines = lines
	return p, nil
}

func (r *PGDoubtfulDebtProvisionRepo) GetProvisionLines(ctx context.Context, provisionID string) ([]domain.DoubtfulDebtProvisionLine, error) {
	var gs []domain.DoubtfulDebtProvisionLineGORM
	if err := r.db.WithContext(ctx).Where("provision_id = ?", provisionID).Order("supplier_name").Find(&gs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.DoubtfulDebtProvisionLine, len(gs))
	for i := range gs {
		out[i] = *ddplFromGORM(&gs[i])
	}
	return out, nil
}

func (r *PGDoubtfulDebtProvisionRepo) ListProvisions(ctx context.Context, companyID string, limit, offset int) ([]domain.DoubtfulDebtProvision, int, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&domain.DoubtfulDebtProvisionGORM{}).Where("company_id = ?", companyID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var gs []domain.DoubtfulDebtProvisionGORM
	if err := q.Order("as_of_date DESC, created_at DESC").Limit(limit).Offset(offset).Find(&gs).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.DoubtfulDebtProvision, len(gs))
	for i := range gs {
		out[i] = *ddpFromGORM(&gs[i])
	}
	return out, int(total), nil
}

// ─── Requisition ─────────────────────────────────────────────────────────

type PGRequisitionRepo struct{ db *gorm.DB }

func NewPGRequisitionRepo(db *gorm.DB) *PGRequisitionRepo {
	return &PGRequisitionRepo{db: db}
}

func reqToGORM(r *domain.PurchaseRequisition) *domain.RequisitionGORM {
	return &domain.RequisitionGORM{
		ID:                r.ID,
		CompanyID:         r.CompanyID,
		RequisitionNumber: r.RequisitionNumber,
		RequesterID:       r.RequesterID,
		RequesterName:     r.RequesterName,
		DepartmentID:      r.DepartmentID,
		NeedByDate:        r.NeedByDate,
		Priority:          r.Priority,
		Reason:            r.Reason,
		Status:            string(r.Status),
		TotalEstimated:    r.TotalEstimated,
		ApprovedBy:        r.ApprovedBy,
		ApprovedAt:        r.ApprovedAt,
		RejectedReason:    r.RejectedReason,
		CreatedBy:         r.CreatedBy,
		CreatedAt:         r.CreatedAt,
		UpdatedAt:         r.UpdatedAt,
	}
}

func reqFromGORM(g *domain.RequisitionGORM) *domain.PurchaseRequisition {
	return &domain.PurchaseRequisition{
		ID:                g.ID,
		CompanyID:         g.CompanyID,
		RequisitionNumber: g.RequisitionNumber,
		RequesterID:       g.RequesterID,
		RequesterName:     g.RequesterName,
		DepartmentID:      g.DepartmentID,
		NeedByDate:        g.NeedByDate,
		Priority:          g.Priority,
		Reason:            g.Reason,
		Status:            domain.RequisitionStatus(g.Status),
		TotalEstimated:    g.TotalEstimated,
		ApprovedBy:        g.ApprovedBy,
		ApprovedAt:        g.ApprovedAt,
		RejectedReason:    g.RejectedReason,
		CreatedBy:         g.CreatedBy,
		CreatedAt:         g.CreatedAt,
		UpdatedAt:         g.UpdatedAt,
	}
}

func reqLineToGORM(l *domain.RequisitionItem) *domain.RequisitionItemGORM {
	return &domain.RequisitionItemGORM{
		ID:             l.ID,
		RequisitionID:  l.RequisitionID,
		LineNumber:     l.LineNumber,
		ItemCode:       l.ItemCode,
		ItemName:       l.ItemName,
		Unit:           l.Unit,
		Quantity:       l.Quantity,
		EstimatedPrice: l.EstimatedPrice,
		EstimatedTotal: l.EstimatedTotal,
		AccountID:      l.AccountID,
	}
}

func reqLineFromGORM(g *domain.RequisitionItemGORM) *domain.RequisitionItem {
	return &domain.RequisitionItem{
		ID:             g.ID,
		RequisitionID:  g.RequisitionID,
		LineNumber:     g.LineNumber,
		ItemCode:       g.ItemCode,
		ItemName:       g.ItemName,
		Unit:           g.Unit,
		Quantity:       g.Quantity,
		EstimatedPrice: g.EstimatedPrice,
		EstimatedTotal: g.EstimatedTotal,
		AccountID:      g.AccountID,
	}
}

func (r *PGRequisitionRepo) CreateRequisition(ctx context.Context, req *domain.PurchaseRequisition) error {
	return r.db.WithContext(ctx).Create(reqToGORM(req)).Error
}

func (r *PGRequisitionRepo) CreateRequisitionLines(ctx context.Context, lines []domain.RequisitionItem) error {
	gs := make([]domain.RequisitionItemGORM, len(lines))
	for i := range lines {
		gs[i] = *reqLineToGORM(&lines[i])
	}
	return r.db.WithContext(ctx).Create(&gs).Error
}

func (r *PGRequisitionRepo) GetRequisition(ctx context.Context, id string) (*domain.PurchaseRequisition, error) {
	var g domain.RequisitionGORM
	if err := r.db.WithContext(ctx).First(&g, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrRequisitionNotFound
		}
		return nil, err
	}
	req := reqFromGORM(&g)
	lines, err := r.GetRequisitionLines(ctx, id)
	if err != nil {
		return nil, err
	}
	req.Lines = lines
	return req, nil
}

func (r *PGRequisitionRepo) GetRequisitionLines(ctx context.Context, requisitionID string) ([]domain.RequisitionItem, error) {
	var gs []domain.RequisitionItemGORM
	if err := r.db.WithContext(ctx).Where("requisition_id = ?", requisitionID).Order("line_number").Find(&gs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.RequisitionItem, len(gs))
	for i := range gs {
		out[i] = *reqLineFromGORM(&gs[i])
	}
	return out, nil
}

func (r *PGRequisitionRepo) ListRequisitions(ctx context.Context, filter domain.RequisitionFilter) ([]domain.PurchaseRequisition, int, error) {
	q := r.db.WithContext(ctx).Model(&domain.RequisitionGORM{}).Where("company_id = ?", filter.CompanyID)
	if filter.Status != "" {
		q = q.Where("status = ?", string(filter.Status))
	}
	if filter.RequesterID != "" {
		q = q.Where("requester_id = ?", filter.RequesterID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var gs []domain.RequisitionGORM
	if err := q.Order("created_at DESC").Limit(filter.Limit).Offset(filter.Offset).Find(&gs).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.PurchaseRequisition, len(gs))
	for i := range gs {
		out[i] = *reqFromGORM(&gs[i])
	}
	return out, int(total), nil
}

func (r *PGRequisitionRepo) UpdateRequisition(ctx context.Context, req *domain.PurchaseRequisition) error {
	g := reqToGORM(req)
	return r.db.WithContext(ctx).Model(&domain.RequisitionGORM{}).
		Where("id = ?", req.ID).
		Select("requisition_number", "requester_id", "requester_name", "department_id", "need_by_date", "priority", "reason", "total_estimated", "updated_at").
		Updates(g).Error
}

func (r *PGRequisitionRepo) UpdateRequisitionStatus(ctx context.Context, id string, status domain.RequisitionStatus, approvedBy string, approvedAt time.Time) error {
	updates := map[string]interface{}{"status": string(status), "updated_at": time.Now()}
	if approvedBy != "" {
		updates["approved_by"] = approvedBy
	}
	if !approvedAt.IsZero() {
		updates["approved_at"] = approvedAt
	}
	return r.db.WithContext(ctx).Model(&domain.RequisitionGORM{}).Where("id = ?", id).Updates(updates).Error
}

func (r *PGRequisitionRepo) RejectRequisition(ctx context.Context, id, reason string) error {
	return r.db.WithContext(ctx).Model(&domain.RequisitionGORM{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"status": string(domain.ReqRejected), "rejected_reason": reason, "updated_at": time.Now()}).Error
}

func (r *PGRequisitionRepo) DeleteRequisition(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.RequisitionGORM{}).Error
}
