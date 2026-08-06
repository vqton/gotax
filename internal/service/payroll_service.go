package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"time"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"

	"gotax/internal/domain"
	payrollxml "gotax/internal/payroll"
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
	UpdateDependant(ctx context.Context, d *domain.Dependant) error
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
	DeleteTimekeeping(ctx context.Context, id string) error

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

	// Salary components
	CreateComponent(ctx context.Context, sc *domain.SalaryComponent) error
	GetComponent(ctx context.Context, id string) (*domain.SalaryComponent, error)
	UpdateComponent(ctx context.Context, sc *domain.SalaryComponent) error
	ListComponents(ctx context.Context, companyID string) ([]domain.SalaryComponent, error)
	DeleteComponent(ctx context.Context, id string) error

	// Salary templates
	CreateTemplate(ctx context.Context, t *domain.SalaryTemplate) error
	GetTemplate(ctx context.Context, id string) (*domain.SalaryTemplate, error)
	UpdateTemplate(ctx context.Context, t *domain.SalaryTemplate) error
	ListTemplates(ctx context.Context, companyID string) ([]domain.SalaryTemplate, error)
	DeleteTemplate(ctx context.Context, id string) error

	// Holidays
	CreateHoliday(ctx context.Context, h *domain.PayrollHoliday) error
	ListHolidays(ctx context.Context, companyID string, year int) ([]domain.PayrollHoliday, error)
	DeleteHoliday(ctx context.Context, id string) error
}

// ─── Service ────────────────────────────────────────────────────

// PayrollService implements payroll business logic.
type PayrollService struct {
	repo      PayrollRepository
	companyRepo domain.CompanyRepository
}

// NewPayrollService creates a new PayrollService.
func NewPayrollService(repo PayrollRepository, companyRepo domain.CompanyRepository) *PayrollService {
	return &PayrollService{repo: repo, companyRepo: companyRepo}
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

	if period.Status != domain.PayrollProcessing {
		return domain.ErrPayrollPeriodNotDraft
	}

	now := time.Now()
	period.Status = domain.PayrollApproved
	period.ApprovedBy = approvedBy
	period.ApprovedAt = &now
	period.UpdatedAt = now

	return s.repo.UpdatePeriod(ctx, period)
}

func (s *PayrollService) SubmitPeriod(ctx context.Context, periodID string) error {
	period, err := s.repo.GetPeriod(ctx, periodID)
	if err != nil {
		return err
	}
	if period.Status != domain.PayrollDraft {
		return domain.ErrPayrollPeriodNotDraft
	}
	runs, err := s.repo.ListRunsByPeriod(ctx, periodID)
	if err != nil {
		return err
	}
	if len(runs) == 0 {
		return domain.ErrPayrollNoEmployees
	}
	now := time.Now()
	period.Status = domain.PayrollProcessing
	period.UpdatedAt = now
	return s.repo.UpdatePeriod(ctx, period)
}

func (s *PayrollService) RejectPeriod(ctx context.Context, periodID, rejectedBy string) error {
	period, err := s.repo.GetPeriod(ctx, periodID)
	if err != nil {
		return err
	}
	if period.Status != domain.PayrollProcessing {
		return domain.ErrPayrollPeriodNotDraft
	}
	period.Status = domain.PayrollDraft
	period.ReviewedBy = rejectedBy
	now := time.Now()
	period.ReviewedAt = &now
	period.UpdatedAt = now
	return s.repo.UpdatePeriod(ctx, period)
}

func (s *PayrollService) CalculatePeriod(ctx context.Context, periodID string) error {
	period, err := s.repo.GetPeriod(ctx, periodID)
	if err != nil {
		return err
	}
	infos, err := s.repo.ListEmployeePayrollInfos(ctx, period.CompanyID)
	if err != nil {
		return err
	}
	if len(infos) == 0 {
		return domain.ErrPayrollNoEmployees
	}
	for _, info := range infos {
		run := domain.CalculateEmployeePayroll(info, period, nil)
		run.ID = generateID()
		run.PeriodID = periodID
		run.EmployeeID = info.EmployeeID
		run.CompanyID = period.CompanyID
		run.Status = "DRAFT"
		run.CreatedAt = time.Now()
		run.UpdatedAt = time.Now()
		if err := s.repo.CreateRun(ctx, &run); err != nil {
			return err
		}
	}
	period.Status = domain.PayrollProcessing
	period.UpdatedAt = time.Now()
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

// ─── Net-to-Gross ───────────────────────────────────────────────

func (s *PayrollService) CalcNetToGross(ctx context.Context, input domain.NetToGrossInput) domain.NetToGrossResult {
	return domain.CalcNetToGross(input)
}

// ─── 13th-Month Salary ─────────────────────────────────────────

func (s *PayrollService) CalcThirteenthMonth(ctx context.Context, input domain.ThirteenthMonthInput) domain.ThirteenthMonthResult {
	return domain.CalcThirteenthMonth(input)
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

// ─── Reports ────────────────────────────────────────────────────

func (s *PayrollService) GetInsuranceSummary(ctx context.Context, periodID string) (*domain.InsuranceSummary, error) {
	runs, err := s.repo.ListRunsByPeriod(ctx, periodID)
	if err != nil {
		return nil, err
	}
	summary := &domain.InsuranceSummary{PeriodID: periodID}
	for _, run := range runs {
		summary.TotalEmployeeSI += run.SIDeduction
		summary.TotalEmployeeHI += run.HIDeduction
		summary.TotalEmployeeUI += run.UIDeduction
		summary.TotalEmployerSI += run.EmployerSI
		summary.TotalEmployerHI += run.EmployerHI
		summary.TotalEmployerUI += run.EmployerUI
		summary.TotalSI += run.SIDeduction + run.EmployerSI
		summary.TotalHI += run.HIDeduction + run.EmployerHI
		summary.TotalUI += run.UIDeduction + run.EmployerUI
		summary.EmployeeCount++
	}
	return summary, nil
}

func (s *PayrollService) GetPITSummary(ctx context.Context, periodID string) (*domain.PITSummary, error) {
	runs, err := s.repo.ListRunsByPeriod(ctx, periodID)
	if err != nil {
		return nil, err
	}
	summary := &domain.PITSummary{PeriodID: periodID}
	for _, run := range runs {
		summary.TotalPIT += run.PITAmount
		summary.TotalTaxableIncome += run.GrossSalary - run.SIDeduction - run.HIDeduction - run.UIDeduction - run.TradeUnionDues
		summary.EmployeeCount++
		if run.PITAmount > 0 {
			summary.EmployeesWithPIT++
		}
	}
	return summary, nil
}

func (s *PayrollService) GetOvertimeSummary(ctx context.Context, periodID string) (*domain.OvertimeSummary, error) {
	runs, err := s.repo.ListRunsByPeriod(ctx, periodID)
	if err != nil {
		return nil, err
	}
	summary := &domain.OvertimeSummary{PeriodID: periodID}
	for _, run := range runs {
		summary.TotalOTHours += run.OTHours
		summary.TotalOTPay += run.OTPay
		summary.TotalNightHours += run.NightShiftHours
		summary.TotalNightPay += run.NightShiftPay
		if run.OTHours > 0 {
			summary.EmployeesWithOT++
		}
	}
	return summary, nil
}

func (s *PayrollService) GetLeaveBalanceReport(ctx context.Context, periodID string) ([]domain.LeaveBalanceReport, error) {
	period, err := s.repo.GetPeriod(ctx, periodID)
	if err != nil {
		return nil, err
	}
	runs, err := s.repo.ListRunsByPeriod(ctx, periodID)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var reports []domain.LeaveBalanceReport
	for _, run := range runs {
		if seen[run.EmployeeID] {
			continue
		}
		seen[run.EmployeeID] = true

		leaveTypes := []domain.LeaveType{domain.LeaveAnnual, domain.LeaveSick, domain.LeaveMaternity}
		for _, lt := range leaveTypes {
			balance, err := s.repo.GetLeaveBalance(ctx, run.EmployeeID, period.Year, lt)
			if err != nil {
				continue
			}
			reports = append(reports, domain.LeaveBalanceReport{
				EmployeeID: run.EmployeeID,
				Year:       period.Year,
				LeaveType:  string(lt),
				Entitled:   balance.Entitled,
				Used:       balance.Used,
				Remaining:  balance.Remaining,
			})
		}
	}
	return reports, nil
}

// ─── Dependants ─────────────────────────────────────────────────

func (s *PayrollService) ListDependants(ctx context.Context, employeeID string) ([]domain.Dependant, error) {
	return s.repo.GetDependants(ctx, employeeID)
}

func (s *PayrollService) CreateDependant(ctx context.Context, d *domain.Dependant) error {
	return s.repo.CreateDependant(ctx, d)
}

func (s *PayrollService) DeleteDependant(ctx context.Context, id string) error {
	return s.repo.DeleteDependant(ctx, id)
}

func (s *PayrollService) UpdateDependant(ctx context.Context, d *domain.Dependant) error {
	return s.repo.UpdateDependant(ctx, d)
}

// ─── Timekeeping ────────────────────────────────────────────────

func (s *PayrollService) CreateTimekeeping(ctx context.Context, t *domain.TimekeepingRecord) error {
	return s.repo.CreateTimekeeping(ctx, t)
}

func (s *PayrollService) ListTimekeeping(ctx context.Context, employeeID, startDate, endDate string) ([]domain.TimekeepingRecord, error) {
	return s.repo.ListTimekeeping(ctx, employeeID, startDate, endDate)
}

func (s *PayrollService) BulkCreateTimekeeping(ctx context.Context, records []domain.TimekeepingRecord) error {
	return s.repo.BulkCreateTimekeeping(ctx, records)
}

func (s *PayrollService) UpdateTimekeeping(ctx context.Context, t *domain.TimekeepingRecord) error {
	return s.repo.CreateTimekeeping(ctx, t) // upsert
}

func (s *PayrollService) DeleteTimekeeping(ctx context.Context, id string) error {
	return s.repo.DeleteTimekeeping(ctx, id)
}

func (s *PayrollService) ParseTimekeepingCSV(_ context.Context, r io.Reader) ([]domain.TimekeepingRecord, error) {
	return payrollxml.ParseTimekeepingCSV(r)
}

// ─── Runs ───────────────────────────────────────────────────────

func (s *PayrollService) UpdateRun(ctx context.Context, run *domain.PayrollRun) error {
	return s.repo.UpdateRun(ctx, run)
}

// ─── Config ─────────────────────────────────────────────────────

func (s *PayrollService) GetConfig(ctx context.Context, companyID, key string) (*domain.PayrollConfig, error) {
	return s.repo.GetConfig(ctx, companyID, key)
}

func (s *PayrollService) SetConfig(ctx context.Context, c *domain.PayrollConfig) error {
	return s.repo.SetConfig(ctx, c)
}

// ─── Salary Components ──────────────────────────────────────────

func (s *PayrollService) CreateComponent(ctx context.Context, sc *domain.SalaryComponent) error {
	components, _ := s.repo.ListComponents(ctx, sc.CompanyID)
	for _, c := range components {
		if c.Code == sc.Code {
			return domain.ErrPayrollComponentExists
		}
	}
	sc.ID = generateID()
	sc.IsActive = true
	sc.CreatedAt = time.Now()
	return s.repo.CreateComponent(ctx, sc)
}

func (s *PayrollService) GetComponent(ctx context.Context, id string) (*domain.SalaryComponent, error) {
	return s.repo.GetComponent(ctx, id)
}

func (s *PayrollService) UpdateComponent(ctx context.Context, sc *domain.SalaryComponent) error {
	existing, err := s.repo.GetComponent(ctx, sc.ID)
	if err != nil {
		return err
	}
	_ = existing
	return s.repo.UpdateComponent(ctx, sc)
}

func (s *PayrollService) ListComponents(ctx context.Context, companyID string) ([]domain.SalaryComponent, error) {
	return s.repo.ListComponents(ctx, companyID)
}

func (s *PayrollService) DeleteComponent(ctx context.Context, id string) error {
	_, err := s.repo.GetComponent(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.DeleteComponent(ctx, id)
}

// ─── Salary Templates ───────────────────────────────────────────

func (s *PayrollService) CreateTemplate(ctx context.Context, t *domain.SalaryTemplate) error {
	templates, _ := s.repo.ListTemplates(ctx, t.CompanyID)
	for _, existing := range templates {
		if existing.Name == t.Name {
			return domain.ErrPayrollTemplateExists
		}
	}
	t.ID = generateID()
	t.CreatedAt = time.Now()
	return s.repo.CreateTemplate(ctx, t)
}

func (s *PayrollService) GetTemplate(ctx context.Context, id string) (*domain.SalaryTemplate, error) {
	return s.repo.GetTemplate(ctx, id)
}

func (s *PayrollService) UpdateTemplate(ctx context.Context, t *domain.SalaryTemplate) error {
	_, err := s.repo.GetTemplate(ctx, t.ID)
	if err != nil {
		return err
	}
	return s.repo.UpdateTemplate(ctx, t)
}

func (s *PayrollService) ListTemplates(ctx context.Context, companyID string) ([]domain.SalaryTemplate, error) {
	return s.repo.ListTemplates(ctx, companyID)
}

func (s *PayrollService) DeleteTemplate(ctx context.Context, id string) error {
	_, err := s.repo.GetTemplate(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.DeleteTemplate(ctx, id)
}

// ─── Holidays ───────────────────────────────────────────────────

func (s *PayrollService) CreateHoliday(ctx context.Context, h *domain.PayrollHoliday) error {
	holidays, _ := s.repo.ListHolidays(ctx, h.CompanyID, h.Year)
	for _, existing := range holidays {
		if existing.Date == h.Date && existing.Name == h.Name {
			return domain.ErrPayrollHolidayExists
		}
	}
	h.ID = generateID()
	h.CreatedAt = time.Now()
	return s.repo.CreateHoliday(ctx, h)
}

func (s *PayrollService) ListHolidays(ctx context.Context, companyID string, year int) ([]domain.PayrollHoliday, error) {
	return s.repo.ListHolidays(ctx, companyID, year)
}

func (s *PayrollService) DeleteHoliday(ctx context.Context, id string) error {
	return s.repo.DeleteHoliday(ctx, id)
}

// ─── Payslip PDF ────────────────────────────────────────────────

func (s *PayrollService) GeneratePayslipPDF(ctx context.Context, runID string) ([]byte, error) {
	run, err := s.repo.GetRun(ctx, runID)
	if err != nil {
		return nil, err
	}

	payslip, err := s.repo.GetPayslipByRun(ctx, runID)
	if err != nil {
		return nil, err
	}

	cfg := config.NewBuilder().
		WithPageSize("A4").
		WithLeftMargin(15).
		WithTopMargin(15).
		WithRightMargin(15).
		WithBottomMargin(15).
		Build()

	m := maroto.New(cfg)

	// Header
	m.AddRows(
		row.New(10).Add(
			text.NewCol(12, fmt.Sprintf("PAYSLIP — %s", payslip.PayslipNo),
				props.Text{Style: fontstyle.Bold, Size: 16, Align: align.Center}),
		),
		row.New(1).Add(col.New(12)),
	)

	// Employee info
	m.AddRows(
		row.New(6).Add(
			text.NewCol(4, "Employee:", props.Text{Style: fontstyle.Bold, Size: 9}),
			text.NewCol(8, run.EmployeeID, props.Text{Size: 9}),
		),
		row.New(6).Add(
			text.NewCol(4, "Period:", props.Text{Style: fontstyle.Bold, Size: 9}),
			text.NewCol(8, fmt.Sprintf("%d/%d", 0, 0), props.Text{Size: 9}),
		),
		row.New(1).Add(col.New(12)),
	)

	// Income section
	m.AddRows(
		row.New(6).Add(
			text.NewCol(12, "INCOME", props.Text{Style: fontstyle.Bold, Size: 10}),
		),
		payslipLine("Base Salary", run.BaseSalary),
		payslipLine("Overtime Pay", run.OTPay),
		payslipLine("Night Shift Pay", run.NightShiftPay),
		payslipLine("Leave Pay", run.LeavePay),
		payslipLine("Holiday Pay", run.HolidayPay),
		payslipLine("Allowances", run.Allowances),
		payslipLine("Bonuses", run.Bonuses),
		payslipLine("Other Income", run.OtherIncome),
		payslipLine("GROSS SALARY", run.GrossSalary),
		row.New(1).Add(col.New(12)),
	)

	// Deductions section
	m.AddRows(
		row.New(6).Add(
			text.NewCol(12, "DEDUCTIONS", props.Text{Style: fontstyle.Bold, Size: 10}),
		),
		payslipLine("Social Insurance", run.SIDeduction),
		payslipLine("Health Insurance", run.HIDeduction),
		payslipLine("Unemployment Insurance", run.UIDeduction),
		payslipLine("Trade Union Dues", run.TradeUnionDues),
		payslipLine("Personal Income Tax", run.PITAmount),
		payslipLine("Other Deductions", run.OtherDeductions),
		payslipLine("TOTAL DEDUCTIONS", run.TotalDeductions),
		row.New(1).Add(col.New(12)),
	)

	// Net pay
	m.AddRows(
		row.New(8).Add(
			text.NewCol(12, fmt.Sprintf("NET PAY: %.0f VND", run.NetPay),
				props.Text{Style: fontstyle.Bold, Size: 14, Align: align.Center}),
		),
	)

	doc, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("generate payslip pdf: %w", err)
	}
	return doc.GetBytes(), nil
}

func payslipLine(label string, amount float64) core.Row {
	return row.New(5).Add(
		text.NewCol(8, label, props.Text{Size: 9}),
		text.NewCol(4, fmt.Sprintf("%.0f", amount), props.Text{Size: 9, Align: align.Right}),
	)
}

func (s *PayrollService) SendPayslip(ctx context.Context, runID string) error {
	payslip, err := s.repo.GetPayslipByRun(ctx, runID)
	if err != nil {
		return err
	}
	now := time.Now()
	payslip.SentAt = &now
	return s.repo.CreatePayslip(ctx, payslip)
}

// ─── Helpers ────────────────────────────────────────────────────

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomHex(8)
}

// ─── Declaration XML Generation ─────────────────────────────────

func (s *PayrollService) companyInfo(ctx context.Context, companyID string) (name, taxCode string) {
	if s.companyRepo != nil {
		if c, err := s.companyRepo.GetByID(ctx, companyID); err == nil {
			return c.LegalNameVN, c.TaxCode
		}
	}
	return "Unknown", ""
}

func (s *PayrollService) GenerateD02TS(ctx context.Context, periodID string) ([]byte, error) {
	period, err := s.repo.GetPeriod(ctx, periodID)
	if err != nil {
		return nil, err
	}
	employees, err := s.repo.ListEmployeePayrollInfos(ctx, period.CompanyID)
	if err != nil {
		return nil, err
	}
	runs, err := s.repo.ListRunsByPeriod(ctx, periodID)
	if err != nil {
		return nil, err
	}
	name, taxCode := s.companyInfo(ctx, period.CompanyID)
	quarter := (period.Month-1)/3 + 1
	return payrollxml.GenerateD02TS(name, taxCode, period.Year, quarter, employees, runs)
}

func (s *PayrollService) Generate05KKTNCN(ctx context.Context, periodID string) ([]byte, error) {
	period, err := s.repo.GetPeriod(ctx, periodID)
	if err != nil {
		return nil, err
	}
	employees, err := s.repo.ListEmployeePayrollInfos(ctx, period.CompanyID)
	if err != nil {
		return nil, err
	}
	runs, err := s.repo.ListRunsByPeriod(ctx, periodID)
	if err != nil {
		return nil, err
	}
	name, taxCode := s.companyInfo(ctx, period.CompanyID)
	quarter := (period.Month-1)/3 + 1
	key := fmt.Sprintf("Q%d/%d", quarter, period.Year)
	return payrollxml.Generate05KKTNCN(taxCode, name, key, employees, runs)
}

func (s *PayrollService) GenerateTK3TS(ctx context.Context, periodID string) ([]byte, error) {
	period, err := s.repo.GetPeriod(ctx, periodID)
	if err != nil {
		return nil, err
	}
	employees, err := s.repo.ListEmployeePayrollInfos(ctx, period.CompanyID)
	if err != nil {
		return nil, err
	}
	name, taxCode := s.companyInfo(ctx, period.CompanyID)
	return payrollxml.GenerateTK3TS(name, taxCode, "", employees)
}

func randomHex(n int) string {
	const hex = "0123456789abcdef"
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		// Fallback to timestamp-based (should never happen)
		for i := range b {
			b[i] = hex[time.Now().UnixNano()%16]
			time.Sleep(1)
		}
	}
	for i := range b {
		b[i] = hex[b[i]%16]
	}
	return string(b)
}
