package models

import "time"

type RegisterCurrentStatus struct {
	ID             int64          `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Register       Register       `gorm:"primary_key;" json:"register" form:"register"`
	RegisterStatus RegisterStatus `gorm:"primary_key;" json:"institution" form:"institution"`
	CreatedAt      *time.Time     `json:"created_at" form:"created_at"`
}
