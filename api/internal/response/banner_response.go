package response

import "time"

type BannerResponse struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     *string   `json:"title" form:"title" gorm:"type:varchar(100);null"`
	Url       *string   `json:"url" form:"url" gorm:"type:varchar(255);null"`
	Image     string    `json:"image" form:"image"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (BannerResponse) TableName() string {
	return "banners"
}
