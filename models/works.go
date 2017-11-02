package models

import "time"

type WorksFor struct {
	ID          int64       `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	InCharge    InCharge    `json:"user" form:"user"`
	Institution Institution `json:"institution" form:"institution"`
	Status      string      `json:"status" form:"stauts"`
	CreatedAt   *time.Time  `json:"created_at" form:"created_at"`
	UpdatedAt   *time.Time  `json:"updated_at" form:"updated_at"`
}
