package transport

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/AI-Hackathon-2026/Clients-Service/pkg/jwtutil"
)

type wrappedWriter struct {
	status int
	http.ResponseWriter
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func LogRequest(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrp := &wrappedWriter{
			status:         http.StatusOK,
			ResponseWriter: w,
		}
		next.ServeHTTP(wrp, r)
		duration := time.Since(start)
		attributes := []slog.Attr{
			slog.Int("status", wrp.status),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Duration("duration", duration),
			slog.String("remote_addr", r.RemoteAddr),
		}
		if wrp.status >= 400 {
			logger.Warn("incoming request", slog.Any("details", attributes))
		} else {
			logger.Info("incoming request", slog.Any("details", attributes))
		}
	})
}

func AuthMiddleware(next http.Handler, secret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Unauthorized: invalid header format", http.StatusUnauthorized)
			return
		}

		claims, err := jwtutil.ValidateToken(parts[1], secret)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
