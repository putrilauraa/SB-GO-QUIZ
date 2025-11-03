// handlers/book_handler.go

package handlers

import (
	"net/http"
	"log"
	"time"
	"fmt"
	"github.com/gin-gonic/gin"
	"database/sql"
	"SB-GO-QUIZ/models" 
)

const (
    MIN_RELEASE_YEAR = 1980
    MAX_RELEASE_YEAR = 2024
    THICKNESS_THICK  = "tebal"
    THICKNESS_THIN   = "tipis"
)

func GetAllBooks(c *gin.Context) {
	var books []models.Book

	query := `
		SELECT 
			id, title, description, image_url, release_year, price, total_page, thickness, category_id, created_at, created_by, modified_at, modified_by 
		FROM books
	`
	rows, err := DB.Query(query)
	if err != nil {
		log.Println("Database error in GetAllBooks:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve books"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var book models.Book
		
		err := rows.Scan(
            &book.ID, &book.Title, &book.Description, &book.ImageURL, &book.ReleaseYear, 
            &book.Price, &book.TotalPage, &book.Thickness, &book.CategoryID, &book.CreatedAt,
            &book.CreatedBy, &book.ModifiedAt, &book.ModifiedBy,
        )
		if err != nil {
			log.Println("Error scanning book data:", err)
			continue
		}
		books = append(books, book)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully retrieved all books",
		"data": books,
	})
}

func CreateBook(c *gin.Context) {
	var input models.BookInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid input: %s", err.Error())})
		return
	}

	if input.ReleaseYear < MIN_RELEASE_YEAR || input.ReleaseYear > MAX_RELEASE_YEAR {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Release year must be between %d and %d.", MIN_RELEASE_YEAR, MAX_RELEASE_YEAR),
		})
		return
	}
    
    var thickness string
    if input.TotalPage > 100 {
        thickness = THICKNESS_THICK
    } else {
        thickness = THICKNESS_THIN 
    }

	currentTime := time.Now()
	creator := "admin" 

	query := `
		INSERT INTO books (
			title, description, image_url, release_year, price, total_page, thickness, category_id, 
			created_at, created_by, modified_at, modified_by
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`
	var newID int
	
	err := DB.QueryRow(query, 
		input.Title, input.Description, input.ImageURL, input.ReleaseYear, 
		input.Price, input.TotalPage, thickness, input.CategoryID, 
		currentTime, creator, currentTime, creator,
	).Scan(&newID) 
	
	if err != nil {
		log.Println("Database error when creating book:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book in database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Book created successfully",
		"data": gin.H{
			"id": newID,
			"title": input.Title,
			"thickness": thickness,
		},
	})
}

func GetBookByID(c *gin.Context) {
	idStr := c.Param("id") 
	
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID format"})
		return
	}

	var book models.Book
	
	query := `
		SELECT 
			id, title, description, image_url, release_year, price, total_page, thickness, category_id, 
			created_at, created_by, modified_at, modified_by 
		FROM books 
		WHERE id = $1
	`
	
	err := DB.QueryRow(query, id).Scan(
		&book.ID, &book.Title, &book.Description, &book.ImageURL, &book.ReleaseYear, 
		&book.Price, &book.TotalPage, &book.Thickness, &book.CategoryID, 
		&book.CreatedAt, &book.CreatedBy, &book.ModifiedAt, &book.ModifiedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		log.Println("Database error in GetBookByID:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book detail"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully retrieved book detail",
		"data": book,
	})
}

func DeleteBook(c *gin.Context) {
	idStr := c.Param("id") 
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID format"})
		return
	}

	query := `DELETE FROM books WHERE id = $1`
	
	result, err := DB.Exec(query, id)
	if err != nil {
		log.Println("Database error during book deletion:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking deletion status"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found or already deleted."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Book with ID %d successfully deleted", id)})
}