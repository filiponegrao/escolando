package models

import "time"

type InCharge struct {
	ID              int64        `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	UserId          int64        `json:"user_id" form:"user_id"`
	Institution     Institution  `json:"institution" form:"institution"`
	Name            string       `json:"name" form:"name"`
	Email           string       `json:"email" form:"email"`
	Phone           string       `json:"phone" form:"phone"`
	Role            InChargeRole `json:"role" form:"role"`
	ProfileImageUrl string       `json:"profile_image_url" form:"profile_image_url"`
	CreatedAt       *time.Time   `json:"created_at" form:"created_at"`
	UpdatedAt       *time.Time   `json:"updated_at" form:"updated_at"`
}
