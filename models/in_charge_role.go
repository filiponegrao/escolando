package models

import (
	"time"
)

type InChargeRole struct {
	ID            int64      `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Name          string     `gorm:"not null;unique" json:"name" form:"name"`
	AccessContent string     `gorm:"column:access_content" json:"accessContent" form:"access_content"`
	CreatedAt     *time.Time `json:"createdAt" form:"created_at"`
	UpdatedAt     *time.Time `json:"updatedAt" form:"updated_at"`
}
