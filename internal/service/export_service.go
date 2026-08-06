package service

import (
	"bytes"
	"context"
	"fmt"

	"github.com/xuri/excelize/v2"

	"gotax/internal/domain"
)

type ExportService struct {
	journalRepo domain.JournalRepository
	accountRepo domain.AccountRepository
	periodRepo  domain.PeriodRepository
}

func NewExportService(journalRepo domain.JournalRepository, accountRepo domain.AccountRepository, periodRepo domain.PeriodRepository) *ExportService {
	return &ExportService{journalRepo: journalRepo, accountRepo: accountRepo, periodRepo: periodRepo}
}

// ExportJournalEntries exports posted journal entries for a given year/month to .xlsx
func (s *ExportService) ExportJournalEntries(ctx context.Context, companyID string, year, month int) ([]byte, error) {
	period, err := s.periodRepo.GetByYearMonth(ctx, year, month)
	if err != nil {
		return nil, err
	}

	entries, err := s.journalRepo.GetByPeriod(ctx, period.ID)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "Chứng từ"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"Ngày", "Số CT", "Loại CT", "Mã TK", "Diễn giải", "Nợ", "Có", "Mã đối tượng"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}
	style, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#E8F0FE"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	f.SetCellStyle(sheet, "A1", "H1", style)

	row := 2
	for _, entry := range entries {
		if entry.Status != domain.JournalEntryPosted {
			continue
		}
		for _, line := range entry.Lines {
			dateStr := entry.EntryDate.Format("2006-01-02")
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), dateStr)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), entry.EntryNumber)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), string(entry.VoucherType))
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), line.AccountCode)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), line.Description)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), line.DebitAmount)
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), line.CreditAmount)
			f.SetCellValue(sheet, fmt.Sprintf("H%d", row), line.ObjectID)
			row++
		}
	}

	for i := 1; i <= len(headers); i++ {
		col, _ := excelize.ColumnNumberToName(i)
		f.SetColWidth(sheet, col, col, 18)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ExportTrialBalance exports trial balance for a given year/month to .xlsx
func (s *ExportService) ExportTrialBalance(ctx context.Context, companyID string, year, month int) ([]byte, error) {
	balances, err := s.journalRepo.GetTrialBalance(ctx, fmt.Sprintf("%s-%d-%02d", companyID, year, month))
	if err != nil {
		return nil, err
	}
	for i := range balances {
		balances[i].Calculate()
	}

	f := excelize.NewFile()
	sheet := "Bảng cân đối"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"Mã TK", "TK tổng hợp", "Nợ đầu kỳ", "Có đầu kỳ", "Phát sinh nợ", "Phát sinh có", "Nợ cuối kỳ", "Có cuối kỳ"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}
	style, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#E8F0FE"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	f.SetCellStyle(sheet, "A1", "H1", style)

	row := 2
	for _, b := range balances {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), b.AccountCode)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), string(b.AccountType))
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), b.OpenBalanceDebit)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), b.OpenBalanceCredit)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), b.PeriodDebit)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), b.PeriodCredit)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), b.TotalDebit)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), b.TotalCredit)
		row++
	}

	for i := 1; i <= len(headers); i++ {
		col, _ := excelize.ColumnNumberToName(i)
		f.SetColWidth(sheet, col, col, 18)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
