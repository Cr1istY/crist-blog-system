package model

import (
	"time"

	"github.com/google/uuid"
)

const RootCategoryID = "00000000-0000-0000-0000-000000000000"

type Category struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"type:text;not null" json:"name"`
	Slug        string    `gorm:"type:text;not null" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	ParentID    uuid.UUID `gorm:"type:uuid;" json:"parent_id"`
	DeletedFlag bool      `gorm:"type:boolean;default:false" json:"deleted_flag"`
}

type CreatePostCategory struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	ParentID uuid.UUID `json:"parent_id"`
}

type CreateCategory struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ParentID    uuid.UUID `json:"parent_id"`
}

type UpdateCategory struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ParentID    uuid.UUID `json:"parent_id"`
}

type CategoryToFrontend struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
}

func (C *Category) TableName() string {
	return "blog.categories"
}
