package otel

import (
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func initTracer() *sdktrace.TracerProvider {
	// สร้าง directory ถ้ายังไม่มี
	if err := os.MkdirAll("/app/logs", 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// เปิดไฟล์ log
	f, err := os.OpenFile("/app/logs/traces.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// สร้าง exporter ที่เขียนลงทั้ง file และ stdout
	exporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithWriter(f),
	)
	if err != nil {
		log.Fatalf("Failed to create exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("todo-service"),
			)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp
}
