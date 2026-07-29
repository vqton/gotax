package service

import (
	"context"
	"time"

	"gotax/internal/domain"
)

type SaleService struct {
	custRepo domain.CustomerRepository
	soRepo   domain.SaleOrderRepository
	dnRepo   domain.DeliveryNoteRepository
	invRepo  domain.CustomerInvoiceRepository
	rcptRepo domain.CustomerReceiptRepository
	cnRepo   domain.CreditNoteRepository
	artRepo  domain.ARTransactionRepository
	sqRepo   domain.SalesQuotationRepository
	now      func() time.Time
}

func NewSaleService(
	custRepo domain.CustomerRepository,
	soRepo domain.SaleOrderRepository,
	dnRepo domain.DeliveryNoteRepository,
	invRepo domain.CustomerInvoiceRepository,
	rcptRepo domain.CustomerReceiptRepository,
	cnRepo domain.CreditNoteRepository,
	artRepo domain.ARTransactionRepository,
	sqRepo domain.SalesQuotationRepository,
) *SaleService {
	return &SaleService{
		custRepo: custRepo, soRepo: soRepo, dnRepo: dnRepo,
		invRepo: invRepo, rcptRepo: rcptRepo, cnRepo: cnRepo,
		artRepo: artRepo, sqRepo: sqRepo,
		now: time.Now,
	}
}

// ─── Customer ──────────────────────────────────────────────────────────

func (s *SaleService) CreateCustomer(ctx context.Context, c *domain.Customer) error {
	if err := c.Validate(); err != nil {
		return err
	}
	existing, _ := s.custRepo.GetCustomerByCode(ctx, c.CompanyID, c.Code)
	if existing != nil {
		return domain.ErrCustomerCodeExists
	}
	if c.Status == "" {
		c.Status = domain.CustomerActive
	}
	c.CreatedAt = s.now()
	c.UpdatedAt = c.CreatedAt
	return s.custRepo.CreateCustomer(ctx, c)
}

func (s *SaleService) GetCustomer(ctx context.Context, id string) (*domain.Customer, error) {
	return s.custRepo.GetCustomer(ctx, id)
}

func (s *SaleService) ListCustomers(ctx context.Context, companyID string) ([]domain.Customer, error) {
	return s.custRepo.ListCustomers(ctx, companyID)
}

func (s *SaleService) UpdateCustomer(ctx context.Context, c *domain.Customer) error {
	existing, err := s.custRepo.GetCustomer(ctx, c.ID)
	if err != nil {
		return err
	}
	if existing.Code != c.Code {
		dup, _ := s.custRepo.GetCustomerByCode(ctx, c.CompanyID, c.Code)
		if dup != nil {
			return domain.ErrCustomerCodeExists
		}
	}
	if c.Status == "" {
		c.Status = domain.CustomerActive
	}
	c.UpdatedAt = s.now()
	return s.custRepo.UpdateCustomer(ctx, c)
}

func (s *SaleService) DeleteCustomer(ctx context.Context, id string) error {
	return s.custRepo.DeleteCustomer(ctx, id)
}

// ─── Sales Order ───────────────────────────────────────────────────────

func (s *SaleService) CreateSO(ctx context.Context, so *domain.SalesOrder) error {
	if err := so.Validate(); err != nil {
		return err
	}
	so.CalculateTotals()
	if so.Status == "" {
		so.Status = domain.SODraft
	}
	if so.SONumber == "" {
		return domain.ErrSONumberRequired
	}
	if existing, _ := s.soRepo.GetSOByNumber(ctx, so.CompanyID, so.SONumber); existing != nil {
		return domain.ErrSONumberExists
	}
	if _, err := s.custRepo.GetCustomer(ctx, so.CustomerID); err != nil {
		return domain.ErrSOCustomerRequired
	}
	so.CreatedAt = s.now()
	so.UpdatedAt = so.CreatedAt
	return s.soRepo.CreateSO(ctx, so)
}

func (s *SaleService) GetSO(ctx context.Context, id string) (*domain.SalesOrder, error) {
	return s.soRepo.GetSO(ctx, id)
}

func (s *SaleService) ListSOs(ctx context.Context, filter domain.SalesOrderFilter) ([]domain.SalesOrder, int, error) {
	return s.soRepo.ListSOs(ctx, filter)
}

func (s *SaleService) UpdateSO(ctx context.Context, so *domain.SalesOrder) error {
	existing, err := s.soRepo.GetSO(ctx, so.ID)
	if err != nil {
		return err
	}
	if existing.Status != domain.SODraft {
		return domain.ErrSOCannotUpdate
	}
	so.CalculateTotals()
	so.UpdatedAt = s.now()
	return s.soRepo.UpdateSO(ctx, so)
}

func (s *SaleService) ApproveSO(ctx context.Context, id, approvedBy string) error {
	so, err := s.soRepo.GetSO(ctx, id)
	if err != nil {
		return err
	}
	if so.Status != domain.SODraft {
		return domain.ErrSOAlreadyApproved
	}
	return s.soRepo.ApproveSO(ctx, id, approvedBy, s.now())
}

func (s *SaleService) CancelSO(ctx context.Context, id, reason string) error {
	so, err := s.soRepo.GetSO(ctx, id)
	if err != nil {
		return err
	}
	if so.Status == domain.SOCancelled {
		return domain.ErrSOAlreadyCancelled
	}
	return s.soRepo.CancelSO(ctx, id, reason)
}

func (s *SaleService) CloseSO(ctx context.Context, id string) error {
	so, err := s.soRepo.GetSO(ctx, id)
	if err != nil {
		return err
	}
	if so.Status != domain.SODelivered {
		return domain.ErrSOInvalidTransition
	}
	return s.soRepo.UpdateSOStatus(ctx, id, domain.SOClosed)
}

// ─── Delivery Note ─────────────────────────────────────────────────────

func (s *SaleService) CreateDN(ctx context.Context, dn *domain.DeliveryNote) error {
	if err := dn.Validate(); err != nil {
		return err
	}
	if dn.DNNumber == "" {
		return domain.ErrDNNumberRequired
	}
	if existing, _ := s.dnRepo.GetDNByNumber(ctx, dn.CompanyID, dn.DNNumber); existing != nil {
		return domain.ErrDNNumberExists
	}
	if _, err := s.soRepo.GetSO(ctx, dn.SOID); err != nil {
		return domain.ErrDNSORequired
	}
	if dn.Status == "" {
		dn.Status = domain.DNDraft
	}
	dn.CreatedAt = s.now()
	return s.dnRepo.CreateDN(ctx, dn)
}

func (s *SaleService) GetDN(ctx context.Context, id string) (*domain.DeliveryNote, error) {
	return s.dnRepo.GetDN(ctx, id)
}

func (s *SaleService) ListDNs(ctx context.Context, filter domain.DeliveryNoteFilter) ([]domain.DeliveryNote, int, error) {
	return s.dnRepo.ListDNs(ctx, filter)
}

func (s *SaleService) UpdateDN(ctx context.Context, dn *domain.DeliveryNote) error {
	existing, err := s.dnRepo.GetDN(ctx, dn.ID)
	if err != nil {
		return err
	}
	if existing.Status != domain.DNDraft {
		return domain.ErrDNCannotUpdate
	}
	return s.dnRepo.UpdateDN(ctx, dn)
}

func (s *SaleService) PostDN(ctx context.Context, id string) error {
	dn, err := s.dnRepo.GetDN(ctx, id)
	if err != nil {
		return err
	}
	if dn.Status != domain.DNDraft {
		return domain.ErrDNInvalidTransition
	}
	return s.dnRepo.UpdateDNStatus(ctx, id, domain.DNPosted)
}

func (s *SaleService) CancelDN(ctx context.Context, id string) error {
	dn, err := s.dnRepo.GetDN(ctx, id)
	if err != nil {
		return err
	}
	if dn.Status == domain.DNPosted {
		return domain.ErrDNInvalidTransition
	}
	return s.dnRepo.UpdateDNStatus(ctx, id, domain.DNCancelled)
}

// ─── Customer Invoice ──────────────────────────────────────────────────

func (s *SaleService) CreateInvoice(ctx context.Context, inv *domain.CustomerInvoice) error {
	if err := inv.Validate(); err != nil {
		return err
	}
	if inv.InvoiceNumber == "" {
		return domain.ErrInvNumberRequired
	}
	if existing, _ := s.invRepo.GetInvoiceByNumber(ctx, inv.CompanyID, inv.InvoiceNumber); existing != nil {
		return domain.ErrInvNumberExists
	}
	cust, err := s.custRepo.GetCustomer(ctx, inv.CustomerID)
	if err != nil {
		return domain.ErrInvCustomerRequired
	}
	inv.CustomerName = cust.Name
	inv.CustomerTaxCode = cust.TaxCode
	if inv.Currency == "" {
		inv.Currency = "VND"
	}
	if inv.Status == "" {
		inv.Status = domain.SInvDraft
	}
	inv.CalculateTotals()
	inv.CreatedAt = s.now()
	return s.invRepo.CreateInvoice(ctx, inv)
}

func (s *SaleService) GetInvoice(ctx context.Context, id string) (*domain.CustomerInvoice, error) {
	return s.invRepo.GetInvoice(ctx, id)
}

func (s *SaleService) ListInvoices(ctx context.Context, filter domain.CustomerInvoiceFilter) ([]domain.CustomerInvoice, int, error) {
	return s.invRepo.ListInvoices(ctx, filter)
}

func (s *SaleService) UpdateInvoice(ctx context.Context, inv *domain.CustomerInvoice) error {
	existing, err := s.invRepo.GetInvoice(ctx, inv.ID)
	if err != nil {
		return err
	}
	if existing.Status != domain.SInvDraft {
		return domain.ErrInvCannotUpdate
	}
	if inv.Currency == "" {
		inv.Currency = "VND"
	}
	inv.CalculateTotals()
	return s.invRepo.UpdateInvoice(ctx, inv)
}

func (s *SaleService) PostInvoice(ctx context.Context, id string) error {
	inv, err := s.invRepo.GetInvoice(ctx, id)
	if err != nil {
		return err
	}
	if inv.Status != domain.SInvSigned {
		return domain.ErrInvInvalidTransition
	}
	return s.invRepo.PostInvoice(ctx, id, s.now().UTC())
}

func (s *SaleService) CancelInvoice(ctx context.Context, id string) error {
	inv, err := s.invRepo.GetInvoice(ctx, id)
	if err != nil {
		return err
	}
	if inv.Status == domain.SInvPosted || inv.Status == domain.SInvPaid {
		return domain.ErrInvInvalidTransition
	}
	return s.invRepo.UpdateInvoiceStatus(ctx, id, domain.SInvCancelled)
}

// ─── Receipt ───────────────────────────────────────────────────────────

func (s *SaleService) CreateReceipt(ctx context.Context, rcpt *domain.CustomerReceipt) error {
	if err := rcpt.Validate(); err != nil {
		return err
	}
	if rcpt.ReceiptNumber == "" {
		return domain.ErrRcpNumberRequired
	}
	if existing, _ := s.rcptRepo.GetReceiptByNumber(ctx, rcpt.CompanyID, rcpt.ReceiptNumber); existing != nil {
		return domain.ErrRcpNumberExists
	}
	if _, err := s.custRepo.GetCustomer(ctx, rcpt.CustomerID); err != nil {
		return domain.ErrRcpCustomerRequired
	}
	if rcpt.Currency == "" {
		rcpt.Currency = "VND"
	}
	if rcpt.Status == "" {
		rcpt.Status = domain.RcpDraft
	}
	rcpt.UnallocatedAmount = rcpt.Amount
	rcpt.CreatedAt = s.now()
	return s.rcptRepo.CreateReceipt(ctx, rcpt)
}

func (s *SaleService) GetReceipt(ctx context.Context, id string) (*domain.CustomerReceipt, error) {
	return s.rcptRepo.GetReceipt(ctx, id)
}

func (s *SaleService) ListReceipts(ctx context.Context, filter domain.ReceiptFilter) ([]domain.CustomerReceipt, int, error) {
	return s.rcptRepo.ListReceipts(ctx, filter)
}

func (s *SaleService) PostReceipt(ctx context.Context, id string) error {
	rcpt, err := s.rcptRepo.GetReceipt(ctx, id)
	if err != nil {
		return err
	}
	if rcpt.Status != domain.RcpDraft {
		return domain.ErrRcpInvalidTransition
	}
	return s.rcptRepo.UpdateReceiptStatus(ctx, id, domain.RcpPosted)
}

func (s *SaleService) CancelReceipt(ctx context.Context, id string) error {
	rcpt, err := s.rcptRepo.GetReceipt(ctx, id)
	if err != nil {
		return err
	}
	if rcpt.Status == domain.RcpPosted {
		return domain.ErrRcpInvalidTransition
	}
	return s.rcptRepo.UpdateReceiptStatus(ctx, id, domain.RcpCancelled)
}

// ─── Credit Note ───────────────────────────────────────────────────────

func (s *SaleService) CreateCN(ctx context.Context, cn *domain.CreditNote) error {
	if err := cn.Validate(); err != nil {
		return err
	}
	if cn.CNNumber == "" {
		return domain.ErrCNNumberRequired
	}
	if existing, _ := s.cnRepo.GetCNByNumber(ctx, cn.CompanyID, cn.CNNumber); existing != nil {
		return domain.ErrCNNumberExists
	}
	if _, err := s.invRepo.GetInvoice(ctx, cn.OriginalInvoiceID); err != nil {
		return domain.ErrCNOriginalInvRequired
	}
	if _, err := s.custRepo.GetCustomer(ctx, cn.CustomerID); err != nil {
		return domain.ErrCNCustomerRequired
	}
	if cn.Status == "" {
		cn.Status = domain.CNDraft
	}
	cn.CalculateTotals()
	cn.CreatedAt = s.now()
	return s.cnRepo.CreateCN(ctx, cn)
}

func (s *SaleService) GetCN(ctx context.Context, id string) (*domain.CreditNote, error) {
	return s.cnRepo.GetCN(ctx, id)
}

func (s *SaleService) ListCNs(ctx context.Context, filter domain.CreditNoteFilter) ([]domain.CreditNote, int, error) {
	return s.cnRepo.ListCNs(ctx, filter)
}

func (s *SaleService) PostCN(ctx context.Context, id string) error {
	cn, err := s.cnRepo.GetCN(ctx, id)
	if err != nil {
		return err
	}
	if cn.Status != domain.CNDraft {
		return domain.ErrCNAlreadyPosted
	}
	return s.cnRepo.PostCN(ctx, id, s.now().UTC())
}

func (s *SaleService) CancelCN(ctx context.Context, id string) error {
	cn, err := s.cnRepo.GetCN(ctx, id)
	if err != nil {
		return err
	}
	if cn.Status == domain.CNPosted {
		return domain.ErrCNAlreadyPosted
	}
	return s.cnRepo.UpdateCNStatus(ctx, id, domain.CNCancelled)
}

// ─── AR Transactions ───────────────────────────────────────────────────

func (s *SaleService) CreateARTransaction(ctx context.Context, t *domain.ARTransaction) error {
	if err := t.Validate(); err != nil {
		return err
	}
	if _, err := s.custRepo.GetCustomer(ctx, t.CustomerID); err != nil {
		return domain.ErrARTransCustomerRequired
	}
	if t.Currency == "" {
		t.Currency = "VND"
	}
	t.CreatedAt = s.now()
	return s.artRepo.CreateARTransaction(ctx, t)
}

func (s *SaleService) GetARTransaction(ctx context.Context, id string) (*domain.ARTransaction, error) {
	return s.artRepo.GetARTransaction(ctx, id)
}

func (s *SaleService) ListARTransactions(ctx context.Context, companyID, customerID string) ([]domain.ARTransaction, error) {
	return s.artRepo.ListARTransactions(ctx, companyID, customerID)
}

// ─── Reports ───────────────────────────────────────────────────────────

func (s *SaleService) GetARAgingReport(ctx context.Context, companyID string) ([]domain.ARAgingReport, error) {
	txns, _, err := s.artRepo.ListARTransactionsAll(ctx, companyID, 0, 0)
	if err != nil {
		return nil, err
	}
	byCust := make(map[string]*domain.ARAgingReport)
	custBal := make(map[string]float64)
	for _, t := range txns {
		custBal[t.CustomerID] += t.Amount
	}
	for cid, bal := range custBal {
		rpt := &domain.ARAgingReport{
			CustomerID: cid,
			Buckets:    domain.ARAgingBucket{Bucket0: bal, Bucket30: 0, Bucket60: 0, Bucket90: 0, Total: bal},
		}
		byCust[cid] = rpt
	}
	out := make([]domain.ARAgingReport, 0, len(byCust))
	for _, r := range byCust {
		out = append(out, *r)
	}
	return out, nil
}

func (s *SaleService) GetARSummary(ctx context.Context, companyID string) ([]domain.ARSummary, error) {
	txns, _, err := s.artRepo.ListARTransactionsAll(ctx, companyID, 0, 0)
	if err != nil {
		return nil, err
	}
	custIDs := make([]string, 0, len(txns))
	seen := make(map[string]bool)
	for _, t := range txns {
		if !seen[t.CustomerID] {
			custIDs = append(custIDs, t.CustomerID)
			seen[t.CustomerID] = true
		}
	}
	m := make(map[string]*domain.ARSummary)
	for _, t := range txns {
		if _, ok := m[t.CustomerID]; !ok {
			cust, _ := s.custRepo.GetCustomer(ctx, t.CustomerID)
			name := ""
			tax := ""
			if cust != nil {
				name = cust.Name
				tax = cust.TaxCode
			}
			m[t.CustomerID] = &domain.ARSummary{
				CustomerID: t.CustomerID, CustomerName: name,
				TaxCode: tax, Currency: t.Currency,
			}
		}
		cur := m[t.CustomerID]
		switch t.TransactionType {
		case domain.ARTransInvoice:
			cur.TotalInvoiced += t.Amount
		case domain.ARTransCreditNote, domain.ARTransReceipt, domain.ARTransOffset:
			cur.TotalReceived += t.Amount
		}
		cur.Outstanding = cur.TotalInvoiced - cur.TotalReceived
	}
	out := make([]domain.ARSummary, 0, len(m))
	for _, v := range m {
		out = append(out, *v)
	}
	return out, nil
}

// ─── Sales Quotation (P1) ──────────────────────────────────────────────

func (s *SaleService) CreateSQ(ctx context.Context, sq *domain.SalesQuotation) error {
	if sq.QNNumber == "" {
		return domain.ErrSQNumberRequired
	}
	if sq.Status == "" {
		sq.Status = "DRAFT"
	}
	sq.CreatedAt = s.now()
	return s.sqRepo.CreateSQ(ctx, sq)
}

func (s *SaleService) GetSQ(ctx context.Context, id string) (*domain.SalesQuotation, error) {
	return s.sqRepo.GetSQ(ctx, id)
}

func (s *SaleService) ListSQs(ctx context.Context, companyID string) ([]domain.SalesQuotation, error) {
	return s.sqRepo.ListSQs(ctx, companyID)
}

func (s *SaleService) UpdateSQ(ctx context.Context, sq *domain.SalesQuotation) error {
	existing, err := s.sqRepo.GetSQ(ctx, sq.ID)
	if err != nil || existing == nil {
		return domain.ErrSQNotFound
	}
	return s.sqRepo.UpdateSQ(ctx, sq)
}

func (s *SaleService) NextSONumber(ctx context.Context, companyID, yyyymm string) (string, error) {
	return s.soRepo.NextSONumber(ctx, companyID, yyyymm)
}

func (s *SaleService) NextDNNumber(ctx context.Context, companyID, yyyymm string) (string, error) {
	return s.dnRepo.NextDNNumber(ctx, companyID, yyyymm)
}

// ─── Auto-numbering ────────────────────────────────────────────────────

func (s *SaleService) GetInvNumber(ctx context.Context, companyID, yyyymm string) (string, error) {
	return s.invRepo.NextInvNumber(ctx, companyID, yyyymm)
}
