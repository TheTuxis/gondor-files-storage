package model

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	FolderID    *uint          `gorm:"index" json:"folder_id"`
	CompanyID   uint           `gorm:"index;not null" json:"company_id"`
	UploadedBy  uint           `gorm:"not null" json:"uploaded_by"`
	MimeType    string         `json:"mime_type"`
	Size        int64          `json:"size"`
	S3Key       string         `gorm:"not null" json:"-"`
	Description *string        `json:"description"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Versions    []FileVersion  `gorm:"foreignKey:FileID" json:"versions,omitempty"`
}
