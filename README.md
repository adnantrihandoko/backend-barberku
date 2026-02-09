<!-- Berikut adalah **Software Requirement Specification (SRS)** dokumen lengkap standar profesional untuk aplikasi Antrian Barbershop ("BarberQueue Pro"). Dokumen ini dirancang agar siap digunakan oleh tim pengembang (frontend, backend, QA) dan stakeholder.

Dokumen ini disesuaikan dengan stack pilihan Anda: **Vue.js (Frontend), Python (Backend), dan PostgreSQL (Database).** -->


# SOFTWARE REQUIREMENT SPECIFICATION (SRS)

**Project Name:** BarberQueue Pro System
**Version:** 1.0
**Status:** Draft Akhir

---

## 1. Pendahuluan (Introduction)

### 1.1 Tujuan

Tujuan dari dokumen ini adalah mendefinisikan kebutuhan fungsional dan non-fungsional untuk pengembangan sistem manajemen antrian barbershop berbasis web. Sistem ini bertujuan untuk meminimalisir waktu tunggu pelanggan di lokasi, mengelola jadwal barber, dan mencatat transaksi keuangan secara digital.

### 1.2 Lingkup Masalah (Scope)

Aplikasi akan terdiri dari dua antarmuka utama:

1. **Client Side (Pelanggan):** Untuk melihat antrian, booking jadwal, dan memilih barber.
2. **Admin/Barber Side (Dashboard):** Untuk manajemen antrian operasional, manajemen layanan, dan laporan pendapatan.

### 1.3 Definisi & Istilah

* **Booking:** Reservasi slot waktu di masa depan.
* **Walk-in:** Pelanggan yang datang langsung tanpa booking (admin yang menginput ke sistem).
* **Queue/Antrian:** Urutan pengerjaan layanan berdasarkan waktu booking atau kedatangan.

---

## 2. Deskripsi Keseluruhan (Overall Description)

### 2.1 Perspektif Produk

Sistem ini beroperasi dengan arsitektur **Three-Tier**:

* **Presentation Layer:** Vue.js (SPA - Single Page Application).
* **Application Layer:** Python (Disarankan menggunakan framework seperti FastAPI atau Django REST Framework) sebagai RESTful API.
* **Data Layer:** PostgreSQL.

### 2.2 Karakteristik Pengguna (User Classes)

| Aktor | Deskripsi |
| --- | --- |
| **Pelanggan (Customer)** | Pengguna umum yang mendaftar, melihat antrian real-time, dan melakukan booking. |
| **Barber (Tukang Cukur)** | Staf yang melihat jadwal kerja mereka sendiri dan mengubah status antrian (Start/Finish). |
| **Administrator/Owner** | Memiliki akses penuh ke manajemen user, layanan, harga, dan laporan keuangan. |

---

## 3. Spesifikasi Kebutuhan Fungsional (Functional Requirements)

### 3.1 Modul Autentikasi & Otorisasi

* **FR-AUTH-01:** Sistem harus memungkinkan pelanggan mendaftar menggunakan email/nomor HP.
* **FR-AUTH-02:** Sistem harus menggunakan JSON Web Token (JWT) untuk manajemen sesi.
* **FR-AUTH-03:** Sistem harus membedakan akses menu berdasarkan role (Customer, Barber, Admin).

### 3.2 Modul Pelanggan (Customer Side - Vue.js)

* **FR-CUST-01 (Dashboard Antrian):** Pelanggan dapat melihat status antrian saat ini secara real-time (siapa yang sedang dicukur, berapa orang yang menunggu).
* **FR-CUST-02 (Booking):** Pelanggan dapat memilih layanan, memilih barber spesifik, dan memilih slot waktu yang tersedia.
* **FR-CUST-03 (Estimasi Waktu):** Sistem harus menampilkan estimasi waktu tunggu berdasarkan durasi rata-rata layanan yang dipilih antrian di depannya.
* **FR-CUST-04 (Notifikasi):** (Opsional) Integrasi WhatsApp/Email saat giliran hampir tiba.

### 3.3 Modul Operasional Barber (Barber Side)

* **FR-BARB-01 (Status Update):** Barber dapat mengubah status antrian pelanggan: *Waiting*  *Processing*  *Completed*  *Cancelled*.
* **FR-BARB-02 (Availability):** Barber dapat mematikan ketersediaan mereka (mode istirahat/pulang) yang akan langsung berdampak pada slot booking di sisi pelanggan.

### 3.4 Modul Admin (Management)

* **FR-ADM-01 (Service Management):** CRUD (Create, Read, Update, Delete) data layanan (nama layanan, harga, estimasi durasi).
* **FR-ADM-02 (Walk-in Input):** Admin dapat memasukkan pelanggan *walk-in* secara manual ke dalam slot antrian yang kosong.
* **FR-ADM-03 (Reporting):** Melihat laporan pendapatan harian, mingguan, dan bulanan, serta performa per barber.

---

## 4. Spesifikasi Database (PostgreSQL Schema)

Database harus dirancang untuk menjaga integritas data transaksional (ACID). Berikut adalah entitas utamanya:

### 4.1 Tabel Utama

1. **Users**
* `id` (UUID/BigInt, PK)
* `email`, `password_hash`, `full_name`, `phone`
* `role` (Enum: 'admin', 'barber', 'customer')


2. **Services**
* `id` (PK)
* `name` (Varchar)
* `price` (Decimal)
* `duration_minutes` (Int) - *Penting untuk kalkulasi antrian.*


3. **Appointments (Queue)**
* `id` (PK)
* `customer_id` (FK -> Users)
* `barber_id` (FK -> Users)
* `service_id` (FK -> Services)
* `scheduled_time` (Timestamp)
* `status` (Enum: 'pending', 'confirmed', 'in_progress', 'completed', 'cancelled')
* `payment_status` (Boolean)



---

## 5. Spesifikasi Teknis (Technical Constraints)

### 5.1 Backend (Python)

* **Framework:** Disarankan menggunakan **FastAPI** (untuk performa tinggi dan handling *concurrency* yang baik) atau **Django** (jika butuh admin panel bawaan yang cepat).
* **ORM:** Gunakan **SQLAlchemy** (jika FastAPI) atau **Django ORM**.
* **Validation:** Pydantic untuk validasi data input/output.
* **Scheduler:** Celery (opsional) untuk membatalkan booking yang telat otomatis atau mengirim reminder.

### 5.2 Frontend (Vue.js)

* **Version:** Vue 3 dengan Composition API.
* **State Management:** Pinia (pengganti Vuex untuk Vue 3).
* **UI Framework:** Tailwind CSS atau Vuetify untuk mempercepat styling.
* **HTTP Client:** Axios.
* **Real-time:** Socket.io-client (jika menggunakan WebSocket untuk update antrian live) atau *Polling* setiap 30 detik.

### 5.3 Database (PostgreSQL)

* Menggunakan fitur **Indexing** pada kolom `scheduled_time` dan `barber_id` untuk mempercepat query pencarian slot.
* Menggunakan **Constraints** untuk mencegah *double booking* pada barber yang sama di jam yang sama.

---

## 6. Kebutuhan Non-Fungsional (NFR)

1. **Performance:** API response time harus di bawah 200ms untuk query standar.
2. **Concurrency:** Sistem harus mampu menangani minimal 50 user yang mengakses jadwal secara bersamaan tanpa *race condition* (dua orang booking slot yang sama).
3. **Security:**
* Password di-hash menggunakan **Argon2** atau **BCrypt**.
* Perlindungan terhadap SQL Injection (ditangani oleh ORM).
* CORS policy yang ketat (hanya domain frontend yang boleh akses API).


4. **Usability:** Desain mobile-first, karena mayoritas pelanggan akan mengakses via HP.

---

## 7. Roadmap Pengembangan (Saran)

Untuk memulai proyek ini dengan stack Vue + Python + Postgres, berikut adalah langkah logis pengerjaannya:

1. **Minggu 1:** Desain Database (ERD) & Setup Environment (Docker Compose disarankan untuk menyatukan Python & Postgres).
2. **Minggu 2:** Pengembangan API Backend (Auth, CRUD Services, Booking Logic).
3. **Minggu 3:** Pengembangan Frontend Vue.js (Halaman Login, Dashboard, Booking Form).
4. **Minggu 4:** Integrasi API ke Frontend & Testing alur antrian.

Apakah Anda ingin saya buatkan **Struktur Folder Proyek** (Scaffolding) untuk Vue dan Python-nya agar Anda bisa langsung mulai coding?