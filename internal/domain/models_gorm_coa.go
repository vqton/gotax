package domain

import "time"

type ApprovalRequestGORM struct {
	ID         string     `gorm:"column:id;primaryKey;size:36" json:"id"`
	EntityType string     `gorm:"column:entity_type;not null;size:50;index" json:"entityType"`
	EntityID   string     `gorm:"column:entity_id;not null;size:36;index" json:"entityId"`
	RequestedBy  string   `gorm:"column:requested_by;not null;size:36" json:"requestedBy"`
	ReviewedBy   *string  `gorm:"column:reviewed_by;size:36" json:"reviewedBy"`
	ReviewNote   *string  `gorm:"column:review_note;type:text" json:"reviewNote"`
	Status       string   `gorm:"column:status;not null;size:20;default:'PENDING';index" json:"status"`
	RequestedAt  time.Time `gorm:"column:requested_at;autoCreateTime" json:"requestedAt"`
	ReviewedAt   *time.Time `gorm:"column:reviewed_at" json:"reviewedAt"`
}

func (ApprovalRequestGORM) TableName() string { return "approval_requests" }

type AccountVersionGORM struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	VersionNumber   string    `gorm:"column:version_number;not null;size:50;uniqueIndex" json:"versionNumber"`
	Name            string    `gorm:"column:name;not null;size:255" json:"name"`
	Description     *string   `gorm:"column:description;type:text" json:"description"`
	EffectiveFrom   time.Time `gorm:"column:effective_from;not null;type:date" json:"effectiveFrom"`
	EffectiveTo     *time.Time `gorm:"column:effective_to;type:date" json:"effectiveTo"`
	IsActive        bool      `gorm:"column:is_active;default:true;index" json:"isActive"`
	CreatedBy       string    `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (AccountVersionGORM) TableName() string { return "account_versions" }

type AccountMappingGORM struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SourceRegime   string    `gorm:"column:source_regime;not null;size:50;index:idx_mapping_regime,unique" json:"sourceRegime"`
	TargetRegime   string    `gorm:"column:target_regime;not null;size:50;index:idx_mapping_regime,unique" json:"targetRegime"`
	OldAccountCode string    `gorm:"column:old_account_code;not null;size:20;uniqueIndex" json:"oldAccountCode"`
	NewAccountCode string    `gorm:"column:new_account_code;not null;size:20" json:"newAccountCode"`
	MappingType    string    `gorm:"column:mapping_type;size:50" json:"mappingType"`
	Note           *string   `gorm:"column:note;type:text" json:"note"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (AccountMappingGORM) TableName() string { return "account_mappings" }

type AccountAnalysisGORM struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountCode  string    `gorm:"column:account_code;not null;size:20;uniqueIndex" json:"accountCode"`
	AnalysisType string    `gorm:"column:analysis_type;not null;size:50" json:"analysisType"`
	Data         string    `gorm:"column:data;type:jsonb" json:"data"`
	UpdatedBy    string    `gorm:"column:updated_by;not null;size:36" json:"updatedBy"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (AccountAnalysisGORM) TableName() string { return "account_analysis" }

type IFRSMappingGORM struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	VASCode    string    `gorm:"column:vas_code;not null;size:20;uniqueIndex" json:"vasCode"`
	IFRS       string    `gorm:"column:ifrs_code;not null;size:20" json:"ifrsCode"`
	Note       *string   `gorm:"column:note;type:text" json:"note"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (IFRSMappingGORM) TableName() string { return "ifrs_mappings" }
