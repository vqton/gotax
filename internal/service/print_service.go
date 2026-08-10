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

// PrintType constants
const (
	PrintTypeReceiptVoucher = "receipt"
	PrintTypePaymentVoucher = "payment"
)

// GeneratePrintPDF generates a PDF for the given type and ID.
func (s *service) GeneratePrintPDF(ctx context.Context, printType, id string) ([]byte, error) {
	switch printType {
	case PrintTypeReceiptVoucher:
		return s.generateReceiptPDF(ctx, id)
	case PrintTypePaymentVoucher:
		return s.generatePaymentPDF(ctx, id)
	default:
		return nil, fmt.Errorf("unsupported print type: %s", printType)
	}
}

func (s *service) generateReceiptPDF(ctx context.Context, id string) ([]byte, error) {
	receipt, err := s.cash.GetReceipt(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get receipt: %w", err)
	}

	cfg := config.NewBuilder().
		WithPageSize("A4").
		WithLeftMargin(15).
		WithTopMargin(15).
		WithRightMargin(15).
		WithBottomMargin(15).
		Build()

	m := maroto.New(cfg)
	m.RegisterHeader(buildVoucherHeader("PHIẾU THU", receipt.VoucherNo, receipt.VoucherDate, receipt.CompanyID)...)
	m.AddRows(buildReceiptBody(receipt)...)
	m.AddRows(buildSignatureRow()...)

	doc, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("generate pdf: %w", err)
	}
	return doc.GetBytes(), nil
}

func (s *service) generatePaymentPDF(ctx context.Context, id string) ([]byte, error) {
	payment, err := s.cash.GetPayment(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get payment: %w", err)
	}

	cfg := config.NewBuilder().
		WithPageSize("A4").
		WithLeftMargin(15).
		WithTopMargin(15).
		WithRightMargin(15).
		WithBottomMargin(15).
		Build()

	m := maroto.New(cfg)
	m.RegisterHeader(buildVoucherHeader("PHIẾU CHI", payment.VoucherNo, payment.VoucherDate, payment.CompanyID)...)
	m.AddRows(buildPaymentBody(payment)...)
	m.AddRows(buildSignatureRow()...)

	doc, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("generate pdf: %w", err)
	}
	return doc.GetBytes(), nil
}

// ─── Shared Layout ────────────────────────────────────────────────────

func buildVoucherHeader(title, voucherNo, voucherDate, companyID string) []core.Row {
	return []core.Row{
		row.New(10).Add(
			text.NewCol(12, title,
				props.Text{Style: fontstyle.Bold, Size: 16, Align: align.Center, Top: 2}),
		),
		row.New(1).Add(col.New(12)),
		row.New(6).Add(
			text.NewCol(6, fmt.Sprintf("Số: %s", voucherNo), props.Text{Size: 10}),
			text.NewCol(6, fmt.Sprintf("Ngày: %s", voucherDate), props.Text{Size: 10, Align: align.Right}),
		),
		row.New(6).Add(
			text.NewCol(12, fmt.Sprintf("Đơn vị: %s", companyID), props.Text{Size: 10}),
		),
		row.New(1).Add(col.New(12)),
	}
}

func buildReceiptBody(r *domain.CashReceipt) []core.Row {
	return []core.Row{
		row.New(6).Add(
			text.NewCol(6, fmt.Sprintf("Họ tên người nộp: %s", r.CounterpartName), props.Text{Size: 10}),
			text.NewCol(6, fmt.Sprintf("Đối tượng: %s", r.CounterpartType), props.Text{Size: 10}),
		),
		row.New(6).Add(
			text.NewCol(6, fmt.Sprintf("Lý do: %s", r.Reason), props.Text{Size: 10}),
			text.NewCol(6, fmt.Sprintf("Tiền: %s", r.Currency), props.Text{Size: 10}),
		),
		row.New(1).Add(col.New(12)),
		row.New(6).Add(
			text.NewCol(6, "Thu", props.Text{Style: fontstyle.Bold, Size: 10}),
			text.NewCol(6, fmt.Sprintf("%s %.0f", r.Currency, r.Amount), props.Text{Style: fontstyle.Bold, Size: 10, Align: align.Right}),
		),
		row.New(1).Add(col.New(12)),
		row.New(6).Add(
			text.NewCol(6, fmt.Sprintf("Nợ: %s", r.DebitAccountID), props.Text{Size: 10}),
			text.NewCol(6, fmt.Sprintf("Có: %s", r.CreditAccountID), props.Text{Size: 10}),
		),
		row.New(1).Add(col.New(12)),
		row.New(6).Add(
			text.NewCol(6, fmt.Sprintf("Tỷ giá: %.0f", r.ExchangeRate), props.Text{Size: 10}),
			text.NewCol(6, fmt.Sprintf("Quy đổi: %.0f VND", r.AmountVND), props.Text{Size: 10}),
		),
	}
}

func buildPaymentBody(p *domain.CashPayment) []core.Row {
	return []core.Row{
		row.New(6).Add(
			text.NewCol(6, fmt.Sprintf("Họ tên người nhận: %s", p.PayeeName), props.Text{Size: 10}),
			text.NewCol(6, fmt.Sprintf("Đối tượng: %s", p.PayeeType), props.Text{Size: 10}),
		),
		row.New(6).Add(
			text.NewCol(6, fmt.Sprintf("Lý do: %s", p.Reason), props.Text{Size: 10}),
			text.NewCol(6, fmt.Sprintf("Tiền: %s", p.Currency), props.Text{Size: 10}),
		),
		row.New(1).Add(col.New(12)),
		row.New(6).Add(
			text.NewCol(6, "Chi", props.Text{Style: fontstyle.Bold, Size: 10}),
			text.NewCol(6, fmt.Sprintf("%s %.0f", p.Currency, p.Amount), props.Text{Style: fontstyle.Bold, Size: 10, Align: align.Right}),
		),
		row.New(1).Add(col.New(12)),
		row.New(6).Add(
			text.NewCol(6, fmt.Sprintf("Nợ: %s", p.DebitAccountID), props.Text{Size: 10}),
			text.NewCol(6, fmt.Sprintf("Có: %s", p.CreditAccountID), props.Text{Size: 10}),
		),
		row.New(1).Add(col.New(12)),
		row.New(6).Add(
			text.NewCol(6, fmt.Sprintf("Tỷ giá: %.0f", p.ExchangeRate), props.Text{Size: 10}),
			text.NewCol(6, fmt.Sprintf("Quy đổi: %.0f VND", p.AmountVND), props.Text{Size: 10}),
		),
	}
}

func buildSignatureRow() []core.Row {
	return []core.Row{
		row.New(1).Add(col.New(12)),
		row.New(20).Add(
			text.NewCol(4, "Người lập phiếu", props.Text{Style: fontstyle.Bold, Size: 10, Align: align.Center}),
			text.NewCol(4, "Thủ trưởng đơn vị", props.Text{Style: fontstyle.Bold, Size: 10, Align: align.Center}),
			text.NewCol(4, "Thủ quỹ", props.Text{Style: fontstyle.Bold, Size: 10, Align: align.Center}),
		),
		row.New(20).Add(
			text.NewCol(4, "(Ký, ghi rõ họ tên)", props.Text{Size: 8, Align: align.Center}),
			text.NewCol(4, "(Ký, ghi rõ họ tên)", props.Text{Size: 8, Align: align.Center}),
			text.NewCol(4, "(Ký, ghi rõ họ tên)", props.Text{Size: 8, Align: align.Center}),
		),
	}
}
