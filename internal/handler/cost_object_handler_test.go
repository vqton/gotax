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

func setupCostObjectHandlerTest(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemoryCostObjectRepo()
	svc := service.NewCostObjectService(repo)
	h := NewCostObjectHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
		c.Next()
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterCostObjectRoutes(r, h, noopMW)
	return r
}

func costObjectBody() map[string]interface{} {
	return map[string]interface{}{
		"code":            "SP001",
		"name":            "Product A",
		"type":            "PRODUCT",
		"costing_method":  "SIMPLE",
	}
}

func TestCostObjectCreate(t *testing.T) {
	r := setupCostObjectHandlerTest(t)
	b, _ := json.Marshal(costObjectBody())
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/cost-objects?company_id=COMP1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var resp domain.CostObject
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "SP001", resp.Code)
	assert.Equal(t, "Product A", resp.Name)
	assert.Equal(t, domain.CostObjectType("PRODUCT"), resp.Type)
	assert.Equal(t, domain.CostingMethod("SIMPLE"), resp.CostingMethod)
	assert.Equal(t, "COMP1", resp.CompanyID)
	assert.True(t, resp.IsActive)
	assert.NotEmpty(t, resp.ID)
}

func TestCostObjectCreate_MissingFields(t *testing.T) {
	r := setupCostObjectHandlerTest(t)
	body := map[string]interface{}{}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/cost-objects?company_id=COMP1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestCostObjectCreate_DuplicateCode(t *testing.T) {
	r := setupCostObjectHandlerTest(t)
	b, _ := json.Marshal(costObjectBody())

	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/api/v1/cost-objects?company_id=COMP1", bytes.NewReader(b))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/v1/cost-objects?company_id=COMP1", bytes.NewReader(b))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusConflict, w2.Code)
}

func TestCostObjectGet(t *testing.T) {
	r := setupCostObjectHandlerTest(t)
	b, _ := json.Marshal(costObjectBody())
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/api/v1/cost-objects?company_id=COMP1", bytes.NewReader(b))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	var created domain.CostObject
	json.Unmarshal(w1.Body.Bytes(), &created)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/v1/cost-objects/"+created.ID, nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	var resp domain.CostObject
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &resp))
	assert.Equal(t, "SP001", resp.Code)
	assert.Equal(t, "Product A", resp.Name)
}

func TestCostObjectGet_NotFound(t *testing.T) {
	r := setupCostObjectHandlerTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cost-objects/nonexistent", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCostObjectList(t *testing.T) {
	r := setupCostObjectHandlerTest(t)
	for i := 0; i < 3; i++ {
		body := map[string]interface{}{
			"code":           "SP0" + string(rune('1'+i)),
			"name":           "Product " + string(rune('A'+i)),
			"type":           "PRODUCT",
			"costing_method": "SIMPLE",
		}
		b, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/cost-objects?company_id=COMP1", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cost-objects?company_id=COMP1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var list []domain.CostObject
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &list))
	assert.Len(t, list, 3)
}

func TestCostObjectUpdate(t *testing.T) {
	r := setupCostObjectHandlerTest(t)
	b, _ := json.Marshal(costObjectBody())
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/api/v1/cost-objects?company_id=COMP1", bytes.NewReader(b))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	var created domain.CostObject
	json.Unmarshal(w1.Body.Bytes(), &created)

	updateBody := map[string]interface{}{
		"code":           "SP001",
		"name":           "Product A - Updated",
		"type":           "PRODUCT",
		"costing_method": "SIMPLE",
	}
	ub, _ := json.Marshal(updateBody)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("PUT", "/api/v1/cost-objects/"+created.ID+"?company_id=COMP1", bytes.NewReader(ub))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	var resp domain.CostObject
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &resp))
	assert.Equal(t, "Product A - Updated", resp.Name)
}

func TestCostObjectDelete(t *testing.T) {
	r := setupCostObjectHandlerTest(t)
	b, _ := json.Marshal(costObjectBody())
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/api/v1/cost-objects?company_id=COMP1", bytes.NewReader(b))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	var created domain.CostObject
	json.Unmarshal(w1.Body.Bytes(), &created)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("DELETE", "/api/v1/cost-objects/"+created.ID, nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusNoContent, w2.Code)

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/api/v1/cost-objects/"+created.ID, nil)
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusNotFound, w3.Code)
}
