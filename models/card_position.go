package models

import (
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/models/types"
	"github.com/google/uuid"
)

type CardPostion struct {
	InternalId     int64           `json:"internal_id" gorm:"column:primary_key;autoIncrement"`
	PublicId       uuid.UUID       `json:"public_id" gorm:"type:uuid;not null"`
	ListInternalId int64           `json:"list_internal_id" gorm:"column:list_internal_id;not null"`
	CardOrder      types.UUIDArray `json:"card_order" gorm:"type:uuid[]"`
}
