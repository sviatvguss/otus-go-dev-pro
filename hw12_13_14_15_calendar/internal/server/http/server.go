package internalhttp

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/sviatvguss/otus-go-dev-pro/hw12_13_14_15_calendar/internal/app"
	"github.com/sviatvguss/otus-go-dev-pro/hw12_13_14_15_calendar/internal/logger"

	"github.com/gorilla/mux"
)

var ErrServerNotStarted = errors.New("server not started")

type Server struct {
	logger *logger.Logger
	app    *app.App
	addr   string
	server *http.Server
}

func NewServer(logger *logger.Logger, app *app.App, host string, port string) *Server {
	return &Server{
		logger: logger,
		app:    app,
		addr:   net.JoinHostPort(host, port),
	}
}

func (s *Server) Start(ctx context.Context) error {
	router := mux.NewRouter()
	router.HandleFunc("/", s.RootHandler).Methods("GET")
	router.Use(s.loggingMiddleware)

	s.server = &http.Server{
		Addr:    s.addr,
		Handler: router,
	}

	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return ErrServerNotStarted
	}
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Server) RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Root"))
}
