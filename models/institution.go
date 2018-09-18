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
	AddressPostal     string     `gorm:"column:address_postal;type:text" json:"addressPostal" form:"address_postal"`
	AddressStreet     string     `gorm:"column:address_street;type:text" json:"addressStreet" form:"address_street"`
	AddressNumber     int        `gorm:"column:address_number;type:text" json:"addressNumber" form:"address_number"`
	AddressComplement string     `gorm:"column:address_complement;type:text" json:"addressComplement" form:"address_complement"`
	Cnpj              string     `gorm:"type:text" json:"cnpj" form:"cnpj"`
	Phone1            string     `gorm:"column:phone_1;type:text" json:"phone1" form:"phone1"`
	Phone2            string     `gorm:"column:phone_2;type:text" json:"phone2" form:"phone2"`
	ExtraInfo         string     `gorm:"column:extra_info;type:text" json:"extraInfo" form:"extra_info"`
	CreatedAt         *time.Time `json:"createdAt" form:"created_at"`
	UpdatedAt         *time.Time `json:"updatedAt" form:"updated_at"`
}
