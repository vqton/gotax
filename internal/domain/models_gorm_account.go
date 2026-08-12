package domain

import "time"

type AccountGORM struct {
	Code         string    `gorm:"column:code;primaryKey;size:20" json:"code"`
	Name         string    `gorm:"column:name;not null;size:255" json:"name"`
	Name2        string    `gorm:"column:name2;size:255" json:"name2"`
	Type         string    `gorm:"column:type;not null;size:20" json:"type"`
	ParentCode   string    `gorm:"column:parent_code;size:20;index" json:"parentCode"`
	IsActive     bool      `gorm:"column:is_active;default:true;index" json:"isActive"`
	IsForeign    bool      `gorm:"column:is_foreign;default:false" json:"isForeign"`
	DetailBy     string    `gorm:"column:detail_by;size:50" json:"detailBy"`
	IsParent     bool      `gorm:"column:is_parent;default:false" json:"isParent"`
	Status       string    `gorm:"column:status;size:20" json:"status"`
	FreezeReason string    `gorm:"column:freeze_reason;size:255" json:"freezeReason"`
	ArrearsDays  int       `gorm:"column:arrears_days;default:0" json:"arrearsDays"`
	Note         string    `gorm:"column:note;type:text" json:"note"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (AccountGORM) TableName() string { return "accounts" }
