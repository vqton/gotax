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

func setupPayrollHandlerTest(t *testing.T) (*gin.Engine, *service.PayrollService) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemoryPayrollRepo()
	svc := service.NewPayrollService(repo)
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

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payroll/periods/"+period.ID+"/approve", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPayrollApprovePeriod_NotDraft(t *testing.T) {
	r, svc := setupPayrollHandlerTest(t)
	period, _ := svc.CreatePeriod(testContext(), "CMP001", 2026, 7)
	_ = svc.ApprovePeriod(testContext(), period.ID, "admin")

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
