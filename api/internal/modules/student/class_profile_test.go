package student

import (
	"dashlearn/internal/models"
	"testing"
)

func TestResolveClassProfile_ClassEight(t *testing.T) {
	level := classLevelEight
	done := true
	out, err := resolveClassProfile(nil, UpdateClassProfileInput{
		ClassLevel:          &level,
		OnboardingCompleted: &done,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.ClassLevel != classLevelEight {
		t.Fatalf("class_level = %s", out.ClassLevel)
	}
	if out.HscBatch != nil || out.Department != nil {
		t.Fatalf("expected nil batch/department")
	}
	if !out.OnboardingCompleted {
		t.Fatalf("expected onboarding_completed true")
	}
}

func TestResolveClassProfile_HSCRequiresBatch(t *testing.T) {
	level := classLevelHSCBeyond
	done := true
	_, err := resolveClassProfile(nil, UpdateClassProfileInput{
		ClassLevel:          &level,
		OnboardingCompleted: &done,
	})
	if err == nil {
		t.Fatal("expected error for missing hsc_batch")
	}
}

func TestResolveClassProfile_HSCRequiresDepartment(t *testing.T) {
	level := classLevelHSCBeyond
	batch := hscBatch2027
	done := true
	_, err := resolveClassProfile(nil, UpdateClassProfileInput{
		ClassLevel:          &level,
		HscBatch:            &batch,
		OnboardingCompleted: &done,
	})
	if err == nil {
		t.Fatal("expected error for missing department")
	}
}

func TestResolveClassProfile_AfterHSCStripsDepartment(t *testing.T) {
	level := classLevelHSCBeyond
	batch := hscBatchAfter
	dept := deptScience
	out, err := resolveClassProfile(nil, UpdateClassProfileInput{
		ClassLevel: &level,
		HscBatch:   &batch,
		Department: &dept,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Department != nil {
		t.Fatalf("expected department stripped for after_hsc")
	}
}

func TestResolveClassProfile_StripHSCFieldsForJunior(t *testing.T) {
	level := classLevelEight
	batch := hscBatch2026
	dept := deptScience
	out, err := resolveClassProfile(nil, UpdateClassProfileInput{
		ClassLevel: &level,
		HscBatch:   &batch,
		Department: &dept,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.HscBatch != nil || out.Department != nil {
		t.Fatalf("expected HSC fields stripped")
	}
}

func TestResolveClassProfile_InvalidEnums(t *testing.T) {
	bad := "class_99"
	if _, err := resolveClassProfile(nil, UpdateClassProfileInput{ClassLevel: &bad}); err == nil {
		t.Fatal("expected invalid class_level error")
	}

	level := classLevelHSCBeyond
	batch := "hsc_2099"
	done := true
	if _, err := resolveClassProfile(nil, UpdateClassProfileInput{
		ClassLevel:          &level,
		HscBatch:            &batch,
		OnboardingCompleted: &done,
	}); err == nil {
		t.Fatal("expected invalid hsc_batch error")
	}

	batch = hscBatch2027
	dept := "arts"
	if _, err := resolveClassProfile(nil, UpdateClassProfileInput{
		ClassLevel:          &level,
		HscBatch:            &batch,
		Department:          &dept,
		OnboardingCompleted: &done,
	}); err == nil {
		t.Fatal("expected invalid department error")
	}
}

func TestResolveClassProfile_SlugOnlyCreate(t *testing.T) {
	slug := "class-8"
	out, err := resolveClassProfile(nil, UpdateClassProfileInput{PreferredClassSlug: &slug})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.ClassLevel != classLevelEight {
		t.Fatalf("inferred class_level = %s", out.ClassLevel)
	}
	if out.OnboardingCompleted {
		t.Fatal("slug-only create should leave onboarding incomplete")
	}
}

func TestResolveClassProfile_SlugOnlyHSCIncomplete(t *testing.T) {
	slug := "hsc"
	out, err := resolveClassProfile(nil, UpdateClassProfileInput{PreferredClassSlug: &slug})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.ClassLevel != classLevelHSCBeyond {
		t.Fatalf("inferred class_level = %s", out.ClassLevel)
	}
	if out.HscBatch != nil {
		t.Fatal("incomplete HSC should allow nil batch")
	}
}

func TestResolveClassProfile_SlugOnlyUnmapped(t *testing.T) {
	slug := "class-9"
	_, err := resolveClassProfile(nil, UpdateClassProfileInput{PreferredClassSlug: &slug})
	if err == nil {
		t.Fatal("expected error when slug cannot be mapped")
	}
}

func TestResolveClassProfile_MergeKeepsPrevious(t *testing.T) {
	existing := &models.StudentClassProfile{
		ClassLevel:          classLevelEight,
		PreferredClassSlug:  ptrString("class-8"),
		OnboardingCompleted: true,
	}
	slug := "ssc"
	out, err := resolveClassProfile(existing, UpdateClassProfileInput{PreferredClassSlug: &slug})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.ClassLevel != classLevelEight {
		t.Fatalf("class_level should stay %s, got %s", classLevelEight, out.ClassLevel)
	}
	if out.PreferredClassSlug == nil || *out.PreferredClassSlug != "ssc" {
		t.Fatalf("slug not updated")
	}
	if !out.OnboardingCompleted {
		t.Fatal("onboarding_completed should remain true")
	}
}

func TestResolveClassProfile_HSCHappyPath(t *testing.T) {
	level := classLevelHSCBeyond
	batch := hscBatch2027
	dept := deptScience
	slug := "hsc"
	done := true
	out, err := resolveClassProfile(nil, UpdateClassProfileInput{
		ClassLevel:          &level,
		HscBatch:            &batch,
		Department:          &dept,
		PreferredClassSlug:  &slug,
		OnboardingCompleted: &done,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.HscBatch == nil || *out.HscBatch != hscBatch2027 {
		t.Fatalf("hsc_batch mismatch")
	}
	if out.Department == nil || *out.Department != deptScience {
		t.Fatalf("department mismatch")
	}
}
