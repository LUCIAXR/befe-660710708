# Database Migrations

คู่มือการรัน Database Migrations สำหรับโปรเจกต์ Bookstore API

## 📋 Migration Files

| File | Description |
|------|-------------|
| `002_add_book_fields_up.sql` | เพิ่มฟิลด์ใหม่ให้กับตาราง books |
| `002_add_book_fields_down.sql` | ย้อนกลับการเปลี่ยนแปลง (rollback) |
| `003_seed_books_data.sql` | ข้อมูลตัวอย่างหนังสือ 15 เล่ม |

---

## 🚀 วิธีการรัน Migration

```bash
# Navigate to project directory
cd /workspaces/befe-00000/week10-lab3

# รัน migration (UP)
docker exec -i bookstore_postgres psql -U bookstore_user -d bookstore < migrations/002_add_book_fields_up.sql

# หรือใช้ตัวเลือกนี้
docker exec bookstore_postgres psql -U bookstore_user -d bookstore -f /migrations/002_add_book_fields_up.sql
```

## ⏪ วิธีการ Rollback (ย้อนกลับ)

⚠️ **คำเตือน:** การ rollback จะลบข้อมูลในคอลัมน์ที่เพิ่มเข้ามา

```bash
# Rollback migration
docker exec -i bookstore_postgres psql -U bookstore_user -d bookstore < migrations/002_add_book_fields_down.sql
```

---

## 📊 ตรวจสอบผลลัพธ์

### 1. ดูโครงสร้างตาราง

```bash
docker exec bookstore_postgres psql -U bookstore_user -d bookstore -c "\d books"
```

### 2. ดูคอลัมน์ที่เพิ่มเข้ามา

```bash
docker exec bookstore_postgres psql -U bookstore_user -d bookstore -c "
SELECT column_name, data_type, is_nullable, column_default
FROM information_schema.columns
WHERE table_name = 'books'
ORDER BY ordinal_position;
"
```

### 3. ดูข้อมูลในตาราง

```bash
docker exec bookstore_postgres psql -U bookstore_user -d bookstore -c "SELECT * FROM books LIMIT 5;"
```


## 🌱 เพิ่มข้อมูลตัวอย่าง (Seed Data)

หลังจากรัน migration แล้ว ให้รัน seed data เพื่อเพิ่มข้อมูลหนังสือ 15 เล่ม

```bash
docker exec -i bookstore_postgres psql -U bookstore_user -d bookstore < migrations/003_seed_books_data.sql
```

## UPDATE Go Struct
```go
  type Book struct {
      ID            int       `json:"id"`
      Title         string    `json:"title"`
      Author        string    `json:"author"`
      ISBN          string    `json:"isbn"`
      Year          int       `json:"year"`
      Price         float64   `json:"price"`

      // ฟิลด์ใหม่
      Category      string    `json:"category"`
      OriginalPrice *float64  `json:"original_price,omitempty"`
      Discount      int       `json:"discount"`
      CoverImage    string    `json:"cover_image"`
      Rating        float64   `json:"rating"`
      ReviewsCount  int       `json:"reviews_count"`
      IsNew         bool      `json:"is_new"`
      Pages         *int      `json:"pages,omitempty"`
      Language      string    `json:"language"`
      Publisher     string    `json:"publisher"`
      Description   string    `json:"description"`

      CreatedAt     time.Time `json:"created_at"`
      UpdatedAt     time.Time `json:"updated_at"`
  }
```

  ## ปรับปรุงและเพิ่ม API Endpoints ใหม่
  ### ปรับปรุง
  ```bash
  - GET /api/v1/books?category=fiction - กรองตามหมวดหมู่
  ```
  ### เพิ่ม
  ```bash
  - GET /api/v1/categories - ค้นหา category
  - GET /api/v1/books/search?q=keyword - ค้นหา
  - GET /api/v1/books/featured - หนังสือแนะนำ
  - GET /api/v1/books/new - หนังสือใหม่
  - GET /api/v1/books/discounted - หนังสือลดราคา
  ```
