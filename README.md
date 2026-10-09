# AI-Powered Code Assessment & Autograder Platform

Platform pembelajaran pemrograman interaktif untuk bahasa C dengan evaluasi otomatis, sandbox execution, dan pembuatan soal berbasis AI.

## Daftar Isi
1. [Tech Stack](#tech-stack)
2. [Instalasi](#instalasi)
3. [Konfigurasi](#konfigurasi)
4. [Menjalankan Aplikasi](#menjalankan-aplikasi)
5. [Struktur Proyek](#struktur-proyek)

## Tech Stack
- **Frontend:** ReactJS, Vite, Tailwind CSS, Monaco Editor
- **Backend:** Golang, Gin Gonic, GORM, MySQL
- **AI/Execution:** Google Gemini API, Judge0 API
- **Auth:** JWT, BCrypt

## Instalasi
1. Clone repo: `git clone <url>`
2. Backend: `go mod download`
3. Frontend: `cd web && npm install`

## Konfigurasi
Salin `.env.example` ke `.env` dan isi variabel:
- `DB_DSN`: Koneksi MySQL
- `JWT_SECRET` & `JWT_REFRESH_SECRET`: Secret token
- `GEMINI_API_KEY`: API Key AI
- `JUDGE0_API_URL`: URL Judge0

## Menjalankan Aplikasi
1. **Infrastruktur:** `docker compose up -d`
   - Tunggu MySQL *healthy* (`docker compose ps` → STATUS `Up (healthy)`)
   - Jika gagal koneksi: `docker compose down -v && docker compose up -d`
2. **Backend:** `air` (live reload) atau `go run ./cmd/server`
   - Jika port 8080 bentrok: `taskkill /F /IM air.exe` lalu ulangi
3. **Frontend:** `cd web && npm run dev` (default port 5173)
4. **Admin:** `go run ./cmd/createadmin -name "Admin" -email admin@example.com -password "pass123"`

## Troubleshooting Umum
| Masalah | Solusi |
|---------|--------|
| `dial tcp 127.0.0.1:3306: connectex: refused` | MySQL belum siap. Tunggu `healthy` atau restart container. |
| `bind: Only one usage of each socket address` (port 8080) | `taskkill /F /PID $(Get-NetTCPConnection -LocalPort 8080 -State Listen).OwningProcess` |
| `JUDGE0_API_URL` salah format | Isi URL lengkap, misal `http://localhost:2358` (bukan `gcc:12`). |
| `driver: bad connection` / `unexpected EOF` | Normal saat startup pertama; retry otomatis hingga 10x. |

## Konfigurasi Penting
Perbaiki `.env` sebelum jalan:
- `JUDGE0_API_URL=http://localhost:2358` (URL Judge0 API, bukan image name)
- `DB_DSN` harus cocok dengan `docker-compose.yml` (user/pass/db)

## Struktur Proyek
- `/cmd`: Entry point server & CLI.
- `/internal`: Business logic (controllers, services, repositories).
- `/web`: Frontend React.
- `/migrations`: SQL Schema.

shortcut: dokumentasi minimalis sesuai AGENTS.md, perlu diperbarui jika API bertambah.
