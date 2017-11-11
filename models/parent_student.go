package models

import (
	"time"
)

type ParentStudent struct {
	ID        int64  `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Parent    Parent `gorm:"ForeignKey:ParentID;not null" json:"user" form:"user"`
	ParentID  int64
	Student   Student `gorm:"ForeignKey:StudentID;not null" json:"student" form:"student"`
	StudentID int64
	Kinship   Kinship `gorm:"ForeignKey:KinshipID;not null" json:"user_access_profile" form:"user_access_profile"`
	KinshipID int64
	CreatedAt *time.Time `json:"created_at" form:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" form:"updated_at"`
}
