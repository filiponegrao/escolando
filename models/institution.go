package models

import "time"

type Institution struct {
	ID                int64      `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Name              string     `json:"name" form:"name"`
	Email             string     `json:"email" form:"email"`
	Owner             User       `json:"owner" form:"owner"`
	AddressPostal     string     `json:"address_postal" form:"address_postal"`
	AddressStreet     string     `json:"address_street" form:"address_street"`
	AddressNumber     int        `json:"address_number" form:"address_number"`
	AddressComplement string     `json:"address_complement" form:"address_complement"`
	Cnpj              string     `json:"cnpj" form:"cnpj"`
	Phone1            string     `json:"phone1" form:"phone1"`
	Phone2            string     `json:"phone2" form:"phone2"`
	ExtraInfo         string     `json:"extra_info" form:"extra_info"`
	CreatedAt         *time.Time `json:"created_at" form:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at" form:"updated_at"`
}
