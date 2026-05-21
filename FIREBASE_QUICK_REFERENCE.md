# Quick Reference: Firebase & FCM Setup

## 1. Firebase Console
- URL: https://console.firebase.google.com/
- Project ID: `barberku` (atau ID Anda)

## 2. Download Config Files
| Platform | File | Location |
|----------|------|----------|
| Android | `google-services.json` | `barberku_app/android/app/` |
| iOS | `GoogleService-Info.plist` | `barberku_app/ios/Runner/` |
| Web | `firebase_options.dart` | `barberku_app/lib/` |

## 3. FCM Server Key
- Firebase Console → Project Settings → Cloud Messaging
- Copy **Server key** (bukan Sender ID)
- Paste ke `backend-barber/.env`:
  ```env
  FCM_SERVER_KEY=your-server-key-here
  ```

## 4. Flutter Setup
```bash
# Install FlutterFire CLI
dart pub global activate flutterfire_cli

# Generate firebase_options.dart
cd barberku_app
flutterfire configure --project=barberku
```

## 5. Backend Setup
```bash
cd backend-barber
cp .env.example .env
# Edit .env dan isi FCM_SERVER_KEY
docker-compose up -d
```

## 6. Test
1. Jalankan Flutter app
2. Login sebagai customer
3. Cek log: `FCM Token: <token>`
4. Admin panggil antrian → notifikasi muncul

## 7. .gitignore
```
**/google-services.json
**/GoogleService-Info.plist
**/firebase-credentials.json
**/.env
```

---
**Full guide**: `FIREBASE_SETUP.md`
