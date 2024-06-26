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
	Router  http.Handler
	Srv     *http.Server
	log     *logrus.Logger
	Context *ServerContext
}

// Represents the server context
type ServerContext struct {
	context.Context

	// The service fetcher instance
	ServiceFetcher service.ServiceFetcher

	// The event dispatcher instance
	EventsDispatcher event.EventsDispatcher
}

// Creates and returns a new instance of Server
func NewServer(c *ServerContext, handler http.Handler, log *logrus.Logger) *Server {
	ch := gorilla.CORS(gorilla.AllowedOrigins([]string{"*"}))

	return &Server{
		Context: c,
		Router:  handler,
		log:     log,
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

	ctx, cancel := context.WithTimeout(s.Context, killTime)

	defer cancel()

	s.log.Print("🛑 The sentry server has been shut down.")

	if err := s.Srv.Shutdown(ctx); err != nil {
		s.log.Fatalf("❌ Unexpected interruption to the sentry server's listening: %s\n", err)
	}

	s.log.Print("❎ The sentry server exited properly")

}

// Turn on the HTTP handlers and listen in.
func (s *Server) ListenAndServe() error {
	return s.Srv.ListenAndServe()
}

// Shutdown terminates the handlers for HTTP.
func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Print("🛑 The sentry server was shut down.")
	return s.Srv.Shutdown(ctx)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return gorilla.LoggingHandler(os.Stdout, next)
}

func buildRouter(c *ServerContext) *mux.Router {
	r := mux.NewRouter()

	r.PathPrefix("/health").HandlerFunc(HealthCheckEndpoint).Methods("GET")

	api := r.PathPrefix("/api/{id:[0-9]+}").Subrouter()
	api.Use(AuthKeyMiddleware(c.ServiceFetcher))
	api.PathPrefix("/envelope").
		HandlerFunc(PostEventEndpoint(c.EventsDispatcher)).
		Methods("POST")

	r.PathPrefix("/").HandlerFunc(NoRouteEndpoint)

	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(loggingMiddleware)
	r.Use()

	return r
}

func SetupServer(srv *Server) {
	go srv.ListenAndServe()

	Greetings(srv)

	srv.GraceFulShutDown(time.Second * 5)
}

func Greetings(srv *Server) {
	srv.log.Infof("🐜 Sentry Plugin Sever listening at port %v...", config.SentryPort())
}

func BuildServer(c *ServerContext) *Server {
	r := buildRouter(c)

	srv := NewServer(c, r, logrus.StandardLogger())

	return srv
}
