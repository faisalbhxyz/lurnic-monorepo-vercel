package response

import "time"

type AcademicNoteClassResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	IconLabel   *string   `json:"icon_label"`
	IconColor   *string   `json:"icon_color"`
	IconImage   *string   `json:"icon_image"`
	Position    int       `json:"position"`
	IsPublished bool      `json:"is_published"`
	NoteCount   int       `json:"note_count,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AcademicNoteSubjectResponse struct {
	ID          uint                        `json:"id"`
	ClassID     uint                        `json:"class_id"`
	Title       string                      `json:"title"`
	Slug        string                      `json:"slug"`
	Position    int                         `json:"position"`
	IsPublished bool                        `json:"is_published"`
	NoteCount   int                         `json:"note_count,omitempty"`
	Papers      []AcademicNotePaperResponse `json:"papers,omitempty"`
	CreatedAt   time.Time                   `json:"created_at"`
	UpdatedAt   time.Time                   `json:"updated_at"`
}

type AcademicNotePaperResponse struct {
	ID          uint                     `json:"id"`
	SubjectID   uint                     `json:"subject_id"`
	Title       string                   `json:"title"`
	Slug        string                   `json:"slug"`
	IconLabel   *string                  `json:"icon_label"`
	IconColor   *string                  `json:"icon_color"`
	Position    int                      `json:"position"`
	IsPublished bool                     `json:"is_published"`
	NoteCount   int                      `json:"note_count,omitempty"`
	Notes       []AcademicNoteItemResponse `json:"notes,omitempty"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

type AcademicNoteItemResponse struct {
	ID          uint      `json:"id"`
	PaperID     uint      `json:"paper_id"`
	Title       string    `json:"title"`
	Subtitle    *string   `json:"subtitle"`
	Thumbnail   *string   `json:"thumbnail"`
	PdfURL      string    `json:"pdf_url"`
	PdfFileName *string   `json:"pdf_file_name"`
	Position    int       `json:"position"`
	IsPublished bool      `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AcademicNoteClassDetailResponse struct {
	ID          uint                          `json:"id"`
	Title       string                        `json:"title"`
	Slug        string                        `json:"slug"`
	IconLabel   *string                       `json:"icon_label"`
	IconColor   *string                       `json:"icon_color"`
	IconImage   *string                       `json:"icon_image"`
	Position    int                           `json:"position"`
	IsPublished bool                          `json:"is_published"`
	Subjects    []AcademicNoteSubjectResponse `json:"subjects"`
}

type AcademicNotePaperDetailResponse struct {
	Class   AcademicNoteClassResponse   `json:"class"`
	Subject AcademicNoteSubjectResponse `json:"subject"`
	Paper   AcademicNotePaperResponse   `json:"paper"`
	Notes   []AcademicNoteItemResponse  `json:"notes"`
}

type AcademicNoteClassAdminResponse struct {
	AcademicNoteClassResponse
	Subjects []AcademicNoteSubjectAdminResponse `json:"subjects"`
}

type AcademicNoteSubjectAdminResponse struct {
	AcademicNoteSubjectResponse
	Papers []AcademicNotePaperAdminResponse `json:"papers"`
}

type AcademicNotePaperAdminResponse struct {
	AcademicNotePaperResponse
	Notes []AcademicNoteItemResponse `json:"notes"`
}
