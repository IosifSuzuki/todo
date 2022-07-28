package middleware

import (
	"context"
	"github.com/IosifSuzuki/todo/internall/logger"
	"github.com/IosifSuzuki/todo/internall/utility"
	"net/http"
)

type AuthenticationMiddleware struct {
}

func (a *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenText := r.Header.Get("Authorization")
		if isValidToken, _ := utility.VerifyToken(tokenText); !isValidToken {
			var errorMsg = "token isn't valid"
			logger.Error(errorMsg)
			http.Error(w, errorMsg, http.StatusUnauthorized)
			return
		}
		if isAccessToken, err := utility.VerifyIsAccessToken(tokenText); !isAccessToken {
			var errorMsg = err.Error()
			logger.Error(errorMsg)
			http.Error(w, errorMsg, http.StatusUnauthorized)
			return
		}
		var userId int
		var err error
		if userId, err = utility.GetUserIdByFromToken(tokenText); err != nil {
			http.Error(w, "access denied", http.StatusUnauthorized)
			return
		}
		reqContext := context.WithValue(r.Context(), utility.UserIdKey, userId)
		next.ServeHTTP(w, r.WithContext(reqContext))
	})
}
