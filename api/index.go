
package handler

import (
	"log"
	"net/http"
	"os"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/gin-gonic/gin"
	"SB-GO-QUIZ/middlewares" 
	"SB-GO-QUIZ/handlers"
)

var router *gin.Engine
var DB *sql.DB

func initDB() *sql.DB {
    connStr := os.Getenv("DATABASE_URL")
    if connStr == "" {
        connStr = "host=localhost port=5432 user=postgres password=postgres12345 dbname=go_quiz_db sslmode=disable"
    }

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Error connecting to the database: ", err) 
    }
    db.Ping()
    return db
}

func init() {
	DB = initDB()
	handlers.DB = DB

	router = gin.Default()
	
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
}

func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}