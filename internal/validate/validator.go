package validate

import (
	"gotax/internal/domain"

	"github.com/go-playground/validator/v10"
)

var v *validator.Validate

func init() {
	v = validator.New()
	v.RegisterValidation("fastatus", faStatusValidator)
	v.RegisterValidation("damethod", depMethodValidator)
	v.RegisterValidation("fasource", faSourceValidator)
	v.RegisterValidation("fatrtype", faTrxTypeValidator)
	v.RegisterValidation("disposaltype", disposalTypeValidator)
	v.RegisterValidation("suppstatus", supplierStatusValidator)
	v.RegisterValidation("postatus", poStatusValidator)
	v.RegisterValidation("grnstatus", grnStatusValidator)
	v.RegisterValidation("invstatus", invStatusValidator)
	v.RegisterValidation("vattype", vatTypeValidator)
	v.RegisterValidation("apttype", apTypeValidator)
	v.RegisterValidation("costtype", costTypeValidator)
	v.RegisterValidation("allocmethod", allocMethodValidator)
}

func Validator() *validator.Validate {
	return v
}

func faStatusValidator(fl validator.FieldLevel) bool {
	s := domain.FixedAssetStatus(fl.Field().String())
	switch s {
	case domain.FADraft, domain.FACancelled, domain.FAActive,
		domain.FADepreciating, domain.FASuspended,
		domain.FAFullyDepr, domain.FADisposed, domain.FASold:
		return true
	}
	return false
}

func depMethodValidator(fl validator.FieldLevel) bool {
	m := domain.DepreciationMethod(fl.Field().String())
	switch m {
	case domain.DepStraightLine, domain.DepDecliningBalance, domain.DepProductionBased:
		return true
	}
	return false
}

func faSourceValidator(fl validator.FieldLevel) bool {
	s := domain.FASource(fl.Field().String())
	switch s {
	case domain.FASourcePurchase, domain.FASourceConstruction,
		domain.FASourceLease, domain.FASourceDonation,
		domain.FASourceCapitalContribution, domain.FASourceExchange:
		return true
	}
	return false
}

func faTrxTypeValidator(fl validator.FieldLevel) bool {
	t := domain.FATransactionType(fl.Field().String())
	switch t {
	case domain.FATrxAcquisition, domain.FATrxDepreciation,
		domain.FATrxAdjustment, domain.FATrxTransfer,
		domain.FATrxDisposal, domain.FATrxSale,
		domain.FATrxRevaluation, domain.FATrxImpairment,
		domain.FATrxCIPTransfer, domain.FATrxSuspend, domain.FATrxResume:
		return true
	}
	return false
}

func disposalTypeValidator(fl validator.FieldLevel) bool {
	t := domain.DisposalType(fl.Field().String())
	switch t {
	case domain.DisposalSale, domain.DisposalLiquidation,
		domain.DisposalDonation, domain.DisposalReturn:
		return true
	}
	return false
}

// ─── Purchase validators ──────────────────────────────────────────────

func supplierStatusValidator(fl validator.FieldLevel) bool {
	switch domain.SupplierStatus(fl.Field().String()) {
	case domain.SupplierActive, domain.SupplierSuspended, domain.SupplierBlacklisted:
		return true
	}
	return false
}

func poStatusValidator(fl validator.FieldLevel) bool {
	switch domain.POStatus(fl.Field().String()) {
	case domain.POStatusDraft, domain.POStatusApproved, domain.POStatusSent,
		domain.POStatusPartial, domain.POStatusReceived,
		domain.POStatusCancelled, domain.POStatusClosed:
		return true
	}
	return false
}

func grnStatusValidator(fl validator.FieldLevel) bool {
	switch domain.GRNStatus(fl.Field().String()) {
	case domain.GRNDraft, domain.GRNPosted, domain.GRNCancelled:
		return true
	}
	return false
}

func invStatusValidator(fl validator.FieldLevel) bool {
	switch domain.InvoiceStatus(fl.Field().String()) {
	case domain.InvoiceDraft, domain.InvoiceVerified, domain.InvoicePosted,
		domain.InvoicePaid, domain.InvoiceCancelled:
		return true
	}
	return false
}

func vatTypeValidator(fl validator.FieldLevel) bool {
	switch domain.VATType(fl.Field().String()) {
	case domain.VAT0, domain.VAT5, domain.VAT8, domain.VAT10, domain.VATNonTaxable:
		return true
	}
	return false
}

func apTypeValidator(fl validator.FieldLevel) bool {
	switch domain.APTransactionType(fl.Field().String()) {
	case domain.APTransInvoice, domain.APTransCreditNote, domain.APTransPayment,
		domain.APTransPrepayment, domain.APTransOffset:
		return true
	}
	return false
}

func costTypeValidator(fl validator.FieldLevel) bool {
	switch domain.CostAllocationType(fl.Field().String()) {
	case domain.CostTransport, domain.CostInsurance, domain.CostCustoms, domain.CostInspection:
		return true
	}
	return false
}

func allocMethodValidator(fl validator.FieldLevel) bool {
	switch domain.CostAllocationMethod(fl.Field().String()) {
	case domain.CostAllocByQty, domain.CostAllocByValue, domain.CostAllocByWeight, domain.CostAllocByVolume:
		return true
	}
	return false
}
