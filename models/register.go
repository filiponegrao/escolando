package models

import "time"

type Register struct {
	ID           int64        `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Title        string       `json:"title" form:"title"`
	Text         string       `json:"text" form:"text"`
	RegisterType RegisterType `json:"register_type"`
	SenderId     int64        `json:"sender_id" form:"sender_id"`
	TargetId     int64        `json:"target_id" form:"target_id"`
	Student      Student      `json:"student" form:"student"`
	CreatedAt    *time.Time   `json:"created_at" form:"created_at"`
	UpdatedAt    *time.Time   `json:"updated_at" form:"updated_at"`
}
