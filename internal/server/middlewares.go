package server

import (
	"fmt"
	"github.com/Aegon95/home24-webscraper/internal"
	"go.uber.org/zap"
	"net/http"
)

type MiddlewareManager interface {
	RecoverPanic(next http.Handler) http.Handler
	LogRequest(next http.Handler) http.Handler
	SecureHeaders(next http.Handler) http.Handler
}

func NewMiddlewareManager(logger *zap.SugaredLogger, helper internal.Helper) MiddlewareManager {
	return &middlewares{
		logger,
		helper,
	}
}

type middlewares struct {
	logger *zap.SugaredLogger
	helper internal.Helper
}

// recoverPanic middleware recovers the app, when app panics
func (m *middlewares) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {

			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				m.helper.ServerError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// logRequest middleware logs every request
func (m *middlewares) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Debugf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

// secureHeaders middleware adds secureHeaders to response
func (m *middlewares) SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}
