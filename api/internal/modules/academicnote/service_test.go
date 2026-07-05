package academicnote

import (
	"testing"

	"dashlearn/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, uint) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	stmts := []string{
		`CREATE TABLE IF NOT EXISTS tenants (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			app_key TEXT NOT NULL UNIQUE,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS academic_note_classes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tenant_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			slug TEXT NOT NULL,
			icon_label TEXT,
			icon_color TEXT,
			icon_image TEXT,
			position INTEGER DEFAULT 0,
			is_published INTEGER DEFAULT 1,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS academic_note_subjects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			class_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			slug TEXT NOT NULL,
			position INTEGER DEFAULT 0,
			is_published INTEGER DEFAULT 1,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS academic_note_papers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			subject_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			slug TEXT NOT NULL,
			icon_label TEXT,
			icon_color TEXT,
			position INTEGER DEFAULT 0,
			is_published INTEGER DEFAULT 1,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS academic_notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			paper_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			subtitle TEXT,
			thumbnail TEXT,
			pdf_url TEXT NOT NULL,
			pdf_file_name TEXT,
			position INTEGER DEFAULT 0,
			is_published INTEGER DEFAULT 1,
			created_at DATETIME,
			updated_at DATETIME
		)`,
	}
	for _, stmt := range stmts {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("exec schema: %v", err)
		}
	}

	tenant := models.Tenant{AppKey: "test-key"}
	if err := db.Create(&tenant).Error; err != nil {
		t.Fatalf("create tenant: %v", err)
	}
	return db, tenant.ID
}

func TestAcademicNoteHierarchyAndPublicAPI(t *testing.T) {
	db, tenantID := setupTestDB(t)
	svc := NewService(db)

	published := true
	if _, err := svc.CreateClass(CreateClassInput{
		Title:       "HSC",
		Slug:        "hsc",
		IconLabel:   strPtr("H"),
		IconColor:   strPtr("#E91E63"),
		IsPublished: &published,
	}, tenantID); err != nil {
		t.Fatalf("CreateClass: %v", err)
	}

	classes, err := svc.GetAllClasses(tenantID)
	if err != nil || len(classes) != 1 {
		t.Fatalf("GetAllClasses: err=%v len=%d", err, len(classes))
	}
	classID := classes[0].ID

	if err := svc.CreateSubject(CreateSubjectInput{
		ClassID:     classID,
		Title:       "Bangla",
		Slug:        "bangla",
		IsPublished: &published,
	}, tenantID); err != nil {
		t.Fatalf("CreateSubject: %v", err)
	}

	adminClass, err := svc.GetClassByID(tenantID, uint64(classID))
	if err != nil || len(adminClass.Subjects) != 1 {
		t.Fatalf("GetClassByID: err=%v subjects=%d", err, len(adminClass.Subjects))
	}
	subjectID := adminClass.Subjects[0].ID

	if err := svc.CreatePaper(CreatePaperInput{
		SubjectID:   subjectID,
		Title:       "Bangla 1st Paper",
		Slug:        "bangla-1st-paper",
		IconLabel:   strPtr("1"),
		IconColor:   strPtr("#42A5F5"),
		IsPublished: &published,
	}, tenantID); err != nil {
		t.Fatalf("CreatePaper: %v", err)
	}

	adminClass, err = svc.GetClassByID(tenantID, uint64(classID))
	if err != nil || len(adminClass.Subjects[0].Papers) != 1 {
		t.Fatalf("GetClassByID papers: err=%v papers=%d", err, len(adminClass.Subjects[0].Papers))
	}
	paperID := adminClass.Subjects[0].Papers[0].ID

	subtitle := "Oporichita"
	if err := svc.CreateNote(CreateNoteInput{
		PaperID:     paperID,
		Title:       "Oporichita",
		Subtitle:    &subtitle,
		PdfURL:      "https://cdn.example.com/oporichita.pdf",
		PdfFileName: strPtr("oporichita.pdf"),
		IsPublished: &published,
	}, tenantID); err != nil {
		t.Fatalf("CreateNote: %v", err)
	}

	publicClasses, err := svc.GetPublicClasses(tenantID)
	if err != nil || len(publicClasses) != 1 {
		t.Fatalf("GetPublicClasses: err=%v len=%d", err, len(publicClasses))
	}
	if publicClasses[0].NoteCount != 1 {
		t.Fatalf("expected note_count=1, got %d", publicClasses[0].NoteCount)
	}

	publicClass, err := svc.GetPublicClassBySlug(tenantID, "hsc")
	if err != nil {
		t.Fatalf("GetPublicClassBySlug: %v", err)
	}
	if len(publicClass.Subjects) != 1 || len(publicClass.Subjects[0].Papers) != 1 {
		t.Fatalf("public class tree incomplete")
	}
	if publicClass.Subjects[0].Papers[0].NoteCount != 1 {
		t.Fatalf("paper note_count expected 1, got %d", publicClass.Subjects[0].Papers[0].NoteCount)
	}

	publicNotes, err := svc.GetPublicNotesByPaperSlug(tenantID, "hsc", "bangla", "bangla-1st-paper")
	if err != nil {
		t.Fatalf("GetPublicNotesByPaperSlug: %v", err)
	}
	if len(publicNotes.Notes) != 1 || publicNotes.Notes[0].Title != "Oporichita" {
		t.Fatalf("unexpected notes: %+v", publicNotes.Notes)
	}

	unpublished := false
	if err := svc.UpdateClass(uint64(classID), UpdateClassInput{
		Title:       "HSC",
		Slug:        "hsc",
		IsPublished: &unpublished,
	}, tenantID); err != nil {
		t.Fatalf("UpdateClass: %v", err)
	}

	publicClasses, err = svc.GetPublicClasses(tenantID)
	if err != nil || len(publicClasses) != 0 {
		t.Fatalf("unpublished class should be hidden: len=%d err=%v", len(publicClasses), err)
	}

	if err := svc.DeleteClass(uint64(classID), tenantID); err != nil {
		t.Fatalf("DeleteClass: %v", err)
	}
	classes, err = svc.GetAllClasses(tenantID)
	if err != nil || len(classes) != 0 {
		t.Fatalf("DeleteClass failed: len=%d err=%v", len(classes), err)
	}
}

func TestCreateNoteRequiresPdf(t *testing.T) {
	db, tenantID := setupTestDB(t)
	svc := NewService(db)

	published := true
	_, _ = svc.CreateClass(CreateClassInput{Title: "C", IsPublished: &published}, tenantID)
	classes, _ := svc.GetAllClasses(tenantID)
	_ = svc.CreateSubject(CreateSubjectInput{ClassID: classes[0].ID, Title: "S", IsPublished: &published}, tenantID)
	adminClass, _ := svc.GetClassByID(tenantID, uint64(classes[0].ID))
	_ = svc.CreatePaper(CreatePaperInput{
		SubjectID: adminClass.Subjects[0].ID, Title: "P", IsPublished: &published,
	}, tenantID)
	adminClass, _ = svc.GetClassByID(tenantID, uint64(classes[0].ID))
	paperID := adminClass.Subjects[0].Papers[0].ID

	err := svc.CreateNote(CreateNoteInput{PaperID: paperID, Title: "No PDF"}, tenantID)
	if err == nil {
		t.Fatal("expected error when pdf_url missing")
	}
}

func strPtr(s string) *string { return &s }
