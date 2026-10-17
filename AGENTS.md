# AGENTS.md

## Project Overview

This repository houses an **AI-Powered Code Assessment & Autograder Platform** backend and web application setup. The project serves as an interactive programming learning environment, specifically tailored for C programming courses with automated code evaluation, custom testcase execution via sandboxed execution engines, and AI-assisted problem generation.

---

## Tech Stack

| Category                  | Technology / Library                         | Description                                                  |
| :------------------------ | :------------------------------------------- | :----------------------------------------------------------- |
| **Frontend Core**         | ReactJS (Vite)                               | Single-page client framework powered by Vite                 |
| **Code Editor**           | `@monaco-editor/react`                       | Browser-based code editor with C language syntax support     |
| **Data Fetching**         | TanStack Query + Axios                       | Asynchronous server-state management with HTTP interceptors  |
| **State & Routing**       | Zustand / React Context & TanStack Router    | Global client state management and declarative routing       |
| **Forms & Validation**    | React Hook Form + Zod                        | Client-side form management and strict schema validation     |
| **Backend Framework**     | Golang with Gin Gonic                        | High-performance RESTful Web Framework                       |
| **CORS Middleware**       | `gin-contrib/cors`                           | Cross-Origin Resource Sharing handling for frontend requests |
| **Authentication**        | `golang-jwt/jwt/v5` + `bcrypt`               | Stateless JWT tokens (Access & Refresh) and password hashing |
| **Validation**            | `go-playground/validator/v10`                | DTO struct validation middleware                             |
| **Database & ORM**        | GORM + MySQL Driver (`gorm.io/driver/mysql`) | Relational ORM mapping with MySQL Local                      |
| **Config & Environment**  | `godotenv`                                   | Environment variable reader for `.env` files                 |
| **AI Integration**        | Google GenAI SDK (Gemini API)                | Automated testcase generation engine                         |
| **Code Execution Engine** | Docker (gcc:12)                              | Isolated sandbox container for C code compilation & run      |
| **Excel Generation**      | `github.com/xuri/excelize/v2`                | Streaming XLSX generation for reports/grades                 |
| **Database Seeding**      | Custom CLI (`cmd/seed`)                      | Idempotent dummy data generation for development             |

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
│   ├── middleware/           # JWT auth, CORS, enrollment guard, and implemented middleware
│   ├── models/               # GORM database structs & entities
│   ├── repositories/         # Database access layer (CRUD operations)
│   ├── services/             # Core business logic & external integrations
│   │   ├── ai/               # Gemini prompt processing & parsing
│   │   ├── autograder/       # Tolerance matcher & output evaluation
│   │   └── docker/           # Submission execution in the configured Docker sandbox image
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
├── docker-compose.yml        # Docker compose configuration for MySQL
├── Dockerfile                # Production image for the Go backend (not the C execution sandbox)
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
[ Middleware Layer ] (CORS and JWT Auth Token Verification)
  │
  ▼
[ Controller Layer ] (HTTP Parameter Parsing, DTO Validation)
  │
  ▼
[ Service Layer ] (Business Logic, AI Prompting, Docker Execution, Excel Generation)
  │
  ├─► [ Repository Layer ] (GORM Database Abstraction ──► MySQL)
  ├─► [ AI Engine ]       (Gemini API Test Case Generator)
  ├─► [ Execution Engine ] (Docker gcc:12 Sandbox)
  └─► [ Report Engine ]    (excelize/v2 Streaming XLSX)
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
4. Backend executes `source_code` + `stdin` in **Docker container (gcc:12)**.
5. Docker returns execution `stdout` / `stderr`.
6. **Tolerance Matcher Engine** compares student `stdout` against `expected_output` using sequence & keyword matching logic.
7. Evaluation score & status (`PASSED`, `WRONG_ANSWER`, `COMPILE_ERROR`) are recorded in MySQL.

### 3. Admin Dashboard & Reporting Flow

1. Admin manages Courses, Quizzes, and Student Enrollments via Admin Dashboard.
2. **Bulk Enrollment**: Admin searches students (by NIM/Name/Email), selects multiple, and enrolls them to a course in one action (`POST /admin/courses/:id/enroll-bulk`).
3. **Progress Tracking**: Admin toggles progress view per course to see student completion % and average scores.
4. **Export Grades/Progress**: Admin downloads Excel reports (`GET /admin/courses/:id/export-grades`) generated via `excelize/v2` streaming to HTTP response.
5. **Student Removal**: Admin removes student from course (`DELETE /admin/courses/:id/enrollments/:student_id`) with CASCADE delete handling via DB constraints.

### 4. Database Seeding (Development)

Run `go run ./cmd/seed` to populate DB with:

- 1 Admin user (`admin@cgrader.local` / `admin123`)
- 5 Student users (`password123`)
- 2 Courses with 3 Quizzes & Testcases
- Enrollments & Random Submissions (PASSED/WRONG_ANSWER)

---

## Development Commands

<!-- STREAMING_CHUNK:Listing actual dev commands... -->

### Backend (Golang)

- **Install dependencies:**

  ```bash
  go mod download
  ```

- **Install Air for live reload:**

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

- **Run server with Air (live update):**

  ```bash
  air
  ```

  Air watches the Go source files and automatically rebuilds/restarts the server when code changes. If the repository contains an `.air.toml` configuration file, use that configuration as the source of truth for watch/build settings.

- **Run server without live reload:**

  ```bash
  go run ./cmd/server
  ```

  Note: Use package path `./cmd/server`, not single file `cmd/server/main.go` to include all package files.

- **Build server binary:**
  ```bash
  go build -o build/server ./cmd/server
  ```
- **Run unit tests:**
  ```bash
  go test ./...
  ```
- **Run tests with coverage:**
  ```bash
  go test -v -cover ./...
  ```

### Frontend (ReactJS / Vite)

- **Install frontend dependencies:**
  ```bash
  cd web && npm install
  ```
- **Start Vite dev server:**
  ```bash
  cd web && npm run dev
  ```
- **Build frontend production package:**
  ```bash
  cd web && npm run build
  ```
- **Run frontend tests (if configured):** `cd web && npm test`. The repository must define a `test` script (for example, `vitest run`) before this command is a required check. If no test runner is configured, explicitly report that automated frontend tests are unavailable and perform lint/manual verification.
- **Lint frontend codebase:**
  ```bash
  cd web && npm run lint
  ```

### Local Services (Docker)

- **Start local MySQL container:**
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
- Docker and Docker Compose (recommended for local MySQL)
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
PORT=1111
DB_DSN=autograder:autograder_pass@tcp(127.0.0.1:3306)/autograder_db?parseTime=true&charset=utf8mb4&loc=Local
JWT_SECRET=<your-jwt-access-secret>
JWT_REFRESH_SECRET=<your-jwt-refresh-secret>
GEMINI_API_KEY=<your-gemini-api-key>
EXEC_DOCKER_IMAGE=gcc:12
EXEC_TIMEOUT_SECONDS=5
EXEC_MEMORY_LIMIT_MB=128
EXEC_CPU_LIMIT=0.5
```

The provided `.env.example` contains DB credentials matching `docker-compose.yml`, plus optional
`CORS_ALLOWED_ORIGINS`, `COOKIE_SECURE`, `GEMINI_MODEL`, and sandbox execution settings. Generate each JWT
secret with `openssl rand -hex 32`; use different values for `JWT_SECRET` and `JWT_REFRESH_SECRET`. The
server refuses to start if required secrets are missing or invalid.

Never commit `.env` or real API keys/secrets to the repository.

### 4. Start Local Infrastructure

The recommended local setup uses Docker Compose for MySQL. SQL files in `migrations/` are the schema source of truth for fresh installs and incremental database changes; GORM models must remain aligned with them. Do not rely on automatic GORM `AutoMigrate` as a substitute for reviewed SQL migrations in shared or production environments:

```bash
docker compose up -d
```

Check that the containers are running:

```bash
docker compose ps
```

The `migrations/` directory is mounted into the MySQL container (`/docker-entrypoint-initdb.d`), so
`migrations/001_init.sql` is applied automatically **only the first time** the MySQL volume is created. In this repository's documented fresh-install schema, `001_init.sql` is the initial schema and includes `users.nim` and `enrollments`; these are not described as a later addition. Never edit an already-applied migration to change an existing database: create a new numbered migration for subsequent schema changes.

To apply the initial schema manually (e.g. against an empty database):

```bash
docker compose exec -T mysql sh -c 'mysql -uroot -p"$MYSQL_ROOT_PASSWORD" autograder_db' < migrations/001_init.sql
```

To start from scratch: `docker compose down -v && docker compose up -d` (this deletes all data).

> **Code Execution:** This project uses Docker directly (`gcc:12` image) for sandboxed C code execution. No separate Judge0 or Piston container is required. The root `Dockerfile` builds the Go backend for deployment; the runtime-pulled `gcc:12` (or `EXEC_DOCKER_IMAGE`) image is a separate sandbox for submitted C code. Ensure Docker Desktop (WSL2 backend on Windows) is running.

### 5. Install Backend Dependencies

go mod tidy
go mod download

````

Install Air if it has not been installed yet:

```bash
go install github.com/air-verse/air@latest
export PATH="$PATH:$(go env GOPATH)/bin"
````

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

The backend and frontend should now run independently (default: API on `http://localhost:1111`, web on `http://localhost:5173`). The frontend API base URL defaults to `http://localhost:1111/api/v1`.

Public registration only creates `STUDENT` accounts. Create an `ADMIN` account with:

```bash
go run ./cmd/createadmin -name "Admin" -email admin@example.com -nim ADMIN001 -password "change-me-please"
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

| Variable Name          | Purpose                                                                                          | Required                              | Example Placeholder                                                                                     |
| :--------------------- | :----------------------------------------------------------------------------------------------- | :------------------------------------ | :------------------------------------------------------------------------------------------------------ |
| `PORT`                 | Server listening port                                                                            | No (Default: 1111)                    | `1111`                                                                                                  |
| `DB_DSN`               | MySQL connection string DSN                                                                      | Yes                                   | `autograder:autograder_pass@tcp(127.0.0.1:3306)/autograder_db?parseTime=true&charset=utf8mb4&loc=Local` |
| `JWT_SECRET`           | Secret key for signing Access Tokens; use at least 32 random bytes                               | Yes                                   | `$(openssl rand -hex 32)`                                                                               |
| `JWT_REFRESH_SECRET`   | Secret key for Refresh Tokens; use at least 32 random bytes and keep different from `JWT_SECRET` | Yes                                   | `$(openssl rand -hex 32)`                                                                               |
| `GEMINI_API_KEY`       | Google Gemini AI API Key                                                                         | Yes                                   | `<your-gemini-api-key>`                                                                                 |
| `CORS_ALLOWED_ORIGINS` | Comma-separated trusted frontend origins (never `*`)                                             | No (Default: `http://localhost:5173`) | `http://localhost:5173`                                                                                 |
| `COOKIE_SECURE`        | Mark refresh-token cookie as `Secure` (set `true` behind HTTPS)                                  | No (Default: `false`)                 | `false`                                                                                                 |
| `GEMINI_MODEL`         | Gemini model used for testcase generation                                                        | No (Default: `gemini-2.5-flash`)      | `gemini-2.5-flash`                                                                                      |
| `EXEC_DOCKER_IMAGE`    | Docker image used to compile/run submitted C code                                                | No (Default: `gcc:12`)                | `gcc:12`                                                                                                |
| `EXEC_TIMEOUT_SECONDS` | Maximum execution time per testcase                                                              | No (Default: `5`)                     | `5`                                                                                                     |
| `EXEC_MEMORY_LIMIT_MB` | Memory limit for each sandbox container                                                          | No (Default: `128`)                   | `128`                                                                                                   |
| `EXEC_CPU_LIMIT`       | CPU quota for each sandbox container                                                             | No (Default: `0.5`)                   | `0.5`                                                                                                   |

Frontend (`web/.env`, optional): `VITE_API_URL` (default `http://localhost:1111/api/v1`).

---

## API Conventions

- **Base URL:** `/api/v1`
- **Response Format:**
  ```json
  {
    "success": true,
    "message": "Operation completed successfully",
    "data": {}
  }
  ```
- **Error Format:** `errors` is an array of objects with `field` and `message` string properties. For non-field errors, `field` may be an empty string; use an empty array when there are no structured details.

  ```json
  {
    "success": false,
    "message": "Error details or validation message",
    "errors": [{ "field": "email", "message": "must be a valid email address" }]
  }
  ```

- **Authentication Header:** `Authorization: Bearer <JWT_ACCESS_TOKEN>`

---

## Request Examples

Register a student (public route):

```bash
curl -X POST http://localhost:1111/api/v1/auth/register \\
  -H "Content-Type: application/json" \\
  -d '{"name":"Student Example","email":"student@example.com","password":"change-me-please","nim":"12345678"}'
```

Bulk-enroll students (admin token required):

```bash
curl -X POST http://localhost:1111/api/v1/admin/courses/1/enroll-bulk \\
  -H "Authorization: Bearer $TOKEN" \\
  -H "Content-Type: application/json" \\
  -d '{"student_ids":[12,13]}'
```

Submit a C solution (student token required):

```bash
curl -X POST http://localhost:1111/api/v1/submissions \\
  -H "Authorization: Bearer $TOKEN" \\
  -H "Content-Type: application/json" \\
  -d '{"quiz_id":1,"source_code":"#include <stdio.h>\\nint main(void){printf(\\"Hello\\\\n\\");return 0;}"}'
```

These examples illustrate the API base URL and authentication pattern. Verify request field names and response payloads against the current DTOs/controllers before using them as integration fixtures.

---

## Database

- **Database Engine:** MySQL 8.0+
- **ORM:** GORM (`gorm.io/gorm`)
- **Core Entities:**
  - `User`: Stores credentials, roles (`STUDENT`, `ADMIN`), and timestamps.
  - `Course`: Learning modules containing multiple quizzes.
  - `Quiz`: Problem statements and constraints created by admins.
  - `TestCase`: `stdin` and `expected_output` generated by AI or defined manually.
  - `Submission`: Student code submissions, execution status, and scores.

---

## Authentication & Authorization

- **Passwords:** Hashed using `bcrypt` (default cost: 10 or higher).
- **Token Strategy:**
  - **Access Token:** Short-lived JWT passed via `Authorization` HTTP header.
  - **Refresh Token:** Stored in secure `HttpOnly` Cookie.
- **Role-Based Access Control (RBAC):**
  - `STUDENT`: Can view courses, solve quizzes, submit code.
  - `ADMIN`: Access to Admin Dashboard, create quizzes, trigger AI testcase generation.

---

### Route Permission Matrix

All paths below are relative to `/api/v1`. “Public” means no access token is required; protected routes must enforce the role and enrollment rules shown.

| Method   | Route                                        | Access                                                                                      |
| :------- | :------------------------------------------- | :------------------------------------------------------------------------------------------ |
| `POST`   | `/auth/register`                             | Public; creates `STUDENT` accounts only                                                     |
| `POST`   | `/auth/login`                                | Public                                                                                      |
| `POST`   | `/auth/refresh`                              | Public route; requires a valid refresh-token cookie                                         |
| `POST`   | `/auth/logout`                               | Authenticated                                                                               |
| `GET`    | `/auth/me`                                   | Authenticated                                                                               |
| `GET`    | `/courses`                                   | Authenticated; admin sees all courses, student sees only `APPROVED` enrolled courses          |
| `GET`    | `/courses/me`                                | Authenticated; returns only `APPROVED` courses for the user                                 |
| `GET`    | `/courses/:id`                               | Enrolled student or admin; non-enrolled student gets `403`, missing course gets `404`       |
| `GET`    | `/quizzes/:id`                               | Enrolled student or admin; course resolved from quiz, not client input                      |
| `POST`   | `/submissions`                               | Enrolled student; course resolved from `quiz_id` before executing any code                  |
| `POST`   | `/submissions/run`                           | Enrolled student; course resolved from `quiz_id` before executing any code (does not create submission record) |
| `GET`    | `/submissions`                               | Authenticated; own submissions only                                                         |
| `GET`    | `/submissions/:id`                           | Authenticated; own submission, or admin                                                     |
| `GET`    | `/student/enrollments`                       | Authenticated                                                                               |
| `POST`   | `/admin/courses`                             | Admin                                                                                       |
| `PUT`    | `/admin/courses/:id`                         | Admin                                                                                       |
| `DELETE` | `/admin/courses/:id`                         | Admin                                                                                       |
| `POST`   | `/admin/quizzes`                             | Admin                                                                                       |
| `POST`   | `/admin/quizzes/generate-testcases`          | Admin                                                                                       |
| `GET`    | `/admin/quizzes/:id`                         | Admin                                                                                       |
| `PUT`    | `/admin/quizzes/:id`                         | Admin                                                                                       |
| `DELETE` | `/admin/quizzes/:id`                         | Admin                                                                                       |
| `POST`   | `/admin/quizzes/:id/testcases`               | Admin                                                                                       |
| `DELETE` | `/admin/testcases/:id`                       | Admin                                                                                       |
| `GET`    | `/admin/courses/:id/progress`                | Admin                                                                                       |
| `GET`    | `/admin/courses/:id/export-grades`           | Admin                                                                                       |
| `GET`    | `/admin/courses/:id/enrollments`             | Admin                                                                                       |
| `GET`    | `/admin/students`                            | Admin; search + pagination of `STUDENT` accounts                                            |
| `POST`   | `/admin/courses/:id/enroll-bulk`             | Admin                                                                                       |
| `DELETE` | `/admin/courses/:id/enrollments/:student_id` | Admin                                                                                       |

Every new or changed route must be added to this matrix. Verify exact route registration and middleware in `cmd/server/routes.go`; this table is the source of truth for route permissions.

---

## Validation & Error Handling

- **Request Payload Validation:** Handled via `go-playground/validator/v10` on Go struct DTO tags (`binding:"required,email"`).
- **Global Error Middleware:** Captures unhandled panics and formats standard JSON error responses.

---

## Testing

- **Backend:** Native Go testing framework (`testing` package). Executable via `go test ./...`.
- **Frontend:** ESLint checks via `npm run lint`. If the frontend does not yet have a test runner configured, document that automated frontend tests are not available and perform manual verification; do not claim `npm test` exists until a test script (for example, `vitest run`) is added to `web/package.json`.
- **Autograder:** Add unit tests for the tolerance matcher, including whitespace, sequence, and keyword behavior.
- **Docker execution:** Test correct output, wrong answer, compilation failure, timeout, and Docker unavailable/connection failure. Tests must never execute untrusted C code directly on the host.
- **AI testcase generation:** Test invalid JSON, valid JSON with schema mismatch, empty testcase arrays, and valid schema responses before persistence.
- **Enrollment/access control:** Test enrolled and non-enrolled student access, missing courses, and admin-only routes.

---

## Security

- Cross-Origin Resource Sharing (CORS) restricted to trusted frontend domain.
- Password hashing with bcrypt before database persistence.
- SQL Injection protection enforced via GORM parameterized queries.
- Execution sandboxing offloaded to Docker containers (`gcc:12`) to prevent remote code execution (RCE) on backend server host.

---

## Important Invariants

1. **Sandboxed Code Execution:** NEVER compile or run untrusted user C code directly on the host machine running the Gin server. All executions MUST go through Docker containers (`gcc:12` image).
2. **AI Test Case Output:** AI Testcase responses MUST be validated as proper JSON arrays matching the strict schema before storing in MySQL.
3. **Database Relationships:** Quiz records MUST be deleted with cascading deletes on associated `TestCases` and `Submissions`.
4. **JWT Security:** Access token verification MUST be enforced on all routes except `/auth/login`, `/auth/register`, and `/auth/refresh`.

---

## Common Pitfalls

- **Database Migration Drift:** Keep SQL migrations and GORM models synchronized; use a new numbered migration for schema changes after `001_init.sql` has been applied.
- **Missing CORS Headers:** Forgetting to configure credentials support in `gin-contrib/cors` when using HTTP-Only cookies for refresh tokens.
- **Direct Database Queries:** Bypassing repository interfaces in controller handlers.
- **Parsing Double Escaped Sequences:** Incorrectly parsing newline sequences (`\n`) in `stdin` strings received from Gemini API.

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
- After implementation, review `Project Setup` and update it whenever dependencies, environment variables, migrations, local services, or development commands change. Do not claim setup is complete unless documented commands have been verified.

---

## Verification Checklist

- [ ] All database migration files match GORM models.
- [ ] No hardcoded secrets or API keys exist in source code.
- [ ] JWT authentication middleware wraps protected endpoints.
- [ ] Docker execution handles network timeout and connection failure gracefully.
- [ ] Frontend Axios interceptor correctly handles token refresh on 401 response status.
