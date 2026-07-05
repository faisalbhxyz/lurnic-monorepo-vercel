package course

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"

	"dashlearn/internal/utils"

	"gorm.io/datatypes"
)

type AssignmentAttachmentInput struct {
	URL      string `json:"url"`
	FileName string `json:"file_name"`
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
}

var uploadAssignmentFile = func(file multipart.File, header *multipart.FileHeader) (url, mimeType string, size int64, err error) {
	url, err = utils.UploadToBunny(file, header)
	if err != nil {
		return "", "", 0, err
	}
	buf := make([]byte, 512)
	_, _ = file.Read(buf)
	mimeType = http.DetectContentType(buf)
	return url, mimeType, header.Size, nil
}

func applyAssignmentAttachmentUploads(chapters []CreateCourseChapter, form *multipart.Form) error {
	if form == nil {
		return nil
	}

	for key, files := range form.File {
		if !stringsHasPrefix(key, "assignment_attachments[") {
			continue
		}

		var chapterIndex, assignmentIndex int
		if _, err := fmt.Sscanf(key, "assignment_attachments[%d][%d][]", &chapterIndex, &assignmentIndex); err != nil {
			continue
		}
		if chapterIndex < 0 || chapterIndex >= len(chapters) {
			continue
		}
		if assignmentIndex < 0 || assignmentIndex >= len(chapters[chapterIndex].Assignments) {
			continue
		}

		existing, err := parseAssignmentAttachments(chapters[chapterIndex].Assignments[assignmentIndex].Attachments)
		if err != nil {
			return err
		}

		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				return fmt.Errorf("failed to open assignment attachment: %w", err)
			}

			url, mimeType, size, uploadErr := uploadAssignmentFile(file, fileHeader)
			_ = file.Close()
			if uploadErr != nil {
				return fmt.Errorf("failed to upload assignment attachment: %w", uploadErr)
			}

			existing = append(existing, AssignmentAttachmentInput{
				URL:      url,
				FileName: fileHeader.Filename,
				MimeType: mimeType,
				Size:     size,
			})
		}

		attachmentsJSON, err := assignmentAttachmentsToJSON(existing)
		if err != nil {
			return err
		}
		chapters[chapterIndex].Assignments[assignmentIndex].Attachments = attachmentsJSON
	}

	return nil
}

func parseAssignmentAttachments(raw *datatypes.JSON) ([]AssignmentAttachmentInput, error) {
	if raw == nil {
		return []AssignmentAttachmentInput{}, nil
	}
	var items []AssignmentAttachmentInput
	if err := json.Unmarshal(*raw, &items); err != nil {
		return nil, fmt.Errorf("invalid assignment attachments: %w", err)
	}
	return items, nil
}

func assignmentAttachmentsToJSON(items []AssignmentAttachmentInput) (*datatypes.JSON, error) {
	if len(items) == 0 {
		return nil, nil
	}
	b, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}
	j := datatypes.JSON(b)
	return &j, nil
}

func stringsHasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
