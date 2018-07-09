package models

import (
	"time"
)

type SchoolGrade struct {
	ID            int64       `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Name          string      `gorm:"type:text;not null" json:"name" form:"name"`
	Institution   Institution `gorm:"ForeignKey:InstitutionID;AssociationForeignKey:ID;not null" json:"institution" form:"institution"`
	InstitutionID int64
	CreatedAt     *time.Time `json:"created_at" form:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at" form:"updated_at"`
}
