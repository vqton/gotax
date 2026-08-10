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

	"gotax/internal/repository"
	"gotax/internal/service"
)

func setupReportOptionHandlerTest(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	optRepo := repository.NewMemorySystemOptionRepo()
	svc := service.NewReportOptionService(optRepo)
	h := NewReportOptionHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("role", "admin")
		c.Next()
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterReportOptionRoutes(r, h, noopMW)
	return r
}

func TestReportOptionUpdateAndGet(t *testing.T) {
	r := setupReportOptionHandlerTest(t)

	body := map[string]interface{}{
		"company_name":    "Công ty ABC",
		"company_address": "Hà Nội",
		"company_tax_code": "0123456789",
		"alignment":       "center",
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/report-options?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Get
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/report-options?company_id=comp1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp service.ReportOptions
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "Công ty ABC", resp.CompanyName)
	assert.Equal(t, "Hà Nội", resp.CompanyAddress)
	assert.Equal(t, "center", resp.Alignment)
}
