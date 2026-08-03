package service

import (
	"context"
	"fmt"
	"math"
	"gotax/internal/domain"
	"time"
)

// APGLService posts purchase documents to the GL journal.
// Satisfied by service.Service (CreatePostedEntry).
type APGLService interface {
	CreatePostedEntry(ctx context.Context, entry *domain.JournalEntry, userID string) error
}

const apPayableAccount = "331"

type PurchaseService struct {
	supRepo  domain.SupplierRepository
	poRepo   domain.PurchaseOrderRepository
	grnRepo  domain.GRNRepository
	invRepo  domain.SupplierInvoiceRepository
	aptRepo  domain.APTransactionRepository
	costRepo domain.CostAllocationRepository
	gl       APGLService
	now      func() time.Time
}

func NewPurchaseService(
	supRepo domain.SupplierRepository,
	poRepo domain.PurchaseOrderRepository,
	grnRepo domain.GRNRepository,
	invRepo domain.SupplierInvoiceRepository,
	aptRepo domain.APTransactionRepository,
	costRepo domain.CostAllocationRepository,
	gl APGLService,
) *PurchaseService {
	return &PurchaseService{
		supRepo: supRepo, poRepo: poRepo, grnRepo: grnRepo,
		invRepo: invRepo, aptRepo: aptRepo, costRepo: costRepo,
		gl:  gl,
		now: time.Now,
	}
}

// ─── Supplier ───────────────────────────────────────────────────────────

func (s *PurchaseService) CreateSupplier(ctx context.Context, sup *domain.Supplier) error {
	existing, _ := s.supRepo.GetSupplierByCode(ctx, sup.CompanyID, sup.Code)
	if existing != nil {
		return domain.ErrSupplierCodeExists
	}
	if sup.Currency == "" {
		sup.Currency = "VND"
	}
	if sup.Status == "" {
		sup.Status = domain.SupplierActive
	}
	sup.CreatedAt = s.now()
	sup.UpdatedAt = sup.CreatedAt
	return s.supRepo.CreateSupplier(ctx, sup)
}

func (s *PurchaseService) GetSupplier(ctx context.Context, id string) (*domain.Supplier, error) {
	return s.supRepo.GetSupplier(ctx, id)
}

func (s *PurchaseService) ListSuppliers(ctx context.Context, companyID string, offset, limit int) ([]domain.Supplier, int, error) {
	return s.supRepo.ListSuppliers(ctx, domain.PurchaseOrderFilter{CompanyID: companyID, Limit: limit, Offset: offset})
}

func (s *PurchaseService) UpdateSupplier(ctx context.Context, sup *domain.Supplier) error {
	existing, err := s.supRepo.GetSupplier(ctx, sup.ID)
	if err != nil {
		return err
	}
	if existing.Code != sup.Code {
		dup, _ := s.supRepo.GetSupplierByCode(ctx, sup.CompanyID, sup.Code)
		if dup != nil {
			return domain.ErrSupplierCodeExists
		}
	}
	if sup.Currency == "" {
		sup.Currency = "VND"
	}
	sup.UpdatedAt = s.now()
	return s.supRepo.UpdateSupplier(ctx, sup)
}

func (s *PurchaseService) DeleteSupplier(ctx context.Context, id string) error {
	return s.supRepo.DeleteSupplier(ctx, id)
}

// ─── Purchase Order ─────────────────────────────────────────────────────

func (s *PurchaseService) CreatePO(ctx context.Context, po *domain.PurchaseOrder) error {
	if err := po.Validate(); err != nil {
		return err
	}
	po.CalculateTotals()
	if po.Status == "" {
		po.Status = domain.POStatusDraft
	}
	if _, err := s.supRepo.GetSupplier(ctx, po.SupplierID); err != nil {
		return domain.ErrPOSupplierRequired
	}
	po.CreatedAt = s.now()
	po.UpdatedAt = po.CreatedAt
	if err := s.poRepo.CreatePO(ctx, po); err != nil {
		return err
	}
	for i := range po.Lines {
		po.Lines[i].POID = po.ID
	}
	return s.poRepo.CreatePOLines(ctx, po.Lines)
}

func (s *PurchaseService) GetPO(ctx context.Context, id string) (*domain.PurchaseOrder, error) {
	return s.poRepo.GetPO(ctx, id)
}

func (s *PurchaseService) GetPOByNumber(ctx context.Context, companyID, poNumber string) (*domain.PurchaseOrder, error) {
	return s.poRepo.GetPOByNumber(ctx, companyID, poNumber)
}

func (s *PurchaseService) ListPOs(ctx context.Context, filter domain.PurchaseOrderFilter) ([]domain.PurchaseOrder, int, error) {
	return s.poRepo.ListPOs(ctx, filter)
}

func (s *PurchaseService) UpdatePO(ctx context.Context, po *domain.PurchaseOrder) error {
	existing, err := s.poRepo.GetPO(ctx, po.ID)
	if err != nil {
		return err
	}
	if existing.Status != domain.POStatusDraft {
		return domain.ErrPOCannotUpdate
	}
	po.CalculateTotals()
	po.UpdatedAt = s.now()
	return s.poRepo.UpdatePO(ctx, po)
}

func (s *PurchaseService) ApprovePO(ctx context.Context, id, approvedBy string) error {
	po, err := s.poRepo.GetPO(ctx, id)
	if err != nil {
		return err
	}
	if !po.Status.ValidTransition(domain.POStatusApproved) {
		return domain.ErrPOInvalidTransition
	}
	return s.poRepo.ApprovePO(ctx, id, approvedBy, s.now())
}

func (s *PurchaseService) CancelPO(ctx context.Context, id, reason string) error {
	po, err := s.poRepo.GetPO(ctx, id)
	if err != nil {
		return err
	}
	if !po.Status.ValidTransition(domain.POStatusCancelled) {
		return domain.ErrPOInvalidTransition
	}
	return s.poRepo.CancelPO(ctx, id, reason)
}

func (s *PurchaseService) ClosePO(ctx context.Context, id string) error {
	po, err := s.poRepo.GetPO(ctx, id)
	if err != nil {
		return err
	}
	if po.Status != domain.POStatusReceived {
		return domain.ErrPOInvalidTransition
	}
	po.Status = domain.POStatusClosed
	po.UpdatedAt = s.now()
	return s.poRepo.UpdatePOStatus(ctx, id, domain.POStatusClosed)
}

// ─── GRN ─────────────────────────────────────────────────────────────────

func (s *PurchaseService) CreateGRN(ctx context.Context, grn *domain.GRN) error {
	if err := grn.Validate(); err != nil {
		return err
	}
	if _, err := s.poRepo.GetPO(ctx, grn.POID); err != nil {
		return domain.ErrGRNPORequired
	}
	if grn.Status == "" {
		grn.Status = domain.GRNDraft
	}
	grn.CreatedAt = s.now()
	if err := s.grnRepo.CreateGRN(ctx, grn); err != nil {
		return err
	}
	for i := range grn.Lines {
		grn.Lines[i].GRNID = grn.ID
	}
	return s.grnRepo.CreateGRNLines(ctx, grn.Lines)
}

func (s *PurchaseService) GetGRN(ctx context.Context, id string) (*domain.GRN, error) {
	return s.grnRepo.GetGRN(ctx, id)
}

func (s *PurchaseService) ListGRNs(ctx context.Context, filter domain.GRNFilter) ([]domain.GRN, int, error) {
	return s.grnRepo.ListGRNs(ctx, filter)
}

func (s *PurchaseService) UpdateGRN(ctx context.Context, grn *domain.GRN) error {
	existing, err := s.grnRepo.GetGRN(ctx, grn.ID)
	if err != nil {
		return err
	}
	if existing.Status != domain.GRNDraft {
		return domain.ErrGRNCannotUpdate
	}
	return s.grnRepo.UpdateGRN(ctx, grn)
}

func (s *PurchaseService) PostGRN(ctx context.Context, id string) error {
	grn, err := s.grnRepo.GetGRN(ctx, id)
	if err != nil {
		return err
	}
	if grn.Status != domain.GRNDraft {
		return domain.ErrGRNInvalidTransition
	}
	return s.grnRepo.UpdateGRNStatus(ctx, id, domain.GRNPosted)
}

func (s *PurchaseService) CancelGRN(ctx context.Context, id string) error {
	grn, err := s.grnRepo.GetGRN(ctx, id)
	if err != nil {
		return err
	}
	if grn.Status == domain.GRNPosted {
		return domain.ErrGRNInvalidTransition
	}
	return s.grnRepo.UpdateGRNStatus(ctx, id, domain.GRNCancelled)
}

// ─── Supplier Invoice ───────────────────────────────────────────────────

func (s *PurchaseService) CreateInvoice(ctx context.Context, inv *domain.SupplierInvoice) error {
	sup, err := s.supRepo.GetSupplier(ctx, inv.SupplierID)
	if err != nil {
		return domain.ErrInvoiceSupplierRequired
	}
	inv.SupplierName = sup.Name
	inv.SupplierTaxCode = sup.TaxCode
	if err := inv.Validate(); err != nil {
		return err
	}
	if inv.Currency == "" {
		inv.Currency = "VND"
	}
	if inv.Status == "" {
		inv.Status = domain.InvoiceDraft
	}
	if inv.VATDeductionStatus == "" {
		inv.VATDeductionStatus = domain.VATPending
	}
	inv.CalculateTotals()
	inv.CreatedAt = s.now()
	if err := s.invRepo.CreateInvoice(ctx, inv); err != nil {
		return err
	}
	for i := range inv.Lines {
		inv.Lines[i].InvoiceID = inv.ID
	}
	return s.invRepo.CreateInvoiceLines(ctx, inv.Lines)
}

func (s *PurchaseService) GetInvoice(ctx context.Context, id string) (*domain.SupplierInvoice, error) {
	return s.invRepo.GetInvoice(ctx, id)
}

func (s *PurchaseService) ListInvoices(ctx context.Context, filter domain.SupplierInvoiceFilter) ([]domain.SupplierInvoice, int, error) {
	return s.invRepo.ListInvoices(ctx, filter)
}

func (s *PurchaseService) UpdateInvoice(ctx context.Context, inv *domain.SupplierInvoice) error {
	existing, err := s.invRepo.GetInvoice(ctx, inv.ID)
	if err != nil {
		return err
	}
	if existing.Status != domain.InvoiceDraft {
		return domain.ErrInvoiceCannotUpdate
	}
	inv.CalculateTotals()
	return s.invRepo.UpdateInvoice(ctx, inv)
}

func (s *PurchaseService) VerifyInvoice(ctx context.Context, id string) error {
	inv, err := s.invRepo.GetInvoice(ctx, id)
	if err != nil {
		return err
	}
	if inv.Status != domain.InvoiceDraft {
		return domain.ErrInvoiceInvalidTransition
	}
	return s.invRepo.UpdateInvoiceStatus(ctx, id, domain.InvoiceVerified)
}

func (s *PurchaseService) PostInvoice(ctx context.Context, id string) error {
	inv, err := s.invRepo.GetInvoice(ctx, id)
	if err != nil {
		return err
	}
	if inv.Status != domain.InvoiceVerified {
		return domain.ErrInvoiceInvalidTransition
	}
	if err := s.verifyThreeWayMatch(ctx, inv); err != nil {
		return err
	}
	now := s.now().UTC()
	if s.gl != nil {
		entry := s.buildInvoiceGLEntry(inv)
		if err := s.gl.CreatePostedEntry(ctx, entry, inv.CreatedBy); err != nil {
			return err
		}
		if err := s.invRepo.SetInvoiceGLPosted(ctx, id, now); err != nil {
			return err
		}
	}
	return s.invRepo.PostInvoice(ctx, id, now)
}

// verifyThreeWayMatch enforces PO qty >= GRN qty >= Invoice qty per line and a
// price variance within tolerance. When the invoice is not linked to a PO the
// match is skipped.
func (s *PurchaseService) verifyThreeWayMatch(ctx context.Context, inv *domain.SupplierInvoice) error {
	if inv.POID == "" {
		return nil
	}
	poLines, err := s.poRepo.GetPOLines(ctx, inv.POID)
	if err != nil {
		return err
	}
	poByID := make(map[string]domain.POItem, len(poLines))
	for _, l := range poLines {
		poByID[l.ID] = l
	}
	grnByPOLine := make(map[string]domain.GRNItem)
	if inv.GRNID != "" {
		grnLines, err := s.grnRepo.GetGRNLines(ctx, inv.GRNID)
		if err != nil {
			return err
		}
		for _, l := range grnLines {
			grnByPOLine[l.POLineID] = l
		}
	}
	const tolerancePct = 5.0
	for _, il := range inv.Lines {
		if il.POLineID == "" {
			continue
		}
		pol, ok := poByID[il.POLineID]
		if !ok {
			return domain.ErrInvoice3WayMismatch
		}
		if il.Quantity > pol.Quantity*(1+tolerancePct/100) {
			return domain.ErrInvoice3WayMismatch
		}
		if g, ok := grnByPOLine[il.POLineID]; ok && il.Quantity > g.QuantityReceived*(1+tolerancePct/100) {
			return domain.ErrInvoice3WayMismatch
		}
		if pol.UnitPrice > 0 {
			variance := math.Abs(il.UnitPrice-pol.UnitPrice) / pol.UnitPrice
			if variance > tolerancePct/100 {
				return domain.ErrInvoice3WayMismatch
			}
		}
	}
	return nil
}

func (s *PurchaseService) buildInvoiceGLEntry(inv *domain.SupplierInvoice) *domain.JournalEntry {
	expense := make(map[string]float64)
	vat := make(map[string]float64)
	for _, l := range inv.Lines {
		expense[l.AccountID] += l.LineTotal
		vat[l.VATAccountID] += l.LineVATAmount
	}
	lines := make([]domain.JournalLine, 0, len(expense)+len(vat)+1)
	for acc, amt := range expense {
		lines = append(lines, domain.JournalLine{
			AccountCode: acc, DebitAmount: amt, Description: "Purchase: " + inv.InvoiceNumber,
		})
	}
	for acc, amt := range vat {
		lines = append(lines, domain.JournalLine{
			AccountCode: acc, DebitAmount: amt, Description: "VAT input: " + inv.InvoiceNumber,
		})
	}
	lines = append(lines, domain.JournalLine{
		AccountCode: apPayableAccount, CreditAmount: inv.TotalAmount, Description: "AP: " + inv.InvoiceNumber,
	})
	return &domain.JournalEntry{
		CompanyID:   inv.CompanyID,
		EntryNumber: inv.InvoiceNumber,
		VoucherType: domain.VoucherTypePurchase,
		EntryDate:   inv.InvoiceDate,
		Description: "AP invoice " + inv.InvoiceNumber,
		Lines:       lines,
	}
}

func (s *PurchaseService) CancelInvoice(ctx context.Context, id string) error {
	inv, err := s.invRepo.GetInvoice(ctx, id)
	if err != nil {
		return err
	}
	if inv.Status == domain.InvoicePosted || inv.Status == domain.InvoicePaid {
		return domain.ErrInvoiceInvalidTransition
	}
	return s.invRepo.UpdateInvoiceStatus(ctx, id, domain.InvoiceCancelled)
}

func (s *PurchaseService) ClaimVAT(ctx context.Context, id string) error {
	inv, err := s.invRepo.GetInvoice(ctx, id)
	if err != nil {
		return err
	}
	if inv.VATDeductionStatus == domain.VATRejected {
		return fmt.Errorf("VAT already rejected for invoice %s", inv.InvoiceNumber)
	}
	inv.VATDeductionStatus = domain.VATClaimed
	return s.invRepo.UpdateInvoice(ctx, inv)
}

func (s *PurchaseService) ReceiveEInvoice(ctx context.Context, inv *domain.SupplierInvoice) error {
	if err := inv.Validate(); err != nil {
		return err
	}
	dup, _ := s.invRepo.GetInvoiceByNumber(ctx, inv.CompanyID, inv.InvoiceNumber)
	if dup != nil {
		return fmt.Errorf("invoice %s already exists", inv.InvoiceNumber)
	}
	sup, _ := s.supRepo.GetSupplierByCode(ctx, inv.CompanyID, inv.SupplierTaxCode)
	if sup == nil {
		if err := s.CreateSupplier(ctx, &domain.Supplier{
			CompanyID: inv.CompanyID, Code: fmt.Sprintf("SUP-%s", inv.SupplierTaxCode),
			Name: inv.SupplierName, TaxCode: inv.SupplierTaxCode,
			Status: domain.SupplierActive, Currency: inv.Currency,
			CreatedAt: s.now(), UpdatedAt: s.now(),
		}); err != nil {
			return err
		}
		sup, _ = s.supRepo.GetSupplierByCode(ctx, inv.CompanyID, fmt.Sprintf("SUP-%s", inv.SupplierTaxCode))
	}
	inv.SupplierID = sup.ID
	inv.Status = domain.InvoiceDraft
	inv.VATDeductionStatus = domain.VATPending
	inv.CreatedAt = s.now()
	if err := s.invRepo.CreateInvoice(ctx, inv); err != nil {
		return err
	}
	for i := range inv.Lines {
		inv.Lines[i].InvoiceID = inv.ID
	}
	return s.invRepo.CreateInvoiceLines(ctx, inv.Lines)
}

// ─── AP Transactions ────────────────────────────────────────────────────

func (s *PurchaseService) CreateAPTransaction(ctx context.Context, t *domain.APTransaction) error {
	if err := t.Validate(); err != nil {
		return err
	}
	if _, err := s.supRepo.GetSupplier(ctx, t.SupplierID); err != nil {
		return domain.ErrAPTransSupplierRequired
	}
	if t.Currency == "" {
		t.Currency = "VND"
	}
	t.CreatedAt = s.now()
	return s.aptRepo.CreateAPTransaction(ctx, t)
}

func (s *PurchaseService) GetAPTransaction(ctx context.Context, id string) (*domain.APTransaction, error) {
	return s.aptRepo.GetAPTransaction(ctx, id)
}

func (s *PurchaseService) ListAPTransactions(ctx context.Context, companyID string, offset, limit int) ([]domain.APTransaction, int, error) {
	return s.aptRepo.ListAPTransactions(ctx, companyID, offset, limit)
}

func (s *PurchaseService) ListAPTransactionsBySupplier(ctx context.Context, companyID, supplierID string) ([]domain.APTransaction, error) {
	return s.aptRepo.ListAPTransactionsBySupplier(ctx, companyID, supplierID)
}

func (s *PurchaseService) GetAPAgingReport(ctx context.Context, companyID string) ([]domain.APAgingReport, error) {
	txns, _, err := s.aptRepo.ListAPTransactions(ctx, companyID, 0, 0)
	if err != nil {
		return nil, err
	}
	bySup := make(map[string]*domain.APAgingReport)
	supBal := make(map[string]float64)
	for _, t := range txns {
		supBal[t.SupplierID] += t.Amount
	}
	for sid, bal := range supBal {
		rpt := &domain.APAgingReport{SupplierID: sid, Buckets: domain.APAgingBucket{Bucket0: bal, Bucket30: 0, Bucket60: 0, Bucket90: 0, Total: bal}}
		bySup[sid] = rpt
	}
	out := make([]domain.APAgingReport, 0, len(bySup))
	for _, r := range bySup {
		out = append(out, *r)
	}
	return out, nil
}

func (s *PurchaseService) GetAPSummary(ctx context.Context, companyID string) ([]domain.APSummary, error) {
	txns, _, err := s.aptRepo.ListAPTransactions(ctx, companyID, 0, 0)
	if err != nil {
		return nil, err
	}
	supIDs := make([]string, 0, len(txns))
	seen := make(map[string]bool)
	for _, t := range txns {
		if !seen[t.SupplierID] {
			supIDs = append(supIDs, t.SupplierID)
			seen[t.SupplierID] = true
		}
	}
	allSuppliers, _ := s.supRepo.ListSuppliersByIDs(ctx, supIDs)
	supMap := make(map[string]domain.Supplier, len(allSuppliers))
	for i := range allSuppliers {
		supMap[allSuppliers[i].ID] = allSuppliers[i]
	}
	m := make(map[string]*domain.APSummary)
	for _, t := range txns {
		sup, ok := supMap[t.SupplierID]
		if !ok {
			continue
		}
		if _, ok := m[t.SupplierID]; !ok {
			m[t.SupplierID] = &domain.APSummary{
				SupplierID: t.SupplierID, SupplierName: sup.Name,
				TaxCode: sup.TaxCode, Currency: t.Currency,
			}
		}
		cur := m[t.SupplierID]
		switch t.TransactionType {
		case domain.APTransInvoice:
			cur.TotalInvoiced += t.Amount
		case domain.APTransPayment, domain.APTransCreditNote, domain.APTransOffset:
			cur.TotalPaid += t.Amount
		}
		cur.Outstanding = cur.TotalInvoiced - cur.TotalPaid
	}
	out := make([]domain.APSummary, 0, len(m))
	for _, v := range m {
		out = append(out, *v)
	}
	return out, nil
}

// ─── Cost Allocation ────────────────────────────────────────────────────

func (s *PurchaseService) CreateCostAllocation(ctx context.Context, c *domain.CostAllocation) error {
	if err := c.Validate(); err != nil {
		return err
	}
	return s.costRepo.CreateCostAllocation(ctx, c)
}

func (s *PurchaseService) GetCostAllocation(ctx context.Context, id string) (*domain.CostAllocation, error) {
	return s.costRepo.GetCostAllocation(ctx, id)
}

func (s *PurchaseService) ListCostAllocationsByInvoice(ctx context.Context, invoiceID string) ([]domain.CostAllocation, error) {
	return s.costRepo.ListCostAllocationsByInvoice(ctx, invoiceID)
}
