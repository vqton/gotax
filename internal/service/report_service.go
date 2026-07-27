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
