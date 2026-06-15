Berikut adalah draf dokumentasi README.md yang bersih, terstruktur, dan profesional untuk API manajemen komunitas badminton Anda. Dokumentasi ini memisahkan endpoint publik dengan endpoint yang memerlukan autentikasi JWT (protected), lengkap dengan contoh struktur data (payload) berdasarkan struct Go yang Anda miliki.

Markdown
# 🏸 Badminton Community Management API

API berbasis Go (Gin Framework) dan GORM yang dirancang untuk mengelola komunitas badminton (mabar), penjadwalan lapangan, pencatatan transaksi *shuttlecock*, statistik pertandingan, hingga kalkulasi tagihan iuran bulanan/sesi secara otomatis.

## 🚀 Fitur Utama
* **Autentikasi & Multi-tenant:** Pendaftaran Admin dan login berbasis JWT Token yang mengikat ke `community_id`.
* **Manajemen Lapangan & Jadwal:** Booking lapangan otomatis beserta perhitungan biaya *court*.
* **Logistik Shuttlecock:** Pencatatan pembelian stok baru, pelacakan sisa kok, dan pencatatan kok yang dikembalikan.
* **Pencatatan Pertandingan:** Mendukung penentuan format skor dan relasi *Many-to-Many* pemain dalam satu *match*.
* **Sistem Billing Otomatis:** Kalkulasi biaya kok berdasarkan pemakaian riil individu per sesi mabar dicampur dengan pembagian biaya sewa lapangan.

