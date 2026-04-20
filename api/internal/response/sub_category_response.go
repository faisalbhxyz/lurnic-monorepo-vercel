package response

import "time"

type SubCategoryResponse struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Slug        string            `json:"slug"`
	Description *string           `json:"description,omitempty"`
	Thumbnail   *string           `json:"thumbnail,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	CategoryID  uint              `json:"category_id,omitempty"`
	Category    *CategoryResponse `json:"category,omitempty"`
}
