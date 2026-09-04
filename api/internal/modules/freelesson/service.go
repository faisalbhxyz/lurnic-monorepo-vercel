package freelesson

import (
	"dashlearn/internal/models"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	defaultCatalogLimit = 50
	maxCatalogLimit     = 100
	maxPerRequest       = 3
	maxLibraryTotal     = 20
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// ListCatalog returns published + is_public lessons on public courses.
func (s *Service) ListCatalog(tenantID uint, classSlug string, limit, offset int) ([]CatalogItem, CatalogMeta, error) {
	if limit <= 0 {
		limit = defaultCatalogLimit
	}
	if limit > maxCatalogLimit {
		limit = maxCatalogLimit
	}
	if offset < 0 {
		offset = 0
	}

	base := s.db.Table("course_lessons AS l").
		Joins("JOIN course_chapters AS ch ON ch.id = l.chapter_id AND ch.access = ?", models.Published).
		Joins("JOIN course_details AS c ON c.id = ch.course_id").
		Joins("LEFT JOIN course_general_settings AS gs ON gs.course_id = c.id").
		Joins("LEFT JOIN sub_categories AS sc ON sc.id = gs.sub_category_id").
		Joins("LEFT JOIN categories AS cat ON cat.id = gs.category_id").
		Where("c.tenant_id = ?", tenantID).
		Where("c.visibility = ?", models.Public).
		Where("l.is_published = ?", true).
		Where("l.is_public = ?", true)

	classSlug = strings.TrimSpace(classSlug)
	if classSlug != "" {
		base = base.Where("(sc.slug = ? OR cat.slug = ?)", classSlug, classSlug)
	}

	var total int64
	if err := base.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, CatalogMeta{}, fmt.Errorf("count free lessons: %w", err)
	}

	type row struct {
		LessonID        uint                     `gorm:"column:lesson_id"`
		LessonTitle     string                   `gorm:"column:lesson_title"`
		LessonType      models.LessonType        `gorm:"column:lesson_type"`
		SourceType      models.LessonSourceType  `gorm:"column:source_type"`
		SourceJSON      []byte                   `gorm:"column:source"`
		ChapterID       uint                     `gorm:"column:chapter_id"`
		ChapterTitle    string                   `gorm:"column:chapter_title"`
		CourseID        uint                     `gorm:"column:course_id"`
		CourseSlug      string                   `gorm:"column:course_slug"`
		CourseTitle     string                   `gorm:"column:course_title"`
		FeaturedImage   *string                  `gorm:"column:featured_image"`
		SubCategorySlug *string                  `gorm:"column:sub_category_slug"`
		CategorySlug    *string                  `gorm:"column:category_slug"`
	}

	var rows []row
	err := base.
		Select(`
			l.id AS lesson_id,
			l.title AS lesson_title,
			l.lesson_type,
			l.source_type,
			l.source,
			ch.id AS chapter_id,
			ch.title AS chapter_title,
			c.id AS course_id,
			c.slug AS course_slug,
			c.title AS course_title,
			c.featured_image,
			sc.slug AS sub_category_slug,
			cat.slug AS category_slug
		`).
		Order("l.id ASC").
		Limit(limit).
		Offset(offset).
		Scan(&rows).Error
	if err != nil {
		return nil, CatalogMeta{}, fmt.Errorf("list free lessons: %w", err)
	}

	items := make([]CatalogItem, 0, len(rows))
	for _, r := range rows {
		lesson := models.CourseLesson{}
		lesson.Source.Scan(r.SourceJSON)

		classSlugs := make([]string, 0, 2)
		if r.SubCategorySlug != nil && *r.SubCategorySlug != "" {
			classSlugs = append(classSlugs, *r.SubCategorySlug)
		}
		if r.CategorySlug != nil && *r.CategorySlug != "" {
			already := false
			for _, s := range classSlugs {
				if s == *r.CategorySlug {
					already = true
					break
				}
			}
			if !already {
				classSlugs = append(classSlugs, *r.CategorySlug)
			}
		}

		items = append(items, CatalogItem{
			LessonID:      r.LessonID,
			LessonTitle:   r.LessonTitle,
			ChapterID:     r.ChapterID,
			ChapterTitle:  r.ChapterTitle,
			CourseID:      r.CourseID,
			CourseSlug:    r.CourseSlug,
			CourseTitle:   r.CourseTitle,
			FeaturedImage: r.FeaturedImage,
			LessonType:    r.LessonType,
			SourceType:    r.SourceType,
			Source:        lesson.Source,
			IsPublic:      true,
			ClassSlugs:    classSlugs,
		})
	}

	return items, CatalogMeta{Total: total, Limit: limit, Offset: offset}, nil
}

// ListLibrary returns the student's saved free lessons with watch progress.
func (s *Service) ListLibrary(tenantID, studentID uint) ([]LibraryItem, error) {
	var saved []models.StudentFreeLesson
	if err := s.db.
		Where("tenant_id = ? AND student_id = ?", tenantID, studentID).
		Order("added_at DESC").
		Find(&saved).Error; err != nil {
		return nil, err
	}
	if len(saved) == 0 {
		return []LibraryItem{}, nil
	}

	lessonIDs := make([]uint, len(saved))
	for i, row := range saved {
		lessonIDs[i] = row.LessonID
	}

	type detailRow struct {
		LessonID      uint                    `gorm:"column:lesson_id"`
		LessonTitle   string                  `gorm:"column:lesson_title"`
		SourceType    models.LessonSourceType `gorm:"column:source_type"`
		SourceJSON    []byte                  `gorm:"column:source"`
		ChapterTitle  string                  `gorm:"column:chapter_title"`
		CourseID      uint                    `gorm:"column:course_id"`
		CourseSlug    string                  `gorm:"column:course_slug"`
		CourseTitle   string                  `gorm:"column:course_title"`
		FeaturedImage *string                 `gorm:"column:featured_image"`
	}

	var details []detailRow
	if err := s.db.Table("course_lessons AS l").
		Joins("JOIN course_chapters AS ch ON ch.id = l.chapter_id").
		Joins("JOIN course_details AS c ON c.id = ch.course_id").
		Where("l.id IN ?", lessonIDs).
		Select(`
			l.id AS lesson_id,
			l.title AS lesson_title,
			l.source_type,
			l.source,
			ch.title AS chapter_title,
			c.id AS course_id,
			c.slug AS course_slug,
			c.title AS course_title,
			c.featured_image
		`).
		Scan(&details).Error; err != nil {
		return nil, err
	}

	detailByID := make(map[uint]detailRow, len(details))
	for _, d := range details {
		detailByID[d.LessonID] = d
	}

	progressByLesson := map[uint]models.StudentLessonVideoProgress{}
	var progressRows []models.StudentLessonVideoProgress
	if err := s.db.
		Where("tenant_id = ? AND student_id = ? AND lesson_id IN ?", tenantID, studentID, lessonIDs).
		Find(&progressRows).Error; err != nil {
		return nil, err
	}
	for _, p := range progressRows {
		progressByLesson[p.LessonID] = p
	}

	completedSet := map[uint]bool{}
	var completions []models.StudentLessonCompletion
	if err := s.db.
		Where("tenant_id = ? AND student_id = ? AND lesson_id IN ?", tenantID, studentID, lessonIDs).
		Find(&completions).Error; err != nil {
		return nil, err
	}
	for _, c := range completions {
		completedSet[c.LessonID] = true
	}

	items := make([]LibraryItem, 0, len(saved))
	for _, row := range saved {
		d, ok := detailByID[row.LessonID]
		if !ok {
			continue
		}
		lesson := models.CourseLesson{}
		_ = lesson.Source.Scan(d.SourceJSON)

		prog := progressByLesson[row.LessonID]
		items = append(items, LibraryItem{
			LessonID:        row.LessonID,
			LessonTitle:     d.LessonTitle,
			ChapterTitle:    d.ChapterTitle,
			CourseID:        d.CourseID,
			CourseSlug:      d.CourseSlug,
			CourseTitle:     d.CourseTitle,
			FeaturedImage:   d.FeaturedImage,
			SourceType:      d.SourceType,
			Source:          lesson.Source,
			AddedAt:         row.AddedAt,
			WatchPercent:    prog.ProgressPercent,
			WatchSeconds:    prog.MaxPositionSeconds,
			DurationSeconds: prog.DurationSeconds,
			Completed:       completedSet[row.LessonID],
		})
	}
	return items, nil
}

// AddToLibrary adds free/public lessons to the student's library (max 3 per request, 20 total).
func (s *Service) AddToLibrary(tenantID, studentID uint, lessonIDs []uint) ([]LibraryItem, error) {
	if len(lessonIDs) == 0 {
		return nil, errors.New("lesson_ids is required")
	}
	if len(lessonIDs) > maxPerRequest {
		return nil, fmt.Errorf("max %d lesson_ids per request", maxPerRequest)
	}

	// Deduplicate while preserving order
	seen := map[uint]bool{}
	unique := make([]uint, 0, len(lessonIDs))
	for _, id := range lessonIDs {
		if id == 0 || seen[id] {
			continue
		}
		seen[id] = true
		unique = append(unique, id)
	}
	if len(unique) == 0 {
		return nil, errors.New("lesson_ids is required")
	}

	type eligible struct {
		LessonID uint `gorm:"column:lesson_id"`
		CourseID uint `gorm:"column:course_id"`
	}
	var eligibleRows []eligible
	if err := s.db.Table("course_lessons AS l").
		Joins("JOIN course_chapters AS ch ON ch.id = l.chapter_id AND ch.access = ?", models.Published).
		Joins("JOIN course_details AS c ON c.id = ch.course_id").
		Where("c.tenant_id = ?", tenantID).
		Where("c.visibility = ?", models.Public).
		Where("l.is_published = ?", true).
		Where("l.is_public = ?", true).
		Where("l.id IN ?", unique).
		Select("l.id AS lesson_id, c.id AS course_id").
		Scan(&eligibleRows).Error; err != nil {
		return nil, err
	}

	eligibleMap := map[uint]uint{}
	for _, e := range eligibleRows {
		eligibleMap[e.LessonID] = e.CourseID
	}
	for _, id := range unique {
		if _, ok := eligibleMap[id]; !ok {
			return nil, fmt.Errorf("lesson %d is not a free/public published lesson", id)
		}
	}

	var existingCount int64
	if err := s.db.Model(&models.StudentFreeLesson{}).
		Where("tenant_id = ? AND student_id = ?", tenantID, studentID).
		Count(&existingCount).Error; err != nil {
		return nil, err
	}

	var already []models.StudentFreeLesson
	if err := s.db.
		Where("tenant_id = ? AND student_id = ? AND lesson_id IN ?", tenantID, studentID, unique).
		Find(&already).Error; err != nil {
		return nil, err
	}
	alreadySet := map[uint]bool{}
	for _, a := range already {
		alreadySet[a.LessonID] = true
	}

	newCount := 0
	for _, id := range unique {
		if !alreadySet[id] {
			newCount++
		}
	}
	if int(existingCount)+newCount > maxLibraryTotal {
		return nil, fmt.Errorf("free lesson library limit is %d", maxLibraryTotal)
	}

	now := time.Now().UTC()
	for _, id := range unique {
		if alreadySet[id] {
			continue
		}
		row := models.StudentFreeLesson{
			TenantID:  tenantID,
			StudentID: studentID,
			LessonID:  id,
			CourseID:  eligibleMap[id],
			AddedAt:   now,
		}
		if err := s.db.Create(&row).Error; err != nil {
			return nil, err
		}
	}

	return s.ListLibrary(tenantID, studentID)
}

// RemoveFromLibrary deletes one lesson from the student's free library.
func (s *Service) RemoveFromLibrary(tenantID, studentID, lessonID uint) error {
	res := s.db.
		Where("tenant_id = ? AND student_id = ? AND lesson_id = ?", tenantID, studentID, lessonID).
		Delete(&models.StudentFreeLesson{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
