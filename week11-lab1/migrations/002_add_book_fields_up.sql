-- Migration: Add additional fields to books table
-- Version: 002
-- Description: เพิ่มฟิลด์ category, original_price, discount, cover_image, rating, reviews_count, is_new, pages, language, publisher, description

-- =============================================================================
-- STEP 1: Add new columns
-- =============================================================================

-- เพิ่มฟิลด์หมวดหมู่หนังสือ
ALTER TABLE books ADD COLUMN IF NOT EXISTS category VARCHAR(50);

-- เพิ่มฟิลด์ราคาเต็มก่อนลด (สำหรับแสดงส่วนลด)
ALTER TABLE books ADD COLUMN IF NOT EXISTS original_price DECIMAL(10,2);

-- เพิ่มฟิลด์ส่วนลด (เปอร์เซ็นต์ 0-100)
ALTER TABLE books ADD COLUMN IF NOT EXISTS discount INTEGER DEFAULT 0;

-- เพิ่มฟิลด์ URL รูปปกหนังสือ
ALTER TABLE books ADD COLUMN IF NOT EXISTS cover_image VARCHAR(500);

-- เพิ่มฟิลด์คะแนนรีวิว (0.0 - 5.0)
ALTER TABLE books ADD COLUMN IF NOT EXISTS rating DECIMAL(3,2) DEFAULT 0.0;

-- เพิ่มฟิลด์จำนวนรีวิว
ALTER TABLE books ADD COLUMN IF NOT EXISTS reviews_count INTEGER DEFAULT 0;

-- เพิ่มฟิลด์แฟล็กหนังสือใหม่
ALTER TABLE books ADD COLUMN IF NOT EXISTS is_new BOOLEAN DEFAULT false;

-- เพิ่มฟิลด์จำนวนหน้า
ALTER TABLE books ADD COLUMN IF NOT EXISTS pages INTEGER;

-- เพิ่มฟิลด์ภาษา
ALTER TABLE books ADD COLUMN IF NOT EXISTS language VARCHAR(50);

-- เพิ่มฟิลด์สำนักพิมพ์
ALTER TABLE books ADD COLUMN IF NOT EXISTS publisher VARCHAR(255);

-- เพิ่มฟิลด์คำอธิบายหนังสือ
ALTER TABLE books ADD COLUMN IF NOT EXISTS description TEXT;

-- =============================================================================
-- STEP 2: Create indexes for better query performance
-- =============================================================================

-- Index สำหรับการกรองตามหมวดหมู่
CREATE INDEX IF NOT EXISTS idx_books_category ON books(category);

-- Index สำหรับการเรียงตามคะแนน (DESC เพราะต้องการเรียงจากมากไปน้อย)
CREATE INDEX IF NOT EXISTS idx_books_rating ON books(rating DESC);

-- Index สำหรับการกรองหนังสือใหม่
CREATE INDEX IF NOT EXISTS idx_books_is_new ON books(is_new) WHERE is_new = true;

-- Index สำหรับการกรองหนังสือลดราคา
CREATE INDEX IF NOT EXISTS idx_books_discount ON books(discount DESC) WHERE discount > 0;

-- =============================================================================
-- STEP 3: Add comments to columns (สำหรับเอกสาร)
-- =============================================================================

COMMENT ON COLUMN books.category IS 'หมวดหมู่หนังสือ เช่น fiction, psychology, business';
COMMENT ON COLUMN books.original_price IS 'ราคาเต็มก่อนลด';
COMMENT ON COLUMN books.discount IS 'ส่วนลดเป็นเปอร์เซ็นต์ (0-100)';
COMMENT ON COLUMN books.cover_image IS 'URL ของรูปปกหนังสือ';
COMMENT ON COLUMN books.rating IS 'คะแนนรีวิวเฉลี่ย (0.0 - 5.0)';
COMMENT ON COLUMN books.reviews_count IS 'จำนวนรีวิวทั้งหมด';
COMMENT ON COLUMN books.is_new IS 'หนังสือเล่มใหม่หรือไม่';
COMMENT ON COLUMN books.pages IS 'จำนวนหน้าหนังสือ';
COMMENT ON COLUMN books.language IS 'ภาษาของหนังสือ';
COMMENT ON COLUMN books.publisher IS 'สำนักพิมพ์';
COMMENT ON COLUMN books.description IS 'คำอธิบายหนังสือ';

-- =============================================================================
-- Migration completed successfully!
-- =============================================================================
