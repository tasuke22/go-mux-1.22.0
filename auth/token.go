// auth/token.go
package auth

import (
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
)

type AuthToken struct{}

func NewAuthToken() *AuthToken {
	return &AuthToken{}
}

// ExtractUserIDFromToken はリクエストからJWTトークンを抽出し、そこからuserIDを取り出します。
func (at *AuthToken) ExtractUserIDFromToken(r *http.Request) (string, error) {
	c, err := r.Cookie("token")
	if err != nil {
		return "", err
	}
	tokenString := c.Value

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", http.ErrNoCookie
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", http.ErrNoCookie
	}

	return userID, nil
}
