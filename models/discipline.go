package models

import "time"

type Discipline struct {
	ID          int64      `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Name        string     `json:"name" form:"name"`
	Description string     `json:"description" form:"description"`
	Segment     string     `json:"segment" form:"segment"`
	CreatedAt   *time.Time `json:"created_at" form:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at" form:"updated_at"`
}
