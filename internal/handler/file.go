package handler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/TheTuxis/gondor-files-storage/internal/service"
	"go.uber.org/zap"
)

type FileHandler struct {
	fileSvc *service.FileService
	logger  *zap.Logger
}

func NewFileHandler(fileSvc *service.FileService, logger *zap.Logger) *FileHandler {
	return &FileHandler{fileSvc: fileSvc, logger: logger}
}

func (h *FileHandler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	userID, _ := c.Get("user_id")
	companyID, _ := c.Get("company_id")

	var folderID *uint
	if fid := c.PostForm("folder_id"); fid != "" {
		if v, err := strconv.ParseUint(fid, 10, 32); err == nil {
			u := uint(v)
			folderID = &u
		}
	}

	var description *string
	if desc := c.PostForm("description"); desc != "" {
		description = &desc
	}

	input := service.FileUploadInput{
		Name:        header.Filename,
		FolderID:    folderID,
		CompanyID:   companyID.(uint),
		UploadedBy:  userID.(uint),
		MimeType:    header.Header.Get("Content-Type"),
		Size:        header.Size,
		Description: description,
		Reader:      file,
	}

	result, err := h.fileSvc.Upload(c.Request.Context(), input)
	if err != nil {
		h.logger.Error("failed to upload file", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file"})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *FileHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file id"})
		return
	}

	file, err := h.fileSvc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	c.JSON(http.StatusOK, file)
}

func (h *FileHandler) Download(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file id"})
		return
	}

	reader, file, err := h.fileSvc.Download(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to download file", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	defer reader.Close()

	c.Header("Content-Disposition", "attachment; filename=\""+file.Name+"\"")
	c.Header("Content-Type", file.MimeType)
	c.Header("Content-Length", strconv.FormatInt(file.Size, 10))
	c.Status(http.StatusOK)

	if _, err := io.Copy(c.Writer, reader); err != nil {
		h.logger.Error("failed to stream file", zap.Error(err))
	}
}

func (h *FileHandler) List(c *gin.Context) {
	companyIDStr := c.Query("company_id")
	if companyIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}

	companyID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company_id"})
		return
	}

	var folderID *uint
	if fid := c.Query("folder_id"); fid != "" {
		if v, err := strconv.ParseUint(fid, 10, 32); err == nil {
			u := uint(v)
			folderID = &u
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.fileSvc.List(c.Request.Context(), uint(companyID), folderID, page, pageSize)
	if err != nil {
		h.logger.Error("failed to list files", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list files"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *FileHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file id"})
		return
	}

	if err := h.fileSvc.Delete(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to delete file", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file deleted"})
}
