package academicnote

import (
	"dashlearn/internal/models"
	"dashlearn/internal/response"
	"dashlearn/internal/utils"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type Service interface {
	// Classes
	GetAllClasses(tenantID uint) ([]response.AcademicNoteClassResponse, error)
	GetClassByID(tenantID uint, id uint64) (*response.AcademicNoteClassAdminResponse, error)
	CreateClass(input CreateClassInput, tenantID uint) error
	UpdateClass(id uint64, input UpdateClassInput, tenantID uint) error
	DeleteClass(id uint64, tenantID uint) error

	// Subjects
	CreateSubject(input CreateSubjectInput, tenantID uint) error
	UpdateSubject(id uint64, input UpdateSubjectInput, tenantID uint) error
	DeleteSubject(id uint64, tenantID uint) error

	// Papers
	CreatePaper(input CreatePaperInput, tenantID uint) error
	UpdatePaper(id uint64, input UpdatePaperInput, tenantID uint) error
	DeletePaper(id uint64, tenantID uint) error

	// Notes
	CreateNote(input CreateNoteInput, tenantID uint) error
	UpdateNote(id uint64, input UpdateNoteInput, tenantID uint) error
	DeleteNote(id uint64, tenantID uint) error

	// Public storefront
	GetPublicClasses(tenantID uint) ([]response.AcademicNoteClassResponse, error)
	GetPublicClassBySlug(tenantID uint, classSlug string) (*response.AcademicNoteClassDetailResponse, error)
	GetPublicNotesByPaperSlug(tenantID uint, classSlug, subjectSlug, paperSlug string) (*response.AcademicNotePaperDetailResponse, error)
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}

func (s *service) countNotesForClass(classID uint, publishedOnly bool) (int, error) {
	q := s.db.Table("academic_notes").
		Joins("JOIN academic_note_papers ON academic_note_papers.id = academic_notes.paper_id").
		Joins("JOIN academic_note_subjects ON academic_note_subjects.id = academic_note_papers.subject_id").
		Where("academic_note_subjects.class_id = ?", classID)
	if publishedOnly {
		q = q.Where("academic_notes.is_published = ? AND academic_note_papers.is_published = ? AND academic_note_subjects.is_published = ?", true, true, true)
	}
	var count int64
	err := q.Count(&count).Error
	return int(count), err
}

func (s *service) countNotesForSubject(subjectID uint, publishedOnly bool) (int, error) {
	q := s.db.Table("academic_notes").
		Joins("JOIN academic_note_papers ON academic_note_papers.id = academic_notes.paper_id").
		Where("academic_note_papers.subject_id = ?", subjectID)
	if publishedOnly {
		q = q.Where("academic_notes.is_published = ? AND academic_note_papers.is_published = ?", true, true)
	}
	var count int64
	err := q.Count(&count).Error
	return int(count), err
}

func (s *service) countNotesForPaper(paperID uint, publishedOnly bool) (int, error) {
	q := s.db.Model(&models.AcademicNote{}).Where("paper_id = ?", paperID)
	if publishedOnly {
		q = q.Where("is_published = ?", true)
	}
	var count int64
	err := q.Count(&count).Error
	return int(count), err
}

func resolveSlug(title, slug string, id uint) string {
	if slug != "" {
		return slug
	}
	return utils.Slugify(title) + "-" + strconv.Itoa(int(id))
}

func boolVal(v *bool, fallback bool) bool {
	if v == nil {
		return fallback
	}
	return *v
}

func (s *service) verifyClassOwnership(classID uint, tenantID uint) error {
	var class models.AcademicNoteClass
	if err := s.db.Where("id = ? AND tenant_id = ?", classID, tenantID).First(&class).Error; err != nil {
		return errors.New("class not found")
	}
	return nil
}

func (s *service) verifySubjectOwnership(subjectID uint, tenantID uint) (*models.AcademicNoteSubject, error) {
	var subject models.AcademicNoteSubject
	err := s.db.
		Joins("JOIN academic_note_classes ON academic_note_classes.id = academic_note_subjects.class_id").
		Where("academic_note_subjects.id = ? AND academic_note_classes.tenant_id = ?", subjectID, tenantID).
		First(&subject).Error
	if err != nil {
		return nil, errors.New("subject not found")
	}
	return &subject, nil
}

func (s *service) verifyPaperOwnership(paperID uint, tenantID uint) (*models.AcademicNotePaper, error) {
	var paper models.AcademicNotePaper
	err := s.db.
		Joins("JOIN academic_note_subjects ON academic_note_subjects.id = academic_note_papers.subject_id").
		Joins("JOIN academic_note_classes ON academic_note_classes.id = academic_note_subjects.class_id").
		Where("academic_note_papers.id = ? AND academic_note_classes.tenant_id = ?", paperID, tenantID).
		First(&paper).Error
	if err != nil {
		return nil, errors.New("paper not found")
	}
	return &paper, nil
}

func (s *service) verifyNoteOwnership(noteID uint, tenantID uint) (*models.AcademicNote, error) {
	var note models.AcademicNote
	err := s.db.
		Joins("JOIN academic_note_papers ON academic_note_papers.id = academic_notes.paper_id").
		Joins("JOIN academic_note_subjects ON academic_note_subjects.id = academic_note_papers.subject_id").
		Joins("JOIN academic_note_classes ON academic_note_classes.id = academic_note_subjects.class_id").
		Where("academic_notes.id = ? AND academic_note_classes.tenant_id = ?", noteID, tenantID).
		First(&note).Error
	if err != nil {
		return nil, errors.New("note not found")
	}
	return &note, nil
}

func toNoteItemResponse(n models.AcademicNote) response.AcademicNoteItemResponse {
	return response.AcademicNoteItemResponse{
		ID:          n.ID,
		PaperID:     n.PaperID,
		Title:       n.Title,
		Subtitle:    n.Subtitle,
		Thumbnail:   n.Thumbnail,
		PdfURL:      n.PdfURL,
		PdfFileName: n.PdfFileName,
		Position:    n.Position,
		IsPublished: n.IsPublished,
		CreatedAt:   n.CreatedAt,
		UpdatedAt:   n.UpdatedAt,
	}
}

func (s *service) GetAllClasses(tenantID uint) ([]response.AcademicNoteClassResponse, error) {
	var classes []models.AcademicNoteClass
	if err := s.db.Where("tenant_id = ?", tenantID).Order("position ASC, id ASC").Find(&classes).Error; err != nil {
		return nil, err
	}

	result := make([]response.AcademicNoteClassResponse, 0, len(classes))
	for _, c := range classes {
		count, _ := s.countNotesForClass(c.ID, false)
		result = append(result, response.AcademicNoteClassResponse{
			ID:          c.ID,
			Title:       c.Title,
			Slug:        c.Slug,
			IconLabel:   c.IconLabel,
			IconColor:   c.IconColor,
			Position:    c.Position,
			IsPublished: c.IsPublished,
			NoteCount:   count,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
		})
	}
	return result, nil
}

func (s *service) GetClassByID(tenantID uint, id uint64) (*response.AcademicNoteClassAdminResponse, error) {
	var class models.AcademicNoteClass
	if err := s.db.
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Preload("Subjects", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC, id ASC")
		}).
		Preload("Subjects.Papers", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC, id ASC")
		}).
		Preload("Subjects.Papers.Notes", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC, id ASC")
		}).
		First(&class).Error; err != nil {
		return nil, err
	}

	count, _ := s.countNotesForClass(class.ID, false)
	subjects := make([]response.AcademicNoteSubjectAdminResponse, 0, len(class.Subjects))
	for _, sub := range class.Subjects {
		subCount, _ := s.countNotesForSubject(sub.ID, false)
		papers := make([]response.AcademicNotePaperAdminResponse, 0, len(sub.Papers))
		for _, paper := range sub.Papers {
			paperCount, _ := s.countNotesForPaper(paper.ID, false)
			notes := make([]response.AcademicNoteItemResponse, 0, len(paper.Notes))
			for _, note := range paper.Notes {
				notes = append(notes, toNoteItemResponse(note))
			}
			papers = append(papers, response.AcademicNotePaperAdminResponse{
				AcademicNotePaperResponse: response.AcademicNotePaperResponse{
					ID:          paper.ID,
					SubjectID:   paper.SubjectID,
					Title:       paper.Title,
					Slug:        paper.Slug,
					IconLabel:   paper.IconLabel,
					IconColor:   paper.IconColor,
					Position:    paper.Position,
					IsPublished: paper.IsPublished,
					NoteCount:   paperCount,
					CreatedAt:   paper.CreatedAt,
					UpdatedAt:   paper.UpdatedAt,
				},
				Notes: notes,
			})
		}
		subjects = append(subjects, response.AcademicNoteSubjectAdminResponse{
			AcademicNoteSubjectResponse: response.AcademicNoteSubjectResponse{
				ID:          sub.ID,
				ClassID:     sub.ClassID,
				Title:       sub.Title,
				Slug:        sub.Slug,
				Position:    sub.Position,
				IsPublished: sub.IsPublished,
				NoteCount:   subCount,
				CreatedAt:   sub.CreatedAt,
				UpdatedAt:   sub.UpdatedAt,
			},
			Papers: papers,
		})
	}

	return &response.AcademicNoteClassAdminResponse{
		AcademicNoteClassResponse: response.AcademicNoteClassResponse{
			ID:          class.ID,
			Title:       class.Title,
			Slug:        class.Slug,
			IconLabel:   class.IconLabel,
			IconColor:   class.IconColor,
			Position:    class.Position,
			IsPublished: class.IsPublished,
			NoteCount:   count,
			CreatedAt:   class.CreatedAt,
			UpdatedAt:   class.UpdatedAt,
		},
		Subjects: subjects,
	}, nil
}

func (s *service) CreateClass(input CreateClassInput, tenantID uint) error {
	class := models.AcademicNoteClass{
		TenantID:    tenantID,
		Title:       input.Title,
		Slug:        utils.Slugify(input.Title),
		IconLabel:   input.IconLabel,
		IconColor:   input.IconColor,
		Position:    input.Position,
		IsPublished: boolVal(input.IsPublished, true),
	}
	if err := s.db.Create(&class).Error; err != nil {
		return err
	}
	class.Slug = resolveSlug(input.Title, input.Slug, class.ID)
	return s.db.Model(&class).Update("slug", class.Slug).Error
}

func (s *service) UpdateClass(id uint64, input UpdateClassInput, tenantID uint) error {
	var class models.AcademicNoteClass
	if err := s.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&class).Error; err != nil {
		return err
	}
	slug := resolveSlug(input.Title, input.Slug, class.ID)
	return s.db.Model(&class).Updates(map[string]interface{}{
		"title":        input.Title,
		"slug":         slug,
		"icon_label":   input.IconLabel,
		"icon_color":   input.IconColor,
		"position":     input.Position,
		"is_published": boolVal(input.IsPublished, class.IsPublished),
	}).Error
}

func (s *service) DeleteClass(id uint64, tenantID uint) error {
	var class models.AcademicNoteClass
	if err := s.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&class).Error; err != nil {
		return err
	}
	return s.db.Delete(&class).Error
}

func (s *service) CreateSubject(input CreateSubjectInput, tenantID uint) error {
	if err := s.verifyClassOwnership(input.ClassID, tenantID); err != nil {
		return err
	}
	subject := models.AcademicNoteSubject{
		ClassID:     input.ClassID,
		Title:       input.Title,
		Slug:        utils.Slugify(input.Title),
		Position:    input.Position,
		IsPublished: boolVal(input.IsPublished, true),
	}
	if err := s.db.Create(&subject).Error; err != nil {
		return err
	}
	subject.Slug = resolveSlug(input.Title, input.Slug, subject.ID)
	return s.db.Model(&subject).Update("slug", subject.Slug).Error
}

func (s *service) UpdateSubject(id uint64, input UpdateSubjectInput, tenantID uint) error {
	subject, err := s.verifySubjectOwnership(uint(id), tenantID)
	if err != nil {
		return err
	}
	slug := resolveSlug(input.Title, input.Slug, subject.ID)
	return s.db.Model(subject).Updates(map[string]interface{}{
		"title":        input.Title,
		"slug":         slug,
		"position":     input.Position,
		"is_published": boolVal(input.IsPublished, subject.IsPublished),
	}).Error
}

func (s *service) DeleteSubject(id uint64, tenantID uint) error {
	subject, err := s.verifySubjectOwnership(uint(id), tenantID)
	if err != nil {
		return err
	}
	return s.db.Delete(subject).Error
}

func (s *service) CreatePaper(input CreatePaperInput, tenantID uint) error {
	if _, err := s.verifySubjectOwnership(input.SubjectID, tenantID); err != nil {
		return err
	}
	paper := models.AcademicNotePaper{
		SubjectID:   input.SubjectID,
		Title:       input.Title,
		Slug:        utils.Slugify(input.Title),
		IconLabel:   input.IconLabel,
		IconColor:   input.IconColor,
		Position:    input.Position,
		IsPublished: boolVal(input.IsPublished, true),
	}
	if err := s.db.Create(&paper).Error; err != nil {
		return err
	}
	paper.Slug = resolveSlug(input.Title, input.Slug, paper.ID)
	return s.db.Model(&paper).Update("slug", paper.Slug).Error
}

func (s *service) UpdatePaper(id uint64, input UpdatePaperInput, tenantID uint) error {
	paper, err := s.verifyPaperOwnership(uint(id), tenantID)
	if err != nil {
		return err
	}
	slug := resolveSlug(input.Title, input.Slug, paper.ID)
	return s.db.Model(paper).Updates(map[string]interface{}{
		"title":        input.Title,
		"slug":         slug,
		"icon_label":   input.IconLabel,
		"icon_color":   input.IconColor,
		"position":     input.Position,
		"is_published": boolVal(input.IsPublished, paper.IsPublished),
	}).Error
}

func (s *service) DeletePaper(id uint64, tenantID uint) error {
	paper, err := s.verifyPaperOwnership(uint(id), tenantID)
	if err != nil {
		return err
	}
	return s.db.Delete(paper).Error
}

func (s *service) CreateNote(input CreateNoteInput, tenantID uint) error {
	if _, err := s.verifyPaperOwnership(input.PaperID, tenantID); err != nil {
		return err
	}
	if input.PdfURL == "" {
		return errors.New("pdf file is required")
	}
	note := models.AcademicNote{
		PaperID:     input.PaperID,
		Title:       input.Title,
		Subtitle:    input.Subtitle,
		Thumbnail:   input.Thumbnail,
		PdfURL:      input.PdfURL,
		PdfFileName: input.PdfFileName,
		Position:    input.Position,
		IsPublished: boolVal(input.IsPublished, true),
	}
	return s.db.Create(&note).Error
}

func (s *service) UpdateNote(id uint64, input UpdateNoteInput, tenantID uint) error {
	note, err := s.verifyNoteOwnership(uint(id), tenantID)
	if err != nil {
		return err
	}
	updates := map[string]interface{}{
		"title":        input.Title,
		"subtitle":     input.Subtitle,
		"position":     input.Position,
		"is_published": boolVal(input.IsPublished, note.IsPublished),
	}
	if input.Thumbnail != nil && *input.Thumbnail != "" {
		updates["thumbnail"] = input.Thumbnail
	}
	if input.PdfURL != nil && *input.PdfURL != "" {
		updates["pdf_url"] = *input.PdfURL
	}
	if input.PdfFileName != nil {
		updates["pdf_file_name"] = input.PdfFileName
	}
	return s.db.Model(note).Updates(updates).Error
}

func (s *service) DeleteNote(id uint64, tenantID uint) error {
	note, err := s.verifyNoteOwnership(uint(id), tenantID)
	if err != nil {
		return err
	}
	if note.Thumbnail != nil && *note.Thumbnail != "" {
		if err := utils.DeleteFromBunny(*note.Thumbnail); err != nil {
			fmt.Println("Error deleting thumbnail:", err)
		}
	}
	if note.PdfURL != "" {
		if err := utils.DeleteFromBunny(note.PdfURL); err != nil {
			fmt.Println("Error deleting pdf:", err)
		}
	}
	return s.db.Delete(note).Error
}

func (s *service) GetPublicClasses(tenantID uint) ([]response.AcademicNoteClassResponse, error) {
	var classes []models.AcademicNoteClass
	if err := s.db.Where("tenant_id = ? AND is_published = ?", tenantID, true).
		Order("position ASC, id ASC").
		Find(&classes).Error; err != nil {
		return nil, err
	}
	result := make([]response.AcademicNoteClassResponse, 0, len(classes))
	for _, c := range classes {
		count, _ := s.countNotesForClass(c.ID, true)
		if count == 0 {
			continue
		}
		result = append(result, response.AcademicNoteClassResponse{
			ID:        c.ID,
			Title:     c.Title,
			Slug:      c.Slug,
			IconLabel: c.IconLabel,
			IconColor: c.IconColor,
			Position:  c.Position,
			NoteCount: count,
		})
	}
	return result, nil
}

func (s *service) GetPublicClassBySlug(tenantID uint, classSlug string) (*response.AcademicNoteClassDetailResponse, error) {
	var class models.AcademicNoteClass
	if err := s.db.
		Where("tenant_id = ? AND slug = ? AND is_published = ?", tenantID, classSlug, true).
		Preload("Subjects", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_published = ?", true).Order("position ASC, id ASC")
		}).
		Preload("Subjects.Papers", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_published = ?", true).Order("position ASC, id ASC")
		}).
		First(&class).Error; err != nil {
		return nil, err
	}

	subjects := make([]response.AcademicNoteSubjectResponse, 0, len(class.Subjects))
	for _, sub := range class.Subjects {
		papers := make([]response.AcademicNotePaperResponse, 0, len(sub.Papers))
		for _, paper := range sub.Papers {
			count, _ := s.countNotesForPaper(paper.ID, true)
			if count == 0 {
				continue
			}
			papers = append(papers, response.AcademicNotePaperResponse{
				ID:        paper.ID,
				SubjectID: paper.SubjectID,
				Title:     paper.Title,
				Slug:      paper.Slug,
				IconLabel: paper.IconLabel,
				IconColor: paper.IconColor,
				Position:  paper.Position,
				NoteCount: count,
			})
		}
		if len(papers) == 0 {
			continue
		}
		subCount, _ := s.countNotesForSubject(sub.ID, true)
		subjects = append(subjects, response.AcademicNoteSubjectResponse{
			ID:        sub.ID,
			ClassID:   sub.ClassID,
			Title:     sub.Title,
			Slug:      sub.Slug,
			Position:  sub.Position,
			NoteCount: subCount,
			Papers:    papers,
		})
	}

	return &response.AcademicNoteClassDetailResponse{
		ID:        class.ID,
		Title:     class.Title,
		Slug:      class.Slug,
		IconLabel: class.IconLabel,
		IconColor: class.IconColor,
		Position:  class.Position,
		Subjects:  subjects,
	}, nil
}

func (s *service) GetPublicNotesByPaperSlug(tenantID uint, classSlug, subjectSlug, paperSlug string) (*response.AcademicNotePaperDetailResponse, error) {
	var class models.AcademicNoteClass
	if err := s.db.Where("tenant_id = ? AND slug = ? AND is_published = ?", tenantID, classSlug, true).First(&class).Error; err != nil {
		return nil, errors.New("class not found")
	}

	var subject models.AcademicNoteSubject
	if err := s.db.Where("class_id = ? AND slug = ? AND is_published = ?", class.ID, subjectSlug, true).First(&subject).Error; err != nil {
		return nil, errors.New("subject not found")
	}

	var paper models.AcademicNotePaper
	if err := s.db.Where("subject_id = ? AND slug = ? AND is_published = ?", subject.ID, paperSlug, true).First(&paper).Error; err != nil {
		return nil, errors.New("paper not found")
	}

	var notes []models.AcademicNote
	if err := s.db.Where("paper_id = ? AND is_published = ?", paper.ID, true).
		Order("position ASC, id ASC").
		Find(&notes).Error; err != nil {
		return nil, err
	}

	noteItems := make([]response.AcademicNoteItemResponse, 0, len(notes))
	for _, n := range notes {
		noteItems = append(noteItems, toNoteItemResponse(n))
	}

	subCount, _ := s.countNotesForSubject(subject.ID, true)
	paperCount, _ := s.countNotesForPaper(paper.ID, true)

	return &response.AcademicNotePaperDetailResponse{
		Class: response.AcademicNoteClassResponse{
			ID:        class.ID,
			Title:     class.Title,
			Slug:      class.Slug,
			IconLabel: class.IconLabel,
			IconColor: class.IconColor,
		},
		Subject: response.AcademicNoteSubjectResponse{
			ID:        subject.ID,
			ClassID:   subject.ClassID,
			Title:     subject.Title,
			Slug:      subject.Slug,
			NoteCount: subCount,
		},
		Paper: response.AcademicNotePaperResponse{
			ID:        paper.ID,
			SubjectID: paper.SubjectID,
			Title:     paper.Title,
			Slug:      paper.Slug,
			IconLabel: paper.IconLabel,
			IconColor: paper.IconColor,
			NoteCount: paperCount,
		},
		Notes: noteItems,
	}, nil
}
