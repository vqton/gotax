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

func setupContractHandlerTest(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	contractRepo := repository.NewMemoryContractRepo()
	paymentRepo := repository.NewMemoryContractPaymentRepo()
	svc := service.NewContractService(contractRepo, paymentRepo)
	h := NewContractHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("role", "admin")
		c.Next()
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterContractRoutes(r, h, noopMW)
	return r
}

func TestContractCreate(t *testing.T) {
	r := setupContractHandlerTest(t)
	body := map[string]interface{}{
		"code":              "HD001",
		"name":              "Hợp đồng bán hàng",
		"contract_type":     "SALES",
		"value":             100000000,
		"counterparty_name": "Công ty XYZ",
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/contracts?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var resp domain.Contract
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "HD001", resp.Code)
	assert.Equal(t, "SALES", resp.ContractType)
	assert.Equal(t, "draft", resp.Status)
}

func TestContractList(t *testing.T) {
	r := setupContractHandlerTest(t)

	// Create
	body := map[string]interface{}{"code": "HD001", "name": "Test"}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/contracts?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// List
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/contracts?company_id=comp1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp []domain.Contract
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp, 1)
}

func TestContractAddPayment(t *testing.T) {
	r := setupContractHandlerTest(t)

	// Create contract
	body := map[string]interface{}{"code": "HD001", "name": "Test", "value": 1000000}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/contracts?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	var contract domain.Contract
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &contract))

	// Add payment
	payBody := map[string]interface{}{
		"payment_date": "2026-08-01",
		"amount":       500000,
		"description":  "Đợt 1",
	}
	pb, _ := json.Marshal(payBody)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/contracts/"+contract.ID+"/payments", bytes.NewReader(pb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// List payments
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/contracts/"+contract.ID+"/payments", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var payments []domain.ContractPayment
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &payments))
	assert.Len(t, payments, 1)
	assert.Equal(t, float64(500000), payments[0].Amount)
}
