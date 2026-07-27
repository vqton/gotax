package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountValidate_Valid(t *testing.T) {
	a := &Account{Code: "1111", Name: "Cash", Type: AccountTypeAsset}
	assert.NoError(t, a.Validate())
	assert.Equal(t, AccountStatusActive, a.Status)
}

func TestAccountValidate_Invalid(t *testing.T) {
	tests := []struct {
		name string
		acct Account
		err  error
	}{
		{"empty code", Account{Name: "X", Type: AccountTypeAsset}, ErrAccountCodeRequired},
		{"short code", Account{Code: "11", Name: "X", Type: AccountTypeAsset}, ErrAccountCodeInvalid},
		{"non-digit code", Account{Code: "11AA", Name: "X", Type: AccountTypeAsset}, ErrAccountCodeInvalid},
		{"empty name", Account{Code: "1111", Type: AccountTypeAsset}, ErrAccountNameRequired},
		{"invalid type", Account{Code: "1111", Name: "X", Type: "BAD"}, ErrAccountInvalidType},
		{"invalid status", Account{Code: "1111", Name: "X", Type: AccountTypeAsset, Status: "BAD"}, ErrAccountStatusInvalid},
		{"frozen no reason", Account{Code: "1111", Name: "X", Type: AccountTypeAsset, Status: AccountStatusFrozen}, ErrFreezeReasonRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, tt.acct.Validate(), tt.err)
		})
	}
}

func TestAccountFreeze(t *testing.T) {
	a := &Account{Code: "1111", Name: "Cash", Type: AccountTypeAsset}
	require.NoError(t, a.Validate())

	assert.NoError(t, a.Freeze("fraud detected"))
	assert.Equal(t, AccountStatusFrozen, a.Status)
	assert.Equal(t, "fraud detected", a.FreezeReason)

	assert.ErrorIs(t, a.Freeze("again"), ErrAccountAlreadyFrozen)
	// already frozen, "" reason returns ErrAccountAlreadyFrozen (checked before reason)
	assert.ErrorIs(t, a.Freeze(""), ErrAccountAlreadyFrozen)
}

func TestAccountUnfreeze(t *testing.T) {
	a := &Account{Code: "1111", Name: "Cash", Type: AccountTypeAsset, Status: AccountStatusFrozen, FreezeReason: "audit"}
	assert.NoError(t, a.Unfreeze("audit complete"))
	assert.Equal(t, AccountStatusActive, a.Status)
	assert.Empty(t, a.FreezeReason)

	assert.ErrorIs(t, a.Unfreeze(""), ErrAccountNotFrozen)
}

func TestAccountCanPost(t *testing.T) {
	tests := []struct {
		name string
		acct Account
		err  error
	}{
		{"active", Account{Code: "1111", Name: "X", Type: AccountTypeAsset, Status: AccountStatusActive, IsActive: true}, nil},
		{"frozen", Account{Code: "1111", Name: "X", Type: AccountTypeAsset, Status: AccountStatusFrozen, IsActive: true}, ErrAccountFrozen},
		{"inactive", Account{Code: "1111", Name: "X", Type: AccountTypeAsset, Status: AccountStatusInactive, IsActive: false}, ErrAccountInactive},
		{"active but IsActive=false", Account{Code: "1111", Name: "X", Type: AccountTypeAsset, Status: AccountStatusActive, IsActive: false}, ErrAccountInactive},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err != nil {
				assert.ErrorIs(t, tt.acct.CanPost(), tt.err)
			} else {
				assert.NoError(t, tt.acct.CanPost())
			}
		})
	}
}

func TestJournalEntryTotalDebit(t *testing.T) {
	je := JournalEntry{Lines: []JournalLine{
		{DebitAmount: 100, CreditAmount: 0},
		{DebitAmount: 50, CreditAmount: 0},
	}}
	assert.Equal(t, 150.0, je.TotalDebit())
	assert.Equal(t, 0.0, je.TotalCredit())
	assert.True(t, je.HasDebit())
	assert.False(t, je.HasCredit())
}

func TestJournalEntryValidate_Valid(t *testing.T) {
	je := JournalEntry{
		EntryDate:   time.Now(),
		Description: "Test",
		Lines: []JournalLine{
			{LineNumber: 1, AccountCode: "1111", DebitAmount: 100, CreditAmount: 0},
			{LineNumber: 2, AccountCode: "5111", DebitAmount: 0, CreditAmount: 100},
		},
	}
	assert.NoError(t, je.Validate())
	assert.Equal(t, "VND", je.CurrencyCode)
}

func TestJournalEntryValidate_Invalid(t *testing.T) {
	tests := []struct {
		name string
		je   JournalEntry
		err  error
	}{
		{"no date", JournalEntry{Description: "X", Lines: []JournalLine{{}}}, ErrJournalEntryInvalidDate},
		{"no desc", JournalEntry{EntryDate: time.Now(), Lines: []JournalLine{{}}}, ErrJournalEntryNoDescription},
		{"no lines", JournalEntry{EntryDate: time.Now(), Description: "X"}, ErrJournalEntryNoLines},
		{"invalid amount", JournalEntry{EntryDate: time.Now(), Description: "X", Lines: []JournalLine{
			{DebitAmount: 0, CreditAmount: 0},
		}}, ErrJournalLineInvalidAmount},
		{"negative debit", JournalEntry{EntryDate: time.Now(), Description: "X", Lines: []JournalLine{
			{DebitAmount: -1, CreditAmount: 0},
		}}, ErrJournalLineInvalidAmount},
		{"unbalanced", JournalEntry{EntryDate: time.Now(), Description: "X", Lines: []JournalLine{
			{DebitAmount: 100, CreditAmount: 0},
		}}, ErrJournalEntryUnbalanced},
		{"foreign no rate", JournalEntry{EntryDate: time.Now(), Description: "X", CurrencyCode: "USD", Lines: []JournalLine{
			{DebitAmount: 100, CreditAmount: 0},
			{DebitAmount: 0, CreditAmount: 100},
		}}, ErrInvalidExchangeRate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, tt.je.Validate(), tt.err)
		})
	}
}

func TestPeriodValidate(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		p    Period
		err  error
	}{
		{"valid", Period{Year: 2025, Month: 6, StartDate: now, EndDate: now.Add(24 * time.Hour), Status: PeriodOpen}, nil},
		{"year too low", Period{Year: 1999, Month: 6, StartDate: now, EndDate: now}, ErrPeriodYearOutOfRange},
		{"year too high", Period{Year: 2101, Month: 6, StartDate: now, EndDate: now}, ErrPeriodYearOutOfRange},
		{"bad month", Period{Year: 2025, Month: 13, StartDate: now, EndDate: now}, ErrPeriodMonthInvalid},
		{"no dates", Period{Year: 2025, Month: 6}, ErrPeriodDateRequired},
		{"end before start", Period{Year: 2025, Month: 6, StartDate: now, EndDate: now.Add(-24 * time.Hour)}, ErrPeriodEndBeforeStart},
		{"bad status", Period{Year: 2025, Month: 6, StartDate: now, EndDate: now, Status: "NOPE"}, ErrPeriodStatusInvalid},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err != nil {
				assert.ErrorIs(t, tt.p.Validate(), tt.err)
			} else {
				assert.NoError(t, tt.p.Validate())
			}
		})
	}
}

func TestAccountBalanceCalculate(t *testing.T) {
	b := AccountBalance{
		OpenBalanceDebit: 1000, OpenBalanceCredit: 200,
		PeriodDebit: 500, PeriodCredit: 300,
	}
	b.Calculate()
	assert.Equal(t, 1500.0, b.TotalDebit)
	assert.Equal(t, 500.0, b.TotalCredit)
	assert.Equal(t, 1000.0, b.ClosingBalance)
}

func TestExchangeRateValidate(t *testing.T) {
	assert.ErrorIs(t, (&ExchangeRate{CurrencyCode: "VN", AverageRate: 1}).Validate(), ErrInvalidCurrencyCode)
	assert.ErrorIs(t, (&ExchangeRate{CurrencyCode: "USD", AverageRate: 0}).Validate(), ErrInvalidExchangeRate)
	assert.NoError(t, (&ExchangeRate{CurrencyCode: "USD", AverageRate: 23000}).Validate())
}

func TestUserValidate(t *testing.T) {
	assert.ErrorIs(t, (&User{}).Validate(), ErrUsernameRequired)
	assert.ErrorIs(t, (&User{Username: "test", Role: "bad"}).Validate(), ErrInvalidUserRole)
	assert.NoError(t, (&User{Username: "admin", Role: UserRoleAdmin}).Validate())
}

func TestUserIsLocked(t *testing.T) {
	u := &User{}
	assert.False(t, u.IsLocked())

	future := time.Now().Add(1 * time.Hour)
	u.LockedUntil = &future
	assert.True(t, u.IsLocked())

	past := time.Now().Add(-1 * time.Hour)
	u.LockedUntil = &past
	assert.False(t, u.IsLocked())
}

func TestApprovalStatusIsTerminal(t *testing.T) {
	assert.True(t, ApprovalApproved.IsTerminal())
	assert.True(t, ApprovalRejected.IsTerminal())
	assert.True(t, ApprovalCancelled.IsTerminal())
	assert.True(t, ApprovalExpired.IsTerminal())
	assert.False(t, ApprovalPending.IsTerminal())
}

func TestApprovalRequestValidate(t *testing.T) {
	assert.ErrorIs(t, (&ApprovalRequest{}).Validate(), ErrEntityTypeRequired)
	assert.NoError(t, (&ApprovalRequest{
		EntityType:  "account",
		EntityID:    "1111",
		RequestType: "change",
		Reason:      "restructure",
		RequestedBy: "user1",
	}).Validate())
}

func TestAccountMappingValidate(t *testing.T) {
	tests := []struct {
		name string
		m    AccountMapping
		err  error
	}{
		{"no source", AccountMapping{TargetRegime: "IFRS", OldCode: "111", NewCode: "1111", MappingType: "DIRECT"}, ErrSourceRegimeRequired},
		{"no old", AccountMapping{SourceRegime: "VAS", TargetRegime: "IFRS", NewCode: "1111", MappingType: "DIRECT"}, ErrOldCodeRequired},
		{"bad type", AccountMapping{SourceRegime: "VAS", TargetRegime: "IFRS", OldCode: "111", NewCode: "1111", MappingType: "BAD"}, ErrMappingTypeInvalid},
		{"split no ratio", AccountMapping{SourceRegime: "VAS", TargetRegime: "IFRS", OldCode: "111", NewCode: "1111", MappingType: "SPLIT"}, ErrSplitRatioRequired},
		{"valid", AccountMapping{SourceRegime: "VAS", TargetRegime: "IFRS", OldCode: "111", NewCode: "1111", MappingType: "DIRECT"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err != nil {
				assert.ErrorIs(t, tt.m.Validate(), tt.err)
			} else {
				assert.NoError(t, tt.m.Validate())
			}
		})
	}
}

func TestAccountAnalysisValidate(t *testing.T) {
	assert.ErrorIs(t, (&AccountAnalysis{}).Validate(), ErrAccountCodeRequired)
	assert.NoError(t, (&AccountAnalysis{AccountCode: "1111"}).Validate())
}

func TestIFRSMappingValidate(t *testing.T) {
	assert.ErrorIs(t, (&IFRSMapping{}).Validate(), ErrVASCodeRequired)
	assert.ErrorIs(t, (&IFRSMapping{VASCode: "1111"}).Validate(), ErrIFRSCodeRequired)
	assert.NoError(t, (&IFRSMapping{VASCode: "1111", IFRSCode: "IFRS-1"}).Validate())
}

func TestCompanyValidate(t *testing.T) {
	tests := []struct {
		name string
		c    Company
		err  error
	}{
		{"no tax code", Company{}, ErrCompanyCodeRequired},
		{"bad tax code len", Company{TaxCode: "12345", LegalNameVN: "X", LegalForm: LegalFormJSC, AccountingRegime: RegimeTT99}, ErrInvalidTaxCodeFormat},
		{"no name", Company{TaxCode: "1234567890"}, ErrCompanyNameRequired},
		{"bad legal form", Company{TaxCode: "1234567890", LegalNameVN: "X", LegalForm: "BAD", AccountingRegime: RegimeTT99}, ErrCompanyInvalidLegalForm},
		{"bad regime", Company{TaxCode: "1234567890", LegalNameVN: "X", LegalForm: LegalFormJSC, AccountingRegime: "BAD"}, ErrCompanyInvalidRegime},
		{"valid", Company{TaxCode: "1234567890", LegalNameVN: "Test Corp", LegalForm: LegalFormJSC, AccountingRegime: RegimeTT99}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err != nil {
				assert.ErrorIs(t, tt.c.Validate(), tt.err)
			} else {
				assert.NoError(t, tt.c.Validate())
			}
		})
	}
	// defaults
	c := Company{TaxCode: "1234567890", LegalNameVN: "Test Corp", LegalForm: LegalFormJSC, AccountingRegime: RegimeTT99}
	require.NoError(t, c.Validate())
	assert.Equal(t, CompanyStatusActive, c.Status)
	assert.Equal(t, "VND", c.DefaultCurrency)
	assert.Equal(t, 1, c.FiscalYearStartMonth)
}

func TestValidatePassword(t *testing.T) {
	assert.ErrorIs(t, ValidatePassword("short"), ErrPasswordTooShort)
	assert.ErrorIs(t, ValidatePassword("ALLUPPERCASE12!"), ErrPasswordNoLower)
	assert.ErrorIs(t, ValidatePassword("alllowercase12!"), ErrPasswordNoUpper)
	assert.ErrorIs(t, ValidatePassword("NoDigitHere!!"), ErrPasswordNoDigit)
	assert.ErrorIs(t, ValidatePassword("NoSpecialChar1"), ErrPasswordNoSpecial)
	assert.NoError(t, ValidatePassword("ValidPass123!"))
}

func TestIsPasswordInHistory(t *testing.T) {
	history := []string{"hash1", "hash2", "hash3"}
	assert.True(t, IsPasswordInHistory(history, "hash2"))
	assert.False(t, IsPasswordInHistory(history, "hash4"))
	assert.False(t, IsPasswordInHistory(nil, "hash1"))
}

func TestNormalBalance(t *testing.T) {
	assert.Equal(t, NormalBalanceDebit, AccountTypeAsset.NormalBalance())
	assert.Equal(t, NormalBalanceDebit, AccountTypeExpense.NormalBalance())
	assert.Equal(t, NormalBalanceCredit, AccountTypeLiability.NormalBalance())
	assert.Equal(t, NormalBalanceCredit, AccountTypeEquity.NormalBalance())
	assert.Equal(t, NormalBalanceCredit, AccountTypeRevenue.NormalBalance())
	assert.Equal(t, NormalBalanceDebit, AccountType("UNKNOWN").NormalBalance())
}

func TestAccountTypeConstants(t *testing.T) {
	assert.Equal(t, AccountType("ASSET"), AccountTypeAsset)
	assert.Equal(t, AccountType("LIABILITY"), AccountTypeLiability)
	assert.Equal(t, AccountType("EQUITY"), AccountTypeEquity)
	assert.Equal(t, AccountType("REVENUE"), AccountTypeRevenue)
	assert.Equal(t, AccountType("EXPENSE"), AccountTypeExpense)
}
