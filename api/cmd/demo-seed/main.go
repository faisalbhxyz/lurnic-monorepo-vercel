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
	count := envUint("DEMO_COUNT", 5)

	author, err := findAuthor(tenantID)
	if err != nil {
		log.Fatalf("failed to find author user for demo seed: %v", err)
	}

	categories := ensureCategories(tenantID, count)
	subCategories := ensureSubCategories(tenantID, categories, count)
	instructors := ensureInstructors(tenantID, count)
	courses := ensureCourses(tenantID, author.ID, count)

	for i := uint(1); i <= count; i++ {
		category := categories[i-1]
		subCategory := subCategories[i-1]
		instructor := instructors[i-1]
		course := courses[i-1]

		ensureCourseGeneralSettings(tenantID, course.ID, category.ID, &subCategory.ID, i)
		ensureCourseInstructor(course.ID, instructor.ID)
	}

	students := ensureStudents(tenantID, count)
	for i := uint(1); i <= count; i++ {
		student := students[i-1]
		course := courses[i-1]
		ensureEnrollment(tenantID, student.ID, course.ID)
		ensureOrder(tenantID, student.ID, course.ID, i)
	}

	ensurePaymentMethods(tenantID, count)
	ensureBanners(tenantID, count)
	ensureGeneralSettings(tenantID)

	fmt.Println("Demo seed complete.")
	fmt.Println("Created/ensured demo records for:")
	fmt.Printf("- Categories: %d, Sub-categories: %d\n", count, count)
	fmt.Printf("- Instructors: %d\n", count)
	fmt.Printf("- Courses: %d (+ general settings + instructor mapping)\n", count)
	fmt.Printf("- Students: %d (+ enrollment + order per course)\n", count)
	fmt.Printf("- Payment methods: %d\n", count)
	fmt.Printf("- Banners: %d\n", count)
	fmt.Println("- General settings: 1")
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

func ensureCategories(tenantID, count uint) []*models.Category {
	out := make([]*models.Category, 0, count)
	for i := uint(1); i <= count; i++ {
		out = append(out, ensureCategory(tenantID, i))
	}
	return out
}

func ensureCategory(tenantID, n uint) *models.Category {
	desc := "Demo category for local development."
	thumb := fmt.Sprintf("https://placehold.co/600x400/png?text=Category+%d", n)
	c := models.Category{
		Name:        fmt.Sprintf("Demo Category %d", n),
		Slug:        fmt.Sprintf("demo-category-%d", n),
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

func ensureSubCategories(tenantID uint, categories []*models.Category, count uint) []*models.SubCategory {
	out := make([]*models.SubCategory, 0, count)
	for i := uint(1); i <= count; i++ {
		out = append(out, ensureSubCategory(tenantID, categories[i-1].ID, i))
	}
	return out
}

func ensureSubCategory(tenantID, categoryID, n uint) *models.SubCategory {
	desc := "Demo sub-category for local development."
	thumb := fmt.Sprintf("https://placehold.co/600x400/png?text=Sub+Category+%d", n)
	sc := models.SubCategory{
		CategoryID:  categoryID,
		Name:        fmt.Sprintf("Demo Sub Category %d", n),
		Slug:        fmt.Sprintf("demo-sub-category-%d", n),
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

func ensureInstructors(tenantID, count uint) []*models.Instructor {
	out := make([]*models.Instructor, 0, count)
	for i := uint(1); i <= count; i++ {
		out = append(out, ensureInstructor(tenantID, i))
	}
	return out
}

func ensureInstructor(tenantID, n uint) *models.Instructor {
	pw, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	role := "Instructor"
	designation := "Senior Instructor"
	image := fmt.Sprintf("https://placehold.co/256x256/png?text=Instructor+%d", n)
	ins := models.Instructor{
		UserID:      cuid.New(),
		FirstName:   "Demo",
		LastName:    ptr(fmt.Sprintf("Instructor %d", n)),
		Phone:       ptr(fmt.Sprintf("+88010000000%02d", n)),
		Email:       fmt.Sprintf("instructor+%d@local.dev", n),
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

func ensureCourses(tenantID, authorID, count uint) []*models.CourseDetails {
	out := make([]*models.CourseDetails, 0, count)
	for i := uint(1); i <= count; i++ {
		out = append(out, ensureCourse(tenantID, authorID, i))
	}
	return out
}

func ensureCourse(tenantID, authorID, n uint) *models.CourseDetails {
	featured := fmt.Sprintf("https://placehold.co/1200x630/png?text=Demo+Course+%d", n)
	c := models.CourseDetails{
		Title:         fmt.Sprintf("Demo Course %d: Local Development", n),
		Slug:          fmt.Sprintf("demo-course-local-development-%d", n),
		Summary:       fmt.Sprintf("Demo course %d seeded for local testing. Safe to delete anytime.", n),
		Description:   ptr(fmt.Sprintf("This is demo content created by the demo seeder (course %d).", n)),
		Visibility:    models.Public,
		PricingModel:  models.CoursePricingModelFree,
		FeaturedImage: &featured,
		AuthorID:      authorID,
		TenantID:      tenantID,
		Position:      int64(n),
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
		existing.Position = int64(n)
		_ = utils.DB.Save(&existing).Error
		return &existing
	}

	if err := utils.DB.Create(&c).Error; err != nil {
		log.Fatalf("failed to create course: %v", err)
	}
	return &c
}

func ensureCourseGeneralSettings(tenantID, courseID, categoryID uint, subCategoryID *uint, n uint) {
	lang := "english"
	duration := fmt.Sprintf("%dh %dm", 1+n, 15*n)
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

func ensureStudents(tenantID, count uint) []*models.Student {
	out := make([]*models.Student, 0, count)
	for i := uint(1); i <= count; i++ {
		out = append(out, ensureStudent(tenantID, i))
	}
	return out
}

func ensureStudent(tenantID, n uint) *models.Student {
	pw, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	st := models.Student{
		UserID:    cuid.New(),
		FirstName: "Demo",
		LastName:  ptr(fmt.Sprintf("Student %d", n)),
		Phone:     ptr(fmt.Sprintf("+88019999999%02d", n)),
		Email:     fmt.Sprintf("student+%d@local.dev", n),
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

func ensureOrder(tenantID, studentID, courseID uint, n uint) {
	invoiceID := int64(1000000 + n)
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

func ensurePaymentMethods(tenantID, count uint) {
	for i := uint(1); i <= count; i++ {
		ensurePaymentMethod(tenantID, i)
	}
}

func ensurePaymentMethod(tenantID, n uint) {
	pm := models.PaymentMethod{
		Title:       fmt.Sprintf("Manual Payment (Demo %d)", n),
		Image:       ptr(fmt.Sprintf("https://placehold.co/600x300/png?text=Payment+%d", n)),
		Instruction: fmt.Sprintf("Demo payment method %d for local testing. No real payment is processed.", n),
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

func ensureBanners(tenantID, count uint) {
	for i := uint(1); i <= count; i++ {
		ensureBanner(tenantID, i)
	}
}

func ensureBanner(tenantID, n uint) {
	title := fmt.Sprintf("Demo Banner %d", n)
	url := "https://example.com"
	image := fmt.Sprintf("https://placehold.co/1600x500/png?text=Demo+Banner+%d", n)
	b := models.Banner{
		Title:    &title,
		Url:      &url,
		Image:    image,
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

