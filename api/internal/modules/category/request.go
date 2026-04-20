package category

type CreateCategoryInput struct {
	Name        string  `json:"name" binding:"required"`
	Slug        string  `json:"slug" binding:"required"`
	Description *string `json:"description" binding:"omitempty"`
}
