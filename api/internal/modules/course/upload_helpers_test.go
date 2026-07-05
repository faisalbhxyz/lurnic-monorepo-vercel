package course

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"testing"

	"gorm.io/datatypes"
)

func TestApplyAssignmentAttachmentUploads_MergesExistingAndNew(t *testing.T) {
	existing := []AssignmentAttachmentInput{
		{URL: "https://cdn.example.com/old.pdf", FileName: "old.pdf", MimeType: "application/pdf", Size: 100},
	}
	existingJSON, err := json.Marshal(existing)
	if err != nil {
		t.Fatalf("marshal existing: %v", err)
	}
	j := datatypes.JSON(existingJSON)

	chapters := []CreateCourseChapter{
		{
			Assignments: []CreateAssignmentInput{
				{Title: "A1", Attachments: &j},
			},
		},
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("assignment_attachments[0][0][]", "new.pdf")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := part.Write([]byte("%PDF test")); err != nil {
		t.Fatalf("write: %v", err)
	}
	_ = writer.Close()

	reader := multipart.NewReader(&body, writer.Boundary())
	form, err := reader.ReadForm(1 << 20)
	if err != nil {
		t.Fatalf("read form: %v", err)
	}

	origUpload := uploadAssignmentFile
	t.Cleanup(func() { uploadAssignmentFile = origUpload })
	uploadAssignmentFile = func(_ multipart.File, header *multipart.FileHeader) (string, string, int64, error) {
		return "https://cdn.example.com/" + header.Filename, "application/pdf", int64(len("%PDF test")), nil
	}

	if err := applyAssignmentAttachmentUploads(chapters, form); err != nil {
		t.Fatalf("applyAssignmentAttachmentUploads: %v", err)
	}

	merged, err := parseAssignmentAttachments(chapters[0].Assignments[0].Attachments)
	if err != nil {
		t.Fatalf("parse merged: %v", err)
	}
	if len(merged) != 2 {
		t.Fatalf("expected 2 attachments, got %d", len(merged))
	}
	if merged[0].FileName != "old.pdf" {
		t.Fatalf("expected old attachment preserved, got %q", merged[0].FileName)
	}
	if merged[1].FileName != "new.pdf" {
		t.Fatalf("expected new attachment appended, got %q", merged[1].FileName)
	}
}
