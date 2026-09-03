package course

import (
	"testing"

	"dashlearn/internal/models"
	"dashlearn/internal/utils"
)

func TestBuildLessonDownloadResult_GoogleDrive(t *testing.T) {
	lesson := &models.CourseLesson{
		Title:      "Chapter 1 — Introduction",
		LessonType: models.Video,
		Source: utils.JSONB[models.Source]{
			Data: models.Source{
				Data:   "https://drive.google.com/file/d/1AbCdEfGhIjKlMnOpQrStUvWxYz/view?usp=sharing",
				IsFile: false,
			},
		},
	}

	got, err := buildLessonDownloadResult(lesson)
	if err != nil {
		t.Fatalf("buildLessonDownloadResult() error = %v", err)
	}

	wantURL := "https://drive.usercontent.google.com/download?id=1AbCdEfGhIjKlMnOpQrStUvWxYz&export=download&confirm=t"
	if got.RedirectURL != wantURL {
		t.Fatalf("RedirectURL = %q, want %q", got.RedirectURL, wantURL)
	}
	if got.ContentType != "video/mp4" {
		t.Fatalf("ContentType = %q, want video/mp4", got.ContentType)
	}
	if got.FileName != "chapter-1-introduction.mp4" {
		t.Fatalf("FileName = %q", got.FileName)
	}
}

func TestBuildLessonDownloadResult_DirectMP4(t *testing.T) {
	sourceURL := "https://cdn.example.com/videos/lesson-1.mp4?token=abc"
	lesson := &models.CourseLesson{
		Title:      "Lesson One",
		LessonType: models.Video,
		Source: utils.JSONB[models.Source]{
			Data: models.Source{
				Data:   sourceURL,
				IsFile: true,
			},
		},
	}

	got, err := buildLessonDownloadResult(lesson)
	if err != nil {
		t.Fatalf("buildLessonDownloadResult() error = %v", err)
	}
	if got.RedirectURL != sourceURL {
		t.Fatalf("RedirectURL = %q, want %q", got.RedirectURL, sourceURL)
	}
}

func TestBuildLessonDownloadResult_YouTubeNotDownloadable(t *testing.T) {
	lesson := &models.CourseLesson{
		Title:      "YouTube Lesson",
		LessonType: models.Video,
		Source: utils.JSONB[models.Source]{
			Data: models.Source{
				Data:   "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
				IsFile: false,
			},
		},
	}

	_, err := buildLessonDownloadResult(lesson)
	if err != ErrDownloadNotDownloadable {
		t.Fatalf("expected ErrDownloadNotDownloadable, got %v", err)
	}
}

func TestBuildLessonDownloadResult_YouTubeWithSeparateDriveURL(t *testing.T) {
	driveURL := "https://drive.google.com/file/d/1AbCdEfGhIjKlMnOpQrStUvWxYz/view?usp=sharing"
	lesson := &models.CourseLesson{
		Title:      "YouTube Lesson",
		LessonType: models.Video,
		SourceType: models.Youtube,
		Source: utils.JSONB[models.Source]{
			Data: models.Source{
				Data:     `<iframe src="https://www.youtube.com/embed/vHXkebRshWk"></iframe>`,
				IsFile:   false,
				DriveURL: driveURL,
			},
		},
	}

	got, err := buildLessonDownloadResult(lesson)
	if err != nil {
		t.Fatalf("buildLessonDownloadResult() error = %v", err)
	}
	wantURL := "https://drive.usercontent.google.com/download?id=1AbCdEfGhIjKlMnOpQrStUvWxYz&export=download&confirm=t"
	if got.RedirectURL != wantURL {
		t.Fatalf("RedirectURL = %q, want %q", got.RedirectURL, wantURL)
	}
}

func TestLessonDownloadFileName(t *testing.T) {
	if got := lessonDownloadFileName("  Hello World!  ", ".mp4"); got != "hello-world.mp4" {
		t.Fatalf("unexpected filename %q", got)
	}
}
