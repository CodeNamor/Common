package server

import (
	"context"

	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/CodeNamor/Common/logging" //TODO: Convert to updating logging client
	"github.com/gorilla/mux"
)

// ListenAndServe creates an http.Server, sets up signal handler
// to listen for SIGINT and SIGTERM to perform graceful shutdown,
// and launches the server. If the server returns anything other
// than a normal close, then it is returned, otherwise returns nil.ListenAndServe
// It logs at the Info level
func ListenAndServe(addr string, handler http.Handler) error {
	logging.Info("HTTP Server address " + addr)
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// setup signals channel for Ctrl+C and SIGTERM
	signals := make(chan os.Signal, 1)                    // channel to listen to signals
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM) // forward ctrl-c + SIGTERM to signals chan

	wgServer := sync.WaitGroup{}
	wgServer.Add(2) // wait for signal and listen routines

	var serverError error

	go func() { // routine to wait for any shutdown signal
		defer wgServer.Done()
		killSignal := <-signals // wait for a signal
		switch killSignal {
		case os.Interrupt:
			logging.Info("SIGINT received (Control-C ?)")
		case syscall.SIGTERM:
			logging.Info("SIGTERM received (Kubernetes shutdown?)")
		case nil: // if we sent a nil signal in, then exit now
			return // exit now, probably due to server error
		}

		logging.Info("graceful shutdown initiated...")
		srv.Shutdown(context.Background())
		logging.Info("graceful shutdown complete")
	}()

	go func() { // routine to begin listening and serving
		defer wgServer.Done()
		serverError = srv.ListenAndServe()
		if !isNormalShutdown(serverError) {
			signals <- nil // we need to exit now, tell other routine
		}
	}()

	// wait for rountines to finish
	wgServer.Wait()
	if isNormalShutdown(serverError) {
		return nil
	}

	return serverError
}

// CreateAndHandleReadinessLiveness creates atomic handlers for
// readiness and liveness, it registers them at the readyURLPath
// and liveURLPath and it returns a readyFn and livenessFn updater
// functions which can be used to set the readiness and liveness
// status at any time. The atomicHandlers will return 200 OK or
// 503 Service Unavailable based on the state of the atomicHandler
// from the last update using the UpdateFn (readyFn, livenessFn),
// the default state is false for both.
func CreateAndHandleReadinessLiveness(router *mux.Router, readyURLPath string, liveURLPath string) (UpdateFn, UpdateFn) {
	readyInitial := false
	readyHandler, readyFn := CreateAtomicHandler(readyInitial)
	router.HandleFunc(readyURLPath, readyHandler)

	liveInitial := false
	livenessHandler, livenessFn := CreateAtomicHandler(liveInitial)
	router.HandleFunc(liveURLPath, livenessHandler)

	return readyFn, livenessFn // return updater fns
}

type UpdateFn func(bool)

// CreateAtomicHandler creates an atomic value handler which
// returns 200 OK or 503 Service Unavailable based on the state
// of the atomic value which defaults to initialValue.
// It returns the handler and an update func which can be used
// to atomically update the state.
func CreateAtomicHandler(initialValue bool) (http.HandlerFunc, UpdateFn) {
	atomicValue := &atomic.Value{}
	atomicValue.Store(initialValue)

	handler := func(w http.ResponseWriter, _ *http.Request) {
		if !atomicValue.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		w.Write([]byte("OK"))
	}

	updateFn := func(value bool) {
		atomicValue.Store(value)
	}

	return handler, updateFn
}

func isNormalShutdown(err error) bool {
	return err == http.ErrServerClosed
}

// DefaultHandler sets the Content-Type header to "application/json"
func DefaultHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		fn(rw, req)
	}
}
