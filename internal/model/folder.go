package model

import (
	"time"

	"gorm.io/gorm"
)

type Folder struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	ParentID  *uint          `gorm:"index" json:"parent_id"`
	CompanyID uint           `gorm:"index;not null" json:"company_id"`
	CreatedBy uint           `gorm:"not null" json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Children  []Folder       `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Files     []File         `gorm:"foreignKey:FolderID" json:"files,omitempty"`
}
