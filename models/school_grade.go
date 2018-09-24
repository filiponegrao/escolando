package models

import (
	"time"
)

type SchoolGrade struct {
	ID        int64   `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Name      string  `gorm:"type:text;not null" json:"name" form:"name"`
	Segment   Segment `gorm:"ForeignKey:SegmentId;AssociationForeignKey:ID;not null" json:"segment" form:"segment"`
	SegmentId int64
	CreatedAt *time.Time `json:"createdAt" form:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" form:"updated_at"`
}
