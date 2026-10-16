package models

import "time"

type Role string

const (
	RoleStudent Role = "STUDENT"
	RoleAdmin   Role = "ADMIN"
)

type SubmissionStatus string

const (
	StatusPassed       SubmissionStatus = "PASSED"
	StatusWrongAnswer  SubmissionStatus = "WRONG_ANSWER"
	StatusCompileError SubmissionStatus = "COMPILE_ERROR"
	StatusRuntimeError SubmissionStatus = "RUNTIME_ERROR"
)

const (
	TestCaseSourceAI     = "AI"
	TestCaseSourceManual = "MANUAL"
)

type EnrollmentStatus string

const (
	EnrollmentPending  EnrollmentStatus = "PENDING"
	EnrollmentApproved EnrollmentStatus = "APPROVED"
	EnrollmentRejected EnrollmentStatus = "REJECTED"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:100;not null" json:"name"`
	Email        string    `gorm:"size:191;uniqueIndex;not null" json:"email"`
	NIM          string    `gorm:"size:50;uniqueIndex;not null" json:"nim"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Role         Role      `gorm:"size:20;not null;default:STUDENT" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Course struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:150;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Quizzes     []Quiz    `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"quizzes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Quiz struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	CourseID    uint       `gorm:"not null;index" json:"course_id"`
	Title       string     `gorm:"size:150;not null" json:"title"`
	Description string     `gorm:"type:text;not null" json:"description"`
	Constraints string     `gorm:"type:text" json:"constraints"`
	MaxScore    int        `gorm:"not null;default:100" json:"max_score"`
	TestCases   []TestCase `gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE" json:"test_cases,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type TestCase struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	QuizID         uint      `gorm:"not null;index" json:"quiz_id"`
	Stdin          string    `gorm:"type:text" json:"stdin"`
	ExpectedOutput string    `gorm:"type:text;not null" json:"expected_output"`
	IsHidden       bool      `gorm:"not null;default:true" json:"is_hidden"`
	Source         string    `gorm:"size:20;not null;default:MANUAL" json:"source"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Submission struct {
	ID           uint             `gorm:"primaryKey" json:"id"`
	UserID       uint             `gorm:"not null;index:idx_submissions_user_quiz,priority:1" json:"user_id"`
	QuizID       uint             `gorm:"not null;index:idx_submissions_user_quiz,priority:2" json:"quiz_id"`
	SourceCode   string           `gorm:"type:mediumtext;not null" json:"source_code"`
	Status       SubmissionStatus `gorm:"size:20;not null" json:"status"`
	Score        int              `gorm:"not null;default:0" json:"score"`
	PassedCount  int              `gorm:"not null;default:0" json:"passed_count"`
	TotalCount   int              `gorm:"not null;default:0" json:"total_count"`
	ErrorMessage string           `gorm:"type:text" json:"error_message"`
	ResultsJson  string           `gorm:"type:json" json:"-"`
	CreatedAt    time.Time        `json:"created_at"`

	// Results di-load dari ResultsJson atau diisi langsung saat submit.
	Results []TestResult `gorm:"-" json:"results,omitempty"`
}

// Enrollment merepresentasikan pendaftaran siswa ke course.
type Enrollment struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	UserID    uint             `gorm:"not null;index:idx_enrollments_user_course,unique,priority:1" json:"user_id"`
	CourseID  uint             `gorm:"not null;index:idx_enrollments_user_course,unique,priority:2" json:"course_id"`
	Status    EnrollmentStatus `gorm:"size:20;not null;default:PENDING" json:"status"`
	EnrolledBy *uint           `gorm:"null" json:"enrolled_by"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`

	// Relations
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Course Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

// TestResult adalah hasil satu testcase (bukan tabel database).
// Untuk testcase hidden, input & output yang diharapkan sengaja tidak diisi.
type TestResult struct {
	Index          int    `json:"index"`
	Hidden         bool   `json:"hidden"`
	Passed         bool   `json:"passed"`
	Status         string `json:"status"`
	Stdin          string `json:"stdin,omitempty"`
	ExpectedOutput string `json:"expected_output,omitempty"`
	ActualOutput   string `json:"actual_output,omitempty"`
	Message        string `json:"message,omitempty"`
}
