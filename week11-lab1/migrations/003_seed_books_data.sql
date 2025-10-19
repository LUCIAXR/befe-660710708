-- Seed Data: Sample books for testing
-- Version: 003
-- Description: เพิ่มข้อมูลตัวอย่างหนังสือ 15 เล่ม (จาก Mock Data ของ Frontend)

-- =============================================================================
-- IMPORTANT: รันไฟล์นี้หลังจากรัน 002_add_book_fields_up.sql แล้วเท่านั้น
-- =============================================================================

-- ลบข้อมูลเก่าก่อน (ถ้ามี)
TRUNCATE TABLE books RESTART IDENTITY CASCADE;

-- =============================================================================
-- Insert sample books data
-- =============================================================================

INSERT INTO books (
    title, author, category, price, original_price, discount,
    cover_image, rating, reviews_count, isbn, pages, language,
    publisher, year, description, is_new
) VALUES
-- 1. The Great Gatsby
(
    'The Great Gatsby',
    'F. Scott Fitzgerald',
    'fiction',
    299.00,
    399.00,
    25,
    '/images/books/gatsby.jpg',
    4.5,
    234,
    '978-0-7432-7356-5',
    180,
    'English',
    'Scribner',
    1925,
    'นวนิยายคลาสสิกของอเมริกาที่เล่าเรื่องราวของเจย์ แกตส์บี้ และความฝันอเมริกันในยุค 1920s',
    false
),

-- 2. 1984
(
    '1984',
    'George Orwell',
    'fiction',
    350.00,
    NULL,
    0,
    '/images/books/1984.jpg',
    4.8,
    512,
    '978-0-452-28423-4',
    328,
    'English',
    'Signet Classic',
    1949,
    'นวนิยายดิสโทเปียที่พรรณนาถึงสังคมเผด็จการในอนาคต',
    true
),

-- 3. To Kill a Mockingbird
(
    'To Kill a Mockingbird',
    'Harper Lee',
    'fiction',
    320.00,
    NULL,
    0,
    '/images/books/mockingbird.jpg',
    4.6,
    189,
    '978-0-06-112008-4',
    324,
    'English',
    'Harper Perennial',
    1960,
    'เรื่องราวการเติบโตและความอยุติธรรมทางเชื้อชาติในอเมริกาใต้',
    false
),

-- 4. Sapiens
(
    'Sapiens: A Brief History of Humankind',
    'Yuval Noah Harari',
    'non-fiction',
    450.00,
    550.00,
    18,
    '/images/books/sapiens.jpg',
    4.7,
    892,
    '978-0-06-231609-7',
    464,
    'Thai',
    'Harper',
    2014,
    'ประวัติศาสตร์ของมนุษยชาติตั้งแต่ยุคหินจนถึงปัจจุบัน',
    false
),

-- 5. The Alchemist
(
    'The Alchemist',
    'Paulo Coelho',
    'fiction',
    280.00,
    NULL,
    0,
    '/images/books/alchemist.jpg',
    4.3,
    1523,
    '978-0-06-231500-7',
    208,
    'Thai',
    'HarperOne',
    1988,
    'นวนิยายผจญภัยเชิงปรัชญาเกี่ยวกับการค้นหาโชคชะตาของตนเอง',
    false
),

-- 6. Thinking, Fast and Slow
(
    'Thinking, Fast and Slow',
    'Daniel Kahneman',
    'psychology',
    420.00,
    NULL,
    0,
    '/images/books/thinking.jpg',
    4.4,
    445,
    '978-0-374-53355-7',
    512,
    'English',
    'FSG',
    2011,
    'การสำรวจระบบความคิดสองระบบที่ขับเคลื่อนวิธีที่เราคิด',
    true
),

-- 7. The Art of War
(
    'The Art of War',
    'Sun Tzu',
    'history',
    250.00,
    350.00,
    29,
    '/images/books/artofwar.jpg',
    4.6,
    667,
    '978-1-59030-225-6',
    273,
    'Thai',
    'Shambhala',
    -500,
    'ตำราพิชัยสงครามจีนโบราณที่ยังใช้ได้ในยุคปัจจุบัน',
    false
),

-- 8. Clean Code
(
    'Clean Code',
    'Robert C. Martin',
    'technology',
    580.00,
    NULL,
    0,
    '/images/books/cleancode.jpg',
    4.5,
    234,
    '978-0-13-235088-2',
    464,
    'English',
    'Prentice Hall',
    2008,
    'คู่มือการเขียนโค้ดที่สะอาดและบำรุงรักษาได้',
    false
),

-- 9. The Lean Startup
(
    'The Lean Startup',
    'Eric Ries',
    'business',
    380.00,
    NULL,
    0,
    '/images/books/leanstartup.jpg',
    4.2,
    556,
    '978-0-307-88789-4',
    336,
    'Thai',
    'Crown Business',
    2011,
    'วิธีการสร้างและบริหารสตาร์ทอัพอย่างมีประสิทธิภาพ',
    false
),

-- 10. The Power of Now
(
    'The Power of Now',
    'Eckhart Tolle',
    'psychology',
    320.00,
    NULL,
    0,
    '/images/books/powerofnow.jpg',
    4.4,
    889,
    '978-1-57731-480-6',
    236,
    'Thai',
    'New World Library',
    1997,
    'คู่มือการใช้ชีวิตอยู่กับปัจจุบันและการตื่นรู้ทางจิตวิญญาณ',
    true
),

-- 11. Atomic Habits
(
    'Atomic Habits',
    'James Clear',
    'psychology',
    390.00,
    450.00,
    13,
    '/images/books/atomichabits.jpg',
    4.8,
    2341,
    '978-0-7352-1129-2',
    320,
    'Thai',
    'Avery',
    2018,
    'วิธีสร้างนิสัยที่ดีและกำจัดนิสัยที่ไม่ดี',
    false
),

-- 12. The 7 Habits of Highly Effective People
(
    'The 7 Habits of Highly Effective People',
    'Stephen R. Covey',
    'business',
    420.00,
    NULL,
    0,
    '/images/books/7habits.jpg',
    4.5,
    1456,
    '978-0-7432-6951-3',
    432,
    'Thai',
    'Free Press',
    1989,
    '7 นิสัยสำหรับการพัฒนาตนเองและความสำเร็จ',
    false
),

-- 13. The Subtle Art of Not Giving a F*ck
(
    'The Subtle Art of Not Giving a F*ck',
    'Mark Manson',
    'psychology',
    340.00,
    NULL,
    0,
    '/images/books/subtleart.jpg',
    4.1,
    1789,
    '978-0-06-245771-4',
    224,
    'Thai',
    'HarperOne',
    2016,
    'แนวทางการใช้ชีวิตแบบตรงไปตรงมาเพื่อชีวิตที่ดีขึ้น',
    false
),

-- 14. Rich Dad Poor Dad
(
    'Rich Dad Poor Dad',
    'Robert T. Kiyosaki',
    'business',
    360.00,
    420.00,
    14,
    '/images/books/richdad.jpg',
    4.3,
    2567,
    '978-1-61268-019-0',
    336,
    'Thai',
    'Plata Publishing',
    1997,
    'บทเรียนการเงินจากพ่อสองคนที่มีมุมมองต่างกัน',
    false
),

-- 15. The Da Vinci Code
(
    'The Da Vinci Code',
    'Dan Brown',
    'fiction',
    380.00,
    NULL,
    0,
    '/images/books/davinci.jpg',
    4.0,
    3456,
    '978-0-385-50420-1',
    689,
    'Thai',
    'Doubleday',
    2003,
    'นวนิยายลึกลับระทึกขวัญเกี่ยวกับรหัสลับในงานศิลปะ',
    false
);

-- =============================================================================
-- Verify data insertion
-- =============================================================================

-- แสดงจำนวนหนังสือที่เพิ่มเข้าไป
SELECT 'Total books inserted:' as info, COUNT(*) as count FROM books;

-- แสดงข้อมูลแบบสรุป
SELECT
    category,
    COUNT(*) as total_books,
    ROUND(AVG(rating), 2) as avg_rating,
    ROUND(AVG(price), 2) as avg_price
FROM books
GROUP BY category
ORDER BY total_books DESC;

-- =============================================================================
-- Seed data completed successfully!
-- =============================================================================
