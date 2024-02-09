package service

import (
	"context"
	"database/sql"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/tasuke/go-mux/model"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/crypto/bcrypt"
	"os"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
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
	// データベース操作で使用され、その操作によってオブジェクトが変更される可能性があるため
	newUser := &model.User{
		ID:       uuid.New().String(),
		Name:     signUpRequest.Name,
		Email:    signUpRequest.Email,
		Password: string(hashedPassword),
	}

	err = newUser.Insert(ctx, us.db, boil.Infer())
	if err != nil {
		return SignUpResponse{}, nil
	}

	// 値でOK。小さい構造体、レスポンスするためだけの構造体, 構造体を変更する必要がない,=> コードの意図がより明確
	signUpResponse := SignUpResponse{
		ID:    newUser.ID,
		Name:  newUser.Name,
		Email: newUser.Email,
	}

	return signUpResponse, nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (us *UserService) Login(loginRequest LoginRequest) (string, error) {
	currentUser, err := fetchUserByEmail(us.db, loginRequest.Email)
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

// fetchUserByEmail は与えられたメールアドレスに一致するユーザーをデータベースから検索します。
func fetchUserByEmail(db *sql.DB, email string) (*model.User, error) {
	// context.Background() は、実際のアプリケーションでは適切なコンテキストに置き換えるべきです。
	user, err := model.Users(qm.Where("email=?", email)).One(context.Background(), db)
	if err != nil {
		return nil, err
	}
	return user, nil
}
