package interceptor

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const traceIDKey = "x-trace-id"

func ServerTracingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	tracer := otel.Tracer("example/namedtracer/main")
	var span trace.Span
	ctx, span = tracer.Start(ctx, "operation")
	defer span.End()

	spanContext := span.SpanContext()
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(traceIDKey, spanContext.TraceID().String()))

	header := metadata.New(map[string]string{traceIDKey: spanContext.TraceID().String()})
	err := grpc.SendHeader(ctx, header)
	if err != nil {
		return nil, err
	}

	res, err := handler(ctx, req)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("err", err.Error()))
	}

	return res, err
}
