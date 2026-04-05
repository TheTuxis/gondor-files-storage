package model

import "time"

type FileShare struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	FileID     uint      `gorm:"index;not null" json:"file_id"`
	SharedWith uint      `gorm:"index;not null" json:"shared_with"`
	Permission string    `gorm:"not null;default:read" json:"permission"`
	SharedBy   uint      `gorm:"not null" json:"shared_by"`
	CreatedAt  time.Time `json:"created_at"`
}
