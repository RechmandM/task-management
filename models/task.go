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


// CREATE TABLE tasks (
//     id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
//     title VARCHAR(255) NOT NULL,
//     description TEXT,
//     status VARCHAR(50) NOT NULL DEFAULT 'pending',
//     assignee VARCHAR(100),

//     created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
//     updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
//     deleted_at TIMESTAMP NULL DEFAULT NULL,

//     UNIQUE KEY uk_tasks_title (title),

//     INDEX idx_status (status),
//     INDEX idx_assignee (assignee),
//     INDEX idx_deleted_at (deleted_at),
//     INDEX idx_created_at (created_at)
// );