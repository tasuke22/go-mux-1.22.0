package controller

import (
	"encoding/json"
	"github.com/tasuke/go-mux/usecase"
	"net/http"
	"time"
)

type UserController struct {
	us *usecase.UserUsecase
}

func NewUserController(us *usecase.UserUsecase) *UserController {
	return &UserController{us: us}
}

func (uc *UserController) SignUp(w http.ResponseWriter, r *http.Request) {
	var signUpRequest usecase.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&signUpRequest); err != nil {
		sendErrorResponse(w, "リクエストのデコードに失敗しました", http.StatusBadRequest)
		return
	}

	signUpResponse, err := uc.us.SignUp(r.Context(), signUpRequest)
	if err != nil {
		sendErrorResponse(w, "サインアップに失敗しました", http.StatusInternalServerError)
		return
	}

	// レスポンスを設定
	sendJSONResponse(w, signUpResponse, http.StatusCreated)
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest usecase.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		sendErrorResponse(w, "リクエストのデコードに失敗しました", http.StatusBadRequest)
		return
	}
	tokenString, err := uc.us.Login(loginRequest)
	if err != nil {
		sendErrorResponse(w, "ログインに失敗しました", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)

	sendJSONResponse(w, LoginResponse{Token: tokenString}, http.StatusOK)
}

func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
	sendJSONResponse(w, "ログアウトしました", http.StatusOK)
}
