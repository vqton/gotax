package domain

import "time"

type ApprovalStatus string
const (ApprovalPending ApprovalStatus="PENDING"; ApprovalApproved ApprovalStatus="APPROVED"; ApprovalRejected ApprovalStatus="REJECTED"; ApprovalCancelled ApprovalStatus="CANCELLED"; ApprovalExpired ApprovalStatus="EXPIRED")
func (s ApprovalStatus) IsTerminal() bool { return s==ApprovalApproved||s==ApprovalRejected||s==ApprovalCancelled||s==ApprovalExpired }

type ApprovalRequest struct {
	ID string `json:"id"`
	TenantID string `json:"tenant_id,omitempty"`
	EntityType string `json:"entity_type"`
	EntityID string `json:"entity_id"`
	RequestType string `json:"request_type"`
	OldValue string `json:"old_value,omitempty"`
	NewValue string `json:"new_value"`
	Reason string `json:"reason"`
	Status ApprovalStatus `json:"status"`
	RequestedBy string `json:"requested_by"`
	ReviewedBy string `json:"reviewed_by,omitempty"`
	ReviewNote string `json:"review_note,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	ReviewedAt *time.Time `json:"reviewed_at,omitempty"`
}
func (r *ApprovalRequest) Validate() error {
	if r.EntityType=="" { return ErrEntityTypeRequired }
	if r.EntityID=="" { return ErrEntityIDRequired }
	if r.RequestType=="" { return ErrRequestTypeRequired }
	if r.Reason=="" { return ErrApprovalReasonRequired }
	if r.RequestedBy=="" { return ErrRequesterRequired }
	if r.Status=="" { r.Status=ApprovalPending }
	return nil
}

type AccountVersion struct {
	ID string `json:"id"`
	VersionNumber string `json:"version_number"`
	Snapshot string `json:"snapshot"`
	ChangeSummary string `json:"change_summary,omitempty"`
	ChangeReason string `json:"change_reason,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
type AccountMapping struct {
	ID string `json:"id"`
	SourceRegime string `json:"source_regime"`
	TargetRegime string `json:"target_regime"`
	OldCode string `json:"old_code"`
	NewCode string `json:"new_code"`
	MappingType string `json:"mapping_type"`
	SplitRatio float64 `json:"split_ratio,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
func (m *AccountMapping) Validate() error {
	if m.SourceRegime=="" { return ErrSourceRegimeRequired }
	if m.TargetRegime=="" { return ErrTargetRegimeRequired }
	if m.OldCode=="" { return ErrOldCodeRequired }
	if m.NewCode=="" { return ErrNewCodeRequired }
	switch m.MappingType{case "DIRECT","MERGE","SPLIT","CUSTOM":default: return ErrMappingTypeInvalid}
	if m.MappingType=="SPLIT"&&m.SplitRatio<=0 { return ErrSplitRatioRequired }
	return nil
}
type AccountAnalysis struct {
	AccountCode string `json:"account_code"`
	CostCenterID string `json:"cost_center_id,omitempty"`
	ProfitCenterID string `json:"profit_center_id,omitempty"`
	DepartmentID string `json:"department_id,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
	CustomDim1 string `json:"custom_dim1,omitempty"`
	CustomDim2 string `json:"custom_dim2,omitempty"`
}
func (a *AccountAnalysis) Validate() error {
	if a.AccountCode=="" { return ErrAccountCodeRequired }
	return nil
}
type IFRSMapping struct {
	ID string `json:"id"`
	VASCode string `json:"vas_code"`
	IFRSCode string `json:"ifrs_code"`
	IFRSName string `json:"ifrs_name,omitempty"`
	ReclassificationRule string `json:"reclassification_rule,omitempty"`
	AdjustmentType string `json:"adjustment_type,omitempty"`
	IsActive bool `json:"is_active"`
}
func (m *IFRSMapping) Validate() error {
	if m.VASCode=="" { return ErrVASCodeRequired }
	if m.IFRSCode=="" { return ErrIFRSCodeRequired }
	return nil
}
type AccountUsage struct {
	AccountCode string `json:"account_code"`
	EntryCount int `json:"entry_count"`
	TotalDebit float64 `json:"total_debit"`
	TotalCredit float64 `json:"total_credit"`
	LastUsedDate string `json:"last_used_date,omitempty"`
}
type VersionDiff struct {
	VersionFrom string `json:"version_from"`
	VersionTo string `json:"version_to"`
	Added []Account `json:"added,omitempty"`
	Removed []Account `json:"removed,omitempty"`
	Modified []AccountDiff `json:"modified,omitempty"`
}
type AccountDiff struct {
	Code string `json:"code"`
	Old Account `json:"old"`
	New Account `json:"new"`
	Changes map[string]Change `json:"changes"`
}
type Change struct{ OldValue interface{} `json:"old_value"`; NewValue interface{} `json:"new_value"` }
