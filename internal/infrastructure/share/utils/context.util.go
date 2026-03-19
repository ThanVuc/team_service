package utils

import (
	"context"

	"github.com/wagslane/go-rabbitmq"
	"google.golang.org/grpc/metadata"
)

func GetRequestIDFromOutgoingContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		vals := md.Get("x-request-id")
		if len(vals) > 0 {
			return vals[0]
		}
	}
	return ""
}

func GetUserIDFromOutgoingContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		vals := md.Get("x-user-id")
		if len(vals) > 0 {
			return vals[0]
		}
	}
	return ""
}

func GetRequestID(d rabbitmq.Delivery) string {
	if v, ok := d.Headers["request_id"]; ok {
		switch val := v.(type) {
		case string:
			return val
		case []byte:
			return string(val)
		}
	}
	return ""
}

func GetGroupIDFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		vals := md.Get("x-group-id")
		if len(vals) > 0 {
			return vals[0]
		}
	}
	return ""
}
