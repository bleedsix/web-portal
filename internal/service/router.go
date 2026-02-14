package service

import (
	"net/http"
	"time"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// router налаштовує маршрути та мідлвари
func (s *service) router() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	
	// Наша кастомна мідлвара для логів
	r.Use(s.loggingMiddleware)

	r.Route("/integrations/web-portal", func(r chi.Router) {
		// Тут будуть твої ендпоінти
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		})
	})

	return r
}

// loggingMiddleware логує запити через slog
func (s *service) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		s.log.Info("HTTP Request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", ww.Status()),
			slog.String("request_id", middleware.GetReqID(r.Context())),
			slog.Duration("duration", time.Since(start)),
		)
	})
}