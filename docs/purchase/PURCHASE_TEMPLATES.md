# Purchase Module — Document Templates

**Version:** 1.0
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)
**Compliant with:** Circular 99/2025/TT-BTC, Decree 123/2020

---

## 1. Purchase Order (Đơn đặt hàng)

```
┌─────────────────────────────────────────────────────────────────────┐
│                       ĐƠN ĐẶT HÀNG                                 │
│                       PURCHASE ORDER                               │
│                                                                     │
│  Số/No: PO-202607-0001                         Ngày/Date: 15/07/2026  │
├─────────────────────────────────────────────────────────────────────┤
│ NHÀ CUNG CẤP / SUPPLIER:                 DOANH NGHIỆP / COMPANY:    │
│   Mã/Code: NCC001                         Công ty ABC               │
│   Tên/Name: Công ty XYZ                   MST/Tax: 0123456789       │
│   MST/Tax: 9876543210                     Địa chỉ: 123 Đường ABC    │
│   Địa chỉ: 456 Đường DEF                  Điện thoại: 028.3822.111 │
│                                                                     │
│  Hạn giao / Delivery: 25/07/2026                                    │
│  Hình thức thanh toán / Payment: Net 30                             │
│  Điều kiện giao hàng / Delivery terms: DDP                         │
│  Loại tiền / Currency: VND                                          │
├─────────────────────────────────────────────────────────────────────┤
│ STT │ Mã hàng │ Diễn giải          │ ĐVT │ Số lượng │ Đơn giá │    │
│     │         │ (Description)      │     │ (Qty)    │ (Price) │    │
├─────┼─────────┼─────────────────────┼─────┼──────────┼─────────┤    │
│  1  │ SP001   │ Máy tính Dell 5420 │ Cái │   10     │15,000,000│   │
│  2  │ SP002   │ Màn hình Dell 27"  │ Cái │   10     │ 5,000,000│   │
├─────┼─────────┴─────────────────────┴─────┴──────────┴─────────┤    │
│     │ Thuế suất VAT: 8%                                        │   │
├─────┼───────────────────────────────────────────────────────────┤   │
│     │ Tổng tiền hàng / Subtotal:            200,000,000          │   │
│     │ Chiết khấu / Discount:                         0          │   │
│     │ Tiền thuế / VAT (8%):                  16,000,000          │   │
│     │ Tổng cộng / Grand Total:              216,000,000          │   │
│     │ Số tiền bằng chữ: Hai trăm mười sáu triệu đồng             │   │
├─────────────────────────────────────────────────────────────────┤   │
│ Người lập         Người phê duyệt        Người đại diện          │   │
│ (Prepared by)      (Approved by)          (Authorized)            │   │
│  ___________       ___________            ___________             │   │
└─────────────────────────────────────────────────────────────────┘
```

---

## 2. Goods Receipt Note (Phiếu nhập kho)

```
┌─────────────────────────────────────────────────────────────────────┐
│                    PHIẾU NHẬP KHO                                  │
│                    GOODS RECEIPT NOTE                              │
│                                                                     │
│  Số/No: GRN-202607-0001                    Ngày/Date: 20/07/2026     │
├─────────────────────────────────────────────────────────────────────┤
│ Đơn đặt hàng / PO: PO-202607-0001                                  │
│ Nhà cung cấp / Supplier: Công ty XYZ                               │
│ Kho / Warehouse: KHO01 (Kho chính)                                 │
├─────────────────────────────────────────────────────────────────────┤
│ STT │ Mã hàng │ Diễn giải            │ ĐVT │ SL nhập │ SL hỏng │   │
│     │         │                      │     │ (OK)    │ (Reject)│   │
├─────┼─────────┼──────────────────────┼─────┼─────────┼─────────┤   │
│  1  │ SP001   │ Máy tính Dell 5420  │ Cái │   10    │    0    │   │
│  2  │ SP002   │ Màn hình Dell 27"   │ Cái │   10    │    0    │   │
├─────┼─────────┴──────────────────────┴─────┴─────────┴─────────┤   │
│     │ Tổng số chứng từ: 02                                     │   │
│     │ Ghi chú: Hàng đạt chất lượng                              │   │
├─────────────────────────────────────────────────────────────────┤   │
│ Người giao hàng     Người nhận hàng       Thủ kho                │   │
│ (Deliverer)         (Receiver)            (Warehouse Keeper)     │   │
│  ___________        ___________           ___________            │   │
└─────────────────────────────────────────────────────────────────┘
```

---

## 3. Supplier Invoice Recording (Hóa đơn mua hàng)

```
┌─────────────────────────────────────────────────────────────────────┐
│               GHI NHẬN HÓA ĐƠN MUA HÀNG                            │
│               SUPPLIER INVOICE RECORDING                           │
│                                                                     │
│  Hóa đơn số/Invoice#: INV-12345                 Ngày/Date: 20/07/2026│
│  Mã số thuế/MST người bán: 9876543210                               │
│  Mã số thuế/MST người mua: 0123456789                               │
│                                                                     │
│  Nhà cung cấp / Supplier: Công ty XYZ                               │
│  Đơn đặt hàng / PO: PO-202607-0001                                  │
│  Phiếu nhập kho / GRN: GRN-202607-0001                              │
│                                                                     │
│  Loại hóa đơn: ☐ Có mã GDT  ☒ Không mã GDT                         │
│  Mã của CQT (nếu có): _______________                               │
│  Tình trạng khấu trừ VAT: ☒ Được khấu trừ  ☐ Không được khấu trừ   │
├─────────────────────────────────────────────────────────────────────┤
│ STT │ Diễn giải        │ ĐVT │ SL   │ Đơn giá   │ VAT% │ Thành tiền│
├─────┼──────────────────┼─────┼──────┼───────────┼──────┼───────────┤
│  1  │ Máy tính Dell   │ Cái │  10  │15,000,000 │  8   │150,000,000│
│  2  │ Màn hình Dell   │ Cái │  10  │ 5,000,000 │  8   │ 50,000,000│
├─────┼──────────────────┴─────┴──────┴───────────┴──────┴───────────┤
│     │ Cộng tiền hàng / Subtotal:                      200,000,000  │
│     │ Thuế GTGT / VAT (8%):                           16,000,000  │
│     │ Tổng cộng thanh toán / Grand Total:             216,000,000  │
├─────────────────────────────────────────────────────────────────────┤
│ Hạch toán / GL Posting:                                            │
│   Nợ 156 (Hàng hóa)          200,000,000                            │
│   Nợ 1331 (Thuế GTGT)         16,000,000                            │
│   Có 331 (PT người bán)      216,000,000                            │
├─────────────────────────────────────────────────────────────────────┤
│ Người lập             Kế toán trưởng           Thủ trưởng đơn vị   │
│ (Prepared)            (Chief Accountant)       (Director)           │
│  ___________          ___________              ___________          │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 4. AP Aging Report (Bảng phân tích tuổi nợ phải trả)

```
┌──────────────────────────────────────────────────────────────────────────┐
│          BẢNG PHÂN TÍCH TUỔI NỢ PHẢI TRẢ                                 │
│          ACCOUNTS PAYABLE AGING REPORT                                   │
│                                                                          │
│  Ngày/Date: 31/07/2026                           Đơn vị/VND              │
├──────────┬──────────┬──────────┬──────────┬──────────┬──────────┬────────┤
│ Nhà CC   │ Tổng nợ  │ Chưa đến │ 1-30     │ 31-60    │ 61-90    │ >90    │
│ (Supplier│ (Total)  │ hạn      │ ngày quá │ ngày quá │ ngày quá │ ngày   │
│          │          │ (Current)│ hạn      │ hạn      │ hạn      │ quá hạn│
├──────────┼──────────┼──────────┼──────────┼──────────┼──────────┼────────┤
│ XYZ Corp │216,000,000│0        │216,000,000│0        │0        │0       │
│ ABC Ltd  │550,000,000│200,000,000│150,000,000│100,000,000│50,000,000│50,000,000│
│ DEF Co   │120,000,000│120,000,000│0        │0        │0        │0       │
├──────────┼──────────┼──────────┼──────────┼──────────┼──────────┼────────┤
│ TỔNG     │886,000,000│320,000,000│366,000,000│100,000,000│50,000,000│50,000,000│
│          │          │36%      │41%      │11%      │6%       │6%     │
└──────────┴──────────┴──────────┴──────────┴──────────┴──────────┴────────┘
```

---

## 5. 3-Way Matching Worksheet

```
┌────────────────────────────────────────────────────────────────────────────┐
│                   3-WAY MATCHING WORKSHEET                                 │
├────────────────────────────────────────────────────────────────────────────┤
│ PO: PO-202607-0001      GRN: GRN-202607-0001      INV: INV-12345           │
├────────────────────────────────────────────────────────────────────────────┤
│ Item    │ PO Qty │ PO Price │ GRN Qty │ INV Qty │ INV Price │ Match Status│
├─────────┼────────┼──────────┼─────────┼─────────┼───────────┼─────────────┤
│ SP001   │  10    │15,000,000│   10    │   10    │15,000,000 │ ✓ MATCH      │
│ SP002   │  10    │ 5,000,000│   10    │   10    │ 5,000,000 │ ✓ MATCH      │
├─────────┼────────┼──────────┼─────────┼─────────┼───────────┼─────────────┤
│ Result: │ PO Qty (20) ≥ GRN Qty (20) ≥ INV Qty (20) → PASS                 │
│         │ Price variance: 0% → PASS                                         │
│         │ OVERALL: ✓ VERIFIED                                                │
└────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. S01-DN Purchase Ledger (Sổ chi tiết mua hàng)

```
┌────────────────────────────────────────────────────────────────────────────┐
│              SỔ CHI TIẾT MUA HÀNG                                          │
│              PURCHASE LEDGER                                               │
│                                                                           │
│  Tài khoản/Account: 156 (Hàng hóa)                                         │
│  Kỳ/Period: Tháng 07/2026                                                  │
├──────┬──────────┬──────────┬──────────┬──────────┬──────────┬──────────────┤
│ Ngày │ Số CT    │ Diễn giải│ Nhà CC   │ TK ĐƯ   │ PS Tăng   │ PS Giảm      │
│(Date)│ (Doc#)   │ (Desc)   │(Supplier)│ (O/A)    │ (Increase)│ (Decrease)   │
├──────┼──────────┼──────────┼──────────┼──────────┼───────────┼──────────────┤
│20/07 │GRN-0001  │ Nhập kho │ XYZ      │331       │200,000,000│              │
│      │          │ PO-0001  │          │          │           │              │
│22/07 │INV-12345│ HĐ mua   │ XYZ      │331       │           │200,000,000   │
│      │          │ (điều chỉnh)│        │          │           │(điều ch. giá)│
├──────┼──────────┼──────────┼──────────┼──────────┼───────────┼──────────────┤
│      │          │ Dư đầu kỳ│          │          │  500,000,000│            │
│      │          │ PS trong │          │          │  200,000,000│   200,000,000│
│      │          │ Dư cuối  │          │          │  500,000,000│            │
└──────┴──────────┴──────────┴──────────┴──────────┴───────────┴──────────────┘
```

---

## 7. S02-DN Supplier Detail Ledger (Sổ chi tiết công nợ phải trả)

```
┌────────────────────────────────────────────────────────────────────────────┐
│              SỔ CHI TIẾT THANH TOÁN VỚI NGƯỜI BÁN                         │
│              SUPPLIER DETAIL LEDGER                                       │
│                                                                           │
│  Tài khoản/Account: 331 (Phải trả người bán)                               │
│  Nhà cung cấp/Supplier: XYZ Corp (MST: 9876543210)                         │
│  Kỳ/Period: Tháng 07/2026                                                  │
├──────┬──────────┬──────────┬──────────┬──────────┬──────────┬──────────────┤
│ Ngày │ Số CT    │ Diễn giải│ PS Nợ    │ PS Có    │ Dư Nợ    │ Dư Có        │
│(Date)│ (Doc#)   │ (Desc)   │ (Debit)  │ (Credit) │(Debit Bal)│(Credit Bal) │
├──────┼──────────┼──────────┼──────────┼──────────┼───────────┼──────────────┤
│01/07 │          │ Dư đầu kỳ│          │          │         0│    500,000,000│
│20/07 │GRN-0001  │ Nhập kho │          │200,000,000│         0│    700,000,000│
│22/07 │INV-12345│ HĐ mua    │200,000,000│216,000,000│         0│    716,000,000│
│      │          │ (điều chỉnh)│          │          │          │              │
│31/07 │UNC-0001  │ TT tiền   │216,000,000│          │          │    500,000,000│
├──────┼──────────┼──────────┼──────────┼──────────┼───────────┼──────────────┤
│      │          │ Dư cuối  │          │          │         0│    500,000,000│
└──────┴──────────┴──────────┴──────────┴──────────┴───────────┴──────────────┘
```

---

## 8. E-Invoice XML Structure (Decree 254/2026 Compliant)

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Invoice xmlns="http://gdt.gov.vn/schemas/einvoice/2026">
  <InvoiceHeader>
    <InvoiceNumber>INV-12345</InvoiceNumber>
    <InvoiceDate>2026-07-20</InvoiceDate>
    <InvoiceTime>14:30:00</InvoiceTime>
    <InvoiceType>01GTKT</InvoiceType>
    <InvoiceSeries>AA/26E</InvoiceSeries>
    <CurrencyCode>VND</CurrencyCode>
    <ExchangeRate>1</ExchangeRate>
  </InvoiceHeader>
  <Seller>
    <TaxCode>9876543210</TaxCode>
    <Name>Công ty XYZ</Name>
    <Address>456 Đường DEF, Quận 1, TP.HCM</Address>
    <Phone>0283822333</Phone>
    <BankAccount>0123456789</BankAccount>
    <BankName>Vietcombank</BankName>
  </Seller>
  <Buyer>
    <TaxCode>0123456789</TaxCode>
    <Name>Công ty ABC</Name>
    <Address>123 Đường ABC, Quận 1, TP.HCM</Address>
  </Buyer>
  <InvoiceLines>
    <Line>
      <LineNumber>1</LineNumber>
      <ItemCode>SP001</ItemCode>
      <ItemName>Máy tính Dell 5420</ItemName>
      <Unit>Cái</Unit>
      <Quantity>10</Quantity>
      <UnitPrice>15000000</UnitPrice>
      <VatRate>8</VatRate>
      <VatAmount>12000000</VatAmount>
      <LineTotal>150000000</LineTotal>
    </Line>
    <Line>
      <LineNumber>2</LineNumber>
      <ItemCode>SP002</ItemCode>
      <ItemName>Màn hình Dell 27"</ItemName>
      <Unit>Cái</Unit>
      <Quantity>10</Quantity>
      <UnitPrice>5000000</UnitPrice>
      <VatRate>8</VatRate>
      <VatAmount>4000000</VatAmount>
      <LineTotal>50000000</LineTotal>
    </Line>
  </InvoiceLines>
  <Summary>
    <Subtotal>200000000</Subtotal>
    <Discount>0</Discount>
    <VatAmount>16000000</VatAmount>
    <GrandTotal>216000000</GrandTotal>
    <AmountInWords>Hai trăm mười sáu triệu đồng</AmountInWords>
  </Summary>
  <DigitalSignature>
    <Signature>base64-encoded-signature</Signature>
    <SignDate>2026-07-20</SignDate>
    <CertificateSerial>SN-123456</CertificateSerial>
  </DigitalSignature>
  <QRCode>base64-qr-code-image</QRCode>
</Invoice>
```

---

## 9. Credit Note Template (Trả lại hàng mua)

```
┌─────────────────────────────────────────────────────────────────────┐
│                  HÓA ĐƠN ĐIỀU CHỈNH / CREDIT NOTE                  │
│                                                                     │
│  Số/No: CN-202607-0001                     Ngày/Date: 25/07/2026   │
│                                                                     │
│  Liên quan đến HĐ / Ref Invoice: INV-12345                          │
│  Lý do / Reason: Hàng không đúng quy cách                           │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│ STT │ Diễn giải                    │ SL │ Đơn giá  │ Thành tiền     │
├─────┼──────────────────────────────┼────┼──────────┼────────────────┤
│  1  │ Trả lại Máy tính Dell 5420  │  1 │15,000,000│ 15,000,000     │
├─────┼──────────────────────────────┴────┴──────────┴────────────────┤
│     │ Thuế GTGT / VAT (8%):                          1,200,000     │
│     │ Tổng cộng giảm / Total Reduction:             16,200,000     │
├─────────────────────────────────────────────────────────────────────┤
│ Hạch toán / GL Posting:                                            │
│   Nợ 331 (PT người bán)           16,200,000                        │
│   Có 156 (Hàng hóa)               15,000,000                        │
│   Có 1331 (Thuế GTGT)             1,200,000                         │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 10. Supplier Evaluation Scorecard (P2)

```
┌─────────────────────────────────────────────────────────────────────┐
│               SUPPLIER EVALUATION SCORECARD                         │
│               Supplier: XYZ Corp                    Q3/2026         │
├──────────────┬──────────┬──────────┬───────────┬───────────────────┤
│ Criteria     │ Weight   │ Score    │ Weighted  │ Notes              │
│              │          │ (1-5)    │ Score     │                    │
├──────────────┼──────────┼──────────┼───────────┼───────────────────┤
│ Quality      │ 40%      │ 4        │ 1.6       │ 2% defect rate    │
│ Price        │ 25%      │ 3        │ 0.75      │ Above market avg  │
│ Delivery     │ 20%      │ 5        │ 1.0       │ On-time 100%      │
│ Service      │ 10%      │ 4        │ 0.4       │ Responsive        │
│ Compliance   │ 5%       │ 5        │ 0.25      │ Full docs         │
├──────────────┼──────────┼──────────┼───────────┼───────────────────┤
│ TOTAL        │ 100%     │          │ 4.0/5.0   │ Good standing     │
└──────────────┴──────────┴──────────┴───────────┴───────────────────┘
```