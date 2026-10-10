# ISSUES

Review seluruh codebase pada 2026-10-10, diperbarui setelah dua putaran perbaikan di hari yang sama.

Semua 35 temuan sudah ditangani. Satu di antaranya, M2, ditutup sebagian atas keputusan pemilik project dan sisa risikonya dicatat di bawah.

Semua perbaikan lolos build, `go vet`, unit test, `govulncheck`, dan integration test ke PostgreSQL 18.

## Risiko yang tersisa

### M2. Token dikirim di body JSON

- Level: Medium, diterima atas keputusan pemilik project.
- `access_token` dan `refresh_token` tetap ada di body response untuk semua klien, jadi script di halaman klien (XSS) bisa membacanya.
- Yang sudah dikurangi: umur refresh token turun dari 7 hari menjadi 24 jam, dan plugin jwt merotasi refresh token setiap dipakai.
- Yang belum tertutup: selama 24 jam itu refresh token curian tetap bisa dipakai. Menutupnya berarti tidak mengirim token ke klien browser.

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
| L2 | Access token tidak mengecek masa berlaku session | Middleware `Authorization` menolak bearer token yang session-nya sudah kedaluwarsa. Biayanya satu query tambahan per request bearer. |
| L3 | `/health` membuka info internal | Hanya mengembalikan status. Flag `machine-id` dihapus. |
| L4 | Password super user tidak ikut berubah | `SUPER_USER` menjadi sumber kebenaran: password yang berbeda di-reset saat boot dengan seed. |
| L5 | Nama category tanpa batas | Maksimal 100 karakter. |
| L6 | Sentry tidak diinisialisasi | `sentry.Init` dipanggil saat `SENTRY_DSN` terisi; middleware hanya dipasang saat itu. |
| L7 | Tidak ada proteksi CSRF selain `SameSite=Lax` | Middleware `CrossOrigin` menolak request tulis dari browser di origin lain, termasuk subdomain saudara, kecuali origin itu ada di `CORS_ALLOWED_ORIGINS`. Tanpa token, jadi klien yang ada tidak perlu berubah; berlaku juga untuk route Authula. |
| L8 | Tidak ada pemindaian dependency | Workflow release menjalankan vet, test, dan `govulncheck` sebelum membuat tag; tidak ada workflow yang jalan di setiap push. Pemindaian pertama menemukan 19 kerentanan; semuanya tertutup dengan Go 1.27.2, `x/net` 0.60.0, `amqp091-go` 1.13.0, dan `go-archive` 0.3.0. |
| C1 | Fitur auth belum pernah dijalankan ke database | `cmd/api/integration_test.go` menjalankan boot dan alur auth ke PostgreSQL 18, dan lolos tanpa perlu perbaikan kode. Workflow release menjalankannya dengan service Postgres. |
| C2 | Mirror `public.users` bisa berbeda | `Reconcile` menyalin ulang semua user Authula saat boot. Baris mirror milik user yang sudah dihapus sengaja tidak dihapus. |
| C3 | Petunjuk perbaikan migration salah | Petunjuk dan batas versi diperbaiki; `cmd/migrate force <versi>` dan target `db/migrations/force` ditambahkan. |
| C4 | `make release` tanpa script | Target `release` dihapus; release kini lewat workflow release-it. |
| C5 | `/health` tidak mengecek database | Handler mem-ping database dan menjawab 503 saat gagal. |
| C6 | Token tidak selalu ada di response sign-in | Sign-in ulang saat session aktif sekarang tetap mengembalikan token. `POST /v1/auth/token` memberi token baru untuk session yang sedang dipakai, termasuk session dari login Google dengan redirect. |
| C7 | Sisa nama project lain | Default app name `gofi`, kunci lock `gofi:seed`, variabel handler category dirapikan. |
| C8 | README kosong | README berisi cara menjalankan, env, database, super user, autentikasi, Docker, dan release. |

Yang dicakup integration test (C1):

- Migration naik dan turun, termasuk schema `authula` dan kolom `role`.
- Seeder dijalankan dua kali, lalu reset password super user.
- Mirror `public.users`, role, dan `Reconcile` setelah satu baris dihapus.
- Sign-in dengan token di body, sign-in ulang, `POST /token` lewat cookie dan bearer, refresh, `/me`.
- Route tulis category: 401 tanpa login, 403 untuk user biasa, 201 untuk super user.
- Request lintas situs ditolak, dan token ditolak setelah sign-out.

Yang tidak dicakup: login Google dan magic link (butuh layanan luar), serta pengiriman email.

Kode yang dibuang (11 butir audit ponytail, semuanya selesai):

- 8 rule validator tanpa pemakai; `MaxLen` ditambahkan karena dipakai L1 dan L5.
- `internal/repositories/include.go`, `internal/lib/uuid.go`, `internal/services/`, `internal/dtos/health.go`, `internal/seeders/data/category.json`.
- Pasangan method `exec*` di repository, beserta interface `Executor`.
- `config.Load()`, `loadJSON`, `optionalUUID`, `includeQuery`, `withTx`, `execAffected`, `translateUnique`, `TotalMeta`.
- `constant.CategoryTypes`, `app.Version`, `KindUnprocessable`.
- `--seed=prod|dev` menjadi flag boolean `--seed`.
- Dependency `go-sqlfmt`, diganti `slog.Debug`.

## Perubahan yang memengaruhi pemakaian

- `ENV` yang tidak diset berarti `production`. Set `ENV=development` di `.env` lokal.
- `DEBUG=true` bersama `ENV=production` membuat API menolak start.
- Route tulis category butuh role `super_user`; user biasa mendapat 403.
- Refresh token berlaku 24 jam, bukan 7 hari. Setelah itu klien meminta pasangan baru di `POST /v1/auth/token` selama session-nya masih hidup.
- Request tulis dari browser di origin yang tidak ada di `CORS_ALLOWED_ORIGINS` ditolak dengan 403.
- `cmd/migrate --seed=dev` menjadi `--seed`, dan `--app-pwd`/`machine-id` sudah tidak ada.
- Go minimal 1.27.2.
