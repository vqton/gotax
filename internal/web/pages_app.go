package web

import (
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
	"gotax/internal/service"
)

// Deps bundles the services pages render from.
type Deps struct {
	Svc      service.Service
	Company  service.CompanyService
	Sale     *service.SaleService
	Purchase *service.PurchaseService
}

// NewPages returns the registry of server-rendered pages → loaders.
func NewPages(d Deps) map[string]Page {
	return map[string]Page{
		"/app/dashboard.html":       {Title: "Tổng quan", NavPath: "/app/dashboard.html", Load: dashboardLoad(d)},
		"/app/users.html":           {Title: "Người dùng", NavPath: "/app/users.html", Load: usersLoad(d)},
		"/app/journal-entries.html": {Title: "Chứng từ kế toán", NavPath: "/app/journal-entries.html", Load: journalEntriesLoad(d)},
		"/app/coa.html":             {Title: "Hệ thống tài khoản", NavPath: "/app/coa.html", Load: coaLoad(d)},
		"/app/customers.html":       {Title: "Khách hàng", NavPath: "/app/customers.html", Load: customersLoad(d)},
		"/app/suppliers.html":       {Title: "Nhà cung cấp", NavPath: "/app/suppliers.html", Load: suppliersLoad(d)},
		"/app/company.html":         {Title: "Thông tin công ty", NavPath: "/app/company.html", Load: companyLoad(d)},
		"/app/exchange-rates.html":  {Title: "Tỷ giá hối đoái", NavPath: "/app/exchange-rates.html", Load: exchangeRatesLoad(d)},
		"/app/periods.html":         {Title: "Kỳ kế toán", NavPath: "/app/periods.html", Load: periodsLoad(d)},
		"/app/cash-receipts.html":   {Title: "Phiếu thu", NavPath: "/app/cash-receipts.html", Load: cashReceiptsLoad(d)},
		"/app/cash-payments.html":   {Title: "Phiếu chi", NavPath: "/app/cash-payments.html", Load: cashPaymentsLoad(d)},
		"/app/cash-transfers.html":  {Title: "Chuyển quỹ", NavPath: "/app/cash-transfers.html", Load: cashTransfersLoad(d)},
	}
}

// NewActions returns htmx mutation endpoints: path → action → handler.
func (s *Server) NewActions(d Deps) map[string]map[string]gin.HandlerFunc {
	return map[string]map[string]gin.HandlerFunc{
		"/app/users": {
			"create": s.usersCreate(d),
			"role":   s.usersRole(d),
		},
		"/app/journal-entries": {
			"filter":  s.journalFilter(d),
			"create":  s.journalCreate(d),
			"submit":  s.journalStatus(d, "submit"),
			"approve": s.journalStatus(d, "approve"),
			"post":    s.journalStatus(d, "post"),
			"cancel":  s.journalStatus(d, "cancel"),
		},
		"/app/coa": {
			"save":     s.coaSave(d),
			"delete":   s.coaDelete(d),
			"freeze":   s.coaFreeze(d),
			"unfreeze": s.coaUnfreeze(d),
		},
		"/app/customers": {
			"create": s.customersCreate(d),
			"delete": s.customersDelete(d),
		},
		"/app/suppliers": {
			"create": s.suppliersCreate(d),
			"delete": s.suppliersDelete(d),
		},
		"/app/company": {
			"save": s.companySave(d),
		},
		"/app/exchange-rates": {
			"create": s.exchangeRatesCreate(d),
		},
		"/app/periods": {
			"create": s.periodsCreate(d),
			"close":  s.periodsClose(d),
			"reopen": s.periodsReopen(d),
		},
		"/app/cash-receipts": {
			"create": s.cashReceiptsCreate(d),
			"submit": s.cashReceiptStatus(d, "submit"),
			"approve": s.cashReceiptStatus(d, "approve"),
			"post":   s.cashReceiptStatus(d, "post"),
		},
		"/app/cash-payments": {
			"create": s.cashPaymentsCreate(d),
			"submit": s.cashPaymentStatus(d, "submit"),
			"approve": s.cashPaymentStatus(d, "approve"),
			"post":   s.cashPaymentStatus(d, "post"),
		},
		"/app/cash-transfers": {
			"create": s.cashTransfersCreate(d),
		},
	}
}

// StatusStat is one row of the dashboard status-distribution card.
type StatusStat struct {
	Status     string
	Label      string
	Count      int
	Pct        int
	BadgeClass string
	FillClass  string
}

func dashboardLoad(d Deps) func(c *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
		ctx := c.Request.Context()
		accs, err := d.Svc.GetAllAccounts(ctx, true)
		if err != nil {
			return nil, err
		}
		entries, err := d.Svc.GetAllEntries(ctx)
		if err != nil {
			return nil, err
		}
		if entries == nil {
			entries = []domain.JournalEntry{}
		}
		total := len(entries)
		counts := map[domain.JournalEntryStatus]int{}
		for _, e := range entries {
			counts[e.Status]++
		}
		// Recent = last 10 POSTED entries, with running totals.
		var recent []domain.JournalEntry
		var recentDebit, recentCredit float64
		for i := len(entries) - 1; i >= 0 && len(recent) < 10; i-- {
			e := entries[i]
			if e.Status != domain.JournalEntryPosted {
				continue
			}
			recent = append(recent, e)
			recentDebit += e.TotalDebit()
			recentCredit += e.TotalCredit()
		}
		order := []StatusStat{
			{Status: string(domain.JournalEntryDraft), Label: "Nháp", BadgeClass: "badge-draft", FillClass: "s-fill-draft"},
			{Status: string(domain.JournalEntryReviewing), Label: "Chờ duyệt", BadgeClass: "badge-reviewing", FillClass: "s-fill-reviewing"},
			{Status: string(domain.JournalEntryApproved), Label: "Đã duyệt", BadgeClass: "badge-approved", FillClass: "s-fill-approved"},
			{Status: string(domain.JournalEntryPosted), Label: "Đã ghi sổ", BadgeClass: "badge-posted", FillClass: "s-fill-posted"},
			{Status: string(domain.JournalEntryCancelled), Label: "Hủy", BadgeClass: "badge-cancelled", FillClass: "s-fill-cancelled"},
		}
		stats := make([]StatusStat, 0, len(order))
		for _, s := range order {
			n := counts[domain.JournalEntryStatus(s.Status)]
			pct := 0
			if total > 0 {
				pct = int(math.Round(float64(n) * 100 / float64(total)))
			}
			s.Count, s.Pct = n, pct
			stats = append(stats, s)
		}
		return gin.H{
			"Accounts":          len(accs),
			"Entries":           counts[domain.JournalEntryPosted],
			"StatusStats":       stats,
			"Recent":            recent,
			"RecentTotalDebit":  recentDebit,
			"RecentTotalCredit": recentCredit,
		}, nil
	}
}

func usersLoad(d Deps) func(c *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
		users, err := d.Svc.ListUsers(c.Request.Context())
		if err != nil {
			return nil, err
		}
		return gin.H{"Users": users}, nil
	}
}

// pageCompanyID resolves the tenant for company-scoped pages. Passed via
// ?company_id= like the JSON API; falls back to CMP001 to match the legacy
// app.js default (single-company deployments).
func pageCompanyID(c *gin.Context) string {
	if id := c.Query("company_id"); id != "" {
		return id
	}
	return "CMP001"
}

func customersLoad(d Deps) func(c *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
		customers, err := d.Sale.ListCustomers(c.Request.Context(), pageCompanyID(c))
		if err != nil {
			return nil, err
		}
		if customers == nil {
			customers = []domain.Customer{}
		}
		return gin.H{"Customers": customers}, nil
	}
}

func (s *Server) customersCreate(d Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		cust := &domain.Customer{
			CompanyID: pageCompanyID(c),
			Code:      c.PostForm("code"),
			Name:      c.PostForm("name"),
			TaxCode:   c.PostForm("tax_code"),
			Address:   c.PostForm("address"),
			Phone:     c.PostForm("phone"),
			Email:     c.PostForm("email"),
			Currency:  "VND",
		}
		if err := d.Sale.CreateCustomer(c.Request.Context(), cust); err != nil {
			log.Printf("create customer: %v", err)
			Toast(c, "error", "Không tạo được khách hàng: "+err.Error())
			s.renderCustomersTable(c, d)
			return
		}
		Toast(c, "success", "Đã thêm khách hàng mới.")
		s.renderCustomersTable(c, d)
	}
}

func (s *Server) customersDelete(d Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.PostForm("id")
		if id == "" {
			Toast(c, "error", "Thiếu mã khách hàng.")
			s.renderCustomersTable(c, d)
			return
		}
		if err := d.Sale.DeleteCustomer(c.Request.Context(), id); err != nil {
			log.Printf("delete customer: %v", err)
			Toast(c, "error", "Không xóa được khách hàng: "+err.Error())
			s.renderCustomersTable(c, d)
			return
		}
		Toast(c, "success", "Đã xóa khách hàng.")
		s.renderCustomersTable(c, d)
	}
}

// ── Suppliers ────────────────────────────────────────────────────────────────

func suppliersLoad(d Deps) func(c *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
		suppliers, _, err := d.Purchase.ListSuppliers(c.Request.Context(), pageCompanyID(c), 0, 0)
		if err != nil {
			return nil, err
		}
		if suppliers == nil {
			suppliers = []domain.Supplier{}
		}
		return gin.H{"Suppliers": suppliers}, nil
	}
}

func (s *Server) suppliersCreate(d Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		sup := &domain.Supplier{
			CompanyID: pageCompanyID(c),
			Code:      c.PostForm("code"),
			Name:      c.PostForm("name"),
			TaxCode:   c.PostForm("tax_code"),
			Address:   c.PostForm("address"),
			Phone:     c.PostForm("phone"),
			Email:     c.PostForm("email"),
			Currency:  "VND",
		}
		// Service does not validate Supplier — check required fields here.
		if sup.Code == "" || sup.Name == "" || sup.TaxCode == "" {
			Toast(c, "error", "Vui lòng nhập mã, tên và mã số thuế nhà cung cấp.")
			s.renderSuppliersTable(c, d)
			return
		}
		if err := d.Purchase.CreateSupplier(c.Request.Context(), sup); err != nil {
			log.Printf("create supplier: %v", err)
			Toast(c, "error", "Không tạo được nhà cung cấp: "+err.Error())
			s.renderSuppliersTable(c, d)
			return
		}
		Toast(c, "success", "Đã thêm nhà cung cấp mới.")
		s.renderSuppliersTable(c, d)
	}
}

func (s *Server) suppliersDelete(d Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.PostForm("id")
		if id == "" {
			Toast(c, "error", "Thiếu mã nhà cung cấp.")
			s.renderSuppliersTable(c, d)
			return
		}
		if err := d.Purchase.DeleteSupplier(c.Request.Context(), id); err != nil {
			log.Printf("delete supplier: %v", err)
			Toast(c, "error", "Không xóa được nhà cung cấp: "+err.Error())
			s.renderSuppliersTable(c, d)
			return
		}
		Toast(c, "success", "Đã xóa nhà cung cấp.")
		s.renderSuppliersTable(c, d)
	}
}

func (s *Server) renderSuppliersTable(c *gin.Context, d Deps) {
	suppliers, _, err := d.Purchase.ListSuppliers(c.Request.Context(), pageCompanyID(c), 0, 0)
	if err != nil {
		c.String(500, "load suppliers failed")
		return
	}
	if suppliers == nil {
		suppliers = []domain.Supplier{}
	}
	s.RenderFragment(c, "suppliers", "suppliers-table", gin.H{"Suppliers": suppliers})
}

// ── Master data: company / exchange rates / periods ─────────────────────────

func companyLoad(d Deps) func(c *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
		company, err := d.Company.GetCompany(c.Request.Context(), pageCompanyID(c))
		if err != nil {
			// No company yet → render the empty form so it can be filled in.
			return gin.H{"Company": domain.Company{ID: pageCompanyID(c)}}, nil
		}
		return gin.H{"Company": *company}, nil
	}
}

func (s *Server) companySave(d Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := pageCompanyID(c)
		company, err := d.Company.GetCompany(c.Request.Context(), id)
		if err != nil {
			Toast(c, "error", "Không tìm thấy công ty: "+err.Error())
			s.renderCompanyForm(c, d)
			return
		}
		// Copy — the memory repo returns its live pointer; mutating it would
		// corrupt stored data even when validation rejects the save.
		upd := *company
		upd.LegalNameVN = strings.TrimSpace(c.PostForm("legal_name_vn"))
		upd.LegalNameEN = strings.TrimSpace(c.PostForm("legal_name_en"))
		upd.TaxCode = strings.TrimSpace(c.PostForm("tax_code"))
		upd.RegAddress = strings.TrimSpace(c.PostForm("reg_address"))
		upd.Phone = strings.TrimSpace(c.PostForm("phone"))
		upd.Email = strings.TrimSpace(c.PostForm("email"))
		upd.Website = strings.TrimSpace(c.PostForm("website"))
		upd.LegalRepName = strings.TrimSpace(c.PostForm("legal_rep_name"))
		upd.LegalRepTitle = strings.TrimSpace(c.PostForm("legal_rep_title"))
		if upd.LegalNameVN == "" || upd.TaxCode == "" {
			Toast(c, "error", "Tên công ty và mã số thuế là bắt buộc.")
			s.renderCompanyForm(c, d)
			return
		}
		if err := d.Company.UpdateCompany(c.Request.Context(), &upd); err != nil {
			log.Printf("company save: %v", err)
			Toast(c, "error", "Không lưu được thông tin công ty: "+err.Error())
			s.renderCompanyForm(c, d)
			return
		}
		Toast(c, "success", "Đã lưu thông tin công ty.")
		s.renderCompanyForm(c, d)
	}
}

func (s *Server) renderCompanyForm(c *gin.Context, d Deps) {
	company, err := d.Company.GetCompany(c.Request.Context(), pageCompanyID(c))
	if err != nil {
		company = &domain.Company{ID: pageCompanyID(c)}
	}
	s.RenderFragment(c, "company", "company-form", gin.H{"Company": *company})
}

func exchangeRatesLoad(d Deps) func(c *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
		rates, err := d.Svc.ListExchangeRates(c.Request.Context())
		if err != nil {
			return nil, err
		}
		if rates == nil {
			rates = []domain.ExchangeRate{}
		}
		return gin.H{"Rates": rates}, nil
	}
}

func (s *Server) exchangeRatesCreate(d Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		rateDate, err := time.Parse("2006-01-02", c.PostForm("rate_date"))
		if err != nil {
			Toast(c, "error", "Ngày tỷ giá không hợp lệ.")
			s.renderExchangeRatesTable(c, d)
			return
		}
		rate := &domain.ExchangeRate{
			CurrencyCode: strings.ToUpper(strings.TrimSpace(c.PostForm("currency_code"))),
			RateDate:     rateDate,
			Source:       strings.TrimSpace(c.PostForm("source")),
		}
		if rate.AverageRate, err = strconv.ParseFloat(c.PostForm("average_rate"), 64); err != nil {
			Toast(c, "error", "Tỷ giá trung bình không hợp lệ.")
			s.renderExchangeRatesTable(c, d)
			return
		}
		if v := c.PostForm("buy_rate"); v != "" {
			if rate.BuyRate, err = strconv.ParseFloat(v, 64); err != nil {
				Toast(c, "error", "Tỷ giá mua không hợp lệ.")
				s.renderExchangeRatesTable(c, d)
				return
			}
		}
		if v := c.PostForm("sell_rate"); v != "" {
			if rate.SellRate, err = strconv.ParseFloat(v, 64); err != nil {
				Toast(c, "error", "Tỷ giá bán không hợp lệ.")
				s.renderExchangeRatesTable(c, d)
				return
			}
		}
		// CreateExchangeRate runs rate.Validate() (currency 3 chars, avg > 0).
		if err := d.Svc.CreateExchangeRate(c.Request.Context(), rate); err != nil {
			log.Printf("create exchange rate: %v", err)
			Toast(c, "error", "Không tạo được tỷ giá: "+err.Error())
			s.renderExchangeRatesTable(c, d)
			return
		}
		Toast(c, "success", "Đã thêm tỷ giá "+rate.CurrencyCode+".")
		s.renderExchangeRatesTable(c, d)
	}
}

func (s *Server) renderExchangeRatesTable(c *gin.Context, d Deps) {
	rates, err := d.Svc.ListExchangeRates(c.Request.Context())
	if err != nil {
		c.String(500, "load exchange rates failed")
		return
	}
	if rates == nil {
		rates = []domain.ExchangeRate{}
	}
	s.RenderFragment(c, "exchange-rates", "exchange-rates-table", gin.H{"Rates": rates})
}

func periodsLoad(d Deps) func(c *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
		periods, err := d.Svc.GetAllPeriods(c.Request.Context())
		if err != nil {
			return nil, err
		}
		if periods == nil {
			periods = []domain.Period{}
		}
		sort.Slice(periods, func(i, j int) bool {
			if periods[i].Year != periods[j].Year {
				return periods[i].Year > periods[j].Year
			}
			return periods[i].Month > periods[j].Month
		})
		return gin.H{"Periods": periods}, nil
	}
}

func (s *Server) periodsCreate(d Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		year, err := strconv.Atoi(c.PostForm("year"))
		if err != nil {
			Toast(c, "error", "Năm không hợp lệ.")
			s.renderPeriodsTable(c, d)
			return
		}
		month, err := strconv.Atoi(c.PostForm("month"))
		if err != nil {
			Toast(c, "error", "Tháng không hợp lệ.")
			s.renderPeriodsTable(c, d)
			return
		}
		start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		period := &domain.Period{
			Year:      year,
			Month:     month,
			StartDate: start,
			EndDate:   start.AddDate(0, 1, -1),
			Status:    domain.PeriodOpen,
		}
		// CreatePeriod runs period.Validate() (year/month/date/status range).
		if err := d.Svc.CreatePeriod(c.Request.Context(), period); err != nil {
			log.Printf("create period: %v", err)
			Toast(c, "error", "Không tạo được kỳ kế toán: "+err.Error())
			s.renderPeriodsTable(c, d)
			return
		}
		Toast(c, "success", "Đã tạo kỳ kế toán "+strconv.Itoa(month)+"/"+strconv.Itoa(year)+".")
		s.renderPeriodsTable(c, d)
	}
}

// periodStatusAction dispatches close/reopen from hx-vals id.
func (s *Server) periodStatusAction(d Deps, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.PostForm("id")
		if id == "" {
			Toast(c, "error", "Thiếu mã kỳ kế toán.")
			s.renderPeriodsTable(c, d)
			return
		}
		var err error
		if action == "close" {
			err = d.Svc.ClosePeriod(c.Request.Context(), id)
		} else {
			err = d.Svc.ReopenPeriod(c.Request.Context(), id)
		}
		if err != nil {
			log.Printf("period %s %s: %v", action, id, err)
			Toast(c, "error", "Thao tác thất bại: "+err.Error())
		} else {
			Toast(c, "success", "Đã cập nhật kỳ kế toán.")
		}
		s.renderPeriodsTable(c, d)
	}
}

func (s *Server) periodsClose(d Deps) gin.HandlerFunc  { return s.periodStatusAction(d, "close") }
func (s *Server) periodsReopen(d Deps) gin.HandlerFunc { return s.periodStatusAction(d, "reopen") }

func (s *Server) renderPeriodsTable(c *gin.Context, d Deps) {
	periods, err := d.Svc.GetAllPeriods(c.Request.Context())
	if err != nil {
		c.String(500, "load periods failed")
		return
	}
	if periods == nil {
		periods = []domain.Period{}
	}
	sort.Slice(periods, func(i, j int) bool {
		if periods[i].Year != periods[j].Year {
			return periods[i].Year > periods[j].Year
		}
		return periods[i].Month > periods[j].Month
	})
	s.RenderFragment(c, "periods", "periods-table", gin.H{"Periods": periods})
}

// renderCustomersTable re-renders the table fragment after a mutation, keeping
// the modal + filters (which live outside the fragment) untouched.
func (s *Server) renderCustomersTable(c *gin.Context, d Deps) {
	customers, err := d.Sale.ListCustomers(c.Request.Context(), pageCompanyID(c))
	if err != nil {
		c.String(500, "load customers failed")
		return
	}
	if customers == nil {
		customers = []domain.Customer{}
	}
	s.RenderFragment(c, "customers", "customers-table", gin.H{"Customers": customers})
}

func (s *Server) usersCreate(d Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := &domain.User{
			Username: c.PostForm("username"),
			Email:    c.PostForm("email"),
			Role:     domain.UserRole(c.PostForm("role")),
			IsActive: true,
		}
		password := c.PostForm("password")
		if user.Username == "" || password == "" {
			Toast(c, "error", "Thiếu tên đăng nhập hoặc mật khẩu.")
			s.renderUsersTable(c, d)
			return
		}
		if err := d.Svc.CreateUser(c.Request.Context(), user, password); err != nil {
			log.Printf("create user: %v", err)
			Toast(c, "error", "Không tạo được người dùng: "+err.Error())
			s.renderUsersTable(c, d)
			return
		}
		Toast(c, "success", "Đã tạo người dùng thành công.")
		s.renderUsersTable(c, d)
	}
}

func (s *Server) usersRole(d Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.PostForm("id")
		role := domain.UserRole(c.PostForm("role"))
		if id == "" || role == "" {
			Toast(c, "error", "Thiếu thông tin.")
			s.renderUsersTable(c, d)
			return
		}
		if err := d.Svc.UpdateUserRole(c.Request.Context(), id, role); err != nil {
			Toast(c, "error", "Không cập nhật được vai trò: "+err.Error())
			s.renderUsersTable(c, d)
			return
		}
		Toast(c, "success", "Đã cập nhật vai trò.")
		s.renderUsersTable(c, d)
	}
}

func (s *Server) renderUsersTable(c *gin.Context, d Deps) {
	users, err := d.Svc.ListUsers(c.Request.Context())
	if err != nil {
		c.String(500, "load users failed")
		return
	}
	s.RenderFragment(c, "users", "users-table", gin.H{"Users": users})
}

// ── Journal entries ──────────────────────────────────────────────────────────

// journalFilterValues reads filter inputs from query (page load) or form
// (htmx filter POST). Returns status + date range strings.
func journalFilterValues(c *gin.Context) (status, from, to string) {
	status = c.Query("status")
	from = c.Query("from")
	to = c.Query("to")
	if status == "" && from == "" && to == "" {
		status = c.PostForm("status")
		from = c.PostForm("from")
		to = c.PostForm("to")
	}
	return
}

// journalData resolves the filtered entry list + filter echo for the table.
func journalData(c *gin.Context, d Deps) (gin.H, error) {
	status, from, to := journalFilterValues(c)
	ctx := c.Request.Context()
	var entries []domain.JournalEntry
	var err error
	switch {
	case status != "":
		entries, err = d.Svc.GetEntriesByStatus(ctx, domain.JournalEntryStatus(strings.ToUpper(status)))
	case from != "" && to != "":
		f, ferr := time.Parse("2006-01-02", from)
		t, terr := time.Parse("2006-01-02", to)
		if ferr != nil || terr != nil {
			entries, err = d.Svc.GetAllEntries(ctx)
		} else {
			entries, err = d.Svc.GetEntriesByDateRange(ctx, f, t.Add(24*time.Hour))
		}
	default:
		entries, err = d.Svc.GetAllEntries(ctx)
	}
	if err != nil {
		return nil, err
	}
	if entries == nil {
		entries = []domain.JournalEntry{}
	}
	return gin.H{
		"Entries":      entries,
		"FilterStatus": strings.ToLower(status),
		"FilterFrom":   from,
		"FilterTo":     to,
	}, nil
}

func journalEntriesLoad(d Deps) func(c *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
		return journalData(c, d)
	}
}

func (s *Server) journalFilter(d Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := journalData(c, d)
		if err != nil {
			log.Printf("journal filter: %v", err)
			Toast(c, "error", "Không tải được chứng từ: "+err.Error())
			data = gin.H{"Entries": []domain.JournalEntry{}}
		}
		s.RenderFragment(c, "journal-entries", "journal-table", data)
	}
}

func (s *Server) journalCreate(d Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		entryDate, err := time.Parse("2006-01-02", c.PostForm("entry_date"))
		if err != nil {
			Toast(c, "error", "Ngày chứng từ không hợp lệ.")
			s.renderJournalTable(c, d)
			return
		}
		voucherType := domain.VoucherType(c.PostForm("voucher_type"))
		if voucherType == "" {
			voucherType = domain.VoucherTypeOther
		}
		codes := c.PostFormArray("account_code")
		descs := c.PostFormArray("line_description")
		debits := c.PostFormArray("debit")
		credits := c.PostFormArray("credit")
		var lines []domain.JournalLine
		for i := range codes {
			if strings.TrimSpace(codes[i]) == "" {
				continue
			}
			line := domain.JournalLine{AccountCode: strings.TrimSpace(codes[i])}
			if i < len(descs) {
				line.Description = descs[i]
			}
			if i < len(debits) && strings.TrimSpace(debits[i]) != "" {
				v, perr := strconv.ParseFloat(debits[i], 64)
				if perr != nil {
					Toast(c, "error", "Số tiền Nợ dòng "+strconv.Itoa(i+1)+" không hợp lệ.")
					s.renderJournalTable(c, d)
					return
				}
				line.DebitAmount = v
			}
			if i < len(credits) && strings.TrimSpace(credits[i]) != "" {
				v, perr := strconv.ParseFloat(credits[i], 64)
				if perr != nil {
					Toast(c, "error", "Số tiền Có dòng "+strconv.Itoa(i+1)+" không hợp lệ.")
					s.renderJournalTable(c, d)
					return
				}
				line.CreditAmount = v
			}
			lines = append(lines, line)
		}
		entry := &domain.JournalEntry{
			EntryDate:    entryDate,
			VoucherType:  voucherType,
			Description:  c.PostForm("description"),
			CurrencyCode: "VND",
			Lines:        lines,
		}
		userID := c.GetString("user_id")
		if err := d.Svc.CreateEntry(c.Request.Context(), entry, userID); err != nil {
			log.Printf("journal create: %v", err)
			Toast(c, "error", "Không tạo được chứng từ: "+err.Error())
		} else {
			Toast(c, "success", "Đã tạo chứng từ "+entry.EntryNumber+".")
		}
		s.renderJournalTable(c, d)
	}
}

// journalStatus dispatches workflow transitions. id comes from hx-vals.
func (s *Server) journalStatus(d Deps, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.PostForm("id")
		userID := c.GetString("user_id")
		ctx := c.Request.Context()
		var err error
		switch action {
		case "submit":
			err = d.Svc.SubmitForReview(ctx, id, userID)
		case "approve":
			err = d.Svc.ApproveEntry(ctx, id, userID)
		case "post":
			err = d.Svc.PostEntry(ctx, id)
		case "cancel":
			err = d.Svc.CancelEntry(ctx, id)
		}
		if err != nil {
			log.Printf("journal %s %s: %v", action, id, err)
			Toast(c, "error", "Thao tác thất bại: "+err.Error())
		} else {
			Toast(c, "success", "Đã cập nhật trạng thái.")
		}
		s.renderJournalTable(c, d)
	}
}

// renderJournalTable re-renders the journal table fragment with the current
// filter state (post actions must keep the visible filter).
func (s *Server) renderJournalTable(c *gin.Context, d Deps) {
	data, err := journalData(c, d)
	if err != nil {
		log.Printf("render journal table: %v", err)
		c.String(500, "load journal entries failed")
		return
	}
	s.RenderFragment(c, "journal-entries", "journal-table", data)
}

// ── COA (Hệ thống tài khoản) ─────────────────────────────────────────────

// CoaLoai is one loại node in the left pane + loại header rows of the grid.
type CoaLoai struct {
	Code  string
	Name  string
	Count int
}

// CoaRow is one grid row (loại header or account), precomputed server-side so
// the template holds zero business logic.
type CoaRow struct {
	Code         string
	Name         string
	Name2        string
	ParentCode   string
	LoaiCode     string // top-level loại ancestor code (filter target)
	Type         string // ASSET/LIABILITY/... (form select value)
	TypeLabel    string // TS/NV/VH/DN/CP (badge text)
	TypeName     string // Tài sản/Nguồn vốn/... (tooltip)
	TypeClass    string // badge-type-ts/... (badge color)
	Normal       string // Nợ/Có
	NormalClass  string // debit/credit (legend dot color)
	DetailBy     string // OBJECT/PROJECT/... (form select value)
	DetailLabel  string // Đối tượng/Công trình/... (display)
	IsForeign    bool
	Status       string // ACTIVE/FROZEN/INACTIVE (filter value)
	StatusClass  string // badge-success/warning/default
	StatusLabel  string // Hoạt động/Ngừng TD/Không dùng
	FreezeReason string
	Depth        int
	HasChildren  bool
	IsLoai       bool
	LoaiName     string
	LoaiCount    int
	Search       string // lowercase code|name|name2 for client filter
}

// acctLess orders accounts by code length then code — loại headers (1-digit)
// first, then cấp 1 (3-digit), then deeper levels. Lexicographic compare
// alone would put "9" after "111".
func acctLess(a, b domain.Account) bool {
	if len(a.Code) != len(b.Code) {
		return len(a.Code) < len(b.Code)
	}
	return a.Code < b.Code
}

// coaTypeMeta maps AccountType → short label, full name, badge class.
func coaTypeMeta(t domain.AccountType) (label, name, cls string) {
	switch t {
	case domain.AccountTypeAsset:
		return "TS", "Tài sản", "badge-type-ts"
	case domain.AccountTypeLiability:
		return "NV", "Nguồn vốn", "badge-type-nv"
	case domain.AccountTypeEquity:
		return "VH", "Vốn chủ sở hữu", "badge-type-vh"
	case domain.AccountTypeRevenue:
		return "DN", "Doanh thu", "badge-type-dn"
	case domain.AccountTypeExpense:
		return "CP", "Chi phí", "badge-type-cp"
	}
	return "", "", "badge-default"
}

// coaStatusMeta maps account status → badge class + label (shared with the
// accountStatusBadge helper).
func coaStatusMeta(status string) (class, label string) {
	switch strings.ToUpper(status) {
	case "ACTIVE":
		return "badge-success", "Hoạt động"
	case "FROZEN":
		return "badge-warning", "Ngừng TD"
	case "INACTIVE":
		return "badge-default", "Không dùng"
	}
	return "badge-default", status
}

// coaDetailLabel maps DetailBy → Vietnamese display label.
func coaDetailLabel(d domain.DetailBy) string {
	switch d {
	case domain.DetailByObject:
		return "Đối tượng"
	case domain.DetailByProject:
		return "Công trình"
	case domain.DetailByContract:
		return "Hợp đồng"
	case domain.DetailByCostItem:
		return "Khoản mục"
	case domain.DetailByDepartment:
		return "Phòng ban"
	}
	return "—"
}

// coaData builds the full COA view: left-pane loại list, flat grid rows
// (loại headers + account rows in tree order), parent datalist options.
func coaData(c *gin.Context, d Deps) (gin.H, error) {
	accs, err := d.Svc.GetAllAccounts(c.Request.Context(), false)
	if err != nil {
		return nil, err
	}
	if accs == nil {
		accs = []domain.Account{}
	}
	children := map[string][]domain.Account{}
	for _, a := range accs {
		children[a.ParentCode] = append(children[a.ParentCode], a)
	}
	for p := range children {
		sort.Slice(children[p], func(i, j int) bool { return acctLess(children[p][i], children[p][j]) })
	}

	var subtreeCount func(code string) int
	subtreeCount = func(code string) int {
		n := 0
		for _, ch := range children[code] {
			n += 1 + subtreeCount(ch.Code)
		}
		return n
	}

	loais := []CoaLoai{}
	rows := []CoaRow{}
	parentOpts := []CoaRow{}

	var walk func(a domain.Account, loaiCode string, depth int)
	walk = func(a domain.Account, loaiCode string, depth int) {
		hasChildren := len(children[a.Code]) > 0
		if depth == 0 {
			loaiCode = a.Code
			loais = append(loais, CoaLoai{Code: a.Code, Name: a.Name, Count: subtreeCount(a.Code) + 1})
		} else {
			parentOpts = append(parentOpts, CoaRow{Code: a.Code, Name: a.Name})
		}
		tl, tn, tc := coaTypeMeta(a.Type)
		normal := "Có"
		normalClass := "credit"
		if a.Type.NormalBalance() == domain.NormalBalanceDebit {
			normal = "Nợ"
			normalClass = "debit"
		}
		sc, sl := coaStatusMeta(string(a.Status))
		rows = append(rows, CoaRow{
			Code: a.Code, Name: a.Name, Name2: a.Name2, ParentCode: a.ParentCode,
			LoaiCode:  loaiCode,
			Type:      string(a.Type),
			TypeLabel: tl, TypeName: tn, TypeClass: tc,
			Normal: normal, NormalClass: normalClass,
			DetailBy: string(a.DetailBy), DetailLabel: coaDetailLabel(a.DetailBy),
			IsForeign: a.IsForeign,
			Status:    string(a.Status), StatusClass: sc, StatusLabel: sl,
			FreezeReason: a.FreezeReason,
			Depth:        depth, HasChildren: hasChildren,
			IsLoai: depth == 0, LoaiName: a.Name, LoaiCount: subtreeCount(a.Code) + 1,
			Search: strings.ToLower(a.Code + "|" + a.Name + "|" + a.Name2),
		})
		for _, ch := range children[a.Code] {
			walk(ch, loaiCode, depth+1)
		}
	}
	loaiAccs := children[""]
	sort.Slice(loaiAccs, func(i, j int) bool { return acctLess(loaiAccs[i], loaiAccs[j]) })
	for _, l := range loaiAccs {
		walk(l, "", 0)
	}

	return gin.H{
		"Loais":         loais,
		"Rows":          rows,
		"ParentOptions": parentOpts,
		"Total":         len(rows),
	}, nil
}

func coaLoad(d Deps) func(c *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
		return coaData(c, d)
	}
}

// renderCoaGrid re-renders the grid fragment after a mutation, keeping the
// client-side filter state (search/status/loại live outside the fragment).
func (s *Server) renderCoaGrid(c *gin.Context, d Deps) {
	data, err := coaData(c, d)
	if err != nil {
		log.Printf("render coa grid: %v", err)
		c.String(500, "load accounts failed")
		return
	}
	s.RenderFragment(c, "coa", "coa-grid", data)
}

// coaSave creates or updates one account. mode=create|update from the modal.
// Code edits are not supported (API pins the code to the path param), so the
// modal locks the code field in update mode and submits it as-is.
func (s *Server) coaSave(d Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		acc := &domain.Account{
			Code:       strings.TrimSpace(c.PostForm("code")),
			Name:       strings.TrimSpace(c.PostForm("name")),
			Name2:      strings.TrimSpace(c.PostForm("name2")),
			Type:       domain.AccountType(c.PostForm("type")),
			ParentCode: strings.TrimSpace(c.PostForm("parent_code")),
			DetailBy:   domain.DetailBy(c.PostForm("detail_by")),
			IsForeign:  c.PostForm("is_foreign") == "on" || c.PostForm("is_foreign") == "true",
		}
		ctx := c.Request.Context()
		var err error
		if c.PostForm("mode") == "update" {
			err = d.Svc.UpdateAccount(ctx, acc)
		} else {
			err = d.Svc.CreateAccount(ctx, acc)
		}
		if err != nil {
			log.Printf("coa save %s: %v", acc.Code, err)
			Toast(c, "error", "Không lưu được tài khoản: "+err.Error())
		} else {
			Toast(c, "success", "Đã lưu tài khoản "+acc.Code+".")
		}
		s.renderCoaGrid(c, d)
	}
}

func (s *Server) coaDelete(d Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := strings.TrimSpace(c.PostForm("code"))
		if err := d.Svc.DeleteAccount(c.Request.Context(), code); err != nil {
			log.Printf("coa delete %s: %v", code, err)
			Toast(c, "error", "Không xóa được tài khoản: "+err.Error())
		} else {
			Toast(c, "success", "Đã xóa tài khoản "+code+".")
		}
		s.renderCoaGrid(c, d)
	}
}

// coaFreeze requires a reason (Ngừng theo dõi tài khoản).
func (s *Server) coaFreeze(d Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := strings.TrimSpace(c.PostForm("code"))
		reason := strings.TrimSpace(c.PostForm("reason"))
		if code == "" || reason == "" {
			Toast(c, "error", "Thiếu mã tài khoản hoặc lý do ngừng theo dõi.")
			s.renderCoaGrid(c, d)
			return
		}
		if err := d.Svc.FreezeAccount(c.Request.Context(), code, reason); err != nil {
			log.Printf("coa freeze %s: %v", code, err)
			Toast(c, "error", "Không ngừng theo dõi được: "+err.Error())
		} else {
			Toast(c, "success", "Đã ngừng theo dõi tài khoản "+code+".")
		}
		s.renderCoaGrid(c, d)
	}
}

func (s *Server) coaUnfreeze(d Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := strings.TrimSpace(c.PostForm("code"))
		if code == "" {
			Toast(c, "error", "Thiếu mã tài khoản.")
			s.renderCoaGrid(c, d)
			return
		}
		if err := d.Svc.UnfreezeAccount(c.Request.Context(), code, "Mở khóa tài khoản"); err != nil {
			log.Printf("coa unfreeze %s: %v", code, err)
			Toast(c, "error", "Không mở khóa được tài khoản: "+err.Error())
		} else {
			Toast(c, "success", "Đã mở khóa tài khoản "+code+".")
		}
		s.renderCoaGrid(c, d)
	}
}

// ─── Cash receipts ─────────────────────────────────────────────────

func cashReceiptsLoad(d Deps) func(c *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
		receipts, total, err := d.Svc.ListCashReceipts(c.Request.Context(), domain.CashReceiptFilter{CompanyID: pageCompanyID(c)})
		if err != nil {
			return nil, err
		}
		if receipts == nil {
			receipts = []domain.CashReceipt{}
		}
		return gin.H{"Receipts": receipts, "Total": total}, nil
	}
}

func (s *Server) cashReceiptsCreate(d Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		voucherDate := strings.TrimSpace(c.PostForm("voucher_date"))
		if _, err := time.Parse("2006-01-02", voucherDate); err != nil {
			Toast(c, "error", "Ngày lập phiếu không hợp lệ.")
			s.renderCashReceiptsTable(c, d)
			return
		}
		amount, err := strconv.ParseFloat(c.PostForm("amount"), 64)
		if err != nil || amount <= 0 {
			Toast(c, "error", "Số tiền không hợp lệ.")
			s.renderCashReceiptsTable(c, d)
			return
		}
		// Legacy page mirrors credit_account_id onto debit_account_id.
		credit := strings.TrimSpace(c.PostForm("credit_account_id"))
		debit := strings.TrimSpace(c.PostForm("debit_account_id"))
		if credit == "" {
			credit = debit
		}
		receiptType := domain.ReceiptType(c.PostForm("receipt_type"))
		if receiptType == "" {
			receiptType = domain.ReceiptCustomerPayment
		}
		r := &domain.CashReceipt{
			CompanyID:       pageCompanyID(c),
			VoucherDate:     voucherDate,
			CashAccountID:   strings.TrimSpace(c.PostForm("cash_account_id")),
			CounterpartName: strings.TrimSpace(c.PostForm("counterpart_name")),
			CounterpartType: domain.CounterpartCustomer,
			Currency:        "VND",
			ExchangeRate:    1,
			Amount:          amount,
			DebitAccountID:  debit,
			CreditAccountID: credit,
			Reason:          strings.TrimSpace(c.PostForm("reason")),
			ReceiptType:     receiptType,
		}
		// CreateCashReceipt runs r.Validate() + voucher numbering.
		if err := d.Svc.CreateCashReceipt(c.Request.Context(), r); err != nil {
			log.Printf("create cash receipt: %v", err)
			Toast(c, "error", "Không tạo được phiếu thu: "+err.Error())
			s.renderCashReceiptsTable(c, d)
			return
		}
		Toast(c, "success", "Đã tạo phiếu thu "+r.VoucherNo+".")
		s.renderCashReceiptsTable(c, d)
	}
}

func (s *Server) cashReceiptStatus(d Deps, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.PostForm("id")
		ctx := c.Request.Context()
		var err error
		switch action {
		case "submit":
			err = d.Svc.SubmitCashReceipt(ctx, id, c.GetString("user_id"))
		case "approve":
			err = d.Svc.ApproveCashReceipt(ctx, id, c.GetString("user_id"))
		case "post":
			err = d.Svc.PostCashReceipt(ctx, id, c.GetString("user_id"))
		}
		if err != nil {
			log.Printf("cash receipt %s %s: %v", action, id, err)
			Toast(c, "error", "Thao tác thất bại: "+err.Error())
		} else {
			Toast(c, "success", "Đã cập nhật trạng thái phiếu thu.")
		}
		s.renderCashReceiptsTable(c, d)
	}
}

func (s *Server) renderCashReceiptsTable(c *gin.Context, d Deps) {
	receipts, total, err := d.Svc.ListCashReceipts(c.Request.Context(), domain.CashReceiptFilter{CompanyID: pageCompanyID(c)})
	if err != nil {
		c.String(500, "load cash receipts failed")
		return
	}
	if receipts == nil {
		receipts = []domain.CashReceipt{}
	}
	s.RenderFragment(c, "cash-receipts", "cash-receipts-table", gin.H{"Receipts": receipts, "Total": total})
}

// ─── Cash Payments ─────────────────────────────────────────────────

func cashPaymentsLoad(d Deps) func(c *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
		payments, total, err := d.Svc.ListCashPayments(c.Request.Context(), domain.CashPaymentFilter{CompanyID: pageCompanyID(c)})
		if err != nil {
			return nil, err
		}
		if payments == nil {
			payments = []domain.CashPayment{}
		}
		return gin.H{"Payments": payments, "Total": total}, nil
	}
}

func (s *Server) cashPaymentsCreate(d Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		voucherDate := strings.TrimSpace(c.PostForm("voucher_date"))
		if _, err := time.Parse("2006-01-02", voucherDate); err != nil {
			Toast(c, "error", "Ngày lập phiếu không hợp lệ.")
			s.renderCashPaymentsTable(c, d)
			return
		}
		amount, err := strconv.ParseFloat(c.PostForm("amount"), 64)
		if err != nil || amount <= 0 {
			Toast(c, "error", "Số tiền không hợp lệ.")
			s.renderCashPaymentsTable(c, d)
			return
		}
		// Legacy page mirrors debit_account_id onto credit_account_id.
		credit := strings.TrimSpace(c.PostForm("credit_account_id"))
		debit := strings.TrimSpace(c.PostForm("debit_account_id"))
		if credit == "" {
			credit = debit
		}
		paymentType := domain.PaymentType(c.PostForm("payment_type"))
		if paymentType == "" {
			paymentType = domain.PaymentSupplier
		}
		p := &domain.CashPayment{
			CompanyID:       pageCompanyID(c),
			VoucherDate:     voucherDate,
			CashAccountID:   strings.TrimSpace(c.PostForm("cash_account_id")),
			PayeeName:       strings.TrimSpace(c.PostForm("payee_name")),
			PayeeType:       domain.CounterpartSupplier,
			Currency:        "VND",
			ExchangeRate:    1,
			Amount:          amount,
			DebitAccountID:  debit,
			CreditAccountID: credit,
			Reason:          strings.TrimSpace(c.PostForm("reason")),
			PaymentType:     paymentType,
		}
		// CreateCashPayment runs p.Validate() + voucher numbering.
		if err := d.Svc.CreateCashPayment(c.Request.Context(), p); err != nil {
			log.Printf("create cash payment: %v", err)
			Toast(c, "error", "Không tạo được phiếu chi: "+err.Error())
			s.renderCashPaymentsTable(c, d)
			return
		}
		Toast(c, "success", "Đã tạo phiếu chi "+p.VoucherNo+".")
		s.renderCashPaymentsTable(c, d)
	}
}

func (s *Server) cashPaymentStatus(d Deps, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.PostForm("id")
		ctx := c.Request.Context()
		var err error
		switch action {
		case "submit":
			err = d.Svc.SubmitCashPayment(ctx, id, c.GetString("user_id"))
		case "approve":
			err = d.Svc.ApproveCashPayment(ctx, id, c.GetString("user_id"))
		case "post":
			err = d.Svc.PostCashPayment(ctx, id, c.GetString("user_id"))
		}
		if err != nil {
			log.Printf("cash payment %s %s: %v", action, id, err)
			Toast(c, "error", "Thao tác thất bại: "+err.Error())
		} else {
			Toast(c, "success", "Đã cập nhật trạng thái phiếu chi.")
		}
		s.renderCashPaymentsTable(c, d)
	}
}

func (s *Server) renderCashPaymentsTable(c *gin.Context, d Deps) {
	payments, total, err := d.Svc.ListCashPayments(c.Request.Context(), domain.CashPaymentFilter{CompanyID: pageCompanyID(c)})
	if err != nil {
		c.String(500, "load cash payments failed")
		return
	}
	if payments == nil {
		payments = []domain.CashPayment{}
	}
	s.RenderFragment(c, "cash-payments", "cash-payments-table", gin.H{"Payments": payments, "Total": total})
}

// ─── Cash Transfers ───────────────────────────────────────────────

func cashTransfersLoad(d Deps) func(c *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
		transfers, err := d.Svc.GetCashTransfers(c.Request.Context(), pageCompanyID(c))
		if err != nil {
			return nil, err
		}
		if transfers == nil {
			transfers = []domain.CashTransfer{}
		}
		return gin.H{"Transfers": transfers}, nil
	}
}

func (s *Server) cashTransfersCreate(d Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		transferDate := strings.TrimSpace(c.PostForm("transfer_date"))
		if _, err := time.Parse("2006-01-02", transferDate); err != nil {
			Toast(c, "error", "Ngày chuyển không hợp lệ.")
			s.renderCashTransfersTable(c, d)
			return
		}
		amount, err := strconv.ParseFloat(c.PostForm("amount"), 64)
		if err != nil || amount <= 0 {
			Toast(c, "error", "Số tiền không hợp lệ.")
			s.renderCashTransfersTable(c, d)
			return
		}
		transferType := domain.TransferType(c.PostForm("transfer_type"))
		if transferType == "" {
			transferType = domain.TransferBankWithdrawal
		}
		t := &domain.CashTransfer{
			CompanyID:     pageCompanyID(c),
			TransferDate:  transferDate,
			FromAccountID: strings.TrimSpace(c.PostForm("from_account_id")),
			ToAccountID:   strings.TrimSpace(c.PostForm("to_account_id")),
			Amount:        amount,
			Currency:      "VND",
			ExchangeRate:  1,
			Reason:        strings.TrimSpace(c.PostForm("reason")),
			TransferType:  transferType,
		}
		// CreateCashTransfer mints vouchers R-/P- + journal entry, posts all.
		if err := d.Svc.CreateCashTransfer(c.Request.Context(), t, c.GetString("user_id")); err != nil {
			log.Printf("create cash transfer: %v", err)
			Toast(c, "error", "Không tạo được lệnh chuyển quỹ: "+err.Error())
			s.renderCashTransfersTable(c, d)
			return
		}
		Toast(c, "success", "Đã tạo lệnh chuyển quỹ "+t.ID+".")
		s.renderCashTransfersTable(c, d)
	}
}

func (s *Server) renderCashTransfersTable(c *gin.Context, d Deps) {
	transfers, err := d.Svc.GetCashTransfers(c.Request.Context(), pageCompanyID(c))
	if err != nil {
		c.String(500, "load cash transfers failed")
		return
	}
	if transfers == nil {
		transfers = []domain.CashTransfer{}
	}
	s.RenderFragment(c, "cash-transfers", "cash-transfers-table", gin.H{"Transfers": transfers})
}
