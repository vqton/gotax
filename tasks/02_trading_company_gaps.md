# GoTax — Trading Company Gap Analysis vs MISA SME 2026

Trading company = mua hàng bán hàng, không sản xuất. Không cần: giá thành sản xuất, phân bổ CCDC theo nhà máy, giá xuất kho theo sản xuất.

---

## 🔴 MUST HAVE (thiếu là không chạy được)

### 1. E-Invoice Lifecycle (Hóa đơn điện tử)
| Feature | MISA SME Description |
|---------|---------------------|
| **Invoice book management** | Theo dõi số series, số hóa đơn, trạng thái sử dụng |
| **Invoice adjustment** | Hóa đơn điều chỉnh (sai sót) |
| **Invoice replacement** | Hóa đơn thay thế |
| **Invoice cancellation** | Hủy/void hóa đơn với lý do |
| **Buyer validation** | Kiểm tra MST người mua活性 tại GDT |
| **Missing/damaged invoice report** | Báo cáo hóa đơn mất/hỏa Nghị định 310/2025 |

### 2. Tax Declarations (Tờ khai thuế)
| Feature | MISA SME Description |
|---------|---------------------|
| **HTKK XML export** | Xuất tờ khai theo chuẩn XML HTKK 5.6.0 để nộp thuế |
| **VAT report auto-generation** | Tự động tổng hợp số liệu lên tờ khai GTGT |
| **CIT finalization (03/TNDN)** | Quyết toán TNDN 3 mức thuế suất 15%/17%/20% (Luật 67/2025) |
| **PIT return (05/KK-TNCN)** | Tờ khai TNCNincludes chỉ tiêu 25.1 thu nhập miễn thuế |
| **PIT finalization (05/QT-TNCN)** | Quyết toán TNCN với phụ lục mới nhất |
| **Tax penalty calculator** | Tính tiền phạt chậm nộp, chậm kê khai |

### 3. Data Migration (Cập nhật số dư đầu năm)
| Feature | MISA SME Description |
|---------|---------------------|
| **Prior year balance carry-over** | Tự động cập nhật số dư CCDC, TSCĐ, chi phí trả trước, doanh thu nhận trước từ năm trước |
| **Multi-year data separation** | Tách dữ liệu theo năm tài chính |
| **TT200 → TT99 migration** | Chuyển đổi dữ liệu từ Thông tư 200 sang Thông tư 99 |

### 4. Print Templates (Mẫu in chứng từ)
| Feature | MISA SME Description |
|---------|---------------------|
| **Voucher print in TT99 format** | In phiếu thu, phiếu chi, sổ nhật ký thu tiền theo TT99 |
| **BCTC print** | In báo cáo tài chính theo mẫu quy định mới |
| **Bank payment order templates** | Mẫu ủy nhiệm chi ngân hàng ACB, VietABank |

---

## 🟡 SHOULD HAVE (cần cho hoạt động hiệu quả)

### 5. Bank / E-Banking
| Feature | MISA SME Description |
|---------|---------------------|
| **Real-time bank API** | Kết nối live với BIDV, VCB, ACB qua BANKHUB |
| **Auto-match transactions** | Tự động đối chiếu giao dịch ngân hàng với bút toán kế toán |

### 6. Purchase
| Feature | MISA SME Description |
|---------|---------------------|
| **Supplier activity check** | Kiểm tra MST nhà cung cấp活性/inactive tại GDT |
| **Import invoices from e-invoice** | Tự động nhập liệu từ hóa đơn đầu vào |

### 7. Sales
| Feature | MISA SME Description |
|---------|---------------------|
| **Customer activity check** | Kiểm tra MST người mua tại GDT |
| **Auto-calculate selling price** | Quản lý bảng giá, tính giá bán theo margin |

### 8. Warehouse
| Feature | MISA SME Description |
|---------|---------------------|
| **Lot tracking / expiry dates** | Theo dõi số lô, hạn sử dụng |
| **Auto GRN from purchase** | Tự động tạo phiếu nhập kho từ hóa đơn mua |
| **Auto DN from sales** | Tự động tạo phiếu xuất kho từ hóa đơn bán |
| **Stock level warnings** | Cảnh báo tồn kho tối thiểu/tối đa |

### 9. Payroll
| Feature | MISA SME Description |
|---------|---------------------|
| **Salary allocation by department** | Phân bổ chi phí lương theo phòng ban |
| **Social insurance declaration** | Tạo tờ khai BHXH/BHYT/BHTN tự động |
| **2026 deduction levels** | Cập nhật mức giảm trừ gia cảnh 2026 |

### 10. Tax
| Feature | MISA SME Description |
|---------|---------------------|
| **PIT progressive table** | Biểu thuế lũy tiến theo Luật 109/2025 |
| **CIT exempt income categories** | Danh mục ngành nghề tỷ lệ doanh thu cao nhất + thu nhập miễn thuế |

### 11. Financial Analysis
| Feature | MISA SME Description |
|---------|---------------------|
| **Financial ratio reports** | Báo cáo tỷ số thanh toán, lợi nhuận, đòn bẩy |
| **Budget vs Actual** | Ngân sách so với thực tế theo phòng ban |

---

## 🟢 NICE TO HAVE (tối ưu vận hành)

### 12. Infrastructure
| Feature | MISA SME Description |
|---------|---------------------|
| **Excel import/export** | Nhập xuất hàng loạt từ Excel cho tất cả module |
| **Document attachment** | Đính kèm file PDF/hình ảnh vào chứng từ |
| **Batch operations** | Duyệt/ghi sổ/hủy hàng loạt |
| **Mobile app** | Xem doanh thu, lợi nhuận, công nợ trên điện thoại |
| **Audit trail per transaction** | Lịch sử thay đổi chi tiết từng chứng từ |
| **Multi-currency revaluation** | Đánh giá lại ngoại tệ cuối kỳ tự động |
| **Period-end automation** | Một click: khấu hao, phân bổ, đánh giá lại, tính thuế |

---

## Summary

| Priority | Count | Examples |
|----------|-------|---------|
| **🔴 MUST HAVE** | 16 | E-invoice lifecycle, HTKK XML, data migration, print templates |
| **🟡 SHOULD HAVE** | 16 | Real-time bank, supplier check, auto GRN/DN, salary allocation |
| **🟢 NICE TO HAVE** | 7 | Excel import, mobile app, batch operations |
| **TOTAL** | **39** | |

**Quick wins (implement first):**
1. HTKK XML export — directly usable, high value
2. Invoice book management — mandatory for compliance
3. Prior year data migration — critical for year-end
4. Supplier/customer GDT check — simple API call, high ROI
5. Print templates — users need to print vouchers daily
