package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv" // ต้องเพิ่มเข้ามาเพื่อใช้แปลง string เป็น int
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // ต้องเพิ่ม driver สำหรับ postgresql
)

// Book struct เพื่อแสดงข้อมูลหนังสือ
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

// getEnv ดึงค่าจาก Environment Variable หรือใช้ค่า Default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

var db *sql.DB

// initDB ตั้งค่าและเชื่อมต่อกับฐานข้อมูล PostgreSQL
func initDB() {
	var err error

	host := getEnv("DB_HOST", "localhost") // ใช้ localhost เป็นค่า default
	name := getEnv("DB_NAME", "bookstore")
	user := getEnv("DB_USER", "bookstore_user")
	password := getEnv("DB_PASSWORD", "your_strong_password")
	port := getEnv("DB_PORT", "5432")

	// แก้ไข: db_name ต้องเป็น dbname
	conSt := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)

	// ต้องใช้ driver "postgres" จาก package github.com/lib/pq
	db, err = sql.Open("postgres", conSt)

	// ตั้งค่า Connection Pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err != nil {
		log.Fatalf("failed to open database connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("successfully connected to database")
}

// getAllBooks ดึงรายการหนังสือทั้งหมด หรือกรองตามปี
func getAllBooks(c *gin.Context) {
	yearInput := c.Query("year")
	var rows *sql.Rows
	var err error

	// แก้ไข: ใช้ Prepared Statement ที่ถูกต้อง
	if yearInput != "" {
		// ถ้ามีการระบุปี ให้กรอง
		rows, err = db.Query("SELECT id, title, author, isbn, year, price, created_at, updated_at FROM books WHERE year = $1", yearInput)
	} else {
		// ถ้าไม่ระบุปี ให้ดึงทั้งหมด
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
		// แก้ไข: ชื่อคอลัมน์เป็น updated_at
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			log.Printf("error scanning book: %v", err)
			continue // ข้ามแถวที่มีปัญหา
		}
		books = append(books, book)
	}
	if books == nil {
		books = []Book{} // ส่ง Array ว่างกลับไปแทน nil
	}

	c.JSON(http.StatusOK, books)
}

// getBook ดึงข้อมูลหนังสือตาม ID
func getBook(c *gin.Context) {
	id := c.Param("id")
	var book Book

	// QueryRow ใช้เมื่อคาดว่าจะได้ผลลัพธ์ 0 หรือ 1 แถว
	// แก้ไข: เพิ่มคอลัมน์ที่จำเป็นทั้งหมดลงใน SELECT
	err := db.QueryRow("SELECT id, title, author, isbn, year, price, created_at, updated_at FROM books WHERE id = $1", id).
		Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price, &book.CreatedAt, &book.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

// createBook สร้างหนังสือเล่มใหม่
func createBook(c *gin.Context) {
	var newBook Book

	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

// updateBook อัปเดตข้อมูลหนังสือตาม ID
func updateBook(c *gin.Context) {
	idStr := c.Param("id")
	var updateBook Book

	if err := c.ShouldBindJSON(&updateBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// แก้ไข: ต้องแปลง ID ที่เป็น string กลับเป็น int เพื่อกำหนดให้กับ struct
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var updatedAt time.Time
	// ใช้ QueryRow กับ RETURNING เพื่อให้แน่ใจว่าได้อัปเดตไปแล้ว
	err = db.QueryRow(
		`UPDATE books
		 SET title = $1, author = $2, isbn = $3, year = $4, price = $5, updated_at = NOW()
		 WHERE id = $6
		 RETURNING updated_at`,
		updateBook.Title, updateBook.Author, updateBook.ISBN,
		updateBook.Year, updateBook.Price, idStr,
	).Scan(&updatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// แก้ไข: กำหนดค่า ID และ UpdatedAt กลับเข้าไปใน struct
	updateBook.ID = id
	updateBook.UpdatedAt = updatedAt
	c.JSON(http.StatusOK, updateBook)
}

// deleteBook ลบหนังสือตาม ID
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

func main() {
	// ต้องเรียก initDB ก่อนใช้ db
	initDB()
	defer db.Close()

	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		err := db.Ping()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "unhealthy", "error": err.Error()}) // แก้ไข: ใช้ Error()
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

	// รัน API Server
	r.Run(":8080")

}
