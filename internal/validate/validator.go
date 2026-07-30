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
