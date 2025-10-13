import React, { useState, useEffect } from 'react';
import BookCard from './BookCard';

const NewBooks = () => {
  // กำหนด State สำหรับจัดการข้อมูล
  const [newBooks, setNewBooks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchNewBooks = async () => {
      try {
        setLoading(true);
        
        // เรียก API เพื่อดึงข้อมูลหนังสือใหม่
        const response = await fetch('http://localhost:8080/api/v1/books/new');

        if (!response.ok) {
          throw new Error('Failed to fetch new books');
        }

        const data = await response.json();
        setNewBooks(data);
        setError(null);
        
      } catch (err) {
        setError(err.message);
        console.error('Error fetching new books:', err);
        
      } finally {
        setLoading(false);
      }
    };

    // เรียกใช้ฟังก์ชันดึงข้อมูล
    fetchNewBooks();
  }, []); // [] = dependency array ว่าง = รันครั้งเดียว

  // กรณีกำลังโหลดข้อมูล
  if (loading) {
    return (
      <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
        <div className="text-center py-8 col-span-full">
          Loading new arrivals...
        </div>
      </div>
    );
  }

  // กรณีเกิดข้อผิดพลาด
  if (error) {
    return (
      <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
        <div className="text-center py-8 col-span-full text-red-600">
          Error: {error}
        </div>
      </div>
    );
  }

  // กรณีแสดงผลข้อมูลปกติ
  return (
    <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
      {newBooks.map(book => (
        <BookCard 
          key={book.id} 
          book={book} 
        />
      ))}
    </div>
  );
};

export default NewBooks;