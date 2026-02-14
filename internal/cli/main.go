package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bleedsix/web-portal/internal/config"
	"github.com/bleedsix/web-portal/internal/service"
	cli "github.com/urfave/cli/v2"
	"log/slog"
)

func Run(args []string) bool {
	// Базовий логер для самого CLI, поки конфіг не завантажено
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))

	app := &cli.App{
		Name:    "web-portal",
		Usage:   "My awesome project",
		Version: "0.0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "config.yaml", // Дефолтне ім'я
				Usage:   "Path to config file",
				EnvVars: []string{"CONFIG_PATH"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Run the application",
				Subcommands: []*cli.Command{
					{
						Name:  "service",
						Usage: "Start the API service",
						Action: func(c *cli.Context) error {
							// 1. Завантаження конфігу
							cfg, err := config.Load(c.String("config"))
							if err != nil {
								return err
							}

							// 2. Створення контексту з Graceful Shutdown (Ctrl+C)
							ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
							defer stop()

							// 3. Запуск сервісу
							service.Run(ctx, cfg)
							return nil
						},
					},
				},
			},
			{
				Name:  "migrate",
				Usage: "Database migrations",
				Subcommands: []*cli.Command{
					{
						Name:  "up",
						Usage: "Apply migrations",
						Action: func(c *cli.Context) error {
							cfg, err := config.Load(c.String("config"))
							if err != nil {
								return err
							}
							cfg.Log().Info("Migrating UP...")
							// return MigrateUp(cfg) // Твоя функція міграцій
							return nil
						},
					},
					{
						Name:  "down",
						Usage: "Rollback migrations",
						Action: func(c *cli.Context) error {
							cfg, err := config.Load(c.String("config"))
							if err != nil {
								return err
							}
							cfg.Log().Info("Migrating DOWN...")
							// return MigrateDown(cfg) // Твоя функція міграцій
							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(args); err != nil {
		log.Error("App crashed", "error", err)
		return false
	}
	return true
}