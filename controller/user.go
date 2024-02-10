package controller

import (
	"encoding/json"
	"github.com/tasuke/go-mux/usecase"
	"net/http"
)

type UserController struct {
	us *usecase.UserUsecase
}

func NewUserController(us *usecase.UserUsecase) *UserController {
	return &UserController{us}
}

func (uc *UserController) SignUp(w http.ResponseWriter, r *http.Request) {
	var signUpRequest usecase.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&signUpRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// service層のSignUpメソッドを呼び出し
	// contextはDBと同じく一つのものをー
	signUpResponse, err := uc.us.SignUp(r.Context(), signUpRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// レスポンスを設定
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created ステータスコードを返す
	if err := json.NewEncoder(w).Encode(signUpResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest usecase.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	tokenString, err := uc.us.Login(loginRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// トークンの返却
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(LoginResponse{Token: tokenString}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
