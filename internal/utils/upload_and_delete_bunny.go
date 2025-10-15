package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gosimple/slug"
)

// UploadToBunny uploads a file to BunnyCDN storage and returns its public URL.
func UploadToBunny(file multipart.File, header *multipart.FileHeader) (string, error) {
	bunnyAPIKey := os.Getenv("BUNNY_API_KEY")
	bunnyZone := os.Getenv("BUNNY_STORAGE_ZONE")
	bunnyHost := os.Getenv("BUNNY_STORAGE_HOSTNAME")
	bunnyPullZone := os.Getenv("BUNNY_PULL_ZONE")
	bunnyFolder := os.Getenv("BUNNY_UPLOAD_FOLDER")

	if bunnyAPIKey == "" || bunnyZone == "" || bunnyHost == "" || bunnyPullZone == "" || bunnyFolder == "" {
		return "", fmt.Errorf("missing BunnyCDN configuration")
	}

	// Clean and slugify filename
	re := regexp.MustCompile(`[*#+~()'"!:@]`)
	cleanedName := re.ReplaceAllString(header.Filename, "")
	ext := filepath.Ext(cleanedName)
	nameWithoutExt := strings.TrimSuffix(cleanedName, ext)
	slugifiedName := slug.Make(nameWithoutExt)
	safeName := fmt.Sprintf("%d_%s%s", time.Now().Unix(), slugifiedName, ext)
	fileKey := fmt.Sprintf("%s/%s", bunnyFolder, safeName)

	// Storage URL
	storageURL := fmt.Sprintf("https://%s/%s/%s", bunnyHost, bunnyZone, fileKey)

	// Reset file pointer just in case
	if seeker, ok := file.(io.Seeker); ok {
		seeker.Seek(0, 0)
	}

	// Detect MIME type (first 512 bytes)
	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	mimeType := http.DetectContentType(buf[:n])

	// Reset pointer again before sending
	if seeker, ok := file.(io.Seeker); ok {
		seeker.Seek(0, 0)
	}

	req, err := http.NewRequest("PUT", storageURL, file)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("AccessKey", bunnyAPIKey)
	req.Header.Set("Content-Type", mimeType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("upload failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("bunny upload failed: %d %s", resp.StatusCode, string(body))
	}

	publicURL := fmt.Sprintf("https://%s/%s", bunnyPullZone, fileKey)
	return publicURL, nil
}

// getFileKeyFromURL extracts the file key (path inside the storage zone) from a full CDN URL
func getFileKeyFromURL(cdnURL string) (string, error) {
	parsedURL, err := url.Parse(cdnURL)
	if err != nil {
		return "", fmt.Errorf("invalid CDN URL: %v", err)
	}
	// Remove leading slashes
	fileKey := strings.TrimLeft(parsedURL.Path, "/")
	if fileKey == "" {
		return "", fmt.Errorf("empty file key in URL")
	}
	return fileKey, nil
}

func DeleteFromBunny(fileURL string) error {
	fileKey, err := getFileKeyFromURL(fileURL)
	if err != nil {
		return fmt.Errorf("invalid CDN URL: %v", err)
	}

	bunnyAPIKey := os.Getenv("BUNNY_API_KEY")
	bunnyZone := os.Getenv("BUNNY_STORAGE_ZONE")
	bunnyHost := os.Getenv("BUNNY_STORAGE_HOSTNAME")

	if bunnyAPIKey == "" || bunnyZone == "" || bunnyHost == "" {
		return fmt.Errorf("missing BunnyCDN environment configuration")
	}

	// Build full Bunny Storage URL
	url := fmt.Sprintf("https://%s/%s/%s", bunnyHost, bunnyZone, fileKey)

	// Create DELETE request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %v", err)
	}
	req.Header.Set("AccessKey", bunnyAPIKey)

	// Perform request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bunny delete failed: %d %s", resp.StatusCode, string(body))
	}

	return nil
}
