package main

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/textproto"
	"os"

	"dashlearn/internal/utils"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	tmp, err := os.CreateTemp("", "lurnic-r2-smoke-*.bin")
	if err != nil {
		log.Fatalf("create temp file: %v", err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	payload := bytes.Repeat([]byte("r2-smoke-test\n"), 128)
	if _, err := tmp.Write(payload); err != nil {
		_ = tmp.Close()
		log.Fatalf("write temp file: %v", err)
	}
	if err := tmp.Close(); err != nil {
		log.Fatalf("close temp file: %v", err)
	}

	f, err := os.Open(tmpPath)
	if err != nil {
		log.Fatalf("open temp file: %v", err)
	}
	defer f.Close()

	h := &multipart.FileHeader{
		Filename: "r2-smoke-test.txt",
		Header: textproto.MIMEHeader{
			"Content-Type": []string{"text/plain"},
		},
	}

	publicURL, err := utils.UploadToBunny(f, h)
	if err != nil {
		log.Fatalf("upload failed: %v", err)
	}
	fmt.Println("upload ok:", publicURL)

	if err := utils.DeleteFromBunny(publicURL); err != nil {
		log.Fatalf("delete failed: %v", err)
	}
	fmt.Println("delete ok")
}

