package repositories

import (
	"context"

	"autograder/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	ListStudentsPaginated(ctx context.Context, search string, page, limit int) ([]models.User, int64, error)
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByNIM(ctx context.Context, nim string) (*models.User, error)
	FindByIdentifier(ctx context.Context, identifier string) (*models.User, error)
	FindByID(ctx context.Context, id uint) (*models.User, error)
}

type userRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) UserRepository { return &userRepository{db: db} }

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return translate(r.db.WithContext(ctx).Create(user).Error)
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, translate(err)
	}
	return &user, nil
}

func (r *userRepository) FindByIdentifier(ctx context.Context, identifier string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ? OR nim = ?", identifier, identifier).First(&user).Error; err != nil {
		return nil, translate(err)
	}
	return &user, nil
}

func (r *userRepository) FindByNIM(ctx context.Context, nim string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("nim = ?", nim).First(&user).Error; err != nil {
		return nil, translate(err)
	}
	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, translate(err)
	}
	return &user, nil
}

func (r *userRepository) ListStudentsPaginated(ctx context.Context, search string, page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	q := r.db.WithContext(ctx).Model(&models.User{}).Where("role = ?", models.RoleStudent)

	if search != "" {
		like := "%" + search + "%"
		q = q.Where("name LIKE ? OR email LIKE ? OR nim LIKE ?", like, like, like)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, translate(err)
	}

	offset := (page - 1) * limit
	if err := q.Offset(offset).Limit(limit).Order("id ASC").Find(&users).Error; err != nil {
		return nil, 0, translate(err)
	}

	return users, total, nil
}
