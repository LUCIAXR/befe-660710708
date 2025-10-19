package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

var db *sql.DB

type Book struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	ISBN          string    `json:"isbn"`
	Year          int       `json:"year"`
	Price         float64   `json:"price"`
	Category      string    `json:"category"`
	CoverImage    string    `json:"coverImage,omitempty"`
	Description   string    `json:"description,omitempty"`
	Rating        float64   `json:"rating,omitempty"`
	Reviews       int       `json:"reviews,omitempty"`
	IsNew         bool      `json:"isNew,omitempty"`
	Discount      int       `json:"discount,omitempty"`
	OriginalPrice float64   `json:"originalPrice,omitempty"`
	Created_At    time.Time `json:"created_at"`
	Updated_At    time.Time `json:"updated_at"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func initDB() {
	var err error
	host := getEnv("DB_HOST", "localhost")
	name := getEnv("DB_NAME", "bookstore")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	port := getEnv("DB_PORT", "5432")

	conSt := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)

	db, err = sql.Open("postgres", conSt)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// กำหนดจำนวน Connection สูงสุด
	db.SetMaxOpenConns(25)

	// กำหนดจำนวน Idle connection สูงสุด
	db.SetMaxIdleConns(20)

	// กำหนดอายุของ Connection
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Database connected successfully!")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// @Summary Health check
// @Description Check if the service is running
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func getHealth(c *gin.Context) {
	err := db.Ping()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "Unhealthy", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "healthy"})
}

// @Summary Get all categories
// @Description Get list of all book categories
// @Tags Categories
// @Produce json
// @Success 200 {array} string
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/categories [get]
func getCategories(c *gin.Context) {
	rows, err := db.Query(`
		SELECT DISTINCT category 
		FROM books 
		WHERE category IS NOT NULL AND category != ''
		ORDER BY category
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		categories = append(categories, category)
	}

	if categories == nil {
		categories = []string{}
	}

	c.JSON(http.StatusOK, categories)
}

// @Summary Search books
// @Description Search books by keyword in title or author
// @Tags Books
// @Produce json
// @Param q query string true "Search keyword"
// @Success 200 {array} Book
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/books/search [get]
func searchBooks(c *gin.Context) {
	keyword := c.Query("q")

	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search keyword is required"})
		return
	}

	searchPattern := "%" + keyword + "%"

	rows, err := db.Query(`
		SELECT id, title, author, isbn, year, price, category, 
		       COALESCE(cover_image, '') as cover_image,
		       COALESCE(description, '') as description,
		       COALESCE(rating, 0) as rating,
		       COALESCE(reviews, 0) as reviews,
		       COALESCE(is_new, false) as is_new,
		       COALESCE(discount, 0) as discount,
		       COALESCE(original_price, 0) as original_price,
		       created_at, updated_at
		FROM books 
		WHERE LOWER(title) LIKE LOWER($1) OR LOWER(author) LIKE LOWER($1)
		ORDER BY created_at DESC
	`, searchPattern)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(
			&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price,
			&book.Category, &book.CoverImage, &book.Description, &book.Rating,
			&book.Reviews, &book.IsNew, &book.Discount, &book.OriginalPrice,
			&book.Created_At, &book.Updated_At,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		books = append(books, book)
	}

	if books == nil {
		books = []Book{}
	}

	c.JSON(http.StatusOK, books)
}

// @Summary Get featured books
// @Description Get featured/recommended books
// @Tags Books
// @Produce json
// @Success 200 {array} Book
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/books/featured [get]
func getFeaturedBooks(c *gin.Context) {
	rows, err := db.Query(`
		SELECT id, title, author, isbn, year, price, category, 
		       COALESCE(cover_image, '') as cover_image,
		       COALESCE(description, '') as description,
		       COALESCE(rating, 0) as rating,
		       COALESCE(reviews, 0) as reviews,
		       COALESCE(is_new, false) as is_new,
		       COALESCE(discount, 0) as discount,
		       COALESCE(original_price, 0) as original_price,
		       created_at, updated_at
		FROM books 
		WHERE rating >= 4.0
		ORDER BY rating DESC, reviews DESC
		LIMIT 10
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(
			&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price,
			&book.Category, &book.CoverImage, &book.Description, &book.Rating,
			&book.Reviews, &book.IsNew, &book.Discount, &book.OriginalPrice,
			&book.Created_At, &book.Updated_At,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		books = append(books, book)
	}

	if books == nil {
		books = []Book{}
	}

	c.JSON(http.StatusOK, books)
}

// @Summary Get new books
// @Description Get latest books ordered by created date
// @Tags Books
// @Produce json
// @Param limit query int false "Number of books to return (default 5)"
// @Success 200 {array} Book
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/books/new [get]
func getNewBooks(c *gin.Context) {
	rows, err := db.Query(`
		SELECT id, title, author, isbn, year, price, category, 
		       COALESCE(cover_image, '') as cover_image,
		       COALESCE(description, '') as description,
		       COALESCE(rating, 0) as rating,
		       COALESCE(reviews, 0) as reviews,
		       COALESCE(is_new, false) as is_new,
		       COALESCE(discount, 0) as discount,
		       COALESCE(original_price, 0) as original_price,
		       created_at, updated_at
		FROM books 
		ORDER BY created_at DESC 
		LIMIT 5
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(
			&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price,
			&book.Category, &book.CoverImage, &book.Description, &book.Rating,
			&book.Reviews, &book.IsNew, &book.Discount, &book.OriginalPrice,
			&book.Created_At, &book.Updated_At,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		books = append(books, book)
	}

	if books == nil {
		books = []Book{}
	}

	c.JSON(http.StatusOK, books)
}

// @Summary Get discounted books
// @Description Get books with discount
// @Tags Books
// @Produce json
// @Success 200 {array} Book
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/books/discounted [get]
func getDiscountedBooks(c *gin.Context) {
	rows, err := db.Query(`
		SELECT id, title, author, isbn, year, price, category, 
		       COALESCE(cover_image, '') as cover_image,
		       COALESCE(description, '') as description,
		       COALESCE(rating, 0) as rating,
		       COALESCE(reviews, 0) as reviews,
		       COALESCE(is_new, false) as is_new,
		       COALESCE(discount, 0) as discount,
		       COALESCE(original_price, 0) as original_price,
		       created_at, updated_at
		FROM books 
		WHERE discount > 0
		ORDER BY discount DESC, created_at DESC
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(
			&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price,
			&book.Category, &book.CoverImage, &book.Description, &book.Rating,
			&book.Reviews, &book.IsNew, &book.Discount, &book.OriginalPrice,
			&book.Created_At, &book.Updated_At,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		books = append(books, book)
	}

	if books == nil {
		books = []Book{}
	}

	c.JSON(http.StatusOK, books)
}

// @Summary Get all books
// @Description Get all books or filter by category
// @Tags Books
// @Produce json
// @Param category query string false "Filter by category"
// @Success 200 {array} Book
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/books [get]
func getAllBooks(c *gin.Context) {
	var rows *sql.Rows
	var err error

	categoryQ := c.Query("category")

	if categoryQ == "" {
		rows, err = db.Query(`
			SELECT id, title, author, isbn, year, price, category, 
			       COALESCE(cover_image, '') as cover_image,
			       COALESCE(description, '') as description,
			       COALESCE(rating, 0) as rating,
			       COALESCE(reviews, 0) as reviews,
			       COALESCE(is_new, false) as is_new,
			       COALESCE(discount, 0) as discount,
			       COALESCE(original_price, 0) as original_price,
			       created_at, updated_at
			FROM books
			ORDER BY created_at DESC
		`)
	} else {
		rows, err = db.Query(`
			SELECT id, title, author, isbn, year, price, category, 
			       COALESCE(cover_image, '') as cover_image,
			       COALESCE(description, '') as description,
			       COALESCE(rating, 0) as rating,
			       COALESCE(reviews, 0) as reviews,
			       COALESCE(is_new, false) as is_new,
			       COALESCE(discount, 0) as discount,
			       COALESCE(original_price, 0) as original_price,
			       created_at, updated_at
			FROM books 
			WHERE LOWER(category) = LOWER($1)
			ORDER BY created_at DESC
		`, categoryQ)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(
			&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price,
			&book.Category, &book.CoverImage, &book.Description, &book.Rating,
			&book.Reviews, &book.IsNew, &book.Discount, &book.OriginalPrice,
			&book.Created_At, &book.Updated_At,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		books = append(books, book)
	}

	if books == nil {
		books = []Book{}
	}

	c.JSON(http.StatusOK, books)
}

// @Summary Get book by ID
// @Description Get details of specific book
// @Tags Books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} Book
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/books/{id} [get]
func getBook(c *gin.Context) {
	id := c.Param("id")
	var book Book

	err := db.QueryRow(`
		SELECT id, title, author, isbn, year, price, category, 
		       COALESCE(cover_image, '') as cover_image,
		       COALESCE(description, '') as description,
		       COALESCE(rating, 0) as rating,
		       COALESCE(reviews, 0) as reviews,
		       COALESCE(is_new, false) as is_new,
		       COALESCE(discount, 0) as discount,
		       COALESCE(original_price, 0) as original_price,
		       created_at, updated_at
		FROM books 
		WHERE id = $1
	`, id).Scan(
		&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price,
		&book.Category, &book.CoverImage, &book.Description, &book.Rating,
		&book.Reviews, &book.IsNew, &book.Discount, &book.OriginalPrice,
		&book.Created_At, &book.Updated_At,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

// @Summary Create a new book
// @Description Create a new book
// @Tags Books
// @Accept json
// @Produce json
// @Param book body Book true "Book object"
// @Success 201 {object} Book
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/books [post]
func createBook(c *gin.Context) {
	var newBook Book

	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id int
	var created_At, updated_At time.Time

	err := db.QueryRow(`
		INSERT INTO books (title, author, isbn, year, price, category, 
		                   cover_image, description, rating, reviews, 
		                   is_new, discount, original_price)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at
	`,
		newBook.Title, newBook.Author, newBook.ISBN, newBook.Year, newBook.Price,
		newBook.Category, newBook.CoverImage, newBook.Description, newBook.Rating,
		newBook.Reviews, newBook.IsNew, newBook.Discount, newBook.OriginalPrice,
	).Scan(&id, &created_At, &updated_At)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newBook.ID = id
	newBook.Created_At = created_At
	newBook.Updated_At = updated_At

	c.JSON(http.StatusCreated, newBook)
}

// @Summary Update a book
// @Description Update book details by ID
// @Tags Books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body Book true "Book object"
// @Success 200 {object} Book
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/books/{id} [put]
func updateBook(c *gin.Context) {
	id := c.Param("id")
	var updateBook Book

	if err := c.ShouldBindJSON(&updateBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ID int
	var updatedAt time.Time

	err := db.QueryRow(`
		UPDATE books
		SET title = $1, author = $2, isbn = $3, year = $4, price = $5,
		    category = $6, cover_image = $7, description = $8, rating = $9,
		    reviews = $10, is_new = $11, discount = $12, original_price = $13
		WHERE id = $14
		RETURNING id, updated_at
	`,
		updateBook.Title, updateBook.Author, updateBook.ISBN, updateBook.Year,
		updateBook.Price, updateBook.Category, updateBook.CoverImage,
		updateBook.Description, updateBook.Rating, updateBook.Reviews,
		updateBook.IsNew, updateBook.Discount, updateBook.OriginalPrice, id,
	).Scan(&ID, &updatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updateBook.ID = ID
	updateBook.Updated_At = updatedAt
	c.JSON(http.StatusOK, updateBook)
}

// @Summary Delete a book
// @Description Delete book by ID
// @Tags Books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/books/{id} [delete]
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book deleted successfully"})
}

// @title Book Store API
// @version 1.0
// @description API for managing books in a bookstore
// @host localhost:8080
// @BasePath /api/v1
func main() {
	initDB()
	defer db.Close()

	r := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	// Health check
	r.GET("/health", getHealth)

	// Swagger documentation
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := r.Group("/api/v1")
	{
		// Categories
		api.GET("/categories", getCategories)

		// Books
		api.GET("/books", getAllBooks)                   // Support ?category=fiction
		api.GET("/books/search", searchBooks)            // ?q=keyword
		api.GET("/books/featured", getFeaturedBooks)     // หนังสือแนะนำ
		api.GET("/books/new", getNewBooks)               // หนังสือใหม่
		api.GET("/books/discounted", getDiscountedBooks) // หนังสือลดราคา
		api.GET("/books/:id", getBook)
		api.POST("/books", createBook)
		api.PUT("/books/:id", updateBook)
		api.DELETE("/books/:id", deleteBook)
	}

	log.Println("Server starting on port 8080...")
	r.Run(":8080")
}
