package student

import (
	"dashlearn/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func parseProfileImageUpload(c *gin.Context) (*string, error) {
	fileHeader, err := c.FormFile("profile_image")
	if err != nil {
		return nil, nil
	}

	if fileHeader.Size > 2*1024*1024 {
		return nil, fmt.Errorf("max image size is 2MB")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open image file")
	}
	defer file.Close()

	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return nil, fmt.Errorf("failed to read file content")
	}
	contentType := http.DetectContentType(buffer)

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/jpg":  true,
	}
	if !allowedTypes[contentType] {
		return nil, fmt.Errorf("only PNG, JPG formats are supported")
	}

	url, err := utils.UploadToBunny(file, fileHeader)
	if err != nil {
		return nil, err
	}

	return &url, nil
}
