// main.go

package main

import (
	"database/sql"
	"log"
	"os"
	_ "github.com/lib/pq"
	"github.com/gin-gonic/gin"
)

func initDB() *sql.DB {
    connStr := os.Getenv("DATABASE_URL")
    
    if connStr == "" {
        log.Println("Using local database credentials (DATABASE_URL not set).")
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
    log.Println("Database connection successful!")
    return db
}

func main() {    
    db := initDB()
    defer db.Close()
    
    router := gin.Default()
    router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "API is running locally"})
    })
    
    log.Fatal(router.Run(":8080")) 
}