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

func setupRecurringHandlerTest(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	memRepo := repository.NewMemoryRecurringEntryRepo()
	memJERepo := repository.NewMemoryJournalRepo()
	accRepo := repository.NewMemoryAccountRepo()
	memJERepo.SetAccounts(accRepo.Accounts())
	svc := service.NewRecurringService(memRepo, memJERepo)
	rh := NewRecurringHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
		c.Next()
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterRecurringRoutes(r, rh, noopMW)

	return r
}

func TestRecurringCreate(t *testing.T) {
	r := setupRecurringHandlerTest(t)
	body := map[string]interface{}{
		"template_name": "Monthly Rent",
		"description":   "Office rent",
		"frequency":     "MONTHLY",
		"day_of_month":  5,
		"is_active":     true,
		"lines": []map[string]interface{}{
			{"account_code": "6421", "debit_amount": 10000000, "credit_amount": 0, "description": "Rent expense"},
			{"account_code": "3311", "debit_amount": 0, "credit_amount": 10000000, "description": "Payable"},
		},
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/recurring-entries?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var resp domain.RecurringEntry
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "Monthly Rent", resp.TemplateName)
	assert.NotEmpty(t, resp.ID)
}

func TestRecurringCreate_ValidationError(t *testing.T) {
	r := setupRecurringHandlerTest(t)
	body := map[string]interface{}{
		"template_name": "",
		"frequency":     "MONTHLY",
		"day_of_month":  5,
		"lines":         []map[string]interface{}{},
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/recurring-entries?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRecurringList(t *testing.T) {
	r := setupRecurringHandlerTest(t)
	body := map[string]interface{}{
		"template_name": "Rent",
		"frequency":     "MONTHLY",
		"day_of_month":  1,
		"is_active":     true,
		"lines":         []map[string]interface{}{{"account_code": "6421", "debit_amount": 5000, "credit_amount": 0}},
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/recurring-entries?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/v1/recurring-entries?company_id=comp1", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	var list []domain.RecurringEntry
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &list))
	assert.Len(t, list, 1)
}

func TestRecurringGet(t *testing.T) {
	r := setupRecurringHandlerTest(t)
	body := map[string]interface{}{
		"template_name": "Insurance",
		"frequency":     "YEARLY",
		"day_of_month":  15,
		"is_active":     true,
		"lines":         []map[string]interface{}{{"account_code": "6422", "debit_amount": 50000, "credit_amount": 0}},
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/recurring-entries?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	var created domain.RecurringEntry
	json.Unmarshal(w.Body.Bytes(), &created)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/v1/recurring-entries/"+created.ID, nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestRecurringGet_NotFound(t *testing.T) {
	r := setupRecurringHandlerTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/recurring-entries/nonexistent", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRecurringUpdate(t *testing.T) {
	r := setupRecurringHandlerTest(t)
	body := map[string]interface{}{
		"template_name": "Utilities",
		"frequency":     "MONTHLY",
		"day_of_month":  10,
		"is_active":     true,
		"lines":         []map[string]interface{}{{"account_code": "6423", "debit_amount": 2000, "credit_amount": 0}},
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/recurring-entries?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	var created domain.RecurringEntry
	json.Unmarshal(w.Body.Bytes(), &created)

	body["template_name"] = "Electricity"
	b2, _ := json.Marshal(body)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("PUT", "/api/v1/recurring-entries/"+created.ID+"?company_id=comp1", bytes.NewReader(b2))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	var updated domain.RecurringEntry
	json.Unmarshal(w2.Body.Bytes(), &updated)
	assert.Equal(t, "Electricity", updated.TemplateName)
}

func TestRecurringDelete(t *testing.T) {
	r := setupRecurringHandlerTest(t)
	body := map[string]interface{}{
		"template_name": "To Delete",
		"frequency":     "MONTHLY",
		"day_of_month":  1,
		"is_active":     true,
		"lines":         []map[string]interface{}{{"account_code": "6424", "debit_amount": 1000, "credit_amount": 0}},
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/recurring-entries?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	var created domain.RecurringEntry
	json.Unmarshal(w.Body.Bytes(), &created)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("DELETE", "/api/v1/recurring-entries/"+created.ID, nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/api/v1/recurring-entries/"+created.ID, nil)
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusNotFound, w3.Code)
}

func TestRecurringRunNow(t *testing.T) {
	r := setupRecurringHandlerTest(t)
	body := map[string]interface{}{
		"template_name": "Depreciation",
		"frequency":     "MONTHLY",
		"day_of_month":  28,
		"is_active":     true,
		"lines": []map[string]interface{}{
			{"account_code": "6425", "debit_amount": 100000, "credit_amount": 0},
			{"account_code": "2131", "debit_amount": 0, "credit_amount": 100000},
		},
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/recurring-entries?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	var created domain.RecurringEntry
	json.Unmarshal(w.Body.Bytes(), &created)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/v1/recurring-entries/"+created.ID+"/run", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusCreated, w2.Code)
	var je domain.JournalEntry
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &je))
	assert.Equal(t, domain.JournalEntryDraft, je.Status)
	assert.Len(t, je.Lines, 2)
}

func TestRecurringRunNow_NotActive(t *testing.T) {
	r := setupRecurringHandlerTest(t)
	body := map[string]interface{}{
		"template_name": "Inactive",
		"frequency":     "MONTHLY",
		"day_of_month":  1,
		"is_active":     false,
		"lines":         []map[string]interface{}{{"account_code": "6426", "debit_amount": 1000, "credit_amount": 0}},
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/recurring-entries?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	var created domain.RecurringEntry
	json.Unmarshal(w.Body.Bytes(), &created)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/v1/recurring-entries/"+created.ID+"/run", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

func TestRecurringProcessDue(t *testing.T) {
	r := setupRecurringHandlerTest(t)
	body := map[string]interface{}{
		"template_name": "Due Entry",
		"frequency":     "MONTHLY",
		"day_of_month":  1,
		"is_active":     true,
		"next_run_date": "2020-01-01",
		"lines":         []map[string]interface{}{{"account_code": "6427", "debit_amount": 500, "credit_amount": 0}},
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/recurring-entries?company_id=comp1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/v1/recurring-entries/process-due", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	var resp map[string]int
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &resp))
	assert.GreaterOrEqual(t, resp["processed"], 1)
}
