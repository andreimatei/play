package main

import (
	"context"
	"go.opentelemetry.io/contrib/zpages"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()
	//exporter, err := stdout.NewExporter(
	//	stdout.WithPrettyPrint(),
	//)
	//exporter, err := otlp.NewExporter(ctx, otlpgrpc.NewDriver(otlpgrpc.WithEndpoint("localhost:49161")))
	//exporter, err := otlp.NewExporter(ctx, otlpgrpc.NewDriver(otlpgrpc.WithEndpoint("localhost:55680")))
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}
	//err = exporter.Start(ctx)
	//if err != nil {
	//	log.Fatalf("failed to start exporter: %v", err)
	//}
	defer exporter.Shutdown(ctx)

	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	zpagesSp := zpages.NewSpanProcessor()
	tracezHandler := zpages.NewTracezHandler(zpagesSp)
	http.Handle("/tracez", tracezHandler)

	go func() {
		http.ListenAndServe("localhost:8080", tracezHandler)
		log.Printf("listening on localhost:8080")
	}()

	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(bsp), sdktrace.WithSpanProcessor(zpagesSp), sdktrace.WithSampler(sdktrace.AlwaysSample()))

	otel.SetTracerProvider(tp)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	otel.SetTextMapPropagator(propagator)

	// Handle this error in a sensible manner where possible
	defer func() { _ = tp.Shutdown(ctx) }()

	lemonsKey := attribute.Key("ex.com/lemons")
	anotherKey := attribute.Key("ex.com/another")

	tracer := otel.Tracer("ex.com/basic")
	member, err := baggage.NewMember("foo1", "bar1")
	if err != nil {
		panic(err)
	}
	b, err := baggage.New(member)
	if err != nil {
		panic(err)
	}
	ctx = baggage.ContextWithBaggage(ctx, b)
	//fooKey.String("foo1"),
	//	barKey.String("bar1"),
	//)
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

	time.Sleep(100 * time.Second)
}
