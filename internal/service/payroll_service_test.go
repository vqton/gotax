package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
	"gotax/internal/repository"
)

func setupPayrollService(t *testing.T) (*PayrollService, context.Context) {
	t.Helper()
	repo := repository.NewMemoryPayrollRepo()
	svc := NewPayrollService(repo, nil)
	return svc, context.Background()
}

// ─── EmployeePayrollInfo ────────────────────────────────────────

func TestCreateEmployeePayrollInfo_Success(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	info := &domain.EmployeePayrollInfo{
		EmployeeID:        "NV001",
		ContractType:      domain.ContractIndefinite,
		SalaryType:        domain.SalaryTimeBased,
		BaseSalary:        10_000_000,
		SalaryCoefficient: 2.0,
		PositionAllowance: 1_000_000,
		Region:            domain.RegionI,
	}

	err := svc.CreateEmployeePayrollInfo(ctx, info)
	require.NoError(t, err)
	assert.NotEmpty(t, info.ID)
}

func TestGetEmployeePayrollInfo_Success(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	info := &domain.EmployeePayrollInfo{
		EmployeeID:   "NV001",
		BaseSalary:   10_000_000,
		ContractType: domain.ContractIndefinite,
		Region:       domain.RegionI,
	}
	require.NoError(t, svc.CreateEmployeePayrollInfo(ctx, info))

	got, err := svc.GetEmployeePayrollInfo(ctx, "NV001")
	require.NoError(t, err)
	assert.Equal(t, float64(10_000_000), got.BaseSalary)
	assert.Equal(t, domain.RegionI, got.Region)
}

func TestGetEmployeePayrollInfo_NotFound(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	_, err := svc.GetEmployeePayrollInfo(ctx, "nonexistent")
	require.ErrorIs(t, err, domain.ErrPayrollEmployeeNotFound)
}

func TestUpdateEmployeePayrollInfo_Success(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	info := &domain.EmployeePayrollInfo{
		EmployeeID:   "NV001",
		BaseSalary:   10_000_000,
		ContractType: domain.ContractIndefinite,
		Region:       domain.RegionI,
	}
	require.NoError(t, svc.CreateEmployeePayrollInfo(ctx, info))

	info.BaseSalary = 15_000_000
	require.NoError(t, svc.UpdateEmployeePayrollInfo(ctx, info))

	got, err := svc.GetEmployeePayrollInfo(ctx, "NV001")
	require.NoError(t, err)
	assert.Equal(t, float64(15_000_000), got.BaseSalary)
}

func TestListEmployeePayrollInfos_Success(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	for i, empID := range []string{"NV001", "NV002", "NV003"} {
		info := &domain.EmployeePayrollInfo{
			EmployeeID:   empID,
			BaseSalary:   float64(10_000_000 + i*1_000_000),
			ContractType: domain.ContractIndefinite,
			Region:       domain.RegionI,
		}
		require.NoError(t, svc.CreateEmployeePayrollInfo(ctx, info))
	}

	list, err := svc.ListEmployeePayrollInfos(ctx, "CMP001")
	require.NoError(t, err)
	assert.Len(t, list, 3)
}

// ─── Periods ────────────────────────────────────────────────────

func TestCreatePeriod_Success(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	period, err := svc.CreatePeriod(ctx, "CMP001", 2026, 7)
	require.NoError(t, err)
	assert.NotEmpty(t, period.ID)
	assert.Equal(t, 2026, period.Year)
	assert.Equal(t, 7, period.Month)
	assert.Equal(t, domain.PayrollDraft, period.Status)
}

func TestCreatePeriod_Duplicate(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	_, err := svc.CreatePeriod(ctx, "CMP001", 2026, 7)
	require.NoError(t, err)

	_, err = svc.CreatePeriod(ctx, "CMP001", 2026, 7)
	require.ErrorIs(t, err, domain.ErrPayrollPeriodExists)
}

func TestGetPeriod_Success(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	created, err := svc.CreatePeriod(ctx, "CMP001", 2026, 7)
	require.NoError(t, err)

	got, err := svc.GetPeriod(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, 2026, got.Year)
	assert.Equal(t, 7, got.Month)
}

func TestListPeriods_Success(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	for m := 1; m <= 6; m++ {
		_, err := svc.CreatePeriod(ctx, "CMP001", 2026, m)
		require.NoError(t, err)
	}

	list, err := svc.ListPeriods(ctx, "CMP001")
	require.NoError(t, err)
	assert.Len(t, list, 6)
	// Newest first
	assert.Equal(t, 6, list[0].Month)
	assert.Equal(t, 1, list[5].Month)
}

// ─── Submit Period ──────────────────────────────────────────────

func TestSubmitPeriod_Success(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	period, err := svc.CreatePeriod(ctx, "CMP001", 2026, 7)
	require.NoError(t, err)

	// Create a payroll run so SubmitPeriod validates
	run := &domain.PayrollRun{
		ID: "RUN001", PeriodID: period.ID, EmployeeID: "NV001",
		GrossSalary: 10_000_000, NetPay: 8_500_000,
	}
	require.NoError(t, svc.repo.CreateRun(ctx, run))

	err = svc.SubmitPeriod(ctx, period.ID)
	require.NoError(t, err)

	got, err := svc.GetPeriod(ctx, period.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.PayrollProcessing, got.Status)
}

func TestSubmitPeriod_NoRuns(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	period, err := svc.CreatePeriod(ctx, "CMP001", 2026, 7)
	require.NoError(t, err)

	err = svc.SubmitPeriod(ctx, period.ID)
	require.ErrorIs(t, err, domain.ErrPayrollNoEmployees)
}

func TestSubmitPeriod_AlreadySubmitted(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	period, err := svc.CreatePeriod(ctx, "CMP001", 2026, 7)
	require.NoError(t, err)
	run := &domain.PayrollRun{
		ID: "RUN001", PeriodID: period.ID, EmployeeID: "NV001",
		GrossSalary: 10_000_000, NetPay: 8_500_000,
	}
	require.NoError(t, svc.repo.CreateRun(ctx, run))
	require.NoError(t, svc.SubmitPeriod(ctx, period.ID))

	err = svc.SubmitPeriod(ctx, period.ID)
	require.ErrorIs(t, err, domain.ErrPayrollPeriodNotDraft)
}

// ─── Approve Period ─────────────────────────────────────────────

func TestApprovePeriod_Success(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	period, err := svc.CreatePeriod(ctx, "CMP001", 2026, 7)
	require.NoError(t, err)
	run := &domain.PayrollRun{
		ID: "RUN001", PeriodID: period.ID, EmployeeID: "NV001",
		GrossSalary: 10_000_000, NetPay: 8_500_000,
	}
	require.NoError(t, svc.repo.CreateRun(ctx, run))
	require.NoError(t, svc.SubmitPeriod(ctx, period.ID))

	err = svc.ApprovePeriod(ctx, period.ID, "admin")
	require.NoError(t, err)

	got, err := svc.GetPeriod(ctx, period.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.PayrollApproved, got.Status)
	assert.Equal(t, "admin", got.ApprovedBy)
}

func TestApprovePeriod_NotProcessing(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	period, err := svc.CreatePeriod(ctx, "CMP001", 2026, 7)
	require.NoError(t, err)

	err = svc.ApprovePeriod(ctx, period.ID, "admin")
	require.ErrorIs(t, err, domain.ErrPayrollPeriodNotDraft)
}

// ─── Dependants ─────────────────────────────────────────────────

func TestCreateDependant_Success(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	dep := &domain.Dependant{
		EmployeeID:   "NV001",
		FullName:     "Nguyen Van B",
		Relationship: "CHILD",
		DateOfBirth:  "2015-06-15",
		IsActive:     true,
	}

	err := svc.repo.CreateDependant(ctx, dep)
	require.NoError(t, err)
	assert.NotEmpty(t, dep.ID)
}

func TestGetDependants_Success(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	dep := &domain.Dependant{
		EmployeeID:   "NV001",
		FullName:     "Nguyen Van B",
		Relationship: "CHILD",
		DateOfBirth:  "2015-06-15",
		IsActive:     true,
	}
	require.NoError(t, svc.repo.CreateDependant(ctx, dep))

	deps, err := svc.repo.GetDependants(ctx, "NV001")
	require.NoError(t, err)
	assert.Len(t, deps, 1)
	assert.Equal(t, "Nguyen Van B", deps[0].FullName)
}

func TestDeleteDependant_Success(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	dep := &domain.Dependant{
		EmployeeID:   "NV001",
		FullName:     "Nguyen Van B",
		Relationship: "CHILD",
		DateOfBirth:  "2015-06-15",
		IsActive:     true,
	}
	require.NoError(t, svc.repo.CreateDependant(ctx, dep))

	err := svc.repo.DeleteDependant(ctx, dep.ID)
	require.NoError(t, err)

	deps, err := svc.repo.GetDependants(ctx, "NV001")
	require.NoError(t, err)
	assert.Len(t, deps, 0)
}

func TestDeleteDependant_NotFound(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	err := svc.repo.DeleteDependant(ctx, "nonexistent")
	require.ErrorIs(t, err, domain.ErrPayrollDependantNotFound)
}

// ─── Runs ───────────────────────────────────────────────────────

func TestCreateRun_Success(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	period, _ := svc.CreatePeriod(ctx, "CMP001", 2026, 7)
	run := &domain.PayrollRun{
		PeriodID:   period.ID,
		EmployeeID: "NV001",
		CompanyID:  "CMP001",
		BaseSalary: 10_000_000,
		Status:     "CALCULATED",
	}

	err := svc.CreateRun(ctx, run)
	require.NoError(t, err)
	assert.NotEmpty(t, run.ID)
}

func TestGetRun_Success(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	period, _ := svc.CreatePeriod(ctx, "CMP001", 2026, 7)
	run := &domain.PayrollRun{
		PeriodID:   period.ID,
		EmployeeID: "NV001",
		CompanyID:  "CMP001",
		BaseSalary: 10_000_000,
	}
	require.NoError(t, svc.CreateRun(ctx, run))

	got, err := svc.GetRun(ctx, run.ID)
	require.NoError(t, err)
	assert.Equal(t, float64(10_000_000), got.BaseSalary)
}

func TestListRunsByPeriod_Success(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	period, _ := svc.CreatePeriod(ctx, "CMP001", 2026, 7)
	for _, empID := range []string{"NV001", "NV002", "NV003"} {
		run := &domain.PayrollRun{
			PeriodID:   period.ID,
			EmployeeID: empID,
			CompanyID:  "CMP001",
		}
		require.NoError(t, svc.CreateRun(ctx, run))
	}

	runs, err := svc.ListRunsByPeriod(ctx, period.ID)
	require.NoError(t, err)
	assert.Len(t, runs, 3)
}

// ─── Leave Requests ─────────────────────────────────────────────

func TestRequestLeave_Success(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	lr := &domain.LeaveRequest{
		EmployeeID: "NV001",
		CompanyID:  "CMP001",
		LeaveType:  domain.LeaveAnnual,
		StartDate:  "2026-07-01",
		EndDate:    "2026-07-03",
		Days:       3,
		Reason:     "Vacation",
	}
	err := svc.RequestLeave(ctx, lr)
	require.NoError(t, err)
	assert.NotEmpty(t, lr.ID)
	assert.Equal(t, domain.LeavePending, lr.Status)
}

func TestApproveLeaveRequest_Success(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	lr := &domain.LeaveRequest{
		EmployeeID: "NV001",
		CompanyID:  "CMP001",
		LeaveType:  domain.LeaveAnnual,
		StartDate:  "2026-07-01",
		EndDate:    "2026-07-03",
		Days:       3,
	}
	require.NoError(t, svc.RequestLeave(ctx, lr))

	err := svc.ApproveLeaveRequest(ctx, lr.ID, "admin")
	require.NoError(t, err)
}

func TestApproveLeaveRequest_AlreadyProcessed(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	lr := &domain.LeaveRequest{
		EmployeeID: "NV001",
		CompanyID:  "CMP001",
		LeaveType:  domain.LeaveAnnual,
		StartDate:  "2026-07-01",
		EndDate:    "2026-07-03",
		Days:       3,
	}
	require.NoError(t, svc.RequestLeave(ctx, lr))
	require.NoError(t, svc.ApproveLeaveRequest(ctx, lr.ID, "admin"))

	err := svc.ApproveLeaveRequest(ctx, lr.ID, "admin2")
	require.ErrorIs(t, err, domain.ErrPayrollLeaveAlreadyProcessed)
}

func TestRejectLeaveRequest_Success(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	lr := &domain.LeaveRequest{
		EmployeeID: "NV001",
		CompanyID:  "CMP001",
		LeaveType:  domain.LeaveAnnual,
		StartDate:  "2026-07-01",
		EndDate:    "2026-07-03",
		Days:       3,
	}
	require.NoError(t, svc.RequestLeave(ctx, lr))

	err := svc.RejectLeaveRequest(ctx, lr.ID, "Not enough coverage")
	require.NoError(t, err)
}

func TestListPendingLeaveRequests_Success(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	for _, empID := range []string{"NV001", "NV002"} {
		lr := &domain.LeaveRequest{
			EmployeeID: empID,
			CompanyID:  "CMP001",
			LeaveType:  domain.LeaveAnnual,
			StartDate:  "2026-07-01",
			EndDate:    "2026-07-03",
			Days:       3,
		}
		require.NoError(t, svc.RequestLeave(ctx, lr))
	}

	list, err := svc.ListPendingLeaveRequests(ctx, "CMP001")
	require.NoError(t, err)
	assert.Len(t, list, 2)
}

// ─── Leave Balances ─────────────────────────────────────────────

func TestGetLeaveBalance_Default(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	lb, err := svc.GetLeaveBalance(ctx, "NV001", 2026, domain.LeaveAnnual)
	require.NoError(t, err)
	assert.Equal(t, 0.0, lb.Entitled)
	assert.Equal(t, 0.0, lb.Used)
}

func TestUpdateLeaveBalance_Success(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	lb := &domain.LeaveBalance{
		EmployeeID: "NV001",
		Year:       2026,
		LeaveType:  domain.LeaveAnnual,
		Entitled:   12,
		Used:       3,
		Remaining:  9,
	}
	require.NoError(t, svc.repo.UpdateLeaveBalance(ctx, lb))

	got, err := svc.GetLeaveBalance(ctx, "NV001", 2026, domain.LeaveAnnual)
	require.NoError(t, err)
	assert.Equal(t, 12.0, got.Entitled)
	assert.Equal(t, 3.0, got.Used)
}

// ─── Payslips ───────────────────────────────────────────────────

func TestGeneratePayslip_Success(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	period, _ := svc.CreatePeriod(ctx, "CMP001", 2026, 7)
	run := &domain.PayrollRun{
		PeriodID:   period.ID,
		EmployeeID: "NV001",
		CompanyID:  "CMP001",
		BaseSalary: 10_000_000,
	}
	require.NoError(t, svc.CreateRun(ctx, run))

	payslip, err := svc.GeneratePayslip(ctx, run.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, payslip.ID)
	assert.Equal(t, run.ID, payslip.RunID)
	assert.Equal(t, "NV001", payslip.EmployeeID)
}

func TestGeneratePayslip_Idempotent(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	period, _ := svc.CreatePeriod(ctx, "CMP001", 2026, 7)
	run := &domain.PayrollRun{
		PeriodID:   period.ID,
		EmployeeID: "NV001",
		CompanyID:  "CMP001",
	}
	require.NoError(t, svc.CreateRun(ctx, run))

	p1, err := svc.GeneratePayslip(ctx, run.ID)
	require.NoError(t, err)

	p2, err := svc.GeneratePayslip(ctx, run.ID)
	require.NoError(t, err)
	assert.Equal(t, p1.ID, p2.ID)
}

// ─── Config ─────────────────────────────────────────────────────

func TestSetConfig_Success(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	cfg := &domain.PayrollConfig{
		CompanyID:     "CMP001",
		ConfigKey:     "OT_RATE",
		ConfigValue:   "1.5",
		EffectiveFrom: "2026-01-01",
	}
	err := svc.repo.SetConfig(ctx, cfg)
	require.NoError(t, err)
}

func TestGetConfig_Success(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	cfg := &domain.PayrollConfig{
		CompanyID:     "CMP001",
		ConfigKey:     "OT_RATE",
		ConfigValue:   "1.5",
		EffectiveFrom: "2026-01-01",
	}
	require.NoError(t, svc.repo.SetConfig(ctx, cfg))

	got, err := svc.repo.GetConfig(ctx, "CMP001", "OT_RATE")
	require.NoError(t, err)
	assert.Equal(t, "1.5", got.ConfigValue)
}

func TestGetConfig_NotFound(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	_, err := svc.repo.GetConfig(ctx, "CMP001", "NONEXISTENT")
	require.ErrorIs(t, err, domain.ErrPayrollConfigNotFound)
}

// ─── Summary ────────────────────────────────────────────────────

func TestGetPeriodSummary_Empty(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	period, _ := svc.CreatePeriod(ctx, "CMP001", 2026, 7)

	summary, err := svc.GetPeriodSummary(ctx, period.ID)
	require.NoError(t, err)
	assert.Equal(t, 0, summary.EmployeeCount)
	assert.Equal(t, 0.0, summary.TotalGross)
}

func TestGetPeriodSummary_WithData(t *testing.T) {
	svc, ctx := setupPayrollService(t)
	period, _ := svc.CreatePeriod(ctx, "CMP001", 2026, 7)

	runs := []domain.PayrollRun{
		{PeriodID: period.ID, EmployeeID: "NV001", CompanyID: "CMP001", GrossSalary: 10_000_000, NetPay: 8_500_000, TotalDeductions: 1_500_000, SIDeduction: 800_000, HIDeduction: 150_000, UIDeduction: 100_000, PITAmount: 450_000},
		{PeriodID: period.ID, EmployeeID: "NV002", CompanyID: "CMP001", GrossSalary: 15_000_000, NetPay: 12_500_000, TotalDeductions: 2_500_000, SIDeduction: 1_200_000, HIDeduction: 225_000, UIDeduction: 150_000, PITAmount: 925_000},
	}
	require.NoError(t, svc.repo.BulkCreateRuns(ctx, runs))

	summary, err := svc.GetPeriodSummary(ctx, period.ID)
	require.NoError(t, err)
	assert.Equal(t, 2, summary.EmployeeCount)
	assert.Equal(t, 25_000_000.0, summary.TotalGross)
	assert.Equal(t, 21_000_000.0, summary.TotalNetPay)
	assert.Equal(t, 4_000_000.0, summary.TotalDeductions)
}
