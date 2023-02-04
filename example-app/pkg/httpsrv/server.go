package httpsrv

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Start(svcName string, port int, handler http.Handler) {
	log.SetPrefix(fmt.Sprintf("[%s] ", svcName))

	// Create our application context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cancel the context when we get a signal from the OS
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		s := <-done
		log.Printf("received signal %v to shutdown", s.String())
		cancel()
	}()

	// Start the HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", port),
		Handler: handler,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		// When the context is cancelled, shutdown the server
		<-ctx.Done()
		if err := srv.Close(); err != nil {
			log.Printf("unable to shutdown HTTP server: %v", err)
		}
	}()

	// Start listening for requests
	log.Printf("lisenting for HTTP requests on %d", port)
	err := srv.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Printf("unable to listen for HTTP connections: %v", err)
		}
	}

	os.Exit(0)
}
