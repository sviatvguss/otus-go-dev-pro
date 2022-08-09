package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	ltime "github.com/sviatvguss/otus-go-dev-pro/hw12_13_14_15_calendar/internal/time"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		recorder := &StatusRecorder{
			ResponseWriter: w,
		}
		next.ServeHTTP(recorder, r)

		resp := fmt.Sprintf(
			"%s [%s] %s %s %s %d %v %s",
			r.RemoteAddr,
			now.Format(ltime.LogDateTimeFormat),
			r.Method,
			r.URL.Path,
			r.Proto,
			recorder.Status,
			time.Since(now),
			r.UserAgent(),
		)

		s.logger.Info(resp)
	})
}
