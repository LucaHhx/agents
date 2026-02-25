package global

import (
	"gorm.io/gorm"
)

type HZ_MODEL struct {
	ID        uint           `gorm:"primarykey" json:"ID"` // 主键ID
	CreatedAt DateTime       // 创建时间
	UpdatedAt DateTime       // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}

type HZ_MODEL_D struct {
	ID        uint     `gorm:"primarykey" json:"ID"` // 主键ID
	CreatedAt DateTime // 创建时间
	UpdatedAt DateTime // 更新时间
}
