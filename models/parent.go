package models

import (
	"time"
)

type Parent struct {
	ID              int64      `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	UserId          int64      `gorm:"column:user_id;not null" json:"user_id" form:"user_id"`
	Name            string     `gorm:"not null;type:text" json:"name" form:"name"`
	Email           string     `gorm:"not null;type:text" json:"email" form:"email"`
	Phone           string     `gorm:"type:text" json:"phone" form:"phone"`
	ProfileImageUrl string     `gorm:"column:profile_image_url" json:"profileImageUrl" form:"profile_image_url"`
	CreatedAt       *time.Time `json:"createdAt" form:"created_at"`
	UpdatedAt       *time.Time `json:"updatedAt" form:"updated_at"`
}
