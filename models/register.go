package models

import (
	"time"
)

type RegisterContact struct {
	Name string `json:"name" form:"name"`
	Role string `json:"role" form:"role"`
}

type Register struct {
	ID             int64        `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Title          string       `gorm:"type:text;not null" json:"title" form:"title"`
	Text           string       `gorm:"type:text;not null" json:"text" form:"text"`
	RegisterType   RegisterType `gorm:"ForeignKey:RegisterTypeID;not null" json:"registerType" form:"register_form"`
	RegisterTypeID int64
	Sender         RegisterContact `json:"sender"`
	SenderId       int64           `gorm:"not null" form:"sender_id"`
	Target         RegisterContact `json:"target"`
	TargetId       int64           `gorm:"not null" json:"targetId" form:"target_id"`
	StudentId      int64           `gorm:"not null" json:"studentId" form:"student_id"`
	Status         RegisterStatus  `gorm:"ForeignKey:StatusId;not null" json:"status"`
	StatusId       int64
	CreatedAt      *time.Time `json:"createdAt" form:"created_at"`
	UpdatedAt      *time.Time `json:"updatedAt" form:"updated_at"`
}
