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
	CreatedAt *time.Time `json:"created_at" form:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" form:"updated_at"`
}
