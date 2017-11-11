package models

import (
	"time"
)

type UserAccess struct {
	ID                  int64 `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	User                User  `gorm:"ForeignKey:UserID;AssociationForeignKey:ID;not null" json:"user" form:"user"`
	UserID              int64
	Institution         Institution `gorm:"ForeignKey:InstitutionID;AssociationForeignKey:ID;not null" json:"institution" form:"institution"`
	InstitutionID       int64
	UserAccessProfile   UserAccessProfile `gorm:"ForeignKey:UserAccessProfileID;AssociationForeignKey:ID;not null" json:"user_access_profile" form:"user_access_profile"`
	UserAccessProfileID int64
	CreatedAt           *time.Time `json:"created_at" form:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at" form:"updated_at"`
}
