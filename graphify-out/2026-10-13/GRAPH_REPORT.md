# Autograder Platform - Knowledge Graph Report

## Overview
AI-Powered Code Assessment & Autograder Platform untuk pembelajaran C programming dengan autograding otomatis.

## Architecture Layers

### 1. Presentation Layer
- **Frontend:** ReactJS (Vite) + Monaco Editor @ `/web`
- **Components:** `web/src/components/`, `web/src/pages/`
- **State Management:** Zustand
- **API Client:** Axios with interceptors

### 2. Controller Layer
- **Location:** `/internal/controllers/`
- **Controllers:**
  - `AuthController` - Authentication endpoints
  - `CourseController` - Course management
  - `QuizController` - Quiz operations
  - `SubmissionController` - Code submission

### 3. Service Layer
- **Location:** `/internal/services/`
- **Core Services:**
  - `AuthService` - JWT management
  - `CourseService` - Course business logic
  - `QuizService` - Quiz + AI test case generation
  - `SubmissionService` - Code execution & grading
- **External Integrations:**
  - `ai/Gemini` - AI test case generation
  - `docker/Client` - Code execution sandbox
  - `autograder/Matcher` - Tolerance matching algorithm

### 4. Repository Layer
- **Location:** `/internal/repositories/`
- **Models:** User, Course, Quiz, TestCase, Submission
- **ORM:** GORM with MySQL

### 5. Config & Utils
- **Config:** Environment loading, DB connection
- **Utils:** JWT, password hashing, validation

## Key Relationships

```
User --[has many]--> Submission
Quiz --[has many]--> TestCase
Quiz --[belongs to]--> Course
Submission --[belongs to]--> User, Quiz
```

## External Dependencies
- **Database:** MySQL 8.0+
- **AI:** Google Gemini API
- **Execution:** Docker (gcc:12)
- **Frontend:** Monaco Editor, Tailwind CSS

## Data Flow
1. Admin creates quiz → AI generates test cases → Stored in DB
2. Student submits code → Executed in Docker → Compared with test cases
3. Results saved with tolerance matching algorithm

## Communities Detected
1. **Authentication Community** - Auth controller/service/repo
2. **Course Management** - Course/Quiz related
3. **Autograding Core** - Submission + Matcher + Execution
4. **AI Integration** - Gemini parser/generator
5. **Infrastructure** - Config, Utils, Middleware

## Recommendations
1. Consolidate execution engine (remove judge0/piston if unused)
2. Add migration management for schema changes
3. Implement caching layer for quiz/test cases
4. Add comprehensive test coverage
5. Consider API versioning
