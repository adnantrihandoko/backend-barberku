# BarberKu — Panduan Lengkap Menjalankan dari Awal

---

## 1. Prasyarat

Pastikan tools berikut sudah terinstall:

| Tool | Minimal Versi | Cek Instalasi |
|------|--------------|---------------|
| Flutter | 3.22+ | `flutter --version` |
| Dart | 3.5+ | `dart --version` |
| Docker & Docker Compose | terbaru | `docker --version && docker compose version` |
| Git | terbaru | `git --version` |

---

## 2. Clone Repository

```bash
mkdir ~/Projects/barberku
cd ~/Projects/barberku

# Clone backend
git clone https://github.com/adnantrihandoko/backend-barberku.git backend-barber

# Clone Flutter app
git clone https://github.com/adnantrihandoko/frontend-barberku.git barberku_app
```

---

## 3. Setup Firebase Project

### 3.1 Buat Project

1. Buka https://console.firebase.google.com/
2. Klik **"Add project"**
3. Masukkan nama: `barberku`
4. Matikan **Google Analytics** (tidak diperlukan)
5. Klik **"Create project"**
6. Tunggu hingga selesai, klik **"Continue"**

### 3.2 Dapatkan Project ID dan Sender ID

1. Di Firebase Console, klik ikon gear ⚙️ → **Project Settings**
2. Tab **General**, catat:
   - **Project ID**: `barberku-xxxxx`
   - **Project Number** (digunakan sebagai Sender ID)
3. Tab **Cloud Messaging**, catat:
   - **Server key** (akan dipakai di backend)

### 3.3 Daftarkan Aplikasi Android

1. Di Firebase Console → **Project Overview**
2. Klik ikon Android ➕ (**Add app**)
3. **Android package name**: `com.barberku.app`
4. **App nickname** (opsional): `BarberKu Android`
5. Klik **"Register app"**
6. Klik **"Download google-services.json"**
7. Letakkan file di: `barberku_app/android/app/google-services.json`

> **Untuk production**: Generate SHA-1 fingerprint untuk release:
> ```bash
> keytool -list -v -keystore ~/.android/debug.keystore \
>   -alias androiddebugkey -storepass android -keypass android
> ```
> Salin SHA-1, paste ke Firebase Console → Project Settings → Your apps → Android app

### 3.4 Daftarkan Aplikasi iOS

1. Di Firebase Console → **Project Overview**
2. Klik ikon iOS ➕ (**Add app**)
3. **iOS bundle ID**: `com.barberku.app`
4. **App Store ID**: kosongkan
5. Klik **"Register app"**
6. Klik **"Download GoogleService-Info.plist"**
7. Buka `barberku_app/ios/Runner.xcworkspace` di Xcode
8. Drag file `GoogleService-Info.plist` ke folder **Runner** di Xcode
9. Centang **"Copy items if needed"**, pastikan target **Runner** tercentang

### 3.5 Aktifkan Push Notification Capability (iOS)

1. Buka `barberku_app/ios/Runner.xcworkspace` di Xcode
2. Pilih target **Runner** → tab **Signing & Capabilities**
3. Klik **+ Capability**
4. Cari dan tambahkan **Push Notifications**
5. (Opsional) Tambahkan juga **Background Modes** → centang **Remote notifications**

---

## 4. Konfigurasi Flutter — Firebase

### 4.1 Install FlutterFire CLI

```bash
dart pub global activate flutterfire_cli
```

> Pastikan `~/.pub-cache/bin` ada di PATH. Tambahkan ke `~/.bashrc`:
> ```bash
> export PATH="$PATH:$HOME/.pub-cache/bin"
> ```

### 4.2 Generate Firebase Options

```bash
cd ~/Projects/barberku/barberku_app

flutterfire configure \
  --project=barberku-xxxxx \
  --android-package-name=com.barberku.app \
  --ios-bundle-id=com.barberku.app
```

Ini akan membuat file `lib/firebase_options.dart` secara otomatis.

### 4.3 Install Flutter Dependencies

```bash
cd ~/Projects/barberku/barberku_app
flutter pub get
```

### 4.4 Update main.dart

Buka `lib/main.dart` dan pastikan Firebase diinisialisasi dengan options:

```dart
import 'package:firebase_core/firebase_core.dart';
import 'firebase_options.dart'; // file hasil generate

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  await Firebase.initializeApp(
    options: DefaultFirebaseOptions.currentPlatform, // tambahkan ini
  );

  final fcmService = FCMService();
  await fcmService.initialize();

  runApp(
    ProviderScope(
      overrides: [
        fcmServiceProvider.overrideWithValue(fcmService),
      ],
      child: const BarberKuApp(),
    ),
  );
}
```

### 4.5 Update Android Build Files

Buka `barberku_app/android/build.gradle`:

```gradle
buildscript {
    ext.kotlin_version = '1.9.22'
    dependencies {
        classpath 'com.android.tools.build:gradle:8.1.0'
        classpath "org.jetbrains.kotlin:kotlin-gradle-plugin:$kotlin_version"
        classpath 'com.google.gms:google-services:4.4.0'  // tambahkan ini
    }
}
```

Buka `barberku_app/android/app/build.gradle`, tambahkan di **paling bawah**:

```gradle
apply plugin: 'com.google.gms.google-services'
```

---

## 5. Konfigurasi Backend

### 5.1 Set Environment Variables

```bash
cd ~/Projects/barberku/backend-barber

# Copy env template
cp .env.example .env

# Edit .env
nano .env
```

Isi file `.env`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=barber
DB_PASS=secret
DB_NAME=barbershop
SERVER_PORT=8080
WS_ORIGIN_ALLOWED=*
JWT_SECRET=barberku-secret-key-change-in-production
FCM_SERVER_KEY=AAAA...  # <-- paste Server Key dari Firebase Console
```

### 5.2 Generate Valid PIN Hash untuk Admin

Sebelum menjalankan seed, buat bcrypt hash untuk PIN admin:

**Cara 1: Online** — buka https://bcrypt-generator.com/
- Input: `1234`
- Salin hasil hash (contoh: `$2a$12$LJ3m4ys3Lk0TSwHkJzWn7uJz...`)

**Cara 2: Menggunakan Go** — jalankan perintah ini di terminal:
```bash
cd ~/Projects/barberku/backend-barber
go run -mod=mod << 'EOF'
package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    hash, _ := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
    fmt.Println(string(hash))
}
EOF
```

### 5.3 Update seed.sql dengan PIN Hash yang Valid

Buka `seed.sql`, ganti hash PIN dengan hasil generate:

```sql
INSERT INTO users (id, name, email, phone, role, pin_hash, is_active, created_at, updated_at)
VALUES (
  '00000000-0000-0000-0000-000000000001',
  'Admin BarberKu',
  'admin@barberku.com',
  '+6281234567890',
  'admin',
  '$2a$12$LJ3m4ys3Lk0TSwHk...',  -- <-- ganti dengan hash bcrypt asli
  true,
  NOW(),
  NOW()
) ON CONFLICT (email) DO NOTHING;
```

---

## 6. Jalankan Backend

### 6.1 Start Docker Containers

```bash
cd ~/Projects/barberku/backend-barber
docker compose up -d --build
```

Ini akan:
- Membuat container **PostgreSQL 16** (port 5432)
- Membuat container **backend Go** (port 8080)
- Menjalankan migrasi database otomatis

### 6.2 Cek Status

```bash
# Cek container berjalan
docker compose ps

# Cek log backend
docker compose logs -f backend

# Cek health endpoint
curl http://localhost:8080/health
```

Output yang diharapkan: `{"status":"ok"}`

### 6.3 Seed Database

Salin seed.sql ke container PostgreSQL:
```bash
docker compose cp seed.sql db:/seed.sql
docker compose exec db psql -U barber -d barbershop -f /seed.sql
```

Verifikasi:
```bash
docker compose exec db psql -U barber -d barbershop -c "SELECT name, email, role FROM users;"
docker compose exec db psql -U barber -d barbershop -c "SELECT name FROM services;"
docker compose exec db psql -U barber -d barbershop -c "SELECT name FROM barbers;"
```

### 6.4 Cek API

```bash
# Login admin
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@barberku.com","pin":"1234"}'

# Response diharapkan:
# {"success":true,"message":"OK","data":{"token":"eyJ...","user":{"name":"Admin BarberKu",...}}}
```

---

## 7. Jalankan Flutter App

### 7.1 Pastikan URL Backend Sudah Benar

Cek file `barberku_app/lib/features/auth/presentation/providers/auth_provider.dart` (atau file konfigurasi endpoint) dan pastikan base URL mengarah ke `http://10.0.2.2:8080` (Android emulator) atau `http://localhost:8080` (web / iOS simulator).

### 7.2 Run

```bash
cd ~/Projects/barberku/barberku_app

# Untuk Android emulator
flutter run

# Untuk Chrome (web)
flutter run -d chrome

# Untuk iOS simulator (macOS)
flutter run -d ios

# Untuk device fisik
flutter run -d <device_id>
```

### 7.3 Login

- **Email**: `admin@barberku.com`
- **PIN**: `1234`

---

## 8. Testing End-to-End

### 8.1 Test Admin Flow

1. Buka aplikasi → tap **"Login Admin"**
2. Masukkan: `admin@barberku.com` / `1234`
3. Dashboard admin muncul dengan 3 tab: Waiting, Serving, Completed
4. Tap **+ Walk-In** untuk menambah antrian manual
5. Isi nama customer dan pilih layanan
6. Antrian muncul di tab **Waiting**
7. Tap **Panggil Berikutnya** (ikon phone) di AppBar
8. Antrian pindah ke tab **Serving**
9. Tap **Selesai** → pindah ke tab **Completed**

### 8.2 Test Notification

1. Pastikan FCM Server Key sudah benar di `.env`
2. Login sebagai admin
3. Tambahkan antrian
4. Panggil antrian tersebut
5. Cek log backend:
   ```bash
   docker compose logs -f backend | grep FCM
   ```
   Output: `"sending FCM notification"` atau `"failed to send FCM notification"`

### 8.3 Test API Langsung

```bash
# Dapatkan token JWT
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@barberku.com","pin":"1234"}' | \
  python3 -c "import sys,json;print(json.load(sys.stdin)['data']['token'])")

echo "Token: $TOKEN"

# Lihat daftar antrian
curl http://localhost:8080/api/v1/queue -H "Authorization: Bearer $TOKEN"

# Lihat services
curl http://localhost:8080/api/v1/services -H "Authorization: Bearer $TOKEN"

# Lihat barbers
curl http://localhost:8080/api/v1/barbers -H "Authorization: Bearer $TOKEN"

# Lihat statistik hari ini
curl "http://localhost:8080/api/v1/stats" -H "Authorization: Bearer $TOKEN"

# Lihat pengaturan toko
curl http://localhost:8080/api/v1/settings -H "Authorization: Bearer $TOKEN"
```

---

## 9. Troubleshooting

### 9.1 "Firebase not initialized"

```
Error: FirebaseException ([core/no-app] No Firebase App '[DEFAULT]' has been created)
```

**Solusi**:
1. Pastikan `google-services.json` ada di `android/app/`
2. Pastikan `firebase_options.dart` sudah di-generate
3. Update `main.dart` seperti di langkah 4.4

### 9.2 "Connection refused" di Flutter

```
DioException: SocketException: Connection refused
```

**Solusi**:
- Android emulator: gunakan `10.0.2.2` bukan `localhost`
- iOS simulator: gunakan `localhost`
- Device fisik: gunakan IP komputer (contoh: `192.168.1.10`)

### 9.3 FCM "401 Unauthorized"

```
FCM returned status 401: ...
```

**Solusi**:
1. Verifikasi Server Key di `.env` sudah benar
2. Regenerate Server Key jika perlu (Firebase Console → Cloud Messaging)
3. Restart backend: `docker compose restart backend`

### 9.4 PostgreSQL connection failed

```
failed to connect to database
```

**Solusi**:
1. Cek database container: `docker compose ps`
2. Cek log DB: `docker compose logs db`
3. Pastikan `.env` sesuai dengan docker-compose.yml:
   ```env
   DB_HOST=localhost   # untuk koneksi dari host
   DB_HOST=db          # untuk koneksi dari container backend
   ```
4. Jalankan ulang: `docker compose down && docker compose up -d`

### 9.5 Docker permission denied

```
permission denied while trying to connect to the Docker daemon socket
```

**Solusi**:
```bash
sudo usermod -aG docker $USER
# Logout dan login ulang, atau jalankan:
newgrp docker
```

---

## 10. Perintah Cepat

```bash
# ===== Backend =====

# Start semua container
docker compose up -d

# Lihat log backend
docker compose logs -f backend

# Lihat log database
docker compose logs -f db

# Masuk ke shell PostgreSQL
docker compose exec db psql -U barber -d barbershop

# Seed database
docker compose cp seed.sql db:/seed.sql
docker compose exec db psql -U barber -d barbershop -f /seed.sql

# Restart backend
docker compose restart backend

# Stop semua container
docker compose down

# Rebuild dan start ulang
docker compose up -d --build

# ===== Flutter =====

# Install dependencies
flutter pub get

# Run di emulator
flutter run

# Run di web
flutter run -d chrome

# Analisis kode
flutter analyze
```
