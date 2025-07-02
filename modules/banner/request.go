package banner

type CreateBannerInput struct {
	Title *string `json:"title" form:"title" binding:"omitempty"`
	Url   *string `json:"url" form:"url" binding:"omitempty"`
	Image string  `json:"image" form:"image" binding:"required"`
}

type UpdateBannerInput struct {
	Title *string `json:"title" form:"title" binding:"omitempty"`
	Url   *string `json:"url" form:"url" binding:"omitempty"`
	Image *string
}
