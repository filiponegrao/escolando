package models

import (
	"time"
)

type RegisterContact struct {
	Name     string `json:"name" form:"name"`
	Role     string `json:"role" form:"role"`
	ImageURL string `json:"imageUrl" form:"imageUrl"`
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
	TargetId       int64           `gorm:"not null" json:"targetId" form:"targetId"`
	StudentId      int64           `gorm:"not null" json:"studentId" form:"studentId"`
	Status         RegisterStatus  `gorm:"ForeignKey:StatusId;not null" json:"status"`
	StatusId       int64           `json:""`
	GroupTargetId  string          `gorm:"type:text" json:"groupTargetId" form:"groupTargetId"`
	ResponsesCount int64           `gorm:"type:text" json:"responsesCount" form:"responsesCount"`
	InstitutionId  int64
	CreatedAt      *time.Time `json:"createdAt" form:"createdAt"`
	UpdatedAt      *time.Time `json:"updatedAt" form:"updatedAt"`
}
