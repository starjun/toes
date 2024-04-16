package global

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Trace gin.HandlerFunc = func(ctx *gin.Context) {}

type OpenTelemetry struct {
	Endpoint string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	Headers  struct {
		Authorization string `mapstructure:"authorization" json:"authorization" yaml:"authorization"`
		Organization  string `mapstructure:"organization" json:"organization" yaml:"organization"`
		StreamName    string `mapstructure:"streamName" json:"streamName" yaml:"streamName"`
	} `mapstructure:"headers" json:"headers" yaml:"headers"`
	Tls struct {
		Insecure bool `mapstructure:"insecure" json:"insecure" yaml:"insecure"`
	} `mapstructure:"tls" json:"tls" yaml:"tls"`
}

var Tracer = otel.Tracer("toes")
var serviceName = "toes"

func InitTrace(ctx context.Context) {
	if Cfg.AppName != "" {
		serviceName = Cfg.AppName
	}

	Tracer = otel.Tracer(serviceName)
	Trace = otelgin.Middleware(serviceName)

	_, err := initProvider()
	if err != nil {
		panic(err)
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}

func initProvider() (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	conn, err := grpc.NewClient(Cfg.OpenTelemetry.Endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn), otlptracegrpc.WithHeaders(
		map[string]string{
			"Authorization": Cfg.OpenTelemetry.Headers.Authorization,
			"organization":  Cfg.OpenTelemetry.Headers.Organization,
			"stream-name":   Cfg.OpenTelemetry.Headers.StreamName,
		},
	))
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

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}
