package gl

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestAccount_Freeze(t *testing.T) {
	t.Run("freeze active account", func(t *testing.T) {
		acc := &Account{Code: "1111", Name: "Test", Type: AccountTypeAsset, Status: AccountStatusActive}
		err := acc.Freeze("test freeze reason")
		assert.NoError(t, err)
		assert.Equal(t, AccountStatusFrozen, acc.Status)
		assert.Equal(t, "test freeze reason", acc.FreezeReason)
	})

	t.Run("freeze already frozen account", func(t *testing.T) {
		acc := &Account{Code: "1111", Name: "Test", Type: AccountTypeAsset, Status: AccountStatusFrozen}
		err := acc.Freeze("reason")
		assert.ErrorIs(t, err, ErrAccountAlreadyFrozen)
	})

	t.Run("freeze without reason", func(t *testing.T) {
		acc := &Account{Code: "1111", Name: "Test", Type: AccountTypeAsset, Status: AccountStatusActive}
		err := acc.Freeze("")
		assert.ErrorIs(t, err, ErrFreezeReasonRequired)
	})
}

func TestAccount_Unfreeze(t *testing.T) {
	t.Run("unfreeze frozen account", func(t *testing.T) {
		acc := &Account{Code: "1111", Name: "Test", Type: AccountTypeAsset, Status: AccountStatusFrozen, FreezeReason: "reason"}
		err := acc.Unfreeze("resolved")
		assert.NoError(t, err)
		assert.Equal(t, AccountStatusActive, acc.Status)
		assert.Empty(t, acc.FreezeReason)
	})

	t.Run("unfreeze active account", func(t *testing.T) {
		acc := &Account{Code: "1111", Name: "Test", Type: AccountTypeAsset, Status: AccountStatusActive}
		err := acc.Unfreeze("reason")
		assert.ErrorIs(t, err, ErrAccountNotFrozen)
	})
}

func TestAccount_CanPost(t *testing.T) {
	t.Run("active account can post", func(t *testing.T) {
		acc := &Account{Code: "1111", Name: "Test", Type: AccountTypeAsset, IsActive: true, Status: AccountStatusActive}
		assert.NoError(t, acc.CanPost())
	})

	t.Run("frozen account cannot post", func(t *testing.T) {
		acc := &Account{Code: "1111", Name: "Test", Type: AccountTypeAsset, IsActive: true, Status: AccountStatusFrozen}
		assert.ErrorIs(t, acc.CanPost(), ErrAccountFrozen)
	})

	t.Run("inactive account cannot post", func(t *testing.T) {
		acc := &Account{Code: "1111", Name: "Test", Type: AccountTypeAsset, IsActive: false, Status: AccountStatusInactive}
		assert.ErrorIs(t, acc.CanPost(), ErrAccountInactive)
	})
}

func TestApprovalRequest_Validate(t *testing.T) {
	t.Run("valid approval request", func(t *testing.T) {
		req := &ApprovalRequest{
			EntityType:  "ACCOUNT",
			EntityID:    "1111",
			RequestType: "CREATE",
			Reason:      "Need new account",
			RequestedBy: "user-1",
		}
		assert.NoError(t, req.Validate())
		assert.Equal(t, ApprovalPending, req.Status)
	})

	t.Run("missing entity type", func(t *testing.T) {
		req := &ApprovalRequest{EntityID: "1111", RequestType: "CREATE", Reason: "test", RequestedBy: "user-1"}
		assert.Error(t, req.Validate())
	})

	t.Run("missing reason", func(t *testing.T) {
		req := &ApprovalRequest{EntityType: "ACCOUNT", EntityID: "1111", RequestType: "CREATE", Reason: "", RequestedBy: "user-1"}
		assert.ErrorIs(t, req.Validate(), ErrApprovalReasonRequired)
	})
}

func TestApprovalStatus_IsTerminal(t *testing.T) {
	assert.True(t, ApprovalStatus("APPROVED").IsTerminal())
	assert.True(t, ApprovalStatus("REJECTED").IsTerminal())
	assert.True(t, ApprovalStatus("CANCELLED").IsTerminal())
	assert.True(t, ApprovalStatus("EXPIRED").IsTerminal())
	assert.False(t, ApprovalStatus("PENDING").IsTerminal())
}

func TestAccountMapping_Validate(t *testing.T) {
	t.Run("valid direct mapping", func(t *testing.T) {
		m := &AccountMapping{
			SourceRegime: "TT200",
			TargetRegime: "TT99",
			OldCode:      "1111",
			NewCode:      "1111",
			MappingType:  "DIRECT",
		}
		assert.NoError(t, m.Validate())
	})

	t.Run("invalid mapping type", func(t *testing.T) {
		m := &AccountMapping{
			SourceRegime: "TT200",
			TargetRegime: "TT99",
			OldCode:      "611",
			NewCode:      "632",
			MappingType:  "INVALID",
		}
		assert.Error(t, m.Validate())
	})

	t.Run("split requires ratio", func(t *testing.T) {
		m := &AccountMapping{
			SourceRegime: "TT200",
			TargetRegime: "TT99",
			OldCode:      "1562",
			NewCode:      "156",
			MappingType:  "SPLIT",
			SplitRatio:   0,
		}
		assert.Error(t, m.Validate())
	})

	t.Run("split with valid ratio", func(t *testing.T) {
		m := &AccountMapping{
			SourceRegime: "TT200",
			TargetRegime: "TT99",
			OldCode:      "1562",
			NewCode:      "156",
			MappingType:  "SPLIT",
			SplitRatio:   1.0,
		}
		assert.NoError(t, m.Validate())
	})
}

func TestAccountAnalysis_Validate(t *testing.T) {
	t.Run("valid analysis", func(t *testing.T) {
		a := &AccountAnalysis{AccountCode: "1111", CostCenterID: "CC-001"}
		assert.NoError(t, a.Validate())
	})

	t.Run("missing account code", func(t *testing.T) {
		a := &AccountAnalysis{AccountCode: ""}
		assert.Error(t, a.Validate())
	})
}

func TestIFRSMapping_Validate(t *testing.T) {
	t.Run("valid IFRS mapping", func(t *testing.T) {
		m := &IFRSMapping{VASCode: "1111", IFRSCode: "IFRS-1100"}
		assert.NoError(t, m.Validate())
	})

	t.Run("missing VAS code", func(t *testing.T) {
		m := &IFRSMapping{VASCode: "", IFRSCode: "IFRS-1100"}
		assert.Error(t, m.Validate())
	})

	t.Run("missing IFRS code", func(t *testing.T) {
		m := &IFRSMapping{VASCode: "1111", IFRSCode: ""}
		assert.Error(t, m.Validate())
	})
}

func TestAccount_ValidateWithStatus(t *testing.T) {
	t.Run("no status defaults to active", func(t *testing.T) {
		acc := &Account{Code: "1111", Name: "Test", Type: AccountTypeAsset}
		assert.NoError(t, acc.Validate())
		assert.Equal(t, AccountStatusActive, acc.Status)
	})

	t.Run("frozen requires reason", func(t *testing.T) {
		acc := &Account{Code: "1111", Name: "Test", Type: AccountTypeAsset, Status: AccountStatusFrozen}
		assert.ErrorIs(t, acc.Validate(), ErrFreezeReasonRequired)
	})

	t.Run("frozen with reason passes", func(t *testing.T) {
		acc := &Account{Code: "1111", Name: "Test", Type: AccountTypeAsset, Status: AccountStatusFrozen, FreezeReason: "reason"}
		assert.NoError(t, acc.Validate())
	})

	t.Run("invalid status", func(t *testing.T) {
		acc := &Account{Code: "1111", Name: "Test", Type: AccountTypeAsset, Status: "UNKNOWN"}
		assert.ErrorIs(t, acc.Validate(), ErrAccountStatusInvalid)
	})
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