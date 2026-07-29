package service

import (
	"context"
	"sort"
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
	gl       Service
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
	gl Service,
) *SaleService {
	return &SaleService{
		custRepo: custRepo, soRepo: soRepo, dnRepo: dnRepo,
		invRepo: invRepo, rcptRepo: rcptRepo, cnRepo: cnRepo,
		artRepo: artRepo, sqRepo: sqRepo,
		gl: gl,
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

func (s *SaleService) checkCreditLimit(ctx context.Context, customerID string, newAmount float64) error {
	cust, err := s.custRepo.GetCustomer(ctx, customerID)
	if err != nil {
		return err
	}
	if cust.CreditLimit <= 0 {
		return nil // no limit set
	}
	invoices, _, err := s.invRepo.ListInvoices(ctx, domain.CustomerInvoiceFilter{CustomerID: customerID})
	if err != nil {
		return err
	}
	outstanding := 0.0
	for _, inv := range invoices {
		if inv.Status == domain.SInvCancelled {
			continue
		}
		outstanding += inv.BalanceDue
	}
	if outstanding+newAmount > cust.CreditLimit {
		return domain.ErrCreditLimitExceeded
	}
	return nil
}

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
	if err := s.checkCreditLimit(ctx, so.CustomerID, so.TotalAmount); err != nil {
		return err
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
	if err := s.checkCreditLimit(ctx, inv.CustomerID, inv.TotalAmount); err != nil {
		return err
	}
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
	switch inv.Status {
	case domain.SInvDraft, domain.SInvIssued:
	default:
		return domain.ErrInvInvalidTransition
	}
	now := s.now().UTC()
	if s.gl != nil {
		lines := []domain.JournalLine{}
		lines = append(lines, domain.JournalLine{
			AccountCode: "131", DebitAmount: inv.TotalAmount, CreditAmount: 0,
			Description: "AR: " + inv.InvoiceNumber,
		})
		groupByRev := make(map[string]float64)
		groupByVAT := make(map[string]float64)
		for _, l := range inv.Lines {
			groupByRev[l.RevenueAccount] += l.LineTotal
			groupByVAT[l.VATAccountID] += l.LineVATAmount
		}
		for acc, amt := range groupByRev {
			lines = append(lines, domain.JournalLine{
				AccountCode: acc, DebitAmount: 0, CreditAmount: amt,
				Description: "Revenue: " + inv.InvoiceNumber,
			})
		}
		for acc, amt := range groupByVAT {
			lines = append(lines, domain.JournalLine{
				AccountCode: acc, DebitAmount: 0, CreditAmount: amt,
				Description: "VAT: " + inv.InvoiceNumber,
			})
		}
		entry := &domain.JournalEntry{
			CompanyID:   inv.CompanyID,
			EntryNumber: inv.InvoiceNumber,
			VoucherType: domain.VoucherTypeSale,
			EntryDate:   inv.InvoiceDate,
			Description: "AR invoice " + inv.InvoiceNumber,
			Lines:       lines,
		}
		if err := s.gl.CreatePostedEntry(ctx, entry, inv.CreatedBy); err != nil {
			return err
		}
		if err := s.invRepo.SetInvoiceGLPosted(ctx, id, now); err != nil {
			return err
		}
	}
	return s.invRepo.PostInvoice(ctx, id, now)
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
	now := s.now().UTC()
	if s.gl != nil {
		cashAcc := "1111"
		switch rcpt.PaymentMethod {
		case "bank_transfer", "cheque", "credit_card":
			cashAcc = "1121"
		}
		lines := []domain.JournalLine{
			{AccountCode: cashAcc, DebitAmount: rcpt.Amount, CreditAmount: 0, Description: "Receipt: " + rcpt.ReceiptNumber},
			{AccountCode: "131", DebitAmount: 0, CreditAmount: rcpt.Amount, Description: "Receipt: " + rcpt.ReceiptNumber},
		}
		entry := &domain.JournalEntry{
			CompanyID:   rcpt.CompanyID,
			EntryNumber: rcpt.ReceiptNumber,
			VoucherType: domain.VoucherTypeReceipt,
			EntryDate:   rcpt.ReceiptDate,
			Description: "AR receipt " + rcpt.ReceiptNumber,
			Lines:       lines,
		}
		if err := s.gl.CreatePostedEntry(ctx, entry, rcpt.CreatedBy); err != nil {
			return err
		}
		if err := s.rcptRepo.SetReceiptGLPosted(ctx, id, now); err != nil {
			return err
		}
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

func (s *SaleService) AllocateReceipt(ctx context.Context, receiptID, invoiceID string, amount float64) error {
	rcpt, err := s.rcptRepo.GetReceipt(ctx, receiptID)
	if err != nil {
		return err
	}
	if amount > rcpt.UnallocatedAmount {
		return domain.ErrRcpAllocExceedsReceipt
	}
	inv, err := s.invRepo.GetInvoice(ctx, invoiceID)
	if err != nil {
		return err
	}
	if amount > inv.BalanceDue {
		return domain.ErrRcpAllocExceedsBalance
	}
	if err := s.invRepo.AllocateToInvoice(ctx, invoiceID, amount); err != nil {
		return err
	}
	rcpt.UnallocatedAmount -= amount
	if err := s.rcptRepo.UpdateReceipt(ctx, rcpt); err != nil {
		return err
	}
	alloc := domain.RcpAllocation{
		ReceiptID: receiptID, InvoiceID: invoiceID, AllocatedAmount: amount,
	}
	return s.rcptRepo.CreateReceiptAllocations(ctx, []domain.RcpAllocation{alloc})
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
	now := s.now().UTC()
	if s.gl != nil {
		revByAccount := make(map[string]float64)
		if cn.ReturnType == domain.RetFull || cn.ReturnType == domain.RetPartial {
			origInv, err := s.invRepo.GetInvoice(ctx, cn.OriginalInvoiceID)
			if err == nil && len(origInv.Lines) > 0 {
				totalOrig := 0.0
				for _, l := range origInv.Lines {
					totalOrig += l.LineTotal
				}
				if totalOrig > 0 {
					for _, l := range origInv.Lines {
						revByAccount[l.RevenueAccount] += l.LineTotal / totalOrig * cn.Subtotal
					}
				}
			}
		}
		if len(revByAccount) == 0 {
			revByAccount["5213"] = cn.Subtotal
		}
		lines := []domain.JournalLine{}
		for acc, amt := range revByAccount {
			lines = append(lines, domain.JournalLine{
				AccountCode: acc, DebitAmount: amt, CreditAmount: 0,
				Description: "CN revenue reversal: " + cn.CNNumber,
			})
		}
		if cn.TaxAmount > 0 {
			lines = append(lines, domain.JournalLine{
				AccountCode: "3331", DebitAmount: cn.TaxAmount, CreditAmount: 0,
				Description: "CN VAT reversal: " + cn.CNNumber,
			})
		}
		lines = append(lines, domain.JournalLine{
			AccountCode: "131", DebitAmount: 0, CreditAmount: cn.TotalAmount,
			Description: "CN: " + cn.CNNumber,
		})
		entry := &domain.JournalEntry{
			CompanyID:   cn.CompanyID,
			EntryNumber: cn.CNNumber,
			VoucherType: domain.VoucherTypeSale,
			EntryDate:   cn.ReturnDate,
			Description: "AR credit note " + cn.CNNumber,
			Lines:       lines,
		}
		if err := s.gl.CreatePostedEntry(ctx, entry, cn.CreatedBy); err != nil {
			return err
		}
		if err := s.cnRepo.SetCNGLPosted(ctx, id, now); err != nil {
			return err
		}
	}
	return s.cnRepo.PostCN(ctx, id, now)
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
	now := s.now().UTC().Truncate(24 * time.Hour)
	invoices, _, err := s.invRepo.ListInvoices(ctx, domain.CustomerInvoiceFilter{CompanyID: companyID})
	if err != nil {
		return nil, err
	}
	byCust := make(map[string]*domain.ARAgingReport)
	for _, inv := range invoices {
		if inv.BalanceDue <= 0.001 {
			continue
		}
		rpt, ok := byCust[inv.CustomerID]
		if !ok {
			rpt = &domain.ARAgingReport{CustomerID: inv.CustomerID, CustomerName: inv.CustomerName, TaxCode: inv.CustomerTaxCode}
			byCust[inv.CustomerID] = rpt
		}
		var daysOverdue int
		if inv.DueDate != nil && inv.DueDate.Before(now) {
			daysOverdue = int(now.Sub(*inv.DueDate).Hours() / 24)
		}
		switch {
		case daysOverdue <= 0:
			rpt.Buckets.Bucket0 += inv.BalanceDue
		case daysOverdue <= 30:
			rpt.Buckets.Bucket30 += inv.BalanceDue
		case daysOverdue <= 60:
			rpt.Buckets.Bucket60 += inv.BalanceDue
		case daysOverdue <= 90:
			rpt.Buckets.Bucket90 += inv.BalanceDue
		default:
			rpt.Buckets.Bucket120 += inv.BalanceDue
		}
		rpt.Buckets.Total += inv.BalanceDue
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

func (s *SaleService) GetCustomerStatement(ctx context.Context, customerID, fromDate, toDate string) (*domain.CustomerStatement, error) {
	cust, err := s.custRepo.GetCustomer(ctx, customerID)
	if err != nil {
		return nil, err
	}
	invFilter := domain.CustomerInvoiceFilter{CustomerID: customerID, FromDate: fromDate, ToDate: toDate}
	invoices, _, err := s.invRepo.ListInvoices(ctx, invFilter)
	if err != nil {
		return nil, err
	}
	rcptFilter := domain.ReceiptFilter{CustomerID: customerID, FromDate: fromDate, ToDate: toDate}
	receipts, _, err := s.rcptRepo.ListReceipts(ctx, rcptFilter)
	if err != nil {
		return nil, err
	}
	cnFilter := domain.CreditNoteFilter{CustomerID: customerID, FromDate: fromDate, ToDate: toDate}
	cns, _, err := s.cnRepo.ListCNs(ctx, cnFilter)
	if err != nil {
		return nil, err
	}

	type stmtItem struct {
		date    time.Time
		refType string
		refNum  string
		desc    string
		debit   float64
		credit  float64
	}
	items := []stmtItem{}

	for _, inv := range invoices {
		items = append(items, stmtItem{
			date: inv.InvoiceDate, refType: "Invoice", refNum: inv.InvoiceNumber,
			desc: "AR invoice " + inv.InvoiceNumber, debit: inv.TotalAmount, credit: 0,
		})
	}
	for _, rcpt := range receipts {
		if rcpt.Status == domain.RcpCancelled {
			continue
		}
		items = append(items, stmtItem{
			date: rcpt.ReceiptDate, refType: "Receipt", refNum: rcpt.ReceiptNumber,
			desc: "Receipt " + rcpt.ReceiptNumber, debit: 0, credit: rcpt.Amount,
		})
	}
	for _, cn := range cns {
		if cn.Status == domain.CNCancelled {
			continue
		}
		items = append(items, stmtItem{
			date: cn.ReturnDate, refType: "Credit Note", refNum: cn.CNNumber,
			desc: "CN " + cn.CNNumber, debit: 0, credit: cn.TotalAmount,
		})
	}

	sort.Slice(items, func(i, j int) bool { return items[i].date.Before(items[j].date) })

	stmt := &domain.CustomerStatement{
		Customer: *cust,
		FromDate: fromDate,
		ToDate:   toDate,
		Lines:    make([]domain.CustomerStatementLine, 0, len(items)),
	}
	bal := 0.0
	for _, it := range items {
		bal += it.debit - it.credit
		stmt.Lines = append(stmt.Lines, domain.CustomerStatementLine{
			Date: it.date, RefType: it.refType, RefNumber: it.refNum,
			Description: it.desc, Debit: it.debit, Credit: it.credit, Balance: bal,
		})
	}
	stmt.ClosingBal = bal
	return stmt, nil
}

func (s *SaleService) GetARGLReconciliation(ctx context.Context, companyID, periodID string) (*domain.ARGLReconciliation, error) {
	invoices, _, err := s.invRepo.ListInvoices(ctx, domain.CustomerInvoiceFilter{CompanyID: companyID})
	if err != nil {
		return nil, err
	}
	subledgerTotal := 0.0
	byCust := make(map[string]*domain.ARGLReconDetail)
	for _, inv := range invoices {
		if inv.Status == domain.SInvCancelled {
			continue
		}
		subledgerTotal += inv.BalanceDue
		d, ok := byCust[inv.CustomerID]
		if !ok {
			d = &domain.ARGLReconDetail{CustomerID: inv.CustomerID, CustomerName: inv.CustomerName}
			byCust[inv.CustomerID] = d
		}
		d.BalanceDue += inv.BalanceDue
	}
	glBalance := 0.0
	if s.gl != nil {
		bal, err := s.gl.GetAccountBalance(ctx, "131", periodID)
		if err == nil && bal != nil {
			glBalance = bal.ClosingBalance
		}
	}
	details := make([]domain.ARGLReconDetail, 0, len(byCust))
	for _, d := range byCust {
		details = append(details, *d)
	}
	return &domain.ARGLReconciliation{
		PeriodID:       periodID,
		SubledgerTotal: subledgerTotal,
		GLBalance:      glBalance,
		Variance:       glBalance - subledgerTotal,
		Details:        details,
	}, nil
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
