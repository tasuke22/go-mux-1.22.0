package middleware

import (
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
)

func JWTMiddleware(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cookieに設定してあるトークンを取得
		c, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString := c.Value

		// jwtトークンの検証
		token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			http.Error(w, "認証トークンが見つかりません。", http.StatusUnauthorized)
			return
		}

		if err != nil || !token.Valid {
			http.Error(w, "認証トークンが無効です。", http.StatusUnauthorized)
			return
		}

		// 次のハンドラを呼び出す
		next(w, r)
	})
}
