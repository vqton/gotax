package einvoice

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"gotax/internal/domain"
)

type BIDVParser struct{}

func NewBIDVParser() *BIDVParser {
	return &BIDVParser{}
}

func (p *BIDVParser) GetBankCode() string {
	return "BIDV"
}

func (p *BIDVParser) Validate(data []byte) error {
	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("invalid CSV format: %w", err)
	}
	if len(records) < 2 {
		return fmt.Errorf("CSV file must have header and at least one data row")
	}
	header := records[0]
	if len(header) < 6 {
		return fmt.Errorf("BIDV CSV must have at least 6 columns")
	}
	return nil
}

func (p *BIDVParser) Parse(data []byte) ([]domain.BankTransaction, error) {
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
		if len(record) < 6 {
			continue
		}
		txn := domain.BankTransaction{
			TransactionDate: record[0],
			TransactionID:   record[1],
			Description:     record[2],
			Reference:       record[5],
		}
		debit, err := parseAmountBIDV(record[3])
		if err != nil {
			return nil, fmt.Errorf("invalid debit amount at row %d: %w", i+1, err)
		}
		txn.Debit = debit
		credit, err := parseAmountBIDV(record[4])
		if err != nil {
			return nil, fmt.Errorf("invalid credit amount at row %d: %w", i+1, err)
		}
		txn.Credit = credit
		transactions = append(transactions, txn)
	}
	return transactions, nil
}

func parseAmountBIDV(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, ".", "")
	return strconv.ParseFloat(s, 64)
}
