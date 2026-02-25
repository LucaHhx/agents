package common

import (
	"server/global"
)

type CrudModel interface {
	TableName() string
	GetModel() *HZ_CRUD
}

type CrudModelDel interface {
	GetDelModel() *HZ_CRUD_DEL
}

type HZ_CRUD struct {
	ID        uint            `gorm:"primarykey" json:"id"` // 主键ID
	CreatedAt global.DateTime `json:"createdAt"`
	UpdatedAt global.DateTime `json:"updatedAt"`
	CreateBy  uint            `json:"createBy"`
	UpdateBy  uint            `json:"updateBy"`
}

type HZ_CRUD_DEL struct {
	DeletedAt global.DateTime `json:"deletedAt"`
	DeletedBy uint            `json:"deletedBy"`
}
