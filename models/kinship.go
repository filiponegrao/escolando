package models

import (
	"time"
)

type Kinship struct {
	ID            int64      `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Name          string     `gorm:"type:text;not null" json:"name" form:"name"`
	AccessContent string     `gorm:"column:access_content" json:"access_content" form:"access_content"`
	CreatedAt     *time.Time `json:"created_at" form:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at" form:"updated_at"`
}
