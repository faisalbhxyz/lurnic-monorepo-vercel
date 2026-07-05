package academicnote

import (
	"dashlearn/internal/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type Handler struct {
	service Service
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{service: NewService(db)}
}

func (h *Handler) GetAllClasses(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	classes, err := h.service.GetAllClasses(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": classes})
}

func (h *Handler) GetClassByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID"})
		return
	}
	class, err := h.service.GetClassByID(c.GetUint("tenant_id"), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": class})
}

func (h *Handler) CreateClass(c *gin.Context) {
	var input CreateClassInput
	if err := c.ShouldBindWith(&input, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	if iconURL, err := uploadFile(c, "icon_image", 2*1024*1024); err == nil {
		input.IconImage = iconURL
	}
	id, err := h.service.CreateClass(input, c.GetUint("tenant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Class created successfully",
		"data":    gin.H{"id": id},
	})
}

func (h *Handler) UpdateClass(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID"})
		return
	}
	var input UpdateClassInput
	if err := c.ShouldBindWith(&input, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	if iconURL, err := uploadFile(c, "icon_image", 2*1024*1024); err == nil {
		input.IconImage = iconURL
	}
	if err := h.service.UpdateClass(id, input, c.GetUint("tenant_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Class updated successfully"})
}

func (h *Handler) DeleteClass(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID"})
		return
	}
	if err := h.service.DeleteClass(id, c.GetUint("tenant_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Class deleted successfully"})
}

func (h *Handler) CreateSubject(c *gin.Context) {
	var input CreateSubjectInput
	if err := c.ShouldBindWith(&input, binding.Form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	if err := h.service.CreateSubject(input, c.GetUint("tenant_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Subject created successfully"})
}

func (h *Handler) UpdateSubject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subject ID"})
		return
	}
	var input UpdateSubjectInput
	if err := c.ShouldBindWith(&input, binding.Form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	if err := h.service.UpdateSubject(id, input, c.GetUint("tenant_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Subject updated successfully"})
}

func (h *Handler) DeleteSubject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subject ID"})
		return
	}
	if err := h.service.DeleteSubject(id, c.GetUint("tenant_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Subject deleted successfully"})
}

func (h *Handler) CreatePaper(c *gin.Context) {
	var input CreatePaperInput
	if err := c.ShouldBindWith(&input, binding.Form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	if err := h.service.CreatePaper(input, c.GetUint("tenant_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Paper created successfully"})
}

func (h *Handler) UpdatePaper(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid paper ID"})
		return
	}
	var input UpdatePaperInput
	if err := c.ShouldBindWith(&input, binding.Form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	if err := h.service.UpdatePaper(id, input, c.GetUint("tenant_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Paper updated successfully"})
}

func (h *Handler) DeletePaper(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid paper ID"})
		return
	}
	if err := h.service.DeletePaper(id, c.GetUint("tenant_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Paper deleted successfully"})
}

func uploadFile(c *gin.Context, field string, maxSize int64) (*string, error) {
	header, err := c.FormFile(field)
	if err != nil {
		return nil, err
	}
	if header.Size > maxSize {
		return nil, fmt.Errorf("%s too large", field)
	}
	file, err := header.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	url, err := utils.UploadToBunny(file, header)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (h *Handler) CreateNote(c *gin.Context) {
	var input CreateNoteInput
	if err := c.ShouldBindWith(&input, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	if thumbURL, err := uploadFile(c, "thumbnail", 2*1024*1024); err == nil {
		input.Thumbnail = thumbURL
	}

	pdfHeader, err := c.FormFile("pdf")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PDF file is required"})
		return
	}
	if pdfHeader.Size > 20*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PDF too large (max 20MB)"})
		return
	}
	pdfFile, err := pdfHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open PDF"})
		return
	}
	defer pdfFile.Close()
	pdfURL, err := utils.UploadToBunny(pdfFile, pdfHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	input.PdfURL = pdfURL
	input.PdfFileName = &pdfHeader.Filename

	if err := h.service.CreateNote(input, c.GetUint("tenant_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Note created successfully"})
}

func (h *Handler) UpdateNote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}
	var input UpdateNoteInput
	if err := c.ShouldBindWith(&input, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	if thumbURL, err := uploadFile(c, "thumbnail", 2*1024*1024); err == nil {
		input.Thumbnail = thumbURL
	}

	if pdfHeader, err := c.FormFile("pdf"); err == nil {
		if pdfHeader.Size > 20*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "PDF too large (max 20MB)"})
			return
		}
		pdfFile, err := pdfHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open PDF"})
			return
		}
		defer pdfFile.Close()
		pdfURL, err := utils.UploadToBunny(pdfFile, pdfHeader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		input.PdfURL = &pdfURL
		input.PdfFileName = &pdfHeader.Filename
	}

	if err := h.service.UpdateNote(id, input, c.GetUint("tenant_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Note updated successfully"})
}

func (h *Handler) DeleteNote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}
	if err := h.service.DeleteNote(id, c.GetUint("tenant_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}

func (h *Handler) GetPublicClasses(c *gin.Context) {
	classes, err := h.service.GetPublicClasses(c.GetUint("tenant_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": classes})
}

func (h *Handler) GetPublicClassBySlug(c *gin.Context) {
	class, err := h.service.GetPublicClassBySlug(c.GetUint("tenant_id"), c.Param("classSlug"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": class})
}

func (h *Handler) GetPublicNotesByPaperSlug(c *gin.Context) {
	data, err := h.service.GetPublicNotesByPaperSlug(
		c.GetUint("tenant_id"),
		c.Param("classSlug"),
		c.Param("subjectSlug"),
		c.Param("paperSlug"),
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}
