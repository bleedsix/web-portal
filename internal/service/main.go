package service

import (
	"context"
	"errors" // Додали для errors.Is
	"log/slog"
	"net/http" // Додали для http.ServerClosed та http.Server
	"time"

	"github.com/bleedsix/web-portal/internal/config"
)

type service struct {
	cfg *config.Config
	log *slog.Logger
}

func newService(cfg *config.Config) *service {
	return &service{
		cfg: cfg,
		log: cfg.Log(),
	}
}

// run запускає HTTP сервер
func (s *service) run(ctx context.Context) error {
	s.log.Info("Service starting", "addr", s.cfg.ListenerAddr())

	// router() ми беремо з файлу router.go, він належить цій же структурі service
	srv := &http.Server{
		Addr:    s.cfg.ListenerAddr(),
		Handler: s.router(),
	}

	// Запуск сервера в горутині, щоб не блокувати головний потік
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Error("Server failed", "error", err)
		}
	}()

	// Чекаємо на сигнал зупинки (Ctrl+C)
	<-ctx.Done()
	s.log.Info("Service stopping...")

	// Graceful shutdown: даємо сервером 5 секунд на завершення поточних запитів
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	
	return nil
}

// Run - публічна точка входу
func Run(ctx context.Context, cfg *config.Config) {
	if err := newService(cfg).run(ctx); err != nil {
		// Логуємо помилку, якщо Graceful Shutdown не вдався
		cfg.Log().Error("Fatal service error", "error", err)
	}
}