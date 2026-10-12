# AI-Powered Code Assessment & Autograder Platform

Platform pembelajaran pemrograman interaktif untuk bahasa C dengan evaluasi otomatis, sandbox execution, dan pembuatan soal berbasis AI.

## Struktur & Dokumentasi
Untuk detail arsitektur, spesifikasi teknis, panduan setup, dan aturan pengembangan, silakan merujuk ke **[AGENTS.md](AGENTS.md)**.

## Quick Start
1. **Clone & Setup:**
   ```bash
   git clone <repository-url>
   cd <project-directory>
   cp .env.example .env
   # Update .env dengan konfigurasi yang sesuai (DB_DSN, JWT_SECRET, GEMINI_API_KEY, PISTON_API_URL)
   ```

2. **Infrastruktur (Docker):**
   ```bash
   docker compose up -d
   ```

3. **Backend:**
   ```bash
   go mod download
   # Install air (live-reload): go install github.com/air-verse/air@latest
   air
   ```

4. **Frontend:**
   ```bash
   cd web
   npm install
   npm run dev
   ```

5. **Admin Setup:**
   ```bash
   go run ./cmd/createadmin -name "Admin" -email admin@example.com -password "change-me-please"
   ```

---
*Proyek ini mengikuti pola Clean Architecture. Silakan baca [AGENTS.md](AGENTS.md) sebelum memulai pengembangan atau modifikasi kode.*
