package models

import (
	"time"
)

type TeacherClass struct {
	ID           int64    `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Teacher      InCharge `gorm:"ForeignKey:InChargeID;not null" json:"teacher" form:"teacher"`
	InChargeID   int64
	Discipline   Discipline `gorm:"ForeignKey:DisciplineID" json:"discipline" form:"discipline"`
	DisciplineID int64
	Class        Class `gorm:"ForeignKey:ClassID" json:"class" form:"class"`
	ClassID      int64
	WeekDay      string     `gorm:"type:text;not null" json:"week_day" form:"week_day"`
	StartHour    int        `gorm:"not null;default:0" json:"start_hour" form:"start_hour"`
	StartMinutes int        `gorm:"not null;default:0" json:"start_minutes" form:"start_minutes"`
	EndHour      int        `gorm:"not null;default:0" json:"end_hour" form:"end_hour"`
	EndMinutes   int        `gorm:"not null;default:0" json:"en_minutes" form:"end_minutes"`
	CreatedAt    *time.Time `json:"created_at" form:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at" form:"updated_at"`
}
