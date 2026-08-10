package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/repository"
	"gotax/internal/service"
)

func setupFiscalYearHandlerTest(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	perRepo := repository.NewMemoryPeriodRepo()
	optRepo := repository.NewMemorySystemOptionRepo()
	svc := service.NewFiscalYearService(perRepo, optRepo)
	h := NewFiscalYearHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("role", "admin")
		c.Next()
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterFiscalYearRoutes(r, h, noopMW)
	return r
}

func TestFiscalYearCreateYear(t *testing.T) {
	r := setupFiscalYearHandlerTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/fiscal-years/2026/periods?company_id=comp1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, float64(12), resp["count"])
}

func TestFiscalYearListPeriods(t *testing.T) {
	r := setupFiscalYearHandlerTest(t)

	// Create year first
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/fiscal-years/2026/periods?company_id=comp1", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	// List periods
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/fiscal-years/periods?company_id=comp1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var periods []interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &periods))
	assert.Equal(t, 12, len(periods))
}
