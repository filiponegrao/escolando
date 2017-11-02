package models

import "time"

type TeacherClass struct {
	ID           int64      `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Teacher      InCharge   `json:"teacher" form:"teacher"`
	Discipline   Discipline `json:"discipline" form:"discipline"`
	Class        Class      `json:"class" form:"class"`
	WeekDay      string     `json:"week_day" form:"week_day"`
	StartHour    int        `json:"start_hour" form:"start_hour"`
	StartMinutes int        `json:"start_minutes" form:"start_minutes"`
	EndHour      int        `json:"end_hour" form:"end_hour"`
	EndMinutes   int        `json:"en_minutes" form:"end_minutes"`
	CreatedAt    *time.Time `json:"created_at" form:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at" form:"updated_at"`
}
