package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gotax/internal/domain"
)

// ─── CostObject ────────────────────────────────────────────────────────────

type pgCostObject struct {
	ID            string    `gorm:"column:id;primaryKey"`
	CompanyID     string    `gorm:"column:company_id"`
	Code          string    `gorm:"column:code"`
	Name          string    `gorm:"column:name"`
	Type          string    `gorm:"column:type"`
	CostingMethod string    `gorm:"column:costing_method"`
	CostCenterID  *string   `gorm:"column:cost_center_id"`
	UnitOfMeasure *string   `gorm:"column:unit_of_measure"`
	StandardCost  float64   `gorm:"column:standard_cost"`
	PlanQuantity  float64   `gorm:"column:plan_quantity"`
	IsActive      bool      `gorm:"column:is_active"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (pgCostObject) TableName() string { return "cost_objects" }

type PGCostObjectRepo struct{ db *gorm.DB }

func NewPGCostObjectRepo(db *gorm.DB) *PGCostObjectRepo {
	return &PGCostObjectRepo{db: db}
}

func toPGCostObject(co *domain.CostObject) *pgCostObject {
	var created, updated time.Time
	if co.CreatedAt != "" {
		created, _ = time.Parse("2006-01-02T15:04:05Z", co.CreatedAt)
	} else {
		created = time.Now()
	}
	if co.UpdatedAt != "" {
		updated, _ = time.Parse("2006-01-02T15:04:05Z", co.UpdatedAt)
	} else {
		updated = time.Now()
	}
	m := &pgCostObject{
		ID:            co.ID,
		CompanyID:     co.CompanyID,
		Code:          co.Code,
		Name:          co.Name,
		Type:          string(co.Type),
		CostingMethod: string(co.CostingMethod),
		StandardCost:  co.StandardCost,
		PlanQuantity:  co.PlanQuantity,
		IsActive:      co.IsActive,
		CreatedAt:     created,
		UpdatedAt:     updated,
	}
	if co.CostCenterID != "" {
		m.CostCenterID = &co.CostCenterID
	}
	if co.UnitOfMeasure != "" {
		m.UnitOfMeasure = &co.UnitOfMeasure
	}
	return m
}

func toDomainCostObject(m *pgCostObject) *domain.CostObject {
	co := &domain.CostObject{
		ID:            m.ID,
		CompanyID:     m.CompanyID,
		Code:          m.Code,
		Name:          m.Name,
		Type:          domain.CostObjectType(m.Type),
		CostingMethod: domain.CostingMethod(m.CostingMethod),
		StandardCost:  m.StandardCost,
		PlanQuantity:  m.PlanQuantity,
		IsActive:      m.IsActive,
		CreatedAt:     m.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     m.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if m.CostCenterID != nil {
		co.CostCenterID = *m.CostCenterID
	}
	if m.UnitOfMeasure != nil {
		co.UnitOfMeasure = *m.UnitOfMeasure
	}
	return co
}

func (r *PGCostObjectRepo) Create(ctx context.Context, co *domain.CostObject) error {
	return r.db.WithContext(ctx).Create(toPGCostObject(co)).Error
}

func (r *PGCostObjectRepo) GetByID(ctx context.Context, id string) (*domain.CostObject, error) {
	var m pgCostObject
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, domain.ErrCostObjectNotFound
	}
	return toDomainCostObject(&m), nil
}

func (r *PGCostObjectRepo) GetByCode(ctx context.Context, companyID, code string) (*domain.CostObject, error) {
	var m pgCostObject
	if err := r.db.WithContext(ctx).Where("company_id = ? AND code = ?", companyID, code).First(&m).Error; err != nil {
		return nil, domain.ErrCostObjectNotFound
	}
	return toDomainCostObject(&m), nil
}

func (r *PGCostObjectRepo) List(ctx context.Context, companyID string) ([]domain.CostObject, error) {
	var rows []pgCostObject
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("code").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.CostObject, len(rows))
	for i := range rows {
		out[i] = *toDomainCostObject(&rows[i])
	}
	return out, nil
}

func (r *PGCostObjectRepo) Update(ctx context.Context, co *domain.CostObject) error {
	return r.db.WithContext(ctx).Save(toPGCostObject(co)).Error
}

func (r *PGCostObjectRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&pgCostObject{}, "id = ?", id).Error
}

// ─── CostPool ──────────────────────────────────────────────────────────────

type pgCostPool struct {
	ID            string    `gorm:"column:id;primaryKey"`
	CompanyID     string    `gorm:"column:company_id"`
	PeriodID      string    `gorm:"column:period_id"`
	GLAccountCode string    `gorm:"column:gl_account_code"`
	Name          string    `gorm:"column:name"`
	Status        string    `gorm:"column:status"`
	TotalAmount   float64   `gorm:"column:total_amount"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (pgCostPool) TableName() string { return "cost_pools" }

type PGCostPoolRepo struct{ db *gorm.DB }

func NewPGCostPoolRepo(db *gorm.DB) *PGCostPoolRepo {
	return &PGCostPoolRepo{db: db}
}

func toPGCostPool(cp *domain.CostPool) *pgCostPool {
	var created, updated time.Time
	if cp.CreatedAt != "" {
		created, _ = time.Parse("2006-01-02T15:04:05Z", cp.CreatedAt)
	} else {
		created = time.Now()
	}
	if cp.UpdatedAt != "" {
		updated, _ = time.Parse("2006-01-02T15:04:05Z", cp.UpdatedAt)
	} else {
		updated = time.Now()
	}
	return &pgCostPool{
		ID:            cp.ID,
		CompanyID:     cp.CompanyID,
		PeriodID:      cp.PeriodID,
		GLAccountCode: cp.GLAccountCode,
		Name:          cp.Name,
		Status:        string(cp.Status),
		TotalAmount:   cp.TotalAmount,
		CreatedAt:     created,
		UpdatedAt:     updated,
	}
}

func toDomainCostPool(m *pgCostPool) *domain.CostPool {
	return &domain.CostPool{
		ID:            m.ID,
		CompanyID:     m.CompanyID,
		PeriodID:      m.PeriodID,
		GLAccountCode: m.GLAccountCode,
		Name:          m.Name,
		Status:        domain.CostPoolStatus(m.Status),
		TotalAmount:   m.TotalAmount,
		CreatedAt:     m.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     m.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func (r *PGCostPoolRepo) Create(ctx context.Context, cp *domain.CostPool) error {
	return r.db.WithContext(ctx).Create(toPGCostPool(cp)).Error
}

func (r *PGCostPoolRepo) GetByID(ctx context.Context, id string) (*domain.CostPool, error) {
	var m pgCostPool
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, domain.ErrCostPoolNotFound
	}
	return toDomainCostPool(&m), nil
}

func (r *PGCostPoolRepo) ListByPeriod(ctx context.Context, companyID, periodID string) ([]domain.CostPool, error) {
	var rows []pgCostPool
	if err := r.db.WithContext(ctx).Where("company_id = ? AND period_id = ?", companyID, periodID).Order("gl_account_code").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.CostPool, len(rows))
	for i := range rows {
		out[i] = *toDomainCostPool(&rows[i])
	}
	return out, nil
}

func (r *PGCostPoolRepo) Update(ctx context.Context, cp *domain.CostPool) error {
	return r.db.WithContext(ctx).Save(toPGCostPool(cp)).Error
}

func (r *PGCostPoolRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&pgCostPool{}, "id = ?", id).Error
}

// ─── CostPoolLine ──────────────────────────────────────────────────────────

type pgCostPoolLine struct {
	ID           string    `gorm:"column:id;primaryKey"`
	PoolID       string    `gorm:"column:pool_id"`
	SourceType   string    `gorm:"column:source_type"`
	SourceID     *string   `gorm:"column:source_id"`
	Description  string    `gorm:"column:description"`
	Amount       float64   `gorm:"column:amount"`
	CostCenterID *string   `gorm:"column:cost_center_id"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (pgCostPoolLine) TableName() string { return "cost_pool_lines" }

type PGCostPoolLineRepo struct{ db *gorm.DB }

func NewPGCostPoolLineRepo(db *gorm.DB) *PGCostPoolLineRepo {
	return &PGCostPoolLineRepo{db: db}
}

func toPGCostPoolLine(l *domain.CostPoolLine) *pgCostPoolLine {
	var created time.Time
	if l.CreatedAt != "" {
		created, _ = time.Parse("2006-01-02T15:04:05Z", l.CreatedAt)
	} else {
		created = time.Now()
	}
	m := &pgCostPoolLine{
		ID:          l.ID,
		PoolID:      l.PoolID,
		SourceType:  l.SourceType,
		Description: l.Description,
		Amount:      l.Amount,
		CreatedAt:   created,
	}
	if l.SourceID != "" {
		m.SourceID = &l.SourceID
	}
	if l.CostCenterID != "" {
		m.CostCenterID = &l.CostCenterID
	}
	return m
}

func toDomainCostPoolLine(m *pgCostPoolLine) *domain.CostPoolLine {
	l := &domain.CostPoolLine{
		ID:          m.ID,
		PoolID:      m.PoolID,
		SourceType:  m.SourceType,
		Description: m.Description,
		Amount:      m.Amount,
		CreatedAt:   m.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if m.SourceID != nil {
		l.SourceID = *m.SourceID
	}
	if m.CostCenterID != nil {
		l.CostCenterID = *m.CostCenterID
	}
	return l
}

func (r *PGCostPoolLineRepo) Create(ctx context.Context, line *domain.CostPoolLine) error {
	return r.db.WithContext(ctx).Create(toPGCostPoolLine(line)).Error
}

func (r *PGCostPoolLineRepo) ListByPool(ctx context.Context, poolID string) ([]domain.CostPoolLine, error) {
	var rows []pgCostPoolLine
	if err := r.db.WithContext(ctx).Where("pool_id = ?", poolID).Order("created_at").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.CostPoolLine, len(rows))
	for i := range rows {
		out[i] = *toDomainCostPoolLine(&rows[i])
	}
	return out, nil
}

func (r *PGCostPoolLineRepo) DeleteByPool(ctx context.Context, poolID string) error {
	return r.db.WithContext(ctx).Delete(&pgCostPoolLine{}, "pool_id = ?", poolID).Error
}

// ─── CostingPeriod ─────────────────────────────────────────────────────────

type pgCostingPeriod struct {
	ID        string    `gorm:"column:id;primaryKey"`
	CompanyID string    `gorm:"column:company_id"`
	Year      int       `gorm:"column:year"`
	Month     int       `gorm:"column:month"`
	Status    string    `gorm:"column:status"`
	ClosedBy  *string   `gorm:"column:closed_by"`
	ClosedAt  *time.Time `gorm:"column:closed_at"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (pgCostingPeriod) TableName() string { return "costing_periods" }

type PGCostingPeriodRepo struct{ db *gorm.DB }

func NewPGCostingPeriodRepo(db *gorm.DB) *PGCostingPeriodRepo {
	return &PGCostingPeriodRepo{db: db}
}

func toPGCostingPeriod(cp *domain.CostingPeriod) *pgCostingPeriod {
	var created time.Time
	if cp.CreatedAt != "" {
		created, _ = time.Parse("2006-01-02T15:04:05Z", cp.CreatedAt)
	} else {
		created = time.Now()
	}
	m := &pgCostingPeriod{
		ID:        cp.ID,
		CompanyID: cp.CompanyID,
		Year:      cp.Year,
		Month:     cp.Month,
		Status:    cp.Status,
		CreatedAt: created,
	}
	if cp.ClosedBy != "" {
		m.ClosedBy = &cp.ClosedBy
	}
	if cp.ClosedAt != "" {
		t, _ := time.Parse("2006-01-02T15:04:05Z", cp.ClosedAt)
		if !t.IsZero() {
			m.ClosedAt = &t
		}
	}
	return m
}

func toDomainCostingPeriod(m *pgCostingPeriod) *domain.CostingPeriod {
	cp := &domain.CostingPeriod{
		ID:        m.ID,
		CompanyID: m.CompanyID,
		Year:      m.Year,
		Month:     m.Month,
		Status:    m.Status,
		CreatedAt: m.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if m.ClosedBy != nil {
		cp.ClosedBy = *m.ClosedBy
	}
	if m.ClosedAt != nil {
		cp.ClosedAt = m.ClosedAt.Format("2006-01-02T15:04:05Z")
	}
	return cp
}

func (r *PGCostingPeriodRepo) Create(ctx context.Context, cp *domain.CostingPeriod) error {
	return r.db.WithContext(ctx).Create(toPGCostingPeriod(cp)).Error
}

func (r *PGCostingPeriodRepo) GetByID(ctx context.Context, id string) (*domain.CostingPeriod, error) {
	var m pgCostingPeriod
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, domain.ErrCostingPeriodNotFound
	}
	return toDomainCostingPeriod(&m), nil
}

func (r *PGCostingPeriodRepo) GetByYearMonth(ctx context.Context, companyID string, year, month int) (*domain.CostingPeriod, error) {
	var m pgCostingPeriod
	if err := r.db.WithContext(ctx).Where("company_id = ? AND year = ? AND month = ?", companyID, year, month).First(&m).Error; err != nil {
		return nil, domain.ErrCostingPeriodNotFound
	}
	return toDomainCostingPeriod(&m), nil
}

func (r *PGCostingPeriodRepo) List(ctx context.Context, companyID string) ([]domain.CostingPeriod, error) {
	var rows []pgCostingPeriod
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("year DESC, month DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.CostingPeriod, len(rows))
	for i := range rows {
		out[i] = *toDomainCostingPeriod(&rows[i])
	}
	return out, nil
}

func (r *PGCostingPeriodRepo) Update(ctx context.Context, cp *domain.CostingPeriod) error {
	return r.db.WithContext(ctx).Save(toPGCostingPeriod(cp)).Error
}

// ─── CostingResult ─────────────────────────────────────────────────────────

type pgCostingResult struct {
	ID             string    `gorm:"column:id;primaryKey"`
	CompanyID      string    `gorm:"column:company_id"`
	PeriodID       string    `gorm:"column:period_id"`
	CostObjectID   string    `gorm:"column:cost_object_id"`
	CostingMethod  string    `gorm:"column:costing_method"`
	TotalDirectMat float64   `gorm:"column:total_direct_mat"`
	TotalDirectLab float64   `gorm:"column:total_direct_lab"`
	TotalOverhead  float64   `gorm:"column:total_overhead"`
	TotalCost      float64   `gorm:"column:total_cost"`
	OutputQuantity float64   `gorm:"column:output_quantity"`
	UnitCost       float64   `gorm:"column:unit_cost"`
	WIPBegin       float64   `gorm:"column:wip_begin"`
	WIPEnd         float64   `gorm:"column:wip_end"`
	Status         string    `gorm:"column:status"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (pgCostingResult) TableName() string { return "costing_results" }

type PGCostingResultRepo struct{ db *gorm.DB }

func NewPGCostingResultRepo(db *gorm.DB) *PGCostingResultRepo {
	return &PGCostingResultRepo{db: db}
}

func toPGCostingResult(cr *domain.CostingResult) *pgCostingResult {
	var created, updated time.Time
	if cr.CreatedAt != "" {
		created, _ = time.Parse("2006-01-02T15:04:05Z", cr.CreatedAt)
	} else {
		created = time.Now()
	}
	if cr.UpdatedAt != "" {
		updated, _ = time.Parse("2006-01-02T15:04:05Z", cr.UpdatedAt)
	} else {
		updated = time.Now()
	}
	return &pgCostingResult{
		ID:             cr.ID,
		CompanyID:      cr.CompanyID,
		PeriodID:       cr.PeriodID,
		CostObjectID:   cr.CostObjectID,
		CostingMethod:  string(cr.CostingMethod),
		TotalDirectMat: cr.TotalDirectMat,
		TotalDirectLab: cr.TotalDirectLab,
		TotalOverhead:  cr.TotalOverhead,
		TotalCost:      cr.TotalCost,
		OutputQuantity: cr.OutputQuantity,
		UnitCost:       cr.UnitCost,
		WIPBegin:       cr.WIPBegin,
		WIPEnd:         cr.WIPEnd,
		Status:         cr.Status,
		CreatedAt:      created,
		UpdatedAt:      updated,
	}
}

func toDomainCostingResult(m *pgCostingResult) *domain.CostingResult {
	return &domain.CostingResult{
		ID:             m.ID,
		CompanyID:      m.CompanyID,
		PeriodID:       m.PeriodID,
		CostObjectID:   m.CostObjectID,
		CostingMethod:  domain.CostingMethod(m.CostingMethod),
		TotalDirectMat: m.TotalDirectMat,
		TotalDirectLab: m.TotalDirectLab,
		TotalOverhead:  m.TotalOverhead,
		TotalCost:      m.TotalCost,
		OutputQuantity: m.OutputQuantity,
		UnitCost:       m.UnitCost,
		WIPBegin:       m.WIPBegin,
		WIPEnd:         m.WIPEnd,
		Status:         m.Status,
		CreatedAt:      m.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:      m.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func (r *PGCostingResultRepo) Create(ctx context.Context, cr *domain.CostingResult) error {
	return r.db.WithContext(ctx).Create(toPGCostingResult(cr)).Error
}

func (r *PGCostingResultRepo) GetByID(ctx context.Context, id string) (*domain.CostingResult, error) {
	var m pgCostingResult
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, domain.ErrCostingResultNotFound
	}
	return toDomainCostingResult(&m), nil
}

func (r *PGCostingResultRepo) ListByPeriod(ctx context.Context, companyID, periodID string) ([]domain.CostingResult, error) {
	var rows []pgCostingResult
	if err := r.db.WithContext(ctx).Where("company_id = ? AND period_id = ?", companyID, periodID).Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.CostingResult, len(rows))
	for i := range rows {
		out[i] = *toDomainCostingResult(&rows[i])
	}
	return out, nil
}

func (r *PGCostingResultRepo) Update(ctx context.Context, cr *domain.CostingResult) error {
	return r.db.WithContext(ctx).Save(toPGCostingResult(cr)).Error
}

func (r *PGCostingResultRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&pgCostingResult{}, "id = ?", id).Error
}

// ─── CostingResultLine ─────────────────────────────────────────────────────

type pgCostingResultLine struct {
	ID              string    `gorm:"column:id;primaryKey"`
	ResultID        string    `gorm:"column:result_id"`
	CostCategory    string    `gorm:"column:cost_category"`
	GLAccountCode   string    `gorm:"column:gl_account_code"`
	Description     string    `gorm:"column:description"`
	PlannedAmount   float64   `gorm:"column:planned_amount"`
	ActualAmount    float64   `gorm:"column:actual_amount"`
	AllocatedAmount float64   `gorm:"column:allocated_amount"`
	Coefficient     float64   `gorm:"column:coefficient"`
	CreatedAt       time.Time `gorm:"column:created_at"`
}

func (pgCostingResultLine) TableName() string { return "costing_result_lines" }

type PGCostingResultLineRepo struct{ db *gorm.DB }

func NewPGCostingResultLineRepo(db *gorm.DB) *PGCostingResultLineRepo {
	return &PGCostingResultLineRepo{db: db}
}

func toPGCostingResultLine(l *domain.CostingResultLine) *pgCostingResultLine {
	var created time.Time
	if l.CreatedAt != "" {
		created, _ = time.Parse("2006-01-02T15:04:05Z", l.CreatedAt)
	} else {
		created = time.Now()
	}
	return &pgCostingResultLine{
		ID:              l.ID,
		ResultID:        l.ResultID,
		CostCategory:    l.CostCategory,
		GLAccountCode:   l.GLAccountCode,
		Description:     l.Description,
		PlannedAmount:   l.PlannedAmount,
		ActualAmount:    l.ActualAmount,
		AllocatedAmount: l.AllocatedAmount,
		Coefficient:     l.Coefficient,
		CreatedAt:       created,
	}
}

func toDomainCostingResultLine(m *pgCostingResultLine) *domain.CostingResultLine {
	return &domain.CostingResultLine{
		ID:              m.ID,
		ResultID:        m.ResultID,
		CostCategory:    m.CostCategory,
		GLAccountCode:   m.GLAccountCode,
		Description:     m.Description,
		PlannedAmount:   m.PlannedAmount,
		ActualAmount:    m.ActualAmount,
		AllocatedAmount: m.AllocatedAmount,
		Coefficient:     m.Coefficient,
		CreatedAt:       m.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func (r *PGCostingResultLineRepo) Create(ctx context.Context, line *domain.CostingResultLine) error {
	return r.db.WithContext(ctx).Create(toPGCostingResultLine(line)).Error
}

func (r *PGCostingResultLineRepo) ListByResult(ctx context.Context, resultID string) ([]domain.CostingResultLine, error) {
	var rows []pgCostingResultLine
	if err := r.db.WithContext(ctx).Where("result_id = ?", resultID).Order("created_at").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.CostingResultLine, len(rows))
	for i := range rows {
		out[i] = *toDomainCostingResultLine(&rows[i])
	}
	return out, nil
}

func (r *PGCostingResultLineRepo) DeleteByResult(ctx context.Context, resultID string) error {
	return r.db.WithContext(ctx).Delete(&pgCostingResultLine{}, "result_id = ?", resultID).Error
}
