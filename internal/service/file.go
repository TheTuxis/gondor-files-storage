package service

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/TheTuxis/gondor-files-storage/internal/model"
	"github.com/TheTuxis/gondor-files-storage/internal/repository"
)

type FileUploadInput struct {
	Name        string
	FolderID    *uint
	CompanyID   uint
	UploadedBy  uint
	MimeType    string
	Size        int64
	Description *string
	Reader      io.Reader
}

type PaginatedResult struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

type FileService struct {
	repo    *repository.FileRepository
	storage *StorageService
}

func NewFileService(repo *repository.FileRepository, storage *StorageService) *FileService {
	return &FileService{repo: repo, storage: storage}
}

func (s *FileService) Upload(ctx context.Context, input FileUploadInput) (*model.File, error) {
	s3Key := fmt.Sprintf("%d/%s/%s", input.CompanyID, uuid.New().String(), input.Name)

	if err := s.storage.Upload(ctx, s3Key, input.Reader, input.MimeType); err != nil {
		return nil, fmt.Errorf("failed to upload to storage: %w", err)
	}

	file := &model.File{
		Name:        input.Name,
		FolderID:    input.FolderID,
		CompanyID:   input.CompanyID,
		UploadedBy:  input.UploadedBy,
		MimeType:    input.MimeType,
		Size:        input.Size,
		S3Key:       s3Key,
		Description: input.Description,
		IsActive:    true,
	}

	if err := s.repo.Create(file); err != nil {
		return nil, fmt.Errorf("failed to create file record: %w", err)
	}

	return file, nil
}

func (s *FileService) Download(ctx context.Context, fileID uint) (io.ReadCloser, *model.File, error) {
	file, err := s.repo.GetByID(fileID)
	if err != nil {
		return nil, nil, fmt.Errorf("file not found: %w", err)
	}

	reader, err := s.storage.Download(ctx, file.S3Key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to download from storage: %w", err)
	}

	return reader, file, nil
}

func (s *FileService) GetByID(ctx context.Context, id uint) (*model.File, error) {
	file, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}
	return file, nil
}

func (s *FileService) List(ctx context.Context, companyID uint, folderID *uint, page, pageSize int) (*PaginatedResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	files, total, err := s.repo.List(companyID, folderID, offset, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PaginatedResult{
		Items:      files,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *FileService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
