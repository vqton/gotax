package validate

import (
	"errors"

	"gotax/internal/domain"

	"github.com/go-playground/validator/v10"
)

func FixedAsset(a *domain.FixedAsset) error {
	if a.DepreciationMethod == "" {
		a.DepreciationMethod = domain.DepStraightLine
	}
	if a.Status == "" {
		a.Status = domain.FADraft
	}
	if a.Source == "" {
		a.Source = domain.FASourcePurchase
	}
	if err := v.Struct(a); err != nil {
		return mapAssetError(err)
	}
	return nil
}

func FixedAssetCategory(c *domain.FixedAssetCategory) error {
	if c.DefaultDepreciationMethod == "" {
		c.DefaultDepreciationMethod = domain.DepStraightLine
	}
	if c.DefaultUsefulLifeMonths <= 0 {
		c.DefaultUsefulLifeMonths = 120
	}
	if err := v.Struct(c); err != nil {
		return mapCategoryError(err)
	}
	if c.ParentID != nil && *c.ParentID == c.ID {
		return domain.ErrFACategorySelfParent
	}
	return nil
}

func mapAssetError(err error) error {
	var verr validator.ValidationErrors
	if !errors.As(err, &verr) || len(verr) == 0 {
		return err
	}
	fe := verr[0]
	switch fe.Field() {
	case "Code":
		return domain.ErrFACodeRequired
	case "Name":
		return domain.ErrFANameRequired
	case "CategoryID":
		return domain.ErrFACategoryRequired
	case "OriginalCost":
		return domain.ErrFAOriginalCostInvalid
	case "UsefulLifeMonths":
		return domain.ErrFAUsefulLifeInvalid
	case "DepreciationMethod":
		return domain.ErrFADepreciationMethodInvalid
	case "Status":
		return domain.ErrFAStatusInvalid
	case "AcquisitionDate":
		return domain.ErrFAAcquisitionDateRequired
	case "DepartmentID":
		return domain.ErrFADepartmentRequired
	case "Source":
		return domain.ErrFASourceInvalid
	case "AssetAccountID":
		return domain.ErrFAAssetAccountRequired
	case "DepreciationAccountID":
		return domain.ErrFADepreciationAccountRequired
	case "ExpenseAccountID":
		return domain.ErrFAExpenseAccountRequired
	default:
		return err
	}
}

func mapCategoryError(err error) error {
	var verr validator.ValidationErrors
	if !errors.As(err, &verr) || len(verr) == 0 {
		return err
	}
	fe := verr[0]
	switch fe.Field() {
	case "Code":
		return domain.ErrFACategoryCodeRequired
	case "Name":
		return domain.ErrFACategoryNameRequired
	case "Level":
		return domain.ErrFACategoryLevelInvalid
	case "AssetAccountID":
		return domain.ErrFAAssetAccountRequired
	case "DepreciationAccountID":
		return domain.ErrFADepreciationAccountRequired
	case "ExpenseAccountID":
		return domain.ErrFAExpenseAccountRequired
	default:
		return err
	}
}
