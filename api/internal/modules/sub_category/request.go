package subcategory

type CreateSubCategoryInput struct {
	CategoryID  uint    `json:"category_id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Slug        string  `json:"slug" binding:"required"`
	Description *string `json:"description" binding:"omitempty"`
}
