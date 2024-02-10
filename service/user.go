package service

import (
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/tasuke/go-mux/model"
	"github.com/tasuke/go-mux/repository"
	"golang.org/x/crypto/bcrypt"
	"os"
)

type UserService struct {
	ur repository.UserRepository
}

func NewUserService(ur repository.UserRepository) *UserService {
	return &UserService{ur: ur}
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

// サービスから直接エラーを返すようにしました。エラーメッセージの書き込みはコントローラに任せるべきです。
func (us *UserService) SignUp(ctx context.Context, signUpRequest SignUpRequest) (SignUpResponse, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signUpRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return SignUpResponse{}, nil
	}

	newUser := &model.User{
		ID:       uuid.New().String(),
		Name:     signUpRequest.Name,
		Email:    signUpRequest.Email,
		Password: string(hashedPassword),
	}

	signedUpUser, err := us.ur.SignUp(ctx, newUser)
	if err != nil {
		return SignUpResponse{}, nil
	}

	signUpResponse := SignUpResponse{
		ID:    signedUpUser.ID,
		Name:  signedUpUser.Name,
		Email: signedUpUser.Email,
	}

	return signUpResponse, nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (us *UserService) Login(loginRequest LoginRequest) (string, error) {
	currentUser, err := us.ur.GetUserByEmail(context.Background(), loginRequest.Email)
	if err != nil {
		return "", err
	}
	// パスワードの検証
	err = bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(loginRequest.Password))
	if err != nil {
		return "", err
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
