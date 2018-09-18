package models

import (
	"time"
)

type StudentEnrollment struct {
	ID        int64   `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Student   Student `gorm:"ForeignKey:StudentID;not null" json:"student" form:"student"`
	StudentID int64
	Class     Class `gorm:"ForeignKey:ClassID;not null" json:"class" form:"class"`
	ClassID   int64
	CreatedAt *time.Time `json:"createdAt" form:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" form:"updated_at"`
}
