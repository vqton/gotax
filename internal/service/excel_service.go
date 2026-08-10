package service

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"

	"gotax/internal/domain"
)

type ExcelImportService struct {
	supplierRepo domain.SupplierRepository
	customerRepo domain.CustomerRepository
	itemRepo     domain.ItemRepository
}

func NewExcelImportService(supplierRepo domain.SupplierRepository, customerRepo domain.CustomerRepository, itemRepo domain.ItemRepository) *ExcelImportService {
	return &ExcelImportService{supplierRepo: supplierRepo, customerRepo: customerRepo, itemRepo: itemRepo}
}

type ImportResult struct {
	Imported int      `json:"imported"`
	Skipped  int      `json:"skipped"`
	Errors   []string `json:"errors,omitempty"`
}

func (s *ExcelImportService) ImportSuppliers(ctx context.Context, companyID string, data []byte) (*ImportResult, error) {
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("invalid xlsx: %w", err)
	}
	defer f.Close()

	sheet := f.GetSheetName(0)
	if sheet == "" {
		return nil, fmt.Errorf("no sheets found")
	}
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}
	result := &ImportResult{}
	if len(rows) < 2 {
		return result, nil
	}
	for i, row := range rows[1:] {
		if len(row) < 2 {
			continue
		}
		supplier := &domain.Supplier{
			CompanyID: companyID,
			Code:      cellStr(row, 0),
			Name:      cellStr(row, 1),
			TaxCode:   cellStr(row, 2),
			Address:   cellStr(row, 3),
			Phone:     cellStr(row, 4),
			Email:     cellStr(row, 5),
		}
		if supplier.Code == "" || supplier.Name == "" {
			result.Skipped++
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: missing code or name", i+2))
			continue
		}
		if len(row) > 6 {
			supplier.BankName = cellStr(row, 6)
		}
		if len(row) > 7 {
			supplier.BankAccountNumber = cellStr(row, 7)
		}
		if err := s.supplierRepo.CreateSupplier(ctx, supplier); err != nil {
			result.Skipped++
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: %v", i+2, err))
			continue
		}
		result.Imported++
	}
	return result, nil
}

func (s *ExcelImportService) ImportCustomers(ctx context.Context, companyID string, data []byte) (*ImportResult, error) {
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("invalid xlsx: %w", err)
	}
	defer f.Close()
	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}
	result := &ImportResult{}
	if len(rows) < 2 {
		return result, nil
	}
	for i, row := range rows[1:] {
		if len(row) < 2 {
			continue
		}
		customer := &domain.Customer{
			CompanyID: companyID,
			Code:      cellStr(row, 0),
			Name:      cellStr(row, 1),
			TaxCode:   cellStr(row, 2),
			Address:   cellStr(row, 3),
			Phone:     cellStr(row, 4),
			Email:     cellStr(row, 5),
		}
		if customer.Code == "" || customer.Name == "" {
			result.Skipped++
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: missing code or name", i+2))
			continue
		}
		if len(row) > 6 {
			customer.BankName = cellStr(row, 6)
		}
		if len(row) > 7 {
			customer.BankAccountNumber = cellStr(row, 7)
		}
		if err := s.customerRepo.CreateCustomer(ctx, customer); err != nil {
			result.Skipped++
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: %v", i+2, err))
			continue
		}
		result.Imported++
	}
	return result, nil
}

func (s *ExcelImportService) ImportItems(ctx context.Context, companyID string, data []byte) (*ImportResult, error) {
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("invalid xlsx: %w", err)
	}
	defer f.Close()
	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}
	result := &ImportResult{}
	if len(rows) < 2 {
		return result, nil
	}
	for i, row := range rows[1:] {
		if len(row) < 2 {
			continue
		}
		item := &domain.Item{
			CompanyID: companyID,
			Code:      cellStr(row, 0),
			Name:      cellStr(row, 1),
		}
		if item.Code == "" || item.Name == "" {
			result.Skipped++
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: missing code or name", i+2))
			continue
		}
		if len(row) > 2 {
			item.Unit = cellStr(row, 2)
		}
		if len(row) > 3 {
			item.CategoryID = cellStr(row, 3)
		}
		if err := s.itemRepo.CreateItem(ctx, item); err != nil {
			result.Skipped++
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: %v", i+2, err))
			continue
		}
		result.Imported++
	}
	return result, nil
}

func cellStr(row []string, idx int) string {
	if idx >= len(row) {
		return ""
	}
	return row[idx]
}

func parseDate(s string) time.Time {
	for _, layout := range []string{"2006-01-02", "02/01/2006", "02-01-2006"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
