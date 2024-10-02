package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/vmarunin/mts2024golang/seminar2/internal/pkg/config"
	"github.com/vmarunin/mts2024golang/seminar2/internal/pkg/delivery/rest"
	"github.com/vmarunin/mts2024golang/seminar2/internal/pkg/service/pinger"
)

type App struct {
	cfg      *config.Config
	log      *slog.Logger
	quitChan chan os.Signal
}

func NewApp(cfg *config.Config) *App {
	logFile, err := os.OpenFile(cfg.GetLogPath(), os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	logH := slog.NewJSONHandler(logFile, &slog.HandlerOptions{})
	return &App{
		cfg: cfg,
		log: slog.New(logH),
	}
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var err error

	pingerSrv := pinger.New(a.cfg, a.log)

	srv, err := rest.NewServer(a.cfg, a.log, pingerSrv)
	if err != nil {
		a.log.ErrorContext(ctx, "server creation error", slog.Any("error", err))
		return err
	}

	go func() {
		a.quitChan = make(chan os.Signal, 1)
		signal.Notify(a.quitChan, syscall.SIGTERM, syscall.SIGINT)
		<-a.quitChan
		cancel()
		a.log.InfoContext(ctx, "Shutting down server...")
	}()

	a.log.InfoContext(ctx, "Starting server...")
	err = srv.Run(ctx)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.log.ErrorContext(ctx, "server run error", slog.Any("error", err))
		return err
	}
	a.log.InfoContext(ctx, "Server shutdown completed")

	return nil
}

func (a *App) Stop(_ context.Context) error {
	a.quitChan <- os.Interrupt
	return nil
}
