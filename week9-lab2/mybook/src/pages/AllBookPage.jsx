import { useEffect, useState } from "react";
import { useNavigate } from 'react-router-dom';

const AllBookPage = () => {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchBooks = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await fetch('/api/v1/books/');
      if (!response.ok) {
        throw new Error('Failed to fetch books');
      }
      const fetchedData = await response.json();
      setData(fetchedData);
      console.log('Books data:', fetchedData);
    } catch (err) {
      console.error('Error fetching books:', err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }
  const navigate = useNavigate();

  const handleNavigateAddBook = () => {
    navigate('/store-manager/add-book'); // <-- ใส่ลิงก์ของคุณที่นี่
  };
  useEffect(() => {
    fetchBooks();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        Loading...
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen text-red-600">
        Error: {error}
      </div>
    );
  }

  // แสดงผลข้อมูลเมื่อสำเร็จ
  return (
    <div className="min-h-screen bg-gray-100 flex flex-col">
      {/* 1. Header (เหมือนเดิม) */}
      <header className="bg-green-600 text-white p-4 shadow-md">
        <h1 className="text-2xl font-bold">Manage List Book</h1>
      </header>

      {/* 2. Main Content (เหมือนเดิม) */}
      <main className="flex-1 container mx-auto p-6">

        {/* 2a. Content Header (เหมือนเดิม) */}
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-3xl font-semibold text-gray-800">List Book</h2>
          <button
          onClick={handleNavigateAddBook} 
          className="bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-4 rounded transition duration-200">
            Add Book
          </button>
        </div>

        {/* 2b. Book List (นี่คือส่วนที่แก้ไข) */}
        <div className="bg-white shadow-lg rounded-lg overflow-hidden">
          <ul>
            {data.length > 0 ? (
              data.map((book) => (

                // *** 1. แก้ไข li ให้ใช้ Flexbox ***
                <li
                  key={book.id}
                  className="p-4 border-b border-gray-200 last:border-b-0 flex justify-between items-center"
                >
                  {/* ห่อข้อมูลหนังสือไว้ใน div ด้านซ้าย */}
                  <div>
                    <p className="text-sm text-gray-500">ID: {book.id}</p>
                    <p className="text-lg font-semibold">{book.title}</p>
                    <p className='text-sm text-gray-500'>โดย {book.author}</p>
                    <p className="text-sm text-gray-700">ราคา: ฿{book.price}</p>
                    <p className="text-sm text-gray-700">รายละเอียด: {book.description}</p>
                    <p className="text-sm text-gray-700">หมวดหมู่: {book.category}</p>
                    <p className="text-sm text-gray-700">ปีที่เผยแพร่: {book.year}</p>
                    <p className="text-sm text-gray-700">สถานะ: {book.status}</p>
                  </div>

                  {/* *** 2. ย้ายปุ่มมาไว้ใน div ด้านขวา *** */}
                  {/* flex-shrink-0 เพื่อไม่ให้ปุ่มหดตัว */}
                  <div className="flex space-x-2 flex-shrink-0">

                    {/* 1. แก้จาก rounded-l-lg เป็น rounded-lg */}
                    <button className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-lg transition duration-200">
                      Edit
                    </button>

                    {/* 2. ลบ <span className="mx-2"></span> ที่ไม่จำเป็นออก */}

                    {/* 3. แก้จาก rounded-r-lg เป็น rounded-lg */}
                    <button className="bg-red-600 hover:bg-red-700 text-white font-bold py-2 px-4 rounded-lg transition duration-200">
                      Delete
                    </button>
                  </div>

                </li>
              ))
            ) : (
              <li className="p-6 text-center text-gray-500">
                ไม่พบข้อมูลหนังสือ
              </li>
            )}

            {/* *** 3. ลบปุ่มที่อยู่ท้าย <ul> ออก *** */}

          </ul>
        </div>
      </main>

      {/* 3. Footer (เหมือนเดิม) */}
      <footer className="bg-green-600 p-4 mt-auto">
      </footer>
    </div>
  );
}

export default AllBookPage;