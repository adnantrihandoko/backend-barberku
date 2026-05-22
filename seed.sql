-- Seed data untuk BarberKu
-- Jalankan: psql -U barber -d barbershop -f seed.sql

-- Admin default (PIN: 1234)
-- PIN hash bcrypt: $2a$10$e0MYzXyIkJVowR0sXqE0ZuP9qZQZQZQZQZQZQZQZQZQZQZQZQZQZQ
-- Gunakan bcrypt online atau library untuk generate hash yang valid
INSERT INTO users (id, name, email, phone, role, pin_hash, is_active, created_at, updated_at)
VALUES (
  '00000000-0000-0000-0000-000000000001',
  'Admin BarberKu',
  'admin@barberku.com',
  '+6281234567890',
  'admin',
  '$2b$12$yerCgCD6vYVSdGGEcUgUGOTUA1LhRTOAHfjeEAEjEyk65710s09du', -- bcrypt hash dari "1234"
  true,
  NOW(),
  NOW()
) ON CONFLICT (email) DO NOTHING;

-- Layanan default
INSERT INTO services (id, name, description, price, duration, is_active, created_at, updated_at)
VALUES
  (gen_random_uuid(), 'Potong Rambut', 'Potong rambut standar', 35000, 30, true, NOW(), NOW()),
  (gen_random_uuid(), 'Cukur Jenggot', 'Cukur jenggot rapi', 25000, 20, true, NOW(), NOW()),
  (gen_random_uuid(), 'Potong + Cukur', 'Paket lengkap potong rambut dan cukur jenggot', 55000, 45, true, NOW(), NOW()),
  (gen_random_uuid(), 'Keramas + Potong', 'Keramas dulu baru potong', 45000, 40, true, NOW(), NOW()),
  (gen_random_uuid(), 'Hair Coloring', 'Pewarnaan rambut', 150000, 90, true, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Barber default
INSERT INTO barbers (id, name, specialty, is_active, created_at, updated_at)
VALUES
  (gen_random_uuid(), 'Ahmad', 'Senior Barber', true, NOW(), NOW()),
  (gen_random_uuid(), 'Budi', 'Junior Barber', true, NOW(), NOW()),
  (gen_random_uuid(), 'Candra', 'Stylist', true, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Store settings default
INSERT INTO store_settings (open_hour, close_hour, max_queue_size, updated_at)
VALUES (9, 21, 50, NOW())
ON CONFLICT DO NOTHING;
