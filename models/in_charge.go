package models

import (
	"time"
)

type InCharge struct {
	ID              int64       `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	UserId          int64       `gorm:"column:user_id;not null" json:"user_id" form:"user_id"`
	Institution     Institution `gorm:"ForeignKey:InstitutionID;not null" json:"institution" form:"institution"`
	InstitutionID   int64
	Name            string       `gorm:"not null;type:text" json:"name" form:"name"`
	Email           string       `gorm:"not null;type:text" json:"email" form:"email"`
	Phone           string       `gorm:"type:text" json:"phone" form:"phone"`
	Role            InChargeRole `gorm:"ForeignKey:RoleID;not null" json:"role" form:"role"`
	RoleID          int64
	ProfileImageUrl string     `gorm:"column:profile_image_url" json:"profile_image_url" form:"profile_image_url"`
	CreatedAt       *time.Time `json:"created_at" form:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at" form:"updated_at"`
}
