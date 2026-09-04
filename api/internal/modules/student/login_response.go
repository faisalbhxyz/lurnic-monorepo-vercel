package student

import "dashlearn/internal/models"

func studentLoginUserResponse(user models.Student, classProfile *models.StudentClassProfile) map[string]interface{} {
	status := "inactive"
	if user.Status {
		status = "active"
	}

	return map[string]interface{}{
		"id":            user.ID,
		"user_id":       user.UserID,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"email":         user.Email,
		"phone":         user.Phone,
		"profile_image": user.ProfileImage,
		"status":        status,
		"class_profile": toClassProfileResponse(classProfile),
	}
}
