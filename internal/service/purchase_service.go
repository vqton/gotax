package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"

	"gotax/internal/domain"
	"gotax/internal/einvoice"
	"gotax/internal/validate"
)

// APGLService posts purchase documents to the GL journal.
// Satisfied by service.Service (CreatePostedEntry, GetExchangeRate).
type APGLService interface {
	CreatePostedEntry(ctx context.Context, entry *domain.JournalEntry, userID string) error
	GetExchangeRate(ctx context.Context, currencyCode string, rateDate time.Time) (*domain.ExchangeRate, error)
}

const apPayableAccount = "331"
const importDutyAccount = "3333"
const importVATAccount = "33312"
const fxGainAccount = "515"
const fxLossAccount = "635"

type PurchaseService struct {
	supRepo  domain.SupplierRepository
	poRepo   domain.PurchaseOrderRepository
	grnRepo  domain.GRNRepository
	invRepo  domain.SupplierInvoiceRepository
	aptRepo  domain.APTransactionRepository
	costRepo domain.CostAllocationRepository
	provRepo domain.DoubtfulDebtProvisionRepository
	reqRepo  domain.RequisitionRepository
	fxRepo   domain.FXRevaluationRepository
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
	provRepo domain.DoubtfulDebtProvisionRepository,
	reqRepo domain.RequisitionRepository,
	fxRepo domain.FXRevaluationRepository,
	gl APGLService,
) *PurchaseService {
	return &PurchaseService{
		supRepo: supRepo, poRepo: poRepo, grnRepo: grnRepo,
		invRepo: invRepo, aptRepo: aptRepo, costRepo: costRepo,
		provRepo: provRepo, reqRepo: reqRepo, fxRepo: fxRepo,
		gl:  gl,
		now: time.Now,
	}
}

// ─── Supplier ───────────────────────────────────────────────────────────

func (s *PurchaseService) CreateSupplier(ctx context.Context, sup *domain.Supplier) error {
	// PG suppliers.id is TEXT PRIMARY KEY with no DB default — without an
	// explicit ID the row is created with id='' and by-id lookups miss it
	// (same defect as customers, fixed in sale service).
	if sup.ID == "" {
		sup.ID = uuid.NewString()
	}
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
	if err := validate.PurchaseOrder(po); err != nil {
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
	if err := validate.GRN(grn); err != nil {
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
	if err := validate.SupplierInvoice(inv); err != nil {
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
	isCreditNote := inv.InvoiceType == domain.InvoiceTypeCreditNote
	if isCreditNote {
		if err := s.verifyCreditNote(ctx, inv); err != nil {
			return err
		}
	} else if err := s.verifyThreeWayMatch(ctx, inv); err != nil {
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
	if err := s.invRepo.PostInvoice(ctx, id, now); err != nil {
		return err
	}
	if isCreditNote {
		return s.applyCreditNote(ctx, inv, now)
	}
	return nil
}

func (s *PurchaseService) verifyCreditNote(ctx context.Context, inv *domain.SupplierInvoice) error {
	original, err := s.invRepo.GetInvoice(ctx, inv.OriginalInvoiceID)
	if err != nil {
		return domain.ErrCreditNoteOriginalRequired
	}
	if original.Status != domain.InvoicePosted && original.Status != domain.InvoicePaid {
		return domain.ErrCreditNoteOriginalNotPosted
	}
	if original.SupplierID != inv.SupplierID {
		return domain.ErrCreditNoteSupplierMismatch
	}
	if math.Abs(inv.TotalAmount) > original.BalanceDue+0.001 {
		return domain.ErrCreditNoteExceedsBalance
	}
	return nil
}

func (s *PurchaseService) applyCreditNote(ctx context.Context, inv *domain.SupplierInvoice, now time.Time) error {
	amount := math.Abs(inv.TotalAmount)
	original, err := s.invRepo.GetInvoice(ctx, inv.OriginalInvoiceID)
	if err != nil {
		return err
	}
	if err := s.invRepo.ReduceInvoiceBalance(ctx, original.ID, amount); err != nil {
		return err
	}
	txn := &domain.APTransaction{
		CompanyID:       inv.CompanyID,
		SupplierID:      inv.SupplierID,
		InvoiceID:       original.ID,
		TransactionType: domain.APTransCreditNote,
		TransactionDate: now,
		Amount:          amount,
		Currency:        inv.Currency,
		ReferenceType:   "invoice",
		ReferenceID:     inv.ID,
		CreatedAt:       now,
	}
	return s.aptRepo.CreateAPTransaction(ctx, txn)
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
	creditNote := inv.InvoiceType == domain.InvoiceTypeCreditNote
	rate := 1.0
	if inv.Currency != "" && inv.Currency != "VND" && inv.ExchangeRate > 0 {
		rate = inv.ExchangeRate
	}
	for _, l := range inv.Lines {
		expense[l.AccountID] += l.LineTotal * rate
		vat[l.VATAccountID] += l.LineVATAmount * rate
	}
	lines := make([]domain.JournalLine, 0, len(expense)+len(vat)+1)
	for acc, amt := range expense {
		amt = math.Abs(amt)
		if amt <= 0.001 {
			continue
		}
		if creditNote {
			lines = append(lines, domain.JournalLine{
				AccountCode: acc, CreditAmount: amt, Description: "Credit note: " + inv.InvoiceNumber,
			})
		} else {
			lines = append(lines, domain.JournalLine{
				AccountCode: acc, DebitAmount: amt, Description: "Purchase: " + inv.InvoiceNumber,
			})
		}
	}
	for acc, amt := range vat {
		amt = math.Abs(amt)
		if amt <= 0.001 {
			continue
		}
		if creditNote {
			lines = append(lines, domain.JournalLine{
				AccountCode: acc, CreditAmount: amt, Description: "VAT input reversal: " + inv.InvoiceNumber,
			})
		} else {
			lines = append(lines, domain.JournalLine{
				AccountCode: acc, DebitAmount: amt, Description: "VAT input: " + inv.InvoiceNumber,
			})
		}
	}
	total := math.Abs(inv.TotalAmount) * rate
	if creditNote {
		lines = append(lines, domain.JournalLine{
			AccountCode: apPayableAccount, DebitAmount: total, Description: "AP reversal: " + inv.InvoiceNumber,
		})
	} else {
		lines = append(lines, domain.JournalLine{
			AccountCode: apPayableAccount, CreditAmount: total, Description: "AP: " + inv.InvoiceNumber,
		})
	}
	if inv.InvoiceType == domain.InvoiceTypeImport {
		inventoryAcc, vatAcc := "152", "1331"
		if len(inv.Lines) > 0 {
			if inv.Lines[0].AccountID != "" {
				inventoryAcc = inv.Lines[0].AccountID
			}
			if inv.Lines[0].VATAccountID != "" {
				vatAcc = inv.Lines[0].VATAccountID
			}
		}
		if inv.ImportDuty > 0 {
			lines = append(lines, domain.JournalLine{
				AccountCode: inventoryAcc, DebitAmount: inv.ImportDuty, Description: "Import duty: " + inv.InvoiceNumber,
			})
			lines = append(lines, domain.JournalLine{
				AccountCode: importDutyAccount, CreditAmount: inv.ImportDuty, Description: "Customs duty payable: " + inv.InvoiceNumber,
			})
		}
		if inv.ImportVAT > 0 {
			lines = append(lines, domain.JournalLine{
				AccountCode: vatAcc, DebitAmount: inv.ImportVAT, Description: "Import VAT input: " + inv.InvoiceNumber,
			})
			lines = append(lines, domain.JournalLine{
				AccountCode: importVATAccount, CreditAmount: inv.ImportVAT, Description: "Import VAT payable: " + inv.InvoiceNumber,
			})
		}
	}
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
		var err error
		sup, err = s.supRepo.GetSupplierByCode(ctx, inv.CompanyID, fmt.Sprintf("SUP-%s", inv.SupplierTaxCode))
		if err != nil {
			return fmt.Errorf("lookup supplier after create: %w", err)
		}
	}
	inv.SupplierID = sup.ID
	if err := validate.SupplierInvoice(inv); err != nil {
		return err
	}
	dup, err := s.invRepo.GetInvoiceByNumber(ctx, inv.CompanyID, inv.InvoiceNumber)
	if err != nil && !errors.Is(err, domain.ErrSupplierInvoiceNotFound) {
		return fmt.Errorf("check duplicate invoice: %w", err)
	}
	if dup != nil {
		return fmt.Errorf("invoice %s already exists", inv.InvoiceNumber)
	}
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
	if err := validate.APTransaction(t); err != nil {
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
	if err := validate.CostAllocation(c); err != nil {
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

// ─── Doubtful Debt Provisions (Circular 99) ─────────────────────────────

// CalculateDoubtfulDebtProvision computes provision lines from outstanding
// supplier prepayments aged by months overdue. Circular 99 tiers: 30/50/70/100%.
func (s *PurchaseService) CalculateDoubtfulDebtProvision(ctx context.Context, companyID string, asOfDate time.Time) ([]domain.DoubtfulDebtProvisionLine, error) {
	txns, _, err := s.aptRepo.ListAPTransactions(ctx, companyID, 0, 0)
	if err != nil {
		return nil, err
	}
	bal := make(map[string]float64)
	oldest := make(map[string]time.Time)
	for _, t := range txns {
		if t.TransactionType != domain.APTransPrepayment {
			continue
		}
		bal[t.SupplierID] += t.Amount
		if d, ok := oldest[t.SupplierID]; !ok || t.TransactionDate.Before(d) {
			oldest[t.SupplierID] = t.TransactionDate
		}
	}
	if len(bal) == 0 {
		return nil, domain.ErrProvisionNoPrepayments
	}

	supIDs := make([]string, 0, len(bal))
	for sid := range bal {
		supIDs = append(supIDs, sid)
	}
	suppliers, _ := s.supRepo.ListSuppliersByIDs(ctx, supIDs)
	supMap := make(map[string]domain.Supplier, len(suppliers))
	for i := range suppliers {
		supMap[suppliers[i].ID] = suppliers[i]
	}

	lines := make([]domain.DoubtfulDebtProvisionLine, 0, len(bal))
	for sid, outstanding := range bal {
		age := ageInMonths(oldest[sid], asOfDate)
		rate := domain.DoubtfulDebtRate(age)
		if rate <= 0 {
			continue
		}
		line := domain.DoubtfulDebtProvisionLine{
			SupplierID:        sid,
			OutstandingAmount: round2(outstanding),
			AgeMonths:         age,
			RatePct:           rate,
			ProvisionAmount:   round2(outstanding * rate),
		}
		if sup, ok := supMap[sid]; ok {
			line.SupplierName = sup.Name
			line.TaxCode = sup.TaxCode
		}
		lines = append(lines, line)
	}
	return lines, nil
}

func (s *PurchaseService) CreateDoubtfulDebtProvision(ctx context.Context, p *domain.DoubtfulDebtProvision) error {
	if p.AsOfDate == "" {
		return domain.ErrProvisionDateRequired
	}
	if _, err := time.Parse(time.DateOnly, p.AsOfDate); err != nil {
		return domain.ErrProvisionDateRequired
	}
	if p.Status == "" {
		p.Status = domain.ProvisionDraft
	}
	var totalOut, totalProv float64
	for i := range p.Lines {
		if p.Lines[i].RatePct <= 0 {
			p.Lines[i].RatePct = domain.DoubtfulDebtRate(p.Lines[i].AgeMonths)
		}
		p.Lines[i].ProvisionAmount = round2(p.Lines[i].OutstandingAmount * p.Lines[i].RatePct)
		totalOut += p.Lines[i].OutstandingAmount
		totalProv += p.Lines[i].ProvisionAmount
	}
	if len(p.Lines) == 0 {
		return domain.ErrProvisionNoLines
	}
	p.TotalOutstanding = round2(totalOut)
	p.TotalProvision = round2(totalProv)
	p.CreatedAt = s.now()
	if err := s.provRepo.CreateProvision(ctx, p); err != nil {
		return err
	}
	for i := range p.Lines {
		p.Lines[i].ProvisionID = p.ID
	}
	return s.provRepo.CreateProvisionLines(ctx, p.Lines)
}

func (s *PurchaseService) GetDoubtfulDebtProvision(ctx context.Context, id string) (*domain.DoubtfulDebtProvision, error) {
	return s.provRepo.GetProvision(ctx, id)
}

func (s *PurchaseService) ListDoubtfulDebtProvisions(ctx context.Context, companyID string, offset, limit int) ([]domain.DoubtfulDebtProvision, int, error) {
	return s.provRepo.ListProvisions(ctx, companyID, limit, offset)
}

// ─── Regulatory Reports (Circular 99) ───────────────────────────────────

// GetPurchaseLedger builds S01-DN (sổ chi tiết mua hàng): goods-receipt
// postings in the period as increases, plus opening/closing balances.
func (s *PurchaseService) GetPurchaseLedger(ctx context.Context, companyID, fromDate, toDate string) (*domain.PurchaseLedgerReport, error) {
	rpt := &domain.PurchaseLedgerReport{CompanyID: companyID, FromDate: fromDate, ToDate: toDate}

	open, err := s.grnTotalBefore(ctx, companyID, fromDate)
	if err != nil {
		return nil, err
	}
	rpt.Opening = open

	grns, _, err := s.grnRepo.ListGRNs(ctx, domain.GRNFilter{
		CompanyID: companyID, FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		return nil, err
	}
	for _, g := range grns {
		if g.Status != domain.GRNPosted {
			continue
		}
		supplierName := s.supplierNameForGRN(ctx, g)
		inc := grnTotal(&g)
		rpt.Increase += inc
		rpt.Rows = append(rpt.Rows, domain.PurchaseLedgerRow{
			Date: g.ReceiptDate.Format(time.DateOnly), DocNumber: g.GRNNumber,
			Description: "Nhập kho", SupplierID: "", SupplierName: supplierName,
			Increase: round2(inc),
		})
	}
	rpt.Closing = round2(rpt.Opening + rpt.Increase - rpt.Decrease)
	return rpt, nil
}

// GetSupplierLedger builds S02-DN (sổ chi tiết công nợ phải trả) for a supplier.
func (s *PurchaseService) GetSupplierLedger(ctx context.Context, companyID, supplierID, fromDate, toDate string) (*domain.SupplierLedgerReport, error) {
	sup, err := s.supRepo.GetSupplier(ctx, supplierID)
	if err != nil {
		return nil, err
	}
	rpt := &domain.SupplierLedgerReport{
		SupplierID: sup.ID, SupplierName: sup.Name, TaxCode: sup.TaxCode,
		FromDate: fromDate, ToDate: toDate,
	}

	all, _, err := s.aptRepo.ListAPTransactions(ctx, companyID, 0, 0)
	if err != nil {
		return nil, err
	}
	f := func(dt time.Time) bool { return !dt.IsZero() }
	fromT, _ := time.Parse(time.DateOnly, fromDate)
	toT, _ := time.Parse(time.DateOnly, toDate)

	balance := 0.0
	for _, t := range all {
		if t.SupplierID != supplierID {
			continue
		}
		if !f(t.TransactionDate) {
			continue
		}
		if t.TransactionDate.Before(fromT) {
			balance += ledgerDelta(t)
			continue
		}
		if t.TransactionDate.After(toT) {
			continue
		}
		delta := ledgerDelta(t)
		balance += delta
		rpt.Rows = append(rpt.Rows, domain.SupplierLedgerRow{
			Date: t.TransactionDate.Format(time.DateOnly), DocNumber: t.ReferenceID,
			Description: string(t.TransactionType), Debit: ledgerDebit(t), Credit: ledgerCredit(t),
			Balance: round2(balance),
		})
	}
	rpt.Opening = round2(rpt.Opening)
	rpt.Closing = round2(balance)
	return rpt, nil
}

// GetGoodsPurchaseReport builds S03-DN (sổ chi tiết hàng hóa): line-level goods purchases.
func (s *PurchaseService) GetGoodsPurchaseReport(ctx context.Context, companyID, fromDate, toDate string) (*domain.GoodsPurchaseReport, error) {
	rpt := &domain.GoodsPurchaseReport{CompanyID: companyID, FromDate: fromDate, ToDate: toDate}
	invs, _, err := s.invRepo.ListInvoices(ctx, domain.SupplierInvoiceFilter{
		CompanyID: companyID, FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		return nil, err
	}
	for _, inv := range invs {
		if inv.Status == domain.InvoiceCancelled {
			continue
		}
		for _, l := range inv.Lines {
			row := domain.GoodsPurchaseRow{
				ItemName: l.ItemName, Unit: l.Unit, Quantity: l.Quantity,
				UnitPrice: l.UnitPrice, LineTotal: l.LineTotal,
				VATRate: l.VATRate, VATAmount: l.LineVATAmount, AccountID: l.AccountID,
			}
			rpt.TotalQuantity += l.Quantity
			rpt.TotalAmount += l.LineTotal
			rpt.TotalVAT += l.LineVATAmount
			rpt.Rows = append(rpt.Rows, row)
		}
	}
	rpt.TotalAmount = round2(rpt.TotalAmount)
	rpt.TotalVAT = round2(rpt.TotalVAT)
	return rpt, nil
}

// GetVATInputReport builds the VAT input tracking report (bảng kê hóa đơn VAT đầu vào).
func (s *PurchaseService) GetVATInputReport(ctx context.Context, companyID, fromDate, toDate string) (*domain.VATInputReport, error) {
	rpt := &domain.VATInputReport{CompanyID: companyID, FromDate: fromDate, ToDate: toDate}
	invs, _, err := s.invRepo.ListInvoices(ctx, domain.SupplierInvoiceFilter{
		CompanyID: companyID, FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		return nil, err
	}
	for _, inv := range invs {
		if inv.Status == domain.InvoiceCancelled {
			continue
		}
		// group lines by VAT rate
		type vatGroup struct{ subtotal, vat float64 }
		byRate := make(map[float64]vatGroup)
		for _, l := range inv.Lines {
			v := byRate[l.VATRate]
			v.subtotal += l.LineTotal
			v.vat += l.LineVATAmount
			byRate[l.VATRate] = v
		}
		for rate, v := range byRate {
			rpt.TotalSubtotal += v.subtotal
			rpt.TotalVAT += v.vat
			rpt.Rows = append(rpt.Rows, domain.VATInputRow{
				InvoiceNumber: inv.InvoiceNumber, InvoiceDate: inv.InvoiceDate.Format(time.DateOnly),
				SupplierName: inv.SupplierName, SupplierTaxCode: inv.SupplierTaxCode,
				VATRate: rate, Subtotal: round2(v.subtotal), VATAmount: round2(v.vat),
				DeductionStatus: string(inv.VATDeductionStatus),
			})
		}
	}
	rpt.TotalSubtotal = round2(rpt.TotalSubtotal)
	rpt.TotalVAT = round2(rpt.TotalVAT)
	return rpt, nil
}

// GetUninvoicedReceipts lists posted GRNs not yet matched to an invoice.
func (s *PurchaseService) GetUninvoicedReceipts(ctx context.Context, companyID string) ([]domain.UninvoicedReceiptRow, error) {
	invs, _, err := s.invRepo.ListInvoices(ctx, domain.SupplierInvoiceFilter{CompanyID: companyID})
	if err != nil {
		return nil, err
	}
	invoicedGRNs := make(map[string]bool)
	for _, inv := range invs {
		if inv.GRNID != "" {
			invoicedGRNs[inv.GRNID] = true
		}
	}
	grns, _, err := s.grnRepo.ListGRNs(ctx, domain.GRNFilter{CompanyID: companyID})
	if err != nil {
		return nil, err
	}
	var rows []domain.UninvoicedReceiptRow
	for _, g := range grns {
		if g.Status != domain.GRNPosted || invoicedGRNs[g.ID] {
			continue
		}
		supplierName := s.supplierNameForGRN(ctx, g)
		for _, l := range g.Lines {
			rows = append(rows, domain.UninvoicedReceiptRow{
				GRNNumber: g.GRNNumber, ReceiptDate: g.ReceiptDate.Format(time.DateOnly),
				SupplierName: supplierName, POID: g.POID, ItemName: l.ItemName,
				Unit: l.Unit, Quantity: l.QuantityReceived, UnitPrice: l.UnitPrice,
				LineTotal: round2(l.LineTotal),
			})
		}
	}
	return rows, nil
}

func ledgerDelta(t domain.APTransaction) float64 {
	switch t.TransactionType {
	case domain.APTransInvoice, domain.APTransCreditNote:
		return t.Amount
	default:
		return -t.Amount
	}
}

func ledgerDebit(t domain.APTransaction) float64 {
	switch t.TransactionType {
	case domain.APTransInvoice, domain.APTransCreditNote:
		return 0
	default:
		return round2(t.Amount)
	}
}

func ledgerCredit(t domain.APTransaction) float64 {
	switch t.TransactionType {
	case domain.APTransInvoice, domain.APTransCreditNote:
		return round2(t.Amount)
	default:
		return 0
	}
}

func grnTotal(g *domain.GRN) float64 {
	var sum float64
	for _, l := range g.Lines {
		sum += l.LineTotal
	}
	return sum
}

func (s *PurchaseService) grnTotalBefore(ctx context.Context, companyID, fromDate string) (float64, error) {
	grns, _, err := s.grnRepo.ListGRNs(ctx, domain.GRNFilter{CompanyID: companyID, ToDate: fromDate})
	if err != nil {
		return 0, err
	}
	var sum float64
	for _, g := range grns {
		if g.Status != domain.GRNPosted {
			continue
		}
		sum += grnTotal(&g)
	}
	return round2(sum), nil
}

func (s *PurchaseService) supplierNameForGRN(ctx context.Context, g domain.GRN) string {
	if g.POID == "" {
		return ""
	}
	po, err := s.poRepo.GetPO(ctx, g.POID)
	if err != nil || po == nil {
		return ""
	}
	sup, err := s.supRepo.GetSupplier(ctx, po.SupplierID)
	if err != nil || sup == nil {
		return ""
	}
	return sup.Name
}

func ageInMonths(from, to time.Time) int {
	if to.Before(from) {
		return 0
	}
	return (to.Year()-from.Year())*12 + int(to.Month()) - int(from.Month())
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

// ─── Requisition ─────────────────────────────────────────────────────────

func (s *PurchaseService) CreateRequisition(ctx context.Context, r *domain.PurchaseRequisition) error {
	if err := validate.Requisition(r); err != nil {
		return err
	}
	r.CalculateTotals()
	r.CreatedAt = s.now()
	r.UpdatedAt = r.CreatedAt
	if err := s.reqRepo.CreateRequisition(ctx, r); err != nil {
		return err
	}
	for i := range r.Lines {
		r.Lines[i].RequisitionID = r.ID
		r.Lines[i].LineNumber = i + 1
	}
	return s.reqRepo.CreateRequisitionLines(ctx, r.Lines)
}

func (s *PurchaseService) GetRequisition(ctx context.Context, id string) (*domain.PurchaseRequisition, error) {
	return s.reqRepo.GetRequisition(ctx, id)
}

func (s *PurchaseService) ListRequisitions(ctx context.Context, filter domain.RequisitionFilter) ([]domain.PurchaseRequisition, int, error) {
	return s.reqRepo.ListRequisitions(ctx, filter)
}

func (s *PurchaseService) UpdateRequisition(ctx context.Context, r *domain.PurchaseRequisition) error {
	existing, err := s.reqRepo.GetRequisition(ctx, r.ID)
	if err != nil {
		return err
	}
	if existing.Status != domain.ReqDraft {
		return domain.ErrRequisitionCannotUpdate
	}
	if err := validate.Requisition(r); err != nil {
		return err
	}
	r.CalculateTotals()
	r.UpdatedAt = s.now()
	return s.reqRepo.UpdateRequisition(ctx, r)
}

func (s *PurchaseService) SubmitRequisition(ctx context.Context, id, submittedBy string) error {
	r, err := s.reqRepo.GetRequisition(ctx, id)
	if err != nil {
		return err
	}
	if !r.Status.ValidTransition(domain.ReqPending) {
		return domain.ErrRequisitionInvalidTransition
	}
	return s.reqRepo.UpdateRequisitionStatus(ctx, id, domain.ReqPending, "", time.Time{})
}

func (s *PurchaseService) ApproveRequisition(ctx context.Context, id, approvedBy string) error {
	r, err := s.reqRepo.GetRequisition(ctx, id)
	if err != nil {
		return err
	}
	if !r.Status.ValidTransition(domain.ReqApproved) {
		return domain.ErrRequisitionInvalidTransition
	}
	return s.reqRepo.UpdateRequisitionStatus(ctx, id, domain.ReqApproved, approvedBy, s.now())
}

func (s *PurchaseService) RejectRequisition(ctx context.Context, id, reason string) error {
	r, err := s.reqRepo.GetRequisition(ctx, id)
	if err != nil {
		return err
	}
	if !r.Status.ValidTransition(domain.ReqRejected) {
		return domain.ErrRequisitionInvalidTransition
	}
	return s.reqRepo.RejectRequisition(ctx, id, reason)
}

func (s *PurchaseService) ConvertRequisitionToPO(ctx context.Context, requisitionID, supplierID string) (*domain.PurchaseOrder, error) {
	r, err := s.reqRepo.GetRequisition(ctx, requisitionID)
	if err != nil {
		return nil, err
	}
	if r.Status != domain.ReqApproved {
		return nil, domain.ErrRequisitionNotApproved
	}
	if _, err := s.supRepo.GetSupplier(ctx, supplierID); err != nil {
		return nil, domain.ErrPOSupplierRequired
	}
	po := &domain.PurchaseOrder{
		CompanyID:     r.CompanyID,
		PONumber:      "PO-" + r.RequisitionNumber,
		SupplierID:    supplierID,
		RequisitionID: r.ID,
		OrderDate:     s.now(),
		Currency:      "VND",
		Status:        domain.POStatusDraft,
		CreatedBy:     r.CreatedBy,
	}
	for i, l := range r.Lines {
		po.Lines = append(po.Lines, domain.POItem{
			LineNumber:   i + 1,
			ItemCode:     l.ItemCode,
			ItemName:     l.ItemName,
			Unit:         l.Unit,
			Quantity:     l.Quantity,
			UnitPrice:    l.EstimatedPrice,
			VATRate:      10,
			VATType:      domain.VAT10,
			AccountID:    l.AccountID,
			VATAccountID: "1331",
		})
	}
	if err := s.CreatePO(ctx, po); err != nil {
		return nil, err
	}
	if err := s.reqRepo.UpdateRequisitionStatus(ctx, r.ID, domain.ReqOrdered, "", time.Time{}); err != nil {
		return nil, err
	}
	return po, nil
}

func (s *PurchaseService) DeleteRequisition(ctx context.Context, id string) error {
	r, err := s.reqRepo.GetRequisition(ctx, id)
	if err != nil {
		return err
	}
	if r.Status != domain.ReqDraft {
		return domain.ErrRequisitionCannotUpdate
	}
	return s.reqRepo.DeleteRequisition(ctx, id)
}

// ─── Purchase Return / Credit Note (P2-2) ────────────────────────────────

func (s *PurchaseService) CreateReturnGRN(ctx context.Context, grn *domain.GRN) error {
	if grn.ReturnOfGRNID == "" {
		return domain.ErrReturnGRNRequired
	}
	original, err := s.grnRepo.GetGRN(ctx, grn.ReturnOfGRNID)
	if err != nil {
		return domain.ErrReturnGRNRequired
	}
	if original.Status != domain.GRNPosted {
		return domain.ErrReturnGRNNotPosted
	}
	if original.ReturnOfGRNID != "" {
		return domain.ErrReturnGRNReturnAgain
	}
	origLines, err := s.grnRepo.GetGRNLines(ctx, original.ID)
	if err != nil {
		return err
	}
	origByPOLine := make(map[string]float64, len(origLines))
	for _, l := range origLines {
		origByPOLine[l.POLineID] += l.QuantityReceived
	}
	returned, _, err := s.grnRepo.ListGRNs(ctx, domain.GRNFilter{CompanyID: grn.CompanyID, ReturnOfGRNID: original.ID})
	if err != nil {
		return err
	}
	alreadyReturned := make(map[string]float64)
	for _, rg := range returned {
		if rg.Status != domain.GRNPosted {
			continue
		}
		lines, _ := s.grnRepo.GetGRNLines(ctx, rg.ID)
		for _, l := range lines {
			alreadyReturned[l.POLineID] += l.QuantityReceived
		}
	}
	grn.POID = original.POID
	grn.WarehouseID = original.WarehouseID
	if err := validate.GRN(grn); err != nil {
		return err
	}
	for i := range grn.Lines {
		l := &grn.Lines[i]
		if l.QuantityReceived <= 0 {
			return domain.ErrReturnQtyExceeds
		}
		balance := origByPOLine[l.POLineID] - alreadyReturned[l.POLineID]
		if l.QuantityReceived > balance+0.001 {
			return domain.ErrReturnQtyExceeds
		}
	}
	now := s.now()
	grn.Status = domain.GRNPosted
	grn.PostedAt = now
	grn.CreatedAt = now
	if err := s.grnRepo.CreateGRN(ctx, grn); err != nil {
		return err
	}
	for i := range grn.Lines {
		grn.Lines[i].GRNID = grn.ID
	}
	return s.grnRepo.CreateGRNLines(ctx, grn.Lines)
}

func (s *PurchaseService) CreateCreditNote(ctx context.Context, inv *domain.SupplierInvoice) error {
	if inv.InvoiceType != domain.InvoiceTypeCreditNote {
		return domain.ErrCreditNoteTypeInvalid
	}
	if inv.OriginalInvoiceID == "" {
		return domain.ErrCreditNoteOriginalRequired
	}
	original, err := s.invRepo.GetInvoice(ctx, inv.OriginalInvoiceID)
	if err != nil {
		return domain.ErrCreditNoteOriginalRequired
	}
	if original.Status != domain.InvoicePosted && original.Status != domain.InvoicePaid {
		return domain.ErrCreditNoteOriginalNotPosted
	}
	inv.CompanyID = original.CompanyID
	inv.SupplierID = original.SupplierID
	inv.SupplierName = original.SupplierName
	inv.SupplierTaxCode = original.SupplierTaxCode
	if inv.Currency == "" {
		inv.Currency = original.Currency
	}
	if err := validate.SupplierInvoice(inv); err != nil {
		return err
	}
	for i := range inv.Lines {
		inv.Lines[i].Quantity = -math.Abs(inv.Lines[i].Quantity)
	}
	inv.CalculateTotals()
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

// ─── FX Revaluation (P2-4) ───────────────────────────────────────────────

func (s *PurchaseService) RevalueAP(ctx context.Context, companyID string, asOfDate time.Time) (*domain.FXRevaluation, error) {
	if s.gl == nil {
		return nil, domain.ErrFXRevaluationRateMissing
	}
	invoices, _, err := s.invRepo.ListInvoices(ctx, domain.SupplierInvoiceFilter{CompanyID: companyID, Status: domain.InvoicePosted})
	if err != nil {
		return nil, err
	}
	reval := &domain.FXRevaluation{CompanyID: companyID, RevaluationDate: asOfDate, CreatedAt: s.now()}
	for _, inv := range invoices {
		if inv.Currency == "" || inv.Currency == "VND" || inv.BalanceDue <= 0.001 || inv.ExchangeRate <= 0 {
			continue
		}
		rate, err := s.gl.GetExchangeRate(ctx, inv.Currency, asOfDate)
		if err != nil {
			return nil, domain.ErrFXRevaluationRateMissing
		}
		balanceVNDOld := inv.BalanceDue * inv.ExchangeRate
		balanceVNDNow := inv.BalanceDue * rate.AverageRate
		diff := balanceVNDNow - balanceVNDOld
		line := domain.FXRevaluationLine{
			InvoiceID:       inv.ID,
			InvoiceNumber:   inv.InvoiceNumber,
			SupplierID:      inv.SupplierID,
			SupplierName:    inv.SupplierName,
			Currency:        inv.Currency,
			BalanceDue:      inv.BalanceDue,
			OriginalRate:    inv.ExchangeRate,
			RevaluationRate: rate.AverageRate,
		}
		if diff > 0 {
			line.FxLoss = diff
			reval.TotalLoss += diff
		} else if diff < 0 {
			line.FxGain = -diff
			reval.TotalGain += -diff
		} else {
			continue
		}
		reval.Lines = append(reval.Lines, line)
	}
	if len(reval.Lines) == 0 {
		return nil, domain.ErrFXRevaluationEmpty
	}
	if err := validate.FXRevaluation(reval); err != nil {
		return nil, err
	}
	if err := s.fxRepo.CreateRevaluation(ctx, reval); err != nil {
		return nil, err
	}
	for i := range reval.Lines {
		reval.Lines[i].RevaluationID = reval.ID
	}
	if err := s.fxRepo.CreateRevaluationLines(ctx, reval.Lines); err != nil {
		return nil, err
	}
	return reval, nil
}

func (s *PurchaseService) GetFXRevaluation(ctx context.Context, id string) (*domain.FXRevaluation, error) {
	return s.fxRepo.GetRevaluation(ctx, id)
}

func (s *PurchaseService) ListFXRevaluations(ctx context.Context, companyID string, offset, limit int) ([]domain.FXRevaluation, int, error) {
	return s.fxRepo.ListRevaluations(ctx, companyID, limit, offset)
}

func (s *PurchaseService) PostFXRevaluation(ctx context.Context, id string) error {
	reval, err := s.fxRepo.GetRevaluation(ctx, id)
	if err != nil {
		return err
	}
	if reval.Status != domain.FXRevalDraft {
		return domain.ErrFXRevaluationAlreadyPosted
	}
	lines := make([]domain.JournalLine, 0, len(reval.Lines)*2)
	for _, l := range reval.Lines {
		if l.FxGain > 0 {
			lines = append(lines,
				domain.JournalLine{AccountCode: apPayableAccount, DebitAmount: l.FxGain, Description: "FX gain reval: " + l.InvoiceNumber},
				domain.JournalLine{AccountCode: fxGainAccount, CreditAmount: l.FxGain, Description: "FX gain reval: " + l.InvoiceNumber},
			)
		}
		if l.FxLoss > 0 {
			lines = append(lines,
				domain.JournalLine{AccountCode: fxLossAccount, DebitAmount: l.FxLoss, Description: "FX loss reval: " + l.InvoiceNumber},
				domain.JournalLine{AccountCode: apPayableAccount, CreditAmount: l.FxLoss, Description: "FX loss reval: " + l.InvoiceNumber},
			)
		}
	}
	now := s.now().UTC()
	entry := &domain.JournalEntry{
		CompanyID:   reval.CompanyID,
		EntryNumber: "FXR-" + reval.ID,
		VoucherType: domain.VoucherTypePurchase,
		EntryDate:   reval.RevaluationDate,
		Description: "AP FX revaluation " + reval.RevaluationDate.Format("2006-01-02"),
		Lines:       lines,
	}
	if s.gl != nil {
		if err := s.gl.CreatePostedEntry(ctx, entry, reval.CreatedBy); err != nil {
			return err
		}
		if err := s.fxRepo.SetRevaluationGLPosted(ctx, id, now); err != nil {
			return err
		}
	}
	return s.fxRepo.UpdateRevaluationStatus(ctx, id, domain.FXRevalPosted)
}

// ─── E-Invoice GDT XML (P2-5) ────────────────────────────────────────────

// ReceiveEInvoiceXML parses a GDT e-invoice XML document into a draft
// supplier invoice and persists it (supplier auto-created on first receipt).
// The raw XML is stored on the invoice for audit.
func (s *PurchaseService) ReceiveEInvoiceXML(ctx context.Context, companyID string, raw []byte) (*domain.SupplierInvoice, error) {
	inv, err := einvoice.Parse(raw)
	if err != nil {
		return nil, err
	}
	if inv.InvoiceType == domain.InvoiceTypeCreditNote {
		return nil, domain.ErrEinvoiceCreditNoteUnsupported
	}
	inv.CompanyID = companyID
	inv.EInvoiceData = string(raw)
	if err := s.ReceiveEInvoice(ctx, inv); err != nil {
		return nil, err
	}
	return inv, nil
}

// GenerateEInvoiceXML renders a stored supplier invoice as a GDT e-invoice
// XML document.
func (s *PurchaseService) GenerateEInvoiceXML(ctx context.Context, id string) (string, error) {
	inv, err := s.invRepo.GetInvoice(ctx, id)
	if err != nil {
		return "", err
	}
	raw, err := einvoice.Generate(inv)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
