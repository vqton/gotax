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

type ccdcTestSetup struct {
	r   *gin.Engine
	svc service.CCDCServiceInterface
	compID string
}

func setupCCDCTest(t *testing.T) *ccdcTestSetup {
	t.Helper()
	gin.SetMode(gin.TestMode)

	catRepo := repository.NewMemoryCCDCCategoryRepo()
	itemRepo := repository.NewMemoryCCDCItemRepo()
	svc := service.NewCCDCService(catRepo, itemRepo)
	h := NewCCDCHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterCCDCRoutes(r, h, noopMW)

	return &ccdcTestSetup{r: r, svc: svc, compID: "CMP001"}
}

// ─── Categories ──────────────────────────────────────────────────────

func TestCCDCCreateCategory(t *testing.T) {
	ts := setupCCDCTest(t)
	body := `{"code":"DCGN","name":"Dụng cụ giản nhip"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/ccdc/categories?company_id="+ts.compID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp domain.ToolEquipmentCategory
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "DCGN", resp.Code)
	assert.NotEmpty(t, resp.ID)
}

func TestCCDCCreateCategory_Empty(t *testing.T) {
	ts := setupCCDCTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/ccdc/categories?company_id="+ts.compID, strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCCDCListCategories(t *testing.T) {
	ts := setupCCDCTest(t)
	// create one first
	ts.r.ServeHTTP(httptest.NewRecorder(), func() *http.Request {
		req, _ := http.NewRequest("POST", "/api/v1/ccdc/categories?company_id="+ts.compID, strings.NewReader(`{"code":"X","name":"X"}`))
		req.Header.Set("Content-Type", "application/json")
		return req
	}())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ccdc/categories?company_id="+ts.compID, nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp []domain.ToolEquipmentCategory
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp, 1)
}

func TestCCDCGetCategory(t *testing.T) {
	ts := setupCCDCTest(t)
	// create
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/ccdc/categories?company_id="+ts.compID, strings.NewReader(`{"code":"Y","name":"Y"}`))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	var created domain.ToolEquipmentCategory
	json.Unmarshal(w.Body.Bytes(), &created)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/v1/ccdc/categories/"+created.ID, nil)
	ts.r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestCCDCDeleteCategory(t *testing.T) {
	ts := setupCCDCTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/ccdc/categories?company_id="+ts.compID, strings.NewReader(`{"code":"Z","name":"Z"}`))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	var created domain.ToolEquipmentCategory
	json.Unmarshal(w.Body.Bytes(), &created)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("DELETE", "/api/v1/ccdc/categories/"+created.ID, nil)
	ts.r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

// ─── Items ───────────────────────────────────────────────────────────

func TestCCDCCreateItem(t *testing.T) {
	ts := setupCCDCTest(t)
	body := `{"code":"CCDC-001","name":"Máy khoan","purchase_date":"2025-01-15","purchase_cost":5000000}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/ccdc?company_id="+ts.compID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp domain.ToolEquipment
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "CCDC-001", resp.Code)
	assert.Equal(t, domain.CCDCActive, resp.Status)
	assert.Equal(t, 5000000.0, resp.CurrentCost)
}

func TestCCDCCreateItem_Duplicate(t *testing.T) {
	ts := setupCCDCTest(t)
	body := `{"code":"DUP","name":"X","purchase_cost":100}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/ccdc?company_id="+ts.compID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/v1/ccdc?company_id="+ts.compID, strings.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusConflict, w2.Code)
}

func TestCCDCListItem(t *testing.T) {
	ts := setupCCDCTest(t)
	body := `{"code":"L1","name":"A","purchase_cost":100}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/ccdc?company_id="+ts.compID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/v1/ccdc?company_id="+ts.compID, nil)
	ts.r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	var resp []domain.ToolEquipment
	json.Unmarshal(w2.Body.Bytes(), &resp)
	assert.Len(t, resp, 1)
}

func TestCCDCGetItem(t *testing.T) {
	ts := setupCCDCTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/ccdc?company_id="+ts.compID, strings.NewReader(`{"code":"G1","name":"B","purchase_cost":200}`))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	var created domain.ToolEquipment
	json.Unmarshal(w.Body.Bytes(), &created)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/v1/ccdc/"+created.ID, nil)
	ts.r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestCCDCUpdateItem(t *testing.T) {
	ts := setupCCDCTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/ccdc?company_id="+ts.compID, strings.NewReader(`{"code":"U1","name":"C","purchase_cost":300}`))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	var created domain.ToolEquipment
	json.Unmarshal(w.Body.Bytes(), &created)

	body := `{"code":"U1","name":"C Updated","status":"DISPOSED"}`
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("PUT", "/api/v1/ccdc/"+created.ID, strings.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestCCDCDeleteItem(t *testing.T) {
	ts := setupCCDCTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/ccdc?company_id="+ts.compID, strings.NewReader(`{"code":"D1","name":"D"}`))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	var created domain.ToolEquipment
	json.Unmarshal(w.Body.Bytes(), &created)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("DELETE", "/api/v1/ccdc/"+created.ID, nil)
	ts.r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}
