package rest

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/vmarunin/mts2024golang/seminar2/internal/pkg/config"
	"github.com/vmarunin/mts2024golang/seminar2/internal/pkg/service/pinger"
)

type Server struct {
	srv    *http.Server
	log    *slog.Logger
	cfg    *config.Config
	pinger *pinger.Service
}

func NewServer(
	cfg *config.Config, log *slog.Logger,
	pinger *pinger.Service,
) (*Server, error) {
	timeout, err := time.ParseDuration(cfg.GetTimeout())
	if err != nil {
		return nil, fmt.Errorf("%w: %w timeout value: %s", config.ErrBadConfig, err, cfg.GetTimeout())
	}
	maxHeaderBytes, err := strconv.ParseInt(cfg.GetMaxHeaderBytes(), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: %w max_header_bytes value: %s", config.ErrBadConfig, err, cfg.GetMaxHeaderBytes())
	}

	server := &Server{
		cfg:    cfg,
		log:    log.With(slog.String("component", "rest")),
		pinger: pinger,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/task", server.CreateTask)
	mux.HandleFunc("GET /api/v1/log", server.GetLog)

	mux.HandleFunc("GET /api/v1/simple/name", server.SimpleName)
	mux.HandleFunc("GET /api/v1/auth/name", server.AuthName)
	mux.HandleFunc("GET /start", server.Start)

	server.srv = &http.Server{
		Addr: ":" + cfg.GetPort(),
		Handler: server.logMiddleware(
			server.recoverMiddleware(
				mux,
			),
		),
		ReadTimeout:    timeout,
		WriteTimeout:   timeout,
		MaxHeaderBytes: int(maxHeaderBytes),
	}

	return server, nil
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		_ = s.srv.Shutdown(context.Background())
	}()
	return s.srv.ListenAndServe()
}
