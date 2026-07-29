# Sale Module — Document Templates

**Version:** 1.0
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)
**Regulatory Basis:** Circular 99/2025/TT-BTC, Decree 123/2020/ND-CP, Decree 254/2026/ND-CP

---

## 1. Sales Order (SO)

```
────────────────────────────────────────────────────────────────────
                         SALES ORDER
                        (ĐƠN ĐẶT HÀNG)

  SO Number:  SO-202607-0001                  Date: 15/07/2026
  Customer:   ABC Corporation                 Tax Code: 0123456789
  Address:    123 Nguyen Hue, Q1, HCMC        Customer ID: KH001
  Contact:    Mr. Tan                         Phone: 0909xxxxxx

  Payment Terms: Net 30 Days
  Delivery Terms: DDP (Delivered Duty Paid)
  Delivery Date: 20/07/2026

  # | Item Code | Description                | Unit | Qty   | Unit Price | Disc% | VAT% | Total
  --+-----------+----------------------------+------+-------+------------+-------+------+------------
  1 | SP001     | Office Chair X200          | pcs  | 100   | 1,200,000  | 5%    | 10%  | 114,000,000
  2 | SP002     | Desk Y500                  | pcs  | 50    | 2,500,000  | 0%    | 10%  | 133,375,000
  3 | SP003     | Filing Cabinet Z100        | pcs  | 20    | 800,000    | 10%   | 10%  | 15,552,000
  --+-----------+----------------------------+------+-------+------------+-------+------+------------
                                        Subtotal:                                   247,522,000
                                        Discount:                                    (4,400,000)
                                        VAT (10%):                                   24,752,200
                                    ─────────────────────────────────────────────────────────────
                                        TOTAL:                                      267,874,200

  Amount in words: Two hundred sixty-seven million eight hundred seventy-four thousand
  two hundred Vietnam Dong

  Notes: Delivery to main office, contact Mr. Tan before delivery

  ────────────────────────────────────────────────────────────────
  Created By:     ___________      Approved By:     ___________
  Sales Manager                    Manager
────────────────────────────────────────────────────────────────────
```

---

## 2. Delivery Note (DN)

```
────────────────────────────────────────────────────────────────────
                       DELIVERY NOTE
                    (PHIẾU XUẤT KHO / BIÊN BẢN GIAO HÀNG)

  DN Number:     DN-202607-0001                 Date: 18/07/2026
  SO Reference:  SO-202607-0001
  Customer:      ABC Corporation                Tax Code: 0123456789
  Delivery Addr: 123 Nguyen Hue, Q1, HCMC
  Warehouse:     WH-01-A (Main Warehouse)
  Shipping:      FastShip Courier

  # | Item Code | Description           | Unit | Qty Ordered | Qty Delivered | Unit Price | Total
  --+-----------+-----------------------+------+-------------+---------------+------------+--------
  1 | SP001     | Office Chair X200     | pcs  | 100         | 80            | 1,200,000  | 96,000,000
  2 | SP002     | Desk Y500             | pcs  | 50          | 50            | 2,500,000  | 125,000,000
  3 | SP003     | Filing Cabinet Z100   | pcs  | 20          | 20            | 800,000    | 16,000,000
  --+-----------+-----------------------+------+-------------+---------------+------------+--------
                                            Total Delivered Quantity: 150 pcs
                                            Total Amount: 237,000,000

  Remarks: Chair quantity 20 pcs back-ordered, expected delivery 25/07/2026

  ────────────────────────────────────────────────────────────────
  Deliverer:     ___________      Receiver:      ___________
  Warehouse Keeper                Customer Signature

  Received Date: ___/___/2026
────────────────────────────────────────────────────────────────────
```

---

## 3. Customer Invoice (E-Invoice)

```
────────────────────────────────────────────────────────────────────
                    ELECTRONIC INVOICE
                   (HÓA ĐƠN ĐIỆN TỬ)
                   Decree 123/2020/ND-CP

  Invoice No:   INV-202607-0001
  GDT Code:     7F8A2B3C-XXXX-XXXX-XXXX-XXXXXXXXXXXX
  Issue Date:   18/07/2026
  Invoice Type: Domestic (VAT taxable)

  Seller:
    Company Name:  GoTax Software Co., Ltd
    Tax Code:      0987654321
    Address:       456 Le Loi, District 1, HCMC
    Phone:         028 xxxxxxx
    Bank Account:  1234567890 at Vietcombank

  Buyer:
    Company Name:  ABC Corporation
    Tax Code:      0123456789
    Address:       123 Nguyen Hue, Q1, HCMC
    Phone:         0909xxxxxx

  # | Description              | Unit | Qty | Unit Price    | Disc% | Total before VAT
  --+--------------------------+------+-----+---------------+-------+-----------------
  1 | Office Chair X200        | pcs  | 80  | 1,200,000     | 5%    | 91,200,000
  2 | Desk Y500                | pcs  | 50  | 2,500,000     | 0%    | 125,000,000
  3 | Filing Cabinet Z100      | pcs  | 20  | 800,000       | 10%   | 14,400,000
  --+--------------------------+------+-----+---------------+-------+-----------------

  ────────────────────────────────────────────────────────────────────────
  VAT Rate | Subtotal         | VAT Amount         | Total
  ─────────┼──────────────────┼────────────────────┼────────────────────
  10%      │ 230,600,000      │ 23,060,000         │ 253,660,000
  ─────────┴──────────────────┴────────────────────┴────────────────────
  Grand Total: 253,660,000

  Amount in words: Two hundred fifty-three million six hundred sixty thousand
  Vietnam Dong

  QR Code: [QR Code linking to GDT invoice verification]

  ────────────────────────────────────────────────────────────────
  Digital Signature:
  [Base64 XMLDSig Signature]
  Signed by: GoTax Software Co., Ltd
  Certificate: GDT-issued digital certificate
────────────────────────────────────────────────────────────────────
```

---

## 4. Customer Receipt

```
────────────────────────────────────────────────────────────────────
                      PAYMENT RECEIPT
                    (PHIẾU THU TIỀN)

  Receipt No:    RCP-202607-0001              Date: 25/07/2026
  Customer:      ABC Corporation               Tax Code: 0123456789
  Payment Method: Bank Transfer
  Bank Reference: BTR-20260725-001
  Bank Account:   Vietcombank HCMC

  ────────────────────────────────────────────────────────────────
  Invoice                  Total        Amount Paid      Balance
  ────────────────────────────────────────────────────────────────
  INV-202607-0001          253,660,000  253,660,000       0
  ────────────────────────────────────────────────────────────────
  Total Received: 253,660,000

  Amount in words: Two hundred fifty-three million six hundred sixty thousand
  Vietnam Dong

  Status: ⬜ Draft    ✅ Posted    ⬜ Cancelled    ⬜ Reconciled

  ────────────────────────────────────────────────────────────────
  Received By:     ___________      Printed: ___/___/2026
  AR Accountant
────────────────────────────────────────────────────────────────────
```

---

## 5. Credit Note

```
────────────────────────────────────────────────────────────────────
                    ELECTRONIC CREDIT NOTE
                   (HÓA ĐƠN ĐIỀU CHỈNH / GIẢM)
                   Decree 123/2020/ND-CP

  Credit Note No:   CN-202607-0001
  GDT Code:         CN-7F8A2B3C-XXXX-XXXX-XXXX-XXXXXXXXXXXX
  Issue Date:       20/07/2026
  Return Type:      Partial Return (defective goods)

  ────────────────────────────────────────────────────────────────
  Original Invoice: INV-202607-0001 (18/07/2026)
  Original GDT Code: 7F8A2B3C-XXXX-XXXX-XXXX-XXXXXXXXXXXX
  ────────────────────────────────────────────────────────────────

  Seller:  GoTax Software Co., Ltd   MST: 0987654321
  Buyer:   ABC Corporation           MST: 0123456789

  Reason for adjustment: Goods returned due to manufacturing defect

  # | Description        | Return Qty | Unit Price    | Disc% | Adjusted Amount
  --+--------------------+------------+---------------+-------+-----------------
  1 | Office Chair X200  | 5          | 1,200,000     | 5%    | (5,700,000)
  --+--------------------+------------+---------------+-------+-----------------

  ────────────────────────────────────────────────────────────────────────
  VAT Rate | Subtotal          | VAT Amount         | Total
  ─────────┼───────────────────┼────────────────────┼────────────────────
  10%      │ (5,700,000)       │ (570,000)          │ (6,270,000)
  ─────────┴───────────────────┴────────────────────┴────────────────────
  Credit Amount: (6,270,000)

  Amount in words: Minus six million two hundred seventy thousand
  Vietnam Dong

  QR Code: [QR Code]

  Digital Signature:
  [Base64 XMLDSig Signature]
────────────────────────────────────────────────────────────────────
```

---

## 6. AR Aging Report

```
────────────────────────────────────────────────────────────────────
                    ACCOUNTS RECEIVABLE AGING
                As at 31 July 2026 (VND)

  Customer      | Total AR     | Current   | 1-30 Days | 31-60 Days | 61-90 Days | 90+ Days
  ──────────────+──────────────+───────────+───────────+────────────+────────────+───────────
  ABC Corp      | 253,660,000  | 253,660,000| 0         | 0          | 0          | 0
  DEF Co        | 405,000,000  | 250,000,000| 100,000,00| 55,000,000 | 0          | 0
                |              |            | 0         |            |            |
  GHI Ltd       | 180,500,000  | 0          | 80,500,000| 50,000,000 | 50,000,000 | 0
  JKL Corp      | 350,000,000  | 0          | 0         | 0          | 100,000,00 | 250,000,000
                |              |            |           |            | 0          |
  MNO SA        | 92,000,000   | 92,000,000 | 0         | 0          | 0          | 0
  ──────────────+──────────────+───────────+───────────+────────────+────────────+───────────
  TOTAL         | 1,281,160,000| 595,660,00| 180,500,00| 105,000,00 | 150,000,00 | 250,000,000
                |              | 0         | 0         | 0          | 0          |
  ──────────────+──────────────+───────────+───────────+────────────+────────────+───────────
  % of Total    | 100%         | 46.5%     | 14.1%     | 8.2%       | 11.7%      | 19.5%

  Aging Buckets:
    Current (not due):   595,660,000
    1-30 Days overdue:   180,500,000
    31-60 Days overdue:  105,000,000
    61-90 Days overdue:  150,000,000
    90+ Days overdue:    250,000,000
  ────────────────────────────────────────────────────────────────────────────────
  DSO (Days Sales Outstanding): 42 days

  Prepared By: ___________     Approved By: ___________
  Date: ___/___/2026
────────────────────────────────────────────────────────────────────
```

---

## 7. Customer Statement

```
────────────────────────────────────────────────────────────────────
                  CUSTOMER STATEMENT OF ACCOUNT
                (BẢNG SAO KÊ CÔNG NỢ PHẢI THU)

  Customer:      ABC Corporation
  Tax Code:      0123456789
  Statement Period: July 2026
  Opening Balance:  150,000,000

  Date       | Ref            | Description                      | Debit     | Credit    | Balance
  ───────────+────────────────+──────────────────────────────────+───────────+───────────+────────────
  01/07/2026 |                | Opening Balance                  |           |           | 150,000,000
  05/07/2026 | INV-0624       | Invoice (Jun goods)              | 150,000,00|           | 300,000,000
             |                |                                 | 0         |           |
  10/07/2026 | RCP-0610       | Payment received                 |           | 150,000,00| 150,000,000
             |                |                                 |           | 0         |
  18/07/2026 | INV-0001       | Invoice (July goods)             | 253,660,00|           | 403,660,000
             |                |                                 | 0         |           |
  20/07/2026 | CN-0001        | Credit Note (defective return)   |           | 6,270,000 | 397,390,000
  25/07/2026 | RCP-0001       | Payment received                 |           | 253,660,00| 143,730,000
             |                |                                 |           | 0         |
  ───────────+────────────────+──────────────────────────────────+───────────+───────────+────────────
  Closing Balance: 143,730,000
  Due Date: 17/08/2026 (Net 30 days from INV-0001)

  Outstanding Items:
  INV-202607-0001: 253,660,000 (issued 18/07) — Received 253,660,000 on 25/07
  CN-202607-0001:   (6,270,000)  (issued 20/07)
  ────────────────────────────────────────────────────────────────
  Balance Due: 0 VND

  Prepared By: ___________     Date: ___/___/2026
────────────────────────────────────────────────────────────────────
```

---

## 8. S01-BH: Sales Ledger (per customer per period)

```
────────────────────────────────────────────────────────────────────
                S01-BH — SALES LEDGER
         (SỔ CHI TIẾT BÁN HÀNG)
         Circular 99/2025/TT-BTC

  Period: July 2026
  Customer: ABC Corporation (KH001)

  Date       | Ref            | Description                  | Revenue(5111)| VAT(3331)  | Total
  ───────────+────────────────+──────────────────────────────+──────────────+────────────+────────────
  01/07/2026 | OB             | Opening Balance              | 0            | 0          | 0
  18/07/2026 | INV-0001       | Office furniture sale        | 230,600,000  | 23,060,000 | 253,660,000
  20/07/2026 | CN-0001        | Return - defective chairs    | (5,700,000)  | (570,000)  | (6,270,000)
  ───────────+────────────────+──────────────────────────────+──────────────+────────────+────────────
  Total                        | 224,900,000  | 22,490,000    | 247,390,000
────────────────────────────────────────────────────────────────────
```

---

## 9. S02-BH: Customer Detail Ledger

```
────────────────────────────────────────────────────────────────────
                S02-BH — CUSTOMER DETAIL LEDGER
         (SỔ CHI TIẾT CÔNG NỢ PHẢI THU)
         Circular 99/2025/TT-BTC

  Period: July 2026
  Customer: ABC Corporation (KH001)

  Date       | Ref            | Description                  | Debit(131↑)| Credit(131↓)| Balance
  ───────────+────────────────+──────────────────────────────+────────────+─────────────+───────────
  01/07/2026 |                | Opening Balance              |            |             | 150,000,000
  10/07/2026 | RCP-0610       | Payment received             |            | 150,000,000 | 0
  18/07/2026 | INV-0001       | Invoice issued               | 253,660,000|             | 253,660,000
  20/07/2026 | CN-0001        | Credit note issued           |            | 6,270,000   | 247,390,000
  25/07/2026 | RCP-0001       | Payment received             |            | 253,660,000 | (6,270,000)
  ───────────+────────────────+──────────────────────────────+────────────+─────────────+───────────
  Closing Balance: (6,270,000) — Credit balance (overpayment)
────────────────────────────────────────────────────────────────────
```

---

## 10. E-Invoice TXML XML Structure (Decree 254/2026)

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Invoice xmlns="http://gdt.gov.vn/invoice">
  <InvoiceHeader>
    <InvoiceNumber>INV-202607-0001</InvoiceNumber>
    <InvoiceDate>2026-07-18T14:30:00+07:00</InvoiceDate>
    <InvoiceType>01GTGT</InvoiceType>
    <CurrencyCode>VND</CurrencyCode>
    <ExchangeRate>1</ExchangeRate>
  </InvoiceHeader>

  <Seller>
    <TaxCode>0987654321</TaxCode>
    <Name>GoTax Software Co., Ltd</Name>
    <Address>456 Le Loi, District 1, HCMC</Address>
    <Phone>028xxxxxxxx</Phone>
    <BankAccount>
      <AccountNumber>1234567890</AccountNumber>
      <BankName>Vietcombank HCMC</BankName>
    </BankAccount>
  </Seller>

  <Buyer>
    <TaxCode>0123456789</TaxCode>
    <Name>ABC Corporation</Name>
    <Address>123 Nguyen Hue, Q1, HCMC</Address>
    <Phone>0909xxxxxx</Phone>
    <BankAccount>
      <AccountNumber>9876543210</AccountNumber>
      <BankName>Techcombank</BankName>
    </BankAccount>
  </Buyer>

  <InvoiceLines>
    <Line>
      <LineNumber>1</LineNumber>
      <ItemCode>SP001</ItemCode>
      <ItemName>Office Chair X200</ItemName>
      <Unit>pcs</Unit>
      <Quantity>80</Quantity>
      <UnitPrice>1200000</UnitPrice>
      <DiscountPct>5</DiscountPct>
      <VatRate>10</VatRate>
      <VatType>VAT_10</VatType>
      <LineTotal>91200000</LineTotal>
      <VatAmount>9120000</VatAmount>
    </Line>
    <!-- Additional lines... -->
  </InvoiceLines>

  <Summary>
    <Subtotal>230600000</Subtotal>
    <TotalDiscount>0</TotalDiscount>
    <VatRateSummary>
      <Rate>
        <VatRate>10</VatRate>
        <Subtotal>230600000</Subtotal>
        <VatAmount>23060000</VatAmount>
      </Rate>
    </VatRateSummary>
    <GrandTotal>253660000</GrandTotal>
    <AmountInWords>...</AmountInWords>
  </Summary>

  <DigitalSignature>
    <SignatureType>XMLDSig</SignatureType>
    <SignatureValue>base64-encoded-signature</SignatureValue>
    <CertificateSerial>12345ABCDE</CertificateSerial>
    <SignedDate>2026-07-18T14:30:00+07:00</SignedDate>
  </DigitalSignature>

  <QRCode>QR code data or reference</QRCode>
  <GDTInvoiceCode>assigned-by-GDT-after-submission</GDTInvoiceCode>
</Invoice>
```

---

## 11. Sales Performance Scorecard (P2)

| Metric | Formula | Target | Frequency |
|--------|---------|--------|-----------|
| Revenue | Sum of 5111 credit | Budget ±5% | Monthly |
| Revenue by customer | Group by customer | N/A | Monthly |
| Revenue by product | Group by item | N/A | Monthly |
| Order count | Count of SOs | Trend +15% YoY | Monthly |
| Avg order value | Total revenue / order count | Trend +10% YoY | Monthly |
| Quote-to-order rate | Converted quotes / total quotes | >60% | Monthly |
| DSO | (AR / Revenue) × days | <35 days | Monthly |
| AR over 60 days | Aging > 60 days / total AR | <10% | Weekly |
| Collection rate | Collected / due this period | >95% | Monthly |
| E-invoice compliance | Issued within 24h / total | >99% | Daily |
| Return rate | Credit note amount / revenue | <3% | Monthly |
