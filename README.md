# AGENTS.md

## Project Overview

This repository houses an **AI-Powered Code Assessment & Autograder Platform** backend and web application setup. The project serves as an interactive programming learning environment, specifically tailored for C programming courses with automated code evaluation, custom testcase execution via sandboxed execution engines, and AI-assisted problem generation.

Access to course content is **admin-controlled**: a student can only see and work on the courses an admin has explicitly enrolled them in (see [Enrollment System](#enrollment-system)).

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
| **AI Integration**        | Google GenAI SDK (Gemini API) / OpenAI SDK   | Automated testcase generation engine                         |
| **Code Execution Engine** | Docker (gcc:12)                              | Isolated sandbox container for C code compilation & run      |
| **Excel Generation**      | `github.com/xuri/excelize/v2`                | Streaming XLSX generation for reports/grades                 |
| **Database Seeding**      | Custom CLI (`cmd/seed`)                      | Idempotent dummy data generation for development             |

---

## Repository Structure

```text
.
├── cmd/
│   ├── server/
│   │   ├── main.go           # Application bootstrap & server initialization
│   │   └── routes.go         # Route registration & middleware wiring
│   ├── createadmin/
│   │   └── main.go           # CLI to create ADMIN accounts
│   └── seed/
│       └── main.go           # Idempotent dummy data seeder (development only)
├── internal/
│   ├── config/               # Environment & app configuration handlers
│   ├── controllers/          # HTTP request handlers & JSON response encoders
│   │   ├── enrollment_controller.go   # Admin enrollment endpoints (bulk enroll, remove, list)
│   │   └── course_controller.go       # Course CRUD + student "my courses" + progress/export
│   ├── middleware/           # JWT Auth, RBAC, enrollment guard, CORS, logging, & rate limiters
│   │   ├── auth.go                    # JWT verification + RequireRole
│   │   └── enrollment.go              # RequireEnrollment / RequireQuizEnrollment / RequireEnrollmentForSubmission
│   ├── models/               # GORM database structs & entities (models.go)
│   ├── repositories/         # Database access layer (CRUD operations)
│   │   └── enrollment_repository.go   # Enrollment queries (FindByUserAndCourse, ListByCourse, ...)
│   ├── services/             # Core business logic & external integrations
│   │   ├── ai/               # Gemini prompt processing & parsing
│   │   ├── autograder/       # Tolerance matcher & output evaluation
│   │   ├── docker/           # Submission execution via Docker (gcc:12)
│   │   └── enrollment_service.go      # Enrollment business rules (bulk, remove, access check)
│   └── utils/                # Helper functions, JWT signing, password hashing
├── migrations/               # Database SQL schema migration files (001_init.sql)
├── web/                      # ReactJS Frontend (Vite workspace)
│   ├── src/
│   │   ├── components/       # Reusable UI components & Monaco Editor wrapper
│   │   ├── hooks/            # TanStack Query custom hooks (incl. useEnrollments, useMyCourses)
│   │   ├── pages/            # View components (Workspace, Admin, Auth)
│   │   │   ├── admin/CourseEnrollments   # Admin: enrolled list, student picker, bulk enroll, remove
│   │   │   └── student/MyCourses         # Student: list of enrolled courses only
│   │   ├── services/         # Axios instance with interceptors
│   │   └── store/            # Zustand global state management
│   ├── package.json          # Frontend dependency specifications
│   ├── tailwind.config.js    # Tailwind CSS configuration
│   ├── eslint.config.js      # ESLint flat config
│   └── vite.config.js        # Vite bundler configuration
├── .air.toml                 # Air live-reload configuration
├── .env.example              # Sample environment variables template
├── docker-compose.yml        # Docker compose configuration for MySQL
├── Dockerfile                # Production Docker build container
├── go.mod                    # Go module dependencies
└── go.sum                    # Go checksum verification (generated by `go mod tidy`)
```

---

## Architecture

The project follows a **Layered Clean Architecture** pattern with clear separation of concerns between HTTP transport, domain business logic, data access, and third-party API integrations:

```text
Request
  │
  ▼
[ Middleware Layer ] (CORS, JWT Auth Token Verification, Role Check, Enrollment Guard, Rate Limiting)
  │
  ▼
[ Controller Layer ] (HTTP Parameter Parsing, DTO Validation)
  │
  ▼
[ Service Layer ] (Business Logic, Enrollment Rules, AI Prompting, Docker Execution, Excel Generation)
  │
  ├─► [ Repository Layer ] (GORM Database Abstraction ──► MySQL)
  ├─► [ AI Engine ]       (Gemini API Test Case Generator)
  ├─► [ Execution Engine ] (Docker gcc:12 Sandbox)
  └─► [ Report Engine ]    (excelize/v2 Streaming XLSX)
```

---

## Enrollment System

This section is the **single source of truth** for how students get access to courses. Any change to course, quiz, or submission endpoints MUST respect the rules below.

### 1. Core Principle

> **A student can only access a course (and its quizzes and submissions) if an admin has enrolled that student in that course.**

- Registering an account does **not** grant access to any course. A new `STUDENT` starts with zero courses.
- Students can **never** enroll themselves. There is no student-facing enroll/join endpoint.
- Only `ADMIN` can create or delete enrollments.
- Access control is enforced **in the backend** (middleware + service + query filtering). Hiding menus in the frontend is only UX, never security.

### 2. Roles & Permissions Matrix

| Action                                      |     `ADMIN`     |   `STUDENT` (enrolled)    | `STUDENT` (not enrolled) |
| :------------------------------------------ | :-------------: | :-----------------------: | :----------------------: |
| Create / edit / delete course               |       Yes       |            No             |            No            |
| Enroll / bulk-enroll / remove students      |       Yes       |            No             |            No            |
| List all students & enrollments of a course |       Yes       |            No             |            No            |
| List courses                                |   All courses   | Only own enrolled courses |        Empty list        |
| View a course detail                        |       Yes       |            Yes            |         **403**          |
| View a quiz in a course                     |       Yes       |            Yes            |         **403**          |
| Submit code to a quiz                       | No (not needed) |            Yes            |         **403**          |
| View own submissions in a course            |       n/a       |            Yes            |         **403**          |

### 3. Data Model

**`users`** (existing, extended): add `nim` (student ID number, `VARCHAR`, `UNIQUE`, nullable for `ADMIN`). `nim`, `name`, and `email` are searchable by admins when picking students.

**`enrollments`** (new): the join table between `courses` and `users` (students).

| Column        | Type            | Constraint                                        | Purpose                                               |
| :------------ | :-------------- | :------------------------------------------------ | :---------------------------------------------------- |
| `id`          | BIGINT UNSIGNED | PK, auto increment                                | Enrollment ID                                         |
| `course_id`   | BIGINT UNSIGNED | FK → `courses.id` **ON DELETE CASCADE**, NOT NULL | Which course                                          |
| `student_id`  | BIGINT UNSIGNED | FK → `users.id` **ON DELETE CASCADE**, NOT NULL   | Which student                                         |
| `enrolled_by` | BIGINT UNSIGNED | FK → `users.id`, NOT NULL                         | Admin who performed the enrollment (audit)            |
| `enrolled_at` | DATETIME        | NOT NULL, default now                             | When enrolled                                         |
| —             | —               | **UNIQUE (`course_id`, `student_id`)**            | A student cannot be enrolled twice in the same course |
| —             | —               | INDEX (`student_id`)                              | Fast "my courses" lookup                              |

Relationships:

```text
users (STUDENT) 1 ──── N enrollments N ──── 1 courses 1 ──── N quizzes 1 ──── N testcases
                                                                    │
users (STUDENT) 1 ───────────────────────── N submissions N ────────┘
```

Rules:

- Deleting a **course** cascades to its enrollments and quizzes. Deleting a **user** cascades to their enrollments.
- Removing an enrollment (see `DELETE` below) is done inside **one DB transaction**: delete the enrollment row, then delete that student's submissions for quizzes of that course. No orphaned submissions may remain.
- The UNIQUE constraint is the final guard against duplicates; the service layer must still check first so it can return a clean `409`/skip result instead of a raw DB error.

### 4. API Endpoints (CRUD)

All routes are under `/api/v1` and require a valid access token.

#### Admin (middleware: `Auth` + `RequireRole("ADMIN")`)

| CRUD           | Method & Path                                        | Description                                                                                      |
| :------------- | :--------------------------------------------------- | :----------------------------------------------------------------------------------------------- |
| **Courses**    | `POST/PUT/DELETE /admin/courses[/:id]`               | Create / update / delete a course (quizzes cascade on delete).                                   |
| **Quizzes**    | `POST/PUT/DELETE /admin/quizzes[/:id]`               | Create / update / delete a quiz (test cases & submissions cascade).                              |
| **Test cases** | `POST /admin/quizzes/:id/testcases`                  | Add a test case. Body: `{ "stdin", "expected_output", "is_hidden" }`.                            |
| **Test cases** | `DELETE /admin/testcases/:id`                        | Remove a single test case.                                                                       |
| **Generate**   | `POST /admin/quizzes/generate-testcases`             | Body: `{ "quiz_id", "problem_statement?", "count?" }`. AI-generates test cases via Gemini.       |
| **Enroll**     | `POST /admin/courses/:id/enroll-bulk`                | Body: `{ "student_ids": [12, 13] }`. Upserts each enrollment to `APPROVED`.                     |
| **List students** | `GET /admin/students?search=&page=&limit=`         | Search students by NIM / name / email (pagination). Used by the "Add students" dialog.           |
| **Enrollments** | `GET /admin/courses/:id/enrollments`                | Enrolled students of the course.                                                                 |
| **Remove**     | `DELETE /admin/courses/:id/enrollments/:student_id`  | Removes the student from the course.                                                             |
| **Progress**   | `GET /admin/courses/:id/progress`                    | Per-student completion % and average score.                                                      |
| **Export**     | `GET /admin/courses/:id/export-grades`               | Excel report (`.xlsx`) of enrolled students.                                                     |

> There is intentionally **no Update** for enrollments. An enrollment is just a link; to "move" a student, delete it and create a new one.

#### Student (middleware: `Auth`; enrollment enforced per route as described)

| Method & Path               | Guard                   | Description                                                                                                 |
| :-------------------------- | :---------------------- | :---------------------------------------------------------------------------------------------------------- |
| `GET /courses`              | Query filter            | Returns **only** courses the student is enrolled in (`APPROVED`, via SQL `JOIN enrollments`); admin sees all. |
| `GET /courses/me`           | Query filter            | Returns the user's `APPROVED` enrolled courses.                                                             |
| `GET /courses/:id`          | `RequireEnrollment`     | Course detail and its quiz list. `403` if not enrolled.                                                      |
| `GET /quizzes/:id`          | `RequireQuizEnrollment` | Resolve `quiz.course_id`, then verify enrollment. `403` if not enrolled.                                    |
| `POST /submissions`         | Middleware              | Resolve `quiz_id → course_id`, then verify enrollment **before** executing any code. `403` if not enrolled. |
| `GET /submissions?quiz_id=` | Owned data              | Own submissions only (admin may view any).                                                                  |
| `GET /student/enrollments`  | Owned data              | The user's own enrollment rows.                                                                             |

### 5. How Access Is Enforced (Implementation Guide)

Use **three layers**, so one forgotten check does not become a data leak:

1. **Role check (middleware):** `RequireRole("ADMIN")` on every `/admin/*` route. `ADMIN` is never blocked by enrollment checks.
2. **Enrollment guard (middleware + service):**
   - Service method (single source of truth):
     ```go
     // internal/services/enrollment_service.go
     // Returns true when allowed. ADMIN always allowed; STUDENT needs an APPROVED enrollment row.
     func (s *EnrollmentService) CheckAccess(ctx context.Context, userID, courseID uint) (bool, error)
     ```
   - Repository: `FindByUserAndCourse(ctx, userID, courseID uint) (*Enrollment, error)` looks up the enrollment row; `CheckAccess` returns `status == APPROVED`.
   - Middleware `RequireEnrollment` reads `:id` (course id) from the route and aborts with `403` on failure. `RequireQuizEnrollment` and `RequireEnrollmentForSubmission` first load the quiz's `course_id` — never trust a `course_id` sent by the client.
3. **Query-level filtering:** list endpoints for students must never fetch all rows and filter in Go. Use a `JOIN`/`WHERE` on `enrollments.student_id = ?` so non-enrolled data never leaves the database.

Request flow for a student opening a quiz and submitting code:

```text
Student ──► POST /submissions { quiz_id, source_code }
   │
   ▼  Auth middleware: valid JWT → user_id, role=STUDENT
   ▼  Controller: validate DTO
   ▼  Submission service:
        1. load quiz → course_id
        2. CheckAccess(user, course_id)  ── not enrolled ──► 403 (stop, nothing executed)
        3. fetch testcases → run in Docker gcc:12 → tolerance matcher
        4. save submission (PASSED / WRONG_ANSWER / COMPILE_ERROR)
```

Frontend behavior (UX only):

- Student dashboard calls `GET /courses` and shows only the returned courses. An empty list shows: "You are not enrolled in any course yet. Please contact your admin."
- If any request returns `403`, the Axios interceptor redirects to the student dashboard with a toast message. Do not retry or refresh the token on `403` (only on `401`).
- Admin enrollment page: enrolled-students table (search + pagination + remove button with confirmation) and an "Add students" dialog backed by `GET /admin/students` with multi-select, which submits to `enroll-bulk`.

### 6. Enrollment Lifecycle (End to End)

1. Student registers via `POST /auth/register` (role forced to `STUDENT`, `nim` required) → **no course access yet**.
2. Admin creates a course and its quizzes.
3. Admin opens the course → "Add students" → searches by NIM/name/email → selects students → `POST /admin/courses/:id/enroll-bulk`.
4. Student logs in → `GET /courses` now returns the enrolled course → opens quizzes → submits code.
5. If the admin removes the student (`DELETE .../enrollments/:student_id`), the student immediately loses access (the next request returns `403`; the course disappears from `GET /courses`). Past submissions remain stored but are no longer reachable through the course.

### 7. Enrollment Rules (Must Not Break)

1. Never expose a student-facing endpoint that creates an enrollment.
2. Every student-facing route that touches a course, quiz, testcase, or submission MUST call `CheckAccess` (or be filtered by `enrollments`). New endpoints must be added to the permission matrix above.
3. Only users with role `STUDENT` can be enrolled. Reject enrolling `ADMIN` accounts.
4. Bulk enrollment must be idempotent and transactional.
5. `enrollments` has `UNIQUE (course_id, student_id)` and cascading foreign keys. Keep migrations and GORM model tags in sync.
6. Do not rely on `course_id` provided by the client for authorization; derive it from the quiz or the route resource.
7. Never return student-private data (emails, scores) to a non-admin outside their own records.

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
3. Backend resolves the quiz's course and **verifies the student is enrolled** (`CheckAccess`). If not, respond `403` and stop.
4. Backend fetches problem testcases from database.
5. Backend executes `source_code` + `stdin` in **Docker container (gcc:12)**.
6. Docker returns execution `stdout` / `stderr`.
7. **Tolerance Matcher Engine** compares student `stdout` against `expected_output` using sequence & keyword matching logic.
8. Evaluation score & status (`PASSED`, `WRONG_ANSWER`, `COMPILE_ERROR`) are recorded in MySQL.

### 3. Admin Dashboard & Reporting Flow

1. Admin manages Courses, Quizzes, and Student Enrollments via Admin Dashboard (enrollment details: see [Enrollment System](#enrollment-system)).
2. **Bulk Enrollment**: Admin searches students (by NIM/Name/Email) from the candidates list, selects multiple, and enrolls them to a course in one action (`POST /admin/courses/:id/enroll-bulk`).
3. **Progress Tracking**: Admin toggles progress view per course to see student completion % and average scores.
4. **Export Grades/Progress**: Admin downloads Excel reports (`GET /admin/courses/:id/export-grades`) generated via `excelize/v2` streaming to HTTP response.
5. **Student Removal**: Admin removes student from course (`DELETE /admin/courses/:id/enrollments/:student_id`); the enrollment and the student's submissions for that course are deleted in one transaction.

### 4. Database Seeding (Development)

Run `go run ./cmd/seed/main.go` to populate DB with:

- 1 Admin user (`admin@cgrader.local` / `admin123`)
- 5 Student users (`password123`), each with a `nim`
- 2 Courses with 3 Quizzes & Testcases
- Enrollments (not every student is enrolled in every course, so access control can be tested) & Random Submissions (PASSED/WRONG_ANSWER) only for enrolled student/course pairs

---

## Development Commands

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
  go build -o build/server cmd/server/main.go
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
```

The provided `.env.example` already contains DB credentials matching `docker-compose.yml`, plus the optional
`CORS_ALLOWED_ORIGINS`, `COOKIE_SECURE`, and `GEMINI_MODEL`. `JWT_SECRET` and `JWT_REFRESH_SECRET` must be
different values (e.g. `openssl rand -hex 32`). The server refuses to start otherwise.

Never commit `.env` or real API keys/secrets to the repository.

### 4. Start Local Infrastructure

The recommended local setup uses Docker Compose for MySQL:

```bash
docker compose up -d
```

Check that the containers are running:

```bash
docker compose ps
```

The `migrations/` directory is mounted into the MySQL container (`/docker-entrypoint-initdb.d`), so
`migrations/001_init.sql` (which includes the `enrollments` table and the `users.nim` column) is applied automatically **only the first time** the MySQL volume is created.

To apply it manually (e.g. against an existing database; the file is idempotent):

```bash
docker compose exec -T mysql sh -c 'mysql -uroot -p"$MYSQL_ROOT_PASSWORD" autograder_db' < migrations/001_init.sql
```

To start from scratch: `docker compose down -v && docker compose up -d` (this deletes all data).

> **Code Execution:** This project uses Docker directly (`gcc:12` image) for sandboxed C code execution. No separate Judge0 or Piston container is required. Ensure Docker Desktop (WSL2 backend on Windows) is running.

### 5. Install Backend Dependencies

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

The backend and frontend should now run independently (default: API on `http://localhost:1111`, web on `http://localhost:5173`).

Public registration only creates `STUDENT` accounts (without any course access). Create an `ADMIN` account with:

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
- Enrollment access control works: an enrolled student can open the course, a non-enrolled student gets `403`, and `GET /courses` only lists enrolled courses.

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

| Variable Name          | Purpose                                                         | Required                              | Example Placeholder                                                                                     |
| :--------------------- | :-------------------------------------------------------------- | :------------------------------------ | :------------------------------------------------------------------------------------------------------ |
| `PORT`                 | Server listening port                                           | No (Default: 1111)                    | `1111`                                                                                                  |
| `DB_DSN`               | MySQL connection string DSN                                     | Yes                                   | `autograder:autograder_pass@tcp(127.0.0.1:3306)/autograder_db?parseTime=true&charset=utf8mb4&loc=Local` |
| `JWT_SECRET`           | Secret key for signing Access Tokens                            | Yes                                   | `<your-jwt-access-secret>`                                                                              |
| `JWT_REFRESH_SECRET`   | Secret key for Refresh Tokens                                   | Yes                                   | `<your-jwt-refresh-secret>`                                                                             |
| `GEMINI_API_KEY`       | Google Gemini AI API Key                                        | Yes                                   | `<your-gemini-api-key>`                                                                                 |
| `EXEC_DOCKER_IMAGE`    | Docker image used to compile/run submitted C code               | No (Default: `gcc:12`)                | `gcc:12`                                                                                                |
| `EXEC_TIMEOUT_SECONDS` | Maximum execution time per testcase                             | No (Default: `5`)                     | `5`                                                                                                     |
| `EXEC_MEMORY_LIMIT_MB` | Memory limit for each sandbox container                         | No (Default: `128`)                   | `128`                                                                                                   |
| `EXEC_CPU_LIMIT`       | CPU quota for each sandbox container                            | No (Default: `0.5`)                   | `0.5`                                                                                                   |
| `CORS_ALLOWED_ORIGINS` | Comma-separated trusted frontend origins (never `*`)            | No (Default: `http://localhost:5173`) | `http://localhost:5173`                                                                                 |
| `COOKIE_SECURE`        | Mark refresh-token cookie as `Secure` (set `true` behind HTTPS) | No (Default: `false`)                 | `false`                                                                                                 |
| `GEMINI_MODEL`         | Gemini model used for testcase generation                       | No (Default: `gemini-2.5-flash`)      | `gemini-2.5-flash`                                                                                      |

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
- **Error Format:**
  ```json
  {
    "success": false,
    "message": "Error details or validation message",
    "errors": []
  }
  ```
- **Authentication Header:** `Authorization: Bearer <JWT_ACCESS_TOKEN>`
- **Status Codes:** `200` OK, `201` Created, `400` validation error, `401` missing/invalid/expired token (frontend refreshes the token), `403` authenticated but not allowed (wrong role or **not enrolled**; frontend must NOT refresh the token), `404` resource not found, `409` conflict (e.g. already enrolled).

---

## Database

- **Database Engine:** MySQL 8.0+
- **ORM:** GORM (`gorm.io/gorm`)
- **Core Entities:**
  - `User`: Stores credentials, `nim` (students), roles (`STUDENT`, `ADMIN`), and timestamps.
  - `Course`: Learning modules containing multiple quizzes.
  - `Enrollment`: Link between a `STUDENT` and a `Course`, created only by an `ADMIN`. Unique per (`course_id`, `student_id`). Defines which courses a student may access.
  - `Quiz`: Problem statements and constraints created by admins.
  - `TestCase`: `stdin` and `expected_output` generated by AI or defined manually.
  - `Submission`: Student code submissions, execution status, and scores. Only exists for enrolled student/course pairs.

---

## Authentication & Authorization

- **Passwords:** Hashed using `bcrypt` (default cost: 10 or higher).
- **Token Strategy:**
  - **Access Token:** Short-lived JWT passed via `Authorization` HTTP header.
  - **Refresh Token:** Stored in secure `HttpOnly` Cookie.
- **Role-Based Access Control (RBAC):**
  - `STUDENT`: Can view **only the courses they are enrolled in**, solve those quizzes, and submit code.
  - `ADMIN`: Access to Admin Dashboard, manage courses and enrollments, create quizzes, trigger AI testcase generation.
- **Resource-level authorization:** Role alone is not enough for students. Course-scoped access additionally requires an enrollment (see [Enrollment System](#enrollment-system)). Do not store enrollment in the JWT; check the database so removal takes effect immediately.

---

## Validation & Error Handling

- **Request Payload Validation:** Handled via `go-playground/validator/v10` on Go struct DTO tags (`binding:"required,email"`).
- **Global Error Middleware:** Captures unhandled panics and formats standard JSON error responses.

---

## Testing

- **Backend:** Native Go testing framework (`testing` package). Executable via `go test ./...`.
  - Enrollment tests that MUST exist: enroll (single/bulk/duplicate skipped), remove (past submissions retained), `CheckAccess` for admin / enrolled student / non-enrolled student, and `GET /courses` returning only enrolled courses.
- **Frontend:** ESLint checks via `npm run lint`.

---

## Security

- Cross-Origin Resource Sharing (CORS) restricted to trusted frontend domain.
- Password hashing with bcrypt before database persistence.
- SQL Injection protection enforced via GORM parameterized queries.
- Execution sandboxing offloaded to Docker containers (`gcc:12`) to prevent remote code execution (RCE) on backend server host.
- Authorization (including enrollment checks) is enforced server-side on every request; never trust client-supplied `course_id` or role.

---

## Important Invariants

1. **Sandboxed Code Execution:** NEVER compile or run untrusted user C code directly on the host machine running the Gin server. All executions MUST go through Docker containers (`gcc:12` image).
2. **AI Test Case Output:** AI Testcase responses MUST be validated as proper JSON arrays matching the strict schema before storing in MySQL.
3. **Database Relationships:** Quiz records MUST be deleted with cascading deletes on associated `TestCases` and `Submissions`. Course deletion MUST cascade to `Enrollments`.
4. **JWT Security:** Access token verification MUST be enforced on all routes except `/auth/login`, `/auth/register`, and `/auth/refresh`.
5. **Admin-Controlled Enrollment:** Only `ADMIN` can create or delete enrollments. A `STUDENT` MUST NOT be able to read, submit to, or even fetch details of a course they are not enrolled in. Authorization is checked **before** any code execution or data fetching.

---

## Common Pitfalls

- **Missing CORS Headers:** Forgetting to configure credentials support in `gin-contrib/cors` when using HTTP-Only cookies for refresh tokens.
- **Direct Database Queries:** Bypassing repository interfaces in controller handlers.
- **Parsing Double Escaped Sequences:** Incorrectly parsing newline sequences (`\n`) in `stdin` strings received from Gemini API.
- **Forgotten Enrollment Check:** Adding a new student-facing endpoint (e.g. a leaderboard or testcase preview) that takes a `quiz_id` or `course_id` without calling `CheckAccess`. This is an access-control leak (IDOR).
- **Frontend-Only Protection:** Hiding a course in the UI but leaving the API open. Always enforce on the backend.
- **Treating `403` like `401`:** The Axios interceptor must refresh tokens only on `401`; refreshing on `403` causes infinite retry loops.
- **Orphaned Data on Removal:** Deleting an enrollment without removing that student's submissions for the course (do it in one transaction).

---

## Rules for AI Agents

### Before Modifying Code

- Inspect related files in `internal/controllers/`, `internal/services/`, and `web/src/`.
- Ensure new endpoints follow existing `/api/v1` path conventions.
- Verify environment variable requirements in `internal/config/`.
- For anything touching courses, quizzes, or submissions, re-read the [Enrollment System](#enrollment-system) section first.

### While Modifying Code

- Maintain separation of concerns: Do not place database logic in controllers.
- Use GORM parameterized queries to prevent SQL injection.
- Keep frontend components modular and follow Tailwind CSS styling patterns.
- Ensure all Go code handles errors explicitly without ignoring return values.
- Every new student-facing, course-scoped endpoint MUST enforce enrollment through `CheckAccess` (or an `enrollments` JOIN for list queries) and be added to the permission matrix.

### After Modifying Code

- Run `go test ./...` to verify backend integrity.
- Run `cd web && npm run lint` to ensure frontend syntax compliance.
- Confirm schema consistency across GORM models and MySQL migration files.
- Update the `Project Setup` section when implementation changes affect dependencies, environment variables, database setup, services, or development commands.
- Ensure the setup documentation describes the final implemented project, not only the original project structure.
- Update the [Enrollment System](#enrollment-system) section if endpoints, tables, or rules change.

---

## Verification Checklist

- [ ] All database migration files match GORM models (including `enrollments` and `users.nim`).
- [ ] No hardcoded secrets or API keys exist in source code.
- [ ] JWT authentication middleware wraps protected endpoints.
- [ ] Docker execution handles network timeout and connection failure gracefully.
- [ ] Frontend Axios interceptor correctly handles token refresh on 401 response status (and not on 403).
- [ ] Only `ADMIN` can call enrollment create/delete endpoints.
- [ ] A non-enrolled student receives `403` on course detail, quiz detail, and submission endpoints.
- [ ] `GET /courses` for a student returns only enrolled courses.
- [ ] Bulk enrollment is idempotent (re-running it does not create duplicates or errors).
- [ ] Removing an enrollment also removes that student's submissions for the course.
