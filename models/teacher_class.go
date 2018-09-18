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
	WeekDays     string     `gorm:"type:text;not null" json:"weekDays" form:"week_day"`
	StartHour    int        `gorm:"not null;default:0" json:"startHour" form:"start_hour"`
	StartMinutes int        `gorm:"not null;default:0" json:"startMinutes" form:"start_minutes"`
	EndHour      int        `gorm:"not null;default:0" json:"endHour" form:"end_hour"`
	EndMinutes   int        `gorm:"not null;default:0" json:"endMinutes" form:"end_minutes"`
	CreatedAt    *time.Time `json:"createdAt" form:"created_at"`
	UpdatedAt    *time.Time `json:"updatedAt" form:"updated_at"`
}
