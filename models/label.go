package models

import "github.com/google/uuid"

type Label struct {
	InternalId int64     `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicId   uuid.UUID `json:"public_id" db:"public_id"`
	Name       string    `json:"name" db:"name"`
	Color      string    `json:"color" db:"color"`
}
