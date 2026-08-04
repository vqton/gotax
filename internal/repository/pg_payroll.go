package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gotax/internal/domain"
)

type PGPayrollRepo struct {
	db *gorm.DB
}

func NewPGPayrollRepo(db *gorm.DB) *PGPayrollRepo {
	return &PGPayrollRepo{db: db}
}

// ─── EmployeePayrollInfo ────────────────────────────────────────

func (r *PGPayrollRepo) CreateEmployeePayrollInfo(ctx context.Context, info *domain.EmployeePayrollInfo) error {
	info.CreatedAt = time.Now()
	info.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Create(info).Error
}

func (r *PGPayrollRepo) GetEmployeePayrollInfo(ctx context.Context, employeeID string) (*domain.EmployeePayrollInfo, error) {
	var info domain.EmployeePayrollInfo
	if err := r.db.WithContext(ctx).Where("employee_id = ?", employeeID).First(&info).Error; err != nil {
		return nil, err
	}
	return &info, nil
}

func (r *PGPayrollRepo) UpdateEmployeePayrollInfo(ctx context.Context, info *domain.EmployeePayrollInfo) error {
	info.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Where("employee_id = ?", info.EmployeeID).Updates(info).Error
}

func (r *PGPayrollRepo) ListEmployeePayrollInfos(ctx context.Context, companyID string) ([]domain.EmployeePayrollInfo, error) {
	var infos []domain.EmployeePayrollInfo
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Find(&infos).Error; err != nil {
		return nil, err
	}
	return infos, nil
}

// ─── Dependants ─────────────────────────────────────────────────

func (r *PGPayrollRepo) CreateDependant(ctx context.Context, d *domain.Dependant) error {
	d.CreatedAt = time.Now()
	return r.db.WithContext(ctx).Create(d).Error
}

func (r *PGPayrollRepo) GetDependants(ctx context.Context, employeeID string) ([]domain.Dependant, error) {
	var deps []domain.Dependant
	if err := r.db.WithContext(ctx).Where("employee_id = ?", employeeID).Find(&deps).Error; err != nil {
		return nil, err
	}
	return deps, nil
}

func (r *PGPayrollRepo) DeleteDependant(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Dependant{}).Error
}

// ─── Periods ────────────────────────────────────────────────────

func (r *PGPayrollRepo) CreatePeriod(ctx context.Context, p *domain.PayrollPeriod) error {
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *PGPayrollRepo) GetPeriod(ctx context.Context, id string) (*domain.PayrollPeriod, error) {
	var p domain.PayrollPeriod
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PGPayrollRepo) UpdatePeriod(ctx context.Context, p *domain.PayrollPeriod) error {
	p.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Where("id = ?", p.ID).Updates(p).Error
}

func (r *PGPayrollRepo) ListPeriods(ctx context.Context, companyID string) ([]domain.PayrollPeriod, error) {
	var periods []domain.PayrollPeriod
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).
		Order("year DESC, month DESC").Find(&periods).Error; err != nil {
		return nil, err
	}
	return periods, nil
}

func (r *PGPayrollRepo) GetPeriodByYearMonth(ctx context.Context, companyID string, year, month int) (*domain.PayrollPeriod, error) {
	var p domain.PayrollPeriod
	if err := r.db.WithContext(ctx).
		Where("company_id = ? AND year = ? AND month = ?", companyID, year, month).
		First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// ─── Runs ───────────────────────────────────────────────────────

func (r *PGPayrollRepo) CreateRun(ctx context.Context, run *domain.PayrollRun) error {
	now := time.Now()
	run.CreatedAt = now
	run.UpdatedAt = now
	return r.db.WithContext(ctx).Create(run).Error
}

func (r *PGPayrollRepo) GetRun(ctx context.Context, id string) (*domain.PayrollRun, error) {
	var run domain.PayrollRun
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&run).Error; err != nil {
		return nil, err
	}
	return &run, nil
}

func (r *PGPayrollRepo) UpdateRun(ctx context.Context, run *domain.PayrollRun) error {
	run.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Where("id = ?", run.ID).Updates(run).Error
}

func (r *PGPayrollRepo) ListRunsByPeriod(ctx context.Context, periodID string) ([]domain.PayrollRun, error) {
	var runs []domain.PayrollRun
	if err := r.db.WithContext(ctx).Where("period_id = ?", periodID).
		Order("employee_id").Find(&runs).Error; err != nil {
		return nil, err
	}
	return runs, nil
}

func (r *PGPayrollRepo) BulkCreateRuns(ctx context.Context, runs []domain.PayrollRun) error {
	if len(runs) == 0 {
		return nil
	}
	now := time.Now()
	for i := range runs {
		runs[i].CreatedAt = now
		runs[i].UpdatedAt = now
	}
	return r.db.WithContext(ctx).Create(&runs).Error
}

// ─── Timekeeping ────────────────────────────────────────────────

func (r *PGPayrollRepo) CreateTimekeeping(ctx context.Context, t *domain.TimekeepingRecord) error {
	t.CreatedAt = time.Now()
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *PGPayrollRepo) ListTimekeeping(ctx context.Context, employeeID, startDate, endDate string) ([]domain.TimekeepingRecord, error) {
	var records []domain.TimekeepingRecord
	q := r.db.WithContext(ctx).Where("employee_id = ?", employeeID)
	if startDate != "" {
		q = q.Where("date >= ?", startDate)
	}
	if endDate != "" {
		q = q.Where("date <= ?", endDate)
	}
	if err := q.Order("date").Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (r *PGPayrollRepo) BulkCreateTimekeeping(ctx context.Context, records []domain.TimekeepingRecord) error {
	if len(records) == 0 {
		return nil
	}
	now := time.Now()
	for i := range records {
		records[i].CreatedAt = now
	}
	return r.db.WithContext(ctx).Create(&records).Error
}

// ─── Leave Requests ─────────────────────────────────────────────

func (r *PGPayrollRepo) CreateLeaveRequest(ctx context.Context, lr *domain.LeaveRequest) error {
	lr.CreatedAt = time.Now()
	return r.db.WithContext(ctx).Create(lr).Error
}

func (r *PGPayrollRepo) ApproveLeaveRequest(ctx context.Context, id, approvedBy string) error {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&domain.LeaveRequest{}).
		Where("id = ? AND status = ?", id, domain.LeavePending).
		Updates(map[string]interface{}{
			"status":      domain.LeaveApproved,
			"approved_by": approvedBy,
			"approved_at": now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrPayrollLeaveNotFound
	}
	return nil
}

func (r *PGPayrollRepo) RejectLeaveRequest(ctx context.Context, id, _ string) error {
	result := r.db.WithContext(ctx).Model(&domain.LeaveRequest{}).
		Where("id = ? AND status = ?", id, domain.LeavePending).
		Update("status", domain.LeaveRejected)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrPayrollLeaveNotFound
	}
	return nil
}

func (r *PGPayrollRepo) ListPendingLeaveRequests(ctx context.Context, companyID string) ([]domain.LeaveRequest, error) {
	var requests []domain.LeaveRequest
	if err := r.db.WithContext(ctx).
		Where("company_id = ? AND status = ?", companyID, domain.LeavePending).
		Order("created_at").Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

// ─── Leave Balances ─────────────────────────────────────────────

func (r *PGPayrollRepo) GetLeaveBalance(ctx context.Context, employeeID string, year int, leaveType domain.LeaveType) (*domain.LeaveBalance, error) {
	var lb domain.LeaveBalance
	if err := r.db.WithContext(ctx).
		Where("employee_id = ? AND year = ? AND leave_type = ?", employeeID, year, leaveType).
		First(&lb).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &domain.LeaveBalance{
				EmployeeID: employeeID,
				Year:       year,
				LeaveType:  leaveType,
			}, nil
		}
		return nil, err
	}
	return &lb, nil
}

func (r *PGPayrollRepo) UpdateLeaveBalance(ctx context.Context, lb *domain.LeaveBalance) error {
	if lb.ID == "" {
		return r.db.WithContext(ctx).Create(lb).Error
	}
	return r.db.WithContext(ctx).Where("id = ?", lb.ID).Save(lb).Error
}

// ─── Payslips ───────────────────────────────────────────────────

func (r *PGPayrollRepo) CreatePayslip(ctx context.Context, p *domain.Payslip) error {
	p.CreatedAt = time.Now()
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *PGPayrollRepo) GetPayslipByRun(ctx context.Context, runID string) (*domain.Payslip, error) {
	var p domain.Payslip
	if err := r.db.WithContext(ctx).Where("run_id = ?", runID).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PGPayrollRepo) ListPayslipsByPeriod(ctx context.Context, periodID string) ([]domain.Payslip, error) {
	var payslips []domain.Payslip
	if err := r.db.WithContext(ctx).Where("period_id = ?", periodID).
		Order("created_at").Find(&payslips).Error; err != nil {
		return nil, err
	}
	return payslips, nil
}

// ─── Config ─────────────────────────────────────────────────────

func (r *PGPayrollRepo) GetConfig(ctx context.Context, companyID, key string) (*domain.PayrollConfig, error) {
	var c domain.PayrollConfig
	if err := r.db.WithContext(ctx).
		Where("company_id = ? AND config_key = ?", companyID, key).
		First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *PGPayrollRepo) SetConfig(ctx context.Context, c *domain.PayrollConfig) error {
	var existing domain.PayrollConfig
	if err := r.db.WithContext(ctx).
		Where("company_id = ? AND config_key = ?", c.CompanyID, c.ConfigKey).
		First(&existing).Error; err == nil {
		c.ID = existing.ID
		c.CreatedAt = existing.CreatedAt
		return r.db.WithContext(ctx).Save(c).Error
	}
	c.CreatedAt = time.Now()
	return r.db.WithContext(ctx).Create(c).Error
}

// ─── Salary Components ──────────────────────────────────────────

func (r *PGPayrollRepo) CreateComponent(ctx context.Context, sc *domain.SalaryComponent) error {
	sc.CreatedAt = time.Now()
	return r.db.WithContext(ctx).Create(sc).Error
}

func (r *PGPayrollRepo) GetComponent(ctx context.Context, id string) (*domain.SalaryComponent, error) {
	var c domain.SalaryComponent
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *PGPayrollRepo) UpdateComponent(ctx context.Context, sc *domain.SalaryComponent) error {
	return r.db.WithContext(ctx).Where("id = ?", sc.ID).Updates(sc).Error
}

func (r *PGPayrollRepo) ListComponents(ctx context.Context, companyID string) ([]domain.SalaryComponent, error) {
	var components []domain.SalaryComponent
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).
		Order("code").Find(&components).Error; err != nil {
		return nil, err
	}
	return components, nil
}

func (r *PGPayrollRepo) DeleteComponent(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.SalaryComponent{}).Error
}

// ─── Salary Templates ───────────────────────────────────────────

func (r *PGPayrollRepo) CreateTemplate(ctx context.Context, t *domain.SalaryTemplate) error {
	t.CreatedAt = time.Now()
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *PGPayrollRepo) GetTemplate(ctx context.Context, id string) (*domain.SalaryTemplate, error) {
	var t domain.SalaryTemplate
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *PGPayrollRepo) UpdateTemplate(ctx context.Context, t *domain.SalaryTemplate) error {
	return r.db.WithContext(ctx).Where("id = ?", t.ID).Updates(t).Error
}

func (r *PGPayrollRepo) ListTemplates(ctx context.Context, companyID string) ([]domain.SalaryTemplate, error) {
	var templates []domain.SalaryTemplate
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).
		Order("name").Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *PGPayrollRepo) DeleteTemplate(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.SalaryTemplate{}).Error
}
