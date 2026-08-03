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
	return NewPurchaseService(repo, repo, repo, repo, repo, repo, repo, repo, nil), context.Background()
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
	svc := NewPurchaseService(repo, repo, repo, repo, repo, repo, repo, repo, gl)
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

func TestGetPOByNumber_NotFound(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "S004B")
	po := makePO(t, svc, ctx, sup.ID, 1, 100)
	_, err := svc.GetPOByNumber(ctx, "c1", "NO-SUCH-PO")
	assert.ErrorIs(t, err, domain.ErrPONotFound)
	got, err := svc.GetPOByNumber(ctx, "c1", po.PONumber)
	require.NoError(t, err)
	assert.Equal(t, po.ID, got.ID)
}

func TestUpdatePO_AfterApprovedFails(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "S004C")
	po := makePO(t, svc, ctx, sup.ID, 1, 100)
	require.NoError(t, svc.ApprovePO(ctx, po.ID, "user"))
	err := svc.UpdatePO(ctx, &domain.PurchaseOrder{
		ID: po.ID, CompanyID: "c1", PONumber: po.PONumber, SupplierID: sup.ID,
		OrderDate: time.Now(), Currency: "VND",
		Lines: []domain.POItem{{ItemName: "Widget", Unit: "pcs", Quantity: 2, UnitPrice: 200, AccountID: "152", VATAccountID: "1331"}},
	})
	assert.ErrorIs(t, err, domain.ErrPOCannotUpdate)
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

// ─── Doubtful Debt Provisions ───────────────────────────────────────────

func makePrepayment(t *testing.T, svc *PurchaseService, ctx context.Context, supID string, date time.Time, amount float64) {
	t.Helper()
	txn := &domain.APTransaction{
		CompanyID: "c1", SupplierID: supID, TransactionType: domain.APTransPrepayment,
		TransactionDate: date, Amount: amount, Currency: "VND",
	}
	require.NoError(t, svc.CreateAPTransaction(ctx, txn))
}

func TestProvision_CalculateTiers(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	asOf := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	sup5 := makeSupplier(t, svc, ctx, "S5")   // 5 months -> 0%
	sup7 := makeSupplier(t, svc, ctx, "S7")   // 7 months -> 30%
	sup13 := makeSupplier(t, svc, ctx, "S13") // 13 months -> 50%
	sup25 := makeSupplier(t, svc, ctx, "S25") // 25 months -> 70%
	sup40 := makeSupplier(t, svc, ctx, "S40") // 40 months -> 100%

	makePrepayment(t, svc, ctx, sup5.ID, asOf.AddDate(0, -5, 0), 1000)
	makePrepayment(t, svc, ctx, sup7.ID, asOf.AddDate(0, -7, 0), 1000)
	makePrepayment(t, svc, ctx, sup13.ID, asOf.AddDate(0, -13, 0), 2000)
	makePrepayment(t, svc, ctx, sup25.ID, asOf.AddDate(0, -25, 0), 1000)
	makePrepayment(t, svc, ctx, sup40.ID, asOf.AddDate(0, -40, 0), 1000)

	lines, err := svc.CalculateDoubtfulDebtProvision(ctx, "c1", asOf)
	require.NoError(t, err)
	require.Len(t, lines, 4)
	bySup := map[string]domain.DoubtfulDebtProvisionLine{}
	for _, l := range lines {
		bySup[l.SupplierID] = l
	}
	_, ok := bySup[sup5.ID]
	assert.False(t, ok, "5-month prepayment must be excluded (0% rate)")
	assert.Equal(t, 7, bySup[sup7.ID].AgeMonths)
	assert.Equal(t, 0.30, bySup[sup7.ID].RatePct)
	assert.Equal(t, 300.0, bySup[sup7.ID].ProvisionAmount)
	assert.Equal(t, 0.50, bySup[sup13.ID].RatePct)
	assert.Equal(t, 1000.0, bySup[sup13.ID].ProvisionAmount)
	assert.Equal(t, 0.70, bySup[sup25.ID].RatePct)
	assert.Equal(t, 700.0, bySup[sup25.ID].ProvisionAmount)
	assert.Equal(t, 1.0, bySup[sup40.ID].RatePct)
	assert.Equal(t, 1000.0, bySup[sup40.ID].ProvisionAmount)
	assert.Equal(t, "Supplier S7", bySup[sup7.ID].SupplierName)
}

func TestProvision_CalculateNoPrepayments(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "SX")
	makePrepayment(t, svc, ctx, sup.ID, time.Now().AddDate(0, -2, 0), 500)
	lines, err := svc.CalculateDoubtfulDebtProvision(ctx, "c1", time.Now())
	require.NoError(t, err)
	assert.Empty(t, lines)

	_, err = svc.CalculateDoubtfulDebtProvision(ctx, "other-company", time.Now())
	assert.ErrorIs(t, err, domain.ErrProvisionNoPrepayments)
}

func TestProvision_CalculateAggregatesPerSupplier(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	asOf := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	sup := makeSupplier(t, svc, ctx, "SAGG")
	makePrepayment(t, svc, ctx, sup.ID, asOf.AddDate(0, -20, 0), 600)
	makePrepayment(t, svc, ctx, sup.ID, asOf.AddDate(0, -19, 0), 400)
	lines, err := svc.CalculateDoubtfulDebtProvision(ctx, "c1", asOf)
	require.NoError(t, err)
	require.Len(t, lines, 1)
	assert.Equal(t, 1000.0, lines[0].OutstandingAmount)
	assert.Equal(t, 20, lines[0].AgeMonths)
	assert.Equal(t, 0.50, lines[0].RatePct)
	assert.Equal(t, 500.0, lines[0].ProvisionAmount)
}

func TestProvision_CreateAndGet(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "SC")
	prov := &domain.DoubtfulDebtProvision{
		CompanyID: "c1", AsOfDate: "2026-08-01", CreatedBy: "u1",
		Lines: []domain.DoubtfulDebtProvisionLine{{
			SupplierID: sup.ID, SupplierName: sup.Name, OutstandingAmount: 1000,
			AgeMonths: 18, RatePct: 0.50,
		}},
	}
	require.NoError(t, svc.CreateDoubtfulDebtProvision(ctx, prov))
	assert.NotEmpty(t, prov.ID)
	assert.Equal(t, 500.0, prov.TotalProvision)
	assert.Equal(t, 1000.0, prov.TotalOutstanding)

	loaded, err := svc.GetDoubtfulDebtProvision(ctx, prov.ID)
	require.NoError(t, err)
	require.Len(t, loaded.Lines, 1)
	assert.Equal(t, prov.ID, loaded.Lines[0].ProvisionID)
	assert.Equal(t, 500.0, loaded.Lines[0].ProvisionAmount)

	provisions, total, err := svc.ListDoubtfulDebtProvisions(ctx, "c1", 0, 10)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, provisions, 1)
}

func TestProvision_CreateValidation(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	err := svc.CreateDoubtfulDebtProvision(ctx, &domain.DoubtfulDebtProvision{
		CompanyID: "c1", AsOfDate: "not-a-date", Lines: []domain.DoubtfulDebtProvisionLine{{SupplierID: "x"}},
	})
	assert.ErrorIs(t, err, domain.ErrProvisionDateRequired)

	err = svc.CreateDoubtfulDebtProvision(ctx, &domain.DoubtfulDebtProvision{
		CompanyID: "c1", AsOfDate: "2026-08-01",
	})
	assert.ErrorIs(t, err, domain.ErrProvisionNoLines)
}

func TestDoubtfulDebtRate(t *testing.T) {
	cases := []struct {
		months int
		want   float64
	}{
		{0, 0}, {5, 0}, {6, 0.30}, {11, 0.30}, {12, 0.50},
		{23, 0.50}, {24, 0.70}, {35, 0.70}, {36, 1.0}, {120, 1.0},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, domain.DoubtfulDebtRate(c.months), "months=%d", c.months)
	}
}

// ─── Regulatory Reports ──────────────────────────────────────────────────

func TestReport_PurchaseLedger(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "SRL")
	po := makePO(t, svc, ctx, sup.ID, 10, 1000)
	grn := makeGRN(t, svc, ctx, po, 10)
	require.NoError(t, svc.PostGRN(ctx, grn.ID))

	from := time.Now().AddDate(0, -1, 0).Format(time.DateOnly)
	to := time.Now().AddDate(0, 1, 0).Format(time.DateOnly)
	rpt, err := svc.GetPurchaseLedger(ctx, "c1", from, to)
	require.NoError(t, err)
	require.Len(t, rpt.Rows, 1)
	assert.Equal(t, 10000.0, rpt.Increase)
	assert.Equal(t, sup.Name, rpt.Rows[0].SupplierName)
	assert.Equal(t, rpt.Opening, rpt.Closing-rpt.Increase+rpt.Decrease)
}

func TestReport_SupplierLedger(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "SSL")
	// invoice -> credit, payment -> debit
	inv := &domain.SupplierInvoice{
		CompanyID: "c1", InvoiceNumber: "INV-SL", SupplierID: sup.ID,
		SupplierName: sup.Name, SupplierTaxCode: sup.TaxCode, InvoiceDate: time.Now(),
		Currency: "VND", Lines: []domain.SupplierInvoiceLine{{
			ItemName: "Widget", Unit: "pcs", Quantity: 1, UnitPrice: 5000,
			VATRate: 10, VATType: domain.VAT10, AccountID: "152", VATAccountID: "1331",
		}},
	}
	require.NoError(t, svc.CreateInvoice(ctx, inv))
	require.NoError(t, svc.CreateAPTransaction(ctx, &domain.APTransaction{
		CompanyID: "c1", SupplierID: sup.ID, InvoiceID: inv.ID,
		TransactionType: domain.APTransInvoice, TransactionDate: time.Now(), Amount: 5500, Currency: "VND",
	}))
	require.NoError(t, svc.CreateAPTransaction(ctx, &domain.APTransaction{
		CompanyID: "c1", SupplierID: sup.ID,
		TransactionType: domain.APTransPayment, TransactionDate: time.Now(), Amount: 2500, Currency: "VND",
	}))

	from := time.Now().AddDate(0, -1, 0).Format(time.DateOnly)
	to := time.Now().AddDate(0, 1, 0).Format(time.DateOnly)
	rpt, err := svc.GetSupplierLedger(ctx, "c1", sup.ID, from, to)
	require.NoError(t, err)
	require.Len(t, rpt.Rows, 2)
	var invRow, payRow domain.SupplierLedgerRow
	for _, r := range rpt.Rows {
		switch r.Description {
		case string(domain.APTransInvoice):
			invRow = r
		case string(domain.APTransPayment):
			payRow = r
		}
	}
	assert.Equal(t, 5500.0, invRow.Credit)
	assert.Equal(t, 2500.0, payRow.Debit)
	assert.Equal(t, 3000.0, rpt.Closing)
}

func TestReport_GoodsPurchase(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "SGP")
	inv := &domain.SupplierInvoice{
		CompanyID: "c1", InvoiceNumber: "INV-GP", SupplierID: sup.ID,
		SupplierName: sup.Name, SupplierTaxCode: sup.TaxCode, InvoiceDate: time.Now(),
		Currency: "VND", Lines: []domain.SupplierInvoiceLine{{
			ItemName: "Widget A", Unit: "pcs", Quantity: 2, UnitPrice: 1000,
			VATRate: 10, VATType: domain.VAT10, AccountID: "152", VATAccountID: "1331",
		}, {
			ItemName: "Widget B", Unit: "box", Quantity: 3, UnitPrice: 500,
			VATRate: 8, VATType: domain.VAT8, AccountID: "642", VATAccountID: "1331",
		}},
	}
	require.NoError(t, svc.CreateInvoice(ctx, inv))

	from := time.Now().AddDate(0, -1, 0).Format(time.DateOnly)
	to := time.Now().AddDate(0, 1, 0).Format(time.DateOnly)
	rpt, err := svc.GetGoodsPurchaseReport(ctx, "c1", from, to)
	require.NoError(t, err)
	require.Len(t, rpt.Rows, 2)
	assert.Equal(t, 3500.0, rpt.TotalAmount)
	assert.Equal(t, 200.0+120.0, rpt.TotalVAT)
	assert.Equal(t, "Widget A", rpt.Rows[0].ItemName)
}

func TestReport_VATInput(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "SVI")
	inv := &domain.SupplierInvoice{
		CompanyID: "c1", InvoiceNumber: "INV-VI", SupplierID: sup.ID,
		SupplierName: sup.Name, SupplierTaxCode: sup.TaxCode, InvoiceDate: time.Now(),
		Currency: "VND", VATDeductionStatus: domain.VATClaimed,
		Lines: []domain.SupplierInvoiceLine{{
			ItemName: "Widget", Unit: "pcs", Quantity: 2, UnitPrice: 1000,
			VATRate: 10, VATType: domain.VAT10, AccountID: "152", VATAccountID: "1331",
		}, {
			ItemName: "Service", Unit: "hr", Quantity: 1, UnitPrice: 1000,
			VATRate: 10, VATType: domain.VAT10, AccountID: "642", VATAccountID: "1331",
		}},
	}
	require.NoError(t, svc.CreateInvoice(ctx, inv))

	from := time.Now().AddDate(0, -1, 0).Format(time.DateOnly)
	to := time.Now().AddDate(0, 1, 0).Format(time.DateOnly)
	rpt, err := svc.GetVATInputReport(ctx, "c1", from, to)
	require.NoError(t, err)
	require.Len(t, rpt.Rows, 1)
	assert.Equal(t, 10.0, rpt.Rows[0].VATRate)
	assert.Equal(t, 3000.0, rpt.Rows[0].Subtotal)
	assert.Equal(t, 300.0, rpt.Rows[0].VATAmount)
	assert.Equal(t, "claimed", rpt.Rows[0].DeductionStatus)
	assert.Equal(t, 300.0, rpt.TotalVAT)
}

func TestReport_UninvoicedReceipts(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "SUI")
	po := makePO(t, svc, ctx, sup.ID, 5, 1000)
	grn := makeGRN(t, svc, ctx, po, 5)
	require.NoError(t, svc.PostGRN(ctx, grn.ID))

	rows, err := svc.GetUninvoicedReceipts(ctx, "c1")
	require.NoError(t, err)
	require.Len(t, rows, 1)
	assert.Equal(t, grn.GRNNumber, rows[0].GRNNumber)
	assert.Equal(t, 5.0, rows[0].Quantity)
	assert.Equal(t, sup.Name, rows[0].SupplierName)
}

// ─── Requisition ─────────────────────────────────────────────────────────

func makeRequisition(t *testing.T, svc *PurchaseService, ctx context.Context, code string) *domain.PurchaseRequisition {
	t.Helper()
	r := &domain.PurchaseRequisition{
		CompanyID: "c1", RequisitionNumber: "REQ-" + code, RequesterID: "user-" + code,
		RequesterName: "Requester", Priority: "high", CreatedBy: "user",
		Lines: []domain.RequisitionItem{{
			ItemName: "Widget", Unit: "pcs", Quantity: 10,
			EstimatedPrice: 1000, AccountID: "152",
		}},
	}
	require.NoError(t, svc.CreateRequisition(ctx, r))
	return r
}

func TestCreateRequisition(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	r := makeRequisition(t, svc, ctx, "R1")
	require.NotEmpty(t, r.ID)
	loaded, err := svc.GetRequisition(ctx, r.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.ReqDraft, loaded.Status)
	assert.InDelta(t, 10000, loaded.TotalEstimated, 0.001)
	require.Len(t, loaded.Lines, 1)
	assert.Equal(t, r.ID, loaded.Lines[0].RequisitionID)
}

func TestCreateRequisition_Invalid(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	err := svc.CreateRequisition(ctx, &domain.PurchaseRequisition{
		CompanyID: "c1", RequisitionNumber: "REQ-X", RequesterID: "u",
	})
	assert.ErrorIs(t, err, domain.ErrRequisitionLinesRequired)
}

func TestRequisition_Workflow(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	r := makeRequisition(t, svc, ctx, "R2")
	sup := makeSupplier(t, svc, ctx, "SR2")

	require.NoError(t, svc.SubmitRequisition(ctx, r.ID, "u"))
	require.NoError(t, svc.ApproveRequisition(ctx, r.ID, "boss"))
	po, err := svc.ConvertRequisitionToPO(ctx, r.ID, sup.ID)
	require.NoError(t, err)
	assert.Equal(t, r.ID, po.RequisitionID)
	assert.Equal(t, sup.ID, po.SupplierID)

	loaded, _ := svc.GetRequisition(ctx, r.ID)
	assert.Equal(t, domain.ReqOrdered, loaded.Status)
}

func TestRequisition_ApproveWithoutSubmitFails(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	r := makeRequisition(t, svc, ctx, "R3")
	err := svc.ApproveRequisition(ctx, r.ID, "boss")
	assert.ErrorIs(t, err, domain.ErrRequisitionInvalidTransition)
}

func TestRequisition_Reject(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	r := makeRequisition(t, svc, ctx, "R4")
	require.NoError(t, svc.SubmitRequisition(ctx, r.ID, "u"))
	require.NoError(t, svc.RejectRequisition(ctx, r.ID, "budget cut"))
	loaded, _ := svc.GetRequisition(ctx, r.ID)
	assert.Equal(t, domain.ReqRejected, loaded.Status)
	assert.Equal(t, "budget cut", loaded.RejectedReason)
}

func TestConvertRequisitionToPO_NotApproved(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	r := makeRequisition(t, svc, ctx, "R5")
	sup := makeSupplier(t, svc, ctx, "SR5")
	_, err := svc.ConvertRequisitionToPO(ctx, r.ID, sup.ID)
	assert.ErrorIs(t, err, domain.ErrRequisitionNotApproved)
}

func TestRequisition_UpdateDraftOnly(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	r := makeRequisition(t, svc, ctx, "R6")
	require.NoError(t, svc.SubmitRequisition(ctx, r.ID, "u"))
	err := svc.UpdateRequisition(ctx, &domain.PurchaseRequisition{ID: r.ID, RequisitionNumber: r.RequisitionNumber, RequesterID: r.RequesterID})
	assert.ErrorIs(t, err, domain.ErrRequisitionCannotUpdate)
}

func TestListRequisitions_FilterByStatus(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	makeRequisition(t, svc, ctx, "R7A")
	r2 := makeRequisition(t, svc, ctx, "R7B")
	require.NoError(t, svc.SubmitRequisition(ctx, r2.ID, "u"))

	list, total, err := svc.ListRequisitions(ctx, domain.RequisitionFilter{CompanyID: "c1", Status: domain.ReqPending})
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	require.Len(t, list, 1)
	assert.Equal(t, r2.ID, list[0].ID)
}

// ─── Purchase Return / Credit Note (P2-2) ────────────────────────────────

func makePostedInvoice(t *testing.T, svc *PurchaseService, ctx context.Context) (*domain.SupplierInvoice, *domain.GRN, *domain.PurchaseOrder) {
	t.Helper()
	sup := makeSupplier(t, svc, ctx, "PR-SUP")
	po := makePO(t, svc, ctx, sup.ID, 10, 1000)
	grn := makeGRN(t, svc, ctx, po, 10)
	require.NoError(t, svc.PostGRN(ctx, grn.ID))
	inv := makeInvoice(t, svc, ctx, sup.ID, po.ID, grn.ID, po.Lines[0].ID, 10, 1000)
	require.NoError(t, verifyAndPost(t, svc, ctx, inv))
	return inv, grn, po
}

func TestCreateReturnGRN(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	_, grn, _ := makePostedInvoice(t, svc, ctx)
	rg := &domain.GRN{
		CompanyID: "c1", GRNNumber: "GRN-R1", ReturnOfGRNID: grn.ID, ReceiptDate: time.Now(),
		Lines: []domain.GRNItem{{POLineID: grn.Lines[0].POLineID, ItemName: "Widget", Unit: "pcs", QuantityReceived: 2}},
	}
	require.NoError(t, svc.CreateReturnGRN(ctx, rg))
	loaded, err := svc.GetGRN(ctx, rg.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.GRNPosted, loaded.Status)
	assert.Equal(t, grn.ID, loaded.ReturnOfGRNID)
}

func TestCreateReturnGRN_QtyExceeds(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	_, grn, _ := makePostedInvoice(t, svc, ctx)
	rg := &domain.GRN{
		CompanyID: "c1", GRNNumber: "GRN-R2", ReturnOfGRNID: grn.ID, ReceiptDate: time.Now(),
		Lines: []domain.GRNItem{{POLineID: grn.Lines[0].POLineID, ItemName: "Widget", Unit: "pcs", QuantityReceived: 99}},
	}
	assert.ErrorIs(t, svc.CreateReturnGRN(ctx, rg), domain.ErrReturnQtyExceeds)
}

func TestCreateReturnGRN_OriginalNotPosted(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	sup := makeSupplier(t, svc, ctx, "PR-S2")
	po := makePO(t, svc, ctx, sup.ID, 5, 1000)
	grn := makeGRN(t, svc, ctx, po, 5)
	rg := &domain.GRN{
		CompanyID: "c1", GRNNumber: "GRN-R3", ReturnOfGRNID: grn.ID, ReceiptDate: time.Now(),
		Lines: []domain.GRNItem{{POLineID: po.Lines[0].ID, ItemName: "Widget", Unit: "pcs", QuantityReceived: 1}},
	}
	assert.ErrorIs(t, svc.CreateReturnGRN(ctx, rg), domain.ErrReturnGRNNotPosted)
}

func TestCreateReturnGRN_ReturnAgain(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	_, grn, _ := makePostedInvoice(t, svc, ctx)
	rg := &domain.GRN{
		CompanyID: "c1", GRNNumber: "GRN-R4", ReturnOfGRNID: grn.ID, ReceiptDate: time.Now(),
		Lines: []domain.GRNItem{{POLineID: grn.Lines[0].POLineID, ItemName: "Widget", Unit: "pcs", QuantityReceived: 2}},
	}
	require.NoError(t, svc.CreateReturnGRN(ctx, rg))
	rg2 := &domain.GRN{
		CompanyID: "c1", GRNNumber: "GRN-R5", ReturnOfGRNID: rg.ID, ReceiptDate: time.Now(),
		Lines: []domain.GRNItem{{POLineID: grn.Lines[0].POLineID, ItemName: "Widget", Unit: "pcs", QuantityReceived: 1}},
	}
	assert.ErrorIs(t, svc.CreateReturnGRN(ctx, rg2), domain.ErrReturnGRNReturnAgain)
}

func TestCreateCreditNote_RequiresOriginal(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	cn := &domain.SupplierInvoice{InvoiceNumber: "CN-1", SupplierID: "s1", InvoiceDate: time.Now(), InvoiceType: domain.InvoiceTypeCreditNote}
	err := svc.CreateCreditNote(ctx, cn)
	assert.ErrorIs(t, err, domain.ErrCreditNoteOriginalRequired)
}

func TestCreateCreditNote_InvalidType(t *testing.T) {
	svc, ctx := newPurchaseTestSvc()
	cn := &domain.SupplierInvoice{InvoiceNumber: "CN-2", SupplierID: "s1", InvoiceDate: time.Now(), OriginalInvoiceID: "INV-1"}
	err := svc.CreateCreditNote(ctx, cn)
	assert.ErrorIs(t, err, domain.ErrCreditNoteTypeInvalid)
}

func TestCreateCreditNote(t *testing.T) {
	svc, gl, ctx := setupPurchaseGL(t)
	inv, _, _ := makePostedInvoice(t, svc, ctx)
	cn := &domain.SupplierInvoice{
		CompanyID: "c1", InvoiceNumber: "CN-3", OriginalInvoiceID: inv.ID,
		InvoiceType: domain.InvoiceTypeCreditNote, InvoiceDate: time.Now(),
		Lines: []domain.SupplierInvoiceLine{{
			ItemName: "Widget", Unit: "pcs", Quantity: 2, UnitPrice: 1000,
			VATRate: 10, VATType: domain.VAT10, AccountID: "152", VATAccountID: "1331",
		}},
	}
	require.NoError(t, svc.CreateCreditNote(ctx, cn))
	loaded, err := svc.GetInvoice(ctx, cn.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.InvoiceDraft, loaded.Status)
	assert.Equal(t, inv.SupplierID, loaded.SupplierID)
	assert.True(t, loaded.TotalAmount < 0)
	assert.InDelta(t, -2200, loaded.TotalAmount, 0.001)
	_ = gl
}

func TestPostCreditNote(t *testing.T) {
	svc, gl, ctx := setupPurchaseGL(t)
	inv, _, _ := makePostedInvoice(t, svc, ctx)
	cn := &domain.SupplierInvoice{
		CompanyID: "c1", InvoiceNumber: "CN-4", OriginalInvoiceID: inv.ID,
		InvoiceType: domain.InvoiceTypeCreditNote, InvoiceDate: time.Now(),
		Lines: []domain.SupplierInvoiceLine{{
			ItemName: "Widget", Unit: "pcs", Quantity: 2, UnitPrice: 1000,
			VATRate: 10, VATType: domain.VAT10, AccountID: "152", VATAccountID: "1331",
		}},
	}
	require.NoError(t, svc.CreateCreditNote(ctx, cn))
	require.NoError(t, svc.VerifyInvoice(ctx, cn.ID))
	require.NoError(t, svc.PostInvoice(ctx, cn.ID))

	original, _ := svc.GetInvoice(ctx, inv.ID)
	assert.InDelta(t, 11000-2200, original.BalanceDue, 0.001)

	txns, err := svc.ListAPTransactionsBySupplier(ctx, "c1", inv.SupplierID)
	require.NoError(t, err)
	found := false
	for _, txn := range txns {
		if txn.TransactionType == domain.APTransCreditNote {
			found = true
			assert.InDelta(t, 2200, txn.Amount, 0.001)
		}
	}
	assert.True(t, found, "credit note AP transaction should exist")

	entries, err := gl.GetEntriesByDateRange(ctx, time.Now().Add(-time.Hour), time.Now().Add(time.Hour))
	require.NoError(t, err)
	entry := findEntry(entries, "CN-4")
	require.NotNil(t, entry)
	var debit float64
	for _, l := range entry.Lines {
		if l.AccountCode == "331" {
			debit = l.DebitAmount
		}
	}
	assert.InDelta(t, 2200, debit, 0.001)
}

func TestPostCreditNote_ExceedsBalance(t *testing.T) {
	svc, _, ctx := setupPurchaseGL(t)
	inv, _, _ := makePostedInvoice(t, svc, ctx)
	cn := &domain.SupplierInvoice{
		CompanyID: "c1", InvoiceNumber: "CN-5", OriginalInvoiceID: inv.ID,
		InvoiceType: domain.InvoiceTypeCreditNote, InvoiceDate: time.Now(),
		Lines: []domain.SupplierInvoiceLine{{
			ItemName: "Widget", Unit: "pcs", Quantity: 100, UnitPrice: 1000,
			VATRate: 10, VATType: domain.VAT10, AccountID: "152", VATAccountID: "1331",
		}},
	}
	require.NoError(t, svc.CreateCreditNote(ctx, cn))
	require.NoError(t, svc.VerifyInvoice(ctx, cn.ID))
	assert.ErrorIs(t, svc.PostInvoice(ctx, cn.ID), domain.ErrCreditNoteExceedsBalance)
}

func findEntry(entries []domain.JournalEntry, number string) *domain.JournalEntry {
	for i := range entries {
		if entries[i].EntryNumber == number {
			return &entries[i]
		}
	}
	return nil
}
