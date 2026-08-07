package repository

import (
	"context"
	"fmt"
	"gotax/internal/domain"
	"sort"
	"sync"
	"time"
)

type MemoryPayrollRepo struct {
	mu sync.RWMutex

	// Employee payroll info
	empPayrollInfo map[string]*domain.EmployeePayrollInfo
	empByCompany   map[string][]string

	// Dependants
	dependants    map[string]*domain.Dependant
	depByEmployee map[string][]string

	// Periods
	periods       map[string]*domain.PayrollPeriod
	periodByComp  map[string][]string
	periodByKey   map[string]string // "companyID:year:month" → periodID

	// Runs
	runs         map[string]*domain.PayrollRun
	runsByPeriod map[string][]string

	// Timekeeping
	timekeeping  map[string]*domain.TimekeepingRecord
	tkByEmployee map[string][]string

	// Leave requests
	leaveRequests map[string]*domain.LeaveRequest
	leaveByEmp   map[string][]string

	// Leave balances
	leaveBalances map[string]*domain.LeaveBalance

	// Payslips
	payslips       map[string]*domain.Payslip
	payslipByRun   map[string]string
	payslipByPeriod map[string][]string

	// Config
	configs map[string]*domain.PayrollConfig

	// Salary components
	components   map[string]*domain.SalaryComponent
	compByCompany map[string][]string

	// Salary templates
	templates   map[string]*domain.SalaryTemplate
	tmplByCompany map[string][]string

	// Holidays
	holidays      map[string]*domain.PayrollHoliday
	holidayByComp map[string][]string

	// Salary Grades
	salaryGrades   map[string]*domain.SalaryGrade
	gradeByCompany map[string][]string

	// Salary Scales
	salaryScales   map[string]*domain.SalaryScale
	scaleByCompany map[string][]string
	scaleByGrade   map[string][]string

	// Employee Salary Grade Assignments
	empSalaryGrade map[string]*domain.EmployeeSalaryGrade
	esgByCompany   map[string][]string
}

func NewMemoryPayrollRepo() *MemoryPayrollRepo {
	return &MemoryPayrollRepo{
		empPayrollInfo: make(map[string]*domain.EmployeePayrollInfo),
		empByCompany:   make(map[string][]string),
		dependants:     make(map[string]*domain.Dependant),
		depByEmployee:  make(map[string][]string),
		periods:        make(map[string]*domain.PayrollPeriod),
		periodByComp:   make(map[string][]string),
		periodByKey:    make(map[string]string),
		runs:           make(map[string]*domain.PayrollRun),
		runsByPeriod:   make(map[string][]string),
		timekeeping:    make(map[string]*domain.TimekeepingRecord),
		tkByEmployee:   make(map[string][]string),
		leaveRequests:  make(map[string]*domain.LeaveRequest),
		leaveByEmp:     make(map[string][]string),
		leaveBalances:  make(map[string]*domain.LeaveBalance),
		payslips:       make(map[string]*domain.Payslip),
		payslipByRun:   make(map[string]string),
		payslipByPeriod: make(map[string][]string),
		configs:         make(map[string]*domain.PayrollConfig),
		components:      make(map[string]*domain.SalaryComponent),
		compByCompany:   make(map[string][]string),
		templates:       make(map[string]*domain.SalaryTemplate),
		tmplByCompany:   make(map[string][]string),
		holidays:        make(map[string]*domain.PayrollHoliday),
		holidayByComp:   make(map[string][]string),
		salaryGrades:    make(map[string]*domain.SalaryGrade),
		gradeByCompany:  make(map[string][]string),
		salaryScales:    make(map[string]*domain.SalaryScale),
		scaleByCompany:  make(map[string][]string),
		scaleByGrade:    make(map[string][]string),
		empSalaryGrade:  make(map[string]*domain.EmployeeSalaryGrade),
		esgByCompany:    make(map[string][]string),
	}
}

func payrollUUID() string {
	return fmt.Sprintf("PR-%d", time.Now().UnixNano())
}

// ─── EmployeePayrollInfo ────────────────────────────────────────

func (r *MemoryPayrollRepo) CreateEmployeePayrollInfo(_ context.Context, info *domain.EmployeePayrollInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *info
	if cp.ID == "" {
		cp.ID = payrollUUID()
	}
	r.empPayrollInfo[cp.EmployeeID] = &cp
	info.ID = cp.ID
	return nil
}

func (r *MemoryPayrollRepo) GetEmployeePayrollInfo(_ context.Context, employeeID string) (*domain.EmployeePayrollInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	info, ok := r.empPayrollInfo[employeeID]
	if !ok {
		return nil, domain.ErrPayrollEmployeeNotFound
	}
	cp := *info
	return &cp, nil
}

func (r *MemoryPayrollRepo) UpdateEmployeePayrollInfo(_ context.Context, info *domain.EmployeePayrollInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.empPayrollInfo[info.EmployeeID]; !ok {
		return domain.ErrPayrollEmployeeNotFound
	}
	cp := *info
	r.empPayrollInfo[info.EmployeeID] = &cp
	return nil
}

func (r *MemoryPayrollRepo) ListEmployeePayrollInfos(_ context.Context, companyID string) ([]domain.EmployeePayrollInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.EmployeePayrollInfo
	for _, info := range r.empPayrollInfo {
		result = append(result, *info)
	}
	return result, nil
}

// ─── Dependants ─────────────────────────────────────────────────

func (r *MemoryPayrollRepo) CreateDependant(_ context.Context, d *domain.Dependant) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *d
	if cp.ID == "" {
		cp.ID = payrollUUID()
	}
	r.dependants[cp.ID] = &cp
	r.depByEmployee[cp.EmployeeID] = append(r.depByEmployee[cp.EmployeeID], cp.ID)
	d.ID = cp.ID
	return nil
}

func (r *MemoryPayrollRepo) GetDependants(_ context.Context, employeeID string) ([]domain.Dependant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := r.depByEmployee[employeeID]
	var result []domain.Dependant
	for _, id := range ids {
		if d, ok := r.dependants[id]; ok {
			result = append(result, *d)
		}
	}
	return result, nil
}

func (r *MemoryPayrollRepo) DeleteDependant(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	d, ok := r.dependants[id]
	if !ok {
		return domain.ErrPayrollDependantNotFound
	}
	empID := d.EmployeeID
	delete(r.dependants, id)
	ids := r.depByEmployee[empID]
	for i, did := range ids {
		if did == id {
			r.depByEmployee[empID] = append(ids[:i], ids[i+1:]...)
			break
		}
	}
	return nil
}

func (r *MemoryPayrollRepo) UpdateDependant(_ context.Context, d *domain.Dependant) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.dependants[d.ID]; !ok {
		return domain.ErrPayrollDependantNotFound
	}
	cp := *d
	r.dependants[d.ID] = &cp
	return nil
}

// ─── Periods ────────────────────────────────────────────────────

func (r *MemoryPayrollRepo) CreatePeriod(_ context.Context, p *domain.PayrollPeriod) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *p
	if cp.ID == "" {
		cp.ID = payrollUUID()
	}
	r.periods[cp.ID] = &cp
	r.periodByComp[cp.CompanyID] = append(r.periodByComp[cp.CompanyID], cp.ID)
	key := fmt.Sprintf("%s:%d:%d", cp.CompanyID, cp.Year, cp.Month)
	r.periodByKey[key] = cp.ID
	p.ID = cp.ID
	return nil
}

func (r *MemoryPayrollRepo) GetPeriod(_ context.Context, id string) (*domain.PayrollPeriod, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.periods[id]
	if !ok {
		return nil, domain.ErrPayrollPeriodNotFound
	}
	cp := *p
	return &cp, nil
}

func (r *MemoryPayrollRepo) UpdatePeriod(_ context.Context, p *domain.PayrollPeriod) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.periods[p.ID]; !ok {
		return domain.ErrPayrollPeriodNotFound
	}
	cp := *p
	r.periods[p.ID] = &cp
	return nil
}

func (r *MemoryPayrollRepo) ListPeriods(_ context.Context, companyID string) ([]domain.PayrollPeriod, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := r.periodByComp[companyID]
	var result []domain.PayrollPeriod
	for _, id := range ids {
		if p, ok := r.periods[id]; ok {
			result = append(result, *p)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Year != result[j].Year {
			return result[i].Year > result[j].Year
		}
		return result[i].Month > result[j].Month
	})
	return result, nil
}

func (r *MemoryPayrollRepo) GetPeriodByYearMonth(_ context.Context, companyID string, year, month int) (*domain.PayrollPeriod, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := fmt.Sprintf("%s:%d:%d", companyID, year, month)
	id, ok := r.periodByKey[key]
	if !ok {
		return nil, domain.ErrPayrollPeriodNotFound
	}
	p := r.periods[id]
	cp := *p
	return &cp, nil
}

// ─── Runs ───────────────────────────────────────────────────────

func (r *MemoryPayrollRepo) CreateRun(_ context.Context, run *domain.PayrollRun) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *run
	if cp.ID == "" {
		cp.ID = payrollUUID()
	}
	r.runs[cp.ID] = &cp
	r.runsByPeriod[cp.PeriodID] = append(r.runsByPeriod[cp.PeriodID], cp.ID)
	run.ID = cp.ID
	return nil
}

func (r *MemoryPayrollRepo) GetRun(_ context.Context, id string) (*domain.PayrollRun, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	run, ok := r.runs[id]
	if !ok {
		return nil, domain.ErrPayrollRunNotFound
	}
	cp := *run
	return &cp, nil
}

func (r *MemoryPayrollRepo) UpdateRun(_ context.Context, run *domain.PayrollRun) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.runs[run.ID]; !ok {
		return domain.ErrPayrollRunNotFound
	}
	cp := *run
	r.runs[run.ID] = &cp
	return nil
}

func (r *MemoryPayrollRepo) ListRunsByPeriod(_ context.Context, periodID string) ([]domain.PayrollRun, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := r.runsByPeriod[periodID]
	var result []domain.PayrollRun
	for _, id := range ids {
		if run, ok := r.runs[id]; ok {
			result = append(result, *run)
		}
	}
	return result, nil
}

func (r *MemoryPayrollRepo) BulkCreateRuns(ctx context.Context, runs []domain.PayrollRun) error {
	for i := range runs {
		if err := r.CreateRun(ctx, &runs[i]); err != nil {
			return err
		}
	}
	return nil
}

// ─── Timekeeping ────────────────────────────────────────────────

func (r *MemoryPayrollRepo) CreateTimekeeping(_ context.Context, t *domain.TimekeepingRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *t
	if cp.ID == "" {
		cp.ID = payrollUUID()
	}
	r.timekeeping[cp.ID] = &cp
	r.tkByEmployee[cp.EmployeeID] = append(r.tkByEmployee[cp.EmployeeID], cp.ID)
	t.ID = cp.ID
	return nil
}

func (r *MemoryPayrollRepo) ListTimekeeping(_ context.Context, employeeID, startDate, endDate string) ([]domain.TimekeepingRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := r.tkByEmployee[employeeID]
	var result []domain.TimekeepingRecord
	for _, id := range ids {
		tk, ok := r.timekeeping[id]
		if !ok {
			continue
		}
		if startDate != "" && tk.Date < startDate {
			continue
		}
		if endDate != "" && tk.Date > endDate {
			continue
		}
		result = append(result, *tk)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date < result[j].Date
	})
	return result, nil
}

func (r *MemoryPayrollRepo) BulkCreateTimekeeping(ctx context.Context, records []domain.TimekeepingRecord) error {
	for i := range records {
		if err := r.CreateTimekeeping(ctx, &records[i]); err != nil {
			return err
		}
	}
	return nil
}

func (r *MemoryPayrollRepo) DeleteTimekeeping(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	tk, ok := r.timekeeping[id]
	if !ok {
		return domain.ErrPayrollRunNotFound
	}
	delete(r.timekeeping, id)
	ids := r.tkByEmployee[tk.EmployeeID]
	for i, uid := range ids {
		if uid == id {
			r.tkByEmployee[tk.EmployeeID] = append(ids[:i], ids[i+1:]...)
			break
		}
	}
	return nil
}

// ─── Leave Requests ─────────────────────────────────────────────

func (r *MemoryPayrollRepo) CreateLeaveRequest(_ context.Context, lr *domain.LeaveRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *lr
	if cp.ID == "" {
		cp.ID = payrollUUID()
	}
	r.leaveRequests[cp.ID] = &cp
	r.leaveByEmp[cp.EmployeeID] = append(r.leaveByEmp[cp.EmployeeID], cp.ID)
	lr.ID = cp.ID
	return nil
}

func (r *MemoryPayrollRepo) ApproveLeaveRequest(_ context.Context, id, approvedBy string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	lr, ok := r.leaveRequests[id]
	if !ok {
		return domain.ErrPayrollLeaveNotFound
	}
	if lr.Status != domain.LeavePending {
		return domain.ErrPayrollLeaveAlreadyProcessed
	}
	now := time.Now()
	lr.Status = domain.LeaveApproved
	lr.ApprovedBy = approvedBy
	lr.ApprovedAt = &now
	return nil
}

func (r *MemoryPayrollRepo) RejectLeaveRequest(_ context.Context, id, _ string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	lr, ok := r.leaveRequests[id]
	if !ok {
		return domain.ErrPayrollLeaveNotFound
	}
	if lr.Status != domain.LeavePending {
		return domain.ErrPayrollLeaveAlreadyProcessed
	}
	lr.Status = domain.LeaveRejected
	return nil
}

func (r *MemoryPayrollRepo) ListPendingLeaveRequests(_ context.Context, companyID string) ([]domain.LeaveRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.LeaveRequest
	for _, lr := range r.leaveRequests {
		if lr.CompanyID == companyID && lr.Status == domain.LeavePending {
			result = append(result, *lr)
		}
	}
	return result, nil
}

// ─── Leave Balances ─────────────────────────────────────────────

func (r *MemoryPayrollRepo) GetLeaveBalance(_ context.Context, employeeID string, year int, leaveType domain.LeaveType) (*domain.LeaveBalance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := fmt.Sprintf("%s:%d:%s", employeeID, year, leaveType)
	lb, ok := r.leaveBalances[key]
	if !ok {
		return &domain.LeaveBalance{
			EmployeeID: employeeID,
			Year:       year,
			LeaveType:  leaveType,
		}, nil
	}
	cp := *lb
	return &cp, nil
}

func (r *MemoryPayrollRepo) UpdateLeaveBalance(_ context.Context, lb *domain.LeaveBalance) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := fmt.Sprintf("%s:%d:%s", lb.EmployeeID, lb.Year, lb.LeaveType)
	cp := *lb
	r.leaveBalances[key] = &cp
	return nil
}

// ─── Payslips ───────────────────────────────────────────────────

func (r *MemoryPayrollRepo) CreatePayslip(_ context.Context, p *domain.Payslip) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *p
	if cp.ID == "" {
		cp.ID = payrollUUID()
	}
	r.payslips[cp.ID] = &cp
	r.payslipByRun[cp.RunID] = cp.ID
	r.payslipByPeriod[cp.PeriodID] = append(r.payslipByPeriod[cp.PeriodID], cp.ID)
	p.ID = cp.ID
	return nil
}

func (r *MemoryPayrollRepo) GetPayslipByRun(_ context.Context, runID string) (*domain.Payslip, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.payslipByRun[runID]
	if !ok {
		return nil, domain.ErrPayrollPayslipNotFound
	}
	p := r.payslips[id]
	cp := *p
	return &cp, nil
}

func (r *MemoryPayrollRepo) ListPayslipsByPeriod(_ context.Context, periodID string) ([]domain.Payslip, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := r.payslipByPeriod[periodID]
	var result []domain.Payslip
	for _, id := range ids {
		if p, ok := r.payslips[id]; ok {
			result = append(result, *p)
		}
	}
	return result, nil
}

// ─── Config ─────────────────────────────────────────────────────

func (r *MemoryPayrollRepo) GetConfig(_ context.Context, companyID, key string) (*domain.PayrollConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	k := fmt.Sprintf("%s:%s", companyID, key)
	c, ok := r.configs[k]
	if !ok {
		return nil, domain.ErrPayrollConfigNotFound
	}
	cp := *c
	return &cp, nil
}

func (r *MemoryPayrollRepo) SetConfig(_ context.Context, c *domain.PayrollConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	k := fmt.Sprintf("%s:%s", c.CompanyID, c.ConfigKey)
	cp := *c
	r.configs[k] = &cp
	return nil
}

// ─── Salary Components ──────────────────────────────────────────

func (r *MemoryPayrollRepo) CreateComponent(_ context.Context, sc *domain.SalaryComponent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *sc
	r.components[cp.ID] = &cp
	r.compByCompany[cp.CompanyID] = append(r.compByCompany[cp.CompanyID], cp.ID)
	return nil
}

func (r *MemoryPayrollRepo) GetComponent(_ context.Context, id string) (*domain.SalaryComponent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.components[id]
	if !ok {
		return nil, domain.ErrPayrollComponentNotFound
	}
	cp := *c
	return &cp, nil
}

func (r *MemoryPayrollRepo) UpdateComponent(_ context.Context, sc *domain.SalaryComponent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.components[sc.ID]; !ok {
		return domain.ErrPayrollComponentNotFound
	}
	cp := *sc
	r.components[sc.ID] = &cp
	return nil
}

func (r *MemoryPayrollRepo) ListComponents(_ context.Context, companyID string) ([]domain.SalaryComponent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := r.compByCompany[companyID]
	var result []domain.SalaryComponent
	for _, id := range ids {
		if c, ok := r.components[id]; ok {
			result = append(result, *c)
		}
	}
	return result, nil
}

func (r *MemoryPayrollRepo) DeleteComponent(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.components[id]
	if !ok {
		return domain.ErrPayrollComponentNotFound
	}
	delete(r.components, id)
	ids := r.compByCompany[c.CompanyID]
	for i, uid := range ids {
		if uid == id {
			r.compByCompany[c.CompanyID] = append(ids[:i], ids[i+1:]...)
			break
		}
	}
	return nil
}

// ─── Salary Templates ───────────────────────────────────────────

func (r *MemoryPayrollRepo) CreateTemplate(_ context.Context, t *domain.SalaryTemplate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	tp := *t
	r.templates[tp.ID] = &tp
	r.tmplByCompany[tp.CompanyID] = append(r.tmplByCompany[tp.CompanyID], tp.ID)
	return nil
}

func (r *MemoryPayrollRepo) GetTemplate(_ context.Context, id string) (*domain.SalaryTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.templates[id]
	if !ok {
		return nil, domain.ErrPayrollTemplateNotFound
	}
	tp := *t
	return &tp, nil
}

func (r *MemoryPayrollRepo) UpdateTemplate(_ context.Context, t *domain.SalaryTemplate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.templates[t.ID]; !ok {
		return domain.ErrPayrollTemplateNotFound
	}
	tp := *t
	r.templates[t.ID] = &tp
	return nil
}

func (r *MemoryPayrollRepo) ListTemplates(_ context.Context, companyID string) ([]domain.SalaryTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := r.tmplByCompany[companyID]
	var result []domain.SalaryTemplate
	for _, id := range ids {
		if t, ok := r.templates[id]; ok {
			result = append(result, *t)
		}
	}
	return result, nil
}

func (r *MemoryPayrollRepo) DeleteTemplate(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.templates[id]
	if !ok {
		return domain.ErrPayrollTemplateNotFound
	}
	delete(r.templates, id)
	ids := r.tmplByCompany[t.CompanyID]
	for i, uid := range ids {
		if uid == id {
			r.tmplByCompany[t.CompanyID] = append(ids[:i], ids[i+1:]...)
			break
		}
	}
	return nil
}

// ─── Holidays ───────────────────────────────────────────────────

func (r *MemoryPayrollRepo) CreateHoliday(_ context.Context, h *domain.PayrollHoliday) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	hp := *h
	r.holidays[hp.ID] = &hp
	r.holidayByComp[hp.CompanyID] = append(r.holidayByComp[hp.CompanyID], hp.ID)
	return nil
}

func (r *MemoryPayrollRepo) ListHolidays(_ context.Context, companyID string, year int) ([]domain.PayrollHoliday, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := r.holidayByComp[companyID]
	var result []domain.PayrollHoliday
	for _, id := range ids {
		h, ok := r.holidays[id]
		if !ok {
			continue
		}
		if h.Year != year {
			continue
		}
		result = append(result, *h)
	}
	return result, nil
}

func (r *MemoryPayrollRepo) DeleteHoliday(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	h, ok := r.holidays[id]
	if !ok {
		return domain.ErrPayrollHolidayNotFound
	}
	delete(r.holidays, id)
	ids := r.holidayByComp[h.CompanyID]
	for i, uid := range ids {
		if uid == id {
			r.holidayByComp[h.CompanyID] = append(ids[:i], ids[i+1:]...)
			break
		}
	}
	return nil
}

// ─── Salary Grades ─────────────────────────────────────────────

func (r *MemoryPayrollRepo) CreateSalaryGrade(_ context.Context, g *domain.SalaryGrade) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *g
	if cp.ID == "" {
		cp.ID = payrollUUID()
	}
	if cp.MidSalary == 0 && cp.MinSalary > 0 && cp.MaxSalary > 0 {
		cp.MidSalary = (cp.MinSalary + cp.MaxSalary) / 2
	}
	now := time.Now()
	cp.CreatedAt = now
	cp.UpdatedAt = now
	r.salaryGrades[cp.ID] = &cp
	r.gradeByCompany[cp.CompanyID] = append(r.gradeByCompany[cp.CompanyID], cp.ID)
	g.ID = cp.ID
	return nil
}

func (r *MemoryPayrollRepo) GetSalaryGrade(_ context.Context, id string) (*domain.SalaryGrade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.salaryGrades[id]
	if !ok {
		return nil, fmt.Errorf("salary grade not found")
	}
	cp := *g
	return &cp, nil
}

func (r *MemoryPayrollRepo) UpdateSalaryGrade(_ context.Context, g *domain.SalaryGrade) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.salaryGrades[g.ID]; !ok {
		return fmt.Errorf("salary grade not found")
	}
	g.UpdatedAt = time.Now()
	r.salaryGrades[g.ID] = g
	return nil
}

func (r *MemoryPayrollRepo) DeleteSalaryGrade(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	g, ok := r.salaryGrades[id]
	if !ok {
		return fmt.Errorf("salary grade not found")
	}
	delete(r.salaryGrades, id)
	ids := r.gradeByCompany[g.CompanyID]
	for i, uid := range ids {
		if uid == id {
			r.gradeByCompany[g.CompanyID] = append(ids[:i], ids[i+1:]...)
			break
		}
	}
	return nil
}

func (r *MemoryPayrollRepo) ListSalaryGrades(_ context.Context, companyID string) ([]domain.SalaryGrade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := r.gradeByCompany[companyID]
	var result []domain.SalaryGrade
	for _, id := range ids {
		if g, ok := r.salaryGrades[id]; ok {
			result = append(result, *g)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Level < result[j].Level })
	return result, nil
}

// ─── Salary Scales ─────────────────────────────────────────────

func (r *MemoryPayrollRepo) CreateSalaryScale(_ context.Context, s *domain.SalaryScale) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *s
	if cp.ID == "" {
		cp.ID = payrollUUID()
	}
	now := time.Now()
	cp.CreatedAt = now
	cp.UpdatedAt = now
	r.salaryScales[cp.ID] = &cp
	r.scaleByCompany[cp.CompanyID] = append(r.scaleByCompany[cp.CompanyID], cp.ID)
	r.scaleByGrade[cp.GradeID] = append(r.scaleByGrade[cp.GradeID], cp.ID)
	s.ID = cp.ID
	return nil
}

func (r *MemoryPayrollRepo) GetSalaryScale(_ context.Context, id string) (*domain.SalaryScale, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.salaryScales[id]
	if !ok {
		return nil, fmt.Errorf("salary scale not found")
	}
	cp := *s
	return &cp, nil
}

func (r *MemoryPayrollRepo) UpdateSalaryScale(_ context.Context, s *domain.SalaryScale) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.salaryScales[s.ID]; !ok {
		return fmt.Errorf("salary scale not found")
	}
	s.UpdatedAt = time.Now()
	r.salaryScales[s.ID] = s
	return nil
}

func (r *MemoryPayrollRepo) DeleteSalaryScale(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, ok := r.salaryScales[id]
	if !ok {
		return fmt.Errorf("salary scale not found")
	}
	delete(r.salaryScales, id)
	ids := r.scaleByCompany[s.CompanyID]
	for i, uid := range ids {
		if uid == id {
			r.scaleByCompany[s.CompanyID] = append(ids[:i], ids[i+1:]...)
			break
		}
	}
	gradeIDs := r.scaleByGrade[s.GradeID]
	for i, uid := range gradeIDs {
		if uid == id {
			r.scaleByGrade[s.GradeID] = append(gradeIDs[:i], gradeIDs[i+1:]...)
			break
		}
	}
	return nil
}

func (r *MemoryPayrollRepo) ListSalaryScales(_ context.Context, companyID string) ([]domain.SalaryScale, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := r.scaleByCompany[companyID]
	var result []domain.SalaryScale
	for _, id := range ids {
		if s, ok := r.salaryScales[id]; ok {
			result = append(result, *s)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Level < result[j].Level })
	return result, nil
}

func (r *MemoryPayrollRepo) ListSalaryScalesByGrade(_ context.Context, gradeID string) ([]domain.SalaryScale, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := r.scaleByGrade[gradeID]
	var result []domain.SalaryScale
	for _, id := range ids {
		if s, ok := r.salaryScales[id]; ok {
			result = append(result, *s)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Level < result[j].Level })
	return result, nil
}

// ─── Employee Salary Grade Assignments ──────────────────────────

func (r *MemoryPayrollRepo) CreateEmployeeSalaryGrade(_ context.Context, e *domain.EmployeeSalaryGrade) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *e
	if cp.ID == "" {
		cp.ID = payrollUUID()
	}
	now := time.Now()
	cp.CreatedAt = now
	r.empSalaryGrade[cp.EmployeeID] = &cp
	r.esgByCompany[cp.CompanyID] = append(r.esgByCompany[cp.CompanyID], cp.ID)
	e.ID = cp.ID
	return nil
}

func (r *MemoryPayrollRepo) GetEmployeeSalaryGrade(_ context.Context, employeeID string) (*domain.EmployeeSalaryGrade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.empSalaryGrade[employeeID]
	if !ok {
		return nil, fmt.Errorf("employee salary grade not found")
	}
	cp := *e
	return &cp, nil
}

func (r *MemoryPayrollRepo) ListEmployeeSalaryGrades(_ context.Context, companyID string) ([]domain.EmployeeSalaryGrade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := r.esgByCompany[companyID]
	var result []domain.EmployeeSalaryGrade
	for _, id := range ids {
		for _, e := range r.empSalaryGrade {
			if e.ID == id {
				result = append(result, *e)
				break
			}
		}
	}
	return result, nil
}
