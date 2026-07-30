package repository

import (
	"context"
	"fmt"
	"gotax/internal/domain"
	"strconv"
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
		TaxCode:           strPtr(s.TaxCode),
		Address:           strPtr(s.Address),
		Phone:             strPtr(s.Phone),
		Email:             strPtr(s.Email),
		BankAccountName:   strPtr(s.BankAccountName),
		BankAccountNumber: strPtr(s.BankAccountNumber),
		BankName:          strPtr(s.BankName),
		PaymentTerms:      strPtr(string(s.PaymentTerms)),
		CreditLimit:       float64Ptr(s.CreditLimit),
		IsActive:          s.Status == domain.SupplierActive,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}
}

func supplierFromGORM(g *domain.SupplierGORM) *domain.Supplier {
	s := &domain.Supplier{
		ID:                g.ID,
		CompanyID:         g.CompanyID,
		Code:              g.Code,
		Name:              g.Name,
		BankAccountName:   "",
		BankAccountNumber: "",
		BankName:          "",
		PaymentTerms:      "",
		Notes:             "",
		Status:            domain.SupplierActive,
		CreatedAt:         g.CreatedAt,
		UpdatedAt:         g.UpdatedAt,
		Currency:          "VND",
	}
	if g.TaxCode != nil {
		s.TaxCode = *g.TaxCode
	}
	if g.Address != nil {
		s.Address = *g.Address
	}
	if g.Phone != nil {
		s.Phone = *g.Phone
	}
	if g.Email != nil {
		s.Email = *g.Email
	}
	if g.BankAccountName != nil {
		s.BankAccountName = *g.BankAccountName
	}
	if g.BankAccountNumber != nil {
		s.BankAccountNumber = *g.BankAccountNumber
	}
	if g.BankName != nil {
		s.BankName = *g.BankName
	}
	if g.PaymentTerms != nil {
		s.PaymentTerms = domain.PaymentTerms(*g.PaymentTerms)
	}
	if g.CreditLimit != nil {
		s.CreditLimit = *g.CreditLimit
	}
	if !g.IsActive {
		s.Status = domain.SupplierSuspended
	}
	return s
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
		"is_active":           g.IsActive,
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
	g := &domain.PurchaseOrderGORM{
		ID:          po.ID,
		CompanyID:   po.CompanyID,
		PONumber:    po.PONumber,
		OrderDate:   po.OrderDate,
		SupplierID:  po.SupplierID,
		Subtotal:    po.Subtotal,
		TaxAmount:   po.TaxAmount,
		TotalAmount: po.TotalAmount,
		Currency:    po.Currency,
		Status:      string(po.Status),
		ApprovedBy:  strPtr(po.ApprovedBy),
		CreatedBy:   po.CreatedBy,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
		Notes:       strPtr(po.Notes),
		PaymentTerms: strPtr(po.PaymentTerms),
	}
	if po.ExpectedDate != nil {
		g.DeliveryDate = po.ExpectedDate
	}
	if po.ApprovedAt != nil {
		g.ApprovedAt = po.ApprovedAt
	}
	return g
}

func poFromGORM(g *domain.PurchaseOrderGORM) *domain.PurchaseOrder {
	po := &domain.PurchaseOrder{
		ID:          g.ID,
		CompanyID:   g.CompanyID,
		PONumber:    g.PONumber,
		OrderDate:   g.OrderDate,
		SupplierID:  g.SupplierID,
		Subtotal:    g.Subtotal,
		TaxAmount:   g.TaxAmount,
		TotalAmount: g.TotalAmount,
		Currency:    g.Currency,
		Status:      domain.POStatus(g.Status),
		CreatedBy:   g.CreatedBy,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}
	if g.ApprovedBy != nil {
		po.ApprovedBy = *g.ApprovedBy
	}
	if g.DeliveryDate != nil {
		po.ExpectedDate = g.DeliveryDate
	}
	if g.ApprovedAt != nil {
		po.ApprovedAt = g.ApprovedAt
	}
	if g.Notes != nil {
		po.Notes = *g.Notes
	}
	if g.PaymentTerms != nil {
		po.PaymentTerms = *g.PaymentTerms
	}
	return po
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
		"supplier_id":  g.SupplierID,
		"order_date":   g.OrderDate,
		"currency":     g.Currency,
		"delivery_date": g.DeliveryDate,
		"payment_terms": g.PaymentTerms,
		"notes":        g.Notes,
	}).Error
}

func (r *PGPurchaseOrderRepo) UpdatePOStatus(ctx context.Context, id string, status domain.POStatus) error {
	return r.db.WithContext(ctx).Model(&domain.PurchaseOrderGORM{}).Where("id = ?", id).Update("status", string(status)).Error
}

func (r *PGPurchaseOrderRepo) ApprovePO(ctx context.Context, id string, approvedBy string, approvedAt time.Time) error {
	return r.db.WithContext(ctx).Model(&domain.PurchaseOrderGORM{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      string(domain.POStatusApproved),
		"approved_by": approvedBy,
		"approved_at": approvedAt,
	}).Error
}

func (r *PGPurchaseOrderRepo) CancelPO(ctx context.Context, id string, cancelReason string) error {
	return r.db.WithContext(ctx).Model(&domain.PurchaseOrderGORM{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":   string(domain.POStatusCancelled),
		"cancelled_reason": cancelReason,
	}).Error
}

func poItemToGORM(l *domain.POItem) *domain.POItemGORM {
	disc := l.DiscountPct
	return &domain.POItemGORM{
		POID:       l.POID,
		LineNumber: l.LineNumber,
		ItemCode:   l.ItemCode,
		ItemName:   l.ItemName,
		Quantity:   l.Quantity,
		UnitPrice:  l.UnitPrice,
		LineTotal:  l.LineTotal,
		Discount:   &disc,
		TaxRate:    &l.VATRate,
	}
}

func poItemFromGORM(g *domain.POItemGORM) *domain.POItem {
	l := &domain.POItem{
		ID:        g.ID,
		POID:      g.POID,
		LineNumber: g.LineNumber,
		ItemCode:  g.ItemCode,
		ItemName:  g.ItemName,
		Unit:      "",
		Quantity:  g.Quantity,
		UnitPrice: g.UnitPrice,
		LineTotal: g.LineTotal,
		VATRate:   10,
		AccountID: "",
	}
	if g.Discount != nil {
		l.DiscountPct = *g.Discount
	}
	if g.TaxRate != nil {
		l.VATRate = *g.TaxRate
	}
	return l
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
	gg := &domain.GRNGORM{
		ID:        g.ID,
		CompanyID: g.CompanyID,
		GRNNumber: g.GRNNumber,
		GRNDate:   g.ReceiptDate,
		POID:      strPtr(g.POID),
		Status:    string(g.Status),
		CreatedBy: g.CreatedBy,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
		Notes:     strPtr(g.Notes),
	}
	return gg
}

func grnFromGORM(gg *domain.GRNGORM) *domain.GRN {
	g := &domain.GRN{
		ID:          gg.ID,
		CompanyID:   gg.CompanyID,
		GRNNumber:   gg.GRNNumber,
		ReceiptDate: gg.GRNDate,
		Status:      domain.GRNStatus(gg.Status),
		CreatedBy:   gg.CreatedBy,
		CreatedAt:   gg.CreatedAt,
		UpdatedAt:   gg.UpdatedAt,
	}
	if gg.POID != nil {
		g.POID = *gg.POID
	}
	if gg.Notes != nil {
		g.Notes = *gg.Notes
	}
	return g
}

func (r *PGGRNRepo) CreateGRN(ctx context.Context, g *domain.GRN) error {
	gg := grnToGORM(g)
	if err := r.db.WithContext(ctx).Create(gg).Error; err != nil {
		return err
	}
	g.ID = gg.ID
	g.CreatedAt = gg.CreatedAt
	g.UpdatedAt = gg.UpdatedAt
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
		q = q.Where("grn_date >= ?", filter.FromDate)
	}
	if filter.ToDate != "" {
		q = q.Where("grn_date <= ?", filter.ToDate)
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
		dq = dq.Where("grn_date >= ?", filter.FromDate)
	}
	if filter.ToDate != "" {
		dq = dq.Where("grn_date <= ?", filter.ToDate)
	}
	dq = dq.Order("grn_date DESC")
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
		"grn_date":      gg.GRNDate,
		"po_id":         gg.POID,
		"status":        gg.Status,
		"notes":         gg.Notes,
	}).Error
}

func (r *PGGRNRepo) UpdateGRNStatus(ctx context.Context, id string, status domain.GRNStatus) error {
	return r.db.WithContext(ctx).Model(&domain.GRNGORM{}).Where("id = ?", id).Update("status", string(status)).Error
}

func grnItemToGORM(l *domain.GRNItem) *domain.GRNItemGORM {
	return &domain.GRNItemGORM{
		GRNID:      l.GRNID,
		POLineID:   strPtr(l.POLineID),
		ItemCode:   l.ItemCode,
		ItemName:   l.ItemName,
		QtyShipped: l.QuantityReceived,
		UnitPrice:  l.UnitPrice,
		LineTotal:  l.LineTotal,
	}
}

func grnItemFromGORM(g *domain.GRNItemGORM) *domain.GRNItem {
	l := &domain.GRNItem{
		ID:               g.ID,
		GRNID:            g.GRNID,
		ItemCode:         g.ItemCode,
		ItemName:         g.ItemName,
		QuantityReceived: g.QtyShipped,
		UnitPrice:        g.UnitPrice,
		LineTotal:        g.LineTotal,
		Unit:             "",
	}
	if g.POLineID != nil {
		l.POLineID = *g.POLineID
	}
	return l
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
	g := &domain.SupplierInvoiceGORM{
		ID:            inv.ID,
		CompanyID:     inv.CompanyID,
		InvoiceNumber: inv.InvoiceNumber,
		InvoiceDate:   inv.InvoiceDate,
		SupplierID:    inv.SupplierID,
		Subtotal:      inv.Subtotal,
		TaxAmount:     inv.TaxAmount,
		TotalAmount:   inv.TotalAmount,
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

func sinvFromGORM(g *domain.SupplierInvoiceGORM) *domain.SupplierInvoice {
	inv := &domain.SupplierInvoice{
		ID:            g.ID,
		CompanyID:     g.CompanyID,
		InvoiceNumber: g.InvoiceNumber,
		InvoiceDate:   g.InvoiceDate,
		SupplierID:    g.SupplierID,
		Subtotal:      g.Subtotal,
		TaxAmount:     g.TaxAmount,
		TotalAmount:   g.TotalAmount,
		Currency:      g.Currency,
		Status:        domain.InvoiceStatus(g.Status),
		GLPosted:      g.GLPosted,
		CreatedBy:     g.CreatedBy,
		CreatedAt:     g.CreatedAt,
		AmountPaid:    0,
		BalanceDue:    g.TotalAmount,
		InvoiceType:   "",
		ExchangeRate:  1,
		SupplierName:  "",
		SupplierTaxCode: "",
	}
	if g.DueDate != nil {
		inv.DueDate = g.DueDate
	}
	if g.PostedAt != nil {
		inv.GLPostedAt = g.PostedAt
	}
	return inv
}

func (r *PGSupplierInvoiceRepo) CreateInvoice(ctx context.Context, inv *domain.SupplierInvoice) error {
	g := sinvToGORM(inv)
	if err := r.db.WithContext(ctx).Create(g).Error; err != nil {
		return err
	}
	inv.ID = g.ID
	inv.CreatedAt = g.CreatedAt
	inv.UpdatedAt = g.UpdatedAt
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
	return r.db.WithContext(ctx).Model(&domain.SupplierInvoiceGORM{}).Where("id = ?", inv.ID).Updates(map[string]interface{}{
		"invoice_date": inv.InvoiceDate,
		"subtotal":     inv.Subtotal,
		"tax_amount":   inv.TaxAmount,
		"grand_total":  inv.TotalAmount,
		"due_date":     inv.DueDate,
		"notes":        inv.Notes,
	}).Error
}

func (r *PGSupplierInvoiceRepo) UpdateInvoiceStatus(ctx context.Context, id string, status domain.InvoiceStatus) error {
	return r.db.WithContext(ctx).Model(&domain.SupplierInvoiceGORM{}).Where("id = ?", id).Update("status", string(status)).Error
}

func (r *PGSupplierInvoiceRepo) PostInvoice(ctx context.Context, id string, postedAt time.Time) error {
	return r.db.WithContext(ctx).Model(&domain.SupplierInvoiceGORM{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":    string(domain.InvoicePosted),
		"gl_posted": true,
		"posted_at": postedAt,
	}).Error
}

func sinvLineToGORM(l *domain.SupplierInvoiceLine) *domain.SupplierInvoiceLineGORM {
	return &domain.SupplierInvoiceLineGORM{
		InvoiceID:  l.InvoiceID,
		LineNumber: 0,
		ItemCode:   l.ItemCode,
		ItemName:   l.ItemName,
		Quantity:   l.Quantity,
		UnitPrice:  l.UnitPrice,
		LineTotal:  l.LineTotal,
		TaxRate:    float64Ptr(l.VATRate),
		TaxAmount:  float64Ptr(l.LineVATAmount),
	}
}

func sinvLineFromGORM(g *domain.SupplierInvoiceLineGORM) *domain.SupplierInvoiceLine {
	l := &domain.SupplierInvoiceLine{
		ID:        strconv.Itoa(int(g.ID)),
		InvoiceID: g.InvoiceID,
		ItemCode:  g.ItemCode,
		ItemName:  g.ItemName,
		Quantity:  g.Quantity,
		UnitPrice: g.UnitPrice,
		LineTotal: g.LineTotal,
		Unit:      "",
		AccountID: "",
	}
	if g.TaxRate != nil {
		l.VATRate = *g.TaxRate
	}
	if g.TaxAmount != nil {
		l.LineVATAmount = *g.TaxAmount
	}
	if g.AccountCode != nil {
		l.AccountID = *g.AccountCode
	}
	return l
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
		ID:        t.ID,
		CompanyID: t.CompanyID,
		SupplierID: t.SupplierID,
		Amount:    t.Amount,
		Currency:  t.Currency,
		Status:    "OPEN",
		CreatedBy: "",
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.CreatedAt,
	}
}

func aptFromGORM(g *domain.APTransactionGORM) *domain.APTransaction {
	return &domain.APTransaction{
		ID:              g.ID,
		CompanyID:       g.CompanyID,
		SupplierID:      g.SupplierID,
		TransactionType: domain.APTransInvoice,
		Amount:          g.Amount,
		Currency:        g.Currency,
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
		CostCenter: strPtr(c.InvoiceID),
		AllocPct:   100,
		AllocAmount: c.CostAmount,
	}
}

func caFromGORM(g *domain.CostAllocationGORM) *domain.CostAllocation {
	c := &domain.CostAllocation{
		CostAmount: g.AllocAmount,
	}
	if g.CostCenter != nil {
		c.InvoiceID = *g.CostCenter
	}
	return c
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
	if err := r.db.WithContext(ctx).Where("cost_center = ?", invoiceID).Find(&gs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.CostAllocation, len(gs))
	for i := range gs {
		out[i] = *caFromGORM(&gs[i])
	}
	return out, nil
}
