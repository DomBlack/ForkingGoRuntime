package rest

import (
	"context"
	"errors"
	"fmt"
	golog "log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/DomBlack/ForkingGoRuntime/example-app/pkg/tracing"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/DomBlack/ForkingGoRuntime/example-app/pkg/tracing"
)

// Server is a simple HTTP server with registered handlers
type Server struct {
	name   string
	port   int
	Log    zerolog.Logger
	router *httprouter.Router
}

// NewServer creates a new Server
func NewServer(name string, port int) *Server {
	tracing.Init(name)

	router := httprouter.New()

	router.NotFound = notFoundHandler
	router.MethodNotAllowed = methodNotAllowedHandler
	router.PanicHandler = panicHandler

	return &Server{
		name:   name,
		port:   port,
		Log:    log.With().Str("service", name).Logger(),
		router: router,
	}
}

// init setups our default logging
func init() {
	log.Logger = log.With().Caller().Stack().Timestamp().Logger().Output(zerolog.NewConsoleWriter())
}

// Start starts the server, and blocks until a OS signal is received
// to shutdown the application
func (s *Server) Start() {
	// Set the default go logger to the servers logger
	goDefaultLogger := golog.Default()
	goDefaultLogger.SetFlags(0)
	goDefaultLogger.SetOutput(s.Log.With().CallerWithSkipFrameCount(4).Logger())

	// Create our application context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cancel the context when we get a signal from the OS
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		signalReceived := <-done
		s.Log.Warn().Str("signal", signalReceived.String()).Msg("received signal to shutdown")
		cancel()
	}()

	// Start the HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", s.port),
		Handler: s.router,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		// When the context is cancelled, shutdown the server
		<-ctx.Done()
		if err := srv.Close(); err != nil {
			log.Err(err).Msg("unable to shutdown HTTP server")
		}
	}()

	// Start listening for requests
	s.Log.Info().Int("port", s.port).Msg("listening for HTTP requests")
	err := srv.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			s.Log.Err(err).Msg("unable to listen for HTTP connections")
		}
	}

	s.Log.Info().Msg("shutdown cleanly")
	os.Exit(0)
}
