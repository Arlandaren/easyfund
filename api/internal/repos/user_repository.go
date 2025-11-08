package repos

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/Arlandaren/easyfund/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetRandomUser(ctx context.Context) (*models.User, error)
	ListUsers(ctx context.Context, limit int) ([]models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	UpdatePasswordHash(ctx context.Context, userID uuid.UUID, hash string) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type userRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) CreateUser(ctx context.Context, user *models.User) error {
	const q = `
		INSERT INTO users (user_id, full_name, email, phone, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, q,
		user.UserID, user.FullName, user.Email, user.Phone, user.PasswordHash, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *userRepositoryImpl) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	const q = `
		SELECT user_id, full_name, email, phone, password_hash, created_at, updated_at
		FROM users WHERE user_id = $1
	`
	u := &models.User{}
	err := r.db.QueryRowContext(ctx, q, userID).Scan(
		&u.UserID, &u.FullName, &u.Email, &u.Phone, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	const q = `
		SELECT user_id, full_name, email, phone, password_hash, created_at, updated_at
		FROM users WHERE email = $1
	`
	u := &models.User{}
	err := r.db.QueryRowContext(ctx, q, email).Scan(
		&u.UserID, &u.FullName, &u.Email, &u.Phone, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepositoryImpl) GetRandomUser(ctx context.Context) (*models.User, error) {
	const q = `
		SELECT user_id, full_name, email, phone, password_hash, created_at, updated_at
		FROM users ORDER BY RANDOM() LIMIT 1
	`
	u := &models.User{}
	err := r.db.QueryRowContext(ctx, q).Scan(
		&u.UserID, &u.FullName, &u.Email, &u.Phone, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepositoryImpl) ListUsers(ctx context.Context, limit int) ([]models.User, error) {
	const q = `
		SELECT user_id, full_name, email, phone, password_hash, created_at, updated_at
		FROM users LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.UserID, &u.FullName, &u.Email, &u.Phone, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, rows.Err()
}

func (r *userRepositoryImpl) UpdateUser(ctx context.Context, user *models.User) error {
	const q = `
		UPDATE users SET full_name=$1, email=$2, phone=$3, updated_at=$4
		WHERE user_id=$5
	`
	_, err := r.db.ExecContext(ctx, q, user.FullName, user.Email, user.Phone, user.UpdatedAt, user.UserID)
	return err
}

func (r *userRepositoryImpl) UpdatePasswordHash(ctx context.Context, userID uuid.UUID, hash string) error {
	const q = `UPDATE users SET password_hash=$1, updated_at=NOW() WHERE user_id=$2`
	_, err := r.db.ExecContext(ctx, q, hash, userID)
	return err
}

func (r *userRepositoryImpl) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	const q = `DELETE FROM users WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, q, userID)
	return err
}
