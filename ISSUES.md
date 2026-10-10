# ISSUES

Review seluruh codebase pada 2026-10-10, diperbarui setelah perbaikan di hari yang sama.

Dari 35 temuan, 30 sudah diperbaiki. Lima sisanya ada di bawah: satu menunggu keputusan, tiga diterima sebagai risiko, satu masih terbuka.

Semua perbaikan lolos build, `go vet`, unit test, dan `govulncheck`. Tidak ada yang diuji ke server atau database hidup (lihat C1).

## Masih terbuka

### C1. Fitur auth belum pernah dijalankan ke database

- Status: **terbuka**. Tidak ada Postgres maupun Docker di mesin tempat perbaikan dikerjakan.
- Belum teruji: schema `authula`, seeder super user, mirror `public.users` dan reconcile-nya, kolom `role`, token JWT, bearer token.
- Paling berisiko: SQL `Reconcile` di `internal/repositories/user.go` dan migration `000002`/`000003`, karena hanya dibaca, belum pernah dieksekusi.
- Langkah berikutnya: jalankan sekali dengan `MIGRATE_ON_BOOT=true` dan `SEED_ON_BOOT=true`, lalu coba sign-in dan satu route tulis category.

## Menunggu keputusan

### M2. Refresh token 7 hari dikirim di body JSON

- Level: Medium.
- Lokasi: `cmd/api/authula.go` (`addTokensToResponse`).
- Masalah: `access_token` dan `refresh_token` ada di body response sign-in, jadi bisa dibaca JavaScript. Script hasil XSS bisa mengambil refresh token dan memakainya selama 7 hari, walaupun cookie session `HttpOnly`.
- Kenapa belum diubah: field ini memang diminta untuk ketiga jenis login. Menutup risikonya berarti mengubah fitur itu.
- Pilihan: kirim token hanya untuk klien non-browser, atau perpendek umur refresh token.

## Diterima sebagai risiko

### L2. Access token tidak mengecek masa berlaku session

- Level: Low.
- Validasi token memastikan session-nya masih ada, bukan belum kedaluwarsa.
- Diterima karena jendelanya paling lama 15 menit (umur access token) setelah session kedaluwarsa, sedangkan menutupnya menambah satu query di setiap request bearer.
- Sign-out tetap langsung mencabut token, karena session-nya dihapus.

### L7. Tidak ada proteksi CSRF selain `SameSite=Lax`

- Level: Low.
- Cukup untuk serangan lintas situs biasa; tidak cukup kalau ada subdomain yang tidak dipercaya di domain yang sama.
- Diterima karena plugin `csrf` mewajibkan semua klien cookie mengirim header token, yang merusak klien yang sudah ada.
- Pasang plugin `csrf` Authula begitu ada subdomain semacam itu.

### C6. Token tidak selalu ada di response sign-in

- Sign-in ulang saat session masih aktif tidak membuat token baru; klien memakai `POST /v1/auth/token/refresh`.
- Login Google yang dimulai dengan `redirect_to` dijawab dengan redirect, yang membawa cookie tetapi tidak token.
- Diterima sebagai perilaku Authula dan sudah didokumentasikan di `README.md`.

## Sudah diperbaiki

| ID | Temuan | Perbaikan |
|----|--------|-----------|
| H1 | Semua user terdaftar bisa mengubah category | Kolom `users.role` (migration `000003`), seeder memberi super user role `super_user`, dan middleware `SuperUser` menjaga semua route tulis category. |
| M1 | `ENV` default `development` | Default menjadi `production`; nilai divalidasi (`development`, `staging`, `production`). |
| M3 | Secret lewat argumen CLI | `make` meng-`export` `.env`; `make run` dan target migration tidak lagi mengirim flag. `cmd/migrate` membaca `DB_DSN` dari environment. |
| M4 | Server tanpa timeout | `ReadTimeout` 15 detik, `WriteTimeout` 30 detik, `IdleTimeout` 60 detik. |
| M5 | `DEBUG=true` mencetak nilai query | `DEBUG=true` ditolak saat `ENV=production`; contoh di `.env.example` menjadi `false`. Log query repository hanya berisi teks SQL. |
| M6 | `AUTHULA_DATABASE_URL` melewati schema | API menolak start kalau env itu terisi. |
| M7 | Migration down tanpa pengaman | `down` dan `refresh` menolak database di luar localhost kecuali diberi `--allow-remote`. |
| L1 | Rate limit hanya di auth | Limiter 120 request/menit per IP di `/v1`; `search` dibatasi 100 karakter. Limiter masih per instance (ditandai komentar `ponytail:`). |
| L3 | `/health` membuka info internal | Hanya mengembalikan status. Flag `machine-id` dihapus. |
| L4 | Password super user tidak ikut berubah | `SUPER_USER` menjadi sumber kebenaran: password yang berbeda di-reset saat boot dengan seed. |
| L5 | Nama category tanpa batas | Maksimal 100 karakter. |
| L6 | Sentry tidak diinisialisasi | `sentry.Init` dipanggil saat `SENTRY_DSN` terisi; middleware hanya dipasang saat itu. |
| L8 | Tidak ada pemindaian dependency | CI di `.github/workflows/ci.yml` menjalankan build, vet, test, dan `govulncheck`. Pemindaian pertama menemukan 19 kerentanan; semuanya tertutup dengan Go 1.27.2, `x/net` 0.60.0, `amqp091-go` 1.13.0, dan `go-archive` 0.3.0. |
| C2 | Mirror `public.users` bisa berbeda | `Reconcile` menyalin ulang semua user Authula saat boot. Baris mirror milik user yang sudah dihapus sengaja tidak dihapus. |
| C3 | Petunjuk perbaikan migration salah | Petunjuk dan batas versi diperbaiki; `cmd/migrate force <versi>` dan target `db/migrations/force` ditambahkan. |
| C4 | `make release` tanpa script | Target `release` dihapus. |
| C5 | `/health` tidak mengecek database | Handler mem-ping database dan menjawab 503 saat gagal. |
| C7 | Sisa nama project lain | Default app name `gofi`, kunci lock `gofi:seed`, variabel handler category dirapikan. |
| C8 | README kosong | README berisi cara menjalankan, env, database, super user, dan autentikasi. |

Kode yang dibuang (11 butir audit ponytail, semuanya selesai):

- 8 rule validator tanpa pemakai; `MaxLen` ditambahkan karena dipakai L1 dan L5.
- `internal/repositories/include.go`, `internal/lib/uuid.go`, `internal/services/`, `internal/dtos/health.go`, `internal/seeders/data/category.json`.
- Pasangan method `exec*` di repository, beserta interface `Executor`.
- `config.Load()`, `loadJSON`, `optionalUUID`, `includeQuery`, `withTx`, `execAffected`, `translateUnique`, `TotalMeta`.
- `constant.CategoryTypes`, `app.Version`, `KindUnprocessable`.
- `--seed=prod|dev` menjadi flag boolean `--seed`.
- Dependency `go-sqlfmt`, diganti `slog.Debug`.

Yang sengaja dipertahankan: `ErrForbidden` dan `KindForbidden` (sekarang dipakai H1) serta `forceVersion` (sekarang dipakai C3).

## Perubahan yang memengaruhi pemakaian

- `ENV` yang tidak diset sekarang berarti `production`. Set `ENV=development` di `.env` lokal.
- `DEBUG=true` bersama `ENV=production` membuat API menolak start.
- Route tulis category sekarang butuh role `super_user`; user biasa mendapat 403.
- `cmd/migrate --seed=dev` menjadi `--seed`, dan `--app-pwd`/`machine-id` sudah tidak ada.
- Go minimal menjadi 1.27.2.
