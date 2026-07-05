package academicnote

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"dashlearn/internal/middleware"
	"dashlearn/internal/models"
	"dashlearn/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func setupHandlerTest(t *testing.T) (*gin.Engine, uint, string) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	secret := "test-jwt-secret-for-academic-notes"
	os.Setenv("JWT_SECRET", secret)
	middleware.SecretKey = []byte(secret)

	db, tenantID := setupTestDB(t)
	utils.DB = db

	user := models.User{
		UserID:   "admin-test-user",
		Name:     "Admin",
		Email:    "admin@test.dev",
		Password: "hash",
		Status:   true,
		TenantID: tenantID,
	}
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		phone TEXT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		status INTEGER DEFAULT 1,
		tenant_id INTEGER NOT NULL,
		role_id INTEGER,
		created_at DATETIME,
		updated_at DATETIME
	)`).Error; err != nil {
		t.Fatalf("users schema: %v", err)
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	r := gin.New()
	v1 := r.Group("/v1")
	RegisterRoutes(v1)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "admin-test-user",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	return r, tenantID, signed
}

func TestHTTPPublicAndPrivateRoutes(t *testing.T) {
	r, tenantID, adminToken := setupHandlerTest(t)
	svc := NewService(utils.DB)

	published := true
	if err := svc.CreateClass(CreateClassInput{
		Title: "HSC", Slug: "hsc", IsPublished: &published,
	}, tenantID); err != nil {
		t.Fatalf("seed class: %v", err)
	}
	classes, _ := svc.GetAllClasses(tenantID)
	_ = svc.CreateSubject(CreateSubjectInput{
		ClassID: classes[0].ID, Title: "Bangla", Slug: "bangla", IsPublished: &published,
	}, tenantID)
	adminClass, _ := svc.GetClassByID(tenantID, uint64(classes[0].ID))
	_ = svc.CreatePaper(CreatePaperInput{
		SubjectID: adminClass.Subjects[0].ID, Title: "1st Paper", Slug: "1st-paper", IsPublished: &published,
	}, tenantID)
	adminClass, _ = svc.GetClassByID(tenantID, uint64(classes[0].ID))
	subtitle := "sub"
	_ = svc.CreateNote(CreateNoteInput{
		PaperID: adminClass.Subjects[0].Papers[0].ID,
		Title:   "Note One", Subtitle: &subtitle,
		PdfURL: "https://example.com/a.pdf", IsPublished: &published,
	}, tenantID)

	req := httptest.NewRequest(http.MethodGet, "/v1/academic-notes", nil)
	req.Header.Set("app-key", "test-key")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("public list status=%d body=%s", w.Code, w.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/v1/academic-notes/hsc/bangla/1st-paper", nil)
	req.Header.Set("app-key", "test-key")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("public notes status=%d body=%s", w.Code, w.Body.String())
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Note One")) {
		t.Fatalf("missing note in response: %s", w.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/v1/private/academic-notes/classes", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("admin list status=%d body=%s", w.Code, w.Body.String())
	}

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	_ = mw.WriteField("title", "Class 8")
	_ = mw.WriteField("slug", "class-8")
	_ = mw.WriteField("is_published", "true")
	_ = mw.Close()

	req = httptest.NewRequest(http.MethodPost, "/v1/private/academic-notes/classes/create", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create class status=%d body=%s", w.Code, w.Body.String())
	}

	var listResp struct {
		Data []models.AcademicNoteClass `json:"data"`
	}
	req = httptest.NewRequest(http.MethodGet, "/v1/private/academic-notes/classes", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	_ = json.Unmarshal(w.Body.Bytes(), &listResp)
	if len(listResp.Data) < 2 {
		t.Fatalf("expected >=2 classes after create, got %d", len(listResp.Data))
	}
}
