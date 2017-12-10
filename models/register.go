package models

import (
	"time"
)

type Register struct {
	ID             int64        `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Title          string       `gorm:"type:text;not null" json:"title" form:"title"`
	Text           string       `gorm:"type:text;not null" json:"text" form:"text"`
	RegisterType   RegisterType `gorm:"ForeignKey:RegisterTypeID;not null" json:"register_type"`
	RegisterTypeID int64
	SenderId       int64          `gorm:"not null" json:"sender_id" form:"sender_id"`
	TargetId       int64          `gorm:"not null" json:"target_id" form:"target_id"`
	StudentId      int64          `gorm:"not null" json:"student_id" form:"student_id"`
	Status         RegisterStatus `gorm:"ForeignKey:StatusId;not null" json:"status"`
	StatusId       int64
	CreatedAt      *time.Time `json:"created_at" form:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at" form:"updated_at"`
}
