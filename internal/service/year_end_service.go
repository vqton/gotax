package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gotax/internal/domain"
)

const retainedEarningsAccount = "421" // Lợi nhuận chưa phân phối (Retained Earnings)

type YearEndCloseResult struct {
	ClosedRevenueCount int                     `json:"closed_revenue_count"`
	ClosedExpenseCount int                     `json:"closed_expense_count"`
	TotalRevenue       float64                 `json:"total_revenue"`
	TotalExpense       float64                 `json:"total_expense"`
	NetIncome          float64                 `json:"net_income"`
	CarryForwardCount  int                     `json:"carry_forward_count"`
	MappingApplied     int                     `json:"mapping_applied"`
	ClosingEntryID     string                  `json:"closing_entry_id"`
	CarryForwardLog    *domain.CarryForwardLog `json:"carry_forward_log"`
}

// YearEndClose performs a complete year-end close:
// 1. Closes Revenue/Expense accounts to Retained Earnings (421)
// 2. Carries forward Balance Sheet accounts to new period
// 3. Applies TT200→TT99 account mapping if configured
func (s *service) YearEndClose(ctx context.Context, companyID, fromPeriodID, toPeriodID, fromYear, toYear, userID string) (*YearEndCloseResult, error) {
	result := &YearEndCloseResult{}

	// Step 1: Get all account balances for the closing period
	balances, err := s.ob.List(ctx, domain.OBListFilter{
		CompanyID: companyID,
		PeriodID:  fromPeriodID,
		Status:    domain.OBStatusApproved,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list balances: %w", err)
	}
	if len(balances) == 0 {
		return nil, domain.ErrOpeningBalanceNotFound
	}

	// Step 2: Separate Balance Sheet vs Income Statement accounts
	var bsBalances []domain.OpeningBalance
	var revenueTotal, expenseTotal float64
	var closingLines []domain.JournalLine

	for _, b := range balances {
		acct, err := s.accounts.GetByCode(ctx, b.AccountCode)
		if err != nil {
			continue // skip unknown accounts
		}

		switch acct.Type {
		case domain.AccountTypeRevenue:
			// Close revenue: Debit Revenue for (Credit - Debit) if normal credit balance
			// Or Credit Revenue for (Debit - Credit) if unusual debit balance
			result.ClosedRevenueCount++
			if b.CreditAmount > b.DebitAmount {
				// Normal: revenue has credit balance, close with debit
				revenueTotal += b.CreditAmount - b.DebitAmount
				closingLines = append(closingLines, domain.JournalLine{
					AccountCode: b.AccountCode,
					DebitAmount: b.CreditAmount - b.DebitAmount,
					Description: fmt.Sprintf("Đóng doanh thu tài khoản %s", b.AccountCode),
				})
			} else if b.DebitAmount > b.CreditAmount {
				// Unusual: revenue has debit balance, close with credit
				revenueTotal += b.DebitAmount - b.CreditAmount
				closingLines = append(closingLines, domain.JournalLine{
					AccountCode:  b.AccountCode,
					CreditAmount: b.DebitAmount - b.CreditAmount,
					Description:  fmt.Sprintf("Đóng doanh thu tài khoản %s", b.AccountCode),
				})
			}

		case domain.AccountTypeExpense:
			// Close expense: Credit Expense for (Debit - Credit) if normal debit balance
			// Or Debit Expense for (Credit - Debit) if unusual credit balance
			result.ClosedExpenseCount++
			if b.DebitAmount > b.CreditAmount {
				// Normal: expense has debit balance, close with credit
				expenseTotal += b.DebitAmount - b.CreditAmount
				closingLines = append(closingLines, domain.JournalLine{
					AccountCode: b.AccountCode,
					CreditAmount: b.DebitAmount - b.CreditAmount,
					Description: fmt.Sprintf("Đóng chi phí tài khoản %s", b.AccountCode),
				})
			} else if b.CreditAmount > b.DebitAmount {
				// Unusual: expense has credit balance, close with debit
				expenseTotal += b.CreditAmount - b.DebitAmount
				closingLines = append(closingLines, domain.JournalLine{
					AccountCode: b.AccountCode,
					DebitAmount: b.CreditAmount - b.DebitAmount,
					Description: fmt.Sprintf("Đóng chi phí tài khoản %s", b.AccountCode),
				})
			}

		default:
			// Balance Sheet accounts — carry forward
			bsBalances = append(bsBalances, b)
		}
	}

	result.TotalRevenue = revenueTotal
	result.TotalExpense = expenseTotal
	result.NetIncome = revenueTotal - expenseTotal

	// Step 3: Create closing journal entry (Revenue → Retained Earnings → Expense)
	if len(closingLines) > 0 {
		netIncome := revenueTotal - expenseTotal
		if netIncome > 0 {
			// Profit: Credit Retained Earnings
			closingLines = append(closingLines, domain.JournalLine{
				AccountCode:  retainedEarningsAccount,
				CreditAmount: netIncome,
				Description:  "Lợi nhuận năm",
			})
		} else if netIncome < 0 {
			// Loss: Debit Retained Earnings
			closingLines = append(closingLines, domain.JournalLine{
				AccountCode: retainedEarningsAccount,
				DebitAmount:  -netIncome,
				Description:  "Lỗ lũy kế năm",
			})
		}

		closingEntry := &domain.JournalEntry{
			CompanyID:   companyID,
			EntryDate:   time.Now(),
			Status:      domain.JournalEntryPosted,
			VoucherType: domain.VoucherTypeClosing,
			Description: fmt.Sprintf("Kết chuyển năm %s → %s", fromYear, toYear),
			Lines:       closingLines,
		}
		if err := s.journals.Create(ctx, closingEntry); err != nil {
			return nil, fmt.Errorf("failed to create closing entry: %w", err)
		}
		result.ClosingEntryID = closingEntry.ID
	}

	// Step 4: Apply TT200→TT99 mapping if configured
	mappings, _ := s.ob.ListCircular99Mappings(ctx)
	mappingMap := make(map[string]string)
	for _, m := range mappings {
		mappingMap[m.OldAccountCode] = m.NewAccountCode
	}

	// Step 5: Carry forward Balance Sheet accounts with mapping
	var carryBalances []domain.OpeningBalance
	for _, b := range bsBalances {
		accountCode := b.AccountCode
		if newCode, ok := mappingMap[accountCode]; ok {
			accountCode = newCode
			result.MappingApplied++
		}

		carryBalances = append(carryBalances, domain.OpeningBalance{
			CompanyID:      companyID,
			PeriodID:       toPeriodID,
			FiscalYearID:   toYear,
			AccountCode:    accountCode,
			CurrencyCode:   b.CurrencyCode,
			OriginalAmount: b.OriginalAmount,
			DebitAmount:    b.DebitAmount,
			CreditAmount:   b.CreditAmount,
			ExchangeRate:   b.ExchangeRate,
			Status:         domain.OBStatusApproved,
			SourceType:     "YEAR_END_CLOSE",
			Reason:         fmt.Sprintf("Year-end close from %s", fromPeriodID),
			CreatedBy:      userID,
		})
	}

	if err := s.ob.BulkCreate(ctx, carryBalances); err != nil {
		return nil, fmt.Errorf("failed to carry forward balances: %w", err)
	}
	result.CarryForwardCount = len(carryBalances)

	// Step 6: Log carry-forward
	fyFrom, err := strconv.Atoi(fromYear)
	if err != nil {
		return nil, fmt.Errorf("invalid from_year %q: %w", fromYear, err)
	}
	fyTo, err := strconv.Atoi(toYear)
	if err != nil {
		return nil, fmt.Errorf("invalid to_year %q: %w", toYear, err)
	}
	var totalDebit, totalCredit float64
	for _, b := range carryBalances {
		totalDebit += b.DebitAmount
		totalCredit += b.CreditAmount
	}

	log := &domain.CarryForwardLog{
		CompanyID:      companyID,
		FromPeriodID:   fromPeriodID,
		ToPeriodID:     toPeriodID,
		FromFiscalYear: fyFrom,
		ToFiscalYear:   fyTo,
		AccountCount:   len(carryBalances),
		TotalDebit:     totalDebit,
		TotalCredit:    totalCredit,
		Status:         "COMPLETED",
		ExecutedBy:     userID,
	}
	if err := s.ob.CreateCarryForwardLog(ctx, log); err != nil {
		return nil, fmt.Errorf("failed to create carry-forward log: %w", err)
	}
	result.CarryForwardLog = log

	return result, nil
}

// ExportYearEndBalances exports all account balances at year-end for reporting.
func (s *service) ExportYearEndBalances(ctx context.Context, companyID, periodID string) ([]domain.OpeningBalance, error) {
	return s.ob.List(ctx, domain.OBListFilter{
		CompanyID: companyID,
		PeriodID:  periodID,
		Status:    domain.OBStatusApproved,
	})
}
