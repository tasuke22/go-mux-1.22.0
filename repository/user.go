package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/tasuke/go-mux/model"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type UserRepository interface {
	SignUp(ctx context.Context, newUser *model.User) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) SignUp(ctx context.Context, newUser *model.User) (*model.User, error) {
	if err := newUser.Insert(ctx, ur.db, boil.Infer()); err != nil {
		return nil, err
	}
	return newUser, nil
}

func (ur *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := model.Users(model.UserWhere.Email.EQ(email)).One(ctx, ur.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
