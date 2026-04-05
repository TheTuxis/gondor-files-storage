package repository

import (
	"github.com/TheTuxis/gondor-files-storage/internal/model"
	"gorm.io/gorm"
)

type FileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Create(file *model.File) error {
	return r.db.Create(file).Error
}

func (r *FileRepository) GetByID(id uint) (*model.File, error) {
	var file model.File
	err := r.db.Preload("Versions").First(&file, id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *FileRepository) List(companyID uint, folderID *uint, offset, limit int) ([]model.File, int64, error) {
	var files []model.File
	var total int64

	query := r.db.Where("company_id = ?", companyID)
	if folderID != nil {
		query = query.Where("folder_id = ?", *folderID)
	} else {
		query = query.Where("folder_id IS NULL")
	}

	if err := query.Model(&model.File{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

func (r *FileRepository) Delete(id uint) error {
	return r.db.Delete(&model.File{}, id).Error
}

func (r *FileRepository) CreateVersion(version *model.FileVersion) error {
	return r.db.Create(version).Error
}
