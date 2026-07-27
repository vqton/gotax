package gl

import (
	"errors"
	"testing"
	"time"
)

func TestAccount_Validate(t *testing.T) {
	tests := []struct {
		name    string
		account Account
		wantErr error
	}{
		{
			name: "valid asset account",
			account: Account{
				Code:       "1111",
				Name:       "Tiền mặt VND",
				Type:       AccountTypeAsset,
				IsActive:   true,
				ParentCode: "111",
				DetailBy:   DetailByNone,
			},
			wantErr: nil,
		},
		{
			name: "invalid - empty code",
			account: Account{
				Code:     "",
				Name:     "Test",
				Type:     AccountTypeAsset,
				IsActive: true,
			},
			wantErr: ErrAccountCodeRequired,
		},
		{
			name: "invalid - empty name",
			account: Account{
				Code:     "1111",
				Name:     "",
				Type:     AccountTypeAsset,
				IsActive: true,
			},
			wantErr: ErrAccountNameRequired,
		},
		{
			name: "invalid - unknown type",
			account: Account{
				Code:     "1111",
				Name:     "Test",
				Type:     "UNKNOWN",
				IsActive: true,
			},
			wantErr: ErrAccountInvalidType,
		},
		{
			name: "invalid - code too short (less than 3 chars)",
			account: Account{
				Code:     "11",
				Name:     "Test",
				Type:     AccountTypeAsset,
				IsActive: true,
			},
			wantErr: ErrAccountCodeInvalid,
		},
		{
			name: "invalid - code with non-digit chars",
			account: Account{
				Code:     "11A1",
				Name:     "Test",
				Type:     AccountTypeAsset,
				IsActive: true,
			},
			wantErr: ErrAccountCodeInvalid,
		},
		{
			name: "valid - detail by object",
			account: Account{
				Code:       "1311",
				Name:       "Phải thu KH",
				Type:       AccountTypeAsset,
				IsActive:   true,
				DetailBy:   DetailByObject,
				IsForeign:  true,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.account.Validate()
			if tt.wantErr == nil && err != nil {
				t.Errorf("expected nil, got %v", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestJournalEntry_Validate(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		je      JournalEntry
		wantErr error
	}{
		{
			name: "valid simple entry",
			je: JournalEntry{
				EntryDate:   now,
				Description: "Mua hàng trả tiền mặt",
				Lines: []JournalLine{
					{AccountCode: "1561", DebitAmount: 10000000, CreditAmount: 0, Description: "Mua hàng"},
					{AccountCode: "1111", DebitAmount: 0, CreditAmount: 10000000, Description: "Trả tiền"},
				},
			},
			wantErr: nil,
		},
		{
			name: "invalid - no lines",
			je: JournalEntry{
				EntryDate:   now,
				Description: "Test",
				Lines:       []JournalLine{},
			},
			wantErr: ErrJournalEntryNoLines,
		},
		{
			name: "invalid - unbalanced debit credit",
			je: JournalEntry{
				EntryDate:   now,
				Description: "Test",
				Lines: []JournalLine{
					{AccountCode: "1561", DebitAmount: 10000000, CreditAmount: 0},
					{AccountCode: "1111", DebitAmount: 0, CreditAmount: 5000000},
				},
			},
			wantErr: ErrJournalEntryUnbalanced,
		},
		{
			name: "invalid - zero date",
			je: JournalEntry{
				Description: "Test",
				Lines: []JournalLine{
					{AccountCode: "1561", DebitAmount: 10000000, CreditAmount: 0},
					{AccountCode: "1111", DebitAmount: 0, CreditAmount: 10000000},
				},
			},
			wantErr: ErrJournalEntryInvalidDate,
		},
		{
			name: "invalid - empty description",
			je: JournalEntry{
				EntryDate:   now,
				Description: "",
				Lines: []JournalLine{
					{AccountCode: "1561", DebitAmount: 10000000, CreditAmount: 0},
					{AccountCode: "1111", DebitAmount: 0, CreditAmount: 10000000},
				},
			},
			wantErr: ErrJournalEntryNoDescription,
		},
		{
			name: "invalid - line with both zero amounts",
			je: JournalEntry{
				EntryDate:   now,
				Description: "Test",
				Lines: []JournalLine{
					{AccountCode: "1561", DebitAmount: 0, CreditAmount: 0},
				},
			},
			wantErr: ErrJournalLineInvalidAmount,
		},
		{
			name: "valid - multi-line entry",
			je: JournalEntry{
				EntryDate:   now,
				Description: "Phân bổ chi phí",
				Lines: []JournalLine{
					{AccountCode: "6421", DebitAmount: 5000000, CreditAmount: 0, Description: "CP nhân viên"},
					{AccountCode: "6422", DebitAmount: 3000000, CreditAmount: 0, Description: "CP vật tư"},
					{AccountCode: "1111", DebitAmount: 0, CreditAmount: 8000000, Description: "Tổng trả"},
				},
			},
			wantErr: nil,
		},
		{
			name: "invalid - line with negative amount",
			je: JournalEntry{
				EntryDate:   now,
				Description: "Test",
				Lines: []JournalLine{
					{AccountCode: "1561", DebitAmount: -1000, CreditAmount: 0},
					{AccountCode: "1111", DebitAmount: 0, CreditAmount: 1000},
				},
			},
			wantErr: ErrJournalLineInvalidAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.je.Validate()
			if tt.wantErr == nil && err != nil {
				t.Errorf("expected nil, got %v", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestAccountBalance_CalculateBalance(t *testing.T) {
	tests := []struct {
		name            string
		balance         AccountBalance
		expectedDebit   float64
		expectedCredit  float64
		expectedBalance float64
	}{
		{
			name: "asset account - debit normal",
			balance: AccountBalance{
				AccountType:     AccountTypeAsset,
				OpenBalanceDebit:  10000000,
				OpenBalanceCredit: 0,
				PeriodDebit:       5000000,
				PeriodCredit:      3000000,
			},
			expectedDebit:   15000000,
			expectedCredit:  3000000,
			expectedBalance: 12000000,
		},
		{
			name: "liability account - credit normal",
			balance: AccountBalance{
				AccountType:      AccountTypeLiability,
				OpenBalanceDebit:  0,
				OpenBalanceCredit: 50000000,
				PeriodDebit:       10000000,
				PeriodCredit:      20000000,
			},
			expectedDebit:   10000000,
			expectedCredit:  70000000,
			expectedBalance: -60000000,
		},
		{
			name: "expense account - debit normal, no open",
			balance: AccountBalance{
				AccountType:     AccountTypeExpense,
				PeriodDebit:     15000000,
				PeriodCredit:    2000000,
			},
			expectedDebit:   15000000,
			expectedCredit:  2000000,
			expectedBalance: 13000000,
		},
		{
			name: "revenue account - credit normal",
			balance: AccountBalance{
				AccountType:      AccountTypeRevenue,
				PeriodDebit:       5000000,
				PeriodCredit:     50000000,
			},
			expectedDebit:   5000000,
			expectedCredit:  50000000,
			expectedBalance: -45000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.balance.Calculate()

			if tt.balance.TotalDebit != tt.expectedDebit {
				t.Errorf("TotalDebit: got %.0f, want %.0f", tt.balance.TotalDebit, tt.expectedDebit)
			}
			if tt.balance.TotalCredit != tt.expectedCredit {
				t.Errorf("TotalCredit: got %.0f, want %.0f", tt.balance.TotalCredit, tt.expectedCredit)
			}
			if tt.balance.ClosingBalance != tt.expectedBalance {
				t.Errorf("ClosingBalance: got %.0f, want %.0f", tt.balance.ClosingBalance, tt.expectedBalance)
			}
		})
	}
}

func TestAccountType_NormalBalance(t *testing.T) {
	tests := []struct {
		accountType AccountType
		want        NormalBalance
	}{
		{AccountTypeAsset, NormalBalanceDebit},
		{AccountTypeExpense, NormalBalanceDebit},
		{AccountTypeLiability, NormalBalanceCredit},
		{AccountTypeEquity, NormalBalanceCredit},
		{AccountTypeRevenue, NormalBalanceCredit},
		{AccountType("UNKNOWN"), NormalBalanceDebit},
	}

	for _, tt := range tests {
		t.Run(string(tt.accountType), func(t *testing.T) {
			if got := tt.accountType.NormalBalance(); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJournalEntry_HasDebitCredit(t *testing.T) {
	je := JournalEntry{
		Lines: []JournalLine{
			{AccountCode: "1111", DebitAmount: 100000, CreditAmount: 0},
			{AccountCode: "5111", DebitAmount: 0, CreditAmount: 100000},
		},
	}

	if !je.HasDebit() {
		t.Error("expected HasDebit true")
	}
	if !je.HasCredit() {
		t.Error("expected HasCredit true")
	}
}

func TestJournalEntry_TotalDebitCredit(t *testing.T) {
	je := JournalEntry{
		Lines: []JournalLine{
			{AccountCode: "1111", DebitAmount: 50000, CreditAmount: 0},
			{AccountCode: "1121", DebitAmount: 30000, CreditAmount: 0},
			{AccountCode: "5111", DebitAmount: 0, CreditAmount: 80000},
		},
	}

	if je.TotalDebit() != 80000 {
		t.Errorf("TotalDebit: got %.0f, want 80000", je.TotalDebit())
	}
	if je.TotalCredit() != 80000 {
		t.Errorf("TotalCredit: got %.0f, want 80000", je.TotalCredit())
	}
}