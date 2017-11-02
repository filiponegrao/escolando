package models

import "time"

type UserAccess struct {
	ID                int64             `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	User              User              `json:"user" form:"user"`
	Institution       Institution       `json:"institution" form:"institution"`
	UserAccessProfile UserAccessProfile `json:"user_access_profile" form:"user_access_profile"`
	CreatedAt         *time.Time        `json:"created_at" form:"created_at"`
	UpdatedAt         *time.Time        `json:"updated_at" form:"updated_at"`
}
