package repository

import (
	"context"
	"fmt"
	"gotax/internal/domain"
	"sort"
	"sync"
	"time"
)

var purchaseSeq int64

func purchaseID(prefix string) string {
	purchaseSeq++
	return prefix + time.Now().Format("20060102150405") + fmt.Sprintf("%03d", purchaseSeq%1000)
}

type MemoryPurchaseRepo struct {
	mu            sync.RWMutex
	suppliers     map[string]*domain.Supplier
	supByCode     map[string]map[string]*domain.Supplier
	pos           map[string]*domain.PurchaseOrder
	poByNumber    map[string]map[string]*domain.PurchaseOrder
	poLines       map[string][]domain.POItem
	grns          map[string]*domain.GRN
	grnByNumber   map[string]map[string]*domain.GRN
	grnLines      map[string][]domain.GRNItem
	invoices      map[string]*domain.SupplierInvoice
	invByNumber   map[string]map[string]*domain.SupplierInvoice
	invLines      map[string][]domain.SupplierInvoiceLine
	apTxns        map[string]*domain.APTransaction
	costAllocs    map[string]*domain.CostAllocation
	provisions    map[string]*domain.DoubtfulDebtProvision
	provLines     map[string][]domain.DoubtfulDebtProvisionLine
	reqs          map[string]*domain.PurchaseRequisition
	reqLines      map[string][]domain.RequisitionItem
}

func NewMemoryPurchaseRepo() *MemoryPurchaseRepo {
	return &MemoryPurchaseRepo{
		suppliers:   make(map[string]*domain.Supplier),
		supByCode:   make(map[string]map[string]*domain.Supplier),
		pos:         make(map[string]*domain.PurchaseOrder),
		poByNumber:  make(map[string]map[string]*domain.PurchaseOrder),
		poLines:     make(map[string][]domain.POItem),
		grns:        make(map[string]*domain.GRN),
		grnByNumber: make(map[string]map[string]*domain.GRN),
		grnLines:    make(map[string][]domain.GRNItem),
		invoices:    make(map[string]*domain.SupplierInvoice),
		invByNumber: make(map[string]map[string]*domain.SupplierInvoice),
		invLines:    make(map[string][]domain.SupplierInvoiceLine),
		apTxns:      make(map[string]*domain.APTransaction),
		costAllocs:  make(map[string]*domain.CostAllocation),
		provisions:  make(map[string]*domain.DoubtfulDebtProvision),
		provLines:   make(map[string][]domain.DoubtfulDebtProvisionLine),
		reqs:        make(map[string]*domain.PurchaseRequisition),
		reqLines:    make(map[string][]domain.RequisitionItem),
	}
}

// ─── Supplier ────────────────────────────────────────────────────────────

func (r *MemoryPurchaseRepo) CreateSupplier(_ context.Context, s *domain.Supplier) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *s
	if cp.ID == "" {
		cp.ID = purchaseID("SUP-")
	}
	if cp.Currency == "" {
		cp.Currency = "VND"
	}
	if cp.Status == "" {
		cp.Status = domain.SupplierActive
	}
	now := time.Now()
	cp.CreatedAt, cp.UpdatedAt = now, now
	if r.supByCode[cp.CompanyID] == nil {
		r.supByCode[cp.CompanyID] = make(map[string]*domain.Supplier)
	}
	if prev, ok := r.supByCode[cp.CompanyID][cp.Code]; ok && prev != nil {
		return domain.ErrSupplierCodeExists
	}
	r.suppliers[cp.ID] = &cp
	r.supByCode[cp.CompanyID][cp.Code] = r.suppliers[cp.ID]
	s.ID = cp.ID
	return nil
}

func (r *MemoryPurchaseRepo) GetSupplier(_ context.Context, id string) (*domain.Supplier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.suppliers[id]
	if !ok {
		return nil, domain.ErrSupplierNotFound
	}
	cp := *s
	return &cp, nil
}

func (r *MemoryPurchaseRepo) GetSupplierByCode(_ context.Context, companyID, code string) (*domain.Supplier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cm, ok := r.supByCode[companyID]
	if !ok {
		return nil, domain.ErrSupplierNotFound
	}
	s, ok := cm[code]
	if !ok {
		return nil, domain.ErrSupplierNotFound
	}
	cp := *s
	return &cp, nil
}

func (r *MemoryPurchaseRepo) ListSuppliers(_ context.Context, filter domain.PurchaseOrderFilter) ([]domain.Supplier, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.Supplier
	for _, s := range r.suppliers {
		if filter.CompanyID != "" && s.CompanyID != filter.CompanyID {
			continue
		}
		if filter.SupplierID != "" && s.ID != filter.SupplierID {
			continue
		}
		out = append(out, *s)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Code < out[j].Code })
	return out, len(out), nil
}

func (r *MemoryPurchaseRepo) UpdateSupplier(_ context.Context, s *domain.Supplier) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.suppliers[s.ID]
	if !ok {
		return domain.ErrSupplierNotFound
	}
	if r.supByCode[s.CompanyID] != nil {
		if other, exists := r.supByCode[s.CompanyID][s.Code]; exists && other.ID != s.ID {
			return domain.ErrSupplierCodeExists
		}
	}
	cp := *s
	cp.UpdatedAt = time.Now()
	r.suppliers[s.ID] = &cp
	if r.supByCode[existing.CompanyID] != nil {
		delete(r.supByCode[existing.CompanyID], existing.Code)
	}
	if r.supByCode[cp.CompanyID] == nil {
		r.supByCode[cp.CompanyID] = make(map[string]*domain.Supplier)
	}
	r.supByCode[cp.CompanyID][cp.Code] = r.suppliers[s.ID]
	return nil
}

func (r *MemoryPurchaseRepo) DeleteSupplier(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, ok := r.suppliers[id]
	if !ok {
		return domain.ErrSupplierNotFound
	}
	delete(r.suppliers, id)
	if r.supByCode[s.CompanyID] != nil {
		delete(r.supByCode[s.CompanyID], s.Code)
	}
	return nil
}

func (r *MemoryPurchaseRepo) ListSuppliersByIDs(_ context.Context, ids []string) ([]domain.Supplier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]domain.Supplier, 0, len(ids))
	for _, id := range ids {
		if s, ok := r.suppliers[id]; ok {
			out = append(out, *s)
		}
	}
	return out, nil
}

// ─── Purchase Order ──────────────────────────────────────────────────────

func (r *MemoryPurchaseRepo) CreatePO(_ context.Context, po *domain.PurchaseOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *po
	if cp.ID == "" {
		cp.ID = purchaseID("PO-")
	}
	if cp.Status == "" {
		cp.Status = domain.POStatusDraft
	}
	now := time.Now()
	cp.CreatedAt, cp.UpdatedAt = now, now
	if r.poByNumber[cp.CompanyID] == nil {
		r.poByNumber[cp.CompanyID] = make(map[string]*domain.PurchaseOrder)
	}
	r.pos[cp.ID] = &cp
	r.poByNumber[cp.CompanyID][cp.PONumber] = r.pos[cp.ID]
	po.ID = cp.ID
	return nil
}

func (r *MemoryPurchaseRepo) GetPO(_ context.Context, id string) (*domain.PurchaseOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	po, ok := r.pos[id]
	if !ok {
		return nil, domain.ErrPONotFound
	}
	cp := *po
	if lines, has := r.poLines[id]; has {
		cp.Lines = make([]domain.POItem, len(lines))
		copy(cp.Lines, lines)
	}
	return &cp, nil
}

func (r *MemoryPurchaseRepo) GetPOByNumber(_ context.Context, companyID, poNumber string) (*domain.PurchaseOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cm, ok := r.poByNumber[companyID]
	if !ok {
		return nil, domain.ErrPONotFound
	}
	po, ok := cm[poNumber]
	if !ok {
		return nil, domain.ErrPONotFound
	}
	return r.loadPO(po)
}

func (r *MemoryPurchaseRepo) ListPOs(_ context.Context, filter domain.PurchaseOrderFilter) ([]domain.PurchaseOrder, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.PurchaseOrder
	for _, po := range r.pos {
		if filter.CompanyID != "" && po.CompanyID != filter.CompanyID {
			continue
		}
		if filter.SupplierID != "" && po.SupplierID != filter.SupplierID {
			continue
		}
		if filter.Status != "" && po.Status != filter.Status {
			continue
		}
		if filter.FromDate != "" && po.OrderDate.Format("2006-01-02") < filter.FromDate {
			continue
		}
		if filter.ToDate != "" && po.OrderDate.Format("2006-01-02") > filter.ToDate {
			continue
		}
		loaded, _ := r.loadPO(po)
		if loaded != nil {
			out = append(out, *loaded)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].OrderDate.After(out[j].OrderDate) })
	total := len(out)
	if filter.Offset > 0 && filter.Offset < total {
		out = out[filter.Offset:]
	}
	if filter.Limit > 0 && filter.Limit < len(out) {
		out = out[:filter.Limit]
	}
	return out, total, nil
}

func (r *MemoryPurchaseRepo) UpdatePO(_ context.Context, po *domain.PurchaseOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.pos[po.ID]
	if !ok {
		return domain.ErrPONotFound
	}
	if existing.Status != domain.POStatusDraft {
		return domain.ErrPOCannotUpdate
	}
	cp := *po
	cp.UpdatedAt = time.Now()
	r.pos[po.ID] = &cp
	if r.poByNumber[cp.CompanyID] == nil {
		r.poByNumber[cp.CompanyID] = make(map[string]*domain.PurchaseOrder)
	}
	r.poByNumber[cp.CompanyID][cp.PONumber] = r.pos[po.ID]
	if len(cp.Lines) > 0 {
		r.poLines[po.ID] = make([]domain.POItem, len(cp.Lines))
		copy(r.poLines[po.ID], cp.Lines)
	}
	return nil
}

func (r *MemoryPurchaseRepo) UpdatePOStatus(_ context.Context, id string, status domain.POStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	po, ok := r.pos[id]
	if !ok {
		return domain.ErrPONotFound
	}
	po.Status = status
	po.UpdatedAt = time.Now()
	return nil
}

func (r *MemoryPurchaseRepo) ApprovePO(_ context.Context, id string, approvedBy string, approvedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	po, ok := r.pos[id]
	if !ok {
		return domain.ErrPONotFound
	}
	po.Status = domain.POStatusApproved
	po.ApprovedBy = approvedBy
	po.ApprovedAt = &approvedAt
	po.UpdatedAt = approvedAt
	return nil
}

func (r *MemoryPurchaseRepo) CancelPO(_ context.Context, id string, cancelReason string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	po, ok := r.pos[id]
	if !ok {
		return domain.ErrPONotFound
	}
	po.Status = domain.POStatusCancelled
	po.CancelledReason = cancelReason
	po.UpdatedAt = time.Now()
	return nil
}

func (r *MemoryPurchaseRepo) GetPOLines(_ context.Context, poID string) ([]domain.POItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	lines, ok := r.poLines[poID]
	if !ok {
		return nil, nil
	}
	out := make([]domain.POItem, len(lines))
	copy(out, lines)
	return out, nil
}

func (r *MemoryPurchaseRepo) CreatePOLines(_ context.Context, items []domain.POItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(items) == 0 {
		return nil
	}
	poID := items[0].POID
	total := len(r.poLines[poID]) + len(items)
	next := make([]domain.POItem, total)
	copy(next, r.poLines[poID])
	for i, it := range items {
		cp := it
		if cp.ID == "" {
			cp.ID = purchaseID("POLI-")
		}
		next[len(r.poLines[poID])+i] = cp
	}
	r.poLines[poID] = next
	return nil
}

func (r *MemoryPurchaseRepo) UpdatePOLines(_ context.Context, items []domain.POItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(items) == 0 {
		return nil
	}
	poID := items[0].POID
	r.poLines[poID] = make([]domain.POItem, len(items))
	copy(r.poLines[poID], items)
	return nil
}

func (r *MemoryPurchaseRepo) NextPONumber(_ context.Context, companyID, yyyymm string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cm := r.poByNumber[companyID]
	if cm == nil {
		return fmt.Sprintf("PO-%s-0001", yyyymm), nil
	}
	max := 0
	for num := range cm {
		var sfx int
		if _, err := fmt.Sscanf(num, "PO-"+yyyymm+"-%d", &sfx); err == nil && sfx > max {
			max = sfx
		}
	}
	return fmt.Sprintf("PO-%s-%04d", yyyymm, max+1), nil
}

// ─── GRN ─────────────────────────────────────────────────────────────────

func (r *MemoryPurchaseRepo) CreateGRN(_ context.Context, g *domain.GRN) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *g
	if cp.ID == "" {
		cp.ID = purchaseID("GRN-")
	}
	if cp.Status == "" {
		cp.Status = domain.GRNDraft
	}
	now := time.Now()
	cp.CreatedAt = now
	if r.grnByNumber[cp.CompanyID] == nil {
		r.grnByNumber[cp.CompanyID] = make(map[string]*domain.GRN)
	}
	r.grns[cp.ID] = &cp
	r.grnByNumber[cp.CompanyID][cp.GRNNumber] = r.grns[cp.ID]
	g.ID = cp.ID
	return nil
}

func (r *MemoryPurchaseRepo) GetGRN(_ context.Context, id string) (*domain.GRN, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.grns[id]
	if !ok {
		return nil, domain.ErrGRNNotFound
	}
	return r.loadGRN(g)
}

func (r *MemoryPurchaseRepo) GetGRNByNumber(_ context.Context, companyID, grnNumber string) (*domain.GRN, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cm, ok := r.grnByNumber[companyID]
	if !ok {
		return nil, domain.ErrGRNNotFound
	}
	g, ok := cm[grnNumber]
	if !ok {
		return nil, domain.ErrGRNNotFound
	}
	return r.loadGRN(g)
}

func (r *MemoryPurchaseRepo) ListGRNs(_ context.Context, filter domain.GRNFilter) ([]domain.GRN, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.GRN
	for _, g := range r.grns {
		if filter.CompanyID != "" && g.CompanyID != filter.CompanyID {
			continue
		}
		if filter.POID != "" && g.POID != filter.POID {
			continue
		}
		if filter.ReturnOfGRNID != "" && g.ReturnOfGRNID != filter.ReturnOfGRNID {
			continue
		}
		if filter.Status != "" && g.Status != filter.Status {
			continue
		}
		if filter.FromDate != "" && g.ReceiptDate.Format("2006-01-02") < filter.FromDate {
			continue
		}
		if filter.ToDate != "" && g.ReceiptDate.Format("2006-01-02") > filter.ToDate {
			continue
		}
		loaded, _ := r.loadGRN(g)
		if loaded != nil {
			out = append(out, *loaded)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ReceiptDate.After(out[j].ReceiptDate) })
	total := len(out)
	if filter.Offset > 0 && filter.Offset < total {
		out = out[filter.Offset:]
	}
	if filter.Limit > 0 && filter.Limit < len(out) {
		out = out[:filter.Limit]
	}
	return out, total, nil
}

func (r *MemoryPurchaseRepo) UpdateGRN(_ context.Context, g *domain.GRN) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.grns[g.ID]; !ok {
		return domain.ErrGRNNotFound
	}
	cp := *g
	cp.POID = g.POID
	cp.UpdatedAt = time.Now()
	r.grns[g.ID] = &cp
	if len(cp.Lines) > 0 {
		r.grnLines[g.ID] = make([]domain.GRNItem, len(cp.Lines))
		copy(r.grnLines[g.ID], cp.Lines)
	}
	return nil
}

func (r *MemoryPurchaseRepo) UpdateGRNStatus(_ context.Context, id string, status domain.GRNStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	g, ok := r.grns[id]
	if !ok {
		return domain.ErrGRNNotFound
	}
	g.Status = status
	return nil
}

func (r *MemoryPurchaseRepo) GetGRNLines(_ context.Context, grnID string) ([]domain.GRNItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	lines, ok := r.grnLines[grnID]
	if !ok {
		return nil, nil
	}
	out := make([]domain.GRNItem, len(lines))
	copy(out, lines)
	return out, nil
}

func (r *MemoryPurchaseRepo) CreateGRNLines(_ context.Context, items []domain.GRNItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(items) == 0 {
		return nil
	}
	grnID := items[0].GRNID
	total := len(r.grnLines[grnID]) + len(items)
	next := make([]domain.GRNItem, total)
	copy(next, r.grnLines[grnID])
	for i, it := range items {
		cp := it
		if cp.ID == "" {
			cp.ID = purchaseID("GRNL-")
		}
		next[len(r.grnLines[grnID])+i] = cp
	}
	r.grnLines[grnID] = next
	return nil
}

func (r *MemoryPurchaseRepo) UpdateGRNLines(_ context.Context, items []domain.GRNItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(items) == 0 {
		return nil
	}
	grnID := items[0].GRNID
	r.grnLines[grnID] = make([]domain.GRNItem, len(items))
	copy(r.grnLines[grnID], items)
	return nil
}

func (r *MemoryPurchaseRepo) NextGRNNumber(_ context.Context, companyID, yyyymm string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cm := r.grnByNumber[companyID]
	if cm == nil {
		return fmt.Sprintf("GRN-%s-0001", yyyymm), nil
	}
	max := 0
	for num := range cm {
		var sfx int
		if _, err := fmt.Sscanf(num, "GRN-"+yyyymm+"-%d", &sfx); err == nil && sfx > max {
			max = sfx
		}
	}
	return fmt.Sprintf("GRN-%s-%04d", yyyymm, max+1), nil
}

// ─── Supplier Invoice ────────────────────────────────────────────────────

func (r *MemoryPurchaseRepo) CreateInvoice(_ context.Context, inv *domain.SupplierInvoice) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *inv
	if cp.ID == "" {
		cp.ID = purchaseID("INV-")
	}
	if cp.Status == "" {
		cp.Status = domain.InvoiceDraft
	}
	if cp.VATDeductionStatus == "" {
		cp.VATDeductionStatus = domain.VATPending
	}
	now := time.Now()
	cp.CreatedAt = now
	if r.invByNumber[cp.CompanyID] == nil {
		r.invByNumber[cp.CompanyID] = make(map[string]*domain.SupplierInvoice)
	}
	r.invoices[cp.ID] = &cp
	r.invByNumber[cp.CompanyID][cp.InvoiceNumber] = r.invoices[cp.ID]
	inv.ID = cp.ID
	return nil
}

func (r *MemoryPurchaseRepo) GetInvoice(_ context.Context, id string) (*domain.SupplierInvoice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	inv, ok := r.invoices[id]
	if !ok {
		return nil, domain.ErrSupplierInvoiceNotFound
	}
	return r.loadInvoice(inv)
}

func (r *MemoryPurchaseRepo) GetInvoiceByNumber(_ context.Context, companyID, invoiceNumber string) (*domain.SupplierInvoice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cm, ok := r.invByNumber[companyID]
	if !ok {
		return nil, domain.ErrSupplierInvoiceNotFound
	}
	inv, ok := cm[invoiceNumber]
	if !ok {
		return nil, domain.ErrSupplierInvoiceNotFound
	}
	return r.loadInvoice(inv)
}

func (r *MemoryPurchaseRepo) ListInvoices(_ context.Context, filter domain.SupplierInvoiceFilter) ([]domain.SupplierInvoice, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.SupplierInvoice
	for _, inv := range r.invoices {
		if filter.CompanyID != "" && inv.CompanyID != filter.CompanyID {
			continue
		}
		if filter.SupplierID != "" && inv.SupplierID != filter.SupplierID {
			continue
		}
		if filter.OriginalInvoiceID != "" && inv.OriginalInvoiceID != filter.OriginalInvoiceID {
			continue
		}
		if filter.Status != "" && inv.Status != filter.Status {
			continue
		}
		if filter.FromDate != "" && inv.InvoiceDate.Format("2006-01-02") < filter.FromDate {
			continue
		}
		if filter.ToDate != "" && inv.InvoiceDate.Format("2006-01-02") > filter.ToDate {
			continue
		}
		loaded, _ := r.loadInvoice(inv)
		if loaded != nil {
			out = append(out, *loaded)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].InvoiceDate.After(out[j].InvoiceDate) })
	total := len(out)
	if filter.Offset > 0 && filter.Offset < total {
		out = out[filter.Offset:]
	}
	if filter.Limit > 0 && filter.Limit < len(out) {
		out = out[:filter.Limit]
	}
	return out, total, nil
}

func (r *MemoryPurchaseRepo) UpdateInvoice(_ context.Context, inv *domain.SupplierInvoice) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.invoices[inv.ID]
	if !ok {
		return domain.ErrSupplierInvoiceNotFound
	}
	cp := *inv
	cp.GLPosted = existing.GLPosted
	cp.GLPostedAt = existing.GLPostedAt
	r.invoices[inv.ID] = &cp
	if len(cp.Lines) > 0 {
		r.invLines[inv.ID] = make([]domain.SupplierInvoiceLine, len(cp.Lines))
		copy(r.invLines[inv.ID], cp.Lines)
	}
	return nil
}

func (r *MemoryPurchaseRepo) UpdateInvoiceStatus(_ context.Context, id string, status domain.InvoiceStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	inv, ok := r.invoices[id]
	if !ok {
		return domain.ErrSupplierInvoiceNotFound
	}
	inv.Status = status
	return nil
}

func (r *MemoryPurchaseRepo) ReduceInvoiceBalance(_ context.Context, id string, amount float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	inv, ok := r.invoices[id]
	if !ok {
		return domain.ErrSupplierInvoiceNotFound
	}
	cp := *inv
	cp.BalanceDue -= amount
	cp.AmountPaid += amount
	r.invoices[id] = &cp
	return nil
}

func (r *MemoryPurchaseRepo) PostInvoice(_ context.Context, id string, postedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	inv, ok := r.invoices[id]
	if !ok {
		return domain.ErrSupplierInvoiceNotFound
	}
	inv.Status = domain.InvoicePosted
	inv.BalanceDue = inv.TotalAmount - inv.AmountPaid
	return nil
}

func (r *MemoryPurchaseRepo) SetInvoiceGLPosted(_ context.Context, id string, postedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	inv, ok := r.invoices[id]
	if !ok {
		return domain.ErrSupplierInvoiceNotFound
	}
	inv.GLPosted = true
	inv.GLPostedAt = &postedAt
	return nil
}

func (r *MemoryPurchaseRepo) GetInvoiceLines(_ context.Context, invoiceID string) ([]domain.SupplierInvoiceLine, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	lines, ok := r.invLines[invoiceID]
	if !ok {
		return nil, nil
	}
	out := make([]domain.SupplierInvoiceLine, len(lines))
	copy(out, lines)
	return out, nil
}

func (r *MemoryPurchaseRepo) CreateInvoiceLines(_ context.Context, items []domain.SupplierInvoiceLine) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(items) == 0 {
		return nil
	}
	invID := items[0].InvoiceID
	total := len(r.invLines[invID]) + len(items)
	next := make([]domain.SupplierInvoiceLine, total)
	copy(next, r.invLines[invID])
	for i, it := range items {
		cp := it
		if cp.ID == "" {
			cp.ID = purchaseID("INVL-")
		}
		next[len(r.invLines[invID])+i] = cp
	}
	r.invLines[invID] = next
	return nil
}

func (r *MemoryPurchaseRepo) UpdateInvoiceLines(_ context.Context, items []domain.SupplierInvoiceLine) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(items) == 0 {
		return nil
	}
	invID := items[0].InvoiceID
	r.invLines[invID] = make([]domain.SupplierInvoiceLine, len(items))
	copy(r.invLines[invID], items)
	return nil
}

// ─── AP Transaction ──────────────────────────────────────────────────────

func (r *MemoryPurchaseRepo) CreateAPTransaction(_ context.Context, t *domain.APTransaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *t
	if cp.ID == "" {
		cp.ID = purchaseID("APT-")
	}
	if cp.Currency == "" {
		cp.Currency = "VND"
	}
	now := time.Now()
	cp.CreatedAt = now
	r.apTxns[cp.ID] = &cp
	t.ID = cp.ID
	return nil
}

func (r *MemoryPurchaseRepo) GetAPTransaction(_ context.Context, id string) (*domain.APTransaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.apTxns[id]
	if !ok {
		return nil, domain.ErrAPTransNotFound
	}
	cp := *t
	return &cp, nil
}

func (r *MemoryPurchaseRepo) ListAPTransactionsBySupplier(_ context.Context, companyID, supplierID string) ([]domain.APTransaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]domain.APTransaction, 0)
	for _, t := range r.apTxns {
		if t.SupplierID != supplierID {
			continue
		}
		if companyID != "" {
			sup := r.suppliers[t.SupplierID]
			if sup == nil || sup.CompanyID != companyID {
				continue
			}
		}
		out = append(out, *t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].TransactionDate.After(out[j].TransactionDate) })
	return out, nil
}

func (r *MemoryPurchaseRepo) ListAPTransactions(_ context.Context, companyID string, offset, limit int) ([]domain.APTransaction, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var filtered []domain.APTransaction
	for _, t := range r.apTxns {
		if companyID != "" {
			sup := r.suppliers[t.SupplierID]
			if sup == nil || sup.CompanyID != companyID {
				continue
			}
		}
		filtered = append(filtered, *t)
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].TransactionDate.After(filtered[j].TransactionDate) })
	total := len(filtered)
	if offset > 0 && offset < total {
		filtered = filtered[offset:]
	}
	if limit > 0 && limit < len(filtered) {
		filtered = filtered[:limit]
	}
	return filtered, total, nil
}

// ─── Cost Allocation ─────────────────────────────────────────────────────

func (r *MemoryPurchaseRepo) CreateCostAllocation(_ context.Context, c *domain.CostAllocation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *c
	if cp.ID == "" {
		cp.ID = purchaseID("CA-")
	}
	r.costAllocs[cp.ID] = &cp
	c.ID = cp.ID
	return nil
}

func (r *MemoryPurchaseRepo) GetCostAllocation(_ context.Context, id string) (*domain.CostAllocation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.costAllocs[id]
	if !ok {
		return nil, domain.ErrAPTransNotFound
	}
	cp := *c
	return &cp, nil
}

func (r *MemoryPurchaseRepo) ListCostAllocationsByInvoice(_ context.Context, invoiceID string) ([]domain.CostAllocation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.CostAllocation
	for _, c := range r.costAllocs {
		if c.InvoiceID == invoiceID {
			out = append(out, *c)
		}
	}
	return out, nil
}

// ─── Doubtful Debt Provisions ───────────────────────────────────────────

func (r *MemoryPurchaseRepo) CreateProvision(_ context.Context, p *domain.DoubtfulDebtProvision) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *p
	if cp.ID == "" {
		cp.ID = purchaseID("PROV-")
	}
	cp.CreatedAt = time.Now()
	r.provisions[cp.ID] = &cp
	p.ID = cp.ID
	return nil
}

func (r *MemoryPurchaseRepo) CreateProvisionLines(_ context.Context, lines []domain.DoubtfulDebtProvisionLine) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range lines {
		cp := lines[i]
		if cp.ID == "" {
			cp.ID = purchaseID("PROVL-")
		}
		r.provLines[cp.ProvisionID] = append(r.provLines[cp.ProvisionID], cp)
		lines[i].ID = cp.ID
	}
	return nil
}

func (r *MemoryPurchaseRepo) GetProvision(_ context.Context, id string) (*domain.DoubtfulDebtProvision, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.provisions[id]
	if !ok {
		return nil, domain.ErrProvisionNotFound
	}
	cp := *p
	if lines, ok := r.provLines[id]; ok {
		cp.Lines = make([]domain.DoubtfulDebtProvisionLine, len(lines))
		copy(cp.Lines, lines)
	}
	return &cp, nil
}

func (r *MemoryPurchaseRepo) GetProvisionLines(_ context.Context, provisionID string) ([]domain.DoubtfulDebtProvisionLine, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if lines, ok := r.provLines[provisionID]; ok {
		out := make([]domain.DoubtfulDebtProvisionLine, len(lines))
		copy(out, lines)
		return out, nil
	}
	return []domain.DoubtfulDebtProvisionLine{}, nil
}

func (r *MemoryPurchaseRepo) ListProvisions(_ context.Context, companyID string, limit, offset int) ([]domain.DoubtfulDebtProvision, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var all []domain.DoubtfulDebtProvision
	for _, p := range r.provisions {
		if p.CompanyID == companyID {
			all = append(all, *p)
		}
	}
	sort.Slice(all, func(i, j int) bool { return all[i].CreatedAt.After(all[j].CreatedAt) })
	if offset > len(all) {
		offset = len(all)
	}
	if limit <= 0 {
		return all, len(all), nil
	}
	end := offset + limit
	if end > len(all) {
		end = len(all)
	}
	return all[offset:end], len(all), nil
}

// ─── Requisition ─────────────────────────────────────────────────────────

func (r *MemoryPurchaseRepo) CreateRequisition(_ context.Context, req *domain.PurchaseRequisition) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cr := *req
	if cr.ID == "" {
		cr.ID = purchaseID("REQ-")
	}
	cr.CreatedAt = time.Now()
	cr.UpdatedAt = cr.CreatedAt
	r.reqs[cr.ID] = &cr
	req.ID = cr.ID
	return nil
}

func (r *MemoryPurchaseRepo) CreateRequisitionLines(_ context.Context, lines []domain.RequisitionItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range lines {
		cl := lines[i]
		if cl.ID == "" {
			cl.ID = purchaseID("REQL-")
		}
		r.reqLines[cl.RequisitionID] = append(r.reqLines[cl.RequisitionID], cl)
		lines[i].ID = cl.ID
	}
	return nil
}

func (r *MemoryPurchaseRepo) GetRequisition(_ context.Context, id string) (*domain.PurchaseRequisition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	req, ok := r.reqs[id]
	if !ok {
		return nil, domain.ErrRequisitionNotFound
	}
	return r.loadRequisition(req)
}

func (r *MemoryPurchaseRepo) GetRequisitionLines(_ context.Context, requisitionID string) ([]domain.RequisitionItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if lines, ok := r.reqLines[requisitionID]; ok {
		out := make([]domain.RequisitionItem, len(lines))
		copy(out, lines)
		return out, nil
	}
	return []domain.RequisitionItem{}, nil
}

func (r *MemoryPurchaseRepo) ListRequisitions(_ context.Context, filter domain.RequisitionFilter) ([]domain.PurchaseRequisition, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var all []domain.PurchaseRequisition
	for _, req := range r.reqs {
		if req.CompanyID != filter.CompanyID {
			continue
		}
		if filter.Status != "" && req.Status != filter.Status {
			continue
		}
		if filter.RequesterID != "" && req.RequesterID != filter.RequesterID {
			continue
		}
		if filter.FromDate != nil && req.CreatedAt.Before(*filter.FromDate) {
			continue
		}
		if filter.ToDate != nil && req.CreatedAt.After(*filter.ToDate) {
			continue
		}
		loaded, err := r.loadRequisition(req)
		if err != nil {
			return nil, 0, err
		}
		all = append(all, *loaded)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].CreatedAt.After(all[j].CreatedAt) })
	total := len(all)
	offset := 0
	if filter.Offset > 0 {
		offset = filter.Offset
	}
	limit := filter.Limit
	if limit <= 0 {
		return all, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	if offset > total {
		offset = total
	}
	return all[offset:end], total, nil
}

func (r *MemoryPurchaseRepo) UpdateRequisition(_ context.Context, req *domain.PurchaseRequisition) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.reqs[req.ID]
	if !ok {
		return domain.ErrRequisitionNotFound
	}
	cr := *req
	cr.CreatedAt = existing.CreatedAt
	cr.UpdatedAt = time.Now()
	r.reqs[req.ID] = &cr
	return nil
}

func (r *MemoryPurchaseRepo) UpdateRequisitionStatus(_ context.Context, id string, status domain.RequisitionStatus, approvedBy string, approvedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	req, ok := r.reqs[id]
	if !ok {
		return domain.ErrRequisitionNotFound
	}
	cr := *req
	cr.Status = status
	if approvedBy != "" {
		cr.ApprovedBy = approvedBy
	}
	if !approvedAt.IsZero() {
		cr.ApprovedAt = &approvedAt
	}
	cr.UpdatedAt = time.Now()
	r.reqs[id] = &cr
	return nil
}

func (r *MemoryPurchaseRepo) RejectRequisition(_ context.Context, id, reason string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	req, ok := r.reqs[id]
	if !ok {
		return domain.ErrRequisitionNotFound
	}
	cr := *req
	cr.Status = domain.ReqRejected
	cr.RejectedReason = reason
	cr.UpdatedAt = time.Now()
	r.reqs[id] = &cr
	return nil
}

func (r *MemoryPurchaseRepo) DeleteRequisition(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.reqs[id]; !ok {
		return domain.ErrRequisitionNotFound
	}
	delete(r.reqs, id)
	delete(r.reqLines, id)
	return nil
}

// ─── helpers ─────────────────────────────────────────────────────────────

func (r *MemoryPurchaseRepo) loadRequisition(req *domain.PurchaseRequisition) (*domain.PurchaseRequisition, error) {
	cp := *req
	if lines, ok := r.reqLines[req.ID]; ok {
		cp.Lines = make([]domain.RequisitionItem, len(lines))
		copy(cp.Lines, lines)
	}
	return &cp, nil
}

func (r *MemoryPurchaseRepo) loadPO(po *domain.PurchaseOrder) (*domain.PurchaseOrder, error) {
	cp := *po
	if lines, ok := r.poLines[po.ID]; ok {
		cp.Lines = make([]domain.POItem, len(lines))
		copy(cp.Lines, lines)
	}
	return &cp, nil
}

func (r *MemoryPurchaseRepo) loadGRN(g *domain.GRN) (*domain.GRN, error) {
	cp := *g
	if lines, ok := r.grnLines[g.ID]; ok {
		cp.Lines = make([]domain.GRNItem, len(lines))
		copy(cp.Lines, lines)
	}
	return &cp, nil
}

func (r *MemoryPurchaseRepo) loadInvoice(inv *domain.SupplierInvoice) (*domain.SupplierInvoice, error) {
	cp := *inv
	if lines, ok := r.invLines[inv.ID]; ok {
		cp.Lines = make([]domain.SupplierInvoiceLine, len(lines))
		copy(cp.Lines, lines)
	}
	return &cp, nil
}
