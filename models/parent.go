package models

import (
	"time"
)

type Parent struct {
	ID              int64       `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	UserId          int64       `gorm:"column:user_id;not null" json:"userId" form:"userId"`
	Institution     Institution `gorm:"ForeignKey:InstitutionID;AssociationForeignKey:ID;not null;association_autoupdate:false;association_autocreate:false" json:"institution" form:"institution"`
	InstitutionID   int64
	Name            string     `gorm:"not null;type:text" json:"name" form:"name"`
	Email           string     `gorm:"not null;type:text" json:"email" form:"email"`
	Phone           string     `gorm:"type:text" json:"phone" form:"phone"`
	ProfileImageUrl string     `gorm:"column:profile_image_url" json:"profileImageUrl" form:"profileImageUrl"`
	CreatedAt       *time.Time `json:"createdAt" form:"createdAt"`
	UpdatedAt       *time.Time `json:"updatedAt" form:"updatedAt"`
}
