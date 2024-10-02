// Package graceful stores utility functions that deal with graceful termination on selected OS signals.
package graceful

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// WithSignalShutdown returns a new ctx derived from the one given, which will be canceled when this process
// receives a SIGTERM or SIGQUIT signal.
func WithSignalShutdown(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

		<-c

		cancel()
	}()

	return ctx
}

// ListenAndServe wraps the equivalent http.Server method with graceful shutdown logic when this process receives a SIGTERM or SIGQUIT signal.
//
// Once the signal is received, shutdown process begins, giving the server additional 10 seconds to complete the remaining requests.
func ListenAndServe(srv *http.Server) error {
	return listenAndServe(srv, nil)
}

// ListenAndServeTLS wraps the equivalent http.Server method with graceful shutdown logic when this process receives a SIGTERM or SIGQUIT signal.
//
// Once the signal is received, shutdown process begins, giving the server additional 10 seconds to complete the remaining requests.
func ListenAndServeTLS(srv *http.Server, certFile, keyFile string) error {
	return listenAndServe(
		srv,
		&tlsCertFiles{
			certFile: certFile,
			keyFile:  keyFile,
		},
	)
}

type tlsCertFiles struct {
	certFile string
	keyFile  string
}

func listenAndServe(srv *http.Server, certAndKeyFiles *tlsCertFiles) error {
	idleConnsClosed := make(chan struct{})

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

		<-c

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
		close(idleConnsClosed)
	}()

	var err error

	if certAndKeyFiles != nil {
		err = srv.ListenAndServeTLS(certAndKeyFiles.certFile, certAndKeyFiles.keyFile)
	} else {
		err = srv.ListenAndServe()
	}

	if err != http.ErrServerClosed {
		return err
	}

	<-idleConnsClosed
	return nil
}
