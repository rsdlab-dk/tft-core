package http

import (
	"context"
	"net/http"
	"time"

	"github.com/rsdlab-dk/tft-core/logger"
	"github.com/rsdlab-dk/tft-core/ratelimit"
	"go.uber.org/zap"
)

func WithRequestID(log *logger.Logger) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			requestID := logger.GenerateRequestID()
			ctx := logger.WithRequestID(r.Context(), requestID)
			
			w.Header().Set("X-Request-ID", requestID)
			
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

func WithRateLimit(limiter ratelimit.Limiter, endpoint string, log *logger.Logger) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			config := ratelimit.NewConfig()
			rule := config.GetRule(endpoint)
			
			key := endpoint + ":" + r.RemoteAddr
			allowed, err := limiter.Allow(r.Context(), key, rule.Rate, rule.Window)
			if err != nil {
				log.WithContext(r.Context()).Error("rate limiter error",
					zap.String("endpoint", endpoint),
					zap.Error(err))
				WriteInternalError(w, log, r)
				return
			}
			
			if !allowed {
				log.WithContext(r.Context()).Warn("rate limit exceeded",
					zap.String("endpoint", endpoint),
					zap.String("key", key))
				WriteError(w, "RATE_LIMIT_EXCEEDED", "Rate limit exceeded", http.StatusTooManyRequests, log, r)
				return
			}
			
			next.ServeHTTP(w, r)
		}
	}
}

func WithLogging(log *logger.Logger) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			
			next.ServeHTTP(wrapped, r)
			
			duration := time.Since(start)
			
			log.WithContext(r.Context()).Info("request completed",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", wrapped.statusCode),
				zap.Duration("duration", duration),
				zap.String("user_agent", r.UserAgent()),
			)
		}
	}
}

func WithCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Request-ID")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}