package service

import (
	"context"
	"fmt"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"

	"gotax/internal/domain"
)

func (s *service) GenerateOpeningBalancePDF(ctx context.Context, companyID, periodID string) ([]byte, error) {
	filter := domain.OBListFilter{CompanyID: companyID, PeriodID: periodID}
	balances, err := s.ob.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list balances: %w", err)
	}

	cfg := config.NewBuilder().
		WithPageSize("A4").
		WithLeftMargin(15).
		WithTopMargin(15).
		WithRightMargin(15).
		WithBottomMargin(15).
		Build()

	m := maroto.New(cfg)

	m.RegisterHeader(buildReportHeader(companyID, periodID)...)

	for _, b := range balances {
		m.AddRows(buildBalanceRow(b)...)
	}

	debit, credit, _ := s.ob.GetTotals(ctx, companyID, periodID)
	m.AddRows(buildTotalsRow(debit, credit)...)

	doc, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("generate pdf: %w", err)
	}
	return doc.GetBytes(), nil
}

func buildReportHeader(companyID, periodID string) []core.Row {
	return []core.Row{
		row.New(8).Add(
			text.NewCol(12, fmt.Sprintf("Opening Balance Report — %s / %s", companyID, periodID),
				props.Text{Style: fontstyle.Bold, Size: 14, Align: align.Center, Top: 2, Bottom: 4}),
		),
		row.New(1).Add(
			col.New(12),
		),
		row.New(6).Add(
			text.NewCol(2, "Account", props.Text{Style: fontstyle.Bold, Size: 9}),
			text.NewCol(2, "Currency", props.Text{Style: fontstyle.Bold, Size: 9}),
			text.NewCol(2, "Debit", props.Text{Style: fontstyle.Bold, Size: 9, Align: align.Right}),
			text.NewCol(2, "Credit", props.Text{Style: fontstyle.Bold, Size: 9, Align: align.Right}),
			text.NewCol(2, "Source", props.Text{Style: fontstyle.Bold, Size: 9}),
			text.NewCol(2, "Status", props.Text{Style: fontstyle.Bold, Size: 9}),
		),
		row.New(1).Add(col.New(12)),
	}
}

func buildBalanceRow(b domain.OpeningBalance) []core.Row {
	return []core.Row{
		row.New(5).Add(
			text.NewCol(2, b.AccountCode, props.Text{Size: 8}),
			text.NewCol(2, b.CurrencyCode, props.Text{Size: 8}),
			text.NewCol(2, fmt.Sprintf("%.0f", b.DebitAmount), props.Text{Size: 8, Align: align.Right}),
			text.NewCol(2, fmt.Sprintf("%.0f", b.CreditAmount), props.Text{Size: 8, Align: align.Right}),
			text.NewCol(2, b.SourceType, props.Text{Size: 8}),
			text.NewCol(2, string(b.Status), props.Text{Size: 8}),
		),
	}
}

func buildTotalsRow(debit, credit float64) []core.Row {
	return []core.Row{
		row.New(1).Add(col.New(12)),
		row.New(6).Add(
			text.NewCol(4, "", props.Text{Size: 8}),
			text.NewCol(2, "Total:", props.Text{Style: fontstyle.Bold, Size: 9, Align: align.Right}),
			text.NewCol(2, fmt.Sprintf("%.0f", debit), props.Text{Style: fontstyle.Bold, Size: 9, Align: align.Right}),
			text.NewCol(2, fmt.Sprintf("%.0f", credit), props.Text{Style: fontstyle.Bold, Size: 9, Align: align.Right}),
			text.NewCol(2, "", props.Text{Size: 8}),
		),
	}
}

type CashFlowResult struct {
	OpeningCash float64         `json:"opening_cash"`
	Operating   CashFlowSection `json:"operating"`
	Investing   CashFlowSection `json:"investing"`
	Financing   CashFlowSection `json:"financing"`
	NetChange   float64         `json:"net_change"`
	ClosingCash float64         `json:"closing_cash"`
}

type CashFlowSection struct {
	Inflows  []CashFlowLine `json:"inflows"`
	Outflows []CashFlowLine `json:"outflows"`
	Net      float64        `json:"net"`
}

type CashFlowLine struct {
	AccountCode string  `json:"account_code"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

func (s *service) CashFlowStatement(ctx context.Context, companyID string, year, month int) (*CashFlowResult, error) {
	period, err := s.periods.GetByYearMonth(ctx, year, month)
	if err != nil {
		return nil, fmt.Errorf("get period: %w", err)
	}
	if period == nil {
		return nil, fmt.Errorf("period not found for %d/%02d", year, month)
	}

	// Opening cash: balance of accounts 111, 112 at end of previous month.
	var openingCash float64
	prevYear, prevMonth := year, month-1
	if prevMonth < 1 {
		prevMonth = 12
		prevYear--
	}
	prevPeriod, err := s.periods.GetByYearMonth(ctx, prevYear, prevMonth)
	if err == nil && prevPeriod != nil {
		for _, code := range []string{"111", "112"} {
			b, err := s.journals.GetBalance(ctx, code, prevPeriod.ID)
			if err == nil && b != nil {
				b.Calculate()
				openingCash += b.ClosingBalance
			}
		}
	}

	// Fetch all posted entries for the period.
	entries, err := s.journals.GetByPeriod(ctx, period.ID)
	if err != nil {
		return nil, fmt.Errorf("get entries: %w", err)
	}

	// Fetch account metadata for descriptions.
	acctMap := make(map[string]string)
	if accounts, err := s.accounts.GetAll(ctx, false); err == nil {
		for _, a := range accounts {
			acctMap[a.Code] = a.Name
		}
	}

	type lineAgg struct {
		inflow  float64
		outflow float64
	}
	operatingLines := make(map[string]*lineAgg)
	investingLines := make(map[string]*lineAgg)
	financingLines := make(map[string]*lineAgg)

	for _, e := range entries {
		if e.Status != domain.JournalEntryPosted {
			continue
		}
		// Check if this entry touches a cash account (111, 112).
		hasCash := false
		for _, l := range e.Lines {
			if len(l.AccountCode) >= 3 && (l.AccountCode[:3] == "111" || l.AccountCode[:3] == "112") {
				hasCash = true
				break
			}
		}
		if !hasCash {
			continue
		}

		for _, l := range e.Lines {
			code := l.AccountCode
			prefix := code
			if len(code) > 3 {
				prefix = code[:3]
			}
			// Skip cash accounts themselves.
			if prefix == "111" || prefix == "112" {
				continue
			}

			var target map[string]*lineAgg
			switch {
			case prefix == "511" || prefix == "515" || prefix == "711" ||
				prefix == "632" || prefix == "641" || prefix == "642" ||
				prefix == "131" || prefix == "152" || prefix == "331":
				target = operatingLines
			case prefix == "211" || prefix == "128" || prefix == "228":
				target = investingLines
			case prefix == "411":
				target = financingLines
			default:
				continue
			}

			ag, ok := target[code]
			if !ok {
				ag = &lineAgg{}
				target[code] = ag
			}
			ag.inflow += l.CreditAmount
			ag.outflow += l.DebitAmount
		}
	}

	operating := CashFlowSection{}
	investing := CashFlowSection{}
	financing := CashFlowSection{}

	for code, ag := range operatingLines {
		desc := acctMap[code]
		if desc == "" {
			desc = code
		}
		if ag.inflow > 0 {
			operating.Inflows = append(operating.Inflows, CashFlowLine{AccountCode: code, Description: desc, Amount: ag.inflow})
		}
		if ag.outflow > 0 {
			operating.Outflows = append(operating.Outflows, CashFlowLine{AccountCode: code, Description: desc, Amount: ag.outflow})
		}
	}

	for code, ag := range investingLines {
		desc := acctMap[code]
		if desc == "" {
			desc = code
		}
		if ag.inflow > 0 {
			investing.Inflows = append(investing.Inflows, CashFlowLine{AccountCode: code, Description: desc, Amount: ag.inflow})
		}
		if ag.outflow > 0 {
			investing.Outflows = append(investing.Outflows, CashFlowLine{AccountCode: code, Description: desc, Amount: ag.outflow})
		}
	}

	for code, ag := range financingLines {
		desc := acctMap[code]
		if desc == "" {
			desc = code
		}
		if ag.inflow > 0 {
			financing.Inflows = append(financing.Inflows, CashFlowLine{AccountCode: code, Description: desc, Amount: ag.inflow})
		}
		if ag.outflow > 0 {
			financing.Outflows = append(financing.Outflows, CashFlowLine{AccountCode: code, Description: desc, Amount: ag.outflow})
		}
	}

	for i := range operating.Inflows {
		operating.Net += operating.Inflows[i].Amount
	}
	for i := range operating.Outflows {
		operating.Net -= operating.Outflows[i].Amount
	}
	for i := range investing.Inflows {
		investing.Net += investing.Inflows[i].Amount
	}
	for i := range investing.Outflows {
		investing.Net -= investing.Outflows[i].Amount
	}
	for i := range financing.Inflows {
		financing.Net += financing.Inflows[i].Amount
	}
	for i := range financing.Outflows {
		financing.Net -= financing.Outflows[i].Amount
	}

	netChange := operating.Net + investing.Net + financing.Net

	return &CashFlowResult{
		OpeningCash: openingCash,
		Operating:   operating,
		Investing:   investing,
		Financing:   financing,
		NetChange:   netChange,
		ClosingCash: openingCash + netChange,
	}, nil
}
