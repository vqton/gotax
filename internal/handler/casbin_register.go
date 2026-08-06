package handler

import "github.com/gin-gonic/gin"

// RegisterRoutesWithCompanyOpt is the opt-in Casbin variant.
// Call as: handler.RegisterRoutesWithCompanyOpt(r, h, ch, ..., authMW, adminMW, casbinMW)
// When casbinMW is omitted it falls back to the base RegisterRoutesWithCompany.
func RegisterRoutesWithCompanyOpt(r *gin.Engine, h *Handler, ch *CompanyHandler, th *TaxHandler, cashH *CashHandler, bankH *BankHandler, purchaseH *PurchaseHandler, saleH *SaleHandler, whH *WarehouseHandler, faH *FAHandler, pwH *PayrollHandler, recH *RecurringHandler, budH *BudgetHandler, ccdcH *CCDCHandler, authMW gin.HandlerFunc, adminMW gin.HandlerFunc, extra ...gin.HandlerFunc) {
	RegisterRoutesWithCompany(r, h, ch, th, cashH, bankH, purchaseH, saleH, whH, faH, pwH, recH, budH, ccdcH, authMW, adminMW)
	for _, mw := range extra {
		r.Use(mw)
	}
}
