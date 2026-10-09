package repositories

import (
	"autograder/internal/models"

	"gorm.io/gorm"
)

type Repo struct{ db *gorm.DB }

func New(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) CreateUser(u *models.User) error { return r.db.Create(u).Error }

func (r *Repo) UserByEmail(email string) (*models.User, error) {
	var u models.User
	return &u, r.db.Where("email = ?", email).First(&u).Error
}

func (r *Repo) UserByID(id uint) (*models.User, error) {
	var u models.User
	return &u, r.db.First(&u, id).Error
}

func (r *Repo) ListCourses() (cs []models.Course, err error) {
	err = r.db.Order("id").Find(&cs).Error
	return
}

func (r *Repo) CreateCourse(c *models.Course) error { return r.db.Create(c).Error }

func (r *Repo) ListQuizzes() (qs []models.Quiz, err error) {
	err = r.db.Order("id").Find(&qs).Error
	return
}

func (r *Repo) QuizByID(id uint) (*models.Quiz, error) {
	var q models.Quiz
	return &q, r.db.Preload("TestCases").First(&q, id).Error
}

func (r *Repo) CreateQuiz(q *models.Quiz) error { return r.db.Create(q).Error }

// DeleteQuiz removes the quiz together with its test cases and submissions.
func (r *Repo) DeleteQuiz(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("quiz_id = ?", id).Delete(&models.TestCase{}).Error; err != nil {
			return err
		}
		if err := tx.Where("quiz_id = ?", id).Delete(&models.Submission{}).Error; err != nil {
			return err
		}
		res := tx.Delete(&models.Quiz{}, id)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func (r *Repo) CreateTestCases(tcs []models.TestCase) error { return r.db.Create(&tcs).Error }

func (r *Repo) CreateSubmission(s *models.Submission) error { return r.db.Create(s).Error }

func (r *Repo) SubmissionsByUser(uid uint) (ss []models.Submission, err error) {
	err = r.db.Where("user_id = ?", uid).Order("id desc").Limit(50).Find(&ss).Error
	return
}
