package course

import (
	"testing"

	"dashlearn/internal/models"
)

func TestExtractDriveFileID(t *testing.T) {
	tests := []struct {
		name   string
		raw    string
		wantID string
		wantOK bool
	}{
		{
			name:   "file view link",
			raw:    "https://drive.google.com/file/d/1AbCdEfGhIjKlMnOpQrStUvWxYz/view?usp=sharing",
			wantID: "1AbCdEfGhIjKlMnOpQrStUvWxYz",
			wantOK: true,
		},
		{
			name:   "open id link",
			raw:    "https://drive.google.com/open?id=1XyZ987654321",
			wantID: "1XyZ987654321",
			wantOK: true,
		},
		{
			name:   "uc export link",
			raw:    "https://drive.google.com/uc?export=download&id=abc123_-XYZ",
			wantID: "abc123_-XYZ",
			wantOK: true,
		},
		{
			name:   "youtube url",
			raw:    "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			wantOK: false,
		},
		{
			name:   "empty",
			raw:    "   ",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotOK := ExtractDriveFileID(tt.raw)
			if gotOK != tt.wantOK {
				t.Fatalf("ExtractDriveFileID(%q) ok = %v, want %v", tt.raw, gotOK, tt.wantOK)
			}
			if gotID != tt.wantID {
				t.Fatalf("ExtractDriveFileID(%q) id = %q, want %q", tt.raw, gotID, tt.wantID)
			}
		})
	}
}

func TestOfflineDownloadable(t *testing.T) {
	driveURL := "https://drive.google.com/file/d/abc123/view"
	fileID := "abc123"

	tests := []struct {
		name       string
		lessonType models.LessonType
		source     models.Source
		want       bool
	}{
		{
			name:       "video with drive url",
			lessonType: models.Video,
			source:     models.Source{Data: driveURL, IsFile: false},
			want:       true,
		},
		{
			name:       "video with direct mp4",
			lessonType: models.Video,
			source:     models.Source{Data: "https://cdn.example.com/lesson.mp4", IsFile: true},
			want:       true,
		},
		{
			name:       "text lesson with drive url",
			lessonType: models.Text,
			source:     models.Source{Data: driveURL},
			want:       false,
		},
		{
			name:       "video with youtube url",
			lessonType: models.Video,
			source:     models.Source{Data: "https://www.youtube.com/watch?v=test"},
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OfflineDownloadable(tt.lessonType, tt.source); got != tt.want {
				t.Fatalf("OfflineDownloadable() = %v, want %v", got, tt.want)
			}
		})
	}

	_, normalizedSource := NormalizeLessonSource(models.Video, models.UploadFile, models.Source{
		Data:   driveURL,
		IsFile: true,
	})
	if normalizedSource.DriveFileID == nil || *normalizedSource.DriveFileID != fileID {
		t.Fatalf("expected drive_file_id %q, got %#v", fileID, normalizedSource.DriveFileID)
	}
}

func TestValidateLessonSource(t *testing.T) {
	if err := ValidateLessonSource(models.GoogleDrive, models.Source{Data: "https://drive.google.com/file/d/abc/view"}); err != nil {
		t.Fatalf("expected valid google_drive source, got %v", err)
	}
	if err := ValidateLessonSource(models.GoogleDrive, models.Source{Data: "https://example.com/video.mp4"}); err == nil {
		t.Fatal("expected validation error for invalid google_drive source")
	}
}
