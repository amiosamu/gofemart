package transport

import (
	"context"
	"net/http"
	"time"

	"github.com/amiosamu/gofemart/internal/domain"
	log "github.com/sirupsen/logrus"
)

type (
	ResponseData struct {
		Status int
		Size   int
	}

	LoggingResponseWriter struct {
		http.ResponseWriter
		ResponseData *ResponseData
	}
)

func (r *LoggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.ResponseData.Size += size
	return size, err
}

func (r *LoggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.ResponseData.Status = statusCode
}

func logFields(handler string) log.Fields {
	return log.Fields{
		"handler": handler,
	}
}

func logError(handler string, err error) {
	log.WithFields(logFields(handler)).Error(err)
}


func withLogging(next http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &ResponseData{
			Status: 0,
			Size:   0,
		}
		lw := LoggingResponseWriter{
			ResponseWriter: w,
			ResponseData:   responseData,
		}

		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		log.WithFields(log.Fields{
			"uri":      r.RequestURI,
			"method":   r.Method,
			"duration": duration,
			"status":   responseData.Status,
			"size":     responseData.Size,
		}).Info("request details: ")
	}
	return http.HandlerFunc(logFn)
}


func (s *APIServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getTokenFromRequest(r)
		if err != nil {
			logError("authMiddleware", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userID, err := s.users.ParseToken(r.Context(), token)
		if err != nil {
			logError("authMiddleware", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), domain.UserIDKeyForContext, userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func getTokenFromRequest(r *http.Request) (string, error) {
	token, err := r.Cookie("token")
	if err != nil {
		return "", err

	}
	return token.Value, nil
}
