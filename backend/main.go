package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite" // Import the pure Go SQLite driver
)

type Note struct {
	gorm.Model
	Title   string `json:"title"`
	Content string `json:"content"`
}

func main() {
	// Connect to SQLite database using modernc.org/sqlite
	db, err := sql.Open("sqlite", "./notes.db")
	if err != nil {
		panic("Failed to connect to database")
	}

	// Wrap the *sql.DB with GORM
	gormDB, err := gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        "./notes.db",
		Conn:       db,
	}, &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Auto-migrate the schema
	gormDB.AutoMigrate(&Note{})

	// Initialize Gin router
	r := gin.Default()

	// Serve static files (frontend)
	r.Static("/static", "./frontend")

	// Route to serve the frontend HTML file
	r.GET("/", func(c *gin.Context) {
		c.File("./frontend/index.html")
	})

	// CORS middleware (to allow frontend to communicate with backend)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Error handling middleware
	r.Use(errorHandler())

	// Create a note
	r.POST("/notes", func(c *gin.Context) {
		var note Note
		if err := c.ShouldBindJSON(&note); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate title and content
		if note.Title == "" || note.Content == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Title and content are required"})
			return
		}

		gormDB.Create(&note)
		c.JSON(http.StatusOK, note)
	})

	// Get all notes with search functionality
	r.GET("/notes", func(c *gin.Context) {
		query := c.Query("q")
		var notes []Note

		if query != "" {
			gormDB.Where("title LIKE ? OR content LIKE ?", "%"+query+"%", "%"+query+"%").Find(&notes)
		} else {
			gormDB.Find(&notes)
		}

		c.JSON(http.StatusOK, notes)
	})

	// Update a note
	r.PUT("/notes/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" || id == "undefined" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
			return
		}

		var note Note
		if err := gormDB.First(&note, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
			return
		}

		var updatedNote Note
		if err := c.ShouldBindJSON(&updatedNote); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate title and content
		if updatedNote.Title == "" || updatedNote.Content == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Title and content are required"})
			return
		}

		note.Title = updatedNote.Title
		note.Content = updatedNote.Content
		gormDB.Save(&note)
		c.JSON(http.StatusOK, note)
	})

	// Delete a note
	r.DELETE("/notes/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" || id == "undefined" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
			return
		}

		var note Note
		if err := gormDB.First(&note, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
			return
		}

		gormDB.Delete(&note)
		c.JSON(http.StatusOK, gin.H{"message": "Note deleted"})
	})

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if PORT is not set
	}
	fmt.Printf("Server is running on http://localhost:%s\n", port)
	r.Run(":" + port)
}

// Error handling middleware
func errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": c.Errors.Last().Error()})
		}
	}
}
