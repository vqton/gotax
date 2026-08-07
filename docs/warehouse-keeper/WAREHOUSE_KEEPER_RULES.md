# Business Rules — Thủ Kho (Warehouse Keeper) Module

**Version:** 1.0
**Date:** August 2026

---

## Assignment Rules

| ID | Rule | Validation | Source |
|----|------|-----------|--------|
| AR-001 | One active keeper per warehouse per period | Check overlap on create/update | MISA benchmark |
| AR-002 | Keeper must be an active user | FK validation | System integrity |
| AR-003 | Warehouse must be active | FK validation | System integrity |
| AR-004 | EffectiveFrom <= EffectiveTo | Date validation | Logic |
| AR-005 | Assignment history preserved (no hard delete) | Soft delete | Audit trail |
| AR-006 | Manager role can also be assigned to same warehouse | Role enum check | MISA pattern |

## Recording Rules

| ID | Rule | Validation | Source |
|----|------|-----------|--------|
| RR-001 | Only assigned keeper can record for their warehouse | Permission check | Segregation of duties |
| RR-002 | Cannot record already recorded slip | Status check | Idempotency |
| RR-003 | Cannot record in closed period | Period check | Circular 99 |
| RR-004 | BalanceQty must equal previous balance + receipt - issue | Calculation check | Accounting integrity |
| RR-005 | Cost prices hidden from keeper if config enabled | View filter | Internal control |
| RR-006 | Recording creates immutable audit trail | System | Law on Accounting Art. 6 |
| RR-007 | Bulk recording allowed (select multiple slips) | UI feature | MISA benchmark |

## Un-recording Rules

| ID | Rule | Validation | Source |
|----|------|-----------|--------|
| UR-001 | Reason is mandatory | Field validation | Internal control |
| UR-002 | Cannot un-record in closed period | Period check | Circular 99 |
| UR-003 | Cannot un-record reconciled entries | Status check | Data integrity |
| UR-004 | Only the original recorder or admin can un-record | Permission check | Audit trail |
| UR-005 | Un-recording logged with user, timestamp, reason | System | Audit trail |

## Stock Ledger Rules

| ID | Rule | Validation | Source |
|----|------|-----------|--------|
| SL-001 | Ledger entries are read-only after creation | System | Immutability |
| SL-002 | Entries ordered by date, then creation time | Sort | Chronological |
| SL-003 | Running balance must be non-negative (if allow_negative=false) | Balance check | Business rule |
| SL-004 | Ledger supports filtering by all fields | Query | Usability |
| SL-005 | Print format follows S06-DN template | Template | TT 99 |

## Inventory Count Rules

| ID | Rule | Validation | Source |
|----|------|-----------|--------|
| IC-001 | Minimum 2 participants (keeper + accountant/manager) | Count validation | Law on Accounting Art. 6 |
| IC-002 | Blind count optional (configurable per count) | Config | MISA benchmark |
| IC-003 | Variance > tolerance triggers re-count | Calculation | Industry standard |
| IC-004 | Tolerance: ±1 unit or ±0.5% of value (whichever larger) | Calculation | MISA/FAST benchmark |
| IC-005 | Maximum 3 counts before force close | Counter | Practical limit |
| IC-006 | Count minutes require 3 signatures (keeper, accountant, manager) | Document | Regulatory |

## Reconciliation Rules

| ID | Rule | Validation | Source |
|----|------|-----------|--------|
| RC-001 | Reconciliation compares StockLedgerEntry vs StockBalance | Query join | Cross-check principle |
| RC-002 | Variance > $1M VND highlighted | Display | Materiality |
| RC-003 | Variance > $10M VND requires investigation note | Validation | Materiality |
| RC-004 | Report supports Excel export | Feature | Usability |

## Configuration Rules

| ID | Rule | Validation | Source |
|----|------|-----------|--------|
| CF-001 | Module can be enabled/disabled per company | Config | MISA toggle pattern |
| CF-002 | Cost price visibility configurable per company | Config | Internal control |
| CF-003 | Config changes logged with timestamp | Audit | System |
| CF-004 | Disabling module hides tabs but preserves data | System | Data integrity |
