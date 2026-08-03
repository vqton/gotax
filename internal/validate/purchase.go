package validate

import (
	"errors"

	"gotax/internal/domain"

	"github.com/go-playground/validator/v10"
)

func Supplier(s *domain.Supplier) error {
	if s.Status == "" {
		s.Status = domain.SupplierActive
	}
	if err := v.Struct(s); err != nil {
		return mapSupplierError(err)
	}
	return nil
}

func PurchaseOrder(po *domain.PurchaseOrder) error {
	if po.Status == "" {
		po.Status = domain.POStatusDraft
	}
	if err := v.Struct(po); err != nil {
		return mapPOError(err)
	}
	return nil
}

func GRN(g *domain.GRN) error {
	if g.Status == "" {
		g.Status = domain.GRNDraft
	}
	if err := v.Struct(g); err != nil {
		return mapGRNError(err)
	}
	return nil
}

func SupplierInvoice(inv *domain.SupplierInvoice) error {
	if inv.Status == "" {
		inv.Status = domain.InvoiceDraft
	}
	if err := v.Struct(inv); err != nil {
		return mapInvoiceError(err)
	}
	return nil
}

func APTransaction(t *domain.APTransaction) error {
	if t.Currency == "" {
		t.Currency = "VND"
	}
	if err := v.Struct(t); err != nil {
		return mapAPTransError(err)
	}
	return nil
}

func CostAllocation(c *domain.CostAllocation) error {
	return mapCostAllocError(v.Struct(c))
}

func Requisition(r *domain.PurchaseRequisition) error {
	if r.Status == "" {
		r.Status = domain.ReqDraft
	}
	if err := v.Struct(r); err != nil {
		return mapRequisitionError(err)
	}
	return nil
}

func firstField(err error) (string, bool) {
	var verr validator.ValidationErrors
	if !errors.As(err, &verr) || len(verr) == 0 {
		return "", false
	}
	return verr[0].Field(), true
}

func mapSupplierError(err error) error {
	f, ok := firstField(err)
	if !ok {
		return err
	}
	switch f {
	case "Code":
		return domain.ErrSupplierCodeRequired
	case "Name":
		return domain.ErrSupplierNameRequired
	case "TaxCode":
		return domain.ErrSupplierTaxCodeRequired
	case "Status":
		return domain.ErrSupplierStatusInvalid
	}
	return err
}

func mapPOError(err error) error {
	f, ok := firstField(err)
	if !ok {
		return err
	}
	switch f {
	case "PONumber":
		return domain.ErrPONumberRequired
	case "SupplierID":
		return domain.ErrPOSupplierRequired
	case "OrderDate":
		return domain.ErrPODateRequired
	case "Status":
		return domain.ErrPOStatusInvalid
	case "Lines":
		return domain.ErrPOLinesRequired
	case "Lines[0].ItemName":
		return domain.ErrPOItemNameRequired
	case "Lines[0].Unit":
		return domain.ErrPOItemUnitRequired
	case "Lines[0].Quantity":
		return domain.ErrPOItemQuantityRequired
	case "Lines[0].UnitPrice":
		return domain.ErrPOItemPriceRequired
	case "Lines[0].AccountID":
		return domain.ErrPOItemAccountRequired
	case "Lines[0].VATAccountID":
		return domain.ErrPOItemVATAccountRequired
	}
	return err
}

func mapGRNError(err error) error {
	f, ok := firstField(err)
	if !ok {
		return err
	}
	switch f {
	case "GRNNumber":
		return domain.ErrGRNNumberRequired
	case "POID":
		return domain.ErrGRNPORequired
	case "ReceiptDate":
		return domain.ErrGRNDateRequired
	case "Status":
		return domain.ErrGRNStatusInvalid
	case "Lines":
		return domain.ErrGRNLinesRequired
	case "Lines[0].ItemName":
		return domain.ErrGRNItemNameRequired
	case "Lines[0].POLineID":
		return domain.ErrGRNItemPOLineRequired
	case "Lines[0].QuantityReceived":
		return domain.ErrGRNItemQtyRequired
	}
	return err
}

func mapInvoiceError(err error) error {
	f, ok := firstField(err)
	if !ok {
		return err
	}
	switch f {
	case "InvoiceNumber":
		return domain.ErrInvoiceNumberRequired
	case "SupplierID":
		return domain.ErrInvoiceSupplierRequired
	case "SupplierName":
		return domain.ErrInvoiceSupplierNameRequired
	case "SupplierTaxCode":
		return domain.ErrInvoiceSupplierTaxCodeRequired
	case "InvoiceDate":
		return domain.ErrInvoiceDateRequired
	case "Status":
		return domain.ErrInvoiceStatusInvalidPurchase
	case "Lines":
		return domain.ErrInvoiceLinesRequired
	case "Lines[0].ItemName":
		return domain.ErrInvoiceItemNameRequired
	case "Lines[0].Quantity":
		return domain.ErrInvoiceItemQtyRequired
	case "Lines[0].UnitPrice":
		return domain.ErrInvoiceItemPriceRequired
	case "Lines[0].AccountID":
		return domain.ErrInvoiceItemAccountRequired
	case "Lines[0].VATAccountID":
		return domain.ErrInvoiceItemVATAccountRequired
	}
	return err
}

func mapAPTransError(err error) error {
	f, ok := firstField(err)
	if !ok {
		return err
	}
	switch f {
	case "SupplierID":
		return domain.ErrAPTransSupplierRequired
	case "TransactionDate":
		return domain.ErrAPTransDateRequired
	case "Amount":
		return domain.ErrAPTransAmountRequired
	case "TransactionType":
		return domain.ErrAPTransTypeInvalid
	}
	return err
}

func mapCostAllocError(err error) error {
	f, ok := firstField(err)
	if !ok {
		return err
	}
	switch f {
	case "InvoiceID":
		return domain.ErrCostAllocInvoiceRequired
	case "CostAmount":
		return domain.ErrCostAllocAmountRequired
	case "CostType":
		return domain.ErrCostAllocTypeInvalid
	case "AllocationMethod":
		return domain.ErrCostAllocMethodInvalid
	}
	return err
}

func mapRequisitionError(err error) error {
	f, ok := firstField(err)
	if !ok {
		return err
	}
	switch f {
	case "RequisitionNumber":
		return domain.ErrRequisitionNumberRequired
	case "RequesterID":
		return domain.ErrRequisitionRequesterRequired
	case "Status":
		return domain.ErrRequisitionStatusInvalid
	case "Lines":
		return domain.ErrRequisitionLinesRequired
	case "Lines[0].ItemName":
		return domain.ErrRequisitionItemNameRequired
	case "Lines[0].Quantity":
		return domain.ErrRequisitionItemQtyRequired
	case "Lines[0].EstimatedPrice":
		return domain.ErrRequisitionItemPriceRequired
	case "Lines[0].AccountID":
		return domain.ErrRequisitionItemAccountRequired
	}
	return err
}

func FXRevaluation(r *domain.FXRevaluation) error {
	if r.Status == "" {
		r.Status = domain.FXRevalDraft
	}
	if err := v.Struct(r); err != nil {
		return mapFXRevaluationError(err)
	}
	return nil
}

func mapFXRevaluationError(err error) error {
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return err
	}
	for _, fe := range verrs {
		switch fe.Field() {
		case "RevaluationDate":
			return domain.ErrFXRevaluationDateRequired
		case "Status":
			return domain.ErrFXRevaluationStatusInvalid
		case "Lines":
			return domain.ErrFXRevaluationLinesRequired
		}
	}
	return err
}
