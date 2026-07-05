package models

import "time"

type AcademicNoteClass struct {
	ID          uint                 `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID    uint                 `gorm:"column:tenant_id" json:"-"`
	Title       string               `json:"title"`
	Slug        string               `json:"slug"`
	IconLabel   *string              `gorm:"column:icon_label" json:"icon_label"`
	IconColor   *string              `gorm:"column:icon_color" json:"icon_color"`
	Position    int                  `gorm:"default:0" json:"position"`
	IsPublished bool                 `gorm:"column:is_published;default:true" json:"is_published"`
	CreatedAt   time.Time            `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time            `gorm:"autoUpdateTime" json:"updated_at"`
	Subjects    []AcademicNoteSubject `gorm:"foreignKey:ClassID;references:ID" json:"subjects,omitempty"`
}

type AcademicNoteSubject struct {
	ID          uint               `gorm:"primaryKey;autoIncrement" json:"id"`
	ClassID     uint               `gorm:"column:class_id" json:"class_id"`
	Title       string             `json:"title"`
	Slug        string             `json:"slug"`
	Position    int                `gorm:"default:0" json:"position"`
	IsPublished bool               `gorm:"column:is_published;default:true" json:"is_published"`
	CreatedAt   time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
	Papers      []AcademicNotePaper `gorm:"foreignKey:SubjectID;references:ID" json:"papers,omitempty"`
}

type AcademicNotePaper struct {
	ID          uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	SubjectID   uint          `gorm:"column:subject_id" json:"subject_id"`
	Title       string        `json:"title"`
	Slug        string        `json:"slug"`
	IconLabel   *string       `gorm:"column:icon_label" json:"icon_label"`
	IconColor   *string       `gorm:"column:icon_color" json:"icon_color"`
	Position    int           `gorm:"default:0" json:"position"`
	IsPublished bool          `gorm:"column:is_published;default:true" json:"is_published"`
	CreatedAt   time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	Notes       []AcademicNote `gorm:"foreignKey:PaperID;references:ID" json:"notes,omitempty"`
}

type AcademicNote struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PaperID     uint      `gorm:"column:paper_id" json:"paper_id"`
	Title       string    `json:"title"`
	Subtitle    *string   `json:"subtitle"`
	Thumbnail   *string   `json:"thumbnail"`
	PdfURL      string    `gorm:"column:pdf_url" json:"pdf_url"`
	PdfFileName *string   `gorm:"column:pdf_file_name" json:"pdf_file_name"`
	Position    int       `gorm:"default:0" json:"position"`
	IsPublished bool      `gorm:"column:is_published;default:true" json:"is_published"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
