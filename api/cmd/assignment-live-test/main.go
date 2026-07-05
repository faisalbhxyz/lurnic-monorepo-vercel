package main

import (
	"bytes"
	"encoding/json"
	"errors"
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
	"github.com/lucsky/cuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

	course, assignment, err := ensurePublishedAssignment(tenantID)
	if err != nil {
		log.Fatalf("assignment fixture: %v", err)
	}

	student, err := ensureEnrolledStudent(tenantID, course.ID)
	if err != nil {
		log.Fatalf("student fixture: %v", err)
	}

	var oldSubmissionIDs []uint
	utils.DB.Model(&models.AssignmentSubmission{}).
		Where("assignment_id = ? AND student_id = ?", assignment.ID, student.ID).
		Pluck("id", &oldSubmissionIDs)
	if len(oldSubmissionIDs) > 0 {
		_ = utils.DB.Where("submission_id IN ?", oldSubmissionIDs).Delete(&models.AssignmentSubmissionFile{})
		_ = utils.DB.Where("id IN ?", oldSubmissionIDs).Delete(&models.AssignmentSubmission{})
	}

	engine, flush, err := server.NewEngine("assignment-live-test")
	if err != nil {
		log.Fatalf("engine: %v", err)
	}
	if flush != nil {
		defer flush(2 * time.Second)
	}

	ts := httptest.NewServer(engine)
	defer ts.Close()
	base := ts.URL + "/v1"

	fmt.Println("== Assignment live test ==")
	fmt.Printf("API: %s\n", base)
	fmt.Printf("Course slug: %s (id=%d)\n", course.Slug, course.ID)
	fmt.Printf("Assignment: %s (id=%d)\n", assignment.Title, assignment.ID)
	fmt.Printf("Student: %s\n", student.Email)

	studentToken, err := loginStudent(base, appKey, student.Email, "password123")
	if err != nil {
		log.Fatalf("student login: %v", err)
	}
	fmt.Println("OK student login")

	getURL := fmt.Sprintf("%s/course/%s/assignments/%d", base, course.Slug, assignment.ID)
	getRes, err := requestJSON(http.MethodGet, getURL, appKey, studentToken, nil)
	if err != nil {
		log.Fatalf("get assignment: %v", err)
	}
	if getRes.StatusCode != http.StatusOK {
		log.Fatalf("get assignment status %d: %s", getRes.StatusCode, string(getRes.Body))
	}
	fmt.Println("OK GET student assignment")

	submitURL := fmt.Sprintf("%s/course/%s/assignments/%d/submit", base, course.Slug, assignment.ID)
	var submitRes *httpResult
	if r2Configured() {
		submitRes, err = submitAssignment(submitURL, appKey, studentToken, "<p>Live test answer</p>", "answer.txt", []byte("assignment live test file"))
	} else {
		fmt.Println("NOTE: R2 not configured — submitting text only")
		submitRes, err = submitTextOnly(submitURL, appKey, studentToken, "<p>Live test answer</p>")
	}
	if err != nil {
		log.Fatalf("submit assignment: %v", err)
	}
	if submitRes.StatusCode != http.StatusCreated {
		log.Fatalf("submit status %d: %s", submitRes.StatusCode, string(submitRes.Body))
	}

	var submitPayload struct {
		Data struct {
			ID     uint   `json:"id"`
			Status string `json:"status"`
		} `json:"data"`
	}
	if err := json.Unmarshal(submitRes.Body, &submitPayload); err != nil {
		log.Fatalf("parse submit response: %v", err)
	}
	if submitPayload.Data.ID == 0 || submitPayload.Data.Status != "pending_review" {
		log.Fatalf("unexpected submit payload: %s", string(submitRes.Body))
	}
	fmt.Printf("OK POST submit (submission id=%d)\n", submitPayload.Data.ID)

	dupRes, err := submitAssignment(submitURL, appKey, studentToken, "again", "x.txt", []byte("x"))
	if err != nil {
		log.Fatalf("duplicate submit request: %v", err)
	}
	if dupRes.StatusCode != http.StatusBadRequest {
		log.Fatalf("expected duplicate submit 400, got %d: %s", dupRes.StatusCode, string(dupRes.Body))
	}
	fmt.Println("OK duplicate submit blocked")

	adminToken, err := loginAdmin(base, "admin@local.dev", "password123")
	if err != nil {
		log.Fatalf("admin login: %v", err)
	}

	listURL := fmt.Sprintf("%s/private/course/%d/assignment-submissions", base, course.ID)
	listRes, err := requestJSON(http.MethodGet, listURL, "", adminToken, nil)
	if err != nil {
		log.Fatalf("admin list: %v", err)
	}
	if listRes.StatusCode != http.StatusOK {
		log.Fatalf("admin list status %d: %s", listRes.StatusCode, string(listRes.Body))
	}
	fmt.Println("OK admin list submissions")

	gradeURL := fmt.Sprintf("%s/private/course/%d/assignment-submissions/%d/grade", base, course.ID, submitPayload.Data.ID)
	gradeBody, _ := json.Marshal(map[string]any{
		"score":    8,
		"feedback": "Live test grade",
	})
	gradeRes, err := requestJSON(http.MethodPost, gradeURL, "", adminToken, gradeBody)
	if err != nil {
		log.Fatalf("grade: %v", err)
	}
	if gradeRes.StatusCode != http.StatusOK {
		log.Fatalf("grade status %d: %s", gradeRes.StatusCode, string(gradeRes.Body))
	}
	fmt.Println("OK admin grade submission")

	historyURL := base + "/student/assignment-submissions?course_id=" + fmt.Sprint(course.ID)
	historyRes, err := requestJSON(http.MethodGet, historyURL, appKey, studentToken, nil)
	if err != nil {
		log.Fatalf("student history: %v", err)
	}
	if historyRes.StatusCode != http.StatusOK {
		log.Fatalf("student history status %d: %s", historyRes.StatusCode, string(historyRes.Body))
	}
	fmt.Println("OK student submission history")

	fmt.Println("\nAll assignment live checks passed.")
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
		return fmt.Errorf("GOOSE_DBSTRING not set (check repo .env)")
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
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if filepath.Base(cwd) == "api" {
		return cwd
	}
	return filepath.Join(cwd, "api")
}

func resolveTenant() (appKey string, tenantID uint, err error) {
	appKey = strings.TrimSpace(os.Getenv("SEED_APP_KEY"))
	if appKey == "" {
		appKey = "local-dev"
	}
	var tenant models.Tenant
	if err := utils.DB.Where("app_key = ?", appKey).First(&tenant).Error; err != nil {
		return "", 0, fmt.Errorf("tenant with app_key %q not found (run: go run ./cmd/seed)", appKey)
	}
	return appKey, tenant.ID, nil
}

func ensurePublishedAssignment(tenantID uint) (*models.CourseDetails, *models.CourseAssignment, error) {
	var assignment models.CourseAssignment
	err := utils.DB.Where("is_published = ? AND course_id IN (SELECT id FROM course_details WHERE tenant_id = ?)", true, tenantID).
		Order("id DESC").First(&assignment).Error
	if err == nil {
		var course models.CourseDetails
		if err := utils.DB.First(&course, assignment.CourseID).Error; err != nil {
			return nil, nil, err
		}
		return &course, &assignment, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, err
	}

	var course models.CourseDetails
	if err := utils.DB.Where("tenant_id = ?", tenantID).Order("id DESC").First(&course).Error; err != nil {
		return nil, nil, fmt.Errorf("no course found for tenant %d", tenantID)
	}

	var chapter models.CourseChapter
	if err := utils.DB.Where("course_id = ?", course.ID).Order("id ASC").First(&chapter).Error; err != nil {
		chapter = models.CourseChapter{
			Title:    "Live Test Chapter",
			Access:   models.Published,
			CourseID: course.ID,
			Position: 0,
		}
		if err := utils.DB.Create(&chapter).Error; err != nil {
			return nil, nil, err
		}
	}

	assignment = models.CourseAssignment{
		CourseID:         course.ID,
		ChapterID:        chapter.ID,
		Title:            "Live Test Assignment",
		Instructions:     "<p>Submit anything for live test</p>",
		IsPublished:      true,
		TimeLimit:        1,
		TimeLimitOption:  models.CourseAssignmentTimeLimitOptionWeek,
		FileUploadLimit:  2,
		TotalMarks:       10,
		MinimumPassMarks: 6,
	}
	if err := utils.DB.Create(&assignment).Error; err != nil {
		return nil, nil, err
	}
	return &course, &assignment, nil
}

func ensureEnrolledStudent(tenantID, courseID uint) (*models.Student, error) {
	email := "student+1@local.dev"
	var student models.Student
	err := utils.DB.Where("email = ? AND tenant_id = ?", email, tenantID).First(&student).Error
	if err == gorm.ErrRecordNotFound {
		pw, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		student = models.Student{
			UserID:    cuid.New(),
			FirstName: "Live",
			LastName:  ptr("Tester"),
			Email:     email,
			Password:  string(pw),
			Status:    true,
			TenantID:  tenantID,
		}
		if err := utils.DB.Create(&student).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	var count int64
	utils.DB.Model(&models.Enrollment{}).
		Where("tenant_id = ? AND student_id = ? AND course_id = ?", tenantID, student.ID, courseID).
		Count(&count)
	if count == 0 {
		if err := utils.DB.Create(&models.Enrollment{
			TenantID:  tenantID,
			StudentID: student.ID,
			CourseID:  courseID,
		}).Error; err != nil {
			return nil, err
		}
	}
	return &student, nil
}

func loginStudent(base, appKey, email, password string) (string, error) {
	body, _ := json.Marshal(map[string]string{"email": email, "password": password})
	res, err := requestJSON(http.MethodPost, base+"/student/login", appKey, "", body)
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

func submitAssignment(url, appKey, bearer, responseText, fileName string, fileBody []byte) (*httpResult, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	_ = writer.WriteField("response_text", responseText)
	part, err := writer.CreateFormFile("files", fileName)
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(fileBody); err != nil {
		return nil, err
	}
	_ = writer.Close()

	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("app-key", appKey)
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

func submitTextOnly(url, appKey, bearer, responseText string) (*httpResult, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	_ = writer.WriteField("response_text", responseText)
	_ = writer.Close()

	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("app-key", appKey)
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
