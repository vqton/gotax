package domain

// CostCenter represents a organizational unit for cost tracking.
type CostCenter struct {
	ID          string `json:"id"`
	CompanyID   string `json:"company_id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	ParentID    string `json:"parent_id,omitempty"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CostingMethod defines the method used to calculate cost.
type CostingMethod string

const (
	CostingMethodSimple      CostingMethod = "SIMPLE"       // Giản đơn
	CostingMethodCoefficient CostingMethod = "COEFFICIENT"  // Hệ số
	CostingMethodProportion  CostingMethod = "PROPORTION"   // Tỷ lệ
	CostingMethodStandard    CostingMethod = "STANDARD"     // Định mức
	CostingMethodProcess     CostingMethod = "PROCESS"      // Phân bước liên tục
	CostingMethodByProduct   CostingMethod = "BY_PRODUCT"   // Loại trừ sản phẩm phụ
)

// CostObjectType classifies what the cost object tracks.
type CostObjectType string

const (
	CostObjectTypeProduct    CostObjectType = "PRODUCT"     // Thành phẩm
	CostObjectTypeService    CostObjectType = "SERVICE"     // Dịch vụ
	CostObjectTypeProject    CostObjectType = "PROJECT"     // Công trình / hạng mục
	CostObjectTypeOrder      CostObjectType = "ORDER"       // Đơn đặt hàng
	CostObjectTypeBatch      CostObjectType = "BATCH"       // Lô / batch sản xuất
)

// CostPoolStatus tracks the lifecycle of a cost pool.
type CostPoolStatus string

const (
	CostPoolStatusOpen   CostPoolStatus = "OPEN"
	CostPoolStatusClosed CostPoolStatus = "CLOSED"
)

// CostObject is the unit to which costs are allocated (sản phẩm, dịch vụ, công trình).
type CostObject struct {
	ID                string          `json:"id"`
	CompanyID         string          `json:"company_id"`
	Code              string          `json:"code"`
	Name              string          `json:"name"`
	Type              CostObjectType  `json:"type"`
	CostingMethod     CostingMethod   `json:"costing_method"`
	CostCenterID      string          `json:"cost_center_id,omitempty"`
	UnitOfMeasure     string          `json:"unit_of_measure,omitempty"`
	StandardCost      float64         `json:"standard_cost,omitempty"`
	StandardMaterial  float64         `json:"standard_material,omitempty"`
	StandardLabor     float64         `json:"standard_labor,omitempty"`
	StandardOverhead  float64         `json:"standard_overhead,omitempty"`
	PlanQuantity      float64         `json:"plan_quantity,omitempty"`
	CompletedUnits    float64         `json:"completed_units,omitempty"`
	WIPUnits          float64         `json:"wip_units,omitempty"`
	CompletionPct     float64         `json:"completion_pct,omitempty"`
	IsActive          bool            `json:"is_active"`
	CreatedAt         string          `json:"created_at"`
	UpdatedAt         string          `json:"updated_at"`
}

// CostPool collects costs before allocation to cost objects (tài khoản tập hợp chi phí).
type CostPool struct {
	ID            string         `json:"id"`
	CompanyID     string         `json:"company_id"`
	PeriodID      string         `json:"period_id"`
	GLAccountCode string         `json:"gl_account_code"` // e.g. 621, 622, 627
	Name          string         `json:"name"`
	Status        CostPoolStatus `json:"status"`
	TotalAmount   float64        `json:"total_amount"`
	CreatedAt     string         `json:"created_at"`
	UpdatedAt     string         `json:"updated_at"`
}

// CostPoolLine is a single cost entry within a pool.
type CostPoolLine struct {
	ID            string  `json:"id"`
	PoolID        string  `json:"pool_id"`
	SourceType    string  `json:"source_type"`    // JOURNAL, PAYROLL, DEPRECIATION, MANUAL
	SourceID      string  `json:"source_id,omitempty"`
	Description   string  `json:"description"`
	Amount        float64 `json:"amount"`
	CostCenterID  string  `json:"cost_center_id,omitempty"`
	CreatedAt     string  `json:"created_at"`
}

// CostingPeriod tracks monthly costing cycles.
type CostingPeriod struct {
	ID        string `json:"id"`
	CompanyID string `json:"company_id"`
	Year      int    `json:"year"`
	Month     int    `json:"month"`
	Status    string `json:"status"` // OPEN, PROCESSING, CLOSED
	ClosedBy  string `json:"closed_by,omitempty"`
	ClosedAt  string `json:"closed_at,omitempty"`
	CreatedAt string `json:"created_at"`
}

// CostingResult stores the output of a costing run for one cost object.
type CostingResult struct {
	ID             string          `json:"id"`
	CompanyID      string          `json:"company_id"`
	PeriodID       string          `json:"period_id"`
	CostObjectID   string          `json:"cost_object_id"`
	CostingMethod  CostingMethod   `json:"costing_method"`
	TotalDirectMat float64         `json:"total_direct_materials"`
	TotalDirectLab float64         `json:"total_direct_labor"`
	TotalOverhead  float64         `json:"total_overhead"`
	TotalCost      float64         `json:"total_cost"`
	OutputQuantity float64         `json:"output_quantity"`
	UnitCost       float64         `json:"unit_cost"`
	WIPBegin       float64         `json:"wip_beginning"`
	WIPEnd         float64         `json:"wip_ending"`
	Status        string          `json:"status"` // DRAFT, FINAL, REVERSED
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
}

// CostCategory classifies the type of cost in a costing result line.
type CostCategory string

const (
	CostCategoryDirectMaterial    CostCategory = "DIRECT_MATERIAL"
	CostCategoryDirectLabor       CostCategory = "DIRECT_LABOR"
	CostCategoryOverhead          CostCategory = "OVERHEAD"
	CostCategoryVariance          CostCategory = "VARIANCE"
	CostCategoryProcessInfo       CostCategory = "PROCESS_INFO"
	CostCategoryByProductDeduction CostCategory = "BY_PRODUCT_DEDUCTION"
)

// CostingResultLine is a detailed cost breakdown within a costing result.
type CostingResultLine struct {
	ID              string       `json:"id"`
	ResultID        string       `json:"result_id"`
	CostCategory    CostCategory `json:"cost_category"`
	GLAccountCode   string       `json:"gl_account_code"`
	Description     string       `json:"description"`
	PlannedAmount   float64      `json:"planned_amount"`
	ActualAmount    float64      `json:"actual_amount"`
	AllocatedAmount float64      `json:"allocated_amount"`
	Coefficient     float64      `json:"coefficient,omitempty"`
	CreatedAt       string       `json:"created_at"`
}
