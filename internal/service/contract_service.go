package service

import (
	"context"
	"fmt"
	"time"

	"gotax/internal/domain"
)

type ContractService struct {
	contractRepo domain.ContractRepository
	paymentRepo  domain.ContractPaymentRepository
}

func NewContractService(contractRepo domain.ContractRepository, paymentRepo domain.ContractPaymentRepository) *ContractService {
	return &ContractService{contractRepo: contractRepo, paymentRepo: paymentRepo}
}

func (s *ContractService) Create(ctx context.Context, c *domain.Contract) error {
	if c.CompanyID == "" {
		return fmt.Errorf("company_id is required")
	}
	if c.Code == "" {
		return fmt.Errorf("code is required")
	}
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	existing, err := s.contractRepo.GetByCode(ctx, c.CompanyID, c.Code)
	if err == nil && existing != nil {
		return domain.ErrContractExists
	}
	if c.ID == "" {
		c.ID = fmt.Sprintf("CON-%d", time.Now().UnixNano())
	}
	now := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	c.CreatedAt = now
	c.UpdatedAt = now
	if c.Status == "" {
		c.Status = "draft"
	}
	return s.contractRepo.Create(ctx, c)
}

func (s *ContractService) GetByID(ctx context.Context, id string) (*domain.Contract, error) {
	return s.contractRepo.GetByID(ctx, id)
}

func (s *ContractService) List(ctx context.Context, companyID string) ([]domain.Contract, error) {
	return s.contractRepo.List(ctx, companyID)
}

func (s *ContractService) Update(ctx context.Context, c *domain.Contract) error {
	if c.CompanyID == "" {
		return fmt.Errorf("company_id is required")
	}
	existing, err := s.contractRepo.GetByID(ctx, c.ID)
	if err != nil {
		return err
	}
	c.CreatedAt = existing.CreatedAt
	c.UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	return s.contractRepo.Update(ctx, c)
}

func (s *ContractService) Delete(ctx context.Context, id string) error {
	return s.contractRepo.Delete(ctx, id)
}

func (s *ContractService) AddPayment(ctx context.Context, p *domain.ContractPayment) error {
	if p.ContractID == "" {
		return fmt.Errorf("contract_id is required")
	}
	if p.Amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if p.ID == "" {
		p.ID = fmt.Sprintf("CP-%d", time.Now().UnixNano())
	}
	return s.paymentRepo.Create(ctx, p)
}

func (s *ContractService) ListPayments(ctx context.Context, contractID string) ([]domain.ContractPayment, error) {
	return s.paymentRepo.ListByContract(ctx, contractID)
}
