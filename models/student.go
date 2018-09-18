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
	ProfileImageUrl string     `gorm:"column:profile_image_url" json:"profileImageUrl" form:"profile_image_url"`
	ExtraInfo       string     `gorm:"column:extra_info;type:text" json:"extraInfo" form:"extra_info"`
	CreatedAt       *time.Time `json:"createdAt" form:"created_at"`
	UpdatedAt       *time.Time `json:"updatedAt" form:"updated_at"`
}
