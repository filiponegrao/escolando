package models

import (
	"time"
)

type Student struct {
	ID              int64       `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Institution     Institution `gorm:"ForeignKey:InstitutionID;not null" json:"institution" form:"institution"`
	InstitutionID   int64
	Responsible     Parent `gorm:"ForeignKey:ParentID;not null" json:"responsible" form:"responsible"`
	ParentID        int64
	Name            string     `gorm:"type:text;not null" json:"name" form:"name"`
	ProfileImageUrl string     `gorm:"column:profile_image_url" json:"profile_image_url" form:"profile_image_url"`
	ExtraInfo       string     `gorm:"column:extra_info;type:text" json:"extra_info" form:"extra_info"`
	CreatedAt       *time.Time `json:"created_at" form:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at" form:"updated_at"`
}
