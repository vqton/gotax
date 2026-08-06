package handler

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
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

func setupPayrollHandlerTest(t *testing.T) (*gin.Engine, *service.PayrollService) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemoryPayrollRepo()
	svc := service.NewPayrollService(repo, nil)
	ph := NewPayrollHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
		c.Next()
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterPayrollRoutes(r, ph, noopMW)

	return r, svc
}

func testContext() *gin.Context {
	c, _ := gin.CreateTestContext(nil)
	return c
}

// ─── Employee Payroll Info ──────────────────────────────────────

func TestPayrollGetEmployeePayrollInfo_NotFound(t *testing.T) {
	r, _ := setupPayrollHandlerTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/employees/nonexistent", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPayrollUpdateEmployeePayrollInfo(t *testing.T) {
	r, _ := setupPayrollHandlerTest(t)

	body := `{"base_salary":15000000,"region":"II"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/payroll/employees/NV001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, w.Code)
}

// ─── Periods ────────────────────────────────────────────────────

func TestPayrollCreatePeriod_Success(t *testing.T) {
	r, _ := setupPayrollHandlerTest(t)
	body := `{"company_id":"CMP001","year":2026,"month":7}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/periods", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var period domain.PayrollPeriod
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &period))
	assert.Equal(t, 2026, period.Year)
	assert.Equal(t, 7, period.Month)
	assert.Equal(t, domain.PayrollDraft, period.Status)
}

func TestPayrollCreatePeriod_Duplicate(t *testing.T) {
	r, _ := setupPayrollHandlerTest(t)
	body := `{"company_id":"CMP001","year":2026,"month":7}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/periods", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/v1/payroll/periods", strings.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusConflict, w2.Code)
}

func TestPayrollListPeriods_Success(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	_, _ = svc.CreatePeriod(testContext(), "CMP001", 2026, 6)
	_, _ = svc.CreatePeriod(testContext(), "CMP001", 2026, 7)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/periods?company_id=CMP001", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var periods []domain.PayrollPeriod
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &periods))
	assert.Len(t, periods, 2)
}

func TestPayrollGetPeriod_Success(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	period, _ := svc.CreatePeriod(testContext(), "CMP001", 2026, 7)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/periods/"+period.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPayrollGetPeriod_NotFound(t *testing.T) {
	r, _ := setupPayrollHandlerTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/periods/nonexistent", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ─── Approve Period ─────────────────────────────────────────────

func TestPayrollApprovePeriod_Success(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	period, _ := svc.CreatePeriod(testContext(), "CMP001", 2026, 7)
	_ = svc.CreateRun(testContext(), &domain.PayrollRun{
		PeriodID: period.ID, EmployeeID: "NV001", CompanyID: "CMP001", BaseSalary: 10_000_000,
	})
	_ = svc.SubmitPeriod(testContext(), period.ID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/periods/"+period.ID+"/approve", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPayrollApprovePeriod_NotProcessing(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	period, _ := svc.CreatePeriod(testContext(), "CMP001", 2026, 7)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/periods/"+period.ID+"/approve", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ─── Runs ───────────────────────────────────────────────────────

func TestPayrollListRuns_Success(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	period, _ := svc.CreatePeriod(testContext(), "CMP001", 2026, 7)
	_ = svc.CreateRun(testContext(), &domain.PayrollRun{
		PeriodID:   period.ID,
		EmployeeID: "NV001",
		CompanyID:  "CMP001",
		BaseSalary: 10_000_000,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/periods/"+period.ID+"/runs", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var runs []domain.PayrollRun
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &runs))
	assert.Len(t, runs, 1)
}

// ─── Summary ────────────────────────────────────────────────────

func TestPayrollGetPeriodSummary_Empty(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	period, _ := svc.CreatePeriod(testContext(), "CMP001", 2026, 7)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/periods/"+period.ID+"/summary", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var summary domain.PayrollSummary
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &summary))
	assert.Equal(t, 0, summary.EmployeeCount)
}

// ─── Timekeeping ────────────────────────────────────────────────

func TestPayrollCreateTimekeeping_Success(t *testing.T) {
	r, _ := setupPayrollHandlerTest(t)
	body := `{"employee_id":"NV001","company_id":"CMP001","date":"2026-07-01","hours_worked":8,"ot_hours":2}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/timekeeping", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPayrollListTimekeeping_Success(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	_ = svc.CreateTimekeeping(testContext(), &domain.TimekeepingRecord{
		EmployeeID:  "NV001",
		CompanyID:   "CMP001",
		Date:        "2026-07-01",
		HoursWorked: 8,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/timekeeping?employee_id=NV001", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var records []domain.TimekeepingRecord
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &records))
	assert.Len(t, records, 1)
}

// ─── Leave ──────────────────────────────────────────────────────

func TestPayrollRequestLeave_Success(t *testing.T) {
	r, _ := setupPayrollHandlerTest(t)
	body := `{"employee_id":"NV001","company_id":"CMP001","leave_type":"ANNUAL","start_date":"2026-07-01","end_date":"2026-07-03","days":3}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/leave", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPayrollListPendingLeaveRequests_Success(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	_ = svc.RequestLeave(testContext(), &domain.LeaveRequest{
		EmployeeID: "NV001",
		CompanyID:  "CMP001",
		LeaveType:  domain.LeaveAnnual,
		StartDate:  "2026-07-01",
		EndDate:    "2026-07-03",
		Days:       3,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/leave/pending?company_id=CMP001", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var requests []domain.LeaveRequest
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &requests))
	assert.Len(t, requests, 1)
}

func TestPayrollApproveLeaveRequest_Success(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	lr := &domain.LeaveRequest{
		EmployeeID: "NV001",
		CompanyID:  "CMP001",
		LeaveType:  domain.LeaveAnnual,
		StartDate:  "2026-07-01",
		EndDate:    "2026-07-03",
		Days:       3,
	}
	_ = svc.RequestLeave(testContext(), lr)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/leave/"+lr.ID+"/approve", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPayrollRejectLeaveRequest_Success(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	lr := &domain.LeaveRequest{
		EmployeeID: "NV001",
		CompanyID:  "CMP001",
		LeaveType:  domain.LeaveAnnual,
		StartDate:  "2026-07-01",
		EndDate:    "2026-07-03",
		Days:       3,
	}
	_ = svc.RequestLeave(testContext(), lr)

	body := `{"reason":"Not enough coverage"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/leave/"+lr.ID+"/reject", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ─── Config ─────────────────────────────────────────────────────

func TestPayrollSetConfig_Success(t *testing.T) {
	r, _ := setupPayrollHandlerTest(t)
	body := `{"company_id":"CMP001","config_key":"OT_RATE","config_value":"1.5","effective_from":"2026-01-01"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/payroll/config", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPayrollGetConfig_NotFound(t *testing.T) {
	r, _ := setupPayrollHandlerTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/config?company_id=CMP001&key=NONEXISTENT", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ─── Dependants ─────────────────────────────────────────────────

func TestPayrollListDependants_Success(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	_ = svc.CreateDependant(testContext(), &domain.Dependant{
		EmployeeID:   "NV001",
		FullName:     "Nguyen Van B",
		Relationship: "CHILD",
		DateOfBirth:  "2015-06-15",
		IsActive:     true,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/employees/NV001/dependants", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var deps []domain.Dependant
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &deps))
	assert.Len(t, deps, 1)
}

func TestPayrollCreateDependant_Success(t *testing.T) {
	r, _ := setupPayrollHandlerTest(t)
	body := `{"full_name":"Nguyen Van B","relationship":"CHILD","date_of_birth":"2015-06-15","is_active":true}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/employees/NV001/dependants", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPayrollDeleteDependant_Success(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	dep := &domain.Dependant{
		EmployeeID:   "NV001",
		FullName:     "Nguyen Van B",
		Relationship: "CHILD",
		DateOfBirth:  "2015-06-15",
		IsActive:     true,
	}
	_ = svc.CreateDependant(testContext(), dep)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/payroll/employees/NV001/dependants/"+dep.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestPayrollDeleteDependant_NotFound(t *testing.T) {
	r, _ := setupPayrollHandlerTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/payroll/employees/NV001/dependants/nonexistent", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ─── Salary Components ──────────────────────────────────────────

func TestPayrollCreateComponent_Success(t *testing.T) {
	r, _ := setupPayrollHandlerTest(t)
	body := `{"company_id":"CMP001","code":"BS","name":"Lương cơ bản","type":"INCOME","calculation":"FIXED","default_value":10000000,"is_taxable":true,"is_insurable":true}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/components", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var comp domain.SalaryComponent
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &comp))
	assert.Equal(t, "BS", comp.Code)
	assert.Equal(t, "INCOME", comp.Type)
}

func TestPayrollCreateComponent_Duplicate(t *testing.T) {
	r, _ := setupPayrollHandlerTest(t)
	body := `{"company_id":"CMP001","code":"BS","name":"Lương cơ bản","type":"INCOME","calculation":"FIXED"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/components", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/v1/payroll/components", strings.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusConflict, w2.Code)
}

func TestPayrollListComponents_Success(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	_ = svc.CreateComponent(testContext(), &domain.SalaryComponent{
		CompanyID: "CMP001", Code: "BS", Name: "Lương cơ bản",
		Type: "INCOME", Calculation: "FIXED", DefaultValue: 10_000_000,
	})
	_ = svc.CreateComponent(testContext(), &domain.SalaryComponent{
		CompanyID: "CMP001", Code: "TA", Name: "Phụ cấp ăn trưa",
		Type: "INCOME", Calculation: "FIXED", DefaultValue: 1_000_000,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/components?company_id=CMP001", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var comps []domain.SalaryComponent
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &comps))
	assert.Len(t, comps, 2)
}

func TestPayrollUpdateComponent_Success(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	comp := &domain.SalaryComponent{
		CompanyID: "CMP001", Code: "BS", Name: "Lương cơ bản",
		Type: "INCOME", Calculation: "FIXED", DefaultValue: 10_000_000,
	}
	_ = svc.CreateComponent(testContext(), comp)

	body := `{"name":"Lương CB updated","default_value":12000000}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/payroll/components/"+comp.ID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPayrollDeleteComponent_Success(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	comp := &domain.SalaryComponent{
		CompanyID: "CMP001", Code: "BS", Name: "Lương cơ bản",
		Type: "INCOME", Calculation: "FIXED",
	}
	_ = svc.CreateComponent(testContext(), comp)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/payroll/components/"+comp.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestPayrollDeleteComponent_NotFound(t *testing.T) {
	r, _ := setupPayrollHandlerTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/payroll/components/nonexistent", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ─── Salary Templates ───────────────────────────────────────────

func TestPayrollCreateTemplate_Success(t *testing.T) {
	r, _ := setupPayrollHandlerTest(t)
	body := `{"company_id":"CMP001","name":"Office Staff"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/templates", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var tmpl domain.SalaryTemplate
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &tmpl))
	assert.Equal(t, "Office Staff", tmpl.Name)
}

func TestPayrollCreateTemplate_Duplicate(t *testing.T) {
	r, _ := setupPayrollHandlerTest(t)
	body := `{"company_id":"CMP001","name":"Office Staff"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/templates", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/v1/payroll/templates", strings.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusConflict, w2.Code)
}

func TestPayrollListTemplates_Success(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	_ = svc.CreateTemplate(testContext(), &domain.SalaryTemplate{
		CompanyID: "CMP001", Name: "Office Staff",
	})
	_ = svc.CreateTemplate(testContext(), &domain.SalaryTemplate{
		CompanyID: "CMP001", Name: "Factory Worker",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/templates?company_id=CMP001", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var tmpls []domain.SalaryTemplate
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &tmpls))
	assert.Len(t, tmpls, 2)
}

func TestPayrollUpdateTemplate_Success(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	tmpl := &domain.SalaryTemplate{CompanyID: "CMP001", Name: "Office Staff"}
	_ = svc.CreateTemplate(testContext(), tmpl)

	body := `{"name":"Office Staff v2"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/payroll/templates/"+tmpl.ID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPayrollDeleteTemplate_Success(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	tmpl := &domain.SalaryTemplate{CompanyID: "CMP001", Name: "Office Staff"}
	_ = svc.CreateTemplate(testContext(), tmpl)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/payroll/templates/"+tmpl.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestPayrollDeleteTemplate_NotFound(t *testing.T) {
	r, _ := setupPayrollHandlerTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/payroll/templates/nonexistent", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ─── Payslip PDF + Send ─────────────────────────────────────────

func TestPayrollGetPayslipPDF_Success(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	period, _ := svc.CreatePeriod(testContext(), "CMP001", 2026, 7)
	run := &domain.PayrollRun{
		PeriodID:   period.ID,
		EmployeeID: "NV001",
		CompanyID:  "CMP001",
		BaseSalary: 10_000_000,
		NetPay:     8_500_000,
	}
	_ = svc.CreateRun(testContext(), run)
	payslip, _ := svc.GeneratePayslip(testContext(), run.ID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/payslips/"+payslip.RunID+"/pdf", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
	assert.True(t, w.Body.Len() > 100) // PDF has content
}

func TestPayrollGetPayslipPDF_NotFound(t *testing.T) {
	r, _ := setupPayrollHandlerTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/payslips/nonexistent/pdf", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPayrollSendPayslip_Success(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	period, _ := svc.CreatePeriod(testContext(), "CMP001", 2026, 7)
	run := &domain.PayrollRun{
		PeriodID:   period.ID,
		EmployeeID: "NV001",
		CompanyID:  "CMP001",
		BaseSalary: 10_000_000,
	}
	_ = svc.CreateRun(testContext(), run)
	payslip, _ := svc.GeneratePayslip(testContext(), run.ID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/payslips/"+payslip.RunID+"/send", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "payslip sent")
}

func TestPayrollSendPayslip_NotFound(t *testing.T) {
	r, _ := setupPayrollHandlerTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/payslips/nonexistent/send", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ─── Reports ────────────────────────────────────────────────────

func TestPayrollGetInsuranceSummary_Empty(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	period, _ := svc.CreatePeriod(testContext(), "CMP001", 2026, 7)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/reports/insurance?period_id="+period.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var summary domain.InsuranceSummary
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &summary))
	assert.Equal(t, 0, summary.EmployeeCount)
}

func TestPayrollGetPITSummary_Empty(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	period, _ := svc.CreatePeriod(testContext(), "CMP001", 2026, 7)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/reports/pit?period_id="+period.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var summary domain.PITSummary
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &summary))
	assert.Equal(t, 0, summary.EmployeeCount)
}

func TestPayrollGetOvertimeSummary_Empty(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	period, _ := svc.CreatePeriod(testContext(), "CMP001", 2026, 7)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/reports/overtime?period_id="+period.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var summary domain.OvertimeSummary
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &summary))
	assert.Equal(t, 0, summary.EmployeesWithOT)
}

func TestPayrollGetLeaveBalanceReport_Empty(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	period, _ := svc.CreatePeriod(testContext(), "CMP001", 2026, 7)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/reports/leave-balance?period_id="+period.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var reports []domain.LeaveBalanceReport
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &reports))
	assert.Len(t, reports, 0)
}

func TestPayrollGetInsuranceSummary_WithData(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	period, _ := svc.CreatePeriod(testContext(), "CMP001", 2026, 7)
	_ = svc.CreateRun(testContext(), &domain.PayrollRun{
		PeriodID:   period.ID,
		EmployeeID: "NV001",
		CompanyID:  "CMP001",
		BaseSalary: 10_000_000,
		SIDeduction: 800_000,
		HIDeduction: 150_000,
		UIDeduction: 100_000,
		EmployerSI:  1_750_000,
		EmployerHI:  300_000,
		EmployerUI:  100_000,
	})
	_ = svc.CreateRun(testContext(), &domain.PayrollRun{
		PeriodID:   period.ID,
		EmployeeID: "NV002",
		CompanyID:  "CMP001",
		BaseSalary: 15_000_000,
		SIDeduction: 1_200_000,
		HIDeduction: 225_000,
		UIDeduction: 150_000,
		EmployerSI:  2_625_000,
		EmployerHI:  450_000,
		EmployerUI:  150_000,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/reports/insurance?period_id="+period.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var summary domain.InsuranceSummary
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &summary))
	assert.Equal(t, 2, summary.EmployeeCount)
	assert.Equal(t, 2_000_000.0, summary.TotalEmployeeSI)
	assert.Equal(t, 4_375_000.0, summary.TotalEmployerSI)
}

// ─── Holiday Tests ──────────────────────────────────────────────

func TestCreateHoliday_Success(t *testing.T) {
	r, _ := setupPayrollHandlerTest(t)

	body := `{"name":"Tết Nguyên Đán","date":"2026-01-29","year":2026,"company_id":"CMP001"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/holidays", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var hol domain.PayrollHoliday
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &hol))
	assert.Equal(t, "Tết Nguyên Đán", hol.Name)
	assert.Equal(t, 2026, hol.Year)
}

func TestListHolidays_Success(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)

	_ = svc.CreateHoliday(testContext(), &domain.PayrollHoliday{
		CompanyID: "CMP001", Name: "Giỗ Tổ Hùng Vương", Date: "2026-04-02", Year: 2026,
	})
	_ = svc.CreateHoliday(testContext(), &domain.PayrollHoliday{
		CompanyID: "CMP001", Name: "Giải phóng", Date: "2026-04-30", Year: 2026,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/holidays?company_id=CMP001&year=2026", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var holidays []domain.PayrollHoliday
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &holidays))
	assert.Equal(t, 2, len(holidays))
}

func TestDeleteHoliday_Success(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)

	hol := &domain.PayrollHoliday{CompanyID: "CMP001", Name: "Test", Date: "2026-01-01", Year: 2026}
	_ = svc.CreateHoliday(testContext(), hol)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/payroll/holidays/"+hol.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteHoliday_NotFound(t *testing.T) {
	r, _ := setupPayrollHandlerTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/payroll/holidays/nonexistent", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ─── Declaration Tests ──────────────────────────────────────────

func TestGenerateD02TS_Success(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)

	period, _ := svc.CreatePeriod(testContext(), "CMP001", 2026, 1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/declarations/d02-ts?period_id="+period.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "D02TS")
}

func TestGenerate05KKTNCN_Success(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)

	period, _ := svc.CreatePeriod(testContext(), "CMP001", 2026, 1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/declarations/05-kk-tncn?period_id="+period.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "KK_TNCN")
}

func TestGenerateTK3TS_Success(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)

	period, _ := svc.CreatePeriod(testContext(), "CMP001", 2026, 1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payroll/declarations/tk3-ts?period_id="+period.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "TK3TS")
}

// ─── CSV Import Test ────────────────────────────────────────────

func TestImportTimekeepingCSV_Success(t *testing.T) {
	r, _ := setupPayrollHandlerTest(t)

	csvContent := "employee_code,date,clock_in,clock_out,ot_hours,night_hours,leave_type,notes\nEMP001,2026-01-15,08:00,17:00,2,1,,overtime\n"
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "timekeeping.csv")
	part.Write([]byte(csvContent))
	writer.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/timekeeping/import?company_id=CMP001", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestImportTimekeepingCSV_NoFile(t *testing.T) {
	r, _ := setupPayrollHandlerTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/timekeeping/import?company_id=CMP001", nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestImportTimekeepingCSV_NoCompanyID(t *testing.T) {
	r, _ := setupPayrollHandlerTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/timekeeping/import", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
