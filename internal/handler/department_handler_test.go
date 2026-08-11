package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gotax/internal/domain"
	"gotax/internal/repository"
	"gotax/internal/service"
)

func setupDepartmentTest(t *testing.T) (*gin.Engine, service.CompanyService, *domain.Company) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemoryCompanyRepo()
	svc := service.NewCompanyService(repo)
	dh := NewDepartmentHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterDepartmentRoutes(r, dh, noopMW)

	comp := &domain.Company{
		TenantID:        "tenant-1",
		TaxCode:         "1234567890",
		LegalNameVN:     "Test Co",
		LegalForm:       domain.LegalFormJSC,
		AccountingRegime: domain.RegimeTT99,
	}
	svc.CreateCompany(nil, comp)

	return r, svc, comp
}

func TestDepartmentCreate(t *testing.T) {
	r, _, comp := setupDepartmentTest(t)
	body := `{"code":"D001","name":"Engineering"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/departments?company_id="+comp.ID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	var result domain.Department
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "D001", result.Code)
	assert.Equal(t, comp.ID, result.CompanyID)
}

func TestDepartmentCreate_MissingCompanyID(t *testing.T) {
	r, _, _ := setupDepartmentTest(t)
	body := `{"code":"D001","name":"Engineering"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/departments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestDepartmentCreate_BadRequest(t *testing.T) {
	r, _, comp := setupDepartmentTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/departments?company_id="+comp.ID, strings.NewReader(`not-json`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestDepartmentList(t *testing.T) {
	r, svc, comp := setupDepartmentTest(t)
	svc.CreateDepartment(nil, &domain.Department{CompanyID: comp.ID, Code: "D001", Name: "Eng"})
	svc.CreateDepartment(nil, &domain.Department{CompanyID: comp.ID, Code: "D002", Name: "HR"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/departments?company_id="+comp.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var result []domain.Department
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Len(t, result, 2)
}

func TestDepartmentList_MissingCompanyID(t *testing.T) {
	r, _, _ := setupDepartmentTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/departments", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestDepartmentGet(t *testing.T) {
	r, svc, comp := setupDepartmentTest(t)
	dept := &domain.Department{CompanyID: comp.ID, Code: "D001", Name: "Eng"}
	svc.CreateDepartment(nil, dept)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/departments/"+dept.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var result domain.Department
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "D001", result.Code)
}

func TestDepartmentGet_NotFound(t *testing.T) {
	r, _, _ := setupDepartmentTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/departments/nonexistent", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestDepartmentUpdate(t *testing.T) {
	r, svc, comp := setupDepartmentTest(t)
	dept := &domain.Department{CompanyID: comp.ID, Code: "D001", Name: "Eng"}
	svc.CreateDepartment(nil, dept)

	body := `{"code":"D001","name":"Engineering Updated"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/departments/"+dept.ID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var result domain.Department
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "Engineering Updated", result.Name)
}

func TestDepartmentUpdate_BadRequest(t *testing.T) {
	r, svc, comp := setupDepartmentTest(t)
	dept := &domain.Department{CompanyID: comp.ID, Code: "D001", Name: "Eng"}
	svc.CreateDepartment(nil, dept)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/departments/"+dept.ID, strings.NewReader(`not-json`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestDepartmentUpdate_EmptyCode(t *testing.T) {
	r, svc, comp := setupDepartmentTest(t)
	dept := &domain.Department{CompanyID: comp.ID, Code: "D001", Name: "Eng"}
	svc.CreateDepartment(nil, dept)

	body := `{"code":"","name":"Engineering"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/departments/"+dept.ID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestDepartmentUpdate_EmptyName(t *testing.T) {
	r, svc, comp := setupDepartmentTest(t)
	dept := &domain.Department{CompanyID: comp.ID, Code: "D001", Name: "Eng"}
	svc.CreateDepartment(nil, dept)

	body := `{"code":"D001","name":""}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/departments/"+dept.ID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestDepartmentDeactivate(t *testing.T) {
	r, svc, comp := setupDepartmentTest(t)
	dept := &domain.Department{CompanyID: comp.ID, Code: "D001", Name: "Eng"}
	svc.CreateDepartment(nil, dept)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/departments/"+dept.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestDepartmentDeactivate_NotFound(t *testing.T) {
	r, _, _ := setupDepartmentTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/departments/nonexistent", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}
