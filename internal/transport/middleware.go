package transport

import (
	"log"
	"net/http"
	"time"
)

type wrappedWriter struct {
    status int
    http.ResponseWriter
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
    w.status = statusCode
    w.ResponseWriter.WriteHeader(statusCode)
}

func LogRequest(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        wrp := &wrappedWriter{
            status: http.StatusOK,
            ResponseWriter: w,
        }
        
        next.ServeHTTP(wrp, r)
        log.Println(wrp.status, r.Method, r.URL.Path, time.Since(start))
    })
}
