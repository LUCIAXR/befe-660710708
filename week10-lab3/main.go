package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	_ "week10-lab3/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type Book struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	ISBN      string    `json:"isbn"`
	Year      int       `json:"year"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

var db *sql.DB

func initDB() {
	var err error

	host := getEnv("DB_HOST", "")
	name := getEnv("DB_NAME", "")
	user := getEnv("DB_USER", "")
	password := getEnv("DB_PASSWORD", "")
	port := getEnv("DB_PORT", "")

	conSt := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)
	//fmt.Println(conSt)
	db, err = sql.Open("postgres", conSt)
	if err != nil {
		log.Fatal("failed to open database")
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("failed to connect to database")
	}
	log.Println("successfully connect to database")
}

// @Summary Get all book
// @Description Get details of a books
// @Tags Books
// @Produce  json
// @Success 200  {array}  Book
// @Failure 500  {object}  ErrorResponse
// @Router  /books [get]
func getAllBooks(c *gin.Context) {
	yearInput := c.Query("year")
	var rows *sql.Rows
	var err error
	// ลูกค้าถาม "มีหนังสืออะไรบ้าง"

	if yearInput != "" {
		rows, err = db.Query("SELECT id, title, author, isbn, year, price, created_at, updated_at FROM books WHERE year = $1", yearInput)
	} else {
		rows, err = db.Query("SELECT id, title, author, isbn, year, price, created_at, updated_at FROM books")
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close() // ต้องปิด rows เสมอ เพื่อคืน Connection กลับ pool

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			// handle error
		}
		books = append(books, book)
	}
	if books == nil {
		books = []Book{}
	}

	c.JSON(http.StatusOK, books)
}

// @Summary Get book by ID
// @Description Get details of a book by ID
// @Tags Books
// @Produce  json
// @Param   id   path      int     true  "book ID"
// @Success 200  {object}  Book
// @Failure 500  {object}  ErrorResponse
// @Router  /books/{id} [get]
func getBook(c *gin.Context) {
	id := c.Param("id")
	var book Book

	// QueryRow ใช้เมื่อคาดว่าจะได้ผลลัพธ์ 0 หรือ 1 แถว
	err := db.QueryRow("SELECT id, title, author, isbn, year, price FROM books WHERE id = $1", id).
		Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

// @Summary create book
// @Description create details of a book
// @Tags Books
// @Produce  json
// @Param   book  body  Book  true  "book data"
// @Success 200  {object}  Book
// @Failure 500  {object}  ErrorResponse
// @Router  /books [post]
func createBook(c *gin.Context) {
	var newBook Book

	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ใช้ RETURNING เพื่อดึงค่าที่ database generate (id, timestamps)
	var id int
	var createdAt, updatedAt time.Time

	err := db.QueryRow(
		`INSERT INTO books (title, author, isbn, year, price)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING id, created_at, updated_at`,
		newBook.Title, newBook.Author, newBook.ISBN, newBook.Year, newBook.Price,
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newBook.ID = id
	newBook.CreatedAt = createdAt
	newBook.UpdatedAt = updatedAt

	c.JSON(http.StatusCreated, newBook) // ใช้ 201 Created
}

// @Summary update book by ID
// @Description update details of a book by ID
// @Tags Books
// @Produce  json
// @Param   id    path   int   true  "book ID"
// @Param   book  body  Book  true  "update book data"
// @Success 200  {object}  Book
// @Failure 500  {object}  ErrorResponse
// @Router  /books/{id} [put]
func updateBook(c *gin.Context) {
	var ID int
	id := c.Param("id")
	var updateBook Book

	if err := c.ShouldBindJSON(&updateBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedAt time.Time
	err := db.QueryRow(
		`UPDATE books
         SET title = $1, author = $2, isbn = $3, year = $4, price = $5
         WHERE id = $6
         RETURNING id, updated_at`,
		updateBook.Title, updateBook.Author, updateBook.ISBN,
		updateBook.Year, updateBook.Price, id,
	).Scan(&ID, &updatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	updateBook.ID = ID
	updateBook.UpdatedAt = updatedAt
	c.JSON(http.StatusOK, updateBook)
}

// @Summary delete book by ID
// @Description delete details of a book by ID
// @Tags Books
// @Produce  json
// @Param   id   path      int     true  "book ID"
// @Success 200  {object}  Book
// @Failure 500  {object}  ErrorResponse
// @Router  /books/{id} [delete]
func deleteBook(c *gin.Context) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM books WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book deleted successfully"})
}

// @title           Simple API Example
// @version         1.0
// @description     This is a simple example of using Gin with Swagger.
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	initDB()
	defer db.Close()
	r := gin.Default()
	r.Use(cors.Default())

	// Swagger endpoint
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", func(c *gin.Context) {
		err := db.Ping()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "unhealthy", "error": err})
			return
		}
		c.JSON(200, gin.H{"message": "healthy"})
	})

	api := r.Group("/api/v1")
	{
		api.GET("/books", getAllBooks)
		api.GET("/books/:id", getBook)
		api.POST("/books", createBook)
		api.PUT("/books/:id", updateBook)
		api.DELETE("/books/:id", deleteBook)
	}

	r.Run(":8080")
}
