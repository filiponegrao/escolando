package models

import "time"

type Student struct {
	ID              int64       `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Institution     Institution `json:"institution" form:"institution"`
	Responsible     Parent      `json:"responsible" form:"responsible"`
	Name            string      `json:"name" form:"name"`
	ProfileImageUrl string      `json:"profile_image_url" form:"profile_image_url"`
	ExtraInfo       string      `json:"extra_info" form:"extra_info"`
	CreatedAt       *time.Time  `json:"created_at" form:"created_at"`
	UpdatedAt       *time.Time  `json:"updated_at" form:"updated_at"`
}
