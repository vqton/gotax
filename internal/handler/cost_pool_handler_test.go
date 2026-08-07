package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
	"gotax/internal/repository"
	"gotax/internal/service"
)

func setupCostPoolHandlerTest(t *testing.T) (*gin.Engine, *domain.CostPool) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	poolRepo := repository.NewMemoryCostPoolRepo()
	lineRepo := repository.NewMemoryCostPoolLineRepo()
	svc := service.NewCostPoolService(poolRepo, lineRepo)
	h := NewCostPoolHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
		c.Next()
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterCostPoolRoutes(r, h, noopMW)

	body := map[string]interface{}{
		"period_id":       "PER001",
		"gl_account_code": "621",
		"name":            "Direct Materials Pool",
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/cost-pools?company_id=COMP1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
	var created domain.CostPool
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &created))
	return r, &created
}

func TestCostPoolCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	poolRepo := repository.NewMemoryCostPoolRepo()
	lineRepo := repository.NewMemoryCostPoolLineRepo()
	svc := service.NewCostPoolService(poolRepo, lineRepo)
	h := NewCostPoolHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
		c.Next()
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterCostPoolRoutes(r, h, noopMW)

	body := map[string]interface{}{
		"period_id":       "PER001",
		"gl_account_code": "621",
		"name":            "Direct Materials Pool",
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/cost-pools?company_id=COMP1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var resp domain.CostPool
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "PER001", resp.PeriodID)
	assert.Equal(t, "621", resp.GLAccountCode)
	assert.Equal(t, "Direct Materials Pool", resp.Name)
	assert.Equal(t, "COMP1", resp.CompanyID)
	assert.Equal(t, domain.CostPoolStatusOpen, resp.Status)
	assert.NotEmpty(t, resp.ID)
}

func TestCostPoolCreate_MissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	poolRepo := repository.NewMemoryCostPoolRepo()
	lineRepo := repository.NewMemoryCostPoolLineRepo()
	svc := service.NewCostPoolService(poolRepo, lineRepo)
	h := NewCostPoolHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
		c.Next()
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterCostPoolRoutes(r, h, noopMW)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/cost-pools?company_id=COMP1", strings.NewReader(`"invalid"`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCostPoolGet(t *testing.T) {
	r, created := setupCostPoolHandlerTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cost-pools/"+created.ID, nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp domain.CostPool
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "Direct Materials Pool", resp.Name)
	assert.Equal(t, "621", resp.GLAccountCode)
}

func TestCostPoolGet_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	poolRepo := repository.NewMemoryCostPoolRepo()
	lineRepo := repository.NewMemoryCostPoolLineRepo()
	svc := service.NewCostPoolService(poolRepo, lineRepo)
	h := NewCostPoolHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
		c.Next()
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterCostPoolRoutes(r, h, noopMW)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cost-pools/nonexistent", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCostPoolList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	poolRepo := repository.NewMemoryCostPoolRepo()
	lineRepo := repository.NewMemoryCostPoolLineRepo()
	svc := service.NewCostPoolService(poolRepo, lineRepo)
	h := NewCostPoolHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
		c.Next()
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterCostPoolRoutes(r, h, noopMW)

	for i := 0; i < 3; i++ {
		body := map[string]interface{}{
			"period_id":       "PER001",
			"gl_account_code": "621",
			"name":            fmt.Sprintf("Pool %d", i+1),
		}
		b, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/cost-pools?company_id=COMP1", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cost-pools?company_id=COMP1&period_id=PER001", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var list []domain.CostPool
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &list))
	assert.Len(t, list, 3)
}

func TestCostPoolDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	poolRepo := repository.NewMemoryCostPoolRepo()
	lineRepo := repository.NewMemoryCostPoolLineRepo()
	svc := service.NewCostPoolService(poolRepo, lineRepo)
	h := NewCostPoolHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
		c.Next()
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterCostPoolRoutes(r, h, noopMW)

	body := map[string]interface{}{
		"period_id":       "PER001",
		"gl_account_code": "621",
		"name":            "Direct Materials Pool",
	}
	b, _ := json.Marshal(body)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/api/v1/cost-pools?company_id=COMP1", bytes.NewReader(b))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	var created domain.CostPool
	json.Unmarshal(w1.Body.Bytes(), &created)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("DELETE", "/api/v1/cost-pools/"+created.ID, nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusNoContent, w2.Code)

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/api/v1/cost-pools/"+created.ID, nil)
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusNotFound, w3.Code)
}

func TestCostPoolAddLine(t *testing.T) {
	r, created := setupCostPoolHandlerTest(t)

	lineBody := map[string]interface{}{
		"source_type":  "JOURNAL",
		"description":  "Material from GRN-001",
		"amount":       5000000,
	}
	b, _ := json.Marshal(lineBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/cost-pools/"+created.ID+"/lines", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var resp domain.CostPoolLine
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "JOURNAL", resp.SourceType)
	assert.Equal(t, "Material from GRN-001", resp.Description)
	assert.Equal(t, float64(5000000), resp.Amount)
	assert.NotEmpty(t, resp.ID)
}

func TestCostPoolAddLine_AmountZero(t *testing.T) {
	r, created := setupCostPoolHandlerTest(t)

	lineBody := map[string]interface{}{
		"source_type":  "JOURNAL",
		"description":  "Zero amount line",
		"amount":       0,
	}
	b, _ := json.Marshal(lineBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/cost-pools/"+created.ID+"/lines", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCostPoolListLines(t *testing.T) {
	r, created := setupCostPoolHandlerTest(t)

	for i := 0; i < 2; i++ {
		lineBody := map[string]interface{}{
			"source_type":  "JOURNAL",
			"description":  fmt.Sprintf("Line %d", i+1),
			"amount":       1000000 * (i + 1),
		}
		b, _ := json.Marshal(lineBody)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/cost-pools/"+created.ID+"/lines", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cost-pools/"+created.ID+"/lines", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var lines []domain.CostPoolLine
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &lines))
	assert.Len(t, lines, 2)
}

func TestCostPoolClose(t *testing.T) {
	r, created := setupCostPoolHandlerTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/cost-pools/"+created.ID+"/close", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "closed", resp["status"])

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/v1/cost-pools/"+created.ID, nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	var pool domain.CostPool
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &pool))
	assert.Equal(t, domain.CostPoolStatusClosed, pool.Status)
}
