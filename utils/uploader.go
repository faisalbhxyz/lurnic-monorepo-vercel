package utils

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	endpoints "github.com/aws/smithy-go/endpoints"
	"github.com/gosimple/slug"
)

type customResolver struct{}

func (r customResolver) ResolveEndpoint(ctx context.Context, params s3.EndpointParameters) (endpoints.Endpoint, error) {
	bucket := os.Getenv("DO_BUCKET_NAME")
	region := os.Getenv("DO_REGION")
	if region == "" || bucket == "" {
		return endpoints.Endpoint{}, fmt.Errorf("DO_REGION or DO_BUCKET_NAME not set")
	}

	// This is the correct endpoint: bucket.region.digitaloceanspaces.com
	endpointURL := fmt.Sprintf("https://%s.%s.digitaloceanspaces.com", bucket, region)
	parsedURL, err := url.Parse(endpointURL)
	if err != nil {
		return endpoints.Endpoint{}, fmt.Errorf("invalid endpoint URL: %w", err)
	}

	props := smithy.Properties{}
	props.Set("SigningRegion", region)

	return endpoints.Endpoint{
		URI:        *parsedURL,
		Properties: props,
		Headers:    map[string][]string{},
	}, nil
}

func initS3(ctx context.Context) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(os.Getenv("DO_REGION")),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("DO_ACCESS_KEY"),
			os.Getenv("DO_SECRET"),
			"",
		)),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.EndpointResolverV2 = customResolver{}
	})

	return client, nil
}
func UploadFile(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileBytes := new(bytes.Buffer)
	_, err = fileBytes.ReadFrom(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	s3Client, err := initS3(ctx)
	if err != nil {
		return "", err
	}

	fileKey := fmt.Sprintf("%s/%d_%s%s",
		os.Getenv("DO_UPLOAD_FOLDER"),
		time.Now().Unix(),
		slug.Make(fileHeader.Filename),
		filepath.Ext(fileHeader.Filename),
	)

	input := &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("DO_BUCKET_NAME")),
		Key:         aws.String(fileKey),
		Body:        bytes.NewReader(fileBytes.Bytes()),
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
		ACL:         s3types.ObjectCannedACLPublicRead,
	}

	_, err = s3Client.PutObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	url := fmt.Sprintf("https://%s.%s.digitaloceanspaces.com/%s",
		os.Getenv("DO_BUCKET_NAME"),
		os.Getenv("DO_REGION"),
		fileKey,
	)

	return url, nil
}
