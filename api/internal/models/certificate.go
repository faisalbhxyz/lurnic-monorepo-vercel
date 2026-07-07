package models

import "time"

type CourseCertificateSettings struct {
	ID                  uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID            uint      `gorm:"column:course_id;uniqueIndex" json:"course_id"`
	IsEnabled           bool      `gorm:"column:is_enabled;default:false" json:"is_enabled"`
	CompletionPercent   uint8     `gorm:"column:completion_percent;default:100" json:"completion_percent"`
	CountLessons        bool      `gorm:"column:count_lessons;default:true" json:"count_lessons"`
	CountQuizzes        bool      `gorm:"column:count_quizzes;default:true" json:"count_quizzes"`
	CountAssignments    bool      `gorm:"column:count_assignments;default:true" json:"count_assignments"`
	TemplatePath        string    `gorm:"column:template_path;default:'/templates/minar-academy'" json:"template_path"`
	Title               *string   `json:"title"`
	SubtitleOne         *string   `gorm:"column:subtitle_one" json:"subtitle_one"`
	SubtitleTwo         *string   `gorm:"column:subtitle_two" json:"subtitle_two"`
	BrandLogo           *string   `gorm:"column:brand_logo" json:"brand_logo"`
	WatermarkImage      *string   `gorm:"column:watermark_image" json:"watermark_image"`
	WatermarkOpacity    uint8     `gorm:"column:watermark_opacity;default:30" json:"watermark_opacity"`
	OrganizationName    *string   `gorm:"column:organization_name" json:"organization_name"`
	SignerName          *string   `gorm:"column:signer_name" json:"signer_name"`
	SignerRole          *string   `gorm:"column:signer_role" json:"signer_role"`
	SignerOrg           *string   `gorm:"column:signer_org" json:"signer_org"`
	DualSignersEnabled  bool      `gorm:"column:dual_signers_enabled;default:false" json:"dual_signers_enabled"`
	Signer2Name         *string   `gorm:"column:signer2_name" json:"signer2_name"`
	Signer2Role         *string   `gorm:"column:signer2_role" json:"signer2_role"`
	Signer2Org          *string   `gorm:"column:signer2_org" json:"signer2_org"`
	OwnerSignature      *string   `gorm:"column:owner_signature" json:"owner_signature"`
	InstructorSignature *string   `gorm:"column:instructor_signature" json:"instructor_signature"`
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type StudentCertificate struct {
	ID                  uint               `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID            uint               `gorm:"column:tenant_id" json:"-"`
	StudentID           uint               `gorm:"column:student_id" json:"student_id"`
	CourseID            uint               `gorm:"column:course_id" json:"course_id"`
	CertificateNumber   string             `gorm:"column:certificate_number;uniqueIndex" json:"certificate_number"`
	StudentName         string             `gorm:"column:student_name" json:"student_name"`
	CourseTitle         string             `gorm:"column:course_title" json:"course_title"`
	ProgressPercent     float32            `gorm:"column:progress_percent" json:"progress_percent"`
	TemplatePath        string             `gorm:"column:template_path" json:"template_path"`
	Title               *string            `json:"title"`
	SubtitleOne         *string            `gorm:"column:subtitle_one" json:"subtitle_one"`
	SubtitleTwo         *string            `gorm:"column:subtitle_two" json:"subtitle_two"`
	BrandLogo           *string            `gorm:"column:brand_logo" json:"brand_logo"`
	WatermarkImage      *string            `gorm:"column:watermark_image" json:"watermark_image"`
	WatermarkOpacity    uint8              `gorm:"column:watermark_opacity;default:30" json:"watermark_opacity"`
	OrganizationName    *string            `gorm:"column:organization_name" json:"organization_name"`
	SignerName          *string            `gorm:"column:signer_name" json:"signer_name"`
	SignerRole          *string            `gorm:"column:signer_role" json:"signer_role"`
	SignerOrg           *string            `gorm:"column:signer_org" json:"signer_org"`
	DualSignersEnabled  bool               `gorm:"column:dual_signers_enabled;default:false" json:"dual_signers_enabled"`
	Signer2Name         *string            `gorm:"column:signer2_name" json:"signer2_name"`
	Signer2Role         *string            `gorm:"column:signer2_role" json:"signer2_role"`
	Signer2Org          *string            `gorm:"column:signer2_org" json:"signer2_org"`
	PricingModel        CoursePricingModel `gorm:"column:pricing_model;default:'free'" json:"pricing_model"`
	OwnerSignature      *string            `gorm:"column:owner_signature" json:"owner_signature"`
	InstructorSignature *string            `gorm:"column:instructor_signature" json:"instructor_signature"`
	IssuedAt            time.Time          `gorm:"column:issued_at" json:"issued_at"`
	CreatedAt           time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
}
