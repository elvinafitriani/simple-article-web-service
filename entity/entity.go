package entity

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Author string `json:"author" binding:"required"`
	Title  string `json:"title" binding:"required"`
	Body   string `json:"body" binding:"required"`
}
