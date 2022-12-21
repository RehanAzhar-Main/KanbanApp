package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	var userData entity.User

	// get user by id
	if err := r.db.Table("users").
		Where("id = ?", id).
		Find(&userData).Error; err != nil {
		return entity.User{}, err
	}

	// user not found
	if userData == (entity.User{}) {
		return entity.User{}, nil
	}

	return userData, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var userData entity.User

	// get user by email
	if err := r.db.WithContext(ctx).Table("users").
		Where("email = ?", email).
		Find(&userData).Error; err != nil {
		return entity.User{}, err
	}

	// if userData.ID != 0 {
	// 	return entity.User{}, nil
	// }

	// user not found
	if userData == (entity.User{}) {
		return entity.User{}, nil
	}

	return userData, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	if err := r.db.WithContext(ctx).Table("users").
		Create(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	if err := r.db.WithContext(ctx).Table("users").Where("id = ?", user.ID).
		Updates(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	var user entity.User
	if err := r.db.WithContext(ctx).Delete(&user, id).Error; err != nil {
		return err
	}

	return nil
}
