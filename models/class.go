package models

import (
	"time"
)

type Class struct {
	ID            int64       `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Institution   Institution `gorm:"ForeignKey:InstitutionID;not null" json:"institution" form:"institution"`
	InstitutionID int64
	InCharge      InCharge `gorm:"ForeignKey:InChargeID;not null" json:"in_charge" form:"in_charge"`
	InChargeID    int64
	Name          string      `gorm:"type:text;not null" json:"name" form:"name"`
	Capacity      int         `gorm:"not null;default:0" json:"capacity" form:"capacity"`
	Enrolled      int         `gorm:"not null;default:0" json:"enrolled" form:"enrolled"`
	SchoolGrade   SchoolGrade `gorm:"ForeignKey:SchoolGradeID;not null" json:"school_grade" form:"school_grade"`
	SchoolGradeID int64
	CreatedAt     *time.Time `json:"created_at" form:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at" form:"updated_at"`
}
