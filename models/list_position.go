package models

import (
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/models/types"
	"github.com/google/uuid"
)

type ListPosition struct {
	InternalID      int64     `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID        uuid.UUID `json:"public_id" db:"public_id" gorm:"column:public_id"`
	BoardInternalID int64     `json:"board_internal_id" db:"board_internal_id" gorm:"column:board_internal_id"`
	//ListOrder     []uuid.UUID 		`json:"list_order" db:"list_order" gorm:"type:json"` //array uuid {uuid1,uud2}
	ListOrder types.UUIDArray `json:"list_order"`
	//gorm sudah ada di types.UUIDArray

}
