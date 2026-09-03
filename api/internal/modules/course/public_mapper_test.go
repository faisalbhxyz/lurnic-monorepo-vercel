package course

import (
	"testing"
	"time"

	"dashlearn/internal/models"
	"dashlearn/internal/response"
	"dashlearn/internal/utils"
)

func TestBuildCourseDetailsPublicResponse_IncludesLessonResources(t *testing.T) {
	desc := "Lesson notes"
	course := &models.CourseDetails{
		ID:    5,
		Title: "Test Course",
		Slug:  "test-course",
		Chapters: []models.CourseChapter{
			{
				ID:       10,
				Title:    "Chapter 1",
				Access:   models.Published,
				CourseID: 5,
				Lessons: []models.CourseLesson{
					{
						ID:          42,
						Title:       "Intro",
						Description: &desc,
						LessonType:  models.Video,
						SourceType:  models.Youtube,
						IsPublished: true,
						IsPublic:    false,
						ChapterID:   10,
						Resources: []models.LessonResource{
							{
								ID:        7,
								CourseID:  5,
								LessonID:  42,
								MimeType:  "application/pdf",
								Title:     "chapter-notes.pdf",
								FilePath:  "https://cdn.example.com/chapter-notes.pdf",
								Position:  0,
								FileSize:  20480,
								CreatedAt: time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC),
								UpdatedAt: time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC),
							},
						},
					},
				},
			},
		},
	}

	got := BuildCourseDetailsPublicResponse(course)
	if len(got.Chapters) != 1 {
		t.Fatalf("expected 1 chapter, got %d", len(got.Chapters))
	}
	if len(got.Chapters[0].Lessons) != 1 {
		t.Fatalf("expected 1 lesson, got %d", len(got.Chapters[0].Lessons))
	}

	resources := got.Chapters[0].Lessons[0].Resources
	if len(resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(resources))
	}

	res := resources[0]
	if res.MimeType != "application/pdf" {
		t.Fatalf("expected mime_type application/pdf, got %q", res.MimeType)
	}
	if res.Title != "chapter-notes.pdf" {
		t.Fatalf("unexpected title %q", res.Title)
	}
	if res.FilePath != "https://cdn.example.com/chapter-notes.pdf" {
		t.Fatalf("unexpected file_path %q", res.FilePath)
	}
	if res.URL != res.FilePath {
		t.Fatalf("expected url alias of file_path, got url=%q file_path=%q", res.URL, res.FilePath)
	}
}

func TestBuildCourseDetailsPublicResponse_IncludesGoogleDriveOfflineFields(t *testing.T) {
	driveURL := "https://drive.google.com/file/d/1AbCdEfGhIjKlMnOpQrStUvWxYz/view?usp=sharing"
	fileID := "1AbCdEfGhIjKlMnOpQrStUvWxYz"
	course := &models.CourseDetails{
		ID:    5,
		Title: "Drive Course",
		Slug:  "drive-course",
		Chapters: []models.CourseChapter{
			{
				ID:       10,
				Title:    "Chapter 1",
				Access:   models.Published,
				CourseID: 5,
				Lessons: []models.CourseLesson{
					{
						ID:         42,
						Title:      "Intro",
						LessonType: models.Video,
						SourceType: models.GoogleDrive,
						Source: utils.JSONB[models.Source]{
							Data: models.Source{
								Data:        driveURL,
								IsFile:      false,
								DriveFileID: &fileID,
							},
						},
						IsPublished: true,
						IsPublic:    false,
						ChapterID:   10,
					},
				},
			},
		},
	}

	got := BuildCourseDetailsPublicResponse(course)
	lesson := got.Chapters[0].Lessons[0]

	if lesson.Source.Data.Data != driveURL {
		t.Fatalf("expected drive url preserved, got %q", lesson.Source.Data.Data)
	}
	if !lesson.OfflineDownloadable {
		t.Fatal("expected offline_downloadable=true for drive video lesson")
	}
	if lesson.DownloadURL == nil || *lesson.DownloadURL != driveURL {
		t.Fatalf("expected download_url=%q, got %#v", driveURL, lesson.DownloadURL)
	}
}

func TestBuildCourseDetailsPublicResponse_YouTubeKeepsPlaybackWithDriveURL(t *testing.T) {
	driveURL := "https://drive.google.com/file/d/1AbCdEfGhIjKlMnOpQrStUvWxYz/view?usp=sharing"
	fileID := "1AbCdEfGhIjKlMnOpQrStUvWxYz"
	youtubeEmbed := `<iframe src="https://www.youtube.com/embed/vHXkebRshWk"></iframe>`
	course := &models.CourseDetails{
		ID:    5,
		Title: "Mixed Source Course",
		Slug:  "mixed-source-course",
		Chapters: []models.CourseChapter{
			{
				ID:       10,
				Title:    "Chapter 1",
				Access:   models.Published,
				CourseID: 5,
				Lessons: []models.CourseLesson{
					{
						ID:         42,
						Title:      "Intro",
						LessonType: models.Video,
						SourceType: models.Youtube,
						Source: utils.JSONB[models.Source]{
							Data: models.Source{
								Data:        youtubeEmbed,
								IsFile:      false,
								DriveURL:    driveURL,
								DriveFileID: &fileID,
							},
						},
						IsPublished: true,
						ChapterID:   10,
					},
				},
			},
		},
	}

	got := BuildCourseDetailsPublicResponse(course)
	lesson := got.Chapters[0].Lessons[0]
	if lesson.SourceType != models.Youtube {
		t.Fatalf("expected youtube source type, got %q", lesson.SourceType)
	}
	if lesson.Source.Data.Data != youtubeEmbed {
		t.Fatalf("expected youtube embed preserved, got %q", lesson.Source.Data.Data)
	}
	if !lesson.OfflineDownloadable {
		t.Fatal("expected offline_downloadable=true when drive_url is set")
	}
	if lesson.DownloadURL == nil || *lesson.DownloadURL != driveURL {
		t.Fatalf("expected download_url=%q, got %#v", driveURL, lesson.DownloadURL)
	}
}

func TestMapLessonResources_EmptyReturnsNil(t *testing.T) {
	if got := response.MapLessonResources(nil); got != nil {
		t.Fatalf("expected nil for empty input, got %#v", got)
	}
}
