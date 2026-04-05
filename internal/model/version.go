package model

import "time"

type FileVersion struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	FileID     uint      `gorm:"index;not null" json:"file_id"`
	Version    int       `gorm:"not null" json:"version"`
	S3Key      string    `gorm:"not null" json:"-"`
	Size       int64     `json:"size"`
	UploadedBy uint      `json:"uploaded_by"`
	CreatedAt  time.Time `json:"created_at"`
}
