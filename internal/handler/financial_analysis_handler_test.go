package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gotax/internal/repository"
	"gotax/internal/service"
)

func setupFinancialAnalysisTest(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	jeRepo := repository.NewMemoryJournalRepo()
	budRepo := repository.NewMemoryBudgetRepo()
	compRepo := repository.NewMemoryCompanyRepo()
	svc := service.NewFinancialAnalysisService(jeRepo, budRepo, compRepo)
	fh := NewFinancialAnalysisHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterFinancialAnalysisRoutes(r, fh, noopMW)

	return r
}

func TestGetFinancialRatios_MissingCompanyID(t *testing.T) {
	r := setupFinancialAnalysisTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/reports/financial-ratios", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetFinancialRatios_EmptyData(t *testing.T) {
	r := setupFinancialAnalysisTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/reports/financial-ratios?company_id=CMP001&year=2026", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var ratios []service.FinancialRatio
	json.Unmarshal(w.Body.Bytes(), &ratios)
	assert.Empty(t, ratios) // No data → no ratios
}

func TestGetBudgetVsActual_MissingCompanyID(t *testing.T) {
	r := setupFinancialAnalysisTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/reports/budget-vs-actual", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetBudgetVsActual_EmptyData(t *testing.T) {
	r := setupFinancialAnalysisTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/reports/budget-vs-actual?company_id=CMP001&year=2026", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var result service.BudgetVsActualResult
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "CMP001", result.CompanyID)
	assert.Empty(t, result.Items)
}
