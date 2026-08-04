package service

import (
	"context"
	"time"

	"gotax/internal/domain"
)

// ─── Repository ─────────────────────────────────────────────────

// PayrollRepository defines payroll data access.
type PayrollRepository interface {
	// Employee payroll info
	CreateEmployeePayrollInfo(ctx context.Context, info *domain.EmployeePayrollInfo) error
	GetEmployeePayrollInfo(ctx context.Context, employeeID string) (*domain.EmployeePayrollInfo, error)
	UpdateEmployeePayrollInfo(ctx context.Context, info *domain.EmployeePayrollInfo) error
	ListEmployeePayrollInfos(ctx context.Context, companyID string) ([]domain.EmployeePayrollInfo, error)

	// Dependants
	CreateDependant(ctx context.Context, d *domain.Dependant) error
	GetDependants(ctx context.Context, employeeID string) ([]domain.Dependant, error)
	DeleteDependant(ctx context.Context, id string) error

	// Payroll periods
	CreatePeriod(ctx context.Context, p *domain.PayrollPeriod) error
	GetPeriod(ctx context.Context, id string) (*domain.PayrollPeriod, error)
	UpdatePeriod(ctx context.Context, p *domain.PayrollPeriod) error
	ListPeriods(ctx context.Context, companyID string) ([]domain.PayrollPeriod, error)
	GetPeriodByYearMonth(ctx context.Context, companyID string, year, month int) (*domain.PayrollPeriod, error)

	// Payroll runs
	CreateRun(ctx context.Context, r *domain.PayrollRun) error
	GetRun(ctx context.Context, id string) (*domain.PayrollRun, error)
	UpdateRun(ctx context.Context, r *domain.PayrollRun) error
	ListRunsByPeriod(ctx context.Context, periodID string) ([]domain.PayrollRun, error)
	BulkCreateRuns(ctx context.Context, runs []domain.PayrollRun) error

	// Timekeeping
	CreateTimekeeping(ctx context.Context, t *domain.TimekeepingRecord) error
	ListTimekeeping(ctx context.Context, employeeID, startDate, endDate string) ([]domain.TimekeepingRecord, error)
	BulkCreateTimekeeping(ctx context.Context, records []domain.TimekeepingRecord) error

	// Leave requests
	CreateLeaveRequest(ctx context.Context, lr *domain.LeaveRequest) error
	ApproveLeaveRequest(ctx context.Context, id, approvedBy string) error
	RejectLeaveRequest(ctx context.Context, id, reason string) error
	ListPendingLeaveRequests(ctx context.Context, companyID string) ([]domain.LeaveRequest, error)

	// Leave balances
	GetLeaveBalance(ctx context.Context, employeeID string, year int, leaveType domain.LeaveType) (*domain.LeaveBalance, error)
	UpdateLeaveBalance(ctx context.Context, lb *domain.LeaveBalance) error

	// Payslips
	CreatePayslip(ctx context.Context, p *domain.Payslip) error
	GetPayslipByRun(ctx context.Context, runID string) (*domain.Payslip, error)
	ListPayslipsByPeriod(ctx context.Context, periodID string) ([]domain.Payslip, error)

	// Config
	GetConfig(ctx context.Context, companyID, key string) (*domain.PayrollConfig, error)
	SetConfig(ctx context.Context, c *domain.PayrollConfig) error
}

// ─── Service ────────────────────────────────────────────────────

// PayrollService implements payroll business logic.
type PayrollService struct {
	repo PayrollRepository
}

// NewPayrollService creates a new PayrollService.
func NewPayrollService(repo PayrollRepository) *PayrollService {
	return &PayrollService{repo: repo}
}

// ─── Period Management ──────────────────────────────────────────

func (s *PayrollService) CreatePeriod(ctx context.Context, companyID string, year, month int) (*domain.PayrollPeriod, error) {
	existing, _ := s.repo.GetPeriodByYearMonth(ctx, companyID, year, month)
	if existing != nil {
		return nil, domain.ErrPayrollPeriodExists
	}

	period := &domain.PayrollPeriod{
		ID:        generateID(),
		CompanyID: companyID,
		Year:      year,
		Month:     month,
		Status:    domain.PayrollDraft,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreatePeriod(ctx, period); err != nil {
		return nil, err
	}
	return period, nil
}

func (s *PayrollService) GetPeriod(ctx context.Context, id string) (*domain.PayrollPeriod, error) {
	return s.repo.GetPeriod(ctx, id)
}

func (s *PayrollService) ListPeriods(ctx context.Context, companyID string) ([]domain.PayrollPeriod, error) {
	return s.repo.ListPeriods(ctx, companyID)
}

// ─── Run Management ─────────────────────────────────────────────

func (s *PayrollService) CreateRun(ctx context.Context, run *domain.PayrollRun) error {
	run.ID = generateID()
	run.CreatedAt = time.Now()
	run.UpdatedAt = time.Now()
	return s.repo.CreateRun(ctx, run)
}

func (s *PayrollService) GetRun(ctx context.Context, id string) (*domain.PayrollRun, error) {
	return s.repo.GetRun(ctx, id)
}

func (s *PayrollService) ListRunsByPeriod(ctx context.Context, periodID string) ([]domain.PayrollRun, error) {
	return s.repo.ListRunsByPeriod(ctx, periodID)
}

func (s *PayrollService) ApprovePeriod(ctx context.Context, periodID, approvedBy string) error {
	period, err := s.repo.GetPeriod(ctx, periodID)
	if err != nil {
		return err
	}

	if period.Status != domain.PayrollDraft {
		return domain.ErrPayrollPeriodNotDraft
	}

	now := time.Now()
	period.Status = domain.PayrollApproved
	period.ApprovedBy = approvedBy
	period.ApprovedAt = &now
	period.UpdatedAt = now

	return s.repo.UpdatePeriod(ctx, period)
}

// ─── Employee Payroll Info ──────────────────────────────────────

func (s *PayrollService) CreateEmployeePayrollInfo(ctx context.Context, info *domain.EmployeePayrollInfo) error {
	info.ID = generateID()
	info.CreatedAt = time.Now()
	info.UpdatedAt = time.Now()
	return s.repo.CreateEmployeePayrollInfo(ctx, info)
}

func (s *PayrollService) GetEmployeePayrollInfo(ctx context.Context, employeeID string) (*domain.EmployeePayrollInfo, error) {
	return s.repo.GetEmployeePayrollInfo(ctx, employeeID)
}

func (s *PayrollService) UpdateEmployeePayrollInfo(ctx context.Context, info *domain.EmployeePayrollInfo) error {
	info.UpdatedAt = time.Now()
	return s.repo.UpdateEmployeePayrollInfo(ctx, info)
}

func (s *PayrollService) ListEmployeePayrollInfos(ctx context.Context, companyID string) ([]domain.EmployeePayrollInfo, error) {
	return s.repo.ListEmployeePayrollInfos(ctx, companyID)
}

// ─── Leave Management ───────────────────────────────────────────

func (s *PayrollService) RequestLeave(ctx context.Context, req *domain.LeaveRequest) error {
	req.ID = generateID()
	req.Status = domain.LeavePending
	req.CreatedAt = time.Now()
	return s.repo.CreateLeaveRequest(ctx, req)
}

func (s *PayrollService) ApproveLeaveRequest(ctx context.Context, id, approvedBy string) error {
	return s.repo.ApproveLeaveRequest(ctx, id, approvedBy)
}

func (s *PayrollService) RejectLeaveRequest(ctx context.Context, id, reason string) error {
	return s.repo.RejectLeaveRequest(ctx, id, reason)
}

func (s *PayrollService) ListPendingLeaveRequests(ctx context.Context, companyID string) ([]domain.LeaveRequest, error) {
	return s.repo.ListPendingLeaveRequests(ctx, companyID)
}

func (s *PayrollService) GetLeaveBalance(ctx context.Context, employeeID string, year int, leaveType domain.LeaveType) (*domain.LeaveBalance, error) {
	return s.repo.GetLeaveBalance(ctx, employeeID, year, leaveType)
}

// ─── Payslip Management ────────────────────────────────────────

func (s *PayrollService) GeneratePayslip(ctx context.Context, runID string) (*domain.Payslip, error) {
	existing, _ := s.repo.GetPayslipByRun(ctx, runID)
	if existing != nil {
		return existing, nil
	}

	run, err := s.repo.GetRun(ctx, runID)
	if err != nil {
		return nil, err
	}

	payslip := &domain.Payslip{
		ID:         generateID(),
		RunID:      runID,
		EmployeeID: run.EmployeeID,
		PeriodID:   run.PeriodID,
		PayslipNo:  "PS-" + run.PeriodID + "-" + run.EmployeeID,
		CreatedAt:  time.Now(),
	}

	if err := s.repo.CreatePayslip(ctx, payslip); err != nil {
		return nil, err
	}
	return payslip, nil
}

func (s *PayrollService) GetPayslipByRun(ctx context.Context, runID string) (*domain.Payslip, error) {
	return s.repo.GetPayslipByRun(ctx, runID)
}

func (s *PayrollService) ListPayslipsByPeriod(ctx context.Context, periodID string) ([]domain.Payslip, error) {
	return s.repo.ListPayslipsByPeriod(ctx, periodID)
}

// ─── Summary ────────────────────────────────────────────────────

func (s *PayrollService) GetPeriodSummary(ctx context.Context, periodID string) (*domain.PayrollSummary, error) {
	runs, err := s.repo.ListRunsByPeriod(ctx, periodID)
	if err != nil {
		return nil, err
	}

	summary := &domain.PayrollSummary{PeriodID: periodID}
	for _, run := range runs {
		summary.EmployeeCount++
		summary.TotalGross += run.GrossSalary
		summary.TotalDeductions += run.TotalDeductions
		summary.TotalNetPay += run.NetPay
		summary.TotalEmployerCost += run.TotalEmployerCost
		summary.TotalSI += run.SIDeduction + run.EmployerSI
		summary.TotalHI += run.HIDeduction + run.EmployerHI
		summary.TotalUI += run.UIDeduction + run.EmployerUI
		summary.TotalPIT += run.PITAmount
	}
	return summary, nil
}

// ─── Helpers ────────────────────────────────────────────────────

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomHex(8)
}

func randomHex(n int) string {
	const hex = "0123456789abcdef"
	b := make([]byte, n)
	for i := range b {
		b[i] = hex[time.Now().UnixNano()%16]
		time.Sleep(1)
	}
	return string(b)
}
