package interceptor

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ctxKey string

const requestIDKey ctxKey = "request_id"

func UnaryLoggingInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()

		ctx, req_id := getRequestID(ctx)

		log := log.With(
			slog.String("method", info.FullMethod),
			slog.Time("received_at", start),
			slog.String("request_id", req_id),
			slog.Any("request", req),
		)

		resp, err := handler(ctx, req)

		log.Info("gRPC request completed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("response", resp),
			slog.Any("error", err),
		)

		return resp, err
	}
}

func getRequestID(ctx context.Context) (context.Context, string) {
	md, ok := metadata.FromIncomingContext(ctx)
	var traceID string
	if ok {
		if vals := md.Get("x-trace-id"); len(vals) > 0 && vals[0] != "" {
			traceID = vals[0]
		}
	}
	if traceID == "" {
		traceID = strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	ctx = context.WithValue(ctx, requestIDKey, traceID)
	return ctx, traceID
}
