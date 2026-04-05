package service

import (
	"context"
	"fmt"

	"github.com/TheTuxis/gondor-files-storage/internal/model"
	"github.com/TheTuxis/gondor-files-storage/internal/repository"
)

type FolderCreateInput struct {
	Name      string
	ParentID  *uint
	CompanyID uint
	CreatedBy uint
}

type FolderUpdateInput struct {
	Name     *string
	ParentID *uint
}

type FolderService struct {
	repo *repository.FolderRepository
}

func NewFolderService(repo *repository.FolderRepository) *FolderService {
	return &FolderService{repo: repo}
}

func (s *FolderService) Create(ctx context.Context, input FolderCreateInput) (*model.Folder, error) {
	folder := &model.Folder{
		Name:      input.Name,
		ParentID:  input.ParentID,
		CompanyID: input.CompanyID,
		CreatedBy: input.CreatedBy,
	}

	if err := s.repo.Create(folder); err != nil {
		return nil, fmt.Errorf("failed to create folder: %w", err)
	}

	return folder, nil
}

func (s *FolderService) GetByID(ctx context.Context, id uint) (*model.Folder, error) {
	folder, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("folder not found: %w", err)
	}
	return folder, nil
}

func (s *FolderService) List(ctx context.Context, companyID uint, parentID *uint) ([]model.Folder, error) {
	folders, err := s.repo.List(companyID, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list folders: %w", err)
	}
	return folders, nil
}

func (s *FolderService) Update(ctx context.Context, id uint, input FolderUpdateInput) (*model.Folder, error) {
	folder, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("folder not found: %w", err)
	}

	if input.Name != nil {
		folder.Name = *input.Name
	}
	if input.ParentID != nil {
		folder.ParentID = input.ParentID
	}

	if err := s.repo.Update(folder); err != nil {
		return nil, fmt.Errorf("failed to update folder: %w", err)
	}

	return folder, nil
}

func (s *FolderService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete folder: %w", err)
	}
	return nil
}
