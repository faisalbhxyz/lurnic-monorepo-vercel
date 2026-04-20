package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gosimple/slug"
)

func getR2SigningRegion() string {
	// R2 is not an AWS region, but SDK signing still requires a region string.
	if v := os.Getenv("R2_REGION"); v != "" {
		return v
	}
	return "us-east-1"
}

func initR2(ctx context.Context) (*s3.Client, error) {
	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	if accountID == "" || accessKeyID == "" || secretAccessKey == "" {
		return nil, fmt.Errorf("missing R2 configuration (R2_ACCOUNT_ID, R2_ACCESS_KEY_ID, R2_SECRET_ACCESS_KEY)")
	}

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(getR2SigningRegion()),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config for R2: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID))
		// R2 is S3-compatible but rejects some AWS-specific headers/behavior.
		// Enabling path-style addressing avoids extra host rewriting and reduces
		// unexpected request parameters for non-AWS endpoints.
		o.UsePathStyle = true
	})

	return client, nil
}

func normalizePublicBaseURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("R2_PUBLIC_BASE_URL not set")
	}
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		raw = "https://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("invalid R2_PUBLIC_BASE_URL: %w", err)
	}
	u.Path = strings.TrimRight(u.Path, "/")
	return u.String(), nil
}

func getR2Bucket() (string, error) {
	b := strings.TrimSpace(os.Getenv("R2_BUCKET"))
	if b == "" {
		return "", fmt.Errorf("R2_BUCKET not set")
	}
	return b, nil
}

func getR2UploadPrefix() string {
	return strings.Trim(strings.TrimSpace(os.Getenv("R2_UPLOAD_PREFIX")), "/")
}

func buildObjectKey(prefix, filename string) string {
	if prefix == "" {
		return filename
	}
	return prefix + "/" + filename
}

// UploadToBunny uploads a file to object storage (Cloudflare R2) and returns its public URL.
func UploadToBunny(file multipart.File, header *multipart.FileHeader) (string, error) {
	ctx := context.Background()

	// Clean and slugify filename
	re := regexp.MustCompile(`[*#+~()'"!:@]`)
	cleanedName := re.ReplaceAllString(header.Filename, "")
	ext := filepath.Ext(cleanedName)
	nameWithoutExt := strings.TrimSuffix(cleanedName, ext)
	slugifiedName := slug.Make(nameWithoutExt)
	safeName := fmt.Sprintf("%d_%s%s", time.Now().Unix(), slugifiedName, ext)

	// Reset file pointer just in case
	if seeker, ok := file.(io.Seeker); ok {
		seeker.Seek(0, 0)
	}

	fileBytes := new(bytes.Buffer)
	if _, err := fileBytes.ReadFrom(file); err != nil {
		return "", fmt.Errorf("failed to read upload file: %w", err)
	}

	publicBaseURL, err := normalizePublicBaseURL(os.Getenv("R2_PUBLIC_BASE_URL"))
	if err != nil {
		return "", err
	}
	bucket, err := getR2Bucket()
	if err != nil {
		return "", err
	}

	objectKey := buildObjectKey(getR2UploadPrefix(), safeName)

	r2Client, err := initR2(ctx)
	if err != nil {
		return "", err
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = r2Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(fileBytes.Bytes()),
		ContentType: aws.String(contentType),
		ContentLength: aws.Int64(int64(fileBytes.Len())),
	}, func(opts *s3.Options) {
		// R2 can reject SDK-managed checksum behavior for PutObject; disable it.
		opts.RequestChecksumCalculation = aws.RequestChecksumCalculation(0)
		opts.ResponseChecksumValidation = aws.ResponseChecksumValidation(0)
	})
	if err != nil {
		return "", fmt.Errorf("r2 upload failed: %w", err)
	}

	return publicBaseURL + "/" + objectKey, nil
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
	ctx := context.Background()

	fileKey, err := getFileKeyFromURL(fileURL)
	if err != nil {
		return fmt.Errorf("invalid CDN URL: %v", err)
	}

	bucket, err := getR2Bucket()
	if err != nil {
		return err
	}

	r2Client, err := initR2(ctx)
	if err != nil {
		return err
	}

	_, err = r2Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		return fmt.Errorf("r2 delete failed: %w", err)
	}

	return nil
}
