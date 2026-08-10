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

func setupSystemOptionHandlerTest(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemorySystemOptionRepo()
	svc := service.NewSystemOptionService(repo)
	h := NewSystemOptionHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
		c.Next()
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterSystemOptionRoutes(r, h, noopMW)
	return r
}

func TestSystemOptionUpsert(t *testing.T) {
	r := setupSystemOptionHandlerTest(t)
	body := map[string]interface{}{
		"category": "global",
		"key":      "accounting_standard",
		"value":    "circular99",
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/system-options?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp domain.SystemOption
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "global", resp.Category)
	assert.Equal(t, "accounting_standard", resp.Key)
	assert.Equal(t, "circular99", resp.Value)
	assert.Equal(t, "comp1", resp.CompanyID)
}

func TestSystemOptionGetByCategory(t *testing.T) {
	r := setupSystemOptionHandlerTest(t)

	// Upsert two options
	for _, tc := range []map[string]interface{}{
		{"category": "global", "key": "fiscal_year_start", "value": "1"},
		{"category": "global", "key": "base_currency", "value": "VND"},
	} {
		b, _ := json.Marshal(tc)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/system-options?company_id=comp1", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	}

	// Get by category
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/system-options/global?company_id=comp1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp []domain.SystemOption
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.GreaterOrEqual(t, len(resp), 2)
}

func TestSystemOptionUpsert_MissingCompany(t *testing.T) {
	r := setupSystemOptionHandlerTest(t)
	body := map[string]interface{}{
		"category": "global",
		"key":      "test",
		"value":    "val",
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/system-options", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNumberingRuleCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	numRepo := repository.NewMemoryNumberingRuleRepo()
	svc := service.NewNumberingRuleService(numRepo)
	h := NewNumberingRuleHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("role", "admin")
		c.Next()
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterNumberingRuleRoutes(r, h, noopMW)

	body := map[string]interface{}{
		"voucher_type":  "BH",
		"prefix":        "BH",
		"number_length": 5,
		"scope":         "company",
		"reset_rule":    "never",
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/numbering-rules?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var resp domain.NumberingRule
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "BH", resp.VoucherType)
	assert.Equal(t, 5, resp.NumberLength)
	assert.Equal(t, 0, resp.CurrentNum)
}

func TestNumberingRuleGetNext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	numRepo := repository.NewMemoryNumberingRuleRepo()
	svc := service.NewNumberingRuleService(numRepo)
	h := NewNumberingRuleHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("role", "admin")
		c.Next()
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterNumberingRuleRoutes(r, h, noopMW)

	// Create rule
	body := map[string]interface{}{
		"voucher_type":  "MB",
		"prefix":        "MB",
		"number_length": 5,
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/numbering-rules?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	// Get next number three times
	for i := 1; i <= 3; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/numbering-rules/next/MB?company_id=comp1", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, float64(i), resp["current_num"])
	}
}
