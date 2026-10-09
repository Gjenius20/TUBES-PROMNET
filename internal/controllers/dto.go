package controllers

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email,max=191"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=191"`
	Password string `json:"password" binding:"required,max=72"`
}

type CourseRequest struct {
	Title       string `json:"title" binding:"required,min=2,max=150"`
	Description string `json:"description" binding:"max=5000"`
}

type QuizCreateRequest struct {
	CourseID    uint   `json:"course_id" binding:"required"`
	Title       string `json:"title" binding:"required,min=2,max=150"`
	Description string `json:"description" binding:"required,max=10000"`
	Constraints string `json:"constraints" binding:"max=5000"`
}

type QuizUpdateRequest struct {
	Title       string `json:"title" binding:"required,min=2,max=150"`
	Description string `json:"description" binding:"required,max=10000"`
	Constraints string `json:"constraints" binding:"max=5000"`
}

type TestCaseRequest struct {
	Stdin          string `json:"stdin" binding:"max=10000"`
	ExpectedOutput string `json:"expected_output" binding:"required,max=10000"`
	IsHidden       bool   `json:"is_hidden"`
}

type GenerateTestCasesRequest struct {
	QuizID           uint   `json:"quiz_id" binding:"required"`
	ProblemStatement string `json:"problem_statement" binding:"max=10000"`
	Count            int    `json:"count" binding:"omitempty,min=1,max=10"`
}

type SubmissionRequest struct {
	QuizID     uint   `json:"quiz_id" binding:"required"`
	SourceCode string `json:"source_code" binding:"required,max=50000"`
}

// Response khusus siswa: tidak membocorkan testcase hidden.
type SampleTestCase struct {
	Stdin          string `json:"stdin"`
	ExpectedOutput string `json:"expected_output"`
}

type StudentQuizResponse struct {
	ID              uint             `json:"id"`
	CourseID        uint             `json:"course_id"`
	Title           string           `json:"title"`
	Description     string           `json:"description"`
	Constraints     string           `json:"constraints"`
	SampleTestCases []SampleTestCase `json:"sample_test_cases"`
}
