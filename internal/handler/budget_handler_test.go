package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
	"gotax/internal/repository"
	"gotax/internal/service"
)

func setupBudgetHandlerTest(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	budRepo := repository.NewMemoryBudgetRepo()
	jeRepo := repository.NewMemoryJournalRepo()
	svc := service.NewBudgetService(budRepo, jeRepo)
	bh := NewBudgetHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
		c.Next()
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterBudgetRoutes(r, bh, noopMW)

	return r
}

func TestBudgetCreate(t *testing.T) {
	r := setupBudgetHandlerTest(t)
	body := map[string]interface{}{
		"account_code": "6421",
		"period_year":  2026,
		"period_month": 1,
		"budgeted":     50000000,
		"actual":       45000000,
		"notes":        "Rent budget",
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/budgets?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var resp domain.Budget
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "6421", resp.AccountCode)
	assert.Equal(t, 2026, resp.PeriodYear)
	assert.Equal(t, 1, resp.PeriodMonth)
	assert.Equal(t, 50000000.0, resp.Budgeted)
	assert.Equal(t, 45000000.0, resp.Actual)
	assert.Equal(t, -5000000.0, resp.Variance)
	assert.NotEmpty(t, resp.ID)
}

func TestBudgetCreate_MissingCompany(t *testing.T) {
	r := setupBudgetHandlerTest(t)
	body := map[string]interface{}{
		"account_code": "6421",
		"period_year":  2026,
		"period_month": 1,
		"budgeted":     50000000,
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/budgets", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBudgetCreate_Duplicate(t *testing.T) {
	r := setupBudgetHandlerTest(t)
	body := map[string]interface{}{
		"account_code": "6421",
		"period_year":  2026,
		"period_month": 1,
		"budgeted":     50000000,
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/budgets?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/v1/budgets?company_id=comp1", bytes.NewReader(b))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusConflict, w2.Code)
}

func TestBudgetList(t *testing.T) {
	r := setupBudgetHandlerTest(t)
	body := map[string]interface{}{
		"account_code": "6421",
		"period_year":  2026,
		"period_month": 1,
		"budgeted":     50000000,
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/budgets?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/v1/budgets?company_id=comp1&year=2026", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	var list []domain.Budget
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &list))
	assert.Len(t, list, 1)
}

func TestBudgetGet(t *testing.T) {
	r := setupBudgetHandlerTest(t)
	body := map[string]interface{}{
		"account_code": "6421",
		"period_year":  2026,
		"period_month": 1,
		"budgeted":     50000000,
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/budgets?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	var created domain.Budget
	json.Unmarshal(w.Body.Bytes(), &created)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/v1/budgets/"+created.ID, nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestBudgetGet_NotFound(t *testing.T) {
	r := setupBudgetHandlerTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/budgets/nonexistent", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBudgetUpdate(t *testing.T) {
	r := setupBudgetHandlerTest(t)
	body := map[string]interface{}{
		"account_code": "6421",
		"period_year":  2026,
		"period_month": 1,
		"budgeted":     50000000,
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/budgets?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	var created domain.Budget
	json.Unmarshal(w.Body.Bytes(), &created)

	body["budgeted"] = 60000000
	b2, _ := json.Marshal(body)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("PUT", "/api/v1/budgets/"+created.ID+"?company_id=comp1", bytes.NewReader(b2))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	var updated domain.Budget
	json.Unmarshal(w2.Body.Bytes(), &updated)
	assert.Equal(t, 60000000.0, updated.Budgeted)
}

func TestBudgetDelete(t *testing.T) {
	r := setupBudgetHandlerTest(t)
	body := map[string]interface{}{
		"account_code": "6421",
		"period_year":  2026,
		"period_month": 1,
		"budgeted":     50000000,
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/budgets?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	var created domain.Budget
	json.Unmarshal(w.Body.Bytes(), &created)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("DELETE", "/api/v1/budgets/"+created.ID, nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/api/v1/budgets/"+created.ID, nil)
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusNotFound, w3.Code)
}

func TestBudgetBulkUpsert(t *testing.T) {
	r := setupBudgetHandlerTest(t)
	body := []map[string]interface{}{
		{"account_code": "6421", "period_year": 2026, "period_month": 1, "budgeted": 50000000},
		{"account_code": "6422", "period_year": 2026, "period_month": 1, "budgeted": 30000000},
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/budgets/upsert?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]int
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 2, resp["upserted"])
}

func TestBudgetVarianceReport(t *testing.T) {
	r := setupBudgetHandlerTest(t)
	body := map[string]interface{}{
		"account_code": "6421",
		"period_year":  2026,
		"period_month": 1,
		"budgeted":     50000000,
		"actual":       45000000,
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/budgets?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/v1/budgets/report?company_id=comp1&year=2026&month=1", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	var report service.BudgetVarianceReport
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &report))
	assert.Len(t, report.Items, 1)
	assert.Equal(t, -5000000.0, report.TotalVariance)
}
