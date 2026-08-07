package handler

import (
	"context"
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

func setupSalaryGradeHandlerTest(t *testing.T) (*gin.Engine, *service.PayrollService) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemoryPayrollRepo()
	svc := service.NewPayrollService(repo, nil)
	ph := NewPayrollHandler(svc)

	r := gin.New()
	authMW := func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
		c.Next()
	}

	RegisterPayrollRoutes(r, ph, authMW)
	return r, svc
}

// ─── Salary Grade Tests ─────────────────────────────────────────

func TestCreateSalaryGrade_Success(t *testing.T) {
	r, _ := setupSalaryGradeHandlerTest(t)
	body := `{"company_id":"CMP001","code":"G1","name":"Nhân viên","min_salary":5000000,"max_salary":10000000,"level":1}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/salary-grades", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateSalaryGrade_DuplicateCode(t *testing.T) {
	r, _ := setupSalaryGradeHandlerTest(t)
	body := `{"company_id":"CMP001","code":"G1","name":"Grade 1","min_salary":5000000,"max_salary":10000000,"level":1}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/salary-grades", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Duplicate
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/v1/payroll/salary-grades", strings.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusConflict, w2.Code)
}

func TestListSalaryGrades(t *testing.T) {
	r, _ := setupSalaryGradeHandlerTest(t)
	// Create two grades
	body := `{"company_id":"CMP001","code":"G1","name":"Grade 1","min_salary":5000000,"max_salary":10000000,"level":1}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/salary-grades", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	body2 := `{"company_id":"CMP001","code":"G2","name":"Grade 2","min_salary":10000000,"max_salary":20000000,"level":2}`
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/v1/payroll/salary-grades", strings.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)

	// List
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/api/v1/payroll/salary-grades?company_id=CMP001", nil)
	r.ServeHTTP(w3, req3)

	assert.Equal(t, http.StatusOK, w3.Code)
	assert.Contains(t, w3.Body.String(), "G1")
	assert.Contains(t, w3.Body.String(), "G2")
}

func TestGetSalaryGrade(t *testing.T) {
	r, svc := setupSalaryGradeHandlerTest(t)
	grade := &domain.SalaryGrade{
		CompanyID: "CMP001", Code: "G1", Name: "Grade 1",
		MinSalary: 5_000_000, MaxSalary: 10_000_000, Level: 1, IsActive: true,
	}
	require.NoError(t, svc.CreateSalaryGrade(context.Background(), grade))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/salary-grades/"+grade.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "G1")
}

func TestUpdateSalaryGrade(t *testing.T) {
	r, svc := setupSalaryGradeHandlerTest(t)
	grade := &domain.SalaryGrade{
		CompanyID: "CMP001", Code: "G1", Name: "Grade 1",
		MinSalary: 5_000_000, MaxSalary: 10_000_000, Level: 1, IsActive: true,
	}
	require.NoError(t, svc.CreateSalaryGrade(context.Background(), grade))

	body := `{"company_id":"CMP001","code":"G1","name":"Grade 1 Updated","min_salary":6000000,"max_salary":12000000,"level":1}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/payroll/salary-grades/"+grade.ID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteSalaryGrade(t *testing.T) {
	r, svc := setupSalaryGradeHandlerTest(t)
	grade := &domain.SalaryGrade{
		CompanyID: "CMP001", Code: "G1", Name: "Grade 1",
		MinSalary: 5_000_000, MaxSalary: 10_000_000, Level: 1, IsActive: true,
	}
	require.NoError(t, svc.CreateSalaryGrade(context.Background(), grade))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/payroll/salary-grades/"+grade.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ─── Salary Scale Tests ─────────────────────────────────────────

func TestCreateSalaryScale_Success(t *testing.T) {
	r, svc := setupSalaryGradeHandlerTest(t)
	grade := &domain.SalaryGrade{
		CompanyID: "CMP001", Code: "G1", Name: "Grade 1",
		MinSalary: 5_000_000, MaxSalary: 10_000_000, Level: 1, IsActive: true,
	}
	require.NoError(t, svc.CreateSalaryGrade(context.Background(), grade))

	body := `{"company_id":"CMP001","grade_id":"` + grade.ID + `","code":"G1.1","name":"Bậc 1","base_salary":5000000,"level":1,"min_years":0}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/salary-scales", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateSalaryScale_InvalidGrade(t *testing.T) {
	r, _ := setupSalaryGradeHandlerTest(t)
	body := `{"company_id":"CMP001","grade_id":"nonexistent","code":"G1.1","name":"Bậc 1","base_salary":5000000,"level":1}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/salary-scales", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestListSalaryScalesByGrade(t *testing.T) {
	r, svc := setupSalaryGradeHandlerTest(t)
	grade := &domain.SalaryGrade{
		CompanyID: "CMP001", Code: "G1", Name: "Grade 1",
		MinSalary: 5_000_000, MaxSalary: 10_000_000, Level: 1, IsActive: true,
	}
	require.NoError(t, svc.CreateSalaryGrade(context.Background(), grade))

	// Create 2 scales
	scale1 := &domain.SalaryScale{
		CompanyID: "CMP001", GradeID: grade.ID, Code: "G1.1", Name: "Bậc 1",
		BaseSalary: 5_000_000, Level: 1, MinYears: 0, IsActive: true,
	}
	scale2 := &domain.SalaryScale{
		CompanyID: "CMP001", GradeID: grade.ID, Code: "G1.2", Name: "Bậc 2",
		BaseSalary: 7_000_000, Level: 2, MinYears: 2, IsActive: true,
	}
	require.NoError(t, svc.CreateSalaryScale(context.Background(), scale1))
	require.NoError(t, svc.CreateSalaryScale(context.Background(), scale2))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/salary-grades/"+grade.ID+"/scales", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "G1.1")
	assert.Contains(t, w.Body.String(), "G1.2")
}

// ─── Employee Salary Grade Tests ────────────────────────────────

func TestAssignEmployeeSalaryGrade(t *testing.T) {
	r, svc := setupSalaryGradeHandlerTest(t)
	grade := &domain.SalaryGrade{
		CompanyID: "CMP001", Code: "G1", Name: "Grade 1",
		MinSalary: 5_000_000, MaxSalary: 10_000_000, Level: 1, IsActive: true,
	}
	require.NoError(t, svc.CreateSalaryGrade(context.Background(), grade))

	scale := &domain.SalaryScale{
		CompanyID: "CMP001", GradeID: grade.ID, Code: "G1.1", Name: "Bậc 1",
		BaseSalary: 5_000_000, Level: 1, MinYears: 0, IsActive: true,
	}
	require.NoError(t, svc.CreateSalaryScale(context.Background(), scale))

	body := `{"employee_id":"NV001","company_id":"CMP001","grade_id":"` + grade.ID + `","scale_id":"` + scale.ID + `","base_salary":7500000,"effective_from":"2026-01-01"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/employee-salary-grades", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetEmployeeSalaryGrade(t *testing.T) {
	r, svc := setupSalaryGradeHandlerTest(t)
	grade := &domain.SalaryGrade{
		CompanyID: "CMP001", Code: "G1", Name: "Grade 1",
		MinSalary: 5_000_000, MaxSalary: 10_000_000, Level: 1, IsActive: true,
	}
	require.NoError(t, svc.CreateSalaryGrade(context.Background(), grade))

	scale := &domain.SalaryScale{
		CompanyID: "CMP001", GradeID: grade.ID, Code: "G1.1", Name: "Bậc 1",
		BaseSalary: 5_000_000, Level: 1, MinYears: 0, IsActive: true,
	}
	require.NoError(t, svc.CreateSalaryScale(context.Background(), scale))

	assign := &domain.EmployeeSalaryGrade{
		EmployeeID: "NV001", CompanyID: "CMP001", GradeID: grade.ID,
		ScaleID: scale.ID, BaseSalary: 7_500_000, EffectiveFrom: "2026-01-01",
	}
	require.NoError(t, svc.AssignEmployeeSalaryGrade(context.Background(), assign))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/employee-salary-grades/NV001", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "NV001")
	assert.Contains(t, w.Body.String(), "7500000")
}
