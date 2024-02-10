package repository

import (
	"context"
	"database/sql"
	"github.com/tasuke/go-mux/model"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) SignUp(ctx context.Context, newUser *model.User) (*model.User, error) {
	if err := newUser.Insert(ctx, ur.db, boil.Infer()); err != nil {
		return nil, err
	}
	return newUser, nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := model.Users(model.UserWhere.Email.EQ(email)).One(ctx, ur.db)
	if err != nil {
		return nil, err
	}
	return user, nil
}
