# Cost Accounting (Giá thành) Module — Specification

**Document ID:** GT-SPEC-CA-001
**Version:** 1.0.0
**Date:** 2026-08-07
**Status:** DRAFT
**Author:** GoTax Engineering

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Regulatory Compliance](#2-regulatory-compliance)
3. [Costing Methods](#3-costing-methods)
4. [Business Requirements (BRD)](#4-business-requirements)
5. [Domain Model](#5-domain-model)
6. [GL Account Mapping](#6-gl-account-mapping)
7. [Use Cases](#7-use-cases)
8. [Data Flow Diagrams](#8-data-flow-diagrams)
9. [API Endpoints](#9-api-endpoints)
10. [Database Schema](#10-database-schema)
11. [Workflow Diagrams](#11-workflow-diagrams)
12. [User Journeys](#12-user-journeys)
13. [Report Templates](#13-report-templates)
14. [Integration Points](#14-integration-points)
15. [Implementation Priority](#15-implementation-priority)
16. [Risk Analysis](#16-risk-analysis)
17. [Appendix](#17-appendix)

---

## 1. Executive Summary

### 1.1 Module Purpose

The Cost Accounting module (Giá thành sản phẩm/dịch vụ) calculates the actual cost of products, services, and work-in-progress (WIP) by collecting direct and indirect costs from across the ERP system and allocating them to cost objects. It replaces manual spreadsheet-based costing with automated, auditable, regulation-compliant cost calculation.

### 1.2 Target Users

| Role | Usage |
|------|-------|
| Cost Accountant (Kế toán giá thành) | Primary operator. Runs monthly costing, reviews allocations, closes periods |
| CFO / Management Accountant | Reviews cost reports, variance analysis, profitability by product |
| Warehouse Keeper (Thủ kho) | Issues materials to cost objects via warehouse module |
| Payroll Accountant | Ensures labor costs flow into cost pools |
| External Auditor | Verifies cost allocation methodology and WIP valuation |

### 1.3 Business Value

- **Compliance**: Mandatory under Circular 99/2025/TT-BTC for manufacturing and construction enterprises
- **Accuracy**: Eliminates manual allocation errors across cost pools and objects
- **Auditability**: Full journal trail from source transactions to final cost calculation
- **Decision Support**: Real-time cost visibility per product/service/project
- **MISA SME 2026 Parity**: Matches the 6 costing methods supported by Vietnam's leading ERP

### 1.4 PROD Readiness Assessment — NOT READY

| Gap | Impact | Phase |
|-----|--------|-------|
| No CostObject entity | Cannot track costs per product/service/project | P1 |
| No CostPool entity | Cannot collect costs into pools before allocation | P1 |
| No allocation engine | Cannot distribute overhead to cost objects | P1 |
| No costing period management | No period-end closing for cost calculation | P1 |
| No WIP valuation | Cannot value unfinished production | P1 |
| No journal entry templates for costing | Missing accounting entries at period-end | P1 |
| No TK 621-627 cost pool accounts in GL | Cost pool balances not tracked | P1 |
| TK 154 sub-accounts not structured | WIP accounts not properly segmented | P1 |
| No integration hooks | Warehouse material issuance doesn't flow to costing | P2 |
| No payroll labor cost integration | Direct labor not captured automatically | P2 |
| No depreciation allocation | Fixed asset depreciation not allocated to cost pools | P2 |
| No cost variance analysis | Cannot compare actual vs standard costs | P3 |
| No by-product handling | Cannot deduct by-product value from main product cost | P4 |
| No cost reports | No Thẻ tính giá thành or summary reports | P4 |

**Conclusion:** Module requires full build-out. Currently only the Cost Centers hierarchy exists (CRUD). Zero costing logic, zero allocation engine, zero GL integration.

---

## 2. Regulatory Compliance

### 2.1 Thông tư 99/2025/TT-BTC Requirements

Circular 99 (replacing TT200/2014 and TT133/2016) mandates:

1. **Cost object tracking** — Enterprises must maintain detailed cost records per product, service, construction package, or batch
2. **Cost pool aggregation** — Indirect costs must be collected in designated cost pools (TK 621-627) before allocation
3. **Allocation methodology** — Enterprises choose one of 6 approved methods and must apply it consistently
4. **WIP valuation** — Unfinished products must be valued at period-end, carried forward to next period
5. **Cost accounting period** — Must align with fiscal period (monthly minimum)
6. **Audit trail** — Every cost allocation must have a journal entry with supporting documentation

### 2.2 TK 154 — WIP Cost Structure

Account 154 (Chi phí sản xuất, kinh doanh dở dang) tracks costs accumulated in work-in-progress:

| Sub-account | Code | Name | Notes |
|-------------|------|------|-------|
| 1541 | 1541 | Chi phí xây lắp | Construction/installation costs |
| 1542 | 1542 | Chi phí sản phẩm khác | Other product costs |
| 1543 | 1543 | Chi phí dịch vụ | Service costs |
| 1544 | 1544 | Chi phí bảo hành | Warranty costs |

Each 154x sub-account uses supplementary accounting by cost element:
- 154x.1 — Direct materials (nguyên vật liệu trực tiếp)
- 154x.2 — Direct labor (nhân công trực tiếp)
- 154x.3 — Direct costs (chi phí trực tiếp khác)
- 154x.4 — Manufacturing overhead (chi phí sản xuất chung)

### 2.3 TK 621-627 — Cost Pool Accounts

| Account | Code | Name | Purpose |
|---------|------|------|---------|
| 621 | 621 | Chi phí nguyên vật liệu trực tiếp | Direct materials cost pool |
| 622 | 622 | Chi phí nhân công trực tiếp | Direct labor cost pool |
| 623 | 623 | Chi phí sử dụng máy thi công | Construction machinery costs |
| 624 | 624 | Chi phí dụng cụ công tác | Tools and equipment costs |
| 625 | 625 | Chi phí khấu hao TSCĐ | Fixed asset depreciation pool |
| 626 | 626 | Chi phí BHLD & BHXH nhân công trực tiếp | Social insurance on direct labor |
| 627 | 627 | Chi phí sản xuất chung | Manufacturing overhead pool |

### 2.4 WIP Evaluation Methods per Regulation

Per Circular 99, Article 28:

1. **Actual cost method** — WIP valued at accumulated actual costs in TK 154
2. **Standard cost method** — WIP valued at standard cost, variance recognized separately
3. **Net realizable value** — Where NRV < cost, write down to NRV (prudence principle)

### 2.5 Industry-Specific Requirements

| Industry | Method(s) | Special Rules |
|----------|-----------|---------------|
| Manufacturing (sản xuất) | Coefficient, Ratio, Standard | Must distinguish product lines; batch tracking for 1542 |
| Construction (xây lắp) | Simple, Ratio | Revenue-proportional allocation for long-term contracts; TK 1541 |
| Services (dịch vụ) | Simple, Standard | TK 1543; labor-heavy cost structures |
| Trading (buôn bán) | N/A | Use purchase cost + direct expenses; no manufacturing costing |

---

## 3. Costing Methods

### 3.1 Phương pháp giản đơn (Direct/Simple Method)

**Description:** Costs are assigned directly to cost objects without an allocation step. Only direct costs (materials, direct labor, direct expenses) are used. No overhead allocation.

**When to use:**
- Small enterprises with a single product or service
- Construction projects with clear cost tracing
- Job-order costing where each job is distinct

**Formula:**
```
Unit Cost = (Direct Materials + Direct Labor + Direct Expenses) / Units Produced
```

**Accounting entries (period-end):**

Debit materials from warehouse:
```
Dr  621 — Chi phí NVL trực tiếp
    Cr  155 — Hàng tồn kho thành phẩm (or 154 WIP)
```

Transfer completed goods:
```
Dr  155 — Hàng tồn kho thành phẩm
    Cr  154 — Chi phí SXKD dở dang
```

**Example:**
- Product A: Direct materials = 50,000,000 VND, Direct labor = 20,000,000 VND, Direct expenses = 5,000,000 VND
- Units produced = 100
- Unit cost = (50M + 20M + 5M) / 100 = 750,000 VND/unit
- COGS = 750,000 × sold quantity

### 3.2 Phương pháp hệ số (Coefficient Method)

**Description:** Total costs collected in cost pools are allocated to cost objects using a coefficient (tỷ lệ phân bổ) based on a base (cơ sở phân bổ) such as direct material value, direct labor cost, or machine hours.

**When to use:**
- Manufacturing enterprises with multiple products sharing common resources
- Where direct costs are a good proxy for resource consumption
- Most common method in Vietnamese manufacturing

**Formula:**
```
Allocation Coefficient = Total Indirect Cost Pool / Total Base Amount
Allocated Cost = Coefficient × Base Amount for Each Cost Object
Unit Cost = Total Cost / Units Produced
```

**Base options:**
- Direct material value (giá trị NVL trực tiếp)
- Direct labor cost (chi phí nhân công trực tiếp)
- Direct cost total (tổng chi phí trực tiếp)
- Machine hours (số giờ máy)

**Accounting entries:**
```
Dr  154x.4 — Chi phí SX chung (per cost object)
    Cr  627 — Chi phí sản xuất chung
```

**Example:**
- Total overhead in TK 627 = 300,000,000 VND
- Total direct material (base) = 500,000,000 VND
- Coefficient = 300M / 500M = 0.6

| Product | Direct Material | Overhead Allocated | Units | Unit Overhead |
|---------|-----------------|-------------------|-------|---------------|
| A | 200,000,000 | 120,000,000 | 200 | 600,000 |
| B | 150,000,000 | 90,000,000 | 150 | 600,000 |
| C | 150,000,000 | 90,000,000 | 100 | 900,000 |

### 3.3 Phương pháp tỷ lệ (Proportion/Ratio Method)

**Description:** Costs are allocated based on the proportion of each cost object's direct costs to total direct costs across all objects. Similar to coefficient but expressed as ratios.

**When to use:**
- When products have comparable cost structures
- Construction projects allocating shared site costs
- When management prefers ratio-based allocation

**Formula:**
```
Allocation Ratio = Object Direct Costs / Total Direct Costs (all objects)
Allocated Overhead = Ratio × Total Overhead Pool
```

**Accounting entries:**
```
Dr  154x.4 — Chi phí SX chung
    Cr  627 — Chi phí sản xuất chung
```

**Example:**
- Total overhead = 240,000,000 VND
- Total direct costs (all products) = 1,200,000,000 VND

| Product | Direct Costs | Ratio | Overhead Share |
|---------|-------------|-------|----------------|
| X | 600,000,000 | 50% | 120,000,000 |
| Y | 360,000,000 | 30% | 72,000,000 |
| Z | 240,000,000 | 20% | 48,000,000 |

### 3.4 Phương pháp định mức (Standard/Norm Method)

**Description:** Uses predetermined standard costs per unit. Variances between actual and standard are tracked and analyzed separately. Most complex but most useful for variance analysis and control.

**When to use:**
- Large manufacturing enterprises with stable production processes
- When management needs variance reports for cost control
- Continuous production with predictable resource consumption

**Formula:**
```
Standard Cost per Unit = Standard Material + Standard Labor + Standard Overhead
Variance = Actual Cost - Standard Cost
Favorable Variance = Standard > Actual (cost saving)
Unfavorable Variance = Actual > Standard (cost overrun)
```

**Accounting entries — recording standard cost:**
```
Dr  155 — Hàng tồn kho thành phẩm (at standard)
    Cr  154 — Chi phí SXKD dở dang (at standard)
```

**Variance recognition:**
```
Dr  627 — Chi phí sản xuất chung (if unfavorable)
    Cr  154 — Chi phí SXKD dở dang (if favorable)
```

**Example:**
- Standard material = 500,000 VND/unit, Standard labor = 200,000 VND/unit, Standard overhead = 150,000 VND/unit
- Standard cost = 850,000 VND/unit
- Actual cost = 900,000 VND/unit
- Unfavorable variance = 50,000 VND/unit × 1,000 units = 50,000,000 VND

### 3.5 Phương pháp phân bước liên tục (Continuous Process Method)

**Description:** For industries with continuous flow production (oil refining, chemicals, food processing). Costs are accumulated by process step (bước) and assigned at each transfer point. Uses equivalent units of production.

**When to use:**
- Continuous process industries (petrochemical, food & beverage, cement)
- No distinct batches — production flows continuously
- Multiple processing departments

**Formula:**
```
Equivalent Units = Completed Units + (WIP Units × % Completion)
Cost per Equivalent Unit = Total Cost / Equivalent Units
```

**Accounting entries — cost accumulation per process:**
```
Dr  1542.3 — Chi phí bước 1
    Cr  621 — Chi phí NVL trực tiếp
    Cr  622 — Chi phí nhân công trực tiếp
    Cr  627 — Chi phí sản xuất chung
```

**Transfer between steps:**
```
Dr  1542.4 — Chi phí bước 2
    Cr  1542.3 — Chi phí bước 1
```

**Example:**
- Step 1: 500 units started, 400 completed, 100 at 60% completion
- Equivalent units = 400 + (100 × 0.6) = 460
- Step 1 cost = 230,000,000 VND
- Cost per equivalent unit = 230M / 460 = 500,000 VND

### 3.6 Phương pháp loại trừ sản phẩm phụ (By-product Exclusion Method)

**Description:** Deducts the estimated value of by-products from the total production cost before allocating to the main product. By-products receive no separate cost allocation.

**When to use:**
- Manufacturing with significant by-products (e.g., sawmills producing lumber + sawdust, oil refining producing gasoline + asphalt)
- When by-product value is relatively small compared to main product

**Formula:**
```
Main Product Cost = Total Production Cost - By-product Net Realizable Value
Unit Cost = Main Product Cost / Main Product Units
```

**Accounting entries — by-product recovery:**
```
Dr  155 — Hàng tồn kho thành phẩm (by-product)
    Cr  154 — Chi phí SXKD dở dang (deduction from main product cost)
```

**Example:**
- Total production cost = 1,000,000,000 VND
- By-product (sawdust) NRV = 50,000,000 VND
- Main product (lumber) cost = 950,000,000 VND
- Main product units = 1,000
- Unit cost = 950,000 VND/unit

---

## 4. Business Requirements

### 4.1 Functional Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-CA-001 | System shall support creating Cost Objects with code, name, type (product/service/project/batch) | P1 |
| FR-CA-002 | System shall support Cost Object hierarchy (parent-child) | P1 |
| FR-CA-003 | System shall create Cost Pools mapped to TK 621-627 accounts | P1 |
| FR-CA-004 | System shall collect direct material costs from warehouse issuance documents | P1 |
| FR-CA-005 | System shall collect direct labor costs from payroll module | P1 |
| FR-CA-006 | System shall collect manufacturing overhead from GL journal entries | P1 |
| FR-CA-007 | System shall allocate overhead using Coefficient method (default) | P1 |
| FR-CA-008 | System shall allocate overhead using Simple method | P1 |
| FR-CA-009 | System shall allocate overhead using Ratio method | P2 |
| FR-CA-010 | System shall allocate overhead using Standard/Norm method | P2 |
| FR-CA-011 | System shall support Continuous Process costing | P3 |
| FR-CA-012 | System shall support By-product exclusion costing | P3 |
| FR-CA-013 | System shall generate journal entries for cost allocation (Dr 154x, Cr 62x) | P1 |
| FR-CA-014 | System shall value WIP at period-end | P1 |
| FR-CA-015 | System shall transfer completed product costs from WIP (154) to finished goods (155) | P1 |
| FR-CA-016 | System shall calculate COGS (632) when goods are sold | P1 |
| FR-CA-017 | System shall support costing periods aligned with fiscal periods | P1 |
| FR-CA-018 | System shall support closing a costing period (lock against further changes) | P1 |
| FR-CA-019 | System shall prevent costing calculations on closed periods | P1 |
| FR-CA-020 | System shall generate Thẻ tính giá thành (Cost Calculation Sheet) per Circular 99 | P2 |
| FR-CA-021 | System shall generate cost summary report by cost object | P2 |
| FR-CA-022 | System shall generate WIP valuation report | P2 |
| FR-CA-023 | System shall track cost variances (actual vs standard) | P3 |
| FR-CA-024 | System shall support multiple costing methods per company (one active per period) | P1 |
| FR-CA-025 | System shall validate all allocation bases have positive values before calculation | P1 |
| FR-CA-026 | System shall support recalculation (reverse + re-run) for a costing period | P1 |
| FR-CA-027 | System shall maintain audit log of all costing calculations and journal entries | P1 |
| FR-CA-028 | System shall support cost allocation by multiple bases (material value, labor cost, custom) | P2 |
| FR-CA-029 | System shall support period-end depreciation allocation to cost pools | P2 |
| FR-CA-030 | System shall support warranty cost tracking (TK 1544) | P3 |
| FR-CA-031 | System shall support construction project costing (TK 1541) | P2 |
| FR-CA-032 | System shall support service costing (TK 1543) | P2 |

### 4.2 Non-Functional Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| NFR-CA-001 | Costing calculation for 10,000 cost objects shall complete in < 30 seconds | P1 |
| NFR-CA-002 | Period-end closing shall be reversible (audit-friendly) | P1 |
| NFR-CA-003 | All monetary values stored as DECIMAL(19,4) — no floating point | P1 |
| NFR-CA-004 | Concurrent costing runs for different companies shall not interfere | P1 |
| NFR-CA-005 | Costing results immutable once period is closed | P1 |
| NFR-CA-006 | Full rollback support — if costing fails, no partial state persists | P1 |

### 4.3 Integration Requirements

| ID | Module | Integration Point |
|----|--------|-------------------|
| IR-CA-001 | GL | Journal entries for cost allocation and WIP transfer |
| IR-CA-002 | Warehouse | Material issuance documents → direct material cost collection |
| IR-CA-003 | Payroll | Direct labor costs → labor cost pool |
| IR-CA-004 | Fixed Assets | Depreciation → overhead cost pool (TK 625) |
| IR-CA-005 | Purchase | Purchase price variance → standard cost variance |
| IR-CA-006 | Sale | Sales quantity → units sold for COGS calculation |
| IR-CA-007 | Budget | Budgeted costs → standard cost baseline |
| IR-CA-008 | Cost Centers | Cost center hierarchy → cost object allocation |

---

## 5. Domain Model

### 5.1 New Entities

#### CostObject — The unit of cost tracking (product, service, project, batch)

```go
type CostObject struct {
    ID               string            `json:"id"`
    CompanyID        string            `json:"company_id"`
    Code             string            `json:"code"`
    Name             string            `json:"name"`
    Type             CostObjectType    `json:"type"`              // product, service, project, batch
    CostCenterID     string            `json:"cost_center_id,omitempty"`
    ParentID         string            `json:"parent_id,omitempty"`
    GLAccountCode    string            `json:"gl_account_code"`   // TK 1541/1542/1543/1544
    UnitOfMeasure    string            `json:"unit_of_measure"`
    StandardMaterial float64           `json:"standard_material,omitempty"`
    StandardLabor    float64           `json:"standard_labor,omitempty"`
    StandardOverhead float64           `json:"standard_overhead,omitempty"`
    IsActive         bool              `json:"is_active"`
    CreatedAt        string            `json:"created_at"`
    UpdatedAt        string            `json:"updated_at"`
}

type CostObjectType string

const (
    CostObjectTypeProduct   CostObjectType = "product"
    CostObjectTypeService   CostObjectType = "service"
    CostObjectTypeProject   CostObjectType = "project"
    CostObjectTypeBatch     CostObjectType = "batch"
    CostObjectTypeByProduct CostObjectType = "by_product"
)
```

#### CostPool — Collects costs from a specific category before allocation

```go
type CostPool struct {
    ID              string       `json:"id"`
    CompanyID       string       `json:"company_id"`
    Code            string       `json:"code"`
    Name            string       `json:"name"`
    GLAccountCode   string       `json:"gl_account_code"` // 621, 622, 623, 624, 625, 626, 627
    PoolType        CostPoolType `json:"pool_type"`
    AllocationBase  string       `json:"allocation_base,omitempty"` // material_value, labor_cost, machine_hours, custom
    IsActive        bool         `json:"is_active"`
    CreatedAt       string       `json:"created_at"`
    UpdatedAt       string       `json:"updated_at"`
}

type CostPoolType string

const (
    CostPoolDirectMaterial  CostPoolType = "direct_material"  // 621
    CostPoolDirectLabor     CostPoolType = "direct_labor"     // 622
    CostPoolMachinery       CostPoolType = "machinery"        // 623
    CostPoolTools           CostPoolType = "tools"            // 624
    CostPoolDepreciation    CostPoolType = "depreciation"     // 625
    CostPoolSocialInsurance CostPoolType = "social_insurance" // 626
    CostPoolOverhead        CostPoolType = "overhead"         // 627
)
```

#### CostAllocation — Records a single allocation run result

```go
type CostAllocation struct {
    ID              string  `json:"id"`
    CompanyID       string  `json:"company_id"`
    CostingPeriodID string  `json:"costing_period_id"`
    CostPoolID      string  `json:"cost_pool_id"`
    CostObjectID    string  `json:"cost_object_id"`
    AllocatedAmount float64 `json:"allocated_amount"`
    AllocationRatio float64 `json:"allocation_ratio"`
    AllocationBase  string  `json:"allocation_base"`
    BaseAmount      float64 `json:"base_amount"`
    CreatedAt       string  `json:"created_at"`
}
```

#### CostingPeriod — Manages the costing lifecycle per fiscal period

```go
type CostingPeriod struct {
    ID              string              `json:"id"`
    CompanyID       string              `json:"company_id"`
    PeriodNumber    int                 `json:"period_number"`        // 1-13
    Year            int                 `json:"year"`
    FiscalYearID    string              `json:"fiscal_year_id"`
    Status          CostingPeriodStatus `json:"status"`
    CostingMethod   CostingMethod       `json:"costing_method"`
    AllocatedTotal  float64             `json:"allocated_total"`
    ClosedAt        string              `json:"closed_at,omitempty"`
    ClosedBy        string              `json:"closed_by,omitempty"`
    CreatedAt       string              `json:"created_at"`
    UpdatedAt       string              `json:"updated_at"`
}

type CostingPeriodStatus string

const (
    CostingPeriodStatusOpen    CostingPeriodStatus = "open"
    CostingPeriodStatusRunning CostingPeriodStatus = "running"
    CostingPeriodStatusClosed  CostingPeriodStatus = "closed"
)

type CostingMethod string

const (
    CostingMethodSimple      CostingMethod = "simple"
    CostingMethodCoefficient CostingMethod = "coefficient"
    CostingMethodRatio       CostingMethod = "ratio"
    CostingMethodStandard    CostingMethod = "standard"
    CostingMethodProcess     CostingMethod = "process"
    CostingMethodByProduct   CostingMethod = "by_product"
)
```

#### CostingResult — Final calculated cost per cost object per period

```go
type CostingResult struct {
    ID                string  `json:"id"`
    CompanyID         string  `json:"company_id"`
    CostingPeriodID   string  `json:"costing_period_id"`
    CostObjectID      string  `json:"cost_object_id"`
    DirectMaterials   float64 `json:"direct_materials"`
    DirectLabor       float64 `json:"direct_labor"`
    DirectExpenses    float64 `json:"direct_expenses"`
    AllocatedOverhead float64 `json:"allocated_overhead"`
    TotalCost         float64 `json:"total_cost"`
    UnitsProduced     float64 `json:"units_produced"`
    UnitCost          float64 `json:"unit_cost"`
    WIPValue          float64 `json:"wip_value"`
    CompletedValue    float64 `json:"completed_value"`
    COGS              float64 `json:"cogs"`
    Variance          float64 `json:"variance,omitempty"` // actual - standard
    CreatedAt         string  `json:"created_at"`
}
```

#### CostAllocationBasis — Defines the base for overhead allocation

```go
type CostAllocationBasis struct {
    ID            string `json:"id"`
    CompanyID     string `json:"company_id"`
    Name          string `json:"name"`
    BasisType     string `json:"basis_type"`     // material_value, labor_cost, machine_hours, direct_cost_total, custom
    GLAccountCode string `json:"gl_account_code,omitempty"` // for custom: which account to sum
    IsActive      bool   `json:"is_active"`
    CreatedAt     string `json:"created_at"`
}
```

### 5.2 Extended Entities

#### CostCenter — Add `cost_type` field

```go
// Addition to existing CostCenter struct
type CostCenter struct {
    // ... existing fields ...
    CostType  string `json:"cost_type,omitempty"` // production, service, construction, admin
}
```

### 5.3 Relationships

```
CostCenter  1──N  CostObject
CostPool    1──N  CostAllocation
CostObject  1──N  CostAllocation
CostingPeriod  1──N  CostAllocation
CostingPeriod  1──N  CostingResult
CostObject  1──1  CostingResult (per period)
```

---

## 6. GL Account Mapping

### 6.1 Account to Cost Pool Mapping

| GL Account | Code | Name | Pool Type | Allocation Base |
|------------|------|------|-----------|-----------------|
| 621 | 621 | Chi phí NVL trực tiếp | Direct Material | Not allocated — traced directly |
| 622 | 622 | Chi phí nhân công trực tiếp | Direct Labor | Not allocated — traced directly |
| 623 | 623 | Chi phí sử dụng máy thi công | Machinery | Machine hours |
| 624 | 624 | Chi phí dụng cụ công tác | Tools | Direct cost proportion |
| 625 | 625 | Chi phí khấu hao TSCĐ | Depreciation | Machine hours or floor area |
| 626 | 626 | Chi phí BHLD & BHXH | Social Insurance | Direct labor proportion |
| 627 | 627 | Chi phí sản xuất chung | Overhead | Material value or labor cost |

### 6.2 Journal Entry Templates

#### 6.2.1 Collecting Direct Material Costs (from warehouse issuance)

```
Date: Period-end
Description: Phân bổ chi phí NVL trực tiếp cho đối tượng {code}

Dr  621          {amount}    Chi phí NVL trực tiếp
    Cr  154x.1   {amount}    Chi phí NVL — {cost_object_name}

Supporting: Warehouse issuance document list by cost object
```

#### 6.2.2 Collecting Direct Labor Costs (from payroll)

```
Date: Period-end
Description: Phân bổ chi phí nhân công trực tiếp cho đối tượng {code}

Dr  622          {amount}    Chi phí nhân công trực tiếp
    Cr  154x.2   {amount}    Chi phí nhân công — {cost_object_name}

Supporting: Payroll summary by cost object
```

#### 6.2.3 Allocating Manufacturing Overhead (Coefficient method)

```
Date: Period-end
Description: Phân bổ chi phí sản xuất chung cho đối tượng {code}

Dr  154x.4      {amount}    Chi phí SX chung — {cost_object_name}
    Cr  627      {amount}    Chi phí sản xuất chung

Allocation base: {base_name} = {base_amount}, Coefficient = {coefficient}
```

#### 6.2.4 Allocating Depreciation to Cost Pool

```
Date: Period-end
Description: Phân bổ khấu hao TSCĐ vào chi phí sản xuất

Dr  625          {amount}    Chi phí khấu hao TSCĐ
    Cr  154      {amount}    Chi phí SXKD dở dang
```

#### 6.2.5 Transfer Completed Products from WIP to Finished Goods

```
Date: Period-end
Description: Giá thành sản phẩm hoàn thành — {cost_object_name}

Dr  155          {amount}    Hàng tồn kho thành phẩm
    Cr  154x     {amount}    Chi phí SXKD dở dang — {cost_object_name}

Amount = Direct Materials + Direct Labor + Direct Expenses + Allocated Overhead
```

#### 6.2.6 COGS Recognition (when goods sold)

```
Date: Sale date
Description: Giá vốn hàng bán — {product_name}

Dr  632          {amount}    Giá vốn hàng bán
    Cr  155      {amount}    Hàng tồn kho thành phẩm
```

#### 6.2.7 WIP Valuation (unfinished work carried forward)

```
Date: Period-end
Description: Đánh giá dở dang cuối kỳ — {cost_object_name}

Dr  154x         {wip_amount}   Chi phí SXKD dở dang (WIP)
    Cr  900      {wip_amount}   Thu nhập chưa thực hiện

Note: Reverse in next period:
Dr  900          {wip_amount}   Thu nhập chưa thực hiện
    Cr  154x     {wip_amount}   Chi phí SXKD dở dang (WIP)
```

---

## 7. Use Cases

### UC-CA-001: Create Cost Object

| Field | Value |
|-------|-------|
| **UC编号** | UC-CA-001 |
| **Name** | Create Cost Object |
| **Actor** | Cost Accountant |
| **Precondition** | User logged in, has cost accounting permission, Cost Center exists |
| **Happy Path** | 1. Navigate to Cost Accounting > Cost Objects > Create |
| | 2. Enter code, name, type (product/service/project/batch) |
| | 3. Select parent cost object (optional) |
| | 4. Select cost center (optional) |
| | 5. Enter standard costs (optional, for standard method) |
| | 6. Select GL account (1541/1542/1543/1544) |
| | 7. Enter unit of measure |
| | 8. Submit |
| **Alternative Paths** | 3a. No parent — leave empty for top-level object |
| | 6a. System auto-suggests GL account based on type |
| **Exception Paths** | Code already exists → 409 Conflict |
| | Missing required fields → 400 Bad Request |
| **Postcondition** | Cost object created, visible in cost object list |

### UC-CA-002: Create Cost Pool

| Field | Value |
|-------|-------|
| **UC编号** | UC-CA-002 |
| **Name** | Create Cost Pool |
| **Actor** | Cost Accountant |
| **Precondition** | User logged in, GL accounts 621-627 exist |
| **Happy Path** | 1. Navigate to Cost Accounting > Cost Pools > Create |
| | 2. Enter code, name |
| | 3. Select GL account (621-627) |
| | 4. Select pool type |
| | 5. Select allocation base (for overhead pools) |
| | 6. Submit |
| **Alternative Paths** | 4a. Pool type = direct_material → no allocation base needed (traced directly) |
| **Exception Paths** | GL account already mapped → 409 Conflict |
| **Postcondition** | Cost pool created, ready to collect costs |

### UC-CA-003: Collect Direct Material Costs

| Field | Value |
|-------|-------|
| **UC编号** | UC-CA-003 |
| **Name** | Collect Direct Material Costs |
| **Actor** | System (automated) / Cost Accountant (trigger) |
| **Precondition** | Costing period open, warehouse issuance documents exist for period, cost objects defined |
| **Happy Path** | 1. Open Costing Period |
| | 2. System queries warehouse issuance documents for period |
| | 3. System groups issuances by cost object |
| | 4. System creates cost allocations: pool = 621, object = cost_object, amount = sum of issuances |
| | 5. Display summary: total materials collected per cost object |
| **Alternative Paths** | 3a. Some issuances not assigned to cost object → flag as "unallocated" |
| **Exception Paths** | No issuance documents → warning, proceed with zero |
| **Postcondition** | Direct material costs collected into 621 pool, ready for allocation |

### UC-CA-004: Collect Direct Labor Costs

| Field | Value |
|-------|-------|
| **UC编号** | UC-CA-004 |
| **Name** | Collect Direct Labor Costs |
| **Actor** | System (automated) / Cost Accountant (trigger) |
| **Precondition** | Costing period open, payroll processed for period, cost objects defined |
| **Happy Path** | 1. System queries payroll data for period |
| | 2. System filters direct labor entries (labor cost type = "direct") |
| | 3. System groups by cost object |
| | 4. Creates cost allocations: pool = 622, object = cost_object, amount = sum of labor |
| | 5. Display summary |
| **Alternative Paths** | 2a. Some labor not assigned → flag as "unallocated direct labor" |
| **Exception Paths** | Payroll not processed → error "Process payroll first" |
| **Postcondition** | Direct labor costs collected into 622 pool |

### UC-CA-005: Collect Manufacturing Overhead

| Field | Value |
|-------|-------|
| **UC编号** | UC-CA-005 |
| **Name** | Collect Manufacturing Overhead |
| **Actor** | System (automated) / Cost Accountant (trigger) |
| **Precondition** | Costing period open, depreciation run completed, utilities/rent journals posted |
| **Happy Path** | 1. System queries GL journal entries for period |
| | 2. Filters entries posted to accounts 623-627 |
| | 3. Groups by cost pool |
| | 4. Displays total per cost pool: 623, 624, 625, 626, 627 |
| | 5. Cost accountant confirms collection |
| **Alternative Paths** | 4a. Cost accountant can manually add/adjust overhead amounts |
| **Exception Paths** | No overhead entries → warning, proceed with zero |
| **Postcondition** | Overhead costs collected into pools, ready for allocation |

### UC-CA-006: Allocate Overhead to Cost Objects (Coefficient Method)

| Field | Value |
|-------|-------|
| **UC编号** | UC-CA-006 |
| **Name** | Allocate Overhead Using Coefficient Method |
| **Actor** | Cost Accountant |
| **Precondition** | Direct costs collected, cost pools have balances, allocation base selected |
| **Happy Path** | 1. Select costing period |
| | 2. Select overhead pool(s) to allocate |
| | 3. Select allocation base (e.g., direct material value) |
| | 4. System calculates: coefficient = total pool / total base |
| | 5. System allocates: each object share = coefficient × object's base amount |
| | 6. Display allocation summary per cost object |
| | 7. Cost accountant confirms |
| | 8. System creates journal entries (Dr 154x.4, Cr 627) |
| **Alternative Paths** | 4a. Coefficient rounds to 6 decimal places |
| | 5a. Rounding differences adjusted to largest cost object |
| **Exception Paths** | Total base = 0 → error "No allocation base available" |
| | Pool balance = 0 → skip, no allocation needed |
| **Postcondition** | Overhead allocated, journal entries posted, cost pools zeroed |

### UC-CA-007: Calculate Costing

| Field | Value |
|-------|-------|
| **UC编号** | UC-CA-007 |
| **Name** | Calculate Total Cost per Cost Object |
| **Actor** | Cost Accountant |
| **Precondition** | All cost pools allocated, period open |
| **Happy Path** | 1. Select costing period |
| | 2. System sums: direct_materials + direct_labor + direct_expenses + allocated_overhead per cost object |
| | 3. System calculates: unit_cost = total_cost / units_produced |
| | 4. System identifies WIP vs completed quantities |
| | 5. System calculates: completed_value, wip_value |
| | 6. Display CostingResult per cost object |
| | 7. Cost accountant confirms |
| | 8. System creates CostingResult records |
| **Alternative Paths** | 5a. If no WIP at period-end → wip_value = 0 |
| **Exception Paths** | Units produced = 0 → error "No units to cost" |
| | Cost pools not fully allocated → error "Complete allocation first" |
| **Postcondition** | Costing results saved, ready for WIP evaluation and period closing |

### UC-CA-008: Evaluate Work-in-Progress (WIP)

| Field | Value |
|-------|-------|
| **UC编号** | UC-CA-008 |
| **Name** | WIP Valuation at Period-End |
| **Actor** | Cost Accountant |
| **Precondition** | Costing calculated, WIP quantities entered |
| **Happy Path** | 1. For each cost object with WIP: enter completion percentage |
| | 2. System calculates: wip_value = total_cost × completion% × wip_units |
| | 3. System generates WIP transfer journal (Dr 154x, Cr 900) |
| | 4. Display WIP valuation summary |
| | 5. Cost accountant confirms |
| **Alternative Paths** | 2a. Standard method: use standard cost for WIP |
| | 2b. Process method: use equivalent units |
| **Exception Paths** | Completion % not entered → error "Enter completion percentage" |
| | Completion % > 100% → error |
| **Postcondition** | WIP valued and journal entries posted |

### UC-CA-009: Close Costing Period

| Field | Value |
|-------|-------|
| **UC编号** | UC-CA-009 |
| **Name** | Close Costing Period |
| **Actor** | Cost Accountant |
| **Precondition** | All calculations complete, all allocations posted, WIP valued |
| **Happy Path** | 1. Review period summary (total costs, allocations, results) |
| | 2. Verify all journal entries posted |
| | 3. Verify no pending allocations |
| | 4. Click "Close Period" |
| | 5. System sets period status = "closed" |
| | 6. System records closed_at timestamp and closed_by user |
| | 7. Period locked — no further changes allowed |
| **Alternative Paths** | 4a. Can reopen if needed (reverse close, requires permission) |
| **Exception Paths** | Pending allocations → error "Complete all allocations first" |
| | Missing WIP valuation → error "Value WIP first" |
| **Postcondition** | Period closed, costing results locked |

### UC-CA-010: Recalculate Costing (Reverse + Re-run)

| Field | Value |
|-------|-------|
| **UC编号** | UC-CA-010 |
| **Name** | Recalculate Costing for Period |
| **Actor** | Cost Accountant |
| **Precondition** | Period is open, previous calculation exists |
| **Happy Path** | 1. Select open period with existing calculations |
| | 2. Click "Recalculate" |
| | 3. System reverses all allocation journal entries for period |
| | 4. System deletes existing CostAllocation and CostingResult records |
| | 5. System re-runs collection and allocation |
| | 6. Display updated results |
| **Alternative Paths** | 2a. Partial recalculate: select specific pool only |
| **Exception Paths** | Period closed → error "Cannot recalculate closed period" |
| | Reversal entries fail → transaction rollback |
| **Postcondition** | Costing recalculated with updated data |

### UC-CA-011: Generate Cost Calculation Sheet (Thẻ tính giá thành)

| Field | Value |
|-------|-------|
| **UC编号** | UC-CA-011 |
| **Name** | Generate Cost Calculation Sheet |
| **Actor** | Cost Accountant |
| **Precondition** | Costing calculated for period |
| **Happy Path** | 1. Select costing period and cost object(s) |
| | 2. Click "Generate Report" |
| | 3. System generates Thẻ tính giá thành per Circular 99 |
| | 4. Report shows: materials, labor, overhead, total, unit cost per object |
| | 5. Export as PDF or Excel |
| **Alternative Paths** | 3a. Single object report |
| | 3b. All objects report |
| **Exception Paths** | No costing data → error "Run costing first" |
| **Postcondition** | Report generated, downloadable |

### UC-CA-012: Generate Cost Variance Report

| Field | Value |
|-------|-------|
| **UC编号** | UC-CA-012 |
| **Name** | Cost Variance Analysis Report |
| **Actor** | CFO / Management Accountant |
| **Precondition** | Standard costing method used, actual costing completed |
| **Happy Path** | 1. Select period and cost object(s) |
| | 2. System compares: actual cost vs standard cost |
| | 3. System calculates: material variance, labor variance, overhead variance |
| | 4. System identifies favorable (F) and unfavorable (U) variances |
| | 5. Display variance report with drill-down |
| | 6. Export as PDF |
| **Alternative Paths** | 4a. Variance by cost element (material, labor, overhead separately) |
| **Exception Paths** | Standard costs not defined → error "Set standard costs first" |
| **Postcondition** | Variance report generated |

### UC-CA-013: Warehouse Keeper Issues Materials to Cost Object

| Field | Value |
|-------|-------|
| **UC编号** | UC-CA-013 |
| **Name** | Issue Materials to Cost Object |
| **Actor** | Warehouse Keeper |
| **Precondition** | Cost object exists, inventory available |
| **Happy Path** | 1. Warehouse keeper creates issuance document |
| | 2. Select cost object from dropdown |
| | 3. Select materials, enter quantities |
| | 4. System creates issuance journal (Dr 154x.1, Cr 155/156) |
| | 5. System updates inventory |
| | 6. Cost appears in cost object's material collection |
| **Alternative Paths** | 2a. Cost object not selected → materials go to unallocated pool |
| **Exception Paths** | Insufficient inventory → error |
| **Postcondition** | Materials issued, cost recorded against cost object |

### UC-CA-014: View Cost Object Cost Summary

| Field | Value |
|-------|-------|
| **UC编号** | UC-CA-014 |
| **Name** | View Cost Object Cost Summary |
| **Actor** | Cost Accountant / CFO |
| **Precondition** | Costing calculated for period |
| **Happy Path** | 1. Navigate to Cost Objects > Select object > Cost Summary |
| | 2. System displays: direct materials, direct labor, direct expenses, allocated overhead |
| | 3. System shows total cost, unit cost |
| | 4. System shows trend across periods (last 12 months) |
| | 5. Drill-down into each cost component |
| **Alternative Paths** | 4a. Compare across cost objects (multi-select) |
| **Exception Paths** | No data for selected period → display "No costing data" |
| **Postcondition** | Cost summary displayed |

### UC-CA-015: Configure Allocation Bases

| Field | Value |
|-------|-------|
| **UC编号** | UC-CA-015 |
| **Name** | Configure Cost Allocation Bases |
| **Actor** | Cost Accountant |
| **Precondition** | User has admin cost accounting permission |
| **Happy Path** | 1. Navigate to Cost Accounting > Settings > Allocation Bases |
| | 2. Create/edit allocation bases |
| | 3. For each base: name, type (material_value, labor_cost, machine_hours, custom) |
| | 4. For custom: select GL account to sum |
| | 5. Assign base to cost pool(s) |
| | 6. Save |
| **Alternative Paths** | 4a. Machine hours: enter manually per period |
| | 4b. Floor area: enter square meters per cost object |
| **Exception Paths** | Duplicate base name → 409 Conflict |
| **Postcondition** | Allocation bases configured, usable in allocation runs |

### UC-CA-016: Run Continuous Process Costing

| Field | Value |
|-------|-------|
| **UC编号** | UC-CA-016 |
| **Name** | Run Process Costing for Continuous Production |
| **Actor** | Cost Accountant |
| **Precondition** | Process costing method selected, production steps defined |
| **Happy Path** | 1. Define production steps (e.g., Mixing → Baking → Packaging) |
| | 2. For each step: enter units started, units completed, WIP units, % completion |
| | 3. System calculates equivalent units per step |
| | 4. System accumulates costs per step |
| | 5. System calculates cost per equivalent unit |
| | 6. System transfers costs between steps |
| | 7. Display cost per unit at each step |
| **Alternative Paths** | 3a. Weighted average method |
| | 3b. FIFO method |
| **Exception Paths** | Step not defined → error "Define production steps first" |
| **Postcondition** | Process costing completed |

### UC-CA-017: Handle By-Product Costing

| Field | Value |
|-------|-------|
| **UC编号** | UC-CA-017 |
| **Name** | Handle By-Product in Costing |
| **Actor** | Cost Accountant |
| **Precondition** | By-product cost objects defined (type = by_product) |
| **Happy Path** | 1. Define by-product with NRV (net realizable value) |
| | 2. System runs normal costing for all products |
| | 3. System deducts by-product NRV from main product cost |
| | 4. System adjusts main product unit cost |
| | 5. Display: total cost, by-product deduction, adjusted main product cost |
| **Alternative Paths** | 3a. Multiple by-products: deduct each separately |
| **Exception Paths** | By-product NRV > main product cost → error "By-product value exceeds main cost" |
| **Postcondition** | By-product value deducted, main product cost adjusted |

---

## 8. Data Flow Diagrams

### 8.1 Cost Collection Flow

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Warehouse   │     │   Payroll   │     │ Fixed Assets│
│  Module      │     │   Module    │     │  Module     │
└──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                    │                    │
       │ Material Issuance  │ Direct Labor       │ Depreciation
       │ Documents          │ Costs              │ Allocation
       ▼                    ▼                    ▼
┌──────────────────────────────────────────────────────────┐
│                    COST COLLECTION                        │
│                                                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────────────────┐   │
│  │  TK 621  │  │  TK 622  │  │  TK 625 / 627       │   │
│  │  Direct  │  │  Direct  │  │  Overhead / Depr.    │   │
│  │  Material│  │  Labor   │  │                      │   │
│  └────┬─────┘  └────┬─────┘  └──────────┬───────────┘   │
│       │              │                    │               │
└───────┼──────────────┼────────────────────┼───────────────┘
        │              │                    │
        ▼              ▼                    ▼
┌──────────────────────────────────────────────────────────┐
│                   COST POOLS                             │
│                                                          │
│  621 → Direct Material Pool                              │
│  622 → Direct Labor Pool                                 │
│  623 → Machinery Pool                                    │
│  624 → Tools Pool                                        │
│  625 → Depreciation Pool                                 │
│  626 → Social Insurance Pool                             │
│  627 → Manufacturing Overhead Pool                       │
└────────────────────────┬─────────────────────────────────┘
                         │
                         ▼
┌──────────────────────────────────────────────────────────┐
│              ALLOCATION ENGINE                           │
│                                                          │
│  For each Cost Pool:                                     │
│    1. Determine allocation base                          │
│    2. Calculate coefficient or ratio                     │
│    3. Allocate to each Cost Object                       │
│    4. Generate journal entry                             │
│                                                          │
│  Direct costs (621, 622) → traced directly to objects    │
│  Indirect costs (623-627) → allocated via base           │
└────────────────────────┬─────────────────────────────────┘
                         │
                         ▼
┌──────────────────────────────────────────────────────────┐
│                COST OBJECTS                              │
│                                                          │
│  ┌────────┐  ┌────────┐  ┌────────┐  ┌────────┐        │
│  │Prod A  │  │Prod B  │  │Proj C  │  │Svc D   │        │
│  │1542.1  │  │1542.1  │  │1541.1  │  │1543.1  │        │
│  │Material│  │Material│  │Material│  │Material│        │
│  │1542.2  │  │1542.2  │  │1541.2  │  │1543.2  │        │
│  │Labor   │  │Labor   │  │Labor   │  │Labor   │        │
│  │1542.4  │  │1542.4  │  │1541.4  │  │1543.4  │        │
│  │OH      │  │OH      │  │OH      │  │OH      │        │
│  └────────┘  └────────┘  └────────┘  └────────┘        │
└──────────────────────────────────────────────────────────┘
```

### 8.2 Period-End Costing Flow

```
┌─────────────────────────────────────────────────────────────┐
│                   PERIOD-END COSTING                         │
│                                                             │
│  ┌─────────────────────────────────────────────────────────┐│
│  │ STEP 1: Open Costing Period                              ││
│  │ - Create CostingPeriod record (status = open)            ││
│  │ - Link to fiscal period                                  ││
│  └────────────────────────────┬────────────────────────────┘│
│                               │                             │
│  ┌────────────────────────────▼────────────────────────────┐│
│  │ STEP 2: Collect Costs                                   ││
│  │ - Pull material issuances → TK 621                       ││
│  │ - Pull labor costs → TK 622                              ││
│  │ - Pull overhead entries → TK 623-627                     ││
│  │ - Create CostPool balances                               ││
│  └────────────────────────────┬────────────────────────────┘│
│                               │                             │
│  ┌────────────────────────────▼────────────────────────────┐│
│  │ STEP 3: Allocate Indirect Costs                         ││
│  │ - For each indirect cost pool:                           ││
│  │   - Read allocation base per cost object                 ││
│  │   - Calculate coefficient/ratio                          ││
│  │   - Create CostAllocation records                        ││
│  │   - Generate journal entries (Dr 154x, Cr 62x)          ││
│  └────────────────────────────┬────────────────────────────┘│
│                               │                             │
│  ┌────────────────────────────▼────────────────────────────┐│
│  │ STEP 4: Calculate Costing                               ││
│  │ - Sum all costs per cost object                          ││
│  │ - Calculate unit cost                                    ││
│  │ - Identify WIP vs completed                              ││
│  │ - Create CostingResult records                           ││
│  └────────────────────────────┬────────────────────────────┘│
│                               │                             │
│  ┌────────────────────────────▼────────────────────────────┐│
│  │ STEP 5: WIP Valuation                                   ││
│  │ - Enter WIP quantities and completion %                  ││
│  │ - Calculate WIP value                                    ││
│  │ - Generate WIP transfer journal                          ││
│  └────────────────────────────┬────────────────────────────┘│
│                               │                             │
│  ┌────────────────────────────▼────────────────────────────┐│
│  │ STEP 6: Close Period                                    ││
│  │ - Verify all entries posted                              ││
│  │ - Set period status = closed                             ││
│  │ - Lock all records                                       ││
│  └─────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────┘
```

### 8.3 WIP Evaluation Flow

```
┌────────────────────────────────────────────────────────────┐
│                  WIP EVALUATION FLOW                        │
│                                                            │
│  ┌──────────────────┐                                      │
│  │ Start of Period  │                                      │
│  │ WIP Opening Bal  │                                      │
│  └────────┬─────────┘                                      │
│           │                                                │
│           ▼                                                │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ Add: Current Period Costs                            │  │
│  │                                                      │  │
│  │  Direct Materials (from 621)                         │  │
│  │  + Direct Labor (from 622)                           │  │
│  │  + Direct Expenses (from 623, 624)                   │  │
│  │  + Allocated Overhead (from 627)                     │  │
│  │                                                      │  │
│  │  = Total WIP Cost                                    │  │
│  └────────┬─────────────────────────────────────────────┘  │
│           │                                                │
│           ▼                                                │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ Allocate: WIP → Completed & WIP Closing              │  │
│  │                                                      │  │
│  │  Completed units → Transfer to TK 155                │  │
│  │  (Dr 155, Cr 154)                                    │  │
│  │                                                      │  │
│  │  WIP closing → Carry forward to next period          │  │
│  │  (Dr 900, Cr 154)                                    │  │
│  └────────┬─────────────────────────────────────────────┘  │
│           │                                                │
│           ▼                                                │
│  ┌──────────────────┐                                      │
│  │ End of Period    │                                      │
│  │ WIP Closing Bal  │                                      │
│  │ = Opening for    │                                      │
│  │   Next Period    │                                      │
│  └──────────────────┘                                      │
└────────────────────────────────────────────────────────────┘
```

---

## 9. API Endpoints

### 9.1 Cost Objects

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/cost-objects` | Create cost object |
| GET | `/api/v1/cost-objects?company_id=X` | List cost objects |
| GET | `/api/v1/cost-objects/:id` | Get cost object by ID |
| GET | `/api/v1/cost-objects/hierarchy?company_id=X` | Get cost object hierarchy |
| PUT | `/api/v1/cost-objects/:id` | Update cost object |
| DELETE | `/api/v1/cost-objects/:id` | Delete cost object |

**POST /api/v1/cost-objects — Request:**
```json
{
  "code": "SP-001",
  "name": "Product A",
  "type": "product",
  "cost_center_id": "CC-123",
  "parent_id": "",
  "gl_account_code": "1542",
  "unit_of_measure": "piece",
  "standard_material": 500000,
  "standard_labor": 200000,
  "standard_overhead": 150000
}
```

**Response:**
```json
{
  "id": "CO-1691234567890",
  "company_id": "COMP-001",
  "code": "SP-001",
  "name": "Product A",
  "type": "product",
  "cost_center_id": "CC-123",
  "parent_id": "",
  "gl_account_code": "1542",
  "unit_of_measure": "piece",
  "standard_material": 500000,
  "standard_labor": 200000,
  "standard_overhead": 150000,
  "is_active": true,
  "created_at": "2026-08-07T10:00:00Z",
  "updated_at": "2026-08-07T10:00:00Z"
}
```

### 9.2 Cost Pools

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/cost-pools` | Create cost pool |
| GET | `/api/v1/cost-pools?company_id=X` | List cost pools |
| GET | `/api/v1/cost-pools/:id` | Get cost pool |
| PUT | `/api/v1/cost-pools/:id` | Update cost pool |
| DELETE | `/api/v1/cost-pools/:id` | Delete cost pool |

**POST /api/v1/cost-pools — Request:**
```json
{
  "code": "POOL-OH",
  "name": "Manufacturing Overhead Pool",
  "gl_account_code": "627",
  "pool_type": "overhead",
  "allocation_base": "material_value"
}
```

### 9.3 Allocation Bases

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/allocation-bases` | Create allocation base |
| GET | `/api/v1/allocation-bases?company_id=X` | List allocation bases |
| PUT | `/api/v1/allocation-bases/:id` | Update allocation base |
| DELETE | `/api/v1/allocation-bases/:id` | Delete allocation base |

### 9.4 Costing Periods

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/costing-periods` | Create/open costing period |
| GET | `/api/v1/costing-periods?company_id=X` | List costing periods |
| GET | `/api/v1/costing-periods/:id` | Get costing period |
| POST | `/api/v1/costing-periods/:id/collect` | Collect costs into pools |
| POST | `/api/v1/costing-periods/:id/allocate` | Run cost allocation |
| POST | `/api/v1/costing-periods/:id/calculate` | Calculate costing results |
| POST | `/api/v1/costing-periods/:id/wip` | Run WIP valuation |
| POST | `/api/v1/costing-periods/:id/close` | Close costing period |
| POST | `/api/v1/costing-periods/:id/reopen` | Reopen costing period |
| POST | `/api/v1/costing-periods/:id/recalculate` | Recalculate (reverse + re-run) |

**POST /api/v1/costing-periods — Request:**
```json
{
  "period_number": 8,
  "year": 2026,
  "fiscal_year_id": "FY-2026",
  "costing_method": "coefficient"
}
```

**POST /api/v1/costing-periods/:id/collect — Request:**
```json
{
  "collect_direct_materials": true,
  "collect_direct_labor": true,
  "collect_overhead": true
}
```

**POST /api/v1/costing-periods/:id/allocate — Request:**
```json
{
  "pool_ids": ["POOL-MAT", "POOL-LAB", "POOL-OH"],
  "base_id": "BASE-MAT-VAL"
}
```

**Response:**
```json
{
  "period_id": "CP-2026-08",
  "status": "allocated",
  "allocations_created": 45,
  "total_allocated": 350000000,
  "journal_entries_posted": 45,
  "details": [
    {
      "pool_id": "POOL-OH",
      "pool_name": "Manufacturing Overhead",
      "total_pool_balance": 300000000,
      "allocation_base": "material_value",
      "total_base": 500000000,
      "coefficient": 0.6,
      "objects_allocated": 15
    }
  ]
}
```

**POST /api/v1/costing-periods/:id/calculate — Response:**
```json
{
  "period_id": "CP-2026-08",
  "results": [
    {
      "cost_object_id": "CO-001",
      "code": "SP-001",
      "name": "Product A",
      "direct_materials": 100000000,
      "direct_labor": 40000000,
      "direct_expenses": 10000000,
      "allocated_overhead": 60000000,
      "total_cost": 210000000,
      "units_produced": 1000,
      "unit_cost": 210000,
      "wip_value": 0,
      "completed_value": 210000000
    }
  ],
  "summary": {
    "total_direct_materials": 500000000,
    "total_direct_labor": 200000000,
    "total_direct_expenses": 50000000,
    "total_allocated_overhead": 300000000,
    "grand_total": 1050000000
  }
}
```

### 9.5 Costing Results

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/costing-results?company_id=X&period_id=Y` | Get costing results for period |
| GET | `/api/v1/costing-results/:id` | Get single costing result |
| GET | `/api/v1/costing-results/summary?company_id=X&period_id=Y` | Cost summary report |
| GET | `/api/v1/costing-results/variance?company_id=X&period_id=Y` | Variance analysis |

### 9.6 Reports

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/cost-reports/thẻ-tính-giá-thành?period_id=X&object_id=Y` | Cost Calculation Sheet (PDF) |
| GET | `/api/v1/cost-reports/wip-valuation?period_id=X` | WIP Valuation Report |
| GET | `/api/v1/cost-reports/cost-summary?period_id=X` | Cost Summary by Object |
| GET | `/api/v1/cost-reports/cost-trend?object_id=X&months=12` | Cost Trend Report |
| GET | `/api/v1/cost-reports/period-comparison?period_from=X&period_to=Y` | Period Comparison |

### 9.7 Warehouse Integration

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/warehouse/issuance` | Issue materials to cost object (existing + extension) |

**POST /api/v1/warehouse/issuance — Request:**
```json
{
  "warehouse_id": "WH-001",
  "cost_object_id": "CO-001",
  "items": [
    {
      "material_code": "VLU-001",
      "quantity": 100,
      "unit_cost": 50000
    }
  ],
  "description": "Issue materials for Product A production"
}
```

---

## 10. Database Schema

### 10.1 Migration File: `000032_cost_accounting.up.sql`

```sql
-- Cost Accounting (Giá thành) module
-- Migration: 000032_cost_accounting

-- ─── Cost Objects ────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS cost_objects (
    id              TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    company_id      TEXT NOT NULL,
    code            TEXT NOT NULL,
    name            TEXT NOT NULL,
    type            TEXT NOT NULL DEFAULT 'product', -- product, service, project, batch, by_product
    cost_center_id  TEXT,
    parent_id       TEXT,
    gl_account_code TEXT NOT NULL DEFAULT '1542', -- 1541, 1542, 1543, 1544
    unit_of_measure TEXT NOT NULL DEFAULT 'piece',
    standard_material NUMERIC(19,4) DEFAULT 0,
    standard_labor    NUMERIC(19,4) DEFAULT 0,
    standard_overhead NUMERIC(19,4) DEFAULT 0,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, code)
);

CREATE INDEX IF NOT EXISTS idx_cost_objects_company ON cost_objects(company_id);
CREATE INDEX IF NOT EXISTS idx_cost_objects_parent ON cost_objects(parent_id);
CREATE INDEX IF NOT EXISTS idx_cost_objects_type ON cost_objects(type);
CREATE INDEX IF NOT EXISTS idx_cost_objects_gl ON cost_objects(gl_account_code);

-- ─── Cost Pools ──────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS cost_pools (
    id              TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    company_id      TEXT NOT NULL,
    code            TEXT NOT NULL,
    name            TEXT NOT NULL,
    gl_account_code TEXT NOT NULL, -- 621-627
    pool_type       TEXT NOT NULL, -- direct_material, direct_labor, machinery, tools, depreciation, social_insurance, overhead
    allocation_base TEXT,          -- material_value, labor_cost, machine_hours, custom
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, code)
);

CREATE INDEX IF NOT EXISTS idx_cost_pools_company ON cost_pools(company_id);
CREATE INDEX IF NOT EXISTS idx_cost_pools_gl ON cost_pools(gl_account_code);

-- ─── Allocation Bases ────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS allocation_bases (
    id              TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    company_id      TEXT NOT NULL,
    name            TEXT NOT NULL,
    basis_type      TEXT NOT NULL, -- material_value, labor_cost, machine_hours, direct_cost_total, custom
    gl_account_code TEXT,          -- for custom: which GL account to sum
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_allocation_bases_company ON allocation_bases(company_id);

-- ─── Costing Periods ─────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS costing_periods (
    id              TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    company_id      TEXT NOT NULL,
    period_number   INT NOT NULL, -- 1-13
    year            INT NOT NULL,
    fiscal_year_id  TEXT NOT NULL,
    status          TEXT NOT NULL DEFAULT 'open', -- open, running, closed
    costing_method  TEXT NOT NULL DEFAULT 'coefficient',
    allocated_total NUMERIC(19,4) DEFAULT 0,
    closed_at       TIMESTAMPTZ,
    closed_by       TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, period_number, year)
);

CREATE INDEX IF NOT EXISTS idx_costing_periods_company ON costing_periods(company_id);
CREATE INDEX IF NOT EXISTS idx_costing_periods_status ON costing_periods(status);
CREATE INDEX IF NOT EXISTS idx_costing_periods_year ON costing_periods(year);

-- ─── Cost Allocations ────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS cost_allocations (
    id                  TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    company_id          TEXT NOT NULL,
    costing_period_id   TEXT NOT NULL,
    cost_pool_id        TEXT NOT NULL,
    cost_object_id      TEXT NOT NULL,
    allocated_amount    NUMERIC(19,4) NOT NULL DEFAULT 0,
    allocation_ratio    NUMERIC(10,6) DEFAULT 0,
    allocation_base     TEXT,
    base_amount         NUMERIC(19,4) DEFAULT 0,
    journal_entry_id    TEXT, -- FK to journal_entries
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (costing_period_id) REFERENCES costing_periods(id),
    FOREIGN KEY (cost_pool_id) REFERENCES cost_pools(id),
    FOREIGN KEY (cost_object_id) REFERENCES cost_objects(id)
);

CREATE INDEX IF NOT EXISTS idx_cost_allocations_period ON cost_allocations(costing_period_id);
CREATE INDEX IF NOT EXISTS idx_cost_allocations_pool ON cost_allocations(cost_pool_id);
CREATE INDEX IF NOT EXISTS idx_cost_allocations_object ON cost_allocations(cost_object_id);
CREATE INDEX IF NOT EXISTS idx_cost_allocations_company ON cost_allocations(company_id);

-- ─── Costing Results ─────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS costing_results (
    id                  TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    company_id          TEXT NOT NULL,
    costing_period_id   TEXT NOT NULL,
    cost_object_id      TEXT NOT NULL,
    direct_materials    NUMERIC(19,4) DEFAULT 0,
    direct_labor        NUMERIC(19,4) DEFAULT 0,
    direct_expenses     NUMERIC(19,4) DEFAULT 0,
    allocated_overhead  NUMERIC(19,4) DEFAULT 0,
    total_cost          NUMERIC(19,4) DEFAULT 0,
    units_produced      NUMERIC(19,4) DEFAULT 0,
    unit_cost           NUMERIC(19,4) DEFAULT 0,
    wip_value           NUMERIC(19,4) DEFAULT 0,
    completed_value     NUMERIC(19,4) DEFAULT 0,
    cogs                NUMERIC(19,4) DEFAULT 0,
    variance            NUMERIC(19,4) DEFAULT 0,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(costing_period_id, cost_object_id),
    FOREIGN KEY (costing_period_id) REFERENCES costing_periods(id),
    FOREIGN KEY (cost_object_id) REFERENCES cost_objects(id)
);

CREATE INDEX IF NOT EXISTS idx_costing_results_period ON costing_results(costing_period_id);
CREATE INDEX IF NOT EXISTS idx_costing_results_object ON costing_results(cost_object_id);
CREATE INDEX IF NOT EXISTS idx_costing_results_company ON costing_results(company_id);

-- ─── Cost Pool Balances (per period) ─────────────────────────────────

CREATE TABLE IF NOT EXISTS cost_pool_balances (
    id                  TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    company_id          TEXT NOT NULL,
    costing_period_id   TEXT NOT NULL,
    cost_pool_id        TEXT NOT NULL,
    opening_balance     NUMERIC(19,4) DEFAULT 0,
    period_debits       NUMERIC(19,4) DEFAULT 0,
    period_credits      NUMERIC(19,4) DEFAULT 0,
    closing_balance     NUMERIC(19,4) DEFAULT 0,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(costing_period_id, cost_pool_id),
    FOREIGN KEY (costing_period_id) REFERENCES costing_periods(id),
    FOREIGN KEY (cost_pool_id) REFERENCES cost_pools(id)
);

CREATE INDEX IF NOT EXISTS idx_cost_pool_balances_period ON cost_pool_balances(costing_period_id);

-- ─── Process Steps (for continuous process costing) ──────────────────

CREATE TABLE IF NOT EXISTS process_steps (
    id                  TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    company_id          TEXT NOT NULL,
    cost_object_id      TEXT NOT NULL,
    step_number         INT NOT NULL,
    step_name           TEXT NOT NULL,
    units_started       NUMERIC(19,4) DEFAULT 0,
    units_completed     NUMERIC(19,4) DEFAULT 0,
    wip_units           NUMERIC(19,4) DEFAULT 0,
    completion_percent  NUMERIC(5,2) DEFAULT 0,
    equivalent_units    NUMERIC(19,4) DEFAULT 0,
    cost_assigned       NUMERIC(19,4) DEFAULT 0,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(cost_object_id, step_number),
    FOREIGN KEY (cost_object_id) REFERENCES cost_objects(id)
);

-- ─── Add cost_type to cost_centers ──────────────────────────────────

ALTER TABLE cost_centers ADD COLUMN IF NOT EXISTS cost_type TEXT DEFAULT '';
```

### 10.2 Migration File: `000032_cost_accounting.down.sql`

```sql
ALTER TABLE cost_centers DROP COLUMN IF EXISTS cost_type;

DROP TABLE IF EXISTS process_steps;
DROP TABLE IF EXISTS cost_pool_balances;
DROP TABLE IF EXISTS costing_results;
DROP TABLE IF EXISTS cost_allocations;
DROP TABLE IF EXISTS costing_periods;
DROP TABLE IF EXISTS allocation_bases;
DROP TABLE IF EXISTS cost_pools;
DROP TABLE IF EXISTS cost_objects;
```

---

## 11. Workflow Diagrams

### 11.1 Period-End Costing Workflow

```
START
  │
  ▼
[1] Open Costing Period
  │   - Create record with period_number, year, method
  │   - Validate no other open period for same period+year
  │
  ▼
[2] Collect Direct Material Costs
  │   - Query warehouse issuance docs for period
  │   - Group by cost_object_id
  │   - Sum quantities × unit_costs
  │   - Create CostAllocation: pool=621, per object
  │
  ▼
[3] Collect Direct Labor Costs
  │   - Query payroll entries for period
  │   - Filter labor_type = "direct"
  │   - Group by cost_object_id
  │   - Create CostAllocation: pool=622, per object
  │
  ▼
[4] Collect Manufacturing Overhead
  │   - Query GL journal entries for accounts 623-627
  │   - Group by account_code
  │   - Create CostPoolBalances
  │
  ▼
[5] Allocate Indirect Costs
  │   ┌─────────────────────────────────┐
  │   │ For each indirect cost pool:    │
  │   │   1. Read allocation base data  │
  │   │   2. Calculate coefficient      │
  │   │   3. Allocate to each object    │
  │   │   4. Create CostAllocation      │
  │   │   5. Post journal entry         │
  │   └─────────────────────────────────┘
  │
  ▼
[6] Calculate Costing Results
  │   - Sum: DM + DL + DE + AllocatedOH per object
  │   - Calculate unit_cost = total / units_produced
  │   - Identify WIP vs completed
  │   - Create CostingResult per object
  │
  ▼
[7] WIP Valuation
  │   - Enter WIP completion percentages
  │   - Calculate: wip_value = total_cost × % × wip_units
  │   - Post WIP transfer journal (Dr 154x, Cr 900)
  │
  ▼
[8] Close Period
  │   - Verify all entries posted
  │   - Set status = "closed"
  │   - Lock all records
  │
  ▼
END
```

### 11.2 Cost Allocation Workflow

```
START
  │
  ▼
[1] Select Cost Pool(s) to allocate
  │
  ▼
[2] Validate pool has balance > 0
  │   IF balance = 0 → SKIP
  │
  ▼
[3] Select Allocation Base
  │   - material_value: sum of direct material per object
  │   - labor_cost: sum of direct labor per object
  │   - machine_hours: manual entry per object
  │   - custom: sum of specified GL account per object
  │
  ▼
[4] Calculate Total Base Amount
  │   total_base = Σ(base_amount for all objects)
  │   IF total_base = 0 → ERROR "No allocation base"
  │
  ▼
[5] Calculate Coefficient
  │   coefficient = pool_balance / total_base
  │
  ▼
[6] Allocate to Each Object
  │   ┌──────────────────────────────────────┐
  │   │ For each cost object:                │
  │   │   ratio = object_base / total_base   │
  │   │   amount = coefficient × object_base │
  │   │   Create CostAllocation record       │
  │   └──────────────────────────────────────┘
  │
  ▼
[7] Handle Rounding Differences
  │   IF Σallocated ≠ pool_balance:
  │     Adjust largest allocation by difference
  │
  ▼
[8] Generate Journal Entry
  │   Dr  154x.4  {total_amount}  Chi phí SX chung
  │       Cr  627  {total_amount}  Chi phí sản xuất chung
  │
  ▼
[9] Zero Pool Balance
  │   pool.closing_balance = 0
  │
  ▼
END
```

### 11.3 WIP Valuation Workflow

```
START
  │
  ▼
[1] Identify Cost Objects with WIP
  │   - Query: units_started > units_completed
  │   - Or: wip_flag = true
  │
  ▼
[2] For Each WIP Cost Object:
  │
  │   [2a] Enter completion percentage
  │        completion% = 0 to 100
  │
  │   [2b] Calculate WIP units
  │        wip_units = units_started - units_completed
  │
  │   [2c] Calculate equivalent units (if process method)
  │        equiv_units = completed + (wip_units × completion%)
  │
  │   [2d] Calculate WIP value
  │        wip_value = total_cost × (wip_units / equiv_units)
  │
  │   [2e] Calculate completed value
  │        completed_value = total_cost - wip_value
  │
  │   [2f] Post journal entry
  │        Dr  154x  {wip_value}   Chi phí SXKD dở dang
  │            Cr  900  {wip_value}   Thu nhập chưa thực hiện
  │
  ▼
[3] Update CostingResult.wip_value for each object
  │
  ▼
[4] Display WIP summary report
  │
  ▼
END
```

---

## 12. User Journeys

### 12.1 Cost Accountant — Monthly Costing Process

```
DAY 1-25: ONGOING COST COLLECTION
  │
  ├── Warehouse issues materials to cost objects (via warehouse module)
  ├── Payroll processes direct labor entries
  ├── Fixed assets depreciation runs
  └── Overhead journals posted to GL
  │
DAY 26-28: MONTH-END COSTING
  │
  ├── [1] Open costing period (Period 8, 2026, coefficient method)
  ├── [2] Collect costs: review material issuances, labor costs, overhead
  │       - Fix any unallocated items
  │       - Confirm cost pool balances
  ├── [3] Run allocation: select base, calculate coefficients
  │       - Review allocation results
  │       - Confirm journal entries
  ├── [4] Calculate costing: review unit costs per product
  │       - Compare with previous period
  │       - Flag anomalies
  ├── [5] WIP valuation: enter completion percentages
  │       - Review WIP values
  ├── [6] Generate reports: Thẻ tính giá thành, cost summary
  ├── [7] Close period
  └── [8] Send reports to CFO for review
  │
DAY 29-30: REVIEW & CLOSE
  │
  ├── CFO reviews cost reports
  ├── Variance analysis review
  └── Period locked
```

### 12.2 CFO — Cost Analysis Review

```
MONTHLY ROUTINE:
  │
  ├── [1] Open Cost Accounting dashboard
  │       - View total costs per period
  │       - View cost trends (12 months)
  ├── [2] Compare periods
  │       - Period 8 vs Period 7
  │       - Which products had cost increases?
  ├── [3] Variance analysis (if standard method)
  │       - Material variance: price vs quantity
  │       - Labor variance: rate vs efficiency
  │       - Overhead variance: spending vs volume
  ├── [4] Drill-down into specific products
  │       - Product A: why did unit cost increase 5%?
  │       - Compare: direct materials, labor, overhead
  ├── [5] Review profitability by product
  │       - Unit cost vs selling price
  │       - Margin analysis
  └── [6] Export reports for board meeting
```

### 12.3 Warehouse Keeper — Material Issuance to Cost Object

```
DAILY ROUTINE:
  │
  ├── [1] Receive production request
  │       - "Issue 500kg of VLU-001 to Product A (CO-001)"
  ├── [2] Open Warehouse > Issue Materials
  │       - Select warehouse
  │       - Select cost object: Product A (CO-001)
  │       - Select material: VLU-001
  │       - Enter quantity: 500
  │       - System shows: available 2,000kg, unit cost 50,000 VND
  ├── [3] Confirm issuance
  │       - System creates: Dr 1542.1, Cr 155
  │       - System updates inventory
  │       - System records: CO-001 consumed 500kg × 50,000 = 25,000,000 VND
  └── [4] Print issuance slip (if needed)
```

---

## 13. Report Templates

### 13.1 Thẻ tính giá thành (Cost Calculation Sheet)

Per Circular 99, Appendix format:

```
╔══════════════════════════════════════════════════════════════════════╗
║                    THẺ TÍNH GIÁ THÀNH                              ║
║                    (Cost Calculation Sheet)                         ║
║                                                                     ║
║  Company: {company_name}     Period: {month}/{year}                ║
║  Product: {product_code} - {product_name}                          ║
║  Unit: {unit_of_measure}                                           ║
╠══════════════════════════════════════════════════════════════════════╣
║                                                                     ║
║  I. CHI PHÍ TRỰC TIẾP (Direct Costs)                              ║
║  ┌─────────────────────────────────┬───────────────┬────────────┐  ║
║  │ Hạng mục                        │ Kỳ này        │ Kỳ trước   │  ║
║  ├─────────────────────────────────┼───────────────┼────────────┤  ║
║  │ 1. Nguyên vật liệu trực tiếp    │ {dm_amount}   │ {dm_prev}  │  ║
║  │ 2. Nhân công trực tiếp          │ {dl_amount}   │ {dl_prev}  │  ║
║  │ 3. Chi phí trực tiếp khác       │ {de_amount}   │ {de_prev}  │  ║
║  │    Tổng chi phí trực tiếp       │ {total_direct}│ {dir_prev} │  ║
║  └─────────────────────────────────┴───────────────┴────────────┘  ║
║                                                                     ║
║  II. CHI PHÍ SẢN XUẤT CHUNG (Manufacturing Overhead)              ║
║  ┌─────────────────────────────────┬───────────────┬────────────┐  ║
║  │ Hạng mục                        │ Kỳ này        │ Kỳ trước   │  ║
║  ├─────────────────────────────────┼───────────────┼────────────┤  ║
║  │ 4. Chi phí sử dụng máy thi công│ {machinery}   │ {mach_prev}│  ║
║  │ 5. Chi phí dụng cụ công tác    │ {tools}       │ {tools_prev│  ║
║  │ 6. Khấu hao TSCĐ               │ {depreciation}│ {dep_prev} │  ║
║  │ 7. BHLD & BHXH                  │ {insurance}   │ {ins_prev} │  ║
║  │ 8. Chi phí sản xuất chung khác │ {other_oh}    │ {oth_prev} │  ║
║  │    Tổng chi phí SX chung       │ {total_oh}    │ {oh_prev}  │  ║
║  └─────────────────────────────────┴───────────────┴────────────┘  ║
║                                                                     ║
║  III. TỔNG CHI PHÍ GIÁ THÀNH (Total Cost)                         ║
║  ┌─────────────────────────────────┬───────────────┬────────────┐  ║
║  │ Tổng chi phí trực tiếp          │ {total_direct}│            │  ║
║  │ + Tổng chi phí SX chung        │ {total_oh}    │            │  ║
║  │ = Tổng chi phí giá thành       │ {grand_total} │            │  ║
║  │ Số lượng sản xuất               │ {units}       │            │  ║
║  │ Đơn giá thành phẩm              │ {unit_cost}   │            │  ║
║  └─────────────────────────────────┴───────────────┴────────────┘  ║
║                                                                     ║
║  Approved by: ________________    Prepared by: ________________    ║
╚══════════════════════════════════════════════════════════════════════╝
```

### 13.2 Cost Summary by Object

```
╔═══════════════════════════════════════════════════════════════════════════╗
║              TỔNG HỢP GIÁ THÀNH THEO ĐỐI TƯỢNG                        ║
║              (Cost Summary by Cost Object)                               ║
║                                                                          ║
║  Company: {company_name}     Period: {month}/{year}                     ║
╠═══════════════════════════════════════════════════════════════════════════╣
║                                                                          ║
║  ┌──────────┬────────────┬────────────┬────────────┬────────┬──────────┐║
║  │ ĐTTC     │ NVL TT     │ NC TT      │ CPC        │ TPCX   │ Tổng CP  │║
║  ├──────────┼────────────┼────────────┼────────────┼────────┼──────────┤║
║  │ SP-001   │ 100,000,000│ 40,000,000 │ 10,000,000 │ 60M    │ 210M     │║
║  │ SP-002   │ 150,000,000│ 60,000,000 │ 15,000,000 │ 90M    │ 315M     │║
║  │ SP-003   │ 50,000,000 │ 20,000,000 │ 5,000,000  │ 30M    │ 105M     │║
║  ├──────────┼────────────┼────────────┼────────────┼────────┼──────────┤║
║  │ TỔNG     │ 300,000,000│ 120,000,000│ 30,000,000 │ 180M   │ 630M     │║
║  └──────────┴────────────┴────────────┴────────────┴────────┴──────────┘║
║                                                                          ║
║  ĐTTC = Đối tượng tính chi phí    NVL TT = Nguyên vật liệu trực tiếp    ║
║  NC TT = Nhân công trực tiếp       CPC = Chi phí trực tiếp khác          ║
║  TPCX = Tổng phân bổ chung         Tổng CP = Tổng chi phí giá thành     ║
╚═══════════════════════════════════════════════════════════════════════════╝
```

### 13.3 WIP Valuation Report

```
╔═══════════════════════════════════════════════════════════════════════════╗
║              BÁO CÁO ĐÁNH GIÁ DỞ ĐANG CUỐI KỲ                         ║
║              (WIP Valuation Report)                                      ║
║                                                                          ║
║  Company: {company_name}     Period: {month}/{year}                     ║
╠═══════════════════════════════════════════════════════════════════════════╣
║                                                                          ║
║  ┌──────────┬────────┬────────┬───────────┬──────────┬─────────────────┐║
║  │ ĐTTC     │ SK BD  │ SK HT  │ SL DĐ     │ % HT     │ Giá trị DĐ     │║
║  ├──────────┼────────┼────────┼───────────┼──────────┼─────────────────┤║
║  │ SP-001   │ 500    │ 400    │ 100       │ 60%      │ 12,600,000      │║
║  │ SP-002   │ 300    │ 200    │ 100       │ 50%      │ 15,750,000      │║
║  ├──────────┼────────┼────────┼───────────┼──────────┼─────────────────┤║
║  │ TỔNG     │ 800    │ 600    │ 200       │          │ 28,350,000      │║
║  └──────────┴────────┴────────┴───────────┴──────────┴─────────────────┘║
║                                                                          ║
║  SK BD = Sản lượng bắt đầu    SK HT = Sản lượng hoàn thành             ║
║  SL DĐ = Số lượng dở dang     % HT = Hoàn thành                        ║
║  Giá trị DĐ = Giá trị dở dang cuối kỳ                                  ║
╚═══════════════════════════════════════════════════════════════════════════╝
```

### 13.4 Cost Variance Analysis

```
╔═══════════════════════════════════════════════════════════════════════════╗
║              PHÂN TÍCH CHÊNH LỆCH GIÁ THÀNH                            ║
║              (Cost Variance Analysis)                                    ║
║                                                                          ║
║  Company: {company_name}     Period: {month}/{year}                     ║
║  Method: Standard Costing                                                ║
╠═══════════════════════════════════════════════════════════════════════════╣
║                                                                          ║
║  ┌──────────┬────────────┬────────────┬────────────┬──────────┬────────┐║
║  │ ĐTTC     │ Chi phí TB │ Chi phí TT │ Chênh lệch │ Loại     │ %      │║
║  ├──────────┼────────────┼────────────┼────────────┼──────────┼────────┤║
║  │ SP-001   │ 200,000,000│ 210,000,000│ 10,000,000 │ Bất lợi  │ +5.0%  │║
║  │ SP-002   │ 320,000,000│ 315,000,000│ -5,000,000 │ Hữu lợi  │ -1.6%  │║
║  │ SP-003   │ 100,000,000│ 105,000,000│ 5,000,000  │ Bất lợi  │ +5.0%  │║
║  ├──────────┼────────────┼────────────┼────────────┼──────────┼────────┤║
║  │ TỔNG     │ 620,000,000│ 630,000,000│ 10,000,000 │          │ +1.6%  │║
║  └──────────┴────────────┴────────────┴────────────┴──────────┴────────┘║
║                                                                          ║
║  Chi tiết SP-001:                                                       ║
║  ┌──────────────────┬────────────┬────────────┬────────────┬──────────┐ ║
║  │ Hạng mục         │ Chi phí TB │ Chi phí TT │ Chênh lệch │ %        │ ║
║  ├──────────────────┼────────────┼────────────┼────────────┼──────────┤ ║
║  │ Nguyên vật liệu  │ 95,000,000 │ 100,000,000│ 5,000,000  │ +5.3%    │ ║
║  │ Nhân công        │ 40,000,000 │ 40,000,000 │ 0          │ 0.0%     │ ║
║  │ Phân bổ chung    │ 65,000,000 │ 70,000,000 │ 5,000,000  │ +7.7%    │ ║
║  └──────────────────┴────────────┴────────────┴────────────┴──────────┘ ║
║                                                                          ║
║  Chi phí TB = Chi phí tiêu chuẩn    Chi phí TT = Chi phí thực tế        ║
╚═══════════════════════════════════════════════════════════════════════════╝
```

---

## 14. Integration Points

### 14.1 GL Module Integration

```
┌──────────────────────────────────────────────────────────────┐
│                     GL INTEGRATION                           │
│                                                              │
│  Cost Accounting → GL:                                       │
│  ┌────────────────────────────────────────────────────────┐  │
│  │ Journal Entry: Cost Allocation                         │  │
│  │   Dr  154x.4  {amount}  Chi phí SX chung              │  │
│  │       Cr  627  {amount}  Chi phí sản xuất chung        │  │
│  │                                                        │  │
│  │ Journal Entry: WIP Transfer                            │  │
│  │   Dr  155     {amount}  Hàng thành phẩm                │  │
│  │       Cr  154x {amount}  Chi phí SXKD dở dang          │  │
│  │                                                        │  │
│  │ Journal Entry: WIP Valuation                           │  │
│  │   Dr  154x    {wip}     Chi phí SXKD dở dang           │  │
│  │       Cr  900  {wip}     Thu nhập chưa thực hiện       │  │
│  │                                                        │  │
│  │ Journal Entry: COGS                                    │  │
│  │   Dr  632     {amount}  Giá vốn hàng bán               │  │
│  │       Cr  155  {amount}  Hàng thành phẩm                │  │
│  └────────────────────────────────────────────────────────┘  │
│                                                              │
│  GL → Cost Accounting:                                       │
│  ┌────────────────────────────────────────────────────────┐  │
│  │ Read GL balances for accounts 621-627                  │  │
│  │ Read GL journal entries for period                     │  │
│  │ Read account structure for 154 sub-accounts            │  │
│  └────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────┘
```

### 14.2 Warehouse Module Integration

```
┌──────────────────────────────────────────────────────────────┐
│                   WAREHOUSE INTEGRATION                       │
│                                                              │
│  Warehouse → Cost Accounting:                                │
│  ┌────────────────────────────────────────────────────────┐  │
│  │ Issuance documents (pxk) for period                    │  │
│  │ Fields: cost_object_id, material_code, qty, unit_cost  │  │
│  │ → Cost Accounting collects as direct material cost     │  │
│  └────────────────────────────────────────────────────────┘  │
│                                                              │
│  Cost Accounting → Warehouse:                                │
│  ┌────────────────────────────────────────────────────────┐  │
│  │ Cost object code displayed in issuance records         │  │
│  │ Cost data visible in warehouse reporting               │  │
│  └────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────┘
```

### 14.3 Payroll Module Integration

```
┌──────────────────────────────────────────────────────────────┐
│                     PAYROLL INTEGRATION                       │
│                                                              │
│  Payroll → Cost Accounting:                                  │
│  ┌────────────────────────────────────────────────────────┐  │
│  │ Payroll journal entries for direct labor               │  │
│  │ Fields: cost_object_id, employee_id, labor_cost        │  │
│  │ → Cost Accounting collects as direct labor cost        │  │
│  │                                                        │  │
│  │ Social insurance (BHXH, BHYT) on direct labor          │  │
│  │ → Goes to TK 626 cost pool                             │  │
│  └────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────┘
```

### 14.4 Fixed Assets Module Integration

```
┌──────────────────────────────────────────────────────────────┐
│                  FIXED ASSETS INTEGRATION                     │
│                                                              │
│  Fixed Assets → Cost Accounting:                             │
│  ┌────────────────────────────────────────────────────────┐  │
│  │ Depreciation journal entries per period                 │  │
│  │ → Manufacturing depreciation → TK 625 cost pool        │  │
│  │ → Construction machinery → TK 623 cost pool            │  │
│  │                                                        │  │
│  │ Allocation by: machine hours, usage hours, or floor    │  │
│  │ area per cost object                                   │  │
│  └────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────┘
```

### 14.5 Budget Module Integration

```
┌──────────────────────────────────────────────────────────────┐
│                    BUDGET INTEGRATION                         │
│                                                              │
│  Budget → Cost Accounting:                                   │
│  ┌────────────────────────────────────────────────────────┐  │
│  │ Budgeted costs per cost object                         │  │
│  │ → Used as standard cost baseline                       │  │
│  │ → Budget vs actual comparison                          │  │
│  └────────────────────────────────────────────────────────┘  │
│                                                              │
│  Cost Accounting → Budget:                                   │
│  ┌────────────────────────────────────────────────────────┐  │
│  │ Actual costs per cost object per period                │  │
│  │ → Budget variance analysis                             │  │
│  │ → Budget adjustment recommendations                    │  │
│  └────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────┘
```

---

## 15. Implementation Priority

### Phase 1 — Core (Weeks 1-4) — MUST HAVE

| Week | Deliverables |
|------|-------------|
| 1 | Domain models: CostObject, CostPool, CostingPeriod, CostAllocation, CostingResult |
| 1 | Repository interfaces in `internal/domain/interfaces.go` |
| 1 | PostgreSQL + Memory repo implementations |
| 2 | Costing period management (open, close, reopen) |
| 2 | Direct cost collection (material + labor) |
| 3 | Simple costing method (giản đơn) |
| 3 | Coefficient costing method (hệ số) — default |
| 4 | Period-end journal entries |
| 4 | Basic Thẻ tính giá thành report |
| 4 | Migration: `000032_cost_accounting` |

**Exit criteria:** Cost accountant can run full monthly costing using simple or coefficient method.

### Phase 2 — Extended Methods (Weeks 5-8)

| Week | Deliverables |
|------|-------------|
| 5 | Ratio costing method (tỷ lệ) |
| 5 | Allocation base configuration |
| 6 | Overhead allocation from GL (TK 623-627) |
| 6 | WIP valuation with completion percentages |
| 7 | Cost summary reports (by object, by period) |
| 7 | WIP valuation report |
| 8 | Recalculate capability (reverse + re-run) |
| 8 | Period comparison report |

**Exit criteria:** Cost accountant can use 3 methods, WIP valuation working.

### Phase 3 — Advanced (Weeks 9-12)

| Week | Deliverables |
|------|-------------|
| 9 | Standard costing method (định mức) |
| 9 | Variance analysis (material, labor, overhead) |
| 10 | Continuous process costing method |
| 10 | Process step management (equivalent units) |
| 11 | Depreciation allocation integration |
| 11 | Social insurance allocation (TK 626) |
| 12 | Cost trend reports, drill-down analysis |

**Exit criteria:** All 4 core methods operational, variance analysis working.

### Phase 4 — Completeness (Weeks 13-16)

| Week | Deliverables |
|------|-------------|
| 13 | By-product costing method |
| 13 | Construction project costing (TK 1541) |
| 14 | Service costing (TK 1543) |
| 14 | Warranty cost tracking (TK 1544) |
| 15 | PDF report generation (Thẻ tính giá thành) |
| 15 | Excel export for all reports |
| 16 | End-to-end integration testing |
| 16 | Documentation + user guide |

**Exit criteria:** All 6 methods operational, all reports functional, full integration.

---

## 16. Risk Analysis

### 16.1 Technical Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Allocation rounding errors accumulate across thousands of objects | Medium | Medium | Use DECIMAL(19,4), adjust largest allocation for rounding difference, validate sum allocated = pool_balance |
| Period-end calculation takes too long with large datasets | Low | High | Batch processing, indexed queries, async job queue |
| Journal entry reversal fails mid-transaction | Low | High | Wrap in database transaction, idempotent reversal keys |
| Concurrent costing runs for same company | Low | High | Advisory lock on company_id during costing run |
| WIP completion percentages entered incorrectly | Medium | Medium | Validation rules, warning on >100%, comparison with previous period |
| Memory backend doesn't handle complex allocation logic | Medium | Low | In-memory implementation follows exact same algorithm as PG |

### 16.2 Regulatory Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Circular 99 interpretation changes | Low | High | Modular allocation engine, easy to swap methods |
| WIP valuation method not compliant | Medium | High | Follow Article 28 of Circular 99 exactly, maintain 3 valuation options |
| GL account mapping doesn't match new Chart of Accounts | Medium | Medium | Configurable mapping, not hardcoded |
| Audit finding on cost allocation methodology | Medium | High | Full journal trail, allocation coefficients documented, consistency checks |
| Industry-specific rules not covered | Medium | Medium | Phase 4 covers construction (1541), services (1543), warranty (1544) |

---

## 17. Appendix

### 17.1 Glossary

| Vietnamese | English | Definition |
|------------|---------|------------|
| Giá thành | Cost/Unit Cost | Total cost to produce one unit of product or service |
| Đối tượng tính chi phí | Cost Object | Unit of cost tracking (product, service, project, batch) |
| Hạng mục chi phí | Cost Pool | Collection point for costs of the same category before allocation |
| Cơ sở phân bổ | Allocation Base | Driver used to distribute overhead to cost objects |
| Hệ số phân bổ | Allocation Coefficient | Ratio used to allocate overhead = pool / total base |
| Chi phí trực tiếp | Direct Cost | Cost traceable to a specific cost object |
| Chi phí gián tiếp | Indirect Cost/Overhead | Cost shared across multiple cost objects |
| Dở dang | Work-in-Progress (WIP) | Unfinished production at period-end |
| Thành phẩm | Finished Goods | Completed products ready for sale |
| Nguyên vật liệu | Raw Materials | Direct materials used in production |
| Nhân công trực tiếp | Direct Labor | Labor directly involved in production |
| Khấu hao | Depreciation | Allocation of fixed asset cost over useful life |
| Giá vốn hàng bán | Cost of Goods Sold (COGS) | Cost of products sold during the period |
| Thẻ tính giá thành | Cost Calculation Sheet | Standard format for cost reporting per Circular 99 |
| Phân bổ | Allocation | Distribution of costs from pool to objects |
| Định mức | Standard/Norm | Predetermined cost per unit for budgeting and variance analysis |
| Sản phẩm phụ | By-product | Secondary product from same production process |
| Phân bước liên tục | Continuous Process | Production flow without distinct batches |
| Chênh lệch giá thành | Cost Variance | Difference between actual and standard cost |
| Thu nhập chưa thực hiện | Unearned Revenue | WIP valuation carry-forward account (TK 900) |

### 17.2 Reference Documents

| Document | Code | Relevance |
|----------|------|-----------|
| Thông tư 99/2025/TT-BTC | TT99/2025 | Enterprise Accounting Regime — primary regulation |
| Nghị định 123/2020/NĐ-CP | ND123/2020 | VAT and invoice regulations |
| Luật Kế toán 2015 | Law 88/2015/QH13 | Accounting Law |
| Luật Doanh nghiệp 2020 | Law 59/2020/QH14 | Enterprise Law |
| MISA SME 2026 | — | Benchmark for 6 costing methods |
| IAS 2 — Inventories | — | International standard for inventory valuation |
| IAS 11 — Construction Contracts | — | Construction costing (now IFRS 15) |

### 17.3 GL Account List from Circular 99

**WIP Accounts (TK 154):**
| Account | Name |
|---------|------|
| 154 | Chi phí SX, kinh doanh dở dang |
| 1541 | Chi phí xây lắp |
| 1542 | Chi phí sản phẩm khác |
| 1543 | Chi phí dịch vụ |
| 1544 | Chi phí bảo hành |

**Cost Pool Accounts (TK 62x):**
| Account | Name |
|---------|------|
| 621 | Chi phí nguyên vật liệu trực tiếp |
| 622 | Chi phí nhân công trực tiếp |
| 623 | Chi phí sử dụng máy thi công |
| 624 | Chi phí dụng cụ công tác |
| 625 | Chi phí khấu hao TSCĐ |
| 626 | Chi phí BHLD & BHXH nhân công trực tiếp |
| 627 | Chi phí sản xuất chung |

**Inventory Accounts:**
| Account | Name |
|---------|------|
| 155 | Hàng tồn kho thành phẩm |
| 156 | Hàng tồn kho bán thành phẩm |
| 157 | Hàng tồn kho nguyên vật liệu |
| 158 | Hàng tồn kho dụng cụ |
| 159 | Hàng tồn kho công cụ, dụng cụ |

**COGS:**
| Account | Name |
|---------|------|
| 632 | Giá vốn hàng bán |

---

**END OF DOCUMENT**
