package handlers

import (
	"go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer = otel.GetTracerProvider().Tracer("handlers")