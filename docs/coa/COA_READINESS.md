# GoTax COA Module - Production Readiness Assessment

## Assessment Date: 2026-07-27 | Assessor: BA Lead + Chief Accountant (20+ yrs)

## VERDICT: ❌ NOT PRODUCTION READY

Current COA module is functional for basic GL operations but lacks enterprise-grade COA management features required for production deployment in Vietnamese enterprises competing with MISA/Fast/BravoERP.

---

## Gap Analysis vs Circular 99/2025/TT-BTC

| Requirement (Circular 99) | Current State | Target State | Gap |
|---|---|---|---|
| 71 Level-1 accounts per Appendix II | 71 level-1 seeded correctly | Must match TT99 exactly | OK |
| Loai 0 (Off-balance sheet) accounts | N/A — TT99 eliminated Loai 0 (001-009); not in Appendix II | N/A | OK |
| TK 215 (Biological assets) | Seeded (incl. 2151/2152/2153) | Required by TT99 Art. 4 | OK |
| TK 332 (Dividends payable) | Seeded | Required by TT99 | OK |
| TK 158 (Bonded warehouse) | Seeded | Required by TT99 | OK |
| TK 171 (Govt bond repurchase) | Seeded | Required by TT99 | OK |
| TK 82112 (Global min tax / Pillar 2) | Seeded | Required by TT99 Art. 5 | OK |
| Official TT99 detail accounts | Migration 000041 adds 22 missing (1362/1363/3362/3363/3531-3534/356-357/41111/41112/4118/6234/6238/21511-215122), renumbers 3385→3386 | Match Appendix II + enterprise detail | OK |
| Diacritic official names | Migration 000041 rewrites all names to official diacritic Vietnamese | Audit/tax-inspection-safe names | OK |
| Enterprise self-modify COA (no MoF approval) | Basic CRUD only | Full COA customization IU | **HIGH** |
| Internal Accounting Policy (IGAP) issuance | NOT implemented | Must document customizations | **HIGH** |
| Financial statement line-item preservation | NOT validated | Cross-checks on modification | **HIGH** |
| Multi-regime support (TT99, TT133, TT58) | TT99 only | Toggle between regimes | **HIGH** |

## Gap Analysis vs Enterprise ERP (MISA/Fast/Bravo)

| Feature | MISA | Fast | Bravo | GoTax | Gap |
|---|---|---|---|---|---|
| COA import/export (CSV/Excel) | Yes | Yes | Yes | **No** | **CRITICAL** |
| Mass account update/deactivate | Yes | Yes | Yes | **No** | **CRITICAL** |
| COA versioning & history | Yes | Partial | Yes | **No** | **CRITICAL** |
| Old-to-new account mapping | Yes | Yes | Yes | **No** | **CRITICAL** |
| Bank account linking | Yes | Yes | Yes | **No** | **HIGH** |
| Budget control per account | Yes | Yes | Yes | **No** | **HIGH** |
| Analysis codes (cost/profit center) | Yes | Yes | Yes | **No** | **HIGH** |
| Account approval workflow | Yes | Yes | Yes | **No** | **HIGH** |
| Auto account code generation | Yes | Yes | Yes | **No** | **MEDIUM** |
| Account freeze/lock | Yes | Yes | Yes | **No** | **HIGH** |
| Balance drill-down inquiry | Yes | Yes | Yes | **No** | **HIGH** |
| IFRS account mapping layer | Yes | No | Yes | **No** | **MEDIUM** |
| Account usage history report | Yes | Yes | Yes | **No** | **MEDIUM** |
| Inter-company account mapping | Yes | Partial | Yes | **No** | **MEDIUM** |
| 5-level deep hierarchy support | Yes | Yes | Yes | Unlimited | OK |
| REST API for COA management | Yes | Partial | Yes | Partial | **MEDIUM** |
| Multi-currency accounts | Yes | Yes | Yes | Basic | **MEDIUM** |
| Audit trail for COA changes | Yes | Yes | Yes | Basic | **MEDIUM** |

## Security & Compliance Gaps

| Requirement | Current | Target | Gap |
|---|---|---|---|
| AUTH module integration | Basic JWT | Full RBAC + VNeID | **CRITICAL** (per AUTH assessment) |
| Account modification audit | Basic | Immutable + who/what/when | **MEDIUM** |
| Sensitive account flagging | None | PII/tax data tagging | **HIGH** |
| COA change approval (4-eyes) | None | Dual-control required | **HIGH** |

## Risk Assessment

| Risk | Severity | Description | Mitigation |
|---|---|---|---|
| Data migration failure | HIGH | No COA import = manual data entry errors | Implement CSV import Phase 1 |
| Circular 99 non-compliance | LOW | Seed aligned to Appendix II (migration 000041); Loai 0 no longer exists in TT99 | Resolved — keep diacritic names |
| Enterprise rejection | HIGH | Cannot match MISA/Fast COA flexibility | Feature gap closure roadmap |
| Audit finding | MEDIUM | No IGAP feature = auditor concern | Add IGAP document generator |
| IFRS conversion blocked | MEDIUM | No IFRS mapping layer | Design for Phase 2 |

## Implementation Priority (COA Module)

| Priority | Feature | Est. Effort | Dependencies |
|---|---|---|---|
| P0 | Add missing Circular 99 accounts (215, 332, 158, 171, 82112) | DONE — seeded in 000001 + 000041 | DB migration |
| P0 | COA CSV/Excel import | 1 week | File parsing lib |
| P0 | Mass account update/deactivate | 0.5 week | Batch operation |
| P0 | Account freeze/lock (independent of period) | 1 week | Schema change |
| P0 | Account balance inquiry with drill-down | 1.5 weeks | Reporting engine |
| P1 | COA versioning & old/new mapping | 2 weeks | Version table |
| P1 | Account approval workflow (4-eyes) | 2 weeks | Workflow engine |
| P1 | Budget control per account | 2 weeks | Budget module |
| P1 | Analysis codes (cost/profit center) | 1 week | Schema extension |
| P1 | IGAP document generation | 1 week | Template engine |
| P1 | COA export (CSV/Excel/PDF) | 0.5 week | Export engine |
| P2 | IFRS account mapping layer | 2 weeks | Mapping table |
| P2 | Multi-regime toggle (TT99/TT133/TT58) | 2 weeks | Regime engine |
| P2 | Bank account integration | 1 week | External API |
| P2 | Auto account code generation | 0.5 week | Code generator |

**Total estimated effort: 16-22 weeks with 2-3 engineers**

## Recommendation

Phase 1 (P0 + P1): 10-14 weeks → Minimum viable enterprise COA capable of PROD
Phase 2 (P2): 6-8 weeks → Full enterprise parity with MISA/Bravo

---

## Verification Log — 2026-08-12 (migrations 000041–000044 applied)

All items below verified against live PostgreSQL after restart (server auto-migrates on boot).

### Source verification (primary vs mirror)

- **Phụ lục II (TT99 official cấp 1 list)** verified against the gazette metadata on
  `congbao.chinhphu.vn` (file `103-2026-TT-BTC_000.rtf` from the MoF) and the mirrored
  HTML appendix at `luatvietnam.vn` (`31f6e4ef1`). Both agree on all 71 cấp 1 accounts.
- Earlier claims that 161–168/358/831 appear anywhere in TT99 are **retracted** — those
  were TT200/other-regime residues; grep of the mirrored appendix confirms zero
  occurrences. 3xx-series detail (356/357/353x) per TT99 Appendix II.4 (funds by type).
- **Điều 11** (verified in mirrored TT99 Art. 11): enterprises may freely open cấp 2+
  detail accounts. Cấp 1 list itself is closed (Appendix II) — **no new cấp 1** may be
  added. 1111/1112, 415-style constructs are only legal as cấp 2+ detail, never cấp 1.

### Seed corrections applied

| Item | Before | After | Where |
|---|---|---|---|
| 7 / 9 (loại headers) | missing | created, 711→7, 911→9 | 000043 |
| 415/417 (non-TT99 cấp 1) | seeded | deleted (unreferenced) / deactivated (referenced) | 000044 |
| 1385 | "Phải thu về cho mượn TSCĐ" (TT200 residue) | "Phải thu về cổ phần hóa" | 000044 |
| 3385 | "BHTN phải nộp" | renumbered 3386, name "Bảo hiểm thất nghiệp" | 000041 |
| 711 | absent | "Thu nhập khác" REVENUE under 7 | 000043 |
| 911 | absent | "Xác định kết quả kinh doanh" under 9 | 000043 |

### Verified DB invariants (live query)

- 220 accounts, all active; all 71 official cấp 1 present; zero orphans (child w/o parent);
  zero parent-prefix violations; every 3+ digit child code starts with its parent code.
- Seed codes 3386/1385 carry correct names, types, parents.

### Service hardening (code, both backends)

- `CreateAccount` verifies parent exists (`GetByCode`) before insert.
- `UpdateAccount` now runs `Account.Validate()` — previously a PUT could bypass
  validation and corrupt hierarchy.
- `DeleteAccount` rejects deletion of accounts with **POSTED** journal usage
  (`GetAccountUsage().EntryCount > 0`) — both PG and memory impls count POSTED only.
  Unposted/deleted-entry usage does not block deletion. Accounts referenced by
  `account_balances`/`account_analysis`/`closing_template_lines` are protected in PG by
  FK (raw error) — memory backend has no FK; journal check is the main guard.
- `Account.Validate()` enforces parent-prefix hierarchy (cấp 2+ code must start with
  parent code; 1-digit loại headers exempt; min 3-digit, max 6-digit codes).

### VAT/declaration cross-check (tax_service)

- `vatRevenueAccounts` = 5111/5112/5113 **only** — 711 removed. GTGT02 (direct method)
  taxes SalesTotal × 5%; revenue excluded 711 means declared sales revenue stays
  consistent with goods/services.
- `vatOutputAccounts` = 33311/33312 — 33313 removed (no such cấp 2 in TT99; 3331 has
  only 33311/33312 per Appendix II).
- `vatExpenseAccounts` — 611 removed (eliminated by TT99); 152/153/156/621/627/641/642
  kept per Appendix I.
