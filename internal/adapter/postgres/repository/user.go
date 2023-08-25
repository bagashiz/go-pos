package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/bagashiz/go-pos/internal/adapter/postgres"
	"github.com/bagashiz/go-pos/internal/core/domain"
)

var (
	/**
	 * psql holds a reference to squirrel.StatementBuilderType
	 * which is used to build SQL queries that compatible with PostgreSQL syntax
	 */
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

/**
 * UserRepository implements port.UserRepository interface
 * and provides an access to the postgres database
 */
type UserRepository struct {
	db *postgres.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *postgres.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// CheckUserExists checks if a user exists in the database using the email
func (ur *UserRepository) CheckUserExists(ctx context.Context, email string) (bool, error) {
	query := psql.Select("COUNT(*)").
		From("users").
		Where(sq.Eq{"email": email}).
		Limit(1).
		RunWith(ur.db)

	var count int

	err := query.QueryRowContext(ctx).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CreateUser creates a new user in the database
func (ur *UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := psql.Insert("users").
		Columns("name", "email", "password").
		Values(user.Name, user.Email, user.Password).
		Suffix("RETURNING *").
		RunWith(ur.db)

	err := query.QueryRowContext(ctx).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID gets a user by ID from the database
func (ur *UserRepository) GetUserByID(ctx context.Context, id uint64) (*domain.User, error) {
	query := psql.Select("*").
		From("users").
		Where(sq.Eq{"id": id}).
		Limit(1).
		RunWith(ur.db)

	var user domain.User

	err := query.QueryRowContext(ctx).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByEmailAndPassword gets a user by email from the database
func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := psql.Select("*").
		From("users").
		Where(sq.Eq{"email": email}).
		Limit(1).
		RunWith(ur.db)

	var user domain.User

	err := query.QueryRowContext(ctx).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ListUsers lists all users from the database
func (ur *UserRepository) ListUsers(ctx context.Context, pageId, pageSize uint64) ([]*domain.User, error) {
	query := psql.Select("*").
		From("users").
		OrderBy("id").
		Limit(pageSize).
		Offset((pageId - 1) * pageSize).
		RunWith(ur.db)

	var users []*domain.User

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.User

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

// UpdateUser updates a user by ID in the database
func (ur *UserRepository) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := psql.Update("users").
		Set("name", user.Name).
		Set("email", user.Email).
		Set("password", user.Password).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": user.ID}).
		Suffix("RETURNING *").
		RunWith(ur.db)

	err := query.QueryRowContext(ctx).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user by ID from the database
func (ur *UserRepository) DeleteUser(ctx context.Context, id uint64) error {
	query := psql.Delete("users").
		Where(sq.Eq{"id": id}).
		RunWith(ur.db)

	_, err := query.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}
