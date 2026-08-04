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
		configs:        make(map[string]*domain.PayrollConfig),
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
