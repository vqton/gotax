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

func setupGLEnabled(t *testing.T) (*SaleService, Service, context.Context) {
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

	repo := repository.NewMemorySaleRepo()
	svc := NewSaleService(repo, repo, repo, repo, repo, repo, repo, repo, gl)
	return svc, gl, context.Background()
}

func seedGLAccounts(t *testing.T, gl Service, ctx context.Context) {
	t.Helper()
	for _, acc := range []domain.Account{
		{Code: "131", Name: "Phai thu khach hang", Type: domain.AccountTypeAsset, IsActive: true},
		{Code: "5111", Name: "Doanh thu ban hang hoa", Type: domain.AccountTypeRevenue, IsActive: true},
		{Code: "3331", Name: "Thue GTGT phai nop", Type: domain.AccountTypeLiability, IsActive: true},
		{Code: "1111", Name: "Tien mat VND", Type: domain.AccountTypeAsset, IsActive: true},
		{Code: "1121", Name: "Tien gui ngan hang VND", Type: domain.AccountTypeAsset, IsActive: true},
		{Code: "5213", Name: "Chiet khau thanh toan", Type: domain.AccountTypeExpense, IsActive: true},
		{Code: "515", Name: "Doanh thu hoat dong tai chinh", Type: domain.AccountTypeRevenue, IsActive: true},
		{Code: "635", Name: "Chi phi tai chinh", Type: domain.AccountTypeExpense, IsActive: true},
		{Code: "642", Name: "Chi phi quan ly doanh nghiep", Type: domain.AccountTypeExpense, IsActive: true},
		{Code: "229", Name: "Du phong phai thu kho doi", Type: domain.AccountTypeLiability, IsActive: true},
		{Code: "632", Name: "Gia von hang ban", Type: domain.AccountTypeExpense, IsActive: true},
		{Code: "156", Name: "Hang hoa", Type: domain.AccountTypeAsset, IsActive: true},
	} {
		err := gl.CreateAccount(ctx, &acc)
		require.NoError(t, err)
	}
}

func seedPeriod(t *testing.T, gl Service, ctx context.Context) {
	t.Helper()
	now := time.Now().UTC()
	_, m, _ := now.Date()
	y := now.Year()
	err := gl.CreatePeriod(ctx, &domain.Period{
		Year: y, Month: int(m),
		StartDate: time.Date(y, m, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(y, m+1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second),
		Status:    domain.PeriodOpen,
	})
	require.NoError(t, err)
}

func TestPostInvoice_CreatesGLEntry(t *testing.T) {
	svc, gl, ctx := setupGLEnabled(t)
	seedGLAccounts(t, gl, ctx)
	seedPeriod(t, gl, ctx)

	cust := &domain.Customer{
		CompanyID: "c1", Code: "C001", Name: "TestCo",
		TaxCode: "1234567890", Currency: "VND",
	}
	err := svc.CreateCustomer(ctx, cust)
	require.NoError(t, err)

	now := time.Now().UTC().Truncate(24 * time.Hour)
	inv := &domain.CustomerInvoice{
		CompanyID: "c1", InvoiceNumber: "INV-001",
		InvoiceDate: now, CustomerID: cust.ID,
		CustomerName: "TestCo", CustomerTaxCode: "1234567890",
		CustomerAddress: "addr", InvoiceType: "domestic",
		Currency: "VND", Status: domain.SInvDraft,
		Lines: []domain.InvLine{
			{ItemName: "Item A", Quantity: 2, UnitPrice: 100000,
				VATRate: 10, VATType: "VAT_10",
				RevenueAccount: "5111", VATAccountID: "3331"},
		},
		CreatedBy: "test-user",
	}
	err = svc.CreateInvoice(ctx, inv)
	require.NoError(t, err)
	require.NotEmpty(t, inv.ID)

	// Post invoice — should trigger GL entry
	err = svc.PostInvoice(ctx, inv.ID)
	require.NoError(t, err)

	// Verify invoice GL flags
	posted, err := svc.GetInvoice(ctx, inv.ID)
	require.NoError(t, err)
	assert.True(t, posted.GLPosted)
	assert.NotNil(t, posted.GLPostedAt)

	// Verify journal entry exists in GL
	entries, err := gl.GetEntriesByDateRange(ctx, now, now.Add(24*time.Hour))
	require.NoError(t, err)
	require.NotEmpty(t, entries, "should have at least one journal entry")

	// Find the AR entry
	var arEntry *domain.JournalEntry
	for i := range entries {
		if entries[i].Description == "AR invoice INV-001" {
			arEntry = &entries[i]
			break
		}
	}
	require.NotNil(t, arEntry, "AR journal entry not found")
	assert.Equal(t, domain.JournalEntryPosted, arEntry.Status)

	// Verify lines: Dr 131 (total), Cr 5111 (subtotal), Cr 3331 (VAT)
	var totalDebit, totalCredit float64
	var hasAR, hasRevenue, hasVAT bool
	for _, l := range arEntry.Lines {
		totalDebit += l.DebitAmount
		totalCredit += l.CreditAmount
		if l.AccountCode == "131" && l.DebitAmount > 0 {
			hasAR = true
		}
		if l.AccountCode == "5111" && l.CreditAmount > 0 {
			hasRevenue = true
		}
		if l.AccountCode == "3331" && l.CreditAmount > 0 {
			hasVAT = true
		}
	}
	assert.True(t, hasAR, "must have Dr 131")
	assert.True(t, hasRevenue, "must have Cr 5111")
	assert.True(t, hasVAT, "must have Cr 3331")
	assert.InDelta(t, totalDebit, totalCredit, 0.001, "debit must equal credit")
	assert.InDelta(t, 220000, totalDebit, 0.001, "total = 200K + 20K VAT")
}

func TestPostReceipt_CreatesGLEntry(t *testing.T) {
	svc, gl, ctx := setupGLEnabled(t)
	seedGLAccounts(t, gl, ctx)
	seedPeriod(t, gl, ctx)

	cust := &domain.Customer{
		CompanyID: "c1", Code: "C001", Name: "TestCo",
		TaxCode: "1234567890", Currency: "VND",
	}
	err := svc.CreateCustomer(ctx, cust)
	require.NoError(t, err)

	now := time.Now().UTC().Truncate(24 * time.Hour)
	rcpt := &domain.CustomerReceipt{
		CompanyID:     "c1",
		ReceiptNumber: "RCP-001",
		CustomerID:    cust.ID,
		ReceiptDate:   now,
		PaymentMethod: "bank_transfer",
		Currency:      "VND",
		Amount:        110000,
		Status:        domain.RcpDraft,
	}
	err = svc.CreateReceipt(ctx, rcpt)
	require.NoError(t, err)

	// Post receipt — should trigger GL entry
	err = svc.PostReceipt(ctx, rcpt.ID)
	require.NoError(t, err)

	// Verify journal entry exists
	entries, err := gl.GetEntriesByDateRange(ctx, now, now.Add(24*time.Hour))
	require.NoError(t, err)

	var arEntry *domain.JournalEntry
	for i := range entries {
		if entries[i].Description == "AR receipt RCP-001" {
			arEntry = &entries[i]
			break
		}
	}
	require.NotNil(t, arEntry, "receipt journal entry not found")
	assert.Equal(t, domain.JournalEntryPosted, arEntry.Status)

	var totalDebit, totalCredit float64
	var hasBank, hasAR bool
	for _, l := range arEntry.Lines {
		totalDebit += l.DebitAmount
		totalCredit += l.CreditAmount
		if l.AccountCode == "1121" && l.DebitAmount > 0 {
			hasBank = true
		}
		if l.AccountCode == "131" && l.CreditAmount > 0 {
			hasAR = true
		}
	}
	assert.True(t, hasBank, "must have Dr bank account (1121 for bank_transfer)")
	assert.True(t, hasAR, "must have Cr 131 (AR)")
	assert.InDelta(t, totalDebit, totalCredit, 0.001)
	assert.InDelta(t, 110000, totalDebit, 0.001)
}

func TestPostReceipt_CashMethod(t *testing.T) {
	svc, gl, ctx := setupGLEnabled(t)
	seedGLAccounts(t, gl, ctx)
	seedPeriod(t, gl, ctx)

	cust := &domain.Customer{
		CompanyID: "c1", Code: "C002", Name: "CashCo",
		TaxCode: "1234567890", Currency: "VND",
	}
	err := svc.CreateCustomer(ctx, cust)
	require.NoError(t, err)

	now := time.Now().UTC().Truncate(24 * time.Hour)
	rcpt := &domain.CustomerReceipt{
		CompanyID: "c1", ReceiptNumber: "RCP-CASH",
		CustomerID: cust.ID, ReceiptDate: now,
		PaymentMethod: "cash", Currency: "VND",
		Amount: 50000, Status: domain.RcpDraft,
	}
	err = svc.CreateReceipt(ctx, rcpt)
	require.NoError(t, err)
	err = svc.PostReceipt(ctx, rcpt.ID)
	require.NoError(t, err)

	entries, err := gl.GetEntriesByDateRange(ctx, now, now.Add(24*time.Hour))
	require.NoError(t, err)

	var hasCash bool
	for i := range entries {
		for _, l := range entries[i].Lines {
			if l.AccountCode == "1111" && l.DebitAmount > 0 {
				hasCash = true
			}
		}
	}
	assert.True(t, hasCash, "cash receipt should Dr 1111")
}

func TestPostInvoice_NoGLEntryIfGLNil(t *testing.T) {
	repo := repository.NewMemorySaleRepo()
	svc := NewSaleService(repo, repo, repo, repo, repo, repo, repo, repo, nil)
	ctx := context.Background()

	cust := &domain.Customer{
		CompanyID: "c1", Code: "C002", Name: "TestCo",
		TaxCode: "1234567890", Currency: "VND",
	}
	err := svc.CreateCustomer(ctx, cust)
	require.NoError(t, err)

	inv := &domain.CustomerInvoice{
		CompanyID: "c1", InvoiceNumber: "INV-NOGL",
		InvoiceDate: time.Now(), CustomerID: cust.ID,
		CustomerName: "TestCo", CustomerTaxCode: "1234567890",
		CustomerAddress: "addr", InvoiceType: "domestic",
		Currency: "VND", Status: domain.SInvDraft,
		Lines: []domain.InvLine{{ItemName: "item", Quantity: 1, UnitPrice: 100,
			VATRate: 10, VATType: "VAT_10",
			RevenueAccount: "5111", VATAccountID: "3331"}},
		CreatedBy: "test-user",
	}
	err = svc.CreateInvoice(ctx, inv)
	require.NoError(t, err)

	// simulate enough status transitions to reach POSTED
	// need e-invoice status to allow post
	err = svc.invRepo.PostInvoice(ctx, inv.ID, time.Now().UTC())
	require.NoError(t, err)

	// should not crash even with nil GL — GLPosted stays false
	posted, err := svc.GetInvoice(ctx, inv.ID)
	require.NoError(t, err)
	assert.False(t, posted.GLPosted)
}

func TestPostCN_CreatesGLEntry(t *testing.T) {
	svc, gl, ctx := setupGLEnabled(t)
	seedGLAccounts(t, gl, ctx)
	seedPeriod(t, gl, ctx)

	cust := &domain.Customer{
		CompanyID: "c1", Code: "C003", Name: "CNCo",
		TaxCode: "1234567890", Currency: "VND",
	}
	err := svc.CreateCustomer(ctx, cust)
	require.NoError(t, err)

	inv := &domain.CustomerInvoice{
		CompanyID: "c1", InvoiceNumber: "INV-CN",
		InvoiceDate: time.Now(), CustomerID: cust.ID,
		CustomerName: "CNCo", CustomerTaxCode: "1234567890",
		CustomerAddress: "addr", InvoiceType: "domestic",
		Currency: "VND", Status: domain.SInvDraft,
		Lines: []domain.InvLine{{ItemName: "item", Quantity: 1, UnitPrice: 100,
			VATRate: 10, VATType: "VAT_10",
			RevenueAccount: "5111", VATAccountID: "3331"}},
		CreatedBy: "test-user",
	}
	err = svc.CreateInvoice(ctx, inv)
	require.NoError(t, err)
	// post invoice first
	err = svc.PostInvoice(ctx, inv.ID)
	require.NoError(t, err)

	now := time.Now().UTC().Truncate(24 * time.Hour)
	cn := &domain.CreditNote{
		CompanyID:         "c1",
		CNNumber:          "CN-001",
		OriginalInvoiceID: inv.ID,
		CustomerID:        cust.ID,
		ReturnDate:        now,
		ReturnReason:      "defective",
		ReturnType:        domain.RetFull,
		Subtotal:          100,
		TaxAmount:         10,
		TotalAmount:       110,
		Status:            domain.CNDraft,
		CreatedBy:         "test-user",
		Lines: []domain.CNLine{
			{ItemName: "item", Quantity: 1, UnitPrice: 100, VATRate: 10, LineTotal: 100, LineVATAmount: 10},
		},
	}
	err = svc.CreateCN(ctx, cn)
	require.NoError(t, err)

	err = svc.PostCN(ctx, cn.ID)
	require.NoError(t, err)

	entries, err := gl.GetEntriesByDateRange(ctx, now, now.Add(24*time.Hour))
	require.NoError(t, err)

	var cnEntry *domain.JournalEntry
	for i := range entries {
		if entries[i].Description == "AR credit note CN-001" {
			cnEntry = &entries[i]
			break
		}
	}
	require.NotNil(t, cnEntry)
	assert.Equal(t, domain.JournalEntryPosted, cnEntry.Status)

	var totalDebit, totalCredit float64
	var hasAR, hasRev, hasVAT bool
	for _, l := range cnEntry.Lines {
		totalDebit += l.DebitAmount
		totalCredit += l.CreditAmount
		if l.AccountCode == "131" && l.CreditAmount > 0 {
			hasAR = true
		}
		if l.AccountCode == "5111" && l.DebitAmount > 0 {
			hasRev = true
		}
		if l.AccountCode == "3331" && l.DebitAmount > 0 {
			hasVAT = true
		}
	}
	assert.True(t, hasAR, "must Cr 131")
	assert.True(t, hasRev, "must Dr revenue account")
	assert.True(t, hasVAT, "must Dr 3331 (VAT)")
	assert.InDelta(t, totalDebit, totalCredit, 0.001)
	assert.InDelta(t, 110, totalDebit, 0.001)
}

func TestARGLReconciliation_Matches(t *testing.T) {
	svc, gl, ctx := setupGLEnabled(t)
	seedGLAccounts(t, gl, ctx)
	seedPeriod(t, gl, ctx)

	cust := &domain.Customer{
		CompanyID: "c1", Code: "RECON01", Name: "ReconCo",
		TaxCode: "999", Currency: "VND",
	}
	err := svc.CreateCustomer(ctx, cust)
	require.NoError(t, err)

	// create invoice that will post GL entry Dr 131
	inv := &domain.CustomerInvoice{
		CompanyID: "c1", InvoiceNumber: "INV-RECON",
		InvoiceDate: time.Now(), CustomerID: cust.ID,
		CustomerName: "ReconCo", CustomerTaxCode: "999",
		CustomerAddress: "addr", InvoiceType: "domestic",
		Currency: "VND", Status: domain.SInvDraft,
		Lines: []domain.InvLine{{ItemName: "svc", Quantity: 1, UnitPrice: 500,
			VATRate: 10, VATType: "VAT_10",
			RevenueAccount: "5111", VATAccountID: "3331"}},
		CreatedBy: "user",
	}
	err = svc.CreateInvoice(ctx, inv)
	require.NoError(t, err)
	err = svc.PostInvoice(ctx, inv.ID)
	require.NoError(t, err)

	period, _ := gl.GetPeriodByYearMonth(ctx, time.Now().Year(), int(time.Now().Month()))

	// run reconciliation
	recon, err := svc.GetARGLReconciliation(ctx, "c1", period.ID)
	require.NoError(t, err)
	require.NotNil(t, recon)
	assert.Equal(t, period.ID, recon.PeriodID)
	assert.InDelta(t, 550.0, recon.SubledgerTotal, 0.001) // 500 + 50 VAT
	assert.InDelta(t, 550.0, recon.GLBalance, 0.001)      // Dr 131 posted
	assert.InDelta(t, 0, recon.Variance, 0.001)
}

func TestARGLReconciliation_Variance(t *testing.T) {
	svc, gl, ctx := setupGLEnabled(t)
	seedPeriod(t, gl, ctx)

	cust := &domain.Customer{
		CompanyID: "c1", Code: "RECON02", Name: "VarianceCo",
		TaxCode: "999", Currency: "VND",
	}
	err := svc.CreateCustomer(ctx, cust)
	require.NoError(t, err)

	// create invoice WITHOUT GL posting (nil GL in sale service)
	// but we still have GL account 131 with zero balance
	repo := repository.NewMemorySaleRepo()
	noGLSvc := NewSaleService(repo, repo, repo, repo, repo, repo, repo, repo, nil)

	err = noGLSvc.CreateCustomer(ctx, cust)
	require.NoError(t, err)
	inv := &domain.CustomerInvoice{
		CompanyID: "c1", InvoiceNumber: "INV-VAR",
		InvoiceDate: time.Now(), CustomerID: cust.ID,
		CustomerName: "VarianceCo", CustomerTaxCode: "999",
		CustomerAddress: "addr", InvoiceType: "domestic",
		Currency: "VND", Status: domain.SInvDraft,
		Lines: []domain.InvLine{{ItemName: "svc", Quantity: 1, UnitPrice: 300,
			VATRate: 0, VATType: "NON_TAX",
			RevenueAccount: "5111", VATAccountID: "3331"}},
		CreatedBy: "user",
	}
	err = noGLSvc.CreateInvoice(ctx, inv)
	require.NoError(t, err)
	err = noGLSvc.invRepo.PostInvoice(ctx, inv.ID, time.Now().UTC())
	require.NoError(t, err)

	period, _ := gl.GetPeriodByYearMonth(ctx, time.Now().Year(), int(time.Now().Month()))

	recon, err := noGLSvc.GetARGLReconciliation(ctx, "c1", period.ID)
	require.NoError(t, err)
	assert.InDelta(t, 300, recon.SubledgerTotal, 0.001)
	assert.InDelta(t, 0, recon.GLBalance, 0.001) // no GL posting → balance is 0
	assert.InDelta(t, -300, recon.Variance, 0.001)
}

// ─── S5: COGS GL on delivery ───────────────────────────────────────────

func TestPostDN_CreatesCOGSGLEntry(t *testing.T) {
	svc, gl, ctx := setupGLEnabled(t)
	seedGLAccounts(t, gl, ctx)
	seedPeriod(t, gl, ctx)
	seedCust(t, svc, ctx, "c1", "co1")
	so := seedSO(t, svc, ctx, "SO-COGS-1", 100)

	dn := &domain.DeliveryNote{
		CompanyID: "co1", DNNumber: "DN-COGS-1", SOID: so.ID,
		DeliveryDate: time.Now().UTC(), Status: domain.DNDraft,
		Lines: []domain.DNLine{{
			SOLineID: so.Lines[0].ID, ItemName: "Widget", Unit: "pcs",
			QtyDelivered: 100, UnitPrice: 100, LineTotal: 10000, CostPrice: 60,
		}},
	}
	require.NoError(t, svc.CreateDN(ctx, dn))
	require.NoError(t, svc.PostDN(ctx, dn.ID))

	entries, err := gl.GetEntriesByDateRange(ctx, time.Now().Add(-24*time.Hour), time.Now().Add(24*time.Hour))
	require.NoError(t, err)
	var entry *domain.JournalEntry
	for i := range entries {
		if entries[i].EntryNumber == "DN-COGS-1" {
			entry = &entries[i]
			break
		}
	}
	require.NotNil(t, entry, "COGS entry not found")
	require.Len(t, entry.Lines, 2)
	assert.Equal(t, "632", entry.Lines[0].AccountCode)
	assert.Equal(t, 6000.0, entry.Lines[0].DebitAmount) // 60×100
	assert.Equal(t, "156", entry.Lines[1].AccountCode)
	assert.Equal(t, 6000.0, entry.Lines[1].CreditAmount)
}

func TestPostDN_SkipsZeroCostLines(t *testing.T) {
	svc, gl, ctx := setupGLEnabled(t)
	seedGLAccounts(t, gl, ctx)
	seedPeriod(t, gl, ctx)
	seedCust(t, svc, ctx, "c1", "co1")
	so := seedSO(t, svc, ctx, "SO-COGS-2", 100)

	dn := &domain.DeliveryNote{
		CompanyID: "co1", DNNumber: "DN-COGS-2", SOID: so.ID,
		DeliveryDate: time.Now().UTC(), Status: domain.DNDraft,
		Lines: []domain.DNLine{
			{SOLineID: so.Lines[0].ID, ItemName: "Widget", Unit: "pcs",
				QtyDelivered: 50, UnitPrice: 100, LineTotal: 5000, CostPrice: 0},
			{SOLineID: so.Lines[0].ID, ItemName: "Widget2", Unit: "pcs",
				QtyDelivered: 50, UnitPrice: 200, LineTotal: 10000, CostPrice: 120},
		},
	}
	require.NoError(t, svc.CreateDN(ctx, dn))
	require.NoError(t, svc.PostDN(ctx, dn.ID))

	entries, err := gl.GetEntriesByDateRange(ctx, time.Now().Add(-24*time.Hour), time.Now().Add(24*time.Hour))
	require.NoError(t, err)
	var entry *domain.JournalEntry
	for i := range entries {
		if entries[i].EntryNumber == "DN-COGS-2" {
			entry = &entries[i]
			break
		}
	}
	require.NotNil(t, entry, "COGS entry not found")
	require.Len(t, entry.Lines, 2) // only costed line
	assert.Equal(t, 6000.0, entry.Lines[0].DebitAmount) // 120×50
}

func TestPostDN_NoGLWhenNil(t *testing.T) {
	repo := repository.NewMemorySaleRepo()
	svc := NewSaleService(repo, repo, repo, repo, repo, repo, repo, repo, nil)
	ctx := context.Background()
	seedCust(t, svc, ctx, "c1", "co1")
	so := seedSO(t, svc, ctx, "SO-COGS-3", 100)

	dn := &domain.DeliveryNote{
		CompanyID: "co1", DNNumber: "DN-COGS-3", SOID: so.ID,
		DeliveryDate: time.Now().UTC(), Status: domain.DNDraft,
		Lines: []domain.DNLine{{
			SOLineID: so.Lines[0].ID, ItemName: "Widget", Unit: "pcs",
			QtyDelivered: 100, UnitPrice: 100, LineTotal: 10000, CostPrice: 60,
		}},
	}
	require.NoError(t, svc.CreateDN(ctx, dn))
	require.NoError(t, svc.PostDN(ctx, dn.ID))

	// no GL → no journal entry, but DN still posted
	dnAfter, err := svc.dnRepo.GetDN(ctx, dn.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.DNPosted, dnAfter.Status)
}

// failOnceDNRepo wraps MemorySaleRepo: UpdateDNStatus fails once, then works.
type failOnceDNRepo struct {
	*repository.MemorySaleRepo
	fail bool
}

func (r *failOnceDNRepo) UpdateDNStatus(ctx context.Context, id string, status domain.DNStatus) error {
	if r.fail {
		r.fail = false
		return domain.ErrDNInvalidTransition
	}
	return r.MemorySaleRepo.UpdateDNStatus(ctx, id, status)
}

// ─── Bug fix: no duplicate COGS GL on partial failure ────────────────

func TestPostDN_NoDuplicateCOGSGLOnStatusFailure(t *testing.T) {
	_, gl, ctx := setupGLEnabled(t)
	seedGLAccounts(t, gl, ctx)
	seedPeriod(t, gl, ctx)

	repo := &failOnceDNRepo{MemorySaleRepo: repository.NewMemorySaleRepo()}
	svc2 := NewSaleService(repo, repo, repo, repo, repo, repo, repo, repo, gl)
	seedCust(t, svc2, ctx, "c1", "co1")
	so2 := seedSO(t, svc2, ctx, "SO-FAIL-1", 100)

	dn := &domain.DeliveryNote{
		CompanyID: "co1", DNNumber: "DN-FAIL-1", SOID: so2.ID,
		DeliveryDate: time.Now().UTC(), Status: domain.DNDraft,
		Lines: []domain.DNLine{{
			SOLineID: so2.Lines[0].ID, ItemName: "Widget", Unit: "pcs",
			QtyDelivered: 100, UnitPrice: 100, LineTotal: 10000, CostPrice: 60,
		}},
	}
	require.NoError(t, svc2.CreateDN(ctx, dn))

	// first post: GL written, then status write fails
	repo.fail = true
	err := svc2.PostDN(ctx, dn.ID)
	require.Error(t, err)

	// retry: GL must NOT be written twice; DN reaches POSTED
	require.NoError(t, svc2.PostDN(ctx, dn.ID))

	entries, err := gl.GetEntriesByDateRange(ctx, time.Now().Add(-24*time.Hour), time.Now().Add(24*time.Hour))
	require.NoError(t, err)
	cogs := 0
	for i := range entries {
		if entries[i].EntryNumber == "DN-FAIL-1" {
			cogs++
		}
	}
	assert.Equal(t, 1, cogs, "COGS entry written exactly once across retry")

	dnAfter, err := svc2.dnRepo.GetDN(ctx, dn.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.DNPosted, dnAfter.Status)
}
