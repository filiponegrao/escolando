package models

import (
	"time"
)

type Institution struct {
	ID                int64  `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Name              string `gorm:"not null;type:text" json:"name" form:"name"`
	Email             string `gorm:"not null;type:text" json:"email" form:"email"`
	Owner             User   `gorm:"ForeignKey:UserID;AssociationForeignKey:ID;not null" json:"owner" form:"owner"`
	UserID            int64
	AddressPostal     string     `gorm:"column:address_postal;type:text" json:"address_postal" form:"address_postal"`
	AddressStreet     string     `gorm:"column:address_street;type:text" json:"address_street" form:"address_street"`
	AddressNumber     int        `gorm:"column:address_number;type:text" json:"address_number" form:"address_number"`
	AddressComplement string     `gorm:"column:address_complement;type:text" json:"address_complement" form:"address_complement"`
	Cnpj              string     `gorm:"type:text" json:"cnpj" form:"cnpj"`
	Phone1            string     `gorm:"column:phone_1;type:text" json:"phone1" form:"phone1"`
	Phone2            string     `gorm:"column:phone_2;type:text" json:"phone2" form:"phone2"`
	ExtraInfo         string     `gorm:"column:extra_info;type:text" json:"extra_info" form:"extra_info"`
	CreatedAt         *time.Time `json:"created_at" form:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at" form:"updated_at"`
}
