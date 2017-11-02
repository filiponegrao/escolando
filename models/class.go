package models

import "time"

type Class struct {
	ID          int64       `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Institution Institution `json:"institution" form:"institution"`
	InCharge    InCharge    `json:"in_charge" form:"in_charge"`
	Name        string      `json:"name" form:"name"`
	Capacity    int         `json:"capacity" form:"capacity"`
	Enrolled    int         `json:"enrolled" form:"enrolled"`
	SchoolGrade SchoolGrade `json:"school_grade" form:"school_grade"`
	CreatedAt   *time.Time  `json:"created_at" form:"created_at"`
	UpdatedAt   *time.Time  `json:"updated_at" form:"updated_at"`
}
