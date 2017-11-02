package models

import "time"

type UserAccessProfile struct {
	ID            int64      `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Name          string     `json:"name" form:"name"`
	AccessContent string     `json:"access_content" form:"access_content"`
	CreatedAt     *time.Time `json:"created_at" form:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at" form:"updated_at"`
}
