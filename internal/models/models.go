package models

import "time"

const (
	RoleStudent = "STUDENT"
	RoleAdmin   = "ADMIN"

	StatusPassed       = "PASSED"
	StatusWrongAnswer  = "WRONG_ANSWER"
	StatusCompileError = "COMPILE_ERROR"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Email     string    `gorm:"size:191;uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Role      string    `gorm:"size:20;not null;default:STUDENT" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Course struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:191;not null" json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Quiz struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	CourseID    *uint        `json:"course_id"`
	Title       string       `gorm:"size:191;not null" json:"title"`
	Statement   string       `gorm:"type:text;not null" json:"statement"`
	TestCases   []TestCase   `gorm:"constraint:OnDelete:CASCADE" json:"test_cases,omitempty"`
	Submissions []Submission `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type TestCase struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	QuizID         uint   `gorm:"index;not null" json:"quiz_id"`
	Stdin          string `gorm:"type:text" json:"stdin"`
	ExpectedOutput string `gorm:"type:text;not null" json:"expected_output"`
}

type Submission struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	QuizID     uint      `gorm:"index;not null" json:"quiz_id"`
	SourceCode string    `gorm:"type:mediumtext;not null" json:"source_code"`
	Status     string    `gorm:"size:20;not null" json:"status"`
	Score      float64   `json:"score"`
	CreatedAt  time.Time `json:"created_at"`
}
