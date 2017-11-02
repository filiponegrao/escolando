package models

import "time"

type StudentEnrollment struct {
	ID        int64      `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Student   Student    `json:"student" form:"student"`
	Class     Class      `json:"class" form:"class"`
	CreatedAt *time.Time `json:"created_at" form:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" form:"updated_at"`
}
