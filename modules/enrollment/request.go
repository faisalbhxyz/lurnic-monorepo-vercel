package enrollment

type Enrollment struct {
	CourseID  uint `json:"course_id" binding:"required"`
	StudentID uint `json:"student_id" binding:"required"`
}
