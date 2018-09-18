package models

import (
	"time"
)

type UserAccessProfile struct {
	ID            int64      `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Name          string     `gorm:"type:text;not null" json:"name" form:"name"`
	AccessContent string     `gorm:"column:access_content" json:"accessContentW" form:"access_content"`
	CreatedAt     *time.Time `json:"createdAt" form:"created_at"`
	UpdatedAt     *time.Time `json:"updatedAt" form:"updated_at"`
}
