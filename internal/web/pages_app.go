package web

import (
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
	"gotax/internal/service"
)

// Deps bundles the services pages render from.
type Deps struct {
	Svc service.Service
}

// NewPages returns the registry of server-rendered pages → loaders.
func NewPages(d Deps) map[string]Page {
	return map[string]Page{
		"/app/dashboard.html":       {Title: "Tổng quan", NavPath: "/app/dashboard.html", Load: dashboardLoad(d)},
		"/app/users.html":           {Title: "Người dùng", NavPath: "/app/users.html", Load: usersLoad(d)},
		"/app/journal-entries.html": {Title: "Chứng từ kế toán", NavPath: "/app/journal-entries.html", Load: journalEntriesLoad(d)},
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
			"Accounts":         len(accs),
			"Entries":          counts[domain.JournalEntryPosted],
			"StatusStats":      stats,
			"Recent":           recent,
			"RecentTotalDebit": recentDebit,
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
