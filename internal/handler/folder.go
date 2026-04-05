package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/TheTuxis/gondor-files-storage/internal/service"
	"go.uber.org/zap"
)

type FolderHandler struct {
	folderSvc *service.FolderService
	logger    *zap.Logger
}

func NewFolderHandler(folderSvc *service.FolderService, logger *zap.Logger) *FolderHandler {
	return &FolderHandler{folderSvc: folderSvc, logger: logger}
}

type createFolderRequest struct {
	Name      string `json:"name" binding:"required"`
	ParentID  *uint  `json:"parent_id"`
	CompanyID uint   `json:"company_id" binding:"required"`
}

type updateFolderRequest struct {
	Name     *string `json:"name"`
	ParentID *uint   `json:"parent_id"`
}

func (h *FolderHandler) Create(c *gin.Context) {
	var req createFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	input := service.FolderCreateInput{
		Name:      req.Name,
		ParentID:  req.ParentID,
		CompanyID: req.CompanyID,
		CreatedBy: userID.(uint),
	}

	folder, err := h.folderSvc.Create(c.Request.Context(), input)
	if err != nil {
		h.logger.Error("failed to create folder", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create folder"})
		return
	}

	c.JSON(http.StatusCreated, folder)
}

func (h *FolderHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder id"})
		return
	}

	folder, err := h.folderSvc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
		return
	}

	c.JSON(http.StatusOK, folder)
}

func (h *FolderHandler) List(c *gin.Context) {
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

	var parentID *uint
	if pid := c.Query("parent_id"); pid != "" {
		if v, err := strconv.ParseUint(pid, 10, 32); err == nil {
			u := uint(v)
			parentID = &u
		}
	}

	folders, err := h.folderSvc.List(c.Request.Context(), uint(companyID), parentID)
	if err != nil {
		h.logger.Error("failed to list folders", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list folders"})
		return
	}

	c.JSON(http.StatusOK, folders)
}

func (h *FolderHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder id"})
		return
	}

	var req updateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := service.FolderUpdateInput{
		Name:     req.Name,
		ParentID: req.ParentID,
	}

	folder, err := h.folderSvc.Update(c.Request.Context(), uint(id), input)
	if err != nil {
		h.logger.Error("failed to update folder", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update folder"})
		return
	}

	c.JSON(http.StatusOK, folder)
}

func (h *FolderHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder id"})
		return
	}

	if err := h.folderSvc.Delete(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to delete folder", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete folder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "folder deleted"})
}
