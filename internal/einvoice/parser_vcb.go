package einvoice

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gotax/internal/domain"
)

type VCBParser struct{}

func NewVCBParser() *VCBParser {
	return &VCBParser{}
}

func (p *VCBParser) GetBankCode() string {
	return "VCB"
}

func (p *VCBParser) Validate(data []byte) error {
	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("invalid CSV format: %w", err)
	}
	if len(records) < 2 {
		return fmt.Errorf("CSV file must have header and at least one data row")
	}
	header := records[0]
	if len(header) < 5 {
		return fmt.Errorf("VCB CSV must have at least 5 columns")
	}
	return nil
}

func (p *VCBParser) Parse(data []byte) ([]domain.BankTransaction, error) {
	if err := p.Validate(data); err != nil {
		return nil, err
	}
	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}
	var transactions []domain.BankTransaction
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) < 5 {
			continue
		}
		txn := domain.BankTransaction{
			TransactionDate: record[0],
			TransactionID:   record[1],
			Description:     record[2],
			Reference:       record[4],
		}
		debit, err := parseAmount(record[3])
		if err != nil {
			return nil, fmt.Errorf("invalid debit amount at row %d: %w", i+1, err)
		}
		txn.Debit = debit
		credit, err := parseAmount(record[4])
		if err != nil {
			return nil, fmt.Errorf("invalid credit amount at row %d: %w", i+1, err)
		}
		txn.Credit = credit
		transactions = append(transactions, txn)
	}
	return transactions, nil
}

func parseAmount(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, ".", "")
	return strconv.ParseFloat(s, 64)
}

func parseVCBDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	formats := []string{
		"02/01/2006",
		"2006-01-02",
		"01/02/2006",
	}
	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date format: %s", s)
}
