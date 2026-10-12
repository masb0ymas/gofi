# ISSUES

Audit seluruh codebase pada 2026-10-12, dengan Authula v1.46.0. Versi ISSUES.md sebelumnya (35 temuan dari 2026-10-10, semuanya sudah ditutup) ada di riwayat git.

**Apa yang dikerjakan repo ini:** GoFi adalah starter API Go di atas Fiber dan PostgreSQL. Login (email & password, Google, magic link) ditangani Authula yang ditanam sebagai library di `/v1/auth`; aplikasinya sendiri baru punya CRUD category, dengan baca untuk publik dan tulis hanya untuk role `super_user`.

**Beban yang diasumsikan:** API publik dengan sign-up terbuka, 1 sampai 3 replica di belakang reverse proxy (nginx atau load balancer), ratusan sampai ribuan user. Tidak ada uang atau data pribadi selain akun.

**Yang sudah dicek:** `go vet` dan `go test ./...` lolos. `govulncheck` tidak menemukan kerentanan yang terpanggil oleh kode. Saat audit, integration test belum dijalankan karena tidak ada PostgreSQL; setelah perbaikan ia dijalankan dan lolos (lihat di bawah).

## Ringkasan per level keamanan

| Level  | No  | Temuan                                                                                                     | Status                             |
| ------ | --- | ---------------------------------------------------------------------------------------------------------- | ---------------------------------- |
| High   | 1   | Email yang belum diverifikasi langsung bisa login, sehingga akun bisa dibajak sebelum pemiliknya mendaftar | Diperbaiki, sisa risiko di Authula |
| High   | 2   | Rate limit bisa dilewati dengan header `X-Forwarded-For` palsu                                             | Diperbaiki                         |
| High   | 3   | Reset password dan rotasi `SUPER_USER` tidak mengakhiri session yang sudah ada                             | Diperbaiki                         |
| Medium | 4   | `POST /v1/auth/token` mengubah access token curian (15 menit) menjadi akses selama 7 hari                  | Diperbaiki                         |
| Medium | 5   | Workflow release tidak pernah menjalankan integration test, padahal README bilang menjalankannya           | Diperbaiki                         |
| Low    | 6   | Action pihak ketiga dengan izin tulis dipasang lewat tag, bukan commit SHA                                 | Diperbaiki                         |
| Low    | 7   | Satu replica membuka sampai 50 koneksi database                                                            | Diperbaiki                         |
| Low    | 8   | Nama category berisi spasi saja diterima                                                                   | Terbuka                            |
| Info   | 9   | Sisa demo: route `GET /`, folder `public/`, folder `templates/`                                            | Terbuka                            |
| Info   | 10  | Dua pembungkus dengan satu pemakai                                                                         | Terbuka                            |

## Perbaikan 2026-10-12

| No  | Perbaikan                                                                                                                                                                                                                                                                                                                                                        |
| --- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | `AutoSignIn: false`, jadi sign-up tidak lagi memberi session. Hook `BeforeSignIn` (`refuseUnverified`) menolak sign-in password untuk email yang belum terverifikasi, dengan 401 dan pesan yang sama seperti password salah. Seeder memverifikasi akun super user yang didaftarkan orang lain lebih dulu, lalu me-reset password-nya dan mengakhiri session-nya. |
| 2   | `EnableIPValidation: true` di `fiber.Config`. IP klien diambil dari entri `X-Forwarded-For` paling kanan yang bukan proxy tepercaya.                                                                                                                                                                                                                             |
| 3   | Hook `AfterChangePassword` menghapus semua session user setelah reset password. Seeder super user menghapus semua session setiap kali ia mengganti password atau memverifikasi akun.                                                                                                                                                                             |
| 4   | `POST /v1/auth/token` hanya menerima cookie session, tidak lagi bearer.                                                                                                                                                                                                                                                                                          |
| 5   | Job `prepare` punya service `postgres:18`, dan step Test mengisi `TEST_DB_DSN`.                                                                                                                                                                                                                                                                                  |
| 6   | `peter-evans/create-pull-request` dipasang ke commit `22a9089` (v7.0.11), dan govulncheck ke `v1.8.0`. v1.9.0 dirilis 2026-10-11, terlalu baru untuk dipakai.                                                                                                                                                                                                    |
| 7   | `DB_MAX_OPEN_CONNS` dan `DB_MAX_IDLE_CONNS` menjadi anggaran total per replica: Authula mendapat separuh, pool API sisanya (`authPoolSize`, `apiPoolSize`).                                                                                                                                                                                                      |

Integration test (`cmd/api/integration_test.go`) sekarang juga mencakup:

- sign-up tanpa session, dan sign-in yang ditolak selama email belum terverifikasi;
- `POST /token` dengan bearer saja dijawab 401;
- rotasi `SUPER_USER` mematikan cookie dan access token yang lama;
- 31 sign-in dengan prefix `X-Forwarded-For` yang berganti-ganti tetap kena 429.

Setiap assertion baru dicek gagal pada perilaku lama dan lolos setelah perbaikan.

**Sisa risiko:**

- **No. 1:** di Authula, login Google atau magic link masih bisa menempel ke akun yang dibuat orang lain dengan email yang sama. Sekarang password penyerang tidak bisa dipakai selama email itu belum terverifikasi. Tapi kalau pemiliknya memverifikasi lewat magic link, password itu ikut berlaku lagi. Ini perlu diperbaiki di upstream.
- **No. 3:** hook reset password tidak dicakup integration test, karena token reset hanya dikirim lewat email.
- **No. 4:** klien tanpa cookie (mobile) harus login ulang setelah refresh token 24 jam habis.

---

## High

### 1. Email yang belum diverifikasi langsung bisa login (`cmd/api/authula.go:L345-353`)

- **Apa ini:** Plugin email-password memakai `RequireEmailVerification: true` dan `AutoSignIn: true`.
- **Masalah:** Di Authula v1.46.0, `RequireEmailVerification` hanya mengirim email; ia tidak menolak login. Sign-up dengan email siapa pun langsung memberi cookie session dan JWT, dan sign-in juga tidak mengecek `EmailVerified`. Akibatnya ada dua serangan:
  - **Pembajakan sebelum daftar:** penyerang daftar dengan `korban@gmail.com` dan password miliknya. Saat korban kemudian login lewat Google atau magic link, Authula mencari user berdasarkan email dan menempelkan login itu ke akun milik penyerang (`plugins/oauth2/usecases/callback_usecase.go:L69`, `plugins/magic-link/usecases/sign_in_usecase.go:L39`). Password penyerang tetap berlaku, jadi ia bisa masuk ke akun korban.
  - **Merebut super user:** kalau API sempat jalan dengan `SEED_ON_BOOT=false` (nilai di `.env.example`), penyerang yang menebak email super user bisa mendaftar lebih dulu. Seeder kemudian menemukan user itu, me-reset password-nya, lalu memberinya role `super_user`, sementara session penyerang dari sign-up masih hidup (lihat no. 3).
- **Perbaikan:** Di config plugin, set `AutoSignIn: false` dan tolak user yang belum terverifikasi lewat hook:
  ```go
  ServiceHooks: &emailpasswordplugintypes.EmailPasswordServiceHooksConfig{
      SignIn: &emailpasswordplugintypes.SignInServiceHooksConfig{
          BeforeSignIn: func(_ context.Context, u *authulamodels.User) error {
              if !u.EmailVerified {
                  return errors.New("invalid credentials") // sama dengan password salah, supaya tidak membocorkan email
              }
              return nil
          },
      },
  },
  ```
  Konsekuensinya, user yang link verifikasinya kedaluwarsa harus login lewat magic link, yang sekaligus memverifikasi email. Sisa risikonya ada di Authula: kalau akun korban terverifikasi lewat magic link, password penyerang ikut hidup lagi. Laporkan ke upstream bahwa menempelkan login Google atau magic link ke akun yang belum terverifikasi harus membuang akun password-nya.
- **Kalau dibiarkan:** Siapa pun bisa menyiapkan akun atas nama orang lain dan masuk ke sana setelah pemiliknya memakai Google atau magic link.

### 2. Rate limit bisa dilewati dengan `X-Forwarded-For` palsu (`cmd/api/server.go:L49-59`, `internal/middlewares/rate_limit.go:L13-24`)

- **Apa ini:** Saat `TRUSTED_PROXIES` diisi, Fiber membaca IP klien dari header `X-Forwarded-For` (header berisi daftar IP yang ditambahkan proxy). Limiter memakai IP itu sebagai kunci, 30 request per menit di `/v1/auth` dan 120 di `/v1`.
- **Masalah:** `EnableIPValidation` tidak diset, jadi `c.IP()` mengembalikan seluruh isi header apa adanya (`fiber/v3 req.go:L729-731`). nginx dengan `$proxy_add_x_forwarded_for` menambahkan IP asli di belakang nilai kiriman klien. Penyerang tinggal mengirim `X-Forwarded-For: <angka acak>` di setiap request, sehingga setiap request mendapat kunci baru dan tidak pernah kena limit.
- **Perbaikan:** Tambahkan satu baris di `fiber.Config`: `EnableIPValidation: true`. Dengan itu Fiber membaca header dari kanan, melewati IP proxy tepercaya, dan mengambil IP pertama yang tidak tepercaya, yaitu IP klien yang sebenarnya.
- **Kalau dibiarkan:** Tebakan password di `/sign-in` tidak dibatasi, dan siapa pun bisa membanjiri kotak masuk orang dengan email magic link atau reset password lewat Resend, yang juga membebani biaya.

### 3. Reset password tidak mengakhiri session yang sudah ada (`cmd/api/authula.go:L345-353`, `internal/seeders/super_user.go:L72-88`)

- **Apa ini:** Ada dua jalur ganti password: reset lewat email (`/email-password/change-password`) dan seeder super user yang me-reset password saat `SUPER_USER` diganti. README menyebut jalur kedua untuk "merotasi password yang bocor".
- **Masalah:** Tidak satu pun jalur menghapus session. Authula juga tidak melakukannya (`change_password_usecase.go` hanya memperbarui hash). Session berlaku 7 hari dan diperpanjang setiap hari dipakai, jadi penyerang yang sudah login dengan password bocor tetap masuk dan tetap bisa mengubah category setelah admin mengganti `SUPER_USER` dan restart.
- **Perbaikan:** Hapus semua session user setiap kali password-nya diganti.
  - Di plugin email-password: `ServiceHooks.PasswordChange.AfterChangePassword` memanggil `auth.CoreServices().SessionService.DeleteAllByUserID(ctx, user.ID)`.
  - Di `SuperUserSeeder.Seed`: panggil `s.Auth.SessionService.DeleteAllByUserID(ctx, user.ID)` di cabang yang menulis hash baru.
    Access token JWT ikut mati karena middleware `Authorization` dan plugin jwt mengecek bahwa session-nya masih ada.
- **Kalau dibiarkan:** Langkah pemulihan setelah password bocor tidak mengusir penyerang.

## Medium

### 4. `POST /v1/auth/token` memperpanjang access token curian sampai 7 hari (`cmd/api/authula.go:L56-59`, `L247-262`)

- **Apa ini:** Route ini membuat pasangan access token dan refresh token baru untuk session pemanggil. Mapping-nya menerima cookie maupun `Authorization: Bearer`.
- **Masalah:** Dengan bearer, access token 15 menit bisa ditukar terus-menerus dengan pasangan baru selama session-nya hidup, yaitu 7 hari. Bearer tidak memperpanjang session, tapi juga tidak dibatasi 24 jam. Ini lebih panjang dari sisa risiko yang diterima di versi sebelumnya (M2: refresh token curian berlaku 24 jam).
- **Perbaikan:** Hapus `bearerAuthOptional` dari mapping `POST:/token`, sehingga route itu hanya menerima cookie session. Klien tanpa cookie (mobile) memperbarui token lewat `/token/refresh` dan login ulang setelah 24 jam. Kalau perilaku sekarang memang disengaja, catat sisa risiko 7 hari ini di README.
- **Kalau dibiarkan:** Satu access token yang bocor lewat XSS atau log memberi akses selama seminggu, bukan 15 menit.

### 5. Gerbang release tidak menjalankan integration test (`.github/workflows/release-prepare.yml:L53-54`)

- **Apa ini:** Satu-satunya test yang mencakup login, token, role `super_user`, CSRF, dan idempotency ada di `cmd/api/integration_test.go`. Test ini dilewati (skip) kalau `TEST_DB_DSN` kosong.
- **Masalah:** Workflow menjalankan `go test ./...` tanpa service Postgres dan tanpa `TEST_DB_DSN`, jadi test itu selalu dilewati. README ("The release workflow runs all of it") dan ISSUES.md lama (C1) menyatakan sebaliknya. Perubahan yang merusak cek role tetap bisa lolos ke release.
- **Perbaikan:** Di job `prepare`, tambahkan `services: postgres: image: postgres:18` dengan `POSTGRES_PASSWORD` dan `POSTGRES_DB`, lalu set `TEST_DB_DSN` di step Test. Itu sekitar 12 baris YAML, tanpa perubahan kode.
- **Kalau dibiarkan:** Regresi di bagian auth baru ketahuan di production.

## Low

### 6. Action pihak ketiga dengan izin tulis tidak dipasang ke commit (`.github/workflows/release-prepare.yml:L57`, `L84`)

- **Apa ini:** `peter-evans/create-pull-request@v7` berjalan dengan `contents: write` dan `pull-requests: write`. `govulncheck@latest` diunduh setiap kali jalan.
- **Masalah:** Tag seperti `v7` bisa dipindah oleh pemilik action, atau oleh penyerang yang mengambil alih repo action itu. Kode baru langsung berjalan dengan token tulis repo ini, yang bisa membuat branch, tag, atau memicu build image.
- **Perbaikan:** Pasang action pihak ketiga ke commit SHA (`uses: peter-evans/create-pull-request@<sha> # v7`) dan kunci versi govulncheck (`govulncheck@v1.x.y`). Action resmi `actions/*` dan `docker/*` boleh tetap memakai tag.
- **Kalau dibiarkan:** Satu action yang dibobol bisa menyuntikkan kode ke release.

### 7. Satu replica membuka sampai 50 koneksi database (`cmd/api/authula.go:L286-291`, `cmd/api/flag.go:L90-91`)

- **Apa ini:** Aplikasi dan Authula masing-masing punya pool koneksi dengan batas yang sama, `DB_MAX_OPEN_CONNS=25`.
- **Masalah:** Satu replica bisa memakai 50 koneksi. Dua replica, ditambah migration atau `psql`, melewati batas default PostgreSQL yang 100 koneksi. Request berikutnya, termasuk login, gagal dengan `too many clients`.
- **Perbaikan:** Beri Authula batasnya sendiri yang lebih kecil, misalnya separuh dari `DB_MAX_OPEN_CONNS`, dan tulis di README bahwa satu replica memakai dua pool.
- **Kalau dibiarkan:** Menambah replica justru bisa menjatuhkan login.

### 8. Nama category berisi spasi saja diterima (`internal/lib/validator/rules.go:L13-21`, `internal/dtos/category.go:L14`)

- **Apa ini:** `Required` hanya menolak nilai kosong `""`.
- **Masalah:** `{"name": "   "}` lolos validasi, sehingga category tanpa nama tersimpan.
- **Perbaikan:** Di `Required`, anggap string kosong setelah `strings.TrimSpace` sebagai tidak diisi. Cukup ganti satu kondisi.
- **Kalau dibiarkan:** Data sampah muncul di daftar category dan harus dihapus manual.

## Info (lean)

### 9. Sisa demo yang tidak dipakai (`cmd/api/routes.go:L24-28`, `cmd/api/server.go:L22`, `L100`, `Dockerfile:L23-25`, `templates/`)

- **Apa ini:** Ada route `GET /` yang menjawab "Hello, World!", middleware `static` yang melayani `public/static/images/golang.jpg`, dan folder `templates/` yang kosong.
- **Masalah:** Tidak ada kode, test, atau dokumen yang memakainya. Middleware static tetap memeriksa filesystem di setiap request sebelum routing.
- **Perbaikan:** Hapus route, import dan `Use` static, `COPY public` di Dockerfile, serta folder `public/` dan `templates/`. Kalau folder itu sengaja disediakan untuk pemakai starter, sebutkan di README.
- **Kalau dibiarkan:** Kode tambahan yang ikut dibaca dan dipelihara tanpa manfaat.

### 10. Dua pembungkus dengan satu pemakai (`internal/dtos/common.go:L18-29`, `internal/lib/context.go:L22-25`)

- **Apa ini:** `ListQuery.base` hanya dipanggil `ListQuery.Validate`, dan komentarnya menyebut DTO lain yang belum ada. `lib.ContextParamUUID` hanya dipanggil `pathID` dan isinya satu baris `uuid.Parse(c.Params(key))`.
- **Masalah:** Ada lapisan tambahan yang harus diikuti saat membaca, tanpa pemakai kedua.
- **Perbaikan:** Pindahkan isi `base` ke `Validate`, lalu ganti panggilan `ContextParamUUID` dengan `uuid.Parse(c.Params("id"))` di `pathID`.
- **Kalau dibiarkan:** Kerugiannya kecil; ini hanya soal kerapian.

---

**Verdict:** No. 1 sampai 7 sudah diperbaiki. Yang tersisa bukan masalah keamanan (no. 8 sampai 10), ditambah sisa risiko no. 1 yang harus ditutup di Authula.

**Lean:** sekitar -20 baris, 0 dependency yang bisa dibuang.

**Tidak dicek:**

- Login Google, magic link, dan pengiriman email tidak dicoba sungguhan. Integration test dijalankan ke PostgreSQL 18 setelah perbaikan, tapi workflow release belum pernah dijalankan di GitHub.
- Di dalam Authula hanya alur sign-up, sign-in, reset password, callback OAuth2, magic link, dan session yang dibaca, bukan seluruh library.
- `graft/`, `.omc/`, `.claude/`, `docs/agents/`, dan `CHANGELOG.md` tidak dibaca.
