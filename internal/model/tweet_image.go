package model

type TweetImage struct {
	TweetID      string `gorm:"type:uuid;primaryKey;not null" json:"tweet_id"`
	ImageID      string `gorm:"type:uuid;primaryKey;not null" json:"image_id"`
	DisplayOrder int    `gorm:"type:integer;not null;default:0" json:"display_order"`
	Tweet        Tweet  `gorm:"foreignKey:TweetID;constraint:OnDelete:CASCADE" json:"-"`
	Image        Image  `gorm:"foreignKey:ImageID;constraint:OnDelete:CASCADE" json:"-"`
}

func (TweetImage) TableName() string {
	return "blog.tweet_images"
}
