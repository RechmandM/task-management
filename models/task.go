package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Title       string `gorm:"type:varchar(255);uniqueIndex;not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	Status      string `gorm:"type:varchar(50);default:pending" json:"status"`
	Assignee    string `gorm:"type:varchar(100)" json:"assignee"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}