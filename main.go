package main

import (
	"os"
	"log"
	"database/sql"
    _ "github.com/lib/pq"
	"github.com/gin-gonic/gin"
	"SB-GO-QUIZ/middlewares"
	"SB-GO-QUIZ/handlers"
)

func initDB() *sql.DB {
	connStr := os.Getenv("DATABASE_URL") 
    if connStr == "" {
        connStr = "host=localhost port=5432 user=postgres password=postgres12345 dbname=go_quiz_db sslmode=disable"
    } else {
        log.Println("Using production database credentials from DATABASE_URL.")
	}
	
	db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Error connecting to the database: ", err)
    }
    err = db.Ping()
    if err != nil {
        log.Fatal("Error pinging database: ", err)
    }
    return db
}

func main() {
	db := initDB()
    handlers.DB = db
    defer db.Close()

	router := gin.Default()
	apiGroup := router.Group("/api")
	
	apiGroup.Use(middlewares.BasicAuthMiddleware())
	{
		apiGroup.GET("/categories", handlers.GetAllCategories) 
		apiGroup.POST("/categories", handlers.CreateCategory)
		apiGroup.GET("/categories/:id", handlers.GetCategoryByID)
		apiGroup.DELETE("/categories/:id", handlers.DeleteCategory)
		apiGroup.GET("/categories/:id/books", handlers.GetBooksByCategory)

		apiGroup.GET("/books", handlers.GetAllBooks) 
		apiGroup.POST("/books", handlers.CreateBook)
		apiGroup.GET("/books/:id", handlers.GetBookByID)
		apiGroup.DELETE("/books/:id", handlers.DeleteBook)
	}

	router.Run(":8080")
}