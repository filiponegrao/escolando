package models

import (
	"time"
)

type Discipline struct {
	ID          int64      `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Name        string     `gorm:"not null;type:text" json:"name" form:"name"`
	Description string     `gorm:"not null;type:text" json:"description" form:"description"`
	Segment     string     `json:"segment" form:"segment"`
	CreatedAt   *time.Time `json:"createdAt" form:"created_at"`
	UpdatedAt   *time.Time `json:"updatedAt" form:"updated_at"`
}
