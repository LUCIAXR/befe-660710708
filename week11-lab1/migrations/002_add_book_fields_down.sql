-- Rollback Migration: Remove additional fields from books table
-- Version: 002
-- Description: ลบฟิลด์ที่เพิ่มเข้ามาใน migration 002 (ย้อนกลับสู่สถานะเดิม)

-- =============================================================================
-- WARNING: การรัน rollback จะลบข้อมูลในคอลัมน์เหล่านี้ทิ้ง!
-- กรุณา backup ข้อมูลก่อนรัน rollback
-- =============================================================================

-- =============================================================================
-- STEP 1: Drop indexes
-- =============================================================================

DROP INDEX IF EXISTS idx_books_category;
DROP INDEX IF EXISTS idx_books_rating;
DROP INDEX IF EXISTS idx_books_is_new;
DROP INDEX IF EXISTS idx_books_discount;

-- =============================================================================
-- STEP 2: Drop columns (ข้อมูลในคอลัมน์เหล่านี้จะถูกลบทิ้ง!)
-- =============================================================================

ALTER TABLE books DROP COLUMN IF EXISTS category;
ALTER TABLE books DROP COLUMN IF EXISTS original_price;
ALTER TABLE books DROP COLUMN IF EXISTS discount;
ALTER TABLE books DROP COLUMN IF EXISTS cover_image;
ALTER TABLE books DROP COLUMN IF EXISTS rating;
ALTER TABLE books DROP COLUMN IF EXISTS reviews_count;
ALTER TABLE books DROP COLUMN IF EXISTS is_new;
ALTER TABLE books DROP COLUMN IF EXISTS pages;
ALTER TABLE books DROP COLUMN IF EXISTS language;
ALTER TABLE books DROP COLUMN IF EXISTS publisher;
ALTER TABLE books DROP COLUMN IF EXISTS description;

-- =============================================================================
-- Rollback completed successfully!
-- Table structure has been reverted to version 001
-- =============================================================================
