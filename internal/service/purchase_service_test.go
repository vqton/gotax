package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
	"gotax/internal/repository"
)

func newPurchaseTestSvc() (*PurchaseService, context.Context) {
	repo := repository.NewMemoryPurchaseRepo()
	return NewPurchaseService(repo, repo, repo, repo, repo, repo, nil), context.Background()
}

func setupPurchaseGL(t *testing.T) (*PurchaseService, Service, context.Context) {
	t.Helper()
	accRepo := repository.NewMemoryAccountRepo()
	jeRepo := repository.NewMemoryJournalRepo()
	jeRepo.SetAccounts(accRepo.Accounts())
	perRepo := repository.NewMemoryPeriodRepo()
	userRepo := repository.NewMemoryUserRepo()
	auditRepo := repository.NewMemoryAuditLogRepo()
	rateRepo := repository.NewMemoryExchangeRateRepo()
	templateRepo := repository.NewMemoryClosingTemplateRepo()
	approvalRepo := repository.NewMemoryApprovalRepo()
	versionRepo := repository.NewMemoryAccountVersionRepo()
	mappingRepo := repository.NewMemoryAccountMappingRepo()
	analysisRepo := repository.NewMemoryAccountAnalysisRepo()
	ifrsRepo := repository.NewMemoryIFRSMappingRepo()
	refreshRepo := repository.NewMemoryRefreshTokenRepo()
	resetRepo := repository.NewMemoryPasswordResetTokenRepo()
	obRepo := repository.NewMemoryOpeningBalanceRepo()
	cashRepo := repository.NewMemoryCashRepo()
	gl := NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo,
		approvalRepo, versionRepo, mappingRepo, analysisRepo, ifrsRepo, refreshRepo, resetRepo, obRepo, cashRepo)

	ctx := context.Background()
	for _, acc := range []domain.Account{
		{Code: "331", Name: "Phai tra nguoi ban", Type: domain.AccountTypeLiability, IsActive: true},
		{Code: "1331", Name: "Thue GTGT duoc khau tru", Type: domain.AccountTypeAsset, IsActive: true},
		{Code: "152", Name: "Nguyen vat lieu", Type: domain.AccountTypeAsset, IsActive: true},
		{Code: "642", Name: "Chi phi quan ly doanh nghiep", Type: domain.AccountTypeExpense, IsActive: true},
	} {
		require.NoError(t, gl.CreateAccount(ctx, &acc))
	}

	repo := repository.NewMemoryPurchaseRepo()
	svc := NewPurchaseService(repo, repo, repo, repo, repo, repo, gl)
	return svc, gl, ctx
}

func makeSupplier(t *testing.T, svc *PurchaseService, ctx context.Context, code string) *domain.Supplier {
	t.Helper()
	sup := &domain.Supplier{CompanyID: "c1", Code: code, Name: "Supplier " + code, TaxCode: code + "-TX"}
	require.NoError(t, svc.CreateSupplier(ctx, sup))
	return sup
}

func makePO(t *testing.T, svc *PurchaseService, ctx context.Context, supID string, qty, price float64) *domain.PurchaseOrder {
	t.Helper()
	po := &domain.PurchaseOrder{
		CompanyID: "c1", PONumber: "PO-" + supID + "-1", SupplierID: supID,
		OrderDate: time.Now(), Currency: "VND",
		Lines: []domain.POItem{{
			ItemName: "Widget", Unit: "pcs", Quantity: qty, UnitPrice: price,
			VATRate: 10, VATType: domain.VAT10,
			AccountID: "152", VATAccountID: "1331",
		}},
	}
	require.NoError(t, svc.CreatePO(ctx, po))
	loaded, err := svc.GetPO(ctx, po.ID)
	require.NoError(t, err)
	return loaded
}

func makeGRN(t *testing.T, svc *PurchaseService, ctx context.Context, po *domain.PurchaseOrder, qty float64) *domain.GRN {
	t.Helper()
	grn := &domain.GRN{
		CompanyID: "c1", GRNNumber: "GRN-" + po.PONumber, POID: po.ID,
		ReceiptDate: time.Now(),
		Lines: []domain.GRNItem{{
			POLineID: po.Lines[0].ID, ItemName: "Widget", Unit: "pcs",
			QuantityReceived: qty, UnitPrice: po.Lines[0].UnitPrice,
			LineTotal: qty * po.Lines[0].UnitPrice,
		}},
	}
	require.NoError(t, svc.CreateGRN(ctx, grn))
	loaded, err := svc.GetGRN(ctx, grn.ID)
	require.NoError(t, err)
	return loaded
}

func makeInvoice(t *testing.T, svc *PurchaseService, ctx context.Context, supID string, poID, grnID, poLineID string, qty, price float64) *domain.SupplierInvoice {
	t.Helper()
	inv := &domain.SupplierInvoice{
		CompanyID: "c1", InvoiceNumber: "INV-" + supID, SupplierID: supID,
		POID: poID, GRNID: grnID, InvoiceDate: time.Now(),
		Lines: []domain.SupplierInvoiceLine{{
			POLineID: poLineID, ItemName: "Widget", Unit: "pcs",
			Quantity: qty, UnitPrice: price, VATRate: 10, VATType: domain.VAT10,
			AccountID: "152", VATAccountID: "1331",
		}},
	}
	require.NoError(t, svc.CreateInvoice(ctx, inv))
	return inv
}

func verifyAndPost(t *testing.T, svc *PurchaseService, ctx context.Context, inv *domain.SupplierInvoice) error {
	t.Helper()
	if err := svc.VerifyInvoice(ctx, inv.ID); err != nil {
		return err
	}
	return svc.PostInvoice(ctx, inv.ID)
}

// ─── Supplier ────────────────────────────────────────────────────────────

func TestPurchaseSupplier_Defaults(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "S001")
	require.Equal(t, "VND", sup.Currency)
	require.Equal(t, domain.SupplierActive, sup.Status)
}

func TestPurchaseSupplier_DuplicateCode(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	makeSupplier(t, svc, ctx, "S001")
	err := svc.CreateSupplier(ctx, &domain.Supplier{
		CompanyID: "c1", Code: "S001", Name: "Dup", TaxCode: "TX",
	})
	assert.ErrorIs(t, err, domain.ErrSupplierCodeExists)
}

// ─── Purchase Order ──────────────────────────────────────────────────────

func TestCreatePO_ValidatesLinesRequired(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "S002")
	po := &domain.PurchaseOrder{
		CompanyID: "c1", PONumber: "PO-2", SupplierID: sup.ID, OrderDate: time.Now(),
	}
	err := svc.CreatePO(ctx, po)
	assert.ErrorIs(t, err, domain.ErrPOLinesRequired)
}

func TestCreatePO_CalculatesTotals(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "S003")
	po := &domain.PurchaseOrder{
		CompanyID: "c1", PONumber: "PO-3", SupplierID: sup.ID, OrderDate: time.Now(),
		Lines: []domain.POItem{{
			ItemName: "Widget", Unit: "pcs", Quantity: 10, UnitPrice: 100000,
			VATRate: 10, VATType: domain.VAT10, AccountID: "152", VATAccountID: "1331",
		}},
	}
	require.NoError(t, svc.CreatePO(ctx, po))
	assert.InDelta(t, 1000000, po.Subtotal, 0.001)
	assert.InDelta(t, 100000, po.TaxAmount, 0.001)
	assert.InDelta(t, 1100000, po.TotalAmount, 0.001)
}

func TestPO_ApproveTwiceFails(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "S004")
	po := &domain.PurchaseOrder{
		CompanyID: "c1", PONumber: "PO-4", SupplierID: sup.ID, OrderDate: time.Now(),
		Lines: []domain.POItem{{
			ItemName: "Widget", Unit: "pcs", Quantity: 1, UnitPrice: 100,
			AccountID: "152", VATAccountID: "1331",
		}},
	}
	require.NoError(t, svc.CreatePO(ctx, po))
	require.NoError(t, svc.ApprovePO(ctx, po.ID, "user"))
	err := svc.ApprovePO(ctx, po.ID, "user")
	assert.ErrorIs(t, err, domain.ErrPOInvalidTransition)
	loaded, _ := svc.GetPO(ctx, po.ID)
	assert.Equal(t, domain.POStatusApproved, loaded.Status)
}

// ─── GRN ─────────────────────────────────────────────────────────────────

func TestCreateGRN_RequiresPO(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	grn := &domain.GRN{
		CompanyID: "c1", GRNNumber: "GRN-1", POID: "missing", ReceiptDate: time.Now(),
		Lines: []domain.GRNItem{{ItemName: "W", Unit: "pcs", QuantityReceived: 1}},
	}
	err := svc.CreateGRN(ctx, grn)
	assert.ErrorIs(t, err, domain.ErrGRNPORequired)
}

func TestGRN_PostThenRepostFails(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "S005")
	po := makePO(t, svc, ctx, sup.ID, 10, 100)
	grn := makeGRN(t, svc, ctx, po, 10)
	require.NoError(t, svc.PostGRN(ctx, grn.ID))
	err := svc.PostGRN(ctx, grn.ID)
	assert.ErrorIs(t, err, domain.ErrGRNInvalidTransition)
	err = svc.CancelGRN(ctx, grn.ID)
	assert.ErrorIs(t, err, domain.ErrGRNInvalidTransition)
}

// ─── Invoice totals & state machine ──────────────────────────────────────

func TestCreateInvoice_CalculatesTotals(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "S006")
	inv := makeInvoice(t, svc, ctx, sup.ID, "", "", "", 1, 1000000)
	assert.InDelta(t, 1000000, inv.Subtotal, 0.001)
	assert.InDelta(t, 100000, inv.TaxAmount, 0.001)
	assert.InDelta(t, 1100000, inv.TotalAmount, 0.001)
}

func TestInvoice_StateTransitions(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "S007")
	inv := makeInvoice(t, svc, ctx, sup.ID, "", "", "", 1, 1000)

	err := svc.PostInvoice(ctx, inv.ID)
	assert.ErrorIs(t, err, domain.ErrInvoiceInvalidTransition)

	require.NoError(t, svc.VerifyInvoice(ctx, inv.ID))
	err = svc.VerifyInvoice(ctx, inv.ID)
	assert.ErrorIs(t, err, domain.ErrInvoiceInvalidTransition)

	require.NoError(t, svc.PostInvoice(ctx, inv.ID))
	err = svc.PostInvoice(ctx, inv.ID)
	assert.ErrorIs(t, err, domain.ErrInvoiceInvalidTransition)

	loaded, _ := svc.GetInvoice(ctx, inv.ID)
	assert.Equal(t, domain.InvoicePosted, loaded.Status)
}

func TestCancelInvoice_PostedBlocked(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "S008")
	inv := makeInvoice(t, svc, ctx, sup.ID, "", "", "", 1, 1000)
	require.NoError(t, svc.VerifyInvoice(ctx, inv.ID))
	require.NoError(t, svc.PostInvoice(ctx, inv.ID))
	err := svc.CancelInvoice(ctx, inv.ID)
	assert.ErrorIs(t, err, domain.ErrInvoiceInvalidTransition)
}

// ─── 3-way matching ──────────────────────────────────────────────────────

func TestThreeWayMatch_OK(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "S010")
	po := makePO(t, svc, ctx, sup.ID, 10, 100000)
	grn := makeGRN(t, svc, ctx, po, 10)
	inv := makeInvoice(t, svc, ctx, sup.ID, po.ID, grn.ID, po.Lines[0].ID, 10, 100000)
	require.NoError(t, verifyAndPost(t, svc, ctx, inv))
}

func TestThreeWayMatch_QtyExceedsPO(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "S011")
	po := makePO(t, svc, ctx, sup.ID, 10, 100000)
	grn := makeGRN(t, svc, ctx, po, 10)
	inv := makeInvoice(t, svc, ctx, sup.ID, po.ID, grn.ID, po.Lines[0].ID, 11, 100000)
	err := verifyAndPost(t, svc, ctx, inv)
	assert.ErrorIs(t, err, domain.ErrInvoice3WayMismatch)
}

func TestThreeWayMatch_QtyExceedsGRN(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "S012")
	po := makePO(t, svc, ctx, sup.ID, 10, 100000)
	grn := makeGRN(t, svc, ctx, po, 5)
	inv := makeInvoice(t, svc, ctx, sup.ID, po.ID, grn.ID, po.Lines[0].ID, 6, 100000)
	err := verifyAndPost(t, svc, ctx, inv)
	assert.ErrorIs(t, err, domain.ErrInvoice3WayMismatch)
}

func TestThreeWayMatch_PriceVarianceExceeds(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "S013")
	po := makePO(t, svc, ctx, sup.ID, 10, 100000)
	grn := makeGRN(t, svc, ctx, po, 10)
	inv := makeInvoice(t, svc, ctx, sup.ID, po.ID, grn.ID, po.Lines[0].ID, 10, 110000)
	err := verifyAndPost(t, svc, ctx, inv)
	assert.ErrorIs(t, err, domain.ErrInvoice3WayMismatch)
}

func TestThreeWayMatch_PriceWithinTolerance(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "S014")
	po := makePO(t, svc, ctx, sup.ID, 10, 100000)
	grn := makeGRN(t, svc, ctx, po, 10)
	inv := makeInvoice(t, svc, ctx, sup.ID, po.ID, grn.ID, po.Lines[0].ID, 10, 104000)
	require.NoError(t, verifyAndPost(t, svc, ctx, inv))
}

func TestThreeWayMatch_UnknownPOLine(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "S015")
	po := makePO(t, svc, ctx, sup.ID, 10, 100000)
	grn := makeGRN(t, svc, ctx, po, 10)
	inv := makeInvoice(t, svc, ctx, sup.ID, po.ID, grn.ID, "no-such-line", 10, 100000)
	err := verifyAndPost(t, svc, ctx, inv)
	assert.ErrorIs(t, err, domain.ErrInvoice3WayMismatch)
}

func TestThreeWayMatch_NoPO_NoMatchRequired(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "S016")
	inv := makeInvoice(t, svc, ctx, sup.ID, "", "", "", 1, 100000)
	require.NoError(t, verifyAndPost(t, svc, ctx, inv))
}

func TestThreeWayMatch_MismatchDoesNotChangeStatus(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "S017")
	po := makePO(t, svc, ctx, sup.ID, 10, 100000)
	grn := makeGRN(t, svc, ctx, po, 10)
	inv := makeInvoice(t, svc, ctx, sup.ID, po.ID, grn.ID, po.Lines[0].ID, 11, 100000)
	require.NoError(t, svc.VerifyInvoice(ctx, inv.ID))
	err := svc.PostInvoice(ctx, inv.ID)
	assert.ErrorIs(t, err, domain.ErrInvoice3WayMismatch)
	loaded, _ := svc.GetInvoice(ctx, inv.ID)
	assert.Equal(t, domain.InvoiceVerified, loaded.Status)
}

// ─── GL auto-posting ─────────────────────────────────────────────────────

func TestPurchasePostInvoice_CreatesGLEntry(t *testing.T) {
	svc, gl, ctx := setupPurchaseGL(t)
	sup := makeSupplier(t, svc, ctx, "S020")
	inv := makeInvoice(t, svc, ctx, sup.ID, "", "", "", 2, 100000)
	require.NoError(t, verifyAndPost(t, svc, ctx, inv))

	loaded, _ := svc.GetInvoice(ctx, inv.ID)
	assert.Equal(t, domain.InvoicePosted, loaded.Status)
	assert.True(t, loaded.GLPosted)
	require.NotNil(t, loaded.GLPostedAt)

	entries, err := gl.GetEntriesByDateRange(ctx, inv.InvoiceDate, inv.InvoiceDate.Add(24*time.Hour))
	require.NoError(t, err)
	var apEntry *domain.JournalEntry
	for i := range entries {
		if entries[i].Description == "AP invoice "+inv.InvoiceNumber {
			apEntry = &entries[i]
			break
		}
	}
	require.NotNil(t, apEntry, "AP journal entry not found")
	assert.Equal(t, domain.JournalEntryPosted, apEntry.Status)

	var totalDebit, totalCredit float64
	var hasInv, hasVAT, hasAP bool
	for _, l := range apEntry.Lines {
		totalDebit += l.DebitAmount
		totalCredit += l.CreditAmount
		if l.AccountCode == "152" && l.DebitAmount > 0 {
			hasInv = true
		}
		if l.AccountCode == "1331" && l.DebitAmount > 0 {
			hasVAT = true
		}
		if l.AccountCode == "331" && l.CreditAmount > 0 {
			hasAP = true
		}
	}
	assert.True(t, hasInv, "must Dr 152 (inventory)")
	assert.True(t, hasVAT, "must Dr 1331 (VAT input)")
	assert.True(t, hasAP, "must Cr 331 (AP)")
	assert.InDelta(t, totalDebit, totalCredit, 0.001)
	assert.InDelta(t, 220000, totalDebit, 0.001)
}

func TestPostInvoice_NoGL_NoJournalEntry(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "S021")
	inv := makeInvoice(t, svc, ctx, sup.ID, "", "", "", 1, 100000)
	require.NoError(t, verifyAndPost(t, svc, ctx, inv))

	loaded, _ := svc.GetInvoice(ctx, inv.ID)
	assert.Equal(t, domain.InvoicePosted, loaded.Status)
	assert.False(t, loaded.GLPosted, "GLPosted must stay false without GL service")
}

func TestPostInvoice_GLFailureKeepsInvoiceUnposted(t *testing.T) {
	svc, _, ctx := setupPurchaseGL(t)
	sup := makeSupplier(t, svc, ctx, "S022")
	inv := makeInvoice(t, svc, ctx, sup.ID, "", "", "", 1, 100000)
	inv.Lines[0].AccountID = "811"
	require.NoError(t, svc.CreateInvoice(ctx, inv))
	require.NoError(t, svc.VerifyInvoice(ctx, inv.ID))

	err := svc.PostInvoice(ctx, inv.ID)
	require.Error(t, err)

	loaded, _ := svc.GetInvoice(ctx, inv.ID)
	assert.Equal(t, domain.InvoiceVerified, loaded.Status)
	assert.False(t, loaded.GLPosted)
}

// ─── AP reports ──────────────────────────────────────────────────────────

func TestGetAPAgingReport(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "S030")
	require.NoError(t, svc.CreateAPTransaction(ctx, &domain.APTransaction{
		CompanyID: "c1", SupplierID: sup.ID, TransactionType: domain.APTransInvoice,
		TransactionDate: time.Now(), Amount: 5000000, Currency: "VND",
	}))
	report, err := svc.GetAPAgingReport(ctx, "c1")
	require.NoError(t, err)
	require.Len(t, report, 1)
	assert.Equal(t, sup.ID, report[0].SupplierID)
	assert.InDelta(t, 5000000, report[0].Buckets.Total, 0.001)
	assert.InDelta(t, 5000000, report[0].Buckets.Bucket0, 0.001)
}

func TestGetAPSummary(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "S031")
	now := time.Now()
	require.NoError(t, svc.CreateAPTransaction(ctx, &domain.APTransaction{
		CompanyID: "c1", SupplierID: sup.ID, TransactionType: domain.APTransInvoice,
		TransactionDate: now, Amount: 1000000, Currency: "VND",
	}))
	require.NoError(t, svc.CreateAPTransaction(ctx, &domain.APTransaction{
		CompanyID: "c1", SupplierID: sup.ID, TransactionType: domain.APTransPayment,
		TransactionDate: now, Amount: 400000, Currency: "VND",
	}))
	summary, err := svc.GetAPSummary(ctx, "c1")
	require.NoError(t, err)
	require.Len(t, summary, 1)
	assert.InDelta(t, 1000000, summary[0].TotalInvoiced, 0.001)
	assert.InDelta(t, 400000, summary[0].TotalPaid, 0.001)
	assert.InDelta(t, 600000, summary[0].Outstanding, 0.001)
}

// ─── Cost allocation ─────────────────────────────────────────────────────

func TestCostAllocation_InvalidMethod(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	c := &domain.CostAllocation{
		CompanyID: "c1", InvoiceID: "inv-1", CostType: domain.CostTransport,
		CostAmount: 1000, AllocationMethod: "by_weight_kg",
	}
	err := svc.CreateCostAllocation(ctx, c)
	assert.ErrorIs(t, err, domain.ErrCostAllocMethodInvalid)
}

func TestCostAllocation_CreateAndList(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	c := &domain.CostAllocation{
		CompanyID: "c1", InvoiceID: "inv-1", CostType: domain.CostTransport,
		CostAmount: 1000, AllocationMethod: domain.CostAllocByValue,
		AllocatedLines: "l1,l2",
	}
	require.NoError(t, svc.CreateCostAllocation(ctx, c))
	allocs, err := svc.ListCostAllocationsByInvoice(ctx, "inv-1")
	require.NoError(t, err)
	require.Len(t, allocs, 1)
	assert.Equal(t, c.ID, allocs[0].ID)
}
