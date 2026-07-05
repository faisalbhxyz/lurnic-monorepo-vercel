package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"dashlearn/internal/models"
	"dashlearn/internal/utils"
	"dashlearn/pkg/server"

	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	if err := utils.ConnectDatabase(); err != nil {
		log.Fatalf("database connect failed: %v", err)
	}
	if err := runMigrations(); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	appKey, tenantID, err := resolveTenant()
	if err != nil {
		log.Fatalf("tenant: %v", err)
	}

	cleanupTestData(tenantID)

	engine, flush, err := server.NewEngine("academic-note-live-test")
	if err != nil {
		log.Fatalf("engine: %v", err)
	}
	if flush != nil {
		defer flush(2 * time.Second)
	}

	ts := httptest.NewServer(engine)
	defer ts.Close()
	base := ts.URL + "/v1"

	fmt.Println("== Academic note live test ==")
	fmt.Printf("API: %s\n", base)
	fmt.Printf("Tenant: %d app-key: %s\n", tenantID, appKey)

	adminToken, err := loginAdmin(base, "admin@local.dev", "password123")
	if err != nil {
		log.Fatalf("admin login: %v (run: go run ./cmd/seed)", err)
	}
	fmt.Println("OK admin login")

	classID, classSlug, err := createClass(base, adminToken)
	if err != nil {
		log.Fatalf("create class: %v", err)
	}
	fmt.Printf("OK create class id=%d slug=%s\n", classID, classSlug)

	subjectID, subjectSlug, err := createSubject(base, adminToken, classID)
	if err != nil {
		log.Fatalf("create subject: %v", err)
	}
	fmt.Printf("OK create subject id=%d slug=%s\n", subjectID, subjectSlug)

	paperID, paperSlug, err := createPaper(base, adminToken, subjectID)
	if err != nil {
		log.Fatalf("create paper: %v", err)
	}
	fmt.Printf("OK create paper id=%d slug=%s\n", paperID, paperSlug)

	if err := ensureNoteViaDB(paperID); err != nil {
		log.Fatalf("seed note: %v", err)
	}
	fmt.Println("OK seed note via DB")

	adminTree, err := requestJSON(http.MethodGet, fmt.Sprintf("%s/private/academic-notes/classes/%d", base, classID), "", adminToken, nil)
	if err != nil {
		log.Fatalf("admin get class: %v", err)
	}
	if adminTree.StatusCode != http.StatusOK {
		log.Fatalf("admin get class status %d: %s", adminTree.StatusCode, string(adminTree.Body))
	}
	if !strings.Contains(string(adminTree.Body), "Live Test Note") {
		log.Fatalf("admin tree missing note: %s", string(adminTree.Body))
	}
	fmt.Println("OK admin class tree")

	publicList, err := requestJSON(http.MethodGet, base+"/academic-notes", appKey, "", nil)
	if err != nil {
		log.Fatalf("public list: %v", err)
	}
	if publicList.StatusCode != http.StatusOK {
		log.Fatalf("public list status %d: %s", publicList.StatusCode, string(publicList.Body))
	}
	if !strings.Contains(string(publicList.Body), classSlug) {
		log.Fatalf("public list missing class: %s", string(publicList.Body))
	}
	fmt.Println("OK public class list")

	publicClass, err := requestJSON(http.MethodGet, base+"/academic-notes/"+classSlug, appKey, "", nil)
	if err != nil {
		log.Fatalf("public class: %v", err)
	}
	if publicClass.StatusCode != http.StatusOK {
		log.Fatalf("public class status %d: %s", publicClass.StatusCode, string(publicClass.Body))
	}
	if !strings.Contains(string(publicClass.Body), subjectSlug) || !strings.Contains(string(publicClass.Body), paperSlug) {
		log.Fatalf("public class missing subject/paper: %s", string(publicClass.Body))
	}
	fmt.Println("OK public class detail")

	publicNotes, err := requestJSON(http.MethodGet, fmt.Sprintf("%s/academic-notes/%s/%s/%s", base, classSlug, subjectSlug, paperSlug), appKey, "", nil)
	if err != nil {
		log.Fatalf("public notes: %v", err)
	}
	if publicNotes.StatusCode != http.StatusOK {
		log.Fatalf("public notes status %d: %s", publicNotes.StatusCode, string(publicNotes.Body))
	}
	if !strings.Contains(string(publicNotes.Body), "Live Test Note") {
		log.Fatalf("public notes missing item: %s", string(publicNotes.Body))
	}
	fmt.Println("OK public notes list")

	if r2Configured() {
		if err := testNoteUploadAPI(base, adminToken, paperID); err != nil {
			log.Fatalf("note upload API: %v", err)
		}
		fmt.Println("OK note upload via API (R2)")
	} else {
		fmt.Println("SKIP note upload API (R2 not configured)")
	}

	cleanupTestData(tenantID)
	fmt.Println("\nAll academic note live checks passed.")
}

func cleanupTestData(tenantID uint) {
	utils.DB.Where("tenant_id = ? AND title LIKE ?", tenantID, "Live Test Class%").Delete(&models.AcademicNoteClass{})
}

func createClass(base, token string) (uint, string, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("title", "Live Test Class HSC")
	_ = w.WriteField("slug", "live-test-hsc")
	_ = w.WriteField("icon_label", "H")
	_ = w.WriteField("icon_color", "#E91E63")
	_ = w.WriteField("position", "0")
	_ = w.WriteField("is_published", "true")
	_ = w.Close()

	res, err := postMultipart(base+"/private/academic-notes/classes/create", token, &buf, w.FormDataContentType())
	if err != nil {
		return 0, "", err
	}
	if res.StatusCode != http.StatusCreated {
		return 0, "", fmt.Errorf("status %d: %s", res.StatusCode, string(res.Body))
	}

	var classes []models.AcademicNoteClass
	if err := utils.DB.Where("slug = ?", "live-test-hsc").Find(&classes).Error; err != nil || len(classes) == 0 {
		return 0, "", fmt.Errorf("class not found in DB")
	}
	return classes[0].ID, classes[0].Slug, nil
}

func createSubject(base, token string, classID uint) (uint, string, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("class_id", fmt.Sprint(classID))
	_ = w.WriteField("title", "Live Test Bangla")
	_ = w.WriteField("slug", "live-test-bangla")
	_ = w.WriteField("position", "0")
	_ = w.WriteField("is_published", "true")
	_ = w.Close()

	res, err := postMultipart(base+"/private/academic-notes/subjects/create", token, &buf, w.FormDataContentType())
	if err != nil {
		return 0, "", err
	}
	if res.StatusCode != http.StatusCreated {
		return 0, "", fmt.Errorf("status %d: %s", res.StatusCode, string(res.Body))
	}

	var subject models.AcademicNoteSubject
	if err := utils.DB.Where("class_id = ? AND slug = ?", classID, "live-test-bangla").First(&subject).Error; err != nil {
		return 0, "", err
	}
	return subject.ID, subject.Slug, nil
}

func createPaper(base, token string, subjectID uint) (uint, string, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("subject_id", fmt.Sprint(subjectID))
	_ = w.WriteField("title", "Live Test 1st Paper")
	_ = w.WriteField("slug", "live-test-1st-paper")
	_ = w.WriteField("icon_label", "1")
	_ = w.WriteField("icon_color", "#42A5F5")
	_ = w.WriteField("position", "0")
	_ = w.WriteField("is_published", "true")
	_ = w.Close()

	res, err := postMultipart(base+"/private/academic-notes/papers/create", token, &buf, w.FormDataContentType())
	if err != nil {
		return 0, "", err
	}
	if res.StatusCode != http.StatusCreated {
		return 0, "", fmt.Errorf("status %d: %s", res.StatusCode, string(res.Body))
	}

	var paper models.AcademicNotePaper
	if err := utils.DB.Where("subject_id = ? AND slug = ?", subjectID, "live-test-1st-paper").First(&paper).Error; err != nil {
		return 0, "", err
	}
	return paper.ID, paper.Slug, nil
}

func ensureNoteViaDB(paperID uint) error {
	pdfURL := "https://example.com/live-test-note.pdf"
	subtitle := "Live subtitle"
	return utils.DB.Create(&models.AcademicNote{
		PaperID:     paperID,
		Title:       "Live Test Note",
		Subtitle:    &subtitle,
		PdfURL:      pdfURL,
		PdfFileName: ptr("live-test-note.pdf"),
		Position:    0,
		IsPublished: true,
	}).Error
}

func testNoteUploadAPI(base, token string, paperID uint) error {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("paper_id", fmt.Sprint(paperID))
	_ = w.WriteField("title", "Uploaded Note")
	_ = w.WriteField("position", "1")
	_ = w.WriteField("is_published", "true")
	part, err := w.CreateFormFile("pdf", "upload-test.pdf")
	if err != nil {
		return err
	}
	if _, err := part.Write([]byte("%PDF-1.4 live test")); err != nil {
		return err
	}
	_ = w.Close()

	res, err := postMultipart(base+"/private/academic-notes/notes/create", token, &buf, w.FormDataContentType())
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("status %d: %s", res.StatusCode, string(res.Body))
	}
	return nil
}

func loadEnv() {
	if cwd, err := os.Getwd(); err == nil {
		if filepath.Base(cwd) == "api" {
			_ = godotenv.Load(filepath.Join(cwd, "..", ".env"))
		}
		_ = godotenv.Load(filepath.Join(cwd, ".env"))
	}
	_ = godotenv.Load()
}

func runMigrations() error {
	dsn := os.Getenv("GOOSE_DBSTRING")
	if dsn == "" {
		return fmt.Errorf("GOOSE_DBSTRING not set")
	}
	cmd := exec.Command("go", "run", "github.com/pressly/goose/v3/cmd/goose@v3.26.0",
		"-dir", "migrations", "mysql", dsn, "up")
	cmd.Dir = mustAPIRoot()
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "no change") || strings.Contains(string(out), "OK") {
			return nil
		}
		return fmt.Errorf("%v: %s", err, out)
	}
	return nil
}

func mustAPIRoot() string {
	cwd, _ := os.Getwd()
	if filepath.Base(cwd) == "api" {
		return cwd
	}
	return filepath.Join(cwd, "api")
}

func resolveTenant() (string, uint, error) {
	appKey := strings.TrimSpace(os.Getenv("SEED_APP_KEY"))
	if appKey == "" {
		appKey = "local-dev"
	}
	var tenant models.Tenant
	if err := utils.DB.Where("app_key = ?", appKey).First(&tenant).Error; err != nil {
		return "", 0, fmt.Errorf("tenant %q not found (run: go run ./cmd/seed)", appKey)
	}
	return appKey, tenant.ID, nil
}

func loginAdmin(base, email, password string) (string, error) {
	body, _ := json.Marshal(map[string]string{"email": email, "password": password})
	res, err := requestJSON(http.MethodPost, base+"/user/login", "", "", body)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d: %s", res.StatusCode, string(res.Body))
	}
	var payload struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(res.Body, &payload); err != nil {
		return "", err
	}
	if payload.Token == "" {
		return "", fmt.Errorf("missing token")
	}
	return payload.Token, nil
}

type httpResult struct {
	StatusCode int
	Body       []byte
}

func requestJSON(method, url, appKey, bearer string, body []byte) (*httpResult, error) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if appKey != "" {
		req.Header.Set("app-key", appKey)
	}
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &httpResult{StatusCode: resp.StatusCode, Body: data}, nil
}

func postMultipart(url, bearer string, body *bytes.Buffer, contentType string) (*httpResult, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+bearer)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &httpResult{StatusCode: resp.StatusCode, Body: data}, nil
}

func r2Configured() bool {
	return os.Getenv("R2_ACCOUNT_ID") != "" &&
		os.Getenv("R2_ACCESS_KEY_ID") != "" &&
		os.Getenv("R2_SECRET_ACCESS_KEY") != ""
}

func ptr(s string) *string { return &s }
