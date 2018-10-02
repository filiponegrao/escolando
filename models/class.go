package models

import (
	"time"
)

type Class struct {
	ID            int64       `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Institution   Institution `gorm:"ForeignKey:InstitutionID;not null" json:"institution"`
	InstitutionId int64
	InCharge      InCharge `gorm:"ForeignKey:InChargeID;not null" json:"inCharge" form:"in_charge"`
	InChargeID    int64
	Name          string      `gorm:"type:text;not null" json:"name" form:"name"`
	Capacity      int         `gorm:"not null;default:0" json:"capacity" form:"capacity"`
	Enrolled      int         `gorm:"not null;default:0" json:"enrolled" form:"enrolled"`
	SchoolGrade   SchoolGrade `gorm:"ForeignKey:SchoolGradeID;not null" json:"schoolGrade" form:"schoolGrade"`
	SchoolGradeID int64
	CreatedAt     *time.Time `json:"createdAt" form:"created_at"`
	UpdatedAt     *time.Time `json:"updatedAt" form:"updated_at"`
}
