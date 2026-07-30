package db

import (
	"context"
	"fmt"
	"log"

	"gotax/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormConfig struct {
	DSN           string
	MaxOpenConns  int
	MaxIdleConns  int
	ConnMaxLifetimeMS int
}

func DefaultGormConfig() GormConfig {
	dsn := "host=db port=5432 user=user password=pass dbname=gotax sslmode=disable"
	return GormConfig{
		DSN:              dsn,
		MaxOpenConns:     20,
		MaxIdleConns:     5,
		ConnMaxLifetimeMS: 1800000,
	}
}

func NewGorm(ctx context.Context, cfg GormConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  cfg.DSN,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm open: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	if cfg.ConnMaxLifetimeMS > 0 {
		sqlDB.SetConnMaxLifetime(0)
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping: %w", err)
	}
	return db, nil
}

func CloseGorm(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func RunGormMigrations(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&domain.AccountGORM{},
		&domain.JournalEntryGORM{},
		&domain.JournalLineGORM{},
		&domain.PeriodGORM{},
		&domain.UserGORM{},
		&domain.RefreshTokenGORM{},
		&domain.PasswordResetTokenGORM{},
		&domain.AuditEntryGORM{},
		&domain.ExchangeRateGORM{},
		&domain.ClosingTemplateGORM{},
		&domain.ApprovalRequestGORM{},
		&domain.AccountVersionGORM{},
		&domain.AccountMappingGORM{},
		&domain.AccountAnalysisGORM{},
		&domain.IFRSMappingGORM{},
		&domain.CompanyGORM{},
		&domain.CompanyBranchGORM{},
		&domain.FiscalYearGORM{},
		&domain.PeriodV2GORM{},
		&domain.DepartmentGORM{},
		&domain.EmployeeGORM{},
		&domain.CompanyBankAccountGORM{},
		&domain.EInvoicePatternGORM{},
		&domain.DigitalSignatureGORM{},
		&domain.IntegrationProfileGORM{},
		&domain.OpeningBalanceGORM{},
		&domain.OpeningBalanceDetailGORM{},
		&domain.CarryForwardLogGORM{},
		&domain.Circular99MappingGORM{},
		&domain.BalanceMigrationGORM{},
		&domain.BankStatementGORM{},
		&domain.BankStatementLineGORM{},
		&domain.BankReconciliationGORM{},
		&domain.BankReconciliationMatchGORM{},
		&domain.PaymentOrderGORM{},
		&domain.PaymentOrderBatchGORM{},
		&domain.LoanAgreementGORM{},
		&domain.LoanDisbursementGORM{},
		&domain.LoanRepaymentGORM{},
		&domain.TermDepositGORM{},
		&domain.CashReceiptGORM{},
		&domain.CashReceiptLineGORM{},
		&domain.CashPaymentGORM{},
		&domain.CashPaymentLineGORM{},
		&domain.CashTransferGORM{},
		&domain.PettyCashFundGORM{},
		&domain.CashInventoryGORM{},
		&domain.AdvanceRequestGORM{},
		&domain.AdvanceSettlementGORM{},
		&domain.SupplierGORM{},
		&domain.PurchaseOrderGORM{},
		&domain.POItemGORM{},
		&domain.GRNGORM{},
		&domain.GRNItemGORM{},
		&domain.SupplierInvoiceGORM{},
		&domain.SupplierInvoiceLineGORM{},
		&domain.APTransactionGORM{},
		&domain.CostAllocationGORM{},
		&domain.CustomerGORM{},
		&domain.SalesOrderGORM{},
		&domain.SOLineGORM{},
		&domain.DeliveryNoteGORM{},
		&domain.DNLineGORM{},
		&domain.CustomerInvoiceGORM{},
		&domain.InvLineGORM{},
		&domain.CustomerReceiptGORM{},
		&domain.RcpAllocationGORM{},
		&domain.CreditNoteGORM{},
		&domain.CNLineGORM{},
		&domain.ARTransactionGORM{},
		&domain.SalesQuotationGORM{},
		&domain.WarehouseGORM{},
		&domain.ItemCategoryGORM{},
		&domain.ItemGORM{},
		&domain.StockBalanceGORM{},
		&domain.InventoryTransactionGORM{},
		&domain.StockTransferGORM{},
		&domain.TransferItemGORM{},
		&domain.StockAdjustmentGORM{},
		&domain.AdjItemGORM{},
		&domain.StockTakeGORM{},
		&domain.TakeItemGORM{},
		&domain.InventoryValuationRunGORM{},
		&domain.TaxDeclarationGORM{},
		&domain.TaxDeclarationLineGORM{},
		&domain.TaxRateGORM{},
		&domain.TaxPaymentGORM{},
		&domain.EInvoiceGORM{},
		&domain.EInvoiceLineGORM{},
		&domain.TaxCalendarGORM{},
		&domain.TaxAlertGORM{},
		&domain.TaxAuditCaseGORM{},
	); err != nil {
		return fmt.Errorf("auto-migrate: %w", err)
	}
	log.Println("gorm: auto-migration complete")
	return nil
}
