package models

import "time"

type Book struct {
    ID           int       `json:"id"`
    Title        string    `json:"title"`
    Description  string    `json:"description"`
    ImageURL     string    `json:"image_url"`
    ReleaseYear  int       `json:"release_year"`
    Price        int       `json:"price"`
    TotalPage    int       `json:"total_page"`
    Thickness    string    `json:"thickness"`
    CategoryID   int       `json:"category_id"`
    CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`   
    ModifiedAt   time.Time `json:"modified_at"` 
    ModifiedBy   string    `json:"modified_by"`
}

type BookInput struct {
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description"`
	ImageURL     string `json:"image_url"`
	ReleaseYear  int    `json:"release_year" binding:"required"`
	Price        int    `json:"price" binding:"required"`
	TotalPage    int    `json:"total_page" binding:"required"`
	CategoryID   int    `json:"category_id" binding:"required"`
}