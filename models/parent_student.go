package models

import "time"

type ParentStudent struct {
	ID        int64      `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Parent    Parent     `json:"user" form:"user"`
	Student   Student    `json:"student" form:"student"`
	Kinship   Kinship    `json:"user_access_profile" form:"user_access_profile"`
	CreatedAt *time.Time `json:"created_at" form:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" form:"updated_at"`
}
