package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gotax/internal/domain"
)

// ─── Contract ──────────────────────────────────────────────────────

type pgContract struct {
	ID                  string    `gorm:"column:id;primaryKey"`
	CompanyID           string    `gorm:"column:company_id"`
	Code                string    `gorm:"column:code"`
	Name                string    `gorm:"column:name"`
	ContractType        string    `gorm:"column:contract_type"`
	Status              string    `gorm:"column:status"`
	Value               float64   `gorm:"column:value"`
	StartDate           *time.Time `gorm:"column:start_date"`
	EndDate             *time.Time `gorm:"column:end_date"`
	CounterpartyName    *string   `gorm:"column:counterparty_name"`
	CounterpartyTaxCode *string   `gorm:"column:counterparty_tax_code"`
	Notes               *string   `gorm:"column:notes"`
	CreatedAt           time.Time `gorm:"column:created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at"`
}

func (pgContract) TableName() string { return "contracts" }

type PGContractRepo struct{ db *gorm.DB }

func NewPGContractRepo(db *gorm.DB) *PGContractRepo {
	return &PGContractRepo{db: db}
}

func toPGContract(c *domain.Contract) *pgContract {
	var created, updated time.Time
	if c.CreatedAt != "" {
		created, _ = time.Parse("2006-01-02T15:04:05Z", c.CreatedAt)
	} else {
		created = time.Now()
	}
	if c.UpdatedAt != "" {
		updated, _ = time.Parse("2006-01-02T15:04:05Z", c.UpdatedAt)
	} else {
		updated = time.Now()
	}
	m := &pgContract{
		ID:           c.ID,
		CompanyID:    c.CompanyID,
		Code:         c.Code,
		Name:         c.Name,
		ContractType: c.ContractType,
		Status:       c.Status,
		Value:        c.Value,
		CreatedAt:    created,
		UpdatedAt:    updated,
	}
	if c.StartDate != "" {
		t, _ := time.Parse("2006-01-02", c.StartDate)
		m.StartDate = &t
	}
	if c.EndDate != "" {
		t, _ := time.Parse("2006-01-02", c.EndDate)
		m.EndDate = &t
	}
	if c.CounterpartyName != "" {
		m.CounterpartyName = &c.CounterpartyName
	}
	if c.CounterpartyTaxCode != "" {
		m.CounterpartyTaxCode = &c.CounterpartyTaxCode
	}
	if c.Notes != "" {
		m.Notes = &c.Notes
	}
	return m
}

func toDomainContract(m *pgContract) *domain.Contract {
	c := &domain.Contract{
		ID:           m.ID,
		CompanyID:    m.CompanyID,
		Code:         m.Code,
		Name:         m.Name,
		ContractType: m.ContractType,
		Status:       m.Status,
		Value:        m.Value,
		CreatedAt:    m.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    m.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if m.StartDate != nil {
		c.StartDate = m.StartDate.Format("2006-01-02")
	}
	if m.EndDate != nil {
		c.EndDate = m.EndDate.Format("2006-01-02")
	}
	if m.CounterpartyName != nil {
		c.CounterpartyName = *m.CounterpartyName
	}
	if m.CounterpartyTaxCode != nil {
		c.CounterpartyTaxCode = *m.CounterpartyTaxCode
	}
	if m.Notes != nil {
		c.Notes = *m.Notes
	}
	return c
}

func (r *PGContractRepo) Create(ctx context.Context, c *domain.Contract) error {
	m := toPGContract(c)
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	if m.Status == "" {
		m.Status = "draft"
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGContractRepo) GetByID(ctx context.Context, id string) (*domain.Contract, error) {
	var m pgContract
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, domain.ErrContractNotFound
	}
	return toDomainContract(&m), nil
}

func (r *PGContractRepo) GetByCode(ctx context.Context, companyID, code string) (*domain.Contract, error) {
	var m pgContract
	err := r.db.WithContext(ctx).Where("company_id = ? AND code = ?", companyID, code).First(&m).Error
	if err != nil {
		return nil, domain.ErrContractNotFound
	}
	return toDomainContract(&m), nil
}

func (r *PGContractRepo) List(ctx context.Context, companyID string) ([]domain.Contract, error) {
	var rows []pgContract
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Contract, len(rows))
	for i, row := range rows {
		out[i] = *toDomainContract(&row)
	}
	return out, nil
}

func (r *PGContractRepo) Update(ctx context.Context, c *domain.Contract) error {
	m := toPGContract(c)
	m.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *PGContractRepo) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&pgContract{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return domain.ErrContractNotFound
	}
	return result.Error
}

// ─── ContractPayment ───────────────────────────────────────────────

type pgContractPayment struct {
	ID          string    `gorm:"column:id;primaryKey"`
	ContractID  string    `gorm:"column:contract_id"`
	PaymentDate time.Time `gorm:"column:payment_date"`
	Amount      float64   `gorm:"column:amount"`
	Description *string   `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (pgContractPayment) TableName() string { return "contract_payments" }

type PGContractPaymentRepo struct{ db *gorm.DB }

func NewPGContractPaymentRepo(db *gorm.DB) *PGContractPaymentRepo {
	return &PGContractPaymentRepo{db: db}
}

func (r *PGContractPaymentRepo) Create(ctx context.Context, p *domain.ContractPayment) error {
	m := &pgContractPayment{
		ID:         p.ID,
		ContractID: p.ContractID,
		Amount:     p.Amount,
		CreatedAt:  time.Now(),
	}
	if p.PaymentDate != "" {
		t, _ := time.Parse("2006-01-02", p.PaymentDate)
		m.PaymentDate = t
	}
	if p.Description != "" {
		m.Description = &p.Description
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGContractPaymentRepo) ListByContract(ctx context.Context, contractID string) ([]domain.ContractPayment, error) {
	var rows []pgContractPayment
	if err := r.db.WithContext(ctx).Where("contract_id = ?", contractID).Order("payment_date ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.ContractPayment, len(rows))
	for i, row := range rows {
		cp := domain.ContractPayment{
			ID:         row.ID,
			ContractID: row.ContractID,
			PaymentDate: row.PaymentDate.Format("2006-01-02"),
			Amount:     row.Amount,
			CreatedAt:  row.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
		if row.Description != nil {
			cp.Description = *row.Description
		}
		out[i] = cp
	}
	return out, nil
}

func (r *PGContractPaymentRepo) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&pgContractPayment{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return domain.ErrContractNotFound
	}
	return result.Error
}
