package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
	"gotax/internal/repository"
	"gotax/internal/service"
)

func setupKeeperTest(t *testing.T) (*gin.Engine, *service.WarehouseKeeperService) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	keeperRepo := repository.NewMemoryWarehouseKeeperRepo()
	whRepo := repository.NewMemoryWarehouseRepo()
	itemRepo := repository.NewMemoryItemRepo()
	keeperSvc := service.NewWarehouseKeeperService(keeperRepo, whRepo, itemRepo)
	keeperH := NewWarehouseKeeperHandler(keeperSvc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
		c.Next()
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterKeeperRoutes(r, keeperH, noopMW)

	return r, keeperSvc
}

const testCompanyID = "CMP001"

// ─── Assignment ────────────────────────────────────────────────────────

func TestKeeperCreateAssignment(t *testing.T) {
	r, _ := setupKeeperTest(t)
	body := `{"warehouse_id":"WH001","user_id":"USR001","role":"keeper","effective_from":"2026-01-01"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/warehouse/keeper/assignments?company_id="+testCompanyID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	var a domain.WarehouseKeeperAssignment
	json.Unmarshal(w.Body.Bytes(), &a)
	assert.NotEmpty(t, a.ID)
	assert.Equal(t, domain.KeeperRoleKeeper, a.Role)
}

func TestKeeperListAssignments(t *testing.T) {
	r, svc := setupKeeperTest(t)
	svc.CreateAssignment(t.Context(), &domain.WarehouseKeeperAssignment{
		CompanyID: testCompanyID, WarehouseID: "WH001", UserID: "USR001",
		Role: domain.KeeperRoleKeeper, EffectiveFrom: mustTime("2026-01-01"),
	}, "test-user")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/warehouse/keeper/assignments?company_id="+testCompanyID, nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var result []domain.WarehouseKeeperAssignment
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Len(t, result, 1)
}

func TestKeeperGetAssignment(t *testing.T) {
	r, svc := setupKeeperTest(t)
	a := &domain.WarehouseKeeperAssignment{
		CompanyID: testCompanyID, WarehouseID: "WH001", UserID: "USR001",
		Role: domain.KeeperRoleKeeper, EffectiveFrom: mustTime("2026-01-01"),
	}
	svc.CreateAssignment(t.Context(), a, "test-user")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/warehouse/keeper/assignments/"+a.ID, nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestKeeperDeleteAssignment(t *testing.T) {
	r, svc := setupKeeperTest(t)
	a := &domain.WarehouseKeeperAssignment{
		CompanyID: testCompanyID, WarehouseID: "WH001", UserID: "USR001",
		Role: domain.KeeperRoleKeeper, EffectiveFrom: mustTime("2026-01-01"),
	}
	svc.CreateAssignment(t.Context(), a, "test-user")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/warehouse/keeper/assignments/"+a.ID, nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

// ─── Stock Ledger ──────────────────────────────────────────────────────

func TestKeeperRecordSlips(t *testing.T) {
	r, svc := setupKeeperTest(t)
	svc.CreateAssignment(t.Context(), &domain.WarehouseKeeperAssignment{
		CompanyID: testCompanyID, WarehouseID: "WH001", UserID: "test-user",
		Role: domain.KeeperRoleKeeper, EffectiveFrom: mustTime("2026-01-01"),
	}, "test-user")

	body := `{"warehouse_id":"WH001","slips":[{"item_id":"ITM001","voucher_type":"receipt","voucher_no":"GRN-001","receipt_qty":10}]}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/warehouse/keeper/ledger/record?company_id="+testCompanyID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestKeeperListLedgerEntries(t *testing.T) {
	r, svc := setupKeeperTest(t)
	svc.CreateAssignment(t.Context(), &domain.WarehouseKeeperAssignment{
		CompanyID: testCompanyID, WarehouseID: "WH001", UserID: "test-user",
		Role: domain.KeeperRoleKeeper, EffectiveFrom: mustTime("2026-01-01"),
	}, "test-user")
	svc.RecordSlips(t.Context(), testCompanyID, "WH001", []service.RecordSlipRequest{
		{ItemID: "ITM001", VoucherType: domain.VoucherReceipt, VoucherNo: "GRN-001", ReceiptQty: 10},
	}, "test-user")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/warehouse/keeper/ledger?company_id="+testCompanyID+"&warehouse_id=WH001", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestKeeperUnrecordEntry(t *testing.T) {
	r, svc := setupKeeperTest(t)
	svc.CreateAssignment(t.Context(), &domain.WarehouseKeeperAssignment{
		CompanyID: testCompanyID, WarehouseID: "WH001", UserID: "test-user",
		Role: domain.KeeperRoleKeeper, EffectiveFrom: mustTime("2026-01-01"),
	}, "test-user")
	svc.RecordSlips(t.Context(), testCompanyID, "WH001", []service.RecordSlipRequest{
		{ItemID: "ITM001", VoucherType: domain.VoucherReceipt, VoucherNo: "GRN-001", ReceiptQty: 10},
	}, "test-user")

	entries, _, _ := svc.ListLedgerEntries(t.Context(), domain.LedgerFilter{
		CompanyID: testCompanyID, WarehouseID: "WH001",
	})
	require.NotEmpty(t, entries)

	body := `{"entry_id":"` + entries[0].ID + `","reason":"wrong entry"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/warehouse/keeper/ledger/unrecord", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestKeeperGetLedgerBalance(t *testing.T) {
	r, svc := setupKeeperTest(t)
	svc.CreateAssignment(t.Context(), &domain.WarehouseKeeperAssignment{
		CompanyID: testCompanyID, WarehouseID: "WH001", UserID: "test-user",
		Role: domain.KeeperRoleKeeper, EffectiveFrom: mustTime("2026-01-01"),
	}, "test-user")
	svc.RecordSlips(t.Context(), testCompanyID, "WH001", []service.RecordSlipRequest{
		{ItemID: "ITM001", VoucherType: domain.VoucherReceipt, VoucherNo: "GRN-001", ReceiptQty: 10},
	}, "test-user")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/warehouse/keeper/ledger/balance?company_id="+testCompanyID+"&warehouse_id=WH001&item_id=ITM001", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var result struct {
		Balance float64 `json:"balance"`
	}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, float64(10), result.Balance)
}

// ─── Pending Slips ─────────────────────────────────────────────────────

func TestKeeperPendingSlips(t *testing.T) {
	r, _ := setupKeeperTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/warehouse/keeper/pending-slips?company_id="+testCompanyID+"&warehouse_id=WH001", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestKeeperPendingSlipsCount(t *testing.T) {
	r, _ := setupKeeperTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/warehouse/keeper/pending-slips/count?company_id="+testCompanyID+"&warehouse_id=WH001", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

// ─── Reconciliation ────────────────────────────────────────────────────

func TestKeeperReconciliation(t *testing.T) {
	r, _ := setupKeeperTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/warehouse/keeper/reconciliation?company_id="+testCompanyID+"&warehouse_id=WH001&from=2026-01-01&to=2026-12-31", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

// ─── Stock Card ────────────────────────────────────────────────────────

func TestKeeperStockCard(t *testing.T) {
	r, _ := setupKeeperTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/warehouse/keeper/stock-card?company_id="+testCompanyID+"&warehouse_id=WH001&item_id=ITM001&period=2026-01", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

// ─── Reports ───────────────────────────────────────────────────────────

func TestKeeperInventorySummary(t *testing.T) {
	r, _ := setupKeeperTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/warehouse/keeper/reports/inventory-summary?company_id="+testCompanyID+"&warehouse_id=WH001", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

// ─── Validation ────────────────────────────────────────────────────────

func TestKeeperCreateAssignmentInvalidRole(t *testing.T) {
	r, _ := setupKeeperTest(t)
	body := `{"warehouse_id":"WH001","user_id":"USR001","role":"invalid","effective_from":"2026-01-01"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/warehouse/keeper/assignments?company_id="+testCompanyID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestKeeperCreateAssignmentMissingCompany(t *testing.T) {
	r, _ := setupKeeperTest(t)
	body := `{"warehouse_id":"WH001","user_id":"USR001","role":"keeper","effective_from":"2026-01-01"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/warehouse/keeper/assignments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code) // company_id required via query param
}

// ─── Helpers ───────────────────────────────────────────────────────────

func mustTime(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}
