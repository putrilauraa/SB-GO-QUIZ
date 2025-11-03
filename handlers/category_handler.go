package handlers

import (
	"fmt"
	"database/sql"
	"net/http"
	"time"
	"log"
	"github.com/gin-gonic/gin"
	"SB-GO-QUIZ/models" 
)

var DB *sql.DB 

func GetAllCategories(c *gin.Context) {
	var categories []models.Category

	rows, err := DB.Query("SELECT id, name, created_at, created_by, modified_at, modified_by FROM categories")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve categories"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var cat models.Category
		err := rows.Scan(&cat.ID, &cat.Name, &cat.CreatedAt, &cat.CreatedBy, &cat.ModifiedAt, &cat.ModifiedBy)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning category data"})
			return
		}
		categories = append(categories, cat)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully retrieved all categories",
		"data": categories,
	})
}

func CreateCategory(c *gin.Context) {
	var input models.CategoryInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentTime := time.Now()
	creator := "admin" 

	query := `
		INSERT INTO categories (name, created_at, created_by, modified_at, modified_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var newID int
	
	err := DB.QueryRow(query, 
		input.Name, 
		currentTime, 
		creator, 
		currentTime, 
		creator,
	).Scan(&newID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category in database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Category created successfully",
		"data": gin.H{
			"id": newID,
			"name": input.Name,
		},
	})
}

func GetCategoryByID(c *gin.Context) {
	idStr := c.Param("id") 
	
	var id int
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID format"})
		return
	}

	var cat models.Category
	
	query := `
		SELECT id, name, created_at, created_by, modified_at, modified_by 
		FROM categories 
		WHERE id = $1
	`
	err = DB.QueryRow(query, id).Scan(
		&cat.ID, &cat.Name, &cat.CreatedAt, &cat.CreatedBy, &cat.ModifiedAt, &cat.ModifiedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve category detail"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully retrieved category detail",
		"data": cat,
	})
}

func DeleteCategory(c *gin.Context) {
	idStr := c.Param("id") 
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID format"})
		return
	}

	query := `DELETE FROM categories WHERE id = $1`
	
	result, err := DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking deletion status"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found or already deleted."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Category with ID %d successfully deleted", id)})
}

func GetBooksByCategory(c *gin.Context) {
	categoryIDStr := c.Param("id")
	var categoryID int
	if _, err := fmt.Sscanf(categoryIDStr, "%d", &categoryID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID format"})
		return
	}

	var books []models.Book
	
	query := `
		SELECT id, title, description, image_url, release_year, price, total_page, thickness, category_id, created_at 
		FROM books 
		WHERE category_id = $1
	`
	rows, err := DB.Query(query, categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve books"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var book models.Book
		err := rows.Scan(
            &book.ID, &book.Title, &book.Description, &book.ImageURL, &book.ReleaseYear, 
            &book.Price, &book.TotalPage, &book.Thickness, &book.CategoryID, &book.CreatedAt,
        )
		if err != nil {
			log.Println("Error scanning book data:", err) 
			continue 
		}
		books = append(books, book)
	}
    
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully retrieved books for category ID %d", categoryID),
		"data": books,
	})
}