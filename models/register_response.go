package models

import (
	"time"
)

type RegisterResponse struct {
	ID         int64    `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Text       string   `gorm:"type:text;not null" json:"text" form:"text"`
	Register   Register `gorm:"ForeignKey:RegisterID;not null" json:"register" form:"register"`
	RegisterID int64
	Sender     RegisterContact `json:"sender"`
	SenderId   int64           `gorm:"not null" form:"sender_id"`
	Status     RegisterStatus  `gorm:"ForeignKey:StatusId;not null" json:"status"`
	StatusId   int64           `json:""`
	CreatedAt  *time.Time      `json:"createdAt" form:"created_at"`
}
