package repository

import (
	"context"
	"fmt"
	"gotax/internal/domain"
	"sort"
	"sync"
	"time"
)

var saleSeq int64

func saleID(prefix string) string {
	saleSeq++
	return prefix + time.Now().Format("20060102150405") + fmt.Sprintf("%03d", saleSeq%1000)
}

type MemorySaleRepo struct {
	mu            sync.RWMutex
	customers     map[string]*domain.Customer
	custByCode    map[string]map[string]*domain.Customer
	sos           map[string]*domain.SalesOrder
	soByNumber    map[string]map[string]*domain.SalesOrder
	soLines       map[string][]domain.SOLine
	dns           map[string]*domain.DeliveryNote
	dnByNumber    map[string]map[string]*domain.DeliveryNote
	dnLines       map[string][]domain.DNLine
	invoices      map[string]*domain.CustomerInvoice
	invByNumber   map[string]map[string]*domain.CustomerInvoice
	invLines      map[string][]domain.InvLine
	receipts      map[string]*domain.CustomerReceipt
	rcptByNumber  map[string]map[string]*domain.CustomerReceipt
	rcptAllocs    map[string][]domain.RcpAllocation
	cns           map[string]*domain.CreditNote
	cnByNumber    map[string]map[string]*domain.CreditNote
	cnLines       map[string][]domain.CNLine
	arTxns        map[string]*domain.ARTransaction
	sqs           map[string]*domain.SalesQuotation
}

func NewMemorySaleRepo() *MemorySaleRepo {
	return &MemorySaleRepo{
		customers:    make(map[string]*domain.Customer),
		custByCode:   make(map[string]map[string]*domain.Customer),
		sos:          make(map[string]*domain.SalesOrder),
		soByNumber:   make(map[string]map[string]*domain.SalesOrder),
		soLines:      make(map[string][]domain.SOLine),
		dns:          make(map[string]*domain.DeliveryNote),
		dnByNumber:   make(map[string]map[string]*domain.DeliveryNote),
		dnLines:      make(map[string][]domain.DNLine),
		invoices:     make(map[string]*domain.CustomerInvoice),
		invByNumber:  make(map[string]map[string]*domain.CustomerInvoice),
		invLines:     make(map[string][]domain.InvLine),
		receipts:     make(map[string]*domain.CustomerReceipt),
		rcptByNumber: make(map[string]map[string]*domain.CustomerReceipt),
		rcptAllocs:   make(map[string][]domain.RcpAllocation),
		cns:          make(map[string]*domain.CreditNote),
		cnByNumber:   make(map[string]map[string]*domain.CreditNote),
		cnLines:      make(map[string][]domain.CNLine),
		arTxns:       make(map[string]*domain.ARTransaction),
		sqs:          make(map[string]*domain.SalesQuotation),
	}
}

// ─── Customer ───────────────────────────────────────────────────────────

func (r *MemorySaleRepo) CreateCustomer(_ context.Context, c *domain.Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *c
	if cp.ID == "" {
		cp.ID = saleID("CUS-")
	}
	if cp.Status == "" {
		cp.Status = domain.CustomerActive
	}
	now := time.Now()
	cp.CreatedAt, cp.UpdatedAt = now, now
	if r.custByCode[cp.CompanyID] == nil {
		r.custByCode[cp.CompanyID] = make(map[string]*domain.Customer)
	}
	if prev, ok := r.custByCode[cp.CompanyID][cp.Code]; ok && prev != nil {
		return domain.ErrCustomerCodeExists
	}
	r.customers[cp.ID] = &cp
	r.custByCode[cp.CompanyID][cp.Code] = r.customers[cp.ID]
	c.ID = cp.ID
	return nil
}

func (r *MemorySaleRepo) GetCustomer(_ context.Context, id string) (*domain.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.customers[id]
	if !ok {
		return nil, domain.ErrCustomerNotFound
	}
	cp := *c
	return &cp, nil
}

func (r *MemorySaleRepo) GetCustomerByCode(_ context.Context, companyID, code string) (*domain.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cm, ok := r.custByCode[companyID]
	if !ok {
		return nil, domain.ErrCustomerNotFound
	}
	c, ok := cm[code]
	if !ok {
		return nil, domain.ErrCustomerNotFound
	}
	cp := *c
	return &cp, nil
}

func (r *MemorySaleRepo) ListCustomers(_ context.Context, companyID string) ([]domain.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.Customer
	for _, c := range r.customers {
		if c.CompanyID == companyID {
			out = append(out, *c)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Code < out[j].Code })
	return out, nil
}

func (r *MemorySaleRepo) UpdateCustomer(_ context.Context, c *domain.Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.customers[c.ID]
	if !ok {
		return domain.ErrCustomerNotFound
	}
	if r.custByCode[c.CompanyID] != nil {
		if other, exists := r.custByCode[c.CompanyID][c.Code]; exists && other.ID != c.ID {
			return domain.ErrCustomerCodeExists
		}
	}
	cp := *c
	cp.UpdatedAt = time.Now()
	r.customers[c.ID] = &cp
	if r.custByCode[existing.CompanyID] != nil {
		delete(r.custByCode[existing.CompanyID], existing.Code)
	}
	if r.custByCode[cp.CompanyID] == nil {
		r.custByCode[cp.CompanyID] = make(map[string]*domain.Customer)
	}
	r.custByCode[cp.CompanyID][cp.Code] = r.customers[c.ID]
	return nil
}

func (r *MemorySaleRepo) DeleteCustomer(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.customers[id]
	if !ok {
		return domain.ErrCustomerNotFound
	}
	delete(r.customers, id)
	if r.custByCode[c.CompanyID] != nil {
		delete(r.custByCode[c.CompanyID], c.Code)
	}
	return nil
}

// ─── Sales Order ────────────────────────────────────────────────────────

func (r *MemorySaleRepo) CreateSO(_ context.Context, so *domain.SalesOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *so
	if cp.ID == "" {
		cp.ID = saleID("SO-")
	}
	if cp.Status == "" {
		cp.Status = domain.SODraft
	}
	now := time.Now()
	cp.CreatedAt, cp.UpdatedAt = now, now
	if r.soByNumber[cp.CompanyID] == nil {
		r.soByNumber[cp.CompanyID] = make(map[string]*domain.SalesOrder)
	}
	r.sos[cp.ID] = &cp
	r.soByNumber[cp.CompanyID][cp.SONumber] = r.sos[cp.ID]
	if len(cp.Lines) > 0 {
		r.soLines[cp.ID] = make([]domain.SOLine, len(cp.Lines))
		copy(r.soLines[cp.ID], cp.Lines)
	}
	so.ID = cp.ID
	return nil
}

func (r *MemorySaleRepo) GetSO(_ context.Context, id string) (*domain.SalesOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	so, ok := r.sos[id]
	if !ok {
		return nil, domain.ErrSONotFound
	}
	return r.loadSO(so)
}

func (r *MemorySaleRepo) GetSOByNumber(_ context.Context, companyID, soNumber string) (*domain.SalesOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cm, ok := r.soByNumber[companyID]
	if !ok {
		return nil, domain.ErrSONotFound
	}
	so, ok := cm[soNumber]
	if !ok {
		return nil, domain.ErrSONotFound
	}
	return r.loadSO(so)
}

func (r *MemorySaleRepo) ListSOs(_ context.Context, filter domain.SalesOrderFilter) ([]domain.SalesOrder, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.SalesOrder
	for _, so := range r.sos {
		if filter.CompanyID != "" && so.CompanyID != filter.CompanyID {
			continue
		}
		if filter.CustomerID != "" && so.CustomerID != filter.CustomerID {
			continue
		}
		if filter.Status != "" && so.Status != filter.Status {
			continue
		}
		if filter.FromDate != "" && so.OrderDate.Format("2006-01-02") < filter.FromDate {
			continue
		}
		if filter.ToDate != "" && so.OrderDate.Format("2006-01-02") > filter.ToDate {
			continue
		}
		loaded, _ := r.loadSO(so)
		if loaded != nil {
			out = append(out, *loaded)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].OrderDate.After(out[j].OrderDate) })
	total := len(out)
	if filter.Offset > 0 && filter.Offset < len(out) {
		out = out[filter.Offset:]
	} else if filter.Offset >= len(out) {
		out = nil
	}
	if filter.Limit > 0 && filter.Limit < len(out) {
		out = out[:filter.Limit]
	}
	return out, total, nil
}

func (r *MemorySaleRepo) UpdateSO(_ context.Context, so *domain.SalesOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.sos[so.ID]
	if !ok {
		return domain.ErrSONotFound
	}
	if existing.Status != domain.SODraft {
		return domain.ErrSOCannotUpdate
	}
	cp := *so
	cp.UpdatedAt = time.Now()
	r.sos[so.ID] = &cp
	if r.soByNumber[cp.CompanyID] == nil {
		r.soByNumber[cp.CompanyID] = make(map[string]*domain.SalesOrder)
	}
	r.soByNumber[cp.CompanyID][cp.SONumber] = r.sos[so.ID]
	if len(cp.Lines) > 0 {
		r.soLines[so.ID] = make([]domain.SOLine, len(cp.Lines))
		copy(r.soLines[so.ID], cp.Lines)
	}
	return nil
}

func (r *MemorySaleRepo) UpdateSOStatus(_ context.Context, id string, status domain.SOStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	so, ok := r.sos[id]
	if !ok {
		return domain.ErrSONotFound
	}
	so.Status = status
	so.UpdatedAt = time.Now()
	return nil
}

func (r *MemorySaleRepo) ApproveSO(_ context.Context, id, approvedBy string, approvedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	so, ok := r.sos[id]
	if !ok {
		return domain.ErrSONotFound
	}
	if so.Status != domain.SODraft {
		return domain.ErrSOAlreadyApproved
	}
	so.Status = domain.SOApproved
	so.ApprovedBy = approvedBy
	so.ApprovedAt = &approvedAt
	so.UpdatedAt = time.Now()
	return nil
}

func (r *MemorySaleRepo) CancelSO(_ context.Context, id, cancelReason string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	so, ok := r.sos[id]
	if !ok {
		return domain.ErrSONotFound
	}
	if so.Status == domain.SOCancelled {
		return domain.ErrSOAlreadyCancelled
	}
	so.Status = domain.SOCancelled
	so.CancelledReason = cancelReason
	so.UpdatedAt = time.Now()
	return nil
}

func (r *MemorySaleRepo) GetSOLines(_ context.Context, soID string) ([]domain.SOLine, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	lines, ok := r.soLines[soID]
	if !ok {
		return nil, nil
	}
	out := make([]domain.SOLine, len(lines))
	copy(out, lines)
	return out, nil
}

func (r *MemorySaleRepo) CreateSOLines(_ context.Context, items []domain.SOLine) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(items) == 0 {
		return nil
	}
	soID := items[0].SOID
	r.soLines[soID] = make([]domain.SOLine, len(items))
	copy(r.soLines[soID], items)
	return nil
}

func (r *MemorySaleRepo) UpdateSOLines(_ context.Context, items []domain.SOLine) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(items) == 0 {
		return nil
	}
	soID := items[0].SOID
	r.soLines[soID] = make([]domain.SOLine, len(items))
	copy(r.soLines[soID], items)
	return nil
}

func (r *MemorySaleRepo) NextSONumber(_ context.Context, companyID, yyyymm string) (string, error) {
	return fmt.Sprintf("SO-%s-%05d", yyyymm, len(r.sos)+1), nil
}

// ─── Delivery Note ──────────────────────────────────────────────────────

func (r *MemorySaleRepo) CreateDN(_ context.Context, dn *domain.DeliveryNote) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *dn
	if cp.ID == "" {
		cp.ID = saleID("DN-")
	}
	if cp.Status == "" {
		cp.Status = domain.DNDraft
	}
	now := time.Now()
	cp.CreatedAt = now
	if r.dnByNumber[cp.CompanyID] == nil {
		r.dnByNumber[cp.CompanyID] = make(map[string]*domain.DeliveryNote)
	}
	r.dns[cp.ID] = &cp
	r.dnByNumber[cp.CompanyID][cp.DNNumber] = r.dns[cp.ID]
	if len(cp.Lines) > 0 {
		r.dnLines[cp.ID] = make([]domain.DNLine, len(cp.Lines))
		copy(r.dnLines[cp.ID], cp.Lines)
	}
	dn.ID = cp.ID
	return nil
}

func (r *MemorySaleRepo) GetDN(_ context.Context, id string) (*domain.DeliveryNote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	dn, ok := r.dns[id]
	if !ok {
		return nil, domain.ErrDNNotFound
	}
	return r.loadDN(dn)
}

func (r *MemorySaleRepo) GetDNByNumber(_ context.Context, companyID, dnNumber string) (*domain.DeliveryNote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cm, ok := r.dnByNumber[companyID]
	if !ok {
		return nil, domain.ErrDNNotFound
	}
	dn, ok := cm[dnNumber]
	if !ok {
		return nil, domain.ErrDNNotFound
	}
	return r.loadDN(dn)
}

func (r *MemorySaleRepo) ListDNs(_ context.Context, filter domain.DeliveryNoteFilter) ([]domain.DeliveryNote, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.DeliveryNote
	for _, dn := range r.dns {
		if filter.CompanyID != "" && dn.CompanyID != filter.CompanyID {
			continue
		}
		if filter.SOID != "" && dn.SOID != filter.SOID {
			continue
		}
		if filter.Status != "" && dn.Status != filter.Status {
			continue
		}
		if filter.FromDate != "" && dn.DeliveryDate.Format("2006-01-02") < filter.FromDate {
			continue
		}
		if filter.ToDate != "" && dn.DeliveryDate.Format("2006-01-02") > filter.ToDate {
			continue
		}
		loaded, _ := r.loadDN(dn)
		if loaded != nil {
			out = append(out, *loaded)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].DeliveryDate.After(out[j].DeliveryDate) })
	total := len(out)
	if filter.Offset > 0 && filter.Offset < len(out) {
		out = out[filter.Offset:]
	} else if filter.Offset >= len(out) {
		out = nil
	}
	if filter.Limit > 0 && filter.Limit < len(out) {
		out = out[:filter.Limit]
	}
	return out, total, nil
}

func (r *MemorySaleRepo) UpdateDN(_ context.Context, dn *domain.DeliveryNote) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.dns[dn.ID]
	if !ok {
		return domain.ErrDNNotFound
	}
	if existing.Status != domain.DNDraft {
		return domain.ErrDNCannotUpdate
	}
	cp := *dn
	r.dns[dn.ID] = &cp
	if r.dnByNumber[cp.CompanyID] == nil {
		r.dnByNumber[cp.CompanyID] = make(map[string]*domain.DeliveryNote)
	}
	r.dnByNumber[cp.CompanyID][cp.DNNumber] = r.dns[dn.ID]
	if len(cp.Lines) > 0 {
		r.dnLines[dn.ID] = make([]domain.DNLine, len(cp.Lines))
		copy(r.dnLines[dn.ID], cp.Lines)
	}
	return nil
}

func (r *MemorySaleRepo) UpdateDNStatus(_ context.Context, id string, status domain.DNStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	dn, ok := r.dns[id]
	if !ok {
		return domain.ErrDNNotFound
	}
	if status == domain.DNPosted && dn.Status != domain.DNDraft {
		return domain.ErrDNInvalidTransition
	}
	dn.Status = status
	return nil
}

func (r *MemorySaleRepo) GetDNLines(_ context.Context, dnID string) ([]domain.DNLine, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	lines, ok := r.dnLines[dnID]
	if !ok {
		return nil, nil
	}
	out := make([]domain.DNLine, len(lines))
	copy(out, lines)
	return out, nil
}

func (r *MemorySaleRepo) CreateDNLines(_ context.Context, items []domain.DNLine) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(items) == 0 {
		return nil
	}
	dnID := items[0].DNID
	r.dnLines[dnID] = make([]domain.DNLine, len(items))
	copy(r.dnLines[dnID], items)
	return nil
}

func (r *MemorySaleRepo) UpdateDNLines(_ context.Context, items []domain.DNLine) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(items) == 0 {
		return nil
	}
	dnID := items[0].DNID
	r.dnLines[dnID] = make([]domain.DNLine, len(items))
	copy(r.dnLines[dnID], items)
	return nil
}

func (r *MemorySaleRepo) NextDNNumber(_ context.Context, companyID, yyyymm string) (string, error) {
	return fmt.Sprintf("DN-%s-%05d", yyyymm, len(r.dns)+1), nil
}

func (r *MemorySaleRepo) NextInvNumber(_ context.Context, companyID, yyyymm string) (string, error) {
	return fmt.Sprintf("INV-%s-%05d", yyyymm, len(r.invoices)+1), nil
}

// ─── Customer Invoice ───────────────────────────────────────────────────

func (r *MemorySaleRepo) CreateInvoice(_ context.Context, inv *domain.CustomerInvoice) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *inv
	if cp.ID == "" {
		cp.ID = saleID("INV-")
	}
	if cp.Status == "" {
		cp.Status = domain.SInvDraft
	}
	now := time.Now()
	cp.CreatedAt = now
	if r.invByNumber[cp.CompanyID] == nil {
		r.invByNumber[cp.CompanyID] = make(map[string]*domain.CustomerInvoice)
	}
	r.invoices[cp.ID] = &cp
	r.invByNumber[cp.CompanyID][cp.InvoiceNumber] = r.invoices[cp.ID]
	if len(cp.Lines) > 0 {
		r.invLines[cp.ID] = make([]domain.InvLine, len(cp.Lines))
		copy(r.invLines[cp.ID], cp.Lines)
	}
	inv.ID = cp.ID
	return nil
}

func (r *MemorySaleRepo) GetInvoice(_ context.Context, id string) (*domain.CustomerInvoice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	inv, ok := r.invoices[id]
	if !ok {
		return nil, domain.ErrInvNotFound
	}
	return r.loadInvoice(inv)
}

func (r *MemorySaleRepo) GetInvoiceByNumber(_ context.Context, companyID, invoiceNumber string) (*domain.CustomerInvoice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cm, ok := r.invByNumber[companyID]
	if !ok {
		return nil, domain.ErrInvNotFound
	}
	inv, ok := cm[invoiceNumber]
	if !ok {
		return nil, domain.ErrInvNotFound
	}
	return r.loadInvoice(inv)
}

func (r *MemorySaleRepo) ListInvoices(_ context.Context, filter domain.CustomerInvoiceFilter) ([]domain.CustomerInvoice, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.CustomerInvoice
	for _, inv := range r.invoices {
		if filter.CompanyID != "" && inv.CompanyID != filter.CompanyID {
			continue
		}
		if filter.CustomerID != "" && inv.CustomerID != filter.CustomerID {
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
	if filter.Offset > 0 && filter.Offset < len(out) {
		out = out[filter.Offset:]
	} else if filter.Offset >= len(out) {
		out = nil
	}
	if filter.Limit > 0 && filter.Limit < len(out) {
		out = out[:filter.Limit]
	}
	return out, total, nil
}

func (r *MemorySaleRepo) UpdateInvoice(_ context.Context, inv *domain.CustomerInvoice) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.invoices[inv.ID]
	if !ok {
		return domain.ErrInvNotFound
	}
	if existing.Status != domain.SInvDraft {
		return domain.ErrInvCannotUpdate
	}
	cp := *inv
	r.invoices[inv.ID] = &cp
	if r.invByNumber[cp.CompanyID] == nil {
		r.invByNumber[cp.CompanyID] = make(map[string]*domain.CustomerInvoice)
	}
	r.invByNumber[cp.CompanyID][cp.InvoiceNumber] = r.invoices[inv.ID]
	if len(cp.Lines) > 0 {
		r.invLines[inv.ID] = make([]domain.InvLine, len(cp.Lines))
		copy(r.invLines[inv.ID], cp.Lines)
	}
	return nil
}

func (r *MemorySaleRepo) UpdateInvoiceStatus(_ context.Context, id string, status domain.SaleInvoiceStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	inv, ok := r.invoices[id]
	if !ok {
		return domain.ErrInvNotFound
	}
	inv.Status = status
	return nil
}

func (r *MemorySaleRepo) PostInvoice(_ context.Context, id string, postedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	inv, ok := r.invoices[id]
	if !ok {
		return domain.ErrInvNotFound
	}
	inv.Status = domain.SInvPosted
	return nil
}

func (r *MemorySaleRepo) SetInvoiceGLPosted(_ context.Context, id string, postedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	inv, ok := r.invoices[id]
	if !ok {
		return domain.ErrInvNotFound
	}
	inv.GLPosted = true
	inv.GLPostedAt = &postedAt
	return nil
}

func (r *MemorySaleRepo) AllocateToInvoice(_ context.Context, invoiceID string, amount float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	inv, ok := r.invoices[invoiceID]
	if !ok {
		return domain.ErrInvNotFound
	}
	inv.AmountReceived += amount
	inv.BalanceDue -= amount
	return nil
}

func (r *MemorySaleRepo) GetInvoiceLines(_ context.Context, invoiceID string) ([]domain.InvLine, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	lines, ok := r.invLines[invoiceID]
	if !ok {
		return nil, nil
	}
	out := make([]domain.InvLine, len(lines))
	copy(out, lines)
	return out, nil
}

func (r *MemorySaleRepo) CreateInvoiceLines(_ context.Context, items []domain.InvLine) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(items) == 0 {
		return nil
	}
	invID := items[0].InvoiceID
	r.invLines[invID] = make([]domain.InvLine, len(items))
	copy(r.invLines[invID], items)
	return nil
}

func (r *MemorySaleRepo) UpdateInvoiceLines(_ context.Context, items []domain.InvLine) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(items) == 0 {
		return nil
	}
	invID := items[0].InvoiceID
	r.invLines[invID] = make([]domain.InvLine, len(items))
	copy(r.invLines[invID], items)
	return nil
}

// ─── Receipt ────────────────────────────────────────────────────────────

func (r *MemorySaleRepo) CreateReceipt(_ context.Context, rcpt *domain.CustomerReceipt) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *rcpt
	if cp.ID == "" {
		cp.ID = saleID("RCP-")
	}
	if cp.Status == "" {
		cp.Status = domain.RcpDraft
	}
	now := time.Now()
	cp.CreatedAt = now
	if r.rcptByNumber[cp.CompanyID] == nil {
		r.rcptByNumber[cp.CompanyID] = make(map[string]*domain.CustomerReceipt)
	}
	r.receipts[cp.ID] = &cp
	r.rcptByNumber[cp.CompanyID][cp.ReceiptNumber] = r.receipts[cp.ID]
	if len(cp.Allocations) > 0 {
		r.rcptAllocs[cp.ID] = make([]domain.RcpAllocation, len(cp.Allocations))
		copy(r.rcptAllocs[cp.ID], cp.Allocations)
	}
	rcpt.ID = cp.ID
	return nil
}

func (r *MemorySaleRepo) GetReceipt(_ context.Context, id string) (*domain.CustomerReceipt, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rcpt, ok := r.receipts[id]
	if !ok {
		return nil, domain.ErrRcpNotFound
	}
	cp := *rcpt
	if allocs, has := r.rcptAllocs[id]; has {
		cp.Allocations = make([]domain.RcpAllocation, len(allocs))
		copy(cp.Allocations, allocs)
	}
	return &cp, nil
}

func (r *MemorySaleRepo) GetReceiptByNumber(_ context.Context, companyID, receiptNumber string) (*domain.CustomerReceipt, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cm, ok := r.rcptByNumber[companyID]
	if !ok {
		return nil, domain.ErrRcpNotFound
	}
	rcpt, ok := cm[receiptNumber]
	if !ok {
		return nil, domain.ErrRcpNotFound
	}
	cp := *rcpt
	if allocs, has := r.rcptAllocs[rcpt.ID]; has {
		cp.Allocations = make([]domain.RcpAllocation, len(allocs))
		copy(cp.Allocations, allocs)
	}
	return &cp, nil
}

func (r *MemorySaleRepo) ListReceipts(_ context.Context, filter domain.ReceiptFilter) ([]domain.CustomerReceipt, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.CustomerReceipt
	for _, rcpt := range r.receipts {
		if filter.CompanyID != "" && rcpt.CompanyID != filter.CompanyID {
			continue
		}
		if filter.CustomerID != "" && rcpt.CustomerID != filter.CustomerID {
			continue
		}
		if filter.Status != "" && rcpt.Status != filter.Status {
			continue
		}
		if filter.FromDate != "" && rcpt.ReceiptDate.Format("2006-01-02") < filter.FromDate {
			continue
		}
		if filter.ToDate != "" && rcpt.ReceiptDate.Format("2006-01-02") > filter.ToDate {
			continue
		}
		cp := *rcpt
		if allocs, has := r.rcptAllocs[rcpt.ID]; has {
			cp.Allocations = make([]domain.RcpAllocation, len(allocs))
			copy(cp.Allocations, allocs)
		}
		out = append(out, cp)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ReceiptDate.After(out[j].ReceiptDate) })
	total := len(out)
	if filter.Offset > 0 && filter.Offset < len(out) {
		out = out[filter.Offset:]
	} else if filter.Offset >= len(out) {
		out = nil
	}
	if filter.Limit > 0 && filter.Limit < len(out) {
		out = out[:filter.Limit]
	}
	return out, total, nil
}

func (r *MemorySaleRepo) UpdateReceipt(_ context.Context, rcpt *domain.CustomerReceipt) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.receipts[rcpt.ID]
	if !ok {
		return domain.ErrRcpNotFound
	}
	cp := *rcpt
	r.receipts[rcpt.ID] = &cp
	if len(cp.Allocations) > 0 {
		r.rcptAllocs[rcpt.ID] = make([]domain.RcpAllocation, len(cp.Allocations))
		copy(r.rcptAllocs[rcpt.ID], cp.Allocations)
	}
	return nil
}

func (r *MemorySaleRepo) UpdateReceiptStatus(_ context.Context, id string, status domain.ReceiptStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	rcpt, ok := r.receipts[id]
	if !ok {
		return domain.ErrRcpNotFound
	}
	rcpt.Status = status
	return nil
}

func (r *MemorySaleRepo) SetReceiptGLPosted(_ context.Context, id string, postedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	rcpt, ok := r.receipts[id]
	if !ok {
		return domain.ErrRcpNotFound
	}
	rcpt.GLPosted = true
	rcpt.GLPostedAt = &postedAt
	return nil
}

func (r *MemorySaleRepo) CreateReceiptAllocations(_ context.Context, allocs []domain.RcpAllocation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(allocs) == 0 {
		return nil
	}
	rcptID := allocs[0].ReceiptID
	r.rcptAllocs[rcptID] = make([]domain.RcpAllocation, len(allocs))
	copy(r.rcptAllocs[rcptID], allocs)
	return nil
}

func (r *MemorySaleRepo) GetReceiptAllocations(_ context.Context, receiptID string) ([]domain.RcpAllocation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	allocs, ok := r.rcptAllocs[receiptID]
	if !ok {
		return nil, nil
	}
	out := make([]domain.RcpAllocation, len(allocs))
	copy(out, allocs)
	return out, nil
}

// ─── Credit Note ────────────────────────────────────────────────────────

func (r *MemorySaleRepo) CreateCN(_ context.Context, cn *domain.CreditNote) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *cn
	if cp.ID == "" {
		cp.ID = saleID("CN-")
	}
	if cp.Status == "" {
		cp.Status = domain.CNDraft
	}
	now := time.Now()
	cp.CreatedAt = now
	if r.cnByNumber[cp.CompanyID] == nil {
		r.cnByNumber[cp.CompanyID] = make(map[string]*domain.CreditNote)
	}
	r.cns[cp.ID] = &cp
	r.cnByNumber[cp.CompanyID][cp.CNNumber] = r.cns[cp.ID]
	if len(cp.Lines) > 0 {
		r.cnLines[cp.ID] = make([]domain.CNLine, len(cp.Lines))
		copy(r.cnLines[cp.ID], cp.Lines)
	}
	cn.ID = cp.ID
	return nil
}

func (r *MemorySaleRepo) GetCN(_ context.Context, id string) (*domain.CreditNote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cn, ok := r.cns[id]
	if !ok {
		return nil, domain.ErrCNNotFound
	}
	return r.loadCN(cn)
}

func (r *MemorySaleRepo) GetCNByNumber(_ context.Context, companyID, cnNumber string) (*domain.CreditNote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cm, ok := r.cnByNumber[companyID]
	if !ok {
		return nil, domain.ErrCNNotFound
	}
	cn, ok := cm[cnNumber]
	if !ok {
		return nil, domain.ErrCNNotFound
	}
	return r.loadCN(cn)
}

func (r *MemorySaleRepo) ListCNs(_ context.Context, filter domain.CreditNoteFilter) ([]domain.CreditNote, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.CreditNote
	for _, cn := range r.cns {
		if filter.CompanyID != "" && cn.CompanyID != filter.CompanyID {
			continue
		}
		if filter.CustomerID != "" && cn.CustomerID != filter.CustomerID {
			continue
		}
		if filter.Status != "" && cn.Status != filter.Status {
			continue
		}
		if filter.FromDate != "" && cn.ReturnDate.Format("2006-01-02") < filter.FromDate {
			continue
		}
		if filter.ToDate != "" && cn.ReturnDate.Format("2006-01-02") > filter.ToDate {
			continue
		}
		loaded, _ := r.loadCN(cn)
		if loaded != nil {
			out = append(out, *loaded)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ReturnDate.After(out[j].ReturnDate) })
	total := len(out)
	if filter.Offset > 0 && filter.Offset < len(out) {
		out = out[filter.Offset:]
	} else if filter.Offset >= len(out) {
		out = nil
	}
	if filter.Limit > 0 && filter.Limit < len(out) {
		out = out[:filter.Limit]
	}
	return out, total, nil
}

func (r *MemorySaleRepo) UpdateCN(_ context.Context, cn *domain.CreditNote) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.cns[cn.ID]
	if !ok {
		return domain.ErrCNNotFound
	}
	cp := *cn
	r.cns[cn.ID] = &cp
	if len(cp.Lines) > 0 {
		r.cnLines[cn.ID] = make([]domain.CNLine, len(cp.Lines))
		copy(r.cnLines[cn.ID], cp.Lines)
	}
	return nil
}

func (r *MemorySaleRepo) UpdateCNStatus(_ context.Context, id string, status domain.CNStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cn, ok := r.cns[id]
	if !ok {
		return domain.ErrCNNotFound
	}
	cn.Status = status
	return nil
}

func (r *MemorySaleRepo) PostCN(_ context.Context, id string, postedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cn, ok := r.cns[id]
	if !ok {
		return domain.ErrCNNotFound
	}
	cn.Status = domain.CNPosted
	return nil
}

func (r *MemorySaleRepo) SetCNGLPosted(_ context.Context, id string, postedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cn, ok := r.cns[id]
	if !ok {
		return domain.ErrCNNotFound
	}
	cn.GLPosted = true
	cn.GLPostedAt = &postedAt
	return nil
}

func (r *MemorySaleRepo) GetCNLines(_ context.Context, cnID string) ([]domain.CNLine, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	lines, ok := r.cnLines[cnID]
	if !ok {
		return nil, nil
	}
	out := make([]domain.CNLine, len(lines))
	copy(out, lines)
	return out, nil
}

func (r *MemorySaleRepo) CreateCNLines(_ context.Context, items []domain.CNLine) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(items) == 0 {
		return nil
	}
	cnID := items[0].CNID
	r.cnLines[cnID] = make([]domain.CNLine, len(items))
	copy(r.cnLines[cnID], items)
	return nil
}

// ─── AR Transaction ─────────────────────────────────────────────────────

func (r *MemorySaleRepo) CreateARTransaction(_ context.Context, t *domain.ARTransaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *t
	if cp.ID == "" {
		cp.ID = saleID("ART-")
	}
	if cp.Currency == "" {
		cp.Currency = "VND"
	}
	now := time.Now()
	cp.CreatedAt = now
	r.arTxns[cp.ID] = &cp
	t.ID = cp.ID
	return nil
}

func (r *MemorySaleRepo) GetARTransaction(_ context.Context, id string) (*domain.ARTransaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.arTxns[id]
	if !ok {
		return nil, domain.ErrARTransNotFound
	}
	cp := *t
	return &cp, nil
}

func (r *MemorySaleRepo) ListARTransactions(_ context.Context, companyID, customerID string) ([]domain.ARTransaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.ARTransaction
	for _, t := range r.arTxns {
		if companyID != "" {
			cust, ok := r.customers[t.CustomerID]
			if !ok || cust.CompanyID != companyID {
				continue
			}
		}
		if customerID != "" && t.CustomerID != customerID {
			continue
		}
		out = append(out, *t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].TransactionDate.After(out[j].TransactionDate) })
	return out, nil
}

func (r *MemorySaleRepo) ListARTransactionsAll(_ context.Context, companyID string, offset, limit int) ([]domain.ARTransaction, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var filtered []domain.ARTransaction
	for _, t := range r.arTxns {
		if companyID != "" {
			cust, ok := r.customers[t.CustomerID]
			if !ok || cust.CompanyID != companyID {
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

// ─── Sales Quotation ────────────────────────────────────────────────────

func (r *MemorySaleRepo) CreateSQ(_ context.Context, sq *domain.SalesQuotation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *sq
	if cp.ID == "" {
		cp.ID = saleID("SQ-")
	}
	now := time.Now()
	cp.CreatedAt = now
	r.sqs[cp.ID] = &cp
	sq.ID = cp.ID
	return nil
}

func (r *MemorySaleRepo) GetSQ(_ context.Context, id string) (*domain.SalesQuotation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sq, ok := r.sqs[id]
	if !ok {
		return nil, nil
	}
	cp := *sq
	return &cp, nil
}

func (r *MemorySaleRepo) ListSQs(_ context.Context, companyID string) ([]domain.SalesQuotation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.SalesQuotation
	for _, sq := range r.sqs {
		if sq.CompanyID == companyID {
			out = append(out, *sq)
		}
	}
	return out, nil
}

func (r *MemorySaleRepo) UpdateSQ(_ context.Context, sq *domain.SalesQuotation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.sqs[sq.ID]
	if !ok {
		return nil
	}
	cp := *sq
	r.sqs[sq.ID] = &cp
	return nil
}

// ─── helpers ────────────────────────────────────────────────────────────

func (r *MemorySaleRepo) loadSO(so *domain.SalesOrder) (*domain.SalesOrder, error) {
	cp := *so
	if lines, ok := r.soLines[so.ID]; ok {
		cp.Lines = make([]domain.SOLine, len(lines))
		copy(cp.Lines, lines)
	}
	return &cp, nil
}

func (r *MemorySaleRepo) loadDN(dn *domain.DeliveryNote) (*domain.DeliveryNote, error) {
	cp := *dn
	if lines, ok := r.dnLines[dn.ID]; ok {
		cp.Lines = make([]domain.DNLine, len(lines))
		copy(cp.Lines, lines)
	}
	return &cp, nil
}

func (r *MemorySaleRepo) loadInvoice(inv *domain.CustomerInvoice) (*domain.CustomerInvoice, error) {
	cp := *inv
	if lines, ok := r.invLines[inv.ID]; ok {
		cp.Lines = make([]domain.InvLine, len(lines))
		copy(cp.Lines, lines)
	}
	return &cp, nil
}

func (r *MemorySaleRepo) loadCN(cn *domain.CreditNote) (*domain.CreditNote, error) {
	cp := *cn
	if lines, ok := r.cnLines[cn.ID]; ok {
		cp.Lines = make([]domain.CNLine, len(lines))
		copy(cp.Lines, lines)
	}
	return &cp, nil
}
