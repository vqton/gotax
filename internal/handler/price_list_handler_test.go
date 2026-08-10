package handler

import (
	"encoding/json"
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

func setupPriceListTest(t *testing.T) (*gin.Engine, *service.PriceListService) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemoryPriceListRepo()
	svc := service.NewPriceListService(repo)
	ph := NewPriceListHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterPriceListRoutes(r, ph, noopMW)

	return r, svc
}

func TestCreatePriceList(t *testing.T) {
	r, _ := setupPriceListTest(t)
	body := `{"code":"PL01","name":"Standard Prices","currency":"VND","is_active":true}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/price-lists?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var pl domain.PriceList
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &pl))
	assert.Equal(t, "PL01", pl.Code)
	assert.Equal(t, "Standard Prices", pl.Name)
}

func TestGetPriceList(t *testing.T) {
	r, svc := setupPriceListTest(t)
	pl := &domain.PriceList{CompanyID: "CMP001", Code: "PL01", Name: "Standard", Currency: "VND"}
	require.NoError(t, svc.CreatePriceList(nil, pl))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/price-lists/"+pl.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got domain.PriceList
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, "PL01", got.Code)
}

func TestListPriceLists(t *testing.T) {
	r, svc := setupPriceListTest(t)
	svc.CreatePriceList(nil, &domain.PriceList{CompanyID: "CMP001", Code: "PL01", Name: "P1", Currency: "VND"})
	svc.CreatePriceList(nil, &domain.PriceList{CompanyID: "CMP001", Code: "PL02", Name: "P2", Currency: "VND"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/price-lists?company_id=CMP001", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var lists []domain.PriceList
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &lists))
	assert.Len(t, lists, 2)
}

func TestAddLines(t *testing.T) {
	r, svc := setupPriceListTest(t)
	pl := &domain.PriceList{CompanyID: "CMP001", Code: "PL01", Name: "Standard", Currency: "VND"}
	require.NoError(t, svc.CreatePriceList(nil, pl))

	body := `[{"item_code":"ITM001","item_name":"Widget","unit":"pcs","unit_price":100000}]`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/price-lists/"+pl.ID+"/lines", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	lines, _ := svc.GetLines(nil, pl.ID)
	require.Len(t, lines, 1)
	assert.Equal(t, "ITM001", lines[0].ItemCode)
}

func TestCalculateSellingPrice(t *testing.T) {
	r, svc := setupPriceListTest(t)
	pl := &domain.PriceList{CompanyID: "CMP001", Code: "PL01", Name: "Standard", Currency: "VND"}
	require.NoError(t, svc.CreatePriceList(nil, pl))
	svc.AddLines(nil, pl.ID, []domain.PriceListLine{
		{ItemCode: "ITM001", ItemName: "Widget", Unit: "pcs", UnitPrice: 100000},
	})

	body := `{"price_list_id":"` + pl.ID + `","item_code":"ITM001","markup_pct":10}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/price-lists/calculate-price", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]float64
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.InDelta(t, 110000.0, resp["unit_price"], 0.01) // 100000 * 1.10
}
