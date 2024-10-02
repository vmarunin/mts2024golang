package rest

import (
	"log/slog"
	"net/http"
	"time"
)

func (s *Server) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		rw := &mwResponseWriter{writer, 0}
		startTime := time.Now()
		next.ServeHTTP(rw, request)
		s.log.Info("request",
			slog.String("method", request.Method),
			slog.String("uri", request.RequestURI),
			slog.Int("code", rw.statusCode),
			slog.Int("duration", int(time.Since(startTime).Milliseconds())),
		)
	})
}

func (s *Server) recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				s.log.Error("panic recovered", slog.Any("error", rvr))
				writer.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(writer, request)
	})
}

type mwResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (m *mwResponseWriter) WriteHeader(statusCode int) {
	m.statusCode = statusCode
	m.ResponseWriter.WriteHeader(statusCode)
}
