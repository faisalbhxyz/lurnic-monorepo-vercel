package student

import (
	"dashlearn/internal/apiresponse"
	"dashlearn/internal/models"
	"dashlearn/internal/utils"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	classLevelFive      = "class_five"
	classLevelSix       = "class_six"
	classLevelSeven     = "class_seven"
	classLevelEight     = "class_eight"
	classLevel910SSC    = "class_9_10_ssc"
	classLevelHSCBeyond = "hsc_beyond"

	hscBatch2026  = "hsc_2026"
	hscBatch2027  = "hsc_2027"
	hscBatch2028  = "hsc_2028"
	hscBatchAfter = "after_hsc"

	deptScience         = "science"
	deptBusinessStudies = "business_studies"
	deptHumanities      = "humanities"
)

var validClassLevels = map[string]struct{}{
	classLevelFive:      {},
	classLevelSix:       {},
	classLevelSeven:     {},
	classLevelEight:     {},
	classLevel910SSC:    {},
	classLevelHSCBeyond: {},
}

var validHscBatches = map[string]struct{}{
	hscBatch2026:  {},
	hscBatch2027:  {},
	hscBatch2028:  {},
	hscBatchAfter: {},
}

var validDepartments = map[string]struct{}{
	deptScience:         {},
	deptBusinessStudies: {},
	deptHumanities:      {},
}

// Default preferred_class_slug when client omits it (applied only if CMS slug exists).
var defaultSlugByClassLevel = map[string]string{
	classLevelFive:      "class-5",
	classLevelSix:       "class-6",
	classLevelSeven:     "class-7",
	classLevelEight:     "class-8",
	classLevel910SSC:    "ssc",
	classLevelHSCBeyond: "hsc",
}

// Reverse map for slug-only create (home picker before onboarding).
var classLevelByDefaultSlug = map[string]string{
	"class-5": classLevelFive,
	"class-6": classLevelSix,
	"class-7": classLevelSeven,
	"class-8": classLevelEight,
	"ssc":     classLevel910SSC,
	"hsc":     classLevelHSCBeyond,
}

// UpdateClassProfileInput supports full upsert and partial merge (nil pointer = omitted).
type UpdateClassProfileInput struct {
	ClassLevel          *string `json:"class_level"`
	HscBatch            *string `json:"hsc_batch"`
	Department          *string `json:"department"`
	PreferredClassSlug  *string `json:"preferred_class_slug"`
	OnboardingCompleted *bool   `json:"onboarding_completed"`
}

// ClassProfileResponse is the storefront shape for GET/PUT and nested login/details.
type ClassProfileResponse struct {
	ClassLevel          string  `json:"class_level"`
	HscBatch            *string `json:"hsc_batch"`
	Department          *string `json:"department"`
	PreferredClassSlug  *string `json:"preferred_class_slug"`
	OnboardingCompleted bool    `json:"onboarding_completed"`
	UpdatedAt           string  `json:"updated_at"`
}

func toClassProfileResponse(p *models.StudentClassProfile) *ClassProfileResponse {
	if p == nil {
		return nil
	}
	return &ClassProfileResponse{
		ClassLevel:          p.ClassLevel,
		HscBatch:            p.HscBatch,
		Department:          p.Department,
		PreferredClassSlug:  p.PreferredClassSlug,
		OnboardingCompleted: p.OnboardingCompleted,
		UpdatedAt:           p.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func loadClassProfile(db *gorm.DB, studentID uint) (*models.StudentClassProfile, error) {
	var profile models.StudentClassProfile
	err := db.Where("student_id = ?", studentID).First(&profile).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func academicNoteSlugExists(db *gorm.DB, tenantID uint, slug string) (bool, error) {
	var count int64
	err := db.Model(&models.AcademicNoteClass{}).
		Where("tenant_id = ? AND slug = ?", tenantID, slug).
		Count(&count).Error
	return count > 0, err
}

func ptrString(s string) *string { return &s }

type classProfileResolved struct {
	ClassLevel          string
	HscBatch            *string
	Department          *string
	PreferredClassSlug  *string
	OnboardingCompleted bool
}

func resolveClassProfile(
	existing *models.StudentClassProfile,
	input UpdateClassProfileInput,
) (*classProfileResolved, error) {
	hasClassLevel := input.ClassLevel != nil && strings.TrimSpace(*input.ClassLevel) != ""
	hasSlug := input.PreferredClassSlug != nil && strings.TrimSpace(*input.PreferredClassSlug) != ""
	hasAny := hasClassLevel ||
		input.HscBatch != nil ||
		input.Department != nil ||
		input.PreferredClassSlug != nil ||
		input.OnboardingCompleted != nil

	if !hasAny {
		return nil, errors.New("at least one class profile field is required")
	}

	out := &classProfileResolved{}

	if existing != nil {
		out.ClassLevel = existing.ClassLevel
		out.HscBatch = existing.HscBatch
		out.Department = existing.Department
		out.PreferredClassSlug = existing.PreferredClassSlug
		out.OnboardingCompleted = existing.OnboardingCompleted
	}

	if hasClassLevel {
		out.ClassLevel = strings.TrimSpace(*input.ClassLevel)
	}

	if input.HscBatch != nil {
		v := strings.TrimSpace(*input.HscBatch)
		if v == "" {
			out.HscBatch = nil
		} else {
			out.HscBatch = &v
		}
	}

	if input.Department != nil {
		v := strings.TrimSpace(*input.Department)
		if v == "" {
			out.Department = nil
		} else {
			out.Department = &v
		}
	}

	if input.PreferredClassSlug != nil {
		v := strings.TrimSpace(*input.PreferredClassSlug)
		if v == "" {
			out.PreferredClassSlug = nil
		} else {
			out.PreferredClassSlug = &v
		}
	}

	if input.OnboardingCompleted != nil {
		out.OnboardingCompleted = *input.OnboardingCompleted
	}

	// Slug-only create: infer class_level when possible.
	if existing == nil && out.ClassLevel == "" {
		if !hasSlug {
			return nil, errors.New("class_level is required")
		}
		slug := strings.TrimSpace(*input.PreferredClassSlug)
		inferred, ok := classLevelByDefaultSlug[slug]
		if !ok {
			return nil, errors.New("class_level is required when preferred_class_slug cannot be mapped")
		}
		out.ClassLevel = inferred
		out.OnboardingCompleted = false
	}

	if out.ClassLevel == "" {
		return nil, errors.New("class_level is required")
	}

	if _, ok := validClassLevels[out.ClassLevel]; !ok {
		return nil, errors.New("invalid class_level")
	}

	// Non–hsc_beyond: strip batch/department.
	if out.ClassLevel != classLevelHSCBeyond {
		out.HscBatch = nil
		out.Department = nil
		return out, nil
	}

	// Incomplete onboarding may leave HSC details empty until full form is submitted.
	if !out.OnboardingCompleted {
		if out.HscBatch != nil && *out.HscBatch != "" {
			if _, ok := validHscBatches[*out.HscBatch]; !ok {
				return nil, errors.New("invalid hsc_batch")
			}
			if *out.HscBatch == hscBatchAfter {
				out.Department = nil
			} else if out.Department != nil && *out.Department != "" {
				if _, ok := validDepartments[*out.Department]; !ok {
					return nil, errors.New("invalid department")
				}
			}
		} else {
			out.HscBatch = nil
			out.Department = nil
		}
		return out, nil
	}

	if out.HscBatch == nil || *out.HscBatch == "" {
		return nil, errors.New("hsc_batch is required when class_level is hsc_beyond")
	}
	if _, ok := validHscBatches[*out.HscBatch]; !ok {
		return nil, errors.New("invalid hsc_batch")
	}

	if *out.HscBatch == hscBatchAfter {
		out.Department = nil
		return out, nil
	}

	if out.Department == nil || *out.Department == "" {
		return nil, errors.New("department is required when hsc_batch is not after_hsc")
	}
	if _, ok := validDepartments[*out.Department]; !ok {
		return nil, errors.New("invalid department")
	}

	return out, nil
}

func upsertClassProfile(db *gorm.DB, tenantID, studentID uint, input UpdateClassProfileInput) (*models.StudentClassProfile, error) {
	existing, err := loadClassProfile(db, studentID)
	if err != nil {
		return nil, err
	}

	resolved, err := resolveClassProfile(existing, input)
	if err != nil {
		return nil, err
	}

	slugWasProvided := input.PreferredClassSlug != nil
	if !slugWasProvided && resolved.PreferredClassSlug == nil {
		if def, ok := defaultSlugByClassLevel[resolved.ClassLevel]; ok {
			exists, err := academicNoteSlugExists(db, tenantID, def)
			if err != nil {
				return nil, err
			}
			if exists {
				resolved.PreferredClassSlug = ptrString(def)
			}
		}
	}

	if resolved.PreferredClassSlug != nil {
		exists, err := academicNoteSlugExists(db, tenantID, *resolved.PreferredClassSlug)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("preferred_class_slug is not a valid academic-notes class slug")
		}
	}

	now := time.Now().UTC()
	row := models.StudentClassProfile{
		StudentID:           studentID,
		ClassLevel:          resolved.ClassLevel,
		HscBatch:            resolved.HscBatch,
		Department:          resolved.Department,
		PreferredClassSlug:  resolved.PreferredClassSlug,
		OnboardingCompleted: resolved.OnboardingCompleted,
		UpdatedAt:           now,
	}

	if existing == nil {
		row.CreatedAt = now
		if err := db.Create(&row).Error; err != nil {
			return nil, err
		}
		return &row, nil
	}

	if err := db.Model(&models.StudentClassProfile{}).
		Where("student_id = ?", studentID).
		Updates(map[string]interface{}{
			"class_level":          row.ClassLevel,
			"hsc_batch":            row.HscBatch,
			"department":           row.Department,
			"preferred_class_slug": row.PreferredClassSlug,
			"onboarding_completed": row.OnboardingCompleted,
			"updated_at":           now,
		}).Error; err != nil {
		return nil, err
	}

	row.CreatedAt = existing.CreatedAt
	return &row, nil
}

func isClassProfileValidationError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "required") ||
		strings.Contains(msg, "invalid") ||
		strings.Contains(msg, "preferred_class_slug") ||
		strings.Contains(msg, "cannot be mapped") ||
		strings.Contains(msg, "at least one")
}

// GetClassProfile handles GET /student/class-profile.
func GetClassProfile(c *gin.Context) {
	studentID := c.GetUint("user_id")
	if studentID == 0 {
		apiresponse.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
		return
	}

	profile, err := loadClassProfile(utils.DB, studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load class profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toClassProfileResponse(profile)})
}

// PutClassProfile handles PUT /student/class-profile (idempotent upsert / merge).
func PutClassProfile(c *gin.Context) {
	studentID := c.GetUint("user_id")
	tenantID := c.GetUint("tenant_id")
	if studentID == 0 {
		apiresponse.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
		return
	}

	var input UpdateClassProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiresponse.Validation(c, "Invalid JSON body")
		return
	}

	profile, err := upsertClassProfile(utils.DB, tenantID, studentID, input)
	if err != nil {
		if isClassProfileValidationError(err) {
			apiresponse.Validation(c, err.Error())
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save class profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toClassProfileResponse(profile)})
}

// studentDetailsPayload nests class_profile on the storefront details object.
func studentDetailsPayload(user models.StudentDetailsRes, classProfile *models.StudentClassProfile) map[string]interface{} {
	return map[string]interface{}{
		"id":            user.ID,
		"user_id":       user.UserID,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"phone":         user.Phone,
		"email":         user.Email,
		"profile_image": user.ProfileImage,
		"status":        user.Status,
		"created_at":    user.CreatedAt,
		"updated_at":    user.UpdatedAt,
		"enrollments":   user.Enrollments,
		"class_profile": toClassProfileResponse(classProfile),
	}
}
