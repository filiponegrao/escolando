package models

import (
	"time"
)

type User struct {
	ID                int64      `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Name              string     `gorm:"type:text;not null" json:"name" form:"name"`
	Email             string     `gorm:"type:text;not null unique" json:"email" form:"email"`
	Password          string     `gorm:"type:text;not null" json:"password" form:"password"`
	FacebookId        string     `gorm:"column:facebook_id;type:text" json:"facebookId" form:"facebook_id"`
	AddressPostal     string     `gorm:"column:address_postal;type:text" json:"addressPostal" form:"address_postal"`
	AddressStreet     string     `gorm:"column:address_street;type:text" json:"addressStreet" form:"address_street"`
	AddressNumber     int        `gorm:"column:address_number;type:text" json:"addressNumber" form:"address_number"`
	AddressComplement string     `gorm:"column:address_complement;type:text" json:"addressComplement" form:"address_complement"`
	Cpf               string     `gorm:"type:text" json:"cpf" form:"cpf"`
	Rg                string     `gorm:"type:text" json:"rg" form:"rg"`
	Phone1            string     `gorm:"column:phone_1;type:text" json:"phone1" form:"phone1"`
	Phone2            string     `gorm:"column:phone_2;type:text" json:"phone2" form:"phone2"`
	ProfileImageUrl   string     `gorm:"column:profile_image_url;type:text" json:"profileImageUrl" form:"profile_image_url"`
	ExtraInfo         string     `gorm:"column:extra_info;type:text" json:"extraInfo" form:"extra_info"`
	CreatedAt         *time.Time `json:"createdAt" form:"created_at"`
	UpdatedAt         *time.Time `json:"updatedAt" form:"updated_at"`
}
