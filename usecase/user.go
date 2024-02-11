package usecase

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/tasuke/go-mux/model"
	"github.com/tasuke/go-mux/repository"
	"golang.org/x/crypto/bcrypt"
	"os"
)

type UserUsecase struct {
	ur repository.UserRepository
}

func NewUserUsecase(ur repository.UserRepository) *UserUsecase {
	return &UserUsecase{ur: ur}
}

type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (us *UserUsecase) SignUp(ctx context.Context, signUpRequest SignUpRequest) (SignUpResponse, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signUpRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return SignUpResponse{}, err
	}

	newUser := &model.User{
		ID:       uuid.New().String(),
		Name:     signUpRequest.Name,
		Email:    signUpRequest.Email,
		Password: string(hashedPassword),
	}

	signedUpUser, err := us.ur.SignUp(ctx, newUser)
	if err != nil {
		return SignUpResponse{}, err
	}

	return SignUpResponse{
		ID:    signedUpUser.ID,
		Name:  signedUpUser.Name,
		Email: signedUpUser.Email,
	}, nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (us *UserUsecase) Login(loginRequest LoginRequest) (string, error) {
	currentUser, err := us.ur.GetUserByEmail(context.Background(), loginRequest.Email)
	if err != nil {
		return "", err
	}
	// パスワードの検証
	err = bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(loginRequest.Password))
	if err != nil {
		return "", errors.New("無効なパスワードです") // より具体的なエラーメッセージ
	}
	// JWTトークンの生成準備
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": currentUser.ID,
		"sub":     currentUser.ID,
	})
	// JWTトークンの生成(署名)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}
