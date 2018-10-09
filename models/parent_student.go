package models

import (
	"time"
)

type ParentStudent struct {
	ID        int64  `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Parent    Parent `gorm:"ForeignKey:ParentID;not null" json:"parent" form:"parent"`
	ParentID  int64
	Student   Student `gorm:"ForeignKey:StudentID;not null" json:"student" form:"student"`
	StudentID int64
	Kinship   Kinship `gorm:"ForeignKey:KinshipID;not null;association_autoupdate:false;association_autocreate:false" json:"kinship" form:"kinship"`
	KinshipID int64
	CreatedAt *time.Time `json:"createdAt" form:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" form:"updated_at"`
}
