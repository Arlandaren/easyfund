package repository

import (
	"context"
	"github.com/Arlandaren/easyfund/ent"
	"github.com/Arlandaren/easyfund/ent/user"
)

type UserRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{client: client}
}

func (r *UserRepository) Create(ctx context.Context, email, passwordHash, role string) (*ent.User, error) {
	return r.client.User.
		Create().
		SetEmail(email).
		SetPasswordHash(passwordHash).
		SetRole(user.Role(role)).
		Save(ctx)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	return r.client.User.Query().Where(user.EmailEQ(email)).First(ctx)
}
