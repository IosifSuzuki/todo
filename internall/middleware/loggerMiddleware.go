package middleware

import (
	"github.com/IosifSuzuki/todo/internall/logger"
	"go.uber.org/zap"
	"net/http"
)

type LoggerMiddleware struct {
}

func (l *LoggerMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("made request with parameters",
			zap.String("url", r.URL.String()),
			zap.String("method", r.Method),
		)
		next.ServeHTTP(w, r)
	})
}
