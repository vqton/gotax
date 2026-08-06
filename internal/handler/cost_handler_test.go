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

func setupCostCenterHandlerTest(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemoryCostCenterRepo()
	svc := service.NewCostCenterService(repo)
	ch := NewCostCenterHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
		c.Next()
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterCostCenterRoutes(r, ch, noopMW)
	return r
}

func TestCostCenterCreate(t *testing.T) {
	r := setupCostCenterHandlerTest(t)
	body := map[string]interface{}{
		"code": "CC01",
		"name": "Phòng Kế toán",
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/cost-centers?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var resp domain.CostCenter
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "CC01", resp.Code)
	assert.Equal(t, "Phòng Kế toán", resp.Name)
	assert.Equal(t, "comp1", resp.CompanyID)
	assert.True(t, resp.IsActive)
	assert.NotEmpty(t, resp.ID)
}

func TestCostCenterCreate_MissingCompany(t *testing.T) {
	r := setupCostCenterHandlerTest(t)
	body := map[string]interface{}{
		"code": "CC01",
		"name": "Phòng Kế toán",
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/cost-centers", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCostCenterCreate_MissingCode(t *testing.T) {
	r := setupCostCenterHandlerTest(t)
	body := map[string]interface{}{
		"name": "Phòng Kế toán",
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/cost-centers?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCostCenterCreate_DuplicateCode(t *testing.T) {
	r := setupCostCenterHandlerTest(t)
	body := map[string]interface{}{
		"code": "CC01",
		"name": "Phòng Kế toán",
	}
	b, _ := json.Marshal(body)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/api/v1/cost-centers?company_id=comp1", bytes.NewReader(b))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/v1/cost-centers?company_id=comp1", bytes.NewReader(b))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusConflict, w2.Code)
}

func TestCostCenterGet(t *testing.T) {
	r := setupCostCenterHandlerTest(t)
	body := map[string]interface{}{
		"code": "CC01",
		"name": "Phòng Kế toán",
	}
	b, _ := json.Marshal(body)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/api/v1/cost-centers?company_id=comp1", bytes.NewReader(b))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	var created domain.CostCenter
	json.Unmarshal(w1.Body.Bytes(), &created)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/v1/cost-centers/"+created.ID, nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	var resp domain.CostCenter
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &resp))
	assert.Equal(t, "CC01", resp.Code)
}

func TestCostCenterGet_NotFound(t *testing.T) {
	r := setupCostCenterHandlerTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cost-centers/nonexistent", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCostCenterList(t *testing.T) {
	r := setupCostCenterHandlerTest(t)
	for i := 0; i < 3; i++ {
		body := map[string]interface{}{
			"code": "CC0" + string(rune('1'+i)),
			"name": "Phòng " + string(rune('A'+i)),
		}
		b, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/cost-centers?company_id=comp1", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cost-centers?company_id=comp1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var list []domain.CostCenter
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &list))
	assert.Len(t, list, 3)
}

func TestCostCenterUpdate(t *testing.T) {
	r := setupCostCenterHandlerTest(t)
	body := map[string]interface{}{
		"code": "CC01",
		"name": "Phòng Kế toán",
	}
	b, _ := json.Marshal(body)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/api/v1/cost-centers?company_id=comp1", bytes.NewReader(b))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	var created domain.CostCenter
	json.Unmarshal(w1.Body.Bytes(), &created)

	updateBody := map[string]interface{}{
		"code": "CC01",
		"name": "Phòng Kế toán - Update",
	}
	ub, _ := json.Marshal(updateBody)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("PUT", "/api/v1/cost-centers/"+created.ID+"?company_id=comp1", bytes.NewReader(ub))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	var resp domain.CostCenter
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &resp))
	assert.Equal(t, "Phòng Kế toán - Update", resp.Name)
}

func TestCostCenterDelete(t *testing.T) {
	r := setupCostCenterHandlerTest(t)
	body := map[string]interface{}{
		"code": "CC01",
		"name": "Phòng Kế toán",
	}
	b, _ := json.Marshal(body)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/api/v1/cost-centers?company_id=comp1", bytes.NewReader(b))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	var created domain.CostCenter
	json.Unmarshal(w1.Body.Bytes(), &created)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("DELETE", "/api/v1/cost-centers/"+created.ID, nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusNoContent, w2.Code)

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/api/v1/cost-centers/"+created.ID, nil)
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusNotFound, w3.Code)
}

func TestCostCenterListHierarchy(t *testing.T) {
	r := setupCostCenterHandlerTest(t)
	body := map[string]interface{}{
		"code": "CC01",
		"name": "Phòng Kế toán",
	}
	b, _ := json.Marshal(body)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/api/v1/cost-centers?company_id=comp1", bytes.NewReader(b))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/v1/cost-centers/hierarchy?company_id=comp1", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	var list []domain.CostCenter
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &list))
	assert.Len(t, list, 1)
}
