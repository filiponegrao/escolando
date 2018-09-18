package models

import (
	"time"
)

type Segment struct {
	ID          int64       `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Institution Institution `gorm:"ForeignKey:InstitutionID"`
	Name        string      `gorm:"not null;type:text" json:"name" form:"name"`
	Description string      `gorm:"not null;type:text" json:"description" form:"description"`
	CreatedAt   *time.Time  `json:"createdAt" form:"created_at"`
	UpdatedAt   *time.Time  `json:"updatedAt" form:"updated_at"`
}
