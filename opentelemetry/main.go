package main

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	//exporter, err := stdout.NewExporter(
	//	stdout.WithPrettyPrint(),
	//)
	//exporter, err := otlp.NewExporter(ctx, otlpgrpc.NewDriver(otlpgrpc.WithEndpoint("localhost:49161")))
	//exporter, err := otlp.NewExporter(ctx, otlpgrpc.NewDriver(otlpgrpc.WithEndpoint("localhost:55680")))
	exporter, err := otlp.NewExporter(ctx, otlpgrpc.NewDriver(otlpgrpc.WithInsecure()))
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}
	//err = exporter.Start(ctx)
	//if err != nil {
	//	log.Fatalf("failed to start exporter: %v", err)
	//}
	defer exporter.Shutdown(ctx)

	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(bsp),sdktrace.WithSampler(sdktrace.AlwaysSample()))

	otel.SetTracerProvider(tp)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	otel.SetTextMapPropagator(propagator)

	// Handle this error in a sensible manner where possible
	defer func() { _ = tp.Shutdown(ctx) }()

	fooKey := attribute.Key("ex.com/foo")
	barKey := attribute.Key("ex.com/bar")
	lemonsKey := attribute.Key("ex.com/lemons")
	anotherKey := attribute.Key("ex.com/another")

	tracer := otel.Tracer("ex.com/basic")
	ctx = baggage.ContextWithValues(ctx,
		fooKey.String("foo1"),
		barKey.String("bar1"),
	)
	log.Printf("1")

	func(ctx context.Context) {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "operation")
		defer span.End()

		span.AddEvent("Nice operation!", trace.WithAttributes(attribute.Int("bogons", 100)))
		span.SetAttributes(anotherKey.String("yes"))

		func(ctx context.Context) {
			var span trace.Span
			ctx, span = tracer.Start(ctx, "Sub operation...")
			defer span.End()

			span.SetAttributes(lemonsKey.String("five"))
			span.AddEvent("Sub span event")
		}(ctx)
	}(ctx)
	log.Printf("2")

	time.Sleep(10 * time.Second)
}