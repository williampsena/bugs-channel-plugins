// This package includes sentry plugin integration
package sentry

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gorilla "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/williampsena/bugs-channel-plugins/pkg/config"
	"github.com/williampsena/bugs-channel-plugins/pkg/event"
	"github.com/williampsena/bugs-channel-plugins/pkg/service"
)

// Represents the HTTP handlers and router instances.
type Server struct {
	Router http.Handler
	Srv    *http.Server
	log    *logrus.Logger
}

// Creates and returns a new instance of Server
func NewServer(handler http.Handler, log *logrus.Logger) *Server {
	ch := gorilla.CORS(gorilla.AllowedOrigins([]string{"*"}))

	return &Server{
		Router: handler,
		log:    log,
		Srv: &http.Server{
			Addr:         fmt.Sprintf(":%v", config.SentryPort()),
			Handler:      ch(handler),
			IdleTimeout:  time.Second * 5,
			ReadTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 5,
		},
	}
}

// Waits for an interrupt signal before gracefully shutting down the handlers.
func (s *Server) GraceFulShutDown(killTime time.Duration) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), killTime)

	defer cancel()

	s.log.Print("🛑 The server has been shut down.")

	if err := s.Srv.Shutdown(ctx); err != nil {
		s.log.Fatalf("❌ Unexpected interruption to the server's listening: %s\n", err)
	}

	s.log.Print("❎ The server exited properly")

}

// Turn on the HTTP handlers and listen in.
func (s *Server) ListenAndServe() error {
	return s.Srv.ListenAndServe()
}

// Shutdown terminates the handlers for HTTP.
func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Print("🛑 The server was shut down.")
	return s.Srv.Shutdown(ctx)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return gorilla.LoggingHandler(os.Stdout, next)
}

func buildRouter(serviceFetcher service.ServiceFetcher) *mux.Router {
	r := mux.NewRouter()

	r.PathPrefix("/health").HandlerFunc(HealthCheckEndpoint).Methods("GET")

	api := r.PathPrefix("/api/{id:[0-9]+}").Subrouter()
	api.Use(AuthKeyMiddleware(serviceFetcher))
	api.PathPrefix("/envelope").
		HandlerFunc(PostEventEndpoint(event.NewLoggerDispatcher())).
		Methods("POST")

	r.PathPrefix("/").HandlerFunc(NoRouteEndpoint)

	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(loggingMiddleware)
	r.Use()

	return r
}

func SetupServer(serviceFetcher service.ServiceFetcher) {
	r := buildRouter(serviceFetcher)

	srv := NewServer(r, logrus.StandardLogger())

	go srv.ListenAndServe()

	srv.log.Infof("🐜 Sentry Plugin Sever listening at port %v...", config.SentryPort())

	srv.GraceFulShutDown(time.Second * 5)
}
