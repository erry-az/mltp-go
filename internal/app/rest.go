package app

import (
	"context"
	"errors"
	"github.com/erry-az/mltp-go/db/query"
	"github.com/erry-az/mltp-go/internal/config"
	"github.com/erry-az/mltp-go/internal/handler"
	"github.com/erry-az/mltp-go/internal/server"
	localMiddleware "github.com/erry-az/mltp-go/internal/server/middleware"
	"github.com/erry-az/mltp-go/internal/service"
	"github.com/erry-az/mltp-go/pkg/myotel"
	"github.com/go-playground/validator/v10"
	deltapprof "github.com/grafana/pyroscope-go/godeltaprof/http/pprof"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"net/http"
	"net/http/pprof"
)

type ValidationError struct {
	Fields []ValidationErrorField `json:"fields"`
}

type ValidationErrorField struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
}

func NewRest(cfg config.Config) error {
	ctx := context.Background()

	otel, err := myotel.New(ctx, cfg.Otel.Endpoint, cfg.ServiceName, cfg.Otel.UsePrometheus)
	if err != nil {
		return err
	}

	dbPool, err := NewDBPool(cfg.DB)
	if err != nil {
		return err
	}

	queries := query.New(dbPool)

	services := service.New(queries, dbPool)

	jwt, err := localMiddleware.NewJWT(cfg.Jwt)
	if err != nil {
		return err
	}

	handlers := handler.NewHandler(services, jwt)

	log := NewLog(cfg.ServiceName, cfg.Log.Url)

	e := echo.New()

	e.Validator = localMiddleware.NewValidator()

	e.Use(otelecho.Middleware(cfg.ServiceName))
	e.Use(localMiddleware.NewWithConfig(log.Logger, localMiddleware.Config{
		WithSpanID:         true,
		WithTraceID:        true,
		WithUserAgent:      false,
		WithRequestID:      false,
		WithRequestHeader:  false,
		WithResponseBody:   false,
		WithResponseHeader: false,
	}))
	e.Use(middleware.Recover())

	e.HTTPErrorHandler = errorHandler(e)

	monitor(e)

	route(e, handlers, jwt)

	return server.StartRest(ctx, cfg.RestServer, e, func(ctx context.Context) {
		otel.Shutdown(ctx)
		_ = dbPool.Close()
		log.Client.Stop()
	})
}

func route(e *echo.Echo, h *handler.Handler, jwt *localMiddleware.Jwt) {
	e.POST("/create_user", h.CreateUser)
	e.GET("/balance_read", h.BalanceRead, jwt.Middleware())
	e.GET("/top_users", h.TopUsers, jwt.Middleware())
	e.GET("/top_transactions_per_user", h.TopTransactions, jwt.Middleware())
	e.POST("/balance_topup", h.TopUp, jwt.Middleware())
	e.POST("/transfer", h.Transfer, jwt.Middleware())
}

func errorHandler(e *echo.Echo) func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		var val validator.ValidationErrors

		switch {
		case errors.As(err, &val):
			fields := make([]ValidationErrorField, 0, len(val))
			for _, fieldError := range val {
				fields = append(fields, ValidationErrorField{
					Field: fieldError.Field(),
					Tag:   fieldError.Tag(),
				})
			}

			_ = c.JSON(http.StatusBadRequest, ValidationError{Fields: fields})
		default:
			e.DefaultHTTPErrorHandler(err, c)
		}
	}
}

func monitor(e *echo.Echo) {
	converter := func(h http.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			h.ServeHTTP(c.Response().Writer, c.Request())
			return nil
		}
	}

	e.GET("/metrics", echoprometheus.NewHandler())

	prefixRouter := e.Group("/debug/pprof")
	{
		prefixRouter.Any("/", converter(pprof.Index))
		prefixRouter.Any("/allocs", converter(pprof.Handler("allocs").ServeHTTP))
		prefixRouter.Any("/block", converter(pprof.Handler("block").ServeHTTP))
		prefixRouter.Any("/cmdline", converter(pprof.Cmdline))
		prefixRouter.Any("/goroutine", converter(pprof.Handler("goroutine").ServeHTTP))
		prefixRouter.Any("/heap", converter(pprof.Handler("heap").ServeHTTP))
		prefixRouter.Any("/mutex", converter(pprof.Handler("mutex").ServeHTTP))
		prefixRouter.Any("/profile", converter(pprof.Profile))
		prefixRouter.Any("/symbol", converter(pprof.Symbol))
		prefixRouter.Any("/threadcreate", converter(pprof.Handler("threadcreate").ServeHTTP))
		prefixRouter.Any("/trace", converter(pprof.Trace))

		prefixRouter.Any("/delta_heap", converter(deltapprof.Heap))
		prefixRouter.Any("/delta_block", converter(deltapprof.Block))
		prefixRouter.Any("/delta_mutex", converter(deltapprof.Mutex))
	}
}
