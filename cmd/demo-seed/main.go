package main

import (
	"dashlearn/internal/models"
	"dashlearn/internal/utils"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/lucsky/cuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()

	if err := utils.ConnectDatabase(); err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	tenantID := envUint("DEMO_TENANT_ID", 1)

	author, err := findAuthor(tenantID)
	if err != nil {
		log.Fatalf("failed to find author user for demo seed: %v", err)
	}

	category := ensureCategory(tenantID)
	subCategory := ensureSubCategory(tenantID, category.ID)
	instructor := ensureInstructor(tenantID)
	course := ensureCourse(tenantID, author.ID)
	ensureCourseGeneralSettings(tenantID, course.ID, category.ID, &subCategory.ID)
	ensureCourseInstructor(course.ID, instructor.ID)

	student := ensureStudent(tenantID)
	ensureEnrollment(tenantID, student.ID, course.ID)
	ensureOrder(tenantID, student.ID, course.ID)

	ensurePaymentMethod(tenantID)
	ensureBanner(tenantID)
	ensureGeneralSettings(tenantID)

	fmt.Println("Demo seed complete.")
	fmt.Println("Created/ensured demo records for:")
	fmt.Println("- Category / Sub-category")
	fmt.Println("- Instructor")
	fmt.Println("- Course (+ settings + instructor mapping)")
	fmt.Println("- Student, Enrollment, Order")
	fmt.Println("- Payment method, Banner, General settings")
}

func envUint(key string, def uint) uint {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	var n uint
	_, err := fmt.Sscanf(v, "%d", &n)
	if err != nil || n == 0 {
		return def
	}
	return n
}

func findAuthor(tenantID uint) (*models.User, error) {
	email := os.Getenv("DEMO_AUTHOR_EMAIL")
	if email == "" {
		email = "admin@local.dev"
	}

	var u models.User
	err := utils.DB.Where("email = ? AND tenant_id = ?", email, tenantID).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("author user not found (email=%s tenant_id=%d). Run: go run ./cmd/seed", email, tenantID)
	}
	return &u, err
}

func ensureCategory(tenantID uint) *models.Category {
	desc := "Demo category for local development."
	thumb := "https://placehold.co/600x400/png?text=Category"
	c := models.Category{
		Name:        "Demo Category",
		Slug:        "demo-category",
		Description: &desc,
		Thumbnail:   &thumb,
		TenantID:    tenantID,
	}

	var existing models.Category
	if err := utils.DB.Where("slug = ? AND tenant_id = ?", c.Slug, tenantID).First(&existing).Error; err == nil {
		// keep it fresh-ish
		existing.Name = c.Name
		existing.Description = c.Description
		existing.Thumbnail = c.Thumbnail
		_ = utils.DB.Save(&existing).Error
		return &existing
	}

	if err := utils.DB.Create(&c).Error; err != nil {
		log.Fatalf("failed to create category: %v", err)
	}
	return &c
}

func ensureSubCategory(tenantID, categoryID uint) *models.SubCategory {
	desc := "Demo sub-category for local development."
	thumb := "https://placehold.co/600x400/png?text=Sub+Category"
	sc := models.SubCategory{
		CategoryID:  categoryID,
		Name:        "Demo Sub Category",
		Slug:        "demo-sub-category",
		Description: &desc,
		Thumbnail:   &thumb,
		TenantID:    tenantID,
	}

	var existing models.SubCategory
	if err := utils.DB.Where("slug = ? AND tenant_id = ?", sc.Slug, tenantID).First(&existing).Error; err == nil {
		existing.CategoryID = categoryID
		existing.Name = sc.Name
		existing.Description = sc.Description
		existing.Thumbnail = sc.Thumbnail
		_ = utils.DB.Save(&existing).Error
		return &existing
	}

	if err := utils.DB.Create(&sc).Error; err != nil {
		log.Fatalf("failed to create sub-category: %v", err)
	}
	return &sc
}

func ensureInstructor(tenantID uint) *models.Instructor {
	pw, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	role := "Instructor"
	designation := "Senior Instructor"
	image := "https://placehold.co/256x256/png?text=Instructor"
	ins := models.Instructor{
		UserID:      cuid.New(),
		FirstName:   "Demo",
		LastName:    ptr("Instructor"),
		Phone:       ptr("+8801000000000"),
		Email:       "instructor@local.dev",
		Password:    string(pw),
		Status:      true,
		TenantID:    tenantID,
		Role:        &role,
		Designation: &designation,
		Image:       &image,
	}

	var existing models.Instructor
	if err := utils.DB.Where("email = ? AND tenant_id = ?", ins.Email, tenantID).First(&existing).Error; err == nil {
		existing.FirstName = ins.FirstName
		existing.LastName = ins.LastName
		existing.Phone = ins.Phone
		existing.Status = true
		existing.Role = ins.Role
		existing.Designation = ins.Designation
		existing.Image = ins.Image
		_ = utils.DB.Save(&existing).Error
		return &existing
	}

	if err := utils.DB.Create(&ins).Error; err != nil {
		log.Fatalf("failed to create instructor: %v", err)
	}
	return &ins
}

func ensureCourse(tenantID, authorID uint) *models.CourseDetails {
	featured := "https://placehold.co/1200x630/png?text=Demo+Course"
	c := models.CourseDetails{
		Title:         "Demo Course: Local Development",
		Slug:          "demo-course-local-development",
		Summary:       "A demo course seeded for local testing. Safe to delete anytime.",
		Description:   ptr("This is demo content created by the demo seeder."),
		Visibility:    models.Public,
		PricingModel:  models.CoursePricingModelFree,
		FeaturedImage: &featured,
		AuthorID:      authorID,
		TenantID:      tenantID,
		Position:      1,
	}

	var existing models.CourseDetails
	if err := utils.DB.Where("slug = ? AND tenant_id = ?", c.Slug, tenantID).First(&existing).Error; err == nil {
		existing.Title = c.Title
		existing.Summary = c.Summary
		existing.Description = c.Description
		existing.Visibility = c.Visibility
		existing.PricingModel = c.PricingModel
		existing.FeaturedImage = c.FeaturedImage
		existing.AuthorID = authorID
		existing.Position = 1
		_ = utils.DB.Save(&existing).Error
		return &existing
	}

	if err := utils.DB.Create(&c).Error; err != nil {
		log.Fatalf("failed to create course: %v", err)
	}
	return &c
}

func ensureCourseGeneralSettings(tenantID, courseID, categoryID uint, subCategoryID *uint) {
	lang := "english"
	duration := "2h 30m"
	level := models.Beginner
	max := int32(0)
	settings := models.CourseGeneralSettings{
		CourseID:        courseID,
		DifficultyLevel: &level,
		MaximumStudent:  &max,
		Language:        &lang,
		CategoryID:      categoryID,
		SubCategoryID:   subCategoryID,
		Duration:        &duration,
	}

	var existing models.CourseGeneralSettings
	if err := utils.DB.Where("course_id = ?", courseID).First(&existing).Error; err == nil {
		existing.DifficultyLevel = settings.DifficultyLevel
		existing.Language = settings.Language
		existing.CategoryID = settings.CategoryID
		existing.SubCategoryID = settings.SubCategoryID
		existing.Duration = settings.Duration
		_ = utils.DB.Save(&existing).Error
		return
	}
	if err := utils.DB.Create(&settings).Error; err != nil {
		log.Fatalf("failed to create course general settings: %v", err)
	}
}

func ensureCourseInstructor(courseID, instructorID uint) {
	var existing models.CourseInstructor
	if err := utils.DB.Where("course_id = ? AND instructor_id = ?", courseID, instructorID).First(&existing).Error; err == nil {
		return
	}
	row := models.CourseInstructor{CourseID: courseID, InstructorID: instructorID}
	if err := utils.DB.Create(&row).Error; err != nil {
		log.Fatalf("failed to create course instructor mapping: %v", err)
	}
}

func ensureStudent(tenantID uint) *models.Student {
	pw, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	st := models.Student{
		UserID:    cuid.New(),
		FirstName: "Demo",
		LastName:  ptr("Student"),
		Phone:     ptr("+8801999999999"),
		Email:     "student@local.dev",
		Password:  string(pw),
		Status:    true,
		TenantID:  tenantID,
	}

	var existing models.Student
	if err := utils.DB.Where("email = ? AND tenant_id = ?", st.Email, tenantID).First(&existing).Error; err == nil {
		existing.FirstName = st.FirstName
		existing.LastName = st.LastName
		existing.Phone = st.Phone
		existing.Status = true
		_ = utils.DB.Save(&existing).Error
		return &existing
	}
	if err := utils.DB.Create(&st).Error; err != nil {
		log.Fatalf("failed to create student: %v", err)
	}
	return &st
}

func ensureEnrollment(tenantID, studentID, courseID uint) {
	var existing models.Enrollment
	if err := utils.DB.Where("student_id = ? AND course_id = ? AND tenant_id = ?", studentID, courseID, tenantID).First(&existing).Error; err == nil {
		return
	}
	e := models.Enrollment{StudentID: studentID, CourseID: courseID, TenantID: tenantID}
	if err := utils.DB.Create(&e).Error; err != nil {
		log.Fatalf("failed to create enrollment: %v", err)
	}
}

func ensureOrder(tenantID, studentID, courseID uint) {
	invoiceID := int64(1000001)
	var existing models.Order
	if err := utils.DB.Where("invoice_id = ? AND tenant_id = ?", invoiceID, tenantID).First(&existing).Error; err == nil {
		return
	}
	total := float64(0)
	o := models.Order{
		StudentID:     studentID,
		CourseID:      courseID,
		DiscountType:  "none",
		Discount:      0,
		Total:         total,
		PaymentStatus: models.OrderPaymentStatusUnpaid,
		InvoiceID:     invoiceID,
		PaymentType:   "manual",
		TenantID:      tenantID,
	}
	if err := utils.DB.Create(&o).Error; err != nil {
		log.Fatalf("failed to create order: %v", err)
	}
}

func ensurePaymentMethod(tenantID uint) {
	pm := models.PaymentMethod{
		Title:       "Manual Payment (Demo)",
		Image:       ptr("https://placehold.co/600x300/png?text=Payment"),
		Instruction: "Use this demo method for local testing. No real payment is processed.",
		Status:      true,
		TenantID:    tenantID,
	}
	var existing models.PaymentMethod
	if err := utils.DB.Where("title = ? AND tenant_id = ?", pm.Title, tenantID).First(&existing).Error; err == nil {
		existing.Instruction = pm.Instruction
		existing.Status = true
		_ = utils.DB.Save(&existing).Error
		return
	}
	if err := utils.DB.Create(&pm).Error; err != nil {
		log.Fatalf("failed to create payment method: %v", err)
	}
}

func ensureBanner(tenantID uint) {
	title := "Demo Banner"
	url := "https://example.com"
	b := models.Banner{
		Title:    &title,
		Url:      &url,
		Image:    "https://placehold.co/1600x500/png?text=Demo+Banner",
		TenantID: tenantID,
	}
	var existing models.Banner
	if err := utils.DB.Where("image = ? AND tenant_id = ?", b.Image, tenantID).First(&existing).Error; err == nil {
		return
	}
	if err := utils.DB.Create(&b).Error; err != nil {
		log.Fatalf("failed to create banner: %v", err)
	}
}

func ensureGeneralSettings(tenantID uint) {
	var existing models.GeneralSettings
	if err := utils.DB.Where("tenant_id = ?", tenantID).First(&existing).Error; err == nil {
		// keep org name stable
		if existing.OrgName == "" {
			existing.OrgName = "Lurnic (Local Demo)"
			_ = utils.DB.Save(&existing).Error
		}
		return
	}

	gs := models.GeneralSettings{
		OrgName:       "Lurnic (Local Demo)",
		StudentPrefix: "S-",
		TeacherPrefix: "T-",
		TenantID:      tenantID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := utils.DB.Create(&gs).Error; err != nil {
		log.Fatalf("failed to create general settings: %v", err)
	}
}

func ptr[T any](v T) *T { return &v }

