package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/erry-az/mltp-go/internal/config"
	_ "github.com/grafana/pyroscope-go/godeltaprof/http/pprof"
	_ "net/http/pprof"
)

// StartRest starts the Rest server.
func StartRest(ctx context.Context, cfg config.RestConfig, handler http.Handler,
	shutdown func(ctx context.Context)) error {
	// Create a new server instance using config.
	// For more information, see https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeOut,
		WriteTimeout: cfg.WriteTimeOut,
	}

	go func() {
		slog.Info("starting server at ", "host", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
		err := server.ListenAndServe()
		if err != nil {
			slog.Error("failed to start rest server", err)
		}
	}()

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	<-s

	if shutdown != nil {
		shutdown(ctx)
	}

	return server.Shutdown(ctx)
}
