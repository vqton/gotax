package service

import (
	"context"
	"fmt"

	"gotax/internal/domain"
	"gotax/internal/validate"
)

type FAServiceInterface interface {
	CreateCategory(ctx context.Context, c *domain.FixedAssetCategory) error
	GetCategory(ctx context.Context, id string) (*domain.FixedAssetCategory, error)
	GetCategoryByCode(ctx context.Context, companyID, code string) (*domain.FixedAssetCategory, error)
	ListCategories(ctx context.Context, filter domain.FACategoryFilter) ([]domain.FixedAssetCategory, error)
	UpdateCategory(ctx context.Context, c *domain.FixedAssetCategory) error
	DeleteCategory(ctx context.Context, id string) error

	CreateAsset(ctx context.Context, a *domain.FixedAsset) error
	GetAsset(ctx context.Context, id string) (*domain.FixedAsset, error)
	GetAssetByCode(ctx context.Context, companyID, code string) (*domain.FixedAsset, error)
	ListAssets(ctx context.Context, filter domain.FAListFilter) ([]domain.FixedAsset, int, error)
	UpdateAsset(ctx context.Context, a *domain.FixedAsset) error
	DeleteAsset(ctx context.Context, id string) error
	ChangeAssetStatus(ctx context.Context, id string, status domain.FixedAssetStatus, userID string) error

	RunDepreciation(ctx context.Context, input domain.FARunDepreciationInput) ([]DepreciationResult, error)
	RunDepreciationForAsset(ctx context.Context, input domain.FARunDepreciationInput, assetID string) (*DepreciationResult, error)
	GetDepreciationByAsset(ctx context.Context, assetID string) ([]domain.DepreciationEntry, error)
	GetDepreciationByPeriod(ctx context.Context, periodID string) ([]domain.DepreciationEntry, error)

	RecordTransaction(ctx context.Context, t *domain.FixedAssetTransaction) error
	GetTransactions(ctx context.Context, assetID string) ([]domain.FixedAssetTransaction, error)

	SetAllocations(ctx context.Context, assetID string, allocs []domain.FixedAssetAllocation) error
	GetAllocations(ctx context.Context, assetID string) ([]domain.FixedAssetAllocation, error)

	DisposeAsset(ctx context.Context, input domain.FADisposalInput) error
	SellAsset(ctx context.Context, input domain.FASaleInput) error
	AdjustAsset(ctx context.Context, input domain.FAAdjustmentInput) error
	RevalueAsset(ctx context.Context, input domain.FARevaluationInput) error
	ImpairAsset(ctx context.Context, input domain.FAImpairmentInput) error
	TransferAsset(ctx context.Context, input domain.FATransferInput) error
	CIPTransfer(ctx context.Context, input domain.FACIPTransferInput) error
	SuspendDepreciation(ctx context.Context, input domain.FASuspendInput) error
	ResumeDepreciation(ctx context.Context, input domain.FAResumeInput) error

	CreateInventoryPlan(ctx context.Context, p *domain.FixedAssetInventoryPlan) error
	GetInventoryPlan(ctx context.Context, id string) (*domain.FixedAssetInventoryPlan, error)
	ListInventoryPlans(ctx context.Context, companyID string) ([]domain.FixedAssetInventoryPlan, error)
	UpdateInventoryPlan(ctx context.Context, p *domain.FixedAssetInventoryPlan) error

	CreateInventoryResult(ctx context.Context, r *domain.FixedAssetInventoryResult) error
	GetInventoryResults(ctx context.Context, planID string) ([]domain.FixedAssetInventoryResult, error)
}

type faService struct {
	faRepo domain.FARepository
	engine *DepreciationEngine
}

func NewFAService(faRepo domain.FARepository) FAServiceInterface {
	return &faService{
		faRepo: faRepo,
		engine: NewDepreciationEngine(faRepo),
	}
}

// ─── Categories ──────────────────────────────────────────────────────

func (s *faService) CreateCategory(ctx context.Context, c *domain.FixedAssetCategory) error {
	if err := validate.FixedAssetCategory(c); err != nil {
		return err
	}
	return s.faRepo.CreateCategory(ctx, c)
}

func (s *faService) GetCategory(ctx context.Context, id string) (*domain.FixedAssetCategory, error) {
	return s.faRepo.GetCategoryByID(ctx, id)
}

func (s *faService) GetCategoryByCode(ctx context.Context, companyID, code string) (*domain.FixedAssetCategory, error) {
	return s.faRepo.GetCategoryByCode(ctx, companyID, code)
}

func (s *faService) ListCategories(ctx context.Context, filter domain.FACategoryFilter) ([]domain.FixedAssetCategory, error) {
	return s.faRepo.ListCategories(ctx, filter)
}

func (s *faService) UpdateCategory(ctx context.Context, c *domain.FixedAssetCategory) error {
	if err := validate.FixedAssetCategory(c); err != nil {
		return err
	}
	return s.faRepo.UpdateCategory(ctx, c)
}

func (s *faService) DeleteCategory(ctx context.Context, id string) error {
	return s.faRepo.DeleteCategory(ctx, id)
}

// ─── Fixed Assets ────────────────────────────────────────────────────

func (s *faService) CreateAsset(ctx context.Context, a *domain.FixedAsset) error {
	if err := validate.FixedAsset(a); err != nil {
		return err
	}
	a.CarryingAmount = a.CalcCarryingAmount()
	return s.faRepo.CreateAsset(ctx, a)
}

func (s *faService) GetAsset(ctx context.Context, id string) (*domain.FixedAsset, error) {
	return s.faRepo.GetAssetByID(ctx, id)
}

func (s *faService) GetAssetByCode(ctx context.Context, companyID, code string) (*domain.FixedAsset, error) {
	return s.faRepo.GetAssetByCode(ctx, companyID, code)
}

func (s *faService) ListAssets(ctx context.Context, filter domain.FAListFilter) ([]domain.FixedAsset, int, error) {
	return s.faRepo.ListAssets(ctx, filter)
}

func (s *faService) UpdateAsset(ctx context.Context, a *domain.FixedAsset) error {
	if err := validate.FixedAsset(a); err != nil {
		return err
	}
	a.CarryingAmount = a.CalcCarryingAmount()
	return s.faRepo.UpdateAsset(ctx, a)
}

func (s *faService) DeleteAsset(ctx context.Context, id string) error {
	return s.faRepo.DeleteAsset(ctx, id)
}

func (s *faService) ChangeAssetStatus(ctx context.Context, id string, status domain.FixedAssetStatus, userID string) error {
	a, err := s.faRepo.GetAssetByID(ctx, id)
	if err != nil {
		return err
	}
	a.Status = status
	return s.faRepo.UpdateAsset(ctx, a)
}

// ─── Depreciation ────────────────────────────────────────────────────

func (s *faService) RunDepreciation(ctx context.Context, input domain.FARunDepreciationInput) ([]DepreciationResult, error) {
	return s.engine.RunForCompany(ctx, input.CompanyID, input)
}

func (s *faService) RunDepreciationForAsset(ctx context.Context, input domain.FARunDepreciationInput, assetID string) (*DepreciationResult, error) {
	a, err := s.faRepo.GetAssetByID(ctx, assetID)
	if err != nil {
		return nil, err
	}
	return s.engine.RunForAsset(ctx, a, input)
}

func (s *faService) GetDepreciationByAsset(ctx context.Context, assetID string) ([]domain.DepreciationEntry, error) {
	return s.faRepo.ListDepreciationByAsset(ctx, assetID)
}

func (s *faService) GetDepreciationByPeriod(ctx context.Context, periodID string) ([]domain.DepreciationEntry, error) {
	return s.faRepo.ListDepreciationByPeriod(ctx, periodID)
}

// ─── Transactions ────────────────────────────────────────────────────

func (s *faService) RecordTransaction(ctx context.Context, t *domain.FixedAssetTransaction) error {
	return s.faRepo.CreateTransaction(ctx, t)
}

func (s *faService) GetTransactions(ctx context.Context, assetID string) ([]domain.FixedAssetTransaction, error) {
	return s.faRepo.ListTransactionsByAsset(ctx, assetID)
}

// ─── Allocations ─────────────────────────────────────────────────────

func (s *faService) SetAllocations(ctx context.Context, assetID string, allocs []domain.FixedAssetAllocation) error {
	var totalPct float64
	for _, a := range allocs {
		totalPct += a.AllocationPct
	}
	if len(allocs) > 0 && totalPct != 100 {
		return domain.ErrFAAllocationPctSum
	}
	return s.faRepo.SetAllocations(ctx, assetID, allocs)
}

func (s *faService) GetAllocations(ctx context.Context, assetID string) ([]domain.FixedAssetAllocation, error) {
	return s.faRepo.GetAllocations(ctx, assetID)
}

// ─── Business Operations ─────────────────────────────────────────────

func (s *faService) DisposeAsset(ctx context.Context, input domain.FADisposalInput) error {
	a, err := s.faRepo.GetAssetByID(ctx, input.FixedAssetID)
	if err != nil {
		return err
	}
	if a.Status != domain.FADepreciating && a.Status != domain.FAFullyDepr {
		return domain.ErrFANotActive
	}
	a.Status = domain.FADisposed
	if err := s.faRepo.UpdateAsset(ctx, a); err != nil {
		return err
	}
	txn := &domain.FixedAssetTransaction{
		CompanyID:       a.CompanyID,
		FixedAssetID:    a.ID,
		TransactionType: domain.FATrxDisposal,
		TransactionDate: input.DisposalDate,
		Amount:          input.Proceeds,
		Description:     input.Description,
		CreatedBy:       input.CreatedBy,
	}
	return s.faRepo.CreateTransaction(ctx, txn)
}

func (s *faService) SellAsset(ctx context.Context, input domain.FASaleInput) error {
	a, err := s.faRepo.GetAssetByID(ctx, input.FixedAssetID)
	if err != nil {
		return err
	}
	if a.Status != domain.FADepreciating && a.Status != domain.FAFullyDepr {
		return domain.ErrFANotActive
	}
	a.Status = domain.FASold
	if err := s.faRepo.UpdateAsset(ctx, a); err != nil {
		return err
	}
	txn := &domain.FixedAssetTransaction{
		CompanyID:       a.CompanyID,
		FixedAssetID:    a.ID,
		TransactionType: domain.FATrxSale,
		TransactionDate: input.SaleDate,
		Amount:          input.Proceeds,
		Description:     input.Description,
		CreatedBy:       input.CreatedBy,
	}
	return s.faRepo.CreateTransaction(ctx, txn)
}

func (s *faService) AdjustAsset(ctx context.Context, input domain.FAAdjustmentInput) error {
	a, err := s.faRepo.GetAssetByID(ctx, input.FixedAssetID)
	if err != nil {
		return err
	}
	oldCarry := a.CarryingAmount
	if input.NewOriginalCost != nil {
		a.OriginalCost = *input.NewOriginalCost
	}
	if input.NewUsefulLifeMonths != nil {
		a.UsefulLifeMonths = *input.NewUsefulLifeMonths
	}
	if input.NewDepreciationMethod != nil {
		a.DepreciationMethod = *input.NewDepreciationMethod
	}
	if input.NewResidualValue != nil {
		a.ResidualValue = *input.NewResidualValue
	}
	a.CarryingAmount = a.CalcCarryingAmount()
	if err := s.faRepo.UpdateAsset(ctx, a); err != nil {
		return err
	}
	txn := &domain.FixedAssetTransaction{
		CompanyID:       a.CompanyID,
		FixedAssetID:    a.ID,
		TransactionType: domain.FATrxAdjustment,
		TransactionDate: input.AdjustmentDate,
		Amount:          a.CarryingAmount - oldCarry,
		OldValue:        oldCarry,
		NewValue:        a.CarryingAmount,
		Description:     input.Reason,
		CreatedBy:       input.CreatedBy,
	}
	return s.faRepo.CreateTransaction(ctx, txn)
}

func (s *faService) RevalueAsset(ctx context.Context, input domain.FARevaluationInput) error {
	if input.FairValue <= 0 {
		return domain.ErrFARevaluationRequiresFairValue
	}
	a, err := s.faRepo.GetAssetByID(ctx, input.FixedAssetID)
	if err != nil {
		return err
	}
	diff := input.FairValue - a.CarryingAmount
	a.OriginalCost += diff
	a.CarryingAmount = a.CalcCarryingAmount()
	if err := s.faRepo.UpdateAsset(ctx, a); err != nil {
		return err
	}
	txn := &domain.FixedAssetTransaction{
		CompanyID:       a.CompanyID,
		FixedAssetID:    a.ID,
		TransactionType: domain.FATrxRevaluation,
		TransactionDate: input.RevaluationDate,
		Amount:          diff,
		OldValue:        a.CarryingAmount - diff,
		NewValue:        a.CarryingAmount,
		CreatedBy:       input.CreatedBy,
	}
	return s.faRepo.CreateTransaction(ctx, txn)
}

func (s *faService) ImpairAsset(ctx context.Context, input domain.FAImpairmentInput) error {
	if input.ImpairmentAmount <= 0 {
		return domain.ErrFAImpairmentAmountInvalid
	}
	a, err := s.faRepo.GetAssetByID(ctx, input.FixedAssetID)
	if err != nil {
		return err
	}
	a.AccumulatedDepreciation += input.ImpairmentAmount
	a.CarryingAmount = a.CalcCarryingAmount()
	if err := s.faRepo.UpdateAsset(ctx, a); err != nil {
		return err
	}
	txn := &domain.FixedAssetTransaction{
		CompanyID:       a.CompanyID,
		FixedAssetID:    a.ID,
		TransactionType: domain.FATrxImpairment,
		TransactionDate: input.ImpairmentDate,
		Amount:          input.ImpairmentAmount,
		OldValue:        a.CarryingAmount + input.ImpairmentAmount,
		NewValue:        a.CarryingAmount,
		Description:     input.Reason,
		CreatedBy:       input.CreatedBy,
	}
	return s.faRepo.CreateTransaction(ctx, txn)
}

func (s *faService) TransferAsset(ctx context.Context, input domain.FATransferInput) error {
	a, err := s.faRepo.GetAssetByID(ctx, input.FixedAssetID)
	if err != nil {
		return err
	}
	if a.DepartmentID == input.DepartmentID {
		return domain.ErrFATransferSameDepartment
	}
	oldDept := a.DepartmentID
	a.DepartmentID = input.DepartmentID
	if err := s.faRepo.UpdateAsset(ctx, a); err != nil {
		return err
	}
	txn := &domain.FixedAssetTransaction{
		CompanyID:       a.CompanyID,
		FixedAssetID:    a.ID,
		TransactionType: domain.FATrxTransfer,
		TransactionDate: input.EffectiveDate,
		Amount:          0,
		Description:     fmt.Sprintf("Transferred from department %s to %s", oldDept, input.DepartmentID),
		CreatedBy:       input.CreatedBy,
	}
	return s.faRepo.CreateTransaction(ctx, txn)
}

func (s *faService) CIPTransfer(ctx context.Context, input domain.FACIPTransferInput) error {
	a, err := s.faRepo.GetAssetByID(ctx, input.FixedAssetID)
	if err != nil {
		return err
	}
	if a.Source != domain.FASourceConstruction {
		return domain.ErrFACIPTransferNotApplicable
	}
	a.OriginalCost = input.TotalCost
	a.CarryingAmount = a.CalcCarryingAmount()
	a.CIPAccountID = input.CIPAccountID
	a.Status = domain.FAActive
	if err := s.faRepo.UpdateAsset(ctx, a); err != nil {
		return err
	}
	txn := &domain.FixedAssetTransaction{
		CompanyID:       a.CompanyID,
		FixedAssetID:    a.ID,
		TransactionType: domain.FATrxCIPTransfer,
		TransactionDate: input.TransferDate,
		Amount:          input.TotalCost,
		CreatedBy:       input.CreatedBy,
	}
	return s.faRepo.CreateTransaction(ctx, txn)
}

func (s *faService) SuspendDepreciation(ctx context.Context, input domain.FASuspendInput) error {
	a, err := s.faRepo.GetAssetByID(ctx, input.FixedAssetID)
	if err != nil {
		return err
	}
	if a.Status != domain.FADepreciating && a.Status != domain.FAActive {
		return domain.ErrFASuspensionNotApplicable
	}
	a.Status = domain.FASuspended
	if err := s.faRepo.UpdateAsset(ctx, a); err != nil {
		return err
	}
	txn := &domain.FixedAssetTransaction{
		CompanyID:       a.CompanyID,
		FixedAssetID:    a.ID,
		TransactionType: domain.FATrxSuspend,
		TransactionDate: input.SuspendDate,
		Description:     input.Reason,
		CreatedBy:       input.CreatedBy,
	}
	return s.faRepo.CreateTransaction(ctx, txn)
}

func (s *faService) ResumeDepreciation(ctx context.Context, input domain.FAResumeInput) error {
	a, err := s.faRepo.GetAssetByID(ctx, input.FixedAssetID)
	if err != nil {
		return err
	}
	if a.Status != domain.FASuspended {
		return domain.ErrFAResumeNotApplicable
	}
	a.Status = domain.FADepreciating
	if err := s.faRepo.UpdateAsset(ctx, a); err != nil {
		return err
	}
	txn := &domain.FixedAssetTransaction{
		CompanyID:       a.CompanyID,
		FixedAssetID:    a.ID,
		TransactionType: domain.FATrxResume,
		TransactionDate: input.ResumeDate,
		CreatedBy:       input.CreatedBy,
	}
	return s.faRepo.CreateTransaction(ctx, txn)
}

// ─── Inventory ───────────────────────────────────────────────────────

func (s *faService) CreateInventoryPlan(ctx context.Context, p *domain.FixedAssetInventoryPlan) error {
	return s.faRepo.CreateInventoryPlan(ctx, p)
}

func (s *faService) GetInventoryPlan(ctx context.Context, id string) (*domain.FixedAssetInventoryPlan, error) {
	return s.faRepo.GetInventoryPlan(ctx, id)
}

func (s *faService) ListInventoryPlans(ctx context.Context, companyID string) ([]domain.FixedAssetInventoryPlan, error) {
	return s.faRepo.ListInventoryPlans(ctx, companyID)
}

func (s *faService) UpdateInventoryPlan(ctx context.Context, p *domain.FixedAssetInventoryPlan) error {
	return s.faRepo.UpdateInventoryPlan(ctx, p)
}

func (s *faService) CreateInventoryResult(ctx context.Context, r *domain.FixedAssetInventoryResult) error {
	return s.faRepo.CreateInventoryResult(ctx, r)
}

func (s *faService) GetInventoryResults(ctx context.Context, planID string) ([]domain.FixedAssetInventoryResult, error) {
	return s.faRepo.GetInventoryResultsByPlan(ctx, planID)
}
