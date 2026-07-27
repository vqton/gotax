package service

import (
	"context"
	"errors"
	"strconv"

	"gotax/internal/domain"
)

// ─── CRUD ──────────────────────────────────────────────────────────

func (s *service) CreateOpeningBalance(ctx context.Context, ob *domain.OpeningBalance) error {
	if err := ob.Validate(); err != nil {
		return err
	}
	return s.ob.Create(ctx, ob)
}

func (s *service) GetOpeningBalance(ctx context.Context, id string) (*domain.OpeningBalance, error) {
	return s.ob.GetByID(ctx, id)
}

func (s *service) ListOpeningBalances(ctx context.Context, filter domain.OBListFilter) ([]domain.OpeningBalance, error) {
	return s.ob.List(ctx, filter)
}

func (s *service) UpdateOpeningBalance(ctx context.Context, ob *domain.OpeningBalance) error {
	existing, err := s.ob.GetByID(ctx, ob.ID)
	if err != nil {
		return err
	}
	if !existing.CanEdit() {
		return domain.ErrOpeningBalanceImmutable
	}
	if err := ob.Validate(); err != nil {
		return err
	}
	return s.ob.Update(ctx, ob)
}

func (s *service) DeleteOpeningBalance(ctx context.Context, id string) error {
	existing, err := s.ob.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !existing.CanEdit() {
		return domain.ErrOpeningBalanceImmutable
	}
	return s.ob.Delete(ctx, id)
}

// ─── Status lifecycle ──────────────────────────────────────────────

func (s *service) SubmitOpeningBalance(ctx context.Context, id, userID string) error {
	existing, err := s.ob.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !existing.CanSubmit() {
		return domain.ErrOpeningBalanceAlreadySubmitted
	}
	return s.ob.UpdateStatus(ctx, id, domain.OBStatusPending, userID)
}

func (s *service) ApproveOpeningBalance(ctx context.Context, id, approverID string) error {
	existing, err := s.ob.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing.Status != domain.OBStatusPending {
		return domain.ErrOpeningBalanceAlreadySubmitted
	}
	return s.ob.UpdateStatus(ctx, id, domain.OBStatusApproved, approverID)
}

func (s *service) CorrectOpeningBalance(ctx context.Context, id, correctedBy, reason string) (*domain.OpeningBalance, error) {
	original, err := s.ob.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if original.Status != domain.OBStatusApproved {
		return nil, domain.ErrOpeningBalanceImmutable
	}
	if reason == "" {
		return nil, domain.ErrCorrectionReasonRequired
	}
	if err := s.ob.UpdateStatus(ctx, id, domain.OBStatusCorrected, correctedBy); err != nil {
		return nil, err
	}
	correction := *original
	correction.ID = ""
	correction.Status = domain.OBStatusDraft
	correction.ApprovedBy = ""
	correction.ApprovedAt = nil
	correction.CorrectionOf = id
	correction.CorrectionReason = reason
	correction.CorrectedBy = correctedBy
	correction.CreatedBy = correctedBy
	correction.CreatedAt = s.now()
	correction.UpdatedAt = s.now()
	if err := s.ob.Create(ctx, &correction); err != nil {
		return nil, err
	}
	return &correction, nil
}

// ─── Bulk Operations ───────────────────────────────────────────────

func (s *service) BulkCreateOpeningBalances(ctx context.Context, balances []domain.OpeningBalance) error {
	for i := range balances {
		if err := balances[i].Validate(); err != nil {
			return err
		}
	}
	return s.ob.BulkCreate(ctx, balances)
}

func (s *service) BulkSubmitOpeningBalances(ctx context.Context, ids []string, userID string) error {
	for _, id := range ids {
		existing, err := s.ob.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if !existing.CanSubmit() {
			return domain.ErrOpeningBalanceAlreadySubmitted
		}
	}
	return s.ob.BulkUpdateStatus(ctx, ids, domain.OBStatusPending, userID)
}

func (s *service) BulkApproveOpeningBalances(ctx context.Context, ids []string, approverID string) error {
	for _, id := range ids {
		existing, err := s.ob.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if existing.Status != domain.OBStatusPending {
			return domain.ErrOpeningBalanceAlreadySubmitted
		}
	}
	return s.ob.BulkUpdateStatus(ctx, ids, domain.OBStatusApproved, approverID)
}

// ─── Details ───────────────────────────────────────────────────────

func (s *service) CreateOpeningBalanceDetail(ctx context.Context, d *domain.OpeningBalanceDetail) error {
	if err := d.Validate(); err != nil {
		return err
	}
	existing, err := s.ob.GetByID(ctx, d.OpeningBalanceID)
	if err != nil {
		return err
	}
	if !existing.CanEdit() {
		return domain.ErrOpeningBalanceImmutable
	}
	return s.ob.CreateDetail(ctx, d)
}

func (s *service) GetOpeningBalanceDetails(ctx context.Context, balanceID string) ([]domain.OpeningBalanceDetail, error) {
	return s.ob.GetDetails(ctx, balanceID)
}

func (s *service) DeleteOpeningBalanceDetail(ctx context.Context, id string) error {
	return s.ob.DeleteDetail(ctx, id)
}

// ─── Totals & balance ──────────────────────────────────────────────

func (s *service) GetOpeningBalanceTotals(ctx context.Context, companyID, periodID string) (float64, float64, error) {
	return s.ob.GetTotals(ctx, companyID, periodID)
}

func (s *service) ValidateOpeningBalancesBalanced(ctx context.Context, companyID, periodID string) (bool, error) {
	return s.ob.ValidateBalanced(ctx, companyID, periodID)
}

// ─── Carry Forward ─────────────────────────────────────────────────

func (s *service) CarryForward(ctx context.Context, companyID, fromPeriodID, toPeriodID, fromFiscalYear, toFiscalYear, executedBy string) (*domain.CarryForwardLog, error) {
	approvedBalances, err := s.ob.List(ctx, domain.OBListFilter{
		CompanyID: companyID,
		PeriodID:  fromPeriodID,
		Status:    domain.OBStatusApproved,
	})
	if err != nil {
		return nil, err
	}
	if len(approvedBalances) == 0 {
		return nil, domain.ErrOpeningBalanceNotFound
	}

	var carryBalances []domain.OpeningBalance
	for _, b := range approvedBalances {
		nb := domain.OpeningBalance{
			CompanyID:    companyID,
			PeriodID:     toPeriodID,
			FiscalYearID: toFiscalYear,
			AccountCode:  b.AccountCode,
			CurrencyCode: b.CurrencyCode,
			OriginalAmount: b.OriginalAmount,
			DebitAmount:  b.DebitAmount,
			CreditAmount: b.CreditAmount,
			ExchangeRate: b.ExchangeRate,
			Status:       domain.OBStatusApproved,
			SourceType:   "CARRY_FORWARD",
			Reason:       "Auto carry-forward from " + fromPeriodID,
			CreatedBy:    executedBy,
		}
		carryBalances = append(carryBalances, nb)
	}
	if err := s.ob.BulkCreate(ctx, carryBalances); err != nil {
		return nil, err
	}

	var totalDebit, totalCredit float64
	for _, b := range carryBalances {
		totalDebit += b.DebitAmount
		totalCredit += b.CreditAmount
	}

	fyFrom, _ := strconv.Atoi(fromFiscalYear)
	fyTo, _ := strconv.Atoi(toFiscalYear)

	log := &domain.CarryForwardLog{
		CompanyID:     companyID,
		FromPeriodID:  fromPeriodID,
		ToPeriodID:    toPeriodID,
		FromFiscalYear: fyFrom,
		ToFiscalYear:  fyTo,
		AccountCount:  len(carryBalances),
		TotalDebit:    totalDebit,
		TotalCredit:   totalCredit,
		Status:        "COMPLETED",
		ExecutedBy:    executedBy,
	}
	if err := s.ob.CreateCarryForwardLog(ctx, log); err != nil {
		return nil, err
	}
	return log, nil
}

func (s *service) GetCarryForwardLogs(ctx context.Context, companyID string) ([]domain.CarryForwardLog, error) {
	return s.ob.GetCarryForwardLogs(ctx, companyID)
}

func (s *service) GetCarryForwardLogByID(ctx context.Context, id string) (*domain.CarryForwardLog, error) {
	return s.ob.GetCarryForwardLogByID(ctx, id)
}

// ─── Circular 99 Mapping ──────────────────────────────────────────

func (s *service) CreateCircular99Mapping(ctx context.Context, m *domain.Circular99Mapping) error {
	existing, err := s.ob.GetCircular99MappingByOldCode(ctx, m.OldAccountCode)
	if err != nil && !errors.Is(err, domain.ErrCircular99MappingNotFound) {
		return err
	}
	if existing != nil {
		return domain.ErrCircular99MappingExists
	}
	return s.ob.CreateCircular99Mapping(ctx, m)
}

func (s *service) ListCircular99Mappings(ctx context.Context) ([]domain.Circular99Mapping, error) {
	return s.ob.ListCircular99Mappings(ctx)
}

func (s *service) GetCircular99MappingByOldCode(ctx context.Context, oldCode string) (*domain.Circular99Mapping, error) {
	return s.ob.GetCircular99MappingByOldCode(ctx, oldCode)
}

// ─── Balance Migration ─────────────────────────────────────────────

func (s *service) CreateBalanceMigration(ctx context.Context, m *domain.BalanceMigration) error {
	if m.Status == "" {
		m.Status = "PENDING"
	}
	return s.ob.CreateMigration(ctx, m)
}

func (s *service) GetBalanceMigrationByID(ctx context.Context, id string) (*domain.BalanceMigration, error) {
	return s.ob.GetMigrationByID(ctx, id)
}

func (s *service) ListBalanceMigrations(ctx context.Context, companyID string) ([]domain.BalanceMigration, error) {
	return s.ob.ListMigrations(ctx, companyID)
}
