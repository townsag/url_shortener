package main

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	// "go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	// "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	// "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	// "go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	// "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func setupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cdleanup will be invoked once.
	// this is a way to aggregate cleanup functions so that one function iteratively calls the other cleanup functions
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	// this is so resources are freed on failure to fully initialize the pipeline
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	// context propagation is used to share attributes between related spans
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up trace provider.
	// tracer provider is the factory object that creates the tracer object
	tracerProvider, err := newTracerProvider(ctx)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Set up meter provider.
	meterProvider, err := newMeterProvider(ctx)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	// Set up logger provider.
	// not sure to what extent I will be using this because my understanding is that logs are emitted by normal application 
	// logging libraries and then captured in another way
	// loggerProvider, err := newLoggerProvider()
	// if err != nil {
	// 	handleErr(err)
	// 	return
	// }
	// shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	// global.SetLoggerProvider(loggerProvider)

	return
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTracerProvider(ctx context.Context) (*trace.TracerProvider, error) {
	traceExporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			trace.WithBatchTimeout(time.Second)),
	)
	return tracerProvider, nil
}

func newMeterProvider(ctx context.Context) (*metric.MeterProvider, error) {
	metricExporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			metric.WithInterval(3*time.Second))),
	)
	return meterProvider, nil
}

// func newLoggerProvider() (*log.LoggerProvider, error) {
// 	logExporter, err := stdoutlog.New()
// 	if err != nil {
// 		return nil, err
// 	}

// 	loggerProvider := log.NewLoggerProvider(
// 		log.WithProcessor(log.NewBatchProcessor(logExporter)),
// 	)
// 	return loggerProvider, nil
// }
