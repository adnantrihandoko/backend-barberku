# Setup Firebase Project & Konfigurasi FCM Server Key

Panduan ini menjelaskan langkah-langkah untuk membuat Firebase project, mengaktifkan Cloud Messaging, dan mengkonfigurasi FCM Server Key di aplikasi BarberKu.

---

## 1. Buat Firebase Project

1. Buka [Firebase Console](https://console.firebase.google.com/)
2. Klik **"Add project"** atau **"Create a project"**
3. Masukkan nama project: `barberku` (atau nama pilihan Anda)
4. (Opsional) Aktifkan Google Analytics jika diperlukan
5. Klik **"Create project"** dan tunggu hingga selesai

---

## 2. Tambahkan Aplikasi Flutter ke Firebase

### Android
1. Di Firebase Console, klik ikon **Android** (`Add app`)
2. Masukkan **Android package name** dari `barberku_app/android/app/build.gradle`:
   ```
   applicationId "com.barberku.app"
   ```
3. (Opsional) Masukkan SHA-1 debugging certificate:
   ```bash
   keytool -list -v -keystore ~/.android/debug.keystore -alias androiddebugkey -storepass android -keypass android
   ```
4. Download `google-services.json`
5. Letakkan file di: `barberku_app/android/app/google-services.json`

### iOS
1. Di Firebase Console, klik ikon **iOS** (`Add app`)
2. Masukkan **iOS bundle ID** dari `barberku_app/ios/Runner.xcodeproj`:
   ```
   com.barberku.app
   ```
3. Download `GoogleService-Info.plist`
4. Buka `barberku_app/ios/Runner.xcworkspace` di Xcode
5. Drag & drop `GoogleService-Info.plist` ke folder `Runner` (pastikan "Copy items if needed" dicentang)

### Web (Opsional)
1. Klik ikon **Web** (`Add app`)
2. Daftarkan aplikasi web
3. Salin konfigurasi Firebase SDK untuk `firebase_options.dart`

---

## 3. Aktifkan Firebase Cloud Messaging (FCM)

FCM sudah otomatis aktif saat project dibuat. Pastikan:

1. Di Firebase Console, buka **Project Settings** (ikon gear)
2. Pilih tab **Cloud Messaging**
3. Pastikan status FCM aktif

---

## 4. Dapatkan FCM Server Key

### Legacy Server Key (untuk HTTP API v1)
1. Di Firebase Console → **Project Settings** → **Cloud Messaging**
2. Di bagian **Cloud Messaging API**, klik **Manage API in Google Cloud Console**
3. Jika diminta, aktifkan **Firebase Cloud Messaging API**
4. Kembali ke Firebase Console
5. Salin **Server key** (bukan Sender ID)

> **Catatan**: Google merekomendasikan menggunakan **FCM HTTP v1 API** dengan service account key (JSON). Namun, implementasi saat ini menggunakan Legacy API untuk kesederhanaan.

### Service Account Key (untuk FCM HTTP v1 API - Recommended)
1. Buka [Google Cloud Console](https://console.cloud.google.com/)
2. Pilih project Firebase Anda
3. Navigasi ke **APIs & Services** → **Credentials**
4. Klik **Create Credentials** → **Service Account**
5. Isi nama: `barberku-fcm`
6. Berikan role: **Firebase Cloud Messaging Admin**
7. Klik **Done**
8. Di daftar service account, klik `barberku-fcm` → **Keys** → **Add Key** → **Create new key**
9. Pilih **JSON** dan download
10. Simpan file sebagai `firebase-credentials.json`

---

## 5. Konfigurasi Environment Variable

### Backend (.env)
Buat file `.env` di `backend-barber/`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=barber
DB_PASS=secret
DB_NAME=barbershop
SERVER_PORT=8080
WS_ORIGIN_ALLOWED=*
JWT_SECRET=barberku-secret-key-change-in-production
FCM_SERVER_KEY=your-legacy-server-key-here
```

Atau jika menggunakan service account JSON:
```env
FCM_CREDENTIAL_PATH=/path/to/firebase-credentials.json
```

### Flutter (Firebase Options)
Jika menggunakan `firebase_core` dengan konfigurasi manual, buat `lib/firebase_options.dart`:

```dart
import 'package:firebase_core/firebase_core.dart' show FirebaseOptions;
import 'package:flutter/foundation.dart';

class DefaultFirebaseOptions {
  static FirebaseOptions get currentPlatform {
    if (kIsWeb) {
      return web;
    }
    switch (defaultTargetPlatform) {
      case TargetPlatform.android:
        return android;
      case TargetPlatform.iOS:
        return ios;
      case TargetPlatform.macOS:
        return macos;
      case TargetPlatform.windows:
        return windows;
      case TargetPlatform.linux:
        throw UnsupportedError('Linux not supported');
      default:
        throw UnsupportedError('Unknown platform');
    }
  }

  static const FirebaseOptions android = FirebaseOptions(
    apiKey: 'YOUR_ANDROID_API_KEY',
    appId: 'YOUR_ANDROID_APP_ID',
    messagingSenderId: 'YOUR_SENDER_ID',
    projectId: 'YOUR_PROJECT_ID',
    storageBucket: 'YOUR_PROJECT_ID.appspot.com',
  );

  static const FirebaseOptions ios = FirebaseOptions(
    apiKey: 'YOUR_IOS_API_KEY',
    appId: 'YOUR_IOS_APP_ID',
    messagingSenderId: 'YOUR_SENDER_ID',
    projectId: 'YOUR_PROJECT_ID',
    storageBucket: 'YOUR_PROJECT_ID.appspot.com',
    iosBundleId: 'com.barberku.app',
  );

  static const FirebaseOptions web = FirebaseOptions(
    apiKey: 'YOUR_WEB_API_KEY',
    appId: 'YOUR_WEB_APP_ID',
    messagingSenderId: 'YOUR_SENDER_ID',
    projectId: 'YOUR_PROJECT_ID',
    storageBucket: 'YOUR_PROJECT_ID.appspot.com',
  );

  static const FirebaseOptions macos = FirebaseOptions(
    apiKey: 'YOUR_MACOS_API_KEY',
    appId: 'YOUR_MACOS_APP_ID',
    messagingSenderId: 'YOUR_SENDER_ID',
    projectId: 'YOUR_PROJECT_ID',
    storageBucket: 'YOUR_PROJECT_ID.appspot.com',
    iosBundleId: 'com.barberku.app',
  );

  static const FirebaseOptions windows = FirebaseOptions(
    apiKey: 'YOUR_WINDOWS_API_KEY',
    appId: 'YOUR_WINDOWS_APP_ID',
    messagingSenderId: 'YOUR_SENDER_ID',
    projectId: 'YOUR_PROJECT_ID',
    storageBucket: 'YOUR_PROJECT_ID.appspot.com',
  );
}
```

> **Alternatif**: Gunakan `flutterfire configure` CLI untuk generate otomatis:
> ```bash
> dart pub global activate flutterfire_cli
> flutterfire configure --project=barberku
> ```

---

## 6. Verifikasi Setup

### Test FCM Token Registration
1. Jalankan aplikasi Flutter
2. Cek log untuk FCM token:
   ```
   FCM Token: <token_string>
   ```
3. Token akan otomatis dikirim ke backend saat login

### Test Push Notification
1. Login sebagai admin
2. Buat antrian baru (customer)
3. Admin memanggil antrian (`CallQueue`)
4. Backend akan mengirim notifikasi FCM ke customer
5. Customer menerima notifikasi:
   - **Foreground**: Local notification muncul
   - **Background**: System notification muncul

### Test via Firebase Console
1. Buka Firebase Console → **Cloud Messaging**
2. Klik **Send your first message**
3. Isi judul dan body
4. Targetkan ke aplikasi `barberku`
5. Klik **Review** → **Publish**

---

## 7. Troubleshooting

| Masalah | Solusi |
|---------|--------|
| `FCM permission denied` | Pastikan permission diminta di `main.dart` |
| `Token not registered` | Cek koneksi internet & konfigurasi Firebase |
| `401 Unauthorized dari FCM` | Verifikasi Server Key di `.env` |
| `Notification tidak muncul di background` | Pastikan `google-services.json` / `GoogleService-Info.plist` benar |
| `Firebase initialization failed` | Jalankan `flutterfire configure` ulang |

---

## 8. Security Best Practices

1. **Jangan commit** `google-services.json`, `GoogleService-Info.plist`, atau `firebase-credentials.json` ke repository
2. Tambahkan ke `.gitignore`:
   ```
   **/google-services.json
   **/GoogleService-Info.plist
   **/firebase-credentials.json
   ```
3. Gunakan environment variable untuk FCM Server Key di production
4. Rotate Server Key secara berkala di Firebase Console
5. Batasi akses service account ke role minimum yang diperlukan

---

## 9. Production Deployment

### Android
1. Generate SHA-1 untuk release keystore:
   ```bash
   keytool -list -v -keystore /path/to/release.keystore -alias your-alias
   ```
2. Tambahkan SHA-1 ke Firebase Console → Project Settings → Your Apps
3. Download ulang `google-services.json` dengan SHA-1 production

### iOS
1. Pastikan **Push Notifications** capability aktif di Xcode
2. Upload APNs certificate ke Firebase Console (Project Settings → Cloud Messaging → iOS app)

### Backend
1. Set environment variable `FCM_SERVER_KEY` di server production
2. Gunakan HTTPS untuk semua endpoint
3. Monitor usage di Firebase Console → Cloud Messaging → Reports

---

Dibutuhkan: 15-30 menit untuk setup lengkap.
