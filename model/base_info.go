package model

import "time"

type BaseInfo struct {
	ID        int64     `gorm:"column:id" json:"id" `                // 自增id
	IsDelete  int       `gorm:"column:is_delete" json:"isDelete" `   // 软删除标记
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt" ` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt" ` // 更新时间
	CreatedBy string    `gorm:"column:created_by" json:"createdBy" ` // 创建人
	UpdatedBy string    `gorm:"column:updated_by" json:"updatedBy"`  // 修改人
}
