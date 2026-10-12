# AI-Powered C Code Autograder

Interactive programming learning environment for C programming with automated code evaluation, sandboxed execution via Docker, and AI-assisted problem generation.

## Tech Stack
*   **Backend**: Golang (Gin Gonic)
*   **Frontend**: ReactJS (Vite)
*   **Database**: MySQL
*   **Execution Engine**: Docker (GCC)

## Setup

### 1. Prerequisites
*   Go 1.21+
*   Node.js 18+
*   Docker & Docker Compose

### 2. Environment Configuration
Create the backend environment file:
```bash
cp .env.example .env
```
Update `.env` with the following configuration:
```env
PORT=9999
DB_DSN=autograder:autograder_pass@tcp(127.0.0.1:3306)/autograder_db?parseTime=true&charset=utf8mb4&loc=Local
JWT_SECRET=<your-jwt-access-secret>
JWT_REFRESH_SECRET=<your-jwt-refresh-secret>
GEMINI_API_KEY=<your-gemini-api-key>
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173
```

### 3. Start Infrastructure
Start the database container:
```bash
docker compose up -d
```

### 4. Running the Application

**Backend:**
```bash
go run ./cmd/server
```
The server will start on `http://localhost:9999`.

**Frontend:**
```bash
cd web
npm install
npm run dev
```
The frontend will start on `http://localhost:5173`.

## Architecture Details
For detailed architectural information, please refer to [AGENTS.md](AGENTS.md).
