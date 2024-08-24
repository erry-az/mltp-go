package app

import (
	"github.com/grafana/loki-client-go/loki"
	slogloki "github.com/samber/slog-loki/v3"
	"log/slog"
)

type Log struct {
	Client *loki.Client
	Logger *slog.Logger
}

func NewLog(serviceName, url string) *Log {
	config, _ := loki.NewDefaultConfig(url)
	config.TenantID = serviceName
	client, _ := loki.New(config)

	logger := slog.New(slogloki.Option{Level: slog.LevelDebug, Client: client}.NewLokiHandler())
	logger = logger.
		With("service_name", serviceName)

	return &Log{
		Client: client,
		Logger: logger,
	}
}
