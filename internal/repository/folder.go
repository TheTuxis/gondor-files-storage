package repository

import (
	"github.com/TheTuxis/gondor-files-storage/internal/model"
	"gorm.io/gorm"
)

type FolderRepository struct {
	db *gorm.DB
}

func NewFolderRepository(db *gorm.DB) *FolderRepository {
	return &FolderRepository{db: db}
}

func (r *FolderRepository) Create(folder *model.Folder) error {
	return r.db.Create(folder).Error
}

func (r *FolderRepository) GetByID(id uint) (*model.Folder, error) {
	var folder model.Folder
	err := r.db.Preload("Children").Preload("Files").First(&folder, id).Error
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *FolderRepository) List(companyID uint, parentID *uint) ([]model.Folder, error) {
	var folders []model.Folder
	query := r.db.Where("company_id = ?", companyID)
	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}
	if err := query.Order("name ASC").Find(&folders).Error; err != nil {
		return nil, err
	}
	return folders, nil
}

func (r *FolderRepository) Update(folder *model.Folder) error {
	return r.db.Save(folder).Error
}

func (r *FolderRepository) Delete(id uint) error {
	return r.db.Delete(&model.Folder{}, id).Error
}
