package academicnote

type CreateClassInput struct {
	Title       string  `json:"title" form:"title" binding:"required"`
	Slug        string  `json:"slug" form:"slug"`
	IconLabel   *string `json:"icon_label" form:"icon_label"`
	IconColor   *string `json:"icon_color" form:"icon_color"`
	Position    int     `json:"position" form:"position"`
	IsPublished *bool   `json:"is_published" form:"is_published"`
}

type UpdateClassInput struct {
	Title       string  `json:"title" form:"title" binding:"required"`
	Slug        string  `json:"slug" form:"slug"`
	IconLabel   *string `json:"icon_label" form:"icon_label"`
	IconColor   *string `json:"icon_color" form:"icon_color"`
	Position    int     `json:"position" form:"position"`
	IsPublished *bool   `json:"is_published" form:"is_published"`
}

type CreateSubjectInput struct {
	ClassID     uint    `json:"class_id" form:"class_id" binding:"required"`
	Title       string  `json:"title" form:"title" binding:"required"`
	Slug        string  `json:"slug" form:"slug"`
	Position    int     `json:"position" form:"position"`
	IsPublished *bool   `json:"is_published" form:"is_published"`
}

type UpdateSubjectInput struct {
	Title       string  `json:"title" form:"title" binding:"required"`
	Slug        string  `json:"slug" form:"slug"`
	Position    int     `json:"position" form:"position"`
	IsPublished *bool   `json:"is_published" form:"is_published"`
}

type CreatePaperInput struct {
	SubjectID   uint    `json:"subject_id" form:"subject_id" binding:"required"`
	Title       string  `json:"title" form:"title" binding:"required"`
	Slug        string  `json:"slug" form:"slug"`
	IconLabel   *string `json:"icon_label" form:"icon_label"`
	IconColor   *string `json:"icon_color" form:"icon_color"`
	Position    int     `json:"position" form:"position"`
	IsPublished *bool   `json:"is_published" form:"is_published"`
}

type UpdatePaperInput struct {
	Title       string  `json:"title" form:"title" binding:"required"`
	Slug        string  `json:"slug" form:"slug"`
	IconLabel   *string `json:"icon_label" form:"icon_label"`
	IconColor   *string `json:"icon_color" form:"icon_color"`
	Position    int     `json:"position" form:"position"`
	IsPublished *bool   `json:"is_published" form:"is_published"`
}

type CreateNoteInput struct {
	PaperID     uint    `json:"paper_id" form:"paper_id" binding:"required"`
	Title       string  `json:"title" form:"title" binding:"required"`
	Subtitle    *string `json:"subtitle" form:"subtitle"`
	Thumbnail   *string `json:"thumbnail"`
	PdfURL      string  `json:"pdf_url" form:"pdf_url"`
	PdfFileName *string `json:"pdf_file_name" form:"pdf_file_name"`
	Position    int     `json:"position" form:"position"`
	IsPublished *bool   `json:"is_published" form:"is_published"`
}

type UpdateNoteInput struct {
	Title       string  `json:"title" form:"title" binding:"required"`
	Subtitle    *string `json:"subtitle" form:"subtitle"`
	Thumbnail   *string `json:"thumbnail"`
	PdfURL      *string `json:"pdf_url" form:"pdf_url"`
	PdfFileName *string `json:"pdf_file_name" form:"pdf_file_name"`
	Position    int     `json:"position" form:"position"`
	IsPublished *bool   `json:"is_published" form:"is_published"`
}
