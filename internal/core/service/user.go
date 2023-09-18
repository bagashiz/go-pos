package service

import (
	"context"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/bagashiz/go-pos/internal/core/util"
)

/**
 * UserService implements port.UserService interface
 * and provides an access to the user repository
 */
type UserService struct {
	repo port.UserRepository
}

// NewUserService creates a new user service instance
func NewUserService(repo port.UserRepository) *UserService {
	return &UserService{
		repo,
	}
}

// Register creates a new user
func (us *UserService) Register(ctx context.Context, user *domain.User) (*domain.User, error) {
	exists, err := us.repo.CheckUserExists(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, port.ErrConflictingData
	}

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword

	return us.repo.CreateUser(ctx, user)
}

// GetUser gets a user by ID
func (us *UserService) GetUser(ctx context.Context, id uint64) (*domain.User, error) {
	return us.repo.GetUserByID(ctx, id)
}

// ListUsers lists all users
func (us *UserService) ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error) {
	return us.repo.ListUsers(ctx, skip, limit)
}

// UpdateUser updates a user's name, email, and password
func (us *UserService) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	existingUser, err := us.repo.GetUserByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	emptyData := user.Name == "" && user.Email == "" && user.Password == ""
	sameData := existingUser.Name == user.Name && existingUser.Email == user.Email
	if emptyData || sameData {
		return nil, port.ErrNoUpdatedData
	}

	var hashedPassword string

	if user.Password != "" {
		hashedPassword, err = util.HashPassword(user.Password)
		if err != nil {
			return nil, err
		}
	}

	user.Password = hashedPassword

	_, err = us.repo.UpdateUser(ctx, user)
	if err != nil {
		if port.IsUniqueConstraintViolationError(err) {
			return nil, port.ErrConflictingData
		}
	}

	return user, nil
}

// DeleteUser deletes a user by ID
func (us *UserService) DeleteUser(ctx context.Context, id uint64) error {
	_, err := us.repo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	return us.repo.DeleteUser(ctx, id)
}
