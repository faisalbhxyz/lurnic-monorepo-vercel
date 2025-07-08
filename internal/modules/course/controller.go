package course

import (
	"context"
	"dashlearn/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type CourseHandler struct {
	service CourseService
}

func NewCourseHandler(db *gorm.DB) *CourseHandler {
	return &CourseHandler{
		service: NewCourseService(db),
	}
}

func (h *CourseHandler) GetAll(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	courses, err := h.service.GetAll(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": courses})
}

func (h *CourseHandler) GetAllLite(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	courses, err := h.service.GetAllLite(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": courses})
}

func (h *CourseHandler) GetAllPublic(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	courses, err := h.service.GetAllPublic(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": courses})
}

func (h *CourseHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	tenantID := c.GetUint("tenant_id")
	course, err := h.service.GetByID(tenantID, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": course})
}

func (h *CourseHandler) Create(c *gin.Context) {
	var input CourseDetailsInput
	var flatInput CreateCourseDetailsInput

	// Step 1: Bind all flat fields (this ignores nested JSON fields like course_chapters)
	if err := c.ShouldBindWith(&flatInput, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_ = copier.Copy(&input, &flatInput)

	// Step 2: Manually parse nested JSON fields from string values
	if chaptersStr := c.PostForm("course_chapters"); chaptersStr != "" {
		var courseChapters []CreateCourseChapter
		if err := json.Unmarshal([]byte(chaptersStr), &courseChapters); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course_chapters: " + err.Error()})
			return
		}
		input.CourseChapters = courseChapters
	}

	if generalSettingsStr := c.PostForm("general_settings"); generalSettingsStr != "" {
		var generalSettings CreateGeneralSettings
		if err := json.Unmarshal([]byte(generalSettingsStr), &generalSettings); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid general_settings: " + err.Error()})
			return
		}
		input.GeneralSettings = generalSettings
	}

	if instructorsStr := c.PostForm("course_instructors"); instructorsStr != "" {
		var instructors []int32
		if err := json.Unmarshal([]byte(instructorsStr), &instructors); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructors: " + err.Error()})
			return
		}
		input.Instructors = instructors
	}

	file, err := c.FormFile("featured_image")
	if err == nil {
		url, err := utils.UploadFile(context.Background(), file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		input.FeaturedImage = &url
	} else {
		input.FeaturedImage = nil
	}

	// if output, err := json.MarshalIndent(input, "", "  "); err == nil {
	// 	fmt.Println("Parsed Input:\n", string(output))
	// }

	// Step 3: Pass the parsed object to the service layer for further processing
	if err := h.service.Create(input, c.GetUint("tenant_id"), c.GetUint("user_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Course created successfully"})
}

func (h *CourseHandler) Update(c *gin.Context) {

	var input CourseDetailsInput
	var flatInput CreateCourseDetailsInput

	courseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	if err := c.ShouldBindWith(&flatInput, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_ = copier.Copy(&input, &flatInput)

	if chaptersStr := c.PostForm("course_chapters"); chaptersStr != "" {
		var chapters []CreateCourseChapter
		if err := json.Unmarshal([]byte(chaptersStr), &chapters); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course_chapters: " + err.Error()})
			return
		}
		input.CourseChapters = chapters
	}

	if settingsStr := c.PostForm("general_settings"); settingsStr != "" {
		var settings CreateGeneralSettings
		if err := json.Unmarshal([]byte(settingsStr), &settings); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid general_settings: " + err.Error()})
			return
		}
		input.GeneralSettings = settings
	}

	if instructorsStr := c.PostForm("course_instructors"); instructorsStr != "" {
		var instructors []int32
		if err := json.Unmarshal([]byte(instructorsStr), &instructors); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructors: " + err.Error()})
			return
		}
		input.Instructors = instructors
	}

	file, err := c.FormFile("featured_image")
	if err == nil {
		url, err := utils.UploadFile(context.Background(), file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		input.FeaturedImage = &url
	}

	// if output, err := json.MarshalIndent(input, "", "  "); err == nil {
	// 	fmt.Println("Parsed Input:\n", string(output))
	// }

	if err := h.service.Update(uint(courseID), c.GetUint("tenant_id"), c.GetUint("user_id"), input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course updated successfully"})
}

func (h *CourseHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	if err := h.service.Delete(uint(id), c.GetUint("tenant_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}
