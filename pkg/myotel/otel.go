package myotel

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

type Otel struct {
	traceProvider  *sdktrace.TracerProvider
	metricProvider *sdkmetric.MeterProvider
}

// New init sdk implementation for otel trace and metric
func New(ctx context.Context, otlpGrpcTarget, serviceName string, usePrometheusExporter bool) (*Otel, error) {
	conn, err := initConn(otlpGrpcTarget)
	if err != nil {
		return nil, err
	}

	res, err := initResource(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	tp, err := initTraceProvider(ctx, res, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics exporter: %w", err)
	}

	var metricReader sdkmetric.Reader

	if usePrometheusExporter {
		metricReader, err = prometheus.New()
		if err != nil {
			return nil, fmt.Errorf("failed to init prometheus exporter: %w", err)
		}
	} else {
		metricReader = sdkmetric.NewPeriodicReader(metricExporter)
	}

	mp, err := initMeterProvider(res, metricReader)
	if err != nil {
		return nil, err
	}

	return &Otel{
		traceProvider:  tp,
		metricProvider: mp,
	}, nil
}

// Shutdown all provider
func (o *Otel) Shutdown(ctx context.Context) {
	wg := sync.WaitGroup{}

	go func() {
		wg.Add(1)

		defer wg.Done()

		_ = o.traceProvider.Shutdown(ctx)
	}()

	go func() {
		wg.Add(1)

		defer wg.Done()

		_ = o.metricProvider.Shutdown(ctx)
	}()

	wg.Wait()
}

// initConn Initialize a gRPC connection to be used by both the tracer and meter
// providers.
func initConn(target string) (*grpc.ClientConn, error) {
	// It connects the OpenTelemetry Collector through local gRPC connection.
	conn, err := grpc.NewClient(target,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	return conn, err
}

// initTraceProvider Initializes an OTLP exporter, and configures the corresponding trace provider.
func initTraceProvider(ctx context.Context, res *resource.Resource, opts ...otlptracegrpc.Option) (*sdktrace.TracerProvider, error) {
	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tracerProvider)

	// Set global propagator to trace context (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider, nil
}

// initMeterProvider an prometheus exporter, and configures the corresponding meter provider.
func initMeterProvider(res *resource.Resource, reader sdkmetric.Reader) (*sdkmetric.MeterProvider, error) {
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetric.WithResource(res),
	)

	otel.SetMeterProvider(meterProvider)

	return meterProvider, nil
}

// initResource Initializes an OTLP exporter, and configures the corresponding meter provider.
func initResource(ctx context.Context, serviceName string) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			// The service name used to display traces in backends
			semconv.ServiceNameKey.String(serviceName),
		),
	)
}
