package service

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"

	"gotax/internal/domain"
)

type OBImportRow struct {
	RowNumber  int
	AccountCode string
	CurrencyCode string
	OriginalAmount float64
	DebitAmount  float64
	CreditAmount float64
	ExchangeRate float64
	SourceType   string
	Reason       string
}

type OBImportError struct {
	Row   int    `json:"row"`
	Error string `json:"error"`
}

type OBImportResult struct {
	Total      int             `json:"total"`
	Success    int             `json:"success"`
	Errors     []OBImportError `json:"errors,omitempty"`
	BalanceIDs []string        `json:"balance_ids,omitempty"`
}

func (s *service) ImportOpeningBalances(ctx context.Context, data []byte, companyID, periodID, createdBy string) (*OBImportResult, error) {
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("open excel: %w", err)
	}
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return nil, fmt.Errorf("read sheet: %w", err)
	}
	if len(rows) < 2 {
		return &OBImportResult{Total: 0, Success: 0}, nil
	}

	header := rows[0]
	colMap := buildColumnMap(header)

	result := &OBImportResult{Total: len(rows) - 1}
	var balances []domain.OpeningBalance

	for i := 1; i < len(rows); i++ {
		ob, errs := parseOBRow(rows[i], colMap, i+1, companyID, periodID, createdBy)
		if len(errs) > 0 {
			result.Errors = append(result.Errors, errs...)
			continue
		}
		balances = append(balances, *ob)
	}

	if len(balances) == 0 {
		return result, nil
	}

	if err := s.ob.BulkCreate(ctx, balances); err != nil {
		return nil, fmt.Errorf("bulk create: %w", err)
	}
	result.Success = len(balances)
	for i := range balances {
		result.BalanceIDs = append(result.BalanceIDs, balances[i].ID)
	}
	return result, nil
}

func buildColumnMap(header []string) map[string]int {
	m := make(map[string]int)
	for i, col := range header {
		key := strings.ToLower(strings.TrimSpace(col))
		switch key {
		case "account_code", "account code", "accountcode", "mã tk", "ma tk":
			m["account_code"] = i
		case "currency_code", "currency code", "currencycode", "tiền tệ", "tien te", "loại tiền", "loai tien":
			m["currency_code"] = i
		case "original_amount", "original amount", "originalamount", "số dư đầu kỳ", "so du dau ky":
			m["original_amount"] = i
		case "debit_amount", "debit amount", "debitamount", "số phát sinh nợ", "so phat sinh no", "nợ", "no":
			m["debit_amount"] = i
		case "credit_amount", "credit amount", "creditamount", "số phát sinh có", "so phat sinh co", "có", "co":
			m["credit_amount"] = i
		case "exchange_rate", "exchange rate", "exchangerate", "tỷ giá", "ty gia":
			m["exchange_rate"] = i
		case "source_type", "source type", "sourcetype", "nguồn", "nguon":
			m["source_type"] = i
		case "reason", "lý do", "ly do", "ghi chú", "ghi chu":
			m["reason"] = i
		}
	}
	return m
}

func parseOBRow(row []string, colMap map[string]int, rowNum int, companyID, periodID, createdBy string) (*domain.OpeningBalance, []OBImportError) {
	var errs []OBImportError
	addErr := func(msg string) {
		errs = append(errs, OBImportError{Row: rowNum, Error: msg})
	}

	col := func(key string) string {
		idx, ok := colMap[key]
		if !ok || idx >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[idx])
	}

	ob := &domain.OpeningBalance{
		CompanyID:  companyID,
		PeriodID:   periodID,
		SourceType: "IMPORT",
		CreatedBy:  createdBy,
	}

	ob.AccountCode = col("account_code")
	if ob.AccountCode == "" {
		addErr("account_code is required")
	}

	ob.CurrencyCode = col("currency_code")
	if ob.CurrencyCode == "" {
		ob.CurrencyCode = "VND"
	}

	if v := col("original_amount"); v != "" {
		ob.OriginalAmount, _ = strconv.ParseFloat(v, 64)
	}
	if v := col("debit_amount"); v != "" {
		ob.DebitAmount, _ = strconv.ParseFloat(v, 64)
	}
	if v := col("credit_amount"); v != "" {
		ob.CreditAmount, _ = strconv.ParseFloat(v, 64)
	}
	if v := col("exchange_rate"); v != "" {
		ob.ExchangeRate, _ = strconv.ParseFloat(v, 64)
	}
	if v := col("source_type"); v != "" {
		ob.SourceType = v
	}
	if v := col("reason"); v != "" {
		ob.Reason = v
	}

	if ob.DebitAmount == 0 && ob.CreditAmount == 0 {
		addErr("debit_amount or credit_amount required")
	}

	if len(errs) > 0 {
		return nil, errs
	}
	if err := ob.Validate(); err != nil {
		addErr(err.Error())
	}
	if len(errs) > 0 {
		return nil, errs
	}
	return ob, nil
}
