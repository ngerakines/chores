package chores

import (
	"log"
	"net/http"
	"time"
)

func loggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer func() {
				logger.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), time.Now().Sub(start))
			}()
			next.ServeHTTP(w, r)
		})
	}
}
