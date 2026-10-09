<!-- STREAMING_CHUNK:Writing project overview and tech stack... -->
# AGENTS.md

## Project Overview

This repository houses an **AI-Powered Code Assessment & Autograder Platform** backend and web application setup. The project serves as an interactive programming learning environment, specifically tailored for C programming courses with automated code evaluation, custom testcase execution via sandboxed execution engines, and AI-assisted problem generation.

---

## Tech Stack

| Category | Technology / Library | Description |
| :--- | :--- | :--- |
| **Frontend Core** | ReactJS (Vite) | Single-page client framework powered by Vite |
| **Code Editor** | `@monaco-editor/react` | Browser-based code editor with C language syntax support |
| **Data Fetching** | TanStack Query + Axios | Asynchronous server-state management with HTTP interceptors |
| **State & Routing**| Zustand / React Context & TanStack Router | Global client state management and declarative routing |
| **Forms & Validation** | React Hook Form + Zod | Client-side form management and strict schema validation |
| **Backend Framework** | Golang with Gin Gonic | High-performance RESTful Web Framework |
| **CORS Middleware** | `gin-contrib/cors` | Cross-Origin Resource Sharing handling for frontend requests |
| **Authentication** | `golang-jwt/jwt/v5` + `bcrypt` | Stateless JWT tokens (Access & Refresh) and password hashing |
| **Validation** | `go-playground/validator/v10` | DTO struct validation middleware |
| **Database & ORM** | GORM + MySQL Driver (`gorm.io/driver/mysql`) | Relational ORM mapping with MySQL Local |
| **Config & Environment** | `godotenv` | Environment variable reader for `.env` files |
| **AI Integration** | Google GenAI SDK (Gemini API) / OpenAI SDK | Automated testcase generation engine |
| **Code Execution Engine**| Judge0 API | Isolated sandbox container API for C code compilation & run |

---

## Repository Structure

<!-- STREAMING_CHUNK:Defining repository tree structure... -->
```text
.
├── cmd/
│   ├── server/
│   │   ├── main.go           # Application bootstrap & server initialization
│   │   └── routes.go         # Route registration & middleware wiring
│   └── createadmin/
│       └── main.go           # CLI untuk membuat akun ADMIN
├── internal/
│   ├── config/               # Environment & app configuration handlers
│   ├── controllers/          # HTTP request handlers & JSON response encoders
│   ├── middleware/           # JWT Auth, CORS, logging, & rate limiters
│   ├── models/               # GORM database structs & entities
│   ├── repositories/         # Database access layer (CRUD operations)
│   ├── services/             # Core business logic & external integrations
│   │   ├── ai/               # Gemini/OpenAI prompt processing & parsing
│   │   ├── autograder/       # Tolerance matcher & output evaluation
│   │   └── judge0/           # Submission proxy to Judge0 Sandbox API
│   └── utils/                # Helper functions, JWT signing, password hashing
├── migrations/               # Database SQL schema migration files (001_init.sql)
├── web/                      # ReactJS Frontend (Vite workspace)
│   ├── src/
│   │   ├── components/       # Reusable UI components & Monaco Editor wrapper
│   │   ├── hooks/            # TanStack Query custom hooks
│   │   ├── pages/            # View components (Workspace, Admin, Auth)
│   │   ├── services/         # Axios instance with interceptors
│   │   └── store/            # Zustand global state management
│   ├── package.json          # Frontend dependency specifications
│   ├── tailwind.config.js    # Tailwind CSS configuration
│   ├── eslint.config.js      # ESLint flat config
│   └── vite.config.js        # Vite bundler configuration
├── .air.toml                 # Air live-reload configuration
├── .env.example              # Sample environment variables template
├── docker-compose.yml        # Docker compose configuration for MySQL & Judge0
├── Dockerfile                # Production Docker build container
├── go.mod                    # Go module dependencies
└── go.sum                    # Go checksum verification (dibuat oleh `go mod tidy`)
```

---

## Architecture

The project follows a **Layered Clean Architecture** pattern with clear separation of concerns between HTTP transport, domain business logic, data access, and third-party API integrations:

```text
Request
  │
  ▼
[ Middleware Layer ] (CORS, JWT Auth Token Verification, Rate Limiting)
  │
  ▼
[ Controller Layer ] (HTTP Parameter Parsing, DTO Validation)
  │
  ▼
[ Service Layer ] (Business Logic, AI Prompting, Judge0 Proxying)
  │
  ├─► [ Repository Layer ] (GORM Database Abstraction ──► MySQL)
  ├─► [ AI Engine ]       (Gemini API Test Case Generator)
  └─► [ Execution Engine ] (Judge0 API Compiler Sandbox)
```

---

## Application Flow

### 1. Test Case Generation Flow (Admin)
1. Admin submits problem statement via Frontend Admin Dashboard.
2. Request hits `/api/v1/admin/quizzes/generate-testcases` via Gin Handler.
3. **AI Engine Service** sends prompt template to Gemini API.
4. Gemini returns a strict JSON array of input standard strings (`stdin`) and expected outputs (`expected_output`).
5. Response is stored in MySQL via GORM `TestCases` repository.

### 2. Code Autograding Submission Flow (Student)
1. Student writes C code in **Monaco Editor** and submits solution.
2. Request hits `/api/v1/submissions` with `quiz_id` and `source_code`.
3. Backend fetches problem testcases from database.
4. Backend proxies `source_code` + `stdin` payload to **Judge0 Sandbox API**.
5. Judge0 returns execution `stdout` / `stderr`.
6. **Tolerance Matcher Engine** compares student `stdout` against `expected_output` using sequence & keyword matching logic.
7. Evaluation score & status (`PASSED`, `WRONG_ANSWER`, `COMPILE_ERROR`) are recorded in MySQL.

---

## Development Commands

<!-- STREAMING_CHUNK:Listing actual dev commands... -->
### Backend (Golang)
* **Install dependencies:**
  ```bash
  go mod download
  ```

* **Install Air for live reload:**
  ```bash
  go install github.com/air-verse/air@latest
  ```
  Ensure the Go binary directory is available in your `PATH`:
  ```bash
  export PATH="$PATH:$(go env GOPATH)/bin"
  ```
  Verify the installation:
  ```bash
  air -v
  ```

* **Run server with Air (live update):**
  ```bash
  air
  ```
  Air watches the Go source files and automatically rebuilds/restarts the server when code changes. If the repository contains an `.air.toml` configuration file, use that configuration as the source of truth for watch/build settings.

* **Run server without live reload:**
  ```bash
  go run ./cmd/server
  ```
  Note: Use package path `./cmd/server`, not single file `cmd/server/main.go` to include all package files.

* **Build server binary:**
  ```bash
  go build -o build/server cmd/server/main.go
  ```
* **Run unit tests:**
  ```bash
  go test ./...
  ```
* **Run tests with coverage:**
  ```bash
  go test -v -cover ./...
  ```

### Frontend (ReactJS / Vite)
* **Install frontend dependencies:**
  ```bash
  cd web && npm install
  ```
* **Start Vite dev server:**
  ```bash
  cd web && npm run dev
  ```
* **Build frontend production package:**
  ```bash
  cd web && npm run build
  ```
* **Lint frontend codebase:**
  ```bash
  cd web && npm run lint
  ```

### Local Services (Docker)
* **Start local MySQL & Judge0 containers:**
  ```bash
  docker-compose up -d
  ```

---

## Project Setup

Follow these steps when setting up the project from a fresh clone. After implementing or changing project code, update this section if the setup process, dependencies, environment variables, migrations, or development commands change.

### 1. Prerequisites

Install the following tools:

- Go (version compatible with `go.mod`)
- Node.js and npm
- MySQL 8.0+
- Docker and Docker Compose (recommended for local MySQL and Judge0)
- Git
- Air for Go live reload during development

Verify the required tools:

```bash
go version
node --version
npm --version
mysql --version
docker --version
docker compose version
git --version
```

### 2. Clone the Repository

```bash
git clone <repository-url>
cd <project-directory>
```

### 3. Configure Backend Environment

Create the backend environment file from the provided template:

```bash
cp .env.example .env
```

Edit `.env` and provide the required values:

```env
PORT=8080
DB_DSN=user:pass@tcp(127.0.0.1:3306)/autograder_db?parseTime=true
JWT_SECRET=<your-jwt-access-secret>
JWT_REFRESH_SECRET=<your-jwt-refresh-secret>
GEMINI_API_KEY=<your-gemini-api-key>
JUDGE0_API_URL=http://localhost:2358
```

The provided `.env.example` already contains DB credentials matching `docker-compose.yml`, plus the optional
`CORS_ALLOWED_ORIGINS`, `COOKIE_SECURE`, and `GEMINI_MODEL`. `JWT_SECRET` and `JWT_REFRESH_SECRET` must be
different values (e.g. `openssl rand -hex 32`). The server refuses to start otherwise.

Never commit `.env` or real API keys/secrets to the repository.

### 4. Start Local Infrastructure

The recommended local setup uses Docker Compose for MySQL and Judge0:

```bash
docker compose up -d
```

Check that the containers are running:

```bash
docker compose ps
```

The `migrations/` directory is mounted into the MySQL container (`/docker-entrypoint-initdb.d`), so
`migrations/001_init.sql` is applied automatically **only the first time** the MySQL volume is created.

To apply it manually (e.g. against an existing database; the file is idempotent):

```bash
docker compose exec -T mysql sh -c 'mysql -uroot -p"$MYSQL_ROOT_PASSWORD" autograder_db' < migrations/001_init.sql
```

To start from scratch: `docker compose down -v && docker compose up -d` (this deletes all data).

> Judge0 1.13.x requires cgroup v1. On hosts using only cgroup v2 (recent Docker Desktop / Linux kernels),
> the Judge0 containers may fail to execute code; see the Judge0 documentation for the required kernel/GRUB setting.

### 5. Install Backend Dependencies

From the repository root. On the first setup, `go.sum` does not exist yet; generate it (requires internet):

```bash
go mod tidy
go mod download
```

Install Air if it has not been installed yet:

```bash
go install github.com/air-verse/air@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```

Verify:

```bash
air -v
```

### 6. Install Frontend Dependencies

```bash
cd web
npm install
cd ..
```

### 7. Run the Project in Development

Start the backend with live reload:

```bash
air
```

In a separate terminal, start the frontend:

```bash
cd web
npm run dev
```

The backend and frontend should now run independently (default: API on `http://localhost:8080`, web on `http://localhost:5173`).

Public registration only creates `STUDENT` accounts. Create an `ADMIN` account with:

```bash
go run ./cmd/createadmin -name "Admin" -email admin@example.com -password "change-me-please"
```

### 8. Verify the Setup

Run backend tests:

```bash
go test ./...
```

Run frontend linting:

```bash
cd web
npm run lint
```

Before considering the setup complete, verify that:

- MySQL is reachable using `DB_DSN`.
- Judge0 is reachable using `JUDGE0_API_URL`.
- Required AI credentials are present.
- The backend starts successfully with `air`.
- The frontend starts successfully with `npm run dev`.
- Authentication and protected API routes work as expected.

### 9. Setup Documentation Requirement for AI Agents

When implementation is complete, the AI agent MUST review this setup documentation and update it when necessary.

The documentation must reflect:

1. Required software and tool versions.
2. Dependency installation commands.
3. Environment variables and `.env` setup.
4. Database creation and migration steps.
5. Docker or other required local services.
6. Air/live-reload configuration for the Go backend.
7. Frontend installation and development commands.
8. Test, lint, build, and verification commands.
9. Any additional setup steps introduced by the implementation.

Do not claim that the project setup is complete unless the documented commands have been verified against the implemented project.

---

## Environment Variables

The backend relies on the following configuration keys (template available in `.env.example`):

| Variable Name | Purpose | Required | Example Placeholder |
| :--- | :--- | :--- | :--- |
| `PORT` | Server listening port | No (Default: 8080) | `8080` |
| `DB_DSN` | MySQL connection string DSN | Yes | `user:pass@tcp(127.0.0.1:3306)/autograder_db?parseTime=true` |
| `JWT_SECRET` | Secret key for signing Access Tokens | Yes | `<your-jwt-access-secret>` |
| `JWT_REFRESH_SECRET` | Secret key for Refresh Tokens | Yes | `<your-jwt-refresh-secret>` |
| `GEMINI_API_KEY` | Google Gemini AI API Key | Yes | `<your-gemini-api-key>` |
| `JUDGE0_API_URL` | Base URL of Judge0 execution sandbox | Yes | `http://localhost:2358` |
| `CORS_ALLOWED_ORIGINS` | Comma-separated trusted frontend origins (never `*`) | No (Default: `http://localhost:5173`) | `http://localhost:5173` |
| `COOKIE_SECURE` | Mark refresh-token cookie as `Secure` (set `true` behind HTTPS) | No (Default: `false`) | `false` |
| `GEMINI_MODEL` | Gemini model used for testcase generation | No (Default: `gemini-2.5-flash`) | `gemini-2.5-flash` |

Frontend (`web/.env`, optional): `VITE_API_URL` (default `http://localhost:8080/api/v1`).

---

## API Conventions

* **Base URL:** `/api/v1`
* **Response Format:**
  ```json
  {
    "success": true,
    "message": "Operation completed successfully",
    "data": {}
  }
  ```
* **Error Format:**
  ```json
  {
    "success": false,
    "message": "Error details or validation message",
    "errors": []
  }
  ```
* **Authentication Header:** `Authorization: Bearer <JWT_ACCESS_TOKEN>`

---

## Database

* **Database Engine:** MySQL 8.0+
* **ORM:** GORM (`gorm.io/gorm`)
* **Core Entities:**
  * `User`: Stores credentials, roles (`STUDENT`, `ADMIN`), and timestamps.
  * `Course`: Learning modules containing multiple quizzes.
  * `Quiz`: Problem statements and constraints created by admins.
  * `TestCase`: `stdin` and `expected_output` generated by AI or defined manually.
  * `Submission`: Student code submissions, execution status, and scores.

---

## Authentication & Authorization

* **Passwords:** Hashed using `bcrypt` (default cost: 10 or higher).
* **Token Strategy:**
  * **Access Token:** Short-lived JWT passed via `Authorization` HTTP header.
  * **Refresh Token:** Stored in secure `HttpOnly` Cookie.
* **Role-Based Access Control (RBAC):**
  * `STUDENT`: Can view courses, solve quizzes, submit code.
  * `ADMIN`: Access to Admin Dashboard, create quizzes, trigger AI testcase generation.

---

## Validation & Error Handling

* **Request Payload Validation:** Handled via `go-playground/validator/v10` on Go struct DTO tags (`binding:"required,email"`).
* **Global Error Middleware:** Captures unhandled panics and formats standard JSON error responses.

---

## Testing

* **Backend:** Native Go testing framework (`testing` package). Executable via `go test ./...`.
* **Frontend:** ESLint checks via `npm run lint`.

---

## Security

* Cross-Origin Resource Sharing (CORS) restricted to trusted frontend domain.
* Password hashing with bcrypt before database persistence.
* SQL Injection protection enforced via GORM parameterized queries.
* Execution sandboxing offloaded to Judge0 API to prevent remote code execution (RCE) on backend server host.

---

## Important Invariants

1. **Sandboxed Code Execution:** NEVER compile or run untrusted user C code directly on the host machine running the Gin server. All executions MUST go through Judge0 API.
2. **AI Test Case Output:** AI Testcase responses MUST be validated as proper JSON arrays matching the strict schema before storing in MySQL.
3. **Database Relationships:** Quiz records MUST be deleted with cascading deletes on associated `TestCases` and `Submissions`.
4. **JWT Security:** Access token verification MUST be enforced on all routes except `/auth/login`, `/auth/register`, and `/auth/refresh`.

---

## Common Pitfalls

* **Missing CORS Headers:** Forgetting to configure credentials support in `gin-contrib/cors` when using HTTP-Only cookies for refresh tokens.
* **Direct Database Queries:** Bypassing repository interfaces in controller handlers.
* **Parsing Double Escaped Sequences:** Incorrectly parsing newline sequences (`\n`) in `stdin` strings received from Gemini API.

---

## Rules for AI Agents

<!-- STREAMING_CHUNK:Documenting AI agent operational guidelines... -->
### Before Modifying Code
- Inspect related files in `internal/controllers/`, `internal/services/`, and `web/src/`.
- Ensure new endpoints follow existing `/api/v1` path conventions.
- Verify environment variable requirements in `internal/config/`.

### While Modifying Code
- Maintain separation of concerns: Do not place database logic in controllers.
- Use GORM parameterized queries to prevent SQL injection.
- Keep frontend components modular and follow Tailwind CSS styling patterns.
- Ensure all Go code handles errors explicitly without ignoring return values.

### After Modifying Code
- Run `go test ./...` to verify backend integrity.
- Run `cd web && npm run lint` to ensure frontend syntax compliance.
- Confirm schema consistency across GORM models and MySQL migration files.
- Update the `Project Setup` section when implementation changes affect dependencies, environment variables, database setup, services, or development commands.
- Ensure the setup documentation describes the final implemented project, not only the original project structure.

---

## Verification Checklist

- [ ] All database migration files match GORM models.
- [ ] No hardcoded secrets or API keys exist in source code.
- [ ] JWT authentication middleware wraps protected endpoints.
- [ ] Judge0 proxy handles network timeout and connection failure gracefully.
- [ ] Frontend Axios interceptor correctly handles token refresh on 401 response status.